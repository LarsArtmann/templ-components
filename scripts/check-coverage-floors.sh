#!/usr/bin/env bash
# check-coverage-floors.sh — per-package coverage floors (M19/F085).
#
# The total-coverage gate (70%) hides package-level regressions: a package
# can drop 10 points while the total barely moves. This script aggregates
# coverage.out per package and enforces the floor recorded for each in
# scripts/coverage-floors.txt ("<pkg> <min%>"). Floors start at current-2
# and ratchet UP only — lowering a floor requires an owner decision in the
# PR that does it.
#
# Usage: scripts/check-coverage-floors.sh coverage.out

set -euo pipefail

PROFILE="${1:-coverage.out}"
FLOORS="$(dirname "$0")/coverage-floors.txt"

[ -f "$PROFILE" ] || {
	echo "usage: $0 <coverage.out> (run go test -coverprofile first)" >&2
	exit 2
}
[ -f "$FLOORS" ] || {
	echo "error: $FLOORS missing — floors must exist; regenerate with -update" >&2
	exit 2
}

# Aggregate: covered statements / total statements per package dir.
# Profile lines: "pkg/file.go:3.10,5.2 <numStmts> <count>"
declare -A COVERED TOTAL
while IFS=' ' read -r loc stmts count; do
	[ -n "${count:-}" ] || continue # header or blank line
	file="${loc%%:*}"
	pkg="$(dirname "$file")"
	pkg="${pkg#github.com/larsartmann/templ-components/}" # floors use short names
	[ -n "${TOTAL[$pkg]:-}" ] || TOTAL[$pkg]=0
	[ -n "${COVERED[$pkg]:-}" ] || COVERED[$pkg]=0
	TOTAL[$pkg]=$((TOTAL[$pkg] + stmts))
	if [ "$count" -gt 0 ]; then
		COVERED[$pkg]=$((COVERED[$pkg] + stmts))
	fi
done <"$PROFILE"

FAILED=0
while read -r pkg floor; do
	[ -z "${pkg:-}" ] && continue
	case "$pkg" in \#*) continue ;; esac

	total="${TOTAL[$pkg]:-}"
	if [ -z "$total" ]; then
		echo "error: floor exists for $pkg but the profile has no blocks — stale floor?" >&2
		FAILED=1
		continue
	fi

	actual=$((COVERED[$pkg] * 100 / total))
	if [ "$actual" -lt "$floor" ]; then
		echo "error: $pkg coverage ${actual}% below floor ${floor}% (ratchet down requires an owner decision)" >&2
		FAILED=1
	else
		printf 'ok: %-12s %s%% (floor %s%%)\n' "$pkg" "$actual" "$floor"
	fi
done <"$FLOORS"

exit "$FAILED"
