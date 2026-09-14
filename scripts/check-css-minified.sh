#!/usr/bin/env bash
# Guard 6: compiled-CSS minification check
#
# Every compiled CSS distribution target (scripts/compiled-css-targets.txt —
# the single source shared with release.sh and TestCompiledCSSInventory) must
# be MINIFIED (<= 5 lines). BuildFlow's tailwind-build provider regenerates
# them WITHOUT --minify, so daemon auto-commits repeatedly land ~5000-line
# un-minified CSS on master (documented recurrences: 2026-09-03, 2026-09-08 —
# TODO #125, root fix lives in larsartmann/buildflow). This <50ms commit-time
# guard blocks the bad commit instead of letting CI's CSS Freshness job go red
# later. Re-minify with: nix run .#css
set -euo pipefail

TARGETS_FILE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/compiled-css-targets.txt"

if [ ! -f "$TARGETS_FILE" ]; then
	echo "ERROR: $TARGETS_FILE missing — the minification guard owns this fixture (fail loud, never skip)." >&2
	exit 1
fi

FAILED=0

while read -r _INPUT OUTPUT; do
	case "$OUTPUT" in "" | \#*) continue ;; esac
	if [ ! -f "$OUTPUT" ]; then
		continue
	fi
	LINES=$(wc -l <"$OUTPUT")
	if [ "$LINES" -gt 5 ]; then
		echo "ERROR: $OUTPUT has $LINES lines — it is NOT minified."
		echo "BuildFlow's tailwind-build provider rewrote it without --minify"
		echo "(TODO #125 recurrence). Re-minify and re-stage:"
		echo "  nix run .#css && git add $OUTPUT"
		FAILED=1
	fi
done < <(grep -vE '^\s*(#|$)' "$TARGETS_FILE")

exit "$FAILED"
