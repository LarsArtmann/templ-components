#!/usr/bin/env bash
# Session tripwire (TODO_LIST #337): end-of-session diff review against the
# session-start SHA. The auto-commit daemon interleaves raw working-tree
# snapshots with real commits, re-tracks deleted files, and races pushes —
# this review is the scripted habit that catches what the eye skips:
#
#   scripts/session-tripwire.sh <session-start-sha>
#
# Checks:
#   1. Working tree clean (uncommitted work is not reviewable).
#   2. Full diff stat vs the session-start SHA.
#   3. Deleted-then-reappeared files (the styles.css re-tracking class).
#   4. Daemon commits whose files were also touched by real commits
#      (the torn-snapshot signature — verify those files byte-by-byte).
#   5. Empty daemon commits (0-file snapshots are noise, but list them).
set -euo pipefail

start="${1:?usage: session-tripwire.sh <session-start-sha>}"
git rev-parse --verify --quiet "${start}^{commit}" >/dev/null ||
	{ echo "FATAL: ${start} is not a commit" >&2; exit 1; }

echo "== 1. working tree =="
if [ -n "$(git status --porcelain)" ]; then
	echo "DIRTY — commit or stash before the review:"
	git status --short
	exit 1
fi
echo "clean"

echo
echo "== 2. diff stat ${start}..HEAD =="
git diff --stat "${start}"..HEAD | tail -15
files="$(git diff --name-status "${start}"..HEAD)"
if [ -z "$files" ]; then
	echo "no changes in range — nothing to review"
	exit 0
fi
total="$(echo "$files" | wc -l | tr -d ' ')"
echo "($total files changed in range)"

echo
echo "== 3. deleted-then-reappeared files =="
deleted="$(echo "$files" | awk '$1 == "D" { print $2 }')"
reappeared=0
for f in $deleted; do
	if git cat-file -e "HEAD:$f" 2>/dev/null; then
		echo "RETRACKED: $f (deleted in range, present at HEAD)"
		reappeared=$((reappeared + 1))
	fi
done
[ "$reappeared" -eq 0 ] && echo "none"

echo
echo "== 4. daemon commits overlapping real commits =="
overlap=0
daemon_shas="$(git log --format='%H %s' "${start}"..HEAD | awk '/chore: auto-commit/ { print $1 }')"
for dsha in $daemon_shas; do
	dfiles="$(git show --name-only --format='' "$dsha")"
	[ -z "$dfiles" ] && continue
	count="$(git show --format='' --name-only "$dsha" | wc -l | tr -d ' ')"
	if [ "$count" -eq 0 ]; then
		echo "EMPTY daemon commit: $(git rev-parse --short "$dsha")"
		continue
	fi
	for dfile in $dfiles; do
		real_touches="$(git log --format='%h %s' "${start}"..HEAD -- "$dfile" | grep -v 'chore: auto-commit' | wc -l | tr -d ' ')"
		if [ "$real_touches" -gt 0 ]; then
			echo "OVERLAP: $(git rev-parse --short "$dsha") + $real_touches real commit(s) both touch $dfile — diff-review this file"
			overlap=$((overlap + 1))
			break
		fi
	done
done
[ "$overlap" -eq 0 ] && echo "none"

echo
echo "== verdict =="
if [ "$reappeared" -eq 0 ] && [ "$overlap" -eq 0 ]; then
	echo "CLEAN: no re-tracked deletions, no daemon/real overlap — spot-check the stat above and you are done."
else
	echo "REVIEW REQUIRED: resolve the flags above (deleted files: decide keep/delete; overlaps: diff each file against what you authored)."
fi
