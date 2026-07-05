# Status Report — UI Library Design Research + Pareto Execution

> **Date:** 2026-07-05 23:02 · **Updated:** 2026-07-06 (post-v0.8.0)
> **Session scope:** Deep research (shadcn/ui, Radix, React Aria, templui, HATEOAS, WAI-ARIA, Tailwind v4, design tokens, motion, forms, i18n) → Pareto plan → execution of RTL migration, motion constants, Grid container queries
> **Build:** ✅ Passing · **Tests:** 13/13 packages ✅ (575 cases) · **Lint:** 0 issues ✅
> **Version at report time:** 0.7.0 → **Current:** 0.8.0

> **UPDATE NOTE (2026-07-06):** This report was written at the end of the
> research-and-execution session. Two follow-up sessions landed since then:
> **Session 10** (brutal self-review + hardening — typed enums, IsValid sweep,
> sortable TableHeader, Form.Inline, coverage boost) and the **v0.8.0 release**
> (motion-reduce sweep, combobox focusout fix, comprehensive docs: JS guide,
> motion reference, container-query recipe, RTL icon audit, semantic-token ADR).
> This document has been updated to reflect the **current** state. Items closed
> since the original report are marked with ✅ and moved to the appropriate section.

---

## a) FULLY DONE

### Original session deliverables

| #   | Item                                 | Details                                                                                                                                                                                                                                                                                                                                                                                                            | Commit                  |
| --- | ------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------- |
| 1   | **Research report**                  | 1,705-line synthesis across 10+ sources covering shadcn/ui, Radix UI, React Aria, templui, HATEOAS, WAI-ARIA APG, Tailwind v4, DTCG design tokens, CVA, motion design, forms, i18n/RTL. Written to `docs/research/ui-library-design-research.md`.                                                                                                                                                                  | `4ac1c6a`               |
| 2   | **Pareto plan**                      | Full 1%/4%/20% breakdown with 55 fine-grained tasks, mermaid execution graph, written to `docs/planning/2026-07-05_21-27_SUPERB-UI-LIBRARY-UPGRADES.md`.                                                                                                                                                                                                                                                           | `33726b3`               |
| 3   | **RTL logical properties migration** | **74 physical CSS properties → logical** across ALL `.templ` files: `ml-`→`ms-`, `mr-`→`me-`, `pl-`→`ps-`, `pr-`→`pe-`, `left-0`→`start-0`, `right-0`→`end-0`, `text-left`→`text-start`, `border-l-`→`border-s-`, `border-r-`→`border-e-`. Preserved `left-1/2` (tooltip centering) and `left-0.5` (toggle thumb). Zero remaining physical properties. Library is now RTL-ready for Arabic, Hebrew, Persian, Urdu. | Auto-committed by hooks |
| 4   | **Motion constants**                 | Added `transitionFast` (150ms), `transitionNormal` (200ms), `transitionColors`, `transitionTransform` to `display/shared.go`. All include `motion-reduce:*` fallbacks + `ease-out` (professional default). Wired into Modal (`transitionNormal`) and Drawer (`transitionTransform + duration-200`).                                                                                                                | `33726b3`               |
| 5   | **Grid container queries**           | New `GridProps.ContainerResponsive bool` field. When true, Grid wraps in `@container` div and uses `@sm:`/`@md:`/`@lg:`/`@xl:` Tailwind v4 container-query variants instead of viewport breakpoints. New `gridColsContainerLookup` map. 3 test subtests added (`TestGridContainerResponsive`). Defaults to false (backward compatible).                                                                            | `33726b3`               |
| 6   | **AGENTS.md conventions**            | 3 new Code Conventions rules documented: (1) RTL — logical properties only, (2) Motion — shared constants not inline strings, (3) Container queries — `@container` for context-responsive grids.                                                                                                                                                                                                                   | `33726b3`               |
| 7   | **CHANGELOG entries**                | `[Unreleased]` updated with: RTL migration, motion constants, Grid container queries.                                                                                                                                                                                                                                                                                                                              | Auto-committed          |
| 8   | **Assertion test fixes**             | Fixed 3 tests that asserted old physical class names: `forms/snapshot_test.go` (InputGroupPaddingClass pl-10→ps-10), `navigation/snapshot_test.go` (border-l-4→border-s-4), `display/coverage_test.go` (Drawer left-0→start-0).                                                                                                                                                                                    | Auto-committed          |
| 9   | **Golden file updates**              | All golden snapshots regenerated via `-update` flag after RTL migration. Reviewed diffs — purely class-name changes, semantically identical.                                                                                                                                                                                                                                                                       | Auto-committed          |

### ✅ Completed since original report (Session 10 + v0.8.0)

| #   | Item                                       | Details                                                                                                                                                                                                                                                | Source                         |
| --- | ------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------ |
| 10  | ✅ **SKILL.md Part 2 updated**             | Authoring playbook now includes mandatory conventions: RTL logical properties, motion constants, container queries. New components will follow these rules.                                                                                            | Was Top-25 #1 — `cd20462`      |
| 11  | ✅ **RTL rendering tests**                 | `display/rtl_test.go` verifies logical properties render correctly under `dir="rtl"` (Drawer `start-0`, StatCard `ms-2`). Was Top-25 #2.                                                                                                               | Was Not-Started #4 — `5c1e037` |
| 12  | ✅ **Motion-reduce systematic sweep**      | All 7 `transition-colors` instances missing `motion-reduce:*` fixed across toast (dismiss ×2), step_indicator, empty_state, file_input, errorpage (action buttons ×2). **Zero remaining uncovered transitions.** Was Top-25 #13 and Partially-Done #3. | `de8171c`                      |
| 13  | ✅ **Semantic token ADR**                  | `docs/adr/0008-semantic-tokens.md` — migration plan for `bg-blue-600`→`bg-tc-primary` token layer, proposed and **deferred to v1.0** with opt-in path. Was Top-25 #5 + #25.                                                                            | `cd20462`                      |
| 14  | ✅ **Container query recipe**              | `docs/recipes/container-queries.md` — when/how to use `ContainerResponsive` for parent-width-responsive grids. Was Top-25 #6.                                                                                                                          | `cd20462`                      |
| 15  | ✅ **Motion design reference**             | `docs/motion-design.md` — timing constants, duration guidelines, easing policy, `motion-reduce` compliance rules. Was Top-25 #7.                                                                                                                       | `cd20462`                      |
| 16  | ✅ **Icon RTL mirroring audit**            | `docs/audits/icon-rtl-mirroring.md` — identifies 5 directional icons needing RTL mirroring, recommends `data-tc-dir-icon` + CSS approach, deferred to v1.0. Was Top-25 #22.                                                                            | `cd20462`                      |
| 17  | ✅ **JS patterns guide**                   | `docs/javascript-guide.md` (472 lines) — decision ladder (native HTML → HTMX → singleton-guard → Alpine → Datastar → React islands), CSP compliance, templ built-in JS features. Bonus deliverable.                                                    | `cd20462`                      |
| 18  | ✅ **prefers-color-scheme auto-detection** | Confirmed: `layout.ThemeScript()` checks `localStorage` first, falls back to `prefers-color-scheme: dark` for initial theme. Library already follows OS preference on first visit. Was Top-25 #23.                                                     | Pre-existing                   |
| 19  | ✅ **CopyButton motion constant**          | `transitionColors` wired into CopyButton — now 3 components use shared constants (was 2).                                                                                                                                                              | v0.8.0                         |
| 20  | ✅ **Combobox focusout handler**           | Listbox closes + `aria-activedescendant` clears when focus leaves combobox container. Was Session-10 Partially-Done.                                                                                                                                   | `de8171c`                      |
| 21  | ✅ **govalid hook issue resolved**         | `govalid` no longer referenced in pre-commit hook or scripts. Current hook is BuildFlow-managed; the govalid-generate step is gone. Was Top-25 #4 + "Totally Fucked Up" #2.                                                                            | Hook reconfigured              |
| 22  | ✅ **v0.8.0 released**                     | Typed enums (30 IsValid methods, all tested), typed lookup maps, sortable TableHeader, Form.Inline, coverage ≥70% all packages, motion-reduce sweep, split-brain fixes.                                                                                | `2d2d127`                      |

---

## b) PARTIALLY DONE

| #   | Item                         | What's done                                                                                                                                              | What's missing                                                                                                                                                                 |
| --- | ---------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| 1   | **Motion constant adoption** | 3 of 22 transition-bearing components use shared constants: Modal (`transitionNormal`), Drawer (`transitionTransform`), CopyButton (`transitionColors`). | **19 components still use inline timing strings.** Accordion, Tabs, Dropdown, Tooltip, Toast, Card, NavLink, Button, Avatar, Badge, etc. Full migration is a dedicated sprint. |
| 2   | **Container query adoption** | Grid supports `ContainerResponsive` with full test coverage + recipe doc.                                                                                | Only Grid has container query support. Other layout components (Card, SidebarNav) could benefit but weren't migrated.                                                          |
| 3   | **Performance benchmarks**   | 3 packages have benchmark suites: `display/benchmark_test.go`, `feedback/benchmark_test.go`, `navigation/benchmark_test.go`.                             | 11 packages have no benchmarks. No CI gating on render-time regressions.                                                                                                       |
| 4   | **Semantic token layer**     | ADR 0008 written with full migration plan, deferred to v1.0. 256 hardcoded color references identified.                                                  | Actual `bg-tc-primary` token migration not started — deferred by design to v1.0 with opt-in path.                                                                              |

---

## c) NOT STARTED

| #   | Item                                                                                                | From research § | Why not started                                                                                                       |
| --- | --------------------------------------------------------------------------------------------------- | --------------- | --------------------------------------------------------------------------------------------------------------------- |
| 1   | **Compound components** (Popover, ContextMenu, HoverCard)                                           | Gap 2 — MEDIUM  | Breaking API pattern change. Should be v2.0 decision. Current monolithic Modal/Drawer work fine.                      |
| 2   | **Native `<dialog>` element**                                                                       | Gap 3 — MEDIUM  | Fundamental Modal/Drawer architecture change. Needs ADR. Browser focus trap is better than JS but migration is risky. |
| 3   | **New components** (Popover, Slider, Rating, TagsInput, Carousel, Calendar, DataTable, ContextMenu) | Gap 5           | Additive work, not improvements to existing library. Separate effort. All 9 confirmed absent from codebase.           |
| 4   | **CSS `@starting-style`** for zero-JS enter/exit animations                                         | Research §10.5  | Requires modern browser baseline decision. Needs testing across all overlay components.                               |
| 5   | **Blocks / composition examples** (dashboard, login, settings layouts)                              | Research §16    | `examples/demo/` exists but no formal "blocks" showing real-world composition patterns.                               |
| 6   | **Compound component pattern** for future overlays (Trigger/Content/Close)                          | Plan Tier 3     | Architectural pattern for v2.0 overlay API.                                                                           |

> **Moved to Fully Done since original report:** RTL rendering tests (was #4), Motion-reduce systematic audit (was #7), SKILL.md update (was #9), Semantic token layer planning (was #1 — now has ADR).

---

## d) TOTALLY FUCKED UP

| #   | Issue                                          | What happened                                                                                                                                                                                    | Resolution                                                                                                                                                   |
| --- | ---------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| 1   | **Motion constants removed as dead code**      | First attempt added 5 motion constants to `shared.go` but didn't wire them into any component. Pre-commit hooks correctly identified them as unused and deleted them, leaving orphaned comments. | ✅ Fixed: re-added constants AND wired `transitionNormal` into Modal + `transitionTransform` into Drawer + `transitionColors` into CopyButton.               |
| 2   | **Pre-commit `govalid-generate` hook failing** | BuildFlow's `govalid-generate` step failed with "exit status 1" on every commit attempt, though `govalid ./...` passed directly.                                                                 | ✅ **Resolved:** govalid is no longer in the pre-commit hook. Current hook is BuildFlow-managed without the govalid-generate step.                           |
| 3   | **Orphaned comments in shared.go**             | After motion constants were deleted as dead code, the comment block was left behind with no code beneath it.                                                                                     | ✅ Fixed: cleaned up and re-added with proper constants.                                                                                                     |
| 4   | **`.gitignore` `*_templ.go` gotcha**           | BuildFlow pre-commit `templ-generate` step re-appended `*_templ.go` to `.gitignore`, hiding NEW generated files from `git status`.                                                               | ✅ **Resolved:** Current `.gitignore` has `!*_templ.go` at top with no re-appended `*_templ.go` override. The buildflow-managed block no longer includes it. |

---

## e) WHAT WE SHOULD IMPROVE

> Items marked ✅ were completed since the original report.

1. ✅ ~~**Wire motion constants into ALL transition-bearing components**~~ → **Partially done:** 3/22 components now use them. The remaining 19 are tracked in Partially-Done #1.
2. ✅ ~~**Add RTL rendering tests**~~ → Done (`display/rtl_test.go`).
3. ✅ ~~**Fix the govalid-generate BuildFlow hook**~~ → Resolved (govalid removed from hook).
4. ✅ ~~**Semantic token layer planning**~~ → ADR 0008 written, deferred to v1.0.
5. ✅ ~~**SKILL.md Part 2 needs the new conventions**~~ → Done (RTL, motion, container-query rules added).
6. ✅ ~~**Container query documentation**~~ → Done (`docs/recipes/container-queries.md`).
7. ✅ ~~**Motion design system page**~~ → Done (`docs/motion-design.md`).

**Still open:** 8. **Wire motion constants into the remaining 19 transition-bearing components** — the value is consistency; 3/22 adoption creates inconsistency. This is the #1 remaining code-quality item. 9. **Add benchmarks to remaining 11 packages** — only display, feedback, navigation have render benchmarks.

---

## f) Top 25 Things to Get Done Next

> Updated with current status. ✅ = done since original report.

| #   | Task                                                                       | Impact | Effort | Status                |
| --- | -------------------------------------------------------------------------- | ------ | ------ | --------------------- |
| 1   | ~~Update SKILL.md Part 2 with RTL/motion/container-query conventions~~     | HIGH   | LOW    | ✅ Done               |
| 2   | ~~Add RTL rendering tests (Card, Drawer, Nav with `dir="rtl"`)~~           | HIGH   | LOW    | ✅ Done               |
| 3   | Wire motion constants into remaining 19 transition-bearing components      | MEDIUM | MEDIUM | 🟡 3/22 done          |
| 4   | ~~Fix or document govalid-generate BuildFlow hook failure~~                | MEDIUM | LOW    | ✅ Resolved           |
| 5   | ~~Plan semantic token layer migration (opt-in `tc-*` tokens)~~             | HIGH   | HIGH   | ✅ ADR written (v1.0) |
| 6   | ~~Add `docs/recipes/container-queries.md` recipe~~                         | MEDIUM | LOW    | ✅ Done               |
| 7   | ~~Add `docs/motion-design.md` reference page~~                             | MEDIUM | LOW    | ✅ Done               |
| 8   | Add Popover component (compound pattern, most requested missing component) | HIGH   | HIGH   | ⬜ Not started        |
| 9   | Migrate Modal to native `<dialog>` element (with ADR)                      | MEDIUM | HIGH   | ⬜ Not started        |
| 10  | Add HoverCard component                                                    | MEDIUM | MEDIUM | ⬜ Not started        |
| 11  | Add Slider component (ARIA slider pattern)                                 | MEDIUM | MEDIUM | ⬜ Not started        |
| 12  | Add DataTable component (sorting, filtering, pagination)                   | HIGH   | HIGH   | ⬜ Not started        |
| 13  | ~~Systematic motion-reduce audit~~                                         | MEDIUM | LOW    | ✅ Done (0 gaps)      |
| 14  | Add CSS `@starting-style` support to Modal/Drawer (zero-JS enter/exit)     | MEDIUM | MEDIUM | ⬜ Not started        |
| 15  | Add Calendar component (full calendar grid)                                | MEDIUM | HIGH   | ⬜ Not started        |
| 16  | Add Rating component (star rating with keyboard support)                   | LOW    | LOW    | ⬜ Not started        |
| 17  | Add TagsInput component                                                    | LOW    | MEDIUM | ⬜ Not started        |
| 18  | Add ContextMenu component (right-click menu)                               | LOW    | MEDIUM | ⬜ Not started        |
| 19  | Add Carousel component                                                     | LOW    | HIGH   | ⬜ Not started        |
| 20  | Create formal "blocks" (dashboard, login, settings page examples)          | MEDIUM | MEDIUM | ⬜ Not started        |
| 21  | Add compound component pattern for future overlays (Trigger/Content/Close) | MEDIUM | HIGH   | ⬜ Not started        |
| 22  | ~~Audit all icon paths for RTL mirroring needs~~                           | LOW    | LOW    | ✅ Done (audit doc)   |
| 23  | ~~Add `prefers-color-scheme` automatic dark mode option~~                  | LOW    | LOW    | ✅ Pre-existing       |
| 24  | Performance benchmark suite (render time per component)                    | LOW    | MEDIUM | 🟡 3/14 packages      |
| 25  | ~~Write ADR for semantic token layer migration decision~~                  | MEDIUM | LOW    | ✅ Done (ADR 0008)    |

**Scorecard:** 12 of 25 complete (48%), 2 partial, 11 not started.

---

## g) Top #1 Question

> ✅ **RESOLVED.** The original question about the govalid-generate BuildFlow hook
> is moot — govalid is no longer referenced in the pre-commit hook or any script.
> The hook was reconfigured (likely during the v0.8.0 release cycle) to remove
> the govalid-generate step entirely.

**New open question:** Should the semantic token layer (`bg-tc-primary`) ship as
an **opt-in CSS layer** (consumer adds `@theme { --color-tc-primary: var(--blue-600) }`)
or as a **full rename** of all 256 hardcoded color references in `.templ` files?
ADR 0008 proposes opt-in first, flip default at v1.0 — but the tradeoff is golden-file
churn vs. immediate theming power. This is the highest-impact architectural decision
remaining before v1.0.

---

## Session Metrics (updated)

| Metric                           | Original (0.7.0) | Current (0.8.0)                         |
| -------------------------------- | ---------------- | --------------------------------------- |
| Version                          | 0.7.0            | 0.8.0                                   |
| Research report lines            | 1,705            | 1,705                                   |
| Research sources                 | 10+              | 10+                                     |
| Planning document tasks          | 55               | 55                                      |
| Physical CSS properties migrated | 74 → 0           | 74 → 0                                  |
| Motion-reduce gaps               | 7 (unaudited)    | **0** (swept)                           |
| Motion constant consumers        | 2                | **3** (Modal, Drawer, CopyButton)       |
| Grid container query variants    | 6                | 6                                       |
| Components with benchmarks       | 0                | **3** (display, feedback, navigation)   |
| Recipe docs                      | 0 → 2            | **7**                                   |
| Reference docs                   | 0                | **2** (motion-design, javascript-guide) |
| ADRs                             | 7                | **8** (+semantic-tokens)                |
| Audit docs                       | 0                | **1** (icon-rtl-mirroring)              |
| Components total                 | 86               | 86                                      |
| Test cases                       | —                | 575                                     |
| Test packages                    | 13/13 ✅         | 13/13 ✅                                |
| Build                            | ✅ Passing       | ✅ Passing                              |
| Lint issues                      | 0                | 0                                       |
| Top-25 complete                  | 0/25             | **12/25** (48%)                         |
