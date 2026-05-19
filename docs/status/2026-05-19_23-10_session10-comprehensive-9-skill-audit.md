# Status Report — templ-components

**Date:** 2026-05-19 23:10 | **Session:** 10 | **Type:** Comprehensive 9-Skill Audit

---

## Executive Summary

Session 10 ran a full 9-skill audit: **code-quality-scan**, **features-audit**, **BDD-testing**, **full-code-review**, **improve-codebase-architecture**, **architecture-review**, **architecture-visualization**, **go-modularize**, and **todo-list-builder**.

**The project is in excellent shape.** Build is clean, lint is clean, all 9 packages pass tests at 71.8% coverage. The library has strong type safety (17 typed enums), clean package DAG (0 circular imports), only 2 external dependencies, and comprehensive accessibility.

However, the audit surfaced **51 actionable items** across 5 priority tiers. The most critical are: a broken demo app, silent icon fallback bugs, JS architecture fragmentation (222 lines across 7 files), and AlertType/ToastType type duplication.

**No code changes were made in this session.** This was a pure audit/analysis session producing documentation artifacts only (FEATURES.md corrections, TODO_LIST.md rewrite, architecture diagrams, planning documents).

---

## A) FULLY DONE ✅

### Build & Quality

- Build: `go build ./...` — **0 errors**
- Tests: `go test ./...` — **9/9 packages pass**
- Lint: `golangci-lint run` — **0 issues**
- Coverage: **71.8%** (range: 70.5%–89.5%)
- CI: GitHub Actions with Go 1.26, lint+build+test

### Architecture Foundation

- **Clean DAG dependency graph** — `utils ← all`, no circular imports
- **`utils.BaseProps` embedding** — All 29 component props structs compose through shared base
- **Map-based style lookups** — All visual variant lookups use maps, not switches
- **`feedbackStyleSet` + generic `lookupFeedbackStyle[T]()`** — Shared style infrastructure
- **`iconPathData` map** — Data-driven icon rendering for 44 icons
- **Typed enums** — 17 string-based enums making impossible states unrepresentable
- **`internal/svg`** — Properly internal, not importable by consumers
- **CSP compliance** — All inline scripts use nonce attribute
- **Comprehensive a11y** — Modal focus trap, dropdown keyboard nav, tabs ARIA, tooltip role
- **Generated `*_templ.go` committed** — 31 files, correct for publishable library

### Documentation

- `AGENTS.md` — Comprehensive project knowledge
- `CONTEXT.md` — Architecture context and patterns
- `FEATURES.md` — Full feature inventory (updated this session)
- `TODO_LIST.md` — Complete prioritized task list (rewritten this session)
- `CHANGELOG.md` — Full changelog with breaking changes
- `docs/migration/v0.1-to-v0.2.md` — Migration guide
- `docs/adr/` — 3 architecture decision records

### Session 10 Audit Artifacts (All Produced)

- `docs/architecture-understanding/2026-05-19_22-56-architecture-review.md` — Full review with scalability, composability, shallow module analysis
- `docs/architecture-understanding/2026-05-19_22-56-current-state.d2/.svg` — Current dependency graph
- `docs/architecture-understanding/2026-05-19_22-56-target-state-improved.d2/.svg` — Target architecture
- `docs/modularization/ANALYSIS-2026-05-19.md` — Verdict: NOT recommended
- `docs/planning/2026-05-19_22-56-session10-comprehensive-audit.md` — Pareto execution plan

---

## B) PARTIALLY DONE 🔨

| Item                   | What's Done                                                                     | What's Missing                                                                                                                    |
| ---------------------- | ------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| **Tooltip a11y (#41)** | `role="tooltip"` + deterministic ID on tooltip div                              | Trigger element lacks `aria-describedby` linkage pointing to tooltip ID                                                           |
| **Icon system**        | 44 icons with map-driven rendering, `IconPathJS()` for JS consumption           | Silent fallback to clock for unknown names; `IconPathJS` stroke-width mismatch (2 vs 1.5); deprecated `Exclamation` still present |
| **Test suite**         | 71.8% coverage, BDD + snapshot + a11y + benchmark across all packages           | 60-80% assertion overlap between bdd_test/snapshot_test/a11y_test; test data duplicated 3× in navigation                          |
| **JS architecture**    | Global singleton pattern for accordion/dropdown, shared dismiss for alert/toast | No shared init strategy; 222 lines fragmented across 7 files; HTMX swap breaks event listeners; ThemeToggle multi-instance bug    |
| **Demo app**           | Builds and runs                                                                 | Doesn't use `layout.Base`, uses Tailwind v2 (library targets v4), no HTMX loaded, `_ = props` discards constructed PageProps      |

---

## C) NOT STARTED ⬜

### P0 — Critical (4 items)

1. **Fix demo app to use `layout.Base`** — Currently raw HTML, wrong Tailwind version, no HTMX
2. **Remove or validate unknown icon names** — `Name("typo")` silently renders clock icon
3. **Fix `IconPathJS` stroke-width mismatch** — JS version uses `stroke-width="2"` vs templ `stroke-width="1.5"`
4. **Delete deprecated `Exclamation` icon** — Identical SVG path to `ExclamationCircle`, dead code

### P1 — Architecture (16 items)

5. Unify `AlertType`/`ToastType` into shared type
6. Merge `alertStyleMap`/`toastStyleMap`
7. Make `SimpleCard` compose through `Card`
8. Add `ComponentProps` interface to `utils.BaseProps`
9. Use stable IDs in modal JS instead of CSS selectors
10. Fix ThemeToggle multi-instance bug
11. Use icon system in Breadcrumbs chevron
12. Add `BaseProps` to `StepIndicatorProps`
13. Convert `LoadingOverlay` from positional params to props struct
14. Validate `SwapOOB` swapStyle parameter
15. Validate `Pagination` CurrentPage > 0
16. Clamp `ProgressBar` percent to [0, 100]
17. Consolidate inline JS into shared init strategy
18. Fix HTMX swap event listener re-attachment
19. Fix `GlobalErrorHandling` shared retry counter race condition
20. Make `GlobalErrorHandling` config values configurable
21. Consolidate modal per-instance JS into single function

### P2 — Quality (10 items)

22. Replace `DropdownItem` empty-Href discrimination with typed enum
23. Change `FillIcon` variadic `rotate ...bool` to `rotate bool`
24. Audit `tailwind-merge-go` thread safety, remove mutex
25. Replace `BoolString()` with `strconv.FormatBool`
26. Validate `SelectOption` contradiction (Disabled+Selected)
27. Use `net/url` for pagination URL construction
28. Add `uint` type for Pagination fields
29. Eliminate 4-way icon list split brain
30. Validate `|` separator doesn't appear in SVG paths
31. Document 20×20 fill vs 24×24 stroke convention

### P2 — Test Cleanup (6 items)

32. Consolidate test files: 1-2 per package
33. Remove unused `badgeTextLive` constant
34. Delete `TestPtr` in `utils_test.go`
35. Replace `splitSpace`/`splitClasses` with `strings.Fields`
36. Move `BenchmarkHotPaths` out of `a11y_test.go`
37. Remove duplicate test data declarations in navigation

### P2 — Documentation (4 items)

38. Update TODO #11 note: `dropdownSafeID` was removed
39. Update CONTRIBUTING.md: remove `dropdownSafeID` reference
40. Document `htmx` → `feedback` runtime JS dependency
41. Fix Tooltip `aria-describedby` linkage on trigger element

### P3 — Deferred (6 items)

42. Convert snapshot tests to golden file comparison
43. Move test helpers out of `utils/`
44. Documentation site generation
45. Add Radio, File input, Toggle/Switch form components
46. Add client-side JS tab switching
47. Make `PageProps` zero-value safe

---

## D) TOTALLY FUCKED UP ❌

### Demo App — An Anti-Advertisement for the Library

`examples/demo/main.go` is the **worst file in the project**. It:

- Constructs `layout.DefaultPageProps()` then **immediately discards it** with `_ = props`
- Writes raw `<html><head>...` strings instead of using `layout.Base`
- Uses Tailwind CSS **v2.2.19** from CDN while the library targets **Tailwind v4** classes
- Doesn't load HTMX — none of the HTMX components can work
- Has no CSP nonce — defeats the library's entire security model
- Is the only thing consumers see if they clone and run — **it's the wrong first impression**

**Severity: This must be fixed before any public release.**

### Silent Icon Fallback — Data Corruption by Default

`icons/icon_paths.go` silently returns a clock icon for unknown `Name` values. A typo like `icons.Name("hoem")` compiles fine and renders a clock with zero indication something is wrong. In a UI component library, rendering the **wrong icon silently** is worse than crashing.

**Severity: High. Every consumer will hit this eventually.**

### `IconPathJS` Stroke-Width Split Brain

The same icon renders differently depending on whether it's rendered via templ (`stroke-width="1.5"`) or via the JS toast system (`stroke-width="2"`). This is a visible inconsistency consumers will notice when comparing in-page icons vs toast icons.

**Severity: Medium-High. Visible rendering difference.**

---

## E) WHAT WE SHOULD IMPROVE

### Architecture

1. **JS runtime consolidation** — 222 lines of inline JS across 7 files with no shared init strategy. Create a `tc-init.js` module loaded once by `layout.Base` that provides `tc.register(id, initFn)`, `tc.dismiss(selector)`, and `tc.onSwap(callback)`.
2. **`ComponentProps` interface** — 29 props structs all embed `BaseProps` but share no common interface. Add `GetBaseProps() BaseProps` for generic handling.
3. **`FeedbackLevel` shared type** — `AlertType` and `ToastType` are identical 4-value enums. Merge into one type.
4. **DropdownItem sum type** — Currently uses empty-string discrimination (Href="" = button). Should be a typed variant.

### Test Suite

5. **Consolidate 37 → ~15 test files** — Every package has bdd_test + snapshot_test + a11y_test with 60-80% overlap. Merge into component_test.go + a11y_test.go.
6. **Icon list split brain** — The same icon enumeration exists in 4 places (constants, map keys, `allIconNames()`, BDD inline list). Auto-generate from map.
7. **Dead test code** — `TestPtr` tests Go built-in `new()`, not library code. Remove.

### API Consistency

8. **`StepIndicatorProps` missing `BaseProps`** — Only feedback component without it.
9. **`LoadingOverlay` uses positional params** — Only feedback component not using props struct.
10. **`SimpleCard` duplicates `Card` shell** — Should compose through `Card` instead.
11. **`BoolString()` duplicates `strconv.FormatBool`** — Only one consumer. Replace.
12. **`ProgressBar` no upper-bound clamp** — `width: 105%+` overflows visually.
13. **`Pagination` accepts negative/zero `CurrentPage`** — Renders broken page links.

### Performance

14. **`utils.Class()` global `sync.Mutex`** — Every `Class()` call locks. If `tailwind-merge-go` is stateless, the mutex is unnecessary overhead.

### Demo

15. **Fix to use the actual library** — Use `layout.Base`, Tailwind v4, HTMX, proper nonce. Make it a showcase, not an embarrassment.

---

## F) Top 25 Things to Do Next

Sorted by impact × effort (Pareto):

| #   | Task                                                               | Effort | Impact   | Priority |
| --- | ------------------------------------------------------------------ | ------ | -------- | -------- |
| 1   | Fix demo app to use `layout.Base` + Tailwind v4 + HTMX             | 45min  | CRITICAL | P0       |
| 2   | Error/panic on unknown icon names instead of silent clock fallback | 15min  | HIGH     | P0       |
| 3   | Fix `IconPathJS` stroke-width mismatch (2 → 1.5)                   | 10min  | HIGH     | P0       |
| 4   | Delete deprecated `Exclamation` icon                               | 10min  | MEDIUM   | P0       |
| 5   | Clamp `ProgressBar` percent to [0, 100]                            | 5min   | MEDIUM   | P1       |
| 6   | Validate `Pagination.CurrentPage` > 0                              | 5min   | MEDIUM   | P1       |
| 7   | Remove unused `badgeTextLive` constant                             | 2min   | LOW      | P1       |
| 8   | Delete `TestPtr` in utils_test.go                                  | 2min   | LOW      | P1       |
| 9   | Unify `AlertType`/`ToastType` into shared `FeedbackLevel`          | 30min  | HIGH     | P1       |
| 10  | Merge `alertStyleMap`/`toastStyleMap`                              | 20min  | MEDIUM   | P1       |
| 11  | Add `BaseProps` to `StepIndicatorProps`                            | 15min  | MEDIUM   | P1       |
| 12  | Convert `LoadingOverlay` to props struct                           | 20min  | MEDIUM   | P1       |
| 13  | Use stable IDs in modal JS (`props.ID + "-panel"`)                 | 20min  | MEDIUM   | P1       |
| 14  | Use `icons.Icon`/`svg.FillIcon` in Breadcrumbs chevron             | 15min  | MEDIUM   | P1       |
| 15  | Change `FillIcon` variadic `rotate ...bool` → `rotate bool`        | 10min  | LOW      | P1       |
| 16  | Make `SimpleCard` compose through `Card`                           | 30min  | MEDIUM   | P2       |
| 17  | Add `ComponentProps` interface to `utils.BaseProps`                | 30min  | MEDIUM   | P2       |
| 18  | Audit `tailwind-merge-go` thread safety, remove mutex              | 30min  | MEDIUM   | P2       |
| 19  | Fix ThemeToggle multi-instance bug                                 | 15min  | MEDIUM   | P2       |
| 20  | Consolidate test files (37 → ~15)                                  | 120min | HIGH     | P2       |
| 21  | Consolidate inline JS into shared init strategy                    | 90min  | HIGH     | P2       |
| 22  | Eliminate icon list 4-way split brain                              | 30min  | MEDIUM   | P2       |
| 23  | Replace `DropdownItem` empty-Href with typed variant               | 45min  | MEDIUM   | P2       |
| 24  | Fix Tooltip `aria-describedby` on trigger element                  | 15min  | MEDIUM   | P2       |
| 25  | Document `htmx` → `feedback` runtime JS coupling                   | 5min   | LOW      | P2       |

---

## G) Top #1 Question I Cannot Figure Out Myself

**Is `tailwind-merge-go` v0.2.1's `twmerge.Merge()` function thread-safe (stateless/pure)?**

The `utils.Class()` function wraps every call in a `sync.Mutex` because the developer wasn't sure. If the merge function is stateless (pure function with no mutable state), the mutex is unnecessary overhead that serializes all template rendering under load. If it IS stateful, the mutex is critical.

I cannot determine this without either:

1. Reading the `tailwind-merge-go` source code to verify no mutable state exists
2. Finding explicit documentation about thread safety
3. Running concurrent benchmarks to prove/disprove contention

**Why it matters:** `utils.Class()` is called by every component on every render. Under high-throughput server-side rendering, a global mutex here is a serializing bottleneck. Removing it (if safe) is a zero-risk performance win.

---

## Metrics Snapshot

| Metric           | Value                           |
| ---------------- | ------------------------------- |
| Packages         | 9                               |
| Templ components | 53                              |
| Icon names       | 44                              |
| Typed enums      | 17                              |
| Props structs    | 29                              |
| Test files       | 37                              |
| Test coverage    | 71.8%                           |
| External deps    | 2                               |
| Circular imports | 0                               |
| Inline JS lines  | 222 (across 7 files)            |
| Max file size    | 548 lines (`forms/bdd_test.go`) |
| Source lines     | ~18K                            |
| Build status     | ✅ Clean                        |
| Lint status      | ✅ 0 issues                     |
| Test status      | ✅ 9/9 pass                     |

---

## Files Modified This Session (Audit Only — No Code Changes)

| File                    | Change                                                                                           |
| ----------------------- | ------------------------------------------------------------------------------------------------ |
| `AGENTS.md`             | Updated date                                                                                     |
| `FEATURES.md`           | Fixed: 42→44 icons, removed stale Ptr, updated coverage, fixed icon list, removed dead IconAttrs |
| `TODO_LIST.md`          | Complete rewrite: 51 items across 5 priority tiers, verified against code                        |
| `*_templ.go` (31 files) | Regenerated by `templ generate` — import formatting only, no logic changes                       |

## New Files This Session

| File                                                                         | Purpose                                  |
| ---------------------------------------------------------------------------- | ---------------------------------------- |
| `docs/architecture-understanding/2026-05-19_22-56-architecture-review.md`    | Full architecture review                 |
| `docs/architecture-understanding/2026-05-19_22-56-current-state.d2`          | Current dependency graph source          |
| `docs/architecture-understanding/2026-05-19_22-56-current-state.svg`         | Current dependency graph rendered        |
| `docs/architecture-understanding/2026-05-19_22-56-target-state-improved.d2`  | Target architecture source               |
| `docs/architecture-understanding/2026-05-19_22-56-target-state-improved.svg` | Target architecture rendered             |
| `docs/modularization/ANALYSIS-2026-05-19.md`                                 | Modularization analysis: NOT recommended |
| `docs/planning/2026-05-19_22-56-session10-comprehensive-audit.md`            | Pareto execution plan                    |

---

_Arte in Aeternum_
