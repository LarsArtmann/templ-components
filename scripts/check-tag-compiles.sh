#!/usr/bin/env bash
# scripts/check-tag-compiles.sh — prove a RELEASED TAG compiles for consumers.
#
# WHY: the Go module proxy serves the git tag's tree AS-IS (no `templ
# generate`, no workspace, no replaces). The one catastrophic failure mode —
# a tag missing *_templ.go or otherwise uncompilable — is currently caught
# only by consumers. This script reproduces the consumer experience exactly:
# a throwaway module, `go get` of the tag from the proxy, `go build` of every
# import path (root + all 6 sub-modules).
#
# Usage:
#   scripts/check-tag-compiles.sh            # latest root tag (git describe)
#   scripts/check-tag-compiles.sh v1.17.0    # explicit tag
#
# Exit 0 = a fresh consumer can compile every package at that tag.
# Run locally after pushing tags, and in CI via the post-release
# workflow_dispatch job (TODO #208 / Pareto M03).

set -euo pipefail

TAG="${1:-}"
if [ -z "$TAG" ]; then
	TAG="$(git describe --tags --abbrev=0 --match 'v*' 2>/dev/null)" || {
		echo "Error: no v* tag found and no tag argument given." >&2
		exit 1
	}
fi

echo "==> tag under test: $TAG"

WORKDIR="$(mktemp -d)"
trap 'rm -rf "$WORKDIR"' EXIT

cd "$WORKDIR"
go mod init "tc-tag-smoke" >/dev/null

# Consumer-realistic flags: proxy first, and jsonv2 (the library requires it).
export GOEXPERIMENT=jsonv2
export GOPROXY="${GOPROXY:-https://proxy.golang.org,direct}"

# Every import path a consumer can reach. `go get` resolves the module and
# its transitive requirements from the proxy — exactly what a fresh consumer
# build does. Sub-modules use directory-prefixed tags.
PATHS=(
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/navigation"
	"github.com/larsartmann/templ-components/recipes"
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/errorpage"
	"github.com/larsartmann/templ-components/charts/echarts"
	"github.com/larsartmann/templ-components/htmx"
	"github.com/larsartmann/templ-components/datastar"
)

SUBMODULES=(
	"utils"
	"icons"
	"errorpage"
	"charts/echarts"
	"datastar"
	"htmx"
)

echo "==> go get root module at ${TAG}"
go get "github.com/larsartmann/templ-components@${TAG}" >/dev/null

echo "==> go get sub-modules at their directory-prefixed tags"
for sub in "${SUBMODULES[@]}"; do
	go get "github.com/larsartmann/templ-components/${sub}@${TAG}" >/dev/null
done

FAILURES=0
for path in "${PATHS[@]}"; do
	if go build "$path"; then
		echo "ok: $path"
	else
		echo "FAIL: $path" >&2
		FAILURES=$((FAILURES + 1))
	fi
done

if [ "$FAILURES" -gt 0 ]; then
	echo "Error: $FAILURES import path(s) failed to compile at ${TAG} — the proxy is now serving a broken release; see docs/release-checklist.md (never retag; cut a patch)." >&2
	exit 1
fi

echo "All import paths compile at ${TAG}. Consumers are safe."
