# Datastar ↔ HTMX Parity — Execution Session 3 Status (2026-10-03, ~03:30 CEST)

**Master plan:** `docs/planning/2026-10-02_11-05_datastar-htmx-parity-master-plan.md`
**Prior reports:** execution-1 (`2026-10-02_12-43`), execution-2 (`2026-10-03_00-55`).

## a) What this session executed

All remaining non-owner-gated plan tasks closed. 14 plan tasks landed this
session (L1-06 collateral, L1-07/08/09 completion, L1-12, L1-13, L1-14,
L1-15, L1-16, L1-17, L1-18, L1-19, L1-20, L1-25, L1-26) plus the L2 debts
from sessions 1–2 (dedicated wire tests, demo cards, docs surfaces).

### a.1 — The RED items from session 2, fixed (and root-caused deeper)

1. **Demo busy card (the session-2 RED test)** — endpoint migrated to
   `wire.Handler(PatchTarget{Selector: "#wire-busy-datastar-out"})`, the
   Datastar busy button's inert `Selector` removed, stale comment + 2 tests
   corrected. Fixed BUT the browser e2e exposed TWO deeper demo bugs:
2. **`/wire` and `/kanban` never loaded the Datastar runtime** — only
   `/datastar` rendered `SDKScript`, so EVERY `data-on:*` attribute on the
   wire page was inert markup (every Datastar button dead on the live demo;
   string + HTTP-contract tests cannot see this class). Fix:
   `demoPageMeta.NeedsDatastar` loads the SDK in the shared shell head;
   pinned by `TestWirePageLoadsDatastarRuntime` and browser-proven by
   `TestDemoWireBusyCardBothTransports` (both dialects).
3. **Demo CSP lacked `'unsafe-eval'`** — the pinned Datastar runtime compiles
   action expressions via string eval (`GenerateExpression`; the documented
   requirement in `docs/datastar-runtime-facts.md`), so even WITH the runtime
   loaded every fetch action died with EvalError. The chromedp symptom
   (trusted clicks no-op while synthetic `.click()` worked) cost several
   probe iterations; the console listener nailed it. Fix: `'unsafe-eval'`
   added to `demoCSP` + the firebase.json `/demo/**` override (kept equal).

### a.2 — Plan tasks landed

| Task   | What shipped                                                                                                                                        |
| ------ | --------------------------------------------------------------------------------------------------------------------------------------------------- |
| L2-07/08/09 | `utils/wire/builders_test.go`: constructors table, non-mutation chain pin, `WithFormDefaults` override table, throttle both-dialect rendering, `WithTransport` (NEW chainer), attribute-string contract |
| L1-13  | `display.Tabs` `Wire` — server-side tab switching: `{tab}` URL placeholder (Calendar convention), per-tab clones with `PreventDefault`, empty `Target` defaults to `#ID`, htmx `outerHTML settle:0s` self-swap, `ClientSide` ignored when wired; 6 subtests + 2 goldens + non-mutation pin |
| L1-14  | SimpleNav Wire verified as ALREADY inherited through `NavLinkProps` (probe-proven both dialects); golden `simple_nav_wired_links` pins it            |
| L1-15  | `navigation.LoadMore` migrated to the typed language: `Swap: PatchModeOuter` rendered by the contract (raw `hx-swap` hand render deleted from the wired path), `InfiniteScroll` → typed `Reveal` — infinite scroll now WORKS under Datastar (`data-on-intersect__once`) instead of being ignored; golden renamed `_ignored` → `_reveal` |
| L1-16  | Calendar `settle:0s` rationale confirmed already documented (code godoc + wire guide modifier paragraph) — no change needed                          |
| L1-17  | Datastar confirm: **NO-GO, bundle-verified** (zero confirm tokens in v0.6.1) — `ConfirmDelete` godoc + wire-guide scope row document the asymmetry |
| L1-12  | View transitions DECODED (`useViewTransition` = kebab-cased response-header dataline) and SHIPPED: `PatchTarget.UseViewTransitions` → `wire.Handler` stamps `Datastar-Use-View-Transition: true`; 3-case header test; guide + facts + FEATURES rows |
| L1-18/19 | Demo swap cards: "Swap styles" (append mode, same Action both dialects) + "Remove mode" (region retraction, `hx-swap="delete"` ↔ `PatchModeRemove`) with 2 new endpoints (`/api/wire/swap-line`, `/api/wire/remove-region`) |
| L1-20  | `TestWireDemoSwapCards` + `TestWireSwapEndpoints` (string + contract) and `TestDemoWireSwapCardsBothTransports` (real Chromium: appends twice, retracts the region, both dialects) |
| L1-25  | Real fuzz runs witnessed: `FuzzAction` 2.06M execs, `FuzzDecodeForm` 539k, `FuzzFormWireAttributes` 677k — zero failures                         |
| L1-26  | `TestWireDemoSingleOptionsObject`: demo-render-level guard — every `data-on:*` expression carries at most one `{` (the multi-object regression class) |

### a.3 — Session-1 debts the verify run surfaced (all fixed)

- `cmd/tc` `packageDeps["datastar"]` missing `loading_button.go` +
  `polled_region.go` (the scaffolder's deps listing must match what a
  vendored component needs).
- `internal/contract` `componentTypes()` missing the two new datastar props
  structs (the ComponentProps interface contract must cover every declared
  props type).
- tc `_sources` mirror drifted (tabs/helpers/loadmore) — re-synced via the
  guard's `--fix` (3 files).
- Demo hero count 121 → 123 (`TestHeroCountsMatchFeatures`), index route
  goldens re-captured for the intended hero text change.
- Lint: dead `datastarSwapMode` deleted, wsl/mnd/gocognit/golines/
  exhaustruct findings in the new code fixed (golangci-lint clean on all
  touched modules), mnd constants extracted in `triggers.go`.

### a.4 — Documentation surfaces updated

`docs/transport-wiring.md` (VT header + confirm NO-GO + triggers paragraph +
website mirror), `docs/datastar-runtime-facts.md` (v0.6.1 provenance row —
sha256 re-verified by hand, byte-identical to v0.5.0; VT decode; confirm
NO-GO), `docs/DOMAIN_LANGUAGE.md` (Patch Mode 8 modes + Swap row; the STALE
"Selector overrides response headers" row REWRITTEN to form-lookup; typed
Trigger row added), `skill/SKILL.md` (Tabs row, wire row, interop notes),
`FEATURES.md` (Tabs/Handler rows + golden count 267→270), `CHANGELOG.md`
[Unreleased] (session 1–3 entries: targeting correction as the loud bugfix
entry per the standing recommendation, new APIs, demo MPA infra fixes,
`go 1.26.0` workspace normalization), `AGENTS.md` (daemon-flip note REWRITTEN
— the toolchain ORDERS `1.26` BELOW `1.26.0`, so `.0` is workspace-canonical;
demo runtime-loading + unsafe-eval lessons; wire bullet extended).

## b) Discoveries this session (ranked by blast radius)

1. **`go 1.26` sorts BELOW `go 1.26.0` and the toolchain enforces it twice**:
   (a) workspace mode requires go.work ≥ every module's directive — one
   module at `.0` next to bare-`1.26` go.work kills EVERY workspace build;
   (b) `go mod tidy -diff` demands the `.0` form in some module graphs
   (visualtest's), breaking the cross-module demo-binary build. Session 1's
   "canonical = bare 1.26, flips are harmless" call was wrong at the BUILD
   level (it was true only for the normalized guards). Every go.mod +
   go.work now read `go 1.26.0`. **Supersedes the 2026-10-02 AGENTS note.**
2. **The demo's Datastar surface was entirely dead** (runtime not loaded on
   /wire + /kanban, then eval-blocked by CSP) — invisible to every
   string/HTTP test. Only a real-Chromium click-through catches this class;
   the new e2e tests are the permanent guards.
3. **chromedp trusted-vs-synthetic click divergence** under CSP-eval
   failures: synthetic `.click()` succeeded (expression compiled via a
   non-eval fallback path) while real CDP mouse events threw EvalError.
   Lesson recorded: when clicks "do nothing", capture the console FIRST.
4. **`useViewTransition` is response-header-driven per patch** (not an
   attribute/fetch option) — and the header name is mechanically the
   kebab-cased dataline name, same converter as selector/mode.

## c) Verification state (witnessed)

- `nix run .#verify` — **exit 0, "All checks passed"** (generate + build +
  workspace tests + lint).
- `nix run .#visual -- -parallel 4` — **exit 0, full pass** (the only delta
  was the intended index-hero golden: "121" → "123"; re-captured, then the
  complete suite re-ran green).
- Per-module test suites (utils, htmx, navigation, display, forms,
  examples/demo, cmd/tc, internal/contract, datastar, website): GREEN.
- `scripts/ci-repro.sh --lint --website` — **VERDICT: PASS (exit 0)** at the
  current tip (03:38 CEST). First attempt failed on nix eval-cache
  contention (a concurrent visual run held the SQLite cache — sqlite-busy
  errors corrupted the govulncheck/golangci-lint lanes) plus the visualtest
  go.sum tidy drift and the stale website sales-page golden; all three
  remediated and the clean re-run passed. The lesson (never run ci-repro
  concurrently with `nix run .#visual` — they fight over the nix eval-cache)
  is queued for AGENTS.md.

## d) Owner questions (carried + one new)

1. **History policy (carried, UNANSWERED):** 12+ unpushed daemon snapshot
   commits now mix sessions 1–3. Folding them via rebase is safe (unpushed)
   but needs your explicit OK. Default: leave as-is, push linearly.
2. **Release framing (carried, standing recommendation active):** ship the
   targeting correction + new APIs as the next v1.x release with the loud
   CHANGELOG entry already written. No in-repo consumer relied on the old
   (never-worked) behavior.
3. **htmx.PolledRegion (carried, standing recommendation):** stays
   htmx-native; `datastar.PolledRegion` is its documented twin.
4. **L1-11 (Swap naming — owner-gated) NOT executed:** the plan gates the
   alias/rename decision on you. Current state: `Action.Swap` is htmx-side
   only with godocs saying so everywhere; `PatchTarget.Mode` owns the
   Datastar mode. If you want a different name (e.g. `ClientSwap`), say the
   word and I'll do the rename across code/tests/docs in one pass.

## e) Remaining plan items

- **L1-30 (commit hygiene)**: pending owner answer to (1).
- **ci-repro witness**: queued behind the visual pass (will be witnessed
  before any push; nothing has been pushed this session).
- Post-release items intentionally deferred: `tc` doctor census counts
  (auto-derived), art-dupl baseline re-record if the new test helpers trip
  the clone detector (checked: no new actionable groups — helpers live in
  existing excluded/exempt shapes).
