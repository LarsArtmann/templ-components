#!/usr/bin/env bash
# fleet-css-check.sh — one-command stale-CSS audit across the consumer fleet.
#
# Reproduces the #410c sweep's grep-the-artifact membership check as a
# standing script: for every fleet repo, collect the library's tc-* hook
# classes that can actually render (the consumer's own sources PLUS the
# templ-components module version pinned in its go.mod, resolved
# REPLACE-AWARE per SKILL.md rule 2), then verify every hook reaches the
# repo's compiled CSS artifact. A hook missing from the artifact renders
# unstyled outside local dev — the pbx-lens (tc-log*) and nsfw-classifier
# (tc-drawer/tc-toast-container) outages.
#
# Usage:
#   scripts/fleet-css-check.sh                 # audit the whole fleet
#   scripts/fleet-css-check.sh pbx-lens crm    # audit specific repos
#   FLEET_ROOT=/path/to/projects scripts/fleet-css-check.sh
#
# Exit 0 = every checked repo covers its hooks (or nothing to check);
# exit 1 = at least one MISS. WARN lines are informational, never fatal.
set -uo pipefail

LIB_MODULE="github.com/larsartmann/templ-components"
PROJECTS_ROOT="${FLEET_ROOT:-$HOME/projects}"

# The #410c fleet. Repos not checked out locally are skipped with a note.
DEFAULT_FLEET=(
	nsfw-classifier pbx-lens KeyHolderAI PapDashboard sperrmuell-direct
	go-taskqueue webphone Rolls-Royce GmbH Zlota44 crm InboxClean
	browser-history DiscordSync cqrs-htmx dnsblockd mr-sync
)

FLEET=("$@")
((${#FLEET[@]} == 0)) && FLEET=("${DEFAULT_FLEET[@]}")

miss_total=0
checked=0

# library_dirs <repo>: print source dirs of every templ-components module
# version the repo builds against — the replace target when one exists
# (bridge checkouts), else the module-cache extraction of the pinned
# version. Prints nothing (caller warns) when the version is not in the
# local module cache.
library_dirs() {
	local repo=$1
	(
		cd "$repo" || exit 0
		GOWORK=off go list -m -f '{{.Path}}{{"\t"}}{{if .Replace}}{{.Replace.Dir}}{{else}}{{.Dir}}{{end}}' all 2>/dev/null |
			rg "^${LIB_MODULE}(/|\t)" | cut -f2
	) 2>/dev/null | rg -v '^$'
}

# hooks_in <dir>...: every distinct tc-* hook literal in templ/Go sources.
# third_party copies are excluded — they are scan inventories, not sources.
hooks_in() {
	rg -o 'tc-[a-z0-9-]+' -g '*.templ' -g '*.go' --no-filename \
		-g '!**/third_party/**' -g '!**/node_modules/**' "$@" 2>/dev/null |
		sort -u || true
}

# artifact_for <repo>: the repo's tracked, compiled CSS artifact (the file
# a real deployment serves), or empty when there is none to check.
artifact_for() {
	local repo=$1
	git -C "$repo" ls-files '*.css' 2>/dev/null |
		rg -v 'node_modules|third_party|templates/' |
		rg '(\.out\.css|app\.css|main\.css|styles\.css|\.min\.css)$' |
		head -1
}

for repo_name in "${FLEET[@]}"; do
	repo="$PROJECTS_ROOT/$repo_name"
	if [[ ! -d $repo/.git && ! -f $repo/go.mod ]]; then
		echo "SKIP  $repo_name (not checked out under $PROJECTS_ROOT)"
		continue
	fi

	artifact_rel=$(artifact_for "$repo")
	if [[ -z $artifact_rel ]]; then
		echo "WARN  $repo_name: no tracked compiled CSS artifact found — nothing to check"
		continue
	fi
	artifact="$repo/$artifact_rel"

	# Hooks that can render = consumer's own sources + the pinned library
	# modules (the tc-* hooks live in the LIBRARY's generated files; a
	# module-importing consumer's own sources never mention them).
	mapfile -t lib_dirs < <(library_dirs "$repo")
	if ((${#lib_dirs[@]} == 0)); then
		echo "WARN  $repo_name: templ-components not resolvable (go mod download?) — checking consumer sources only"
	fi

	hooks=$(hooks_in "$repo" "${lib_dirs[@]}")
	if [[ -z $hooks ]]; then
		echo "OK    $repo_name (no tc-* hooks in scope)"
		continue
	fi

	missing=0
	hook_count=0
	while IFS= read -r hook; do
		[[ -z $hook ]] && continue
		hook_count=$((hook_count + 1))
		if ! rg -qF -- "$hook" "$artifact"; then
			echo "MISS  $repo_name: $hook not in $artifact_rel"
			missing=$((missing + 1))
		fi
	done <<<"$hooks"

	checked=$((checked + 1))
	if ((missing == 0)); then
		echo "OK    $repo_name: $artifact_rel covers all $hook_count hooks"
	else
		miss_total=$((miss_total + missing))
	fi
done

echo
echo "Fleet audit: $checked repo(s) checked, $miss_total missing hook(s)."
((miss_total == 0)) || exit 1
