#!/usr/bin/env bash
# Changelog-guard policy (#133): which diffs must warm CHANGELOG.md.
#
# Policy (decided 2026-09-09, supersedes the open owner question):
#   - Library CODE changes (*.templ, or non-test .go under the library
#     packages) REQUIRE a CHANGELOG.md diff in the same PR.
#   - Test-only diffs (_test.go, visualtest goldens/fixtures), docs-only
#     diffs, examples/, cmd/, and website/ diffs are EXEMPT — they change
#     no consumer-observable behaviour, so a release note would be noise.
#     Release-time CHANGELOG entries for notable test-infrastructure work
#     are welcome but voluntary.
#   - The release script itself still refuses to cut with an empty
#     [Unreleased], independent of this guard.
#
# Usage: check-changelog-guard.sh <changed-file-lines...>
#        (stdin also accepted; used by CI with the PR file list)
# Exit 0 = guard satisfied; exit 1 = component code changed without CHANGELOG.md.

set -euo pipefail

files=()
if [ $# -gt 0 ]; then
	files=("$@")
else
	while IFS= read -r line; do
		[ -n "$line" ] && files+=("$line")
	done
fi

lib_code=""
changelog=0

for f in "${files[@]}"; do
	case "$f" in
	CHANGELOG.md)
		changelog=1
		;;
	*.templ)
		lib_code="$lib_code$f"$'\n'
		;;
	*.go)
		case "$f" in
		*_test.go) ;;                                     # test-only diffs are exempt (#133)
		examples/* | visualtest/* | cmd/* | website/*) ;; # non-library surfaces
		*)
			lib_code="$lib_code$f"$'\n'
			;;
		esac
		;;
	esac
done

echo "Changed library code files:"
if [ -n "$lib_code" ]; then
	printf '%s' "$lib_code"
else
	echo "  (none)"
fi

if [ -n "$lib_code" ] && [ "$changelog" -eq 0 ]; then
	echo "::error::This PR changes library components but does not touch CHANGELOG.md." >&2
	echo "::error::House rule: every feature/fix change warms [Unreleased] in the same change." >&2
	echo "::error::Test-only, docs-only, examples, cmd, and website diffs are exempt." >&2
	exit 1
fi

echo "Changelog guard satisfied."
