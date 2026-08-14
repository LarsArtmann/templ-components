#!/usr/bin/env bash
# Rehearse the full CI pipeline locally before pushing.
#
# Mirrors .github/workflows/ci.yaml job-for-job so a red CI never surprises
# you after the push:
#   1. Fast guards (lint-config, templ-sync, version-sync, module-sync, DAG)
#   2. Lint          (golangci-lint via `nix run .#lint`, actionlint if present)
#   3. Build & Test  (generate, tidy-drift, vet, build, race tests + coverage
#                     gate, per-module isolation, visualtest compile, examples)
#   4. CSS Freshness (recompile via `nix run .#css`, diff against committed)
#   5. Visual        (nix run .#visual; fails on skip or failure)
#
# Requires: nix (flake), templ on PATH (nix develop), Go 1.26.
# Usage: scripts/verify-local.sh

set -euo pipefail

REPO_ROOT=$(cd "$(dirname "$0")/.." && pwd)
cd "$REPO_ROOT"

export GOEXPERIMENT=jsonv2

step() { printf '\n\033[1m==> %s\033[0m\n' "$1"; }
fail() {
	printf '\033[31mFAILED: %s\033[0m\n' "$1" >&2
	exit 1
}

# --- 1. Fast guards (mirrors CI "Lint" job prologue) -----------------------

step "Fast guards"
scripts/check-lint-config.sh
scripts/check-templ-sync.sh
scripts/check-version-sync.sh
scripts/check-module-sync.sh
scripts/check-module-layers.sh

# --- 2. Lint ----------------------------------------------------------------

step "Lint (golangci-lint)"
nix run .#lint || fail "golangci-lint"

if command -v actionlint >/dev/null 2>&1; then
	step "Lint (actionlint)"
	actionlint .github/workflows/ci.yaml .github/workflows/website.yml || fail "actionlint"
else
	echo "actionlint not on PATH; skipping (CI will run it)."
fi

# --- 3. Build & Test ----------------------------------------------------------

step "templ generate"
command -v templ >/dev/null 2>&1 ||
	fail "templ not on PATH; enter the dev shell with: nix develop"
tree_before=$(git status --porcelain)
templ generate ./...

step "Verify generated *_templ.go are tracked"
# Unlike CI (clean checkout), a local run starts with the developer's own
# uncommitted changes. Only generate/tidy-introduced drift should fail, so
# compare the tree state before vs after instead of diffing against HEAD.
UNTRACKED=$(git status --porcelain -- '*_templ.go' | grep '^??' || true)
if [ -n "$UNTRACKED" ]; then
	echo "$UNTRACKED"
	fail "untracked *_templ.go files; git add them (library consumers need them committed)"
fi

step "Go mod tidy (all modules)"
go mod tidy
for mod in utils icons errorpage charts/echarts datastar htmx; do
	(cd "$mod" && go mod tidy)
done

step "Verify no generate/tidy drift"
if [ "$tree_before" != "$(git status --porcelain)" ]; then
	diff <(echo "$tree_before") <(git status --porcelain) || true
	fail "generate/tidy modified the tree (see + lines above)"
fi

step "Go vet"
go vet ./...

step "Build"
go build ./...

step "Test (race + coverage)"
go list ./... | grep -v examples | xargs go test -race -coverprofile=coverage.out

step "Per-module isolation tests (GOWORK=off)"
for mod in utils icons errorpage charts/echarts datastar htmx; do
	echo "  ==> $mod"
	(cd "$mod" && GOWORK=off go test -race -count=1 ./...)
done

step "Compile visualtest module"
(cd visualtest && GOWORK=off GOEXPERIMENT=jsonv2 go test -count=1 ./...)

step "Docs-health drift guard"
(cd utils && go test ./... -run TestDocsCountDrift)

step "Coverage threshold (70%)"
COVERAGE=$(go tool cover -func=coverage.out | awk '/^total:/ { sub("%", "", $3); print $3 }')
echo "Total coverage: ${COVERAGE}%"
awk -v c="$COVERAGE" 'BEGIN { exit (c < 70) ? 1 : 0 }' ||
	fail "coverage ${COVERAGE}% is below 70% threshold"

step "Build examples"
go build ./examples/...

# --- 4. CSS Freshness ---------------------------------------------------------

step "CSS Freshness"
cp examples/demo/static/app.css /tmp/committed-app.css
nix run .#css
if ! diff -q /tmp/committed-app.css examples/demo/static/app.css >/dev/null; then
	diff /tmp/committed-app.css examples/demo/static/app.css | head -30
	fail "committed CSS is stale; recompile with: nix run .#css"
fi
echo "CSS is fresh."

# --- 5. Visual Regression -----------------------------------------------------

step "Visual Regression (nix run .#visual)"
set +e
OUTPUT=$(nix run .#visual 2>&1)
STATUS=$?
set -e
echo "$OUTPUT"
[ "$STATUS" -eq 0 ] || fail "visual tests exited $STATUS (see log above)"
if echo "$OUTPUT" | grep -qiE '(SKIP|no Chromium binary found)'; then
	fail "visual tests were skipped; Chromium must be available"
fi

printf '\n\033[32mAll CI stages passed locally. Ready to push.\033[0m\n'
