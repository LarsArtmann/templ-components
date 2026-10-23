#!/usr/bin/env bash
# test-coverage-library-gate.sh — proves the library coverage gate can FAIL.
#
# A gate that has never been observed failing is a hope, not a gate (TODO #439).
# Every case feeds a synthetic merged profile through `--profile`, so this runs
# in well under a second and never compiles the library.
#
# Cases:
#   1. below the floor                -> exit 1 + the measured percentage
#   2. at/above the floor             -> exit 0 + PASS line
#   3. generated *_templ.go excluded  -> hand-written 20% still FAILS at gate 90
#                                        (raw would be 92.7% and would pass)
#   4. examples/ excluded            -> 100% hand-written PASSES despite a
#                                        fully-uncovered examples block
#   5. cmd/tc/ excluded              -> same trick for the CLI
#   6. bad arguments                 -> exit 2, never a silent pass
set -euo pipefail

cd "$(dirname "$0")/.."

SCRIPT="scripts/coverage-library.sh"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

out=""

fail() {
	printf '%s\n' "$*" >&2
	echo "FAILED: ${CURRENT:-coverage-library gate self-test}" >&2
	exit 1
}

ok() { printf 'ok: %s\n' "$*"; }

profile() { # profile <file> <block...>
	local file="$1"
	shift
	{
		echo "mode: set"
		printf '%s\n' "$@"
	} >"$file"
}

# run <expected-exit> <label> <command...>
run() {
	local want="$1" label="$2"
	shift 2
	CURRENT="$label"

	local got=0
	out="$("$@" 2>&1)" || got=$?

	[ "$got" = "$want" ] || fail "$label: exit $got, want $want"$'\n'"$out"
	ok "$label"
}

PKG="github.com/larsartmann/templ-components"

# 1. Below the floor: 10 statements, 5 covered = 50.0%.
profile "$TMP/below.out" \
	"$PKG/display/sample.go:1.1,2.2 5 1" \
	"$PKG/display/sample.go:3.1,4.2 5 0"
run 1 "below-floor profile fails the gate" "$SCRIPT" --profile "$TMP/below.out" --gate 94
case "$out" in
*"50.0%"*) ;;
*) fail "below-floor run did not report 50.0%:"$'\n'"$out" ;;
esac
case "$out" in
*"FAIL: hand-written coverage 50.0% is below the 94% gate"*) ;;
*) fail "below-floor run did not print the FAIL line:"$'\n'"$out" ;;
esac

# 2. Above the floor: 10 of 10 covered.
profile "$TMP/above.out" \
	"$PKG/display/sample.go:1.1,2.2 5 1" \
	"$PKG/display/sample.go:3.1,4.2 5 1"
run 0 "above-floor profile passes the gate" "$SCRIPT" --profile "$TMP/above.out" --gate 94
case "$out" in
*"PASS: hand-written coverage 100.0% meets the 94% gate"*) ;;
*) fail "above-floor run did not print the PASS line:"$'\n'"$out" ;;
esac

# 3. Generated templ output must not rescue the hand-written number: 2 of 10
#    hand-written statements covered (20%), 100 of 100 generated covered.
#    Raw would be 102/110 = 92.7% and pass gate 90; the gate must still fail.
profile "$TMP/generated.out" \
	"$PKG/display/sample.go:1.1,2.2 10 2" \
	"$PKG/display/sample_templ.go:1.1,2.2 100 100"
run 1 "generated _templ.go does not lift the hand-written number" \
	"$SCRIPT" --profile "$TMP/generated.out" --gate 90
case "$out" in
*"hand-written (generated excluded):     2/10 = 20.0%"*) ;;
*) fail "generated-exclusion case reported the wrong hand-written total:"$'\n'"$out" ;;
esac

# 4. examples/ is not library code: fully uncovered examples must not fail it.
profile "$TMP/examples.out" \
	"$PKG/display/sample.go:1.1,2.2 10 10" \
	"$PKG/examples/demo/pages.go:1.1,2.2 100 0"
run 0 "examples/ is excluded from the gate" "$SCRIPT" --profile "$TMP/examples.out" --gate 94

# 5. Same for the tc scaffold CLI.
profile "$TMP/cmdtc.out" \
	"$PKG/display/sample.go:1.1,2.2 10 10" \
	"$PKG/cmd/tc/main.go:1.1,2.2 100 0"
run 0 "cmd/tc/ is excluded from the gate" "$SCRIPT" --profile "$TMP/cmdtc.out" --gate 94

# 6. Bad arguments must never look like a pass.
run 2 "--gate rejects a non-numeric value" "$SCRIPT" --gate abc
run 2 "--profile rejects a missing file" "$SCRIPT" --profile "$TMP/does-not-exist.out"
profile "$TMP/empty.out"
run 2 "an empty profile is rejected" "$SCRIPT" --profile "$TMP/empty.out"

echo "PASS: coverage-library gate self-test (all cases behaved as specified)"
