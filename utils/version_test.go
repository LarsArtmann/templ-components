package utils

import (
	"os"
	"strings"
	"testing"
)

// TestVersionMatchesChangelog ensures utils.Version stays in sync with the
// latest released version declared in CHANGELOG.md. When you bump Version for
// a release, add the matching CHANGELOG entry (or vice versa) — this test
// fails if they drift apart.
func TestVersionMatchesChangelog(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("../CHANGELOG.md")
	if err != nil {
		t.Skipf("CHANGELOG.md not found (running outside repo root?): %v", err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "## [") {
			continue
		}
		// Skip the [Unreleased] section — look for the first real version.
		if strings.Contains(line, "[Unreleased]") {
			continue
		}
		// line looks like: ## [0.4.0] — 2026-06-27
		start := strings.Index(line, "[") + 1
		end := strings.Index(line, "]")
		if start < 1 || end <= start {
			t.Fatalf("could not parse version heading: %q", line)
		}
		want := line[start:end]
		if want != Version {
			t.Errorf("utils.Version = %q, but CHANGELOG.md latest release is %q", Version, want)
		}
		return
	}
	t.Fatalf("no released version heading found in CHANGELOG.md")
}
