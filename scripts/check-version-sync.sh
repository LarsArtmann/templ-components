#!/usr/bin/env bash
# Fast (<50ms) pre-commit guard: blocks commits if the version string has
# drifted across the three files that must always agree:
#   1. utils/version.go   (const Version = "X.Y.Z")
#   2. CHANGELOG.md       (## [X.Y.Z] — latest released heading)
#   3. FEATURES.md        (**Version:** X.Y.Z)
#
# Mirrors the logic of TestVersionMatchesChangelog + TestVersionMatchesFeatures
# but runs in <50ms so it fits within BuildFlow's pre-commit budget.
#
# Usage:
#   scripts/check-version-sync.sh          # exit 1 on drift
#   scripts/check-version-sync.sh --quiet  # suppress output
set -euo pipefail

QUIET="${1:-}"
REPO_ROOT=$(git rev-parse --show-toplevel 2>/dev/null || echo ".")

# --- Extract version from each source ---

VERSION_GO=$(sed -n 's/.*const Version = "\([0-9][0-9.]*\)".*/\1/p' \
	"$REPO_ROOT/utils/version.go" | head -1)

VERSION_CHANGELOG=$(sed -n 's/^## \[\([0-9][0-9.]*\)].*/\1/p' \
	"$REPO_ROOT/CHANGELOG.md" | head -1)

VERSION_FEATURES=$(sed -n 's/.*\*\*Version:\*\*\s*\([0-9][0-9.]*\).*/\1/p' \
	"$REPO_ROOT/FEATURES.md" | head -1)

# --- Validate extraction ---

ERRORS=""

if [ -z "$VERSION_GO" ]; then
	ERRORS+="  Cannot extract version from utils/version.go\n"
fi
if [ -z "$VERSION_CHANGELOG" ]; then
	ERRORS+="  Cannot extract version from CHANGELOG.md (expected '## [X.Y.Z]')\n"
fi
if [ -z "$VERSION_FEATURES" ]; then
	ERRORS+="  Cannot extract version from FEATURES.md (expected '**Version:** X.Y.Z')\n"
fi

if [ -n "$ERRORS" ]; then
	if [ "$QUIET" != "--quiet" ]; then
		echo "" >&2
		echo "BLOCKED: Version extraction failed." >&2
		echo "" >&2
		echo -e "$ERRORS" >&2
		echo "" >&2
	fi
	exit 1
fi

# --- Compare ---

MISMATCH=0
OUTPUT=""

if [ "$VERSION_GO" != "$VERSION_CHANGELOG" ]; then
	OUTPUT+="  utils/version.go says $VERSION_GO but CHANGELOG.md says $VERSION_CHANGELOG\n"
	MISMATCH=1
fi

if [ "$VERSION_GO" != "$VERSION_FEATURES" ]; then
	OUTPUT+="  utils/version.go says $VERSION_GO but FEATURES.md says $VERSION_FEATURES\n"
	MISMATCH=1
fi

if [ "$MISMATCH" -ne 0 ]; then
	if [ "$QUIET" != "--quiet" ]; then
		echo "" >&2
		echo "BLOCKED: Version numbers are out of sync across files." >&2
		echo "" >&2
		echo -e "$OUTPUT" >&2
		echo "Fix: bump all three together — utils/version.go, CHANGELOG.md, FEATURES.md." >&2
		echo "Or run: scripts/release.sh <new-version> \"<summary>\"" >&2
		echo "" >&2
	fi
	exit 1
fi

exit 0
