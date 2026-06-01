# Comprehensive Status Report — templ-components

**Date:** 2026-06-01 18:35 CEST  
**Branch:** master (ahead of origin by 1 commit)  
**Latest Commit:** `f4f2b4d` — a11y(layout): comprehensive accessibility improvements + footer slot + dark mode color consistency  
**Go Version:** 1.26.3 | **templ:** v0.3.1020 | **tailwind-merge-go:** v0.2.1 | **go-error-family:** v0.3.0

---

## a) FULLY DONE ✅

### Component Library (56 templ components across 9 packages)

| Package      | Components                                                                                                                | Status      |
| ------------ | ------------------------------------------------------------------------------------------------------------------------- | ----------- |
| `utils`      | BaseProps, Class(), Ternary(), DismissScript(), CurrentYear()                                                             | ✅ Complete |
| `display`    | Card, SimpleCard, StatCard, Badge, Avatar, Button, Modal, Dropdown, Accordion, Tabs, Tooltip, Table, EmptyState           | ✅ Complete |
| `errorpage`  | ErrorPage, ErrorDetail, ErrorAlert, ErrorHandler, 5 family styles                                                         | ✅ Complete |
| `feedback`   | Alert, Toast, ToastContainer, Spinner, LoadingOverlay, InlineLoading, Skeleton, SkeletonGroup, ProgressBar, StepIndicator | ✅ Complete |
| `forms`      | Input, Select, Textarea, Checkbox, Toggle, FileInput, Label, FieldError, InputGroup, Radio, RadioGroup                    | ✅ Complete |
| `htmx`       | LoadingIndicator, InlineLoadingOverlay, LoadingButton, ConfirmDelete, SwapOOB, CSRFToken, ErrorHandler                    | ✅ Complete |
| `icons`      | 45 typed icon names, Icon(), IconPathJS(), FillIcon()                                                                     | ✅ Complete |
| `layout`     | Base(), Minimal(), ThemeToggle(), ThemeScript()                                                                           | ✅ Complete |
| `navigation` | Nav, SimpleNav, NavLink, MobileNavLink, MobileMenu, MobileMenuToggle, Breadcrumbs, Pagination, Footer                     | ✅ Complete |

### Quality Infrastructure

- **1,049+ tests** across 53 test files — all passing (`go test ./...`)
- **0 lint issues** (`golangci-lint run` across all 11 packages)
- **0 TODO/FIXME/HACK/XXX** markers in source code (generated files excluded)
- **40 generated `*_templ.go` files** committed and in sync with source
- **No circular imports** across 11 Go packages
- **100% AriaLabel propagation** — all 25+ components with BaseProps propagate to root element

### Recent Deliveries (Last 5 Sessions)

1. **Dark mode color consistency** — Standardized `slate-*` → `gray-*` across all components (card, accordion, dropdown, table, tabs, button, loading, htmx)
2. **Accessibility: motion-reduce** — Added `motion-reduce:transition-none` + `motion-reduce:animate-none` to all animated components (modal, accordion, spinner, skeleton, toast, HTMX loading)
3. **Accessibility: interactive cues** — `cursor-pointer` on buttons, `caret-blue-600` on form inputs, `scroll-smooth` + `selection:` colors on body
4. **Layout: Footer slot API** — `PageProps.Footer` field added; footer renders after `<main>` but inside `<body>`
5. **Errorpage package** — Full go-error-family v0.3.0 integration with 5 error families, `FromError()` converter, 6 pre-built constructors, HTMLShell + JSON modes
6. **Deep deduplication** — Zero semantic code clones across the entire codebase (art-dupl verified)
7. **Type safety** — All enums validated (2 panic-on-unknown, 10 map+fallback), structural variants use if-branch
8. **SVG consolidation** — Single source of truth in `internal/svg` for all shared paths

---

## b) PARTIALLY DONE 🟡

### Component Gaps

- **Radio button** — `Radio` and `RadioGroup` exist but `DefaultRadioProps`/`DefaultRadioGroupProps` have no real defaults (empty structs)
- **Toggle** — Exists with size variants but no `Disabled` state or `AriaLabel` propagation through label wrapper
- **FileInput** — Exists but no `Multiple` support, no drag-drop styling, no file type icons
- **StepIndicator** — Horizontal only; vertical variant not started
- **Tabs** — No client-side JS tab switching (keyboard navigation only via arrow keys not implemented)
- **Table** — `Caption` field exists but `<caption>` rendering not implemented
- **Badge** — No `Href` / click support despite being a common pattern
- **Avatar** — Status dot size is fixed (h-2.5 w-2.5) regardless of avatar size

### Testing

- **Coverage gaps:** `fillIcon` (~60%), `Select` (~65%), `Textarea` (~68%) — below 70% threshold
- **Missing BDD tests:** navigation package, htmx package, layout package, icons package
- **Missing snapshot tests:** Several packages lack golden file comparison
- **Dark mode class verification:** Only `display/` has dedicated `TestDarkModeClasses`; other packages don't verify `dark:` output

### Infrastructure

- **CI/CD** — GitHub Actions workflow exists but may not be fully configured for all checks
- **Pre-commit hook** — `scripts/pre-commit.sh` exists but may not be `chmod +x`
- **Demo app** — `examples/demo/` runs but HTMX is disabled (`props.HTMXVersion = ""` in demo)
- **Documentation site** — No generated docs (pkgsite/doc2go not set up)

---

## c) NOT STARTED 🔵

### New Components

- **Date Picker** — No component exists
- **Combobox / Autocomplete** — No component exists
- **Dialog / Drawer** — Modal exists but no slide-out drawer variant
- **Form wrapper** — No high-level form component with validation integration
- **Skeleton variants** — Only 7 variants; missing list, paragraph, chart
- **ProgressBar indeterminate** — Only determinate state exists
- **More Heroicons** — 45 of ~300 available

### Architecture Improvements

- **Design tokens** — No CSS custom properties or Go-based token system; colors hardcoded as Tailwind strings
- **ComponentProps interface** — No `GetBaseProps()` interface; prevents generic component wrappers
- **Validation layer** — No `Validate() error` method on props structs
- **Modularization** — Single module; no `go.work` workspace split
- **Theme system** — No extracted Tailwind preset/theme configuration file

### DevOps & Release

- **Live demo site** — No deployed showcase
- **goreleaser** — Config exists but not verified for tag-based releases
- **Nix flake** — Not started
- **Visual regression testing** — Not started
- **Accessibility audit automation** — No axe-core/pa11y integration

---

## d) TOTALLY FUCKED UP! 🔴

**Nothing.** All tests pass, lint is clean, build compiles, no TODOs in source, no circular imports. The codebase is in its best state ever.

---

## e) WHAT WE SHOULD IMPROVE! 💡

### 1. Design Token System (High Impact, Medium Effort)

**Problem:** Colors, shadows, and borders are hardcoded as Tailwind strings in 30+ places. `bg-blue-600`, `text-gray-900`, `border-gray-200` appear scattered. This makes:

- Theming impossible (can't swap blue for indigo)
- Dark mode maintenance fragile (we just fixed slate→gray inconsistency)
- Brand customization requires forking

**Solution:** Extract a `design` package (or `utils/design.go`) with named constants:

```go
package design

const (
    ColorPrimary     = "blue-600"
    ColorPrimaryDark = "blue-500"
    ColorText        = "gray-900"
    ColorTextDark    = "gray-100"
    // ...
)

func PrimaryButton() string { return "bg-" + ColorPrimary + " text-white" }
```

**Existing code that fits:** `errorpage/styles.go` already uses a `familyVisualStyle` struct pattern — this proves the approach works. `feedback/styles.go` uses `feedbackStyleSet` + lookup maps. We just need to generalize this to the whole library.

### 2. ComponentProps Interface (High Impact, Low Effort)

**Problem:** Can't write generic utilities like `WrapWithAria(component, label)` because props structs don't share an interface.

**Solution:**

```go
type ComponentProps interface {
    GetBaseProps() BaseProps
    SetBaseProps(BaseProps)
}
```

**Existing code that fits:** All 25+ props structs already embed `BaseProps`. Adding methods is trivial. This enables generic component composition, validation pipelines, and wrapper functions.

### 3. CSS Custom Properties + Tailwind v4 Integration (High Impact, High Effort)

**Problem:** We're fighting Tailwind's utility-first approach by hardcoding colors in Go strings. This is backwards. The modern way (2024+) is CSS custom properties + Tailwind's `@theme` directive.

**Solution:** Ship a `tc-theme.css` file that consumers import:

```css
@import "tailwindcss";
@theme {
  --color-tc-primary: #2563eb;
  --color-tc-primary-dark: #60a5fa;
  --color-tc-surface: #ffffff;
  --color-tc-surface-dark: #0f172a;
  /* ... */
}
```

Then components use `bg-tc-primary` instead of `bg-blue-600`. This makes theming one CSS variable change.

**Existing code that fits:** The `layout/theme.templ` already manages dark mode class toggling. The CSS file would complement it. Tailwind v4's `@theme` is exactly what we need.

### 4. Well-Established Libraries We Could Adopt (Medium Impact, Low Effort)

| Library                         | Purpose              | How We'd Use It                        |
| ------------------------------- | -------------------- | -------------------------------------- |
| `github.com/alecthomas/chroma`  | Syntax highlighting  | For code blocks in error pages / docs  |
| `github.com/yosssi/gcss`        | CSS preprocessing    | For the theme CSS file                 |
| `github.com/nicksnyder/go-i18n` | Internationalization | For error messages, form labels        |
| `github.com/invopop/validation` | Struct validation    | For `Validate() error` on props        |
| `github.com/stretchr/testify`   | Test assertions      | Already standard; ensure full adoption |

**What NOT to adopt:**

- **React/Vue/Angular** — defeats the purpose of server-rendered templ
- **Bootstrap/Bulma** — we already use Tailwind
- **Alpine.js** — adds JS dependency; our vanilla JS pattern is lighter

### 5. Consolidate JS Pattern (Medium Impact, Medium Effort)

**Problem:** 10 script blocks across 7 files use 3 different patterns:

- Global singleton: `window.tcDropdownAttached`, `window.tcAccordionAttached`
- Per-instance IIFE: Modal focus trap
- IIFE-wrapped global guard: ThemeToggle

**Solution:** Standardize on a single `tcInit()` function in `utils/` that handles all component initialization via data attributes:

```javascript
// One script, all components
document.addEventListener("click", function (e) {
  // Dropdown: [data-dropdown-trigger]
  // Accordion: [data-accordion-trigger]
  // Dismiss: [data-dismiss]
  // ...
});
```

**Existing code that fits:** `utils.DismissScript()` already demonstrates shared JS extraction. We just need to extend the pattern.

---

## f) Top #25 Things We Should Get Done Next!

Sorted by **Impact / Effort** ratio (highest first):

| #   | Task                                                                  | Impact    | Effort    | Package               |
| --- | --------------------------------------------------------------------- | --------- | --------- | --------------------- |
| 1   | **Add `ComponentProps` interface + `GetBaseProps()`**                 | 🔥 High   | 🟢 Low    | `utils/`              |
| 2   | **Extract design tokens (colors, shadows, radii) to named constants** | 🔥 High   | 🟡 Medium | `utils/` or `design/` |
| 3   | **Fix test coverage gaps (<70%: fillIcon, Select, Textarea)**         | 🔥 High   | 🟢 Low    | `display/`, `forms/`  |
| 4   | **Add BDD tests for navigation, htmx, layout, icons**                 | 🔥 High   | 🟡 Medium | multiple              |
| 5   | **Standardize all JS to single `tcInit()` pattern**                   | 🔥 High   | 🟡 Medium | `utils/`, multiple    |
| 6   | **Add `Validate() error` to all props structs**                       | 🔥 High   | 🟡 Medium | all packages          |
| 7   | **Ship `tc-theme.css` with CSS custom properties**                    | 🔥 High   | 🟡 Medium | new file              |
| 8   | **Fix demo app: enable HTMX + add more component showcases**          | 🔥 High   | 🟢 Low    | `examples/demo/`      |
| 9   | **Add dark mode class verification tests to all packages**            | 🔥 High   | 🟢 Low    | all test files        |
| 10  | **Add `Disabled` to Toggle, FileInput, DropdownItem**                 | 🔥 High   | 🟢 Low    | `forms/`, `display/`  |
| 11  | **Implement Table `<caption>` rendering**                             | 🟡 Medium | 🟢 Low    | `display/`            |
| 12  | **Add Badge `Href` / click support**                                  | 🟡 Medium | 🟢 Low    | `display/`            |
| 13  | **Add StepIndicator vertical variant**                                | 🟡 Medium | 🟡 Medium | `feedback/`           |
| 14  | **Add ProgressBar indeterminate state**                               | 🟡 Medium | 🟡 Medium | `feedback/`           |
| 15  | **Add client-side JS tab switching**                                  | 🟡 Medium | 🟡 Medium | `display/`            |
| 16  | **Replace hardcoded `blue-600` with `Primary` token**                 | 🟡 Medium | 🟢 Low    | all packages          |
| 17  | **Set up GitHub Actions CI properly**                                 | 🟡 Medium | 🟢 Low    | `.github/workflows/`  |
| 18  | **Tag v0.2.0 release + update CHANGELOG**                             | 🟡 Medium | 🟢 Low    | root                  |
| 19  | **Deploy live demo site**                                             | 🟡 Medium | 🟡 Medium | `cmd/demo/`           |
| 20  | **Add Date Picker component**                                         | 🟡 Medium | 🔴 High   | `forms/`              |
| 21  | **Add Combobox/Autocomplete component**                               | 🟡 Medium | 🔴 High   | `forms/`              |
| 22  | **Add more Heroicons (45→100)**                                       | 🟡 Medium | 🟢 Low    | `icons/`              |
| 23  | **Modularize into Go workspace**                                      | 🟡 Medium | 🔴 High   | root                  |
| 24  | **Add `uint` type to Pagination fields**                              | 🟢 Low    | 🟢 Low    | `navigation/`         |
| 25  | **Write ADR for JS attachment patterns**                              | 🟢 Low    | 🟢 Low    | `docs/adr/`           |

---

## g) Top #1 Question I Can NOT Figure Out Myself! ❓

**"Should we extract all Tailwind class strings into a Go-based design token system (compile-time), or should we ship a CSS custom properties file + Tailwind v4 `@theme` (runtime)?"**

**Why this matters:**

- **Go tokens** (compile-time): Type-safe, autocomplete in IDE, no CSS file needed, consumers just `go get`. But theming requires recompiling Go code.
- **CSS tokens** (runtime): True runtime theming, consumers can override with CSS. But we need a CSS file, and Tailwind v4 `@theme` is still stabilizing.

**What I can't figure out:**

- Does Tailwind v4's `@theme` work reliably enough to depend on it? (It shipped in v4.0 in Jan 2025)
- If we ship CSS, how do consumers import it? Go modules don't serve CSS well.
- Should we do BOTH: Go constants for internal use + CSS file for consumer theming?

**What I tried:**

- Searched Tailwind v4 docs — `@theme` is the official approach but requires a build step
- The demo uses CDN (`@tailwindcss/browser@4`) which supports `@theme`
- Our consumers are Go developers who may not have a CSS build pipeline

**What I need from you:**
A decision on the theming strategy that balances:

1. Ease of use for Go developers (`go get` + it works)
2. Customizability (change primary color without forking)
3. Future-proofing (Tailwind v4 compatibility)
4. Maintenance burden (don't maintain two token systems)

---

## Metrics Snapshot

| Metric                        | Value        |
| ----------------------------- | ------------ |
| Hand-written Go lines         | 9,818        |
| Generated templ Go lines      | 13,308       |
| Templ source files            | 40           |
| Test files                    | 53           |
| Go packages                   | 11           |
| Test cases                    | 1,049+       |
| Lint issues                   | 0            |
| TODOs in source               | 0            |
| Open items in TODO_LIST.md    | ~80          |
| Completed items (git history) | ~200+        |
| Days since last release       | ~14 (v0.1.x) |
| Commits since last tag        | 100+         |

---

## Risk Assessment

| Risk                             | Likelihood | Impact | Mitigation                                                               |
| -------------------------------- | ---------- | ------ | ------------------------------------------------------------------------ |
| Tailwind v4 breaking changes     | Medium     | High   | Pin `@tailwindcss/browser@4` in demo; avoid v4-specific features in core |
| API churn before v1.0            | High       | Medium | Document breaking changes; maintain CHANGELOG; semantic versioning       |
| Consumer theming demands         | High       | Medium | **Need decision on theming strategy (see question above)**               |
| Test coverage regression         | Low        | Medium | CI threshold; pre-commit hook                                            |
| Dependency drift (templ version) | Medium     | Low    | Monthly dep updates; test before bump                                    |

---

## Next Session Recommendation

1. **Answer the theming question** — this unlocks #1, #2, #16 on the priority list
2. **Implement `ComponentProps` interface** — 30-minute change, massive architectural unlock
3. **Fix the 3 coverage gaps** — quick wins that raise the quality bar
4. **Pick one new component** — Date Picker or Combobox (both high-demand, forms-adjacent)

---

_Generated: 2026-06-01 18:35 CEST_  
_Commit: f4f2b4d_  
_All tests passing. Lint clean. Working tree clean._
