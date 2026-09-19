#!/usr/bin/env bash
# Build the templ-components website (Go SSG + Tailwind CSS v4).
# Output: dist/ (Firebase Hosting "public" directory).
#
# SITE_SKIP_STARS=1 builds with the no-stars fallback badge instead of the
# live GitHub count. The visual regression tier (flake app .#visual) sets it:
# a live count drifts between dist rebuilds and flakes every route golden
# that contains the badge.
set -euo pipefail
cd "$(dirname "$0")"

rm -rf dist

stars_flag=()
if [ "${SITE_SKIP_STARS:-0}" = "1" ]; then
	stars_flag=(--skip-stars)
fi

GOEXPERIMENT=jsonv2 go run ./cmd/site --out dist --repo-root .. "${stars_flag[@]}"

# Tailwind v4 CLI (available in the nix devShell; CI installs
# @tailwindcss/cli). @source directives in site.css scan the library's
# generated components and this site's own.
tailwindcss -i site.css -o dist/assets/app.css --minify

echo "website: dist/ ready"
