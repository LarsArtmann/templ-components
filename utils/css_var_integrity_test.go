package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// varRefPattern captures CSS custom-property REFERENCES, including those with
// fallback values: var(--x) and var(--x, fallback).
var varRefPattern = regexp.MustCompile(`var\((--[a-zA-Z0-9-]+)`)

// varDefPattern captures CSS custom-property DEFINITIONS: "--x:" (declaration)
// in stylesheets and Go/templ string literals that build inline styles.
var varDefPattern = regexp.MustCompile(`(--[a-zA-Z0-9-]+)\s*:`)

// libraryPackages are the component-source directories whose emitted HTML the
// CSS-var integrity sweep covers.
var libraryPackages = []string{
	"display", "errorpage", "feedback", "forms", "htmx",
	"layout", "navigation", "recipes",
}

// TestHeatmapBrandVarsDefined pins the exact silent-failure class the 2026-09-08
// demo audit caught: display.Heatmap colors cells via
// rgba(var(--ds-brand-rgb), alpha) and ring-[var(--ds-brand)], but the tokens
// were defined NOWHERE — every cell rendered transparent (only the peak ring
// showed). templates/custom.css must define both tokens (light + dark).
func TestHeatmapBrandVarsDefined(t *testing.T) {
	t.Parallel()

	css, err := os.ReadFile(filepath.Join("..", "templates", "custom.css"))
	if err != nil {
		t.Fatalf("read templates/custom.css: %v", err)
	}

	for _, token := range []string{"--ds-brand-rgb", "--ds-brand"} {
		if got := strings.Count(string(css), token+":"); got < 2 {
			t.Errorf(
				"templates/custom.css defines %s %d time(s), want >= 2 (light + dark) — the Heatmap renders transparent without it (2026-09-08 audit f23/f24)",
				token,
				got,
			)
		}
	}
}

// collectVarDefsFromCSS reads one stylesheet and adds every custom-property
// definition it declares. A missing optional file is not an error.
func collectVarDefsFromCSS(defs map[string]bool, file string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // optional file (e.g. compiled CSS absent in fresh clones)
		}

		return fmt.Errorf("read %s: %w", file, err)
	}

	for _, def := range varDefPattern.FindAllStringSubmatch(string(content), -1) {
		defs[def[1]] = true
	}

	return nil
}

// collectComponentVarRefs walks one library package's component sources and
// committed golden renderings, collecting var(--token) references plus any
// inline-style definitions the components emit themselves (e.g. AppShell's
// --tc-sidebar-w).
func collectComponentVarRefs(defs map[string]bool) (map[string][]string, error) {
	refs := map[string][]string{}

	for _, pkg := range libraryPackages {
		err := filepath.Walk(filepath.Join("..", pkg), func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil || info.IsDir() {
				return walkErr
			}

			name := info.Name()
			if !isSweepableSource(name) {
				return nil
			}

			//nolint:gosec // fixed repo-relative test paths walked from a constant list; no untrusted input
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read %s: %w", path, err)
			}

			for _, ref := range varRefPattern.FindAllStringSubmatch(string(content), -1) {
				refs[ref[1]] = append(refs[ref[1]], path)
			}

			if strings.Contains(string(content), "\"--") {
				for _, def := range varDefPattern.FindAllStringSubmatch(string(content), -1) {
					if strings.Contains(string(content), "\""+def[1]+":") {
						defs[def[1]] = true
					}
				}
			}

			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walk %s: %w", pkg, err)
		}
	}

	return refs, nil
}

// isSweepableSource reports whether a walked file is a component source or a
// committed golden rendering (generated and test files are excluded).
func isSweepableSource(name string) bool {
	if strings.HasSuffix(name, "_templ.go") || strings.HasSuffix(name, "_test.go") {
		return false
	}

	return strings.HasSuffix(name, ".templ") ||
		strings.HasSuffix(name, ".go") ||
		strings.HasSuffix(name, ".golden")
}

// TestNoUndefinedCSSVarReferences sweeps every library component source AND
// the committed golden renderings for var(--token) references that no
// stylesheet defines. An undefined custom property silently resolves to
// nothing: colors vanish, layouts collapse — the Heatmap --ds-brand-rgb bug
// class. Tokens may be defined in templates/custom.css, the compiled demo CSS
// (Tailwind theme tokens like --color-amber-600, --font-sans), the optional
// semantic theme file, or inline styles emitted by a component itself
// (--tc-sidebar-w).
func TestNoUndefinedCSSVarReferences(t *testing.T) {
	t.Parallel()

	defs := map[string]bool{}
	for _, file := range []string{
		filepath.Join("..", "templates", "custom.css"),
		filepath.Join("..", "templ-components-theme.css"),
		filepath.Join("..", "examples", "demo", "static", "app.css"),
	} {
		if err := collectVarDefsFromCSS(defs, file); err != nil {
			t.Fatal(err)
		}
	}

	refs, err := collectComponentVarRefs(defs)
	if err != nil {
		t.Fatal(err)
	}

	customCSS, err := os.ReadFile(filepath.Join("..", "templates", "custom.css"))
	if err != nil {
		t.Fatalf("read templates/custom.css: %v", err)
	}

	// custom.css may reference tokens defined by the Tailwind theme or itself.
	for _, ref := range varRefPattern.FindAllStringSubmatch(string(customCSS), -1) {
		refs[ref[1]] = append(refs[ref[1]], "templates/custom.css")
	}

	for token, origins := range refs {
		if strings.HasPrefix(token, "--tw-") || defs[token] {
			continue // Tailwind internal runtime tokens / defined somewhere
		}

		t.Errorf(
			"CSS custom property %s is referenced (e.g. %s) but defined NOWHERE (custom.css, compiled demo CSS, semantic theme, or inline styles) — it silently renders as nothing. Define it in templates/custom.css or emit it inline.",
			token,
			origins[0],
		)
	}
}
