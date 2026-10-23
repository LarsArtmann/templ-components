// Package golden provides golden file comparison for snapshot tests.
//
// Usage:
//
//	golden.Assert(t, "button_primary", output)
//
// For table-driven sweeps:
//
//	golden.AssertSnapshots(t, []golden.Snapshot{
//		{"badge_success", utils.Render(t, Badge(...))},
//		{"badge_error", utils.Render(t, Badge(...))},
//	})
//
// To update golden files:
//
//	go test -update ./...
package golden

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"testing"
)

//nolint:gochecknoglobals // CLI flag for golden file updates
var update = flag.Bool("update", false, "update golden files instead of comparing")

var classRe = regexp.MustCompile(`class="([^"]*)"`)

// autoIDRe matches EnsureID-generated identifiers (utils.EnsureID):
// tc-<prefix>-<digits>-<digits> (fallback: timestamp + counter — tried first because digits are valid hex)
// tc-<prefix>-<16 hex chars> (crypto/rand primary path)
// Prefixes may contain hyphens (e.g. "mobile-menu" in Nav).
// These are non-deterministic and must be normalized for reproducible golden files.
var autoIDRe = regexp.MustCompile(`tc-([a-z][a-z-]*)-(?:\d+-\d+|[a-f0-9]{16})`)

// Assert compares got against the golden file for the given test name.
// If the -update flag is set, it writes got to the golden file instead.
// Golden files are stored in testdata/<name>.golden relative to the test file.
// CSS class attribute order and auto-generated IDs are normalized before comparison.
func Assert(t *testing.T, name, got string) {
	t.Helper()
	assertInDir(t, "testdata", name, got)
}

// Snapshot pairs a test name with its rendered HTML for table-driven golden tests.
type Snapshot struct {
	Name string
	HTML string
}

// AssertSnapshots runs golden assertions for multiple snapshots as named subtests.
// Each snapshot becomes t.Run(name, ...) with t.Parallel(), eliminating the boilerplate
// of individual t.Run + golden.Assert calls. Example:
//
//	golden.AssertSnapshots(t, []golden.Snapshot{
//		{"badge_success", utils.Render(t, Badge(BadgeProps{Type: BadgeSuccess}))},
//		{"badge_error", utils.Render(t, Badge(BadgeProps{Type: BadgeError}))},
//	})
func AssertSnapshots(t *testing.T, snapshots []Snapshot) {
	t.Helper()

	for _, s := range snapshots {
		t.Run(s.Name, func(t *testing.T) {
			t.Parallel()
			Assert(t, s.Name, s.HTML)
		})
	}
}

// assertInDir is the directory-parameterized core of Assert. It exists so the
// golden package's own tests can pass a per-test t.TempDir() instead of sharing
// a single testdata/ directory (which raced under t.Parallel when tests created
// and removed files concurrently).
func assertInDir(t *testing.T, dir, name, got string) {
	t.Helper()

	// Guard against path traversal: name must be a base filename, not a path.
	if strings.ContainsAny(name, `/\`) || filepath.Clean(name) != name {
		t.Fatalf("golden: invalid name %q (must not contain path separators or ..)", name)

		return
	}

	goldenPath := filepath.Join(dir, name+".golden")

	normalized := normalize(got)

	if *update {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o750); err != nil {
			t.Fatalf("golden: create golden dir: %v", err)
		}

		previous, readErr := os.ReadFile(goldenPath)

		if err := os.WriteFile(goldenPath, []byte(normalized), 0o600); err != nil {
			t.Fatalf("golden: write %s: %v", goldenPath, err)
		}

		// Changed-file summary: after an -update sweep the developer wants to
		// know WHICH goldens actually changed (same-edit rule: every changed
		// golden carries its doc-count bump + CHANGELOG line in the same
		// batch). Created files log as new; identical rewrites stay silent.
		switch {
		case readErr != nil:
			t.Logf("golden: created %s", goldenPath)
		case string(previous) != normalized:
			t.Logf("golden: CHANGED %s", goldenPath)
		}

		return
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf(
			"golden: read %s: %v\nHint: run with -update to create golden files",
			goldenPath,
			err,
		)
	}

	if string(want) != normalized {
		t.Errorf(
			"golden: %s mismatch\n--- want\n+++ got\n%s",
			name,
			diff(string(want), normalized),
		)
	}
}

// normalize applies all normalizations to make golden files deterministic:
//   - Auto-generated IDs (tc-<prefix>-<hex>) are replaced with tc-<prefix>-NORMALIZED
//   - CSS class tokens inside class="..." are sorted alphabetically
//
// The order matters: ID normalization runs first so that any ID references
// inside class values (unlikely but possible) are already stable.
func normalize(html string) string {
	return normalizeClasses(normalizeIDs(html))
}

// normalizeIDs replaces non-deterministic EnsureID values with stable placeholders.
// Format tc-modal-a1b2c3d4e5f6a7b8 becomes tc-modal-NORMALIZED so that golden
// files are reproducible regardless of crypto/rand output. All cross-references
// (aria-controls, popovertarget, for, href="#...", inline JS) are also replaced
// because the substitution operates on the full document text.
func normalizeIDs(html string) string {
	return autoIDRe.ReplaceAllString(html, "tc-${1}-NORMALIZED")
}

// normalizeClasses sorts CSS class values in class="" attributes for deterministic comparison.
func normalizeClasses(html string) string {
	return classRe.ReplaceAllStringFunc(html, func(match string) string {
		sub := classRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}

		classes := strings.Fields(sub[1])
		sort.Strings(classes)

		return `class="` + strings.Join(classes, " ") + `"`
	})
}

// diff returns a line-by-line diff between two strings with line numbers.
// Uses LCS (Longest Common Subsequence) alignment so that insertions and
// deletions don't cascade into showing every subsequent line as changed.
// Long replacement lines (templ output is often one huge line per element)
// are additionally windowed around their first difference and followed by a
// word-level diff of just the changed tokens — a bare line pair is
// unreadable during goldens triage otherwise (TODO #339).
//
//nolint:makezero // pre-allocated DP table is intentional; index assignment follows
func diff(want, got string) string {
	wantLines := strings.Split(want, "\n")
	gotLines := strings.Split(got, "\n")
	wantLen, gotLen := len(wantLines), len(gotLines)

	// Build LCS DP table (O(wantLen*gotLen) space — fine for small golden files).
	lcs := make([][]int, wantLen+1)
	for idx := range lcs {
		lcs[idx] = make([]int, gotLen+1)
	}

	for wi := wantLen - 1; wi >= 0; wi-- {
		for gi := gotLen - 1; gi >= 0; gi-- {
			if wantLines[wi] == gotLines[gi] {
				lcs[wi][gi] = lcs[wi+1][gi+1] + 1
			} else {
				lcs[wi][gi] = max(lcs[wi+1][gi], lcs[wi][gi+1])
			}
		}
	}

	// Backtrack through the LCS to emit only changed lines with line numbers.
	return emitDiff(wantLines, gotLines, lcs)
}

const (
	// longLineThreshold is where a diff line stops being printed verbatim.
	longLineThreshold = 160
	// windowRadius is the context kept on each side of a long line's first
	// difference when windowing.
	windowRadius = 80
)

// diffOp is one line-level change: '-' removes wantLines[wantIdx], '+'
// adds gotLines[gotIdx].
type diffOp struct {
	kind    byte
	wantIdx int
	gotIdx  int
}

// emitDiff backtracks the LCS table into change ops and renders them. A '-'
// op immediately followed by a '+' op is a replacement pair: short lines
// print verbatim (the historical format, pinned by tests), long lines get
// windowing around their first difference plus a word diff of the changed
// tokens.
func emitDiff(wantLines, gotLines []string, lcs [][]int) string {
	ops := collectDiffOps(wantLines, gotLines, lcs)

	var b strings.Builder

	for idx := 0; idx < len(ops); idx++ {
		change := ops[idx]

		if change.kind == '-' && idx+1 < len(ops) && ops[idx+1].kind == '+' {
			next := ops[idx+1]
			writeReplacement(&b, wantLines[change.wantIdx], gotLines[next.gotIdx], change.wantIdx, next.gotIdx)

			idx++

			continue
		}

		if change.kind == '-' {
			fmt.Fprintf(&b, "--- [%d] %s\n", change.wantIdx+1, windowLine(wantLines[change.wantIdx]))
		} else {
			fmt.Fprintf(&b, "+++ [%d] %s\n", change.gotIdx+1, windowLine(gotLines[change.gotIdx]))
		}
	}

	return b.String()
}

func writeReplacement(b *strings.Builder, oldLine, newLine string, wantIdx, gotIdx int) {
	if len(oldLine) <= longLineThreshold && len(newLine) <= longLineThreshold {
		fmt.Fprintf(b, "--- [%d] %s\n", wantIdx+1, oldLine)
		fmt.Fprintf(b, "+++ [%d] %s\n", gotIdx+1, newLine)

		return
	}

	oldWindow, newWindow := windowPair(oldLine, newLine)
	fmt.Fprintf(b, "--- [%d] %s\n", wantIdx+1, oldWindow)
	fmt.Fprintf(b, "+++ [%d] %s\n", gotIdx+1, newWindow)
	fmt.Fprintf(b, "    ~~~ word diff: %s\n", strings.TrimSpace(wordDiff(oldLine, newLine)))
}

// lcsAt reads the LCS table defensively: callers walk one row/column past the
// line slices, so a malformed or short table must not panic the diff renderer.
func lcsAt(lcs [][]int, row, col int) int {
	if row < 0 || row >= len(lcs) || col < 0 || col >= len(lcs[row]) {
		return 0
	}

	return lcs[row][col]
}

func collectDiffOps(wantLines, gotLines []string, lcs [][]int) []diffOp {
	var ops []diffOp

	wantIdx, gotIdx := 0, 0

	for wantIdx < len(wantLines) && gotIdx < len(gotLines) {
		switch {
		case wantLines[wantIdx] == gotLines[gotIdx]:
			wantIdx++
			gotIdx++
		case lcsAt(lcs, wantIdx+1, gotIdx) >= lcsAt(lcs, wantIdx, gotIdx+1):
			ops = append(ops, diffOp{kind: '-', wantIdx: wantIdx, gotIdx: -1})
			wantIdx++
		default:
			ops = append(ops, diffOp{kind: '+', wantIdx: -1, gotIdx: gotIdx})
			gotIdx++
		}
	}

	for ; wantIdx < len(wantLines); wantIdx++ {
		ops = append(ops, diffOp{kind: '-', wantIdx: wantIdx, gotIdx: -1})
	}

	for ; gotIdx < len(gotLines); gotIdx++ {
		ops = append(ops, diffOp{kind: '+', wantIdx: -1, gotIdx: gotIdx})
	}

	return ops
}

// windowLine windows a standalone (unpaired) long line from its start — the
// element's opening tag carries the most signal on pure inserts/deletes.
func windowLine(line string) string {
	return windowFrom(line, 0)
}

// windowPair windows both halves of a long replacement pair around their
// first difference so the interesting region stays aligned and visible.
func windowPair(oldLine, newLine string) (string, string) {
	firstDiff := 0
	for firstDiff < len(oldLine) && firstDiff < len(newLine) && oldLine[firstDiff] == newLine[firstDiff] {
		firstDiff++
	}

	start := max(firstDiff-windowRadius, 0)

	return windowFrom(oldLine, start), windowFrom(newLine, start)
}

func windowFrom(line string, start int) string {
	if len(line) <= longLineThreshold {
		return line
	}

	start = min(start, max(len(line)-2*windowRadius, 0))
	end := min(start+2*windowRadius, len(line))

	prefix, suffix := "", ""

	if start > 0 {
		prefix = "…"
	}

	if end < len(line) {
		suffix = fmt.Sprintf("…(+%d chars)", len(line)-end)
	}

	return prefix + line[start:end] + suffix
}

// wordDiff tokenizes two long lines on whitespace and prints only the
// changed tokens (LCS over words) — the readable form of a one-token change
// inside a multi-kilobyte single-line element.
func wordDiff(want, got string) string {
	wantWords := strings.Fields(want)
	gotWords := strings.Fields(got)
	lcs := wordLCS(wantWords, gotWords)

	var b strings.Builder

	wordBacktrack(&b, wantWords, gotWords, lcs)

	return b.String()
}

func wordBacktrack(b *strings.Builder, wantWords, gotWords []string, lcs [][]int) {
	wantAt, gotAt := 0, 0

	for wantAt < len(wantWords) && gotAt < len(gotWords) {
		switch {
		case wantWords[wantAt] == gotWords[gotAt]:
			wantAt++
			gotAt++
		case lcs[wantAt+1][gotAt] >= lcs[wantAt][gotAt+1]:
			fmt.Fprintf(b, "-%s ", wantWords[wantAt])
			wantAt++
		default:
			fmt.Fprintf(b, "+%s ", gotWords[gotAt])
			gotAt++
		}
	}

	for ; wantAt < len(wantWords); wantAt++ {
		fmt.Fprintf(b, "-%s ", wantWords[wantAt])
	}

	for ; gotAt < len(gotWords); gotAt++ {
		fmt.Fprintf(b, "+%s ", gotWords[gotAt])
	}
}

//nolint:makezero // pre-allocated DP table is intentional; index assignment follows
func wordLCS(wantWords, gotWords []string) [][]int {
	lcs := make([][]int, len(wantWords)+1)
	for idx := range lcs {
		lcs[idx] = make([]int, len(gotWords)+1)
	}

	for wi, wantWord := range slices.Backward(wantWords) {
		for gi := len(gotWords) - 1; gi >= 0; gi-- {
			if wantWord == gotWords[gi] {
				lcs[wi][gi] = lcs[wi+1][gi+1] + 1
			} else {
				lcs[wi][gi] = max(lcs[wi+1][gi], lcs[wi][gi+1])
			}
		}
	}

	return lcs
}
