# 100 Improvement Ideas — templ-components

**Created:** 2026-10-01 | **Baseline:** v1.19.4
**Scope:** fresh ideas grounded in the current repo state (7-module workspace, 121 components, dual-transport `wire`, visualtest/axe/golden tiers, `tc` CLI, Go SSG website). Items already tracked in `TODO_LIST.md` / `ROADMAP.md` are marked ⚑ so you can skip re-deciding them.

---

## How this list is sorted

Every idea carries two 1–5 scores:

| Axis       | Meaning                                                                                        |
| ---------- | ---------------------------------------------------------------------------------------------- |
| **Impact** | 5 = transforms adoption/reliability, 4 = significant, 3 = solid, 2 = incremental, 1 = marginal |
| **Effort** | 1 = <1h, 2 = half-day, 3 = 1–2 days, 4 = ~a week, 5 = multi-week/session                       |

Ideas are ranked by a Pareto-lite order: **quick, high-leverage wins first**, then big bets, then solid improvements, then nice-to-haves. Within each tier, cheapest (lowest effort) comes first; ties break toward higher impact.

| Tier | Definition                          | Ranks  |
| ---- | ----------------------------------- | ------ |
| P1   | Quick wins — impact ≥ 4, effort ≤ 2 | 1–6    |
| P2   | Big bets — impact ≥ 4, effort ≥ 3   | 7–22   |
| P3   | Solid improvements — impact = 3     | 23–69  |
| P4   | Nice-to-have / policy — impact ≤ 2  | 70–100 |

⚑ = already referenced in a tracking doc (see the Tracked column).

---

## If you only do ten

1. Submit to `awesome-templ` + `templ.guide` listings (#1, ⚑ #28/#29)
2. Per-module `go build` gate + torn-snapshot tripwire for the auto-commit daemon (#2, #3)
3. Cross-link the GOTH-stack story across the three libraries (#4)
4. Determinism gate (render twice, byte-compare) (#5)
5. ConfirmDialog on native `<dialog>` (#6)
6. Explicit `HTMXOff` (#7, ⚑ #282)
7. `tc doctor` single-shot health check (#8)
8. Post-tag consumer compile smoke in CI (#9, ⚑ #15)
9. `ExampleXxx` for every component (#11)
10. Copy-paste playground + real-world example app (#19, #20)

---

## P1 — Quick wins (impact ≥ 4, effort ≤ 2)

| Rank | Idea                                                                                       | Theme      | Impact | Effort | Tracked   |
| ---- | ------------------------------------------------------------------------------------------ | ---------- | :----: | :----: | --------- |
| 1    | Submit the library to `awesome-templ` and the `templ.guide` directory listing              | Ecosystem  |   4    |   1    | ⚑ #28/#29 |
| 2    | Per-module `go build` gate before any daemon auto-commit                                   | CI/Release |   4    |   2    | ⚑ #232    |
| 3    | Daemon torn-snapshot tripwire: refuse commits when files changed within N seconds          | CI/Release |   4    |   2    | ⚑ #232    |
| 4    | Cross-link the three libraries as a branded "GOTH stack" with badges + one narrative       | Docs       |   4    |   2    | ⚑         |
| 5    | Determinism gate: render twice, byte-compare raw HTML, fail on map-iteration/time drift    | Testing    |   4    |   2    | ⚑ #45     |
| 6    | `ConfirmDialog` on native `<dialog>` using `dialog.returnValue` for promise-style confirms | Components |   4    |   2    |           |

## P2 — Big bets (impact ≥ 4, effort ≥ 3)

| Rank | Idea                                                                                          | Theme      | Impact | Effort | Tracked |
| ---- | --------------------------------------------------------------------------------------------- | ---------- | :----: | :----: | ------- |
| 7    | Explicit `HTMXOff` value for pages with no wired components (kill the CDN-means-on ambiguity) | Wire       |   4    |   3    | ⚑ #282  |
| 8    | `tc doctor`: verify hooks path, templ pin, CSS freshness, go.work, replace directives at once | DX         |   4    |   3    |         |
| 9    | Post-tag consumer compile smoke in CI (throwaway module `go get`-ing the new tag)             | CI/Release |   4    |   3    | ⚑ #15   |
| 10   | Release-lock signal while `release.sh` runs so the daemon cannot race a cut                   | CI/Release |   4    |   3    | ⚑ #232  |
| 11   | `ExampleXxx` for every component so `pkg.go.dev` documents usage                              | Ecosystem  |   4    |   3    |         |
| 12   | Make `Validate() error` a uniform props-struct convention with one shared harness             | API        |   4    |   4    |         |
| 13   | Component registry `tc.Catalog()` (name + signature + package) feeding docs, goldens, CLI     | API        |   4    |   4    |         |
| 14   | Run the axe sweep over component goldens (static HTML audit lane, not just demo routes)       | A11y       |   4    |   4    |         |
| 15   | Firefox visual lane for Popover / `field-sizing` / `base-select` degradation                  | Testing    |   4    |   4    | ⚑ #43   |
| 16   | Prebuilt CDN-class CSS/JS bundle path for no-build consumers                                  | Perf       |   4    |   4    |         |
| 17   | `nix run .#new`: scaffold a fresh project from the starter template                           | DX         |   4    |   4    |         |
| 18   | Command-palette recipe composing Dropdown + Combobox + Popover (no new component)             | Components |   4    |   4    |         |
| 19   | Copy-paste component playground on the site (click, render, copy Go)                          | Docs       |   5    |   5    |         |
| 20   | Real-world example app (templ-components + cqrs-htmx + go-cqrs-lite) doubling as demo site    | Docs       |   5    |   5    |         |
| 21   | Real-runtime `datastar`/`ssetest` e2e module driving the pinned runtime in a browser          | Wire       |   4    |   5    | ⚑       |
| 22   | `tc migrate`: bump consumers across versions and flag breaking changes                        | DX         |   4    |   5    |         |

## P3 — Solid improvements (impact = 3)

### Effort 1

| Rank | Idea                                                                            | Theme     | Impact | Effort | Tracked |
| ---- | ------------------------------------------------------------------------------- | --------- | :----: | :----: | ------- |
| 23   | Attrs-precedence contract test: consumer `Attrs`/`Class` must always win        | API       |   3    |   1    |         |
| 24   | Guard: no `role="combobox"` on non-text inputs across all components            | A11y      |   3    |   1    | ⚑ #274  |
| 25   | Fresh-clone hook check in CI/doctor (`core.hooksPath` actually set)             | CI        |   3    |   1    |         |
| 26   | Discussion template to formalize the consumer adoption-ask / demand-survey loop | Ecosystem |   3    |   1    |         |

### Effort 2

| Rank | Idea                                                                                          | Theme      | Impact | Effort | Tracked |
| ---- | --------------------------------------------------------------------------------------------- | ---------- | :----: | :----: | ------- |
| 27   | Weekly proxy-lag monitor: proxy `@latest` == newest local tag per sub-module                  | CI/Release |   3    |   2    | ⚑       |
| 28   | Publish the "Optimistic UI with htmx and Datastar" case study (ADR-0041)                      | Docs       |   3    |   2    | ⚑       |
| 29   | Scanner guard banning remaining `map[string]X` lookups (typed-enum keys everywhere)           | API        |   3    |   2    |         |
| 30   | Scoped-ID helper so auto IDs never collide across nested components after swaps               | API        |   3    |   2    |         |
| 31   | Segmented control (iOS-style), built on existing radio + peer patterns                        | Components |   3    |   2    |         |
| 32   | `CopyField` (readonly input + CopyButton) for API keys/tokens                                 | Components |   3    |   2    |         |
| 33   | `<output>`-based live computed total for forms                                                | Components |   3    |   2    |         |
| 34   | Bottom-sheet `Drawer` variant for mobile (`data-side="bottom"`)                               | Components |   3    |   2    |         |
| 35   | Extract a reusable flaky-endpoint (stall/fail) visualtest helper for stateful components      | Wire       |   3    |   2    | ⚑ #262  |
| 36   | Reduced-motion visual lane proving animations are actually suppressed                         | A11y       |   3    |   2    |         |
| 37   | Golden-orphan detector: fail on a golden file with no generating test                         | Testing    |   3    |   2    | ⚑ #54   |
| 38   | Promote the art-dupl check to a blocking CI job (after two green advisory runs)               | Testing    |   3    |   2    | ⚑ #295  |
| 39   | CSS size budget guard (bytes per module) that fails on regression                             | Perf       |   3    |   2    |         |
| 40   | `tc add --dry-run` showing exactly which files would be mirrored/copied                       | DX         |   3    |   2    |         |
| 41   | Editor support pack: documented VS Code / Neovim config (templ LSP, tailwind, format-on-save) | DX         |   3    |   2    | ⚑       |

### Effort 3

| Rank | Idea                                                                                         | Theme      | Impact | Effort | Tracked |
| ---- | -------------------------------------------------------------------------------------------- | ---------- | :----: | :----: | ------- |
| 42   | Cross-module compatibility matrix (which `utils`/`icons`/root versions may be mixed) + guard | API        |   3    |   3    |         |
| 43   | Promote the `Card.Body`/`Table.Body` slot pattern into a documented generic `Slot` type      | API        |   3    |   3    |         |
| 44   | Timeline / activity-feed component (common admin pattern, not yet present)                   | Components |   3    |   3    |         |
| 45   | Layout-preserving skeleton set generated from real component markup                          | Components |   3    |   3    |         |
| 46   | Document the optimistic/pending register as a general `wire` pattern                         | Wire       |   3    |   3    | ⚑       |
| 47   | Schedule the htmx v4 event-rename audit as a spike, not an ad-hoc bump                       | Wire       |   3    |   3    | ⚑ #316b |
| 48   | Extend the touch-target audit to every interactive golden, not just demo routes              | A11y       |   3    |   3    |         |
| 49   | Accessible-name computation audit as a golden-tier check                                     | A11y       |   3    |   3    |         |
| 50   | `prefers-contrast` / `forced-colors` visual variants beyond class-level checks               | A11y       |   3    |   3    |         |
| 51   | Localized demo route (`dir="rtl"` + `lang`) to browser-prove logical-property mirroring      | A11y       |   3    |   3    |         |
| 52   | Publish an accessibility conformance statement (WCAG 2.2 AA) derived from the axe ledger     | A11y       |   3    |   3    |         |
| 53   | Chart geometry property-based tests (monotonic ticks, arc closure)                           | Testing    |   3    |   3    | ⚑ #47   |
| 54   | Measure and publish per-component HTML byte cost (tiny/average/heavy)                        | Perf       |   3    |   3    |         |
| 55   | Inline-JS byte budget per component + drift failure                                          | Perf       |   3    |   3    |         |
| 56   | `@layer` ordering so consumer overrides never need `!important`                              | Perf       |   3    |   3    |         |
| 57   | Hot-reload dev app watching templ + Go + Tailwind in one command                             | DX         |   3    |   3    | ⚑       |
| 58   | One-command local Postgres for the starter app (migrations + seed + test DB)                 | DX         |   3    |   3    | ⚑       |
| 59   | Go OG-image generator per page (retire frozen `public/og/*`)                                 | Docs       |   3    |   3    | ⚑       |
| 60   | Docs TOC scroll-spy + active header states + mobile docs nav below `lg`                      | Docs       |   3    |   3    | ⚑       |
| 61   | Firebase per-PR preview channel (also cleanUrls/CSP live-verification vehicle)               | Docs       |   3    |   3    | ⚑       |
| 62   | "Who uses this" consumers section with a version-pin table                                   | Ecosystem  |   3    |   3    |         |

### Effort 4

| Rank | Idea                                                                                    | Theme      | Impact | Effort | Tracked |
| ---- | --------------------------------------------------------------------------------------- | ---------- | :----: | :----: | ------- |
| 63   | Typed Go `Theme`/`Tokens` struct mapping to `@theme` CSS vars (compile-checked theming) | API        |   3    |   4    |         |
| 64   | DataTable column visibility + resize (CSP-safe, client-side)                            | Components |   3    |   4    |         |
| 65   | Focus-order goldens for overlays (capture tab order as data, diff it)                   | A11y       |   3    |   4    |         |
| 66   | Lazy-load / split icon path data so icons-only consumers tree-shake                     | Perf       |   3    |   4    |         |
| 67   | v2 module-path migration ADR + a codemod for the import switch                          | Ecosystem  |   3    |   4    |         |

### Effort 5

| Rank | Idea                                                                                 | Theme | Impact | Effort | Tracked |
| ---- | ------------------------------------------------------------------------------------ | ----- | :----: | :----: | ------- |
| 68   | Typed interval/intersect triggers in `wire.Action` (mini trigger-language ADR)       | Wire  |   3    |   5    | ⚑ #178  |
| 69   | Hosted playground / StackBlitz-style live Go+templ environment for zero-install eval | Docs  |   3    |   5    |         |

## P4 — Nice-to-have / policy (impact ≤ 2)

### Effort 1

| Rank | Idea                                                                             | Theme   | Impact | Effort | Tracked |
| ---- | -------------------------------------------------------------------------------- | ------- | :----: | :----: | ------- |
| 70   | Document drop-vs-queue (`hx-sync`) guidance for rapid multi-action users         | Wire    |   2    |   1    |         |
| 71   | Policy: any newly wired component ships both-dialect goldens in the same plan    | Wire    |   2    |   1    |         |
| 72   | Policy note on test-only dependencies (so future e2e modules don't re-litigate)  | Testing |   2    |   1    |         |
| 73   | Make the demo smoke an explicit gate in the plan-authoring checklist             | Testing |   2    |   1    |         |
| 74   | Make `SITE_SKIP_STARS=1` the default for all non-production dist entry points    | Perf    |   2    |   1    | ⚑ #271  |
| 75   | `tc ls` derived "N components addable" footer (guarded)                          | DX      |   2    |   1    | ⚑ #294  |
| 76   | `.#website` flake app wrapping `build.sh`                                        | Docs    |   2    |   1    | ⚑       |
| 77   | Failure-screenshot naming convention + CI cleanup for `testdata/.fail/`          | CI      |   2    |   1    |         |
| 78   | Retire or wire `templates/styles.css` + `theme.out.css` (dead artifact decision) | CI      |   2    |   1    | ⚑ #211  |

### Effort 2

| Rank | Idea                                                                              | Theme      | Impact | Effort | Tracked |
| ---- | --------------------------------------------------------------------------------- | ---------- | :----: | :----: | ------- |
| 79   | Content-hash the compiled CSS, expose as a `data-css-build` attribute             | API        |   2    |   2    |         |
| 80   | Collapse `DefaultXxxProps` boilerplate with a small `withDefaults` generic        | API        |   2    |   2    |         |
| 81   | Pre-write demand-gated component ADRs (MultiSelect, DateRangePicker, TreeView…)   | Components |   2    |   2    | ⚑ #217  |
| 82   | Fuzz `wire.Action.Attributes` + `wire.DecodeForm` with adversarial input          | Wire       |   2    |   2    | ⚑ #48   |
| 83   | `go test -race` lane for visualtest helpers                                       | Testing    |   2    |   2    |         |
| 84   | Benchmark SVG chart geometry on large series (1k+ points), document limits        | Perf       |   2    |   2    |         |
| 85   | `content-visibility` guidance + demo for long tables/cards                        | Perf       |   2    |   2    |         |
| 86   | `preconnect`/`fetchpriority` audit for the self-hosted htmx path                  | Perf       |   2    |   2    |         |
| 87   | `tc explain <component>` printing docs + a runnable example                       | DX         |   2    |   2    |         |
| 88   | RSS/Atom feed for the changelog                                                   | Docs       |   2    |   2    | ⚑       |
| 89   | Richer JSON-LD (BreadcrumbList/TechArticle) + `rel=prev/next` per docs page       | Docs       |   2    |   2    | ⚑       |
| 90   | CSP violation telemetry (`report-to`) so violations stop being invisible          | Docs       |   2    |   2    | ⚑       |
| 91   | Workflow hardening pass: `permissions: contents: read` + `workflow_dispatch` cron | CI         |   2    |   2    | ⚑       |
| 92   | Publish the "SSE shipped 100% inert for months" audit writeup                     | Ecosystem  |   2    |   2    |         |
| 93   | Public roadmap board / GitHub Projects synced from `ROADMAP.md`                   | Ecosystem  |   2    |   2    |         |

### Effort 3

| Rank | Idea                                                                             | Theme   | Impact | Effort | Tracked |
| ---- | -------------------------------------------------------------------------------- | ------- | :----: | :----: | ------- |
| 94   | Exported `datastar.SSEEventWriter`-style helper for the SSE wire format          | Wire    |   2    |   3    | ⚑       |
| 95   | Screen-reader (or documented proxy) audit of every `aria-live` region end to end | A11y    |   2    |   3    |         |
| 96   | e2e cover the kanban move endpoint's 422 sorted-view rejection                   | Testing |   2    |   3    | ⚑ #224  |
| 97   | Repeatable Lighthouse lane (`.#lighthouse` or chromedp timing budget)            | DX      |   2    |   3    | ⚑ #272  |

### Effort 4

| Rank | Idea                                                                   | Theme   | Impact | Effort | Tracked     |
| ---- | ---------------------------------------------------------------------- | ------- | :----: | :----: | ----------- |
| 98   | Mutation-testing pilot (gremlins) on `utils` with a kill-rate baseline | Testing |   2    |   4    | ⚑ #215      |
| 99   | PR wall-clock budget (+20% comment) and a benchstat PR comment         | CI      |   2    |   4    | ⚑ #213/#214 |

### Effort 5

| Rank | Idea                                                                     | Theme     | Impact | Effort | Tracked |
| ---- | ------------------------------------------------------------------------ | --------- | :----: | :----: | ------- |
| 100  | Cross-language port of `wire` attribute rendering (templ-free templates) | Ecosystem |   1    |   5    |         |

---

## Theme index

| Theme          | Ranks                                         |
| -------------- | --------------------------------------------- |
| API            | 12, 13, 23, 29, 30, 42, 43, 63, 79, 80        |
| Components     | 6, 18, 31, 32, 33, 34, 44, 45, 64, 81         |
| Wire           | 7, 21, 35, 46, 47, 68, 70, 71, 82, 94         |
| A11y           | 14, 24, 36, 48, 49, 50, 51, 52, 65, 95        |
| Testing        | 5, 15, 37, 38, 53, 72, 73, 83, 96, 98         |
| Performance    | 16, 39, 54, 55, 56, 66, 74, 84, 85, 86        |
| DX / Tooling   | 8, 17, 22, 40, 41, 57, 58, 75, 87, 97         |
| Docs / Website | 4, 19, 20, 28, 59, 60, 61, 69, 76, 88, 89, 90 |
| CI / Release   | 2, 3, 9, 10, 25, 27, 77, 78, 91, 99           |
| Ecosystem      | 1, 11, 26, 62, 67, 92, 93, 100                |

---

## Suggested next actions

1. **Claim the P1 rows as TODO entries** (#1–6) with IDs — all are <2h except ConfirmDialog, and four are already half-tracked.
2. **Schedule the P2 big bets as spikes** (#7–22): `tc doctor`, `HTMXOff`, post-tag smoke, release-lock, and the daemon gates are the reliability spine; playground + example app are the adoption spine.
3. **Fold P3/P4 into `ROADMAP.md`** rather than `TODO_LIST.md` until they graduate to actionable work.

_Generated from a repo-wide read (TODO_LIST, ROADMAP, STANDOUT-IDEAS, SUPERB-FOR-PERSONAL-USE, AGENTS, and the module tree) on 2026-10-01._
