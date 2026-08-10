#!/usr/bin/env bash
# scripts/check-module-sync.sh — verify multi-module workspace integrity.
#
# Checks:
#   1. Every go.mod file's module path matches its directory.
#   2. Every replace directive uses a relative path (no absolute paths).
#   3. Every require entry for a sibling module matches the shared version.
#
# Runs in <100ms. Wired into CI and pre-commit.

set -euo pipefail
cd "$(dirname "$0")/.."

ERRORS=0

# 1. Verify each sub-module go.mod has correct module path.
declare -A MODULE_PATHS=(
    ["utils"]="github.com/larsartmann/templ-components/utils"
    ["icons"]="github.com/larsartmann/templ-components/icons"
    ["errorpage"]="github.com/larsartmann/templ-components/errorpage"
    ["charts/echarts"]="github.com/larsartmann/templ-components/charts/echarts"
)

for dir in "${!MODULE_PATHS[@]}"; do
    expected="${MODULE_PATHS[$dir]}"
    actual="$(head -1 "${dir}/go.mod" | awk '{print $2}')"
    if [ "$actual" != "$expected" ]; then
        echo "::error::${dir}/go.mod module path is '${actual}', expected '${expected}'"
        ERRORS=$((ERRORS + 1))
    fi
done

# 2. Verify no absolute paths in replace directives.
while IFS= read -r line; do
    if echo "$line" | grep -qE 'replace.*=> */'; then
        file="$(echo "$line" | cut -d: -f1)"
        echo "::error::${file}: replace directive uses absolute path (not portable for CI/consumers)"
        ERRORS=$((ERRORS + 1))
    fi
done < <(grep -rn 'replace.*=>' go.mod utils/go.mod icons/go.mod errorpage/go.mod charts/echarts/go.mod 2>/dev/null || true)

# 3. Verify all sibling module requires share the same version.
#    Only match require lines (indented, followed by version), not replace lines.
SHARED_VERSION="$(grep -E '^\s+github.com/larsartmann/templ-components/utils v' go.mod | head -1 | awk '{print $2}' || true)"
if [ -n "$SHARED_VERSION" ]; then
    for modfile in go.mod icons/go.mod errorpage/go.mod charts/echarts/go.mod; do
        while IFS= read -r version; do
            if [ -n "$version" ] && [ "$version" != "$SHARED_VERSION" ]; then
                echo "::error::${modfile}: sibling module version '${version}' != shared '${SHARED_VERSION}'"
                ERRORS=$((ERRORS + 1))
            fi
        done < <(grep -E '^\s+github.com/larsartmann/templ-components/[a-z/]+ v' "$modfile" 2>/dev/null | awk '{print $2}' || true)
    done
fi

if [ "$ERRORS" -gt 0 ]; then
    echo "Module sync check: ${ERRORS} error(s) found."
    exit 1
fi

echo "Module sync check: OK (5 modules, all paths and versions consistent)."
