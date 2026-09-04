# Wire SDK Decision Gates — D1/D2/D3 Memos and Outcomes

**Date:** 2026-09-04
**Context:** `docs/planning/2026-09-04_22-22_wire-sdk-pareto-execution-plan.md` § 5 — three gates were left open by the wire-SDK session because they each require an owner-level call. This file records the memos and the outcomes taken under the standing directive to keep the plan moving; each outcome is reversible by a later owner decision (noted per gate).

---

## D1 — Ratify a superseding ADR-0033 with a light-DOM-only, opt-in `wc` module?

**Options:** yes (build the module) / no (the consumer recipe stays final).

**Memo.** The narrow shape is well understood (see the T17.1 draft ADR-0037, status _Proposed_): light-DOM hosts only, never Shadow DOM, opt-in separate module so consumers who don't want custom elements never see it. The costs: one more module in the release lockstep (9→10 tags), a registration-script emitter to keep CSP-safe, and — the real one — **nobody has asked for it**. ADR-0035's core test ("nobody has asked for it") applies verbatim: the sole demand signal on record is the owner's original wire-SDK question, which ADR-0036 already answered by separating the axes — transport composability is real demand, the element-model axis is not. ADR-0033's narrow exception requires "concrete consumer demand" to fire; none exists.

**Outcome: NOT ratified — the WC module is not built.** The consumer recipe in `docs/transport-wiring.md` remains the documented answer; the draft (ADR-0037) is committed with status _Proposed_ so a future "yes" is one word. T17 phases 2+ and T18 stay gated off. _Reversible by:_ owner ratifying ADR-0037; everything else (scaffold, API, tests) then follows the draft's implementation notes.

---

## D2 — Unknown `Transport` value: silent htmx fallback vs loud failure?

**Options:** silent fallback (repo convention) / error path / panic-in-dev.

**Memo.** The repo has one convention here and it is deliberate: "Zero runtime panics in component code," graceful fallback everywhere (unknown badge styles → default, unknown avatar status → no dot, invalid SwapOOB style → `outerHTML`). Component rendering happens on every request; a typo in a consumer's Transport constant must not 500 a page. The loudness argument is real but belongs at the boundary the consumer controls, not inside a render path.

**Outcome: fallback stays, now PINNED by tests.** `wire.Action.transport()` resolves any unknown value to htmx exactly like the zero value; the invariant pack (`utils/wire/invariants_test.go`) pins cross-dialect URL reference and inert-empty-URL behavior _including_ unknown enum values, and `FuzzAction` holds the no-panic/no-empty-attribute invariant under arbitrary input. Consumers who want loudness can check `wire.TransportIsValid` themselves — it is exported for exactly that. No API change. _Reversible by:_ a new ADR; the fallback is behavior, not surface, so changing it would be a minor-version semantics note.

---

## D3 — Rollout breadth for `Wire` fields across the catalogue?

**Options:** catalogue-wide / pilot-only (Button only) / component-by-component where transport-symmetric.

**Memo.** Catalogue-wide would smear `Wire` over components whose semantics are not transport-symmetric (e.g. `htmx.ConfirmDelete`'s `hx-confirm` has no Datastar expression counterpart without new runtime facts — an ADR-0035-class decision each). Pilot-only would leave the obvious wins unwired: `LoadMore` is already an htmx-shaped `hx-get` button — exactly what `wire.Action` describes. The principle the wire contract itself established (ADR-0036: common subset only) scales to the catalogue: add `Wire` where the component's behavior is expressible in both dialects without new runtime facts; leave transport-specific components in their modules.

**Outcome: transport-symmetric components only, case-by-case.** Adopted cases this round: `navigation.LoadMore` gains `Wire` (T11); the forms server-validation demo uses the wire contract through `wire.Handler` (T12); `htmx.ConfirmDelete` stays htmx-only after the ADR-0035-freeze check — `hx-confirm` has no Datastar twin in the pinned runtime, so a "datastar twin" would be new runtime surface, not wiring (T13 decision note). Each future `Wire` field needs: both-dialect rendering tests + goldens, and a scope note in its docs row. _Reversible by:_ ordinary PR.

---

**Standing effect:** with D1 closed (no), D2 closed (fallback, pinned), D3 closed (symmetric-only rule), the gated workstreams resolve: T17.2+/T18 dropped (with cause), T11–T13 unblocked and executed per the outcomes above.
