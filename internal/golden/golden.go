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

		if err := os.WriteFile(goldenPath, []byte(normalized), 0o600); err != nil {
			t.Fatalf("golden: write %s: %v", goldenPath, err)
		}

		t.Logf("golden: updated %s", goldenPath)

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
	var b strings.Builder

	wantIdx, gotIdx := 0, 0

	for wantIdx < wantLen && gotIdx < gotLen {
		if wantLines[wantIdx] == gotLines[gotIdx] {
			wantIdx++
			gotIdx++
		} else if lcs[wantIdx+1][gotIdx] >= lcs[wantIdx][gotIdx+1] {
			fmt.Fprintf(&b, "--- [%d] %s\n", wantIdx+1, wantLines[wantIdx])
			wantIdx++
		} else {
			fmt.Fprintf(&b, "+++ [%d] %s\n", gotIdx+1, gotLines[gotIdx])
			gotIdx++
		}
	}

	for ; wantIdx < wantLen; wantIdx++ {
		fmt.Fprintf(&b, "--- [%d] %s\n", wantIdx+1, wantLines[wantIdx])
	}

	for ; gotIdx < gotLen; gotIdx++ {
		fmt.Fprintf(&b, "+++ [%d] %s\n", gotIdx+1, gotLines[gotIdx])
	}

	return b.String()
}
