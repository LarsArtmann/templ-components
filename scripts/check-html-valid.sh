#!/usr/bin/env bash
# check-html-valid.sh — W3C HTML validation over the golden corpus (M17/F075).
#
# Every *.golden snapshot is wrapped in a minimal document scaffold (unless it
# already is one) and validated with the Nu Html Checker (html5validator from
# nixpkgs). The validator's snapshot lags the HTML spec; every ignore pattern
# below is a DOCUMENTED staleness or framework-dialect class, with the reason.
# When a newer vnu lands in nixpkgs, prune this list.
#
# Usage:
#   nix shell nixpkgs#html5validator -c scripts/check-html-valid.sh
#
# Exit non-zero when any unignored error is found; prints the full report.

set -euo pipefail
cd "$(dirname "$0")/.."

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

PKG_DIRS="display feedback forms layout navigation htmx errorpage charts/echarts datastar"

for g in $PKG_DIRS; do
	d="$WORK/pkg-${g//\//-}"
	mkdir -p "$d"
	i=0
	for f in "$g"/testdata/*.golden; do
		[ -e "$f" ] || continue
		i=$((i + 1))
		out="$d/$(basename "$f" .golden).html"
		if head -c 20 "$f" | grep -qi '<!doctype\|<html'; then
			cp "$f" "$out"
		else
			{
				printf '<!DOCTYPE html><html><head><title>t</title></head><body>'
				cat "$f"
				printf '</body></html>'
			} >"$out"
		fi
	done
done

# Ignored error classes (regex). html5validator 0.4.2's --ignore does plain
# substring matching (verified empirically), so filtering happens here with
# real regexes. Every pattern is a DOCUMENTED validator-staleness or
# framework-dialect class:
#   hx-[a-z-]+" not allowed      htmx dialect attributes — not in any spec by design
#   "data-\*"...NCNames          Datastar's data-on:click colon syntax — deliberate
#   popover|popovertarget        vnu predates the Popover API (Baseline 2024)
#   fetchpriority|enterkeyhint   vnu staleness — standard attributes
#   Element "search" not allowed vnu predates the <search> element (Baseline 2023)
#   svg...as child of summary    spec disagreement; browsers+a11y treat svg in
#                                summary as phrasing content (every accordion chevron)
#   style...child of body        ViewTransitions emits a style fragment position
#   CSS: Parse Error             vnu's CSS parser predates @view-transition
#   Invalid RGB function         vnu's CSS parser predates CSS Color 4 rgb(a b c / d)
#   Stray (start|end) tag        vnu's tree builder predates the customizable
#     (button|selectedcontent)   <select> (button + selectedcontent inside select)
#
# Every quoted token matches BOTH straight (") and curly (“”) quotes: vnu's
# message quoting changed across releases, and CI downloads vnu.jar at
# runtime while local html5validator bundles an older checker — one list
# must serve both snapshots.
IGNORE_RE='Attribute ["“]hx-[a-z-]+["”] not allowed|["“]data-\*["”] attribute names|["“](popover|popovertarget|fetchpriority|enterkeyhint)["”] not allowed|Element ["“]search["”] not allowed|Element ["“]svg["”] not allowed as child of element ["“]summary["”]|Element ["“]style["”] not allowed as child of element ["“]body["”]|CSS: Parse Error|["“]background-color["”]: Invalid RGB function|Stray (start|end) tag ["“](button|selectedcontent)["”]'

# Validator invocation: html5validator (nixpkgs local) or a direct vnu.jar
# via VNU_JAR (CI downloads the jar — no pip/PEP-668 involved). Output
# message format is identical (html5validator wraps the same checker).
if [ -n "${VNU_JAR:-}" ]; then
	report="$(java -jar "$VNU_JAR" $(find "$WORK" -name '*.html') 2>&1 || true)"
else
	report="$(html5validator --root "$WORK" 2>/dev/null || true)"
fi
filtered="$(printf '%s\n' "$report" | grep 'error:' | grep -vE "$IGNORE_RE" || true)"

if [ -n "$filtered" ]; then
	printf 'HTML validation errors (after documented ignores):\n%s\n' "$filtered" >&2
	exit 1
fi

echo "HTML validation clean: $(
	cd "$WORK" && ls pkg-*/*.html 2>/dev/null | wc -l
) golden files, $(printf '%s\n' "$IGNORE_RE" | tr '|' '\n' | wc -l) documented ignore classes."
