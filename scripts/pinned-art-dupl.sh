#!/usr/bin/env bash
# scripts/pinned-art-dupl.sh — build the art-dupl analyzer at the PINNED rev.
#
# Baseline hashes are only comparable within ONE art-dupl binary (ADR-0009,
# tool-version pin). The fork checkout ~/projects/art-dupl moves; building it
# at HEAD yields an analyzer whose hashes match NOTHING in the committed
# baseline — `check` then reports the whole tree as new (the 105-group false
# alarm class, 2026-10-24 near-miss). This script always builds the exact rev
# the baseline was recorded with, echoes its version, and leaves the binary at
# /tmp/art-dupl-fork — the advisory CI lane's first resolution path
# (scripts/ci-repro.sh resolves $ART_DUPL_BIN → /tmp/art-dupl-fork → PATH).
#
# Usage:
#   scripts/pinned-art-dupl.sh            # build + echo version
#   scripts/pinned-art-dupl.sh -- check   # build, then run `check` args
#
# After a DELIBERATE analyzer bump: re-record the baseline with the SAME new
# binary, update ADR-0009 + the PINNED_REV below, and log the version (ritual
# in ADR-0009 "Baseline").

set -euo pipefail

FORK_REPO="${ART_DUPL_REPO:-$HOME/projects/art-dupl}"
PINNED_REV="93432fd910fb" # v0.9.1-0.20261017225421-93432fd910fb — matches .art-dupl-baseline.json
WORKTREE="/tmp/art-dupl-pinned"
OUT="/tmp/art-dupl-fork"

if [ ! -d "$FORK_REPO/.git" ]; then
	echo "ERROR: art-dupl fork not found at $FORK_REPO (set ART_DUPL_REPO)" >&2
	exit 1
fi

git -C "$FORK_REPO" worktree remove --force "$WORKTREE" >/dev/null 2>&1 || true
git -C "$FORK_REPO" worktree add --detach --quiet "$WORKTREE" "$PINNED_REV"
trap 'git -C "$FORK_REPO" worktree remove --force "$WORKTREE" >/dev/null 2>&1 || true' EXIT

GOTOOLCHAIN=auto go -C "$WORKTREE" build -o "$OUT" ./cmd/art-dupl

echo "built $OUT:"
"$OUT" --version

if [ "${1:-}" = "--" ]; then
	shift
	exec "$OUT" "$@"
fi
