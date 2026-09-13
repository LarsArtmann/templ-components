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

# --ignore patterns (regex, matched against vnu messages). Each is a known
# false-positive class; the comment is the reason it cannot be a real bug:
#   hx-.*                htmx dialect attributes — not in any HTML spec by design
#   data-\*.*NCNames     Datastar's data-on:click colon syntax — deliberate
#   popover              vnu predates the Popover API (Baseline 2024)
#   fetchpriority        vnu staleness — standard attribute since 2022
#   enterkeyhint         vnu staleness — standard attribute
#   "search" not allowed vnu predates the <search> element (Baseline 2023)
#   svg.*summary         spec disagreement: browsers + a11y treat svg-in-summary
#                        as phrasing content; used by every accordion chevron
#   style.*body          ViewTransitions emits a style fragment; style-in-body
#                        is universally supported
#   CSS: Parse Error     vnu's CSS parser predates @view-transition
#   Invalid RGB function vnu's CSS parser predates CSS Color 4 space-separated rgb()
#   Stray end tag button vnu's tree builder predates the customizable <select>
#                        (button+selectedcontent inside select)
html5validator --root "$WORK" \
	--ignore \
	'hx-[a-z-]+" not allowed' \
	'"data-\*" attribute names' \
	'"popover" not allowed' \
	'"popovertarget" not allowed' \
	'"fetchpriority" not allowed' \
	'"enterkeyhint" not allowed' \
	'"search" not allowed' \
	'"svg" not allowed as child of element "summary"' \
	'"style" not allowed as child of element "body"' \
	'CSS: Parse Error' \
	'"background-color": Invalid RGB function' \
	'Stray end tag "button"' \
	2>/dev/null
