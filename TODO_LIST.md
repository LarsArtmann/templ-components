# TODO List — templ-components

**Updated:** 2026-05-04

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
| 6   | ⬜     | Unify `alertStyleSet`/`toastStyleSet` into shared `feedbackStyleSet` | P1       | Identical struct types, duplicate lookup functions, duplicate icon sub-templates. ~60 lines saved. |
| 7   | ✅     | Generic `utils.MapEnum[T ~string]`                                   | —        | Done. Used in display/empty_state.go                                                               |
| 8   | ✅     | All Props structs embed `utils.BaseProps`                            | —        | Done.                                                                                              |
| 9   | ✅     | Map-based style lookups (not switches)                               | —        | Done.                                                                                              |
| 10  | ✅     | Rename `layout.BaseProps` → `PageProps`                              | —        | Done.                                                                                              |
| 11  | ⬜     | Deepen icon rendering: path-data map + shared SVG helper             | P1       | 187-line switch → map[Name]string + strokeIcon(). ~100 lines removed.                              |
| 12  | ⬜     | Extract shared form error/aria helper                                | P2       | errorAttrs() + errorSlot() for Input, Select, Textarea, Checkbox. ~30 lines deduped.               |
| 13  | ⬜     | Replace `AvatarProps.Online/Offline bool` with `AvatarStatus` enum   | P2       | Impossible state (both true) is representable.                                                     |
| 14  | ⬜     | Replace `StatCard.positive bool` with `TrendDirection` enum          | P2       | Semantic: Up/Down/None instead of boolean.                                                         |
| 15  | ⬜     | Fix `HTMXSRI string` → `HTMXUseSRI bool` in PageProps                | P2       | Stringly-typed boolean (`"true"`).                                                                 |
| 16  | ⬜     | Fix integer division in ProgressBar percent                          | P2       | `Current * 100 / Total` truncates. Consider float64.                                               |
| 17  | ⬜     | Add `Content templ.Component` to `TableCell`                         | P2       | Currently only `Text string` — cannot render HTML in cells.                                        |
| 18  | ⬜     | Remove dead `TableProps.Bordered` field or implement styling         | P3       | Defined but never rendered.                                                                        |
| 19  | ⬜     | Remove dead `icons.IconAttrs` or add tests/usage                     | P3       | Only exported function in icons, 0% coverage, unused.                                              |
| 20  | ✅     | Add `CONTEXT.md` with architecture decisions                         | —        | Done.                                                                                              |
| 21  | ✅     | Add `docs/adr/` for architecture decision records                    | —        | Done. ADR-0001.                                                                                    |

## Testing

| #   | Status | Task                                              | Priority | Notes                                             |
| --- | ------ | ------------------------------------------------- | -------- | ------------------------------------------------- |
| 22  | ⬜     | Add render tests for breadcrumbs                  | P1       | No test coverage.                                 |
| 23  | ⬜     | Add render tests for nav (Nav, SimpleNav, Footer) | P1       | No test coverage.                                 |
| 24  | ⬜     | Add render tests for mobile_menu                  | P2       | No test coverage.                                 |
| 25  | ⬜     | Add render tests for htmx error_handling          | P2       | No test coverage.                                 |
| 26  | ⬜     | Convert snapshot tests to golden file comparison  | P1       | Currently substring-based assertions.             |
| 27  | ⬜     | Add a11y attribute validation tests               | P1       | Verify `aria-*`, `role` attributes.               |
| 28  | ⬜     | Add dark mode output verification tests           | P2       | Verify `dark:` classes present.                   |
| 29  | ⬜     | Add component composition tests                   | P2       | Nesting components inside each other.             |
| 30  | ⬜     | Add benchmark tests for hot paths                 | P2       | `Class()`, spinner render.                        |
| 31  | ✅     | Add tests for `internal/svg` package              | —        | Done. FillIcon and SpinnerSVG render tests added. |
| 32  | ✅     | Add direct test for `utils.MapEnum[T]`            | —        | Done. 3 subtests: found, missing, empty key.      |
| 33  | ⬜     | Test all `Default*Props()` constructors           | P2       | 5 constructors with 0% coverage.                  |
| 34  | ✅     | Fix TestPtr bug                                   | —        | `new(v)` → `Ptr(v)`.                              |

## Security & CSP

| #   | Status | Task                                                  | Priority | Notes                                                       |
| --- | ------ | ----------------------------------------------------- | -------- | ----------------------------------------------------------- |
| 35  | ✅     | Nonce propagation audit                               | —        | All inline scripts verified.                                |
| 36  | ⬜     | Add `SecurityHeaders` test to layout                  | P2       | Verify meta tags rendered.                                  |
| 37  | ✅     | CSP compliance for all inline scripts                 | —        | All scripts use nonce.                                      |
| 38  | ⬜     | Verify dropdown JS template interpolation is XSS-safe | P2       | `'{{ props.ID }}'` in script block — verify templ escaping. |

## DevOps & Tooling

| #   | Status | Task                                 | Priority | Notes                                             |
| --- | ------ | ------------------------------------ | -------- | ------------------------------------------------- |
| 39  | ✅     | Set up GitHub Actions CI             | —        | Done. Go 1.26, lint+build+test.                   |
| 40  | ⬜     | Release automation (goreleaser)      | P3       | Tag-based releases.                               |
| 41  | ⬜     | Investigate nix flake migration      | P3       | No build system exists.                           |
| 42  | ⬜     | Pre-commit hook for `templ generate` | P2       | Ensure generated files stay in sync.              |
| 43  | ⬜     | Add `layout/sri.go` package comment  | P3       | golangci-lint `revive: package-comments` warning. |

## Documentation

| #   | Status | Task                                 | Priority | Notes                            |
| --- | ------ | ------------------------------------ | -------- | -------------------------------- |
| 44  | ✅     | Create FEATURES.md                   | —        | Comprehensive feature inventory. |
| 45  | ✅     | Create TODO_LIST.md                  | —        | This file.                       |
| 46  | ⬜     | Update CHANGELOG.md with recent work | P1       | Still shows generic v0.1.0.      |
| 47  | ⬜     | Create example/demo app              | P2       | Showcase all components.         |
| 48  | ⬜     | Documentation site generation        | P3       | Auto-generated from source.      |
| 49  | ⬜     | Version migration guides             | P3       | Breaking changes documentation.  |

## Deduplication

| #   | Status | Task                                            | Priority | Notes                                       |
| --- | ------ | ----------------------------------------------- | -------- | ------------------------------------------- |
| 50  | ✅     | Semantic deduplication (13→7 clone groups)      | —        | Extracted 10+ sub-templates.                |
| 51  | ⬜     | Remaining 9 clone groups (dupl)                 | P3       | All in test files — structural only.        |
| 52  | ⬜     | Move `boolString()` to `utils/` and standardize | P3       | Inconsistent with `fmt.Sprintf("%t", ...)`. |

---

## Completed This Session (2026-05-04)

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
- Total: 332 tests passing (up from 66 top-level + subtests)
