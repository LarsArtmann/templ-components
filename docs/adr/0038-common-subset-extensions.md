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

- *Per-component dialect attributes* (`HxGet`-style fields) — already the
  legacy `FilterDropdown` shape; it cannot express Datastar at all and
  forks per component.
- *Signal binding + expression-built URLs* — expression string
  concatenation is an injection surface and cannot be a complete literal.
- *Extending the subset further (polling, indicators, confirm)* — remains
  rejected: those have no common semantics to render; the scope boundary
  stays.

## Related

- `docs/transport-wiring.md` — the living contract and pattern pack.
- `docs/datastar-runtime-facts.md` — bundle evidence for
  `contentType: 'form'` and the `__debounce` spelling.
- `docs/recipes/server-side-validation.md`, `docs/recipes/file-upload.md`
  — the consumer-facing patterns these fields unlock.
