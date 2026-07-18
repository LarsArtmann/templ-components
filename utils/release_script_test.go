package utils

import (
	"os"
	"strings"
	"testing"
)

// TestReleaseScriptInvariants guards scripts/release.sh against regressing the
// three defects fixed on 2026-07-18 (see docs/status/...v0.18.0-release-postmortem.md):
//  1. Commit body must not duplicate the one-line summary — use RELEASE_NOTES.
//  2. Model attribution must be parameterized — no hardcoded MiniMax-M3.
//  3. Release notes must come from --notes-file or CHANGELOG [Unreleased], not
//     a hostile stdin read loop with no editing or file input.
//
// This is a static-analysis drift guard: it reads the script as text and asserts
// the invariant patterns are present (and the regression patterns are absent).
func TestReleaseScriptInvariants(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("../scripts/release.sh")
	if err != nil {
		t.Skipf("scripts/release.sh not found (running outside repo root?): %v", err)
	}

	script := string(data)

	checkContains(t, script, `RELEASE_BODY="${RELEASE_NOTES}`, "commit body must use ${RELEASE_NOTES}")
	checkAbsent(t, script, `RELEASE_BODY="${RELEASE_SUMMARY}`, "body must not duplicate the subject summary")

	checkContains(t, script, "${CRUSH_MODEL", "model attribution must reference ${CRUSH_MODEL}")
	checkAbsent(t, script, "Crush:MiniMax-M3", "model must not be hardcoded to MiniMax-M3")

	checkContains(t, script, "--notes-file", "must support --notes-file FILE")
	checkContains(t, script, "[Unreleased]", "must extract notes from CHANGELOG [Unreleased] by default")
	checkAbsent(t, script, "while IFS= read -r line", "no hostile stdin read loop")

	// Invariants added 2026-07-18 after the v0.18.1 cut exposed three more
	// defects: (4) FEATURES.md version must move with utils.Version + CHANGELOG
	// (utils.TestVersionMatchesFeatures enforces it); (5) a failed verify must
	// not leave the tree dirty — an EXIT-trap rollback restores the version
	// files; (6) AGENTS.md bans `git checkout --` in favor of `git restore`.
	checkContains(
		t,
		script,
		"Bumped FEATURES.md version to",
		"must bump FEATURES.md version alongside utils.Version and CHANGELOG",
	)
	checkContains(
		t,
		script,
		"trap release_rollback EXIT",
		"must install an EXIT-trap rollback so a failed verify restores version files",
	)
	checkContains(
		t,
		script,
		"RELEASE_COMMITTED=1",
		"must flip RELEASE_COMMITTED after the commit so a later tag failure keeps the commit",
	)
	checkAbsent(t, script, "git checkout --", "must use git restore, not git checkout (AGENTS.md ban)")
}

func checkContains(t *testing.T, script, needle, msg string) {
	t.Helper()

	if !strings.Contains(script, needle) {
		t.Errorf("release.sh invariant failed: %s (missing %q)", msg, needle)
	}
}

func checkAbsent(t *testing.T, script, needle, msg string) {
	t.Helper()

	if strings.Contains(script, needle) {
		t.Errorf("release.sh invariant failed: %s (found forbidden %q)", msg, needle)
	}
}
