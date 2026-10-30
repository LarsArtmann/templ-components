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
#   scripts/ci-repro.sh --lint       # also run the Lint job (guards + actionlint + golangci-lint per module)
#   scripts/ci-repro.sh --css        # also run the CSS Freshness job (needs Nix)
#                                    # WATCH (#120): the 2026-08-31 local-vs-CI CSS discrepancy was
#                                    # never root-caused (post-fabd1fb unreproducible; nix run .#css
#                                    # byte-stable, SHA-256 verified twice). If CI's CSS Freshness
#                                    # fires where local --css was green, REOPEN #120 with both logs.
#                                    # 2026-10-14 (NOT #120): CI fired where the ritual had never run
#                                    # --css at all; new .templ classes staled app.css. The canonical
#                                    # ritual command now includes --css.
#   scripts/ci-repro.sh --visual     # also run the Visual Regression job (needs Nix + Chromium)
#   scripts/ci-repro.sh --website     # also run the Website job's local lane:
#                                    # website tests + lint + site build + tidy-check
#                                    # (the tidy-check is the deterministic undo of the
#                                    # daemon's go.mod pin flips — 3 incidents in 4 days)
#   scripts/ci-repro.sh --vuln       # also run the vulnerability gates: govulncheck (all Go modules)
#                                    # + pnpm audit --prod (website/; needs real node — the bun shim breaks pnpm,
#                                    # so run under `nix develop` or with the PATH workaround)
#   scripts/ci-repro.sh --tidy       # run the tidy-probe (Tidy Probe workflow): GOWORK=off tidy
#                                    # per module + fail on go.mod/go.sum drift (needs propagated tags)
#   scripts/ci-repro.sh --html       # also run the W3C HTML validation gate (html5validator;
#                                    # needs Nix: `nix shell nixpkgs#html5validator` is wrapped for you,
#                                    # or set VNU_JAR like CI's html-validation job does)
#   scripts/ci-repro.sh --flake      # also run `nix flake check` (needs Nix; loud skip when absent)
#   scripts/ci-repro.sh --cold       # use a throwaway GOCACHE (simulates a cache-less runner;
#                                    # module downloads still hit the local GOMODCACHE)
#
# Exit code is CI's exit code: 0 means "this push should be green".

set -euo pipefail

RUN_LINT=0
RUN_CSS=0
RUN_VISUAL=0
RUN_VULN=0
RUN_TIDY=0
RUN_WEBSITE=0
RUN_HTML=0
RUN_FLAKE=0
COLD=0
for arg in "$@"; do
	case "$arg" in
	--lint) RUN_LINT=1 ;;
	--css) RUN_CSS=1 ;;
	--visual) RUN_VISUAL=1 ;;
	--vuln) RUN_VULN=1 ;;
	--tidy) RUN_TIDY=1 ;;
	--website) RUN_WEBSITE=1 ;;
	--html) RUN_HTML=1 ;;
	--flake) RUN_FLAKE=1 ;;
	--quiet-diff) QUIET_DIFF=1 ;;
	--cold) COLD=1 ;;
	*)
		echo "Unknown flag: $arg" >&2
		echo "Usage: $0 [--lint] [--css] [--visual] [--vuln] [--tidy] [--website] [--html] [--flake] [--quiet-diff] [--cold]" >&2
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

step() {
	echo ""
	echo "==> $1"
}
MODULES="utils icons errorpage charts/echarts datastar htmx"

step "Generate templ files (pinned templ via PATH — use nix develop for the pinned binary)"
find . -name '*_templ.go' -print0 | xargs -0 rm -f
templ generate
# FileName-form normalization: a root invocation encodes repo-root-relative
# paths into EVERY workspace module's generated files, while the committed
# canonical form is module-relative (BuildFlow's per-module steppers). The
# form is error-metadata only, but the drift gate below compares bytes, so
# regenerate each published sub-module + website from its own directory —
# the same shape CI's ephemeral checkout never needs because it is discarded.
for mod in $MODULES website; do
	(cd "$mod" && templ generate ./...)
done

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
# --quiet-diff: show only that drift EXISTS, not the whole diff (useful for
# scheduled/CI mirror runs where the log is noise). Default prints the diff —
# it is the fastest way to see what a daemon commit missed.
#
# FileName-form tolerance: a generated file's templ.Error FileName strings
# encode the INVOKING DIRECTORY (repo root vs the file's module), and this
# script must work from a tree committed in either form (BuildFlow's
# per-module steppers commit module-relative; a root `templ generate` — what
# CI runs — writes repo-root-relative). FileName strings are error metadata
# only, so both sides are normalized to a placeholder before comparing; any
# OTHER line still counts as drift.
normalize_form() {
	sed 's/FileName: `[^`]*`/FileName: NORM/'
}
drift=0
non_templ_drift=$(git diff --name-only -- . ':(exclude)*_templ.go')
if [ -n "$non_templ_drift" ]; then
	git diff --exit-code -- $non_templ_drift || drift=1
fi
while IFS= read -r gen_file; do
	if ! diff <(git show ":$gen_file" | normalize_form) <(normalize_form <"$gen_file") >/dev/null; then
		echo "REAL content drift: $gen_file"
		git diff -- "$gen_file" | head -40
		drift=1
	fi
done < <(git diff --name-only -- '*_templ.go')
if [ "$drift" = "1" ]; then
	if [ "${QUIET_DIFF:-0}" = "1" ]; then
		echo "ERROR: uncommitted changes present (run without --quiet-diff for the diff)" >&2
	fi
	exit 1
fi

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
# awk, not bc: bc is missing on some dev machines and the check would
# silently pass (CI's ubuntu runner has bc, so this stays CI-equivalent).
if awk "BEGIN{exit !($COVERAGE < 70)}"; then
	echo "Coverage ${COVERAGE}% is below 70% threshold" >&2
	exit 1
fi

step "Coverage gate self-test (synthetic profiles, no compile)"
scripts/test-coverage-library-gate.sh

step "Library coverage (hand-written, generated templ excluded)"
scripts/coverage-library.sh --gate 94

step "Build examples"
go build ./examples/...

if [ "$RUN_LINT" = "1" ]; then
	step "Lint guards"
	scripts/check-lint-config.sh
	scripts/check-templ-sync.sh
	scripts/check-version-sync.sh
	scripts/check-module-sync.sh
	scripts/check-module-layers.sh
	scripts/check-lint-modules.sh
	scripts/check-dupl-pin.sh
	scripts/check-tc-sources-sync.sh
	scripts/test-release-assertions.sh
	scripts/test-settle-staged-go-format.sh

	step "Changelog warmth guard (#133 — local parity with CI's changelog-guard job)"
	# CI runs scripts/check-changelog-guard.sh over the PR file list; locally
	# the same diff is origin/master...HEAD (empty on master = vacuously green).
	CHANGELOG_GUARD_FILES=$(git diff --name-only origin/master...HEAD 2>/dev/null || true)
	if [ -n "$CHANGELOG_GUARD_FILES" ]; then
		printf '%s\n' "$CHANGELOG_GUARD_FILES" | scripts/check-changelog-guard.sh
	else
		echo "  (no branch diff vs origin/master — guard vacuously green)"
	fi

	step "Actionlint (GitHub Actions workflows — same as CI's Actionlint step)"
	if command -v actionlint >/dev/null 2>&1; then
		actionlint
	else
		echo "actionlint not found — run under nix develop (nix run .#lint includes it) or:" >&2
		echo "  nix shell nixpkgs#actionlint -c $0 --lint" >&2
		exit 1
	fi

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

	# Mirror of ci.yaml's Lint job: visualtest has its own go.mod (its tests
	# run in the Visual Regression job, but the LINT lane lives with the
	# others). Without this lane a visualtest-only lint finding passes local
	# CI reproduction and fails the real CI Lint job — exactly how the
	# gocognit finding on runKanbanContractProbes escaped on 2026-09-17.
	echo "---- visualtest (tests run in the Visual Regression job)"
	(cd visualtest && golangci-lint run --timeout=5m ./...)
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

if [ "$RUN_WEBSITE" = "1" ]; then
	step "Website module (mirrors website.yml build job, #241)"
	(
		cd website
		GOWORK=off go build ./...
		GOWORK=off go test ./... -count=1
		GOWORK=off golangci-lint run ./...
		GOEXPERIMENT=jsonv2 go run ./cmd/site --skip-stars --out dist --repo-root ..
		rm -rf dist
	)
	step "Website tidy-check (go mod tidy must be a no-op — kills the pin-flip class)"
	(
		cd website
		GOWORK=off go mod tidy
		if ! git diff --exit-code -- go.mod go.sum >/dev/null; then
			echo "ERROR: website/go.mod drifted after tidy — the daemon flipped the pin again;" >&2
			echo "commit the tidied go.mod (tidy is the deterministic undo, #241)." >&2
			git diff -- go.mod go.sum | head -20 >&2
			exit 1
		fi
	)
	echo "website lane green."
fi

if [ "$RUN_TIDY" = "1" ]; then
	step "tidy probe (Tidy Probe workflow): GOWORK=off tidy per module + drift check"
	failed=0
	for mod in . $MODULES visualtest; do
		echo "---- $mod"
		if ! (cd "$mod" && GOWORK=off go mod tidy); then
			echo "ERROR: $mod tidy failed — release tags may not have propagated yet" >&2
			failed=1
			continue
		fi
		if ! git diff --exit-code -- "$mod/go.mod" "$mod/go.sum" >/dev/null; then
			echo "ERROR: $mod go.mod/go.sum stale after tidy — commit the refresh" >&2
			failed=1
		fi
	done
	if [ "$failed" -ne 0 ]; then
		exit 1
	fi
	echo "tidy probe clean."
fi

if [ "$RUN_VULN" = "1" ]; then
	step "govulncheck (root + sub-modules; symbol-level findings fail the gate)"
	command -v govulncheck >/dev/null || {
		echo "ERROR: govulncheck not found — run under nix develop or:" >&2
		echo "  nix shell nixpkgs#govulncheck -c $0 --vuln" >&2
		exit 1
	}
	govulncheck ./...
	for mod in $MODULES visualtest; do
		echo "---- $mod"
		(cd "$mod" && GOWORK=off govulncheck ./...)
	done

	step "website build (Go SSG + Tailwind; mirrors website.yml build job)"
	(
		cd website
		GOWORK=off go build ./...
		GOWORK=off golangci-lint run ./...
		GOEXPERIMENT=jsonv2 go run ./cmd/site --skip-stars --out dist --repo-root ..
		rm -rf dist
	)
fi

# ADVISORY lane (never fails the run): the art-dupl clone gate is binary-pinned
# (fork at ~/projects/art-dupl — see docs/adr/0009-accepted-clones.md), so a
# bare "exit 1" here could be version skew, not real clones. Running it beside
# every verdict makes drift visible without blocking the ritual.
step "art-dupl clone gate (ADVISORY — reported beside the verdict, never fails the run)"
DUPL_BIN="${ART_DUPL_BIN:-}"
if [ -z "$DUPL_BIN" ] && [ -x /tmp/art-dupl-fork ]; then
	DUPL_BIN=/tmp/art-dupl-fork
elif [ -z "$DUPL_BIN" ]; then
	DUPL_BIN="$(command -v art-dupl 2>/dev/null || true)"
fi
if [ -z "$DUPL_BIN" ]; then
	DUPL_VER="dupl-skip"
	echo "SKIP: no art-dupl binary (set ART_DUPL_BIN or build the fork — see docs/adr/0009-accepted-clones.md)"
else
	DUPL_VER="$("$DUPL_BIN" --version 2>&1 | head -1)"
	echo "binary: $DUPL_BIN ($DUPL_VER)"
	DUPL_PIN="$(sed -n 's/^PINNED_REV="\([^"]*\)".*/\1/p' scripts/pinned-art-dupl.sh 2>/dev/null || true)"
	DUPL_RECORDED="$(sed -n 's/.*"recordedAt": "\([^"]*\)".*/\1/p' .art-dupl-baseline.json 2>/dev/null || true)"
	echo "baseline recordedAt: ${DUPL_RECORDED:-UNKNOWN} (pinned rev: ${DUPL_PIN:-UNKNOWN})"
	if [ -n "$DUPL_PIN" ] && [ -n "$DUPL_VER" ] && ! echo "$DUPL_VER" | grep -q "${DUPL_PIN:0:12}"; then
		echo "WARNING: resolved binary is NOT the pinned rev ($DUPL_PIN) — its hashes may match nothing in the baseline (ADR-0009); build via scripts/pinned-art-dupl.sh" >&2
	fi
	set +e
	DUPL_OUT="$("$DUPL_BIN" check -c .art-dupl.json -t 1 --type-aware 2>&1)"
	DUPL_EXIT=$?
	set -e
	printf '%s\n' "$DUPL_OUT" | tail -3
	if [ "$DUPL_EXIT" -eq 0 ]; then
		echo "art-dupl: GREEN — 0 clone groups outside baseline"
	else
		echo "art-dupl: ADVISORY-FAIL (exit $DUPL_EXIT) — new clone groups outside baseline; classify per ADR-0009 (extract or accept+rebaseline)"
	fi
fi

if [ "$RUN_HTML" = "1" ]; then
	step "W3C HTML validation (golden corpus; mirrors ci.yaml's html-validation job)"
	if [ -n "${VNU_JAR:-}" ]; then
		scripts/check-html-valid.sh
	elif command -v html5validator >/dev/null 2>&1; then
		scripts/check-html-valid.sh
	elif command -v nix >/dev/null 2>&1; then
		nix shell nixpkgs#html5validator -c scripts/check-html-valid.sh
	else
		echo "ERROR: html5validator not found and Nix unavailable." >&2
		echo "  Run under nix develop, or: nix shell nixpkgs#html5validator -c $0 --html" >&2
		echo "  (loud skip would hide a broken page — this lane FAILS instead)" >&2
		exit 1
	fi
fi

if [ "$RUN_FLAKE" = "1" ]; then
	step "nix flake check (treefmt + flake-level checks)"
	if command -v nix >/dev/null 2>&1; then
		nix flake check
	else
		echo "SKIP: nix not on PATH — flake checks did NOT run (TODO: rerun with nix)" >&2
	fi
fi

echo ""
echo ""
echo "ALL STEPS PASSED — working tree matches CI's expectations."
echo "VERDICT: PASS (exit 0) — $(date '+%Y-%m-%d %H:%M:%S %Z') [dupl: $DUPL_VER]"
exit 0
