# Status Report — templ-components

**Date:** 2026-05-03 07:04  
**Branch:** master  
**Commit:** eed8aa0 (latest) → pending dedup commit  
**Go:** 1.26.2 | **templ:** v0.3.1001 | **tailwind-merge-go:** v0.2.1

---

## Metrics at a Glance

| Metric                  | Value                         |
| ----------------------- | ----------------------------- |
| Packages                | 8                             |
| Tests                   | 76 passing, 0 failing         |
| `.templ` files          | 30                            |
| Total `.templ` lines    | 2,958                         |
| `go vet`                | clean                         |
| Clone groups (art-dupl) | 7 (down from 13)              |
| Typed enums             | 16+                           |
| CSP violations          | 0                             |
| Dependencies            | 2 (templ + tailwind-merge-go) |

---

## a) FULLY DONE

### Deduplication Pass (this session)

Reduced semantic clone groups from **13 → 7** (46% reduction), net -74 lines of code:

| File                          | Change                                                                               |
| ----------------------------- | ------------------------------------------------------------------------------------ |
| `htmx/loading.templ`          | 3 inline spinner SVGs → `@feedback.Spinner()`                                        |
| `navigation/pagination.templ` | Extracted `paginationArrowIcon`, `paginationArrow`, `mobilePageButton` sub-templates |
| `feedback/alert.templ`        | Extracted `inlineMessage`, `inlineErrorIcon`, `inlineSuccessIcon` sub-templates      |
| `display/tabs.templ`          | Extracted `tabLink` sub-template                                                     |
| `display/dropdown.templ`      | Extracted `dropdownItemClass` constant + `dropdownItemLink` sub-template             |
| `display/empty_state.templ`   | Extracted `emptyStateActionClass` constant + `emptyStateAction` sub-template         |
| `display/helpers.templ`       | NEW — shared `fillIcon` sub-template for 20×20 filled SVG icons                      |
| `display/accordion.templ`     | Inline SVG → `@fillIcon()`                                                           |
| `display/card.templ`          | 2 inline SVGs → `@fillIcon()`                                                        |
| `icons/icon.templ`            | Spinner case → `@feedback.Spinner()`                                                 |

### Test Cleanup (this session, incidental)

| File                        | Change                                                                         |
| --------------------------- | ------------------------------------------------------------------------------ |
| `display/card_test.go`      | Removed zero-value struct fields in test constructors                          |
| `display/modal_test.go`     | Removed zero-value struct fields in test constructors                          |
| `feedback/helpers_test.go`  | Extracted `assertStyleFunc4` helper, deduplicated toast/alert style assertions |
| `feedback/snapshot_test.go` | Removed zero-value struct fields in test constructors                          |

### Previously Completed (from git history)

- 16+ typed string enums (AlertType, ToastType, TabsStyle, BadgeSize, AvatarSize, etc.)
- CSP compliance for all inline scripts (nonce-based)
- 76 render + unit tests across all 8 packages
- Comprehensive project hygiene (git config, linting, formatting, metadata)
- `utils.Class()` with tailwind-merge-go for intelligent class merging
- `utils.BaseProps` shared across all components
- Icon system with 30+ named icons + typed `Name` enum

---

## b) PARTIALLY DONE

### Deduplication — 7 Remaining Clone Groups

These are **structural similarities** between fundamentally different elements — not safely deduplicable without over-engineering:

| #   | Clone                                                                                            | Why Not Deduped                                                                                                        |
| --- | ------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------- |
| 1   | `modal.templ:36,39` + `theme.templ:28,31` + `theme.templ:31,34`                                  | SVG close-icon in modal vs sun/moon icons in theme toggle — different packages, different semantic meaning             |
| 2   | `breadcrumbs.templ:27,31` + `pagination.templ:142,156`                                           | Active link vs inactive link in different navigation patterns — extracting would create weird cross-package dependency |
| 3   | `forms/label.templ:23,29` + `pagination.templ:59,70`                                             | Conditional `<p>` with class vs conditional `<a>` with class — different HTML elements                                 |
| 4   | `progress.templ:48,56` + `htmx/helpers.templ:7,17`                                               | Progress bar inner div vs ConfirmDelete button — superficial attribute similarity only                                 |
| 5   | `nav_link.templ:37,47` + `nav_link.templ:55,67`                                                  | Desktop NavLink vs MobileNavLink — different element structures (only class similarity)                                |
| 6   | `alert.templ:152,153` + `toast.templ:190,191` + `loading.templ:22,29` + `pagination.templ:91,94` | Path element attribute pattern across 4 files — too thin to extract                                                    |
| 7   | `empty_state.templ:25,32` + `empty_state.templ:33,40`                                            | `<a>` vs `<button>` inside `emptyStateAction` — already the best we can do without dynamic element tags                |

### Test Coverage

- Render tests exist for all components but no **snapshot/golden** tests
- No integration/E2E tests
- No accessibility (a11y) automated checks

---

## c) NOT STARTED

1. **FEATURES.md** — No feature inventory exists
2. **TODO_LIST.md** — No comprehensive TODO tracking
3. **CONTEXT.md** — No project context document
4. **Snapshot/golden file testing** — No rendered HTML snapshots
5. **A11y testing** — No automated accessibility audits
6. **Benchmark tests** — No performance benchmarks
7. **Example app** — No demo/showcase application
8. **CSS extraction** — No shared Tailwind preset/theme file
9. **Cross-package shared helpers** — `fillIcon` exists in display/ but other packages can't use it without import cycles
10. **Documentation site** — No generated docs/website
11. **Migration guide** — No version migration docs
12. **Playground/REPL** — No interactive component preview
13. **CI pipeline** — No `.github/workflows` or equivalent
14. **Release automation** — No goreleaser or similar
15. **Changelog** — No CHANGELOG.md

---

## d) TOTALLY FUCKED UP

1. **`utils/utils_test.go:48`** — LSP warning: `impossible condition: non-nil == nil` — dead test code that needs fixing
2. **No FEATURES.md / TODO_LIST.md** — Previous status reports reference these but they don't exist on disk (deleted or never created)
3. **8 stale status reports** in `docs/status/` — no cleanup, no pruning, growing unbounded
4. **`empty_state.go`** hand-written Go file — potential templ generation conflict
5. **Cross-package import risk** — `icons/icon.templ` now imports `feedback` while `feedback/loading.templ` uses `icons` via toast — circular import hazard lurking

---

## e) WHAT WE SHOULD IMPROVE

1. **Fix the dead test code** in `utils/utils_test.go:48`
2. **Create FEATURES.md** — audit actual code, not aspirational
3. **Create TODO_LIST.md** — track real work items
4. **Extract shared SVG helpers to a `shared/` or `internal/svg/` package** — `fillIcon` and spinner are needed cross-package; currently `display/fillIcon` can't be used by `navigation/` or `feedback/`
5. **Add snapshot testing** — golden files for rendered HTML
6. **Add a11y linting** — at minimum `aria-*` attribute validation
7. **Set up CI** — build + test + vet + lint on push
8. **Prune old status reports** — keep last 2, archive rest
9. **Consistent nonce handling** — some components use `props.Nonce`, some hardcode, some skip
10. **Form validation patterns** — no client-side validation helpers exist
11. **Dark mode testing** — no tests verify `dark:` class output
12. **Component composition tests** — no tests for nesting components inside each other

---

## f) Top 25 Things We Should Get Done Next

| #   | Priority | Task                                                                                    | Impact |
| --- | -------- | --------------------------------------------------------------------------------------- | ------ |
| 1   | P0       | Fix `utils/utils_test.go:48` dead code warning                                          | 5min   |
| 2   | P0       | Create FEATURES.md from actual codebase audit                                           | 30min  |
| 3   | P0       | Create TODO_LIST.md                                                                     | 20min  |
| 4   | P0       | Extract `fillIcon` + `spinnerSVG` to shared package (avoid cross-package import issues) | 1hr    |
| 5   | P1       | Add snapshot/golden file tests for all 30 components                                    | 2hr    |
| 6   | P1       | Set up CI pipeline (build + test + vet)                                                 | 1hr    |
| 7   | P1       | Add a11y attribute validation tests                                                     | 1hr    |
| 8   | P1       | Consistent nonce propagation audit                                                      | 30min  |
| 9   | P1       | Create CHANGELOG.md                                                                     | 20min  |
| 10  | P2       | Dark mode output verification tests                                                     | 1hr    |
| 11  | P2       | Component composition tests (nesting)                                                   | 1hr    |
| 12  | P2       | Prune old status reports (keep last 2)                                                  | 5min   |
| 13  | P2       | Add benchmark tests for hot paths (Class, spinner render)                               | 30min  |
| 14  | P2       | Extract shared Tailwind preset/theme file                                               | 1hr    |
| 15  | P2       | Create example/demo app                                                                 | 2hr    |
| 16  | P3       | Client-side form validation helpers                                                     | 2hr    |
| 17  | P3       | Documentation site generation                                                           | 3hr    |
| 18  | P3       | Interactive playground/REPL                                                             | 4hr    |
| 19  | P3       | Release automation (goreleaser)                                                         | 1hr    |
| 20  | P3       | Version migration guides                                                                | 2hr    |
| 21  | P3       | Add `CONTEXT.md` with architecture decisions                                            | 30min  |
| 22  | P3       | Investigate `empty_state.go` hand-written file for conflicts                            | 15min  |
| 23  | P4       | Cross-package circular import audit (icons ↔ feedback)                                  | 30min  |
| 24  | P4       | Add `docs/adr/` directory for architecture decisions                                    | 15min  |
| 25  | P4       | Explore nix flake migration (replace justfile/shell scripts)                            | 2hr    |

---

## g) Top #1 Question I Cannot Figure Out Myself

**Should `icons` and `feedback` have a cross-package dependency?**

Currently `icons/icon.templ` imports `feedback.Spinner` for the Spinner icon case. This creates `icons → feedback` dependency. Meanwhile `feedback/toast.templ` and `feedback/loading.templ` conceptually want icon support. This is a **potential circular import** if feedback ever needs icons directly.

The proper fix is extracting shared SVG primitives (spinner, fillIcon) into a third package (e.g., `internal/svg` or `shared/`). But that's an architectural decision about package layout that only the project owner should make — it affects the public API surface and import ergonomics.

**Question:** Do you want a `shared/` package for cross-cutting SVG/template helpers, or should we keep the current structure and accept the one-directional `icons → feedback` dependency?

---

## File Change Summary (Uncommitted)

```
 display/accordion.templ     |  13 ++---
 display/card.templ          |   8 +--
 display/card_test.go        |  21 +++-----
 display/dropdown.templ      |  47 +++++++++---------
 display/empty_state.templ   |  42 +++++++++-------
 display/helpers.templ       |  NEW (shared fillIcon)
 display/modal_test.go       |  17 ++-----
 display/tabs.templ          |  65 ++++++++++---------------
 feedback/alert.templ        |  34 ++++++++-----
 feedback/helpers_test.go    |  64 +++++++++---------------
 feedback/snapshot_test.go   |  26 +---------
 htmx/loading.templ          |  17 ++-----
 icons/icon.templ            |   7 ++-
 navigation/pagination.templ | 115 ++++++++++++++++++++++----------------------

 13 files changed, 201 insertions(+), 275 deletions(-)
 Net: -74 lines
```
