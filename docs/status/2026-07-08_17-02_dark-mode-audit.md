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

# Status Report: Dark Mode Audit — 2026-07-08

**Session:** Single-session comprehensive dark mode audit and fix
**Date:** 2026-07-08 17:02 CEST
**Scope:** All packages — display, forms, feedback, navigation, layout, htmx, errorpage, icons, utils, internal/svg, examples/demo

---

## What Was Done

A systematic audit of light AND dark mode support across every `.templ` and `.go` source file
(excluding generated `*_templ.go` and test files). Found and fixed **30+ dark mode gaps** where
components had light-mode color classes without corresponding `dark:` variants.

### Audit Methodology

1. **Palette consistency scan** — searched for `slate-*`, `zinc-*`, `neutral-*`, `stone-*` to verify
   the "gray-only" convention. Result: **clean** — no palette mixing found.
2. **Missing `dark:` variant scan** — searched for neutral color classes (`text-gray-*`, `bg-white`,
   `bg-gray-*`, `border-gray-*`, `ring-gray-*`, `divide-gray-*`, `placeholder:*`) without `dark:`
   counterparts on the same line. Used regex grep + sub-agent verification.
3. **Semantic color scan** — searched for `bg-blue-600`, `bg-red-600`, `text-blue-600`, `text-red-600`,
   `text-green-600`, `text-amber-500`, `text-orange-500`, etc. without `dark:` variants.
4. **Focus ring scan** — searched for `focus-visible:ring-gray-*` without `dark:focus-visible:ring-*`.
5. **Final verification scan** — re-ran all searches after fixes to confirm zero remaining gaps.

### Issues Found and Fixed (by category)

#### Critical: Invisible/Unreadable Elements in Dark Mode

| Component                     | File                                | Issue                                            | Fix                        |
| ----------------------------- | ----------------------------------- | ------------------------------------------------ | -------------------------- |
| Modal/Drawer close button     | `display/shared.templ:47`           | `text-gray-400` no `dark:text-*`                 | Added `dark:text-gray-500` |
| Dropdown chevron              | `display/dropdown.templ:132`        | `text-gray-400` no `dark:text-*`                 | Added `dark:text-gray-500` |
| Toast dismiss button (templ)  | `feedback/toast.templ:155`          | `text-gray-400` no `dark:text-*`                 | Added `dark:text-gray-500` |
| Toast dismiss button (JS)     | `feedback/toast.templ:110`          | `text-gray-400` no `dark:text-*` in JS className | Added `dark:text-gray-500` |
| Breadcrumb separator (custom) | `navigation/breadcrumbs.templ:95`   | `text-gray-400` no `dark:text-*`                 | Added `dark:text-gray-500` |
| Breadcrumb chevron icon       | `navigation/breadcrumbs.templ:97`   | `text-gray-400` no `dark:text-*`                 | Added `dark:text-gray-500` |
| Mobile menu toggle            | `navigation/mobile_menu.templ:15`   | `text-gray-400` no `dark:text-*`                 | Added `dark:text-gray-500` |
| Theme toggle button           | `layout/theme.templ:34`             | `text-gray-400` no `dark:text-*`                 | Added `dark:text-gray-500` |
| Sidebar footer text           | `navigation/sidebar_nav.templ:97`   | `text-gray-500` no `dark:text-*`                 | Added `dark:text-gray-400` |
| Step indicator current circle | `feedback/step_indicator.templ:104` | `bg-white` no `dark:bg-*`                        | Added `dark:bg-gray-800`   |
| Errorpage Infrastructure icon | `errorpage/styles.go:84`            | `text-gray-400` no `dark:text-*`                 | Added `dark:text-gray-500` |
| Errorpage Default icon        | `errorpage/styles.go:94`            | `text-gray-400` no `dark:text-*`                 | Added `dark:text-gray-500` |

#### Important: Inconsistent Dark Mode Styling

| Component                       | File                                                               | Issue                                     | Fix                                                       |
| ------------------------------- | ------------------------------------------------------------------ | ----------------------------------------- | --------------------------------------------------------- |
| Button Primary variant          | `display/button_go.go:81`                                          | `bg-blue-600` no `dark:bg-*`              | Added `dark:bg-blue-500 dark:hover:bg-blue-400`           |
| Button Danger variant           | `display/button_go.go:83`                                          | `bg-red-600` no `dark:bg-*`               | Added `dark:bg-red-500 dark:hover:bg-red-400`             |
| Button Link variant             | `display/button_go.go:85`                                          | `text-blue-600` no `dark:text-*`          | Added `dark:text-blue-400 dark:hover:text-blue-300`       |
| Avatar fallback                 | `display/avatar.templ:121`                                         | `bg-blue-600` no `dark:bg-*`              | Added `dark:bg-blue-500`                                  |
| EmptyState action button        | `display/empty_state.templ:27`                                     | `bg-blue-600` no `dark:bg-*`              | Added `dark:bg-blue-500 dark:hover:bg-blue-400`           |
| Tabs pills active               | `display/tabs.templ:118`                                           | `bg-blue-600` no `dark:bg-*`              | Added `dark:bg-blue-500`                                  |
| Pagination active page          | `navigation/pagination.templ:199`                                  | `bg-blue-600` no `dark:bg-*`              | Added `dark:bg-blue-500`                                  |
| SidebarNav active item          | `navigation/sidebar_nav.templ:83`                                  | `bg-blue-600` no `dark:bg-*`              | Added `dark:bg-blue-500`                                  |
| NotFound404 go-home button      | `errorpage/notfound404.templ:57`                                   | `bg-blue-600` no `dark:bg-*`              | Added `dark:bg-blue-500 dark:hover:bg-blue-400`           |
| NotFound404 go-back button      | `errorpage/notfound404.templ:66`                                   | `focus-visible:ring-gray-400` no `dark:*` | Added `dark:focus-visible:ring-gray-500`                  |
| Errorpage family action buttons | `errorpage/styles.go` (all 5 families + default)                   | `bg-*-600` no `dark:bg-*`                 | Added `dark:bg-*-500 dark:hover:bg-*-400` to all          |
| Errorpage family icon colors    | `errorpage/styles.go` (Rejection, Conflict, Transient, Corruption) | `text-*-500` no `dark:text-*-400`         | Added `dark:text-*-400` to all                            |
| Errorpage Infrastructure action | `errorpage/styles.go:87`                                           | `focus-visible:ring-gray-500` no `dark:*` | Added `dark:focus-visible:ring-gray-400`                  |
| StatCard trend up arrow         | `display/card.templ:276`                                           | `text-green-500` no `dark:text-*`         | Added `dark:text-green-400`                               |
| StatCard trend down arrow       | `display/card.templ:278`                                           | `text-red-500` no `dark:text-*`           | Added `dark:text-red-400`                                 |
| Toggle switch track             | `forms/toggle.templ:94`                                            | `peer-checked:bg-blue-600` no `dark:*`    | Added `dark:peer-checked:bg-blue-500`                     |
| FileInput button                | `forms/file_input.templ:26`                                        | `file:bg-blue-600` no `dark:*`            | Added `dark:file:bg-blue-500 dark:hover:file:bg-blue-400` |
| Checkbox input                  | `forms/input.templ:179`                                            | `text-blue-600` no `dark:text-*`          | Added `dark:text-blue-400 dark:focus:ring-blue-500`       |
| Radio input                     | `forms/radio.templ:7`                                              | `text-blue-600` no `dark:text-*`          | Added `dark:text-blue-400 dark:focus:ring-blue-500`       |
| Required asterisk (Label)       | `forms/label.templ:15`                                             | `text-red-500` no `dark:text-*`           | Added `dark:text-red-400`                                 |
| Required asterisk (RadioGroup)  | `forms/radio.templ:85`                                             | `text-red-500` no `dark:text-*`           | Added `dark:text-red-400`                                 |
| Skip-to-content link            | `layout/base.templ:191`                                            | `focus:bg-blue-600` no `dark:*`           | Added `dark:focus:bg-blue-500`                            |

#### Demo Page Fixes

| Issue                           | File                                   | Fix                                                                   |
| ------------------------------- | -------------------------------------- | --------------------------------------------------------------------- |
| Spinner colors missing dark:    | `examples/demo/demo.templ:105,106,107` | Added `dark:text-blue-400`, `dark:text-gray-400`, `dark:text-red-400` |
| Icon name labels missing dark:  | `examples/demo/demo.templ:93`          | Added `dark:text-gray-400`                                            |
| Search icon addon missing dark: | `examples/demo/demo.templ:165`         | Added `dark:text-gray-500`                                            |

### Files Changed

**54 files total:**

- 25 source `.templ` / `.go` files (the actual fixes)
- 25 generated `*_templ.go` files (auto-regenerated)
- 7 golden test files (updated via `-update` flag)
- 1 unit test file (`feedback/helpers_test.go` — updated expected string)

### Verification

- `go build ./...` — passes
- `go test ./...` — all 14 packages pass
- `golangci-lint run` — 0 issues
- `errorpage/styles.go` auto-reformatted by `golangci-lint --fix` (golines alignment)

---

## a) FULLY DONE

1. **Palette consistency audit** — verified no `slate-*`/`zinc-*`/`neutral-*`/`stone-*` mixing anywhere
2. **Neutral color audit** — all `text-gray-*`, `bg-white`, `bg-gray-*`, `border-gray-*`, `ring-gray-*` now have `dark:` variants
3. **Semantic color audit** — all `bg-blue-600`/`bg-red-600`/`text-blue-600`/`text-red-600`/`text-green-500`/`text-amber-500`/`text-orange-500` now have `dark:` variants
4. **Focus ring audit** — all `focus-visible:ring-gray-*` now have `dark:focus-visible:ring-*`
5. **Demo page** — all spinner colors, icon labels, and addon icons fixed
6. **Golden files updated** — 7 golden snapshots regenerated to match new output
7. **Unit test updated** — `TestStepCircleClass` expected string updated
8. **Build + test + lint** — all green

## b) PARTIALLY DONE

1. **Toast JS string literal** — the `tcShowToast` JS function constructs toast HTML including class strings.
   The dismiss button className was fixed, but the toast container/type colors in the JS are still
   hardcoded (e.g., `style.border`, `style.bg`, `style.text` come from the Go `toastStyleMap` which
   IS properly dark-mode aware). The JS path is fine for templ-rendered toasts but dynamically-created
   toasts rely on the same Go map, so this is actually complete — just worth noting the dual path.

2. **Doc comment examples** — 3 doc comments in `.templ` files still show `text-blue-600` without
   `dark:text-blue-400` (in `htmx/loading.templ:15,28` and `feedback/loading.templ:59`). These are
   godoc examples, not rendered code, but could mislead consumers copying them.

## c) NOT STARTED

1. **Automated dark mode test** — no test exists that renders every component in both light and dark
   mode and asserts contrast/readability. The golden files capture HTML output but don't validate
   that `dark:` classes are present.
2. **Contrast ratio verification** — no WCAG contrast ratio checking was done. The `dark:` variants
   use Tailwind's default palette (e.g., `dark:text-gray-400` on `dark:bg-gray-900`), which should
   be fine, but was not verified programmatically.
3. **Tailwind v4 content scanning verification** — did not verify that all `dark:` classes added
   in Go string literals (e.g., `button_go.go`, `errorpage/styles.go`) are actually picked up by
   Tailwind v4's content scanner. Since these are in `.go` files and Tailwind scans raw text, they
   should be found, but no build-time CSS verification was done.

## d) TOTALLY FUCKED UP

1. **`errorpage/handler.go` unrelated change** — The working tree had a pre-existing uncommitted
   change to `errorpage/handler.go` (`encoding/json` → `encoding/json/v2`) that I did NOT make and
   did NOT notice until reviewing the git diff. It's now mixed into my changes. I did not touch it
   (correctly, per the "don't revert changes you didn't author" rule), but it will show up in any
   commit. **This needs to be separated before committing.**

2. **Initial edit failures** — Two edits failed because I tried to edit files I hadn't read yet in
   this conversation (`forms/input.templ` and `layout/base.templ`). The tool correctly blocked these.
   I had to re-read the files and retry. Not a real fuck-up, but a workflow inefficiency.

## e) WHAT WE SHOULD IMPROVE

1. **No dark mode compliance test** — The repo has `utils.TestMotionReduceCompliance` that scans all
   `.templ` files for missing `motion-reduce:` classes. We need an equivalent `TestDarkModeCompliance`
   that scans for neutral color classes (`text-gray-*`, `bg-white`, `bg-gray-*`, `border-gray-*`)
   without `dark:` variants. This would have caught all 30+ issues automatically.

2. **Inconsistent dark: variant patterns** — Some components use `dark:bg-blue-500` (lighter shade),
   others use `dark:bg-blue-600` (same shade). The convention should be documented: in dark mode,
   primary colors should use the `-500` shade (slightly lighter than `-600` for better visibility
   on dark backgrounds). This was done consistently in this fix, but the convention isn't written
   down anywhere.

3. **Doc comments don't show dark mode** — The godoc examples in `.templ` files show light-mode-only
   color classes. Consumers copying these examples won't get dark mode support. All doc examples
   should include `dark:` variants.

4. **Sidebar is a special case** — `SidebarNav` uses `bg-gray-900` as its base (permanently dark
   sidebar), so its `hover:bg-gray-800` is intentional. This is the only component with this pattern.
   It should be documented as an intentional exception.

5. **Toggle thumb stays white** — The toggle thumb is `bg-white` in both modes (intentional — the
   track changes color instead). This is the only `bg-white` without a `dark:` variant. Should be
   documented as intentional.

6. **`errorpage/styles.go` got reformatted** — The `golangci-lint --fix` reformatted the struct
   literals from multi-field-per-line to one-field-per-line. This is cosmetically different from
   the original style. Not wrong, but worth noting.

## f) Up to 50 Things We Should Get Done Next

### Dark Mode Specific (high priority)

1. Add `utils.TestDarkModeCompliance` — scan all `.templ` files for neutral color classes without `dark:` variants
2. Add `utils.TestDarkModeSemanticColors` — scan for `bg-blue-600`/`bg-red-600`/`text-blue-600` etc. without `dark:` variants
3. Document the dark mode convention in AGENTS.md: "`-600` in light → `-500` in dark for primary/action colors; `-400` for text"
4. Update doc comments in `htmx/loading.templ` and `feedback/loading.templ` to include `dark:` variants in example code
5. Add a dark mode integration test that renders components with `.dark` class wrapper and asserts `dark:` classes appear in output
6. Verify WCAG AA contrast ratios for all dark mode color combinations (gray-400 on gray-900, blue-400 on gray-800, etc.)
7. Add a `darkMode` prop to the demo page so you can toggle between light/dark previews
8. Consider adding `prefers-color-scheme` support as a fallback when `.dark` class is absent
9. Audit all `hover:` variants for missing `dark:hover:` counterparts (separate from this audit which focused on base colors)
10. Audit all `focus:` variants for missing `dark:focus:` counterparts

### Pre-existing Working Tree Cleanup

11. Separate the `errorpage/handler.go` `encoding/json/v2` change from this dark mode work before committing
12. Check if `encoding/json/v2` is actually available in Go 1.26 (it's a Go 1.24+ experimental package — verify)

### General Component Quality

13. Fix the `progressbar.templ:54` templ hint: "if statement can be modernized using max" (pre-existing LSP diagnostic)
14. Add a visual regression test that renders all components in both light and dark mode and compares screenshots
15. Add `dark:` variants to all remaining doc comments across the codebase
16. Consider adding a `Theme` context value so components can programmatically know if they're rendering in dark mode
17. Add a `prefers-reduced-transparency` media query for overlay components (Modal, Drawer) — separate from motion-reduce
18. Audit `ring-offset-*` classes for dark mode (ring-offset-white without dark:ring-offset-gray-900)
19. Audit `shadow-*` classes for dark mode shadows
20. Audit `backdrop-blur-*` opacity values for dark mode visibility

### Testing Infrastructure

21. Add a test that asserts every component with `BaseProps` propagates `dark:` classes from `props.Class`
22. Add a contract test that every props struct with a `Color` field includes dark mode guidance in its godoc
23. Add a golden test variant that renders with `.dark` parent class
24. Add a test for the toast JS-created toast (currently only templ-rendered toasts are golden-tested)
25. Benchmark dark mode class resolution (tailwind-merge-go may have different performance with longer class strings)

### Documentation

26. Update AGENTS.md with the dark mode convention documentation
27. Add a `docs/dark-mode-guide.md` explaining the `-600`→`-500` pattern and gray-only palette rule
28. Update `README.md` dark mode section with the complete convention
29. Update `SKILL.md` Part 2 with dark mode checklist for new components
30. Add an ADR for the dark mode color shade convention

### Accessibility

31. Verify screen reader behavior with `dark:` color changes (should be no-op since SR ignores CSS)
32. Add `color-scheme: dark` CSS property to the `.dark` class (improves native form controls in dark mode)
33. Audit native form controls (checkbox, radio, select) for `color-scheme` support
34. Add `scrollbar-color` CSS for dark mode (default scrollbars are jarring on dark backgrounds)
35. Verify `::selection` colors work in dark mode (already handled in `base.templ` BodyClass)

### Component-Specific

36. `SidebarNav` — consider making the permanently-dark sidebar optional (prop to switch to light sidebar)
37. `CountBadge` — verify `ring-white dark:ring-gray-800` is sufficient (the ring is the badge border on the parent)
38. `Tooltip` — verify tooltip arrow border colors work in dark mode (already has `dark:border-*-gray-700`)
39. `Table` — audit `hover:bg-gray-50` for `dark:hover:bg-gray-800` (may already be done — verify)
40. `Nav` — verify the mobile menu slide animation is visible in dark mode

### Build & Release

41. Commit the dark mode fixes (after separating `handler.go`)
42. Bump `utils.Version` and add CHANGELOG entry for the dark mode audit
43. Add a release note highlighting the dark mode improvements
44. Tag a patch release (this is a bug fix, not a feature)
45. Update `FEATURES.md` to note full dark mode compliance

### Future-Proofing

46. Consider a `darkMode()` helper in `utils` that returns `dark:` prefixed classes for a given light class
47. Consider a CSS-first approach using `@custom-variant dark` + CSS variables instead of per-class `dark:` variants
48. Explore Tailwind v4's `@theme` dark mode tokens for a more maintainable approach
49. Consider adding a `Theme` enum (Light/Dark/Auto) to `layout` for more granular control
50. Add a pre-commit hook that checks for missing `dark:` variants (like the motion-reduce check)

## g) Top 2 Questions I Cannot Answer Myself

### Q1: Should the dark mode compliance test be failing or informational?

The `utils.TestMotionReduceCompliance` test is failing (blocks CI). Should the dark mode equivalent
also be failing, or should it be informational only (like `TestSkillComponentCount`)? A failing test
means every new component MUST include dark mode from day one — which is the right policy, but it's
a judgment call about enforcement level.

### Q2: Is the `encoding/json/v2` change in `handler.go` intentional and ready to commit?

There's a pre-existing uncommitted change in `errorpage/handler.go` switching from `encoding/json`
to `encoding/json/v2`. I did not make this change and don't know its context. It may be from a
previous session or an experiment. Should it be:

- Committed separately as its own change?
- Reverted if it's experimental?
- Left as-is for now?

This needs a human decision before committing the dark mode work.

---

## Summary

**Before this session:** 30+ dark mode gaps across 8 packages — icons invisible, buttons
inconsistent, action buttons without dark hover states, required asterisks invisible in dark mode.

**After this session:** All neutral and semantic color classes have `dark:` variants. All tests
pass. Zero lint issues. Golden files updated. One unrelated `handler.go` change needs separating
before commit.

**The codebase now properly supports light AND dark mode.**
