#!/usr/bin/env bash
# check-html-report-classes.sh — every class="..." token in docs/research/*.html
# must exist as a selector in that report's embedded <style> block.
#
# Why: the deep-dive reports are hand-authored, self-contained HTML that ships
# without a build step. A class token with no CSS rule renders as unstyled
# markup and no test catches it — this script is that test (found 0 misses on
# the 2026-10-08 baseline; the 2 report files are the current population).
#
# Screenshot-pass checklist for HTML deliverables (run before publishing):
#   1. bash scripts/check-html-report-classes.sh   <- this script, class coverage
#   2. headless capture + eyeball:
#      nix shell nixpkgs#chromium -c chromium --headless --no-sandbox \
#        --screenshot=/tmp/out.png --window-size=1440,9000 --virtual-time-budget=4000 \
#        "file://$PWD/docs/research/<report>.html"
#   3. crop/resize and VIEW the PNG (hero, scorecard, tables, footer).
#
# Allowlist: space-separated class names in ALLOWED_MISSING (state classes
# toggled by JS that intentionally have no CSS rule).

set -euo pipefail

cd "$(dirname "$0")/.."

ALLOWED_MISSING="${ALLOWED_MISSING:-}"
fail=0

for html in docs/research/*.html; do
	[ -e "$html" ] || continue

	# Strip <pre>/<code> blocks first: their escaped-HTML samples contain
	# literal class="..." text that is CONTENT, not markup.
	body=$(perl -0777 -pe 's/<pre\b.*?<\/pre>//gs; s/<code\b.*?<\/code>//gs' "$html")
	style=$(sed -n '/<style/,/<\/style>/p' <<<"$body")

	# shellcheck disable=SC2001  # word-splitting over class lists is the point
	tokens=$(grep -o 'class="[^"]*"' <<<"$body" |
		sed 's/class="//;s/"$//' |
		tr ' ' '\n' |
		sed '/^$/d' |
		sort -u)

	missing=""
	while IFS= read -r token; do
		[ -n "$token" ] || continue
		case " $ALLOWED_MISSING " in
		*" $token "*) continue ;;
		esac
		if ! grep -q "\.${token}[^a-zA-Z0-9_-]" <<<"$style" &&
			! grep -q "\.${token}$" <<<"$style"; then
			missing="$missing $token"
		fi
	done <<<"$tokens"

	total=$(wc -l <<<"$tokens" | tr -d ' ')
	if [ -n "$missing" ]; then
		echo "FAIL $html — class tokens with no CSS rule:$missing"
		fail=1
	else
		echo "OK   $html — $total class tokens, all covered"
	fi
done

exit "$fail"
