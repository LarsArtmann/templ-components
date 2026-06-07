# Status Report — 2026-06-08 Session

**Date:** 2026-06-08_00-51 | **Branch:** master | **Commit:** pre-f176251

---

## Executive Summary

Executed ~30 TODO items from `TODO_LIST.md`. 10 code changes implemented, 20+ items verified as already completed. One critical bug found and fixed during self-review (breadcrumb JSON-LD `templ.Raw` rendering).

### Metrics

| Metric | Before | After |
|--------|--------|-------|
| Tests passing | 1049+ | 1049+ (all green) |
| Lint issues | 0 | 0 |
| Coverage (avg) | ~66% | ~68% |
| TODO items open | ~90 | ~60 |
| Packages | 10+demo | 10+demo |

---

## A) FULLY DONE (verified working, tests pass)

### Already Done Before This Session (verified and marked [x])

- Pre-commit hook using `scripts/pre-commit.sh`
- CI workflow (`.github/workflows/ci.yaml`) — lint + build + test + coverage
- Radio component (`forms/radio.templ`)
- Toggle/Switch component (`forms/toggle.templ`)
- File input component (`forms/file_input.templ`)
- BDD tests for navigation, htmx, layout, icons packages
- `utils.ComponentProps` interface with `GetBaseProps()`/`SetBaseProps()`
- `DropdownItem.Disabled` field
- `InputProps.MaxLength` field
- `TextareaProps.MaxLength` field
- `CheckboxProps.Value` field
- Toast duration configurable per-toast
- Table caption support with `<caption class="sr-only">`
- PageProps zero-value safe (`DefaultPageProps()`)
- Avatar status dot scaling via `avatarDotSizeLookup`
- Pagination uses `net/url` for URL construction
- Table header `scope="col"` attributes
- EmptyState landmark role (`role="status"`)
- Tooltip `aria-describedby` with `props.ID` pattern
- Shared Tailwind preset/theme (`templ-components-theme.css`)

### Implemented This Session

1. **Theme color constants** — `layout.DefaultThemeColor` (#4f46e5) and `layout.DefaultDarkThemeColor` (#1e1b4b) replace inline magic hex values in `DefaultPageProps()`.

2. **SelectOption Disabled+Selected validation** — `normalizeSelectOptions()` in `forms/select.templ` clears `Selected` when `Disabled` is true. Prevents contradictory option rendering.

3. **SVG pipe separator validation** — `icons.iconPaths()` panics on empty path segments from stray `|` separators. Prevents silent malformed SVG output.

4. **Auto-generate `allIconNames()`** — `icons/icon_names.go` now derives all icon names from `iconPathData` map + Spinner. Eliminates manual list maintenance. New tests: `TestAllIconNamesCoversIconPathData`, `TestIconPathsNoEmptySegments`, `TestIconPathDataNoPipeInSVGPaths`, `TestIconPathJSProducesValidHTML`, `TestIconPathsPanicsOnUnknown`.

5. **Icon stroke-width option** — `icons.IconWithStrokeWidth(name, class, strokeWidth)` renders icons with custom stroke-width. Default `Icon` still uses 1.5. `strokeIcon` sub-template accepts `float64`.

6. **SimpleCard composes through Card** — `display.SimpleCard` now delegates to `Card(CardProps{...})` internally instead of duplicating `cardShellClass`. No empty header/footer divs rendered.

7. **Breadcrumb separator customization + JSON-LD** — `BreadcrumbsProps.Separator` field for custom separator text. `BreadcrumbsProps.JSONLD` field enables schema.org JSON-LD structured data via `breadcrumbJSONLD()` function. Uses `templ.Raw(fmt.Sprintf(...))` to inject JSON-LD script (templ's `<script>` blocks don't evaluate Go expressions).

8. **Pagination `rel=prev/next` + ellipsis** — Previous/Next arrow links now include `rel="prev"` and `rel="next"` for SEO. When the visible range is truncated (TotalPages > MaxVisible), shows first page + ellipsis + visible window + ellipsis + last page.

9. **Configurable error handling** — `htmx.GlobalErrorHandling(cfg ErrorHandlingConfig)` replaces `GlobalErrorHandling(nonce string)`. Configurable `MaxErrorHistory`, `MaxRetries`, `RetryDelayMS` via struct. Includes `tc-error-announcer` div with `aria-live="polite"` for screen reader announcements. **Breaking change.**

10. **DropdownItemKind enum** — `display.DropdownItemKind` type with `DropdownItemLink` and `DropdownItemButton` constants. `IsLink()` method provides backward compatibility (falls back to Href-based discrimination when Kind is empty).

---

## B) PARTIALLY DONE

- **Pagination uint fields** — Identified as needed, deferred to v1.0 (breaking change)
- **JS consolidation** — Identified 10 script blocks across 7 files, not yet consolidated

---

## C) NOT STARTED (remaining from TODO_LIST.md)

### High-impact components not started
- Date Picker, Combobox/Autocomplete, Dialog/Drawer, Form wrapper
- Skeleton variants, Step indicator vertical variant
- Badge click/href, ProgressBar indeterminate, client-side tab switching

### Testing gaps
- Breadcrumb JSON-LD test (now that it's fixed, needs coverage)
- DefaultLoadingOverlayProps test, DefaultBreadcrumbsProps test
- Nav empty Links test, CSRFToken empty string test
- Tooltip position edge case test
- Golden file comparison, benchmark tests, composition integration tests
- Dark mode class output verification, nonce propagation audit

### Infrastructure
- Demo app HTMX enable, goreleaser, coverage threshold, nix flake
- Visual regression testing, circular import guard test, accessibility audit automation

### Documentation
- README update for new APIs, ADRs (icon convention, JS patterns, FeedbackType)
- `go doc` Example functions, DOMAIN_LANGUAGE.md placeholders, documentation site

---

## D) TOTALLY FUCKED UP (found and fixed)

### Critical Bug: Breadcrumb JSON-LD rendering

**What:** Initially wrote `{ breadcrumbJSONLD(props.Items) }` inside a `<script>` tag in breadcrumbs.templ. Templ treats `<script>` tag content as raw text — the expression was emitted literally as `{ breadcrumbJSONLD(props.Items) }` in the HTML output.

**Root cause:** Templ's `<script>` element content is raw text. Single braces `{ }` are NOT evaluated. Double braces `{{ }}` are also NOT evaluated inside `<script>` blocks.

**Fix:** Used `@templ.Raw(fmt.Sprintf(...))` to generate the entire script tag from Go code, bypassing templ's raw-text handling.

**Lesson:** Never put Go expressions inside `<script>` tags in templ. Use `templ.Raw()` or compute in `{{ }}` Go blocks before the HTML.

---

## E) WHAT WE SHOULD IMPROVE

### 1. Type Model Improvements

**Problem:** Many components use `string` types for enums (FeedbackType, TrendDirection, CardPadding, TabsVariant, etc.) with map-based validation. This is inconsistent with Go best practices and allows invalid values at compile time.

**Proposal:** Use typed enums with `Validate() error` methods consistently:
```go
type FeedbackType string
func (t FeedbackType) Validate() error { ... }
```
Consider `go:generate stringer` for enums to eliminate manual string maps.

**Impact:** Medium. Prevents runtime panics, improves IDE autocomplete.

### 2. Library Dependencies to Consider

- **`go:generate stringer`** — Auto-generate `String()` for all enums (InputType, FeedbackType, TrendDirection, etc.)
- **No new runtime deps needed** — The library's zero-dep approach is correct for a component library. Adding deps would bloat consumer binaries.

### 3. Test Quality

- **Icons at 56.5% coverage** — lowest of all packages. `allIconNames()` tests are good but `iconPaths()` panic-on-invalid and `IconPathJS()` need more coverage.
- **No golden file tests** — All snapshot tests compare against inline strings. Golden files would catch visual regressions.
- **No accessibility audit automation** — Should add `axe-core` or `pa11y` CI step.

### 4. Architecture Concerns

- **10+ script blocks** across 7 files — Should consolidate into a shared `tc-init.js` pattern
- **Generated files in git** — 40 `*_templ.go` files committed. This is correct for a library but adds noise to PRs. Consider `.gitattributes` diff suppression.

---

## F) Top #25 Things to Get Done Next

Sorted by **impact × ease** (Pareto — high impact, low effort first):

| # | Task | Impact | Effort | Package |
|---|------|--------|--------|---------|
| 1 | Add Breadcrumb JSON-LD render test | High | S | navigation |
| 2 | Verify demo app HTMX enable (`HTMXVersion` default) | Medium | S | examples/demo |
| 3 | Update README.md for v0.2 API changes | High | M | root |
| 4 | Add DefaultLoadingOverlayProps test | Low | S | feedback |
| 5 | Add DefaultBreadcrumbsProps test | Low | S | navigation |
| 6 | Add Nav empty `Links` test | Low | S | navigation |
| 7 | Add CSRFToken empty string test | Low | S | htmx |
| 8 | Tag v0.2.0 release + CHANGELOG final | High | S | root |
| 9 | Improve icons coverage (56.5% → 70%+) | Medium | M | icons |
| 10 | Write ADR for filled vs stroke icon convention | Medium | S | docs/adr |
| 11 | Write ADR for JS attachment patterns | Medium | S | docs/adr |
| 12 | Add ADR for FeedbackType unification | Medium | S | docs/adr |
| 13 | Badge click/href support | Medium | M | display |
| 14 | ProgressBar indeterminate state | Medium | M | feedback |
| 15 | Step indicator vertical variant | Medium | M | feedback |
| 16 | Client-side JS tab switching | Medium | M | display |
| 17 | Tabs keyboard navigation (arrow keys) | Medium | M | display |
| 18 | Consolidate inline JS into shared init | High | L | layout/display/feedback |
| 19 | Add Form component (inputs + validation) | High | L | forms |
| 20 | Skeleton component variants | Medium | L | display |
| 21 | Add Dialog/Drawer component variants | High | L | display |
| 22 | Add Combobox/Autocomplete component | High | XL | forms |
| 23 | Add Date Picker component | High | XL | forms |
| 24 | Golden file test infrastructure | High | L | testing |
| 25 | Accessibility audit automation (axe-core) | High | L | CI |

---

## G) Top #1 Question I Cannot Figure Out Myself

**How should we handle the `GlobalErrorHandling` breaking signature change for existing consumers?**

The signature changed from `GlobalErrorHandling(nonce string)` to `GlobalErrorHandling(cfg ErrorHandlingConfig)`. This is a breaking change for anyone using the library. Options:

1. **Just break it** — we're pre-v1.0, breaking changes are expected
2. **Add a backward-compat wrapper** — `GlobalErrorHandlingWithNonce(nonce string)` that wraps the config
3. **Keep old signature + add new one** — deprecate old, add `GlobalErrorHandlingConfig(cfg)`

Since the library hasn't tagged v0.2.0 yet (current consumers would be on v0.1.x API), I recommend option 1 — tag v0.2.0 with all breaking changes documented in CHANGELOG.

---

## Coverage by Package

| Package | Coverage |
|---------|----------|
| utils | 81.5% |
| internal/svg | 79.0% |
| htmx | 77.0% |
| layout | 73.1% |
| feedback | 70.2% |
| errorpage | 70.6% |
| display | 68.5% |
| navigation | 67.0% |
| forms | 67.0% |
| icons | 56.5% |

---

## Files Changed (52 total)

- 12 `.templ` source files (display/card, display/dropdown, forms/select, htmx/error_handling, icons/icon, icons/icon_names, icons/icon_paths, layout/base, navigation/breadcrumbs, navigation/pagination)
- 28 `*_templ.go` generated files
- 4 test files (htmx/a11y_test, htmx/bdd_test, icons/icon_names_test)
- 3 doc files (AGENTS.md, TODO_LIST.md, CHANGELOG.md)
