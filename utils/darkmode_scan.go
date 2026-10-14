package utils

import (
	"path/filepath"
	"regexp"
	"strings"
)

// darkModeException defines a file/line pattern that is intentionally
// styled the same in both light and dark mode.
type darkModeException struct {
	fileSubstring string
	lineSubstring string
}

// darkModeExceptions lists the library's intentional same-in-both-mode
// patterns (documented in AGENTS.md). The table matches on library-relative
// path SUBSTRING, so it only ever fires on the library's own paths — a
// consumer's tree needs the per-line `tc-lint:ignore` allowlist instead.
var darkModeExceptions = []darkModeException{
	// Toggle thumb: bg-white in both modes (track changes color instead)
	{"forms/toggle.templ", "bg-white shadow-sm"},
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

// darkModeNeutralRe matches neutral color classes (text-gray-*, bg-white,
// bg-gray-*, border-gray-*, divide-gray-*, ring-gray-*).
var darkModeNeutralRe = regexp.MustCompile(
	`(text-gray-[0-9]+|bg-white|bg-gray-[0-9]+|border-gray-[0-9]+|divide-gray-[0-9]+|ring-gray-[0-9]+)`,
)

// darkModeSemanticRe matches semantic color classes (bg-blue-600, bg-red-600,
// text-blue-600, text-green-500, etc.).
var darkModeSemanticRe = regexp.MustCompile(
	`(bg-(blue|red|green|amber|orange|gray)-[0-9]+|text-(blue|red|green|amber|orange)-[0-9]+)`,
)

// goCommentLineRe skips Go-style comment lines (both languages' sources are
// swept; class-shaped prose in comments is not markup).
var goCommentLineRe = regexp.MustCompile(`^\s*//`)

// hasDarkVariantRe is the same-line dark: presence check.
var hasDarkVariantRe = regexp.MustCompile(`dark:`)

// darkModeRules pairs each dark-mode rule with its color regex.
var darkModeRules = []struct {
	rule string
	re   *regexp.Regexp
}{
	{rule: "darkmode-neutral", re: darkModeNeutralRe},
	{rule: "darkmode-semantic", re: darkModeSemanticRe},
}

// checkLineForDarkModeGap checks a single line for dark mode gaps.
// Returns true if a violation was found.
func checkLineForDarkModeGap(path, line string, colorRe *regexp.Regexp) bool {
	if goCommentLineRe.MatchString(line) {
		return false
	}

	matches := colorRe.FindAllString(line, -1)
	if len(matches) == 0 {
		return false
	}

	if hasDarkVariantRe.MatchString(line) {
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

// ScanDarkMode walks dir recursively and reports color classes without a
// dark: variant on the same line, in .templ and .go component sources
// (generated *_templ.go and *_test.go files are excluded). Two rules are
// reported: "darkmode-neutral" (grays/white) and "darkmode-semantic"
// (blue/red/green/amber/orange).
//
// It is the engine behind TestDarkModeCompliance + TestDarkModeSemanticColors
// (repo guards) and `tc lint --dark`. The library's documented same-in-both-
// mode exceptions are built in; consumers suppress their own exceptions with
// the per-line `tc-lint:ignore` allowlist. Findings use dir-relative paths.
func ScanDarkMode(dir string) ([]Finding, error) {
	var findings []Finding

	err := scanSources(dir, func(name string) bool {
		return (strings.HasSuffix(name, ".templ") || strings.HasSuffix(name, ".go")) &&
			!strings.HasSuffix(name, "_templ.go") && !strings.HasSuffix(name, "_test.go")
	}, func(rel, content string) error {
		full := filepath.Join(dir, rel)

		for i, line := range strings.Split(content, "\n") {
			for _, r := range darkModeRules {
				if checkLineForDarkModeGap(full, line, r.re) {
					findings = append(findings, Finding{
						File:    rel,
						Line:    i + 1,
						Rule:    r.rule,
						Message: "dark mode gap (color class without a dark: variant on the same line)",
						Snippet: strings.TrimSpace(line),
					})
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return findings, nil
}
