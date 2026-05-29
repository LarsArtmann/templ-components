#!/usr/bin/env bash
# Pre-commit hook for templ-components
# Runs: templ generate, go build, go test, golangci-lint

set -euo pipefail

export GOWORK=off

echo "Running templ-components pre-commit checks..."

# Remove stale generated files and regenerate
find . -name '*_templ.go' -print0 | xargs -0 rm -f
templ generate ./...

# Build
go build ./...

# Test
go test ./...

# Lint (exclude examples/ per .golangci.yml)
golangci-lint run ./display/... ./errorpage/... ./feedback/... ./forms/... ./htmx/... ./icons/... ./layout/... ./navigation/... ./utils/... ./internal/...

echo "All checks passed."
