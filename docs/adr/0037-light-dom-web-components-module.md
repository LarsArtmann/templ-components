# ADR 0037: Light-DOM Web Components Module (`wc`) — Draft for Future Ratification

## Date

2026-09-04

## Status

**Proposed — not ratified.** This ADR is committed as a *draft* so that the
decision it enables is one word, not one session (gate D1 in
`docs/wire-gates-d1-d2-d3.md`). If ratified, it supersedes the "narrow
exception" carve-out of [ADR-0033](0033-web-components-rejection.md) for this
module's scope only. Until then, **nothing here is built** and ADR-0033 stands
in full: the consumer-side light-DOM recipe in `docs/transport-wiring.md` is
the only supported path.

## Context

ADR-0033 permanently rejected the Web Components technology stack (Custom
Elements, Shadow DOM, DSD) for the library, while carving out a narrow
exception: a future **opt-in wrapper project** driven by concrete consumer
demand. Gate D1 asked whether to fire that exception now.

The demand test fails today: the only signal on record is the owner's original
wire-SDK question, and ADR-0036 already answered its real content by
separating the transport axis (real demand — built) from the element-model
axis (no demand). ADR-0035's "nobody has asked for it" test applies verbatim.

## Decision (if ratified)

A new `wc` module, opt-in like `charts/echarts` and `datastar`, with this
narrow shape:

1. **Light-DOM hosts only.** `customElements.define`d elements that render
   `createElement`-style light DOM — never Shadow DOM, never DSD. The Tailwind
   theming model survives because hosts carry no shadow boundary; consumer
   classes and both runtimes' swaps/patches traverse the light DOM unchanged.
2. **Host API, not components.** The module exposes a `Define`/`Host` pair:
   `Define` emits a CSP-safe (nonce'd, singleton-guarded) registration script
   for a named element with a declarative template; `Host` renders the
   element's opening/closing tags with children in between, propagating
   `BaseProps` like every other component. The library ships zero component
   logic inside the elements themselves.
3. **Leaf module.** Depends on `utils` only (BaseProps, wire if needed).
   Layer-1 position in the DAG; released in tag lockstep like every sibling.
4. **Transport-agnostic by construction.** Hosts carry no wiring of their
   own; htmx swaps and Datastar patches inside them work unchanged because
   both runtimes traverse light DOM (ADR-0036's orthogonality result).

## Consequences (if ratified)

**Positive:** the recurring Web Components question gets a final, shipped
answer instead of a recipe; framework-adjacent consumers can adopt the
library's markup inside their own element boundaries.

**Negative:** one more module in the release lockstep (9→10 tags); a
registration-script emitter to keep CSP-clean and singleton-guarded; ongoing
surface to keep honest — for a capability **nobody has asked for**.

## What kills this ADR

If the recipe in `docs/transport-wiring.md` keeps satisfying consumers, this
draft simply expires: close gate D1 as "no" permanently and delete this file.
The recipe documents the identical capability at zero library cost.

## References

- [ADR-0033: Web Components rejection](0033-web-components-rejection.md) —
  the rejection this draft would partially supersede
- [ADR-0036: transport wiring contract](0036-transport-wiring-contract.md) —
  the axes analysis; orthogonality of element model and transport
- `docs/wire-gates-d1-d2-d3.md` — the D1 memo and outcome (not ratified)
- `docs/transport-wiring.md` — the consumer recipe (today's supported path)
