#!/usr/bin/env bash
# Fast (<100ms) pre-commit guard: every template embedded in the `tc new`
# scaffolder (cmd/tc/_sources/**) must be byte-identical to its library twin.
#
# WHY: the scaffolder ships COPIES of real component sources. When a component
# is edited without re-copying, `tc new` scaffolds stale code — and the Go
# guard (cmd/tc TestSourcesMatchPackageFiles) only fires in `go test`, which
# BuildFlow's pre-commit budget cannot run. This shell mirror catches the
# drift at commit time. (Found twice on 2026-09-17: notfound404 + shared.)
#
# Exclusions mirror the Go test:
#   - starter/            curated CSS, no library twin
#   - the datastar bump-protocol doc (scaffolder-owned)
#
# Usage:
#   scripts/check-tc-sources-sync.sh            # exit 1 on drift
#   scripts/check-tc-sources-sync.sh --fix      # re-copy drifted files, then exit 1
set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
cd "$REPO_ROOT"

FIX="${1:-}"
SOURCES_DIR="cmd/tc/_sources"

# The scaffolder-owned datastar doc has no library twin (name kept in sync
# with cmd/tc/main_test.go's exclusion).
DATASTAR_DOC="docs/datastar-dependency-bump-protocol.md"

DRIFT=0

while IFS= read -r -d '' embedded; do
	rel="${embedded#"$SOURCES_DIR"/}"

	# Exclusions (mirror the Go test).
	case "$rel" in
	starter/*) continue ;;
	"src/datastar/$DATASTAR_DOC") continue ;;
	datastar/*docs*) continue ;;
	esac

	# Only .templ files have library twins here; skip anything else that
	# lacks a counterpart.
	if [ ! -f "$rel" ]; then
		continue
	fi

	if ! cmp -s "$embedded" "$rel"; then
		echo "DRIFT: $embedded differs from $rel" >&2
		DRIFT=1
		if [ "$FIX" = "--fix" ]; then
			cp "$rel" "$embedded"
			echo "  fixed: re-copied $rel -> $embedded"
		fi
	fi
done < <(find "$SOURCES_DIR" -type f -name '*.templ' -print0)

if [ "$DRIFT" -ne 0 ]; then
	if [ "$FIX" != "--fix" ]; then
		echo "" >&2
		echo "BLOCKED: cmd/tc/_sources drifted from the library sources." >&2
		echo "Re-copy the drifted files (or run scripts/check-tc-sources-sync.sh --fix)" >&2
		echo "and commit the refreshed scaffolder copies." >&2
	fi
	exit 1
fi

echo "tc scaffolder sources in sync with library twins."
exit 0
