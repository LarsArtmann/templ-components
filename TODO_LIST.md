# TODO List

**Updated:** 2026-06-08
**Coverage:** 70%+ | **Tests:** 1100+ | **Packages:** 10+demo | **Lint:** 0 issues | **Icons:** 75

> Session 3: 8 commits — Badge href, ProgressBar indeterminate, StepIndicator vertical, 30 new icons, JS tab switching, Form component, benchmarks, thread-safety fix.

---

## 🔴 Bugs & Security

- [x] Fix JS re-attachment after HTMX DOM swaps — **by design**: all components use document-level event delegation
- [x] Fix ThemeToggle multi-instance bug — **not a bug**: event delegation handles all instances
- [x] Fix Accordion JS to IIFE-per-instance pattern — **not needed**: event delegation handles dynamic accordions
- [x] Fix Checkbox unconditional `id=""` → conditional `if props.ID != ""` rendering
- [x] Fix pre-commit hook — uses scripts/pre-commit.sh
- [x] Fix tailwind-merge-go thread safety — `sync.Mutex` in `utils.Class()` is REQUIRED (source: utils/utils.go)

## 🟡 Breaking Changes (defer to v1.0)

- [ ] Move test helpers to `internal/testutil/` — breaking change for external consumers (source: utils/test_helpers.go)
- [ ] Spinner BaseProps conversion — `Spinner(size, colorClass)` positional args → `SpinnerProps` struct (source: feedback/loading.templ)
- [ ] SimpleNav BaseProps conversion — `(brandText, brandHref, links, currentPath)` → `SimpleNavProps` struct (source: navigation/nav.templ)
- [ ] Add BaseProps to StepIndicatorProps (source: feedback/progress.templ)
- [ ] Pagination uint fields — `CurrentPage` and `TotalPages` should be `uint` (source: navigation/pagination.templ)

## 🟢 Enhancements

### Validation & Robustness

- [x] Add TrendDirection consistent validation — `normalizeTrend()` in display/card.templ
- [x] Extract hardcoded SVG path strings to named constants — `internal/svg`
- [x] Propagate AriaLabel from BaseProps in all 25 components
- [x] Extract shared dismiss JS — `dismissScript()` in feedback/styles.go
- [x] Unify alertIconName/toastIconName — shared `feedbackIconName()`
- [x] Fix toast JS string builder — `fmt.Fprintf` replaces string concatenation
- [x] Validate SelectOption contradiction (Disabled + Selected simultaneously)
- [x] Validate `|` separator doesn't appear in SVG paths — `iconPaths()` panics on empty segments
- [x] Add `DropdownItem.Disabled` field
- [x] Add `InputProps.MaxLength` field
- [x] Add `TextareaProps.MaxLength` field
- [x] Add `CheckboxProps.Value` field
- [x] Replace `DropdownItem` empty-Href discrimination with typed `DropdownItemKind` enum
- [x] Validate SwapOOB swapStyle parameter — panics on invalid styles (source: htmx/helpers.templ)

### Component Features

- [x] Make GlobalErrorHandling config values configurable — `ErrorHandlingConfig` struct
- [x] Extract error handling magic numbers — configurable via `ErrorHandlingConfig`
- [x] SimpleCard composes through Card internally
- [x] Toast duration configurable per-toast
- [x] Pagination ellipsis rendering for large page ranges
- [x] Table caption support — Caption field + render `<caption>` element
- [x] Avatar status dot scaling — fixed per size via `avatarDotSizeLookup` map
- [x] Breadcrumb separator customization — `Separator` field on `BreadcrumbsProps`
- [x] Use `net/url` for pagination URL construction
- [x] Make `PageProps` zero-value safe — `DefaultPageProps()` provides all defaults
- [x] Magic theme colors — extracted to `DefaultThemeColor` and `DefaultDarkThemeColor` constants
- [x] Eliminate 4-way icon list split brain — `allIconNames()` auto-generated from `iconPathData` map
- [x] Move avatar fallback SVG to icons package — `svg.PathAvatarFill` constant
- [x] `ComponentProps` interface with `GetBaseProps()` / `SetBaseProps()` for all props structs
- [x] Add stroke-width option to `icons.Icon` — `IconWithStrokeWidth(name, class, strokeWidth)`
- [x] Badge click/href support — renders as `<a>` when `Href` is set (source: display/badge.templ)
- [x] ProgressBar indeterminate state — `Indeterminate bool` with `aria-busy` and animated bar (source: feedback/progress.templ)
- [x] Step indicator vertical variant — `StepVertical` orientation with vertical connector lines (source: feedback/progress.templ)
- [x] Client-side JS tab switching — `ClientSide bool` with keyboard nav (ArrowLeft/Right, Home, End) (source: display/tabs.templ)
- [x] Tabs keyboard navigation — ArrowLeft/Right, ArrowUp/Down, Home, End (source: display/tabs.templ)

### New Components

- [x] Add Radio button component (source: forms/radio.templ)
- [x] Add Toggle/Switch component (source: forms/toggle.templ)
- [x] Add File input component (source: forms/file_input.templ)
- [x] Add Form component wrapping inputs + CSRF token (source: forms/form.templ)
- [x] Add more Heroicons — 75 total (was 45), +30 navigation/action/utility icons (source: icons/)
- [ ] Add Date Picker component (source: docs/status/)
- [ ] Add Combobox/Autocomplete component (source: docs/status/)
- [ ] Add Dialog/Drawer component variants (source: docs/status/)

### Accessibility

- [x] Add `aria-live="polite"` directly in HTMX error handling — `tc-error-announcer` div
- [x] Add Table header `scope` attributes — `<th scope="col">` on all header cells
- [x] Add EmptyState landmark role — `role="status"` on outer div
- [x] Add Breadcrumb structured data (JSON-LD) — `JSONLD` field on `BreadcrumbsProps`
- [x] Add Pagination SEO `rel=prev/next`
- [x] Investigate tooltip JS-based `aria-describedby` injection — already done with `props.ID` pattern
- [ ] Consolidate inline JS into shared init strategy — 10 script blocks across 7 files
- [ ] Add `uint` type for Pagination fields (source: navigation/pagination.templ)

### Testing

- [x] BDD tests for navigation, htmx, layout, icons packages
- [x] Dark mode `dark:` class output verification tests (source: display/a11y_test.go)
- [x] Benchmark tests for display, feedback, navigation packages (source: \*/benchmark_test.go)
- [x] Component composition integration tests — Card+Badge, Table+Content, StatCard (source: display/composition_test.go)
- [x] Godoc ExampleXxx() functions for forms package — Form, Input, Select, Textarea (source: forms/example_test.go)
- [ ] Improve coverage for functions below 70%: fillIcon, Select, Textarea
- [ ] Convert snapshot tests to golden file comparison
- [ ] Consistent nonce propagation audit across all components
- [ ] Add accessibility audit automation — axe-core/pa11y

### Infrastructure

- [x] Set up GitHub Actions CI — build + test + lint on push/PR
- [x] Pre-commit hook with `chmod +x`
- [x] Set coverage threshold in CI (60%)
- [x] Add build test for `examples/` in CI
- [x] Audit `tailwind-merge-go` thread safety — `sync.Mutex` IS required (source: utils/utils.go)
- [ ] Verify `go get` works from clean project
- [ ] Set up goreleaser for tag-based releases
- [ ] Modularize into Go workspace (10-module `go.work`)
- [ ] Consider `go:generate stringer` for enums
- [ ] Consider `Validate() error` method on props structs

### Documentation

- [x] Update README.md for new API (source: README.md)
- [x] Write ADR for filled vs stroke icon convention (ADR 0004)
- [x] Write ADR for JS attachment patterns (ADR 0005)
- [x] Add ADR for FeedbackType unification decision (ADR 0006)
- [x] Fill in placeholder terms in DOMAIN_LANGUAGE.md
- [x] Document thread-safety requirement on `utils.Class()` in CONTRIBUTING.md
- [x] Document PageProps convention — why it doesn't embed BaseProps (CONTEXT.md)
- [ ] Documentation site generation — pkgsite, doc2go, or custom

### Release & Discovery

- [ ] Tag v0.2.0 release and update CHANGELOG.md
- [ ] Submit to awesome-templ for discoverability
- [ ] Open PR on templ.guide to get listed
- [ ] Cross-link ecosystem in README — cqrs-htmx, go-cqrs-lite (GOTH stack story)

### Housekeeping

- [x] Prune old status reports — keep last 2, archive rest
- [ ] Investigate gopls QF1003 suppression for generated `*_templ.go` files
- [ ] Extract shared Tailwind preset/theme configuration file
- [ ] Plan v1.0 API freeze scope and timeline
