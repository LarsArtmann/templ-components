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
#   5b. Bumps require versions in every module go.mod (discovered, not hardcoded)
#   6. Moves release notes from [Unreleased] under the new version heading
#   6b. Bumps FEATURES.md version + date (drift-guarded)
#   7. Runs the full verify suite: templ generate + CSS + build/test/lint.
#       CRITICAL: runs with replace directives still in place — see the
#       ordering note below the env exports.
#   7b. Strips replace directives (tagged go.mod files must be consumer-clean)
#   8. Stages and commits as `release: <version> — <summary>` (one-commit convention)
#   9. Creates annotated, SSH-signed tags: root `v<x>` + sub-module `utils/v<x>`, etc.
#  10. Re-adds replace directives in a follow-up commit for local dev
#
# Required: GPG/SSH signing key configured (the tag signing matches v0.5.0).
# Does NOT push. House rule: "NEVER PUSH TO REMOTE" — push manually after review.
#
# NOTE (v1.10.0 incident): the auto-commit daemon raced this script mid-cut —
# it committed the in-flight version bumps, pushed the tags before the script
# finished, and BuildFlow recompiled the CSS distribution targets only AFTER
# the release commit was written. Two blemishes shipped in the v1.10.0 tag: a
# leftover root self-replace directive and compiled CSS one release stale.
# Both are functionally harmless for consumers (dependency replace directives
# are ignored; CSS compiles from source), but the multi-target CSS compile
# below and the step-8b pre-tag tree assertions are the hardening. Pushed tags
# are immutable — verify the release tree BEFORE creating them.

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
	--help | -h)
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

export GOEXPERIMENT=jsonv2
# NOTE: GOWORK is NOT set to off — we use go.work (workspace mode).
#
# CRITICAL ORDERING: verification runs BEFORE replace directives are stripped.
# go1.26.5 workspace mode does NOT preempt module-graph resolution of require
# entries at unpushed versions: with replaces removed and the new tags not yet
# on the proxy, `go build` fails with "unknown revision <sub>/v<version>" for
# every sibling require (proven during the v1.9.0 cut; GOPRIVATE is NOT a
# factor — verified with a 2x2 replaces/GOPRIVATE matrix). With replaces
# present, sibling resolution is purely local, so the full build/test/lint
# suite gates the release first; the strip happens afterwards so the tagged
# go.mod files stay consumer-clean.

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

# 2b. Published sub-modules, derived from the ROOT go.mod's replace
#     directives (replace github.com/larsartmann/templ-components/<path> => ./<path>)
#     - the single source of truth for "what gets tagged". Hardcoding this
#     list is how v1.8.3 shipped as a root-only release: the hardcoded set had
#     drifted from the real module set (missing htmx/ and datastar/), so
#     consumers pinning the sub-modules at the new version found no tags on
#     the proxy and every dependent build broke.
SUBMODULE_PATHS="$(awk '
	$1 == "replace" && $2 ~ /^github\.com\/larsartmann\/templ-components\// {
		path = $2
		sub(/^github\.com\/larsartmann\/templ-components\//, "", path)
		print path
	}
' go.mod)"
if [ -z "$SUBMODULE_PATHS" ]; then
	echo "Error: no templ-components replace directives found in go.mod." >&2
	echo "The release script expects the local-dev replace set to be present" >&2
	echo "(the previous release's follow-up commit re-adds them)." >&2
	exit 1
fi
MODFILES="go.mod"
for submod in $SUBMODULE_PATHS; do
	MODFILES="$MODFILES ${submod}/go.mod"
done
echo "Publishing root + sub-modules: $(echo "$SUBMODULE_PATHS" | tr '\n' ' ')"

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

# Rollback guard: if anything between here and the release commit fails (verify,
# drift-guard, etc.), restore the version files this script mutated so the tree
# is left clean for retry. RELEASE_COMMITTED is flipped to 1 once the commit
# lands, so a later tag failure keeps the commit (only the tag needs retrying).
# NOTE: bash keeps only ONE EXIT trap — a second `trap ... EXIT` silently
# replaces the first (this exact bug left a dirty tree behind when verify
# failed during the v1.9.0 cut). Both cleanup concerns live in one hook.
RELEASE_COMMITTED=0
REPLACE_BACKUP_DIR="$(mktemp -d)"
release_cleanup() {
	if [ "$RELEASE_COMMITTED" = "0" ]; then
		echo "Release aborted; rolling back version files, go.mod files, CHANGELOG, FEATURES..." >&2
		git restore utils/version.go $MODFILES CHANGELOG.md FEATURES.md 2>/dev/null || true
	fi
	rm -rf "$REPLACE_BACKUP_DIR"
}
trap release_cleanup EXIT

# 5. Bump utils.Version.
sed -i.bak -E "s|^(const[[:space:]]+Version[[:space:]]+=[[:space:]]+\")[^\"]+(\")|\1${NEW_VERSION}\2|" utils/version.go
rm -f utils/version.go.bak
echo "Bumped utils.Version to $NEW_VERSION"

# 5b. Bump internal module require versions in all go.mod files.
#     Each sub-module references siblings at the shared version. The Go module
#     proxy resolves these via directory-prefixed tags (utils/v2.0.0, etc.).
#     The replace directives override these for local dev; at publish time the
#     proxy uses the tagged version.
for modfile in $MODFILES; do
	sed -i.bak -E \
		"s|(github.com/larsartmann/templ-components/[a-z/]+) v[0-9]+\.[0-9]+\.[0-9]+|\1 v${NEW_VERSION}|g" \
		"$modfile"
	rm -f "${modfile}.bak"
done
echo "Bumped all internal module require entries to v${NEW_VERSION}."

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
' CHANGELOG.md >"$CHANGELOG_TMP"
mv "$CHANGELOG_TMP" CHANGELOG.md
echo "Updated CHANGELOG.md: moved [Unreleased] body under [${NEW_VERSION}] heading."

# 6b. Bump FEATURES.md version + date.
#     AGENTS.md mandates utils.Version, CHANGELOG heading, and FEATURES.md
#     version move together; utils.TestVersionMatchesFeatures enforces it.
if [ -f FEATURES.md ]; then
	sed -i.bak -E "s#(\*\*Updated:\*\* )[0-9]{4}-[0-9]{2}-[0-9]{2}([[:space:]]*\|[[:space:]]*\*\*Version:\*\* )[0-9]+\.[0-9]+\.[0-9]+#\1${TODAY}\2${NEW_VERSION}#" FEATURES.md
	rm -f FEATURES.md.bak
	echo "Bumped FEATURES.md version to $NEW_VERSION (date $TODAY)."
else
	echo "Warning: FEATURES.md not found; skipped its version bump." >&2
fi

# 7. Run full verify (all modules — replaces still in place; see ordering note).
echo "Running full verify (templ generate + CSS compile + per-module build/test/lint)..."
find . -name '*_templ.go' -print0 | xargs -0 rm -f
templ generate ./...

# Recompile every compiled-CSS distribution target so the release commit
# ships fresh artifacts (commands verified byte-identical to BuildFlow's
# tailwind-build output). v1.10.0 shipped these one release stale because
# BuildFlow only recompiled them after the release commit was written.
if command -v tailwindcss &>/dev/null; then
	tailwindcss -i templates/app.css -o templates/styles.css --minify
	tailwindcss -i templates/templ-components-theme.css -o templates/templ-components-theme.out.css --minify
	tailwindcss -i examples/demo/demo.css -o examples/demo/static/app.css --minify
	tailwindcss -i examples/demo/demo.css -o examples/demo/demo.out.css --minify
	tailwindcss -i website/src/styles/global.css -o website/src/styles/global.out.css --minify
else
	echo "Warning: tailwindcss not found — compiled CSS not recompiled." >&2
	echo "Run 'nix run .#build' or install tailwindcss to recompile." >&2
fi

# Root module (workspace mode — replaces removed, tags don't exist yet).
go build ./...
go test ./... -count=1 -race

# Sub-modules (workspace resolves sibling modules locally).
for mod in utils icons errorpage charts/echarts datastar htmx; do
	echo "==> verify $mod"
	(cd "$mod" && go build ./... && go test ./... -count=1 -race)
done

# Lint (golangci-lint does not support go.work).
golangci-lint run \
	./display/... ./feedback/... ./forms/... \
	./integration/... ./internal/... \
	./layout/... ./navigation/... ./recipes/... ./cmd/...
for mod in utils icons errorpage charts/echarts datastar htmx; do
	echo "==> lint $mod"
	(cd "$mod" && golangci-lint run ./...)
done

# Drift-guard: version files must agree with utils.Version (CHANGELOG heading
# AND FEATURES.md version). The full suite ran above; this surfaces a targeted
# message on mismatch. Rollback is handled by the EXIT trap (release_cleanup),
# so no ad-hoc git restore is needed here.
if ! (cd utils && go test ./... -run 'TestVersionMatches(Changelog|Features)' -count=1 >/dev/null 2>&1); then
	echo "Error: version drift-guard failed. utils.Version, CHANGELOG heading, and FEATURES.md version must all agree." >&2
	exit 1
fi

# 7b. Strip replace directives (remove-at-release strategy).
#     Runs AFTER verification: go1.26.5 workspace mode still resolves require
#     entries at unpushed versions, so the suite above must run while the
#     replaces keep sibling resolution local. The release commit ships WITHOUT
#     replace directives so the tagged go.mod files are clean for consumers;
#     they are re-added after tagging (step 10).
for modfile in $MODFILES; do
	backup_name="$(echo "$modfile" | tr '/' '_')"
	# [ /] matches BOTH sub-module replaces (".../utils => ./utils") and the
	# root self-replace ("...templ-components => ./"). The old trailing-slash
	# pattern leaked the self-replace into the v1.10.0 tag's go.mod.
	grep -E '^replace github\.com/larsartmann/templ-components[ /]' "$modfile" >"${REPLACE_BACKUP_DIR}/${backup_name}" 2>/dev/null || true
	sed -i.bak -E '/^replace github\.com\/larsartmann\/templ-components[ \/]/d' "$modfile"
	rm -f "${modfile}.bak"
	# Clean up resulting consecutive blank lines
	sed -i.bak '/^$/N;/^\n$/d' "$modfile"
	rm -f "${modfile}.bak"
done
echo "Removed replace directives from all go.mod files (backed up for re-addition after tagging)."

# Sanity: every go.mod must still parse after the strip.
for modfile in $MODFILES; do
	go mod edit -json "$modfile" >/dev/null || {
		echo "Error: $modfile failed to parse after stripping replace directives." >&2
		exit 1
	}
done

# 8. Stage and commit.
git add utils/version.go CHANGELOG.md FEATURES.md
git add -u # any verified updates

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

RELEASE_COMMITTED=1 # commit landed; EXIT trap no longer rolls back
RELEASE_COMMIT="$(git rev-parse HEAD)"

# 8b. Assert the release commit's TREE before tagging. The auto-commit
#     daemon can split a cut across intermediate commits (v1.10.0: the
#     version bumps landed in a daemon commit, so the script's own commit
#     only carried the replace strip). Tags are immutable once pushed —
#     verify the tree, not the diff, before creating them.
if [ "$(git show HEAD:utils/version.go | grep -E '^const[[:space:]]+Version' | sed -E 's/.*"([^"]+)".*/\1/')" != "$NEW_VERSION" ]; then
	echo "Error: HEAD tree utils/version.go is not $NEW_VERSION (daemon race?)." >&2
	exit 1
fi
git show HEAD:CHANGELOG.md | grep -q "^## \[${NEW_VERSION}\] — " || {
	echo "Error: HEAD tree CHANGELOG.md lacks the [${NEW_VERSION}] heading." >&2
	exit 1
}
git show HEAD:FEATURES.md | grep -q "\*\*Version:\*\* ${NEW_VERSION}" || {
	echo "Error: HEAD tree FEATURES.md version is not ${NEW_VERSION}." >&2
	exit 1
}
for modfile in $MODFILES; do
	if git show "HEAD:${modfile}" | grep -qE '^replace github\.com/larsartmann/templ-components[ /]'; then
		echo "Error: HEAD tree ${modfile} still contains templ-components replace directives." >&2
		exit 1
	fi
done
echo "Release commit tree verified: version files agree, go.mod files replace-free."

# 9. Annotated, SSH-signed tags (root + all sub-modules).
git tag -s "v${NEW_VERSION}" -m "v${NEW_VERSION}: ${RELEASE_SUMMARY}" "$RELEASE_COMMIT"

# Sub-module tags use directory-prefixed format so the Go module proxy
# resolves them correctly (e.g., utils/v2.0.0 → module .../utils at v2.0.0).
for submod in $SUBMODULE_PATHS; do
	git tag -s "${submod}/v${NEW_VERSION}" -m "${submod}/v${NEW_VERSION}: ${RELEASE_SUMMARY}" "$RELEASE_COMMIT"
done

# 10. Re-add replace directives for local dev (remove-at-release strategy).
#     The tagged commit (step 8) has no replace directives — clean for consumers.
#     This follow-up commit restores them so local development continues to work
#     without proxy round-trips.
for modfile in $MODFILES; do
	backup_name="$(echo "$modfile" | tr '/' '_')"
	backup_file="${REPLACE_BACKUP_DIR}/${backup_name}"
	if [ -s "$backup_file" ]; then
		echo "" >>"$modfile"
		cat "$backup_file" >>"$modfile"
	fi
done
git add $MODFILES
git commit -m "chore: re-add replace directives after v${NEW_VERSION} release

These local replace directives were removed for the release commit so the
tagged go.mod files are clean for consumers. They are re-added here for
local development convenience (go.work provides workspace resolution, but
replace directives are needed for GOWORK=off standalone testing).

Co-Authored-By: Crush <noreply@crush.lars.software>"
echo "Re-added replace directives for local dev."

echo ""
echo "Release v${NEW_VERSION} cut at commit ${RELEASE_COMMIT}."
echo "Tags: v${NEW_VERSION} (root) + $(for s in $SUBMODULE_PATHS; do printf '%s/v%s, ' "$s" "$NEW_VERSION"; done | sed 's/, $//') (sub-modules)"
echo ""
echo "Next steps:"
echo "  1. Review the release: git show v${NEW_VERSION}"
echo "  2. Push (when ready):  git push origin master --follow-tags"
echo "  3. House rule says NEVER PUSH TO REMOTE — confirm with the user first."
