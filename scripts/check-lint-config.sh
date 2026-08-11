#!/usr/bin/env bash
# Fast (<50ms) pre-commit guard: blocks commits if .golangci.yml re-enables
# a documented-disabled linter. This catches the recurring regression where
# a stale working tree silently re-adds ireturn/godoclint/testableexamples.
#
# Root cause (identified 2026-07-28): The AI agent's working tree sometimes
# holds a stale .golangci.yml. When staging broad changes (e.g., `git add -A`
# for docs work), the stale file rides along. BuildFlow's pre-commit hook
# runs with a 60s budget and does NOT run `go test ./...`, so the
# TestGolangciDisabledLinters guard never fires at commit time — only in CI.
# This script closes that gap.
#
# Usage:
#   scripts/check-lint-config.sh          # exit 1 on violation
#   scripts/check-lint-config.sh --quiet  # suppress output
set -euo pipefail

CONFIG=".golangci.yml"
QUIET="${1:-}"

if [ ! -f "$CONFIG" ]; then
	exit 0
fi

# Only flag disabled linters if they appear in the enable: section (not disable:).
# Uses awk to track which YAML subsection we're in.
VIOLATIONS=$(awk '
BEGIN { section = "" }
/^[[:space:]]{2}enable:[[:space:]]*$/ { section = "enable"; next }
/^[[:space:]]{2}[a-z]+:[[:space:]]*$/ { section = "" }
section == "enable" && /[[:space:]]-[[:space:]]*(godoclint|ireturn|testableexamples)[[:space:]]*$/ {
    printf "%d: %s\n", NR, $0
}
' "$CONFIG" 2>/dev/null || true)
IRETURN_BLOCK=$(awk '
BEGIN { section = "" }
/^[[:space:]]{2}enable:[[:space:]]*$/ { section = "enable"; next }
/^[[:space:]]{2}[a-z]+:[[:space:]]*$/ { section = "" }
section == "enable" && /^[[:space:]]+ireturn:/ {
    printf "%d: %s\n", NR, $0
}
' "$CONFIG" 2>/dev/null || true)

if [ -n "$VIOLATIONS" ] || [ -n "$IRETURN_BLOCK" ]; then
	if [ "$QUIET" != "--quiet" ]; then
		echo "" >&2
		echo "BLOCKED: .golangci.yml re-enables a disabled linter." >&2
		echo "" >&2
		echo "The following linters are FUNDAMENTALLY INCOMPATIBLE with a templ library" >&2
		echo "and are documented as disabled in AGENTS.md:" >&2
		if [ -n "$VIOLATIONS" ]; then
			echo "" >&2
			echo "$VIOLATIONS" | while IFS= read -r line; do
				echo "  $line" >&2
			done
		fi
		if [ -n "$IRETURN_BLOCK" ]; then
			echo "" >&2
			echo "  Dead ireturn: settings block found — delete it." >&2
			echo "$IRETURN_BLOCK" | while IFS= read -r line; do
				echo "    $line" >&2
			done
		fi
		echo "" >&2
		echo "Fix: remove these entries from .golangci.yml." >&2
		echo "This regression has occurred 5 times — see TestGolangciDisabledLinters." >&2
		echo "" >&2
	fi
	exit 1
fi

exit 0
