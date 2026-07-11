#!/usr/bin/env bash
# Pre-commit hook for templ-components
# Runs: templ generate, go build, go test, golangci-lint

set -euo pipefail

export GOWORK=off
export GOEXPERIMENT=jsonv2

echo "Running templ-components pre-commit checks..."

# Remove stale generated files and regenerate
find . -name '*_templ.go' -print0 | xargs -0 rm -f
templ generate ./...

# Build
go build ./...

# Test
go test ./...

# Lint (examples/ excluded via .golangci.yml)
golangci-lint run ./...

echo "All checks passed."
