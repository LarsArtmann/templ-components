package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestRTLLogicalProperties verifies that .templ source files in the library
// packages use CSS logical properties (ms-, me-, ps-, pe-, text-start,
// text-end, border-s-, border-e-) instead of physical directional properties
// (ml-, mr-, pl-, pr-, text-left, text-right, border-l-, border-r-).
// Logical properties automatically mirror in RTL (dir="rtl") without code
// changes; physical properties do not and cause broken RTL layouts.
//
// This is a FAILING test — violations block CI.
//
// Run via: go test ./utils/... -run TestRTL.
func TestRTLLogicalProperties(t *testing.T) {
	t.Parallel()

	root := ".."
	dirs := []string{"display", "feedback", "forms", "navigation", "errorpage", "layout", "htmx"}

	// Physical properties that have direct logical equivalents.
	// These have ZERO exceptions — every occurrence is a violation.
	physicalRe := regexp.MustCompile(
		`\b(ml-|mr-|pl-|pr-|text-left|text-right|border-l-|border-r-)`,
	)

	violations := 0

	for _, dir := range dirs {
		err := filepath.Walk(filepath.Join(root, dir), func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}

			if !strings.HasSuffix(path, ".templ") {
				return nil
			}

			data, readErr := os.ReadFile(path) //nolint:gosec // test scans templ files
			if readErr != nil {
				return fmt.Errorf("read file: %w", readErr)
			}

			for line := range strings.SplitSeq(string(data), "\n") {
				if !physicalRe.MatchString(line) {
					continue
				}

				violations++

				t.Errorf("RTL physical-property violation in %s:\n  %s", path, strings.TrimSpace(line))
			}

			return nil
		})
		if err != nil {
			t.Logf("walk error for %s: %v", dir, err)
		}
	}

	if violations > 0 {
		t.Errorf("found %d RTL logical-property violations", violations)
	}
}
