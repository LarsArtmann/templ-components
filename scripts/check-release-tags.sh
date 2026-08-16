#!/usr/bin/env bash
# scripts/check-release-tags.sh — verify a root release tag shipped lockstep
# with every published sub-module tag.
#
# Usage:
#   scripts/check-release-tags.sh [version]    # default: highest root tag
#
# Why this exists: the v1.8.3 release pushed only the root tag. Consumers
# (dnsblockd et al.) pin sub-modules like errorpage and htmx in go.mod; a
# root tag without sibling sub-module tags leaves those pins unresolvable on
# the Go module proxy and breaks every dependent build with "unknown
# revision". The release script now derives the sub-module set dynamically;
# this guard is the belt to that suspenders — run it before pushing tags.
#
# Fails when any published sub-module lacks "<path>/v<version>" or when such
# a tag points at a different commit than the root tag (all release tags are
# created on the same release commit by scripts/release.sh).
set -euo pipefail

cd "$(dirname "$0")/.."

if [ ! -f go.mod ]; then
	echo "Error: run from the repository root (go.mod not found)." >&2
	exit 1
fi

VERSION="${1:-}"
if [ -z "$VERSION" ]; then
	VERSION="$(git tag --list 'v*' | sed 's/^v//' | sort -V | tail -n1)"
	if [ -z "$VERSION" ]; then
		echo "Error: no root tags found; pass a version explicitly." >&2
		exit 1
	fi
	echo "Checking latest root tag: v${VERSION}"
else
	echo "Checking root tag: v${VERSION}"
fi

ROOT_TAG="v${VERSION}"
ROOT_COMMIT="$(git rev-parse -q --verify "${ROOT_TAG}^{commit}")"
if [ -z "$ROOT_COMMIT" ]; then
	echo "Error: root tag ${ROOT_TAG} does not exist." >&2
	exit 1
fi

# Published sub-modules: derived exactly like scripts/release.sh derives them
# (the root go.mod replace directives), so the two can never disagree.
SUBMODULE_PATHS="$(awk '
	$1 == "replace" && $2 ~ /^github\.com\/larsartmann\/templ-components\// {
		path = $2
		sub(/^github\.com\/larsartmann\/templ-components\//, "", path)
		print path
	}
' go.mod)"
if [ -z "$SUBMODULE_PATHS" ]; then
	echo "Error: no templ-components replace directives in go.mod — nothing to check." >&2
	exit 1
fi

FAILED=0
for submod in $SUBMODULE_PATHS; do
	tag="${submod}/${ROOT_TAG}"
	commit="$(git rev-parse -q --verify "refs/tags/${tag}^{commit}")"
	if [ -z "$commit" ]; then
		echo "MISSING: ${tag} — the root tag shipped without its sub-module tag." >&2
		echo "  Fix: git tag -s \"${tag}\" -m \"${tag}\" ${ROOT_COMMIT} (then push it)." >&2
		FAILED=1
	elif [ "$commit" != "$ROOT_COMMIT" ]; then
		echo "DIVERGED: ${tag} points at ${commit}, root tag at ${ROOT_COMMIT}." >&2
		echo "  Release tags must all sit on the release commit." >&2
		FAILED=1
	else
		echo "ok: ${tag}"
	fi
done

if [ "$FAILED" = "1" ]; then
	echo "Release tag set for v${VERSION} is incomplete — push blocked." >&2
	exit 1
fi

echo "All sub-module tags present for v${VERSION}."
