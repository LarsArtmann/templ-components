#!/usr/bin/env bash
# Guards cmd/tc/_sources as a COMPLETE mirror of the library sources the
# `tc add` scaffolder ships.
#
#   Direction 1 (embedded -> library): every embedded file (minus the
#     exclusions below) must have an identical library twin. Catches content
#     drift AND orphans (twin deleted/renamed) at commit time.
#   Direction 2 (library -> embedded): every .templ and *_types.go in the
#     mirrored packages must be embedded. Catches newly added components the
#     scaffolder would silently miss ('tc add' registers .templ files from the
#     embedded FS and copies <base>_types.go from beside the .templ).
#
# WHY: the scaffolder ships COPIES of real component sources. When a component
# is edited without re-copying, `tc add` scaffolds stale code; when a component
# is added without copying, `tc add` cannot scaffold it at all. The Go guard
# (cmd/tc TestSourcesMatchPackageFiles) only fires in `go test`, which
# BuildFlow's pre-commit budget cannot run. Direction 1 caught real drift
# twice on 2026-09-17 (notfound404 + shared).
#
# Exclusions (keep in sync with cmd/tc/main_test.go):
#   - starter/**                         curated scaffolder templates, no twins
#   - datastar/DATASTAR-BUMP-PROTOCOL.md scaffolder-owned checklist
#
# Modes:
#   (default)  check only — print every problem, exit 1
#   --fix      full mirror: re-copy drifted, copy unembedded, remove orphans;
#              exits 0 once the mirror is complete (prints
#              "synced N file(s) (added X, refreshed Y, removed Z)"),
#              1 only if a problem cannot be fixed
#
# Usage:
#   scripts/check-tc-sources-sync.sh            # exit 1 on any gap
#   scripts/check-tc-sources-sync.sh --fix      # mirror everything
set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
cd "$REPO_ROOT"

FIX="${1:-}"
SOURCES_DIR="cmd/tc/_sources"
DATASTAR_DOC="DATASTAR-BUMP-PROTOCOL.md"

# Mirrored packages: the top-level component packages `tc add` ships.
# charts/ (Tier 2 opt-in adapter) and utils/ (leaf plumbing) are intentionally
# not mirrored. Keep in sync with mirroredPackages in cmd/tc/main_test.go.
MIRRORED_PKGS=(datastar display errorpage feedback forms htmx layout navigation recipes)

added=0 refreshed=0 removed=0 problems=0
# Set when a problem exists that --fix must NOT paper over: a missing package
# directory means a broken checkout (or a list error), and auto-removing every
# embedded copy of that package would be destructive rather than healing
# (found by scripts/test-tc-sources-guard.sh 2026-09-23).
unfixable=0

note() { printf '%s\n' "$*" >&2; }

# Packages whose library directory is absent. --fix must not treat their
# embedded copies as removable orphans: a missing package directory means a
# broken checkout (or an MIRRORED_PKGS list error), and auto-removing every
# embedded copy of that package would be destructive rather than healing
# (found by scripts/test-tc-sources-guard.sh 2026-09-23).
missing_pkgs=""
for pkg in "${MIRRORED_PKGS[@]}"; do
	if [ ! -d "$pkg" ]; then
		missing_pkgs="$missing_pkgs $pkg/"
	fi
done

pkg_is_missing() {
	case " $missing_pkgs " in
	*" $1/"*) return 0 ;;
	esac

	return 1
}

# Direction 1: every embedded file has an identical library twin.
while IFS= read -r -d '' embedded; do
	rel="${embedded#"$SOURCES_DIR"/}"

	# Exclusions (mirror the Go test): scaffolder-owned files without twins.
	case "$rel" in
	starter/*) continue ;;
	"datastar/$DATASTAR_DOC") continue ;;
	esac

	pkg="${rel%%/*}"
	if pkg_is_missing "$pkg"; then
		note "ORPHAN-SKIPPED: $embedded (package '$pkg' directory is missing — fix the checkout, not the mirror)"
		problems=$((problems + 1))
		unfixable=1
		continue
	fi

	twin="$rel"
	if [ ! -f "$twin" ]; then
		note "ORPHAN: $embedded has no library twin"
		problems=$((problems + 1))
		if [ "$FIX" = "--fix" ]; then
			# Orphans are normally tracked; git rm keeps history explicit. Fall
			# back to plain deletion only for untracked strays (daemon races).
			if git ls-files --error-unmatch -- "$embedded" >/dev/null 2>&1; then
				git rm -q -- "$embedded"
			else
				rm -f -- "$embedded"
			fi
			note "  fixed: removed $embedded"
			removed=$((removed + 1))
		fi
		continue
	fi

	if ! cmp -s "$embedded" "$twin"; then
		note "DRIFT: $embedded differs from $twin"
		problems=$((problems + 1))
		if [ "$FIX" = "--fix" ]; then
			cp -- "$twin" "$embedded"
			note "  fixed: re-copied $twin -> $embedded"
			refreshed=$((refreshed + 1))
		fi
	fi
done < <(find "$SOURCES_DIR" -type f -print0)

# Direction 2: every mirrorable library file is embedded.
for pkg in "${MIRRORED_PKGS[@]}"; do
	if [ ! -d "$pkg" ]; then
		note "MISSING PACKAGE: $pkg is listed in MIRRORED_PKGS but has no directory"
		problems=$((problems + 1))
		unfixable=1
		continue
	fi

	while IFS= read -r -d '' twin; do
		embedded="$SOURCES_DIR/$twin"
		if [ -f "$embedded" ]; then
			continue # content equality is Direction 1's job
		fi

		note "UNEMBEDDED: $twin is not shipped by 'tc add'"
		problems=$((problems + 1))
		if [ "$FIX" = "--fix" ]; then
			mkdir -p "$(dirname -- "$embedded")"
			cp -- "$twin" "$embedded"
			note "  fixed: copied $twin -> $embedded"
			added=$((added + 1))
		fi
	done < <(find "$pkg" -type f \( -name '*.templ' -o -name '*_types.go' \) -print0)
done

synced=$((added + refreshed + removed))

if [ "$problems" -ne 0 ]; then
	if [ "$unfixable" -eq 1 ]; then
		note "BLOCKED: unfixable problems present (see above); no fixes were applied for them."
		exit 1
	fi

	if [ "$FIX" = "--fix" ] && [ "$synced" -gt 0 ]; then
		printf 'synced %d file(s) (added %d, refreshed %d, removed %d)\n' "$synced" "$added" "$refreshed" "$removed"
		exit 0
	fi

	if [ "$FIX" = "--fix" ]; then
		note "BLOCKED: tc scaffolder mirror has unfixable problems (see above)."
	else
		note ""
		note "BLOCKED: cmd/tc/_sources is not a complete mirror of the library sources."
		note "Run scripts/check-tc-sources-sync.sh --fix, review, and commit the result."
	fi

	exit 1
fi

if [ "$synced" -gt 0 ]; then
	printf 'synced %d file(s) (added %d, refreshed %d, removed %d)\n' "$synced" "$added" "$refreshed" "$removed"
else
	note "tc scaffolder sources are a complete mirror of the library twins."
fi

exit 0
