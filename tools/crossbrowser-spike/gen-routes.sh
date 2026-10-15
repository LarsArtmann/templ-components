#!/usr/bin/env bash
# Regenerates routes.json from examples/demo/pages.go — the demo route
# registry stays the single source of truth; this extracts its Path fields.
set -euo pipefail
cd "$(dirname "$0")/../../.."
rg -o 'Path:\s+"([^"]+)"' -r '$1' examples/demo/pages.go | jq -Rn '{path: .}' \
  | jq -c '[inputs.path | {path: .}]' > tools/crossbrowser-spike/routes.json
echo "routes.json: $(rg -c 'path' tools/crossbrowser-spike/routes.json) routes"
