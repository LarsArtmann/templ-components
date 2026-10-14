package utils

import (
	"path/filepath"
	"testing"
)

// scanDarkModeDirs runs ScanDarkMode over the library package dirs (relative
// to this test's CWD) and reports findings for ONE rule id in the guard
// format.
func scanDarkModeDirs(t *testing.T, dirs []string, rule string) {
	t.Helper()

	violations := 0

	for _, dir := range dirs {
		findings, err := ScanDarkMode(filepath.Join("..", dir))
		if err != nil {
			t.Logf("scan error for %s: %v", dir, err)

			continue
		}

		for _, f := range findings {
			if f.Rule != rule {
				continue
			}

			violations++

			t.Errorf("dark mode gap in %s:\n  %s", f.File, f.Snippet)
		}
	}

	if violations > 0 {
		t.Errorf("found %d dark mode compliance violations", violations)
	}
}

// TestDarkModeCompliance verifies that every neutral color class (text-gray-*,
// bg-white, bg-gray-*, border-gray-*, divide-gray-*, ring-gray-*) in .templ
// and .go source files has a corresponding dark: variant on the same line.
func TestDarkModeCompliance(t *testing.T) {
	t.Parallel()

	dirs := []string{"display", "feedback", "forms", "navigation", "errorpage", "layout", "htmx", "examples/demo"}
	scanDarkModeDirs(t, dirs, "darkmode-neutral")
}

// TestDarkModeSemanticColors verifies that semantic color classes (bg-blue-600,
// bg-red-600, text-blue-600, text-green-500, etc.) in .templ and .go source
// files have corresponding dark: variants on the same line.
func TestDarkModeSemanticColors(t *testing.T) {
	t.Parallel()

	dirs := []string{"display", "feedback", "forms", "navigation", "errorpage", "layout", "htmx", "examples/demo"}
	scanDarkModeDirs(t, dirs, "darkmode-semantic")
}
