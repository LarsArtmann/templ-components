package utils

import (
	"os"
	"strings"
	"testing"
)

// TestTailwindGoSourceScanning verifies that CSS entry points include @source
// directives for .go files, not just .templ files. Tailwind v4's content
// scanner only scans files matched by @source. Many components define Tailwind
// class lookup maps in Go source (buttonVariantLookup, feedbackStyleMap,
// familyStyleMap, modalSizeLookup, etc.). Without scanning .go files, those
// classes produce silently-missing CSS.
//
// History (2026-07-28): errorpage family classes (amber/orange/purple)
// were entirely missing from the compiled CSS because familyStyleMap lived in
// errorpage/styles.go and the CSS only scanned *.templ.
func TestTailwindGoSourceScanning(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		path string
	}{
		{"consumer template", "../templates/app.css"},
		{"demo CSS", "../examples/demo/demo.css"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			data, err := os.ReadFile(tc.path)
			if err != nil {
				t.Fatalf("read %s: %v", tc.path, err)
			}

			src := string(data)

			// Must scan .templ files.
			if !strings.Contains(src, "*.templ") {
				t.Errorf("%s missing @source directive for *.templ files", tc.path)
			}

			// Must scan .go files (lookup maps, constants with Tailwind classes).
			if !strings.Contains(src, "*.go") {
				t.Errorf(
					"%s missing @source directive for *.go files — "+
						"Tailwind classes defined in Go source "+
						"(buttonVariantLookup, feedbackStyleMap, familyStyleMap) "+
						"will not get CSS generated",
					tc.path,
				)
			}
		})
	}
}
