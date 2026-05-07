# Status Report: Code Deduplication Session

**Date:** 2026-05-07 07:34 CEST  
**Branch:** master  
**Commit Base:** d379618 (docs,config,components,examples: comprehensive v0.2 session cleanup)

---

## Executive Summary

Systematic deduplication pass across all production `.templ` files and test files. Net result: **-174 lines deleted**, 0 tests lost, 0 regressions, all green.

| Metric                   | Before | After     | Delta        |
| ------------------------ | ------ | --------- | ------------ |
| Go clone groups (t≥15)   | 32     | 21        | -11          |
| Templ clone groups (t≥5) | 12     | 11        | -1           |
| Total clone groups       | **44** | **32**    | **-12**      |
| Lines changed            | —      | 405+/579- | **-174 net** |
| Tests passing            | 10/10  | 10/10     | 0            |
| Lint issues              | 0      | 0         | 0            |

Note: 14 of the 21 remaining Go clone groups involve new untracked `bdd_test.go` files from session 3 that weren't in the original scan. The original tracked files went from 32→7 meaningful groups.

---

## A) FULLY DONE

### Production Code (templ)

1. **Extract `helpText` shared template** — `forms/label.templ`
   - Created `templ helpText(id, text string)` private helper
   - Replaced 3 identical inline `<p class="mt-1 text-sm...">` blocks in `input.templ`, `select.templ`, `textarea.templ`
   - Eliminates forms clone group: `forms/input.templ:98,100` + `forms/select.templ:77,79` + `forms/textarea.templ:75,77`

2. **Replace `paginationArrowIcon` with `svg.FillIcon`** — `navigation/pagination.templ`
   - Added `internal/svg` import
   - Replaced inline `<svg viewBox="0 0 20 20" fill="currentColor">` with `@svg.FillIcon("h-5 w-5", svgPath)`
   - Reused existing package instead of duplicating SVG wrapper

3. **Extract `menuIconSVG` helper** — `navigation/mobile_menu.templ`
   - Created `templ menuIconSVG(class string, pathData string)` private helper
   - Replaced 2 inline `<svg fill="none" viewBox="0 0 24 24" stroke-width="1.5"...>` blocks
   - Eliminates `mobile_menu.templ:17,20` ↔ `mobile_menu.templ:20,23` clone

4. **Simplify `examples/demo/main.go`**
   - Extracted `alert(title, message, alertType)` helper function
   - Removed 5 zero-value `BaseProps: utils.BaseProps{}` fields
   - Removed 4 `//nolint:exhaustruct` comments (moved to helper)
   - Removed unused `utils` import
   - Eliminated `examples/demo/main.go:37,39` ↔ `41,43` ↔ `65,67` ↔ `108,110` ↔ `140,142` clone group

### Test Code

5. **Merge duplicate nav dark mode tests** — `navigation/a11y_test.go`
   - `TestNavDarkMode`: 2 sub-tests → 1 combined test (both assertions in single render)

6. **Merge duplicate breadcrumb dark mode tests** — `navigation/a11y_test.go`
   - `TestBreadcrumbsDarkMode`: 2 sub-tests → 1 combined test (active + link in single render)

7. **Remove duplicate table tests** — `display/`
   - Removed caption test from `a11y_test.go` (kept in `table_test.go` with added `<caption` check)
   - Removed headers+data+caption tests from `bdd_test.go` (kept striped test)

8. **Remove duplicate dropdown tests** — `display/bdd_test.go`
   - Removed `TestDropdownUserCanSelectAction` entirely (covered by `dropdown_test.go`)

9. **Remove duplicate tabs tests** — `display/bdd_test.go`
   - Removed `TestTabsUserCanSwitchViews` entirely (covered by `tabs_test.go`)

10. **Remove duplicate avatar tests** — `display/`
    - Removed `TestAvatarUserCanIdentifyUsers` from `bdd_test.go` (covered by `avatar_test.go`)
    - Removed `TestCompositionAvatarWithStatus` from `composition_test.go` (covered by `avatar_test.go`)

11. **Remove duplicate feedback tests** — `feedback/bdd_test.go`
    - Removed ToastContainer test (covered by `snapshot_test.go`)
    - Removed InlineLoading test (exact duplicate of `snapshot_test.go:90,95`)
    - Removed StepIndicator test (subset of `snapshot_test.go:137,147`)

12. **Remove duplicate HTMX test** — `htmx/snapshot_test.go`
    - Removed `TestGlobalErrorHandlingRender` (covered by `a11y_test.go`)

13. **Extract shared `allIconNames` list** — `icons/`
    - Added `var allIconNames` in `icon_names_test.go`
    - Replaced inline icon list in `TestAllIconsRender` (snapshot_test.go)
    - Replaced inline icon list in `TestIconCount` (icon_names_test.go)

14. **Convert to table-driven tests**
    - `display/avatar_test.go`: avatar status tests (online/offline) → loop
    - `icons/snapshot_test.go`: IconAttrs tests → loop
    - `layout/integration_test.go`: HTML document assertion checks → loop
    - `utils/utils_test.go`: MapEnum fallback tests → loop
    - `forms/bdd_test.go`: FieldError tests → loop

---

## B) PARTIALLY DONE

1. **Templ cross-package clones remain** — 11 groups still detected:
   - `modal.templ:36,39` ↔ `theme.templ:28,31` ↔ `theme.templ:31,34` (close-X SVG icon in modal vs sun/moon SVG in theme toggle)
   - `loading.templ:46,56` ↔ `htmx/helpers.templ:6,17` (LoadingOverlay vs ConfirmDelete structure)
   - `alert.templ:55,58` ↔ `toast.templ:148,151` ↔ `icon.templ:11,14` (SVG wrapper pattern)
   - These are cross-package structural similarities that would require shared packages or would harm readability if extracted

2. **New untracked `bdd_test.go` files** (from session 3) — `htmx/bdd_test.go`, `icons/bdd_test.go`, `layout/bdd_test.go`, `navigation/bdd_test.go`
   - These were NOT in the original `art-dupl` scan but are now generating 14 clone groups
   - They duplicate existing snapshot/a11y tests heavily
   - Should be either committed and deduplicated or removed

---

## C) NOT STARTED

1. **Cross-package shared SVG icon package** — A shared `internal/svg` extension for close icons, chevron icons etc. used by modal, theme toggle, breadcrumbs, pagination, toast. Would eliminate ~5 templ clone groups but adds import complexity.

2. **Shared dismiss script pattern** — Alert dismiss, toast dismiss, and toast container all use similar `if (!window.tcXxxAttached) { ... }` patterns. Could extract to a shared JS helper but would require nonce propagation refactoring.

3. **Deep dedup of new bdd_test.go files** — The 4 untracked BDD test files heavily overlap with snapshot and a11y tests. Need to decide: keep BDD tests as behavioral layer (value: user-facing language) and remove overlapping snapshot tests, or vice versa.

4. **`feedback/helpers_test.go` duplicates** — `assertStyleFunc4` is called identically in `TestToastStyles` and `TestAlertStyles` with the same structure. Could merge into a generic `testStyleLookup` helper.

---

## D) TOTALLY FUCKED UP

Nothing. All changes verified:

- `go build ./...` ✅
- `go test ./...` — 10/10 packages pass ✅
- `golangci-lint run` — 0 issues ✅
- `templ generate` — 0 updates needed ✅

One pre-existing bug was found and fixed: `navigation/bdd_test.go:164-171` tested `TotalPages=1` with `AssertContains("<nav")` but `Pagination` intentionally skips rendering when `TotalPages <= 1`. Fixed to `AssertNotContains`.

---

## E) WHAT WE SHOULD IMPROVE

1. **Test file proliferation** — 4 packages now have 4+ test files each (snapshot, a11y, bdd, composition). The BDD tests add behavioral language but create massive overlap. Need a clear testing strategy: either BDD replaces snapshot tests, or BDD tests only test unique behaviors not covered elsewhere.

2. **Demo file maintenance** — `examples/demo/main.go` is verbose with raw `w.Write` HTML. Should use `layout.Base` and actual components instead of hand-written HTML strings. Would eliminate the remaining 4 clone groups in that file.

3. **Clone detection in CI** — `art-dupl` should run in CI to prevent regressions. A threshold of ≤15 Go + ≤12 templ groups seems reasonable.

4. **Shared test helpers** — Many test files repeat `utils.Render(t, Component)` + `utils.AssertContains`. A `renderAndAssert(t, component, contains...)` helper would cut test boilerplate.

5. **Icon test dedup** — `TestIconNames` in `icon_names_test.go` duplicates the icon list even with `allIconNames` extracted. The table-driven struct `{name, value}` is redundant since `Name` is a string alias — just iterate `allIconNames`.

---

## F) Top #25 Things We Should Get Done Next

### High Impact (architecture & quality)

| #   | Task                                                                         | Effort | Impact |
| --- | ---------------------------------------------------------------------------- | ------ | ------ |
| 1   | Decide BDD test strategy: keep/slim/remove the 4 untracked bdd_test.go files | 2h     | High   |
| 2   | Rewrite `examples/demo/main.go` to use `layout.Base` + actual components     | 1h     | Medium |
| 3   | Add `art-dupl` threshold check to CI/lint workflow                           | 30min  | High   |
| 4   | Extract shared close-X SVG to `internal/svg` (used by modal, toast, alert)   | 30min  | Medium |
| 5   | Merge `feedback/helpers_test.go` toast/alert style tests into generic helper | 20min  | Low    |
| 6   | Add godoc to all exported types/functions still missing it                   | 2h     | High   |
| 7   | Complete the `templ QF1003` tagged switch fix in `card.templ:406`            | 10min  | Low    |
| 8   | Add `testify` or keep custom assertions? Decide and be consistent            | 30min  | Medium |
| 9   | Snapshot testing: add golden file comparison for all components              | 3h     | High   |
| 10  | Dark mode test coverage audit — ensure every component has a dark mode test  | 1h     | Medium |

### Medium Impact (developer experience)

| #   | Task                                                                             | Effort | Impact |
| --- | -------------------------------------------------------------------------------- | ------ | ------ |
| 11  | Extract `renderAndAssert(t, component, contains...)` test helper                 | 20min  | Medium |
| 12  | Simplify `TestIconNames` to just iterate `allIconNames`                          | 10min  | Low    |
| 13  | Add example tests (`func ExampleXxx()`) for each exported component              | 2h     | Medium |
| 14  | Create `CONTRIBUTING.md` with coding conventions from AGENTS.md                  | 1h     | Medium |
| 15  | Add `go generate` based benchmarks for hot render paths                          | 1h     | Medium |
| 16  | Audit all components for missing `aria-*` attributes                             | 2h     | High   |
| 17  | Add keyboard navigation tests for interactive components (dropdown, modal, tabs) | 2h     | Medium |
| 18  | Verify all components pass WCAG 2.1 AA automated checks                          | 3h     | High   |

### Lower Impact (polish & cleanup)

| #   | Task                                                              | Effort | Impact |
| --- | ----------------------------------------------------------------- | ------ | ------ |
| 19  | Consolidate `navigation/snapshot_test.go` single-line test clones | 15min  | Low    |
| 20  | Add integration test for full page render with all components     | 1h     | Medium |
| 21  | Remove unused `props` variable in `examples/demo/main.go:25-27`   | 5min   | Low    |
| 22  | Add version badge and CI status to README                         | 30min  | Medium |
| 23  | Migrate `examples/demo` to use `layout.Base` properly             | 1h     | Medium |
| 24  | Add changelog generation from conventional commits                | 1h     | Low    |
| 25  | Performance: profile render allocations for top 5 components      | 2h     | Medium |

---

## G) Top #1 Question I Cannot Figure Out Myself

**What is the intended relationship between BDD tests (`bdd_test.go`) and snapshot/render tests?**

Four untracked `bdd_test.go` files appeared from session 3 (`htmx`, `icons`, `layout`, `navigation`). They heavily overlap with existing snapshot and a11y tests. The BDD tests use user-facing language ("user sees...") but test the same component rendering. Should we:

- **A)** Keep BDD as the canonical behavioral test layer and remove overlapping snapshot tests?
- **B)** Keep snapshot tests as the authoritative render verification and remove BDD tests that only duplicate them?
- **C)** Keep both but enforce BDD tests only cover user behaviors NOT tested elsewhere?

This decision affects 4 files (~600 lines) and 14 clone groups.

---

## Files Modified (This Session)

```
 display/a11y_test.go        | -10 lines  (removed duplicate table test)
 display/avatar_test.go      | ~37 lines  (table-driven status tests)
 display/bdd_test.go         | -119 lines (removed duplicate dropdown/tabs/avatar/table tests)
 display/composition_test.go | -24 lines  (removed duplicate avatar status tests)
 display/table_test.go       | +1 line    (added <caption check)
 examples/demo/main.go       | -60/+17   (extracted alert helper, removed zero values)
 feedback/bdd_test.go        | -31 lines  (removed duplicate toast/loading/step tests)
 forms/bdd_test.go           | ~39 lines  (table-driven FieldError, simplified textarea)
 forms/input.templ           | -2/+1      (use helpText template)
 forms/label.templ           | +7 lines   (extracted helpText template)
 forms/select.templ          | -2/+1      (use helpText template)
 forms/textarea.templ        | -2/+1      (use helpText template)
 htmx/snapshot_test.go       | -8 lines   (removed duplicate GlobalErrorHandling test)
 icons/icon_names_test.go    | ~23 lines  (extracted allIconNames, simplified TestIconCount)
 icons/snapshot_test.go      | ~50 lines  (table-driven IconAttrs, use allIconNames)
 layout/integration_test.go  | ~27 lines  (table-driven assertions)
 navigation/a11y_test.go     | ~18 lines  (merged dark mode tests)
 navigation/mobile_menu.templ| ~10 lines  (extracted menuIconSVG helper)
 navigation/pagination.templ | ~5 lines   (use svg.FillIcon)
 utils/utils_test.go         | ~24 lines  (loop-based MapEnum tests)
```

**Net: -174 lines, 0 regressions, 0 lint issues, all tests green.**
