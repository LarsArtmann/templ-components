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

# golangci-lint VERSION lockstep (TODO #443): CI installs one pinned version
# (ci.yaml) while the dev shell ships nixpkgs' — a skew means the local gate
# can demand more (or less) than CI enforces. Skew is not a code bug: fix is
# a dedicated pin-bump commit (workflow + flake input in lockstep), never a
# release cargo. Loud skip when no binary is on PATH (pre-commit may run
# before the dev shell is entered).
CI_GCL_REV="$(sed -n 's|.*golangci-lint/v2/cmd/golangci-lint@\(v[0-9.]*\).*|\1|p' .github/workflows/ci.yaml | head -1)"
if [ -z "$CI_GCL_REV" ]; then
	echo "WARN: could not parse the golangci-lint pin from ci.yaml — lockstep unchecked" >&2
	exit 0
fi

if ! command -v golangci-lint >/dev/null 2>&1; then
	echo "SKIP: golangci-lint not on PATH — version lockstep unchecked (CI pins $CI_GCL_REV)" >&2
	exit 0
fi

# Normalize both sides to bare x.y.z (CI pins carry v; `golangci-lint
# version` prints the bare form in its "has version X" field).
LOCAL_GCL_VER="$(golangci-lint version 2>/dev/null | head -1 | sed -n 's/.* has version v\{0,1\}\([0-9][0-9.]*\) .*/\1/p')"
CI_GCL_VER="${CI_GCL_REV#v}"
if [ -z "$LOCAL_GCL_VER" ]; then
	echo "WARN: could not parse local golangci-lint version — lockstep unchecked" >&2
	exit 0
fi

if [ "$LOCAL_GCL_VER" != "$CI_GCL_VER" ]; then
	echo "BLOCKED: golangci-lint version skew — CI pins $CI_GCL_REV, local is $LOCAL_GCL_VER." >&2
	echo "  Fix the CODE so both versions pass, then align the pins in ONE dedicated" >&2
	echo "  commit (ci.yaml install rev + flake nixpkgs input together — never release cargo)." >&2
	echo "  History: TODO #443 / AGENTS lint-skew note (2026-10-23)." >&2
	exit 1
fi

echo "golangci-lint lockstep OK: $LOCAL_GCL_VER (CI == local == $CI_GCL_REV)"

exit 0
