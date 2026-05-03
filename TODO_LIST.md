# TODO List — templ-components

**Updated:** 2026-05-03

Legend: ✅ Done | 🔨 In Progress | ⬜ Not Started | ❌ Blocked

---

## Architecture

| # | Status | Task | Priority | Notes |
|---|--------|------|----------|-------|
| 1 | ✅ | Extract shared SVG helpers (`fillIcon`, spinner) to `internal/svg/` package | — | Done — commit dc383a4. Breaks `icons → feedback` dependency |
| 2 | ⬜ | Unify `AlertType`/`ToastType` into shared `SemanticLevel` type | P1 | Two identical enums with near-identical style maps |
| 3 | ✅ | Generic `utils.MapEnum[T ~string](map, fallback, key) T` | — | Done — commit 82ed623. Used in `display/empty_state.go` |
| 4 | ✅ | All Props structs embed `utils.BaseProps` | — | Forms and ProgressBarProps were last, now done |
| 5 | ✅ | Map-based style lookups (not switches) | — | alertStyles, badgeColorClass, badgeDotColorClass converted |
| 6 | ✅ | Rename `layout.BaseProps` → `PageProps` | — | Eliminates name collision with `utils.BaseProps` |
| 7 | ⬜ | Cross-package circular import guard test | P4 | Verify `icons → feedback` is one-directional |
| 8 | ✅ | Add `CONTEXT.md` with architecture decisions | — | Done — commit 01d2cde. Package layout, import graph, patterns |
| 9 | ✅ | Add `docs/adr/` for architecture decision records | — | Done — commit 01d2cde. ADR-0001: shared SVG package |

## Testing

| # | Status | Task | Priority | Notes |
|---|--------|------|----------|-------|
| 10 | ⬜ | Add render tests for breadcrumbs | P1 | No test coverage at all |
| 11 | ⬜ | Add render tests for nav (Nav, SimpleNav, Footer) | P1 | No test coverage |
| 12 | ⬜ | Add render tests for mobile_menu | P2 | No test coverage |
| 13 | ⬜ | Add render tests for htmx error_handling | P2 | No test coverage |
| 14 | ⬜ | Add snapshot/golden file tests for all 30 components | P1 | Currently substring-based assertions |
| 15 | ⬜ | Add a11y attribute validation tests | P1 | Verify `aria-*`, `role` attributes |
| 16 | ⬜ | Add dark mode output verification tests | P2 | Verify `dark:` classes present |
| 17 | ⬜ | Add component composition tests | P2 | Nesting components inside each other |
| 18 | ⬜ | Add benchmark tests for hot paths | P2 | `Class()`, spinner render |
| 19 | ✅ | Fix TestPtr bug | — | `new(v)` → `Ptr(v)` |

## Security & CSP

| # | Status | Task | Priority | Notes |
|---|--------|------|----------|-------|
| 20 | ✅ | Nonce propagation audit | — | All inline scripts verified — all use nonce attribute correctly |
| 21 | ⬜ | Add `SecurityHeaders` test to layout | P2 | Verify meta tags rendered when `SecurityHeaders=true` |
| 22 | ✅ | CSP compliance for all inline scripts | — | All scripts use `nonce` attribute |

## DevOps & Tooling

| # | Status | Task | Priority | Notes |
|---|--------|------|----------|-------|
| 23 | ✅ | Set up GitHub Actions CI | — | Done — Go 1.26, lint+build+test, coverage artifact |
| 24 | ⬜ | Release automation (goreleaser) | P3 | Tag-based releases |
| 25 | ⬜ | Investigate nix flake migration | P3 | No build system exists |
| 26 | ⬜ | Pre-commit hook for `templ generate` | P2 | Ensure generated files stay in sync |

## Documentation

| # | Status | Task | Priority | Notes |
|---|--------|------|----------|-------|
| 27 | ✅ | Create FEATURES.md | — | Comprehensive feature inventory |
| 28 | ✅ | Create TODO_LIST.md | — | This file |
| 29 | ⬜ | Create example/demo app | P2 | Showcase all components |
| 30 | ⬜ | Documentation site generation | P3 | Auto-generated from source |
| 31 | ⬜ | Version migration guides | P3 | Breaking changes documentation |

## Deduplication

| # | Status | Task | Notes |
|---|--------|------|-------|
| 32 | ✅ | Semantic deduplication (13→7 clone groups) | Extracted 10+ sub-templates |
| 33 | ⬜ | Remaining 7 clone groups | Structural only — not safely deduplicable |

---

## Completed This Session (2026-05-03, 17 commits)

- Semantic deduplication: 13→7→11 clone groups (extracted 10+ sub-templates)
- Fix TestPtr bug (was testing `new()` not `Ptr()`)
- Rename `layout.BaseProps` → `PageProps`
- Forms embed `utils.BaseProps` (InputProps, SelectProps, TextareaProps, CheckboxProps)
- ProgressBarProps embeds `utils.BaseProps`
- Switch→map style lookups (alertStyles, badgeColorClass, badgeDotColorClass)
- Consolidate SRI hash functions into single `htmxSRI()`
- Extract shared SVG helpers to `internal/svg/` package (FillIcon, SpinnerSVG)
- Add generic `utils.MapEnum[T ~string]` for data-driven enum mapping
- Fix CI Go version: 1.24 → 1.26
- Create FEATURES.md (53 components, 42 types, 42 icons)
- Create TODO_LIST.md (33 items)
- Create CONTEXT.md (architecture, patterns, import graph)
- Create docs/adr/0001-shared-svg-helpers.md (first ADR)
- Nonce propagation audit: all inline scripts verified correct
- Prune 9 old status reports
- All pushed to origin
