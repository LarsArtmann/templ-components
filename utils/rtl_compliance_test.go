package utils

import (
	"path/filepath"
	"testing"
)

// TestRTLLogicalProperties verifies that .templ source files in the library
// packages use CSS logical properties (ms-, me-, ps-, pe-, text-start,
// text-end, border-s-, border-e-) instead of physical directional properties
// (ml-, mr-, pl-, pr-, text-left, text-right, border-l-, border-r-).
// Logical properties automatically mirror in RTL (dir="rtl") without code
// changes; physical properties do not and cause broken RTL layouts.
//
// It also bans the v4.2-DEPRECATED inset aliases start-* / end-* (e.g.
// start-0, end-full): probe-verified 2026-10-08 as output-identical LOGICAL
// aliases of inset-s-* / inset-e-* today, but a future major can drop the
// aliases, so the library stays on the canonical forms. HTML comments are
// stripped before scanning, so prose never trips the class regexes. The
// regex's preceding-character class keeps longer classes safe (inset-s-0,
// text-start, items-start, justify-end, col-start-2, rounded-s-*).
//
// This is a FAILING test — violations block CI.
//
// Run via: go test ./utils/... -run TestRTL.
func TestRTLLogicalProperties(t *testing.T) {
	t.Parallel()

	root := ".."
	dirs := []string{"display", "feedback", "forms", "navigation", "errorpage", "layout", "htmx"}

	violations := 0

	for _, dir := range dirs {
		findings, err := ScanRTL(filepath.Join(root, dir))
		if err != nil {
			t.Logf("scan error for %s: %v", dir, err)

			continue
		}

		for _, f := range findings {
			violations++

			t.Errorf("%s in %s:\n  %s", f.Message, f.File, f.Snippet)
		}
	}

	if violations > 0 {
		t.Errorf("found %d RTL logical-property violations", violations)
	}
}
