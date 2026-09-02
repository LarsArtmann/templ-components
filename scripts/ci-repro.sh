#!/usr/bin/env bash
# scripts/ci-repro.sh — reproduce CI's exact step sequence locally, pre-push.
#
# WHY: `go test ./...` from the repo root (workspace mode) only runs the ROOT
# module's packages — `go list ./...` lists ZERO utils/icons/... packages even
# with go.work active (verified 2026-08-31). The only complete local test form
# is CI's per-module sequence, which this script mirrors step-for-step so a
# push is never the first time CI's order executes.
#
# Usage:
#   scripts/ci-repro.sh              # core: generate + tidy + verify + test (the Build & Test job)
#   scripts/ci-repro.sh --lint       # also run the Lint job (guards + golangci-lint per module)
#   scripts/ci-repro.sh --css        # also run the CSS Freshness job (needs Nix)
#   scripts/ci-repro.sh --visual     # also run the Visual Regression job (needs Nix + Chromium)
#   scripts/ci-repro.sh --cold       # use a throwaway GOCACHE (simulates a cache-less runner;
#                                    # module downloads still hit the local GOMODCACHE)
#
# Exit code is CI's exit code: 0 means "this push should be green".

set -euo pipefail

RUN_LINT=0
RUN_CSS=0
RUN_VISUAL=0
COLD=0
for arg in "$@"; do
	case "$arg" in
	--lint) RUN_LINT=1 ;;
	--css) RUN_CSS=1 ;;
	--visual) RUN_VISUAL=1 ;;
	--cold) COLD=1 ;;
	*)
		echo "Unknown flag: $arg" >&2
		echo "Usage: $0 [--lint] [--css] [--visual] [--cold]" >&2
		exit 2
		;;
	esac
done

export GOEXPERIMENT=jsonv2
cd "$(git rev-parse --show-toplevel 2>/dev/null || echo .)"

if [ "$COLD" = "1" ]; then
	COLD_CACHE="$(mktemp -d)/gocache"
	export GOCACHE="$COLD_CACHE"
	echo "==> Cold cache: $COLD_CACHE"
fi

step() { echo ""; echo "==> $1"; }
MODULES="utils icons errorpage charts/echarts datastar htmx"

step "Generate templ files (pinned templ via PATH — use nix develop for the pinned binary)"
find . -name '*_templ.go' -print0 | xargs -0 rm -f
templ generate

step "Verify all *_templ.go files are tracked"
MISSING=0
while IFS= read -r -d '' templ_file; do
	gen_file="${templ_file%.templ}_templ.go"
	if ! git ls-files --error-unmatch "$gen_file" >/dev/null 2>&1; then
		echo "ERROR: $gen_file is not tracked. Generated files MUST be committed." >&2
		MISSING=1
	fi
done < <(find . -name '*.templ' -not -path './vendor/*' -not -path './cmd/tc/_sources/*' -print0)
[ "$MISSING" = "0" ] || exit 1

step "Go mod tidy (all modules, incl. visualtest with local replaces)"
go mod tidy
for mod in $MODULES visualtest; do
	(cd "$mod" && GOWORK=off go mod tidy)
done

step "Verify no untracked changes (git diff --exit-code)"
git diff --exit-code

step "Go vet (root)"
go vet ./...

step "Build (root)"
go build ./...

step "Test (root, minus examples, race + coverage)"
go list ./... | grep -v examples | xargs go test -race -coverprofile=coverage.out -count=1

step "Per-module isolation tests (GOWORK=off, race)"
for mod in $MODULES; do
	echo "---- $mod"
	(cd "$mod" && GOWORK=off go test -race -count=1 ./...)
done

step "Compile visualtest module (GOWORK=off; tests skip without Chromium)"
(
	cd visualtest
	for attempt in 1 2 3; do
		if GOWORK=off go test -count=1 ./...; then
			break
		fi
		echo "warning: attempt ${attempt}/3 failed — retrying" >&2
		sleep $((attempt * 5))
		[ "$attempt" = "3" ] && exit 1
	done
)

step "Docs-health drift guard"
(cd utils && go test ./... -run TestDocsCountDrift -count=1)

step "Coverage threshold (>= 70%)"
COVERAGE="$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | tr -d '%')"
echo "Total coverage: ${COVERAGE}%"
if [ "$(echo "$COVERAGE < 70" | bc -l)" -eq 1 ]; then
	echo "Coverage ${COVERAGE}% is below 70% threshold" >&2
	exit 1
fi

step "Build examples"
go build ./examples/...

if [ "$RUN_LINT" = "1" ]; then
	step "Lint guards"
	scripts/check-lint-config.sh
	scripts/check-templ-sync.sh
	scripts/check-version-sync.sh
	scripts/check-module-sync.sh
	scripts/check-module-layers.sh

	step "golangci-lint (root module)"
	golangci-lint run --timeout=5m \
		./display/... ./feedback/... ./forms/... \
		./integration/... ./internal/... \
		./layout/... ./navigation/... ./recipes/... ./cmd/...

	step "golangci-lint (sub-modules)"
	for mod in $MODULES; do
		echo "---- $mod"
		(cd "$mod" && golangci-lint run --timeout=5m ./...)
	done
fi

if [ "$RUN_CSS" = "1" ]; then
	step "CSS Freshness (recompile + diff)"
	cp examples/demo/static/app.css /tmp/committed-app.css
	nix run .#css
	if ! diff -q /tmp/committed-app.css examples/demo/static/app.css >/dev/null; then
		echo "ERROR: Committed CSS is stale. Recompile with: nix run .#css" >&2
		diff /tmp/committed-app.css examples/demo/static/app.css | head -30 >&2
		exit 1
	fi
	echo "CSS is fresh."
fi

if [ "$RUN_VISUAL" = "1" ]; then
	step "Visual Regression (Nix Chromium; hard gate in CI)"
	nix run .#visual
fi

echo ""
echo "ALL STEPS PASSED — working tree matches CI's expectations."
