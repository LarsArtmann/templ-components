# TODO List

**Updated:** 2026-06-27
**Coverage:** ~75% | **Tests:** 456 functions + ~650 subtests | **Packages:** 10+demo | **Lint:** 0 issues | **Icons:** 101 (100 path + 1 Spinner)

> Session 5: Modal/Drawer overlay JS extraction + CSP fix, SimpleNav props conversion, Alert/Toast style unification, Form children pattern, Pagination uint, FormFieldWrapper adoption, button/select/toggle map lookups, stdlib htmlEscape.

---

## 🔴 Bugs & Security

- [x] Fix JS re-attachment after HTMX DOM swaps — **by design**: all components use document-level event delegation
- [x] Fix ThemeToggle multi-instance bug — **not a bug**: event delegation handles all instances
- [x] Fix Accordion JS to IIFE-per-instance pattern — **not needed**: event delegation handles dynamic accordions
- [x] Fix Checkbox unconditional `id=""` → conditional `if props.ID != ""` rendering
- [x] Fix pre-commit hook — uses scripts/pre-commit.sh
- [x] Fix tailwind-merge-go thread safety — `sync.Mutex` in `utils.Class()` is REQUIRED (source: utils/utils.go)

## 🟡 Breaking Changes (defer to v1.0)

- [ ] Move test helpers to `internal/testutil/` — breaking change for external consumers (source: utils/test_helpers.go) — **DEFERRED TO v1.0**: 70 test files + external consumers depend on `utils.Render`, `utils.AssertContains`, etc. Moving now would break every consumer with no functional benefit.
- [x] Spinner BaseProps conversion — `SpinnerProps` struct with BaseProps, Size, Color fields (source: feedback/loading.templ)
- [x] SimpleNav BaseProps conversion — `SimpleNav(SimpleNavProps)` with BaseProps embedding (source: navigation/nav.templ)
- [x] Add BaseProps to StepIndicatorProps (source: feedback/step_indicator.templ)
- [x] Pagination uint fields — `CurrentPage`, `TotalPages`, `MaxVisible` converted to `uint` (source: navigation/pagination.templ)

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
- [x] Add 25+ more Heroicons — 100 path icons + 1 Spinner = 101 total (was 75) (source: icons/)
- [x] Add Date Picker component (source: forms/date_picker.templ)
- [x] Add Combobox/Autocomplete component (source: forms/combobox.templ)
- [x] Add Dialog/Drawer component — accessible side panel with left/right slide, focus trap (source: display/drawer.templ)

### Accessibility

- [x] Add `aria-live="polite"` directly in HTMX error handling — `tc-error-announcer` div
- [x] Add Table header `scope` attributes — `<th scope="col">` on all header cells
- [x] Add EmptyState landmark role — `role="status"` on outer div
- [x] Add Breadcrumb structured data (JSON-LD) — `JSONLD` field on `BreadcrumbsProps`
- [x] Add Pagination SEO `rel=prev/next`
- [x] Investigate tooltip JS-based `aria-describedby` injection — already done with `props.ID` pattern
- [x] Consolidate inline JS into shared init strategy — Modal/Drawer overlay JS extracted to shared.go (overlayCloseJS/overlayOpenJS/overlayTrapJS). CSP-safe data-tc-close delegation replaces inline onclick.
- [x] Add `uint` type for Pagination fields (source: navigation/pagination.templ)

### Testing

- [x] BDD tests for navigation, htmx, layout, icons packages
- [x] Dark mode `dark:` class output verification tests (source: display/a11y_test.go)
- [x] Benchmark tests for display, feedback, navigation packages (source: \*/benchmark_test.go)
- [x] Component composition integration tests — Card+Badge, Table+Content, StatCard (source: display/composition_test.go)
- [x] Godoc ExampleXxx() functions for forms package — Form, Input, Select, Textarea (source: forms/example_test.go)
- [x] Improve coverage for functions below 70%: fillIcon, Select, Textarea — Verified: all handwritten Go logic is at 100% coverage. Remaining gaps (71-76%) are in generated `*_templ.go` boilerplate (templ runtime error branches), not business logic. These can't be meaningfully improved without testing templ internals.
- [~] Convert remaining snapshot tests to golden file comparison (feedback/ done, pattern in internal/golden) — Golden files exist for display (4), feedback (7), navigation (1). Converting the remaining 60+ assertion-based snapshot tests would create ~60 golden files for marginal value — the existing AssertContains tests already verify behavior. New components should use golden; old tests work fine as-is.
- [x] Consistent nonce propagation audit across all components — verified 2026-06-28: all 13 components with inline executable scripts use nonce. base.templ (external src) and breadcrumbs.templ (JSON-LD) correctly don't need nonce.
- [ ] Add accessibility audit automation — axe-core/pa11y

### Infrastructure

- [x] Set up GitHub Actions CI — build + test + lint on push/PR
- [x] Pre-commit hook with `chmod +x`
- [x] Set coverage threshold in CI (70%)
- [x] Add build test for `examples/` in CI
- [x] Audit `tailwind-merge-go` thread safety — `sync.Mutex` IS required (source: utils/utils.go)
- [x] Verify `go get` works from clean project — verified 2026-06-28: `go get github.com/larsartmann/templ-components@v0.4.0` builds and runs from empty project. Post-v0.4.0 types (ButtonHTMLType, utils.Version) need v0.5.0 tag.
- [x] Set up goreleaser for tag-based releases — **Not applicable**: This is a library, not a binary application. Go library versioning works via Git tags + the module proxy (`go get @v0.4.0`). No binary to release. CI pipeline already validates builds and runs tests on every push.
- [ ] Modularize into Go workspace (10-module `go.work`)
- [x] Consider `go:generate stringer` for enums — **Not feasible**: Go's `stringer` tool only supports integer-backed constants. All 26 enums in this library are `type X string` (e.g., `type BadgeType string`), which stringer explicitly rejects ("can't handle non-integer constant type"). The enum values ARE already strings, so a String() method would be redundant (they already stringify naturally via `string(myEnum)`). No action needed.
- [ ] Consider `Validate() error` method on props structs — **DEFERRED TO v1.0**: The library's current design philosophy is silent fallback (invalid enum → safe default, never crash). Adding `Validate() error` that returns errors for what currently falls back would change this philosophy. Needs a design decision: should Validate() replace the fallback pattern, or supplement it as an opt-in check? Implementation requires per-component methods (73 components).

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

- [x] Tag v0.2.0 release and update CHANGELOG.md
- [x] Tag v0.3.0 with Priority 2 features (Drawer, ValidationSummary, 25 icons, Spinner BaseProps) — also tagged v0.4.0
- [x] Submit to awesome-templ for discoverability — entry prepared: `- [templ-components](https://github.com/LarsArtmann/templ-components) — 73 accessible server-rendered components with 101 icons, Tailwind v4, HTMX helpers, go-error-family integration.` Ready for PR submission to https://github.com/nelsonlapreu/awesome-templ
- [ ] Open PR on templ.guide to get listed
- [x] Cross-link ecosystem in README — cqrs-htmx, go-cqrs-lite, go-error-family (GOTH stack story)

### Housekeeping

- [x] Prune old status reports — keep last 2, archive rest
- [x] Investigate gopls QF1003 suppression for generated `*_templ.go` files — Decision: leave as-is. The lint rule is correct for handwritten code; generated files are excluded via `.golangci.yml`. Adding per-file suppressions would mask real issues in handwritten code.
- [~] Extract shared Tailwind preset/theme configuration file — PARTIALLY DONE: tailwind.css in tailwind-v4-adoption-guide.md provides the pattern. A standalone preset file is deferred until multiple consumers exist.
- [x] Plan v1.0 API freeze scope and timeline — **v1.0 scope**: (1) Move test helpers to `internal/testutil/` (breaking), (2) Add `Validate() error` to all props structs, (3) Freeze all type names and prop field names, (4) Remove deprecated aliases (AlertType, ToastType). **Timeline**: After cqrs-htmx adminui fully adopts templ-components (in progress) and at least one external consumer. Target: after v0.6.0.
