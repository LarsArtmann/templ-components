#!/usr/bin/env bash
# Shared release-cut tree assertions, extracted from scripts/release.sh step 8b
# so scripts/test-release-assertions.sh can exercise them against fixtures
# without cutting a release.
#
# Callers MUST have `set -o pipefail` active: the SIGPIPE class these
# assertions guard against (v1.12.0 cut — `git show | grep -q` exit 141 once
# the file passed the 64KB pipe buffer) only fires under pipefail, and the
# `grep -c` idiom below is the fix being pinned. Callers MUST run inside a git
# repo whose HEAD commit carries the release tree.

# assert_release_tree <version> <modfile>...
# Returns 1 unless HEAD's tree pins utils/version.go at <version>, carries the
# CHANGELOG heading and FEATURES.md version for <version>, and has no
# templ-components replace directives in any <modfile>.
assert_release_tree() {
	if [ "$#" -lt 1 ]; then
		echo "usage: assert_release_tree <version> <modfile>..." >&2
		return 1
	fi
	local version="$1"
	shift
	if [ "$(git show HEAD:utils/version.go | grep -E '^const[[:space:]]+Version' | sed -E 's/.*"([^"]+)".*/\1/')" != "$version" ]; then
		echo "Error: HEAD tree utils/version.go is not $version (daemon race?)." >&2
		return 1
	fi
	# grep -q exits at the first match, so under `set -o pipefail` the upstream
	# `git show` dies of SIGPIPE once the file exceeds the 64KB pipe buffer and
	# the pipeline reports 141 even when the match exists (v1.12.0 cut:
	# CHANGELOG.md at 155KB tripped the heading assertion and aborted pre-tag).
	# `grep -c` reads the whole input, so the pipeline exit is the match result
	# alone.
	git show HEAD:CHANGELOG.md | grep -c "^## \[${version}\] — " >/dev/null || {
		echo "Error: HEAD tree CHANGELOG.md lacks the [${version}] heading." >&2
		return 1
	}
	git show HEAD:FEATURES.md | grep -c "\*\*Version:\*\* ${version}" >/dev/null || {
		echo "Error: HEAD tree FEATURES.md version is not ${version}." >&2
		return 1
	}
	local modfile
	for modfile in "$@"; do
		if git show "HEAD:${modfile}" | grep -cE '^replace github\.com/larsartmann/templ-components[ /]' >/dev/null; then
			echo "Error: HEAD tree ${modfile} still contains templ-components replace directives." >&2
			return 1
		fi
	done
	return 0
}
