# Status Report — templ-components

**Date:** 2026-06-01 21:14 | **Branch:** master | **Ahead of origin:** 3 commits | **Working tree:** clean

---

## Metrics

| Metric | Value |
|--------|-------|
| Hand-written Go LOC | 1,380 |
| Generated `*_templ.go` LOC | 13,308 |
| `.templ` files | 40 |
| Test files | 56 |
| Passing tests | 1,095 |
| Total coverage | 66.1% |
| Packages | 11 (10 + demo) |
| Components | 56 |
| Dependencies | 3 (templ, tailwind-merge-go, go-error-family) |
| Lint issues | 0 |
| TODOs in source | 0 |
| Panic calls (hand-written) | 1 (icons name validation) |

### Per-Package Coverage

| Package | Coverage | Tests |
|---------|----------|-------|
| `utils` | 81.5% | ~30 |
| `internal/svg` | 79.0% | ~10 |
| `htmx` | 77.3% | ~60 |
| `icons` | 75.0% | ~40 |
| `layout` | 73.1% | ~35 |
| `navigation` | 72.1% | ~50 |
| `errorpage` | 70.6% | ~120 |
| `feedback` | 70.2% | ~70 |
| `display` | 68.6% | ~200 |
| `forms` | 66.9% | ~80 |
| `examples/demo` | 0.0% | 0 |

---

## a) FULLY DONE ✅

### Core Library (All 56 Components)
- **display/ (14):** Accordion, Avatar, Badge, StatusBadge, Card, SimpleCard, StatCard, Dropdown, EmptyState, SimpleEmptyState, Modal, Table, Tabs, Tooltip
- **feedback/ (12):** Alert, Toast, ToastContainer, Spinner, ProgressBar, StepIndicator, 6 Skeleton variants, LoadingOverlay
- **forms/ (6):** Input, Select, Textarea, Checkbox, Radio/RadioGroup, Toggle, FileInput, InputGroup
- **navigation/ (9):** Nav, NavLink, MobileNavLink, Breadcrumbs, Pagination, StepIndicator
- **errorpage/ (3):** ErrorPage, ErrorDetail, ErrorAlert
- **layout/ (4):** Base, Minimal, ThemeScript, ThemeToggle
- **htmx/ (7):** LoadingIndicator, InlineLoadingOverlay, LoadingOverlay, LoadingButton, ErrorHandler, Retry, Swap
- **icons/ (1 + 42 icons):** Icon component with typed Name enum
- **utils/:** BaseProps, ComponentProps interface, Class, Ternary, MapEnum, DismissScript, test helpers

### Architecture & Quality
- All 33+ props structs embed `utils.BaseProps` and auto-satisfy `ComponentProps` interface
- All root elements propagate Class, Attrs, ID, AriaLabel from BaseProps (25/25 verified)
- Zero circular imports in dependency graph
- Consistent pattern: map-based lookups for styles, if-branch for structural variants
- CSP-compliant: all inline scripts use `nonce={ props.Nonce }`

### Accessibility
- `motion-reduce:transition-none` / `motion-reduce:animate-none` on all animated components
- `aria-expanded`, `aria-controls`, `aria-selected`, `aria-modal`, `role="dialog"`, etc. where appropriate
- Skip link in Base layout
- `aria-invalid` / `aria-describedby` on form inputs with errors
- Dark mode class tests in 5 packages (errorpage, feedback, forms, icons, navigation)

### Modern Design Defaults
- Dark mode colors: all `gray-*` (no mixed `slate-*`)
- Interactive: `cursor-pointer` on buttons, `caret-blue-600` on inputs, `scroll-smooth` on body
- Card shadow: `shadow-sm` (upgraded from `shadow-xs`)
- Selection colors: blue-tinted for both light and dark

### Theming
- Tailwind v4 `@theme` strategy documented and shipped (`templ-components-theme.css`)
- Semantic aliases: `--color-tc-primary`, `--color-tc-surface`, etc.
- `@custom-variant dark` for class-based dark mode
- README updated with theming section

### Infrastructure
- **CI** (`.github/workflows/ci.yaml`): lint + build/test with coverage artifact upload
- **Pre-commit** (`scripts/pre-commit.sh`): templ generate → build → test → lint
- **Lint** (`.golangci.yml`): 30+ linters configured

### Documentation
- README with theming section, metrics, installation
- FEATURES.md with full component inventory
- AGENTS.md with architecture, conventions, breaking changes
- Multiple planning/status docs in `docs/`

---

## b) PARTIALLY DONE ⚠️

### Test Coverage Gaps
- **0% coverage functions:** `DefaultRadioProps()`, `DefaultRadioGroupProps()`, `DefaultInputGroupProps()` — simple default constructors, trivially testable
- **<50% coverage:** `dropdownItemLink` (44.7%) — link variant of dropdown items
- **Untested edge cases in forms:**
  - Toggle: disabled state not tested, aria-checked not verified
  - Radio: Disabled option, error state, DefaultRadioProps/DefaultRadioGroupProps
  - FileInput: disabled state, error state, DefaultFileInputProps
  - InputGroup: right addon, both addons simultaneously
- **Accessibility testing:**
  - No keyboard navigation tests anywhere (Escape for Modal/Dropdown, Tab order)
  - `forms/a11y_test.go` only tests dark mode classes — no `aria-*` attribute tests
  - No `aria-checked` test for Toggle
  - No `role="radiogroup"` test for RadioGroup

### Deprecated Type Aliases
- `AlertType = FeedbackType` and `ToastType = FeedbackType` still exported in `feedback/`
- Documented as deprecated but still in public API
- Should plan removal timeline (v0.3 or v1.0?)

### Lookup Map Type Safety
- 7 lookup maps use `map[string]string` instead of typed enum keys
- `badgeSizeLookup`, `avatarSizeLookup`, `avatarDotSizeLookup`, `spinnerSizeLookup`, `progressHeightLookup`, `cardPaddingLookup`, `modalSizeLookup`
- Each requires `string(size)` conversion — unnecessary if key is the typed enum

### `errorpage/handler.go` Size
- 375 lines — the only file over 300 lines
- Could split into `handler.go` (core) + `constructors.go` (6 pre-built constructors)
- 6 `_, _ = fmt.Fprint(bw, ...)` ignored errors (safe: writing to `bytes.Buffer`)

---

## c) NOT STARTED 📋

### New Components
- **DatePicker** — no code, no types defined
- **Combobox** — no code, no types defined
- **Drawer** (slide-out panel) — no code, no types defined
- **Form wrapper** (form-level validation, submit handling) — no code

### Quality
- **BDD tests** for display, feedback, forms, navigation packages (htmx and errorpage already have BDD)
- **Keyboard navigation tests** — zero coverage across all components
- **htmlEscape()** unit test — 50% coverage, only exercised indirectly

### Infrastructure
- **Coverage threshold** in CI — no minimum enforced
- **Documentation site** — no dedicated site or auto-generated API docs
- **v0.2.0 release tag** — not cut yet

---

## d) TOTALLY FUCKED UP 💀

**Nothing.** Zero critical issues. Zero broken tests. Zero lint errors. Zero TODOs.

The closest things to "problems":
- `DismissScript()` shows 0% coverage despite having a test — likely a coverage profiling artifact for a single-return function
- templ version mismatch warning: local v0.3.1036 vs go.mod v0.3.1020 (harmless: local is a dev build not yet on proxy)

---

## e) WHAT WE SHOULD IMPROVE

### Architecture & Type Safety

1. **`ButtonProps.Type` is raw `string`** — should be typed enum (`ButtonType` exists for visual variants, need `HTMLButtonType` for submit/reset/button)
2. **`ProgressBarProps.Color` is raw `string`** — could be typed with presets
3. **`StatusBadge(status string)` accepts any string** — could accept typed values
4. **7 lookup maps use `map[string]string`** — should use typed enum keys (`map[BadgeSize]string`)
5. **Mixed receivers on `PaginationProps`** — `normalize()` is pointer, `pageURL()` is value (cosmetic but inconsistent)
6. **`examples/demo/` has 0% coverage** — acceptable for a demo, but worth noting

### Potential Library Improvements

7. **`go-error-family` coupling** — only used by `errorpage/`, all consumers pull it in. Could be an optional sub-module or interface-based decoupling.
8. **`"Go back"` duplicated** in error constructors (lines 154 and 170 of `handler.go`) — extract to constant like existing `msgGoHome`
9. **`errorpage/handler.go` at 375 lines** — split constructors to separate file
10. **No `ComponentProps` value receiver** — only `*BaseProps` satisfies, document or provide alternative

### Testing Philosophy
11. **No keyboard interaction tests** — this is the biggest a11y gap
12. **forms/ aria tests missing** — `a11y_test.go` only tests dark mode classes
13. **`htmlEscape` not directly tested** — security-relevant function

---

## f) Top 25 Next Tasks (Sorted by Impact/Effort)

### Tier 1: High Impact, Low Effort (< 30 min each)

| # | Task | Impact | Effort | Rationale |
|---|------|--------|--------|-----------|
| 1 | Test `DefaultRadioProps()`, `DefaultRadioGroupProps()`, `DefaultInputGroupProps()` | Coverage | 10min | Three 0% functions, trivial tests |
| 2 | Fix `dropdownItemLink` coverage (44.7% → 80%+) | Coverage | 15min | Only function <50% outside demo |
| 3 | Extract `"Go back"` to constant in `errorpage/handler.go` | Cleanliness | 5min | Duplicated string |
| 4 | Split `errorpage/handler.go` → `handler.go` + `constructors.go` | Maintainability | 10min | Only file >300 lines |
| 5 | Typed enum keys for 7 lookup maps (`map[BadgeSize]string`) | Type safety | 20min | Eliminates `string()` casts |
| 6 | Test Toggle disabled state + aria-checked | A11y/Coverage | 15min | New feature untested |
| 7 | Test Radio disabled option + DefaultRadioProps | A11y/Coverage | 15min | Missing edge cases |
| 8 | Test FileInput disabled + DefaultFileInputProps | Coverage | 10min | Missing edge cases |

### Tier 2: High Impact, Medium Effort (1-2 hours each)

| # | Task | Impact | Effort | Rationale |
|---|------|--------|--------|-----------|
| 9 | Add `HTMLButtonType` enum for `ButtonProps.Type` | Type safety | 30min | Currently raw string, accepts anything |
| 10 | Add `htmlEscape()` direct unit tests | Security | 30min | XSS-relevant function at 50% coverage |
| 11 | Add forms/ `aria-*` attribute tests (aria-checked, role="radiogroup", aria-invalid) | A11y | 1hr | forms/ a11y_test.go only tests dark mode |
| 12 | BDD tests for `forms/` package | Quality | 1.5hr | Behavioral test coverage |
| 13 | BDD tests for `navigation/` package | Quality | 1hr | Behavioral test coverage |
| 14 | Add CI coverage threshold (e.g., 60% minimum) | Quality gate | 15min | Prevent regression |

### Tier 3: Medium Impact, Medium Effort (2-4 hours each)

| # | Task | Impact | Effort | Rationale |
|---|------|--------|--------|-----------|
| 15 | BDD tests for `display/` package | Quality | 2hr | Largest package by component count |
| 16 | BDD tests for `feedback/` package | Quality | 1.5hr | Behavioral test coverage |
| 17 | Keyboard navigation tests (Modal Escape, Dropdown Escape, Tab order) | A11y | 3hr | Biggest a11y gap |
| 18 | `ProgressBarProps.Color` typed enum with presets | Type safety | 1hr | Currently raw string |

### Tier 4: High Impact, High Effort (New Features)

| # | Task | Impact | Effort | Rationale |
|---|------|--------|--------|-----------|
| 19 | **Drawer component** (slide-out panel) | Feature | 4hr | Common UI pattern, Modal-like architecture |
| 20 | **Form wrapper** (form-level validation, submit) | Feature | 6hr | High-value for consumers, ties forms together |
| 21 | **DatePicker component** | Feature | 8hr | Complex, requires date handling library |
| 22 | **Combobox component** | Feature | 6hr | Complex, requires search + keyboard nav |

### Tier 5: Infrastructure & Release

| # | Task | Impact | Effort | Rationale |
|---|------|--------|--------|-----------|
| 23 | Documentation site (auto-generated API docs) | Adoption | 4hr | Helps consumers discover components |
| 24 | Remove deprecated `AlertType`/`ToastType` aliases | API cleanup | 30min | Clean public API for v0.2.0 |
| 25 | Cut v0.2.0 release tag | Milestone | 1hr | First public version with breaking changes |

---

## g) Top #1 Question I Cannot Answer Myself

**Should `go-error-family` be decoupled from `errorpage/` as an optional dependency?**

Currently `errorpage/handler.go` imports `go-error-family` directly. This means every consumer of `templ-components` pulls in `go-error-family` even if they never use error pages. Options:

1. **Keep as-is** — only 3 deps total, go-error-family is lightweight, tight coupling is acceptable
2. **Interface-based decoupling** — define `ErrorClassifier` interface in errorpage, let consumers optionally pass go-error-family adapter
3. **Sub-module split** — `errorpage/` becomes its own Go module with its own `go.mod`

The tradeoff: options 2 and 3 add complexity and indirection for a dependency that's small and same-author. But it violates the "no framework deps" principle stated in our architecture docs. **What's your preference?**

---

## Session History

This report covers work across multiple sessions starting from 2026-04-27. Key milestones:

1. **Alpha development** (Apr-May): All 56 components built, tested, documented
2. **Go 1.26 migration** (May 17): Adopted Go 1.26 features, removed `utils.Ptr`
3. **Deep deduplication** (May 21): Zero clone detection issues, consolidated SVG paths
4. **errorpage package** (May 29): Full go-error-family integration, HTTP handler, JSON/HTML modes
5. **Modern design audit** (Jun 1): Dark mode consistency, motion-reduce, interactive polish
6. **Theming foundation** (Jun 1): Tailwind v4 @theme strategy, ComponentProps interface
7. **This session** (Jun 1): Comprehensive review, status report, execution plan
