# Status Report: Snapshot Test Infrastructure Overhaul

**Date:** 2026-07-30 01:55
**Session Goal:** "How can we better use snap tests? Especially for UI tests?"
**Verdict:** Framework upgraded and coverage expanded significantly, but several gaps remain and a few things were missed entirely.

---

## A) FULLY DONE

### 1. ID Normalization in Golden Package (`internal/golden/golden.go`)

**Problem solved:** `utils.EnsureID` uses `crypto/rand` to generate IDs like `tc-modal-a1b2c3d4e5f6a7b8`. Every render produces different output. The old golden package had NO normalization for this — every EnsureID component required a manual `ID: "something"` in test props or the golden test would fail non-deterministically.

**What was built:**

- `autoIDRe` regex matching both EnsureID formats: `tc-<prefix>-<16hex>` (crypto/rand primary) and `tc-<prefix>-<digits>-<digits>` (fallback timestamp+counter)
- Supports hyphenated prefixes (`tc-mobile-menu-*` used by Nav)
- `normalizeIDs()` replaces all matches with `tc-<prefix>-NORMALIZED`
- Operates on full document text, so cross-references (aria-controls, popovertarget, for, href="#...", inline JS) are all caught
- Explicit consumer IDs (`id="my-custom-id"`) are preserved — regex requires exactly 16 hex chars or the digit-digit pattern
- 6 test cases covering: hex format, fallback format, cross-reference attrs, explicit IDs preserved, short hex rejected, hyphenated prefix

**Files changed:**

- `internal/golden/golden.go` — new `normalize()`, `normalizeIDs()`, `autoIDRe`
- `internal/golden/golden_coverage_test.go` — `TestNormalizeIDs` with 6 subtests, `TestGoldenDiffLCSAlignment`

### 2. LCS-Based Diff Algorithm

**Problem solved:** The old diff was a naive line-by-line comparison (`wantLines[i] != gotLines[i]`). An insertion at line 3 would cascade into showing EVERY subsequent line as "changed" — making diffs nearly unreadable for any non-trivial change.

**What was built:**

- Longest Common Subsequence (LCS) DP table alignment
- Only actually-changed lines appear in the diff
- 1-based line numbers on each diff line (`--- [3] <div>` instead of `--- <div>`)
- `nolint:makezero` for the pre-allocated DP table (intentional index assignment pattern)

**Files changed:**

- `internal/golden/golden.go` — replaced `diff()` and removed `lineAt()`
- `internal/golden/golden_test.go` — updated `TestDiffOutput` assertions for new format
- `internal/golden/golden_coverage_test.go` — `TestGoldenDiffLCSAlignment`

### 3. Table-Driven Golden Helper (`golden.AssertSnapshots`)

**Problem solved:** Each golden test was 10-15 lines of boilerplate: `t.Parallel()` + `t.Run("name", func(t *testing.T) { t.Parallel(); output := utils.Render(...); golden.Assert(t, "name", output) })`. Adding 5 variants of a component meant 75 lines.

**What was built:**

- `golden.Snapshot{Name, HTML}` struct
- `golden.AssertSnapshots(t, []golden.Snapshot{...})` — creates parallel subtests automatically
- Each snapshot becomes `t.Run(name, func(t *testing.T) { t.Parallel(); Assert(t, name, html) })`
- Documented with usage examples in the package doc comment

**Files changed:**

- `internal/golden/golden.go` — new `Snapshot` type + `AssertSnapshots` function

### 4. Golden Test Sweep — 27 Components Covered

**Problem solved:** Only ~38 of 86 components had golden tests. 8 EnsureID components (Accordion, Tabs, Dropdown, Tooltip, Carousel, ContextMenu, Combobox, Nav) had ZERO golden tests because the ID workaround was a known pain point. Entire component categories (Checkbox, Toggle, RadioGroup, Form, ValidationSummary, InlineError, InlineSuccess, Footer, Nav, EndOfList, ErrorAlert, ErrorDetail) were untested.

**What was built — 5 new `golden_sweep_test.go` files:**

| Package      | New Components Covered                                                                                                                                              | New Golden Files |
| ------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------- |
| `display`    | Accordion (2 variants), Tabs (2 variants), Dropdown, Tooltip (2 variants), Carousel, ContextMenu, Avatar (2 variants), EmptyState (2 variants)                      | 14               |
| `forms`      | Checkbox (3 variants), Toggle (3 variants), RadioGroup (2 variants), Combobox, InputGroup, Form, ValidationSummary, FileInput (2 variants), DatePicker (2 variants) | 16               |
| `feedback`   | InlineError, InlineSuccess, SkeletonGroup, ToastContainer                                                                                                           | 4                |
| `navigation` | Nav, SimpleNav, Footer, EndOfList (2 variants)                                                                                                                      | 5                |
| `errorpage`  | ErrorAlert (2 variants), ErrorDetail (2 variants)                                                                                                                   | 4                |

**Total golden files: ~63 → 102 (+39 files, +62%)**

**Files created:**

- `display/golden_sweep_test.go`
- `forms/golden_sweep_test.go`
- `feedback/golden_sweep_test.go`
- `navigation/golden_sweep_test.go`
- `errorpage/golden_sweep_test.go`

### 5. AGENTS.md Documentation Updated

Replaced the single-line golden testing entry with a comprehensive three-tier snapshot strategy section explaining when to use golden tests vs substring assertions vs visual regression, plus a "how to add a new component's golden test" playbook.

---

## B) PARTIALLY DONE

### 1. Visual Regression Test Coverage — NOT TOUCHED

The session focused entirely on HTML golden tests. Visual regression tests (`visualtest`) remain at **10/86 components (12%)** — entirely unchanged. This was the right call for this session (the golden framework needed the ID normalization fix first), but it means the "especially for UI tests" half of the prompt is unaddressed.

**Current visual coverage gaps:**

- 76 components have zero visual tests
- Entire `navigation` (12) and `layout` (10) packages untested visually
- Only 2 RTL variants exist (Button + Card), zero RTL+dark combos
- 2 currently-failing visual tests (`modal/open_light`, `drawer/right_light`) with leftover `.fail/` artifacts — not investigated

### 2. Golden Coverage — Still Gaps

27 new components got golden tests, but some remain uncovered:

- `layout`: AppShell, Split, Stack, Container, Base, Minimal, ThemeToggle, ThemeScript (only `Script` has a golden)
- `htmx`: LoadingIndicator, InlineLoadingOverlay, LoadingButton, ConfirmDelete, SwapOOB, CSRFToken, GlobalErrorHandling — ZERO golden tests
- `display`: SimpleCard, StatCard (has golden but partial), StatusBadge, SimpleEmptyState, HoverCard (only basic variant)
- `forms`: Label, FieldError, FormFieldWrapper, Radio (standalone)
- `recipes`: Dashboard, SettingsLayout, LoginCard — ZERO golden tests

### 3. Substring Test Audit — NOT STARTED

The repo has ~30+ `snapshot_test.go` files using `utils.AssertContains`. Many of these test the SAME components that now have golden tests, creating redundant coverage. No audit was done to identify which substring tests are now redundant and could be replaced by more comprehensive golden tests.

---

## C) NOT STARTED

1. **Visual test infrastructure improvements** — No changes to `visualtest` package at all
2. **Visual test coverage expansion** — Zero new visual tests
3. **Fixing the 2 failing visual tests** (`modal/open_light`, `drawer/right_light`)
4. **Dark-mode golden variants** — Only 1 dark golden exists (`display/dark_golden_test.go`). No systematic dark-mode golden coverage.
5. **RTL golden variants** — Zero RTL golden files exist
6. **Contract test for golden coverage** — No test asserting "every component has at least one golden test"
7. **Migration of existing golden tests** — Old `golden_test.go` / `golden_new_test.go` files still use the old `golden.Assert` pattern individually. Could be refactored to `AssertSnapshots` for consistency.
8. **`layout` package golden tests** — Not started
9. **`htmx` package golden tests** — Not started
10. **`recipes` package golden tests** — Not started
11. **Visual test helper improvements** — No `AssertScreenshotSnapshots` table-driven equivalent was created for visual tests

---

## D) TOTALLY FUCKED UP / MISTAKES MADE

### 1. Left a Typo in a Code Comment (Caught and Fixed)

During the regex edit, accidentally introduced `n//` (stray `n` before comment) in `golden.go`. Caught immediately on next view and fixed. No impact.

### 2. Regex Order Bug (Caught by Test)

Initial regex `tc-([a-z]+)-(?:[a-f0-9]{16}|\d+-\d+)` tried hex first. Since digits are valid hex chars, the fallback format `tc-modal-1753897654321234567-1` partially matched as hex (`NORMALIZED567-1`). Fixed by reordering alternatives: `\d+-\d+` first.

### 3. Hyphenated Prefix Bug (Caught by Test)

Initial regex `[a-z]+` didn't match `mobile-menu` prefix used by Nav. Fixed to `[a-z][a-z-]*`. Caught by the `TestGoldenSweepNav` failure showing raw `tc-mobile-menu-<hex>` IDs in the diff.

### 4. Pre-Existing `.golangci.yml` Regression — NOT FIXED (Not My Change)

`TestGolangciDisabledLinters` fails because `ireturn`, `godoclint`, `testableexamples` re-entered the enable list. This is a documented recurring regression (AGENTS.md T1). I did NOT cause it and did NOT fix it — but it means `go test ./...` is red, which is concerning.

### 5. `visualtest/doc.go` Has Compiler Errors (Pre-Existing, Not Fixed)

`allocCancel` and `sharedAllocCtx` are declared and not used. These are pre-existing gopls compiler errors in the visualtest module. Not my code, not touched.

### 6. Did Not Run `templ generate` or Full Build Verification

Changed only `.go` files (no `.templ` edits), so `templ generate` was not needed. But I also did not run `go build ./...` explicitly — only `go test`. Tests compile the same code, so this is fine, but for completeness the full build should be verified.

### 7. Did Not Run Full Lint Suite

Only ran `golangci-lint` on `internal/golden/...`. Did NOT lint the new `golden_sweep_test.go` files in `display`, `forms`, `feedback`, `navigation`, `errorpage`. There could be lint findings (varnamelen, mnd, etc.) in the new test files.

---

## E) WHAT WE SHOULD IMPROVE

### Architecture & Strategy

1. **The three-tier testing strategy needs formalization.** The tiers (golden → substring → visual) are documented in AGENTS.md but there's no enforcement. A contract test should assert every component has at least a golden test.

2. **Dark-mode golden testing is undersolved.** Golden tests render light-mode HTML. Dark-mode classes (`dark:bg-gray-900`) are IN the output, but we're not testing that they're correct for each component variant. The existing `dark_golden_test.go` only covers Badge, Button, Card. A systematic approach would render each component and assert dark-mode class compliance via golden.

3. **The `layout` and `htmx` packages are testing deserts.** `layout` has 10 components with 1 golden test. `htmx` has 8 components with 0 golden tests. These are foundational packages.

4. **Substring tests and golden tests overlap.** Many components have BOTH a `snapshot_test.go` (substring) AND now a `golden_sweep_test.go`. The substring tests are less comprehensive but faster to understand. Need a policy: golden as the primary, substring for targeted invariants only.

5. **Visual regression tests need the same table-driven treatment.** `visualtest.AssertScreenshot` has no batch equivalent. Adding 76 components × (light + dark) = 152 screenshots individually is impractical.

6. **The 2 failing visual tests are a liability.** `modal/open_light` and `drawer/right_light` have `.fail/` artifacts committed. Failing tests erode trust in the test suite. Either fix them or mark them as known-failures with `t.Skip`.

### Framework

7. **`golden.AssertSnapshots` could support directory names.** Currently names must be flat (`accordion_default`). Supporting `display/accordion_default` would let packages share a single `testdata/` or organize by category.

8. **No `-update` dry-run mode.** Running `-update` silently overwrites. A `--diff` mode that shows what WOULD change without writing would be safer for CI.

9. **No golden file count drift guard.** A test asserting "golden file count >= N" would catch accidentally deleted golden files.

10. **Normalization could be extensible.** Currently hardcoded (classes + IDs). If components start using timestamps or nonces in output, we'd need another regex. A `normalizers []func(string) string` chain would be more future-proof.

### Quality

11. **New test files are NOT lint-clean verified.** Ran `golangci-lint` only on `internal/golden/...`. The 5 new `golden_sweep_test.go` files might have `mnd` (magic number detector) or `varnamelen` findings.

12. **No benchmark for the normalization pipeline.** The LCS diff is O(n*m) — fine for small golden files but could be slow on very large ones. No measurement was done.

13. **The `boolPtr` function in `golden_coverage_test.go` was removed** (it was flagged as unused by gopls). This is correct but means the file's line count changed — if any tooling counts test lines, it'll be off.

---

## F) UP TO 50 THINGS TO DO NEXT

### High Priority (Framework & Coverage Gaps)

1. Run `golangci-lint` on all new `golden_sweep_test.go` files and fix findings
2. Fix the pre-existing `.golangci.yml` regression (ireturn/godoclint/testableexamples re-enabled)
3. Add golden tests for `htmx` package (8 components, 0 golden tests)
4. Add golden tests for `layout` package (AppShell, Split, Stack, Container, Base, Minimal, ThemeToggle)
5. Add golden tests for `recipes` package (Dashboard, SettingsLayout, LoginCard)
6. Add golden tests for remaining `display` components (SimpleCard, StatusBadge, SimpleEmptyState)
7. Add golden tests for remaining `forms` components (Label, FieldError, FormFieldWrapper, Radio standalone)
8. Create a contract test: "every exported component function has at least one golden test"
9. Fix the 2 failing visual tests (`modal/open_light`, `drawer/right_light`)
10. Investigate and clean up `.fail/` artifacts in `visualtest/testdata/`

### Medium Priority (Visual Regression Expansion)

11. Create `AssertScreenshotSnapshots` table-driven helper for visual tests
12. Add visual tests for `feedback` components (Spinner, ProgressBar, Skeleton, StepIndicator)
13. Add visual tests for `forms` components (Checkbox, Toggle, RadioGroup, Combobox, Textarea)
14. Add visual tests for `navigation` components (Pagination, Breadcrumbs, SidebarNav)
15. Add visual tests for `layout` components (AppShell, Container, Split, Stack)
16. Add dark-mode visual variants for Popover and ContextMenu (currently light-only)
17. Add RTL visual variants for more components (currently only Button + Card)
18. Add RTL + dark combo visual tests (currently zero)
19. Add responsive viewport visual tests for Nav (hamburger collapse), Grid, Table
20. Add interaction-state visual tests for Tab switching, Accordion expand/collapse

### Medium Priority (Golden Improvements)

21. Migrate existing `golden_test.go` / `golden_new_test.go` to `AssertSnapshots` pattern
22. Add dark-mode golden variants systematically (render each component, golden the dark output)
23. Add golden tests for component edge cases (empty props, nil children, error states)
24. Add golden tests for all enum variants (BadgeType × BadgeSize matrix, AlertType × Dismissible, etc.)
25. Add golden tests for container-aware variants (`ContainerAware: true`)
26. Normalize nonces in golden output (currently `nonce=""` — if components start embedding actual nonces)
27. Add a golden file count drift guard test
28. Consider HTML normalization (attribute ordering, whitespace normalization) beyond just CSS classes
29. Add `--dry-run` mode to `-update` flag that shows diffs without writing
30. Document the normalization pipeline in `docs/`

### Lower Priority (Cleanup & Polish)

31. Audit substring tests (`snapshot_test.go` files) for redundancy with new golden tests
32. Remove redundant substring tests where golden coverage is strictly more comprehensive
33. Keep substring tests only for targeted invariant checks (aria attributes, specific structural assertions)
34. Add a `TESTING.md` doc explaining the three-tier strategy with examples
35. Benchmark the normalization + diff pipeline on the largest golden file
36. Consider using `htmlparse`/`golang.org/x/net/html` for structural normalization instead of regex
37. Add fuzz test for `normalizeIDs` (random ID-like strings should not crash)
38. Add fuzz test for `normalizeClasses` (malformed class attributes)
39. Add fuzz test for `diff` (arbitrary strings should not crash)
40. Add fuzz test for `AssertSnapshots` (empty names, duplicate names)
41. Add a test that renders every component with `Default*Props()` and asserts no panic
42. Add a test that renders every component with zero-value props and asserts no panic
43. Add golden tests for `htmx` interaction states (HTMX loading, error, retry)
44. Add golden tests for `errorpage` handler JSON mode
45. Add golden tests for `errorpage` `WriteError`/`WriteErrorPage` wrappers
46. Consider snapshot testing for generated CSS (detect unintended style changes)
47. Add visual test for Carousel slide navigation (arrow click state)
48. Add visual test for Tooltip positioning (all 4 positions)
49. Add visual test for Dropdown open state (all positions)
50. Add visual test for Table with sortable headers (ascending/descending indicators)

---

## G) QUESTIONS I CANNOT ANSWER MYSELF

### 1. Should we invest in visual regression test expansion now, or is golden coverage sufficient?

Golden tests catch structural/attribute regressions but NOT visual layout issues (color shifts, spacing, overflow, dark-mode rendering). The visual test suite covers only 12% of components. Expanding it to match golden coverage would require ~150 new screenshots and significant chromedp test infrastructure work. **Is the ROI there, or should we keep visual tests focused on the highest-risk components (overlays, responsive layouts, dark-mode color-sensitive components)?**

### 2. Should we fix the pre-existing `.golangci.yml` regression as part of this work?

The `TestGolangciDisabledLinters` test is failing because `ireturn`, `godoclint`, and `testableexamples` re-entered the enable list. This is a documented recurring issue (AGENTS.md T1, "5 times across sessions"). It's not related to snapshot testing, but it means `go test ./...` is red right now. **Should I fix this in this branch, or leave it for a separate dedicated session?**

### 3. What is the intended relationship between substring tests and golden tests going forward?

Many components now have BOTH a `snapshot_test.go` (substring assertions like `AssertContains(output, "aria-label")`) AND a `golden_sweep_test.go` (full HTML golden). The golden test is strictly more comprehensive — if the golden matches, the structural invariants are implicitly verified. **Should we deprecate/remove substring tests where golden coverage exists, or keep both as defense-in-depth?**
