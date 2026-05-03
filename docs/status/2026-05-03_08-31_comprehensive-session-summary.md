# Comprehensive Status Report — templ-components

**Date:** 2026-05-03 08:31
**Session span:** commits `eed8aa0` → `01d2cde` (17 commits, 3 phases)
**Previous reports:** `2026-05-03_07-54`, `2026-05-03_08-02`
**Branch:** `master` | **Working tree:** clean

---

## At a Glance

| Metric | Value |
|--------|-------|
| Total commits (all time) | 56 |
| Commits this session | 17 |
| Packages | 9 (display, feedback, forms, htmx, icons, internal/svg, layout, navigation, utils) |
| `.templ` files | 31 |
| `.go` files (handwritten) | 31 |
| `.go` files (generated) | 31 |
| Lines in `.templ` | 2,983 |
| Lines in `.go` (handwritten) | 2,201 |
| Net change this session | +1,461 −2,651 (net −1,190 lines) |
| Templ components | 53 |
| Typed enums | 16 |
| Icon constants | 42 |
| Clone groups (art-dupl) | 11 (structural, low severity) |
| Build | ✅ Clean |
| Vet | ✅ Clean |
| Tests | ✅ 75/75 passing |
| Coverage | **57.0%** (was ~55%) |
| CI (GitHub Actions) | ✅ Configured (Go 1.26, lint+build+test) |

---

## a) FULLY DONE

### Phase 1: Semantic Deduplication (commit `7d17413`)

- Ran `art-dupl -t 5 . --semantic --only templ` — found 13 clone groups, reduced to 7
- Extracted 10+ sub-templates:
  - `htmx/loading.templ`: 3 inline spinner SVGs → `@feedback.Spinner()`
  - `navigation/pagination.templ`: `paginationArrowIcon`, `paginationArrow`, `mobilePageButton`
  - `feedback/alert.templ`: `inlineMessage`, `inlineErrorIcon`, `inlineSuccessIcon`
  - `display/tabs.templ`: `tabLink`
  - `display/dropdown.templ`: `dropdownItemClass`, `dropdownItemLink`
  - `display/empty_state.templ`: `emptyStateActionClass`, `emptyStateAction`

### Phase 2: Architecture Improvements (commits `7ac54a4`–`b416c8d`)

- Fixed TestPtr bug: `new(v)` → `Ptr(v)` in `utils/utils_test.go`
- Renamed `layout.BaseProps` → `layout.PageProps` (eliminates collision with `utils.BaseProps`)
- Forms embed `utils.BaseProps`: InputProps, SelectProps, TextareaProps, CheckboxProps
- `feedback.ProgressBarProps` embeds `utils.BaseProps`
- Map-based style lookups: `alertStyles`, `badgeColorClass`, `badgeDotColorClass` (switch → map)
- Consolidated two SRI hash functions into single `htmxSRI(version, ext)`

### Phase 3: Shared SVG + Documentation + Tooling (commits `dc383a4`–`01d2cde`)

- Created `internal/svg` package with `FillIcon` and `SpinnerSVG` — breaks `icons → feedback` dependency
- Created `display/helpers.templ` thin wrapper delegating to `internal/svg`
- Converted `display/accordion.templ`, `display/card.templ`, `display/dropdown.templ` to use `@fillIcon()`
- Added generic `utils.MapEnum[T ~string]` for data-driven enum mapping
- Updated `display/empty_state.go` to use `MapEnum` + map
- Fixed CI Go version: 1.24 → 1.26 in `.github/workflows/ci.yaml`
- Created `FEATURES.md` — 53 components, 42 types, 42 icons
- Created `TODO_LIST.md` — 33 tracked items with priorities
- Created `CONTEXT.md` — package layout, import graph, patterns, naming conventions
- Created `docs/adr/0001-shared-svg-helpers.md` — first ADR
- Pruned 9 old status reports (kept last 2)
- Nonce propagation audit: all inline scripts verified correct
- All 17 commits pushed to origin

---

## b) PARTIALLY DONE

| Item | Status | Detail |
|------|--------|--------|
| TODO_LIST.md accuracy | **Stale** | Items 1, 3, 8, 9, 20 are marked ⬜ but were completed this session (internal/svg, MapEnum, CONTEXT.md, docs/adr/, nonce audit). Needs update. |
| Test coverage | **57%** | Good for a UI component library but `internal/svg` has 0% and several `Default*Props` constructors are untested. |
| Golden/snapshot testing | **Partial** | Snapshot tests exist for 8/9 packages but use substring assertions, not golden file comparison. |
| CI pipeline | **Working** | Build+test+lint+vet passes, but no coverage threshold enforcement, no release automation. |

---

## c) NOT STARTED

| # | Task | Priority | Notes |
|---|------|----------|-------|
| 1 | Unify `AlertType`/`ToastType` into shared `SemanticLevel` | P1 | Two identical 4-value enums with near-identical style maps |
| 2 | Add render tests for `breadcrumbs.templ` | P1 | Zero component-specific tests |
| 3 | Add render tests for `nav.templ` (Nav, SimpleNav, Footer) | P1 | Zero component-specific tests |
| 4 | Add render tests for `mobile_menu.templ` | P2 | Zero component-specific tests |
| 5 | Snapshot/golden file tests for all 30 components | P1 | Replace substring assertions |
| 6 | a11y attribute validation tests | P1 | Verify `aria-*`, `role` in rendered output |
| 7 | Dark mode output verification | P2 | Verify `dark:` Tailwind classes present |
| 8 | Cross-package circular import guard test | P2 | Automated cycle detection |
| 9 | Component composition tests | P2 | Nesting components inside each other |
| 10 | Benchmark tests | P2 | `Class()`, spinner render hot paths |
| 11 | Release automation (goreleaser) | P3 | Tag-based releases |
| 12 | Pre-commit hook for `templ generate` | P2 | Keep generated files in sync |
| 13 | Example/demo app | P2 | Showcase all 53 components |
| 14 | Documentation site generation | P3 | Auto-generated from source |
| 15 | Version migration guides | P3 | Breaking changes docs |
| 16 | `internal/svg` tests | P2 | 0% coverage — `FillIcon` and `SpinnerSVG` untested |

---

## d) TOTALLY FUCKED UP

Nothing is critically broken. Build passes, tests pass, vet is clean. However:

| Issue | Severity | Detail |
|-------|----------|--------|
| TODO_LIST.md is **stale** | Medium | 5 completed items still marked ⬜. Misleading for anyone resuming work. |
| `MapEnum` has **0% coverage** | Low | New function added but not directly tested (used only indirectly via `empty_state.go`). |
| `internal/svg` has **0% coverage** | Low | New package with no test file at all. |
| `Default*Props` constructors untested | Low | `DefaultBadgeProps`, `DefaultCardProps`, `DefaultModalProps`, `DefaultProgressBarProps`, `DefaultPageProps` — all 0%. These are simple constructors but should be verified. |
| `IconAttrs` has **0% coverage** | Low | Only exported function in icons with no test. |
| LSP shows stale errors | Cosmetic | After `.templ` edits, `templ generate` must run before LSP catches up. Not a real error — generated files are gitignored. |

---

## e) WHAT WE SHOULD IMPROVE

### High Impact

1. **Update TODO_LIST.md** — Mark 5 completed items as ✅. Add new items discovered this session. This is the #1 hygiene issue.

2. **Unify AlertType/ToastType** — Two identical enums (`Success, Error, Warning, Info`) with near-identical style maps in `feedback/alert.templ` and `feedback/toast.templ`. Extract to a single `SemanticLevel` type + single shared style map. Saves ~40 lines and eliminates a class of bugs where alert and toast styles diverge.

3. **Test `internal/svg`** — New shared package with 0% coverage. It's the foundation for `display/helpers`, `feedback/loading`, and `icons/icon`. Should have basic render tests.

4. **Golden file testing** — Current substring assertions are fragile. Golden file comparison gives deterministic, reviewable snapshots. Would catch unintended CSS class changes.

### Medium Impact

5. **Navigation test coverage** — `breadcrumbs`, `nav`, `mobile_menu` have no component-specific render tests. Navigation is critical UX — should be tested.

6. **a11y validation** — Library claims accessibility support but no tests verify `aria-*` or `role` attributes. Should be automated.

7. **Coverage threshold** — Set a minimum coverage in CI (e.g., 60%) to prevent regression.

8. **`MapEnum` direct test** — Generic function should have its own unit test with various enum types.

### Lower Impact

9. **Default props tests** — Verify all `Default*Props()` constructors return valid defaults.

10. **Benchmark suite** — `utils.Class()` is called on every render. Should benchmark to catch perf regressions.

---

## f) Top 25 Things We Should Get Done Next

| # | Task | Priority | Estimated Effort | Impact |
|---|------|----------|-----------------|--------|
| 1 | Update TODO_LIST.md — mark 5 completed items ✅, add new items | P0 | 10 min | Hygiene |
| 2 | Unify AlertType/ToastType → shared SemanticLevel + style map | P1 | 1-2 hr | Architecture |
| 3 | Add render tests for `breadcrumbs.templ` | P1 | 30 min | Coverage |
| 4 | Add render tests for `nav.templ` (Nav, SimpleNav, Footer) | P1 | 1 hr | Coverage |
| 5 | Add render tests for `mobile_menu.templ` | P2 | 30 min | Coverage |
| 6 | Add tests for `internal/svg` (FillIcon, SpinnerSVG) | P2 | 30 min | Coverage |
| 7 | Add direct unit test for `MapEnum[T]` | P2 | 15 min | Coverage |
| 8 | Add direct unit test for `IconAttrs` | P2 | 15 min | Coverage |
| 9 | Convert snapshot tests to golden file comparison | P1 | 2-3 hr | Reliability |
| 10 | Add a11y attribute validation tests | P1 | 1-2 hr | Correctness |
| 11 | Add dark mode class verification tests | P2 | 1 hr | Correctness |
| 12 | Add circular import guard test | P2 | 30 min | Safety |
| 13 | Add benchmark tests for `Class()`, spinner, fillIcon | P2 | 1 hr | Performance |
| 14 | Add component composition tests | P2 | 1 hr | Integration |
| 15 | Set coverage threshold in CI (60%) | P2 | 15 min | Quality gate |
| 16 | Add pre-commit hook for `templ generate` | P2 | 30 min | Dev experience |
| 17 | Test `Default*Props` constructors | P3 | 30 min | Coverage |
| 18 | Create example/demo app showcasing all components | P2 | 4-8 hr | Adoption |
| 19 | Add StatCard render test (currently 0% coverage) | P2 | 15 min | Coverage |
| 20 | Extract remaining structural clone groups where safe | P3 | 1-2 hr | Dedup |
| 21 | Set up goreleaser for tag-based releases | P3 | 1-2 hr | Release |
| 22 | Add SecurityHeaders render test to layout | P2 | 15 min | Security |
| 23 | Documentation site generation from source | P3 | 4-8 hr | Adoption |
| 24 | Version migration guide template | P3 | 1 hr | Adoption |
| 25 | Investigate nix flake migration | P3 | 2-4 hr | Tooling |

---

## g) Top #1 Question

**Should `SemanticLevel` (the unified AlertType/ToastType) live in `feedback/` or in `utils/`?**

Arguments:
- **`feedback/`** — Both consumers (Alert, Toast) are in feedback. It's domain-specific to user feedback severity.
- **`utils/`** — If other packages ever need semantic severity (e.g., a `display.NotificationBanner`), having it in utils avoids a cross-package import. `utils` is already the shared type hub (`BaseProps` lives there).

My recommendation: **`utils/`** — consistent with where `BaseProps` lives, and `SemanticLevel` is generic enough (success/error/warning/info) that it's not feedback-specific. But I'm not 100% certain of future usage patterns, so the owner should decide.

---

## Package Coverage Breakdown

| Package | Coverage | Tests | Notes |
|---------|----------|-------|-------|
| `display` | 62.8% | 14 | StatCard 0%, DefaultBadgeProps/CardProps/ModalProps 0% |
| `feedback` | 69.2% | 19 | DefaultProgressBarProps 0% |
| `forms` | 53.1% | 9 | Good coverage for core components |
| `htmx` | 77.3% | 7 | LoadingIndicator at 80% |
| `icons` | 7.3% | 3 | IconAttrs 0%, mostly SVG generation |
| `internal/svg` | 0.0% | 0 | **No test file** |
| `layout` | 68.3% | 6 | DefaultPageProps 0% |
| `navigation` | 71.1% | 10 | Nav at 60% |
| `utils` | 52.8% | 7 | MapEnum 0%, test helpers 0% |
| **Total** | **57.0%** | **75** | |

---

## Clone Groups (11 total, all structural)

| # | Files | Lines | Severity |
|---|-------|-------|----------|
| 1 | `mobile_menu.templ` (self) | 3 lines | Trivial — consecutive similar blocks |
| 2 | `modal.templ`, `theme.templ` (self) | 3 lines | SVG path similarity |
| 3 | `label.templ`, `pagination.templ` | 6 lines | Class string patterns |
| 4 | `breadcrumbs.templ`, `pagination.templ` | 10 lines | Navigation link pattern |
| 5 | `internal/svg/svg.templ`, `pagination.templ` | 3 lines | SVG element pattern |
| 6 | `feedback/loading.templ`, `htmx/helpers.templ` | 10 lines | Loading indicator pattern |
| 7 | `feedback/progress.templ`, `htmx/helpers.templ`, `internal/svg/svg.templ` | 8 lines | SVG element pattern |
| 8 | `nav_link.templ` (self) | 10 lines | Desktop/mobile link pattern |
| 9 | `feedback/alert.templ`, `feedback/toast.templ`, `icons/icon.templ` | 3 lines | Spinner SVG reference |
| 10 | `feedback/alert.templ`, `feedback/toast.templ`, `htmx/loading.templ`, `pagination.templ` | 4 lines | Small structural clone |
| 11 | `empty_state.templ` (self) | 8 lines | Similar action blocks |

All 11 are structural patterns — shared domain concepts (navigation links, SVG elements, loading indicators) that are not safely deduplicable without over-abstraction.

---

## Architecture Summary

```
Import graph (acyclic):
  utils          ← all packages
  internal/svg   ← display, feedback, icons
  icons          ← display (empty_state)

Key patterns:
  - Props embedding: type XProps struct { utils.BaseProps; ... }
  - Style maps: var xMap = map[XType]string{...}
  - Sub-templates: private templ functions for shared rendering within a package
  - CSP: all <script> tags use nonce={ nonce } or nonce={ props.Nonce }
  - Enums: type XxxType string + const constants
  - MapEnum[T ~string] for 1:1 string→enum; many:1 stays as switch
```

---

_Generated at 2026-05-03 08:31 by Parakletos_
