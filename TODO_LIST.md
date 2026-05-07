# TODO List — templ-components

**Updated:** 2026-05-07

Legend: ✅ Done | 🔨 In Progress | ⬜ Not Started | ❌ Blocked

---

## Critical Bugs (Found 2026-05-04)

| #   | Status | Task                                                         | Priority | Notes                                                                             |
| --- | ------ | ------------------------------------------------------------ | -------- | --------------------------------------------------------------------------------- |
| 1   | ✅     | Fix `class="dropdownItemClass"` rendering literal string     | P0       | Was rendering literal text, not Go variable. Fixed: `class={ dropdownItemClass }` |
| 2   | ✅     | Fix `class="emptyStateActionClass"` rendering literal string | P0       | Same bug in empty_state.templ L27,L35. Fixed.                                     |
| 3   | ✅     | Fix `.golangci.yml` Go version mismatch                      | P0       | Was `go: "1.23"`, CI uses `"1.26"`. Fixed.                                        |
| 4   | ✅     | Fix README.md stale `layout.BaseProps` → `PageProps`         | P0       | Two occurrences fixed.                                                            |

## Architecture

| #   | Status | Task                                                                 | Priority | Notes                                                                                              |
| --- | ------ | -------------------------------------------------------------------- | -------- | -------------------------------------------------------------------------------------------------- |
| 5   | ✅     | Extract shared SVG helpers to `internal/svg/`                        | —        | Done. FillIcon, SpinnerSVG.                                                                        |
| 6   | ✅     | Unify `alertStyleSet`/`toastStyleSet` into shared `feedbackStyleSet` | P1       | Done. Shared `feedbackStyleSet` struct + `lookupFeedbackStyle[T]()` generic in `feedback/styles.go`. |
| 7   | ✅     | Generic `utils.MapEnum[T ~string]`                                   | —        | Done. Used in display/empty_state.go                                                               |
| 8   | ✅     | All Props structs embed `utils.BaseProps`                            | —        | Done.                                                                                              |
| 9   | ✅     | Map-based style lookups (not switches)                               | —        | Done.                                                                                              |
| 10  | ✅     | Rename `layout.BaseProps` → `PageProps`                              | —        | Done.                                                                                              |
| 11  | ✅     | Deepen icon rendering: path-data map + shared SVG helper             | P1       | Done. 187-line switch → `iconPathData` map + `strokeIcon()` template. ~100 lines removed.           |
| 12  | ⬜     | Extract shared form error/aria helper                                | P2       | Low impact — FieldError already shared. ARIA attrs can't be extracted from templ attributes.       |
| 13  | ✅     | Replace `AvatarProps.Online/Offline bool` with `AvatarStatus` enum   | P2       | Done. `AvatarStatusOnline`/`AvatarStatusOffline`/`AvatarStatusNone`. Impossible state eliminated.  |
| 14  | ✅     | Replace `StatCard.positive bool` with `TrendDirection` enum          | P2       | Done. `TrendUp`/`TrendDown`/`TrendNone`. Now uses `StatCardProps` struct.                          |
| 15  | ✅     | Fix `HTMXSRI string` → `HTMXUseSRI bool` in PageProps                | P2       | Done. Boolean field with `true` default.                                                           |
| 16  | ✅     | Fix integer division in ProgressBar percent                          | P2       | Done. `float64` division with `%.0f` formatting.                                                   |
| 17  | ✅     | Add `Content templ.Component` to `TableCell`                         | P2       | Done. Falls back to `Text` when `Content` is nil.                                                  |
| 18  | ✅     | Implement `TableProps.Bordered` styling                              | P3       | Done. Adds `border border-gray-200 dark:border-slate-700` to table cells.                           |
| 19  | ✅     | Add tests for `icons.IconAttrs`                                       | P3       | Done. Tests for aria-label and aria-hidden behavior. Plus `TestAllIconsRender` for all 42 icons.    |
| 20  | ✅     | Add `CONTEXT.md` with architecture decisions                         | —        | Done.                                                                                              |
| 21  | ✅     | Add `docs/adr/` for architecture decision records                    | —        | Done. ADR-0001.                                                                                    |

## Testing

| #   | Status | Task                                              | Priority | Notes                                                       |
| --- | ------ | ------------------------------------------------- | -------- | ----------------------------------------------------------- |
| 22  | ✅     | Add render tests for breadcrumbs                  | P1       | Done. a11y, dark mode, edge cases in `navigation/a11y_test.go`. |
| 23  | ✅     | Add render tests for nav (Nav, SimpleNav, Footer) | P1       | Done. Brand, sticky, dark mode, right items in `navigation/a11y_test.go`. |
| 24  | ✅     | Add render tests for mobile_menu                  | P2       | Already covered in `navigation/snapshot_test.go`. Enhanced. |
| 25  | ✅     | Add render tests for htmx error_handling          | P2       | Done. A11y, nonce, event listeners in `htmx/a11y_test.go`.  |
| 26  | ⬜     | Convert snapshot tests to golden file comparison  | P1       | Deferred. Current substring assertions work well for library. |
| 27  | ✅     | Add a11y attribute validation tests               | P1       | Done. Modal, Dropdown, Tabs, Tooltip, Accordion, Avatar, Table in `display/a11y_test.go`. |
| 28  | ✅     | Add dark mode output verification tests           | P2       | Done. Card, Badge, Table, Dropdown, Avatar, Nav, Footer, Base layout. |
| 29  | ⬜     | Add component composition tests                   | P2       | Deferred. Complex to test outside of templ files.             |
| 30  | ✅     | Add benchmark tests for hot paths                 | P2       | Done. `utils.Class()` and Badge render benchmarks in `display/a11y_test.go`. |
| 31  | ✅     | Add tests for `internal/svg` package              | —        | Done. FillIcon and SpinnerSVG render tests added.            |
| 32  | ✅     | Add direct test for `utils.MapEnum[T]`            | —        | Done. 3 subtests: found, missing, empty key.                 |
| 33  | ✅     | Test all `Default*Props()` constructors           | P2       | Done. Card, Badge, Modal, ProgressBar constructors tested.   |
| 34  | ✅     | Fix TestPtr bug                                   | —        | `new(v)` → `Ptr(v)`.                                         |
| 35  | ✅     | Nonce propagation audit                           | —        | All inline scripts verified.                                 |
| 36  | ✅     | Add `SecurityHeaders` test to layout              | P2       | Done. X-Content-Type-Options, Referrer-Policy, skip link in `layout/a11y_test.go`. |
| 37  | ✅     | CSP compliance for all inline scripts             | —        | All scripts use nonce.                                        |
| 38  | ✅     | Verify dropdown JS template interpolation is XSS-safe | P2   | Done. Templ auto-escapes special chars in `display/a11y_test.go`. |

## DevOps & Tooling

| #   | Status | Task                                 | Priority | Notes                                             |
| --- | ------ | ------------------------------------ | -------- | ------------------------------------------------- |
| 39  | ✅     | Set up GitHub Actions CI             | —        | Done. Go 1.26, lint+build+test.                   |
| 40  | ⬜     | Release automation (goreleaser)      | P3       | Tag-based releases.                               |
| 41  | ⬜     | Investigate nix flake migration      | P3       | No build system exists.                           |
| 42  | ✅     | Pre-commit hook for `templ generate` | P2       | Done. `scripts/pre-commit.sh` — auto-stages generated files. |
| 43  | ✅     | Add `layout/sri.go` package comment  | P3       | Done. `// Package layout provides Sub-Resource Integrity hashes...` |

## Documentation

| #   | Status | Task                                 | Priority | Notes                            |
| --- | ------ | ------------------------------------ | -------- | -------------------------------- |
| 44  | ✅     | Create FEATURES.md                   | —        | Comprehensive feature inventory. |
| 45  | ✅     | Create TODO_LIST.md                  | —        | This file.                       |
| 46  | ✅     | Update CHANGELOG.md with recent work | P1       | Done. Full changelog with breaking changes. |
| 47  | ⬜     | Create example/demo app              | P2       | Showcase all components.         |
| 48  | ⬜     | Documentation site generation        | P3       | Auto-generated from source.      |
| 49  | ⬜     | Version migration guides             | P3       | Breaking changes documentation.  |

## Deduplication

| #   | Status | Task                                            | Priority | Notes                                       |
| --- | ------ | ----------------------------------------------- | -------- | ------------------------------------------- |
| 50  | ✅     | Semantic deduplication (13→7 clone groups)      | —        | Extracted 10+ sub-templates.                |
| 51  | ⬜     | Remaining 9 clone groups (dupl)                 | P3       | All in test files — structural only.        |
| 52  | ✅     | Move `boolString()` to `utils/` and standardize | P3       | Done. `utils.BoolString()` replaces local `boolString()`. |

---

## New Items (Discovered 2026-05-07 Audit)

| #   | Status | Task                                                         | Priority | Notes                                                       |
| --- | ------ | ------------------------------------------------------------ | -------- | ----------------------------------------------------------- |
| 53  | ⬜     | Add 11 missing `DefaultXxxProps()` constructors              | P1       | EmptyState, Modal, Loading, Progress, Label, ErrorHandling, HTMXLoading, Breadcrumbs, MobileMenu, NavLink, Base (verify) |
| 54  | ⬜     | Add `aria-required` to form inputs when `Required: true`     | P1       | WCAG requirement. Input, Select, Textarea.                  |
| 55  | ⬜     | Add `<html lang>` to Base layout (new `Lang` field)          | P1       | WCAG 3.1.1. Add `Lang string` to PageProps, default "en".   |
| 56  | ⬜     | Add skip-to-content link in Base layout                      | P1       | WCAG 2.4.1. Hidden link that jumps to `#main-content`.      |
| 57  | ⬜     | Update migration guide with `icons.Name` breaking change     | P1       | `Icon: "folder"` → `Icon: icons.Folder` not in docs.        |
| 58  | ⬜     | Add `alt` text to Avatar `<img>`                             | P1       | WCAG 1.1.1. Add `Alt` field to AvatarProps.                 |
| 59  | ⬜     | Add Table `Caption` field + render `<caption>`               | P2       | WCAG 1.3.1. Tables should have captions for screen readers. |
| 60  | ⬜     | Fix pre-commit hook to be executable                         | P3       | `chmod +x scripts/pre-commit.sh` — every commit shows warning. |
| 61  | ⬜     | Convert 6 remaining components to `utils.Class()`            | P2       | Alert, Loading, Progress, Toast, NavLink (Modal: KV limitation). |
| 62  | ⬜     | Add `aria-live="polite"` to HTMX loading indicators          | P2       | Screen reader announcements for dynamic loading.            |
| 63  | ⬜     | Add `aria-live="polite"` to HTMX error handling              | P2       | Dynamic error announcements for screen readers.             |
| 64  | ⬜     | Improve forms test coverage (58% → 75%+)                     | P2       | Lowest package coverage.                                    |
| 65  | ⬜     | Improve utils test coverage (56% → 75%+)                     | P2       | MergeAttrs, CurrentYear, Deref undertested.                 |
| 66  | ⬜     | Add EmptyState landmark role (`role="region"`)               | P3       | Screen reader navigation aid.                               |
| 67  | ⬜     | Add Table header `scope` attributes                          | P2       | `<th scope="col">` for screen reader column association.    |
| 68  | ⬜     | Document `PageProps` not embedding `BaseProps`               | P3       | Only Props struct that breaks convention. Document why.      |
| 69  | ⬜     | Exclude `examples/` from golangci-lint properly              | P3       | 23 issues in demo/main.go. Can't exclude via config.         |

## Session 3 (2026-05-07) — Class Merge, Type Safety, A11y, Docs
| 1 | `39e237e` | fix: Tailwind class merge, a11y gaps, and form help text IDs |
| 2 | `1c014ef` | feat: add 17 missing DefaultXxxProps() constructors across all packages |
| 3 | `3d5aded` | feat!: type-safe icons — DropdownItem.Icon and EmptyStateProps.Icon now use icons.Name |
| 4 | `c9981d7` | feat(a11y): add modal focus trap and Escape key handler |
| 5 | `a1a81b5` | feat(a11y): add dropdown keyboard navigation |
| 6 | `e8acc3f` | fix(a11y): add proper ARIA linkage to tabs |
| 7 | `e6b403c` | fix(a11y): add id to tooltip div for aria-describedby linkage |
| 8 | `9884910` | docs: add package doc comments for all packages |

Key changes:
- **Silent bug fix**: `utils.Class()` replaces comma-join in forms, Badge, StatCard, Nav — consumer class overrides were silently ignored
- **17 DefaultXxxProps() constructors** added across all packages
- **Type safety**: `Icon` fields changed from `string` to `icons.Name` (breaking change)
- **Modal a11y**: Focus trap, Escape key, focus management on open/close
- **Dropdown a11y**: Arrow key navigation, Escape to close, auto-focus first item
- **Tabs a11y**: Proper id/aria-controls/aria-labelledby/tabindex linkage
- **Tooltip a11y**: Deterministic id on tooltip div for aria-describedby
- **Package docs**: doc.go for display, feedback, htmx, icons, layout, navigation

## Session 2 (2026-05-04) — Architecture, Tests, Demo

- Unified `alertStyleSet`/`toastStyleSet` into shared `feedbackStyleSet` + generic `lookupFeedbackStyle[T]()`
- Deepened icon rendering: 187-line switch → `iconPathData` map + `strokeIcon()` sub-template
- Added `AvatarStatus` enum replacing `Online`/`Offline` boolean fields
- Added `TrendDirection` enum for `StatCardProps` (was `StatCard(value, label, change, positive)`)
- Fixed `HTMXSRI string` → `HTMXUseSRI bool` in PageProps
- Fixed ProgressBar integer division truncation (now uses float64)
- Added `Content templ.Component` to `TableCell`
- Implemented `TableProps.Bordered` styling (was dead code)
- Added `utils.BoolString()` replacing local `boolString()` in accordion
- Added `icons.IconAttrs()` tests + `TestAllIconsRender` for all 42 icons
- Added `layout/sri.go` package comment
- Created pre-commit hook script (`scripts/pre-commit.sh`)
- Updated CHANGELOG.md with comprehensive entries including breaking changes
- Added a11y tests: navigation breadcrumbs/nav/footer, display (modal, dropdown, tabs, tooltip, accordion, avatar, table), htmx error handling, layout security headers
- Added dark mode verification tests across display, navigation, and layout packages
- Added `Default*Props()` constructor tests for Card, Badge, Modal, ProgressBar
- Added dropdown XSS safety test verifying templ auto-escaping
- Added benchmark tests for `utils.Class()` and Badge rendering

## Session 1 (2026-05-03) — Initial Audit

- Fixed critical bug: `class="dropdownItemClass"` rendering literal string (dropdown.templ L41)
- Fixed critical bug: `class="emptyStateActionClass"` rendering literal string (empty_state.templ L27, L35)
- Fixed `.golangci.yml` Go version: 1.23 → 1.26
- Fixed README.md stale `layout.BaseProps` → `PageProps` references
- Full code review: 44 issues found (2 critical, 4 high, 10 medium, 20+ low)
- Features audit: verified 56 components, 44 types, 42 icons
- Code quality scan: build ✓, tests ✓, 1 lint issue, 9 clone groups (all test files)
- Architecture review: 5 deepening opportunities identified
- Architecture visualization: current + target state D2 diagrams rendered
- Added BDD tests: display/bdd_test.go (20 scenarios), feedback/bdd_test.go (16 scenarios), forms/bdd_test.go (16 scenarios)
- Added internal/svg/svg_test.go (5 subtests for FillIcon and SpinnerSVG)
- Added direct utils.MapEnum test (3 subtests)
