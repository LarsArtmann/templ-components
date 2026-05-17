# Status Report — templ-components

**Date:** 2026-05-17 05:01 CEST
**Branch:** `master` (1 commit ahead of origin)
**Total Commits:** 107
**Visibility:** PUBLIC

---

## Executive Summary

Executed Phase 1 (Critical Bugs) and started Phase 2 (Architecture Type-Safety) from the comprehensive action plan. Fixed 4 P0/P1 bugs, replaced pre-commit hook, added lint exclusion, and refactored the Tabs component to eliminate impossible states. The codebase is now in a **broken commit + fix pending** state — 3 test files need committing to align with the Tabs API change.

**Build: PASSING. Tests: PASSING. Lint: 0 ISSUES.**

---

## a) FULLY DONE ✅

### Phase 1: Critical Bugs & DevOps (7 tasks)

| #   | Task                                   | Status     | Details                                                                                                                                                                                                        |
| --- | -------------------------------------- | ---------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | Fix `NavLinkProps.Attrs` shadowing     | ✅         | Removed `Attrs templ.Attributes` from `NavLinkProps`. Now inherits from `BaseProps.Attrs` correctly.                                                                                                           |
| 2   | Fix Dropdown JS XSS                    | ✅         | Added `dropdownSafeID()` using `strconv.Quote()`. `props.ID` no longer raw-interpolated into JS string.                                                                                                        |
| 3   | Fix Accordion state coupling           | ✅         | Replaced `hidden` attribute + `max-h-96` CSS class with `data-open` attribute for JS state. Removed `hidden` attribute — closed panels use `max-h-0` class only. Server-closed accordions now openable via JS. |
| 4   | Validate required ID in Modal/Dropdown | ✅         | Added `validateDropdownID()` and `validateModalID()` — panic with clear error if ID is empty. Tests updated.                                                                                                   |
| 5   | Fix pre-commit hook                    | ✅         | Replaced buildflow with `scripts/pre-commit.sh` (runs `templ generate` on staged .templ files).                                                                                                                |
| 6   | Tag v0.1.0-alpha                       | ⬜ SKIPPED | Not tagged yet — waiting for all changes to settle before tagging.                                                                                                                                             |
| 7   | Exclude examples/ from lint            | ✅         | Added `exclude-dirs: ["examples"]` to `.golangci.yml`.                                                                                                                                                         |

### Phase 2: Architecture Type-Safety (started, 1 of 8 tasks)

| #   | Task                                    | Status | Details                                                                                                                                                                                                                                     |
| --- | --------------------------------------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 8   | Replace `Tab.Active` with `ActiveTabID` | ✅     | **Breaking change.** Removed `Active bool` from `Tab` struct. Added `ActiveTabID string` to `TabsProps`. Active state computed from `isActiveTab(tabID, activeTabID)`. Impossible state (zero/multiple active tabs) is now unrepresentable. |

### Comprehensive Plan Written

- Created `docs/planning/2026-05-17_01-25_alpha-release-critical-path.md` with 62 tasks across 10 phases, mermaid execution graph, and architectural observations.
- Identified **5 items listed as "Not Done" in previous reports that were actually already implemented.**

---

## b) PARTIALLY DONE 🔨

### Phase 2: Architecture (7 remaining)

| #   | Task                                 | Status | Notes                           |
| --- | ------------------------------------ | ------ | ------------------------------- |
| 9   | Merge BadgeDefault into BadgeNeutral | ⬜     | Identical CSS. Breaking change. |
| 10  | Consolidate badge color maps         | ⬜     | Two maps can drift.             |
| 11  | Fix ErrorAttrs dual reference        | ⬜     | Help text ID not linked.        |
| 12  | Extract tooltip position struct      | ⬜     | Two switches on same type.      |
| 13  | Extract card shell CSS               | ⬜     | Repeated 3×.                    |
| 14  | HTMX CDN URL constant                | ⬜     | Repeated 4×.                    |
| 15  | Error handling magic numbers         | ⬜     | Hardcoded values.               |

---

## c) NOT STARTED ⬜

### Phase 3: JS Unification (5 tasks)

| #   | Task                                                                         | Effort |
| --- | ---------------------------------------------------------------------------- | ------ |
| 16  | Accordion IIFE refactor                                                      | 10 min |
| 17  | Modal IIFE refactor                                                          | 10 min |
| 18  | Dropdown strconv.Quote (partially done — Go helper exists, template updated) | 8 min  |
| 19  | Extract shared dismiss JS                                                    | 10 min |
| 20  | Single-source toast icon paths                                               | 12 min |

### Phase 4: A11y Gaps (3 tasks)

| #   | Task                          | Effort |
| --- | ----------------------------- | ------ |
| 21  | aria-live on error handling   | 8 min  |
| 22  | Avatar status dot scaling     | 8 min  |
| 23  | Document tcShowToast coupling | 10 min |

### Phase 5: Dead Code Cleanup (7 tasks)

| #   | Task                                   | Effort |
| --- | -------------------------------------- | ------ |
| 24  | Remove IconAttrs dead code             | 3 min  |
| 25  | Remove no-op DefaultXxxProps           | 5 min  |
| 26  | Move test helpers to internal/testutil | 10 min |
| 27  | Move ProgressBar a11y test             | 5 min  |
| 28  | Fix TestIconCount hardcoded 45         | 3 min  |
| 29  | Consolidate Exclamation aliases        | 5 min  |
| 30  | Minimal positional → props struct      | 10 min |

### Phase 6: HTMX Decoupling (2 tasks)

| #   | Task                                | Effort |
| --- | ----------------------------------- | ------ |
| 31  | Decouple htmx/loading from feedback | 10 min |
| 32  | FillIcon integration decision       | 10 min |

### Phase 7: Demo App (6 tasks)

| #     | Task                                           | Effort |
| ----- | ---------------------------------------------- | ------ |
| 33    | Delete broken demo                             | 2 min  |
| 34    | Create new demo with layout.Base               | 10 min |
| 35-38 | Add display/feedback/forms/nav/htmx components | 40 min |

### Phase 8: Test Coverage (13 tasks)

| #     | Task                                                | Effort  |
| ----- | --------------------------------------------------- | ------- |
| 39-51 | BDD tests for nav, layout, htmx, icons + edge cases | 150 min |

### Phase 9: Documentation & Growth (7 tasks)

| #     | Task                                                                | Effort |
| ----- | ------------------------------------------------------------------- | ------ |
| 52-58 | Cross-link ecosystem, badges, CHANGELOG, awesome-templ, templ.guide | 50 min |

### Phase 10: Long-Term Polish (4 tasks)

| #     | Task                                                  | Effort |
| ----- | ----------------------------------------------------- | ------ |
| 59-62 | ExampleXxx functions, goreleaser, stroke-width option | 50 min |

---

## d) TOTALLY FUCKED UP 💣

### Commit 662392c Is Broken Without Unstaged Fix

The commit `662392c` includes the new `tabs.templ` (with `ActiveTabID`) but the test files in that commit still use the old `Tab.Active` field. The 3 unstaged test files fix this. **They must be committed before anything else.**

### Demo App Still Broken

`examples/demo/main.go` still uses Tailwind v2 CDN, raw `w.Write`, discards `PageProps`. No work done on this yet.

---

## e) WHAT WE SHOULD IMPROVE

### Immediate

1. **Commit the 3 test fixes** — The repo is in a broken state without them.
2. **Continue Phase 2 architecture** — Badge consolidation, tooltip extraction, card shell CSS, HTMX constants. These are all 5-10 min tasks that reduce drift risk.
3. **Then delete and rebuild the demo** — The current demo is a liability. A v2 CDN demo on a v4-exclusive library makes us look amateurish.

### Process

4. **Don't commit half-done refactors** — The Tabs change should have been committed as a unit (templ + tests + docs). The commit hook would have caught this if tests were run.

---

## f) Top 25 Next Actions

| #   | Task                                 | Impact           | Effort | Category     |
| --- | ------------------------------------ | ---------------- | ------ | ------------ |
| 1   | Commit 3 test fixes (unbreak repo)   | BLOCKS ALL       | 2 min  | Critical     |
| 2   | Merge BadgeDefault into BadgeNeutral | Clean API        | 8 min  | Architecture |
| 3   | Consolidate badge color maps         | Drift prevention | 10 min | Architecture |
| 4   | Extract tooltip position struct      | DRY              | 10 min | Architecture |
| 5   | Extract card shell CSS               | DRY              | 8 min  | Architecture |
| 6   | HTMX CDN URL constant                | DRY              | 5 min  | Architecture |
| 7   | Error handling magic numbers         | Readability      | 5 min  | Architecture |
| 8   | Fix ErrorAttrs dual reference        | A11y             | 10 min | Bug          |
| 9   | Accordion IIFE refactor              | Consistency      | 10 min | Architecture |
| 10  | Modal IIFE refactor                  | Consistency      | 10 min | Architecture |
| 11  | Extract shared dismiss JS            | DRY              | 10 min | Architecture |
| 12  | Single-source toast icon paths       | Drift prevention | 12 min | Architecture |
| 13  | Avatar status dot scaling            | A11y             | 8 min  | A11y         |
| 14  | aria-live on error handling          | A11y             | 8 min  | A11y         |
| 15  | Remove IconAttrs dead code           | Cleanliness      | 3 min  | Cleanup      |
| 16  | Remove no-op DefaultXxxProps         | Cleanliness      | 5 min  | Cleanup      |
| 17  | Fix TestIconCount hardcoded 45       | Maint            | 3 min  | Cleanup      |
| 18  | Consolidate Exclamation aliases      | Clarity          | 5 min  | Cleanup      |
| 19  | Delete and rebuild demo app          | Trust            | 50 min | DX           |
| 20  | Decouple htmx/loading from feedback  | Architecture     | 10 min | Architecture |
| 21  | Minimal positional → props struct    | Consistency      | 10 min | Architecture |
| 22  | Tag v0.1.0-alpha                     | Semver           | 3 min  | Release      |
| 23  | Cross-link ecosystem in README       | Differentiation  | 8 min  | Growth       |
| 24  | Submit to awesome-templ              | Discovery        | 10 min | Growth       |
| 25  | BDD tests for navigation             | Coverage         | 36 min | Testing      |

---

## g) My Top #1 Question

**Should we keep `BadgeDefault` as a deprecated alias for `BadgeNeutral`, or remove it entirely?**

`BadgeDefault` and `BadgeNeutral` produce identical CSS. `BadgeDefault` is the zero-value of `BadgeType` and is used by `DefaultBadgeProps()`. Removing it is a breaking change. Options:

1. **Remove `BadgeDefault` entirely** — Replace with `BadgeNeutral` everywhere. Breaking but honest.
2. **Deprecate `BadgeDefault`** — Keep the constant, add deprecation comment, point to `BadgeNeutral`. Non-breaking.
3. **Remove `BadgeDefault`, change default to `BadgeNeutral`** — Breaking but the API becomes honest.

The existing test at `a11y_test.go:146` asserts `props.Type != BadgeDefault`, which would need updating regardless. I recommend option 3 but it's your call on breaking changes.

---

## Session Metrics

| Metric                      | Start of Session | Current           |
| --------------------------- | ---------------- | ----------------- |
| Commits                     | 105              | 107               |
| P0 bugs                     | 2                | 0                 |
| P1 bugs                     | 4                | 1 (demo)          |
| Impossible states           | 1 (Tab.Active)   | 0                 |
| Lint issues                 | 0                | 0                 |
| Tests passing               | 146              | 146               |
| Pre-commit hook             | 💣 buildflow     | ✅ templ generate |
| Examples excluded from lint | ❌               | ✅                |
| Planning doc                | ❌               | ✅ 62 tasks       |
