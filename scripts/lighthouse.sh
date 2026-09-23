#!/usr/bin/env bash
# Lighthouse lane skeleton (backlog #272). EXERCISES the built site dist and
# writes HTML/JSON reports to /tmp — NOT wired into CI (needs Node + a stable
# Chrome; CI's visual lane uses Nix Chromium). Run locally:
#
#   nix develop -c bash scripts/lighthouse.sh            # builds dist if missing
#
# Chosen shape (decision record for #272):
#   - Target the STATIC dist (not the demo): the site is the shipped product
#     surface and has no server-side state to confound scores.
#   - Categories: performance, accessibility, best-practices, SEO. The a11y
#     gate proper stays axe-based (visualtest/axe_sweep_test.go) — Lighthouse
#     a11y is advisory here, a smoke for regressions axe cannot see (audit
#     ordering, viewport-emulation differences).
#   - Thresholds: asserted manually at first (reports in /tmp/lighthouse);
#     promote to assert=true budgets in CI only after 5 consecutive stable
#     local runs, mirroring the dup-gate advisory->blocking ladder (#295).
set -euo pipefail
cd "$(dirname "$0")/.."

DIST="website/dist"
if [ ! -f "$DIST/index.html" ]; then
	echo "dist missing — building (skip-stars default keeps scores deterministic)"
	nix develop -c bash website/build.sh
fi

PORT="${LIGHTHOUSE_PORT:-8903}"
npx --yes serve "$DIST" -l "$PORT" &
SERVE_PID=$!
trap 'kill "$SERVE_PID" 2>/dev/null || true' EXIT
sleep 1

mkdir -p /tmp/lighthouse
for route in / /sales; do
	name="${route//\//_}"
	npx --yes lighthouse "http://localhost:${PORT}${route}" \
		--output=html --output=json \
		--output-path="/tmp/lighthouse/report${name}" \
		--only-categories=performance,accessibility,best-practices,seo \
		--chrome-flags="--headless=new --no-sandbox"
done

echo "reports in /tmp/lighthouse/ — review scores manually; budgets are not asserted yet"
