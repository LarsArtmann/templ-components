package icons

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAnimatedIconCSSExists verifies that every animation type has a corresponding
// CSS class definition in templates/custom.css. Without this test, deleting
// .tc-anim-* rules from custom.css would be undetectable — the TestCustomCSSUtilities
// drift guard in the utils package does not scan the icons module (separate Go module).
func TestAnimatedIconCSSExists(t *testing.T) {
	t.Parallel()

	cssPath := filepath.Join("..", "templates", "custom.css")

	data, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("cannot read templates/custom.css from icons module: %v", err)
	}

	css := string(data)

	// Base wrapper class must exist.
	if !strings.Contains(css, ".tc-anim") {
		t.Error("templates/custom.css missing .tc-anim base class")
	}

	// Every animation type must have a corresponding .tc-anim-<type> rule.
	for _, anim := range AllAnimations() {
		classSelector := ".tc-anim-" + string(anim)
		if !strings.Contains(css, classSelector) {
			t.Errorf(
				"templates/custom.css missing %s rule for animation %q",
				classSelector, anim,
			)
		}
	}
}
