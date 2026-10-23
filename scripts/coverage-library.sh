#!/usr/bin/env bash
# coverage-library.sh — coverage for the PUBLISHED library, split by code origin.
#
# The raw `go test -cover` number is dominated by templ's generated
# `*_templ.go`: every element emits a `WriteString` plus an unreachable
# `if err != nil { return err }` branch. No amount of prop permutation can
# cover those branches, so raw coverage plateaus around 72% no matter how
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
#   scripts/coverage-library.sh              # report both numbers
#   scripts/coverage-library.sh --gate 80    # exit 1 if hand-written < 80%
#
# Sub-modules are discovered from the tree (every go.mod except visualtest,
# website, and examples) so a newly split module is covered automatically.
# Each module is measured in isolation (GOWORK=off) so the result matches what
# a consumer's `go test` sees; profiles are merged by block.
set -euo pipefail

cd "$(dirname "$0")/.."

GATE=""
if [ "${1:-}" = "--gate" ]; then
	GATE="${2:-80}"
	case "$GATE" in
	'' | *[!0-9]*)
		echo "usage: $0 [--gate <percent>]" >&2
		exit 2
		;;
	esac
fi

export GOEXPERIMENT="${GOEXPERIMENT:-jsonv2}"

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

echo "Root module (minus examples)..."
# shellcheck disable=SC2046 # package list must word-split into arguments
go list ./... | grep -v '/examples/' | xargs go test -coverprofile="$TMP/root.out" -count=1 >/dev/null
PROFILES="$TMP/root.out"

for mod in $MODULES; do
	echo "$mod module..."
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

printf '\n%-16s %10s %8s\n' "module" "hand-written" "statements"
printf '%-16s %10s %8s\n' "------" "------------" "----------"

read -r _ root_tot root_pct < <(totals "$TMP/root.out" hand-written)
printf '%-16s %9s%% %8s\n' "root" "$root_pct" "$root_tot"

for mod in $MODULES; do
	name="$(echo "$mod" | tr '/' '-')"
	read -r _ mod_tot mod_pct < <(totals "$TMP/$name.out" hand-written)
	printf '%-16s %9s%% %8s\n' "$mod" "$mod_pct" "$mod_tot"
done

echo
echo "Library coverage (published modules only)"
echo "  raw (generated templ output included): ${RAW_COV}/${RAW_TOT} = ${RAW_PCT}%"
echo "  hand-written (generated excluded):     ${HAND_COV}/${HAND_TOT} = ${HAND_PCT}%"

# The generated share is what it is; covering every hand-written statement
# cannot lift the raw number past this ceiling.
awk -v raw_cov="$RAW_COV" -v raw_tot="$RAW_TOT" -v hand_cov="$HAND_COV" -v hand_tot="$HAND_TOT" '
	END {
		gen_cov = raw_cov - hand_cov
		gen_tot = raw_tot - hand_tot
		printf "  raw ceiling if all hand-written code were covered: %d/%d = %.1f%%\n", gen_cov + hand_tot, raw_tot, (gen_cov + hand_tot) * 100 / raw_tot
		printf "  generated share: %d/%d statements (%.1f%% of the library)\n", gen_tot, raw_tot, gen_tot * 100 / raw_tot
	}' </dev/null

echo
echo "Why the two numbers differ: templ emits an unreachable"
echo "'if err := WriteString(...); err != nil { return err }' branch per element,"
echo "so generated files cap out well below the hand-written share."

if [ -n "$GATE" ]; then
	if awk -v pct="$HAND_PCT" -v gate="$GATE" 'BEGIN { exit !(pct + 0 < gate + 0) }'; then
		echo "FAIL: hand-written coverage ${HAND_PCT}% is below the ${GATE}% gate" >&2
		exit 1
	fi
	echo "PASS: hand-written coverage ${HAND_PCT}% meets the ${GATE}% gate"
fi
