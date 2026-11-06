package tone_test

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils/tone"
)

var docUpdate = flag.Bool("update", false, "rewrite the generated tone sections in the docs")

const (
	themeBridgeDoc = "../../docs/recipes/scoped-theme-bridge.md"
	themeLayerCSS  = "../../templates/templ-components-theme.css"

	paletteBegin = "<!-- BEGIN generated: tone palette table (utils/tone/doc_sync_test.go; go test ./utils/tone/ -update) -->"
	paletteEnd   = "<!-- END generated: tone palette table -->"

	familyBegin = "<!-- BEGIN generated: token family table (utils/tone/doc_sync_test.go; go test ./utils/tone/ -update) -->"
	familyEnd   = "<!-- END generated: token family table -->"
)

// tonePaletteTable renders the "Which colors does the library use?" table
// from the tone package's own tables — the single source of truth for the
// library palette (report H4: this was the fifth hand-maintained copy).
func tonePaletteTable() string {
	chartSeries := map[tone.Hue]int{}
	for i, hue := range tone.ChartHues() {
		chartSeries[hue] = i + 1
	}

	var b strings.Builder
	b.WriteString("| Hue | Tones rendered in it | Chart palette |\n")
	b.WriteString("| --- | --- | --- |\n")

	for _, hue := range tone.Hues() {
		prefix := string(hue) + "-"

		var tones []string

		for _, t := range tone.Tones() {
			if toneUsesHue(t, prefix) {
				tones = append(tones, string(t))
			}
		}

		tonesCell := "—"
		if len(tones) > 0 {
			tonesCell = strings.Join(tones, ", ")
		}

		chartCell := "—"
		if series, ok := chartSeries[hue]; ok {
			chartCell = fmt.Sprintf("yes (series %d)", series)
		}

		fmt.Fprintf(&b, "| `%s` | %s | %s |\n", hue, tonesCell, chartCell)
	}

	return b.String()
}

// toneUsesHue reports whether any surface of the tone references the hue
// (classes carry the hue as a `hue-shade` token prefix).
func toneUsesHue(t tone.Tone, prefix string) bool {
	return toneUsesPrefix(t, prefix)
}

// toneUsesPrefix reports whether any surface of the tone carries the given
// class prefix (a hue like "blue-" or a token family like "accent-").
func toneUsesPrefix(t tone.Tone, prefix string) bool {
	surfaces := tone.For(t)

	for _, pair := range []tone.Pair{
		surfaces.Dot, surfaces.Text, surfaces.TextSoft, surfaces.Border,
		surfaces.BGTint, surfaces.BGSolid, surfaces.Icon,
	} {
		if strings.Contains(pair.Light, prefix) || strings.Contains(pair.Dark, prefix) {
			return true
		}
	}

	return false
}

// toneFamilyTable renders the token-family table (ADR-0049): tones that
// render through token families instead of palette hues. The accent
// family is derived from the tone tables, like the palette table above —
// adding a family means extending this generator, not hand-editing docs.
func toneFamilyTable() string {
	var b strings.Builder
	b.WriteString("| Family | Tones rendered in it |\n")
	b.WriteString("| --- | --- |\n")

	var accentTones []string

	for _, t := range tone.Tones() {
		if toneUsesPrefix(t, "accent-") {
			accentTones = append(accentTones, string(t))
		}
	}

	tonesCell := "—"
	if len(accentTones) > 0 {
		tonesCell = strings.Join(accentTones, ", ")
	}

	fmt.Fprintf(&b, "| `accent` | %s |\n", tonesCell)

	return b.String()
}

// syncGeneratedDocSection verifies (or with -update rewrites) one generated
// section of the theme bridge doc against its derived table. Shared by the
// palette and family table tests; the -update branches read-modify-write
// the same file, so callers must stay sequential (never t.Parallel).
func syncGeneratedDocSection(t *testing.T, begin, end, table string) {
	t.Helper()

	doc, err := os.ReadFile(themeBridgeDoc)
	if err != nil {
		t.Fatalf("read theme bridge doc: %v", err)
	}

	want := begin + "\n" + table + end

	if *docUpdate {
		updated := replaceGeneratedSection(string(doc), begin, end, want)
		if updated != string(doc) {
			if err := os.WriteFile(themeBridgeDoc, []byte(updated), 0o644); err != nil {
				t.Fatalf("rewrite theme bridge doc: %v", err)
			}
		}

		return
	}

	start := strings.Index(string(doc), begin)

	endIdx := strings.Index(string(doc), end)
	if start < 0 || endIdx < start {
		t.Fatalf(
			"theme bridge doc lacks the generated markers %q/%q — run go test ./utils/tone/ -update",
			begin, end,
		)
	}

	have := string(doc)[start : endIdx+len(end)]
	if normalizeGeneratedBlock(have) != normalizeGeneratedBlock(want) {
		t.Errorf(
			"theme bridge doc generated section is stale:\n--- have ---\n%s\n--- want ---\n%s\nRun: go test ./utils/tone/ -update",
			have,
			want,
		)
	}
}

// TestThemeBridgePaletteTableSync keeps the generated palette table in
// docs/recipes/scoped-theme-bridge.md byte-identical to what the tone
// package derives. Any tone/hue table change fails here until the doc is
// regenerated (`go test ./utils/tone/ -run TestThemeBridgePaletteTableSync -update`).
// NOT parallel: the -update branches of both doc-sync tests read-modify-
// write the same file, and interleaved writes clobber each other's fill
// (both tables landed half-filled once — 2026-11-06). Trivial cost.
//
//nolint:paralleltest // sequential on purpose: shared read-modify-write doc target
func TestThemeBridgePaletteTableSync(t *testing.T) {
	syncGeneratedDocSection(t, paletteBegin, paletteEnd, tonePaletteTable())
}

// TestThemeBridgeFamilyTableSync keeps the generated token-family table
// (ADR-0049) in docs/recipes/scoped-theme-bridge.md byte-identical to what
// the tone package derives — same mechanism as the palette table above.
// NOT parallel: shared read-modify-write target (see the palette test).
//
//nolint:paralleltest // sequential on purpose: shared read-modify-write doc target
func TestThemeBridgeFamilyTableSync(t *testing.T) {
	syncGeneratedDocSection(t, familyBegin, familyEnd, toneFamilyTable())
}

// normalizeGeneratedBlock makes the generated-section comparison insensitive to
// markdown table padding. The repo formatter aligns tables; the generator emits
// compact rows. Both express the same table, and comparing padding made the
// guard fail after every formatter run (2026-10-23 recurrence). Semantic drift
// (a missing hue, a changed column value) still fails.
func normalizeGeneratedBlock(block string) string {
	lines := strings.Split(block, "\n")
	normalized := make([]string, 0, len(lines))

	for _, line := range lines {
		if isTableSeparator(line) {
			continue
		}

		line = strings.Join(strings.Fields(line), " ")
		if line == "" {
			continue
		}

		normalized = append(normalized, line)
	}

	return strings.Join(normalized, "\n")
}

// isTableSeparator reports whether the line is a markdown table delimiter row
// (only pipes, dashes, colons, and spaces). Padding differences make the
// delimiter row length vary; the row carries no data, so both sides skip it.
func isTableSeparator(line string) bool {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "|") {
		return false
	}

	hasDash := false

	for _, char := range trimmed {
		switch char {
		case '|', ' ':
		case '-':
			hasDash = true
		case ':':
		default:
			return false
		}
	}

	return hasDash
}

// TestThemeTokenLayerInventory pins the tc-* token inventory of
// templates/templ-components-theme.css (ADR-0008): exactly the four semantic
// ramps (primary/danger/success/warning — the accent, danger, success, and
// warning tones) at every shade 50–950, so a ramp that loses a shade or a
// stray token that appears fails loudly instead of silently re-breaking
// consumer overrides.
func TestThemeTokenLayerInventory(t *testing.T) {
	t.Parallel()

	css, err := os.ReadFile(themeLayerCSS)
	if err != nil {
		t.Fatalf("read theme layer css: %v", err)
	}

	tokenRe := regexp.MustCompile(`--color-tc-([a-z]+)-(\d+):`)
	shades := []string{"50", "100", "200", "300", "400", "500", "600", "700", "800", "900", "950"}
	want := map[string]bool{}

	for _, ramp := range []string{"primary", "danger", "success", "warning"} {
		for _, shade := range shades {
			want[fmt.Sprintf("%s-%s", ramp, shade)] = true
		}
	}

	got := map[string]bool{}
	for _, match := range tokenRe.FindAllStringSubmatch(string(css), -1) {
		got[match[1]+"-"+match[2]] = true
	}

	for token := range want {
		if !got[token] {
			t.Errorf("theme layer css is missing token --color-tc-%s (ramp incomplete)", token)
		}
	}

	for token := range got {
		if !want[token] {
			t.Errorf("theme layer css carries unexpected ramp token --color-tc-%s (inventory drift)", token)
		}
	}
}

func replaceGeneratedSection(doc, begin, end, replacement string) string {
	start := strings.Index(doc, begin)

	finish := strings.Index(doc, end)
	if start < 0 || finish < start {
		return doc + "\n" + replacement + "\n"
	}

	return doc[:start] + replacement + doc[finish+len(end):]
}
