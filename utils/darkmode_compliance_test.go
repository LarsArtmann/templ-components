package utils

import (
	"path/filepath"
	"testing"

	"github.com/larsartmann/templ-components/utils/tone"
)

// TestDarkModeSemanticRegexCoversEveryHue is the regression proof for the
// derived hue list (TODO #427): the hand-maintained alternation had drifted
// and classes like text-emerald-600 or bg-yellow-500 WITHOUT a dark: half
// silently passed the guard. Every non-gray hue in the tone table, in both
// text- and bg- form, must now be FLAGGED when the dark: variant is
// missing — and ignored when it is present on the line.
func TestDarkModeSemanticRegexCoversEveryHue(t *testing.T) {
	t.Parallel()

	for _, h := range tone.Hues() {
		if h == tone.HueGray {
			continue // the neutral rule owns gray
		}

		for _, class := range []string{"text-" + string(h) + "-600", "bg-" + string(h) + "-500"} {
			gapLine := `class="` + class + `"`
			if !checkLineForDarkModeGap("synthetic.templ", gapLine, darkModeSemanticRe) {
				t.Errorf(
					"semantic regex does not flag %q without a dark: variant — the hue list drifted from tone.Hues() again",
					class,
				)
			}

			coveredLine := `class="` + class + ` dark:` + class + `"`
			if checkLineForDarkModeGap("synthetic.templ", coveredLine, darkModeSemanticRe) {
				t.Errorf("semantic regex flags %q although the dark: variant is present", coveredLine)
			}
		}
	}
}

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
