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
		"trap release_cleanup EXIT",
		"must install an EXIT-trap rollback so a failed verify restores version files",
	)
	checkContains(
		t,
		script,
		"release_cleanup()",
		"the rollback hook must exist (renamed from release_rollback when the v1.9.0 fix merged both EXIT traps into one)",
	)
	checkContains(
		t,
		script,
		"RELEASE_COMMITTED=1",
		"must flip RELEASE_COMMITTED after the commit so a later tag failure keeps the commit",
	)
	checkAbsent(t, script, "git checkout --", "must use git restore, not git checkout (AGENTS.md ban)")

	// Invariant added 2026-08-16 after v1.8.3 shipped as a root-only release:
	// the sub-module set (what gets bumped, replace-stripped, and tagged) must
	// be DERIVED from the root go.mod's replace directives, never hardcoded.
	// The hardcoded list had silently drifted (missing htmx/ and datastar/),
	// so those modules never received release tags and every consumer pinning
	// them at the new version failed with "unknown revision".
	checkContains(
		t,
		script,
		"SUBMODULE_PATHS=",
		"must derive the sub-module set from the root go.mod replace directives",
	)
	checkAbsent(
		t,
		script,
		"for submod in utils icons errorpage charts/echarts",
		"must not hardcode the sub-module tag list (v1.8.3 root-only release root cause)",
	)
	checkContains(t, script, "for submod in $SUBMODULE_PATHS", "tagging loop must use the derived set")
	checkContains(t, script, "for modfile in $MODFILES", "go.mod loops must use the derived set")
}

// TestCheckReleaseTagsScriptInvariants guards scripts/check-release-tags.sh:
// the lockstep-tag guard must exist, derive the same sub-module set as the
// release script, and verify tag-to-commit equality (not just existence).
func TestCheckReleaseTagsScriptInvariants(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("../scripts/check-release-tags.sh")
	if err != nil {
		t.Skipf("scripts/check-release-tags.sh not found (running outside repo root?): %v", err)
	}

	script := string(data)

	checkContains(t, script, "SUBMODULE_PATHS=", "guard must derive sub-modules like release.sh does")
	checkContains(t, script, "^{commit}", "guard must resolve tags to commits")
	checkContains(t, script, "MISSING:", "guard must name missing sub-module tags")
	checkContains(t, script, "DIVERGED:", "guard must catch tags pointing at the wrong commit")
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
