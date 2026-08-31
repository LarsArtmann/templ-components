#!/usr/bin/env bash
# scripts/check-master-green.sh — turn a red master into an open issue.
#
# Runs daily from .github/workflows/master-red-alert.yml (and on demand via
# workflow_dispatch). The v1.11.0 lesson: master CI sat red for 9 days with
# nobody noticing because nothing alerts on a failing default branch. This
# script checks the latest completed master run of each watched workflow and
# opens (or updates) a single tracking issue while any of them is red.
#
# Requires: gh (authenticated via GH_TOKEN), jq.

set -euo pipefail

REPO="${GITHUB_REPOSITORY:-}"
if [ -z "$REPO" ]; then
	REPO="$(gh repo view --json nameWithOwner --jq .nameWithOwner)"
fi

WORKFLOWS=("CI" "Website")
declare -a RED=()

for wf in "${WORKFLOWS[@]}"; do
	json="$(gh run list --repo "$REPO" --workflow "$wf" --branch master \
		--status completed --limit 1 --json name,conclusion,createdAt,url 2>/dev/null || echo '[]')"
	conclusion="$(echo "$json" | jq -r '.[0].conclusion // ""')"
	if [ "$conclusion" = "failure" ]; then
		RED+=("$(echo "$json" | jq -r '.[0] | "- \(.name): latest master run failed at \(.createdAt) — \(.url)"')")
	fi
done

if [ "${#RED[@]}" -eq 0 ]; then
	echo "master-green: all watched workflows are green. Nothing to do."
	exit 0
fi

TITLE="master CI is red"
BODY="Automated master health check $(date -u '+%Y-%m-%dT%H:%M:%SZ'):

$(printf '%s\n' "${RED[@]}")

Opened by the master-red-alert workflow. Close once master is green again."

existing="$(gh issue list --repo "$REPO" --state open --json number,title \
	--jq "map(select(.title == \"$TITLE\")) | .[0].number // empty")"

if [ -n "$existing" ]; then
	echo "master-red: updating existing tracking issue #$existing"
	gh issue comment "$existing" --repo "$REPO" --body "$BODY"
else
	echo "master-red: opening new tracking issue"
	gh issue create --repo "$REPO" --title "$TITLE" --body "$BODY"
fi
