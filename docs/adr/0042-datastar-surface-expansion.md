# ADR 0042: Datastar Surface Expansion — Polling, Busy, and Typed Triggers

## Date

2026-10-02

## Status

Accepted. **Supersedes the scope freeze of
[ADR-0035](0035-datastar-scope-freeze.md)** via that ADR's own revisit trigger
(explicit owner request) for the surfaces named below. ADR-0035's other
provisions — opt-in separate module, zero transitive deps beyond
`go-datastar/static`, honest positioning against the HTMX default — remain in
force. [ADR-0036](0036-transport-wiring-contract.md) (wire contract) and
[ADR-0038](0038-common-subset-extensions.md) (common-subset extensions) are
extended, not replaced.

## Context

ADR-0035 (2026-08-14) froze the `datastar` module at four deliverables with
three explicit revisit triggers, requiring a superseding ADR when one fires.
Since then the transport story moved three times (ADR-0036 wiring contract,
ADR-0038 extensions, the 2026-10-02 parity session), and each round made the
same gap louder: **the `htmx` module ships 9 components while `datastar`
ships 4, and Datastar users have no polling primitive at all** — the single
largest interaction class Datastar consumers cannot express with this
library. `htmx.PolledRegion` has no Datastar counterpart; busy-state buttons
(`htmx.LoadingButton`) have none either; and `hx-trigger="every Ns"` /
`"revealed"` have no typed wire-level twin.

On 2026-10-02 the owner requested **"full support for BOTH HTMX and
Datastar — get the BEST out of them"** and authorized execution of the parity
master plan (`docs/planning/2026-10-02_11-05_datastar-htmx-parity-master-plan.md`).
That is consumer demand in ADR-0035's strongest recorded form — the same
trigger class that legitimately lifted the "no attribute-helper surface"
clause in ADR-0036.

## Decision

The ADR-0035 freeze is lifted for exactly four surfaces, each bundle-verified
against the pinned `go-datastar/static` runtime and browser-proven before
merge:

1. **`datastar.PolledRegion`** — interval polling via
   `data-on-interval__duration.<n>[.leading]` (bundle-decoded v0.6.1: default
   1000ms; `duration` parses `ms`/`s`/bare-ms; `leading` must be an element of
   the duration mod set and fires once immediately; each tick is wrapped in
   the runtime's busy markers; `setInterval` + teardown cleanup). The Datastar
   twin of `htmx.PolledRegion`, sharing its props philosophy (URL, interval,
   eager, aria-live, optional timestamp).
2. **`datastar.LoadingButton`** — busy-state button that disables itself and
   swaps its label while its action is in flight, using the runtime's fetching
   signal — the Datastar twin of `htmx.LoadingButton`.
3. **Typed interval/intersect triggers in `wire`** (TODO #178) —
   `every Ns` ↔ `data-on-interval__duration.Ns` and
   `revealed` ↔ `data-on-intersect` render through one typed contract, so
   polling/reveal wiring stops being dialect-local boilerplate at every call
   site. Bundle-decoded `on-intersect` mods: `full`/`half`/`threshold.<0-100>`,
   `exit`, `once`.
4. **Confirm-parity evaluation** — research whether the pinned runtime can
   express `ConfirmDelete`'s gate; ship only if it degrades gracefully without
   JS hacks (expression-level gating, no `confirm()` interception).

**Still out of scope (the freeze survives here):** full HTMX feature parity
for its own sake, signals/state management helpers, OOB-swap equivalents
(response-driven patching stays Datastar's model), and any expansion beyond
the four surfaces above. New Datastar components still require a superseding
ADR or an explicit owner request.

**Quality bar (unchanged by this ADR):** every runtime claim is pinned by a
bundle-contract guard (`datastar.TestPinnedRuntimeBundleContract` lineage) and
proven in real Chromium via `visualtest` before merge; string assertions on
our own output are not proof. Zero-value props stay inert. Zero runtime
panics.

## Consequences

**Easier:** Datastar consumers reach interaction parity on poll/busy/reveal —
the three patterns the HTMX module owns today — without leaving the typed
contract. The wire trigger language removes per-component dialect forks for
the polling class.

**Costs (accepted):** the `datastar` module's "exactly four deliverables"
positioning statement (website guide, ADR-0035 citations) is now "four
deliverables plus the parity surfaces named here" — those texts move with
this ADR. Each new component carries the full lens burden (golden, a11y,
bdd, example, e2e) — roughly double the test surface per component vs the
frozen scope.

**Rejected alternatives:**

- _Stay frozen_ — fails the owner's explicit request and leaves the #1 gap
  open while the wire contract keeps proving that parity is cheap to express.
- _Fold Datastar polling into `wire.Action` as a common-subset field_ —
  rejected: htmx polling (`hx-trigger="every"`) and Datastar polling
  (`data-on-interval`) belong to different elements' wiring and different
  components (a polled REGION, not a single action). The trigger language
  (item 3) covers the shared spelling without pretending one `Action` is the
  right shape for a self-polling region.
- _Deprecate the `datastar` module_ — rejected: the owner's request is the
  opposite direction.

## Related

- [ADR-0035](0035-datastar-scope-freeze.md) — the freeze this ADR supersedes
  (annotation added there).
- [ADR-0036](0036-transport-wiring-contract.md),
  [ADR-0038](0038-common-subset-extensions.md) — the wire contract and its
  extension mechanism.
- [ADR-0043](0043-typed-trigger-language.md) — the full ADR for surface 3
  (typed interval/intersect triggers).
- `docs/datastar-runtime-facts.md` — the bundle evidence for every fact cited
  above.
- `docs/planning/2026-10-02_11-05_datastar-htmx-parity-master-plan.md` — the
  execution plan this ADR unblocks.
