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

	assertRequiredFluidClasses(t, defined)

	scanDirs := []string{
		"display", "feedback", "forms", "navigation",
		"errorpage", "layout", "htmx", "datastar",
		"recipes", "charts/echarts", "examples/demo",
	}

	used := scanTemplForCSSClasses(t, scanDirs)

	assertUsedClassesDefined(t, used, defined)
}

func assertRequiredFluidClasses(t *testing.T, defined map[string]bool) {
	t.Helper()

	required := []string{
		"tc-fluid-display", "tc-fluid-h1", "tc-fluid-h2",
		"tc-fluid-h3", "tc-fluid-h4", "tc-fluid-lead",
	}

	for _, cls := range required {
		if !defined[cls] {
			t.Errorf(
				"required CSS class .%s is missing from templates/custom.css — "+
					"fluid typography classes must not be deleted",
				cls,
			)
		}
	}
}

func assertUsedClassesDefined(
	t *testing.T,
	used map[string][]string,
	defined map[string]bool,
) {
	t.Helper()

	var violations []string

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

// cssTokenRe matches tc-* CSS class names while filtering out data-tc-*
// attributes and --tc-* custom properties. The non-capturing prefix ensures
// the character before tc- is NOT a lowercase letter or hyphen.
var cssTokenRe = regexp.MustCompile(`(?:^|[^a-z-])(tc-[a-z][a-z0-9-]*)`)

var cssCommentRe = regexp.MustCompile(`^\s*//`)

// scanTemplForCSSClasses walks .templ files in the given directories and
// returns a map of tc-* CSS class names to the files where they are used.
func scanTemplForCSSClasses(t *testing.T, dirs []string) map[string][]string {
	t.Helper()

	used := make(map[string]map[string]bool)

	for _, dir := range dirs {
		walkCSSClassDir(t, filepath.Join("..", dir), used)
	}

	return flattenCSSClassMap(used)
}

func walkCSSClassDir(t *testing.T, dirPath string, used map[string]map[string]bool) {
	t.Helper()

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() || !strings.HasSuffix(path, ".templ") {
			return walkErr
		}

		if strings.Contains(path, filepath.Join("cmd", "tc")) {
			return nil
		}

		collectCSSClassesFromFile(path, used)

		return nil
	})
	if err != nil {
		t.Logf("walk error for %s: %v", dirPath, err)
	}
}

func collectCSSClassesFromFile(path string, used map[string]map[string]bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	relPath := strings.TrimPrefix(path, "../")

	for line := range strings.SplitSeq(string(data), "\n") {
		if cssCommentRe.MatchString(line) {
			continue
		}

		for _, match := range cssTokenRe.FindAllStringSubmatch(line, -1) {
			cls := match[1]

			if used[cls] == nil {
				used[cls] = make(map[string]bool)
			}

			used[cls][relPath] = true
		}
	}
}

func flattenCSSClassMap(used map[string]map[string]bool) map[string][]string {
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
	"tc-echarts":            "wrapper class for ECharts container div, styled inline by component",
	"tc-menu-open":          "JS state class toggled by singleton handler, styled via Tailwind utilities",
	"tc-menu-close":         "JS state class toggled by singleton handler, styled via Tailwind utilities",
	"tc-btn-loading":        "HTMX loading state class, styled via Tailwind htmx-request utilities",
	"tc-toast-container":    "element ID (id=), not a CSS class",
	"tc-loading-indicator":  "element ID (id=), not a CSS class",
	"tc-error-announcer":    "element ID (id=), not a CSS class",
	"tc-datastar-announcer": "element ID (id=), not a CSS class",
	"tc-form-loading":       "referenced in comment only, not a CSS class",
}
