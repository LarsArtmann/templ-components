package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strconv"
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

	// The replace-directives tripwire must be wired too: the v1.17.0 release
	// race left all 5 dependent sub-modules without replace blocks because the
	// release script's re-add step vanished under the daemon (2026-09-13).
	// Same pre-BuildFlow requirement as the lint-config guard.
	replaceGuardIdx := strings.Index(src, "check-replace-directives.sh")
	if replaceGuardIdx < 0 {
		t.Errorf(
			"pre-commit hook (%s) no longer calls check-replace-directives.sh — "+
				"re-add the pre-BuildFlow guard (see scripts/check-replace-directives.sh).",
			hookPath,
		)
	} else if buildFlowIdx >= 0 && replaceGuardIdx > buildFlowIdx {
		t.Errorf(
			"pre-commit hook runs check-replace-directives.sh AFTER buildflow — " +
				"the guard must run BEFORE BuildFlow so it is not masked by the 60s budget.",
		)
	}
}

var goDirectiveRe = regexp.MustCompile(`^go (\d+\.\d+(?:\.\d+)?)$`)

// goDirectiveFrom returns the `go X.Y[.Z]` language version of a go.mod or
// go.work source, or "" when absent.
func goDirectiveFrom(src string) string {
	for line := range strings.SplitSeq(src, "\n") {
		if match := goDirectiveRe.FindStringSubmatch(strings.TrimSpace(line)); match != nil {
			return match[1]
		}
	}

	return ""
}

// TestGoDirectiveSkew pins the invariant that every workspace module's go.mod
// `go` directive is <= the go.work `go` directive. On 2026-09-17 the
// auto-commit daemon bumped the ROOT go.mod to a newer Go than the pinned
// toolchain/go.work, and every manual `git commit` failed in the pre-commit
// hook with a cryptic "module . requires go >= X but go.work has go Y"
// go-tool-run error while the daemon itself bypassed the hook. The daemon does
// not run tests. The module set is derived from go.work's `use` lines, so newly
// added modules are covered automatically.
//
// Where it runs: go.work is gitignored (local dev only), so CI checkouts have
// none — the guard SKIPS there (it cannot compare against a go.work it cannot
// read). CI still catches a daemon-bumped directive on its own: the per-module
// build dies with "go.mod requires go >= X" against the pinned runner
// toolchain. The guard's added value is local/dev-shell runs, where it names
// EVERY offending module precisely instead of failing on the first one.
func TestGoDirectiveSkew(t *testing.T) {
	t.Parallel()

	workSrc, err := os.ReadFile("../go.work")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			t.Skip("../go.work absent (CI checkout: go.work is gitignored) — skew guard is local-dev only; CI's own build fails on a bumped directive.")
		}

		t.Fatalf("read ../go.work: %v", err)
	}

	var workGo string

	var useDirs []string

	for line := range strings.SplitSeq(string(workSrc), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "go ") {
			workGo = goDirectiveFrom(line)
		} else if dir, ok := strings.CutPrefix(line, "use "); ok {
			useDirs = append(useDirs, strings.TrimSpace(dir))
		}
	}

	if workGo == "" {
		t.Fatalf("go.work has no `go <version>` directive — cannot verify toolchain skew.")
	}

	for _, dir := range useDirs {
		modSrc, err := os.ReadFile(fmt.Sprintf("../%s/go.mod", strings.TrimPrefix(dir, "./")))
		if err != nil {
			t.Fatalf("read %s/go.mod (listed in go.work `use`): %v", dir, err)
		}

		modGo := goDirectiveFrom(string(modSrc))
		if modGo == "" {
			t.Errorf("%s/go.mod has no `go <version>` directive.", dir)

			continue
		}

		if compareGoVersions(modGo, workGo) > 0 {
			t.Errorf(
				"%s/go.mod requires go %s but go.work pins go %s — bump go.work AND the pinned "+
					"toolchain together, or revert the module directive. This exact skew broke every "+
					"manual commit on 2026-09-17.",
				dir, modGo, workGo,
			)
		}
	}
}

// compareGoVersions compares two Go language-version strings ("1.26",
// "1.26.7"). Returns -1, 0, or 1. A missing patch component counts as 0.
func compareGoVersions(a, b string) int {
	aParts, bParts := parseGoVersion(a), parseGoVersion(b)

	for i := range 3 {
		switch {
		case aParts[i] < bParts[i]:
			return -1
		case aParts[i] > bParts[i]:
			return 1
		}
	}

	return 0
}

func parseGoVersion(v string) [3]int {
	var parts [3]int

	for i, seg := range strings.SplitN(v, ".", 3) {
		if i == 3 {
			break
		}

		n, err := strconv.Atoi(seg)
		if err != nil {
			return parts
		}

		parts[i] = n
	}

	return parts
}
