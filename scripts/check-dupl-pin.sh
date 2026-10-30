#!/usr/bin/env bash
# check-dupl-pin.sh — the art-dupl analyzer rev must be IN LOCKSTEP across
# the three places that must agree (ADR-0009 tool-version pin):
#
#   1. scripts/pinned-art-dupl.sh   PINNED_REV   (local supply)
#   2. .github/workflows/ci.yaml    go install @rev (blocking clone-gate job)
#   3. .art-dupl-baseline.json      recordedAt    (must postdate or equal the
#      rev's commit date — cheap sanity: the baseline file exists and names
#      the same threshold the gate runs at)
#
# Skew means CI enforces a DIFFERENT analyzer than the baseline was recorded
# with → every group reports new (the 105-group false-alarm class).

set -euo pipefail
cd "$(dirname "$0")/.."

PIN_REV="$(sed -n 's/^PINNED_REV="\([^"]*\)".*/\1/p' scripts/pinned-art-dupl.sh)"
CI_REV="$(sed -n 's|.*art-dupl/cmd/art-dupl@\([0-9a-f]*\).*|\1|p' .github/workflows/ci.yaml)"

if [ -z "$PIN_REV" ]; then
	echo "ERROR: could not parse PINNED_REV from scripts/pinned-art-dupl.sh" >&2
	exit 1
fi

if [ -z "$CI_REV" ]; then
	echo "ERROR: could not parse the art-dupl install rev from .github/workflows/ci.yaml" >&2
	exit 1
fi

if [ "$PIN_REV" != "$CI_REV" ]; then
	echo "ERROR: art-dupl rev skew: pinned-art-dupl.sh=$PIN_REV ci.yaml=$CI_REV" >&2
	echo "  Bump BOTH in one commit, then re-record the baseline with that binary" >&2
	echo "  (ritual in docs/adr/0009-accepted-clones.md)." >&2
	exit 1
fi

if [ ! -f .art-dupl-baseline.json ]; then
	echo "ERROR: .art-dupl-baseline.json missing — the clone gate has no baseline" >&2
	exit 1
fi

if ! grep -q '"threshold"' .art-dupl-baseline.json; then
	echo "ERROR: .art-dupl-baseline.json has no threshold field" >&2
	exit 1
fi

echo "art-dupl pin lockstep OK: $PIN_REV (pinned-art-dupl.sh == ci.yaml), baseline present"
