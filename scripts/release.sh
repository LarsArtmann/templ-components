#!/usr/bin/env bash
# scripts/release.sh — cut a templ-components release in one command.
#
# Usage:
#   scripts/release.sh <new-version> <release-summary> [--notes-file FILE]
#
# Examples:
#   scripts/release.sh 0.7.0 "typed HTMX retry policies, Drawer motion-reduce"
#   scripts/release.sh 0.7.0 "typed HTMX retry" --notes-file /tmp/release-notes.md
#
# Release notes source (first found wins):
#   --notes-file FILE   read notes from FILE (markdown)
#   (default)           extract from CHANGELOG.md [Unreleased] section
#
# What it does:
#   1. Validates the working tree is clean and on master
#   2. Confirms the new version is greater than the current one
#   3. Bumps utils.Version
#   4. Collects release notes (--notes-file or CHANGELOG [Unreleased])
#   5. Moves notes from [Unreleased] to a new versioned heading (inserts fresh [Unreleased])
#   6. Regenerates *_templ.go and runs the full verify suite
#   7. Stages and commits as `release: <version> — <summary>` (one-commit convention)
#   8. Creates an annotated, SSH-signed tag `v<version>: <summary>`
#
# Required: GPG/SSH signing key configured (the tag signing matches v0.5.0).
# Does NOT push. House rule: "NEVER PUSH TO REMOTE" — push manually after review.

set -euo pipefail

NEW_VERSION=""
RELEASE_SUMMARY=""
NOTES_FILE=""

while [ $# -gt 0 ]; do
    case "$1" in
        --notes-file)
            if [ $# -lt 2 ]; then
                echo "Error: --notes-file requires a path argument." >&2
                exit 1
            fi
            NOTES_FILE="$2"
            shift 2
            ;;
        --help|-h)
            echo "Usage: $0 <new-version> <release-summary> [--notes-file FILE]"
            echo "Example: $0 0.7.0 'typed HTMX retry policies, Drawer motion-reduce'"
            echo ""
            echo "Release notes source (first found wins):"
            echo "  --notes-file FILE   read notes from FILE (markdown)"
            echo "  (default)           extract from CHANGELOG.md [Unreleased] section"
            exit 0
            ;;
        -*)
            echo "Error: unknown flag: $1" >&2
            exit 1
            ;;
        *)
            if [ -z "$NEW_VERSION" ]; then
                NEW_VERSION="$1"
            elif [ -z "$RELEASE_SUMMARY" ]; then
                RELEASE_SUMMARY="$1"
            else
                echo "Error: unexpected positional argument: $1" >&2
                exit 1
            fi
            shift
            ;;
    esac
done

if [ -z "$NEW_VERSION" ] || [ -z "$RELEASE_SUMMARY" ]; then
    echo "Usage: $0 <new-version> <release-summary> [--notes-file FILE]" >&2
    echo "Example: $0 0.7.0 'typed HTMX retry policies, Drawer motion-reduce'" >&2
    exit 1
fi
TODAY="$(date -u +%Y-%m-%d)"

export GOWORK=off
export GOEXPERIMENT=jsonv2

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

# 4. Collect release notes.
#    Source priority: --notes-file > CHANGELOG [Unreleased] body.
#    Project rule: "[Unreleased] must be warm at all times" — the notes already
#    live in CHANGELOG, so we extract them rather than forcing the user to retype
#    them into a hostile stdin prompt with no editing or file input.
if [ -n "$NOTES_FILE" ]; then
    if [ ! -f "$NOTES_FILE" ]; then
        echo "Error: --notes-file '$NOTES_FILE' does not exist." >&2
        exit 1
    fi
    RELEASE_NOTES="$(cat "$NOTES_FILE")"
    echo "Using release notes from $NOTES_FILE"
else
    RELEASE_NOTES="$(awk '
        /^## \[Unreleased\]$/ { unreleased=1; next }
        unreleased && /^## \[/ { exit }
        unreleased { print }
    ' CHANGELOG.md)"
    if [ -z "$RELEASE_NOTES" ]; then
        echo "Error: [Unreleased] section in CHANGELOG.md is empty." >&2
        echo "Add changelog entries to [Unreleased] before cutting a release," >&2
        echo "or pass --notes-file FILE with the release notes." >&2
        exit 1
    fi
    echo "Extracted release notes from CHANGELOG.md [Unreleased] section."
fi

# Trim leading/trailing blank lines for clean commit body + CHANGELOG formatting.
RELEASE_NOTES="$(printf '%s\n' "$RELEASE_NOTES" | awk 'NF{p=1} p{lines[++n]=$0} END{while(n>0 && lines[n]~/^[[:space:]]*$/) n--; for(i=1;i<=n;i++) print lines[i]}')"

# 5. Bump utils.Version.
sed -i.bak -E "s|^(const[[:space:]]+Version[[:space:]]+=[[:space:]]+\")[^\"]+(\")|\1${NEW_VERSION}\2|" utils/version.go
rm -f utils/version.go.bak
echo "Bumped utils.Version to $NEW_VERSION"

# 6. Move release notes from [Unreleased] to the new version heading.
#    On encountering [Unreleased]: emit fresh-empty [Unreleased], then the new
#    version heading with the (trimmed) notes body. Skip the ORIGINAL [Unreleased]
#    body until the next ## [ heading so notes are not duplicated.
CHANGELOG_TMP="$(mktemp)"
awk -v NEW_VERSION="$NEW_VERSION" -v TODAY="$TODAY" -v RELEASE_NOTES="$RELEASE_NOTES" '
    /^## \[Unreleased\]$/ {
        print; print ""
        printf "## [%s] — %s\n\n", NEW_VERSION, TODAY
        print RELEASE_NOTES
        print ""
        skip=1
        next
    }
    skip && /^## \[/ { skip=0 }
    skip { next }
    { print }
' CHANGELOG.md > "$CHANGELOG_TMP"
mv "$CHANGELOG_TMP" CHANGELOG.md
echo "Updated CHANGELOG.md: moved [Unreleased] body under [${NEW_VERSION}] heading."

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

# Commit body = the release notes (multi-paragraph), NOT a duplicate of the
# one-line summary. The subject already carries the summary; the body carries
# the detail. Model attribution is parameterized so the script works under any
# Crush model (export CRUSH_MODEL before invoking, else falls back to "unknown").
RELEASE_BODY="${RELEASE_NOTES}

💘 Generated with Crush

Assisted-by: Crush:${CRUSH_MODEL:-unknown}"

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
