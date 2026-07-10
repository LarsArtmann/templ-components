#!/usr/bin/env bash
# Pre-commit hook for templ-components
# Runs: templ generate, go build, go test, golangci-lint

set -euo pipefail

export GOWORK=off

echo "Running templ-components pre-commit checks..."

# Guard: encoding/json/v2 is not yet stable (Go 1.27+) and breaks builds.
# Auto-formatters running under GOEXPERIMENT=jsonv2 can rewrite imports.
if grep -rn 'encoding/json/v2' --include='*.go' .; then
    echo "ERROR: encoding/json/v2 import detected. This package is not stable yet (blocked on Go 1.27)."
    echo "Use encoding/json (v1) instead."
    exit 1
fi

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
