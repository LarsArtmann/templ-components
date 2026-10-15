#!/usr/bin/env bash
# Build the templ-components website (Go SSG + Tailwind CSS v4).
# Output: dist/ (Firebase Hosting "public" directory).
#
# Star-badge resolution (inverted 2026-09-23, backlog #271; binary default
# flipped 2026-10-15 so even a bare `go run ./cmd/site` is deterministic):
#   default            -> --skip-stars IS the cmd/site default (deterministic
#                         fallback badge; a live count drifts between dist
#                         rebuilds and flakes every route golden containing
#                         the badge — the default is the SAFE build for tests,
#                         local runs, and tools)
#   SITE_LIVE_STARS=1  -> live GitHub count via --live-stars (production
#                         deploys: the CI Website job sets this for the
#                         artifact it deploys)
#   SITE_SKIP_STARS=1  -> legacy alias for the default (still honored)
set -euo pipefail
cd "$(dirname "$0")"

rm -rf dist

stars_flag=()
if [ "${SITE_LIVE_STARS:-0}" = "1" ] && [ "${SITE_SKIP_STARS:-0}" != "1" ]; then
	stars_flag=(--live-stars)
fi

GOEXPERIMENT=jsonv2 go run ./cmd/site --out dist --repo-root .. "${stars_flag[@]}"

# Tailwind v4 CLI (available in the nix devShell; CI installs
# @tailwindcss/cli). @source directives in site.css scan the library's
# generated components and this site's own.
tailwindcss -i site.css -o dist/assets/app.css --minify

echo "website: dist/ ready"
