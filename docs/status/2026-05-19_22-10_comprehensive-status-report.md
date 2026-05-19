# Comprehensive Status Report — templ-components

**Date:** 2026-05-19 22:10 UTC | **Branch:** master | **Ahead of origin:** 0 commits

## Repository Metrics

| Metric                 | Value                                                                              |
| ---------------------- | ---------------------------------------------------------------------------------- |
| `.templ` files         | 32 files, 3,272 lines                                                              |
| Hand-written `.go`     | 18 files, 562 lines                                                                |
| Test files             | 38 files, 5,350 lines                                                              |
| Generated `*_templ.go` | 32 files                                                                           |
| Packages               | 9 (display, feedback, forms, htmx, icons, layout, navigation, utils, internal/svg) |
| Dependencies           | 2 (templ v0.3.1001, tailwind-merge-go v0.2.1)                                      |
| **Build**              | Clean                                                                              |
| **Tests**              | 9/9 pass (with `-race` flag, no data races)                                        |
| **Lint**               | 0 issues                                                                           |
| **Coverage**           | 68.0% total                                                                        |

---

## a) FULLY DONE ✅

### Session 11 (This Session)

| #   | Task                                                                               | Commit               |
| --- | ---------------------------------------------------------------------------------- | -------------------- |
| 1   | C2: Unify FeedbackType — canonical type with AlertType/ToastType aliases           | `712be15`            |
| 2   | C3: Convert LoadingOverlay to props struct with BaseProps                          | `712be15`            |
| 3   | C3: StepIndicatorProps embeds BaseProps                                            | `712be15`            |
| 4   | C4: FillIcon `rotate ...bool` → `rotate bool` (variadic anti-pattern)              | `712be15`            |
| 5   | C4: Fix all 6 FillIcon call sites (accordion, card, dropdown, helpers, pagination) | `712be15`            |
| 6   | C4: ThemeToggle multi-instance fix (IIFE + global guard)                           | `9cc6d10`            |
| 7   | C4: Modal stable IDs (`id="{props.ID}-panel"` + getElementById)                    | `9cc6d10`            |
| 8   | C4: Tooltip aria-describedby                                                       | `9cc6d10`            |
| 9   | C1: Rewrite demo app with layout.Base, Tailwind v4 CDN                             | `0a2d754`            |
| 10  | C5: Breadcrumbs hardcoded SVG → icons.Icon(icons.ChevronRight)                     | `f19dd42`            |
| 11  | C6: Pagination pageURL → net/url                                                   | `f19dd42`            |
| 12  | C6: Restore sync.Mutex in utils.Class() for thread safety                          | `f19dd42`            |
| 13  | C7: Replace splitClasses with strings.Fields                                       | `a23b17f`            |
| 14  | C7: Inline splitSpace wrapper                                                      | `a23b17f`            |
| 15  | C7: Extract BenchmarkHotPaths to benchmark_test.go                                 | `a23b17f`            |
| 16  | C8: Fix stale dropdownSafeID in CONTRIBUTING.md                                    | `e37a36d`            |
| 17  | AGENTS.md updated with breaking changes + conventions                              | `c883235`            |
| 18  | Session 10 audit artifacts (9 skills, plan, architecture diagrams)                 | `ee61026`, `fc91940` |

### Session 10 (Previous Session)

| #   | Task                                                    | Commit                          |
| --- | ------------------------------------------------------- | ------------------------------- |
| 19  | Remove deprecated Exclamation icon                      | `3847355`                       |
| 20  | Icon unknown name → panic (no silent fallback)          | `3847355`                       |
| 21  | ProgressBar percent clamp [0,100]                       | `3847355`                       |
| 22  | Pagination CurrentPage < 1 → clamp to 1                 | `3847355`                       |
| 23  | EmptyStateProps.Icon → icons.Name                       | `3847355`                       |
| 24  | Extract inline strings to package constants (3 commits) | `08f84c0`, `568bd62`, `7a9ec9e` |

---

## b) PARTIALLY DONE ⚠️

### Components That Embed BaseProps But Don't Propagate ID

These components have `utils.BaseProps` but don't render `id={ props.ID }` on their root element. Users who set `ID` on these get nothing:

| Component   | File                              | Missing             |
| ----------- | --------------------------------- | ------------------- |
| Alert       | `feedback/alert.templ:70`         | ID                  |
| Toast       | `feedback/toast.templ:147`        | ID                  |
| StatCard    | `display/card.templ:147`          | ID + Attrs          |
| Nav         | `navigation/nav.templ:33`         | ID + Attrs          |
| Dropdown    | `display/dropdown.templ:77`       | ID                  |
| ProgressBar | `feedback/progress.templ:50`      | ID                  |
| Breadcrumbs | `navigation/breadcrumbs.templ:17` | No BaseProps at all |

### Components Without BaseProps

| Component   | File                              | Impact                            |
| ----------- | --------------------------------- | --------------------------------- |
| Spinner     | `feedback/loading.templ:47`       | No ID/Class/Attrs                 |
| SimpleNav   | `navigation/nav.templ:69`         | No ID/Class/Attrs                 |
| Breadcrumbs | `navigation/breadcrumbs.templ:17` | No ID/Class/Attrs/Nonce           |
| Accordion   | `display/accordion.templ:30`      | No BaseProps (has AccordionProps) |
| Table       | `display/table.templ:49`          | No BaseProps (has TableProps)     |
| Tabs        | `display/tabs.templ:31`           | No BaseProps (has TabsProps)      |

### Modal Focus Restore

`tcCloseModal` references `document.getElementById(id + '-prev-focus')` but `tcOpenModal` never saves `document.activeElement` to that ID. The element is never created. Focus restore after closing modal is broken — users lose their place in the tab order.

### DropdownItem Typed Variant

`DropdownItem` uses `if item.Href != ""` to choose between `<a>` and `<button>`. The button branch has no built-in action wiring — callers must pass `Attrs` with `onclick`. This is a typed variant gap (should be `DropdownItemKind` enum: `LinkItem`, `ButtonItem`).

---

## c) NOT STARTED 📋

### Architecture Improvements

| #   | Task                                                          | Effort  | Impact                                           |
| --- | ------------------------------------------------------------- | ------- | ------------------------------------------------ |
| 1   | Add Radio component (forms/radio.templ)                       | Medium  | High — gap in form library                       |
| 2   | Extract shared `LookupStyle[T]` generic helper                | Low     | Medium — eliminates 5 duplicate lookup functions |
| 3   | Add `utils.AssertContainsClass` test helper                   | Low     | Medium — eliminates fragile class-string tests   |
| 4   | Extract shared dismiss JS (Alert + Toast duplication)         | Low     | Medium — DRY                                     |
| 5   | Add `DefaultXxxProps()` for all enums without one             | Low     | Low — completeness                               |
| 6   | Add Radio group component                                     | Medium  | Medium                                           |
| 7   | Extract shared panic assertion test helper                    | Low     | Low — test DRY                                   |
| 8   | Add `t.Parallel()` to layout/a11y_test.go TestSecurityHeaders | Trivial | Low                                              |

### Type Safety Improvements

| #   | Task                                                           | Effort | Impact                                      |
| --- | -------------------------------------------------------------- | ------ | ------------------------------------------- |
| 9   | Add `ComponentProps` interface with `GetBaseProps() BaseProps` | Medium | Medium — enables generic component wrapping |
| 10  | Refactor BadgeType to use shared FeedbackType                  | Low    | Low — semantic consistency                  |
| 11  | Add `DefaultTooltipPosition()` etc. for missing enum defaults  | Low    | Low                                         |

### Test Coverage Gaps

| #   | Missing Coverage                                    | Component     |
| --- | --------------------------------------------------- | ------------- |
| 12  | Empty `Title`, empty `Message`, unknown `AlertType` | Alert         |
| 13  | Empty `Message`, unknown `ToastType`                | Toast         |
| 14  | `Total=0`, negative `Current`                       | ProgressBar   |
| 15  | Empty `Steps`, `CurrentStep` out of bounds          | StepIndicator |
| 16  | Empty `Links` array, empty `Href`                   | Nav           |
| 17  | Empty `Items` list, item with both Href and action  | Dropdown      |
| 18  | Modal without `Title`                               | Modal         |
| 19  | Empty token string                                  | CSRFToken     |
| 20  | `CurrentPage > TotalPages`, `TotalPages=0`          | Pagination    |

### Documentation

| #   | Task                                                            | Effort  |
| --- | --------------------------------------------------------------- | ------- |
| 21  | Add `CONTRIBUTING.md` section on thread-safety (`classMu`)      | Low     |
| 22  | Document why `*_templ.go` must be committed                     | Done ✅ |
| 23  | Add architecture decision record for `FeedbackType` unification | Low     |

---

## d) TOTALLY FUCKED UP 🚨

### 1. Mutex Removal (Session 10, Recovered in Session 11)

**What happened:** Removed `sync.Mutex` from `utils.Class()` based on incorrect assumption that tailwind-merge-go's internal LRU cache was thread-safe.

**Reality:** Race detector confirmed data races in the merge pipeline (`create-tailwind-merge.go:36-42`). Concurrent tests panicked with nil pointer dereference in the LRU cache.

**Fix:** Restored mutex in commit `f19dd42`.

**Lesson:** Never remove synchronization primitives without running `go test -race`. The library's LRU cache has mutexes, but the wrapper state above it does not.

### 2. Modal Focus Trap — Broken Restore

`tcCloseModal` has dead code for focus restoration. The `id + '-prev-focus'` element is never created by `tcOpenModal`. Users who close a modal via Escape or close button lose keyboard focus position. This is a WCAG failure.

**Fix needed:** Store `document.activeElement` when opening, restore focus when closing.

### 3. Test Fragility — Exact Class String Matching

38 test files assert on exact CSS class strings. Any change to Tailwind class ordering (from tailwind-merge-go updates) or addition of new classes breaks tests. Examples:

- `display/accordion_test.go:47` — exact `"overflow-hidden transition-all duration-200 max-h-0"`
- `display/table_test.go:96` — `"border border-gray-200"`
- `icons/snapshot_test.go:47` — `"h-5 w-5"`

These are implementation tests, not behavior tests. Should use `AssertContainsClass` (word-level token matching).

### 4. Dead Code in utils.go

Three exported functions with zero references:

- `Deref[T](p *T) T` — Go 1.22 has `cmp.Or`, but this is still unused
- `DerefOr[T](p *T, fallback T) T` — zero references
- `MergeAttrs(m ...templ.Attributes) templ.Attributes` — zero references

They increase API surface without value. Should be removed or at least documented as deprecated.

### 5. fillIcon Wrapper is Pointless

`display/helpers.templ` contains a 3-line wrapper:

```templ
templ fillIcon(class string, path string, rotate bool) {
    @svg.FillIcon(class, path, rotate)
}
```

This adds zero value — just calls through. Every caller could `@svg.FillIcon` directly. The wrapper obscures where the real implementation lives.

### 6. Hardcoded Inline SVGs (3 Remaining)

After fixing breadcrumbs, 3 components still have hardcoded inline SVGs instead of using `icons.Icon`:

- **Alert dismiss button** (`feedback/alert.templ:92`) — X icon, should use `icons.X`
- **Toast dismiss button** (`feedback/toast.templ:126`) — X icon in JS string, should use `icons.IconPathJS(icons.X)`
- **StepIndicator checkmark** (`feedback/progress.templ:123`) — Check icon, should use `icons.Check`

### 7. Duplicate Map-Lookup-With-Fallback Pattern (5x)

Identical logic copy-pasted across 5 files:

```go
func xxxLookup(t T) V {
    if v, ok := xxxMap[t]; ok { return v }
    return xxxDefault
}
```

Locations:

- `display/badge.templ:107` (badgeLookupStyle)
- `display/tooltip.templ:54` (tooltipLookupPosition)
- `display/card.templ:83` (cardPaddingClass)
- `display/modal_go.go:48` (modalSizeClass)
- `display/badge.templ:76` (badgeSizeClass)

---

## e) WHAT WE SHOULD IMPROVE 🔧

### Immediate (Next Session)

1. **Fix Modal focus restore** — store activeElement on open, restore on close
2. **Fix all BaseProps ID propagation gaps** — Alert, Toast, StatCard, Nav, Dropdown, ProgressBar
3. **Add Radio component** — only missing basic form input type
4. **Extract shared Lookup helper** — one generic function replacing 5 copies
5. **Extract shared dismiss JS** — Alert and Toast share identical inline script
6. **Add `utils.AssertContainsClass` test helper** — replace fragile exact-string tests
7. **Remove dead code** — Deref, DerefOr, MergeAttrs
8. **Replace hardcoded SVGs with icons.Icon** — Alert dismiss, StepIndicator checkmark

### Short Term

9. **Add proper focus management to Modal** — `aria-modal`, focus sentinel elements
10. **Add `templ.Component` children support to Table** — currently only takes data, no slot for custom cell rendering
11. **Add loading state to Button** — `forms/button.templ` with spinner
12. **Extract shared test panic assertion** — used in 3 files
13. **Add `Nonce` propagation to all inline script tags** — verify every `<script>` has `nonce={ props.Nonce }`

### Medium Term

14. **DropdownItem typed variant** — `DropdownItemKind` enum (LinkItem, ButtonItem)
15. **SimpleNav → SimpleNavProps** — 4 positional params → struct
16. **Extract shared JS init pattern** — 7 components with inline JS could share a single pattern
17. **Add component-level CSS extraction** — generate static CSS instead of relying on Tailwind CDN
18. **Add dark mode toggle to demo app** — currently has theme support but no visible toggle in demo

---

## f) Top #25 Things to Get Done Next 📋

Sorted by impact / effort ratio (Pareto):

| #   | Task                                                            | Effort | Impact                   | Priority |
| --- | --------------------------------------------------------------- | ------ | ------------------------ | -------- |
| 1   | Fix Modal focus restore (store/restore activeElement)           | 15min  | High (a11y)              | 🔴       |
| 2   | Add `id={ props.ID }` to components missing it (6 components)   | 30min  | High (API consistency)   | 🔴       |
| 3   | Extract shared `LookupStyle` generic (replaces 5 copies)        | 20min  | Medium (DRY)             | 🔴       |
| 4   | Extract shared dismiss JS (Alert + Toast)                       | 20min  | Medium (DRY)             | 🔴       |
| 5   | Add `utils.AssertContainsClass` + refactor fragile tests        | 30min  | High (test reliability)  | 🔴       |
| 6   | Add Radio component                                             | 45min  | High (feature gap)       | 🟡       |
| 7   | Replace hardcoded SVGs with icons.Icon (3 remaining)            | 20min  | Low-Medium (SSOT)        | 🟡       |
| 8   | Remove dead code (Deref, DerefOr, MergeAttrs)                   | 10min  | Low (cleanliness)        | 🟡       |
| 9   | Inline `fillIcon` wrapper (call svg.FillIcon directly)          | 15min  | Low (simplicity)         | 🟡       |
| 10  | Add test coverage for error paths (empty inputs, out-of-bounds) | 60min  | Medium (reliability)     | 🟡       |
| 11  | Extract shared panic assertion test helper                      | 10min  | Low (test DRY)           | 🟡       |
| 12  | Add `t.Parallel()` to TestSecurityHeaders                       | 2min   | Trivial                  | 🟡       |
| 13  | Add `DefaultXxxProps()` for enums missing them                  | 20min  | Low (completeness)       | 🟢       |
| 14  | DropdownItem typed variant                                      | 45min  | Medium (type safety)     | 🟢       |
| 15  | SimpleNav → SimpleNavProps                                      | 30min  | Medium (API consistency) | 🟢       |
| 16  | Add dark mode toggle to demo app                                | 30min  | Low (demo polish)        | 🟢       |
| 17  | Add Button with loading state                                   | 30min  | Medium (feature)         | 🟢       |
| 18  | Extract shared JS init pattern                                  | 60min  | Medium (maintainability) | 🟢       |
| 19  | Add Radio group component                                       | 30min  | Medium (feature)         | 🟢       |
| 20  | BadgeType → FeedbackType unification                            | 30min  | Low (consistency)        | 🔵       |
| 21  | Add `ComponentProps` interface                                  | 60min  | Medium (architecture)    | 🔵       |
| 22  | Document thread-safety in CONTRIBUTING.md                       | 10min  | Low (docs)               | 🔵       |
| 23  | Add ADR for FeedbackType unification                            | 15min  | Low (docs)               | 🔵       |
| 24  | Generate static CSS instead of Tailwind CDN                     | 120min | High (production)        | 🔵       |
| 25  | Add form validation JS                                          | 90min  | Medium (feature)         | 🔵       |

---

## g) Top #1 Question I Cannot Figure Out Myself ❓

**Question: Should `Breadcrumbs` accept `BaseProps` (ID/Class/Attrs) or stay as a simple `[]BreadcrumbItem` parameter?**

The current signature is:

```go
templ Breadcrumbs(items []BreadcrumbItem)
```

Every other navigation component (Nav, NavLink, Pagination) uses props structs with `BaseProps`. Breadcrumbs is the outlier.

**Arguments for adding BaseProps:**

- Consistency with every other component in the library
- Users might want `aria-label`, custom `class`, or `id` on the `<nav>`
- The `navigation` package would have uniform API shape

**Arguments against:**

- Breadcrumbs is conceptually simpler than other components — just a list of items
- Adding `BreadcrumbProps` with `BaseProps` means callers write `Breadcrumbs(BreadcrumbProps{Items: [...]})` instead of `Breadcrumbs([...])`
- The current API is more ergonomic for the common case

**The deeper tension:** This library's design philosophy oscillates between "every component is a props struct with BaseProps" (the stated convention) and "simple components get simple signatures" (pragmatism). If we add BaseProps to Breadcrumbs, do we also add it to Spinner (currently `Spinner(size, colorClass)`)? To SimpleNav (currently `SimpleNav(brandText, brandHref, links, currentPath)`)? Where do we draw the line?

I genuinely don't know where the boundary should be. The convention says "all component props embed BaseProps" but that makes simple components verbose. Breaking the convention for ergonomics creates inconsistency. **What is the intended design philosophy for simple vs. complex components?**
