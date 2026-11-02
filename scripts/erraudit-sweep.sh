#!/usr/bin/env bash
# scripts/erraudit-sweep.sh — Go 1.26+ error-handling modernization sweep.
#
# The 2026-11-01 pass proved the repo 100% on errors.AsType with 0 findings;
# this script makes that verdict ONE repeatable command instead of a
# per-session archaeology dig. It runs, per go.mod module (the repo's
# GOWORK=off convention — golangci-lint does not support go.work either):
#
#   1. `erraudit fix` DRY-RUN (its default) — previews auto-fixable
#      errors.As migrations. Advisory only; never fails the sweep.
#   2. `erraudit lint --type-aware --type legacy_as` — THE GATE: any finding
#      exits non-zero. errors.As migrations are auto-fixable, so a finding
#      is never a judgment call.
#   3. `erraudit lint --type-aware --type legacy_is` — advisory lane:
#      errors.Is findings are REPORTED but never fail. Matching by value is
#      often CORRECT for sentinels (io.EOF, sql.ErrNoRows, context.Canceled);
#      migrating those to AsType by reflex is the documented cargo-cult trap.
#
# JS/TS findings (empty catch, discarded rejections, floating promises) are
# outside the --type-filtered lanes; run `erraudit lint ./...` bare ad hoc.
#
# Usage:
#   scripts/erraudit-sweep.sh          # dry-run preview + gate + advisory
#   scripts/erraudit-sweep.sh --fix    # APPLY auto-fixes (after a red gate)
#
# Exit codes: 0 = gate green; 1 = legacy_as finding(s); 2 = erraudit failed
# to run (missing binary, broken module load).

set -euo pipefail
cd "$(dirname "$0")/.."

FIX_MODE=false
[ "${1:-}" = "--fix" ] && FIX_MODE=true

command -v erraudit >/dev/null 2>&1 || {
	echo "FATAL: erraudit not on PATH" >&2
	exit 2
}
echo "erraudit $(erraudit version 2>/dev/null | head -1)"

export GOWORK=off GOEXPERIMENT=jsonv2

# Same module set as the per-module lint/test loops (AGENTS.md).
MODULES=("." "utils" "icons" "errorpage" "charts/echarts" "htmx" "datastar" "visualtest" "website")

GATE_RED=0

for mod in "${MODULES[@]}"; do
	label="${mod:-. (root)}"
	echo
	echo "=== module: $label"

	cd "$mod"

	fix_out="$(erraudit fix ./... -q 2>&1 || true)"
	if [ -n "$fix_out" ]; then
		echo "--- auto-fix preview (dry-run; rerun with --fix to apply):"
		echo "$fix_out"
	else
		echo "auto-fix preview: clean"
	fi

	if $FIX_MODE && [ -n "$fix_out" ]; then
		echo "--- applying fixes"
		erraudit fix ./... -w -q
	fi

	as_out="$(erraudit lint ./... --type-aware --type legacy_as -q 2>&1 || true)"
	as_count="$(grep -cE '\.go:[0-9]+:[0-9]+:' <<<"$as_out" || true)"
	if [ "$as_count" -gt 0 ]; then
		echo "GATE RED: $as_count errors.As finding(s):"
		echo "$as_out"
		GATE_RED=1
	else
		echo "gate (errors.As): clean"
	fi

	is_out="$(erraudit lint ./... --type-aware --type legacy_is -q 2>&1 || true)"
	is_count="$(grep -cE '\.go:[0-9]+:[0-9]+:' <<<"$is_out" || true)"
	if [ "$is_count" -gt 0 ]; then
		echo "advisory: $is_count errors.Is finding(s) (decide per case; sentinels are often correct):"
		echo "$is_out"
	else
		echo "advisory (errors.Is): clean"
	fi

	cd - >/dev/null
done

echo
if [ "$GATE_RED" -eq 1 ]; then
	echo "VERDICT: RED — run 'scripts/erraudit-sweep.sh --fix', review the diff, test."
	exit 1
fi
echo "VERDICT: PASS — zero errors.As findings across all modules."
