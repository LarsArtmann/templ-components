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
// It also bans the PHYSICAL-positioning inset utilities start-* / end-*
// (e.g. start-0, end-full): they resolve to left-*/right-* and never mirror
// in RTL. Use the logical inset-s-* / inset-e-* forms instead. The regex's
// preceding-character class keeps longer classes safe (inset-s-0,
// text-start, items-start, justify-end, col-start-2, rounded-s-*).
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

	// Deprecated physical-positioning inset utilities: start-*/end-* are
	// v4.2-deprecated ALIASES of inset-s-*/inset-e-* (probe-verified 2026-10-08:
	// identical logical output today — the ban is future-proofing against the
	// aliases' removal, and keeps the canon on one form).
	insetRe := regexp.MustCompile(`(?:^|[\s"'({])(start|end)-[a-z0-9.]`)

	violations := 0

	for _, dir := range dirs {
		err := filepath.Walk(filepath.Join(root, dir), func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}

			if !strings.HasSuffix(path, ".templ") {
				return nil
			}

			data, readErr := os.ReadFile(path)
			if readErr != nil {
				return fmt.Errorf("read file: %w", readErr)
			}

			for line := range strings.SplitSeq(string(data), "\n") {
				switch {
				case physicalRe.MatchString(line):
					violations++

					t.Errorf("RTL physical-property violation in %s:\n  %s", path, strings.TrimSpace(line))
				case insetRe.MatchString(line):
					violations++

					t.Errorf(
						"RTL physical-inset violation (use inset-s-*/inset-e-* instead of start-*/end-*) in %s:\n  %s",
						path,
						strings.TrimSpace(line),
					)
				}
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
