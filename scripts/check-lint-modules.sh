#!/usr/bin/env bash
# Fast (<100ms) pre-commit guard: the golangci-lint module/package sets must be
# IDENTICAL across the three places that define them:
#
#   1. .github/workflows/ci.yaml  (Lint job — the gate that actually runs on PRs)
#   2. scripts/ci-repro.sh        (local pre-push reproduction, --lint lane)
#   3. scripts/pre-commit.sh      (full manual pre-push verify)
#
# WHY: golangci-lint does not support go.work, so each file hand-maintains the
# module list. When a module is added to one file but not the others, the new
# module's lint findings pass locally and surface only in CI (or never) —
# exactly how the 2026-09-17 golines drift stayed invisible for 3 runs and how
# visualtest's 21 findings once escaped local CI reproduction.
#
# Checked sets:
#   - root-module package list  (./display/... ./forms/... etc.)
#   - sub-module list           (utils icons errorpage charts/echarts datastar htmx)
#   - visualtest lint lane      (must exist in ci.yaml + ci-repro.sh)
#
# Usage:
#   scripts/check-lint-modules.sh
set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
cd "$REPO_ROOT"

CI_YAML=".github/workflows/ci.yaml"
CI_REPRO="scripts/ci-repro.sh"
PRE_COMMIT="scripts/pre-commit.sh"

DRIFT=0

fail() {
	echo "MISMATCH: $1" >&2
	DRIFT=1
}

# --- Sub-module set ---------------------------------------------------------
# Canonical source: ci-repro.sh's MODULES= line (single place for the loop).
# Effective lint set = MODULES + visualtest (ci-repro.sh lints visualtest
# right after the MODULES loop; ci.yaml lists it as its own (cd ...) line).
REPRO_MODULES="$(sed -n 's/^MODULES="\(.*\)"$/\1/p' "$CI_REPRO" | head -1)"
if [ -z "$REPRO_MODULES" ]; then
	fail "cannot parse MODULES= from $CI_REPRO"
	REPRO_MODULES="(unparsed)"
fi

# ci.yaml module lint lines: (cd <mod> && golangci-lint run
YAML_MODULES="$(grep -oE '\(cd [a-z/.-]+ && golangci-lint run' "$CI_YAML" |
	sed 's/(cd \(.*\) && golangci-lint run/\1/' | sort | tr '\n' ' ' | sed 's/ $//')"

REPRO_LINT_MODULES="$(for mod in $REPRO_MODULES visualtest; do echo "$mod"; done | sort | tr '\n' ' ' | sed 's/ $//')"

if [ "$YAML_MODULES" != "$REPRO_LINT_MODULES" ]; then
	fail "module sets diverge
    ci.yaml : $YAML_MODULES
    ci-repro: $REPRO_LINT_MODULES
  (also update scripts/pre-commit.sh — both files must carry the same set)"
fi

# visualtest lane must exist in ci.yaml and ci-repro.sh (its tests live in the
# Visual job, but the LINT lane must be in both Lint jobs).
grep -q '(cd visualtest && golangci-lint run' "$CI_YAML" ||
	fail "visualtest lint lane missing from $CI_YAML"
grep -q '(cd visualtest && golangci-lint run' "$CI_REPRO" ||
	fail "visualtest lint lane missing from $CI_REPRO"

# --- Root-module package set ------------------------------------------------
# Extracts ./pkg/... tokens ONLY from the root `golangci-lint run` command
# (the trigger line plus its backslash continuations). A whole-file grep would
# swallow unrelated tokens like `go build ./examples/...` in the Build job.
extract_root_pkgs() {
	awk '
		/golangci-lint run/ && !/\(cd/ { in_cmd = 1 }
		in_cmd {
			line = $0
			while (match(line, /\.\/[a-z-]+\/\.\.\./)) {
				print substr(line, RSTART, RLENGTH)
				line = substr(line, RSTART + RLENGTH)
			}
			if (line !~ /\\[ \t]*$/) in_cmd = 0
		}
	' "$1" | sort -u | tr '\n' ' ' | sed 's/ $//'
}

YAML_ROOT="$(extract_root_pkgs "$CI_YAML")"
REPRO_ROOT="$(extract_root_pkgs "$CI_REPRO")"
PRECOMMIT_ROOT="$(extract_root_pkgs "$PRE_COMMIT")"

if [ "$YAML_ROOT" != "$REPRO_ROOT" ]; then
	fail "root package sets diverge
    ci.yaml : $YAML_ROOT
    ci-repro: $REPRO_ROOT"
fi

if [ "$YAML_ROOT" != "$PRECOMMIT_ROOT" ]; then
	fail "root package sets diverge
    ci.yaml : $YAML_ROOT
    pre-commit: $PRECOMMIT_ROOT"
fi

# Every ci.yaml root package must also be listed in pre-commit.sh and vice
# versa is implied by the equality above.

if [ "$DRIFT" -ne 0 ]; then
	echo "" >&2
	echo "BLOCKED: golangci-lint module/package sets are out of sync across" >&2
	echo "ci.yaml / scripts/ci-repro.sh / scripts/pre-commit.sh." >&2
	echo "Update ALL THREE files with the same set (a module linted in only one" >&2
	echo "place is how lint regressions reach master invisible)." >&2
	exit 1
fi

echo "lint module sets in sync (root + sub-modules + visualtest lane)."
exit 0
