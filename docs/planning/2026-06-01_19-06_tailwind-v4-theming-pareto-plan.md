# Pareto Execution Plan — Tailwind v4 Theming & Modern Design System

**Date:** 2026-06-01 19:06 CEST  
**Commit:** `42e22b4`  
**Context:** Post-comprehensive-status-report, post-Tailwind-v4-docs-research  
**Key Insight:** Tailwind v4's `@theme` directive allows overriding default color variables (`--color-blue-500`) which automatically changes ALL `bg-blue-500`, `text-blue-500`, etc. utilities. This means **we don't need to change component code to enable theming**. Consumers override at the CSS level.

---

## Pareto Analysis: The 80/20 Breakdown

### The 1% that delivers 51% of the result

> **Answer the theming question definitively** + ship the CSS file

Tailwind v4's `@theme` system makes the "Go tokens vs CSS tokens" debate moot. The right answer:

- Components emit standard Tailwind classes (`bg-blue-600`) → no breaking changes
- Ship `templ-components-theme.css` with `@theme` overrides → consumers customize via CSS
- Document the approach → consumers know how to theme

### The 4% that delivers 64% of the result

> 1% + `ComponentProps` interface + fix 3 coverage gaps + add dark mode tests to all packages

These 4 items unlock generic component composition, raise quality bar, and prevent regressions.

### The 20% that delivers 80% of the result

> 4% + BDD tests for 4 packages + standardize JS pattern + fix demo app + add missing component features + CI setup

These make the library production-grade, well-tested, and consumable.

### The remaining 80% (nice-to-haves)

> New components (DatePicker, Combobox), visual regression testing, nix flake, documentation site, modularization, etc.

---

## Scoring Formula

**Priority Score = (Impact × CustomerValue) / Effort**

| Dimension     | Scale                                                          |
| ------------- | -------------------------------------------------------------- |
| Impact        | 1–5 (architectural multiplier, ecosystem unlock)               |
| CustomerValue | 1–5 (direct consumer benefit)                                  |
| Effort        | 1–5 (time/complexity: 1=≤10min, 2=≤30min, 3=≤1h, 4=≤2h, 5=≥2h) |

---

## Comprehensive Task List (All 80+ Items)

### Phase 1: The 1% — Theming Foundation (Score: 25.0)

| #   | Task                                                              | Impact | CV  | Effort | Score    | File(s)     | Time |
| --- | ----------------------------------------------------------------- | ------ | --- | ------ | -------- | ----------- | ---- |
| 1.1 | Create `templ-components-theme.css` with `@theme` color overrides | 5      | 5   | 1      | **25.0** | new         | 10m  |
| 1.2 | Add `@custom-variant dark` to theme CSS                           | 5      | 5   | 1      | **25.0** | new         | 5m   |
| 1.3 | Document theming approach in README ("How to customize colors")   | 5      | 5   | 1      | **25.0** | README.md   | 10m  |
| 1.4 | Add theming section to FEATURES.md                                | 4      | 4   | 1      | **16.0** | FEATURES.md | 5m   |

**Phase 1 Deliverable:** Library consumers can change the primary color by overriding `--color-blue-500` in their own `@theme`. No Go code changes needed.

---

### Phase 2: The 4% — Architecture & Quality (Scores: 20.0–15.0)

| #    | Task                                                        | Impact | CV  | Effort | Score    | File(s)             | Time |
| ---- | ----------------------------------------------------------- | ------ | --- | ------ | -------- | ------------------- | ---- |
| 2.1  | Add `ComponentProps` interface with `GetBaseProps()` method | 5      | 4   | 1      | **20.0** | utils/base_props.go | 10m  |
| 2.2  | Add `SetBaseProps()` to interface                           | 4      | 3   | 1      | **12.0** | utils/base_props.go | 5m   |
| 2.3  | Implement interface on all 25+ props structs                | 4      | 4   | 2      | **8.0**  | all packages        | 20m  |
| 2.4  | Fix `fillIcon` test coverage to >70%                        | 3      | 3   | 1      | **9.0**  | display/            | 10m  |
| 2.5  | Fix `Select` test coverage to >70%                          | 3      | 3   | 1      | **9.0**  | forms/              | 10m  |
| 2.6  | Fix `Textarea` test coverage to >70%                        | 3      | 3   | 1      | **9.0**  | forms/              | 10m  |
| 2.7  | Add `TestDarkModeClasses` to `errorpage` package            | 3      | 3   | 1      | **9.0**  | errorpage/          | 10m  |
| 2.8  | Add `TestDarkModeClasses` to `feedback` package             | 3      | 3   | 1      | **9.0**  | feedback/           | 10m  |
| 2.9  | Add `TestDarkModeClasses` to `forms` package                | 3      | 3   | 1      | **9.0**  | forms/              | 10m  |
| 2.10 | Add `TestDarkModeClasses` to `htmx` package                 | 3      | 3   | 1      | **9.0**  | htmx/               | 10m  |
| 2.11 | Add `TestDarkModeClasses` to `layout` package               | 3      | 3   | 1      | **9.0**  | layout/             | 10m  |
| 2.12 | Add `TestDarkModeClasses` to `navigation` package           | 3      | 3   | 1      | **9.0**  | navigation/         | 10m  |
| 2.13 | Add `TestDarkModeClasses` to `icons` package                | 3      | 3   | 1      | **9.0**  | icons/              | 10m  |

**Phase 2 Deliverable:** Generic component composition unlocked; coverage gaps closed; dark mode verified in every package.

---

### Phase 3: The 20% — Production Readiness (Scores: 12.0–6.0)

| #    | Task                                                                 | Impact | CV  | Effort | Score    | File(s)                      | Time |
| ---- | -------------------------------------------------------------------- | ------ | --- | ------ | -------- | ---------------------------- | ---- |
| 3.1  | BDD tests for `navigation` package (Nav, Pagination, Breadcrumbs)    | 4      | 4   | 2      | **8.0**  | navigation/bdd_test.go       | 20m  |
| 3.2  | BDD tests for `htmx` package (Loading, ErrorHandling)                | 4      | 4   | 2      | **8.0**  | htmx/bdd_test.go             | 20m  |
| 3.3  | BDD tests for `layout` package (Base, Minimal, Theme)                | 4      | 4   | 2      | **8.0**  | layout/bdd_test.go           | 20m  |
| 3.4  | BDD tests for `icons` package (render all, unknown fallback)         | 4      | 4   | 2      | **8.0**  | icons/bdd_test.go            | 20m  |
| 3.5  | Standardize JS: extract `tcInit()` shared function in `utils/`       | 4      | 4   | 2      | **8.0**  | utils/ + multiple            | 20m  |
| 3.6  | Convert Dropdown JS to `tcInit()` pattern                            | 4      | 3   | 2      | **6.0**  | display/dropdown.templ       | 15m  |
| 3.7  | Convert Accordion JS to `tcInit()` pattern                           | 4      | 3   | 2      | **6.0**  | display/accordion.templ      | 15m  |
| 3.8  | Convert Modal JS to `tcInit()` pattern                               | 4      | 3   | 2      | **6.0**  | display/modal.templ          | 15m  |
| 3.9  | Convert Dismiss JS to `tcInit()` pattern                             | 4      | 3   | 2      | **6.0**  | feedback/alert.templ         | 15m  |
| 3.10 | Convert Toast/Alert dismiss to `tcInit()` pattern                    | 4      | 3   | 2      | **6.0**  | feedback/toast.templ         | 15m  |
| 3.11 | Convert MobileMenu JS to `tcInit()` pattern                          | 4      | 3   | 2      | **6.0**  | navigation/mobile_menu.templ | 15m  |
| 3.12 | Fix demo app: enable HTMX (`props.HTMXVersion = "2.0.6"`)            | 4      | 4   | 1      | **16.0** | examples/demo/main.go        | 5m   |
| 3.13 | Add more component showcases to demo (DatePicker missing hint, etc.) | 3      | 4   | 2      | **6.0**  | examples/demo/demo.templ     | 20m  |
| 3.14 | Add `Disabled` to `ToggleProps` + propagate to input                 | 3      | 3   | 1      | **9.0**  | forms/toggle.templ           | 10m  |
| 3.15 | Add `AriaLabel` propagation through Toggle label wrapper             | 3      | 3   | 1      | **9.0**  | forms/toggle.templ           | 10m  |
| 3.16 | Add `Multiple` to `FileInputProps`                                   | 3      | 3   | 1      | **9.0**  | forms/file_input.templ       | 10m  |
| 3.17 | Add `DragDrop` styling to FileInput                                  | 3      | 3   | 2      | **4.5**  | forms/file_input.templ       | 20m  |
| 3.18 | Implement Table `<caption>` rendering (field exists, not wired)      | 3      | 3   | 1      | **9.0**  | display/table.templ          | 10m  |
| 3.19 | Add `Href` + click support to Badge                                  | 3      | 3   | 1      | **9.0**  | display/badge.templ          | 10m  |
| 3.20 | Add `Disabled` to `DropdownItem`                                     | 3      | 3   | 1      | **9.0**  | display/dropdown.templ       | 10m  |
| 3.21 | Scale Avatar status dot with size                                    | 3      | 3   | 1      | **9.0**  | display/avatar.templ         | 10m  |
| 3.22 | Add default values to `DefaultRadioProps` (empty struct currently)   | 3      | 3   | 1      | **9.0**  | forms/radio_go.go            | 5m   |
| 3.23 | Add default values to `DefaultRadioGroupProps`                       | 3      | 3   | 1      | **9.0**  | forms/radio_go.go            | 5m   |
| 3.24 | Set up GitHub Actions CI (build + test + lint on push/PR)            | 4      | 4   | 2      | **8.0**  | .github/workflows/ci.yaml    | 20m  |
| 3.25 | Add `go vet` + staticcheck to CI pipeline                            | 3      | 4   | 2      | **6.0**  | .github/workflows/ci.yaml    | 15m  |
| 3.26 | Set coverage threshold in CI (e.g., 60%)                             | 3      | 4   | 1      | **12.0** | .github/workflows/ci.yaml    | 10m  |
| 3.27 | `chmod +x` on `scripts/pre-commit.sh`                                | 2      | 3   | 1      | **6.0**  | scripts/pre-commit.sh        | 2m   |

**Phase 3 Deliverable:** All packages have BDD tests; JS patterns unified; demo app shows full library; CI gate keeps quality.

---

### Phase 4: The Remaining 80% — Expansion & Polish (Scores: 6.0–2.0)

#### New Components

| #    | Task                                                | Impact | CV  | Effort | Score   | File(s)                 | Time |
| ---- | --------------------------------------------------- | ------ | --- | ------ | ------- | ----------------------- | ---- |
| 4.1  | Add Date Picker component                           | 4      | 4   | 5      | **3.2** | forms/date_picker.templ | 2h+  |
| 4.2  | Add Combobox/Autocomplete component                 | 4      | 4   | 5      | **3.2** | forms/combobox.templ    | 2h+  |
| 4.3  | Add Drawer component (Modal slide-out variant)      | 3      | 3   | 4      | **2.3** | display/drawer.templ    | 1.5h |
| 4.4  | Add Form wrapper component (validation integration) | 3      | 3   | 4      | **2.3** | forms/form.templ        | 1.5h |
| 4.5  | Add Skeleton list/paragraph/chart variants          | 3      | 3   | 2      | **4.5** | feedback/loading.templ  | 30m  |
| 4.6  | Add ProgressBar indeterminate state                 | 3      | 3   | 2      | **4.5** | feedback/progress.templ | 30m  |
| 4.7  | Add StepIndicator vertical variant                  | 3      | 3   | 3      | **3.0** | feedback/progress.templ | 1h   |
| 4.8  | Add client-side JS tab switching                    | 3      | 3   | 3      | **3.0** | display/tabs.templ      | 1h   |
| 4.9  | Add Tabs keyboard navigation (arrow keys)           | 3      | 3   | 2      | **4.5** | display/tabs.templ      | 30m  |
| 4.10 | Add more Heroicons (45 → 100)                       | 3      | 3   | 2      | **4.5** | icons/icon_paths.go     | 30m  |

#### Validation & Robustness

| #    | Task                                                               | Impact                   | CV  | Effort | Score   | File(s)                     | Time                |
| ---- | ------------------------------------------------------------------ | ------------------------ | --- | ------ | ------- | --------------------------- | ------------------- | --- |
| 4.11 | Validate SelectOption contradiction (Disabled + Selected)          | 2                        | 2   | 1      | **4.0** | forms/select.templ          | 10m                 |
| 4.12 | Validate SwapOOB swapStyle parameter                               | 2                        | 2   | 1      | **4.0** | htmx/helpers.templ          | 10m                 |
| 4.13 | Validate `                                                         | ` separator in SVG paths | 2   | 2      | 1       | **4.0**                     | icons/icon_paths.go | 10m |
| 4.14 | Replace DropdownItem empty-Href with typed `DropdownItemKind` enum | 3                        | 2   | 3      | **2.0** | display/dropdown.templ      | 1h                  |
| 4.15 | Make GlobalErrorHandling config configurable                       | 2                        | 2   | 2      | **2.0** | htmx/error_handling.templ   | 30m                 |
| 4.16 | Extract error handling magic numbers                               | 2                        | 2   | 1      | **4.0** | htmx/error_handling.templ   | 10m                 |
| 4.17 | SimpleCard should compose through Card internally                  | 2                        | 2   | 2      | **2.0** | display/card.templ          | 30m                 |
| 4.18 | Toast duration configurable per-toast                              | 2                        | 3   | 2      | **3.0** | feedback/toast.templ        | 30m                 |
| 4.19 | Pagination ellipsis for large ranges                               | 2                        | 3   | 2      | **3.0** | navigation/pagination.templ | 30m                 |
| 4.20 | Use `net/url` for pagination URL construction                      | 2                        | 2   | 2      | **2.0** | navigation/pagination.templ | 30m                 |
| 4.21 | Make `PageProps` zero-value safe                                   | 2                        | 2   | 2      | **2.0** | layout/base.templ           | 30m                 |
| 4.22 | Extract theme colors to named constants (`#4f46e5`, `#1e1b4b`)     | 2                        | 2   | 1      | **4.0** | layout/base.templ           | 10m                 |
| 4.23 | Auto-generate `allIconNames()` from `iconPathData`                 | 2                        | 2   | 2      | **2.0** | icons/icon_paths.go         | 30m                 |
| 4.24 | Replace StepIndicator checkmark SVG with icon system               | 2                        | 2   | 1      | **4.0** | feedback/progress.templ     | 10m                 |
| 4.25 | Add `uint` type to Pagination fields                               | 2                        | 2   | 1      | **4.0** | navigation/pagination.templ | 10m                 |

#### Accessibility

| #    | Task                                            | Impact | CV  | Effort | Score   | File(s)                      | Time |
| ---- | ----------------------------------------------- | ------ | --- | ------ | ------- | ---------------------------- | ---- |
| 4.26 | Add `aria-live="polite"` to HTMX error handling | 3      | 3   | 1      | **9.0** | htmx/error_handling.templ    | 10m  |
| 4.27 | Add Table header `scope` attributes             | 3      | 3   | 1      | **9.0** | display/table.templ          | 10m  |
| 4.28 | Add EmptyState landmark role (`role="region"`)  | 2      | 2   | 1      | **4.0** | display/empty_state.templ    | 5m   |
| 4.29 | Add Breadcrumb structured data (JSON-LD)        | 2      | 2   | 2      | **2.0** | navigation/breadcrumbs.templ | 30m  |
| 4.30 | Add Pagination SEO `rel=prev/next`              | 2      | 2   | 1      | **4.0** | navigation/pagination.templ  | 10m  |
| 4.31 | Tooltip JS-based `aria-describedby` injection   | 2      | 2   | 2      | **2.0** | display/tooltip.templ        | 30m  |

#### Testing Infrastructure

| #    | Task                                                        | Impact | CV  | Effort | Score   | File(s)               | Time |
| ---- | ----------------------------------------------------------- | ------ | --- | ------ | ------- | --------------------- | ---- |
| 4.32 | Convert snapshot tests to golden file comparison            | 2      | 2   | 3      | **1.3** | multiple              | 1h   |
| 4.33 | Add benchmark tests for Icon, Card, Table, Nav              | 2      | 2   | 2      | **2.0** | multiple              | 30m  |
| 4.34 | Component composition integration tests                     | 2      | 3   | 3      | **2.0** | display/              | 1h   |
| 4.35 | Full page render integration test (Base+Nav+Content+Footer) | 2      | 3   | 2      | **3.0** | layout/               | 30m  |
| 4.36 | Consistent nonce propagation audit                          | 2      | 2   | 2      | **2.0** | all packages          | 30m  |
| 4.37 | Circular import guard test                                  | 1      | 1   | 1      | **1.0** | docs/status/          | 10m  |
| 4.38 | Accessibility audit automation (axe-core/pa11y)             | 2      | 3   | 4      | **1.5** | docs/status/          | 2h   |
| 4.39 | Move ProgressBar a11y test from display/ to feedback/       | 1      | 1   | 1      | **1.0** | display/a11y_test.go  | 5m   |
| 4.40 | `AssertContainsClass` test helper                           | 2      | 2   | 1      | **4.0** | utils/test_helpers.go | 10m  |
| 4.41 | Extract shared testNavLinks helper                          | 1      | 1   | 1      | **1.0** | navigation/           | 10m  |

#### Documentation

| #    | Task                                                                      | Impact | CV  | Effort | Score   | File(s)                 | Time |
| ---- | ------------------------------------------------------------------------- | ------ | --- | ------ | ------- | ----------------------- | ---- |
| 4.42 | Update README for new API (AvatarStatus, StatCardProps, BreadcrumbsProps) | 3      | 4   | 2      | **6.0** | README.md               | 30m  |
| 4.43 | Update CONTEXT.md with JS pattern docs                                    | 2      | 2   | 1      | **4.0** | CONTEXT.md              | 10m  |
| 4.44 | ADR: filled vs stroke icon convention                                     | 1      | 1   | 1      | **1.0** | docs/adr/               | 15m  |
| 4.45 | ADR: JS attachment patterns (singleton vs IIFE)                           | 2      | 2   | 1      | **4.0** | docs/adr/               | 15m  |
| 4.46 | ADR: FeedbackType unification                                             | 1      | 1   | 1      | **1.0** | docs/adr/               | 10m  |
| 4.47 | Add `go doc` ExampleXxx() functions                                       | 2      | 3   | 2      | **3.0** | multiple                | 30m  |
| 4.48 | Fill DOMAIN_LANGUAGE.md placeholders                                      | 2      | 2   | 1      | **4.0** | docs/DOMAIN_LANGUAGE.md | 15m  |
| 4.49 | Document thread-safety of `utils.Class()`                                 | 1      | 1   | 1      | **1.0** | CONTRIBUTING.md         | 10m  |
| 4.50 | Document 20×20 fill vs 24×24 stroke convention                            | 2      | 2   | 1      | **4.0** | internal/svg/           | 10m  |
| 4.51 | Document why PageProps doesn't embed BaseProps                            | 1      | 1   | 1      | **1.0** | CONTEXT.md              | 10m  |
| 4.52 | Documentation site generation (pkgsite/doc2go)                            | 2      | 3   | 4      | **1.5** | project root            | 2h   |

#### Release & DevOps

| #    | Task                                                     | Impact | CV  | Effort | Score   | File(s)            | Time |
| ---- | -------------------------------------------------------- | ------ | --- | ------ | ------- | ------------------ | ---- |
| 4.53 | Tag v0.2.0 release + update CHANGELOG                    | 4      | 4   | 2      | **8.0** | root               | 30m  |
| 4.54 | Submit to awesome-templ                                  | 3      | 3   | 1      | **9.0** | GitHub PR          | 15m  |
| 4.55 | Open PR on templ.guide to get listed                     | 3      | 3   | 1      | **9.0** | GitHub PR          | 15m  |
| 4.56 | Build real-world example app (clone-and-run)             | 3      | 4   | 5      | **2.4** | cmd/example/       | 3h+  |
| 4.57 | Deploy live component showcase site                      | 3      | 4   | 4      | **3.0** | cmd/demo/          | 2h   |
| 4.58 | Cross-link ecosystem in README                           | 2      | 2   | 1      | **4.0** | README.md          | 10m  |
| 4.59 | Verify `go get` works from clean project                 | 2      | 3   | 1      | **6.0** | docs/status/       | 15m  |
| 4.60 | Add build test for `examples/` in CI                     | 2      | 3   | 1      | **6.0** | .github/workflows/ | 10m  |
| 4.61 | Set up goreleaser for tag-based releases                 | 2      | 3   | 3      | **2.0** | .goreleaser.yml    | 1h   |
| 4.62 | Investigate nix flake for reproducible builds            | 1      | 2   | 4      | **0.5** | flake.nix          | 2h   |
| 4.63 | Investigate visual regression testing                    | 1      | 2   | 4      | **0.5** | docs/status/       | 2h   |
| 4.64 | Modularize into Go workspace (10-module go.work)         | 2      | 2   | 5      | **0.8** | root               | 3h+  |
| 4.65 | Consider `go:generate stringer` for enums                | 1      | 1   | 2      | **0.5** | docs/status/       | 30m  |
| 4.66 | Audit `tailwind-merge-go` thread safety                  | 1      | 1   | 2      | **0.5** | utils/utils.go     | 30m  |
| 4.67 | Plan v1.0 API freeze scope and timeline                  | 2      | 2   | 2      | **2.0** | docs/status/       | 30m  |
| 4.68 | Cross-package circular import audit                      | 1      | 1   | 2      | **0.5** | icons/, feedback/  | 30m  |
| 4.69 | Prune old status reports (keep last 2)                   | 1      | 1   | 1      | **1.0** | docs/status/       | 10m  |
| 4.70 | Investigate gopls QF1003 suppression for generated files | 1      | 1   | 1      | **1.0** | docs/status/       | 10m  |

---

## D2 Execution Graph

```d2
# Pareto Execution Graph: Tailwind v4 Theming & Design System

direction: right

# Phase 1: The 1% (51% of value)
phase1: Phase 1: Theming Foundation {
  style.fill: "#22c55e"
  style.stroke: "#16a34a"

  theme_css: Create templ-components-theme.css
  dark_variant: Add @custom-variant dark
  readme_doc: Document theming in README
  features_doc: Add to FEATURES.md

  theme_css -> dark_variant
  dark_variant -> readme_doc
  dark_variant -> features_doc
}

# Phase 2: The 4% (64% of value)
phase2: Phase 2: Architecture & Quality {
  style.fill: "#3b82f6"
  style.stroke: "#2563eb"

  component_interface: ComponentProps interface
  coverage_fillicon: Fix fillIcon coverage
  coverage_select: Fix Select coverage
  coverage_textarea: Fix Textarea coverage
  dark_tests_all: Dark mode tests (all packages)

  component_interface -> dark_tests_all
  coverage_fillicon -> dark_tests_all
  coverage_select -> dark_tests_all
  coverage_textarea -> dark_tests_all
}

# Phase 3: The 20% (80% of value)
phase3: Phase 3: Production Readiness {
  style.fill: "#f59e0b"
  style.stroke: "#d97706"

  bdd_nav: BDD navigation
  bdd_htmx: BDD htmx
  bdd_layout: BDD layout
  bdd_icons: BDD icons
  js_unify: Unify JS patterns
  demo_fix: Fix demo app
  component_fixes: Component fixes (Toggle, FileInput, Table, Badge, etc.)
  ci_setup: CI setup

  bdd_nav -> js_unify
  bdd_htmx -> js_unify
  bdd_layout -> js_unify
  bdd_icons -> js_unify
  js_unify -> ci_setup
  demo_fix -> ci_setup
  component_fixes -> ci_setup
}

# Phase 4: The rest
phase4: Phase 4: Expansion & Polish {
  style.fill: "#ef4444"
  style.stroke: "#dc2626"

  new_components: New Components
  validation: Validation
  a11y: Accessibility
  testing: Testing
  docs: Documentation
  release: Release

  new_components -> release
  validation -> release
  a11y -> release
  testing -> release
  docs -> release
}

# Dependencies between phases
phase1 -> phase2: {
  style.stroke-dash: 3
}
phase2 -> phase3: {
  style.stroke-dash: 3
}
phase3 -> phase4: {
  style.stroke-dash: 3
}

# Key insight node
insight: "Key Insight:\nTailwind v4 @theme\nallows overriding\ndefault colors.\nNo component code\nchanges needed!"
insight.style.fill: "#a855f7"
insight.style.stroke: "#9333ea"
insight.shape: oval

insight -> phase1
```

---

## Theming Strategy Decision (Answering the #1 Question)

**Question:** Should we use Go compile-time tokens or CSS custom properties?

**Answer: Neither — and both.**

Tailwind v4's `@theme` system provides a third path that makes the old dichotomy obsolete:

1. **Components emit standard Tailwind classes** (`bg-blue-600`, `text-gray-900`)
   - No breaking changes to existing code
   - No new CSS build pipeline required for basic usage
   - Works with Tailwind CDN or any build setup

2. **Ship `templ-components-theme.css`** with `@theme` overrides:

   ```css
   @import "tailwindcss";
   @theme {
     --color-primary-500: #2563eb;
     --color-surface: #ffffff;
     --color-surface-dark: #0f172a;
   }
   @custom-variant dark (&:where(.dark, .dark *));
   ```

3. **Consumers customize by overriding `--color-*` variables**:

   ```css
   @theme {
     --color-blue-500: #4f46e5; /* Changes ALL blue-500 utilities */
   }
   ```

4. **Future (v1.0):** Introduce `bg-tc-primary` semantic tokens alongside `bg-blue-600`:
   ```css
   @theme {
     --color-tc-primary: var(--color-blue-600);
     --color-tc-surface: var(--color-white);
   }
   ```

This gives us:

- ✅ **Ease of use**: `go get` → works immediately (no CSS required)
- ✅ **Customizability**: Override any color via `@theme`
- ✅ **Future-proofing**: Native Tailwind v4, no custom solutions
- ✅ **Low maintenance**: One optional CSS file, not two token systems

---

## Execution Order Summary

```
Phase 1 (Theming Foundation)        → 4 tasks  → ~30 min  → 51% value
Phase 2 (Architecture & Quality)  → 13 tasks → ~3h     → 64% value
Phase 3 (Production Readiness)     → 27 tasks → ~8h     → 80% value
Phase 4 (Expansion & Polish)       → 70 tasks → ~30h+   → 100% value
```

**Total estimated time:** ~42 hours of focused work  
**Total tasks:** 114 (broken into ≤12-minute chunks where possible)

---

## Next Action

**WAITING FOR INSTRUCTIONS.**

Pick one of:

1. **"Execute Phase 1"** — Ship the theming CSS file + docs (~30 min)
2. **"Execute Phase 2"** — Add ComponentProps interface + fix coverage + dark tests (~3h)
3. **"Execute Phase 3"** — BDD tests + JS unification + demo fixes (~8h)
4. **"Execute all"** — Full execution mode, all phases (~42h)
5. **"Reprioritize"** — Re-sort based on specific customer needs

_Generated: 2026-06-01 19:06 CEST_  
_All tests passing. Lint clean. Working tree clean._
