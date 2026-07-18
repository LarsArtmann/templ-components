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
