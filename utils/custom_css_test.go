package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// TestCustomCSSUtilities asserts that every tc-* CSS class name referenced in
// .templ source files is defined in templates/custom.css.
//
// The tc-* prefix is the library's convention for custom CSS utilities and
// component-specific styles (overlays, scroll-snap, fluid typography, etc.).
// If a class is used in source but missing from custom.css, the component will
// render without its intended styling — and nothing else in the test suite
// will catch it until a consumer reports a visual regression.
//
// Exceptions are classes that serve as JS state hooks, HTMX toggle targets, or
// element IDs rather than custom.css style definitions.
func TestCustomCSSUtilities(t *testing.T) {
	t.Parallel()

	defined := loadCustomCSSClasses(t)

	requiredFluid := []string{
		"tc-fluid-display", "tc-fluid-h1", "tc-fluid-h2",
		"tc-fluid-h3", "tc-fluid-h4", "tc-fluid-lead",
	}

	for _, cls := range requiredFluid {
		if !defined[cls] {
			t.Errorf(
				"required CSS class .%s is missing from templates/custom.css — "+
					"fluid typography classes must not be deleted",
				cls,
			)
		}
	}

	scanDirs := []string{
		"display", "feedback", "forms", "navigation",
		"errorpage", "layout", "htmx", "datastar",
		"recipes", "charts/echarts", "examples/demo",
	}

	used := scanTemplForCSSClasses(t, scanDirs)

	violations := make([]string, 0)

	for cls, files := range used {
		if defined[cls] {
			continue
		}

		if _, ok := cssClassExceptions[cls]; ok {
			continue
		}

		violations = append(violations, fmt.Sprintf(
			"  .%s — used in: %s", cls, strings.Join(files, ", "),
		))
	}

	if len(violations) > 0 {
		sort.Strings(violations)

		t.Errorf(
			"%d tc-* CSS class(es) used in .templ files but not defined in "+
				"templates/custom.css (and not in the exceptions list):\n%s",
			len(violations),
			strings.Join(violations, "\n"),
		)
	}
}

// loadCustomCSSClasses reads templates/custom.css and returns a set of all
// .tc-* class selector names (without the leading dot).
func loadCustomCSSClasses(t *testing.T) map[string]bool {
	t.Helper()

	cssPath := filepath.Join("..", "templates", "custom.css")

	data, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("cannot read templates/custom.css: %v", err)
	}

	selectorRe := regexp.MustCompile(`\.tc-[a-z][a-z0-9-]*`)

	defined := make(map[string]bool)

	for _, match := range selectorRe.FindAllString(string(data), -1) {
		defined[strings.TrimPrefix(match, ".")] = true
	}

	return defined
}

// scanTemplForCSSClasses walks .templ files in the given directories and
// returns a map of tc-* CSS class names to the files where they are used.
//
// Only classes that appear as CSS class references are collected — data-tc-*
// attributes and --tc-* custom properties are filtered out.
func scanTemplForCSSClasses(t *testing.T, dirs []string) map[string][]string {
	t.Helper()

	// (?:^|[^a-z-]) ensures the char before tc- is NOT a lowercase letter
	// or hyphen, which filters out data-tc-* and --tc-* tokens.
	tokenRe := regexp.MustCompile(`(?:^|[^a-z-])(tc-[a-z][a-z0-9-]*)`)

	commentRe := regexp.MustCompile(`^\s*//`)

	used := make(map[string]map[string]bool)

	for _, dir := range dirs {
		dirPath := filepath.Join("..", dir)

		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if info.IsDir() || !strings.HasSuffix(path, ".templ") {
				return nil
			}

			if strings.Contains(path, "cmd"+string(filepath.Separator)+"tc") {
				return nil
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read file: %w", err)
			}

			relPath := strings.TrimPrefix(path, "../")

			for line := range strings.SplitSeq(string(data), "\n") {
				if commentRe.MatchString(line) {
					continue
				}

				for _, match := range tokenRe.FindAllStringSubmatch(line, -1) {
					cls := match[1]

					if used[cls] == nil {
						used[cls] = make(map[string]bool)
					}

					used[cls][relPath] = true
				}
			}

			return nil
		})
		if err != nil {
			t.Logf("walk error for %s: %v", dir, err)
		}
	}

	result := make(map[string][]string, len(used))

	for cls, fileSet := range used {
		files := make([]string, 0, len(fileSet))
		for file := range fileSet {
			files = append(files, file)
		}

		sort.Strings(files)
		result[cls] = files
	}

	return result
}

// cssClassExceptions lists tc-* identifiers that appear in .templ files but
// are NOT CSS class definitions in templates/custom.css. Each has an explicit
// reason documenting why it does not need a custom.css definition.
var cssClassExceptions = map[string]string{
	"tc-echarts":           "wrapper class for ECharts container div, styled inline by component",
	"tc-menu-open":         "JS state class toggled by singleton handler, styled via Tailwind utilities",
	"tc-menu-close":        "JS state class toggled by singleton handler, styled via Tailwind utilities",
	"tc-btn-loading":       "HTMX loading state class, styled via Tailwind htmx-request utilities",
	"tc-toast-container":   "element ID (id=), not a CSS class",
	"tc-loading-indicator": "element ID (id=), not a CSS class",
	"tc-error-announcer":   "element ID (id=), not a CSS class",
	"tc-form-loading":      "referenced in comment only, not a CSS class",
}
