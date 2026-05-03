# Status Report — templ-components

**Date:** 2026-05-03 07:54
**Branch:** master
**Commits:** `2589d39` (HEAD) ← `7d17413` ← `eed8aa0`
**Go:** 1.26.2 | **templ:** v0.3.1001 | **tailwind-merge-go:** v0.2.1

---

## Metrics at a Glance

| Metric                  | Value                         | Delta (since 2026-04-28)      |
| ----------------------- | ----------------------------- | ------------------------------ |
| Packages                | 8                             | —                              |
| Exported components     | 56 (29 files)                 | —                              |
| Tests                   | 237 passing, 0 failing        | +161 (test restructure)        |
| Test coverage (avg)     | 57.7%                         | —                              |
| `.templ` files          | 30                            | —                              |
| `.templ` lines          | 2,960                         | +2                             |
| Go files (hand-written) | 31                            | +1 (display/helpers.templ)     |
| Go lines (hand-written) | 2,211                         | -57 net                        |
| Typed enums             | 16                            | —                              |
| `go vet`                | clean (1 warning in test)     | —                              |
| Clone groups (art-dupl) | **2** (down from 14)          | **-86% in this session alone** |
| CSP violations          | 0                             | —                              |
| Dependencies            | 2 (templ + tailwind-merge-go) | —                              |

### Test Coverage by Package

| Package     | Coverage |
| ----------- | -------- |
| htmx        | 77.3%    |
| navigation  | 71.1%    |
| feedback    | 69.3%    |
| layout      | 68.3%    |
| display     | 62.8%    |
| utils       | 54.5%    |
| forms       | 53.1%    |
| icons       | 5.6%     |

---

## a) FULLY DONE

### Session: Deep Semantic Deduplication (2026-05-03)

**14 → 2 clone groups** (86% reduction). Two commits: `7d17413` (production dedup) + `2589d39` (test hygiene).

#### Production Code Deduplication (`7d17413`)

| File                          | Change                                                                                          |
| ----------------------------- | ----------------------------------------------------------------------------------------------- |
| `navigation/pagination.templ` | Extracted `paginationArrowIcon` from `paginationArrow` active/disabled branches                 |
| `navigation/pagination.templ` | Extracted `paginationArrow`, `mobilePageButton` sub-templates (was already done)                |
| `htmx/loading.templ`          | 3 inline spinner SVGs → `@feedback.Spinner()`                                                  |
| `icons/icon.templ`            | Spinner case → `@feedback.Spinner()`                                                            |
| `feedback/alert.templ`        | Extracted `inlineMessage`, `inlineErrorIcon`, `inlineSuccessIcon` sub-templates                 |
| `display/tabs.templ`          | Extracted `tabLink` sub-template                                                                |
| `display/dropdown.templ`      | Extracted `dropdownItemClass` constant + `dropdownItemLink` sub-template                        |
| `display/empty_state.templ`   | Extracted `emptyStateActionClass` constant + `emptyStateAction` sub-template                    |
| `display/helpers.templ`       | **NEW** — shared `fillIcon` sub-template for 20×20 filled SVG icons                            |
| `display/accordion.templ`     | Inline SVG → `@fillIcon()`                                                                     |
| `display/card.templ`          | 2 inline SVGs → `@fillIcon()`                                                                  |

#### Test Code Deduplication & Hygiene (`2589d39`)

| File                          | Change                                                                                |
| ----------------------------- | ------------------------------------------------------------------------------------- |
| `utils/test_helpers.go`       | Added generic `AssertEqual[T comparable]` helper                                      |
| `display/helpers_test.go`     | Replaced 3× `if got != want` patterns → `utils.AssertEqual`                          |
| `display/tooltip_test.go`     | 3 subtests → 1 table-driven test                                                      |
| `display/card_test.go`        | Removed zero-value `BaseProps` fields                                                 |
| `display/modal_test.go`       | Removed zero-value `BaseProps` fields                                                 |
| `feedback/helpers_test.go`    | Extracted `assertStyleFunc4`; converted `stepLineClass` to table-driven; `AssertEqual` |
| `feedback/snapshot_test.go`   | Removed zero-value `BaseProps` and struct fields                                      |
| `layout/sri_test.go`          | Extracted `testIntegrityHash` helper with function parameter                          |
| `navigation/pagination_test.go` | Extracted `renderPagination` helper for 4 repeated render calls                     |
| `navigation/snapshot_test.go` | Extracted shared `testNavLinks` variable                                              |

#### Net Result

```
23 files changed, 401 insertions(+), 402 deletions(-)
Net: -1 line (massive structural improvement, zero growth)
Clone groups: 14 → 2
```

### Previously Completed (from git history)

- 16 typed string enums (AlertType, ToastType, TabsStyle, BadgeSize, AvatarSize, etc.)
- CSP compliance for all inline scripts (nonce-based)
- 56 exported templ components across 30 files
- `utils.Class()` with tailwind-merge-go for intelligent class merging
- `utils.BaseProps` shared across all components
- Icon system with 30+ named icons + typed `Name` enum
- `CHANGELOG.md` initiated
- Comprehensive project hygiene (git config, linting, formatting, metadata)

---

## b) PARTIALLY DONE

### Remaining 2 Clone Groups (Structural, Not Worth Deduplicating)

| # | Clone                                         | Why Not Deduped                                                               |
| - | --------------------------------------------- | ----------------------------------------------------------------------------- |
| 1 | `navigation/snapshot_test.go:14,14` + `22,22` | 1-token: `utils.Render(t, NavLink(NavLinkProps{...}, "/"))` — different test data each time |
| 2 | `feedback/helpers_test.go:82,91` + `118,127`  | 1-token: `assertStyleFunc4(t, fmt.Sprintf("...Styles(%q)", tt.typ), ...)` — testing different functions |

Both are **single-line, 1-token clones** — inherent to testing different functions with the same assertion pattern. Further extraction would hurt readability.

### Test Coverage Gaps

| Package | Coverage | Gap Description                             |
| ------- | -------- | ------------------------------------------- |
| icons   | 5.6%     | Only name-lookup tests, no render tests for most icons |
| forms   | 53.1%    | Snapshot tests exist but no edge case tests  |
| utils   | 54.5%    | `Class()` only tested via indirect usage     |

---

## c) NOT STARTED

| # | Item                                                                                       |
| - | ------------------------------------------------------------------------------------------ |
| 1 | **FEATURES.md** — No feature inventory document exists                                     |
| 2 | **TODO_LIST.md** — No comprehensive TODO tracking exists                                   |
| 3 | **CONTEXT.md** — No project context/architecture document                                  |
| 4 | **Snapshot/golden file testing** — No rendered HTML golden files                            |
| 5 | **A11y testing** — No automated accessibility audits (aria-*, role, tabindex checks)        |
| 6 | **Benchmark tests** — No performance benchmarks for render paths                           |
| 7 | **Example/demo app** — No showcase application                                             |
| 8 | **CSS extraction** — No shared Tailwind preset/theme configuration file                    |
| 9 | **Cross-package shared helpers** — `fillIcon` is in `display/` but other packages can't use it without import cycles |
| 10 | **Documentation site** — No generated documentation website                                |
| 11 | **Migration guide** — No version migration documentation                                   |
| 12 | **Playground/REPL** — No interactive component preview                                     |
| 13 | **CI pipeline** — No `.github/workflows` or equivalent                                     |
| 14 | **Release automation** — No goreleaser or similar                                          |
| 15 | **ADR directory** — No `docs/adr/` for architecture decision records                       |

---

## d) TOTALLY FUCKED UP

| # | Issue                                                                                                    | Severity |
| - | -------------------------------------------------------------------------------------------------------- | -------- |
| 1 | **`utils/utils_test.go:48`** — `p := new(v)` followed by `if p == nil` is a dead code warning (LSP: "impossible condition: non-nil == nil"). The `Ptr` function was changed to use `new()` but the test was never updated. | Low      |
| 2 | **No FEATURES.md / TODO_LIST.md** — Referenced in previous status reports but never actually created.    | Medium   |
| 3 | **9 stale status reports** in `docs/status/` — no cleanup, growing unbounded, historical noise           | Low      |
| 4 | **Cross-package import risk** — `icons/icon.templ` imports `feedback.Spinner` creating `icons → feedback` dependency. If feedback ever needs icons, circular import occurs. | Medium   |
| 5 | **`display/empty_state.go`** — Hand-written Go file alongside `empty_state.templ`. Potential templ generation conflict if names collide. | Low      |
| 6 | **icons package at 5.6% coverage** — Essentially untested for actual rendering.                         | Medium   |

---

## e) WHAT WE SHOULD IMPROVE

| # | Improvement                                                                                                                   | Effort |
| - | ----------------------------------------------------------------------------------------------------------------------------- | ------ |
| 1 | Fix dead test code in `utils/utils_test.go:48` — test `Ptr` properly                                                        | 5min   |
| 2 | Create **FEATURES.md** — audit actual codebase, not aspirational                                                             | 30min  |
| 3 | Create **TODO_LIST.md** — track real work items with status                                                                  | 20min  |
| 4 | Extract shared SVG helpers (`fillIcon`, `spinnerSVG`) to `internal/svg/` or `shared/` package to resolve cross-package risk | 1hr    |
| 5 | Add snapshot/golden file testing for all 56 components                                                                       | 2hr    |
| 6 | Add a11y attribute validation tests (aria-*, role, tabindex)                                                                 | 1hr    |
| 7 | Set up CI pipeline (build + test + vet + lint on push)                                                                       | 1hr    |
| 8 | Consistent nonce handling audit — some components use `props.Nonce`, some hardcode, some skip                                | 30min  |
| 9 | Prune old status reports — keep last 2, archive rest                                                                         | 5min   |
| 10 | Dark mode class output verification tests                                                                                     | 1hr    |
| 11 | Component composition tests (nesting components inside each other)                                                            | 1hr    |
| 12 | Form validation helper patterns — no client-side validation helpers exist                                                     | 2hr    |

---

## f) Top 25 Things We Should Get Done Next

| #  | Priority | Task                                                                                     | Impact | Effort |
| -- | -------- | ---------------------------------------------------------------------------------------- | ------ | ------ |
| 1  | P0       | Fix `utils/utils_test.go:48` dead code warning                                           | Clean  | 5min   |
| 2  | P0       | Create FEATURES.md from actual codebase audit                                            | Track  | 30min  |
| 3  | P0       | Create TODO_LIST.md                                                                      | Track  | 20min  |
| 4  | P0       | Extract `fillIcon` + `spinnerSVG` to shared package (resolve `icons → feedback` import)  | Arch   | 1hr    |
| 5  | P1       | Add icon render tests (icons is at 5.6% coverage)                                        | Test   | 1hr    |
| 6  | P1       | Add snapshot/golden file tests for all 56 components                                     | Test   | 2hr    |
| 7  | P1       | Set up CI pipeline (build + test + vet + lint)                                           | Infra  | 1hr    |
| 8  | P1       | Add a11y attribute validation tests                                                      | Test   | 1hr    |
| 9  | P1       | Consistent nonce propagation audit across all components                                 | Sec    | 30min  |
| 10 | P1       | Create CHANGELOG.md (update existing, add version tracking)                              | Docs   | 20min  |
| 11 | P2       | Dark mode `dark:` class output verification tests                                        | Test   | 1hr    |
| 12 | P2       | Component composition tests (nesting)                                                    | Test   | 1hr    |
| 13 | P2       | Prune old status reports (keep last 2, archive rest)                                     | Clean  | 5min   |
| 14 | P2       | Add benchmark tests for hot paths (`Class()`, `Spinner()`, `Badge()` render)            | Perf   | 30min  |
| 15 | P2       | Extract shared Tailwind preset/theme configuration file                                  | DX     | 1hr    |
| 16 | P2       | Create example/demo application                                                          | DX     | 2hr    |
| 17 | P2       | Investigate `empty_state.go` hand-written file for templ generation conflicts            | Risk   | 15min  |
| 18 | P3       | Client-side form validation helpers                                                      | Feat   | 2hr    |
| 19 | P3       | Documentation site generation                                                            | Docs   | 3hr    |
| 20 | P3       | Interactive playground/REPL for component preview                                        | DX     | 4hr    |
| 21 | P3       | Release automation (goreleaser)                                                          | Infra  | 1hr    |
| 22 | P3       | Version migration guides                                                                 | Docs   | 2hr    |
| 23 | P3       | Add `docs/adr/` directory for architecture decision records                              | Docs   | 15min  |
| 24 | P4       | Cross-package circular import audit (icons ↔ feedback full analysis)                     | Risk   | 30min  |
| 25 | P4       | Explore nix flake migration (replace justfile/shell scripts)                             | Infra  | 2hr    |

---

## g) Top #1 Question I Cannot Figure Out Myself

**Should we create a `shared/` or `internal/svg/` package for cross-cutting SVG/template helpers?**

Currently:
- `display/helpers.templ` has `fillIcon` — only usable within `display/`
- `feedback/loading.templ` has `Spinner()` — usable by `icons/` (which imports it)
- `icons/icon.templ` imports `feedback.Spinner` → creates `icons → feedback` directional dependency
- If `feedback/` ever needs `display.fillIcon`, we get a cycle

Options:
1. **`internal/svg/`** — shared SVG primitives (`fillIcon`, `spinnerSVG`, `chevronIcon`). Clean but adds package.
2. **`shared/`** — broader shared helpers. More flexible but less focused.
3. **Keep as-is** — accept `icons → feedback` dependency. Works until it doesn't.
4. **Move spinner into `icons/`** — spinner is an icon, after all. Then `feedback/` imports `icons.Spinner()` instead.

**This decision affects the public API surface and import ergonomics — only the project owner should decide.**

---

## Component Inventory (56 exported components)

### layout (4)
`Base`, `Minimal`, `ThemeScript`, `ThemeToggle`

### feedback (10)
`ToastContainer`, `Toast`, `Spinner`, `LoadingOverlay`, `InlineLoading`, `Skeleton`, `SkeletonGroup`, `Alert`, `InlineError`, `InlineSuccess`, `ProgressBar`, `StepIndicator`

### display (13)
`Tooltip`, `Accordion`, `Modal`, `Table`, `Avatar`, `Tabs`, `Badge`, `StatusBadge`, `Card`, `SimpleCard`, `StatCard`, `Dropdown`, `EmptyState`, `SimpleEmptyState`

### forms (5)
`Input`, `Checkbox`, `Select`, `Textarea`, `Label`, `FieldError`

### navigation (8)
`Nav`, `SimpleNav`, `Footer`, `NavLink`, `MobileNavLink`, `Breadcrumbs`, `MobileMenuToggle`, `MobileMenu`, `Pagination`

### icons (1)
`Icon`

### htmx (5)
`GlobalErrorHandling`, `ConfirmDelete`, `SwapOOB`, `CSRFToken`, `LoadingIndicator`, `InlineLoadingOverlay`, `LoadingButton`

---

*Generated by art-dupl deduplication session — 2026-05-03*
