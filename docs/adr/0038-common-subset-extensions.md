# ADR 0038: Common-Subset Extensions — ContentType and DebounceMS

## Date

2026-09-07

## Status

Accepted. Extends [ADR-0036](0036-transport-wiring-contract.md) (the
transport-agnostic wiring contract); ADR-0036's response-driven-targeting
rule for Datastar stands unchanged. Follows ADR-0036's own amendment
mechanism (an extension is legitimate when both dialects can express it and
the common subset stays honest about it).

## Context

ADR-0036 froze the `wire.Action` common subset at five fields: Transport,
Method, URL, Event, Target. Everything else was deliberately dialect-local
(polling, OOB, indicators, SSE).

Shipping production-grade **forms** on both transports (the 2026-09-07
"Forms × Transports" plan) surfaced two capabilities that broke that
boundary from the consumer side — not as nice-to-haves, but because the
#1 and #2 forms patterns cannot be expressed dialect-agnostically without
them:

1. **Form encoding** (`ContentType`). Datastar's default submission is a
   signals-JSON body; htmx serializes the enclosing form natively. Without
   an encoding knob, a Datastar form cannot send `?q=filter` — it sends
   `?datastar={"q":...}` — and whole-form submission (validation gate,
   submitter name/value, enctype-aware bodies, file uploads) is
   impossible. The pinned v1.0.3 bundle accepts `contentType: 'form'`
   (verified against the embedded bytes, not the docs).
2. **Debounced triggers** (`DebounceMS`). Auto-submit filter inputs are the
   second-most-common forms pattern. htmx expresses debounce as the
   `delay:<n>ms` trigger modifier; Datastar as the `__debounce.<n>ms` event
   key modifier (spelling decoded from the pinned bundle's parser:
   attribute name splits on `__`, groups split on `.`). Duplicating that
   rendering per component would fork the trigger builder and drift.

Both facts are recorded in `docs/datastar-runtime-facts.md` with
bundle-level evidence, per the repo's counterparty-artifact rule.

## Decision

`wire.Action` grows two optional fields; the common subset stays
**closed otherwise**:

- `ContentType ContentType` — `json` (runtime default) or `form`.
  `ContentTypeForm` renders `{contentType: 'form'}` in the Datastar
  expression; htmx ignores it (native serialization). Form-level
  components (`Form`, `FilterInput`, `FilterDropdown.Wire`) default it to
  form encoding so field values travel without configuration.
- `DebounceMS int` — htmx renders `delay:<n>ms` on `hx-trigger` (plus
  `changed` for value events: input, change, keyup); Datastar appends
  `__debounce.<n>ms` to the `data-on:` key. Zero (default) renders
  nothing; negatives are inert.

Both fields render through the same expression/trigger builders as the
original five — one source of truth per dialect.

## Consequences

**Easier:** dual-transport forms, filters, and uploads need no dialect
branching in consumer code. One handler serves both runtimes
(`wire.IsDatastar` is the only branch), and the demo proves each pattern
end-to-end.

**Costs (accepted):** the subset is no longer "five fields" —
documentation and invariants that hardcode the count move with this ADR.
`ContentType` is a no-op under htmx: the field exists so ONE props value
describes ONE exchange, not to change htmx behavior.

**Rejected alternatives:**

- _Per-component dialect attributes_ (`HxGet`-style fields) — already the
  legacy `FilterDropdown` shape; it cannot express Datastar at all and
  forks per component.
- _Signal binding + expression-built URLs_ — expression string
  concatenation is an injection surface and cannot be a complete literal.
- _Extending the subset further (polling, indicators, confirm)_ — remains
  rejected: those have no common semantics to render; the scope boundary
  stays.

## Second Extension (2026-09-07): Selector

> **CORRECTED 2026-10-02 — the override claim below is FALSE for the pinned
> bundle.** The `selector` option exists and picks the form under
> `contentType: 'form'`, but it NEVER reaches patch targeting: the "option
> checked first, header second" code path reads an undefined state key and
> never runs (browser-proven; see the Correction section and
> `docs/datastar-runtime-facts.md`). `wire.Action.Selector` now renders only
> under `ContentTypeForm`.

Datastar v1.0.3 added a client-side fetch option our pinned bundle now
carries: `{selector: '<css>'}` — patch the response into the element(s)
matching the selector, **overriding `Datastar-Selector` response-header
targeting when both are present** (decoded from the bundle's response
dispatcher: the option is checked first, the header second). Under
`contentType: 'form'` the same option additionally selects which form
serializes (`querySelector(sel)` over `closest("form")`).

Decision: `wire.Action.Selector` renders that option, Datastar-only — the
exact twin of `Target` (htmx-only). Response-driven targeting stays the
default when `Selector` is empty; ADR-0036's rule is narrowed, not
replaced: _Target renders for htmx only; Selector renders for Datastar
only; empty Selector keeps response-header targeting authoritative._ The
invariant tests pin both directions, and the busy-state demo endpoint now
runs with zero response-header routing to prove the client-side path.

## Third Extension (2026-10-02): Swap — the region-merge style

Swap is the region-merge style of an exchange (how the response is inserted
relative to the target) — the third and last common-subset gap found by the
2026-09-17 `wire.Action` design review (6+ in-repo sites hand-rolled `hx-swap`
outside the contract). Both dialects express it:

- htmx renders `hx-swap` (`innerHTML`, `outerHTML`, `beforebegin`, `afterbegin`,
  `beforeend`, `afterend`, `delete`, `none`).
- Datastar's fetch actions accept a `mode` option (bundle-verified — see
  `docs/datastar-runtime-facts.md`), with the 8-value set
  `remove/outer/inner/replace/prepend/append/before/after`.

  > **CORRECTED 2026-10-02:** the 8-value MODE SET is real but it is a
  > SERVER-side dataline vocabulary only — there is no client `{mode}` fetch
  > option; the runtime reads the mode exclusively from the
  > `Datastar-Mode` response header. See the Correction section.

**Decision: `Action.Swap PatchMode`** — it reuses the server-side merge-mode
vocabulary (the same enum as `PatchTarget.Mode`), so ONE mode word describes
both the client request and the server targeting. A second dialect-local swap
enum was rejected: `utils/wire` is a leaf module and cannot import the htmx
module's `SwapStyle`, and a parallel enum would give the repo three spellings
of one concept.

Mapping and honesty:

- 1:1 — `inner`→innerHTML, `outer`→outerHTML, `before`→beforebegin,
  `prepend`→afterbegin, `append`→beforeend, `after`→afterend,
  `remove`→delete.
- `replace` is Datastar-only (a morphing replaceWith); under htmx it degrades
  to `outerHTML`, the closest style. Dialect-only modifiers (htmx
  `settle:0s`, Datastar `namespace`) stay outside the contract — components
  that need them keep the raw attribute (see the Calendar MonthNav arrows).
- Zero value renders NOTHING in either dialect. This is deliberate parity:
  htmx's default swap is `innerHTML` and `wire.Handler`'s default mode is
  `inner`, so an unspecified `Swap` behaves identically under both runtimes.
  A set `Swap` renders `hx-swap` for htmx and `{mode: …}` for Datastar (the
  option overrides the `Datastar-Mode` response header when both are present).

  > **CORRECTED 2026-10-02:** a set `Swap` renders `hx-swap` for htmx and
  > NOTHING for Datastar — the client `{mode}` option never existed in the
  > runtime; `wire.Handler(PatchTarget.Mode)` owns the Datastar mode.

### Correction (2026-10-02): client fetch options do not reach targeting

The L1-06 browser e2e (the first REAL-browser test of this ADR's claims)
falsified the Selector-override and client-mode readings both ADR sections
above recorded from static bundle decodes. Full decode + probe evidence:

- The response handler's "options override headers" branch merges over an
  UNDEFINED state key (minified `$`, destructured but never assigned), so it
  never executes; patch datalines come exclusively from `datastar-*`
  response headers.
- The client `selector` option is read in exactly one place: the
  form-encoding branch (`querySelector(sel)` over `closest("form")`).
- Browser proof: `{selector: '#region'}` left a header-less patch inert; a
  client `{mode: 'inner'}` could not stop a `Datastar-Mode: append` header
  from appending; a header-less response whose fragment root carried the
  region's id patched by id-match (outer default).

Library corrections (shipped in the same commit):

- `wire.Action.Selector` renders only under `ContentTypeForm`; godoc and
  guide now direct patch targeting to `wire.Handler(PatchTarget)`.
- `wire.Action.Swap` renders htmx-only (`hx-swap`); Datastar stays
  response-driven — the SAME asymmetry `Target` already has, so ADR-0036's
  rule is again narrowed, not replaced.
- `visualtest/wire_e2e_test.go` pins all four behaviors at the browser
  level — string tests alone encoded the wrong model here; only a real
  runtime falsified it.

### Same-day bug fix: ONE Datastar options object

Bundle decoding for this extension surfaced a latent defect: the runtime's
action dispatcher reads EXACTLY ONE options argument
(`apply(ctx, url, opts)`), so emitting `{selector: …}` and `{contentType: …}`
as separate objects silently dropped form encoding whenever a `Selector` was
also set. `wire.datastarActionExpr` now merges every fetch option
(selector, mode, contentType) into one object literal; the invariant and
both-dialect tests pin the shape. `PatchModeRemove` was added at the same
time — the runtime's mode set has 8 values, not 7.

## Related

- `docs/transport-wiring.md` — the living contract and pattern pack.
- `docs/datastar-runtime-facts.md` — bundle evidence for
  `contentType: 'form'` and the `__debounce` spelling.
- `docs/recipes/server-side-validation.md`, `docs/recipes/file-upload.md`
  — the consumer-facing patterns these fields unlock.
- [ADR-0042](0042-datastar-surface-expansion.md) — where polling/reveal
  triggers and the Datastar component twins live (deliberately NOT common-subset
  fields; see the rejected alternative in the second extension above).
