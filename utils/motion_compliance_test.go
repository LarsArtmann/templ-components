package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// motionReduceRules bundles the regexes the motion-reduce sweep applies per line.
type motionReduceRules struct {
	transitionRe      *regexp.Regexp
	animateRe         *regexp.Regexp
	motionReduceRe    *regexp.Regexp
	transitionConstRe *regexp.Regexp
}

// TestMotionReduceCompliance verifies that every transition/animation class
// in .templ source files has a corresponding motion-reduce fallback.
// This prevents accessibility regressions where transitions are added without
// the motion-reduce safety net.
func TestMotionReduceCompliance(t *testing.T) {
	t.Parallel()

	// Directories to scan for .templ files (relative to project root)
	root := ".."
	dirs := []string{"display", "feedback", "forms", "navigation", "errorpage", "layout", "htmx"}

	rules := motionReduceRules{
		transitionRe:   regexp.MustCompile(`transition[-:]`),
		animateRe:      regexp.MustCompile(`animate[-:]`),
		motionReduceRe: regexp.MustCompile(`motion-reduce:`),
		// Shared motion constants that already include motion-reduce fallbacks.
		transitionConstRe: regexp.MustCompile(
			`(utils\.Transition|transitionFast|transitionNormal|transitionColors|transitionTransform)`,
		),
	}

	violations := 0

	for _, dir := range dirs {
		count, err := countMotionGapsInDir(t, filepath.Join(root, dir), rules)
		violations += count

		if err != nil {
			t.Logf("walk error for %s: %v", dir, err)
		}
	}

	if violations > 0 {
		t.Errorf("found %d motion-reduce compliance violations", violations)
	}
}

// countMotionGapsInDir walks one package directory for .templ files and
// reports motion-reduce gaps via t.Errorf, returning the violation count.
func countMotionGapsInDir(t *testing.T, dir string, rules motionReduceRules) (int, error) {
	t.Helper()

	violations := 0

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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

		violations += countMotionGapsInContent(t, path, string(data), rules)

		return nil
	})

	return violations, err
}

// countMotionGapsInContent checks each line that has transition or animate
// classes and reports the ones lacking a motion-reduce fallback.
func countMotionGapsInContent(t *testing.T, path, content string, rules motionReduceRules) int {
	t.Helper()

	violations := 0

	for line := range strings.SplitSeq(content, "\n") {
		if !rules.needsMotionReduce(line) {
			continue
		}

		violations++

		t.Errorf("motion-reduce gap in %s:\n  %s", path, strings.TrimSpace(line))
	}

	return violations
}

// needsMotionReduce reports whether a line carries transition/animate classes
// without an inline motion-reduce fallback or a shared motion-constant
// reference (the constants already include motion-reduce fallbacks).
func (r motionReduceRules) needsMotionReduce(line string) bool {
	if !r.transitionRe.MatchString(line) && !r.animateRe.MatchString(line) {
		return false
	}

	if r.motionReduceRe.MatchString(line) {
		return false
	}

	return !r.transitionConstRe.MatchString(line)
}
