# Visual Testing — Comprehensive Status Report

**Date:** 2026-07-27 20:38
**Session scope:** Adding pixel-level visual regression testing to the templ-components library
**Author:** Crush (assisted session)

---

## TL;DR

Built a **chromedp + pixelmatch** visual regression framework in a **separate
Go module** (`visualtest/`), wired into Nix (`nix run .#visual`) and CI.
15 golden PNGs across Button/Alert/Badge/Card covering light/dark/hover/focus/
disabled/mobile. All green. But honest critique below reveals real gaps.

---

## a) FULLY DONE ✅

These work, are verified, and committed (via auto-commit daemon):

| Item                                                                                        | Evidence                                                                                        |
| ------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- |
| Separate module (`visualtest/go.mod`) with `replace` directive                              | Consumers never pull chromedp — confirmed `go list ./...` from root shows 0 visualtest packages |
| chromedp browser harness                                                                    | `newBrowser()` per-test Chromium process, skip-if-absent, 1-2s startup                          |
| pixelmatch perceptual diff (YIQ + anti-alias skip)                                          | `compare.go` — replaced naive per-channel diff from first pass                                  |
| Compare-layer unit tests (4 cases)                                                          | `compare_test.go` — identical/different/dimension-mismatch/sub-threshold                        |
| `AssertScreenshot` public API with `Options{Dark,RTL,Viewport,MaxMismatch,Threshold,State}` | `harness.go`                                                                                    |
| `-update` flag (mirrors `internal/golden` DX)                                               | `golden.go`                                                                                     |
| `.fail/` failure artifacts (actual + diff PNGs) + auto-cleanup on pass                      | `golden.go` `cleanFailureArtifacts`                                                             |
| HTTP-server page serving (data: URL couldn't hold 82KB CSS)                                 | `capture()` in `harness.go`                                                                     |
| CSS loader (reads `examples/demo/static/app.css` via runtime.Caller)                        | `css.go`                                                                                        |
| Isolated render shell (`#tc-root` wrapper, dark/rtl/dir attrs)                              | `render.go`                                                                                     |
| `StateHover` (mouse-move via bounding rect JS)                                              | `hoverAction()`                                                                                 |
| `StateFocus` (focus first focusable descendant — `#tc-root` div isn't focusable)            | `focusAction()`                                                                                 |
| Nix flake `visual` app (Chromium in runtimeInputs, GOWORK=off, arg passthrough)             | `flake.nix`                                                                                     |
| CI workflow `visual` job (Nix-based for renderer bit-parity, failure artifact upload)       | `.github/workflows/ci.yaml`                                                                     |
| 15 golden PNGs: button (7), alert (3), badge (2), card (3)                                  | `visualtest/testdata/`                                                                          |
| `docs/visual-testing.md` (how-to, options table, failure flow, why-separate-module)         | committed                                                                                       |
| `AGENTS.md` updated (module table, flake commands, conventions)                             | committed                                                                                       |
| `visualtest/.gitignore` excludes `.fail/`                                                   | committed                                                                                       |
| `nix flake check` passes, `nix fmt` clean, `go vet` clean, `gofmt` clean                    | verified                                                                                        |

---

## b) PARTIALLY DONE 🟠

| Item                           | What exists                                                                 | What's missing                                                                                              |
| ------------------------------ | --------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| Component coverage             | 4 of 15+ packages (Button, Alert, Badge, Card)                              | 11 packages untested: forms (21 comps), navigation (12), feedback (10 more), htmx, errorpage, recipes, etc. |
| State coverage                 | Hover, Focus, Disabled on Button only                                       | No hover/focus tests on Alert, Badge, Card, inputs                                                          |
| Dark-mode coverage             | 6 of 15 goldens are `_dark`                                                 | No dark variants for alerts beyond error, badges beyond error_pill                                          |
| Mobile/responsive              | 1 test (card/mobile @ 375px)                                                | No tests at 768px (tablet), 1024px (desktop) breakpoints                                                    |
| RTL coverage                   | `RTL` option **exists and works** but **zero tests use it**                 | Completely untested axis — logical-property mirroring is unverified                                         |
| Failure-detection verification | Manually corrupted 1 golden, confirmed 97.9% mismatch caught + diff emitted | No automated regression test of the harness itself                                                          |

---

## c) NOT STARTED ❌

- **Interactive components** (Modal, Drawer, Dropdown, Combobox, Accordion, Tabs, Tooltip, Popover, ContextMenu) — none have visual tests. These are the highest-regression-risk components because they depend on JS + CSS animations + top-layer positioning.
- **Forms** (21 components) — Input, Select, Textarea, Toggle, Slider, Rating, TagsInput, Calendar, ValidationSummary, FieldError, etc. — zero visual coverage.
- **HTMX loading/error states** — InlineLoadingOverlay, SwapOOB, ViewTransitions — zero visual coverage.
- **Error pages** — ErrorPage, NotFound404, ErrorDetail, ErrorAlert — zero visual coverage.
- **Recipes** (Dashboard, SettingsLayout, LoginCard) — zero visual coverage of the composed screens.
- **Table** (sortable headers, clickable rows, content-visibility) — zero visual coverage.
- **Animation capture** — `@starting-style`, `allow-discrete`, toast enter/exit, spinner, skeletons are animated; screenshots capture one frame. No strategy for animation golden testing.
- **Cross-browser** — only Chromium. No Firefox/WebKit (would need playwright or similar — violates no-Node rule, so likely out of scope, but untested assumption).
- **Font determinism** — tests rely on Chromium's bundled fonts; no explicit font-loading or `-font-render-hinting` verification beyond the flag I set. Real CI may differ from local Nix if nixpkgs revs diverge.
- **Golden drift detection** — no test that flags when goldens get stale (e.g., component changed but golden not regenerated).
- **Visual coverage report** — no metric for "% of components with visual tests" (unlike `TestSkillComponentCount` for unit tests).

---

## d) TOTALLY FUCKED UP 💥

### 1. I overwrote the root `go.mod` in the first 5 minutes

**What happened:** Ran `cat > go.mod` from the repo root instead of `visualtest/`. The `bash` tool's `working_dir` defaulted to CWD. Wiped the root module declaration.
**Impact:** Would have broken the entire library build if committed.
**Fix:** `git restore go.mod` immediately. Lesson: **always use absolute paths for file writes in shell heredocs**.

### 2. The first pixel diff implementation was naive garbage

**What happened:** Raw per-channel `uint8` tolerance comparison. Would have produced false positives on every Chromium version bump due to font anti-aliasing shifts.
**Fix:** Replaced with `orisano/pixelmatch` (YIQ perceptual distance + AA skip) — the same lib chromedp's own test suite uses. But this was caught in self-critique, not during design.

### 3. `data:` URL approach was fundamentally broken

**What happened:** First `capture()` tried `chromedp.Navigate("data:text/html," + page)` with 82KB of inlined CSS. All tests failed with "context canceled" — data URLs have size limits.
**Wasted time:** ~10 minutes debugging before realizing the obvious.
**Fix:** Switched to `httptest.Server`. Should have started there.

### 4. `arguments` is not defined in arrow functions

**What happened:** `hoverAction` JS used `(() => { document.querySelector(arguments[0]) })()` — arrow functions don't bind `arguments`. Caught only at runtime.
**Fix:** Interpolated selector as string literal. Basic JS mistake.

### 5. Tried to `Focus()` a `<div>`

**What happened:** `chromedp.Focus("#tc-root")` — but `#tc-root` is a div, not focusable. "Element is not focusable (-32000)".
**Fix:** `focusAction()` descends to first focusable child. Should have thought about the DOM structure before writing the action.

### 6. Left dead code (`decodePNG`) with a `//nolint:unused` escape hatch

**What happened:** Wrote a "helper for future test introspection" that was never used. Suppressed the linter instead of deleting it.
**Fix:** Removed in the self-critique pass. The `//nolint:unused` pattern is a smell — if it's unused, delete it.

### 7. Lied in the doc comment (`-tags=visual`)

**What happened:** `doc.go` said `go test ./... -tags=visual` but no such build tag existed. Copy-pasted from a pattern that didn't apply.
**Fix:** Corrected to `go test ./...`. Documentation that lies is worse than no documentation.

### 8. The auto-commit daemon created **12+ noisy commits** for this session

**What happened:** The auto-git daemon committed every intermediate state, producing commits like `8f83db5 feat(visualtest): overhaul visual testing harness...` alongside `f4fd54b feat(visualtest): add visual regression testing framework` — many are near-duplicates of partial work.
**Impact:** Git history is noisy. The AGENTS.md says this is "expected behavior" but for a feature this size, a squash would be cleaner. Not my fault (daemon), but worth noting.

---

## e) WHAT WE SHOULD IMPROVE 🎯

### Architecture / Design

1. **The render shell (`#tc-root`) wraps every component in `display: inline-block; padding: 16px`.** This means screenshots include 16px of padding around every component — goldens are larger than the component itself. This is probably wrong for tight-cropped comparison and makes diffs include padding pixels. Should be configurable or removed.

2. **One Chromium process per test is wasteful.** Each test launches + kills Chromium (~1s overhead × 15 tests = 15s overhead). Should share one browser with multiple tabs (my first design tried this but had a context-cancellation bug I didn't debug — I gave up and went per-process). The per-process approach scales poorly: 100 visual tests = 100 Chromium launches.

3. **No CSS regeneration step in the test harness.** If a component adds new Tailwind classes, the tests render with stale CSS unless someone manually recompiles `examples/demo/static/app.css`. The harness should either (a) compile CSS on the fly or (b) fail loud if the CSS is older than the `.templ` files.

4. **`MaxMismatch` default of 0.1% may be too strict or too loose — unvalidated.** I picked it by intuition. No data on what real regressions look like (what % do they affect?) or what anti-aliasing noise looks like across nixpkgs revs. Should calibrate with a deliberate-breakage experiment.

5. **No "acceptable diff" review workflow.** When a golden changes, the only review tool is opening two PNGs side by side. Should integrate with a visual diff tool or produce an overlay/blink comparison HTML.

6. **`StateHover` moves mouse to center of `#tc-root` — but hover styles are on the button, not the wrapper.** This works only because the button fills the wrapper. For components where the interactive element is smaller than the wrapper, hover will miss. Should target the first interactive child (like `focusAction` does).

### Type Model

7. **`Options` uses bool fields (`Dark`, `RTL`) that OR together in `resolveOptions`.** This means you can't turn Dark _off_ if a prior option set it on. Should use a tri-state or pointer or `*bool`. Edge case, but violates "make impossible states unrepresentable."

8. **`InteractionState` is an `int` enum without a `String()` method.** Error messages and logs show numeric values. Should implement `String()` for debuggability.

9. **`Viewport` is a separate type but `Options` has no `ViewportPreset` helper.** Every test manually specifies `{Width: 375, Height: 667}`. Should have `ViewportMobile`, `ViewportTablet`, `ViewportDesktop` constants.

### Testing Gaps

10. **Zero tests for the harness itself** beyond the compare layer. No test that `AssertScreenshot` correctly skips when no browser, no test that `-update` actually writes, no test that failure artifacts are created/ cleaned.

11. **No flakiness testing.** The tests pass 3/3 times in my session, but I never ran them 100× to check for intermittent pixel drift. Visual tests are notoriously flaky; this needs a stability run.

12. **The `hoverAction` JS builds a string with `fmt.Sprintf("%q", sel)` — XSS-safe for our use, but no test proves it.** If `sel` ever came from user input (it won't — it's hardcoded), this would be an injection vector.

### DX / Documentation

13. **No `make visual` or shortcut beyond `nix run .#visual`.** Contributors without Nix are stuck manually setting 3 env vars (GOEXPERIMENT, GOWORK, CHROMEDP_CHROME_PATH). Should have a `scripts/visual.sh` or at least document the manual path more prominently.

14. **`docs/visual-testing.md` doesn't mention the `.fail/` cleanup-on-pass behavior.** Documented here but not in the doc.

15. **AGENTS.md says "tests skip if no browser" but doesn't warn that CI will FAIL if Nix can't provide Chromium.** Different failure modes for local vs CI.

---

## f) Up to 50 things to do next

Sorted by **impact / effort** (P0 = do first):

| #   | Priority | Task                                                                                                |
| --- | -------- | --------------------------------------------------------------------------------------------------- |
| 1   | **P0**   | Add visual tests for **Modal** and **Drawer** (native `<dialog>`, highest regression risk)          |
| 2   | **P0**   | Add visual tests for **Dropdown/Popover/ContextMenu** (Popover API top-layer positioning)           |
| 3   | **P0**   | Add visual tests for **Input** + **Select** (forms, most-used components)                           |
| 4   | **P0**   | Add at least **one RTL test** — the option exists but is 100% unused                                |
| 5   | **P0**   | Fix the shared-browser bug (one Chromium, many tabs) — 15s→2s test suite                            |
| 6   | **P0**   | Fix `StateHover` to target first interactive child, not wrapper center                              |
| 7   | **P1**   | Add visual tests for **Accordion** (`<details>` migration)                                          |
| 8   | **P1**   | Add visual tests for **Tabs** (RTL keyboard nav, active tab focus ring)                             |
| 9   | **P1**   | Add visual tests for **Tooltip** (pure CSS show/hide)                                               |
| 10  | **P1**   | Add visual tests for **Toast** (enter/exit animation frame)                                         |
| 11  | **P1**   | Add visual tests for **Spinner** + **Skeleton** (animation frame)                                   |
| 12  | **P1**   | Add visual tests for **ProgressBar** (clamped aria-valuenow visual)                                 |
| 13  | **P1**   | Add visual tests for **Table** (sortable headers, clickable rows, lazy rows)                        |
| 14  | **P1**   | Add visual tests for **Avatar** + **AvatarStatus** (dot rendering)                                  |
| 15  | **P1**   | Add visual tests for **Carousel** (scroll-snap positioning)                                         |
| 16  | **P1**   | Add visual tests for **Combobox** (dropdown open state)                                             |
| 17  | **P1**   | Add visual tests for **errorpage.ErrorPage** + **NotFound404**                                      |
| 18  | **P1**   | Add a **dark-mode variant for EVERY component** that has semantic colors                            |
| 19  | **P1**   | Calibrate `MaxMismatch` with a deliberate-breakage experiment (change a color, measure %)           |
| 20  | **P2**   | Add `ViewportMobile`/`ViewportTablet`/`ViewportDesktop` presets                                     |
| 21  | **P2**   | Add `String()` to `InteractionState`                                                                |
| 22  | **P2**   | Change `Options.Dark`/`RTL` from `bool` to `*bool` (tri-state) or use a mode enum                   |
| 23  | **P2**   | Remove or make configurable the `#tc-root` 16px padding in the render shell                         |
| 24  | **P2**   | Add visual tests for **Toggle** (peer-checked translate-x animation frame)                          |
| 25  | **P2**   | Add visual tests for **Checkbox** + **RadioGroup** (accent-color)                                   |
| 26  | **P2**   | Add visual tests for **Textarea** (auto-grow field-sizing)                                          |
| 27  | **P2**   | Add visual tests for **Slider** + **Rating**                                                        |
| 28  | **P2**   | Add visual tests for **ValidationSummary** + **FieldError**                                         |
| 29  | **P2**   | Add visual tests for **Breadcrumb** + **Pagination** + **SidebarNav**                               |
| 30  | **P2**   | Add visual tests for **NavLink** + **MobileNavLink** (active/inactive states)                       |
| 31  | **P2**   | Add visual tests for **Footer**                                                                     |
| 32  | **P2**   | Add visual tests for **CountBadge** (overflow "N+", zero-hide)                                      |
| 33  | **P2**   | Add visual tests for **RelativeTime**                                                               |
| 34  | **P2**   | Add visual tests for **EndOfList** + **LoadMore**                                                   |
| 35  | **P2**   | Add visual tests for **recipes.Dashboard** (full composed screen)                                   |
| 36  | **P2**   | Add visual tests for **recipes.SettingsLayout**                                                     |
| 37  | **P2**   | Add visual tests for **recipes.LoginCard**                                                          |
| 38  | **P2**   | Add **768px tablet** viewport tests for responsive components (Grid, Card ContainerAware)           |
| 39  | **P2**   | Add a **visual coverage metric** (like `TestSkillComponentCount`): "% of components with ≥1 golden" |
| 40  | **P2**   | Add a test that `-update` actually rewrites goldens (harness meta-test)                             |
| 41  | **P2**   | Add a test that `.fail/` artifacts are created on mismatch and cleaned on pass                      |
| 42  | **P3**   | Add **CSS-staleness detection**: fail if `app.css` mtime < newest `.templ` mtime                    |
| 43  | **P3**   | Add **animation capture** strategy (multiple frames? settled-state guarantee?)                      |
| 44  | **P3**   | Add **cross-Chromium-version flakiness run** (run 50×, measure variance)                            |
| 45  | **P3**   | Add **overlay/blink-comparison HTML** generator for PR review of golden changes                     |
| 46  | **P3**   | Document the manual (non-Nix) path more prominently in `docs/visual-testing.md`                     |
| 47  | **P3**   | Document `.fail/` cleanup-on-pass behavior in `docs/visual-testing.md`                              |
| 48  | **P3**   | Consider a `scripts/visual.sh` for non-Nix contributors                                             |
| 49  | **P3**   | Add **htmx loading states** visual tests (InlineLoadingOverlay, LoadingButton during request)       |
| 50  | **P3**   | Investigate **Firefox/WebKit** feasibility (likely out of scope per no-Node rule, but document why) |

---

## g) Questions I CANNOT figure out myself

### 1. Should visual tests gate CI (blocking) or be informational (non-blocking)?

**Context:** The CI job I added runs `nix run .#visual` and fails the build on mismatch. But visual tests are notoriously flaky across environments, and the goldens are generated against one nixpkgs revision. If CI runs a different nixpkgs (via `channel:nixos-unstable`), font rendering may drift and false-fail.
**The question:** Do you want visual tests to **block merges** (strict, catches real regressions, may false-fail on nixpkgs drift) or run as a **separate non-required check** (informational, uploads diff artifacts for manual review)?
**Why I can't decide:** This is a team-workflow preference, not a technical decision. I don't know your tolerance for false positives vs. your desire for hard gates.

### 2. What's the target golden coverage — every component, or a representative subset?

**Context:** The library has ~90 components. Visual-testing all of them (with light/dark/hover/focus variants) is 300-500 goldens — a large binary footprint in the repo and significant maintenance when visual changes are intentional.
**The question:** Do you want visual tests for **every component** (exhaustive, high maintenance) or a **representative subset** (e.g., one per component family, covering each styling pattern once)?
**Why I can't decide:** This trades coverage against maintenance burden — a product judgment, not an engineering one.

### 3. Is the `examples/demo/static/app.css` the right CSS source, or should the harness compile its own?

**Context:** The harness reads the demo's pre-compiled CSS. If someone adds Tailwind classes to a component but forgets to recompile the demo CSS, visual tests silently render with stale styles (false green — the test passes but doesn't reflect reality).
**The question:** Should the harness (a) keep using the demo CSS and **fail-loud if stale** (mtime check), (b) **compile its own CSS** on every run (adds ~1s + a tailwindcss dependency), or (c) keep as-is and trust manual recompilation?
**Why I can't decide:** Option (b) violates the "no Node.js runtime" rule unless we use the Nix `tailwindcss_4` binary — which is available but I don't know if you want test runs to invoke it.

---

## Session metrics

| Metric                          | Value                                 |
| ------------------------------- | ------------------------------------- |
| Harness Go files                | 8 files, 846 lines                    |
| Golden PNGs                     | 15                                    |
| Test functions                  | 6 visual + 4 compare-unit             |
| Commits this session            | ~12 (auto-daemon)                     |
| Packages with visual coverage   | 4 of 15 (display, feedback partially) |
| Components with visual coverage | ~8 of ~90                             |
| Test suite runtime              | ~2s (parallel)                        |
| Things fucked up                | 8 (see section d)                     |

---

## Resolution (2026-07-27, later session)

The `visualtest/` framework described in section a **shipped** in CHANGELOG `[Unreleased]` (Added) — separate Go module, chromedp + pixelmatch, `nix run .#visual`, CI job, 15 baseline goldens. The forward work from section f was routed:

| Forward item (this report, P0/P1) | Routed to |
| --- | --- |
| Visual tests for Modal/Drawer/Dropdown/Input/Select (highest regression risk) | TODO_LIST #75 |
| Share one Chromium process across tests (15s→2s) | TODO_LIST #76 |
| First RTL visual test (`Options.RTL` exists, 0 users) | TODO_LIST #77 |
| `StateHover` target first interactive child, not wrapper center (e.6) | Remains open (not yet a TODO) |
| `MaxMismatch` calibration, `Viewport*` presets, `InteractionState.String()`, `*bool` tri-state | Remains as ROADMAP directions |

Coverage stands at 4 of 15 packages (Button/Alert/Badge/Card). The 11 untested packages — especially the interactive overlays (Modal/Drawer/Dropdown/Combobox) that depend on JS + top-layer + animations — are the highest-value next work (TODO #75).
