# Comprehensive Execution Plan — 2026-05-20

**Generated:** 2026-05-20 | **Scope:** All remaining work across docs/planning, docs/status, docs/STANDOUT-IDEAS, docs/SUPERB-FOR-PERSONAL-USE, codebase audit
**Source files analyzed:** 9 planning/status docs + live codebase audit
**Total tasks:** 87 | **Estimated total:** ~17 hours

## Legend

- **Priority:** P0 (critical bug) → P1 (high) → P2 (medium) → P3 (low) → P4 (ecosystem) → P5 (personal-use)
- **Impact:** 🔴 High (affects many users/components) → 🟡 Medium → 🟢 Low
- **Effort:** ⏱️ estimated minutes
- **Status:** ⬜ Not started | ✅ Done (prior session) | ❌ Deferred

## Previously Completed (Session 11 — DO NOT RE-PLAN)

These items are DONE but still marked ⬜ in TODO_LIST.md (stale):

| Done Item | Commit |
|-----------|--------|
| Demo app rewrite with layout.Base + Tailwind v4 | `0a2d754` |
| FeedbackType unification (AlertType/ToastType → aliases) | `712be15` |
| LoadingOverlay → props struct | `712be15` |
| StepIndicator BaseProps | `712be15` |
| FillIcon variadic → bool | `712be15` |
| ThemeToggle multi-instance fix | `9cc6d10` |
| Modal stable IDs | `9cc6d10` |
| Tooltip aria-describedby | `9cc6d10` |
| Breadcrumbs icon system | `f19dd42` |
| Pagination URL builder (net/url) | `f19dd42` |
| Mutex restored (tailwind-merge-go NOT thread-safe) | `f19dd42` |
| Test cleanup (splitClasses, benchmarks) | `a23b17f` |
| CONTRIBUTING.md fix | `e37a36d` |

## Explicitly Deferred (DO NOT INCLUDE IN PLAN)

| Item | Reason |
|------|--------|
| ComponentProps interface (GetBaseProps on 29 structs) | Significant boilerplate, low current ROI |
| DropdownItem typed variant | Breaking API change, low immediate value |
| JS consolidation (shared tc-init.js) | Risky runtime behavior change across 7 components |
| SimpleCard compose through Card | Current implementation is cleaner |
| SwapOOB swapStyle validation | Minimal misuse risk |
| Consolidate test files (37→15) | Post-v1.0 |
| Convert snapshot tests to golden files | Post-v1.0 |
| Move test helpers out of utils/ | Post-v1.0 |
| Documentation site generation | Post-v1.0 |
| Radio/File input/Toggle components | Post-v1.0 |
| Client-side JS tab switching | Post-v1.0 |
| PageProps zero-value safety | Post-v1.0 |
| uint for Pagination fields | Post-v1.0 |
| Icon list auto-gen (split brain) | Post-v1.0 |
| Make GlobalErrorHandling configurable | Post-v1.0 |
| Modularization (go.work) | Analysis concluded NOT recommended |

---

## WAVE 1 — Critical Fixes & Quick Wins (P0/P1)

### T1: Fix Modal Focus Restore — P0 Bug

**Problem:** `tcCloseModal` looks for `id + "-prev-focus"` element that is never created. Focus is never restored after closing modal. **WCAG failure.**

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T1.1 | Save `document.activeElement` in `tcOpenModal` to `dataset.prevFocus` on modal wrapper | 10m | 🔴 |
| T1.2 | In `tcCloseModal`, read saved element ref and call `.focus()` | 8m | 🔴 |
| T1.3 | Add test: open modal → close → verify focus restored to trigger | 10m | 🔴 |
| T1.4 | Run full test suite + lint | 4m | 🟢 |

**Total: ~32min**

### T2: Fix ID Propagation on 6 Components — P1 Consistency

**Problem:** Alert, Toast, StatCard, Nav, Dropdown, ProgressBar embed `BaseProps` but never render `id={ props.ID }` on root element.

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T2.1 | Add `id={ props.ID }` to Alert root `<div>` | 5m | 🟡 |
| T2.2 | Add `id={ props.ID }` to Toast root `<div>` | 5m | 🟡 |
| T2.3 | Add `id={ props.ID }` to StatCard root `<div>` | 5m | 🟡 |
| T2.4 | Add `id={ props.ID }` to Nav root `<nav>` | 5m | 🟡 |
| T2.5 | Add `id={ props.ID }` to Dropdown wrapper `<div>` | 5m | 🟡 |
| T2.6 | Add `id={ props.ID }` to ProgressBar root `<div>` | 5m | 🟡 |
| T2.7 | Add snapshot tests for ID propagation on each | 10m | 🟡 |
| T2.8 | templ generate + build + test + lint | 5m | 🟢 |

**Total: ~45min**

### T3: Fix Icon System Issues — P1 Correctness

**Problem:** IconPathJS stroke-width mismatch, unknown icons silently fallback, deprecated Exclamation icon.

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T3.1 | Fix IconPathJS: use `stroke-width="1.5"` to match templ rendering | 8m | 🔴 |
| T3.2 | Add validation: unknown icon name → panic with clear message | 10m | 🔴 |
| T3.3 | Delete deprecated `Exclamation` from constants, path map, tests | 8m | 🟡 |
| T3.4 | Add test: `Name("nonexistent")` panics | 5m | 🔴 |
| T3.5 | Run tests + lint | 4m | 🟢 |

**Total: ~35min**

### T4: Add BaseProps to Breadcrumbs — P1 API Consistency

**Problem:** Breadcrumbs takes `[]BreadcrumbItem` only — no BaseProps. Can't add ID, class, attrs, nonce.

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T4.1 | Create `BreadcrumbsProps` struct with BaseProps + `Items []BreadcrumbItem` | 10m | 🟡 |
| T4.2 | Create `DefaultBreadcrumbsProps()` constructor | 5m | 🟢 |
| T4.3 | Update `Breadcrumbs(BreadcrumbsProps)` signature | 10m | 🟡 |
| T4.4 | Update template: propagate ID, Class, Attrs on root `<nav>` | 8m | 🟡 |
| T4.5 | Update all call sites (demo app, tests) | 10m | 🟡 |
| T4.6 | templ generate + build + test + lint | 5m | 🟢 |

**Total: ~48min**

### T5: Add BaseProps to SimpleNav — P1 API Consistency

**Problem:** SimpleNav takes individual string args — no BaseProps.

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T5.1 | Create `SimpleNavProps` struct with BaseProps + current fields | 10m | 🟡 |
| T5.2 | Create `DefaultSimpleNavProps()` constructor | 5m | 🟢 |
| T5.3 | Update `SimpleNav(SimpleNavProps)` signature | 10m | 🟡 |
| T5.4 | Update template: propagate ID, Class, Attrs | 8m | 🟡 |
| T5.5 | Update all call sites (demo app, tests) | 10m | 🟡 |
| T5.6 | templ generate + build + test + lint | 5m | 🟢 |

**Total: ~48min**

### T6: Add BaseProps to Spinner — P1 API Consistency

**Problem:** Spinner takes `(size, colorClass)` — no BaseProps.

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T6.1 | Create `SpinnerProps` struct with BaseProps + Size + ColorClass | 8m | 🟡 |
| T6.2 | Create `DefaultSpinnerProps()` constructor | 4m | 🟢 |
| T6.3 | Update `Spinner(SpinnerProps)` + template propagation | 8m | 🟡 |
| T6.4 | Update all call sites (HTMX loading, tests, demo) | 10m | 🟡 |
| T6.5 | templ generate + build + test + lint | 5m | 🟢 |

**Total: ~35min**

---

## WAVE 2 — Code Quality & Safety (P2)

### T7: Remove Dead Code — P2 Cleanup

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T7.1 | Remove `Deref`, `DerefOr`, `MergeAttrs` from utils.go | 5m | 🟢 |
| T7.2 | Remove associated tests | 5m | 🟢 |
| T7.3 | Remove `TestPtr` from utils_test.go | 3m | 🟢 |
| T7.4 | Remove unused `badgeTextLive` from display/badge_test.go | 3m | 🟢 |
| T7.5 | Build + test + lint | 4m | 🟢 |

**Total: ~20min**

### T8: Replace Hardcoded SVGs with Icon System — P2 Consistency

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T8.1 | Alert dismiss X: replace inline SVG with `@icons.Icon(icons.X, ...)` | 8m | 🟡 |
| T8.2 | Toast JS dismiss: build SVG string from `icons.IconPathJS()` | 12m | 🟡 |
| T8.3 | StepIndicator checkmark: replace inline SVG with `@icons.Icon(icons.Check, ...)` | 8m | 🟡 |
| T8.4 | templ generate + build + test + lint | 4m | 🟢 |

**Total: ~32min**

### T9: Input Validation — P2 Safety

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T9.1 | Clamp ProgressBar percent to [0, 100] with max() | 5m | 🟡 |
| T9.2 | Validate Pagination CurrentPage > 0 (clamp to 1) | 5m | 🟡 |
| T9.3 | Validate SelectOption Disabled+Selected contradiction | 8m | 🟡 |
| T9.4 | Validate `\|` separator in SVG icon paths | 5m | 🟢 |
| T9.5 | Add tests for each validation | 10m | 🟡 |
| T9.6 | Build + test + lint | 4m | 🟢 |

**Total: ~37min**

### T10: BoolString → strconv.FormatBool — P2 Dedup

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T10.1 | Replace all BoolString() calls with strconv.FormatBool() | 8m | 🟢 |
| T10.2 | Delete BoolString from utils.go | 3m | 🟢 |
| T10.3 | Update tests | 5m | 🟢 |
| T10.4 | Build + test + lint | 3m | 🟢 |

**Total: ~19min**

### T11: Extract Shared Lookup Helper — P2 Dedup

**Problem:** Map-lookup-with-default-fallback pattern repeated 5× across display package.

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T11.1 | Create generic `LookupStyle[T comparable](m map[T]style, key T, fallback style) style` | 8m | 🟡 |
| T11.2 | Refactor badgeSizeLookup, cardPaddingLookup, etc. to use it | 12m | 🟡 |
| T11.3 | Update tests | 8m | 🟢 |
| T11.4 | Build + test + lint | 4m | 🟢 |

**Total: ~32min**

### T12: Extract Shared Dismiss JS — P2 Dedup

**Problem:** Alert + Toast have duplicate dismiss handler logic.

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T12.1 | Extract shared `tcDismissAttached` handler (already partially done per AGENTS.md) | 8m | 🟡 |
| T12.2 | Verify Alert + Toast both use it | 5m | 🟢 |
| T12.3 | Test dismiss on both components | 8m | 🟢 |

**Total: ~21min**

### T13: Test Coverage Gaps — P2 Quality

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T13.1 | Alert: empty Title, empty Message, unknown AlertType | 8m | 🟡 |
| T13.2 | Toast: empty Message, unknown ToastType | 8m | 🟡 |
| T13.3 | ProgressBar: Total=0, negative Current | 5m | 🟡 |
| T13.4 | StepIndicator: empty Steps, CurrentStep out of bounds | 8m | 🟡 |
| T13.5 | Nav: empty Links, empty Href | 5m | 🟡 |
| T13.6 | Dropdown: empty Items, item with both Href + action | 8m | 🟡 |
| T13.7 | Modal: without Title | 5m | 🟡 |
| T13.8 | CSRFToken: empty token string | 3m | 🟢 |
| T13.9 | Pagination: CurrentPage > TotalPages, TotalPages=0 | 5m | 🟡 |
| T13.10 | Add `t.Parallel()` to TestSecurityHeaders | 3m | 🟢 |

**Total: ~58min**

### T14: Remove Duplicate Test Data — P2 Cleanup

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T14.1 | Extract shared `testNavLinks` to navigation test helper file | 8m | 🟢 |
| T14.2 | Update all navigation test files to use shared helper | 8m | 🟢 |
| T14.3 | Build + test + lint | 4m | 🟢 |

**Total: ~20min**

### T15: Add BaseProps to SimpleEmptyState — P2 Consistency

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T15.1 | Update `SimpleEmptyState(message string)` to `SimpleEmptyState(SimpleEmptyStateProps)` | 8m | 🟢 |
| T15.2 | Create `DefaultSimpleEmptyStateProps()` | 3m | 🟢 |
| T15.3 | Update call sites + tests | 5m | 🟢 |
| T15.4 | Build + test + lint | 3m | 🟢 |

**Total: ~19min**

---

## WAVE 3 — JS Fixes (P2)

### T16: Fix JS Re-attachment After HTMX Swaps — P2 Runtime Bug

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T16.1 | Audit all `window.tc*Attached` guards across 7 components | 8m | 🔴 |
| T16.2 | Replace global guards with per-element `data-tc-initialized` attribute | 12m | 🔴 |
| T16.3 | Update Accordion, Dropdown, Modal, ThemeToggle init scripts | 10m | 🟡 |
| T16.4 | Test: render component → HTMX swap → verify listeners re-attach | 12m | 🔴 |
| T16.5 | Build + test + lint | 4m | 🟢 |

**Total: ~46min**

### T17: Fix GlobalErrorHandling Retry Counter Race — P2 Concurrency

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T17.1 | Scope `retryCount` per-request (use closure or data attribute) | 10m | 🟡 |
| T17.2 | Test concurrent error handling | 8m | 🟡 |
| T17.3 | Build + test + lint | 4m | 🟢 |

**Total: ~22min**

---

## WAVE 4 — Documentation & Polish (P3)

### T18: Update Stale TODO_LIST.md — P3 Accuracy

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T18.1 | Mark all session 11 completed items as ✅ | 8m | 🟢 |
| T18.2 | Remove stale notes (e.g., TODO #11 dropdownSafeID) | 5m | 🟢 |
| T18.3 | Add new items from this plan not yet tracked | 5m | 🟢 |

**Total: ~18min**

### T19: Document Conventions — P3 Knowledge Sharing

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T19.1 | Document htmx→feedback runtime JS coupling (code comment) | 5m | 🟢 |
| T19.2 | Document fill vs stroke 20×20/24×24 convention (code comment) | 5m | 🟢 |
| T19.3 | Document thread-safety requirement in CONTRIBUTING.md | 5m | 🟢 |
| T19.4 | Add ADR for FeedbackType unification (already done, document decision) | 8m | 🟢 |

**Total: ~23min**

### T20: Add go doc Examples — P3 Discoverability

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T20.1 | Add `ExampleAlert()` test function to feedback package | 8m | 🟡 |
| T20.2 | Add `ExampleBadge()` to display package | 5m | 🟡 |
| T20.3 | Add `ExampleCard()` to display package | 5m | 🟡 |
| T20.4 | Add `ExamplePagination()` to navigation package | 8m | 🟡 |
| T20.5 | Add `ExampleIcon()` to icons package | 5m | 🟡 |
| T20.6 | Build + test + verify on pkg.go.dev output | 5m | 🟡 |

**Total: ~36min**

### T21: Update FEATURES.md — P3 Accuracy

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T21.1 | Run features audit against current codebase state | 10m | 🟢 |
| T21.2 | Update FEATURES.md with accurate status indicators | 8m | 🟢 |
| T21.3 | Verify coverage numbers match current | 3m | 🟢 |

**Total: ~21min**

### T22: Update AGENTS.md — P3 Memory

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T22.1 | Update coverage, test count, completed items | 5m | 🟢 |
| T22.2 | Add new conventions discovered during this session | 5m | 🟢 |
| T22.3 | Remove stale entries | 3m | 🟢 |

**Total: ~13min**

---

## WAVE 5 — Ecosystem & Visibility (P4)

### T23: Tag v0.1.0-alpha — P4 Milestone

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T23.1 | Final verification: build + test + lint + coverage | 5m | 🔴 |
| T23.2 | Update CHANGELOG.md or create release notes | 10m | 🟡 |
| T23.3 | Tag v0.1.0-alpha | 3m | 🔴 |

**Total: ~18min**

### T24: Cross-link READMEs — P4 Visibility

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T24.1 | Add templ-components link to cqrs-htmx README | 5m | 🟡 |
| T24.2 | Add cqrs-htmx link to templ-components README | 5m | 🟡 |
| T24.3 | Rewrite templ-components README with ecosystem story | 12m | 🔴 |

**Total: ~22min**

### T25: Get Listed on templ.guide — P4 Distribution

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T25.1 | Create issue/PR on templ docs repo for listing | 10m | 🟡 |
| T25.2 | Submit to awesome-templ | 8m | 🟡 |

**Total: ~18min**

### T26: Deploy Demo Site — P4 Showcase

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T26.1 | Add `cmd/demo/main.go` with HTTP server serving demo page | 10m | 🔴 |
| T26.2 | Add Fly.io/Railway deployment config | 8m | 🟡 |
| T26.3 | Deploy and verify | 5m | 🔴 |

**Total: ~23min**

---

## WAVE 6 — Personal Use Ecosystem (P5)

### T27: Unify Error Handling Across Libs — P5 Integration

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T27.1 | Design shared error type hierarchy across 3 libs | 12m | 🟡 |
| T27.2 | Map cqrs-htmx validation errors to forms.ErrorAttrs pipeline | 10m | 🟡 |
| T27.3 | Write integration test | 8m | 🟡 |

**Total: ~30min**

### T28: Reference Starter App — P5 Ecosystem

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T28.1 | Scaffold Go project with all 3 libs as deps | 10m | 🔴 |
| T28.2 | Add auth pages (login, register) using layout + forms | 12m | 🟡 |
| T28.3 | Add CRUD resource with table + pagination + forms | 12m | 🟡 |
| T28.4 | Add HTMX integration for dynamic updates | 10m | 🟡 |
| T28.5 | Wire up PostgreSQL + migrations | 10m | 🟡 |
| T28.6 | Deploy to Fly.io | 8m | 🟡 |

**Total: ~62min**

### T29: Hot Reload Dev Environment — P5 DX

| # | Sub-task | ⏱️ | Impact |
|---|----------|-----|--------|
| T29.1 | Create air/live-reload config for templ + Go + Tailwind | 10m | 🟡 |
| T29.2 | Document in README | 5m | 🟢 |

**Total: ~15min**

---

## Summary Table — All Tasks Sorted by Priority

| Wave | Task | Description | Priority | Impact | ⏱️ Total | Deps |
|:----:|:----:|-------------|:--------:|:------:|:--------:|:----:|
| 1 | T1 | Fix Modal Focus Restore (WCAG) | P0 | 🔴 | 32m | — |
| 1 | T3 | Fix Icon System Issues | P1 | 🔴 | 35m | — |
| 1 | T2 | Fix ID Propagation (6 components) | P1 | 🟡 | 45m | — |
| 1 | T4 | Add BaseProps to Breadcrumbs | P1 | 🟡 | 48m | — |
| 1 | T5 | Add BaseProps to SimpleNav | P1 | 🟡 | 48m | — |
| 1 | T6 | Add BaseProps to Spinner | P1 | 🟡 | 35m | T6.4 updates HTMX loading |
| 2 | T7 | Remove Dead Code | P2 | 🟢 | 20m | — |
| 2 | T8 | Replace Hardcoded SVGs | P2 | 🟡 | 32m | T3 (icon system must be correct first) |
| 2 | T9 | Input Validation | P2 | 🟡 | 37m | — |
| 2 | T10 | BoolString → strconv.FormatBool | P2 | 🟢 | 19m | — |
| 2 | T11 | Extract Shared Lookup Helper | P2 | 🟡 | 32m | — |
| 2 | T12 | Extract Shared Dismiss JS | P2 | 🟡 | 21m | — |
| 2 | T13 | Test Coverage Gaps | P2 | 🟡 | 58m | — |
| 2 | T14 | Remove Duplicate Test Data | P2 | 🟢 | 20m | — |
| 2 | T15 | Add BaseProps to SimpleEmptyState | P2 | 🟢 | 19m | — |
| 3 | T16 | Fix JS Re-attachment After HTMX | P2 | 🔴 | 46m | — |
| 3 | T17 | Fix Retry Counter Race | P2 | 🟡 | 22m | — |
| 4 | T18 | Update Stale TODO_LIST.md | P3 | 🟢 | 18m | — |
| 4 | T19 | Document Conventions | P3 | 🟢 | 23m | — |
| 4 | T20 | Add go doc Examples | P3 | 🟡 | 36m | — |
| 4 | T21 | Update FEATURES.md | P3 | 🟢 | 21m | — |
| 4 | T22 | Update AGENTS.md | P3 | 🟢 | 13m | — |
| 5 | T23 | Tag v0.1.0-alpha | P4 | 🔴 | 18m | T1-T17 complete |
| 5 | T24 | Cross-link READMEs | P4 | 🟡 | 22m | — |
| 5 | T25 | Get Listed on templ.guide | P4 | 🟡 | 18m | T23 |
| 5 | T26 | Deploy Demo Site | P4 | 🔴 | 23m | T1-T17 complete |
| 6 | T27 | Unify Error Handling Across Libs | P5 | 🟡 | 30m | — |
| 6 | T28 | Reference Starter App | P5 | 🔴 | 62m | T27 |
| 6 | T29 | Hot Reload Dev Environment | P5 | 🟡 | 15m | — |

---

## Execution Order (Recommended)

```
T1  →  T3  →  T2  →  T4  →  T5  →  T6
                                    ↓
T7  →  T8  →  T9  →  T10 →  T11 →  T12 →  T13 →  T14 →  T15
                                                          ↓
                                              T16  →  T17
                                                          ↓
                                              T18 → T19 → T20 → T21 → T22
                                                                        ↓
                                                    T23 → T24 → T25 → T26
                                                                        ↓
                                                    T27 → T28 → T29
```

**Wave 1 can be parallelized** (T1, T2, T3 have no deps on each other).
**Wave 2 can be parallelized** (T7, T9, T10, T11, T12, T14, T15 have no cross-deps). T8 depends on T3.
**Wave 3** is independent JS work.
**Wave 4** is documentation (can be done anytime).
**Wave 5** requires T1-T17 to be complete.
**Wave 6** is independent ecosystem work.

---

## Total Estimates

| Wave | Tasks | Total Time | Avg per sub-task |
|------|:-----:|:----------:|:----------------:|
| 1 — Critical | T1-T6 | ~3.7h | ~8m |
| 2 — Quality | T7-T15 | ~4.3h | ~7m |
| 3 — JS Fixes | T16-T17 | ~1.1h | ~10m |
| 4 — Docs | T18-T22 | ~1.8h | ~6m |
| 5 — Ecosystem | T23-T26 | ~1.3h | ~7m |
| 6 — Personal | T27-T29 | ~1.8h | ~10m |
| **Total** | **29 tasks, 87 sub-tasks** | **~14h** | **~8m** |
