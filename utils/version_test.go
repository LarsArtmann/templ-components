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

	for line := range strings.SplitSeq(string(data), "\n") {
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

// TestVersionMatchesFeatures ensures utils.Version stays in sync with the
// version declared in FEATURES.md. Documentation drift (e.g. FEATURES.md still
// showing an old version after a release) is caught here.
func TestVersionMatchesFeatures(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("../FEATURES.md")
	if err != nil {
		t.Skipf("FEATURES.md not found (running outside repo root?): %v", err)
	}

	const marker = "**Version:**"
	for line := range strings.SplitSeq(string(data), "\n") {
		// marker may appear mid-line, e.g. "**Updated:** ... | **Version:** 0.7.0"
		if _, rest, ok := strings.Cut(line, marker); ok {
			want := strings.TrimSpace(rest)
			if want != Version {
				t.Errorf("utils.Version = %q, but FEATURES.md declares %q", Version, want)
			}

			return
		}
	}

	t.Fatalf("no **Version:** marker found in FEATURES.md")
}

// TestVersionMatchesReadmeBadge ensures the hand-edited README version badge
// stays in sync with utils.Version. The badge URL embeds the version as a
// literal (img.shields.io cannot read Go files), so it needs its own drift
// guard alongside the CHANGELOG and FEATURES checks.
func TestVersionMatchesReadmeBadge(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("../README.md")
	if err != nil {
		t.Fatalf("README.md not found: %v", err)
	}

	// Badge line looks like:
	// [![Version](https://img.shields.io/badge/version-v1.16.0-blue?style=flat-square)]
	const marker = "img.shields.io/badge/version-v"
	for line := range strings.SplitSeq(string(data), "\n") {
		if _, rest, ok := strings.Cut(line, marker); ok {
			// rest starts with "v1.16.0-blue?style=..."
			version, _, hasSuffix := strings.Cut(rest, "-")
			if !hasSuffix {
				version = rest
			}

			version = strings.TrimPrefix(version, "v")
			if version != Version {
				t.Errorf("utils.Version = %q, but README.md badge shows %q", Version, version)
			}

			return
		}
	}

	t.Fatalf("no version badge (marker %q) found in README.md", marker)
}
