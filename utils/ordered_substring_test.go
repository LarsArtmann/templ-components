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

	tailwindTokenRe := regexp.MustCompile(`^[a-z][a-z0-9]*(?::[a-z][a-z0-9]*)*(?:-[a-z0-9\[\]/._%]+)*$`)
	containsCallRe := regexp.MustCompile(`strings\.Contains\s*\(`)
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

			scanLinesForOrderedSubstrings(
				t, path, string(data),
				tailwindTokenRe, containsCallRe, stringLiteralRe,
				&violations,
			)

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

// scanLinesForOrderedSubstrings walks each line of a test file and flags any
// strings.Contains call whose string literal contains multiple Tailwind-class
// tokens in a fixed order.
func scanLinesForOrderedSubstrings(
	t *testing.T,
	path, data string,
	tailwindTokenRe, containsCallRe, stringLiteralRe *regexp.Regexp,
	violations *int,
) {
	t.Helper()

	lineNum := 0

	for line := range strings.SplitSeq(data, "\n") {
		lineNum++

		if !containsCallRe.MatchString(line) {
			continue
		}

		for _, match := range stringLiteralRe.FindAllString(line, -1) {
			literal := match[1 : len(match)-1]
			if !isOrderedTailwindSubstring(literal, tailwindTokenRe) {
				continue
			}

			*violations++

			t.Errorf(
				"ordered Tailwind substring in %s:%d\n  %s\n  literal: %q\n  use AssertContainsAll or single-token AssertContains instead",
				path,
				lineNum,
				strings.TrimSpace(line),
				literal,
			)
		}
	}
}

// isOrderedTailwindSubstring reports whether a string literal contains two or
// more space-separated tokens that all look like Tailwind utility classes, with
// at least one containing a hyphen (excludes all-lowercase English phrases).
func isOrderedTailwindSubstring(literal string, tailwindTokenRe *regexp.Regexp) bool {
	if !strings.Contains(literal, " ") {
		return false
	}

	tokens := strings.Fields(literal)
	if len(tokens) < 2 {
		return false
	}

	hasHyphen := false

	for _, tok := range tokens {
		if !tailwindTokenRe.MatchString(tok) {
			return false
		}

		if strings.Contains(tok, "-") {
			hasHyphen = true
		}
	}

	return hasHyphen
}
