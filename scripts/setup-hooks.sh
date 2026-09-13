#!/usr/bin/env bash
# Activate the tracked git hooks for this clone (idempotent).
#
# The pre-commit guards live in .githooks/ (tracked) instead of .git/hooks/
# (untracked, lost on fresh clones). This script points git at them.
#
# Run once after cloning:
#   scripts/setup-hooks.sh

set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)

if [ "$(git config core.hooksPath)" = ".githooks" ]; then
    echo "hooks already active (core.hooksPath=.githooks)"
else
    git config core.hooksPath .githooks
    echo "activated core.hooksPath=.githooks"
fi

if [ ! -x "$REPO_ROOT/.githooks/pre-commit" ]; then
    echo "ERROR: .githooks/pre-commit missing or not executable" >&2
    exit 1
fi

echo "tracked pre-commit hook ready (fast guards + BuildFlow)."
echo "For the full pre-push verify run: scripts/pre-commit.sh"
