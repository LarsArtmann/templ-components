#!/usr/bin/env bash
# check-module-layers.sh — DAG enforcement for the 7-module workspace (ADR-0034)
#
# Verifies that no module imports from a higher layer (upward dependency).
# The DAG is:
#   Layer 0: utils           (leaf)
#   Layer 1: icons, charts/echarts, datastar, htmx  (depend on utils)
#   Layer 2: errorpage       (depends on utils, icons)
#   Layer 3: root            (depends on all above)
#
# An upward dependency (e.g., utils importing display) would create a cycle
# or break the standalone-build invariant. This script catches it at commit time.

set -euo pipefail
REPO_ROOT=$(git rev-parse --show-toplevel 2>/dev/null || echo ".")
PREFIX="github.com/larsartmann/templ-components"

errors=0

# check_layer <module_dir> <allowed_subpaths...>
# Scans all .go files (excluding generated *_templ.go) in <module_dir> for
# imports of ${PREFIX}/X where X is NOT in the allowed list.
# Allowed subpaths are matched as path prefixes (e.g., "utils" matches "utils/svg").
check_layer() {
	local dir="$1"
	shift
	local -a allowed=("$@")

	while IFS= read -r file; do
		[[ -z "$file" ]] && continue
		# Extract all imports matching our module prefix
		local imports
		imports=$(grep -oE "\"${PREFIX}/[^\"]+\"" "$file" 2>/dev/null | sed 's/"//g' | sort -u || true)

		for imp in $imports; do
			# Strip the prefix to get the subpath
			local subpath="${imp#${PREFIX}/}"

			local ok=false
			for allowed_prefix in "${allowed[@]}"; do
				if [[ "$subpath" == "$allowed_prefix" ]] || [[ "$subpath" == "$allowed_prefix/"* ]]; then
					ok=true
					break
				fi
			done

			if [[ "$ok" != "true" ]]; then
				echo "DAG VIOLATION: $dir imports '$imp' (upward dependency)"
				echo "  File: ${file#$REPO_ROOT/}"
				errors=$((errors + 1))
			fi
		done
	done < <(find "$REPO_ROOT/$dir" -name '*.go' ! -name '*_templ.go' ! -path '*/vendor/*' 2>/dev/null)
}

# Layer 0: utils — leaf, can only import its own sub-packages
check_layer "utils" "utils"

# Layer 1: icons — can import utils (+ self)
check_layer "icons" "utils" "icons"

# Layer 1: charts/echarts — can import utils (+ self)
check_layer "charts/echarts" "utils" "charts/echarts"

# Layer 1: datastar — can import utils (+ self)
check_layer "datastar" "utils" "datastar"

# Layer 1: htmx — can import utils (+ self)
check_layer "htmx" "utils" "htmx"

# Layer 2: errorpage — can import utils, icons (+ self)
check_layer "errorpage" "utils" "icons" "errorpage"

# Root module (layer 3) is not checked — it can import everything.

if [[ $errors -gt 0 ]]; then
	echo ""
	echo "Module layer check: FAILED ($errors violation(s))"
	echo "See ADR-0034 for the module dependency DAG."
	exit 1
fi

echo "Module layer check: OK (no upward dependencies in 6 sub-modules)"
