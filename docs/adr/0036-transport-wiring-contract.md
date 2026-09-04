# ADR 0036: Transport-Agnostic Wiring Contract (`utils/wire`)

## Date

2026-09-04

## Status

Accepted. Supersedes the "no attribute-helper surface" clause of
[ADR-0035](0035-datastar-scope-freeze.md) (via that ADR's own revisit trigger —
see below). [ADR-0033](0033-web-components-rejection.md) is **reaffirmed, not
superseded**.

## Context

The owner asked: "I want to smartly and composably support ALL of Datastar,
HTMX, and Web Components — does that make sense, and how would we build that
into a proper SDK?"

Two accepted ADRs constrain the answer:

- [ADR-0030](0030-datastar-integration-strategy.md) made Datastar a minimal,
  opt-in complement to the HTMX default.
- [ADR-0035](0035-datastar-scope-freeze.md) then froze the `datastar` module at
  four deliverables ("no attribute-helper surface beyond the existing action
  helpers"), with explicit revisit triggers, and required a superseding ADR
  when one fires.
- [ADR-0033](0033-web-components-rejection.md) permanently rejected the Web
  Components technology stack (Custom Elements, Shadow DOM, DSD) for the
  library, while carving out a "narrow exception" for a future opt-in wrapper
  project driven by concrete consumer demand.

The request conflates three axes that must be separated before any of it makes
sense:

1. **The transport axis** — which client runtime executes a hypermedia
   exchange: htmx (`hx-*` attributes) or Datastar (`data-on:*` +
   `@action()` expressions). htmx and Datastar are *dialects of the same
   operation*: method + URL, triggered by a DOM event, patching a target
   region. This axis is real, and the dialect duplication is paid by every
   consumer who supports both runtimes.
2. **The element-model axis** — whether markup lives in native elements or in
   consumer-defined custom elements. This is *orthogonal* to transport, not a
   third transport: both runtimes traverse light-DOM custom elements without
   modification.
3. **The encapsulation axis** — Shadow DOM/DSD style isolation. ADR-0033's
   rejection applies here in full: a shadow boundary cannot receive the
   consumer's Tailwind utility classes, which breaks the library's entire
   theming model.

## Decision

### 1. Adopt `utils/wire` — one typed Action, two dialects

A new `utils/wire` package defines the transport-agnostic wiring contract:

```go
wire.Action{Transport, Method, URL, Event, Target}.Attributes() templ.Attributes
```

- Transport `""`/`htmx` renders `hx-get`/`hx-post`/…, `hx-trigger`, `hx-target`.
- Transport `datastar` renders `data-on:<event>="@<method>('…')"`.
- Zero values resolve to library defaults (htmx per ADR-0030, GET, dialect
  default event); an empty URL renders nil so unwired components stay inert.
- The contract covers only the dialects' **common subset**. Polling, reveal
  triggers, out-of-band swaps, confirm dialogs, indicators, signals, and SSE
  management stay where they belong: in the `htmx` and `datastar` modules'
  components.

This supersedes ADR-0035's "no attribute-helper surface" clause through that
ADR's own revisit trigger — the owner explicitly requested the capability
(2026-09-04), which is consumer demand in its strongest form. The ADR-0035
freeze on the **`datastar` module itself** remains in force: wire adds no
Datastar components, no runtime surface, and no new module dependency — it is
a leaf contract in `utils` that renders attribute strings.

### 2. Target semantics: encode the verified runtime facts, not a fantasy API

The pinned Datastar v1.0.2 bundle (`go-datastar/static` v0.4.0) accepts **no
`target` option** on fetch actions — verified by reading the bundled source
(see `docs/datastar-runtime-facts.md`). Datastar targeting is response-driven:
for non-SSE HTML the runtime reads `Datastar-Selector` / `Datastar-Mode`
response headers; for SSE the patch events carry the selector; a fragment
root whose id matches the target id patches in the default outer mode.

Therefore `Action.Target` renders **only** for htmx. For Datastar, handlers
honor it by echoing it back: the wire package exports
`HeaderDatastarRequest`, `HeaderHXRequest`, `HeaderDatastarSelector`, and
`HeaderDatastarMode` so one handler can serve both transports by branching on
the request header and setting the response headers. The demo's
`/api/wire/fragment` endpoint implements this and its test pins the contract.

### 3. Web Components: ADR-0033 stands; document the consumer recipe

Supporting Web Components **as a library feature** would require superseding
ADR-0033, whose rationale survives this analysis untouched:

- Shadow DOM/DSD remain incompatible with the Tailwind `@theme` theming model.
- The library's zero-JS identity (native `<dialog>`, `<details>`, Popover API,
  scroll-snap, `field-sizing`) is a competitive feature, not a limitation.
- The distribution problem Web Components solve (framework-agnostic widgets
  for JS consumers) does not exist for a vendored Go library.

What *is* true, and is now documented instead of hinted at: **light-DOM custom
elements defined by consumers compose perfectly with everything this library
ships.** A `customElements.define`d element that renders its children in the
light DOM receives the consumer's Tailwind classes, and htmx swaps and
Datastar patches inside it work unchanged — the axes are orthogonal.
`docs/transport-wiring.md` carries the complete recipe. The library ships no
custom-element code; if concrete consumer demand for a wrapper module
materializes (ADR-0033's narrow exception), that is a separate module with its
own superseding ADR — the narrow shape would be light-DOM hosts only, never
Shadow DOM.

## Consequences

**Positive:**

- Consumers supporting both runtimes write the wiring once; switching a page's
  transport is a one-field change, not a rewrite of every component call.
- The Datastar wire-format knowledge (request marker, response-header
  targeting) moved from tribal demo code into typed, tested constants.
- The recurring "should we support Web Components?" question now has a
  documented answer with the exact conditions under which it would change.

**Negative:**

- One more `utils` surface to document and keep honest (mitigated by the
  deliberate scope freeze to the common subset).
- `Action.Target` being transport-asymmetric is honest but requires the
  documented handler recipe; consumers who skip the recipe will wonder why
  Datastar ignores the target (the guide's first FAQ).

## References

- [ADR-0030: Datastar integration strategy](0030-datastar-integration-strategy.md)
- [ADR-0035: Datastar scope freeze](0035-datastar-scope-freeze.md) — the
  superseded clause and the trigger that fired
- [ADR-0033: Web Components rejection](0033-web-components-rejection.md) —
  reaffirmed
- `docs/datastar-runtime-facts.md` — the verified bundle facts behind the
  Target decision
- `docs/transport-wiring.md` — the consumer guide and the WC recipe
