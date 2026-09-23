#!/usr/bin/env bash
# Self-test for scripts/check-tc-sources-sync.sh (TODO_LIST #284a/#284b, 2026-09-23).
#
# Exercises the guard's scenarios against a DETACHED TEMP WORKTREE so the real
# tree (and the auto-commit daemon) are never touched:
#
#   1. clean tree            -> exit 0, "complete mirror"
#   2. drift                 -> check exit 1 (DRIFT), --fix refreshes 1, then clean
#   3. new library file      -> check exit 1 (UNEMBEDDED), --fix adds 1, then clean
#   4. orphan embedded file  -> check exit 1 (ORPHAN), --fix removes 1, then clean
#   5. rename (orphan+new)   -> check exit 1, --fix adds 1 + removes 1, then clean
#   6. missing package dir   -> check exit 1 (MISSING PACKAGE), unfixable by --fix
#   7. idempotence           -> --fix twice after a fix; second run is a no-op
#
# Usage: scripts/test-tc-sources-guard.sh
# Exit 0 = every scenario passed. Run from anywhere; needs only git + bash.
set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
GUARD="$REPO_ROOT/scripts/check-tc-sources-sync.sh"
BASE="$(mktemp -d)"
WORKTREE="$BASE/wt"
nwt=0

cleanup() {
	git -C "$REPO_ROOT" worktree remove --force "$WORKTREE" >/dev/null 2>&1 || true
	git -C "$REPO_ROOT" worktree prune >/dev/null 2>&1 || true
	rm -rf "$BASE"
}
trap cleanup EXIT

passed=0
failed=0

# expect <desc> <want_exit> <cmd...> [grep_pattern]
expect() {
	local desc="$1" want="$2" pattern="$3"
	shift 3
	local out rc
	out="$("$@" 2>&1)" && rc=0 || rc=$?
	if [ "$rc" -ne "$want" ]; then
		printf 'FAIL %s: exit %d, want %d\n%s\n' "$desc" "$rc" "$want" "$out"
		return 1
	fi
	if [ -n "$pattern" ] && ! printf '%s' "$out" | grep -q "$pattern"; then
		printf 'FAIL %s: output missing %q\n%s\n' "$desc" "$pattern" "$out"
		return 1
	fi
	printf 'ok   %s\n' "$desc"
}

new_worktree() {
	git -C "$REPO_ROOT" worktree remove --force "$WORKTREE" >/dev/null 2>&1 || true
	git -C "$REPO_ROOT" worktree prune >/dev/null 2>&1 || true
	nwt=$((nwt + 1))
	WORKTREE="$BASE/wt-$nwt"
	git -C "$REPO_ROOT" worktree add --detach -q "$WORKTREE" HEAD >/dev/null
}

record() {
	if [ "$1" -eq 0 ]; then
		passed=$((passed + 1))
	else
		failed=$((failed + 1))
	fi
}

# --- Scenario 1: clean tree ------------------------------------------------
new_worktree
expect "clean check is green" 0 'complete mirror' bash -c "'$GUARD' 2>&1"
record $?

# --- Scenario 2: drift -----------------------------------------------------
new_worktree
printf '\n' >>"$WORKTREE/display/eyebrow.templ"
expect "drift detected" 1 DRIFT bash -c "cd '$WORKTREE' && '$GUARD' 2>&1"
record $?
out="$(cd "$WORKTREE" && "$GUARD" --fix 2>&1)"
if printf '%s' "$out" | grep -q "refreshed 1"; then
	printf 'ok   drift --fix refreshes exactly 1\n'
	record 0
else
	printf 'FAIL drift --fix: %s\n' "$out"
	record 1
fi
expect "clean after drift fix" 0 '' bash -c "cd '$WORKTREE' && '$GUARD' >/dev/null"
record $?

# --- Scenario 3: new library file ------------------------------------------
new_worktree
printf 'component tcGuardFixture {\n\tauto id\n}\n' >"$WORKTREE/display/tc_guard_fixture.templ"
expect "unembedded detected" 1 UNEMBEDDED bash -c "cd '$WORKTREE' && '$GUARD' 2>&1"
record $?
out="$(cd "$WORKTREE" && "$GUARD" --fix 2>&1)"
if printf '%s' "$out" | grep -q "added 1"; then
	printf 'ok   new file --fix adds exactly 1\n'
	record 0
else
	printf 'FAIL new file --fix: %s\n' "$out"
	record 1
fi
expect "clean after add fix" 0 '' bash -c "cd '$WORKTREE' && '$GUARD' >/dev/null"
record $?

# --- Scenario 4: orphan embedded file --------------------------------------
new_worktree
rm "$WORKTREE/display/eyebrow.templ"
expect "orphan detected" 1 ORPHAN bash -c "cd '$WORKTREE' && '$GUARD' 2>&1"
record $?
out="$(cd "$WORKTREE" && "$GUARD" --fix 2>&1)"
if printf '%s' "$out" | grep -q "removed 1"; then
	printf 'ok   orphan --fix removes exactly 1\n'
	record 0
else
	printf 'FAIL orphan --fix: %s\n' "$out"
	record 1
fi
expect "clean after orphan fix" 0 '' bash -c "cd '$WORKTREE' && '$GUARD' >/dev/null"
record $?

# --- Scenario 5: rename = orphan + new --------------------------------------
new_worktree
mv "$WORKTREE/display/eyebrow.templ" "$WORKTREE/display/eyebrow_renamed.templ"
expect "rename detected (orphan)" 1 ORPHAN bash -c "cd '$WORKTREE' && '$GUARD' 2>&1"
record $?
out="$(cd "$WORKTREE" && "$GUARD" --fix 2>&1)"
if printf '%s' "$out" | grep -q "added 1" && printf '%s' "$out" | grep -q "removed 1"; then
	printf 'ok   rename --fix adds 1 and removes 1\n'
	record 0
else
	printf 'FAIL rename --fix: %s\n' "$out"
	record 1
fi
expect "clean after rename fix" 0 '' bash -c "cd '$WORKTREE' && '$GUARD' >/dev/null"
record $?

# --- Scenario 6: missing package dir ----------------------------------------
new_worktree
mv "$WORKTREE/datastar" "$WORKTREE/datastar.hidden"
expect "missing package detected" 1 'MISSING PACKAGE: datastar' bash -c "cd '$WORKTREE' && '$GUARD' 2>&1"
record $?
expect "missing package is unfixable" 1 '' bash -c "cd '$WORKTREE' && '$GUARD' --fix >/dev/null 2>&1"
record $?
mv "$WORKTREE/datastar.hidden" "$WORKTREE/datastar"

# --- Scenario 7: idempotence ------------------------------------------------
new_worktree
printf '\n' >>"$WORKTREE/forms/slider.templ"
cd "$WORKTREE"
"$GUARD" --fix >/dev/null 2>&1 || true
out2="$("$GUARD" --fix 2>&1)" && rc2=0 || rc2=$?
out3="$("$GUARD" --fix 2>&1)" && rc3=0 || rc3=$?
if [ "$rc2" -eq 0 ] && [ "$rc3" -eq 0 ] && ! printf '%s' "$out3" | grep -q "synced [1-9]"; then
	printf 'ok   idempotent: second and third --fix are no-ops\n'
	record 0
else
	printf 'FAIL idempotence: rc2=%d rc3=%d out3=%s\n' "$rc2" "$rc3" "$out3"
	record 1
fi

printf '\n%d passed, %d failed\n' "$passed" "$failed"
[ "$failed" -eq 0 ]
