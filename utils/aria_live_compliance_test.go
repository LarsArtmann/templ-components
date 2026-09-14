package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAriaLivePoliteness enforces the library's aria-live politeness policy
// (docs/aria-live-politeness.md, M21/F100): no component may render
// aria-live="assertive". Urgent page-integrated errors announce via
// role="alert" — which carries the assertive semantics with clearer intent —
// and everything transient or status-like is polite. An assertive live
// region interrupts the screen reader mid-sentence; for transient content
// (toasts, loading states, batch announcements) that is hostile, and for
// stacked toasts it compounds into chaos.
//
// The sweep is the drift-guard style used by TestMotionReduceCompliance: it
// scans .templ sources in the library packages so a future component cannot
// silently reintroduce assertive politeness.
func TestAriaLivePoliteness(t *testing.T) {
	t.Parallel()

	root := ".."
	dirs := []string{"display", "feedback", "forms", "navigation", "errorpage", "layout", "htmx", "datastar", "recipes"}

	violations := 0

	for _, dir := range dirs {
		dirPath := filepath.Join(root, dir)

		files, err := filepath.Glob(filepath.Join(dirPath, "*.templ"))
		if err != nil {
			t.Fatalf("glob %s: %v", dirPath, err)
		}

		for _, file := range files {
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("read %s: %v", file, err)
			}

			for lineNum, line := range strings.Split(string(data), "\n") {
				if strings.Contains(line, `aria-live="assertive"`) {
					t.Errorf(
						"%s:%d renders aria-live=\"assertive\" — policy violation (docs/aria-live-politeness.md). "+
							"Use role=\"alert\" for urgent page-integrated errors, polite for everything else.",
						filepath.Join(dir, filepath.Base(file)),
						lineNum+1,
					)

					violations++
				}
			}
		}
	}

	if violations == 0 {
		t.Log("aria-live politeness: OK (no assertive live regions in library .templ sources)")
	}
}
