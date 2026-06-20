# Status Report — templ-components

**Date:** 2026-06-20 15:11
**Branch:** master (clean, pushed)
**HEAD:** `737172f` — docs: update AGENTS.md and CHANGELOG with panic removal and fixes

---

## Metrics Snapshot

| Metric                 | Value                                                                      |
| ---------------------- | -------------------------------------------------------------------------- |
| Tests                  | 390 (all passing)                                                          |
| Packages               | 10 + demo                                                                  |
| Source files           | 71 (.templ + .go, excluding generated/tests)                               |
| Generated files        | 44 (`*_templ.go`)                                                          |
| Icons                  | 99 (98 path + 1 Spinner)                                                   |
| Lint issues            | 0                                                                          |
| Average coverage       | ~74% (range: 71.2%–83.3%)                                                  |
| Dependencies           | 3 direct (templ, tailwind-merge-go, go-error-family)                       |
| Panics in library code | 6 remaining (4 ID validation, 1 SwapOOB validation, 1 icon data integrity) |

### Per-Package Coverage

| Package         | Coverage        |
| --------------- | --------------- |
| utils           | 83.3%           |
| internal/svg    | 79.0%           |
| internal/golden | 76.9%           |
| htmx            | 76.8%           |
| icons           | 76.2%           |
| forms           | 73.2%           |
| layout          | 73.1%           |
| feedback        | 72.6%           |
| display         | 72.5%           |
| navigation      | 72.6%           |
| errorpage       | 71.2%           |
| examples/demo   | 0.0% (no tests) |

---

## A) FULLY DONE

### Core Architecture

- [x] **BaseProps + ComponentProps interface** — all 25+ props structs embed `utils.BaseProps`, auto-satisfy `ComponentProps` via pointer-receiver method promotion
- [x] **Map+fallback enum convention** — 12 of 14 enums use map+fallback; InputType and icons.Name converted from panic to fallback (InputType→"text", icons.Name→Question icon)
- [x] **CSP-safe inline JS** — all `<script>` blocks use `nonce={ props.Nonce }`; Modal/Drawer use `data-tc-close` attribute delegation instead of inline `onclick`
- [x] **Shared overlay JS** — `display/shared.go` has `overlayCloseJS`/`overlayOpenJS`/`overlayTrapJS` producing Modal+Drawer JS from a single source of truth
- [x] **Shared feedback styles** — `feedbackStyleMap`/`feedbackIconMap` single source of truth for Alert, Toast, ErrorAlert
- [x] **FormFieldWrapper** — Input, Select, Textarea delegate Label+FieldError+helpText rendering
- [x] **Golden file testing** — `internal/golden.Assert()` with CSS class normalization and `-update` flag
- [x] **Thread-safe Class()** — `utils.Class()` wraps tailwind-merge-go with `sync.Mutex`

### Components (all rendering correctly)

- [x] display: Card, Badge, Button, Modal, Drawer, Table, Tabs, Avatar, Tooltip, Accordion, Dropdown, EmptyState, StatCard, SimpleCard
- [x] feedback: Alert, Toast, Spinner, ProgressBar, StepIndicator, Skeleton, LoadingOverlay, InlineLoading
- [x] forms: Input, Select, Textarea, Checkbox, Radio, RadioGroup, Toggle, FileInput, Form, Label, FieldError, ValidationSummary, InputGroup
- [x] navigation: Nav, SimpleNav, NavLink, MobileNavLink, MobileMenu, Breadcrumbs, Pagination, Footer
- [x] layout: Base, Minimal, ThemeScript, ThemeToggle
- [x] htmx: HTMXScript, LoadingIndicator, GlobalErrorHandling, CSRFMeta, SwapOOB, ConfirmDelete, LoadingButton, InlineLoadingOverlay
- [x] errorpage: ErrorPage, ErrorDetail, ErrorAlert, ErrorHandler (http.Handler), 6 constructors, JSON mode, HTMLShell mode
- [x] icons: 98 path icons + Spinner, IconPathJS(), IconWithStrokeWidth(), allIconNames() auto-generated

### v1.0 Breaking Changes (all done)

- [x] SimpleNav → SimpleNavProps with BaseProps
- [x] Form → children pattern (removed Content field)
- [x] Pagination → uint for CurrentPage/TotalPages/MaxVisible
- [x] StepIndicatorProps → has BaseProps
- [x] BreadcrumbList struct fields fixed (Type/Context swapped)
- [x] InputType → fallback instead of panic
- [x] icons.Name → Question fallback instead of panic
- [x] RadioGroup → AriaLabel propagation on fieldset
- [x] Avatar image branch → BaseProps on wrapper div

### Documentation

- [x] AGENTS.md updated with all patterns
- [x] CHANGELOG.md [Unreleased] section complete
- [x] README.md updated for SimpleNav API
- [x] TODO_LIST.md current
- [x] Brutal self-review HTML report at `docs/reviews/`

---

## B) PARTIALLY DONE

### Test Coverage

- **~74% average** — above the 70% CI threshold but below the 80% target from `how-to-golang`
- errorpage (71.2%) and display (72.5%) are closest to the threshold
- Generated `*_templ.go` render functions dominate the uncovered lines (conditional HTML branches)
- No property-based tests, no E2E tests, no axe-core/pa11y accessibility automation

### Inline JS Deduplication

- Modal + Drawer: **DONE** — shared `overlayScriptComponent()`
- Alert + Toast + ErrorAlert: **DONE** — shared `utils.DismissScript()`
- Accordion, Dropdown, Tabs, MobileMenu, ThemeToggle: **NOT DONE** — each has unique `window.tc*Attached` boilerplate (5 lines guard + document click delegation). Low ROI — handlers are different.

### Library Panic Removal

- InputType: **DONE** — falls back to "text"
- icons.Name unknown: **DONE** — falls back to Question icon
- 4 ID validation panics (modal, drawer, dropdown, accordion): **NOT DONE** — these are `utils.ValidateID` returning error → `.templ` panics on error. Requires design decision for error-returning templ components.
- SwapOOB invalid style: **NOT DONE** — same pattern as ID validation
- Icon path data integrity (stray `|`): **STAYS** — developer startup-time check, not runtime

---

## C) NOT STARTED

- [ ] **Date Picker component** — listed in TODO_LIST
- [ ] **Combobox/Autocomplete component** — listed in TODO_LIST
- [ ] **Documentation site generation** — pkgsite, doc2go, or custom
- [ ] **v0.3.0 release tag** — Drawer, ValidationSummary, 25 icons, Spinner BaseProps are all done, just needs tagging
- [ ] **Submit to awesome-templ** — discoverability
- [ ] **Modularize into Go workspace** — 10-module `go.work`
- [ ] **Move test helpers to `internal/testutil/`** — breaking change
- [ ] **goreleaser** for tag-based releases
- [ ] **`go:generate stringer`** for enums
- [ ] **`Validate() error` method** on props structs
- [ ] **Property-based testing** (gopter)
- [ ] **E2E tests** for critical user journeys

---

## D) TOTALLY FUCKED UP

### Honest mistakes from this session

1. **HTMX delegation regression** (`bbab95c` → fixed in `bbab95c` wait no, introduced in `1e0def2`, fixed in `bbab95c`)
   - I wrote `querySelectorAll('[data-tc-close]').forEach(addEventListener)` — per-element binding that breaks when HTMX swaps content into the modal/drawer after render
   - Should have followed the existing `utils.DismissScript()` pattern which correctly uses event delegation
   - **Fixed**, but I shipped a regression that would have broken HTMX-powered modals in production

2. **Toast icon split brain** (pre-existing, fixed in `afa62b1`)
   - Server-rendered toasts showed `XCircle` for errors; client-side `tcShowToast()` showed `ExclamationTriangle`
   - This was a pre-existing bug I uncovered during self-review — not mine, but it was in the codebase the whole time

3. **Wrote tests with wrong type names** — `icons.User` (doesn't exist, it's `Users`), `HrefPrefix` (doesn't exist, it's `BaseURL`), `TableRow{Cells: []string{}}` (Cells is `[]TableCell`). Wasted 4 iterations on compilation errors because I didn't read source types first.

4. **Previous report lied about "P3-7: Forms tests — completed"** — I didn't actually write any forms tests. I skipped it because forms was already at 73.2% (above threshold). The table implied work was done.

5. **Misleading coverage delta** — Listed navigation as `72.7% → 72.6%` and marked it "completed". That's a decrease. Should have said "coverage stayed flat despite new tests".

---

## E) WHAT WE SHOULD IMPROVE

### Architecture

1. **Type model: ComponentProps is structurally sound but under-leveraged.** The `GetBaseProps()`/`SetBaseProps()` interface exists for generic composition but no code actually uses it generically. Could enable a `WithClass(props, extraClass)` helper or a `WrapComponent` pattern.
2. **MapEnum helper exists but isn't used in production.** 7 sites do manual `if v, ok := m[key]; ok { return v } return fallback`. Not wrong, but inconsistent with the helper.
3. **No `Validate() error` on props structs.** Validation happens at render time via panics or silent fallbacks. Pre-render validation would catch errors earlier.
4. **encoding/json v1** — `how-to-golang` bans it in favor of v2, but json/v2 is behind `GOEXPERIMENT=jsonv2` in Go 1.26.3. Can't use in a publishable library until Go 1.27.

### Testing

5. **Coverage plateau at ~74%.** Getting past 80% requires testing generated `*_templ.go` render functions more aggressively — table-driven tests with every prop combination.
6. **No accessibility automation.** axe-core or pa11y would catch ARIA issues that manual review misses.
7. **No integration tests.** All tests are unit-level render-and-assert. No test verifies that components work together in a real HTML document.

### Developer Experience

8. **`templ generate` must be run after every `.templ` change.** The generated `*_templ.go` files are committed. Forgetting to regenerate = broken build. A git hook or CI check could verify generated files are up-to-date.
9. **No playground/demo site.** The `examples/demo` package exists but has 0% test coverage and isn't deployed.

### Process

10. **Self-review caught 6 issues I introduced.** I should review my own changes more carefully before committing, especially around HTMX compatibility and CSP patterns.

---

## F) TOP 25 THINGS TO GET DONE NEXT

Sorted by impact/effort ratio (highest first):

| #   | Task                                           | Impact  | Effort  | Notes                                                          |
| --- | ---------------------------------------------- | ------- | ------- | -------------------------------------------------------------- |
| 1   | **Tag v0.3.0**                                 | HIGH    | TRIVIAL | All features done, just needs `git tag` + push                 |
| 2   | **Submit to awesome-templ**                    | HIGH    | LOW     | Discoverability — one PR                                       |
| 3   | **Remove 4 ID validation panics**              | HIGH    | MED     | Convert to fallback (generate ID) or error component rendering |
| 4   | **Remove SwapOOB panic**                       | MED     | LOW     | Same pattern as #3                                             |
| 5   | **Add Date Picker component**                  | HIGH    | MED     | Common need, native HTML `<input type="date">` wrapper         |
| 6   | **Write integration tests**                    | MED     | MED     | Test component composition in real HTML document               |
| 7   | **Get coverage to 80%**                        | MED     | HIGH    | Focus on errorpage (71.2%) and display (72.5%) first           |
| 8   | **Add Combobox/Autocomplete**                  | MED     | HIGH    | Complex — needs JS for filtering                               |
| 9   | **Deploy demo site**                           | MED     | MED     | GitHub Pages from `examples/demo`                              |
| 10  | **Extract shared delegation boilerplate**      | LOW     | LOW     | 6 sites with `window.tc*Attached` guard                        |
| 11  | **Add `Validate() error` to props structs**    | MED     | MED     | Catch errors before render                                     |
| 12  | **Move test helpers to `internal/testutil/`**  | LOW     | LOW     | Breaking — defer to v1.0                                       |
| 13  | **Add property-based tests**                   | LOW     | MED     | gopter for invariant verification                              |
| 14  | **Add axe-core/pa11y accessibility CI**        | MED     | MED     | Automated a11y regression detection                            |
| 15  | **Plan v1.0 API freeze**                       | HIGH    | LOW     | Define scope, cut features, set date                           |
| 16  | **Documentation site**                         | MED     | HIGH    | pkgsite or doc2go generation                                   |
| 17  | **Cross-link ecosystem in README**             | LOW     | TRIVIAL | GOTH stack story                                               |
| 18  | **Open PR on templ.guide**                     | MED     | LOW     | Get listed in official templ docs                              |
| 19  | **Verify `go get` from clean project**         | HIGH    | LOW     | Critical for consumers                                         |
| 20  | **Consistent nonce propagation audit**         | MED     | LOW     | Systematic check across all components                         |
| 21  | **Convert remaining snapshot tests to golden** | LOW     | LOW     | Pattern exists, just needs adoption                            |
| 22  | **Add goreleaser**                             | LOW     | MED     | Automated releases on tag                                      |
| 23  | **Consider `go:generate stringer` for enums**  | LOW     | LOW     | Type-safe string representation                                |
| 24  | **Modularize into Go workspace**               | LOW     | HIGH    | 10-module go.work — big refactor                               |
| 25  | **Badge info=indigo vs Feedback info=blue**    | TRIVIAL | TRIVIAL | Minor color consistency decision                               |

---

## G) TOP #1 QUESTION

**How should we handle the 4 remaining ID validation panics (modal, drawer, dropdown, accordion)?**

These components require a non-empty `props.ID` to function — the JS uses `document.getElementById(id)` and `getElementById(id + '-panel')`. Without an ID, the component renders but is non-functional (can't close, can't trap focus).

Options:

1. **Auto-generate ID** (e.g., `fmt.Sprintf("tc-modal-%d", time.Now().UnixNano())`) — convenient but unpredictable, hard to reference from HTMX
2. **Silent no-op** — render without JS, just the HTML shell — misleading, looks functional but isn't
3. **Keep panic** (current) — fails fast and loud, but crashes the entire page
4. **Error component** — render a visible error placeholder instead of the component — explicit but ugly in production

I cannot figure out the right tradeoff here without knowing how consumers use these components. **Do your HTMX-powered pages always provide IDs, or do you sometimes rely on auto-generation?**
