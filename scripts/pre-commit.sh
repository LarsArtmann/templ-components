#!/usr/bin/env bash
# Pre-commit hook: ensure templ-generated files are in sync
# Install: cp scripts/pre-commit.sh .git/hooks/pre-commit
set -e

STAGED_TEMPL=$(git diff --cached --name-only --diff-filter=ACM | grep '\.templ$' || true)

if [ -n "$STAGED_TEMPL" ]; then
	echo "templ files changed, regenerating..."
	templ generate ./...
	
	# Stage the regenerated files
	CHANGED_GO=$(git diff --name-only | grep '_templ\.go$' || true)
	if [ -n "$CHANGED_GO" ]; then
		echo "Staging regenerated files: $CHANGED_GO"
		git add $CHANGED_GO
	fi
fi
