# Nonce/CSP Deepening — Decision Memo (#359, 2026-10-17)

Follow-through decisions on the landed nonce contract. Each item: context,
verdict, and what would reopen it. The landed contract: `utils.ScriptAttrs`
(omit-empty attribute map) + `utils.ScriptComponent` (canonical whole-script
writer, ADR-0047 nonce precedence) pinned by `TestNoEmptyNonceAttribute` /
`TestEmptyNonceStillRendersScripts` (integration) and the demo sweeps
(`TestNoEmptyNonceAcrossDemoPages`, `TestDemoScriptsCarryDemoNonce`).

## 1. #24 Phase 3 — opt-in strict mode (error on empty Nonce)

**Verdict: REJECT for now (⫱G1).** The omit-empty rule exists precisely so
an empty nonce is SAFE: strict-CSP pages block the script either way, and
non-CSP consumers keep it working. A render-time error would break the
documented non-CSP path to guard against a class the output sweeps already
catch. A library panic/error in a render path is the one thing this library
never does (zero-runtime-panics guarantee).

Reopen when: a consumer asks for fail-fast nonce threading in development.
Shape then: a dev-only middleware asserting every script tag on the response
carries the request nonce — test-side, not library-side.

## 2. forbidigo / convention guard banning literal `nonce=""` emitters

**Verdict: REJECT — test-level coverage already owns the class.**
`TestNoEmptyNonceAttribute` sweeps every script-emitting component's OUTPUT;
the demo sweeps cover every rendered page. A forbidigo pattern over
templ-GENERATED Go (`*_templ.go`) is where such a literal would actually
appear — generated code is regenerated, so the lint fights the generator and
the raw `.templ` source rarely spells the literal. The convention is already
enforceable where it matters: new emitters are REQUIRED to render through
`utils.ScriptComponent`/`ScriptAttrs` (AGENTS.md CSP rule), and those two
helpers carry the omit-empty branch in exactly one audited place.

Reopen when: a hand-rolled `nonce=` site actually ships (the demo sweeps
make that a same-day catch, not a slow rot).

## 3. `templ.GetNonce(ctx)` auto-detect when Nonce is empty

**Verdict: RESOLVED (superseded by ADR-0047).** The precedence shipped as
`utils.ResolveNonce(ctx, props.Nonce)`: explicit prop wins → request context
(`ContextWithNonce`) → omit. `ScriptComponent` resolves per render, so one
component value is safe across requests. The Alert-does-others-don't
inconsistency this row flagged is gone at the helper level; any straggler
component that hand-rolls nonce rendering instead of the helpers is the bug
to fix, not the precedence.

## 4. Sweep for uncovered `<style nonce>` emitters

**Verdict: AUDITED — zero sites (2026-10-17).** No library `.templ` source
renders an inline `<style>` tag; `utils.ScriptAttrs` documents the
`<style>` form for when one lands. Nothing uncovered today; the first
`<style>` emitter must render its nonce through `ScriptAttrs` (same
omit-empty rule) and join the integration render table.

## Owner-decision residue

Items 1–2 are closed as rejects-by-default under the standing guardrails
(zero render-path errors; test-level coverage over lint surface). If the
owner wants either flipped, the reopen conditions above are the entry
points — one-line flips, no architectural work.
