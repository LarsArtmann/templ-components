package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestInlineIconPathCompliance enforces the single-source-of-truth rule for
// SVG path data (AGENTS.md "SVG paths"): component sources must never inline
// a path d="M..." literal. Stroke and fill glyphs come from the icons package
// (icons.Icon for markup, icons.IconPathData for JS injection); paths needed
// across module boundaries (utils is a leaf and cannot import icons) live as
// shared constants in utils/svg. The sweep is the drift-guard style used by
// TestMotionReduceCompliance: it scans library sources so a future component
// cannot silently copy-paste path data again (the X glyph was inlined three
// times before this guard existed).
//
// Canonical definitions are exempt: icons/icon_paths.go (the icon table
// itself) and utils/svg (the shared constants + SpinnerSVG primitive).
func TestInlineIconPathCompliance(t *testing.T) {
	t.Parallel()

	root := ".."
	dirs := []string{
		"display", "forms", "feedback", "layout", "navigation",
		"errorpage", "recipes", "htmx", "datastar", "charts/echarts",
		"icons", "utils",
	}

	// File-relative exemptions with the reason each path literal is canonical.
	exemptions := map[string]string{
		"icons/icon_paths.go":                         "canonical icon path table",
		filepath.Join("utils", "svg", "svg.templ"):    "canonical shared constants + SpinnerSVG",
		filepath.Join("utils", "svg", "svg_templ.go"): "generated canonical constants",
	}

	// d="M..." in template/JS text and the d=\"M...\" form inside Go string
	// literals that embed JavaScript.
	pathRe := regexp.MustCompile(`d=\\"M[0-9]|d="M[0-9]`)

	violations := 0

	for _, dir := range dirs {
		dirPath := filepath.Join(root, dir)

		files, err := filepath.Glob(filepath.Join(dirPath, "*.templ"))
		if err != nil {
			t.Fatalf("glob %s: %v", dirPath, err)
		}
		goFiles, err := filepath.Glob(filepath.Join(dirPath, "*.go"))
		if err != nil {
			t.Fatalf("glob %s: %v", dirPath, err)
		}
		files = append(files, goFiles...)

		for _, file := range files {
			base := filepath.Base(file)
			if strings.HasSuffix(base, "_templ.go") || strings.HasSuffix(base, "_test.go") {
				continue
			}
			if _, ok := exemptions[filepath.Join(dir, base)]; ok {
				continue
			}

			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("read %s: %v", file, err)
			}

			for lineNum, line := range strings.Split(string(data), "\n") {
				if pathRe.MatchString(line) {
					t.Errorf(
						"%s:%d inlines SVG path data — single-source violation. "+
							"Render via icons.Icon / icons.IconWithStrokeWidth (markup) or "+
							"icons.IconPathData (JS injection); for cross-module needs add a "+
							"constant in utils/svg.",
						filepath.Join(dir, base),
						lineNum+1,
					)

					violations++
				}
			}
		}
	}

	if violations == 0 {
		t.Log("inline icon path compliance: OK (no inlined path data in library sources)")
	}
}

// TestInlineIconPathRegexDetectsViolations pins the detector itself: both
// literal forms (template attribute and Go-embedded JS) must match, while
// component-sourced paths stay clean.
func TestInlineIconPathRegexDetectsViolations(t *testing.T) {
	t.Parallel()

	pathRe := regexp.MustCompile(`d=\\"M[0-9]|d="M[0-9]`)

	mustMatch := []string{
		`<path d="M6 18L18 6M6 6l12 12"></path>`,
		"btn.innerHTML='<svg><path d=\\\"M6 18L18 6M6 6l12 12\\\"></path></svg>'",
	}
	for _, sample := range mustMatch {
		if !pathRe.MatchString(sample) {
			t.Errorf("detector missed violation sample: %q", sample)
		}
	}

	clean := []string{
		`<path d={ svg.PathXMark }></path>`,
		`<path fill-rule="evenodd" d={ p } clip-rule="evenodd"></path>`,
		`<path d={ fmt.Sprintf("M%s", coords) }></path>`,
	}
	for _, sample := range clean {
		if pathRe.MatchString(sample) {
			t.Errorf("detector flagged clean sample: %q", sample)
		}
	}
}
