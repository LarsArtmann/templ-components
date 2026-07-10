<!-- AUTO-UPDATED 2026-07-10: Retrospective status overlay -->

> ## 🔔 Update Notice — 2026-07-10
>
> This report is **historical**. Many items listed as "open", "todo", or "broken" below
> have since been **fixed and verified**. Do not act on open items without first checking
> [TODO_LIST.md](../../TODO_LIST.md) for current status.
>
> **Key fixes completed since this report:**
>
> - ✅ All 7 P0 bugs fixed (InlineLoadingOverlay a11y, SanitizeID mismatch, FromError fallback,
>   Footer BaseProps, ErrorPage/NotFound404 `<main>` landmark, CSRFTokenName, grid-rows verified)
> - ✅ `encoding/json/v2` purged from all production code + pre-commit guard added
> - ✅ Motion constants centralized in `utils/motion.go`, wired into 13 components
> - ✅ `FamilyFromErrorFamily` → `FromErrorFamily` (old name kept as deprecated alias)
> - ✅ `icons.IconRTL()` + CSS for directional icon RTL mirroring
> - ✅ 33 regression tests added (htmx, errorpage, layout, navigation, feedback, display)
> - ✅ Dark golden test infrastructure (badge/card/button)
> - ✅ CHANGELOG consolidated, ROADMAP updated, migration guide created
> - ✅ All 14 packages pass, 0 lint issues
>
> **Canonical source of truth:** [TODO_LIST.md](../../TODO_LIST.md) (52 items, 37 ✅ done, 12 deferred/blocked)

---

# Status Report — 2026-07-05 03:14 CEST

> **Updated:** 2026-07-06 (post-v0.8.0). Version at report: 0.6.x → **Current:** 0.8.0

**Session goal:** Process `docs/feedback/*` consumer feedback, build a comprehensive
prioritized TODO list, execute every item, verify.

**Done-check:** `nix run .#verify` → **All checks passed** (generate + build + test + lint, 0 issues)

> **UPDATE NOTE (2026-07-06):** This was session 6. Since then, sessions 7–10 + v0.7.0/v0.8.0
> releases landed. Most open items from this report are now resolved. Items closed since are
> marked ✅. Current state: v0.8.0, 575 test cases, 13/13 packages green, 0 lint issues.

---

## a) FULLY DONE (verified green)

| #   | What                                                                                                                                    | Evidence                                       |
| --- | --------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------- |
| 1   | `PageProps` godoc: documented `CSSPath` and `HTMXVersion` silent auto-injects + how to suppress                                         | `layout/base.templ:50-68, 78-88`               |
| 2   | README "Suppressing auto-injected `<head>` tags" subsection with copy-paste example                                                     | `README.md:96-113`                             |
| 3   | `SimpleNavProps.RightItems` field + forwarded to `Nav.RightItems`                                                                       | `navigation/nav.templ:70-105`                  |
| 4   | `StatCardProps.Href` — wraps card in `<a>` with hover/focus/cursor styling; extracted `statCardInner` sub-template                      | `display/card.templ:131-230`                   |
| 5   | `layout.Script(nonce, src, attrs)` CSP-safe script helper                                                                               | `layout/script.templ`                          |
| 6   | `display.Grid` + typed `GridCols` enum (1–6) with map+fallback                                                                          | `display/grid.templ`                           |
| 7   | `feedback.SkeletonCardGrid(count)` responsive loading grid                                                                              | `feedback/loading.templ:207-225`               |
| 8   | Recipe: `docs/migration/play-cdn-to-tailwind-v4.md` (7-step migration)                                                                  | New file                                       |
| 9   | Recipe: `docs/recipes/server-rendered-htmx-error-feedback.md` (3 render modes)                                                          | New file                                       |
| 10  | README discoverability: Grid, SkeletonCardGrid, StatCard.Href, SimpleNav.RightItems examples; recipe cross-links; component count 73→76 | `README.md`                                    |
| 11  | CHANGELOG `[Unreleased]` — comprehensive Added/Changed/Internal                                                                         | `CHANGELOG.md:7-44`                            |
| 12  | `GridProps` registered in contract inventory test                                                                                       | `internal/contract/component_props_test.go:47` |
| 13  | Fixed stale `sidebar_nav.golden` (pre-existing failure — templ runtime cosmetic space change)                                           | `navigation/testdata/sidebar_nav.golden`       |
| 14  | Fixed 4 lint errors in `sri_net_test.go` (pre-existing — errcheck/noctx/paralleltest)                                                   | `layout/sri_net_test.go`                       |
| 15  | 9 new test functions across 4 packages — all pass under `-race`                                                                         | Various `*_test.go`                            |

**Verify output:** 13/13 packages `ok`, 0 lint issues.

---

## b) PARTIALLY DONE

| Item                             | What's done                                                                                      | What's missing                                                                                                                                 |
| -------------------------------- | ------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| Test coverage for new components | Assertion-based unit tests (Grid, Script, SkeletonCardGrid, StatCard.Href, SimpleNav.RightItems) | No golden files, no BDD tests, no a11y tests, no example tests for the 3 new components (skill mandates all test lenses per component)         |
| StatCard refactor                | `statCardInner` extracted, Href branch added, existing tests pass                                | `display/testdata/stat_card_icon.golden` was not checked against the restructured output (tests passed, but no explicit golden regen was done) |
| README accuracy                  | Component count bumped to 76, enum count to 27                                                   | Component-per-package headers not updated (display says "20 components" ✓, but feedback still says "12" — SkeletonCardGrid makes 13)           |

---

## c) NOT STARTED (should have been)

| #   | What                                  | Why it matters                                                                                                                                                                            | Status (2026-07-06)                                                            |
| --- | ------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------ |
| 1   | **AGENTS.md update**                  | 3 new components, 2 new fields, 1 new enum, `statCardInner` sub-template — NONE documented in AGENTS.md conventions. Violates the "Aggressive Update Protocol" from the global AGENTS.md. | ✅ Done — AGENTS.md now documents all conventions comprehensively              |
| 2   | **TODO_LIST.md update**               | None of the 13 completed items were recorded. TODO_LIST still shows the old session 5 header.                                                                                             | ✅ Done — maintained through v0.8.0                                            |
| 3   | **FEATURES.md update**                | New components (Grid, Script, SkeletonCardGrid) and new fields not in the feature inventory.                                                                                              | ✅ Done — version-sync drift-guard test added                                  |
| 4   | **`examples/demo/` update**           | Demo is the canonical "how a consumer assembles a page" reference. Grid not demonstrated there.                                                                                           | ✅ Done — demo updated with Grid + StatCard composition                        |
| 5   | **`integration/composition_test.go`** | Grid composes with Card/StatCard — cross-package integration test not extended.                                                                                                           | ✅ Done — 7 composition integration tests added                                |
| 6   | **Skill `SKILL.md` update**           | New GridCols enum, Grid component, Script helper, SkeletonCardGrid not mentioned in decision trees or conventions.                                                                        | ✅ Done — SKILL.md fully rewritten with use-case table + authoring conventions |

---

## d) TOTALLY FUCKED UP

Nothing is broken — verify passes, tests pass, lint passes. But these are
**judgment failures**, not bugs:

| #   | What                                                                                                                                                                                                                                                                                                                                       | Severity           |
| --- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------ |
| 1   | **Ignored the `templ minmax` diagnostic for the entire session.** `feedback/loading.templ:217` (`if n <= 0`) can be modernized to `max`. It was visible in every tool output. 30-second fix, never done.                                                                                                                                   | Low (cosmetic)     |
| 2   | **`GridCols5` has a poor responsive ladder.** Jumps from `sm:grid-cols-2` straight to `lg:grid-cols-5` — skipping 3 and 4 column intermediate states. Should be `sm:grid-cols-3 lg:grid-cols-5` for a smoother progression.                                                                                                                | Medium (design)    |
| 3   | **`StatCard.Href` has no typed HTMX fields.** Overview's actual use case was `hx-get` + `hx-target`. My implementation renders a plain `<a>` — the consumer must pass hx attributes via `Attrs` (string map). Button has typed `HxGet`/`HxTarget` fields. StatCard.Href is ergonomically inferior for the exact consumer who requested it. | Medium (ergonomic) |
| 4   | **Claimed "76 components" without rigorous counting.** Is `layout.Script` a "component" or a "helper"? Is `SkeletonCardGrid` a new component or a Skeleton variant? The count may be off by 1–2. Previous counts were carefully verified (see v0.5.0 changelog "corrected metrics").                                                       | Low (accuracy)     |

---

## e) WHAT WE SHOULD IMPROVE

> All 7 items below were resolved in sessions 7–10 + v0.7.0/v0.8.0 releases.

1. ✅ **Update AGENTS.md immediately** — Done. All conventions documented and maintained through v0.8.0.
2. ✅ **Add golden tests for new components** — Done. All new components have golden files.
3. ✅ **Add BDD + a11y tests** — Done. Grid, Script, SkeletonCardGrid all have BDD + a11y + example tests.
4. ✅ **Fix the GridCols5 responsive ladder.** — Done (added intermediate md breakpoint).
5. ✅ **Add typed HTMX fields to StatCard** — Done (`HxGet`/`HxTarget`/`HxSwap` added; `HxSwap` typed as `htmx.SwapStyle` in v0.8.0).
6. ✅ **Update `examples/demo/`** — Done. Grid + StatCard composition showcased.
7. ✅ **Fix the `templ minmax` hint** — Done (ProgressBar uses `max(0, min(100, v))`).

---

## f) Up to 25 things to do next

| #   | Task                                                                                     | Impact | Effort | Status (2026-07-06)                                            |
| --- | ---------------------------------------------------------------------------------------- | ------ | ------ | -------------------------------------------------------------- |
| 1   | Update AGENTS.md with new components, fields, enums, sub-templates                       | High   | Low    | ✅ Done                                                        |
| 2   | Update TODO_LIST.md with completed session 6 items                                       | Med    | Low    | ✅ Done                                                        |
| 3   | Update FEATURES.md with Grid, Script, SkeletonCardGrid                                   | Med    | Low    | ✅ Done                                                        |
| 4   | Add golden test for `display.Grid` (all 6 GridCols variants)                             | Med    | Low    | ✅ Done                                                        |
| 5   | Add golden test for `layout.Script`                                                      | Low    | Low    | ✅ Done                                                        |
| 6   | Add golden test for `feedback.SkeletonCardGrid`                                          | Low    | Low    | ✅ Done                                                        |
| 7   | Add BDD test for Grid (user sees responsive grid)                                        | Med    | Low    | ✅ Done                                                        |
| 8   | Add BDD test for StatCard.Href (user clicks stat card → navigates)                       | Med    | Low    | ✅ Done                                                        |
| 9   | Add a11y test for Grid (aria-label propagation, role)                                    | Med    | Low    | ✅ Done                                                        |
| 10  | Add a11y test for Script (nonce always emitted)                                          | Med    | Low    | ✅ Done (CSP nonce-presence test covers all script components) |
| 11  | Add a11y test for SkeletonCardGrid (role=status, motion-reduce)                          | Med    | Low    | ✅ Done                                                        |
| 12  | Add godoc ExampleGrid, ExampleScript, ExampleSkeletonCardGrid                            | Low    | Low    | ✅ Done                                                        |
| 13  | Fix `GridCols5` responsive ladder (`sm:grid-cols-3 lg:grid-cols-5`)                      | Med    | Low    | ✅ Done                                                        |
| 14  | Add typed HTMX fields to StatCard (HxGet, HxTarget, HxSwap) or document Attrs workaround | High   | Med    | ✅ Done (HxSwap typed as htmx.SwapStyle in v0.8.0)             |
| 15  | Fix `templ minmax` diagnostic in `feedback/loading.templ:217`                            | Low    | Low    | ✅ Done                                                        |
| 16  | Add Grid + StatCard composition to `examples/demo/`                                      | Med    | Low    | ✅ Done                                                        |
| 17  | Add Grid to `integration/composition_test.go`                                            | Med    | Low    | ✅ Done                                                        |
| 18  | Update `feedback` README section count (12 → 13 components)                              | Low    | Low    | ✅ Done                                                        |
| 19  | Update skill `SKILL.md` with GridCols enum in decision tree + Script helper              | Med    | Low    | ✅ Done (full rewrite)                                         |
| 20  | Process the 2 new untracked feedback files (DiscordSync, swettyswipper)                  | High   | Med    | ✅ Done (session 7 implemented all feedback)                   |
| 21  | Consider `GridProps.Gap` typed enum (gap-2/gap-4/gap-6/gap-8)                            | Low    | Low    | ⬜ Not done — low priority                                     |
| 22  | Consider `Card.Header` / `Card.Body` slot fields (SEC feedback)                          | Med    | Med    | ✅ Done (Card.Body slot shipped)                               |
| 23  | Verify README component count (76) by actual grep across packages                        | Low    | Low    | ✅ Done (82 components verified)                               |
| 24  | Consider `layout.Stylesheet(nonce, href, attrs)` companion to `layout.Script`            | Low    | Low    | ✅ Done                                                        |
| 25  | Cut v0.7.0 release once items 1–14 are done                                              | High   | Low    | ✅ Done (v0.7.0 + v0.8.0 released)                             |

**Scorecard:** 23 of 25 complete (92%).

---

## g) Top #1 question I cannot figure out myself

> ✅ **RESOLVED.** StatCard got typed HTMX fields (`HxGet`/`HxTarget`/`HxSwap`) in session 7.
> In v0.8.0, `HxSwap` was further typed from `string` to `htmx.SwapStyle` for compile-time safety.
> The design philosophy question is settled: typed HTMX everywhere for components that support it.
