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

# Status Report — Session 11: Comprehensive Hardening + JS Deep Research

> **Updated:** 2026-07-06 (post-v0.8.0 final review). **Current version:** 0.8.0

> **Date:** 2026-07-06 01:36
> **Session scope:** Motion-reduce a11y sweep, combobox focusout handler, motion constant wiring, pagination `rel="canonical"`, JS deep research guide, container query recipe, motion design reference, icon RTL audit, semantic token ADR
> **Build:** ✅ Passing · **Tests:** 13/13 packages ✅ · **Lint:** 0 issues ✅
> **Version:** v0.8.0 (released — tag `v0.8.0`, commit `2d2d127`)
> **Commits this session:** 8 (`a0dbae7` → `2d2d127`)
> **Files changed:** 50 files, +1,279 insertions, -220 deletions

> **UPDATE NOTE (2026-07-06):**
> **FACTUAL CORRECTION:** Item #8 below claims "`transitionNormal` replaces inline string in
> Accordion panel" — this is **incorrect**. The BuildFlow-generated `utils.Class(transitionNormal, templ.KV(...))`
> call was reverted because `templ.KV` can't be passed to `utils.Class`. The accordion now uses
> **inline** `motion-reduce:transition-none motion-reduce:duration-0` classes, not the shared constant.
> Only **3** components use shared motion constants (Modal, Drawer, CopyButton), not 4.
>
> The "Top #1 Question" about motion constant adoption timing remains open — 3 of 22
> transition-bearing components use shared constants.

---

## a) FULLY DONE ✅

### Bug Fixes & Accessibility

| #   | Item                            | Details                                                                                                                                                                                                                                                                                     | Commit    |
| --- | ------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------- |
| 1   | **Motion-reduce a11y sweep**    | 7 `transition-colors` instances missing `motion-reduce` fallbacks fixed across: toast dismiss button (JS string + HTML class ×2), step indicator circle, empty state action class, file input, error page action buttons (×2). **Zero motion-reduce gaps remaining in the entire library.** | `de8171c` |
| 2   | **Combobox `focusout` handler** | Listbox now closes and `aria-activedescendant` clears when focus leaves the combobox container (mouse click outside, Tab away). Previously could remain stale if outside-click handler didn't fire before blur.                                                                             | `de8171c` |
| 3   | **AGENTS.md split-brain fix**   | Corrected false claim that repo uses "multi-module workspace with 6 modules". Master branch is a **single module** — the multi-module `go.work` structure was an abandoned experiment.                                                                                                      | `a0dbae7` |

### Code Quality

| #   | Item                                        | Details                                                                                                                                                                                                                                                                                  | Commit    |
| --- | ------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------- |
| 4   | **StatCardProps.HxSwap typed**              | Changed from raw `string` to `htmx.SwapStyle` — consumers now pass typed constants instead of string literals. Cross-package type safety.                                                                                                                                                | `cc88d41` |
| 5   | **SortDirectionIsValid**                    | Added to complete the enum validation set (was the only enum missing IsValid).                                                                                                                                                                                                           | `cc88d41` |
| 6   | **ButtonHTMLType typed map**                | Converted from `map[X]bool` to `map[X]string` + `utils.Lookup`, matching all other enums.                                                                                                                                                                                                | `cc88d41` |
| 7   | **Feedback/errorpage lookup consolidation** | Replaced `lookupFeedbackStyle()` generic + `feedbackIconName()` helper with direct `utils.Lookup()` calls. Same behavior, less custom code. Same for `FamilyStatusCode`.                                                                                                                 | `d3c8b88` |
| 8   | **Motion constants wired into CopyButton**  | `transitionColors` replaces inline string in CopyButton. Accordion uses inline `motion-reduce:*` classes (the `transitionNormal` constant approach was reverted — `templ.KV` can't be mixed with `utils.Class`). **3 components** total use shared constants: Modal, Drawer, CopyButton. | `de8171c` |

### SEO

| #   | Item                                        | Details                                                                                                                                                                                                          | Commit    |
| --- | ------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------- |
| 9   | **Pagination `rel="canonical"` for page 1** | Page 1 link now carries `rel="canonical"` when rendered as a non-active `<a>`. New `rel` parameter added to `activeSpanOrLink` and `paginationPageItem` signatures. Breadcrumbs caller updated with empty `rel`. | `098f7c3` |

### Testing

| #   | Item                                     | Details                                                                                                                                                  | Commit    |
| --- | ---------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- | --------- |
| 10  | **FEATURES.md version-sync drift-guard** | `TestVersionMatchesFeatures` asserts `FEATURES.md` `**Version:**` marker matches `utils.Version`. Prevents documentation drift like the 0.6.1→0.7.0 gap. | `6e94f93` |
| 11  | **TableHeader sortable golden test**     | Golden snapshot for sortable columns with `aria-sort` and ↑/↓ indicators.                                                                                | `cc88d41` |
| 12  | **StatCard `<a>` motion-reduce**         | Added `motion-reduce:transition-none motion-reduce:duration-0` to linked StatCard variant.                                                               | `cc88d41` |

### Documentation

| #   | Item                                   | Details                                                                                                                                                                                                                                                                                                                                                          | Commit    |
| --- | -------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------- |
| 13  | **JavaScript in Templ Projects guide** | 472-line comprehensive reference at `docs/javascript-guide.md`. Covers decision ladder (native HTML → HTMX → singleton-guard → Alpine → Datastar → React islands), CSP compliance, templ built-in JS features (`OnceHandle`, `JSFuncCall`, `JSONString`, `JSONScript`), TypeScript workflow, View Transitions API, event delegation, anti-patterns. 12 sections. | `cd20462` |
| 14  | **Motion design reference**            | `docs/motion-design.md`: timing constants table, duration guidelines (100-400ms), easing policy (`ease-out` default), motion-reduce compliance rules, per-component adoption table.                                                                                                                                                                              | `cd20462` |
| 15  | **Container queries recipe**           | `docs/recipes/container-queries.md`: when/how to use `ContainerResponsive`, Tailwind v4 `@container` size reference table, manual container query example.                                                                                                                                                                                                       | `cd20462` |
| 16  | **Icon RTL mirroring audit**           | `docs/audits/icon-rtl-mirroring.md`: identifies 5 directional icons (ArrowRight, ArrowLeft, ChevronRight, PathArrowLeft, PathArrowRight) needing RTL mirroring. Recommends `data-tc-dir-icon` + CSS `scaleX(-1)` approach. Deferred to v1.0.                                                                                                                     | `cd20462` |
| 17  | **Semantic token layer ADR**           | `docs/adr/0008-semantic-tokens.md`: proposes `bg-tc-primary` semantic aliases over hardcoded `bg-blue-600`. Three-phase migration (document → opt-in → default). Deferred to v1.0.                                                                                                                                                                               | `cd20462` |
| 18  | **SKILL.md conventions update**        | Added RTL logical properties, motion constants, and container query conventions to mandatory conventions + anti-patterns sections. Added `docs/javascript-guide.md` to progressive disclosure list.                                                                                                                                                              | `cd20462` |
| 19  | **CHANGELOG entries**                  | All code changes and new docs added to `[Unreleased]` section.                                                                                                                                                                                                                                                                                                   | `cd20462` |
| 20  | **v0.8.0 release**                     | Version bumped, CHANGELOG heading cut, tag `v0.8.0` created.                                                                                                                                                                                                                                                                                                     | `2d2d127` |

---

## b) PARTIALLY DONE 🟡

| #   | Item                         | What's done                                                                                                                                                                                                                | What's missing (Status 2026-07-06)                                                                                                                             |
| --- | ---------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **Motion constant adoption** | **✅ COMPLETED (2026-07-10):** All 22 transition-bearing components now use shared constants via `utils/motion.go`. 19 previously-inline components wired.                                                                 | —                                                                                                                                                              |
| 2   | **RTL support**              | All CSS properties migrated to logical (`ms-`, `me-`, `ps-`, `pe-`, `start-`, `end-`). 0 physical properties remain. **✅ Icon RTL implemented (2026-07-10):** `icons.IconRTL()` + CSS `scaleX(-1)` for directional icons. | Keyboard navigation in Dropdown/Tabs maps ArrowLeft/Right without checking `dir` attribute (LOW priority). No `dir="rtl"` golden rendering tests.              |
| 3   | **Container query adoption** | Grid supports `ContainerResponsive` with full test coverage + recipe doc.                                                                                                                                                  | Only Grid has container query support. Other components (Card, SidebarNav) could benefit.                                                                      |
| 4   | **JS documentation**         | 472-line guide covers all patterns                                                                                                                                                                                         | No ADR for "why singleton-guard instead of Alpine.js" — the ADR 0005 exists but predates the comprehensive guide. The guide could be cross-linked from README. |

---

## c) NOT STARTED ⬜

| #   | Item                                                              | Priority      | Effort      | Why not started                                                                                          |
| --- | ----------------------------------------------------------------- | ------------- | ----------- | -------------------------------------------------------------------------------------------------------- |
| 1   | **Wire motion constants into remaining ~16 components**           | MEDIUM        | 90min       | Mechanical work but touches every component. Risk of golden file churn. Best done as a dedicated sprint. |
| 2   | **RTL rendering tests** (`dir="rtl"` golden tests)                | MEDIUM        | 30min       | Was planned but deprioritized — logical properties are in place, visual verification is the gap.         |
| 3   | **Icon RTL mirroring** (`data-tc-dir-icon` attribute)             | LOW           | 45min       | Requires changing `icons.Icon` signature or adding wrapper. Minor breaking change deferred to v1.0.      |
| 4   | **Semantic token layer** (`bg-tc-primary` → `bg-blue-600`)        | HIGH (v1.0)   | 4hrs+       | 256 color references. Massive golden file churn. Needs dedicated major-version migration. ADR written.   |
| 5   | **Native `<dialog>` element** for Modal/Drawer                    | MEDIUM        | 3hrs        | Fundamental architecture change. Needs ADR. Browser focus trap is better than JS but migration is risky. |
| 6   | **New components** (Popover, Slider, Rating, DataTable, Carousel) | MEDIUM        | 2-6hrs each | Additive work, separate effort.                                                                          |
| 7   | **CSS `@starting-style`** for zero-JS enter/exit                  | LOW           | 2hrs        | Requires modern browser baseline decision.                                                               |
| 8   | **`Validate() error` on props structs**                           | MEDIUM (v1.0) | 4hrs        | Design decision needed: replace fallback pattern or supplement as opt-in?                                |
| 9   | **Move test helpers to `internal/testutil/`**                     | LOW (v1.0)    | 2hrs        | Breaking change for external consumers. 70+ test files depend on exported helpers.                       |
| 10  | **Compound component pattern** (Trigger/Content/Close)            | LOW (v2.0)    | 3hrs        | Breaking API change. Current monolithic Modal/Drawer work fine.                                          |

---

## d) TOTALLY FUCKED UP 💥

| #   | Issue                                                | What happened                                                                                                                                                                                                                                                                                                                                                                          | Resolution                                                                                                                               |
| --- | ---------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **Previous session's code edits were silently lost** | The first hardening attempt edited 7 `.templ` files (motion-reduce fixes, combobox focusout, accordion/copy_button wiring). Before those edits could be committed, another session (or BuildFlow hook) committed different changes to the same files, overwriting the edits. The doc files (SKILL.md, CHANGELOG, new docs) survived because they weren't touched by the other session. | ✅ Fixed: re-did all code edits against the current tree. BuildFlow pre-commit hook then auto-committed them cleanly.                    |
| 2   | **AGENTS.md claimed multi-module workspace**         | The repo root `AGENTS.md` described a "6-module go.work structure" that doesn't exist on master. This was an abandoned experiment. The false documentation misled both AI sessions and human contributors about how to build/test the project.                                                                                                                                         | ✅ Fixed: `a0dbae7` corrected to single-module description.                                                                              |
| 3   | **BuildFlow pre-commit hook auto-commits**           | The BuildFlow `pre-commit` hook detects file changes and auto-commits them with its own message before the explicit `git commit` runs. This means the commit history shows multiple small commits instead of one clean commit per logical change. It also means if you edit files and then review `git diff`, you may see nothing — the hook already committed.                        | ⚠️ Known infrastructure behavior. Not harmful, but disorienting. Workaround: check `git log` instead of `git diff` to see what's staged. |

---

## e) WHAT WE SHOULD IMPROVE 🔧

1. **Wire motion constants into ALL transition-bearing components** — 3 of 22 components use shared constants. Partial adoption creates inconsistency. Either commit fully or remove the constants (current state is the worst of both worlds). Note: Accordion can't use the constant via `utils.Class` due to `templ.KV` incompatibility — would need a different pattern or inline approach.

2. **Add `dir="rtl"` rendering tests** — logical properties are in place but zero tests verify RTL rendering output. A golden test on 2-3 key components (Card, Drawer, Nav) would prove the migration.

3. **Fix BuildFlow auto-commit behavior or document it prominently** — The hook silently commits edited `.templ` files before the explicit commit. This caused the "lost edits" incident in this session. Either disable the auto-commit or add a prominent warning in AGENTS.md.

4. **Cross-link `docs/javascript-guide.md` from README** — The 472-line guide is the most comprehensive JS reference for templ projects but is only discoverable via SKILL.md's progressive disclosure list. A README link would help consumers.

5. **Container query documentation in README** — The `ContainerResponsive` prop exists but the README doesn't mention it. The recipe doc helps but discoverability is still low.

6. **Motion constant audit automation** — A grep-based test that asserts every `transition-*` class in `.templ` files uses a shared constant (or has an explicit exemption) would prevent regression.

7. **Pagination keyboard RTL mapping** — Dropdown and Tabs JS handlers map `ArrowRight` → next and `ArrowLeft` → previous without checking `dir` attribute. In RTL, this mapping should swap per WAI-ARIA APG.

---

## f) Top 25 Things to Get Done Next

Sorted by impact × effort × customer value.

| #   | Task                                                                   | Impact | Effort | Est.  |
| --- | ---------------------------------------------------------------------- | ------ | ------ | ----- |
| 1   | Wire motion constants into remaining ~16 transition-bearing components | MEDIUM | MEDIUM | 90m   |
| 2   | Add RTL rendering tests (Card, Drawer, Nav with `dir="rtl"`)           | MEDIUM | LOW    | 30m   |
| 3   | Add motion-constant-compliance test (grep assertion)                   | MEDIUM | LOW    | 15m   |
| 4   | Cross-link `docs/javascript-guide.md` from README                      | LOW    | LOW    | 5m    |
| 5   | Demo app: showcase TableHeader sortable columns                        | MEDIUM | LOW    | 15m   |
| 6   | Demo app: showcase Form.Inline horizontal layout                       | LOW    | LOW    | 10m   |
| 7   | Demo app: showcase Grid.ContainerResponsive                            | LOW    | LOW    | 10m   |
| 8   | Pagination keyboard RTL mapping (ArrowLeft/Right swap in `dir="rtl"`)  | MEDIUM | MEDIUM | 30m   |
| 9   | Icon RTL mirroring (`data-tc-dir-icon` + CSS `scaleX(-1)`)             | LOW    | MEDIUM | 45m   |
| 10  | Plan semantic token layer migration (opt-in `tc-*` tokens)             | HIGH   | HIGH   | 4hrs+ |
| 11  | Add Popover component (compound pattern, most requested)               | HIGH   | HIGH   | 4hrs  |
| 12  | Migrate Modal to native `<dialog>` element (with ADR)                  | MEDIUM | HIGH   | 3hrs  |
| 13  | Add HoverCard component                                                | MEDIUM | MEDIUM | 2hrs  |
| 14  | Add Slider component (ARIA slider pattern)                             | MEDIUM | MEDIUM | 2hrs  |
| 15  | Add DataTable component (sorting, filtering, pagination)               | HIGH   | HIGH   | 6hrs+ |
| 16  | Add CSS `@starting-style` support to Modal/Drawer                      | MEDIUM | MEDIUM | 2hrs  |
| 17  | Add Calendar component (full calendar grid)                            | MEDIUM | HIGH   | 4hrs  |
| 18  | Add Rating component (star rating with keyboard support)               | LOW    | LOW    | 1hr   |
| 19  | Add TagsInput component                                                | LOW    | MEDIUM | 2hrs  |
| 20  | Add ContextMenu component (right-click menu)                           | LOW    | MEDIUM | 2hrs  |
| 21  | Create formal "blocks" (dashboard, login, settings page examples)      | MEDIUM | MEDIUM | 3hrs  |
| 22  | Add compound component pattern for future overlays                     | MEDIUM | HIGH   | 3hrs  |
| 23  | Add `Validate() error` to props structs (v1.0 scope)                   | MEDIUM | HIGH   | 4hrs  |
| 24  | Move test helpers to `internal/testutil/` (v1.0 scope)                 | LOW    | HIGH   | 2hrs  |
| 25  | Performance benchmark suite (render time per component)                | LOW    | MEDIUM | 2hrs  |

---

## g) Top #1 Question I Cannot Figure Out Myself

> ⚠️ **STILL OPEN.** The motion constant adoption question remains unresolved.
> Current state: 3 of 22 transition-bearing components use shared constants (Modal, Drawer,
> CopyButton). Accordion can't use `transitionNormal` via `utils.Class` due to `templ.KV`
> incompatibility — needs a different approach.
>
> **Recommendation:** Defer to a broader "design system" sprint that also addresses semantic
> tokens (256 color refs), easing curves, and view transitions. Doing motion constants now
> means two rounds of golden file churn instead of one. The inconsistency is cosmetic — the
> inline strings are semantically identical to the constants.

---

## Session Metrics

| Metric                                   | Value                                                                                              |
| ---------------------------------------- | -------------------------------------------------------------------------------------------------- |
| Commits                                  | 8                                                                                                  |
| Files changed                            | 50                                                                                                 |
| Lines added                              | +1,279                                                                                             |
| Lines removed                            | -220                                                                                               |
| Motion-reduce gaps fixed                 | 7 → 0 remaining                                                                                    |
| Components using shared motion constants | 3 (was 2) — Modal, Drawer, CopyButton. Accordion uses inline classes (constant approach reverted). |
| New documentation files                  | 6 (JS guide, motion ref, container query recipe, RTL audit, semantic token ADR, hardening plan)    |
| New tests                                | 3 (FEATURES.md drift-guard, table golden, StatCard golden update)                                  |
| Typed enums with IsValid                 | 33 (complete)                                                                                      |
| Packages below 70% coverage              | 0                                                                                                  |
| Test packages green                      | 13/13                                                                                              |
| Lint issues                              | 0                                                                                                  |
| Version released                         | v0.8.0                                                                                             |
| JS guide lines                           | 472                                                                                                |
| Total library size                       | ~24,000 lines Go/templ                                                                             |
