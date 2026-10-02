# Datastar ↔ HTMX parity session — status report

**Date:** 2026-10-02 10:40
**Owner ask:** "Datastar support is so-so; I want full support for BOTH HTMX and
Datastar and get the BEST out of them."
**Scope:** `utils/wire` common-subset extension + a Datastar correctness fix,
verified end to end. One new ADR section (0038 third extension); no new
component, no new dependency.

---

## a) FULLY DONE

### 1. `wire.Action.Swap` — the region-merge style, both dialects

The 2026-09-17 `Action` design review ranked this its #1 gap ("6+ in-repo sites
hand-roll `hx-swap` outside the contract"). Shipped:

- `Action.Swap PatchMode` reuses the **server-side** `PatchMode` vocabulary
  (`PatchTarget.Mode`), so one word describes both what the client requests and
  what the server targets.
- htmx renders `hx-swap` (`innerHTML`/`outerHTML`/`beforebegin`/`afterbegin`/
  `beforeend`/`afterend`/`delete`); Datastar renders the `{mode: '…'}` fetch
  option, which overrides the `Datastar-Mode` response header.
- Zero value renders NOTHING in either dialect — htmx's default (`innerHTML`)
  and `wire.Handler`'s default (`inner`) are both inner, so an unspecified
  `Swap` behaves identically under both runtimes.
- `replace` is Datastar-only (a morphing `replaceWith`) and degrades to htmx
  `outerHTML`; `remove` ↔ htmx `delete`. Dialect-only modifiers (`settle:0s`,
  Datastar `namespace`) stay out of the contract.
- `display.KanbanBoard` migrated off its hand-rolled `hx-swap="outerHTML"` onto
  `Swap` (HTMX path); its Datastar path stays response-driven as designed.

Files: `utils/wire/wire.go`, `utils/wire/wire_test.go`
(`TestActionSwap`), `utils/wire/invariants_test.go`
(`TestSwapDialectIsolation`, `TestUnspecifiedSwapRendersNothing`),
`display/kanban.go`.

### 2. `wire.PatchModeRemove` — the missing 8th mode

The pinned bundle's mode set is
`["remove","outer","inner","replace","prepend","append","before","after"]` — 8
values. The enum had 7 (`remove` was absent, so `wire.Handler` could not
express a region removal). Added `PatchModeRemove` + `PatchModeIsValid` +
tests.

### 3. FIXED — Datastar dropped every fetch option after the first

Bundle decoding for the swap work surfaced a **latent correctness bug**:
the runtime's action dispatcher reads EXACTLY ONE options argument
(`apply(ctx, url, opts)` via `__action`), so `wire.Action` emitting
`{selector: …}` and `{contentType: 'form'}` as two separate objects silently
**dropped form encoding whenever a `Selector` was also set**. All fetch options
(selector, mode, contentType) now merge into ONE object literal:
`@post('/x', {selector: '#out', mode: 'outer', contentType: 'form'})`.
Pinned by `TestActionSwap` ("…share ONE options object") and the extended
`FuzzAction` (now fuzzes Selector + Swap).

### 4. Research closed (Phase 0 of the prior review)

`docs/datastar-runtime-facts.md` gained three bundle-verified facts:

- The action dispatcher's one-options-object signature (the bug above).
- The `mode` fetch option + its 8-value set, and that non-`outer`/`replace`
  modes REQUIRE a selector (`PatchElementsExpectedSelector`).
- `mode` is also a response-header dataline key, so client option and header are
  interchangeable and the option wins.

### 5. Documentation + guard sync

- `docs/adr/0038-common-subset-extensions.md` — "Third Extension: Swap" + the
  bug-fix section.
- `docs/transport-wiring.md` — zero-values table, dialect mapping table, a new
  "Swap / merge mode" section with the mapping table, and the one-object rule.
- `docs/adr/0036-transport-wiring-contract.md` — 7→8 modes.
- `FEATURES.md` — wire table (Selector, Swap, PatchMode rows) + corrected the
  stale "selector option deliberately unadopted" scope note.
- `AGENTS.md` — corrected the same stale claim; documented `Swap` + the
  one-object rule.
- `skill/SKILL.md` — `Action` fields, `PatchMode` 8 modes, one-object note.
- `website/content/docs/guides/transport-wiring.md` +
  `website/content/docs/api-reference.md` — same corrections.
- `CHANGELOG.md` `[Unreleased]` — Added (Swap) + Fixed (options object).

### 6. Repo hygiene: go directive drift repaired

A daemon auto-commit (`67fa7b3b`) had flipped `visualtest/go.mod`
`go 1.26.0` → `go 1.26`, which made the workspace unloadable and failed
`TestGoWorkDirectiveMatchesRootGoMod`. Restored the intended consistent state
(commit `60805f09` had aligned root to `1.26.0`): root `go.mod`, `visualtest/go.mod`,
and `go.work` all `1.26.0`. **The daemon keeps flipping these — re-check after
daemon commits.**

---

## b) VERIFIED

- `go build ./...` (workspace), root `go test ./...`, and the per-module loop
  (`utils icons errorpage charts/echarts datastar htmx`) — all green.
- `visualtest` compile/skip check — green (41s; browser tests skip without the
  nix Chromium env, the demo-binary build itself is exercised).
- `website` module tests — green.
- `golangci-lint` — root module + utils module: 0 issues.
- `gofmt` — clean.

Untouched-by-this-session failures: none observed.

---

## c) NOT DONE (ranked roadmap for true "best of both")

| # | Task | Impact | Effort |
|---|------|--------|--------|
| 1 | `datastar.PolledRegion` — interval polling (`data-on-interval__duration.10s[.leading]`, bundle-decoded) to match `htmx.PolledRegion`. The single biggest Datastar gap; needs a browser e2e (interval firing) BEFORE shipping | High | M |
| 2 | `datastar.LoadingButton` — signal-indicator label swap, parity with `htmx.LoadingButton` | Medium | S |
| 3 | Typed interval/intersect triggers in the wire contract (TODO #178) — `hx-trigger="every Ns"/"revealed"` ↔ `data-on-interval`/`data-on-intersect`; ADR-sized, bundle-verified | High | M |
| 4 | Adopt `Wire` on more transport-symmetric components (D3 rule): `display.Tabs`, `forms.SimpleNav` (TODO #155) | Medium | M |
| 5 | Migrate the remaining `hx-swap` sites: `navigation/loadmore.templ` (self-targeting), leave `forms/calendar_nav.go` (`settle:0s` is htmx-only) | Low | S |
| 6 | `wire.Get/Post/Put/Patch/Delete(url)` constructors + promoted `WithEvent`/`WithContentType`/`WithDebounce` builders (review items 4–5) | Medium | S |
| 7 | `ThrottleMS` sibling of `DebounceMS` (`throttle:Nms` ↔ `__throttle.Nms`, spelling bundle-verified) | Medium | S |
| 8 | Demo card showing `Swap` under both transports | Low | S |

**Waiver note:** item 1 requires a Chromium e2e; the interval spelling is
bundle-decoded but not browser-proven in this session, so it was deliberately
NOT shipped — shipping an unproven runtime claim is the exact inert-integration
failure class (`docs/datastar-runtime-facts.md` provenance rule).

---

## d) RECOMMENDED NEXT SESSION

Start at roadmap #1 (`datastar.PolledRegion`) with a browser e2e, then #3
(trigger language) as its ADR. Both unlock the largest remaining Datastar
depth. #6/#7 are quick ergonomic wins.
