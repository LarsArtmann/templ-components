package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// darkModeException defines a file/line pattern that is intentionally
// styled the same in both light and dark mode.
type darkModeException struct {
	fileSubstring string
	lineSubstring string
}

// darkModeExceptions lists all intentional same-in-both-mode patterns.
var darkModeExceptions = []darkModeException{
	// Toggle thumb: bg-white in both modes (track changes color instead)
	{"forms/toggle.templ", "bg-white shadow-sm"},
	// SidebarNav: permanently dark sidebar
	{"navigation/sidebar_nav.templ", "hover:bg-gray-800 hover:text-white"},
	{"navigation/sidebar_nav.templ", "bg-gray-900 dark:bg-black"},
	{"navigation/sidebar_nav.templ", "text-white"},
	{"navigation/sidebar_nav.templ", "bg-gray-800"},
	// JS class arrays for overlay transitions — not CSS color classes
	{"display/shared.go", "openClasses"},
	{"display/shared.go", "closeClasses"},
	// Toggle track — dark:peer-checked:bg-blue-500 is on the same line
	{"forms/toggle.templ", "peer-checked:bg-blue-600"},
	// Feedback style icons already have dark: variants via lookupFeedbackStyle
	{"feedback/styles.go", "Icon:"},
	// Badge dot colors are decorative — solid background with text on top
	{"display/badge.templ", "bg-green-500"},
	{"display/badge.templ", "bg-yellow-500"},
	{"display/badge.templ", "bg-red-500"},
	{"display/badge.templ", "bg-blue-500"},
	// Count badge background — solid red with white text, readable in both modes
	{"display/count_badge.templ", "bg-red-500"},
	// Action buttons have dark: variants on the same line
	{"errorpage/styles.go", "ActionButton:"},
	// Avatar silhouette icon — decorative, light-on-blue in both modes
	{"display/avatar.templ", "text-blue-200"},
}

// isDarkModeException checks if a file/line matches any exception.
func isDarkModeException(path, line string) bool {
	for _, ex := range darkModeExceptions {
		if strings.Contains(path, ex.fileSubstring) && strings.Contains(line, ex.lineSubstring) {
			return true
		}
	}

	return false
}

// darkPrefixes lists all dark: variant prefixes that can precede a color class.
var darkPrefixes = []string{
	"dark:",
	"dark:hover:",
	"dark:focus:",
	"dark:focus-visible:",
	"dark:placeholder:",
	"dark:group-hover:",
	"dark:peer-checked:",
	"dark:file:",
	"dark:selection:",
	"dark:group-focus:",
}

// isWithinDarkVariant checks if a color class at position idx in line is
// preceded by a dark: variant prefix.
func isWithinDarkVariant(line string, idx int) bool {
	for _, prefix := range darkPrefixes {
		prefixLen := len(prefix)
		if idx >= prefixLen && line[idx-prefixLen:idx] == prefix {
			return true
		}
	}

	return false
}

// allColorsInDarkVariant checks if ALL color matches on a line are inside
// dark: variant prefixes (meaning the line already has proper dark: coverage).
func allColorsInDarkVariant(line string, matches []string) bool {
	for _, m := range matches {
		idx := strings.Index(line, m)
		if idx < 0 {
			continue
		}

		if isWithinDarkVariant(line, idx) {
			continue
		}

		return false
	}

	return true
}

// checkLineForDarkModeGap checks a single line for dark mode gaps.
// Returns true if a violation was found.
func checkLineForDarkModeGap(path, line string, colorRe *regexp.Regexp) bool {
	commentRe := regexp.MustCompile(`^\s*//`)
	if commentRe.MatchString(line) {
		return false
	}

	matches := colorRe.FindAllString(line, -1)
	if len(matches) == 0 {
		return false
	}

	darkRe := regexp.MustCompile(`dark:`)
	if darkRe.MatchString(line) {
		return false
	}

	if isDarkModeException(path, line) {
		return false
	}

	if allColorsInDarkVariant(line, matches) {
		return false
	}

	return true
}

// scanDarkMode scans .templ and .go files (excluding generated/test files)
// for color classes without dark: variants, reporting violations.
func scanDarkMode(t *testing.T, dirs []string, colorRe *regexp.Regexp) {
	t.Helper()

	violations := 0

	for _, dir := range dirs {
		err := filepath.Walk(filepath.Join("..", dir), func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}

			if !strings.HasSuffix(path, ".templ") && !strings.HasSuffix(path, ".go") {
				return nil
			}

			if strings.HasSuffix(path, "_templ.go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}

			data, readErr := os.ReadFile(path) //nolint:gosec // test scans source files
			if readErr != nil {
				return fmt.Errorf("read file: %w", readErr)
			}

			for line := range strings.SplitSeq(string(data), "\n") {
				if checkLineForDarkModeGap(path, line, colorRe) {
					violations++

					t.Errorf("dark mode gap in %s:\n  %s", path, strings.TrimSpace(line))
				}
			}

			return nil
		})
		if err != nil {
			t.Logf("walk error for %s: %v", dir, err)
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
	neutralColorRe := regexp.MustCompile(
		`(text-gray-[0-9]+|bg-white|bg-gray-[0-9]+|border-gray-[0-9]+|divide-gray-[0-9]+|ring-gray-[0-9]+)`,
	)
	scanDarkMode(t, dirs, neutralColorRe)
}

// TestDarkModeSemanticColors verifies that semantic color classes (bg-blue-600,
// bg-red-600, text-blue-600, text-green-500, etc.) in .templ and .go source
// files have corresponding dark: variants on the same line.
func TestDarkModeSemanticColors(t *testing.T) {
	t.Parallel()

	dirs := []string{"display", "feedback", "forms", "navigation", "errorpage", "layout", "htmx", "examples/demo"}
	semanticColorRe := regexp.MustCompile(
		`(bg-(blue|red|green|amber|orange|gray)-[0-9]+|text-(blue|red|green|amber|orange)-[0-9]+)`,
	)
	scanDarkMode(t, dirs, semanticColorRe)
}
