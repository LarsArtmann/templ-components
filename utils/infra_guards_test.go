package utils

import (
	"os"
	"strings"
	"testing"
)

// TestEnvrcConsistency verifies the committed .envrc carries the
// GOEXPERIMENT flag every tool needs (go, gopls, BuildFlow, IDE) — not just
// the `nix develop` shell. .envrc is tracked (no secrets) so direnv activates
// it on every clone. GOWORK is no longer set to off: the go.work file lists
// all sub-modules for multi-module workspace development.
func TestEnvrcConsistency(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("../.envrc")
	if err != nil {
		t.Fatalf("read ../.envrc: %v\nIf .envrc was deliberately removed, delete this test too.", err)
	}

	src := string(data)

	for _, want := range []string{
		"GOEXPERIMENT=jsonv2",
	} {
		if !strings.Contains(src, want) {
			t.Errorf(".envrc is missing %q — every tool outside `nix develop` will misbuild the module.", want)
		}
	}

	// .envrc must never carry secrets or machine-specific paths.
	dangerous := []string{
		"source ", "source\t", "~/.secrets", "/home/", "/Users/",
		"export SECRET", "export TOKEN", "export PASSWORD", "export API_KEY",
	}
	for _, bad := range dangerous {
		if strings.Contains(src, bad) {
			t.Errorf(".envrc contains %q — .envrc is committed and must stay secret-free and machine-independent.", bad)
		}
	}
}

// TestPreCommitHookInstallsGuard verifies the local pre-commit hook still wires
// the lint-config guard BEFORE BuildFlow. BuildFlow's `precommit install`
// regenerates this hook and has, in past sessions, dropped the manual guard —
// silently re-opening the recurring `.golangci.yml` disabled-linter regression
// (5 occurrences). The hook lives in .git/ (not committed), so it is absent on
// a fresh CI clone; the test skips there and only asserts locally where the
// hook is installed.
func TestPreCommitHookInstallsGuard(t *testing.T) {
	t.Parallel()

	hookPath := "../.git/hooks/pre-commit"

	data, err := os.ReadFile(hookPath)
	if err != nil {
		t.Skipf("pre-commit hook not installed at %s (skipping; CI clones have no hooks): %v", hookPath, err)

		return
	}

	src := string(data)

	// The guard must run BEFORE BuildFlow so the 60s BuildFlow budget never
	// masks a disabled-linter re-enablement.
	if !strings.Contains(src, "check-lint-config.sh") {
		t.Errorf(
			"pre-commit hook (%s) no longer calls check-lint-config.sh — "+
				"re-add the pre-BuildFlow guard (see scripts/check-lint-config.sh and .git/hooks/pre-commit).",
			hookPath,
		)
	}

	// The guard must precede the BuildFlow invocation, not follow it. Match the
	// actual execution line (buildflow --build-mode), not the word "buildflow"
	// which also appears in the header comment.
	guardIdx := strings.Index(src, "check-lint-config.sh")
	buildFlowIdx := strings.Index(src, "buildflow --build-mode")

	if guardIdx >= 0 && buildFlowIdx >= 0 && guardIdx > buildFlowIdx {
		t.Errorf(
			"pre-commit hook runs check-lint-config.sh AFTER buildflow — " +
				"the guard must run BEFORE BuildFlow so it is not masked by the 60s budget.",
		)
	}
}
