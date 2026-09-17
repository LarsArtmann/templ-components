#!/usr/bin/env bash
# Fast (<1s) pre-commit guard: blocks commits if *_templ.go files have
# drifted from their .templ source (import mismatch, EITHER direction).
#
# This catches the recurring regression where a stale working tree or the
# BuildFlow daemon commits a *_templ.go that doesn't match the .templ source:
#   - 2026-07-28: source imported encoding/json (v1), generated had v2
#   - 2026-09-17: website/internal/pages/base_templ.go flipped to
#     encoding/json/v2 while base.templ said v1 (invisible because the guard
#     only walked root-module packages)
#   - 2026-09-17 20:52: daemon committed examples/demo/errorpage_demo_templ.go
#     stale relative to errorpage_demo.templ (CI "Verify no untracked changes")
#
# Mirrors the logic of utils/templ_sync_test.go (TestTemplGeneratedInSync)
# but runs without `go test` so it fits within BuildFlow's pre-commit budget.
#
# Checks, REPO-WIDE (all modules incl. website + examples/demo):
#   1. every import in the .templ source must appear in the generated _templ.go
#   2. every non-runtime import in the generated _templ.go must appear in the
#      .templ source (the generator copies source imports; a mismatch means the
#      generated file is stale or hand-edited)
# Excludes github.com/a-h/templ runtime imports (always injected by the
# generator) and cmd/tc/_sources (scaffolding templates with no twins).
#
# Usage:
#   scripts/check-templ-sync.sh          # exit 1 on drift
#   scripts/check-templ-sync.sh --quiet  # suppress output
set -euo pipefail

QUIET="${1:-}"

# Run from the repo root regardless of the caller's CWD.
REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
cd "$REPO_ROOT"

DRIFT_FOUND=0
DRIFT_OUTPUT=""

# extract_templ_imports prints the quoted import paths of a .templ or
# generated .go file (import blocks + single-line imports).
# NOTE: defined before use — bash resolves functions at execution time, and a
# late definition silently yields empty extractions (a false-green guard).
extract_templ_imports() {
	awk '
		/^import[[:space:]]*\(/ { in_block = 1; next }
		in_block && /^\)/ { in_block = 0; next }
		in_block && /"[^"]+"/ {
			match($0, /"[^"]+"/)
			print substr($0, RSTART+1, RLENGTH-2)
		}
		/^import[[:space:]]+"[^"]+"/ {
			match($0, /"[^"]+"/)
			print substr($0, RSTART+1, RLENGTH-2)
		}
	' "$1"
}

report() {
	DRIFT_OUTPUT+="  $1\n"
	DRIFT_FOUND=1
}

# find .templ files repo-wide, excluding scaffold sources and noise dirs.
while IFS= read -r -d '' templ_file; do
	gen_file="${templ_file%.templ}_templ.go"

	if [ ! -f "$gen_file" ]; then
		report "MISSING: $gen_file (not generated yet)"
		continue
	fi

	# Direction 1: source imports must exist in the generated file.
	while IFS= read -r imp; do
		[ -z "$imp" ] && continue

		# Skip templ runtime imports (generator always injects these).
		case "$imp" in
		github.com/a-h/templ*) continue ;;
		esac

		if ! grep -qF "\"$imp\"" "$gen_file" 2>/dev/null; then
			report "DRIFT: $templ_file imports \"$imp\" but $gen_file does not"
		fi
	done < <(
		extract_templ_imports "$templ_file"
	)

	# Direction 2: generated imports must exist in the source (staleness tripwire).
	while IFS= read -r imp; do
		[ -z "$imp" ] && continue

		case "$imp" in
		github.com/a-h/templ*) continue ;;
		esac

		if ! grep -qF "\"$imp\"" "$templ_file" 2>/dev/null; then
			report "STALE: $gen_file imports \"$imp\" but $templ_file does not (regenerate)"
		fi
	done < <(
		extract_templ_imports "$gen_file"
	)
done < <(
	find . -name '*.templ' \
		-not -path './.git/*' \
		-not -path '*/node_modules/*' \
		-not -path '*/dist/*' \
		-not -path '*/testdata/*' \
		-not -path './cmd/tc/_sources/*' \
		-print0
)

if [ "$DRIFT_FOUND" -ne 0 ]; then
	if [ "$QUIET" != "--quiet" ]; then
		echo "" >&2
		echo "BLOCKED: *_templ.go files are out of sync with .templ sources." >&2
		echo "" >&2
		echo -e "$DRIFT_OUTPUT" >&2
		echo "Fix: run 'templ generate ./...' and commit the updated *_templ.go files." >&2
		echo "" >&2
	fi
	exit 1
fi

exit 0
