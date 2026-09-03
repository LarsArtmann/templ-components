#!/usr/bin/env bash
# Fixture test for the release-cut tree assertions
# (scripts/lib-release-assertions.sh, release.sh step 8b).
#
# Builds a throwaway git repo in a tmpdir and runs the real assertion function
# against it — no templ, no Go toolchain, no network; finishes in well under a
# second, so it can run in CI for pennies.
#
# The CHANGELOG fixture deliberately exceeds the 64KB pipe buffer with the
# heading at the TOP (real CHANGELOGs are newest-first), reproducing the
# v1.12.0 SIGPIPE geometry — see the comment at the fixture. The tmpdir is
# left for tmpfs cleanup (the house rule bans rm, and the dir was created
# seconds earlier).
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib-release-assertions.sh
source "${SCRIPT_DIR}/lib-release-assertions.sh"

TMP="$(mktemp -d)"
cd "$TMP"
git init -q
git config user.email test@example.com
git config user.name fixture

VERSION="9.9.9"
mkdir utils
cat >utils/version.go <<EOF
package utils

const Version = "${VERSION}"
EOF
# Heading EARLY, bulk after — real CHANGELOGs are newest-first, and this is
# the v1.12.0 SIGPIPE geometry: grep finds the match in the first bytes and
# exits while `git show` is still pumping ~150KB through the 64KB pipe. A
# `grep -q` regression then kills the pipeline with 141 under pipefail and
# the positive cases below fail. (Heading at the END would let grep read to
# EOF and never SIGPIPE — the fixture would test nothing.)
{
	echo "## [${VERSION}] — 2026-09-03"
	echo ""
	head -c 150000 /dev/zero | tr '\0' '#'
	echo ""
} >CHANGELOG.md
echo "**Version:** ${VERSION}" >FEATURES.md
cat >go.mod <<'EOF'
module github.com/larsartmann/templ-components

go 1.26.7
EOF
cat >utils/go.mod <<EOF
module github.com/larsartmann/templ-components/utils

go 1.26.7
EOF
mkdir -p errorpage
cat >errorpage/go.mod <<EOF
module github.com/larsartmann/templ-components/errorpage

go 1.26.7
EOF
mkdir -p icons
cat >icons/go.mod <<EOF
module github.com/larsartmann/templ-components/icons

go 1.26.7

replace github.com/larsartmann/templ-components/utils => ../utils
EOF
git add -A
git commit -qm fixture

pass=0
fail=0
expect_ok() {
	if assert_release_tree "$@" >/dev/null 2>&1; then
		pass=$((pass + 1))
	else
		fail=$((fail + 1))
		echo "FAIL (expected ok): assert_release_tree $*" >&2
	fi
}
expect_err() {
	if assert_release_tree "$@" >/dev/null 2>&1; then
		fail=$((fail + 1))
		echo "FAIL (expected error): assert_release_tree $*" >&2
	else
		pass=$((pass + 1))
	fi
}

# Positive: clean tree, version files agree. The >64KB CHANGELOG with the
# heading past the pipe buffer is the SIGPIPE regression pin.
expect_ok "$VERSION" go.mod utils/go.mod errorpage/go.mod
expect_ok "$VERSION" go.mod utils/go.mod
# Negative: wrong version everywhere.
expect_err "0.0.1" go.mod utils/go.mod errorpage/go.mod
# Negative: a leaked replace directive in one module file.
expect_err "$VERSION" go.mod utils/go.mod icons/go.mod

echo "release-assertions: ${pass} passed, ${fail} failed"
[ "$fail" -eq 0 ]
