# Status Report: Snapshot Test Infrastructure — Post-Push Reality Check

**Date:** 2026-07-30 02:20
**Session Span:** 2 sessions (initial work + self-review/fix/push)
**Current State:** All code committed and pushed to `origin/master`. Tests green, lint clean.
**Verdict:** Solid framework improvements shipped. Coverage gaps remain. Commit history is messy due to BuildFlow auto-commits.

---

## A) FULLY DONE

### 1. Golden Framework — 3 Major Upgrades (`internal/golden/golden.go`)

**ID Normalization** — `normalizeIDs()` replaces non-deterministic `EnsureID` output (`tc-modal-a1b2c3d4e5f6a7b8`) with stable `tc-modal-NORMALIZED`. Handles both crypto/rand hex format and timestamp+counter fallback. Supports hyphenated prefixes (`tc-mobile-menu-*`). 6 test cases covering all formats including edge cases.

**LCS-Based Diff** — Replaced naive line-by-line diff with Longest Common Subsequence alignment. Insertions/deletions no longer cascade. 1-based line numbers on output. Lint-clean (`makezero` + `varnamelen` nolint with justification).

**Table-Driven Helper** — `golden.AssertSnapshots(t, []golden.Snapshot{{Name, HTML}...})` eliminates per-test `t.Run` + `Assert` boilerplate. Each snapshot runs as a parallel subtest.

### 2. Golden Test Coverage — 27 Components Added

5 new `golden_sweep_test.go` files across 5 packages:

| Package | New Components | New `.golden` Files |
|---------|---------------|-------------------|
| `display` | Accordion, Tabs, Dropdown, Tooltip, Carousel, ContextMenu, Avatar, EmptyState | 14 |
| `forms` | Checkbox, Toggle, RadioGroup, Combobox, InputGroup, Form, ValidationSummary, FileInput, DatePicker | 16 |
| `feedback` | InlineError, InlineSuccess, SkeletonGroup, ToastContainer | 4 |
| `navigation` | Nav, SimpleNav, Footer, EndOfList | 5 |
| `errorpage` | ErrorAlert, ErrorDetail | 4 |

**Total golden files: 63 → 102 (+62%)**

### 3. Pre-Commit Blockers Fixed

| Issue | Root Cause | Fix |
|-------|-----------|-----|
| `.golangci.yml` regression (#6) | ireturn, godoclint, testableexamples re-enabled | Removed from enable list + deleted dead `ireturn:` settings block |
| `visualtest/doc.go` compile error | `:=` shadowed package-level `sharedAllocCtx`/`allocCancel` inside `sync.Once.Do` | Changed to `=` per existing code comment |
| gci import ordering in sweep tests | Wrong import grouping | Fixed to match existing pattern |

### 4. Documentation

Updated AGENTS.md with comprehensive three-tier snapshot testing strategy:
1. HTML golden tests (fast, deterministic — the backbone)
2. Substring assertions (targeted invariant checks)
3. Visual regression (pixel-level, separate module)

### 5. Verification Status

```
go test ./...     → 15/15 packages PASS
golangci-lint     → 0 issues
git status        → clean, pushed to origin/master
```

---

## B) PARTIALLY DONE

### 1. Golden Test Coverage — Still Gaps

**102 golden files across 5 packages, but 3 packages have ZERO golden tests:**

| Package | Components | Golden Files | Coverage |
|---------|-----------|-------------|----------|
| `htmx` | 8 | 0 | **0%** |
| `layout` | 10 | 1 (only `Script`) | **10%** |
| `recipes` | 3 | 0 | **0%** |

**Remaining gaps in covered packages:**
- `display`: SimpleCard, StatusBadge, SimpleEmptyState (3 missing)
- `forms`: Label, FieldError, FormFieldWrapper, Radio standalone (4 missing)

### 2. Visual Regression Coverage — Untouched (12%)

Still 10/86 components have visual tests (12%). The entire session focused on HTML golden tests. Visual coverage gaps:

| Package | Visual Tests | Components |
|---------|-------------|-----------|
| `navigation` | 0/12 | Zero |
| `layout` | 0/10 | Zero |
| `feedback` | 1/13 | Only Alert |
| `forms` | 2/21 | Only Input + Select |
| `display` | 7/30 | Button, Card, Badge, Modal, Drawer, Dropdown, Popover, ContextMenu |

Only 2 RTL variants (Button + Card). Zero RTL+dark combos. 2 failing visual tests with `.fail/` artifacts.

### 3. Commit History Quality — Messy

BuildFlow auto-committed several times with low-quality messages:
- `582a0e3` — "Looking at the changes in the internal/golden/ package, I need to generate a detailed commit message explaining the changes." (hallucinated AI thinking leaked into commit message)
- `3bea5b8` — "test(visualtest): initialize visual testing module with package documentation" (misleading — it was a compile fix, not initialization)
- Multiple commits for what should have been 2-3 clean commits

This is a known BuildFlow daemon issue (documented in AGENTS.md T13).

---

## C) NOT STARTED

1. **Golden tests for `htmx` package** — 8 components, 0 tests
2. **Golden tests for `layout` package** — AppShell, Split, Stack, Container, Base, Minimal, ThemeToggle, ThemeScript, Stylesheet
3. **Golden tests for `recipes` package** — Dashboard, SettingsLayout, LoginCard
4. **Remaining display gaps** — SimpleCard, StatusBadge, SimpleEmptyState
5. **Remaining forms gaps** — Label, FieldError, FormFieldWrapper, Radio
6. **Visual regression expansion** — 76 components untested
7. **Fixing 2 failing visual tests** (`modal/open_light`, `drawer/right_light`)
8. **Cleaning `.fail/` artifacts** in `visualtest/testdata/.fail/`
9. **Dark-mode golden variants** — only `display/dark_golden_test.go` exists (3 components)
10. **RTL golden variants** — zero exist
11. **Golden coverage drift guard** — no contract test asserting every component has a golden
12. **Migrating existing golden tests** to `AssertSnapshots` pattern
13. **Substring test audit** — identify redundant `snapshot_test.go` checks
14. **`AssertScreenshotSnapshots`** — table-driven helper for visual tests
15. **Status report from session 1** (`docs/status/2026-07-30_01-55_*.md`) — content is accurate but wasn't committed until this session

---

## D) TOTALLY FUCKED UP

### 1. Git Index Corruption — Undetected for Entire First Session

The git index was corrupt (`error: index uses ޓ extension`). I didn't run `git status` at the START of the session, so I didn't discover this until the self-review phase. All work from session 1 was uncommitted and could have been lost.

**Fix applied:** `mv .git/index .git/index.corrupt && git read-tree HEAD`

**Lesson:** Always run `git status` as the first action in any session.

### 2. BuildFlow Commit Messages — Quality Disaster

The BuildFlow auto-commit daemon produced:
- A commit message that starts with "Looking at the changes..." (AI thinking leaked)
- Generic "chore: update project configuration" messages
- Misleading "initialize visual testing module" for what was a bug fix

These are now permanently in the git history on `origin/master`. Documented in AGENTS.md T13 but not fixable retroactively without history rewrite (which we won't do).

### 3. Pre-Existing Failures Left Unaddressed Too Long

`TestGolangciDisabledLinters` was failing the entire first session and I reported work as "complete" without noticing. The status report from session 1 falsely claimed everything was done. Only caught in the self-review.

### 4. First Status Report Was Inaccurate

The `docs/status/2026-07-30_01-55_snapshot-test-infrastructure-overhaul.md` report claimed:
- "All tests pass" — FALSE, `TestGolangciDisabledLinters` was failing
- "Framework improvements complete" — technically true but couldn't be committed due to corrupt index
- "Lint clean" — only checked `internal/golden/...`, not the new sweep test files (which had a gci issue)

---

## E) WHAT WE SHOULD IMPROVE

### Process

1. **Run `git status` first.** Always. Before any work. The corrupt index cost an entire session of uncommitted work.
2. **Run `go test ./...` before claiming done.** Not just the package you touched. The full suite. `TestGolangciDisabledLinters` would have been caught.
3. **Commit incrementally.** Don't let 40+ files accumulate uncommitted. Commit after each logical unit.
4. **Don't trust BuildFlow commit messages.** They are hallucinated and misleading. Always commit manually with proper messages.
5. **Verify lint on ALL new files,** not just the framework package.

### Architecture

6. **Golden coverage drift guard.** A test that asserts every exported component function has at least one golden test. Without this, coverage gaps are invisible.
7. **`htmx` and `layout` are testing deserts.** These are foundational packages (layout renders the page shell, htmx handles all HTMX integration). Zero golden coverage is a liability.
8. **The three-tier strategy needs enforcement.** Documentation says "golden first, substring for invariants, visual for layout." But there's no mechanism preventing new components from shipping with zero tests.
9. **Consider `go-cmp` for diff output.** The custom LCS diff is fine for small files, but `go-cmp` produces richer, more readable diffs with path information. It's already a transitive dependency of many Go testing tools. (Counter-argument: keeps golden.go dependency-free for the published module.)
10. **Fuzz tests for normalizations.** `normalizeIDs` and `normalizeClasses` handle arbitrary HTML strings. Fuzz testing would catch regex edge cases (malformed attributes, nested quotes, etc.).

### Quality

11. **`.fail/` artifacts should be gitignored or auto-cleaned.** Committing failure screenshots is noise.
12. **The `visualtest/go.mod` replace version was stale** (v1.3.0 vs v1.3.1). BuildFlow auto-fixed this but it shouldn't have been stale.
13. **Dark-mode golden tests need a systematic approach.** Currently just `display/dark_golden_test.go` covering Badge, Button, Card. Should be every component × dark mode.

---

## F) 50 THINGS TO DO NEXT

### Tier 1: Quick Wins (Low Effort, High Impact)

1. Add golden tests for `htmx` package (8 components — small, self-contained)
2. Add golden tests for remaining `display` components (SimpleCard, StatusBadge, SimpleEmptyState)
3. Add golden tests for remaining `forms` components (Label, FieldError, FormFieldWrapper, Radio)
4. Add a golden coverage drift guard test (contract test: "every component has ≥1 golden")
5. Clean up `.fail/` artifacts in visualtest and add to `.gitignore`
6. Fix the 2 failing visual tests or mark with `t.Skip` + TODO comment
7. Migrate `display/golden_test.go` and `golden_new_test.go` to `AssertSnapshots` pattern
8. Add fuzz test for `normalizeIDs` (arbitrary ID-like strings)
9. Add fuzz test for `normalizeClasses` (malformed class attributes)
10. Add fuzz test for `diff` (arbitrary strings should not crash or hang)

### Tier 2: Medium Effort, High Impact

11. Add golden tests for `layout` package (AppShell, Split, Stack, Container, Base, Minimal)
12. Add golden tests for `recipes` package (Dashboard, SettingsLayout, LoginCard)
13. Create `AssertScreenshotSnapshots` table-driven helper for visual tests
14. Add dark-mode golden variants for all feedback components
15. Add dark-mode golden variants for all forms components
16. Add dark-mode golden variants for all navigation components
17. Add visual tests for feedback components (Spinner, ProgressBar, Skeleton, StepIndicator)
18. Add visual tests for forms components (Checkbox, Toggle, RadioGroup, Combobox)
19. Add visual tests for navigation components (Pagination, Breadcrumbs, SidebarNav)
20. Investigate and fix the modal/drawer visual test failures (likely anti-aliasing)

### Tier 3: Medium Effort, Medium Impact

21. Add golden tests for enum variant matrices (BadgeType × BadgeSize, AlertType × Dismissible)
22. Add golden tests for container-aware variants (`ContainerAware: true`)
23. Add golden tests for edge cases (empty props, nil children, error states)
24. Add RTL golden variants for components with directional classes
25. Add responsive viewport visual tests (Nav hamburger, Grid, Table)
26. Add interaction-state visual tests (Tab switch, Accordion expand/collapse)
27. Audit substring tests for redundancy with golden tests
28. Remove redundant substring tests where golden is strictly more comprehensive
29. Add a `TESTING.md` doc with examples and decision tree
30. Benchmark normalization + diff pipeline on largest golden file

### Tier 4: Polish & Future-Proofing

31. Add visual tests for layout components (AppShell, Container, Split, Stack)
32. Add dark-mode visual variants for Popover and ContextMenu
33. Add RTL visual variants for more components
34. Add RTL + dark combo visual tests
35. Consider HTML structural normalization (attribute ordering via `x/net/html`)
36. Add `-update --dry-run` mode (show what would change without writing)
37. Add golden file count metric to status reports
38. Consider snapshot testing for generated CSS
39. Add visual test for Carousel slide navigation
40. Add visual test for Tooltip positioning (all 4 positions)
41. Add visual test for Dropdown open state (all positions)
42. Add visual test for Table sortable headers (asc/desc indicators)
43. Add golden tests for `htmx` interaction states (loading, error, retry)
44. Add golden tests for `errorpage` handler JSON mode
45. Add golden tests for `errorpage` WriteError/WriteErrorPage wrappers

### Tier 5: Process & Tooling

46. Fix BuildFlow daemon to generate commit messages from `git diff --stat` (requires `larsartmann/buildflow` change)
47. Add pre-push hook that runs `go test ./...` (not just pre-commit which has 60s budget)
48. Add CI check for golden file count trend (alert on sudden drops)
49. Document the normalization pipeline design decisions in `docs/adr/`
50. Consider a `make golden-update` convenience target

---

## G) QUESTIONS I CANNOT ANSWER MYSELF

### 1. Should we invest in visual regression expansion now, or is golden coverage the priority?

Golden tests are fast (<1ms each) and catch structural/attribute regressions. Visual tests are slow (~2s each) but catch layout/color/dark-mode issues. Expanding visual from 12% to 50%+ coverage means ~150 new screenshots and significant chromedp maintenance. **Is the visual test ROI there for a server-rendered component library where consumers style via Tailwind overrides, or should we focus on getting golden coverage to 100% first and keep visual tests targeted at high-risk components (overlays, responsive layouts)?**

### 2. The BuildFlow commit messages on `origin/master` are bad. Rewrite history or live with it?

Commits like `"Looking at the changes in the internal/golden/ package, I need to generate a detailed commit message..."` are permanently pushed. The AGENTS.md says "NEVER force push." But these are actively misleading. **Should we do a one-time history cleanup (interactive rebase + force-with-lease), or accept these as a known BuildFlow limitation and move forward?**

### 3. Should the `htmx` and `layout` packages get golden tests, or are they different enough to need a different testing approach?

`layout` components (Base, Minimal) render full HTML documents with `<head>`, CDN scripts, and embedded CSS — much larger output than typical components. `htmx` components often render only small fragments (hidden inputs, script tags). **Should these use the same `golden.AssertSnapshots` pattern, or do they need specialized handling (e.g., partial golden comparison for large HTML documents, or focusing on specific attributes for tiny fragments)?**
