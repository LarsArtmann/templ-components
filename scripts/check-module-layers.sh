#!/usr/bin/env bash
# check-module-layers.sh — DAG enforcement for the 7-module workspace (ADR-0034)
#
# Verifies that no module imports from a higher layer (upward dependency).
# The DAG is:
#   Layer 0: utils               (leaf)
#   Layer 1: icons, charts/echarts, datastar, htmx   (depend on utils)
#   Layer 2: errorpage           (depends on utils, icons)
#   Layer 3: root                (depends on all above)
#   Layer 4: visualtest          (consumer/test module — depends on all above;
#                                 nothing may depend on it. Its sibling pins
#                                 are local replace directives, so an upward
#                                 import here breaks standalone builds, same
#                                 as for the published modules.)
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

# Root module (layer 3) is not checked per-package — it can import everything.

# Layer 4: visualtest — the consumer/test module. Explicitly modeled (it
# previously passed only by skip): it may import every library module but is
# imported by none. The allow-list is exhaustive on purpose — a new import
# outside it (e.g. a stray cmd/ dependency) fails here and forces a
# conscious DAG decision.
#
# Decision 2026-09-13: "examples" is allowed. visualtest/demo.go builds and
# runs the examples/demo binary as the live e2e server (local replace only,
# never published); examples/demo is a root-module package on the same layer
# as display/forms/... already listed, and imports nothing from visualtest
# (verified — no cycle). First surfaced when the tracked .githooks/pre-commit
# made manual commits run this guard again.
check_layer "visualtest" \
	"utils" "icons" "errorpage" "charts/echarts" "datastar" "htmx" \
	"display" "feedback" "forms" "layout" "navigation" "recipes" \
	"examples" \
	"visualtest"

# Layer 4: website — the consumer/site module (static-site generator for
# templcomponents.lars.software; renders through the library's own
# components; local replace only, never published).
check_layer "website" \n "utils" "icons" "errorpage" "charts/echarts" "datastar" "htmx" \n "display" "feedback" "forms" "layout" "navigation" "recipes" \n "website"

if [[ $errors -gt 0 ]]; then
	echo ""
	echo "Module layer check: FAILED ($errors violation(s))"
	echo "See ADR-0034 for the module dependency DAG."
	exit 1
fi

echo "Module layer check: OK (no upward dependencies in 6 sub-modules + visualtest/website consumers)"
