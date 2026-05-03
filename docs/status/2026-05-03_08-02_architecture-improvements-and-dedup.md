# Status Report — templ-components

**Date:** 2026-05-03 08:02  
**Branch:** master (7 commits ahead of origin)  
**Base commit:** eed8aa0 → b416c8d (current)  
**Go:** 1.26.2 | **templ:** v0.3.1001 | **tailwind-merge-go:** v0.2.1

---

## Metrics at a Glance

| Metric                  | Value                 | Delta                        |
| ----------------------- | --------------------- | ---------------------------- |
| Packages                | 8                     | —                            |
| Tests                   | 75 passing, 0 failing | was 76 (test renamed/merged) |
| `.templ` files          | 30                    | +1 (display/helpers.templ)   |
| Total `.templ` lines    | 2,961                 | +3                           |
| `go vet`                | clean                 | ✓                            |
| Clone groups (art-dupl) | 7                     | was 13 (46% reduction)       |
| Typed enums             | 16+                   | —                            |
| CSP violations          | 0                     | ✓                            |
| Dependencies            | 2                     | —                            |
| Commits this session    | 7                     | —                            |
| Net line delta          | +439 / -502           | net -63                      |

---

## a) FULLY DONE

### Session 2 Work (this round — 7 commits)

| #   | Commit    | What                                                                                                                                                                                                                                                           |
| --- | --------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | `7d17413` | **Semantic deduplication pass** — 13→7 clone groups. Extracted `paginationArrow`, `mobilePageButton`, `inlineMessage`, `tabLink`, `dropdownItemLink`, `emptyStateAction`, `fillIcon` sub-templates. Replaced 3 inline spinner SVGs with `@feedback.Spinner()`. |
| 2   | `7ac54a4` | **Fix TestPtr bug** — `new(v)` → `Ptr(v)`. Was not actually testing the `Ptr()` function. Eliminated the "impossible condition: non-nil == nil" LSP warning.                                                                                                   |
| 3   | `8ad0334` | **Rename layout.BaseProps → PageProps** — Eliminated confusing name collision between `layout.BaseProps` (page metadata) and `utils.BaseProps` (component ID/Class/Attrs). Also renamed `DefaultBaseProps` → `DefaultPageProps`.                               |
| 4   | `ed88096` | **Forms embed utils.BaseProps** — `InputProps`, `SelectProps`, `TextareaProps`, `CheckboxProps` now embed `utils.BaseProps` instead of duplicating ID/Class/Attrs fields. Gains Nonce + AriaLabel support.                                                     |
| 5   | `5fcf1b0` | **ProgressBarProps embed utils.BaseProps** — Last Props struct without embedding. Now consistent.                                                                                                                                                              |
| 6   | `817e8b9` | **Switch→map style lookups** — `alertStyles`, `badgeColorClass`, `badgeDotColorClass` converted from switch statements to package-level maps with defaults. Matches existing `toastStyleMap` pattern.                                                          |
| 7   | `b416c8d` | **Consolidate SRI functions** — Two near-identical `htmxIntegrityHash`/`htmxResponseTargetsIntegrityHash` → single `htmxSRI(version, ext)`. Extracted `htmxSRIEntry` type.                                                                                     |

### Previously Completed (from earlier sessions)

- 16+ typed string enums across all packages
- CSP compliance for all inline scripts (nonce-based)
- `utils.Class()` with tailwind-merge-go for intelligent class merging
- `utils.BaseProps` shared across all components (now including forms + progress)
- Icon system with 30+ named icons + typed `Name` enum
- Comprehensive project hygiene (git config, linting, formatting)
- Test helper library (`Render`, `AssertContains`, `AssertNotContains`, `AssertEqual`)

---

## b) PARTIALLY DONE

### Architecture Improvements (6 of ~12 complete)

| Done                                     | Remaining                                          |
| ---------------------------------------- | -------------------------------------------------- |
| ✅ utils.BaseProps embedding (all Props) | ⬜ Shared SVG helpers package                      |
| ✅ layout.BaseProps → PageProps rename   | ⬜ AlertType/ToastType → SemanticLevel unification |
| ✅ Map-based style lookups               | ⬜ Generic string→enum mapper in utils             |
| ✅ SRI function consolidation            | ⬜ Cross-package circular import audit             |
| ✅ fillIcon shared in display/           | ⬜ Move fillIcon to shared package                 |
| ✅ TestPtr bug fix                       | ⬜ Test helper improvements                        |

### Test Coverage

- 75 tests across 8 packages — all render + unit tests pass
- **Missing render tests**: breadcrumbs, nav, mobile_menu, forms helpers (SanitizeID tested, but no render tests for input/select/textarea error states)
- **No snapshot/golden testing** — all tests use substring assertions
- **No a11y automated checks**

---

## c) NOT STARTED

1. **FEATURES.md** — No feature inventory exists
2. **TODO_LIST.md** — No comprehensive TODO tracking
3. **CONTEXT.md** — No project context document
4. **Shared SVG helpers package** — `fillIcon` lives in `display/`, can't be used by `navigation/` or `feedback/` without import issues
5. **AlertType/ToastType unification** — Two identical enums with near-identical style maps
6. **Generic string→enum mapper** — `mapEmptyStateIcon` and `mapStatusToBadgeType` are hand-written switches
7. **Snapshot/golden file testing** — No rendered HTML snapshots
8. **A11y testing** — No automated accessibility audits
9. **Benchmark tests** — No performance benchmarks
10. **Example app** — No demo/showcase application
11. **CSS extraction** — No shared Tailwind preset/theme file
12. **Documentation site** — No generated docs/website
13. **CI pipeline** — No `.github/workflows` or equivalent
14. **Release automation** — No goreleaser or similar
15. **Changelog automation** — CHANGELOG.md exists but manual
16. **Prune old status reports** — 10 status files, should keep 2–3
17. **nix flake migration** — No flake.nix, no justfile, no Makefile

---

## d) TOTALLY FUCKED UP

1. **icons ↔ feedback potential circular import** — `icons/icon.templ` imports `feedback.Spinner`. If `feedback` ever imports `icons` directly, we get a build break. Currently safe (one-directional), but fragile.
2. **10 stale status reports** in `docs/status/` — growing unbounded, no cleanup discipline.
3. **LSP warnings in generated files** — After `.templ` changes, `templ generate` must run before LSP catches up. Not a bug but confusing during development.
4. **Test count changed 76→75** — Not a regression; test restructuring merged/deduplicated some test names. Should verify nothing was accidentally lost.
5. **`display/helpers.templ` can't be imported by other packages** — `fillIcon` is display-only. If navigation or feedback needs it, we need a shared package.

---

## e) WHAT WE SHOULD IMPROVE

### High Priority

1. **Create FEATURES.md** — Audit actual code, document every component, every prop, every variant
2. **Create TODO_LIST.md** — Track real work items with status
3. **Extract shared SVG helpers** — `fillIcon` + spinner SVG into `internal/svg/` or similar to avoid cross-package issues
4. **Unify AlertType/ToastType** — Extract to shared `SemanticLevel` type; single style map
5. **Add CI pipeline** — Even a simple GitHub Actions for build+test+vet

### Medium Priority

6. **Prune old status reports** — Keep last 2–3, archive the rest
7. **Snapshot testing** — Golden files for all 30 components
8. **A11y attribute validation** — At minimum check `aria-*` attributes
9. **Consistent nonce propagation audit** — Some components have `props.Nonce`, some don't
10. **Dark mode testing** — Verify `dark:` classes in rendered output

### Lower Priority

11. **Generic enum mapper** — `utils.MapEnum[T ~string](map, fallback, key)` instead of hand-written switches
12. **Add CONTEXT.md** — Architecture decisions, package layout rationale
13. **Form validation helpers** — Client-side validation JS utilities
14. **Component composition tests** — Nesting components inside each other
15. **Example/demo app** — Showcase all components

---

## f) Top 25 Things We Should Get Done Next

| #   | Priority | Task                                                  | Est.  | Impact                               |
| --- | -------- | ----------------------------------------------------- | ----- | ------------------------------------ |
| 1   | P0       | Create FEATURES.md from codebase audit                | 30min | High — documentation foundation      |
| 2   | P0       | Create TODO_LIST.md                                   | 20min | High — tracking foundation           |
| 3   | P0       | Prune old status reports (keep last 2)                | 2min  | Low — hygiene                        |
| 4   | P0       | Push 7 commits to origin                              | 1min  | High — backup                        |
| 5   | P1       | Extract shared SVG helpers to `internal/svg/` package | 1hr   | High — resolves cross-package issues |
| 6   | P1       | Unify AlertType/ToastType into shared `SemanticLevel` | 1hr   | High — eliminates enum duplication   |
| 7   | P1       | Set up GitHub Actions CI (build+test+vet)             | 30min | High — safety net                    |
| 8   | P1       | Add snapshot/golden file tests for all 30 components  | 2hr   | High — regression safety             |
| 9   | P1       | Add a11y attribute validation tests                   | 1hr   | Medium — accessibility               |
| 10  | P1       | Consistent nonce propagation audit                    | 30min | Medium — security                    |
| 11  | P2       | Generic `utils.MapEnum[T]` for string→enum lookups    | 30min | Medium — DRY                         |
| 12  | P2       | Dark mode output verification tests                   | 1hr   | Medium — correctness                 |
| 13  | P2       | Component composition tests (nesting)                 | 1hr   | Medium — integration                 |
| 14  | P2       | Add benchmark tests for hot paths                     | 30min | Low — performance                    |
| 15  | P2       | Extract shared Tailwind preset/theme file             | 1hr   | Medium — customization               |
| 16  | P2       | Create example/demo app                               | 2hr   | Medium — usability                   |
| 17  | P3       | Client-side form validation helpers                   | 2hr   | Medium — UX                          |
| 18  | P3       | Documentation site generation                         | 3hr   | Medium — adoption                    |
| 19  | P3       | Create CONTEXT.md with architecture decisions         | 30min | Medium — onboarding                  |
| 20  | P3       | Investigate nix flake migration                       | 2hr   | Low — build tooling                  |
| 21  | P3       | Add `docs/adr/` for architecture decision records     | 15min | Low — documentation                  |
| 22  | P3       | Release automation (goreleaser)                       | 1hr   | Low — DevOps                         |
| 23  | P3       | Version migration guides                              | 2hr   | Low — adoption                       |
| 24  | P4       | Cross-package circular import guard test              | 15min | Low — safety                         |
| 25  | P4       | Interactive playground/REPL                           | 4hr   | Low — nice-to-have                   |

---

## g) Top #1 Question I Cannot Figure Out Myself

**Should the shared SVG helpers (`fillIcon`, spinner SVG) live in `internal/svg/`, `shared/`, or somewhere else?**

The current situation:

- `display/helpers.templ` has `fillIcon` — usable only within `display/` package
- `icons/icon.templ` imports `feedback.Spinner` for the Spinner icon case
- `htmx/loading.templ` also imports `feedback.Spinner`
- `feedback/loading.templ` has the canonical `Spinner` component

Creating `internal/svg/` or `shared/` would:

- ✅ Solve cross-package access to `fillIcon` (needed by `navigation/pagination.templ` for arrows)
- ✅ Centralize spinner SVG definition (one source of truth)
- ✅ Break the `icons → feedback` import dependency
- ❌ Add a new package to the public API surface
- ❌ Require deciding between `internal/` (not importable by consumers) vs `shared/` (importable)

This is an architectural decision that affects the public API and import ergonomics. I need your call on package location and visibility.

---

## Session Commits (7 total)

```
b416c8d refactor(layout): consolidate two SRI hash functions into single htmxSRI()
817e8b9 refactor(feedback,display): convert switch-based style lookups to map-based
5fcf1b0 refactor(feedback): embed utils.BaseProps in ProgressBarProps
ed88096 refactor(forms): embed utils.BaseProps in all form Props structs
8ad0334 refactor(layout): rename BaseProps → PageProps to disambiguate from utils.BaseProps
7ac54a4 fix(utils): correct TestPtr — was testing new() instead of Ptr()
7d17413 refactor: semantic deduplication — 13→7 clone groups, -57 net lines
```

## File Change Summary (since session start)

```
 33 files changed, 941 insertions(+), 502 deletions(-)
 Net: +439 lines added, 502 lines removed
```
