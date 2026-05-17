# TODO List — templ-components

**Updated:** 2026-05-07

Legend: ✅ Done | 🔨 In Progress | ⬜ Not Started | ❌ Blocked

---

## Session 4 (2026-05-07) — Comprehensive Audit

Comprehensive 8-skill audit: code quality scan, features audit, TODO list builder, architecture review, architecture improvement, architecture visualization, full code review, BDD testing.

### Build & Quality

| #   | Status | Task                              | Priority | Notes                                                |
| --- | ------ | --------------------------------- | -------- | ---------------------------------------------------- |
| 1   | ✅     | Build passes (`go build ./...`)   | —        | Clean. Zero issues.                                  |
| 2   | ✅     | All tests pass (`go test ./...`)  | —        | 9 packages, all green.                               |
| 3   | ✅     | Lint passes (`golangci-lint run`) | —        | 0 issues on library packages.                        |
| 4   | ✅     | Features audit completed          | —        | FEATURES.md updated with status for every component. |

---

## Critical Bugs & Type Safety

| #   | Status | Task                                                         | Priority | Notes                                                                                       |
| --- | ------ | ------------------------------------------------------------ | -------- | ------------------------------------------------------------------------------------------- |
| 5   | ✅     | Fix `class="dropdownItemClass"` rendering literal string     | P0       | Fixed in session 1.                                                                         |
| 6   | ✅     | Fix `class="emptyStateActionClass"` rendering literal string | P0       | Fixed in session 1.                                                                         |
| 7   | ✅     | Fix `.golangci.yml` Go version mismatch                      | P0       | Fixed: `1.23` → `1.26`.                                                                     |
| 8   | ✅     | Fix README.md stale `layout.BaseProps` → `PageProps`         | P0       | Fixed in session 1.                                                                         |
| 9   | ✅     | Fix `NavLinkProps.Attrs` shadowing `BaseProps.Attrs`         | —        | Done. Removed shadowing `Attrs` field from `NavLinkProps`. Consumer attributes propagate correctly. |
| 10  | ✅     | Validate required `ID` in Modal and Dropdown                 | —        | Done. `validateModalID()` and `validateDropdownID()` panic on empty ID.                             |
| 11  | ✅     | Fix Dropdown JS XSS vector                                   | —        | Done. `dropdownSafeID()` uses `strconv.Quote` for JS string escaping.                               |
| 12  | ✅     | Fix Accordion state coupling with `max-h-96`                 | —        | Done. Uses `data-open` attribute as single source of truth.                                          |

---

## Architecture

| #   | Status | Task                                                                 | Priority | Notes                                                                                                |
| --- | ------ | -------------------------------------------------------------------- | -------- | ---------------------------------------------------------------------------------------------------- |
| 13  | ✅     | Extract shared SVG helpers to `internal/svg/`                        | —        | Done. FillIcon, SpinnerSVG.                                                                          |
| 14  | ✅     | Unify alert/toast styles into shared `feedbackStyleSet`              | —        | Done. Shared struct + `lookupFeedbackStyle[T]()`.                                                    |
| 15  | ✅     | Generic `utils.MapEnum[T ~string]`                                   | —        | Done.                                                                                                |
| 16  | ✅     | All Props structs embed `utils.BaseProps`                            | —        | Done.                                                                                                |
| 17  | ✅     | Map-based style lookups (not switches)                               | —        | Done.                                                                                                |
| 18  | ✅     | Rename `layout.BaseProps` → `PageProps`                              | —        | Done.                                                                                                |
| 19  | ✅     | Deepen icon rendering: path-data map                                 | —        | Done. `iconPathData` map + `strokeIcon()`.                                                           |
| 20  | ✅     | Consolidate badge color maps into single struct map                  | —        | Done. `badgeStyleMap` with `badgeStyle{BG, Dot}` struct.                                            |
| 21  | ✅     | Merge `BadgeDefault` with `BadgeNeutral`                             | —        | Done. `BadgeDefault` removed, only `BadgeNeutral` remains.                                           |
| 22  | ✅     | Replace `Tab.Active` with `TabsProps.ActiveTabID`                    | —        | Done. Single `ActiveTabID string` on `TabsProps`, impossible state unrepresentable.                  |
| 23  | ⬜     | Unify JS attachment pattern across Accordion/Dropdown/Modal          | P2       | Three different patterns (global flag, IIFE, global functions). Standardize on IIFE-per-instance.    |
| 24  | ⬜     | Extract shared dismiss JS for Alert and Toast                        | P2       | Nearly identical event delegation pattern duplicated.                                                |
| 25  | ⬜     | Make toast icon SVG paths single-source                              | P2       | Paths duplicated in Go (`toastIconPath`) and JS (`tcToastIcons`).                                    |
| 26  | ✅     | Decouple `htmx/loading` from `feedback.Spinner`                      | —        | Done. Accepts `templ.Component` for spinner parameter.                                               |
| 26a | ⬜     | Extract tooltip position/arrow into single struct-returning function | P3       | Already uses struct map but calls lookup twice — cache result.                                       |
| 26b | ⬜     | Extract card shell CSS into `cardShellClass()`                       | P3       | Repeated 3× in Card, StatCard, SimpleCard.                                                           |
| 27  | ✅     | Replace `AvatarProps.Online/Offline bool` with `AvatarStatus` enum   | —        | Done.                                                                                                |
| 28  | ✅     | Replace `StatCard.positive bool` with `TrendDirection` enum          | —        | Done.                                                                                                |
| 29  | ✅     | Fix `HTMXSRI string` → `HTMXUseSRI bool`                             | —        | Done.                                                                                                |

---

## Accessibility

| #   | Status | Task                                                     | Priority | Notes                                                          |
| --- | ------ | -------------------------------------------------------- | -------- | -------------------------------------------------------------- |
| 30  | ⬜     | Add `alt` text to Avatar `<img>`                         | P1       | WCAG 1.1.1. Add `Alt` field to AvatarProps.                    |
| 31  | ⬜     | Add `aria-required` to form inputs when `Required: true` | P1       | WCAG requirement. Input, Select, Textarea.                     |
| 32  | ⬜     | Add `<html lang>` to Base layout                         | P1       | WCAG 3.1.1. Add `Lang` field to PageProps, default "en".       |
| 33  | ⬜     | Add Table header `scope` attributes                      | P2       | `<th scope="col">` for screen reader column association.       |
| 34  | ⬜     | Add `aria-live="polite"` to HTMX loading indicators      | P2       | Screen reader announcements for dynamic loading.               |
| 35  | ⬜     | Add `aria-live="polite"` to HTMX error handling          | P2       | Dynamic error announcements for screen readers.                |
| 36  | ⬜     | Fix `ErrorAttrs` for simultaneous error + help text      | P2       | `aria-describedby` should reference both error and help IDs.   |
| 37  | ⬜     | Scale avatar status dot with avatar size                 | P3       | Dot is fixed `h-2.5 w-2.5`, proportionally huge on XS avatars. |
| 38  | ✅     | Modal focus trap and Escape key handler                  | —        | Done in session 3.                                             |
| 39  | ✅     | Dropdown keyboard navigation                             | —        | Done in session 3.                                             |
| 40  | ✅     | Tabs ARIA linkage                                        | —        | Done in session 3.                                             |
| 41  | ✅     | Tooltip `aria-describedby` linkage                       | —        | Done in session 3.                                             |

---

## Testing

| #   | Status | Task                                               | Priority | Notes                                                                      |
| --- | ------ | -------------------------------------------------- | -------- | -------------------------------------------------------------------------- |
| 42  | ⬜     | Add BDD tests for navigation package               | P1       | No BDD tests exist for Nav, Pagination, Breadcrumbs.                       |
| 43  | ⬜     | Add BDD tests for htmx package                     | P1       | No BDD tests exist for loading indicators, error handling.                 |
| 44  | ⬜     | Add BDD tests for layout package                   | P1       | No BDD tests exist for Base, Minimal, Theme.                               |
| 45  | ⬜     | Add BDD tests for icons package                    | P2       | No BDD tests exist for Icon rendering.                                     |
| 46  | ⬜     | Add tests for Table mismatched header/row lengths  | P2       | No validation exists.                                                      |
| 47  | ⬜     | Add tests for Modal/Dropdown with empty ID         | P2       | Should fail/panic gracefully.                                              |
| 48  | ⬜     | Add test for `mapStatusToBadgeType` boundary cases | P2       | Case sensitivity, whitespace, unknown values.                              |
| 49  | ⬜     | Improve forms test coverage (58% → 75%+)           | P2       | Lowest package coverage.                                                   |
| 50  | ⬜     | Improve utils test coverage (56% → 75%+)           | P2       | MergeAttrs, CurrentYear, Deref undertested.                                |
| 51  | ⬜     | Convert snapshot tests to golden file comparison   | P2       | Current substring assertions work but golden files would be more thorough. |
| 52  | ✅     | Add a11y attribute validation tests                | —        | Done.                                                                      |
| 53  | ✅     | Add dark mode output verification tests            | —        | Done.                                                                      |
| 54  | ✅     | Add benchmark tests for hot paths                  | —        | Done. `utils.Class()` and Badge render benchmarks.                         |

---

## Dead Code & Cleanup

| #   | Status | Task                                         | Priority | Notes                                                              |
| --- | ------ | -------------------------------------------- | -------- | ------------------------------------------------------------------ |
| 55  | ⬜     | Remove or use `icons.IconAttrs`              | P2       | Exported but never called anywhere. Dead code.                     |
| 56  | ⬜     | Remove or use `internal/svg.FillIcon`        | P2       | Only referenced by `display/helpers.templ` proxy.                  |
| 57  | ⬜     | Remove no-op `DefaultXxxProps()` functions   | P3       | Several return zero-value structs (Accordion, Table, Dropdown).    |
| 58  | ⬜     | Move test helpers out of `utils/`            | P3       | `Render`, `AssertContains` etc. should be in `internal/testutil/`. |
| 59  | ⬜     | Move `display/a11y_test.go` ProgressBar test | P3       | Tests `feedback.ProgressBar` from display package.                 |
| 60  | ⬜     | Fix `examples/demo/main.go` syntax error     | P3       | Line 115: `expected operand, found '{'`.                           |

---

## DevOps & Tooling

| #   | Status | Task                                 | Priority | Notes                                                          |
| --- | ------ | ------------------------------------ | -------- | -------------------------------------------------------------- |
| 61  | ✅     | Set up GitHub Actions CI             | —        | Done. Go 1.26, lint+build+test.                                |
| 62  | ⬜     | Release automation (goreleaser)      | P3       | Tag-based releases.                                            |
| 63  | ⬜     | Fix pre-commit hook to be executable | P3       | `chmod +x scripts/pre-commit.sh` — every commit shows warning. |
| 64  | ⬜     | Exclude `examples/` from lint        | P3       | 23 issues in demo/main.go.                                     |

---

## Documentation

| #   | Status | Task                                           | Priority | Notes                                      |
| --- | ------ | ---------------------------------------------- | -------- | ------------------------------------------ |
| 65  | ✅     | Create FEATURES.md                             | —        | Comprehensive feature inventory.           |
| 66  | ✅     | Create TODO_LIST.md                            | —        | This file.                                 |
| 67  | ✅     | Create CONTEXT.md                              | —        | Architecture context.                      |
| 68  | ✅     | Update CHANGELOG.md                            | —        | Full changelog with breaking changes.      |
| 69  | ✅     | Migration guide (v0.1→v0.2)                    | —        | `docs/migration/v0.1-to-v0.2.md`.          |
| 70  | ✅     | Fix example/demo app                           | —        | Builds successfully. Showcases Nav, Alert, StatCard, Icons. |
| 71  | ⬜     | Documentation site generation                  | P3       | Auto-generated from source.                |
| 72  | ⬜     | Document `PageProps` not embedding `BaseProps` | P3       | Only Props struct that breaks convention.  |

---

## Completed (Previous Sessions)

### Session 3 (2026-05-07)

- `utils.Class()` replaces comma-join in forms, Badge, StatCard, Nav
- 17 `DefaultXxxProps()` constructors added across all packages
- Type-safe icons: `Icon string` → `icons.Name` (breaking change)
- Modal a11y: focus trap, Escape key, focus management
- Dropdown a11y: arrow key navigation, Escape to close
- Tabs a11y: proper ARIA linkage
- Tooltip a11y: deterministic id for aria-describedby
- Package doc comments for all packages
- 14 commits, 69.7% avg test coverage

### Session 2 (2026-05-04)

- Unified `feedbackStyleSet` + `lookupFeedbackStyle[T]()`
- Deepened icon rendering: 187-line switch → `iconPathData` map
- Added `AvatarStatus` enum, `TrendDirection` enum
- Fixed `HTMXSRI` → `HTMXUseSRI bool`
- Fixed ProgressBar integer division
- Added `Content templ.Component` to `TableCell`
- Implemented `TableProps.Bordered`
- Added `utils.BoolString()`
- Extensive a11y, dark mode, benchmark, XSS tests

### Session 1 (2026-05-03)

- Fixed 4 critical bugs (literal string rendering, version mismatch, stale docs)
- Semantic deduplication (13→7 clone groups)
- `layout.BaseProps` → `PageProps` rename
- Map-based style lookups
- Created all initial documentation
- Added BDD tests for display, feedback, forms packages
