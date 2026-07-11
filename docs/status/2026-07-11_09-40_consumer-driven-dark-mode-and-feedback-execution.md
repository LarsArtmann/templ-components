# Status Report: Consumer-Driven Dark Mode & Feedback Execution

**Date:** 2026-07-11 09:40
**Session scope:** Dark mode research → packaging fixes → consumer feedback execution → self-review → remediation
**Commits this session:** 8 (`1a078c4` through `865a9e4`)

---

## a) FULLY DONE

### Dark Mode Research & Documentation

- ✅ `docs/dark-mode-research.md` — comprehensive first-principles analysis of all dark mode mechanisms in Tailwind v4 + modern CSS (`prefers-color-scheme`, `color-scheme`, `light-dark()`, `@theme` vs `@theme inline`, explicit `dark:` variants vs token indirection). Includes evaluation matrix, three consumer paths with code, "what NOT to do" section, and action items (all marked complete).
- ✅ Cross-linked from 4 locations: adoption guide, ADR 0011, feedback file, AGENTS.md.

### Theme.css Packaging Fixes (1% tier — 51% impact)

- ✅ Fixed `color-scheme: light` → `color-scheme: light dark` on `:root` — native form controls were rendering light-only in dark mode for ALL consumers.
- ✅ Split `@custom-variant dark` into a clearly commented opt-in section with three options documented inline (OS-following, toggle, CSS-variable). Previously, importing theme.css for color overrides silently forced class-based dark mode.
- ✅ Removed stale `prefers-color-scheme` fallback block (redundant with `color-scheme: light dark`).

### Adoption Guide Rewrite (4% tier — 64% impact)

- ✅ "Setting Class on components" section — documents Go promoted-field struct literal gotcha.
- ✅ "Vendoring templ-components" rewritten — Go module cache `@source` path pattern + `go list -m` tip + troubleshooting.
- ✅ "Theming" section — `@theme` palette override pattern + CSS-variable design system mapping (`--color-white: var(--surface)`).
- ✅ "Dark mode strategies" — three first-class paths (OS-following, toggle, CSS-variable) with code examples and comparison table.
- ✅ "Default constructors" section — documents `DefaultXxxProps()` pattern.
- ✅ `prefers-color-scheme` consumer implications note.
- ✅ Fixed "Why Tailwind v4?" table row (dark mode is default `prefers-color-scheme`, not class-based).
- ✅ Cross-linked `theme-bridge.md` recipe from Theming section.

### Component API Improvements (20% tier — 80% impact)

- ✅ `GridColsAutoFit` + `MinColWidth` — CSS `auto-fit`/`minmax()` grid template for container-width-responsive layouts.
- ✅ `Card.Header` slot — `templ.Component` that replaces entire default header section.
- ✅ `CardPaddingNone` fix — children/body now render without wrapping padding `<div>` for table-in-card layouts.

### Bug Fixes

- ✅ `ContainerResponsive` + `GridColsAutoFit` precedence logic — auto-fit now takes precedence (was silently falling back to default grid).
- ✅ `errorpage/handler.go` — restored `encoding/json` v1 twice (reverted by parallel session, fixed again).
- ✅ `ThemeScript`/`ThemeToggle` godoc — clarified toggle-only usage.

### Tests

- ✅ 11 new tests (10 in `grid_card_feedback_test.go`, 1 in `enums_test.go`).
- ✅ 2 golden tests (`card_header_slot`, `grid_autofit`) with golden files.
- ✅ All existing tests pass, dark mode compliance tests pass, CSP nonce tests pass.

### Process Compliance

- ✅ CHANGELOG `[Unreleased]` updated with all 11 entries (Added/Fixed/Changed sections).
- ✅ FEATURES.md updated (Card, Grid, GridCols rows).
- ✅ AGENTS.md updated with Grid/Header/CardPaddingNone changes.
- ✅ `dark-mode-research.md` action items marked as complete.

### Verification

- ✅ `templ generate` + `go build` + `go test` + `golangci-lint` all green.
- ✅ BuildFlow pre-commit 26/26 passed on every commit.

---

## b) PARTIALLY DONE

### CHANGELOG `[Unreleased]` is warm but NOT released

- The `[Unreleased]` section has 11 entries from this session plus 5 from a prior session. These should be released as v0.15.0 when ready. The version in `utils/version.go` is still `0.14.0`.
- **What's missing:** A release decision (not my call — user decides when to cut).

### Adoption guide dark mode section

- Three paths are documented, but the code examples for Path 2 (toggle) and Path 3 (CSS-variable) are somewhat terse. A consumer brand-new to Tailwind v4 might need more hand-holding. The cross-link to `dark-mode-research.md` covers the gap for now.

---

## c) NOT STARTED

### From the original cqrs-htmx feedback (intentionally deferred)

- ❌ **ADR 0008 Semantic Token Layer Phase 2** — deferred to v1.0 by ADR decision. Components still emit `bg-blue-600`, not `bg-tc-primary`.
- ❌ **BuildFlow gitignore `*_templ.go` re-append** — upstream issue in `larsartmann/buildflow`. Documented in AGENTS.md but not fixed.
- ❌ **Remove deprecated `FamilyFromErrorFamily` alias** — scheduled for v1.0.
- ❌ **Modular workspace split** — post-v1.0 per AGENTS.md.
- ❌ **`go.mod` templ version bump to v0.3.1036** — waiting for official upstream release.
- ❌ **`goBackScript` promotion to utils** — triggered when 2nd package needs it.
- ❌ **`overlayShellProps` sub-struct refactor** — triggered when 3rd overlay type emerges.
- ❌ **pkg.go.dev API reference** — mentioned in feedback as a discoverability gap. pkg.go.dev auto-generates from godoc; no action needed beyond ensuring godoc is good.

### From the dark mode research (intentionally rejected)

- ❌ **Mandatory semantic token layer** — rejected by ADR-001, confirmed by research.
- ❌ **`light-dark()` CSS function adoption** — baseline 2024, not Tailwind-expressible, revisit 2027+.
- ❌ **Surface token redesign (`--tc-surface`)** — `@theme` palette override already solves this.

---

## d) TOTALLY FUCKED UP

### Build regression from parallel session

- **What happened:** My commit `0362b2f` fixed `encoding/json/v2` → `encoding/json` in `errorpage/handler.go`. A parallel session (commit `743057f`) then reverted this fix — adding back `encoding/json/v2` and `encoding/json/jsontext` imports with only a whitespace formatting change as the visible diff. This broke the build for any consumer without `GOEXPERIMENT=jsonv2`.
- **Impact:** Build was broken on `master` from `743057f` until I caught it in the self-review and fixed it in `d5cbc60`.
- **Root cause:** The parallel session didn't read my commit before editing the same file. The `encoding/json/v2` prohibition in AGENTS.md was not respected.
- **Lesson:** When multiple sessions touch the same file, check `git log` for recent commits before editing. The `utils/jsonv2_guard_test.go` exists but it's a test-time check, not a build-time check — the build itself breaks.

### `ContainerResponsive` + `GridColsAutoFit` logic bug

- **What happened:** When I added `GridColsAutoFit`, I placed the auto-fit check before the `ContainerResponsive` check, but both were `if` statements (not `else if`). So when both were set, `ContainerResponsive` silently overrode auto-fit, calling `gridContainerClass(GridColsAutoFit)` which isn't in the container lookup map and fell back to the default 3-column grid.
- **Impact:** Consumers setting both flags would get a 3-column grid instead of auto-fit. No crash, just wrong output.
- **Root cause:** I didn't test the combined case. Caught in self-review.
- **Lesson:** Always test flag combinations, not just individual flags.

---

## e) WHAT WE SHOULD IMPROVE

### Process

1. **Parallel session coordination** — The json/v2 regression happened because a parallel session edited a file without checking recent commits. Consider a pre-edit `git log -1 -- <file>` check in the workflow.
2. **CHANGELOG discipline** — I added features without updating CHANGELOG on the same commit. The AGENTS.md rule says `[Unreleased]` must be warm at all times. I only caught this in self-review.
3. **Test combined flag states** — The `ContainerResponsive` + `GridColsAutoFit` bug would have been caught by a combinatorial test. Add flag-combination tests for all components with multiple boolean fields.
4. **Golden tests for every new feature** — I initially forgot golden tests for the new Grid/Card features. Should be part of the definition of done.

### Architecture

5. **`encoding/json/v2` keeps recurring** — This is the 4th time it's been introduced and had to be removed. The guard test (`utils/jsonv2_guard_test.go`) catches it at test time, but the build breaks before tests run. Consider a build-tag-based guard or a `go.mod` replace directive that makes the import fail at resolve time.
6. **`gridContainerClass` falls back silently for unknown values** — `GridColsAutoFit` isn't in the container lookup map, so `gridContainerClass(GridColsAutoFit)` returns the default. This is graceful degradation, but for a value that SHOULD never be passed to that function, it's a silent logic bug. The `else if` fix prevents it, but the function's API is misleading.

### Documentation

7. **Adoption guide is getting long** — 380+ lines. Consider splitting into a quick-start guide and a reference guide.
8. **Dark mode documentation is spread across 4 files** — `dark-mode-research.md`, `tailwind-v4-adoption-guide.md`, `adr/0011-dark-mode-convention.md`, `templ-components-theme.css`. A single "Dark Mode" index page that links to all of them would help.

---

## f) Up to 50 Things We Should Get Done Next

### Critical / High Impact

| #   | Task                                                                      | Impact | Effort |
| --- | ------------------------------------------------------------------------- | ------ | ------ |
| 1   | Cut v0.15.0 release (11 CHANGELOG entries waiting)                        | High   | 10min  |
| 2   | Add build-time json/v2 guard (build tag or go.mod replace)                | High   | 30min  |
| 3   | Add combinatorial flag tests for all components with 2+ bool fields       | High   | 60min  |
| 4   | Audit all components for `ContainerResponsive`-style silent fallback bugs | Medium | 30min  |

### Component Improvements

| #   | Task                                                            | Impact | Effort |
| --- | --------------------------------------------------------------- | ------ | ------ |
| 5   | Add `DefaultGridProps` doc example showing `MinColWidth` usage  | Low    | 5min   |
| 6   | Add `Card.Header` usage example in godoc                        | Low    | 5min   |
| 7   | Consider `Card.TitleElement` field (h1/h2/h3/h4 selection)      | Medium | 30min  |
| 8   | Add `GridColsAutoFit` to the demo binary                        | Low    | 10min  |
| 9   | Add `Card.Header` to the demo binary                            | Low    | 10min  |
| 10  | Consider `GridProps.MaxColWidth` for `minmax(min, max)` pattern | Low    | 15min  |

### Documentation

| #   | Task                                                                      | Impact | Effort |
| --- | ------------------------------------------------------------------------- | ------ | ------ |
| 11  | Create `docs/dark-mode-index.md` linking all dark mode docs               | Medium | 15min  |
| 12  | Split adoption guide into quick-start + reference                         | Medium | 60min  |
| 13  | Add "Migration from v0.14 to v0.15" guide                                 | Low    | 15min  |
| 14  | Document all `DefaultXxxProps()` constructors in one table                | Low    | 20min  |
| 15  | Add consumer guide section: "Choosing a dark mode strategy" decision tree | Low    | 15min  |
| 16  | Audit godoc across all packages for completeness                          | Medium | 60min  |
| 17  | Add `CHANGELOG.md` entry for the `@source` GOMODCACHE docs                | Low    | 2min   |

### Testing

| #   | Task                                                             | Impact | Effort |
| --- | ---------------------------------------------------------------- | ------ | ------ |
| 18  | Add fuzz test for `gridAutoFitClass()` with edge-case widths     | Low    | 10min  |
| 19  | Add fuzz test for `GridColsAutoFit` + `MinColWidth` combinations | Low    | 10min  |
| 20  | Add test: `Card.Header` + `Card.Body` both set renders both      | Low    | 5min   |
| 21  | Add test: `CardPaddingNone` + `Card.Body` slot works correctly   | Low    | 5min   |
| 22  | Add test: `CardPaddingNone` + `Card.Header` slot works correctly | Low    | 5min   |
| 23  | Add dark mode golden test for `GridColsAutoFit`                  | Low    | 10min  |
| 24  | Add dark mode golden test for `Card.Header` slot                 | Low    | 10min  |
| 25  | Benchmark `gridAutoFitClass` string concatenation                | Low    | 5min   |

### Architecture / Type Model

| #   | Task                                                                                                             | Impact | Effort |
| --- | ---------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 26  | Consider typed `MinColWidth` type (e.g., `GridMinWidth string` with `IsValid()`)                                 | Low    | 15min  |
| 27  | Consider `CardPaddingNone` as a separate code path (not a padding enum value)                                    | Low    | 20min  |
| 28  | Audit all enum lookups for silent-fallback-on-unknown behavior                                                   | Medium | 30min  |
| 29  | Consider a `GridMode` enum (`GridModeFixed` / `GridModeAutoFit` / `GridModeContainer`) instead of separate bools | Medium | 45min  |

### Maintenance / Cleanup

| #   | Task                                                                      | Impact | Effort |
| --- | ------------------------------------------------------------------------- | ------ | ------ |
| 30  | Remove the stale LSP diagnostic for `handler.go` (restart gopls)          | Low    | 1min   |
| 31  | Check `.gitignore` for BuildFlow `*_templ.go` re-append after each commit | Low    | 1min   |
| 32  | Audit all `*_templ.go` files are committed (none hidden by gitignore)     | Low    | 5min   |
| 33  | Run `go mod tidy` to clean up any unused dependencies                     | Low    | 5min   |
| 34  | Update `docs/planning/` plan to mark completed items                      | Low    | 10min  |
| 35  | Clean up `docs/status/` — archive old status reports                      | Low    | 10min  |

### Consumer Experience

| #   | Task                                                              | Impact | Effort |
| --- | ----------------------------------------------------------------- | ------ | ------ |
| 36  | Create a "Quick start" guide (5-line getting started)             | High   | 15min  |
| 37  | Add a "Common patterns" cookbook section                          | Medium | 45min  |
| 38  | Document the `icons.IconPathData` icons-only adoption path        | Medium | 20min  |
| 39  | Add an interactive demo page for dark mode toggle vs OS-following | Low    | 30min  |
| 40  | Create a "Theming cheat sheet" (1-page reference)                 | Medium | 20min  |

### Error Page Package

| #   | Task                                                                              | Impact | Effort |
| --- | --------------------------------------------------------------------------------- | ------ | ------ |
| 41  | Add integration test for `writeJSONError` with json v1                            | Medium | 15min  |
| 42  | Consider using `encoding/json` `Marshal` instead of `Encoder` for error responses | Low    | 10min  |
| 43  | Add test verifying JSON error response shape matches HTMX expectations            | Medium | 20min  |

### Feedback Response

| #   | Task                                                                      | Impact | Effort |
| --- | ------------------------------------------------------------------------- | ------ | ------ |
| 44  | Respond to cqrs-htmx feedback with the `@theme` palette override solution | Medium | 10min  |
| 45  | Update the feedback file with "addressed" annotations                     | Low    | 10min  |
| 46  | Consider a "Consumer Adoption Log" documenting who adopted what and when  | Low    | 15min  |

### Research / Investigation

| #   | Task                                                                                     | Impact | Effort |
| --- | ---------------------------------------------------------------------------------------- | ------ | ------ |
| 47  | Investigate Tailwind v4 `@source` with Go import paths (feature request?)                | Low    | 30min  |
| 48  | Research `light-dark()` browser support timeline for adoption readiness                  | Low    | 15min  |
| 49  | Investigate whether `@theme inline` would simplify the theme.css file                    | Low    | 20min  |
| 50  | Research whether other Go templ UI libraries have solved the dark mode packaging problem | Low    | 20min  |

---

## g) Top 2 Questions I Cannot Answer Myself

### Q1: Should we cut v0.15.0 now, or batch with more features?

The `[Unreleased]` section has 16 entries (11 from this session, 5 from prior). The changes include new component APIs (`GridColsAutoFit`, `Card.Header`, `CardPaddingNone`), critical bug fixes (`color-scheme`, json/v2 regression), and significant documentation rewrites. This feels like a minor version bump. But I don't know if you want to batch more features before cutting, or if you have a release cadence in mind.

### Q2: Should the `encoding/json/v2` problem be solved with a build-tag guard or a `go.mod` replace directive?

The json/v2 import has been introduced and removed 4 times now. The test guard catches it at test time, but the build breaks before tests run. A build-tag-based approach (e.g., a `//go:build !jsonv2` file that imports `encoding/json/v2` and fails) would catch it at build time. Alternatively, a `go.mod` replace directive that redirects `encoding/json/v2` to a non-existent module would make `go mod tidy` fail. I don't know which approach you prefer, or if there's a better solution I'm not seeing. This is a tooling decision that depends on your BuildFlow setup.
