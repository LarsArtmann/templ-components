#!/usr/bin/env bash
# test-settle-staged-go-format.sh — hermetic proof that the pre-commit
# formatter settle step works AND stays out of the worktree (#441).
#
# Cases:
#   1. a misformatted staged Go file is replaced in the INDEX
#   2. the worktree file is untouched (unstaged edits cannot leak into the commit)
#   3. a second run is a no-op (idempotent)
#   4. non-Go staged files are ignored
set -euo pipefail

SCRIPT="$(cd "$(dirname "$0")" && pwd)/settle-staged-go-format.sh"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

fail() {
	printf '%s\n' "$*" >&2
	echo "FAILED: settle-staged-go-format" >&2
	exit 1
}

ok() { printf 'ok: %s\n' "$*"; }

cd "$TMP"
git init -q
git config user.email test@example.com
git config user.name test
git config commit.gpgsign false

# 1 + 2: stage a misformatted file, then edit the worktree WITHOUT staging.
printf 'package sample\n\nfunc  Add(a,b int) int{return a+b}\n' >sample.go
git add sample.go
staged_before="$(git rev-parse :sample.go)"

printf 'package sample\n\nfunc Add(a, b int) int { return a + b } // NOT STAGED\n' >sample.go

SETTLE_REPO_ROOT="$TMP" "$SCRIPT" 2>"$TMP/log"

staged_after="$(git rev-parse :sample.go)"
[ "$staged_before" != "$staged_after" ] || fail "case 1: index blob was not replaced"
ok "case 1: misformatted staged file was settled in the index"

[ "$(cat sample.go)" = "$(printf 'package sample\n\nfunc Add(a, b int) int { return a + b } // NOT STAGED\n')" ] ||
	fail "case 2: worktree file was modified"
ok "case 2: worktree untouched (unstaged edit did not leak)"

if git show :sample.go | gofmt -l | grep -q .; then
	fail "case 1: staged blob is still not gofmt-clean"
fi

if git show :sample.go | grep -q 'NOT STAGED'; then
	fail "case 2: unstaged content reached the index"
fi

# 3: idempotence.
SETTLE_REPO_ROOT="$TMP" "$SCRIPT" 2>>"$TMP/log"
[ "$(git rev-parse :sample.go)" = "$staged_after" ] || fail "case 3: second run rewrote the index"
ok "case 3: second run is a no-op"

# 4: unrelated staged files are left alone.
printf 'not go\n' >notes.txt
git add notes.txt
notes_before="$(git rev-parse :notes.txt)"
SETTLE_REPO_ROOT="$TMP" "$SCRIPT" 2>>"$TMP/log"
[ "$(git rev-parse :notes.txt)" = "$notes_before" ] || fail "case 4: non-Go staged file was rewritten"
ok "case 4: non-Go staged files are ignored"

echo "PASS: settle-staged-go-format self-test"
