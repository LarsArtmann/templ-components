package utils

import (
	"path/filepath"
	"regexp"
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

// TestDarkModeRingShadowColorRules pins TODO #490: the COLOR-carrying
// ring-offset/shadow forms (ring-offset-white|black|gray-N|hue-N,
// shadow-white|black/N|gray-N|hue-N/N) must be flagged without a dark:
// variant and covered by the same dark: half when present. Size/width
// classes (shadow-lg, shadow-sm, ring-offset-1) carry no color and must
// NEVER trip the guard — that is what keeps the documented Toggle-thumb
// exception (bg-white shadow-sm) working without an allowlist entry.
func TestDarkModeRingShadowColorRules(t *testing.T) {
	t.Parallel()

	mustFlag := []string{
		// neutral color-carrying forms
		`class="ring-offset-white"`,
		`class="ring-offset-gray-100"`,
		`class="shadow-gray-500"`,
		`class="shadow-white/10"`,
		`class="shadow-black/20"`,
		// semantic color-carrying forms (every hue via the derived table)
		`class="ring-offset-blue-600"`,
		`class="shadow-rose-500/20"`,
	}

	// each color-carrying line must be flagged by its own rule; spell the
	// rule split explicitly so the table stays honest
	mustFlagPerRule := map[*regexp.Regexp][]string{}

	for _, line := range mustFlag {
		if darkModeNeutralRe.MatchString(line) {
			mustFlagPerRule[darkModeNeutralRe] = append(mustFlagPerRule[darkModeNeutralRe], line)
		}

		if darkModeSemanticRe.MatchString(line) {
			mustFlagPerRule[darkModeSemanticRe] = append(mustFlagPerRule[darkModeSemanticRe], line)
		}
	}

	for rule, lines := range mustFlagPerRule {
		for _, line := range lines {
			if !checkLineForDarkModeGap("synthetic.templ", line, rule) {
				t.Errorf("rule does not flag %q — the ring-offset/shadow color scan regressed", line)
			}
		}
	}

	neverFlag := []string{
		// size/width classes carry no color
		`class="shadow-lg"`,
		`class="shadow-sm"`,
		`class="shadow-xl hover:shadow-md"`,
		`class="ring-offset-1"`,
		`class="ring-offset-2 focus-visible:ring-offset-2"`,
		// color forms WITH their dark: half on the same line
		`class="shadow-lg dark:shadow-black/20"`,
		`class="ring-offset-2 dark:ring-offset-gray-900"`,
		`class="ring-offset-blue-600 dark:ring-offset-blue-400"`,
		`class="shadow-gray-500 dark:shadow-gray-700"`,
		// variants that only color the DARK half are fine (light half is size-only)
		`class="shadow-lg dark:shadow-black/20 hover:shadow-xl"`,
	}
	for _, line := range neverFlag {
		if checkLineForDarkModeGap("synthetic.templ", line, darkModeNeutralRe) ||
			checkLineForDarkModeGap("synthetic.templ", line, darkModeSemanticRe) {
			t.Errorf("rule flags a non-color or properly covered line: %q", line)
		}
	}

	// the documented Toggle thumb line — shadow-sm is a size class, and the
	// bg-white half is covered by the forms/toggle.templ exception (the path
	// must match the exception table, hence these two dedicated cases)
	if checkLineForDarkModeGap("synthetic.templ", `class="w-6 h-6 shadow-sm"`, darkModeNeutralRe) ||
		checkLineForDarkModeGap("synthetic.templ", `class="w-6 h-6 shadow-sm"`, darkModeSemanticRe) {
		t.Errorf("shadow-sm (size-only) flagged outside the Toggle exception")
	}

	if checkLineForDarkModeGap("forms/toggle.templ", `class="w-6 h-6 bg-white shadow-sm"`, darkModeNeutralRe) {
		t.Errorf("Toggle thumb exception no longer covers the real thumb line")
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
