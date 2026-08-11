#!/usr/bin/env bash
# Pre-commit hook for templ-components
# Runs: templ generate, per-module build/test/lint across all 7 modules.
#
# GOWORK=off is intentional: it tests the replace-directive resolution path
# (the same path proxy consumers use at publish time). Workspace mode (go.work)
# is for interactive development only.

set -euo pipefail

export GOWORK=off
export GOEXPERIMENT=jsonv2

echo "Running templ-components pre-commit checks..."

# Remove stale generated files and regenerate
find . -name '*_templ.go' -print0 | xargs -0 rm -f
templ generate ./...

# --- Root module (uses replace directives for sub-modules) ---
go build ./...
go test ./... -count=1

# --- Sub-modules (standalone isolation) ---
for mod in utils icons errorpage charts/echarts datastar htmx; do
	echo "==> $mod"
	(cd "$mod" && go build ./... && go test ./... -count=1)
done

# --- Lint (golangci-lint does not support go.work) ---
echo "==> lint root"
golangci-lint run \
	./display/... ./feedback/... ./forms/... \
	./integration/... ./internal/... \
	./layout/... ./navigation/... ./recipes/... ./cmd/...

for mod in utils icons errorpage charts/echarts datastar htmx; do
	echo "==> lint $mod"
	(cd "$mod" && golangci-lint run ./...)
done

echo "All checks passed."
