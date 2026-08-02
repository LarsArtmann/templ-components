package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestNoOrderedTailwindSubstringsInTests is a drift-guard that prevents
// reintroducing brittle ordered Tailwind-substring assertions in tests.
//
// utils.Class() wraps tailwind-merge-go, which reorders classes
// nondeterministically (depending on internal LRU cache state). Any assertion
// like strings.Contains(out, "flex flex-col") is therefore a latent flake: under
// one cache state the classes appear in the asserted order and the test passes;
// under another they are reordered and the test fails (or, for negative
// assertions, a real regression is masked). The fix is always to use
// AssertContainsAll (set semantics) or AssertContains/AssertNotContains with a
// single class token.
//
// This test scans every *_test.go file in the library packages and fails if it
// finds a strings.Contains call whose search string contains two or more
// space-separated tokens that all look like Tailwind utility classes.
//
// Run via: go test ./utils/... -run TestNoOrderedTailwind.
func TestNoOrderedTailwindSubstringsInTests(t *testing.T) {
	t.Parallel()

	root := ".."
	dirs := []string{
		"display", "feedback", "forms", "navigation", "errorpage",
		"layout", "htmx", "datastar", "icons", "recipes",
		"integration", "internal",
	}

	// A Tailwind-class-like token: lowercase letters/digits with optional
	// variant prefixes (dark:, sm:, hover:), hyphens, brackets, slashes,
	// colons, dots, underscores, and percentages. Crucially NO uppercase and
	// NO spaces — this excludes English phrases ("Page not found") and HTML
	// fragments.
	tailwindTokenRe := regexp.MustCompile(`^[a-z][a-z0-9]*(?::[a-z][a-z0-9]*)*(?:-[a-z0-9\[\]/._%]+)*$`)

	// Find strings.Contains( calls on a line.
	containsCallRe := regexp.MustCompile(`strings\.Contains\s*\(`)

	// Extract double-quoted and backtick string literals.
	stringLiteralRe := regexp.MustCompile("(\"[^\"]*\"|`[^`]*`)")

	violations := 0

	for _, dir := range dirs {
		walkErr := filepath.Walk(filepath.Join(root, dir), func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}
			if !strings.HasSuffix(path, "_test.go") {
				return nil
			}

			data, readErr := os.ReadFile(path) //nolint:gosec // test scans source files
			if readErr != nil {
				return fmt.Errorf("read file: %w", readErr)
			}

			lineNum := 0
			for line := range strings.SplitSeq(string(data), "\n") {
				lineNum++
				if !containsCallRe.MatchString(line) {
					continue
				}

				for _, match := range stringLiteralRe.FindAllString(line, -1) {
					// Strip surrounding quotes/backticks.
					literal := match[1 : len(match)-1]
					if !strings.Contains(literal, " ") {
						continue
					}

					tokens := strings.Fields(literal)
					if len(tokens) < 2 {
						continue
					}

					allTailwind := true
					hasHyphen := false
					for _, tok := range tokens {
						if !tailwindTokenRe.MatchString(tok) {
							allTailwind = false
							break
						}
						if strings.Contains(tok, "-") {
							hasHyphen = true
						}
					}

					// Flag only when every token looks like a Tailwind class
					// and at least one has a hyphen (excludes all-lowercase
					// English phrases with no hyphens).
					if allTailwind && hasHyphen {
						violations++
						t.Errorf("ordered Tailwind substring in %s:%d\n  %s\n  literal: %q\n  use AssertContainsAll or single-token AssertContains instead",
							path, lineNum, strings.TrimSpace(line), literal)
					}
				}
			}
			return nil
		})
		if walkErr != nil {
			t.Fatalf("walk %s: %v", dir, walkErr)
		}
	}

	if violations > 0 {
		t.Fatalf("found %d ordered Tailwind substring assertion(s); see above", violations)
	}
}
