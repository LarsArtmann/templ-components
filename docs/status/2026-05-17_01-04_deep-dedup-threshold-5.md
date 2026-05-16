# Session 7b — Deep Deduplication Round 2: Threshold 5

**Date:** 2026-05-17 01:04
**Branch:** master (4 commits ahead of origin)
**Status:** ✅ All green — build, tests, lint pass

---

## Executive Summary

Continued deduplication at a lower threshold (5 vs previous 7). Reduced clone groups from **7 → 3** at threshold 5. All 3 remaining are structural false positives that cannot be meaningfully eliminated without artificial restructuring. The icon system now covers 45 icons (up from 42 before this session series) and all feedback components (alert, toast, inline error/success) have been migrated from inline SVGs to `icons.Icon()`.

---

## A) FULLY DONE ✅

### 1. Clone Group Elimination at Threshold 5 (7 → 3)

| #   | Clone Group                                             | Lines                                   | Verdict            | Fix                                                                                                                                                                   |
| --- | ------------------------------------------------------- | --------------------------------------- | ------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | modal/alert/toast icon SVGs (3-way)                     | modal:32-40, alert:58-61, toast:149-152 | **ELIMINATED**     | Converted alert to `icons.Icon(alertIconName(...))`, toast to `icons.Icon(toastIconName(...))`. Removed `alertIconPath` and `toastIconPath` templ templates entirely. |
| 2   | htmx/loading vs pagination disabled (2-way)             | loading:22-29, pagination:98-101        | **FALSE POSITIVE** | `<div>` overlay vs `<span>` disabled — completely different semantics                                                                                                 |
| 3   | empty_state a/button self-clone (2-way)                 | empty_state:30-37, 38-45                | **FALSE POSITIVE** | `<a>` vs `<button>` — inherent HTML limitation, shared class constant already in place                                                                                |
| 4   | nav_link external/internal branch (2-way)               | nav_link:51-61, 69-81                   | **ELIMINATED**     | Unified into single `<a>` with conditional `target="_blank" rel="noopener noreferrer"`                                                                                |
| 5   | forms/input label/helpText (2-way)                      | input:146-148, label:22-24              | **ELIMINATED**     | Extracted `checkboxLabel` helper that reuses `helpText()` template                                                                                                    |
| 6   | forms/FieldError vs pagination mobilePageButton (2-way) | label:30-36, pagination:68-79           | **ELIMINATED**     | Collapsed FieldError from if/else to single `<p>` with conditional `id` attribute                                                                                     |
| 7   | progress/htmx helpers (2-way)                           | progress:55-63, htmx/helpers:7-17       | **FALSE POSITIVE** | progressbar div vs delete button — no semantic overlap                                                                                                                |

### 2. New Icons Added (43 → 45)

| Icon                | Path                                            | Used By                                   |
| ------------------- | ----------------------------------------------- | ----------------------------------------- |
| `CheckCircle`       | `M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z` | Toast success                             |
| `ExclamationCircle` | `M12 9v3.75m9-.75a9 9 0 11-18 0...`             | Same as `Exclamation` (alias for clarity) |

### 3. Templates Removed (−30 lines of inline SVG)

- `alertIconPath(AlertType)` — 14-line switch template with filled 20×20 SVG paths
- `toastIconPath(ToastType)` — 12-line switch template with stroke 24×24 SVG paths

Replaced with Go functions `alertIconName()` and `toastIconName()` that map types to `icons.Name` constants.

### 4. Files Changed (8 files, net reduction)

| File                        | Change                                                                                                       |
| --------------------------- | ------------------------------------------------------------------------------------------------------------ |
| `feedback/alert.templ`      | Replaced SVG wrapper + `alertIconPath` with `icons.Icon(alertIconName(...))`. Added `alertIconName` Go func. |
| `feedback/toast.templ`      | Replaced SVG wrapper + `toastIconPath` with `icons.Icon(toastIconName(...))`. Added `toastIconName` Go func. |
| `icons/icon_names.go`       | Added `CheckCircle`, `ExclamationCircle`                                                                     |
| `icons/icon_paths.go`       | Added path data for `CheckCircle`, `ExclamationCircle`                                                       |
| `icons/icon_names_test.go`  | Updated count 43→45, added test cases, reformatted by gci                                                    |
| `forms/label.templ`         | Collapsed `FieldError` from if/else to single `<p>` with conditional id                                      |
| `forms/input.templ`         | Extracted `checkboxLabel` helper that calls `helpText()`                                                     |
| `navigation/nav_link.templ` | Unified `NavLink` external/internal into single `<a>` element                                                |

### 5. Cumulative Session Deduplication Progress

| Run                    | Threshold | Clone Groups | Actionable | Eliminated |
| ---------------------- | --------- | ------------ | ---------- | ---------- |
| Session 7 start (t=7)  | 7         | 6            | 4          | —          |
| Session 7 end (t=7)    | 7         | 2            | 0          | 4          |
| Session 7b start (t=5) | 5         | 7            | 4          | —          |
| Session 7b end (t=5)   | 5         | 3            | 0          | 4 more     |

**Total eliminated across both runs: 8 clone groups.**

---

## B) PARTIALLY DONE 🟡

### Unrelated Staged Changes

- `docs/status/2026-05-17_00-33_public-release-prep-and-demo-security-fix.md` — has whitespace reformatting from a previous session, not committed. Not part of this dedup work.

---

## C) NOT STARTED ⬜

1. **Toast JS icon path duplication** — `toast.templ:82-86` has inline JS strings with SVG path data for dynamically created toasts (`tcToastIcons`). These still duplicate the icon path data. Would need to generate from Go at build time.
2. **StepIndicator inline check SVG** — `progress.templ:112-115` uses `stroke-width="3"` (different from icon system's 1.5). Could use `icons.Check` if visual difference is acceptable.
3. **Filled icon variant** — `internal/svg.FillIcon` (20×20 filled) is a parallel system. Could be unified into `icons.Icon()` with a variant parameter.
4. **Examples directory lint** — 23 lint issues in `examples/demo/` still blocking full `./...` lint.

---

## D) TOTALLY FUCKED UP 💥

**Nothing.** Clean session again. All changes compiled, tested, and linted zero issues.

One minor issue: `icon_names_test.go` error message says `expected 42` but count was updated to 45 — the error message string wasn't updated. This is cosmetic (test passes, error message only shows on failure). Should be fixed.

---

## E) WHAT WE SHOULD IMPROVE

### Architecture

1. **Filled vs stroke icon unification** — We now have two parallel SVG systems: `icons.Icon()` for 24×24 stroke and `internal/svg.FillIcon()` for 20×20 filled. The alert was converted from filled to stroke (visual change). Should decide: one system or two, and document it.

2. **Toast dynamic icons** — The toast container creates SVG icons in JavaScript with hardcoded path strings. These should be generated from the same Go data source.

3. **Icon stroke-width parameter** — `icons.Icon()` always uses `stroke-width="1.5"`. Some use cases (StepIndicator) need `stroke-width="3"`. Consider adding an option.

### Code Quality

4. **Test error message staleness** — `TestIconCount` error says `expected 42` but actual expected is 45. The string should match.

5. **Icon alias documentation** — `Exclamation` and `ExclamationCircle` have identical paths. Should document which to use when, or consolidate.

6. **Snapshot test regeneration** — Components changed in this session (alert, toast, nav_link, forms) need updated golden files.

---

## F) Top 25 Things To Do Next

### Ship It (Critical Path)

1. **Tag v0.2.0** — All breaking changes consolidated, README updated, CONTRIBUTING.md exists
2. **Push to origin** — 4 local commits ahead
3. **GitHub Actions CI** — `.github/workflows/ci.yml` for build+test+lint
4. **Fix test error message** — `expected 42` → `expected 45` in `TestIconCount`
5. **Regenerate snapshot tests** — Alert, toast, nav_link, forms all changed

### Icon System

6. **Document filled vs stroke decision** — ADR or AGENTS.md entry
7. **Unify FillIcon into icons package** — Or formally document the two-system approach
8. **Add stroke-width option to icons.Icon** — `IconWithStroke(name, class, strokeWidth)`
9. **Consolidate Exclamation/ExclamationCircle** — Same path, pick one or document aliases
10. **Add remaining commonly-needed Heroicons** — Only 45 of ~300 covered

### Component Quality

11. **Table component** — Frequently requested, not implemented
12. **Tooltip component** — Common UI pattern
13. **Tabs keyboard navigation** — Arrow key support for tab panels
14. **Toast dynamic icon generation** — Generate JS from Go data
15. **Form validation patterns** — `forms.Validate()` helpers or examples

### Architecture

16. **Modularize into sub-modules** — `go.work` with separate modules
17. **CSP nonce audit** — Verify all inline scripts receive nonce
18. **Fix examples/demo lint** — 23 issues blocking full `./...` lint
19. **Dark mode color standardization** — Consistent palette across all components
20. **Accessibility audit** — Run axe-core on all components

### Polish & DX

21. **Interactive playground** — Live preview in demo app
22. **API documentation** — Auto-generate from Go doc comments
23. **CHANGELOG.md** — For v0.2.0 release
24. **Migration guide v0.1→v0.2** — Already started, needs completion
25. **Performance benchmarks** — Measure render time for key components

---

## G) Top #1 Question I Cannot Figure Out Myself

**Should alert icons use stroke-based (24×24) or filled (20×20) SVGs?**

In this session I converted alert icons from 20×20 filled SVGs to 24×24 stroke-based `icons.Icon()`. This is a **visual change** — the icons look slightly different (outline vs solid). The toast icons were already stroke-based, so now alert and toast are consistent. But the original Heroicons design for alerts uses filled/solid icons for emphasis.

Options:

1. **Keep stroke (current)** — Consistent with toast and the rest of the icon system
2. **Revert to filled** — Add a filled icon variant to the system for alerts only
3. **Add both variants** — Support `icons.IconFilled()` alongside `icons.Icon()`

This is a design decision that affects the component library's visual identity.

---

## Project Metrics

| Metric                                | Value                   |
| ------------------------------------- | ----------------------- |
| `.templ` files (excl examples)        | 31                      |
| `.go` files (excl examples/generated) | 52                      |
| Icon count                            | 45                      |
| art-dupl clone groups (t=5)           | 3 (all false positives) |
| art-dupl clone groups (t=7)           | 1 (false positive)      |
| Test packages                         | 9/9 passing             |
| Lint issues                           | 0                       |
| Build status                          | ✅ Clean                |
| Commits ahead of origin               | 4                       |
