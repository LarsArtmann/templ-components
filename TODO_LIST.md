# TODO List — templ-components

**Updated:** 2026-05-19

Legend: ✅ Done | 🔨 In Progress | ⬜ Not Started | ❌ Blocked

---

## Session 10 (2026-05-19) — Comprehensive 9-Skill Audit

Full audit: code quality scan, features audit, TODO list builder, architecture review, architecture improvement, architecture visualization, full code review, BDD testing, go-modularize analysis.

### Code Quality Baseline

|     | #   | Status                            | Task                            | Notes |
| --- | --- | --------------------------------- | ------------------------------- | ----- |
| 1   | ✅  | Build passes (`go build ./...`)   | Clean. Zero issues.             |
| 2   | ✅  | All tests pass (`go test ./...`)  | 9 packages, all green.          |
| 3   | ✅  | Lint passes (`golangci-lint run`) | 0 issues on library packages.   |
| 4   | ✅  | Coverage: 71.8%                   | Range: 70.5%–89.5% per package. |

---

## Critical Issues (P0)

|     | #   | Status                                 | Task      | Package                                                                                                       | Notes |
| --- | --- | -------------------------------------- | --------- | ------------------------------------------------------------------------------------------------------------- | ----- |
| 5   | ⬜  | Fix demo app to use `layout.Base`      | examples/ | Currently uses raw HTML, wrong Tailwind version (v2 vs v4), no HTMX. Anti-advertisement for the library.      |
| 6   | ⬜  | Remove or validate unknown icon names  | icons/    | `iconPaths()` silently falls back to clock icon. `Name("typo")` renders wrong icon. Error or panic instead.   |
| 7   | ⬜  | Fix `IconPathJS` stroke-width mismatch | icons/    | JS version uses `stroke-width="2"` vs templ `stroke-width="1.5"`. Split brain: same icon renders differently. |
| 8   | ⬜  | Delete deprecated `Exclamation` icon   | icons/    | Identical SVG path to `ExclamationCircle`. Dead code walking. Remove from constants, path map, test lists.    |

---

## Architecture (P1)

|     | #   | Status                                                          | Task        | Package                                                                                                                         | Notes |
| --- | --- | --------------------------------------------------------------- | ----------- | ------------------------------------------------------------------------------------------------------------------------------- | ----- |
| 9   | ⬜  | Unify `AlertType`/`ToastType` into shared type                  | feedback/   | Same 4 values (Success, Error, Warning, Info). Merge into `FeedbackLevel` or keep one type.                                     |
| 10  | ⬜  | Merge `alertStyleMap`/`toastStyleMap`                           | feedback/   | Near-identical maps. Icon color difference (400 vs 600) may be intentional — verify.                                            |
| 11  | ⬜  | Make `SimpleCard` compose through `Card`                        | display/    | Currently duplicates shell rendering. `SimpleCard` should call `Card` with no title/footer.                                     |
| 12  | ⬜  | Add `ComponentProps` interface to `utils.BaseProps`             | utils/      | `GetBaseProps() BaseProps` method enables generic component handling. 29 props structs share BaseProps but no common interface. |
| 13  | ⬜  | Use stable IDs in modal JS instead of CSS selectors             | display/    | `[role="dialog"] > div:last-child` breaks if extra div added. Use `props.ID + "-panel"`.                                        |
| 14  | ⬜  | Fix ThemeToggle multi-instance bug                              | layout/     | `tcThemeToggleAttached` prevents second toggle from working. Remove global guard or use per-ID.                                 |
| 15  | ⬜  | Use `icons.Icon`/`svg.FillIcon` in Breadcrumbs chevron          | navigation/ | Currently hardcoded raw SVG instead of using icon system.                                                                       |
| 16  | ⬜  | Add `BaseProps` to `StepIndicatorProps`                         | feedback/   | Only feedback component without BaseProps. Can't add ID/class/attrs.                                                            |
| 17  | ⬜  | Convert `LoadingOverlay` from positional params to props struct | feedback/   | Only component in feedback not using props struct. Inconsistent API.                                                            |
| 18  | ⬜  | Validate `SwapOOB` swapStyle parameter                          | htmx/       | Accepts any string. HTMX expects specific values. Typo = silent failure.                                                        |
| 19  | ⬜  | Validate `Pagination` CurrentPage > 0                           | navigation/ | `CurrentPage: 0` renders page-0 link. Should validate or clamp.                                                                 |
| 20  | ⬜  | Clamp `ProgressBar` percent to [0, 100]                         | feedback/   | `Current > Total` produces `width: 105%+`, overflows.                                                                           |

---

## JS Architecture (P1)

|     | #   | Status                                                      | Task     | Package                                                                          | Notes |
| --- | --- | ----------------------------------------------------------- | -------- | -------------------------------------------------------------------------------- | ----- |
| 21  | ⬜  | Consolidate inline JS into shared init strategy             | multi    | 222 lines across 7 files. Create `tc-init.js` pattern loaded by `layout.Base`.   |
| 22  | ⬜  | Fix HTMX swap event listener re-attachment                  | multi    | Global `tc*Attached` guards prevent re-attachment after HTMX DOM swaps.          |
| 23  | ⬜  | Fix `GlobalErrorHandling` shared retry counter              | htmx/    | `retryCount` is shared across concurrent requests. Race condition.               |
| 24  | ⬜  | Make `GlobalErrorHandling` retry/config values configurable | htmx/    | MAX_RETRIES, RETRY_DELAY_MS, MAX_ERROR_HISTORY are hardcoded.                    |
| 25  | ⬜  | Consolidate modal per-instance JS into single function      | display/ | Each modal emits identical `tcCloseModal_<id>` function. Use `tcCloseModal(id)`. |

---

## Type Safety & Code Quality (P2)

|     | #   | Status                                                       | Task          | Package                                                                                        | Notes |
| --- | --- | ------------------------------------------------------------ | ------------- | ---------------------------------------------------------------------------------------------- | ----- |
| 26  | ⬜  | Replace `DropdownItem` empty-Href discrimination             | display/      | `Href=""` = button, `Href="url"` = link. Use typed enum or separate structs.                   |
| 27  | ⬜  | Change `FillIcon` variadic `rotate ...bool` to `rotate bool` | internal/svg/ | Variadic allows nonsensical multi-bool. Every call site passes 0 or 1.                         |
| 28  | ⬜  | Audit `tailwind-merge-go` thread safety, remove mutex        | utils/        | `utils.Class()` has global `sync.Mutex`. If twmerge is stateless (pure func), delete the lock. |
| 29  | ⬜  | Replace `BoolString()` with `strconv.FormatBool`             | utils/        | Duplicates stdlib. Only one consumer (accordion).                                              |
| 30  | ⬜  | Validate `SelectOption` contradiction (Disabled+Selected)    | forms/        | Disabled+Selected is impossible state per HTML spec.                                           |
| 31  | ⬜  | Use `net/url` for pagination URL construction                | navigation/   | Current `pageURL` doesn't handle URL fragments correctly.                                      |
| 32  | ⬜  | Add `uint` type for `Pagination.CurrentPage/TotalPages`      | navigation/   | Currently `int` allows negative values.                                                        |

---

## Icon System Cleanup (P2)

|     | #   | Status                                         | Task                                    | Package                                                                                                 | Notes                  |
| --- | --- | ---------------------------------------------- | --------------------------------------- | ------------------------------------------------------------------------------------------------------- | ---------------------- | ------------------------------- |
| 33  | ⬜  | Eliminate 4-way icon list split brain          | icons/                                  | Same data in: constants, path map keys, `allIconNames()`, BDD test inline list. Auto-generate from map. |
| 34  | ⬜  | Validate `                                     | ` separator doesn't appear in SVG paths | icons/                                                                                                  | Multi-path icons use ` | ` delimiter with no validation. |
| 35  | ⬜  | Document 20×20 fill vs 24×24 stroke convention | icons/                                  | `FillIcon` and `strokeIcon` use different viewBox/fill paradigms.                                       |

---

## Test Suite Cleanup (P2)

|     | #   | Status                                                    | Task           | Package                                                                            | Notes |
| --- | --- | --------------------------------------------------------- | -------------- | ---------------------------------------------------------------------------------- | ----- |
| 36  | ⬜  | Consolidate test files: 1-2 per package                   | all            | 37 test files with 60-80% assertion overlap across bdd/snapshot/a11y.              |
| 37  | ⬜  | Remove unused `badgeTextLive` constant                    | display/       | Linter warning: `display/badge_test.go:16`.                                        |
| 38  | ⬜  | Delete `TestPtr` in `utils_test.go`                       | utils/         | Tests Go's built-in `new()`, not library code. Leftover from Ptr removal.          |
| 39  | ⬜  | Replace `splitSpace`/`splitClasses` with `strings.Fields` | utils/, icons/ | Cross-package duplication of trivial string splitter.                              |
| 40  | ⬜  | Move `BenchmarkHotPaths` out of `a11y_test.go`            | display/       | Benchmark test in wrong file. Should be in benchmark_test.go or component_test.go. |
| 41  | ⬜  | Remove duplicate test data declarations in navigation/    | navigation/    | `testNavLinks` declared 3 times in snapshot_test.go.                               |

---

## Documentation (P2)

|     | #   | Status                                                    | Task      | Package                                                                         | Notes |
| --- | --- | --------------------------------------------------------- | --------- | ------------------------------------------------------------------------------- | ----- |
| 42  | ⬜  | Update TODO #11 note: `dropdownSafeID` was removed        | TODO_LIST | The function was removed in later refactor. Note is stale.                      |
| 43  | ⬜  | Update CONTRIBUTING.md: remove `dropdownSafeID` reference | docs/     | Line 42 references removed function.                                            |
| 44  | ⬜  | Document `htmx` → `feedback` runtime JS dependency        | htmx/     | `GlobalErrorHandling` requires `ToastContainer` on page.                        |
| 45  | ⬜  | Fix Tooltip `aria-describedby` linkage                    | display/  | Tooltip has `role="tooltip"` + ID but trigger element lacks `aria-describedby`. |

---

## Deferred (P3 — Post v1.0)

|     | #   | Status                                               | Task                                                                                 | Notes |
| --- | --- | ---------------------------------------------------- | ------------------------------------------------------------------------------------ | ----- |
| 46  | 🔨  | Convert snapshot tests to golden file comparison     | Infrastructure designed. Deprioritized until v1.0 API freeze.                        |
| 47  | 🔨  | Move test helpers out of `utils/`                    | Breaking API change. Planned for v1.0.                                               |
| 48  | 🔨  | Documentation site generation                        | `pkg.go.dev` provides adequate API docs. Custom site is post-v1.0.                   |
| 49  | ⬜  | Add Radio, File input, Toggle/Switch form components | Not yet implemented.                                                                 |
| 50  | ⬜  | Add client-side JS tab switching                     | Tabs currently server-rendered only. Inconsistent with other interactive components. |
| 51  | ⬜  | Make `PageProps` zero-value safe                     | `PageProps{}` produces broken page. `DefaultPageProps()` is opt-in.                  |

---

## Completed (Sessions 1-9)

### Session 9 (2026-05-18) — String Literal Extraction

- Extracted all inline string literals to package constants across all packages
- Display: 12 constants across badge, card, accordion, dropdown, tabs, tooltip tests
- Feedback: 7 constants across alert, toast, loading, progress tests
- Forms: 5 constants across input, helpers tests
- Icons: 2 constants in snapshot tests
- Navigation: 4 constants across nav, pagination tests
- Lint: 0 issues, all tests pass

### Session 3 (2026-05-07)

- `utils.Class()` replaces comma-join in forms, Badge, StatCard, Nav
- 17 `DefaultXxxProps()` constructors added across all packages
- Type-safe icons: `Icon string` → `icons.Name` (breaking change)
- Modal a11y: focus trap, Escape key, focus management
- Dropdown a11y: arrow key navigation, Escape to close
- Tabs a11y: proper ARIA linkage
- Tooltip a11y: deterministic id for aria-describedby
- Package doc comments for all packages
- 14 commits, 69.7% avg test coverage

### Session 2 (2026-05-04)

- Unified `feedbackStyleSet` + `lookupFeedbackStyle[T]()`
- Deepened icon rendering: 187-line switch → `iconPathData` map
- Added `AvatarStatus` enum, `TrendDirection` enum
- Fixed `HTMXSRI` → `HTMXUseSRI bool`
- Fixed ProgressBar integer division
- Added `Content templ.Component` to `TableCell`
- Implemented `TableProps.Bordered`
- Added `utils.BoolString()`
- Extensive a11y, dark mode, benchmark, XSS tests

### Session 1 (2026-05-03)

- Fixed 4 critical bugs (literal string rendering, version mismatch, stale docs)
- Semantic deduplication (13→7 clone groups)
- `layout.BaseProps` → `PageProps` rename
- Map-based style lookups
- Created all initial documentation
- Added BDD tests for display, feedback, forms packages
