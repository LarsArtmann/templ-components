#!/usr/bin/env bash
# Fast (<100ms) pre-commit guard: blocks commits if *_templ.go files have
# drifted from their .templ source (import mismatch). This catches the
# recurring regression where a stale working tree or BuildFlow daemon
# commits a *_templ.go that doesn't match the .templ source.
#
# Mirrors the logic of utils/templ_sync_test.go (TestTemplGeneratedInSync)
# but runs in <100ms so it fits within BuildFlow's pre-commit budget.
#
# Checks: every import in the .templ source must also appear in the
# generated _templ.go file. Excludes github.com/a-h/templ runtime imports
# (the generator always injects these).
#
# Usage:
#   scripts/check-templ-sync.sh          # exit 1 on drift
#   scripts/check-templ-sync.sh --quiet  # suppress output
set -euo pipefail

QUIET="${1:-}"

# Packages with .templ files (same list as TestTemplGeneratedInSync).
PACKAGES=(
	display errorpage feedback forms
	htmx layout navigation recipes
	charts/echarts datastar icons integration
)

DRIFT_FOUND=0
DRIFT_OUTPUT=""

for pkg in "${PACKAGES[@]}"; do
	# Find all .templ files in this package.
	for templ_file in "$pkg"/*.templ; do
		[ -f "$templ_file" ] || continue

		gen_file="${templ_file%.templ}_templ.go"

		if [ ! -f "$gen_file" ]; then
			DRIFT_OUTPUT+="  MISSING: $gen_file (not generated yet)\n"
			DRIFT_FOUND=1
			continue
		fi

		# Extract import paths from .templ source (lines inside import blocks
		# or single-line imports). Match quoted paths, excluding a-h/templ.
		while IFS= read -r imp; do
			[ -z "$imp" ] && continue

			# Skip templ runtime imports (generator always injects these).
			case "$imp" in
			github.com/a-h/templ*) continue ;;
			esac

			# Check if this import appears in the generated file.
			if ! grep -qF "\"$imp\"" "$gen_file" 2>/dev/null; then
				base_templ=$(basename "$templ_file")
				base_gen=$(basename "$gen_file")
				DRIFT_OUTPUT+="  DRIFT: $base_templ imports \"$imp\" but $base_gen does not\n"
				DRIFT_FOUND=1
			fi
		done < <(
			# Extract import paths from .templ file:
			# - Single-line: import "path"
			# - Block: lines with "path" inside import ( )
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
			' "$templ_file"
		)
	done
done

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
