<!-- AUTO-UPDATED 2026-07-10: Retrospective status overlay -->

> ## 🔔 Update Notice — 2026-07-10
>
> This report is **historical**. Many items listed as "open", "todo", or "broken" below
> have since been **fixed and verified**. Do not act on open items without first checking
> [TODO_LIST.md](../../TODO_LIST.md) for current status.
>
> **Key fixes completed since this report:**
>
> - ✅ All 7 P0 bugs fixed (InlineLoadingOverlay a11y, SanitizeID mismatch, FromError fallback,
>   Footer BaseProps, ErrorPage/NotFound404 `<main>` landmark, CSRFTokenName, grid-rows verified)
> - ✅ `encoding/json/v2` purged from all production code + pre-commit guard added
> - ✅ Motion constants centralized in `utils/motion.go`, wired into 13 components
> - ✅ `FamilyFromErrorFamily` → `FromErrorFamily` (old name kept as deprecated alias)
> - ✅ `icons.IconRTL()` + CSS for directional icon RTL mirroring
> - ✅ 33 regression tests added (htmx, errorpage, layout, navigation, feedback, display)
> - ✅ Dark golden test infrastructure (badge/card/button)
> - ✅ CHANGELOG consolidated, ROADMAP updated, migration guide created
> - ✅ All 14 packages pass, 0 lint issues
>
> **Canonical source of truth:** [TODO_LIST.md](../../TODO_LIST.md) (52 items, 37 ✅ done, 12 deferred/blocked)

---

# Dark Mode Compliance Plan — 2026-07-08

**Created:** 2026-07-08 17:15 CEST
**Source:** `docs/status/2026-07-08_17-02_dark-mode-audit.md` — 50 TODOs from dark mode audit session
**Constraints:** Every task ≤ 12 min. Sorted by importance → impact → effort → customer value.
**Decisions incorporated:**

- Dark mode compliance test → **FAILING** (blocks CI)
- `encoding/json/v2` in `handler.go` → **intentional** (migrating to json/v2 everywhere)

---

## Priority Legend

| Priority | Meaning                                                        |
| -------- | -------------------------------------------------------------- |
| **P0**   | Blocks release or prevents regressions — must do before commit |
| **P1**   | High value, low effort — quick wins                            |
| **P2**   | Medium value — good to do, not urgent                          |
| **P3**   | Low priority / future exploration                              |

---

## Master Table

| ID  | Task                                                                                         | Priority | Impact   | Effort | Subtasks | Deps | Category |
| --- | -------------------------------------------------------------------------------------------- | -------- | -------- | ------ | -------- | ---- | -------- |
| 1   | Separate `handler.go` json/v2 change into its own commit                                     | P0       | High     | 5m     | 1        | —    | Cleanup  |
| 2   | Commit dark mode fixes (excluding handler.go)                                                | P0       | High     | 10m    | 2        | 1    | Release  |
| 3   | Add `TestDarkModeCompliance` — scan `.templ` for neutral colors without `dark:`              | P0       | Critical | 12m    | 3        | 2    | Testing  |
| 4   | Add `TestDarkModeSemanticColors` — scan for `bg-blue-600` etc. without `dark:`               | P0       | Critical | 12m    | 3        | 2    | Testing  |
| 5   | Run new tests — verify they PASS (all issues already fixed)                                  | P0       | Critical | 5m     | 1        | 3,4  | Testing  |
| 6   | Add `color-scheme: dark` to `.dark` class in theme CSS                                       | P0       | High     | 8m     | 2        | —    | A11y     |
| 7   | Bump `utils.Version` + CHANGELOG `[Unreleased]` entry                                        | P0       | Medium   | 10m    | 3        | 2    | Release  |
| 8   | Update FEATURES.md to note full dark mode compliance                                         | P1       | Medium   | 5m     | 1        | 2    | Docs     |
| 9   | Document dark mode convention in AGENTS.md                                                   | P1       | High     | 12m    | 3        | —    | Docs     |
| 10  | Fix `progressbar.templ:54` LSP hint (use `max()`)                                            | P1       | Low      | 8m     | 2        | —    | Code     |
| 11  | Fix doc comments in `htmx/loading.templ` (add `dark:` variants)                              | P1       | Low      | 4m     | 1        | —    | Docs     |
| 12  | Fix doc comments in `feedback/loading.templ` (add `dark:` variants)                          | P1       | Low      | 4m     | 1        | —    | Docs     |
| 13  | Fix doc comments in `icons/icon.templ` (add `dark:` variants)                                | P1       | Low      | 4m     | 1        | —    | Docs     |
| 14  | Fix doc comments in `forms/input_group.templ` (add `dark:` variants)                         | P1       | Low      | 4m     | 1        | —    | Docs     |
| 15  | Update SKILL.md Part 2 with dark mode checklist for new components                           | P1       | High     | 10m    | 1        | 9    | Docs     |
| 16  | Add ADR `0011-dark-mode-convention.md`                                                       | P1       | Medium   | 10m    | 1        | 9    | Docs     |
| 17  | Update README.md dark mode section with complete convention                                  | P1       | Medium   | 10m    | 1        | 9    | Docs     |
| 18  | Tag patch release (v0.9.1 or similar)                                                        | P1       | High     | 10m    | 2        | 7,8  | Release  |
| 19  | Add dark mode integration test — render with `.dark` wrapper, assert `dark:` classes present | P2       | Medium   | 12m    | 2        | 3,4  | Testing  |
| 20  | Audit `hover:` variants for missing `dark:hover:`                                            | P2       | Medium   | 12m    | 2        | —    | Audit    |
| 21  | Audit `focus:` variants for missing `dark:focus:`                                            | P2       | Medium   | 12m    | 2        | —    | Audit    |
| 22  | Audit `ring-offset-*` for dark mode                                                          | P2       | Low      | 8m     | 2        | —    | Audit    |
| 23  | Audit `shadow-*` classes for dark mode                                                       | P2       | Low      | 6m     | 2        | —    | Audit    |
| 24  | Audit `backdrop-blur-*` opacity for dark mode                                                | P2       | Low      | 6m     | 1        | —    | Audit    |
| 25  | Add test for toast JS-created toast (dynamic path)                                           | P2       | Medium   | 12m    | 2        | —    | Testing  |
| 26  | Add dark golden test variants (render with `.dark` parent)                                   | P2       | Medium   | 12m    | 2        | 19   | Testing  |
| 27  | Verify WCAG AA contrast ratios for all dark mode color combos                                | P2       | High     | 12m    | 3        | —    | A11y     |
| 28  | Add `scrollbar-color` for dark mode                                                          | P2       | Low      | 6m     | 2        | —    | CSS      |
| 29  | Verify `::selection` colors work in dark mode                                                | P2       | Low      | 4m     | 1        | —    | Audit    |
| 30  | Verify Table `hover:bg-gray-50` has `dark:hover:bg-gray-800`                                 | P2       | Low      | 6m     | 2        | —    | Audit    |
| 31  | Add `darkMode` toggle to demo page                                                           | P2       | Low      | 10m    | 1        | —    | Demo     |
| 32  | Add dark mode release note                                                                   | P2       | Low      | 5m     | 1        | 18   | Docs     |
| 33  | Update CONTRIBUTING.md with dark mode section                                                | P2       | Low      | 10m    | 1        | 9    | Docs     |
| 34  | Add pre-commit hook for dark mode check                                                      | P2       | Medium   | 10m    | 1        | 3,4  | Tooling  |
| 35  | Add test that `BaseProps.Class` propagates `dark:` classes                                   | P2       | Low      | 10m    | 1        | —    | Testing  |
| 36  | Benchmark dark mode class resolution (tailwind-merge-go with longer strings)                 | P2       | Low      | 12m    | 2        | —    | Perf     |
| 37  | Verify CountBadge `ring-white dark:ring-gray-800` is sufficient                              | P3       | Low      | 5m     | 1        | —    | Audit    |
| 38  | Verify Tooltip arrow border colors in dark mode                                              | P3       | Low      | 5m     | 1        | —    | Audit    |
| 39  | Verify Nav mobile menu slide animation visible in dark mode                                  | P3       | Low      | 5m     | 1        | —    | Audit    |
| 40  | Add `prefers-color-scheme` fallback when `.dark` class is absent                             | P3       | Medium   | 10m    | 1        | —    | CSS      |
| 41  | Add `prefers-reduced-transparency` media query for overlays                                  | P3       | Low      | 6m     | 1        | —    | A11y     |
| 42  | Add contract test: props with `Color` field have dark mode godoc                             | P3       | Low      | 10m    | 1        | —    | Testing  |
| 43  | Add test asserting every component has `dark:` classes                                       | P3       | Medium   | 10m    | 1        | 19   | Testing  |
| 44  | Explore Tailwind v4 `@theme` dark mode tokens (CSS-first approach)                           | P3       | Medium   | 12m    | 2        | —    | Research |
| 45  | Add `darkMode()` helper in `utils` (returns `dark:` prefixed classes)                        | P3       | Low      | 10m    | 1        | —    | Code     |
| 46  | Add `Theme` enum (Light/Dark/Auto) to `layout`                                               | P3       | Low      | 12m    | 2        | —    | Code     |
| 47  | Add `SidebarNav` light mode option (prop to switch from dark sidebar)                        | P3       | Low      | 12m    | 2        | —    | Code     |
| 48  | Add visual regression testing (screenshot comparison light/dark)                             | P3       | Medium   | 12m    | 2        | —    | Research |
| 49  | Add screen reader behavior verification (dark: changes are no-op for SR)                     | P3       | Low      | 5m     | 1        | —    | A11y     |
| 50  | Add `color-scheme: light` to `:root` (explicit light mode)                                   | P3       | Low      | 4m     | 1        | 6    | CSS      |
| 51  | Add test for `prefers-reduced-transparency` support                                          | P3       | Low      | 5m     | 1        | 41   | Testing  |
| 52  | Explore `scroll-smooth` in dark mode (verify no issues)                                      | P3       | Low      | 3m     | 1        | —    | Audit    |

---

## Detailed Task Breakdown

### P0 — Critical (blocks release / prevents regressions)

#### Task 1: Separate `handler.go` json/v2 change into its own commit

- **1a** (5m): `git stash` the dark mode changes, commit `handler.go` alone as `refactor: migrate errorpage/handler.go to encoding/json/v2`

#### Task 2: Commit dark mode fixes (excluding handler.go)

- **2a** (5m): Stage all dark mode `.templ`, `.go`, golden, and test files (exclude `handler.go`)
- **2b** (5m): Write commit message: `fix: comprehensive dark mode audit — 30+ missing dark: variants fixed across all packages`

#### Task 3: Add `TestDarkModeCompliance` — neutral colors

- **3a** (5m): Design test in `utils/darkmode_test.go` — read all `.templ` files, regex scan for `text-gray-*`, `bg-white`, `bg-gray-*`, `border-gray-*`, `divide-gray-*`, `ring-gray-*` without `dark:` on same line
- **3b** (5m): Implement test with allowed-exceptions list (Toggle thumb `bg-white`, SidebarNav permanently-dark sidebar)
- **3c** (2m): Run test, verify it PASSES (all issues already fixed in this session)

#### Task 4: Add `TestDarkModeSemanticColors` — semantic colors

- **4a** (5m): Design test — scan for `bg-blue-600`, `bg-red-600`, `bg-green-600`, `bg-amber-600`, `bg-orange-600`, `bg-gray-600`, `text-blue-600`, `text-red-600`, `text-green-600`, `text-amber-500`, `text-orange-500`, `text-blue-500`, `text-red-500`, `text-green-500` without `dark:` on same line
- **4b** (5m): Implement test with exceptions (doc comments in `//` lines)
- **4c** (2m): Run test, verify it PASSES

#### Task 5: Run new tests — verify they PASS

- **5a** (5m): `go test ./utils/... -run TestDarkMode` — confirm zero failures

#### Task 6: Add `color-scheme: dark` to `.dark` class

- **6a** (5m): Add `color-scheme: dark;` to `.dark { }` block in `templ-components-theme.css`
- **6b** (3m): Add `color-scheme: light;` to `:root` or base body styles

#### Task 7: Bump `utils.Version` + CHANGELOG

- **7a** (3m): Add `[Unreleased]` CHANGELOG entry: "Dark mode audit: 30+ missing `dark:` variants fixed across all packages. Added `TestDarkModeCompliance` and `TestDarkModeSemanticColors` regression tests."
- **7b** (3m): Bump `utils.Version` if cutting release (e.g., `0.9.1`)
- **7c** (4m): Update `FEATURES.md` `**Version:**` to match

---

### P1 — High Value, Low Effort (quick wins)

#### Task 8: Update FEATURES.md for dark mode compliance

- **8a** (5m): Add/update FEATURES.md entry: "Full dark mode compliance — all components have `dark:` variants for every neutral and semantic color class"

#### Task 9: Document dark mode convention in AGENTS.md

- **9a** (4m): Add section: "Dark mode color convention: `-600` light → `-500` dark for backgrounds, `-400` for text. Neutral colors: `gray-*` exclusively."
- **9b** (4m): Document exceptions: SidebarNav (permanently dark sidebar, `hover:bg-gray-800` intentional), Toggle thumb (`bg-white` in both modes, track changes color instead)
- **9c** (4m): Document the compliance test: "Every `.templ` file is scanned by `TestDarkModeCompliance` — neutral colors without `dark:` variants fail CI"

#### Task 10: Fix `progressbar.templ:54` LSP hint

- **10a** (3m): Replace `if` statement with `max()` builtin (Go 1.21+)
- **10b** (5m): `templ generate ./... && go test ./feedback/...`

#### Task 11: Fix doc comments in `htmx/loading.templ`

- **11a** (4m): Update 2 doc comments: add `dark:text-blue-400` to spinner Color examples

#### Task 12: Fix doc comments in `feedback/loading.templ`

- **12a** (4m): Update 1 doc comment: add `dark:text-blue-400` to spinner Color example

#### Task 13: Fix doc comments in `icons/icon.templ`

- **13a** (4m): Update 2 doc comments: add `dark:text-gray-400` to icon class examples

#### Task 14: Fix doc comments in `forms/input_group.templ`

- **14a** (4m): Update 1 doc comment: add `dark:text-gray-500` to LeftAddon icon example

#### Task 15: Update SKILL.md Part 2 with dark mode checklist

- **15a** (10m): Add "Dark mode checklist" section to SKILL.md Part 2: "Every neutral color (`text-gray-*`, `bg-white`, `bg-gray-*`, `border-gray-*`) must have `dark:` variant. Every semantic color (`bg-blue-600`, `text-red-600`) must have `dark:` variant. Run `TestDarkModeCompliance` before committing."

#### Task 16: Add ADR `0011-dark-mode-convention.md`

- **16a** (10m): Write ADR documenting: class-based dark mode via `@custom-variant dark`, gray-only palette, `-600`→`-500` shade convention, compliance test enforcement

#### Task 17: Update README.md dark mode section

- **17a** (10m): Update README.md dark mode section with: complete color convention, compliance test mention, exceptions list

#### Task 18: Tag patch release

- **18a** (5m): Run `scripts/release.sh 0.9.1 "dark mode audit — 30+ fixes, compliance tests"` (or appropriate version)
- **18b** (5m): Review release commit and tag with `git show`

---

### P2 — Medium Value (good to do, not urgent)

#### Task 19: Add dark mode integration test

- **19a** (6m): Create `integration/dark_mode_test.go` — render every component, assert output contains `dark:` classes
- **19b** (6m): Add test cases for each package (display, forms, feedback, navigation, layout, htmx, errorpage)

#### Task 20: Audit `hover:` variants for missing `dark:hover:`

- **20a** (6m): Grep for `hover:bg-` and `hover:text-` in `.templ` files without `dark:hover:` on same line
- **20b** (6m): Fix any issues found

#### Task 21: Audit `focus:` variants for missing `dark:focus:`

- **21a** (6m): Grep for `focus:ring-` and `focus:border-` without `dark:focus:` on same line
- **21b** (6m): Fix any issues found

#### Task 22: Audit `ring-offset-*` for dark mode

- **22a** (4m): Grep for `ring-offset-` without `dark:ring-offset-` on same line
- **22b** (4m): Fix if needed (likely `ring-offset-white` → `dark:ring-offset-gray-900`)

#### Task 23: Audit `shadow-*` classes for dark mode

- **23a** (3m): Grep for `shadow-*` in `.templ` files — assess if dark mode shadows needed
- **23b** (3m): Add `dark:shadow-*` if visibility is an issue

#### Task 24: Audit `backdrop-blur-*` opacity for dark mode

- **24a** (6m): Check `bg-white/80 dark:bg-gray-950/80` etc. — verify opacity is sufficient in dark mode

#### Task 25: Add test for toast JS-created toast — **✅ DONE (2026-07-10)**

- **25a** (6m): ~~Write test that simulates `tcShowToast()` JS call and verifies the constructed HTML has `dark:` classes~~ — Done: `feedback/toast_regression_test.go`
- **25b** (6m): ~~Add golden file for JS-created toast~~ — Done: 6 toast regression tests added

#### Task 26: Add dark golden test variants — **✅ DONE (2026-07-10)**

- **26a** (6m): ~~Create test helper that renders component inside `<div class="dark">` wrapper~~ — Done: `display/dark_golden_test.go`
- **26b** (6m): ~~Generate dark-mode golden files for key components~~ — Done: badge_dark, card_dark, button_dark golden files created

#### Task 27: Verify WCAG AA contrast ratios

- **27a** (4m): List all dark mode color combinations (text on bg): `gray-100/gray-900`, `gray-400/gray-900`, `blue-400/gray-800`, `red-400/gray-900`, etc.
- **27b** (4m): Calculate contrast ratios using WCAG formula (need ≥ 4.5:1 for normal text, ≥ 3:1 for large text)
- **27c** (4m): Document results, flag any failures

#### Task 28: Add `scrollbar-color` for dark mode

- **28a** (3m): Add `scrollbar-color: gray-600 gray-900;` to `.dark` in theme CSS
- **28b** (3m): Test in browser (Firefox supports, WebKit needs `-webkit-scrollbar`)

#### Task 29: Verify `::selection` colors in dark mode

- **29a** (4m): Check `selection:bg-blue-100 dark:selection:bg-blue-900 selection:text-blue-900 dark:selection:text-blue-100` in `base.templ` BodyClass — already has dark: variants, verify correctness

#### Task 30: Verify Table `hover:bg-gray-50` has `dark:hover:bg-gray-800`

- **30a** (3m): Check `table.templ` for hover states
- **30b** (3m): Fix if missing

#### Task 31: Add `darkMode` toggle to demo page

- **31a** (10m): Add a ThemeToggle to the demo page so you can preview components in dark mode

#### Task 32: Add dark mode release note

- **32a** (5m): Write release note highlighting: "Comprehensive dark mode audit — all 83 components now have proper `dark:` variants. Added regression tests to prevent future gaps."

#### Task 33: Update CONTRIBUTING.md with dark mode section

- **33a** (10m): Add "Dark mode" section to CONTRIBUTING.md: convention, compliance test, checklist for new components

#### Task 34: Add pre-commit hook for dark mode check

- **34a** (10m): Add `TestDarkModeCompliance` + `TestDarkModeSemanticColors` to `scripts/pre-commit.sh`

#### Task 35: Add test that `BaseProps.Class` propagates `dark:` classes

- **35a** (10m): Write test that passes `Class: "dark:bg-red-500"` to a component and asserts it appears in output

#### Task 36: Benchmark dark mode class resolution

- **36a** (6m): Write benchmark: `utils.Class("bg-blue-600 dark:bg-blue-500 hover:bg-blue-500 dark:hover:bg-blue-400 text-white ...")` vs shorter strings
- **36b** (6m): Compare with existing benchmarks, document if any perf regression

---

### P3 — Lower Priority / Future Exploration

#### Task 37: Verify CountBadge ring colors

- **37a** (5m): Check `ring-white dark:ring-gray-800` on CountBadge — the ring separates the badge from the parent element

#### Task 38: Verify Tooltip arrow border colors

- **38a** (5m): Check `border-t-gray-900 dark:border-t-gray-700` etc. — already has dark: variants, verify visibility

#### Task 39: Verify Nav mobile menu slide animation

- **39a** (5m): Check that the mobile menu slide transition is visible against dark backgrounds

#### Task 40: Add `prefers-color-scheme` fallback

- **40a** (10m): Add `@media (prefers-color-scheme: dark) { :root:not(.light) { ... } }` to theme CSS for users without JS

#### Task 41: Add `prefers-reduced-transparency` for overlays

- **41a** (6m): Add `@media (prefers-reduced-transparency) { .dark { background: solid; } }` for Modal/Drawer overlays

#### Task 42: Add contract test for `Color` field godoc

- **42a** (10m): Write test that checks props structs with a `Color string` field have godoc mentioning dark mode

#### Task 43: Add test asserting every component has `dark:` classes

- **43a** (10m): Exhaustive test: render every component, parse output, assert at least one `dark:` class exists

#### Task 44: Explore Tailwind v4 `@theme` dark tokens

- **44a** (6m): Research: can `@theme` define dark-mode-specific tokens that auto-apply?
- **44b** (6m): Prototype: replace per-class `dark:` with CSS variable approach

#### Task 45: Add `darkMode()` helper in `utils`

- **45a** (10m): Design `func DarkMode(lightClass, darkClass string) string` that returns `lightClass + " dark:" + darkClass`

#### Task 46: Add `Theme` enum to `layout`

- **46a** (6m): Design `type Theme string` with `ThemeLight`, `ThemeDark`, `ThemeAuto`
- **46b** (6m): Implement context propagation so components can detect theme

#### Task 47: Add `SidebarNav` light mode option

- **47a** (6m): Add `SidebarVariant` prop: `SidebarDark` (default) vs `SidebarLight`
- **47b** (6m): Implement light sidebar classes

#### Task 48: Add visual regression testing

- **48a** (6m): Research tools: Playwright, Percy, or simple screenshot diffing for Go/templ
- **48b** (6m): Prototype: render demo page in light + dark, screenshot, diff

#### Task 49: Add screen reader behavior verification

- **49a** (5m): Verify that `dark:` CSS changes don't affect screen reader output (they shouldn't — SR ignores visual CSS)

#### Task 50: Add `color-scheme: light` to `:root`

- **50a** (4m): Explicitly set `color-scheme: light` on `:root` so native form controls render correctly in light mode

#### Task 51: Add test for `prefers-reduced-transparency` support

- **51a** (5m): Test that overlays respect `prefers-reduced-transparency` media query

#### Task 52: Verify `scroll-smooth` in dark mode

- **52a** (3m): Verify `scroll-smooth` on body works correctly in dark mode (no visual issues)

---

## Summary by Priority

| Priority  | Tasks  | Total Effort | Key Outcome                                      |
| --------- | ------ | ------------ | ------------------------------------------------ |
| **P0**    | 7      | ~52m         | Regression tests, clean commit, release-ready    |
| **P1**    | 11     | ~72m         | Documentation, doc comments, ADR, patch release  |
| **P2**    | 18     | ~150m        | Comprehensive testing, audits, WCAG verification |
| **P3**    | 16     | ~115m        | Future exploration, advanced features            |
| **Total** | **52** | **~389m**    | **Full dark mode compliance + prevention**       |

## Summary by Category

| Category | Tasks | Total Effort |
| -------- | ----- | ------------ |
| Testing  | 14    | ~115m        |
| Docs     | 10    | ~72m         |
| Audit    | 10    | ~65m         |
| Release  | 4     | ~30m         |
| Code     | 4     | ~35m         |
| A11y     | 5     | ~35m         |
| CSS      | 4     | ~28m         |
| Perf     | 1     | ~12m         |
| Research | 2     | ~12m         |
| Tooling  | 1     | ~10m         |
| Demo     | 1     | ~10m         |
| Cleanup  | 1     | ~5m          |
