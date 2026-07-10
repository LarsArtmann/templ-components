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

# Comprehensive Execution Plan — Post-Bug-Hunt Sprint

**Created:** 2026-07-07 18:10
**Source:** `docs/status/2026-07-07_18-02_bug-hunt-critical-fixes.md` (items a–g)
**Rule:** Every task ≤ 12 minutes. Sorted by Impact → Customer-Value → Effort (Pareto).
**Total tasks:** 89 across 9 tiers.

---

## Sorting Key

| Impact      | Meaning                                                                         |
| ----------- | ------------------------------------------------------------------------------- |
| 🔴 CRITICAL | Uncommitted work at risk, or unaudited code may contain more ship-breaking bugs |
| 🟠 HIGH     | Prevents recurrence of the bug classes found this session, or blocks release    |
| 🟡 MED      | Improves quality, DX, or coverage                                               |
| ⚪ LOW      | Polish, submissions, nice-to-haves                                              |

| Effort | Range                            |
| ------ | -------------------------------- |
| XS     | ≤ 5 min                          |
| S      | 5–12 min                         |
| M      | 12–30 min (split into sub-tasks) |
| L      | 30+ min (split into sub-tasks)   |

---

## Tier 1 — CRITICAL: Protect & Verify Current Work (7 tasks)

| #   | Task                                                                                                            | Impact | Effort | Est     | Prereq |
| --- | --------------------------------------------------------------------------------------------------------------- | ------ | ------ | ------- | ------ |
| 1   | Commit forms/feedback bug fixes (Toggle, Combobox×3, Select×2, Checkbox, Toast, ProgressBar)                    | 🔴     | S      | 8m      | —      |
| 2   | Commit display bug fixes (overlay aria/inert, dropdown RTL, accordion grid, tabs EnsureID, copybutton, tooltip) | 🔴     | S      | 8m      | 1      |
| 3   | Commit navigation fix (pagination dynamic class)                                                                | 🔴     | XS     | 3m      | 2      |
| 4   | Commit demo + docs (CHANGELOG, AGENTS.md, demo TOC/anchor fixes)                                                | 🔴     | S      | 8m      | 3      |
| 5   | Run `nix run .#verify` to confirm canonical check passes                                                        | 🔴     | XS     | 5m      | 4      |
| 6   | ~~Verify `grid-rows-[0fr]` produces correct CSS in Tailwind v4~~ **✅ VERIFIED (2026-07-10)**                   | 🔴     | S      | ~~10m~~ | —      |
| 7   | Fix `.gitignore` if BuildFlow re-added `*_templ.go` after regen                                                 | 🔴     | XS     | 3m      | 4      |

---

## Tier 2 — HIGH: Audit Remaining Packages (18 tasks)

> These packages have NOT been line-by-line audited. The 3 audited packages yielded 17 bugs.
> Statistically, the remaining 5 packages likely contain more defects.

| #   | Task                                                                                 | Impact | Effort | Est      | Prereq |
| --- | ------------------------------------------------------------------------------------ | ------ | ------ | -------- | ------ |
| 8   | Audit `htmx/loading.templ` (LoadingIndicator, InlineLoadingOverlay) for bugs         | 🟠     | S      | 10m      | —      |
| 9   | Audit `htmx/helpers.templ` (LoadingButton, ConfirmDelete) for bugs                   | 🟠     | S      | 10m      | —      |
| 10  | Audit `htmx/error_handling.templ` (SwapOOB, CSRFToken, GlobalErrorHandling) for bugs | 🟠     | S      | 12m      | —      |
| 11  | Audit `errorpage/errorpage.templ` for rendering bugs                                 | 🟠     | S      | 10m      | —      |
| 12  | Audit `errorpage/notfound404.templ` for rendering bugs                               | 🟠     | S      | 10m      | —      |
| 13  | Audit `errorpage/errordetail.templ` + `erroralert.templ` for bugs                    | 🟠     | S      | 10m      | —      |
| 14  | Audit `errorpage/handler.go` + `fromerror.go` + `constructors.go` for logic bugs     | 🟠     | S      | 12m      | —      |
| 15  | Audit `layout/base.templ` for rendering bugs (head, meta, CSP)                       | 🟠     | S      | 10m      | —      |
| 16  | Audit `layout/theme.templ` (ThemeScript, ThemeToggle) for bugs                       | 🟠     | S      | 8m       | —      |
| 17  | Audit `layout/script.templ` + `stylesheet.templ` for bugs                            | 🟠     | XS     | 5m       | —      |
| 18  | Audit `navigation/nav.templ` + `nav_link.templ` for bugs                             | 🟠     | S      | 10m      | —      |
| 19  | Audit `navigation/mobile_menu.templ` for bugs                                        | 🟠     | XS     | 8m       | —      |
| 20  | Audit `navigation/breadcrumbs.templ` for bugs                                        | 🟠     | XS     | 5m       | —      |
| 21  | Audit `navigation/sidebar_nav.templ` for bugs                                        | 🟠     | XS     | 8m       | —      |
| 22  | Audit `navigation/loadmore.templ` for bugs                                           | 🟠     | XS     | 5m       | —      |
| 23  | Audit `forms/input_group.templ` (was missed in first pass)                           | 🟠     | XS     | 8m       | —      |
| 24  | Audit `forms/date_picker.templ` + `file_input.templ` for bugs                        | 🟠     | XS     | 8m       | —      |
| 25  | Fix any bugs found in Tier 2 audits (estimate 5–10 bugs)                             | 🟠     | M      | 12m each | 8–24   |

---

## Tier 3 — HIGH: Prevent Recurrence (5 tasks)

| #   | Task                                                                                                                 | Impact | Effort | Est | Prereq |
| --- | -------------------------------------------------------------------------------------------------------------------- | ------ | ------ | --- | ------ |
| 26  | Add CI check script: grep for dynamic Tailwind class concatenation (`"peer-"+`, `"hover:"+`, etc.) in `.templ` files | 🟠     | S      | 12m | —      |
| 27  | Wire the dynamic-class CI check into `utils/` as a test (like motion-reduce compliance test)                         | 🟠     | S      | 10m | 26     |
| 28  | Add test: CopyButton `type="submit"` edge case (preventDefault doesn't block form submit)                            | 🟠     | XS     | 5m  | —      |
| 29  | Add test: Combobox full lifecycle (render → type → select → verify hidden value)                                     | 🟠     | S      | 12m | —      |
| 30  | Add test: Accordion with very long content (>1000px) renders without clipping                                        | 🟠     | XS     | 5m  | —      |

---

## Tier 4 — HIGH: Release (4 tasks)

| #   | Task                                                                                                    | Impact | Effort | Est | Prereq |
| --- | ------------------------------------------------------------------------------------------------------- | ------ | ------ | --- | ------ |
| 31  | Verify `[Unreleased]` CHANGELOG has body (release script requirement)                                   | 🟠     | XS     | 3m  | 1–4    |
| 32  | Run `scripts/release.sh 0.10.0 "17 bug fixes: toggle, combobox, select, overlay a11y, accordion, tabs"` | 🟠     | S      | 10m | 5, 31  |
| 33  | Review release commit with `git show v0.10.0`                                                           | 🟠     | XS     | 5m  | 32     |
| 34  | Push to remote (requires user approval)                                                                 | 🟠     | XS     | 2m  | 33     |

---

## Tier 5 — MED: Browser Verification (5 tasks)

| #   | Task                                                                             | Impact | Effort | Est | Prereq |
| --- | -------------------------------------------------------------------------------- | ------ | ------ | --- | ------ |
| 35  | Start demo server, manually verify Toggle thumb slides when checked              | 🟡     | XS     | 5m  | 5      |
| 36  | Manually verify Accordion opens/closes with long content                         | 🟡     | XS     | 5m  | 5      |
| 37  | Manually verify Modal open→close→focus-restore cycle                             | 🟡     | XS     | 8m  | 5      |
| 38  | Manually verify Tooltip text is announced (browser devtools accessibility panel) | 🟡     | XS     | 8m  | 5      |
| 39  | Manually verify Pagination arrows render with correct rounding                   | 🟡     | XS     | 5m  | 5      |

---

## Tier 6 — MED: CI & Testing Infrastructure (8 tasks)

| #   | Task                                                                                | Impact | Effort | Est | Prereq |
| --- | ----------------------------------------------------------------------------------- | ------ | ------ | --- | ------ |
| 40  | Add fuzz tests to CI: `go test -fuzz=. -run=Fuzz ./...` step (30s timeout)          | 🟡     | S      | 10m | —      |
| 41  | Add Dropdown keyboard nav test (full menu item cycle: ArrowDown/Up/Home/End/Escape) | 🟡     | S      | 12m | —      |
| 42  | Add Tabs client-side keyboard nav test (ArrowLeft/Right/Home/End in LTR + RTL)      | 🟡     | S      | 12m | —      |
| 43  | Add Combobox fuzz test (keyboard handling never panics)                             | 🟡     | XS     | 8m  | —      |
| 44  | Add integration test for overlay open→close→open lifecycle (DOM assertions)         | 🟡     | M      | 12m | —      |
| 45  | Add test for Tooltip with non-standard trigger element (custom component)           | 🟡     | XS     | 8m  | —      |
| 46  | Run `art-dupl` on Go sources (`*_templ.go` + handwritten `.go`)                     | 🟡     | XS     | 10m | —      |
| 47  | Verify `inert` attribute renders correctly in templ (boolean attribute, no value)   | 🟡     | XS     | 5m  | —      |

---

## Tier 7 — MED: Documentation & Polish (12 tasks)

| #   | Task                                                                                        | Impact | Effort | Est | Prereq |
| --- | ------------------------------------------------------------------------------------------- | ------ | ------ | --- | ------ |
| 48  | Document CopyButton `<a>` variant navigation suppression (preventDefault behavior)          | 🟡     | XS     | 5m  | —      |
| 49  | Document Tabs `resolveActiveTabID` behavior change in migration guide                       | 🟡     | XS     | 5m  | —      |
| 50  | Document Tooltip `aria-describedby` propagation selector in godoc                           | 🟡     | XS     | 5m  | —      |
| 51  | Fix SKILL.md component count discrepancy (82 vs 83)                                         | 🟡     | XS     | 5m  | —      |
| 52  | Archive completed planning docs — add `STATUS: COMPLETED` headers                           | 🟡     | S      | 10m | —      |
| 53  | Add "bug hunt checklist" doc from lessons learned this session                              | 🟡     | S      | 12m | —      |
| 54  | Create a "bug classes" reference: dynamic Tailwind classes, slice mutation, aria sync, etc. | 🟡     | S      | 12m | —      |
| 55  | Update README with v0.10.0 bug fix highlights (if releasing)                                | 🟡     | XS     | 8m  | 32     |
| 56  | Review all `data-tc-*` attribute names for consistency                                      | 🟡     | S      | 10m | —      |
| 57  | Verify color contrast ratios meet WCAG AA (automated check)                                 | 🟡     | S      | 12m | —      |
| 58  | Add `role="group"` audit for form field grouping a11y                                       | 🟡     | S      | 10m | —      |
| 59  | Write a test for `forms.InputGroup` prefix/suffix rendering                                 | 🟡     | XS     | 8m  | —      |

---

## Tier 8 — MED: New Features & Refactoring (12 tasks)

| #   | Task                                                                                              | Impact | Effort | Est   | Prereq |
| --- | ------------------------------------------------------------------------------------------------- | ------ | ------ | ----- | ------ |
| 60  | Wire motion constants into 5 components (Accordion, Tooltip, Table, Avatar, Badge)                | 🟡     | S      | 12m   | —      |
| 61  | Wire motion constants into 5 more components (Dropdown, Pagination, Nav, Breadcrumbs, SidebarNav) | 🟡     | S      | 12m   | 60     |
| 62  | Wire motion constants into remaining 9 components                                                 | 🟡     | M      | 12m×2 | 61     |
| 63  | Refactor overlay JS generators to use `text/template` (part 1: extract template)                  | 🟡     | S      | 12m   | —      |
| 64  | Refactor overlay JS generators to use `text/template` (part 2: wire + test)                       | 🟡     | S      | 12m   | 63     |
| 65  | Self-host htmx.js in examples/demo (download + commit, no CDN)                                    | ⚪     | XS     | 10m   | —      |
| 66  | Add standalone `/forms` quickstart demo route                                                     | 🟡     | M      | 12m×2 | —      |
| 67  | Add blocks/composition examples (dashboard layout)                                                | 🟡     | M      | 12m×2 | —      |
| 68  | Add blocks/composition examples (login layout)                                                    | 🟡     | S      | 12m   | —      |
| 69  | Add blocks/composition examples (settings layout)                                                 | 🟡     | S      | 12m   | —      |
| 70  | Add "doc reality" CI check (verify AGENTS.md claims match filesystem)                             | 🟡     | M      | 12m×2 | —      |

---

## Tier 9 — LOW: Community & Future (19 tasks)

| #   | Task                                                                                     | Impact | Effort | Est   | Prereq |
| --- | ---------------------------------------------------------------------------------------- | ------ | ------ | ----- | ------ |
| 71  | awesome-templ PR submission (updated component count)                                    | ⚪     | XS     | 5m    | 34     |
| 72  | templ.guide listing submission                                                           | ⚪     | XS     | 5m    | 34     |
| 73  | Configure SSH tag signing (`gpg.ssh.allowedSignersFile`)                                 | ⚪     | XS     | 10m   | —      |
| 74  | Pagination RTL icon visual swap (ArrowLeft/ArrowRight in RTL)                            | 🟡     | XS     | 10m   | —      |
| 75  | Add Popover component — types + props struct                                             | 🟠     | M      | 12m   | —      |
| 76  | Add Popover component — template + JS                                                    | 🟠     | M      | 12m   | 75     |
| 77  | Add Popover component — tests (golden, a11y, BDD, edge)                                  | 🟠     | M      | 12m   | 76     |
| 78  | Add SortableDataTable — types + wrapper logic                                            | 🟠     | M      | 12m   | —      |
| 79  | Add SortableDataTable — template + integration test                                      | 🟠     | M      | 12m   | 78     |
| 80  | Add FilterDropdown component                                                             | 🟡     | M      | 12m×2 | —      |
| 81  | Design `Validate() error` pattern for v1.0 (document interface, pick 3 pilot components) | 🟡     | M      | 12m×2 | —      |
| 82  | Implement `Validate() error` on 3 pilot components                                       | 🟡     | M      | 12m×2 | 81     |
| 83  | Move test helpers to `internal/testutil/` (part 1: create package + move golden)         | 🟡     | M      | 12m   | —      |
| 84  | Move test helpers to `internal/testutil/` (part 2: move Render/AssertContains)           | 🟡     | M      | 12m   | 83     |
| 85  | Move test helpers to `internal/testutil/` (part 3: update all consumers)                 | 🟡     | M      | 12m×2 | 84     |
| 86  | Semantic token layer — document phase (ADR 0008 update)                                  | 🟡     | S      | 12m   | —      |
| 87  | Consumer validation: adopt templ-components in a real project (DiscordSync)              | 🟠     | L      | 12m×5 | 34     |
| 88  | Add Playwright/screenshot test infrastructure (part 1: setup)                            | 🟡     | L      | 12m×3 | —      |
| 89  | Add Playwright/screenshot test infrastructure (part 2: overlay lifecycle tests)          | 🟡     | L      | 12m×2 | 88     |

---

## Summary by Tier

| Tier      | Theme                         | Tasks  | Total Est  | Blocking?                         |
| --------- | ----------------------------- | ------ | ---------- | --------------------------------- |
| 1         | Protect & verify current work | 7      | ~40m       | YES — blocks everything           |
| 2         | Audit remaining packages      | 18     | ~180m      | YES — may find more critical bugs |
| 3         | Prevent recurrence            | 5      | ~44m       | No                                |
| 4         | Release v0.10.0               | 4      | ~20m       | Depends on 1–3                    |
| 5         | Browser verification          | 5      | ~31m       | No                                |
| 6         | CI & test infra               | 8      | ~79m       | No                                |
| 7         | Documentation & polish        | 12     | ~100m      | No                                |
| 8         | New features & refactoring    | 12     | ~140m+     | No                                |
| 9         | Community & future            | 19     | ~250m+     | No                                |
| **Total** |                               | **89** | **~884m+** |                                   |

## Pareto Distribution

- **20% of tasks (Tier 1+2 = 25 tasks)** deliver **80% of value** — commit safety + complete bug audit
- **Tier 3+4 (9 tasks)** prevent the bug classes from recurring and ship the fix to consumers
- **Tier 5–9 (55 tasks)** are quality-of-life, new features, and long-term investments

## Execution Order

```
Tier 1 (commit + verify)
    ↓
Tier 2 (audit remaining packages — fix bugs as found)
    ↓
Tier 3 (prevent recurrence CI checks)
    ↓
Tier 4 (release v0.10.0)
    ↓
Tier 5 (browser verification)
    ↓
Tiers 6–9 (parallel, no dependencies between them)
```
