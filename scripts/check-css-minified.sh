#!/usr/bin/env bash
# Guard 6: Demo CSS minification check
#   examples/demo/static/app.css must be MINIFIED (<= 5 lines). BuildFlow's
#   tailwind-build provider regenerates it WITHOUT --minify, so daemon
#   auto-commits repeatedly land a ~5000-line un-minified CSS on master
#   (documented recurrences: 2026-09-03, 2026-09-08 — TODO #125, root fix
#   lives in larsartmann/buildflow). This <50ms commit-time guard blocks the
#   bad commit instead of letting CI's CSS Freshness job go red later.
#   Re-minify with: nix run .#css
set -euo pipefail

CSS="examples/demo/static/app.css"

if [ ! -f "$CSS" ]; then
	exit 0
fi

LINES=$(wc -l <"$CSS")
if [ "$LINES" -gt 5 ]; then
	echo "ERROR: $CSS has $LINES lines — it is NOT minified."
	echo "BuildFlow's tailwind-build provider rewrote it without --minify"
	echo "(TODO #125 recurrence). Re-minify and re-stage:"
	echo "  nix run .#css && git add $CSS"
	exit 1
fi
