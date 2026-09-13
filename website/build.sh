#!/usr/bin/env bash
# Build the templ-components website (Go SSG + Tailwind CSS v4).
# Output: dist/ (Firebase Hosting "public" directory).
set -euo pipefail
cd "$(dirname "$0")"

rm -rf dist

GOEXPERIMENT=jsonv2 go run ./cmd/site --out dist --repo-root ..

# Tailwind v4 CLI (available in the nix devShell; CI installs
# @tailwindcss/cli). @source directives in site.css scan the library's
# generated components and this site's own.
tailwindcss -i site.css -o dist/assets/app.css --minify

echo "website: dist/ ready"
