// Package golden provides golden file comparison for snapshot tests.
//
// Usage:
//
//	golden.Assert(t, "button_primary", output)
//
// To update golden files:
//
//	go test -update ./...
package golden

import (
	"flag"
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

// Assert compares got against the golden file for the given test name.
// If the -update flag is set, it writes got to the golden file instead.
// Golden files are stored in testdata/<name>.golden relative to the test file.
// CSS class attribute order is normalized before comparison.
func Assert(t *testing.T, name, got string) {
	t.Helper()
	assertInDir(t, "testdata", name, got)
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

	normalized := normalizeClasses(got)

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

	want, err := os.ReadFile(goldenPath) //nolint:gosec // testdata path is test-controlled
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

// diff returns a simple line-by-line diff between two strings.
func diff(want, got string) string {
	wantLines := strings.SplitAfter(want, "\n")
	gotLines := strings.SplitAfter(got, "\n")

	maxLen := max(len(gotLines), len(wantLines))

	var b strings.Builder

	for i := range maxLen {
		w := lineAt(wantLines, i)

		g := lineAt(gotLines, i)
		if w != g {
			b.WriteString("--- ")
			b.WriteString(w)
			b.WriteString("\n+++ ")
			b.WriteString(g)
			b.WriteString("\n")
		}
	}

	return b.String()
}

func lineAt(lines []string, i int) string {
	if i < len(lines) {
		return lines[i]
	}

	return ""
}
