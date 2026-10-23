#!/usr/bin/env bash
# settle-staged-go-format.sh — make the INDEX gofmt-clean before BuildFlow runs.
#
# Why (#441): BuildFlow's pre-commit DAG has no edge from golangci-lint to the
# formatter steps (golangci-lint depends on go-fix + workspace-build-verify,
# both blocklisted in pre-commit). nix-fmt/templ-fmt therefore rewrite staged
# .go files concurrently with lint reading them; the 2026-10-23 commit failed
# golangci-lint on six modules while the formatters were mid-write, and the
# identical run passed once the tree was settled.
#
# Safety: this formats the STAGED BLOB, never the worktree file. Unstaged edits
# cannot leak into the commit — the worktree is untouched and only the index
# entry is replaced (the same effect BuildFlow's own post-format re-stage has).
set -euo pipefail

REPO_ROOT="${SETTLE_REPO_ROOT:-$(cd "$(dirname "$0")/.." && pwd)}"
cd "$REPO_ROOT"

if ! command -v gofmt >/dev/null 2>&1; then
	echo "settle-staged-go-format: gofmt not on PATH; leaving the index untouched" >&2
	exit 0
fi

settled=0

while IFS= read -r file; do
	[ -n "$file" ] || continue

	entry="$(git ls-files -s -- "$file")"
	[ -n "$entry" ] || continue
	mode="${entry%% *}"

	staged="$(git cat-file blob ":$file")"
	formatted="$(printf '%s\n' "$staged" | gofmt)" || continue
	if [ "$formatted" = "$staged" ]; then
		continue
	fi

	blob="$(printf '%s\n' "$formatted" | git hash-object -w --stdin)"
	git update-index --cacheinfo "${mode},${blob},${file}"
	echo "settle-staged-go-format: index now holds the gofmt-clean version of ${file}" >&2
	settled=$((settled + 1))
done < <(git diff --cached --name-only --diff-filter=ACMR -- '*.go')

if [ "$settled" -gt 0 ]; then
	echo "settle-staged-go-format: settled ${settled} staged Go file(s)" >&2
fi
