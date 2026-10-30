#!/usr/bin/env bash
# check-templ-gate.sh — the templ v0.3.1070-migration GATE PROBE (TODO #335,
# plan C18). Read-only: prints whether the gate is OPEN or CLOSED and, when
# open, the full migration recipe. Never fails on a closed gate — a closed
# gate is the CORRECT state while nixpkgs lags the pin; the script exists so
# the re-check is one command instead of archaeology (it was re-run manually
# three times before this script existed).
#
# The gate OPENS when nixpkgs (the flake's locked input) ships a templ NEWER
# than the go.mod pin — the zero-diff-generate invariant demands the dev
# shell's generator == every go.mod pin, so the pin can only move when the
# shell's binary moves with it.
#
# Forward-compat contract (verified 2026-10-29, crush-daily sweep): templ
# ≥ v0.3.1100 DELETED templ.ResolveAttributeValue, which the committed
# v0.3.1020-generated *_templ.go files call 100+ times per file. Until the
# migration window lands, consumers must pair templ-components ≤ v1.25.0
# with templ ≤ v0.3.1070. The migration MUST therefore regenerate + re-record
# goldens as ONE release event (see the recipe below).

set -euo pipefail
cd "$(dirname "$0")/.."

PIN="$(sed -n 's/.*github.com\/a-h\/templ \(v[0-9.]*\).*/\1/p' go.mod | head -1)"
PIN_BARE="${PIN#v}"

if [ -z "$PIN_BARE" ]; then
	echo "ERROR: could not parse the templ pin from go.mod" >&2
	exit 1
fi

NIXPKGS_TEMPL="$(nix eval --inputs-from . nixpkgs#templ.version --raw 2>/dev/null || true)"

echo "go.mod pin:        $PIN"
echo "nixpkgs templ:     ${NIXPKGS_TEMPL:-<nix eval failed — run inside the flake>}"

if [ -z "$NIXPKGS_TEMPL" ]; then
	echo "GATE: UNKNOWN (nix eval failed — probe needs the flake context)"
	exit 0
fi

if [ "$(printf '%s\n' "$PIN_BARE" "$NIXPKGS_TEMPL" | sort -V | head -1)" = "$PIN_BARE" ] &&
	[ "$PIN_BARE" != "$NIXPKGS_TEMPL" ]; then
	echo ""
	echo "GATE: OPEN — nixpkgs now ships templ $NIXPKGS_TEMPL > pin $PIN_BARE."
	echo "Migration recipe (ONE release event, per TODO #335 + ADR docs):"
	echo "  1. Bump the templ pin in ALL go.mod files + go.work (lockstep, one commit)."
	echo "  2. Regenerate: nix develop -c templ generate ./... from the repo root, then"
	echo "     per-module for every sub-module + website (FileName canonical form)."
	echo "  3. Review icon+label sites: v0.3.1070 emits a separating space after"
	echo "     self-closing @icons — trim sources or accept the goldens' new bytes."
	echo "  4. Re-record HTML goldens (-update) AND the pixel visual goldens"
	echo "     (nix run .#visual -- -update) in the SAME event."
	echo "  5. Verify the generated code no longer references templ.ResolveAttributeValue"
	echo "     (deleted in v0.3.1100): rg -l ResolveAttributeValue --glob '*_templ.go' | wc -l → 0"
	echo "  6. Update the AGENTS.md templ-pin section + close TODO #335 in the release commit."
	exit 0
fi

if [ "$PIN_BARE" = "$NIXPKGS_TEMPL" ]; then
	echo "GATE: CLOSED — nixpkgs matches the pin; nothing to migrate."
else
	echo "GATE: CLOSED — nixpkgs ($NIXPKGS_TEMPL) is OLDER than the pin ($PIN_BARE):"
	echo "the pin moved ahead of nixpkgs; restore the shell/go.mod lockstep instead."
fi
