#!/usr/bin/env bash
# Guard 0: hook-adoption check (TODO #209)
#
# The tracked pre-commit guards live in .githooks/ — they are active ONLY when
# `git config core.hooksPath` points there. A clone that copied the hook into
# .git/hooks/ (or ran 'buildflow precommit install') executes a STALE copy
# that drifts from the tracked source of truth. Loud warning, deliberately
# non-fatal: a missing adoption should teach, not block the commit.
set -euo pipefail

CURRENT=$(git config core.hooksPath || true)

if [ "$CURRENT" != ".githooks" ]; then
	echo "WARNING: core.hooksPath is '${CURRENT:-<unset>}' — this clone is NOT running the tracked .githooks/ guards." >&2
	echo "         (You may be executing a stale .git/hooks/ copy.) Fix once:" >&2
	echo "           scripts/setup-hooks.sh" >&2
fi
