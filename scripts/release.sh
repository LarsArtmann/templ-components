#!/usr/bin/env bash
# scripts/release.sh — cut a templ-components release in one command.
#
# Usage:
#   scripts/release.sh <new-version> <release-summary>
#
# Example:
#   scripts/release.sh 0.7.0 "typed HTMX retry policies, Drawer motion-reduce"
#
# What it does:
#   1. Validates the working tree is clean and on master
#   2. Confirms the new version is greater than the current one
#   3. Bumps utils.Version
#   4. Inserts a new CHANGELOG heading (replaces [Unreleased], adds fresh [Unreleased] above)
#   5. Asks for the release notes (multi-line) and writes them into the CHANGELOG entry
#   6. Regenerates *_templ.go and runs the full verify suite
#   7. Stages and commits as `release: <version> — <summary>` (one-commit convention)
#   8. Creates an annotated, SSH-signed tag `v<version>: <summary>`
#
# Required: GPG/SSH signing key configured (the tag signing matches v0.5.0).
# Does NOT push. House rule: "NEVER PUSH TO REMOTE" — push manually after review.

set -euo pipefail

if [ $# -ne 2 ]; then
    echo "Usage: $0 <new-version> <release-summary>" >&2
    echo "Example: $0 0.7.0 'typed HTMX retry policies, Drawer motion-reduce'" >&2
    exit 1
fi

NEW_VERSION="$1"
RELEASE_SUMMARY="$2"
TODAY="$(date -u +%Y-%m-%d)"

cd "$(dirname "$0")/.."

# 1. Working tree must be clean.
if [ -n "$(git status --porcelain)" ]; then
    echo "Error: working tree is not clean. Commit or stash changes first." >&2
    git status --short
    exit 1
fi

# 2. Must be on master.
CURRENT_BRANCH="$(git rev-parse --abbrev-ref HEAD)"
if [ "$CURRENT_BRANCH" != "master" ]; then
    echo "Error: must be on master (currently on $CURRENT_BRANCH)." >&2
    exit 1
fi

# 3. New version must be > current version.
CURRENT_VERSION="$(grep -E '^const[[:space:]]+Version' utils/version.go | sed -E 's/.*"([^"]+)".*/\1/')"
echo "Current version: $CURRENT_VERSION"
echo "New version:     $NEW_VERSION"
if [ "$NEW_VERSION" = "$CURRENT_VERSION" ]; then
    echo "Error: new version is identical to current version." >&2
    exit 1
fi

# Use sort -V to check ordering. New > current.
SORTED_LOWER="$(printf '%s\n%s\n' "$CURRENT_VERSION" "$NEW_VERSION" | sort -V | head -n1)"
if [ "$SORTED_LOWER" != "$CURRENT_VERSION" ]; then
    echo "Error: new version ($NEW_VERSION) is not greater than current ($CURRENT_VERSION)." >&2
    exit 1
fi

# 4. Verify [Unreleased] has content (not just an empty placeholder).
#    Look for the first non-empty, non-heading line after ## [Unreleased].
UNRELEASED_BODY="$(awk '/^## \[Unreleased\]$/ {found=1; next} found && /^## \[/ {exit} found && /^### / {has_heading=1; next} found && has_heading && NF > 0 {print; exit}' CHANGELOG.md)"
if [ -z "$UNRELEASED_BODY" ]; then
    echo "Error: [Unreleased] section in CHANGELOG.md is empty." >&2
    echo "Add changelog entries to [Unreleased] before cutting a release." >&2
    exit 1
fi
echo "[Unreleased] section has content — proceeding."

# 5. Bump utils.Version.
sed -i.bak -E "s|^(const[[:space:]]+Version[[:space:]]+=[[:space:]]+\")[^\"]+(\")|\1${NEW_VERSION}\2|" utils/version.go
rm -f utils/version.go.bak
echo "Bumped utils.Version to $NEW_VERSION"

# 6. Insert CHANGELOG heading. The pattern is:
#      ## [Unreleased]
#      <empty>
#      ## [OLD_VERSION] — <old-date>
#      ...
# We want:
#      ## [Unreleased]
#      <empty>
#      ## [NEW_VERSION] — <TODAY>
#      ### Added
#      ### Changed
#      ### Fixed
#      (release notes inserted below)
#      ...
#      ## [OLD_VERSION] — <old-date>
#      ...

echo "Collecting release notes. Press Ctrl-D on an empty line to finish."
echo "(Markdown formatting OK. Tip: use '### Added', '### Changed', '### Fixed'.)"
RELEASE_NOTES=""
while IFS= read -r line; do
    [ -z "$line" ] && break
    RELEASE_NOTES="${RELEASE_NOTES}${line}"$'\n'
done

# Replace the first occurrence of '## [Unreleased]' with new structure.
CHANGELOG_TMP="$(mktemp)"
{
    awk '/^## \[Unreleased\]$/ {print; print ""; in_unreleased=1; next} in_unreleased && /^$/ {print; printf "## [%s] — %s\n\n", "'"$NEW_VERSION"'", "'"$TODAY"'"; print RELEASE_NOTES; print ""; in_unreleased=0; next} {print}' RELEASE_NOTES="$RELEASE_NOTES" CHANGELOG.md
} > "$CHANGELOG_TMP"
mv "$CHANGELOG_TMP" CHANGELOG.md
echo "Updated CHANGELOG.md with $NEW_VERSION heading"

# 7. Run full verify.
echo "Running full verify (templ generate + build + test + lint)..."
find . -name '*_templ.go' -print0 | xargs -0 rm -f
templ generate ./...
go build ./...
go test ./... -count=1 -race
golangci-lint run ./display/... ./errorpage/... ./feedback/... ./forms/... ./htmx/... ./icons/... ./layout/... ./navigation/... ./utils/... ./internal/...

# Drift-guard test must pass (utils.Version matches CHANGELOG first non-Unreleased).
if ! go test ./utils/... -run TestVersionMatchesChangelog -count=1 >/dev/null 2>&1; then
    echo "Error: version drift-guard test failed. CHANGELOG heading does not match utils.Version." >&2
    git checkout -- utils/version.go CHANGELOG.md
    exit 1
fi

# 8. Stage and commit.
git add utils/version.go CHANGELOG.md
git add -u  # any verified updates

RELEASE_BODY="${RELEASE_SUMMARY}

This release rolls up all changes since v${CURRENT_VERSION}.

💘 Generated with Crush

Assisted-by: Crush:MiniMax-M3"

git commit -m "release: ${NEW_VERSION} — ${RELEASE_SUMMARY}

${RELEASE_BODY}

Co-Authored-By: Crush <noreply@crush.lars.software>"

RELEASE_COMMIT="$(git rev-parse HEAD)"

# 9. Annotated, SSH-signed tag.
git tag -s "v${NEW_VERSION}" -m "v${NEW_VERSION}: ${RELEASE_SUMMARY}" "$RELEASE_COMMIT"

echo ""
echo "Release v${NEW_VERSION} cut at commit ${RELEASE_COMMIT}."
echo "Tag: v${NEW_VERSION} (annotated, SSH-signed)"
echo ""
echo "Next steps:"
echo "  1. Review the release: git show v${NEW_VERSION}"
echo "  2. Push (when ready):  git push origin master --follow-tags"
echo "  3. House rule says NEVER PUSH TO REMOTE — confirm with the user first."
