package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// iconPathExemptions lists the canonical path-data definitions by
// library-relative path SUFFIX with the reason each is canonical.
//
//nolint:gochecknoglobals // declarative exemption table
var iconPathExemptions = []struct {
	suffix string
	why    string
}{
	{"icons/icon_paths.go", "canonical icon path table"},
	{filepath.Join("utils", "svg", "svg.templ"), "canonical shared constants + SpinnerSVG"},
	{filepath.Join("utils", "svg", "svg_templ.go"), "generated canonical constants"},
}

// iconPathInlineRe matches `d="M..."` in template text and the `d=\"M...\"`
// form inside Go string literals that embed JavaScript.
var iconPathInlineRe = regexp.MustCompile(`d=\\"M[0-9]|d="M[0-9]`)

// ScanIconPaths reports inlined SVG `d="M..."` path literals in the DIRECT
// .templ and .go children of dir (generated *_templ.go and *_test.go files
// excluded, canonical definitions exempt). The single-source rule: stroke and
// fill glyphs come from the icons package (icons.Icon for markup,
// icons.IconPathData for JS injection); paths needed across module boundaries
// live as shared constants in utils/svg.
//
// It is the engine behind TestInlineIconPathCompliance (repo guard) and
// `tc lint --icons`. Findings use dir-relative paths.
func ScanIconPaths(dir string) ([]Finding, error) {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", dir, err)
	}

	defer func() { _ = root.Close() }()

	var findings []Finding

	for _, pattern := range []string{"*.templ", "*.go"} {
		matches, globErr := fs.Glob(root.FS(), pattern)
		if globErr != nil {
			return nil, fmt.Errorf("glob %s: %w", pattern, globErr)
		}

		for _, rel := range matches {
			if strings.HasSuffix(rel, "_templ.go") || strings.HasSuffix(rel, "_test.go") {
				continue
			}

			full := filepath.Join(dir, rel)

			if iconPathExempt(full) {
				continue
			}

			data, readErr := root.ReadFile(rel)
			if readErr != nil {
				return nil, fmt.Errorf("read %s: %w", rel, readErr)
			}

			for i, line := range strings.Split(string(data), "\n") {
				if !iconPathInlineRe.MatchString(line) || lintFindingSuppressed(line, "icon-path-inline") {
					continue
				}

				findings = append(findings, Finding{
					File: rel,
					Line: i + 1,
					Rule: "icon-path-inline",
					Message: "inlines SVG path data — render via icons.Icon / icons.IconWithStrokeWidth (markup) " +
						"or icons.IconPathData (JS injection); cross-module needs go in utils/svg",
					Snippet: strings.TrimSpace(line),
				})
			}
		}
	}

	return findings, nil
}

// iconPathExempt reports whether a library-relative path is a canonical
// path-data definition.
func iconPathExempt(path string) bool {
	for _, ex := range iconPathExemptions {
		if strings.HasSuffix(filepath.ToSlash(path), filepath.ToSlash(ex.suffix)) {
			return true
		}
	}

	return false
}
