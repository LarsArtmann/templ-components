#!/usr/bin/env bash
# coverage-library.sh — coverage for the PUBLISHED library, split by code origin.
#
# The raw `go test -cover` number is dominated by templ's generated
# `*_templ.go`: every element emits a `WriteString` plus an unreachable
# `if err != nil { return err }` branch. No amount of prop permutation can
# cover those branches, so raw coverage plateaus around 73% no matter how
# thorough the tests are. The honest adoption signal is the share of
# HAND-WRITTEN library code covered.
#
# This script measures both, over the published modules only (the demo app in
# examples/demo and the tc CLI are not part of the library):
#
#   raw          every .go file, generated templ output included
#   hand-written generated `*_templ.go` excluded
#
# Usage:
#   scripts/coverage-library.sh                 # measure, report, write reports/
#   scripts/coverage-library.sh --gate 80       # exit 1 if hand-written < 80%
#   scripts/coverage-library.sh --profile f.out # reuse a merged profile, no tests
#   scripts/coverage-library.sh --record        # refresh scripts/coverage-baseline.txt
#
# 80% is the adoption threshold (`awesome-templ`'s requirement, see TODO #28).
# CI runs --gate 94: a ratchet floor just below the measured 96.7%, so a
# coverage regression fails the build instead of quietly eroding the claim.
#
# The measured numbers are written to reports/coverage-library.txt (gitignored,
# uploaded as a CI artifact). The figures the docs quote live in
# scripts/coverage-baseline.txt, guarded by utils/TestCoverageBaselineMatchesProse.
#
# Sub-modules are discovered from the tree (every go.mod except visualtest,
# website, and examples) so a newly split module is covered automatically.
# Each module is measured in isolation (GOWORK=off) so the result matches what
# a consumer's `go test` sees; profiles are merged by block.
set -euo pipefail

cd "$(dirname "$0")/.."

BASELINE="scripts/coverage-baseline.txt"
REPORT_DIR="${REPORT_DIR:-reports}"

usage() {
	echo "usage: $0 [--gate <percent>] [--profile <file>] [--record]"
}

GATE=""
PROFILE=""
RECORD=0

while [ $# -gt 0 ]; do
	case "$1" in
	--gate)
		GATE="${2:-}"
		case "$GATE" in
		'' | *[!0-9]*)
			echo "error: --gate needs a percentage, got '${GATE}'" >&2
			usage >&2
			exit 2
			;;
		esac
		shift 2
		;;
	--profile)
		PROFILE="${2:-}"
		if [ -z "$PROFILE" ] || [ ! -f "$PROFILE" ]; then
			echo "error: --profile needs a readable merged profile, got '${PROFILE}'" >&2
			usage >&2
			exit 2
		fi
		shift 2
		;;
	--record)
		RECORD=1
		shift
		;;
	-h | --help)
		usage
		exit 0
		;;
	*)
		echo "error: unknown argument '$1'" >&2
		usage >&2
		exit 2
		;;
	esac
done

export GOEXPERIMENT="${GOEXPERIMENT:-jsonv2}"

# --- output plumbing: every report line goes to stdout and, when measuring,
# --- into reports/coverage-library.txt for the CI artifact.
REPORT_FILE=""
if [ -z "$PROFILE" ]; then
	REPORT_FILE="$REPORT_DIR/coverage-library.txt"
	mkdir -p "$REPORT_DIR"
	: >"$REPORT_FILE"
fi

say() {
	printf '%s\n' "$*"
	if [ -n "$REPORT_FILE" ]; then
		printf '%s\n' "$*" >>"$REPORT_FILE"
	fi
}

sayf() { # sayf <format> [args...]
	local line
	# shellcheck disable=SC2059 # the format is a caller-supplied literal
	line="$(printf "$@")"
	say "$line"
}

if [ -n "$PROFILE" ]; then
	say "Reusing merged profile: $PROFILE (per-module breakdown skipped)"
	MERGED="$PROFILE"
else
	TMP="$(mktemp -d)"
	trap 'rm -rf "$TMP"' EXIT

	# Discover published sub-modules: every go.mod outside the repo-local
	# visualtest/website helpers and the non-module examples tree.
	MODULES=""
	for modfile in */go.mod */*/go.mod; do
		[ -f "$modfile" ] || continue

		dir="$(dirname "$modfile")"
		case "$dir" in
		visualtest | website | examples/*) continue ;;
		esac

		MODULES="$MODULES $dir"
	done

	say "Root module (minus examples)..."
	# shellcheck disable=SC2046 # package list must word-split into arguments
	go list ./... | grep -v '/examples/' | xargs go test -coverprofile="$TMP/root.out" -count=1 >/dev/null
	PROFILES="$TMP/root.out"

	for mod in $MODULES; do
		say "$mod module..."
		name="$(echo "$mod" | tr '/' '-')"
		(cd "$mod" && GOWORK=off go test ./... -coverprofile="$TMP/$name.out" -count=1 >/dev/null)
		PROFILES="$PROFILES $TMP/$name.out"
	done

	MERGED="$TMP/merged.out"
	{
		echo "mode: set"
		for profile in $PROFILES; do
			tail -n +2 "$profile"
		done
	} >"$MERGED"
fi

# totals <profile> <scope>: prints "covered total percent"
#   library       -> drop examples/ and cmd/tc/ (not published library code)
#   hand-written  -> additionally drop generated *_templ.go
totals() {
	awk -v scope="$2" '
		$1 == "mode:" { next }
		scope == "library" && ($1 ~ /\/examples\// || $1 ~ /\/cmd\/tc\//) { next }
		scope == "hand-written" && ($1 ~ /\/examples\// || $1 ~ /\/cmd\/tc\// || $1 ~ /_templ\.go:/) { next }
		{
			total += $2
			if ($3 > 0) covered += $2
		}
		END {
			if (total == 0) { printf "0 0 0.0\n"; exit }
			printf "%d %d %.1f\n", covered, total, covered * 100 / total
		}
	' "$1"
}

read -r RAW_COV RAW_TOT RAW_PCT < <(totals "$MERGED" library)
read -r HAND_COV HAND_TOT HAND_PCT < <(totals "$MERGED" hand-written)

if [ "$RAW_TOT" = "0" ]; then
	echo "error: profile ${MERGED} contains no library statements" >&2
	exit 2
fi

if [ -z "$PROFILE" ]; then
	sayf '%-16s %10s %8s' "module" "hand-written" "statements"
	sayf '%-16s %10s %8s' "------" "------------" "----------"

	read -r _ root_tot root_pct < <(totals "$TMP/root.out" hand-written)
	sayf '%-16s %9s%% %8s' "root" "$root_pct" "$root_tot"

	for mod in $MODULES; do
		name="$(echo "$mod" | tr '/' '-')"
		read -r _ mod_tot mod_pct < <(totals "$TMP/$name.out" hand-written)
		sayf '%-16s %9s%% %8s' "$mod" "$mod_pct" "$mod_tot"
	done

	say ""
fi

say "Library coverage (published modules only)"
say "  raw (generated templ output included): ${RAW_COV}/${RAW_TOT} = ${RAW_PCT}%"
say "  hand-written (generated excluded):     ${HAND_COV}/${HAND_TOT} = ${HAND_PCT}%"

# The generated share is what it is; covering every hand-written statement
# cannot lift the raw number past this ceiling.
read -r CEIL_PCT GEN_TOT GEN_SHARE_PCT < <(awk -v raw_cov="$RAW_COV" -v raw_tot="$RAW_TOT" -v hand_tot="$HAND_TOT" '
	END {
		gen_cov = raw_cov - (raw_tot - hand_tot)
		gen_tot = raw_tot - hand_tot
		printf "%.1f %d %.1f\n", (gen_cov + hand_tot) * 100 / raw_tot, gen_tot, gen_tot * 100 / raw_tot
	}' </dev/null)

say "  raw ceiling if all hand-written code were covered: $((RAW_COV - (RAW_TOT - HAND_TOT) + HAND_TOT))/${RAW_TOT} = ${CEIL_PCT}%"
say "  generated share: ${GEN_TOT}/${RAW_TOT} statements (${GEN_SHARE_PCT}% of the library)"

say ""
say "Why the two numbers differ: templ emits an unreachable"
say "'if err := WriteString(...); err != nil { return err }' branch per element,"
say "so generated files cap out well below the hand-written share."

if [ "$RECORD" = "1" ]; then
	cat >"$BASELINE" <<EOF
# Library coverage baseline — single source of truth for the figures quoted in
# AGENTS.md. Regenerate with: scripts/coverage-library.sh --record
# Guard: utils/TestCoverageBaselineMatchesProse (arithmetic + prose equality).
recorded_at=$(date -u +%Y-%m-%d)
recorded_commit=$(git rev-parse --short HEAD 2>/dev/null || echo unknown)
raw_covered=${RAW_COV}
raw_total=${RAW_TOT}
raw_pct=${RAW_PCT}
hand_covered=${HAND_COV}
hand_total=${HAND_TOT}
hand_pct=${HAND_PCT}
ceiling_pct=${CEIL_PCT}
generated_share_pct=${GEN_SHARE_PCT}
EOF
	say ""
	say "Recorded baseline: ${BASELINE} (${HAND_PCT}% hand-written, ${RAW_PCT}% raw, ${CEIL_PCT}% ceiling)"
fi

# A stale baseline is not fatal (coverage improving is good news), but the docs
# then understate the library — say so instead of letting it rot silently.
if [ "$RECORD" = "0" ] && [ -f "$BASELINE" ]; then
	recorded="$(awk -F= '$1 == "hand_pct" { print $2 }' "$BASELINE")"
	if [ -n "$recorded" ] && ! awk -v a="$recorded" -v b="$HAND_PCT" 'BEGIN { exit !(a + 0 > b + 0) }'; then
		if awk -v a="$recorded" -v b="$HAND_PCT" 'BEGIN { exit !((b - a) >= 0.1) }'; then
			say ""
			say "note: ${BASELINE} records ${recorded}% hand-written, measured ${HAND_PCT}%."
			say "      Re-run with --record to refresh the baseline and the AGENTS.md figures."
		fi
	fi
fi

if [ -n "$GATE" ]; then
	if awk -v pct="$HAND_PCT" -v gate="$GATE" 'BEGIN { exit !(pct + 0 < gate + 0) }'; then
		echo "FAIL: hand-written coverage ${HAND_PCT}% is below the ${GATE}% gate" >&2
		exit 1
	fi
	echo "PASS: hand-written coverage ${HAND_PCT}% meets the ${GATE}% gate"
fi
