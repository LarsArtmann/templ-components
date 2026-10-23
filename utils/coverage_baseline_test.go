package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// TestCoverageBaselineMatchesProse is the #438 guard: the coverage figures the
// docs quote must match a stored measurement, and that measurement must be
// internally consistent.
//
// The numbers live in scripts/coverage-baseline.txt, produced by
// `scripts/coverage-library.sh --record`. Two failure modes are covered:
//   - hand-typed prose drifting away from the recorded measurement
//   - a recorded baseline whose percentages no longer follow from its counts
//
// The baseline also records go_version: Go 1.27.0 and 1.27.3 assign different
// statement counts to identical coverage blocks, so a figure without its
// toolchain is not a fact (#440).
func TestCoverageBaselineMatchesProse(t *testing.T) {
	t.Parallel()

	baseline := readCoverageBaseline(t)
	required := []string{
		"recorded_at", "recorded_commit", "go_version",
		"raw_covered", "raw_total", "raw_pct",
		"hand_covered", "hand_total", "hand_pct",
		"ceiling_pct", "generated_share_pct",
	}

	for _, key := range required {
		if _, ok := baseline[key]; !ok {
			t.Errorf("scripts/coverage-baseline.txt is missing %q (regenerate with scripts/coverage-library.sh --record)", key)
		}
	}

	if t.Failed() {
		return
	}

	goVersion := baseline["go_version"]
	if !regexp.MustCompile(`^go\d+\.\d+\.\d+$`).MatchString(goVersion) {
		t.Errorf("baseline go_version = %q, want a goX.Y.Z toolchain id", goVersion)
	}

	raw := coverageNumbers{
		covered: intField(t, baseline, "raw_covered"),
		total:   intField(t, baseline, "raw_total"),
		pct:     floatField(t, baseline, "raw_pct"),
	}
	hand := coverageNumbers{
		covered: intField(t, baseline, "hand_covered"),
		total:   intField(t, baseline, "hand_total"),
		pct:     floatField(t, baseline, "hand_pct"),
	}
	ceiling := floatField(t, baseline, "ceiling_pct")
	generatedShare := floatField(t, baseline, "generated_share_pct")

	raw.assertPercent(t, "raw_pct")
	hand.assertPercent(t, "hand_pct")

	if hand.covered > raw.covered || hand.total > raw.total {
		t.Errorf("hand-written counts (%d/%d) exceed raw counts (%d/%d)", hand.covered, hand.total, raw.covered, raw.total)
	}

	// Every hand-written statement covered, generated coverage unchanged:
	// (raw_covered - hand_covered + hand_total) / raw_total.
	wantCeiling := percent(raw.covered-hand.covered+hand.total, raw.total)
	assertFloat(t, "ceiling_pct", ceiling, wantCeiling)
	assertFloat(t, "generated_share_pct", generatedShare, percent(raw.total-hand.total, raw.total))

	// The living prose must quote the recorded figures AND name the baseline,
	// otherwise a reader cannot tell which measurement the numbers came from.
	agents := readDoc(t, "AGENTS.md")
	for _, want := range []string{
		"scripts/coverage-baseline.txt",
		formatFigure(hand.pct, "hand-written"),
		formatFigure(raw.pct, "raw"),
		formatFigure(ceiling, "ceiling"),
		formatFigure(generatedShare, "of statements"),
	} {
		if !strings.Contains(string(agents), want) {
			t.Errorf("AGENTS.md coverage bullet does not contain %q; update it with the recorded baseline", want)
		}
	}
}

type coverageNumbers struct {
	covered int
	total   int
	pct     float64
}

func (n coverageNumbers) assertPercent(t *testing.T, label string) {
	t.Helper()

	if n.total <= 0 {
		t.Errorf("%s: total statements = %d, want > 0", label, n.total)
		return
	}

	assertFloat(t, label, n.pct, percent(n.covered, n.total))
}

func percent(covered, total int) float64 {
	if total == 0 {
		return 0
	}

	return float64(covered) * 100 / float64(total)
}

func assertFloat(t *testing.T, label string, got, want float64) {
	t.Helper()

	// The baseline stores one decimal; half a unit of the last place is noise.
	if diff := got - want; diff > 0.051 || diff < -0.051 {
		t.Errorf("%s = %.1f, but the counts imply %.1f", label, got, want)
	}
}

func formatFigure(pct float64, qualifier string) string {
	return strconv.FormatFloat(pct, 'f', 1, 64) + "% " + qualifier
}

func readCoverageBaseline(t *testing.T) map[string]string {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("..", "scripts", "coverage-baseline.txt"))
	if err != nil {
		t.Fatalf("read scripts/coverage-baseline.txt: %v", err)
	}

	fields := make(map[string]string)

	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			t.Errorf("malformed baseline line %q", line)
			continue
		}

		fields[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}

	return fields
}

func intField(t *testing.T, fields map[string]string, key string) int {
	t.Helper()

	value, err := strconv.Atoi(fields[key])
	if err != nil {
		t.Fatalf("%s = %q, want an integer", key, fields[key])
	}

	return value
}

func floatField(t *testing.T, fields map[string]string, key string) float64 {
	t.Helper()

	value, err := strconv.ParseFloat(fields[key], 64)
	if err != nil {
		t.Fatalf("%s = %q, want a number", key, fields[key])
	}

	return value
}
