# ADR 0043: Typed Trigger Language for `wire` — Interval and Reveal

## Date

2026-10-02

## Status

Accepted. Extends [ADR-0036](0036-transport-wiring-contract.md) (the one-spec,
two-dialects rule) onto the trigger axis and executes surface 3 of
[ADR-0042](0042-datastar-surface-expansion.md). Supersedes the "triggers stay
out of the contract for now" research note in
[docs/transport-wiring.md](../transport-wiring.md) (2026-09-04).

## Context

htmx and Datastar both express two trigger sources beyond DOM events —
**interval polling** and **viewport reveal** — with bundle-verified grammars:

| Concept          | htmx 2.0.10                                        | Datastar v0.6.1 (pinned bundle)                                        |
| ---------------- | -------------------------------------------------- | ---------------------------------------------------------------------- |
| Polling          | `hx-trigger="every 5s"`                            | `data-on-interval__duration.5s` (duration lives in the attribute NAME) |
| Reveal once      | `hx-trigger="revealed"` (intersect-once shorthand) | `data-on-intersect__once`                                              |
| Reveal every     | `hx-trigger="intersect"`                           | `data-on-intersect` (no `__once`)                                      |
| Threshold        | `intersect threshold:0.5`                          | `__half` (50) / `__full` (100) / `__threshold.<0-100>` (percent)       |
| Leave viewport   | — (no exit semantics)                              | `__exit` (bundle-verified)                                             |
| Multiple sources | one comma-separated `hx-trigger` value             | one attribute PER source (each plugin owns an attribute)               |

Two runtime facts shape any typed model (see
[docs/datastar-runtime-facts.md](../datastar-runtime-facts.md)):

1. **Neither parser understands minutes or hours.** htmx's `parseInterval`
   and Datastar's duration parser both fall through to `parseFloat` and read
   the leading digits as MILLISECONDS — `"5m"` polls five times a
   millisecond. A typed layer must normalize (`5m` → `300s`) or it launders a
   footgun through both dialects.
2. **The trigger spelling is not attribute-VALUE syntax on the Datastar
   side.** The interval duration is encoded in the attribute NAME via the
   runtime's `__` modifier grammar, so rendering must go through dynamic
   `templ.Attributes` keys, not string interpolation into a fixed attribute.

Until now `wire` deliberately left this axis dialect-local (ADR-0038's
"common subset" stopped at `Event` + `DebounceMS`). That forced every
polling/reveal call site to fork per transport — exactly the boilerplate
ADR-0036 exists to remove.

## Decision

Two typed fields on `wire.Action` — **not** an `Event` enum extension:

```go
Action{
    URL:      "/api/stats",
    Interval: "10s",       // every 10s  <->  data-on-interval__duration.10s
    // or
    Reveal: &wire.Reveal{ThresholdPercent: 50}, // revealed/intersect <-> data-on-intersect
}
```

- **`Action.Interval string`** — any duration the runtime parsers accept,
  with `m`/`h` normalized to seconds by `wire.NormalizedInterval` (exported
  so `datastar.PolledRegion` shares one normalizer). Empty string = no
  polling.
- **`Action.Reveal *Reveal`** — pointer (nil = no reveal trigger), with the
  **inverted `EveryEntry` flag**: the zero value fires ONCE on first entry
  (htmx `"revealed"` / Datastar `__once`), because fire-once lazy-load is the
  common case and the zero value must equal the documented default.
  `ThresholdPercent` maps to `half`/`full`/`threshold.<n>` at the special
  values, `threshold:<fraction>` under htmx (only when > 0). `Exit` renders
  Datastar `__exit`; under htmx it degrades to the entry trigger —
  documented, mirroring `PatchModeReplace` → `outerHTML`.
- **Merging rules:** htmx merges ALL trigger sources (event, reveal, interval)
  into ONE comma-separated `hx-trigger`; Datastar renders one attribute per
  source. An interval-only or reveal-only Action (Event unspecified) renders
  NO click/event attribute — the action must not also fire on click.
- **Debounce keeps its existing semantics** (`delay:Nms` + auto-`changed` on
  input-like events under htmx); Datastar debounce stays out of this ADR —
  the runtime has no declarative debounce modifier, so no cross-dialect
  spelling exists to model.

### Rejected alternatives

1. **Extend `wire.Event`** with `EventInterval`/`EventReveal` — events cannot
   carry parameters (duration, threshold, once/exit). It would collapse into
   magic strings on neighboring fields anyway, with worse validation.
2. **Keep triggers dialect-local** (the 2026-09-04 status quo) — preserves
   per-component forks; the parity plan's whole point is removing them.
3. **A separate `wire.Trigger` struct list** — YAGNI with exactly two
   non-event sources; fields on `Action` keep zero-value rendering inert and
   compose with the existing `htmxTrigger()` merge without a second API
   surface.

## Consequences

**Easier:** polling and reveal wiring is one spec again — `datastar.PolledRegion`,
`navigation.LoadMore` (infinite scroll), and consumer call sites stop
duplicating trigger spelling per dialect. Interval normalization lives in one
place (`wire.NormalizedInterval`), closing the `"5m"` misparse for every
caller at once.

**Costs / limits (documented, deliberate):**

- `Reveal.Exit` is Datastar-only; htmx degrades to entry-firing. Consumers who
  need leave-detection under htmx need custom JS (out of scope).
- `leading` (fire-immediately-then-wait) is intentionally NOT modeled on the
  wire trigger: combined with an outer self-patching region it refetches
  forever (the same loop that killed `datastar.PolledRegion`'s proposed
  `Eager` prop).
- Datastar renders one action expression PER trigger source; an Action with
  both Interval and Reveal performs two independently-triggered exchanges of
  the same endpoint (identical to htmx's comma-joined semantics, but visible
  as two attributes in the DOM).

**Verification:** `utils/wire/triggers_test.go` pins both dialect spellings
token-by-token (dynamic attribute keys included); the runtime grammars are
bundle-decoded in `datastar.TestPinnedRuntimeBundleContract`'s token set
(`on-interval`, `duration`, `leading`, `on-intersect`); browser proof lands
with `visualtest/wire_e2e_test.go` extensions.

## Related

- [ADR-0036](0036-transport-wiring-contract.md) — one Action, two dialects
- [ADR-0038](0038-common-subset-extensions.md) — what else crossed the
  dialect line (selector targeting, swap modes)
- [ADR-0042](0042-datastar-surface-expansion.md) — the four-surface expansion
- [docs/datastar-runtime-facts.md](../datastar-runtime-facts.md) — pinned
  bundle grammar for `on-interval` / `on-intersect`
