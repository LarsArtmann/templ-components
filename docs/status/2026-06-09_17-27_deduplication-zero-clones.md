# Status Report — templ-components

**Date:** 2026-06-09 17:27 | **Session:** 5 (deduplication cleanup) | **Branch:** master

---

## Executive Summary

10/10 clone groups eliminated at threshold 45 (art-dupl semantic). 144 lines of duplicate test code removed. 0 regressions. Full green: 12 packages pass, 1040 tests, 0 lint issues, 0 clone groups.

---

## a) FULLY DONE ✅

### This Session (2026-06-09)

| # | Task | Files Changed | Lines Removed |
|---|------|--------------|---------------|
| 1 | navigation/coverage_test.go — merged 6 pagination sub-tests into 3 table-driven tests (ellipsis, rel, boundary) | `navigation/coverage_test.go` | -80 |
| 2 | navigation/pagination_test.go — absorbed mobile/page-info assertions | `navigation/pagination_test.go` | +8 |
| 3 | forms/edge_cases_test.go — removed duplicate textarea placeholder test | `forms/edge_cases_test.go` | -10 |
| 4 | errorpage/handler_test.go — merged ParseFamily valid+case-insensitive into single table-driven test | `errorpage/handler_test.go` | -16 |
| 5 | errorpage/a11y+bdd — merged dismissible ErrorAlert tests, removed duplicate | `errorpage/a11y_test.go`, `errorpage/bdd_test.go` | -11 |
| 6 | display/coverage_test.go — removed duplicate external link dropdown test | `display/coverage_test.go` | -13 |
| 7 | htmx/a11y_test.go — merged 5 sub-tests into 2 (nonce + combined defaults) | `htmx/a11y_test.go` | -43 |
| 8 | htmx/bdd_test.go — removed 3 duplicate tests (nonce, events, SwapOOB) | `htmx/bdd_test.go` | -30 |

**Net result:** 9 files, 75 insertions, 219 deletions. -144 net lines. 0 → 10 → 0 clone groups.

### Project Overall (All Sessions)

| Metric | Value |
|--------|-------|
| Packages | 10 + demo |
| `.templ` files | 43 (33 source + 10 shared/internal) |
| `*_templ.go` (generated, committed) | 43 |
| Handwritten `.go` files | 25 |
| Test files | 64 |
| Total test cases | ~1040 |
| Test coverage (avg) | 72–79% across all packages |
| Icons | 99 (98 path + 1 Spinner) |
| Components | 56 templ components |
| Typed enums | 18 |
| Lint issues | 0 |
| Clone groups (t=45) | 0 |
| CI | GitHub Actions (build+test+lint+coverage) |

### Feature Completion Matrix

| Package | Components | Status |
|---------|-----------|--------|
| `display` | 14 | Accordion, Avatar, Badge, StatusBadge, Card, SimpleCard, StatCard, Dropdown, EmptyState, SimpleEmptyState, Modal, Table, Tabs, Tooltip, Drawer |
| `errorpage` | 3 | ErrorPage, ErrorDetail, ErrorAlert + HTTP handler + go-error-family bridge |
| `feedback` | 12 | Alert, InlineError, InlineSuccess, Spinner, LoadingOverlay, InlineLoading, Skeleton (7 variants), ProgressBar, StepIndicator, ToastContainer, Toast |
| `forms` | 6 | Input, InputGroup, Select, Textarea, Checkbox, Label, Radio, Toggle, FileInput, Form, ValidationSummary |
| `htmx` | 7 | LoadingIndicator, InlineLoadingOverlay, LoadingButton, ConfirmDelete, CSRFToken, SwapOOB, GlobalErrorHandling |
| `icons` | 1 (99 icons) | Icon, FillIcon, IconWithStrokeWidth, IconPathJS |
| `layout` | 4 | Base, Minimal, Page, ThemeToggle, ThemeScript |
| `navigation` | 9 | Nav, SimpleNav, NavLink, MobileNavLink, Breadcrumbs, Pagination, MobileMenu, StepIndicator |

---

## b) PARTIALLY DONE 🟡

| Task | Status | What's Left |
|------|--------|-------------|
| Coverage improvement | 70–79% | fillIcon, Select, Textarea below 70% |
| Snapshot → golden migration | feedback/ done | Other packages still use raw string comparison |
| Documentation site | Not started | pkgsite/doc2go decision needed |

---

## c) NOT STARTED ⬜

### Breaking Changes (defer to v1.0)

- Move test helpers to `internal/testutil/`
- SimpleNav BaseProps conversion
- Add BaseProps to StepIndicatorProps
- Pagination uint fields

### New Components

- Date Picker
- Combobox/Autocomplete

### Infrastructure

- Verify `go get` from clean project
- Set up goreleaser for tag-based releases
- Modularize into Go workspace (10-module `go.work`)
- Consider `go:generate stringer` for enums
- Consider `Validate() error` method on props structs
- Extract shared Tailwind preset/theme configuration file
- Plan v1.0 API freeze scope and timeline

### Testing

- Consistent nonce propagation audit
- Add accessibility audit automation (axe-core/pa11y)
- Investigate gopls QF1003 suppression for generated `*_templ.go` files

### Release & Discovery

- Tag v0.3.0
- Submit to awesome-templ
- Open PR on templ.guide
- Cross-link ecosystem in README

### Housekeeping

- Consolidate inline JS into shared init strategy (10 script blocks across 7 files)
- Add `uint` type for Pagination fields

---

## d) TOTALLY FUCKED UP 💥

Nothing. Zero regressions this session. All 1040 tests pass, 0 lint issues, 0 clone groups.

---

## e) WHAT WE SHOULD IMPROVE 🔧

### Critical Quality Improvements

1. **JS consolidation** — 10 inline `<script>` blocks across 7 `.templ` files. Each has its own singleton guard. Should extract to a shared `tc-init.js` module pattern.
2. **Coverage gaps** — 3 functions below 70% (fillIcon, Select, Textarea). Quick wins.
3. **Snapshot test consistency** — feedback/ uses golden files, other packages use raw `AssertContains`. Should migrate all to `internal/golden`.
4. **Test helper location** — `utils/test_helpers.go` is exported, but only used internally. Breaking change to move it, but worth it for v1.0.
5. **Nonce audit** — Not all components consistently propagate Nonce to inline scripts. Systematic audit needed.
6. **gopls QF1003** — Generated `*_templ.go` files trigger gopls diagnostics. Need suppression strategy.

### Architecture Improvements

7. **Props validation** — Consider `Validate() error` on props structs instead of panic-on-invalid-pattern. More Go-idiomatic.
8. **Stringer for enums** — `go:generate stringer` for the 18 typed enums instead of manual string constants.
9. **Tailwind preset** — Extract theme configuration into importable preset file so consumers don't have to copy-paste CSS.
10. **Go workspace** — Modularize into sub-modules for independent versioning. But the import graph is well-structured already.

### DX & Release

11. **goreleaser** — Need automated tag-based releases. Currently manual.
12. **`go get` verification** — Need to verify clean-project consumption works (the whole reason `*_templ.go` is committed).
13. **Ecosystem presence** — Not listed on awesome-templ or templ.guide yet. Low effort, high visibility.
14. **Documentation site** — pkgsite or doc2go for browseable API docs.

---

## f) TOP #25 THINGS TO DO NEXT

### Priority 1 — Ship Quality (1–2 sessions)

| # | Task | Effort | Impact |
|---|------|--------|--------|
| 1 | Verify `go get` from clean project works | 1h | Critical — blocks v0.3.0 |
| 2 | Tag v0.3.0 release + update CHANGELOG | 30min | Ships Drawer, ValidationSummary, 25 icons, Spinner BaseProps |
| 3 | Fix coverage gaps: fillIcon, Select, Textarea | 2h | Gets all packages above 70% |
| 4 | Consolidate inline JS into shared init strategy | 4h | Reduces 10 script blocks to 1–2 |
| 5 | Nonce propagation audit across all components | 2h | CSP compliance |

### Priority 2 — Release Readiness (2–3 sessions)

| # | Task | Effort | Impact |
|---|------|--------|--------|
| 6 | Set up goreleaser for tag-based releases | 2h | Automates publishing |
| 7 | Submit to awesome-templ | 30min | Discoverability |
| 8 | Open PR on templ.guide | 30min | Discoverability |
| 9 | Cross-link ecosystem in README (cqrs-htmx, go-cqrs-lite) | 30min | GOTH stack story |
| 10 | Migrate remaining snapshot tests to golden files | 3h | Test consistency |

### Priority 3 — Polish (3–5 sessions)

| # | Task | Effort | Impact |
|---|------|--------|--------|
| 11 | `go:generate stringer` for all 18 enums | 2h | Eliminates manual string constants |
| 12 | `Validate() error` on props structs | 4h | More idiomatic than panic |
| 13 | Extract shared Tailwind preset/theme file | 2h | Consumer DX |
| 14 | Move test helpers to `internal/testutil/` | 1h | API cleanliness (breaking) |
| 15 | Investigate gopls QF1003 suppression | 1h | IDE experience |

### Priority 4 — New Features (5+ sessions)

| # | Task | Effort | Impact |
|---|------|--------|--------|
| 16 | Date Picker component | 8h | Common need |
| 17 | Combobox/Autocomplete component | 8h | Common need |
| 18 | SimpleNav BaseProps conversion | 2h | API consistency |
| 19 | Add BaseProps to StepIndicatorProps | 1h | API consistency |
| 20 | Pagination uint fields | 1h | Type safety |

### Priority 5 — Future

| # | Task | Effort | Impact |
|---|------|--------|--------|
| 21 | Modularize into Go workspace | 8h | Independent versioning |
| 22 | Documentation site (pkgsite/doc2go) | 4h | API discoverability |
| 23 | Accessibility audit automation (axe-core/pa11y) | 4h | A11y compliance |
| 24 | Plan v1.0 API freeze scope and timeline | 2h | Long-term stability |
| 25 | Add `Drawer` variant: full-screen overlay mode | 2h | Common mobile pattern |

---

## g) TOP #1 QUESTION I CANNOT FIGURE OUT MYSELF 🤔

**Should v0.3.0 include the JS consolidation (#4) or ship first and consolidate in v0.4.0?**

The 10 inline `<script>` blocks work correctly but are a maintenance burden. Consolidating them is a significant refactor that touches 7 `.templ` files and changes the JS attachment pattern. Shipping v0.3.0 now gets Drawer, ValidationSummary, 25 new icons, and Spinner BaseProps into consumers' hands sooner. But if we consolidate JS in v0.4.0, that's two breaking JS-pattern changes in rapid succession.

---

## Project Health Dashboard

| Metric | Value | Trend |
|--------|-------|-------|
| Clone groups (t=45) | **0** | 10 → 0 ✅ |
| Lint issues | **0** | Stable ✅ |
| Test pass rate | **100%** (1040/1040) | Stable ✅ |
| Coverage (avg) | **73.2%** | Stable |
| Test files | 64 | Stable |
| Source files | 43 `.templ` + 25 `.go` | Stable |
| Generated (committed) | 43 `*_templ.go` | Stable |
| Lines of code | ~4,600 `.templ` + ~12,500 `.go` | Stable |
| Icons | 99 | Stable |
| Dependencies | 0 framework deps (templ + tailwind-merge-go only) | Stable ✅ |

---

## Files Changed This Session

```
 display/coverage_test.go      |  13 ----
 errorpage/a11y_test.go        |  12 ----
 errorpage/bdd_test.go         |   2 +-
 errorpage/handler_test.go     |  27 ++------
 forms/edge_cases_test.go      |  10 ---
 htmx/a11y_test.go             |  43 +-----------
 htmx/bdd_test.go              |  30 ---------
 navigation/coverage_test.go   | 147 +++++++++++++++++-------------------------
 navigation/pagination_test.go |  10 ++-
 9 files changed, 75 insertions(+), 219 deletions(-)
```
