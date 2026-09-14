#!/usr/bin/env bash
# scripts/check-replace-directives.sh — tripwire for the release-race class.
#
# scripts/release.sh strips the local replace directives from every module's
# go.mod before tagging (published modules must not carry local-path replaces),
# then re-adds them in a follow-up commit. The v1.17.0 cut proved that re-add
# can silently vanish when the BuildFlow daemon races the script: all 5
# dependent sub-modules (icons, errorpage, charts/echarts, htmx, datastar)
# ended up with NO replace blocks, breaking every per-module go test / lint
# run ("missing go.sum entry") until manually repaired (2026-09-13).
#
# check-module-sync.sh only validates the FORM of replaces that exist
# (relative paths, shared versions); it passes when they are absent entirely.
# This guard pins the EXACT expected set per module, so a vanished re-add —
# or a lost/wrong-nested path (charts/echarts needs ../../utils) — fails in
# <50ms instead of surfacing as a mystery lint typecheck failure.
#
# Expected state: master between releases ALWAYS carries the full replace set.
# The release script's strip commit (--no-verify, immediate) is superseded by
# the re-add commit; if the re-add ever fails to land, this guard fires on the
# next CI run / commit.
#
# Limitation: only standalone `replace A => B` lines are recognized (the form
# `go mod edit -replace` produces). A grouped `replace (...)` block makes the
# missing-checks fail loud — update the normalizer if that form is ever used.

set -euo pipefail
cd "$(dirname "$0")/.."

ERRORS=0
MOD="github.com/larsartmann/templ-components"

# replaces_in <go.mod> — emit normalized "module=>path" lines for standalone
# replace directives.
replaces_in() {
	grep -E '^[[:space:]]*replace[[:space:]]' "$1" 2>/dev/null |
		sed -E 's/^[[:space:]]*replace[[:space:]]+//; s/[[:space:]]+=>[[:space:]]+/=>/; s/[[:space:]]+$//' || true
}

# check_module <go.mod> [expected replace ...] — every expected replace must be
# present, and no unexpected templ-components-sibling replace may appear
# (external/local-debug replaces are out of scope).
check_module() {
	local modfile="$1"
	shift
	local actual expected line
	if [ ! -f "$modfile" ]; then
		echo "::error::${modfile} does not exist"
		ERRORS=$((ERRORS + 1))
		return
	fi
	actual="$(replaces_in "$modfile")"

	for expected in "$@"; do
		if ! printf '%s\n' "$actual" | grep -Fxq "$expected"; then
			echo "::error::${modfile}: missing replace '${expected/=>/ => }'"
			ERRORS=$((ERRORS + 1))
		fi
	done

	while IFS= read -r line; do
		[ -n "$line" ] || continue
		case "$line" in
		"$MOD" | "$MOD"/*)
			if ! printf '%s\n' "$@" | grep -Fxq "$line"; then
				echo "::error::${modfile}: unexpected sibling replace '${line/=>/ => }' (update scripts/check-replace-directives.sh if intentional)"
				ERRORS=$((ERRORS + 1))
			fi
			;;
		esac
	done <<<"$actual"
}

# Root module: self + all 6 sub-modules, ./ paths.
check_module go.mod \
	"$MOD=>./" \
	"$MOD/utils=>./utils" \
	"$MOD/icons=>./icons" \
	"$MOD/errorpage=>./errorpage" \
	"$MOD/charts/echarts=>./charts/echarts" \
	"$MOD/datastar=>./datastar" \
	"$MOD/htmx=>./htmx"

# Leaf module: no sibling replaces at all.
check_module utils/go.mod

# Layer 1: utils only. NOTE charts/echarts nests one level deeper (../../utils).
check_module icons/go.mod "$MOD/utils=>../utils"
check_module charts/echarts/go.mod "$MOD/utils=>../../utils"
check_module htmx/go.mod "$MOD/utils=>../utils"
check_module datastar/go.mod "$MOD/utils=>../utils"

# Layer 2: utils + icons.
check_module errorpage/go.mod \
	"$MOD/utils=>../utils" \
	"$MOD/icons=>../icons"

# Dev-only modules (not published, never stripped): root + all 6, ../ paths.
check_module visualtest/go.mod \
	"$MOD=>.." \
	"$MOD/utils=>../utils" \
	"$MOD/icons=>../icons" \
	"$MOD/errorpage=>../errorpage" \
	"$MOD/charts/echarts=>../charts/echarts" \
	"$MOD/datastar=>../datastar" \
	"$MOD/htmx=>../htmx"

check_module website/go.mod \
	"$MOD=>../" \
	"$MOD/utils=>../utils" \
	"$MOD/icons=>../icons" \
	"$MOD/errorpage=>../errorpage" \
	"$MOD/charts/echarts=>../charts/echarts" \
	"$MOD/datastar=>../datastar" \
	"$MOD/htmx=>../htmx"

if [ "$ERRORS" -gt 0 ]; then
	echo "Replace-directives check: ${ERRORS} error(s) found."
	exit 1
fi

echo "Replace-directives check: OK (9 modules, exact replace sets pinned)."
