package utils

import (
	"path/filepath"
	"testing"
)

// TestMotionReduceCompliance verifies that every transition/animation class
// in .templ source files has a corresponding motion-reduce fallback.
// This prevents accessibility regressions where transitions are added without
// the motion-reduce safety net. HTML comments are stripped before scanning,
// so comment prose never trips the class regexes.
func TestMotionReduceCompliance(t *testing.T) {
	t.Parallel()

	// Directories to scan for .templ files (relative to project root)
	root := ".."
	dirs := []string{"display", "feedback", "forms", "navigation", "errorpage", "layout", "htmx"}

	violations := 0

	for _, dir := range dirs {
		findings, err := ScanMotionReduce(filepath.Join(root, dir))
		if err != nil {
			t.Logf("scan error for %s: %v", dir, err)

			continue
		}

		violations += len(findings)

		for _, f := range findings {
			t.Errorf("motion-reduce gap in %s:\n  %s", f.File, f.Snippet)
		}
	}

	if violations > 0 {
		t.Errorf("found %d motion-reduce compliance violations", violations)
	}
}
