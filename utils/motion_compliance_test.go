package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestMotionReduceCompliance verifies that every transition/animation class
// in .templ source files has a corresponding motion-reduce fallback.
// This prevents accessibility regressions where transitions are added without
// the motion-reduce safety net.
func TestMotionReduceCompliance(t *testing.T) {
	t.Parallel()

	// Directories to scan for .templ files (relative to project root)
	root := ".."
	dirs := []string{"display", "feedback", "forms", "navigation", "errorpage", "layout", "htmx"}

	// Patterns that require motion-reduce
	transitionRe := regexp.MustCompile(`transition[-:]`)
	animateRe := regexp.MustCompile(`animate[-:]`)

	// motion-reduce presence check
	motionReduceRe := regexp.MustCompile(`motion-reduce:`)

	// Shared motion constants that already include motion-reduce fallbacks.
	// Lines referencing these are exempt from the inline motion-reduce check.
	transitionConstRe := regexp.MustCompile(
		`(utils\.Transition|transitionFast|transitionNormal|transitionColors|transitionTransform)`,
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

			content := string(data)

			// Check each line that has transition or animate
			for line := range strings.SplitSeq(content, "\n") {
				hasTransition := transitionRe.MatchString(line)

				hasAnimate := animateRe.MatchString(line)
				if !hasTransition && !hasAnimate {
					continue
				}
				// Check if the same line or a nearby line has motion-reduce
				if motionReduceRe.MatchString(line) {
					continue
				}
				// Lines referencing shared motion constants are safe — the
				// constants already include motion-reduce fallbacks.
				if transitionConstRe.MatchString(line) {
					continue
				}
				// Allow multi-line class strings — check if motion-reduce appears
				// within the same class attribute context. For single-line violations:
				violations++

				t.Errorf("motion-reduce gap in %s:\n  %s", path, strings.TrimSpace(line))
			}

			return nil
		})
		if err != nil {
			t.Logf("walk error for %s: %v", dir, err)
		}
	}

	if violations > 0 {
		t.Errorf("found %d motion-reduce compliance violations", violations)
	}
}
