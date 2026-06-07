# TODO List

**Updated:** 2026-06-08
**Coverage:** 66%+ | **Tests:** 1049+ | **Packages:** 10+demo | **Lint:** 0 issues

> Reviewed 2026-06-08: ~30 items completed. See git history for previous state.

---

## 🔴 Bugs & Security

- [x] Fix JS re-attachment after HTMX DOM swaps — **by design**: all components use document-level event delegation; dynamically added elements are handled automatically. Global singleton guards prevent duplicate listeners. (source: display/, feedback/, layout/)
- [x] Fix ThemeToggle multi-instance bug — **not a bug**: event delegation on `document` catches clicks on ALL `[data-theme-toggle]` elements regardless of when they're added. Removed unnecessary IIFE for consistency. (source: layout/theme.templ)
- [x] Fix Accordion JS to IIFE-per-instance pattern — **not needed**: event delegation handles dynamically added accordions; IIFE-per-instance would create redundant listeners. (source: display/accordion.templ)
- [x] Fix Checkbox unconditional `id=""` → conditional `if props.ID != ""` rendering (source: forms/input.templ:160)
- [x] Fix pre-commit hook — uses scripts/pre-commit.sh (source: .git/hooks/pre-commit)

## 🟡 Breaking Changes (defer to v1.0)

- [ ] Move test helpers to `internal/testutil/` — breaking change for external consumers (source: utils/test_helpers.go)
- [ ] Spinner BaseProps conversion — `Spinner(size, colorClass)` positional args → `SpinnerProps` struct (source: feedback/loading.templ)
- [ ] SimpleNav BaseProps conversion — `(brandText, brandHref, links, currentPath)` → `SimpleNavProps` struct (source: navigation/nav.templ)
- [ ] Add BaseProps to StepIndicatorProps (source: feedback/step_indicator.templ)
- [ ] Pagination uint fields — `CurrentPage` and `TotalPages` should be `uint` (source: navigation/pagination.templ)

## 🟢 Enhancements

### Validation & Robustness

- [x] Add TrendDirection consistent validation — `normalizeTrend()` in display/card.templ normalizes invalid values to TrendNone (source: display/card.templ)
- [x] Extract hardcoded SVG path strings to named constants — `internal/svg` now has PathChevronDown, PathChevronSmall, PathArrowUp/Down/Left/Right, PathAvatarFill (source: internal/svg/svg.templ)
- [x] Propagate AriaLabel from BaseProps in all 25 components — conditional `aria-label` on root elements; hardcoded labels use AriaLabel override via `utils.Ternary` (source: all packages)
- [x] Extract shared dismiss JS — `dismissScript()` in feedback/styles.go eliminates duplicate between Alert and Toast (source: feedback/alert.templ, feedback/toast.templ)
- [x] Unify alertIconName/toastIconName — shared `feedbackIconName()` in feedback/styles.go (source: feedback/styles.go)
- [x] Fix toast JS string builder — `fmt.Fprintf` replaces string concatenation in `toastJSStyles()` (source: feedback/toast.templ)
- [x] Validate SelectOption contradiction (Disabled + Selected simultaneously) — `normalizeSelectOptions()` clears Selected when Disabled (source: forms/select.templ)
- [x] Validate `|` separator doesn't appear in SVG paths — `iconPaths()` panics on empty segments (source: icons/icon_paths.go)
- [x] Add `DropdownItem.Disabled` field (source: display/dropdown.templ)
- [x] Add `InputProps.MaxLength` field (source: forms/input.templ)
- [x] Add `TextareaProps.MaxLength` field (source: forms/textarea.templ)
- [x] Add `CheckboxProps.Value` field (source: forms/input.templ)
- [x] Replace `DropdownItem` empty-Href discrimination with typed `DropdownItemKind` enum (LinkItem, ButtonItem) — backward compat via `IsLink()` fallback (source: display/dropdown.templ)
- [ ] Validate SwapOOB swapStyle parameter (source: htmx/helpers.templ)

### Component Features

- [x] Make GlobalErrorHandling config values configurable — `ErrorHandlingConfig` struct with `MaxErrorHistory`, `MaxRetries`, `RetryDelayMS` (source: htmx/error_handling.templ)
- [x] Extract error handling magic numbers — configurable via `ErrorHandlingConfig` (source: htmx/error_handling.templ)
- [x] SimpleCard composes through Card internally (source: display/card.templ)
- [x] Toast duration configurable per-toast (source: feedback/toast.templ)
- [x] Pagination ellipsis rendering for large page ranges — shows first/last page + ellipsis when range is truncated (source: navigation/pagination.templ)
- [x] Table caption support — Caption field + render `<caption>` element (source: display/table.templ)
- [x] Avatar status dot scaling — fixed per size via `avatarDotSizeLookup` map (source: display/avatar.templ)
- [x] Breadcrumb separator customization — `Separator` field on `BreadcrumbsProps` (source: navigation/breadcrumbs.templ)
- [x] Use `net/url` for pagination URL construction (source: navigation/pagination.templ)
- [x] Make `PageProps` zero-value safe — `DefaultPageProps()` provides all defaults (source: layout/base.templ)
- [x] Magic theme colors — extracted to `DefaultThemeColor` and `DefaultDarkThemeColor` constants (source: layout/base.templ)
- [x] Eliminate 4-way icon list split brain — `allIconNames()` auto-generated from `iconPathData` map (source: icons/icon_names.go)
- [x] Move avatar fallback SVG to icons package — `svg.PathAvatarFill` constant (source: display/avatar.templ)
- [x] `ComponentProps` interface with `GetBaseProps()` / `SetBaseProps()` for all props structs (source: utils/base_props.go)
- [x] Add stroke-width option to `icons.Icon` — `IconWithStrokeWidth(name, class, strokeWidth)` (source: icons/icon.templ)
- [ ] Replace hardcoded SVGs with icon system — remaining: StepIndicator (file does not exist yet) (source: feedback/)

### New Components

- [x] Add Radio button component (source: forms/radio.templ)
- [x] Add Toggle/Switch component (source: forms/toggle.templ)
- [x] Add File input component (source: forms/file_input.templ)
- [ ] Add Date Picker component (source: docs/status/)
- [ ] Add Combobox/Autocomplete component (source: docs/status/)
- [ ] Add Dialog/Drawer component variants (source: docs/status/)
- [ ] Add Form component wrapping inputs + validation (source: forms/)
- [ ] Add skeleton component variants (card, table, list) (source: docs/status/)
- [ ] Step indicator vertical variant (source: feedback/step_indicator.templ)
- [ ] Badge click/href support (source: display/badge.templ)
- [ ] Add more Heroicons (currently 45 of ~300) (source: icons/)
- [ ] ProgressBar indeterminate state (source: feedback/progressbar.templ)
- [ ] Add client-side JS tab switching (source: display/tabs.templ)
- [ ] Tabs active tab keyboard navigation (arrow keys) (source: display/tabs.templ)

### Accessibility

- [x] Add `aria-live="polite"` directly in HTMX error handling — `tc-error-announcer` div with sr-only (source: htmx/error_handling.templ)
- [ ] Consolidate inline JS into shared init strategy — 10 script blocks across 7 files (source: layout/, display/, feedback/)
- [x] Add Table header `scope` attributes — `<th scope="col">` on all header cells (source: display/table.templ)
- [x] Add EmptyState landmark role — `role="status"` on outer div (source: display/empty_state.templ)
- [x] Add Breadcrumb structured data (JSON-LD) — `JSONLD` field on `BreadcrumbsProps` (source: navigation/breadcrumbs.templ)
- [x] Add Pagination SEO `rel=prev/next` (source: navigation/pagination.templ)
- [x] Investigate tooltip JS-based `aria-describedby` injection — already done with `props.ID` pattern (source: display/tooltip.templ)
- [ ] Add `uint` type for Pagination fields (source: navigation/pagination.templ)

### Testing

- [ ] Improve coverage for functions below 70%: fillIcon, Select, Textarea (source: multiple files)
- [ ] Add `DefaultLoadingOverlayProps` test (source: feedback/loading.templ)
- [ ] Add `DefaultBreadcrumbsProps` test (source: navigation/breadcrumbs.templ)
- [ ] Add Nav empty `Links` test (source: navigation/nav.templ)
- [ ] Add CSRFToken empty string test (source: htmx/csrf.templ)
- [ ] Add tooltip position edge case test (source: display/tooltip.templ)
- [x] BDD tests for navigation package — Nav, Pagination, Breadcrumbs (source: navigation/)
- [x] BDD tests for htmx package — Loading, ErrorHandling (source: htmx/)
- [x] BDD tests for layout package — Base, Minimal, Theme (source: layout/)
- [x] BDD tests for icons package — all icons render, unknown fallback (source: icons/)
- [ ] Add `utils.AssertContainsClass` test helper — replace fragile exact-string class tests (source: display/)
- [ ] Remove duplicate test data — extract shared testNavLinks helper (source: navigation/)
- [ ] Consolidate test files — eliminate duplication (source: multiple test files)
- [ ] Convert snapshot tests to golden file comparison (source: all packages)
- [ ] Add benchmark tests for Icon, Card, Table, Nav renders (source: multiple test files)
- [ ] Add component composition integration tests — Card+Badge, Nav+Avatar, Table+Dropdown, Modal+Form (source: display/)
- [ ] Add integration test: full page render with Base + Nav + Content + Footer (source: layout/)
- [ ] Dark mode `dark:` class output verification tests (source: all packages)
- [ ] Consistent nonce propagation audit across all components (source: all packages)
- [ ] Add circular import guard test (source: docs/status/)
- [ ] Add accessibility audit automation — axe-core/pa11y (source: docs/status/)
- [ ] Move ProgressBar a11y test from `display/` to `feedback/` package (source: display/a11y_test.go)
- [ ] Add `TableCell` documentation — `Content` takes priority over `Text` (source: display/table.templ)

### Infrastructure

- [x] Fix demo app to enable HTMX — DefaultPageProps already sets HTMXVersion:"2.0.6" (source: examples/demo/main.go)
- [x] Set up GitHub Actions CI — build + test + lint on push/PR (source: .github/workflows/ci.yaml)
- [ ] Verify `go get` works from clean project (source: docs/status/)
- [x] Pre-commit hook needs `chmod +x` (source: scripts/pre-commit.sh)
- [ ] Set up goreleaser for tag-based releases with cross-compilation and checksums (source: .goreleaser.yml)
- [ ] Add `go vet` / staticcheck to CI pipeline (source: docs/status/)
- [ ] Set coverage threshold in CI (e.g., 60%) (source: docs/status/)
- [ ] Add build test for `examples/` in CI (source: examples/demo/)
- [ ] Modularize into Go workspace (10-module `go.work`) (source: docs/modularization/)
- [ ] Deploy demo site — add cmd/demo/main.go with HTTP server + Fly.io/Railway config (source: cmd/demo/)
- [ ] Investigate nix flake for reproducible builds (source: flake.nix)
- [ ] Consider `go:generate stringer` for enums (source: docs/status/)
- [ ] Consider `Validate() error` method on props structs (source: docs/status/)
- [ ] Investigate visual regression testing (source: docs/status/)
- [ ] Audit `tailwind-merge-go` thread safety — remove `sync.Mutex` in `utils.Class()` if stateless (source: utils/utils.go)

### Documentation

- [ ] Update README.md for new API (AvatarStatus, StatCardProps, BreadcrumbsProps) (source: README.md)
- [ ] Update CONTEXT.md with JS pattern decision documentation (source: CONTEXT.md)
- [ ] Write ADR for filled vs stroke icon convention (source: docs/adr/)
- [ ] Write ADR for JS attachment patterns — singleton vs IIFE (source: docs/adr/)
- [ ] Add ADR for FeedbackType unification decision (source: docs/adr/)
- [ ] Add `go doc` ExampleXxx() functions — ExampleAlert, ExampleBadge, ExampleCard, ExamplePagination, ExampleIcon (source: feedback/, display/, navigation/, icons/)
- [ ] Fill in placeholder terms in DOMAIN_LANGUAGE.md (source: docs/DOMAIN_LANGUAGE.md)
- [ ] Document thread-safety requirement on `utils.Class()` in CONTRIBUTING.md (source: utils/utils.go)
- [ ] Document 20×20 fill vs 24×24 stroke icon convention (source: internal/svg/svg.templ)
- [ ] Document PageProps convention — why it doesn't embed BaseProps (source: CONTEXT.md)
- [ ] Documentation site generation — pkgsite, doc2go, or custom (source: project root)

### Release & Discovery

- [ ] Tag v0.2.0 release and update CHANGELOG.md (source: project root)
- [ ] Submit to awesome-templ for discoverability (source: GitHub PR)
- [ ] Open PR on templ.guide to get listed (source: GitHub PR)
- [ ] Cross-link ecosystem in README — cqrs-htmx, go-cqrs-lite (GOTH stack story) (source: README.md)
- [ ] Build real-world example app using all three libs (clone-and-run) (source: docs/STANDOUT-IDEAS.md)
- [ ] Build and deploy live component showcase site (source: docs/STANDOUT-IDEAS.md)

### Housekeeping

- [ ] Prune old status reports — keep last 2, archive rest (source: docs/status/)
- [ ] Investigate gopls QF1003 suppression for generated `*_templ.go` files (source: display/card_templ.go)
- [ ] Extract shared Tailwind preset/theme configuration file (source: project root)
- [ ] Plan v1.0 API freeze scope and timeline (source: docs/status/)
- [ ] Cross-package circular import audit — icons ↔ feedback full analysis (source: icons/, feedback/)
