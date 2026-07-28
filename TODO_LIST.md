# TODO List — templ-components

**Updated:** 2026-07-28 | **Version:** 1.2.0

> Only open, actionable items. Completed work is tracked in [`CHANGELOG.md`](CHANGELOG.md).
> Statuses: ⬜ deferred, ⚫ blocked (needs external resources).

---

## Open — actionable

### Test infrastructure

| #   | Task                                                                 | Notes                                                                                                                                                                                                                                                                                                                                                                                                                 |
| --- | -------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 73  | Finish converting assertion-based tests to golden files              | `display`/`navigation`/`feedback`/`forms`/`layout` now have golden snapshots; remaining assertion-heavy tests (esp. `htmx`, per-component edge cases) still use brittle substring checks. Golden files improve diff readability and reduce brittle `strings.Contains` checks.                                                                                                                                         |
| 79  | Visual regression coverage: expand to high-risk components           | 11 component types have goldens (alert/badge/button/card/contextmenu/drawer/dropdown/input/modal/popover/select = 31 PNGs). Add: Combobox, Tabs, Table, Accordion, Tooltip, Carousel, CopyButton, Spinner, Skeleton, ProgressBar, Toast, Avatar, StepIndicator, ErrorPage, NotFound404. Highest regression risk: interactive (Combobox, Tabs, Accordion) + animated (Carousel, ProgressBar).                          |
| 80  | Human-eyeball the 4 AI-generated overlay goldens                     | `dropdown` (light/dark), `popover`, `contextmenu` PNGs were captured by the agent but never verified by a human. AI cannot read PNGs — confirm no enshrined rendering bug (e.g. wrong top-layer position). `nix run .#visual` then inspect `visualtest/testdata/{dropdown,popover,contextmenu}/`.                                                                                                                     |
| 81  | Audit repo-wide for ordered-Tailwind-substring test assertions       | `utils.Class()` wraps tailwind-merge-go which reorders classes **nondeterministically** (depends on LRU cache state). Any `strings.Contains(out, "flex flex-col")` is a latent flake (~13% under `-race`). Convert each to `utils.AssertContainsAll`. Add `TestNoOrderedTailwindSubstringsInTests` drift-guard. Root cause: `integration/appshell_composition_test.go` was the first found + fixed; others unscanned. |
| 82  | Calibrate `MaxMismatch` for overlay visual tests empirically         | Current overlay `MaxMismatch` (0.02) is a guess. Run each overlay test 10×, set threshold at p99 of observed mismatch. Prevents both false-negatives (real regressions masked) and false-positives (font AA variance).                                                                                                                                                                                                |
| 83  | Fix `StateHover` to target first interactive child                   | `StateHover` moves the mouse over `#tc-root` center. Components whose hover styles target an inner element (button/link) may not have those styles applied during capture. Descend to the first interactive descendant like `StateFocus` already does (`harness.go:214`).                                                                                                                                             |
| 84  | Visualtest API: tri-state optionals + viewport presets + state names | `Options.Dark`/`RTL` are `bool` — zero-value conflates "false" with "unset". Change to `*bool` for tri-state. Add `ViewportMobile`/`ViewportTablet`/`ViewportDesktop` presets. Add `InteractionState.String()` for readable failure messages.                                                                                                                                                                         |
| 85  | Pin Chromium version in `flake.nix`                                  | Visual tests rely on the Nix-provided Chromium; an un-pinned rolling version can shift pixel output across CI runs (font/AA changes). Pin a specific version for reproducibility.                                                                                                                                                                                                                                     |
| 94  | Fix 2 latent visual test failures exposed by `:=` shadowing fix      | After fixing the `visualtest/doc.go` `:=`→`=` shadowing bug (2026-07-28), the visualtest suite actually runs (previously silently passed because the broken `:=` left `sharedAllocCtx` nil and tests skipped). Two tests now fail: `drawer/right_light` (100% mismatch, 32x32 blank — dialog element not found), `modal/open_light` (9.86% mismatch). Both `_dark` variants pass. Likely cause: `<dialog Open=true>` doesn't promote to top-layer without `showModal()` JS — the dialog renders in-flow but isn't visible. The 10:14 report's "31 goldens green" claim was vacuously true. |

### Components & recipes

| #   | Task                                             | Notes                                                                                                                                                                                                                                                                                          |
| --- | ------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 86  | Add popover edge-flipping to `popoverPositionJS` | The positioner (`display/shared.go`) handles 4 positions but does **not flip** when the preferred side clips the viewport. Mirror the proven ContextMenu clamping pattern: detect clip, flip to opposite side. SSR-verified only — no browser harness exists yet to confirm positioning fixes. |
| 87  | Add `recipes.AuthLayout` + `recipes.EmptyState`  | Composition screens companion to the existing Dashboard/SettingsLayout/LoginCard (ADR-0019). AuthLayout: centered card + side-panel split. EmptyState: icon+title+action composition.                                                                                                          |

### Tooling

| #   | Task                                               | Notes                                                                                                                                                                                   |
| --- | -------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 88  | Add `nix run .#css` app                            | Recompiles `examples/demo/static/app.css` via `tailwindcss --minify`. Currently only Docker and `release.sh` rebuild it; local `go run` serves stale CSS after component class changes. |
| 89  | `tc` CLI: `tc version` + `tc add --list-deps` flag | `tc version` prints `utils.Version`. `tc add --list-deps <component>` lists sibling `.go` files the component depends on (the silent-incompleteness gap warned about in v1.2.0).        |
| 67  | Switch treefmt `gofmt` → `gofumpt` in `flake.nix`  | Latent conflict with `.golangci.yml` `gofumpt` linter; deferred to avoid formatting churn across the entire codebase.                                                                   |

### Documentation

| #   | Task                                                              | Notes                                                                                                                                                                                                                          |
| --- | ----------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| 90  | Write `docs/migration/skeletoncardgrid-api-change.md`             | The `[Unreleased]` breaking change (`SkeletonCardGrid(count)` → `SkeletonCardGrid(SkeletonCardGridProps{...})`) has no migration doc. Consumers upgrading will hit a compile error with only the CHANGELOG note to guide them. |
| 91  | Add "Testing" section to README + `docs/testing-guide.md`         | README has zero mention of the test strategy (golden HTML snapshots, pixel-level visual regression, drift-guard scanners). A reader evaluating the library cannot tell how it is verified.                                     |
| 92  | Fix unused `boolPtr` in `internal/golden/golden_coverage_test.go` | Dead helper flagged by linters. Trivial removal or wire-up.                                                                                                                                                                    |

---

## Blocked — External dependencies

| #   | Task                                     | Blocker                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| --- | ---------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 28  | `awesome-templ` PR submission            | Needs upstream maintainer approval.                                                                                                                                                                                                                                                                                                                                                                                                                               |
| 29  | `templ.guide` listing submission         | Needs upstream maintainer approval.                                                                                                                                                                                                                                                                                                                                                                                                                               |
| 93  | BuildFlow daemon: honest commit messages | 6+ sessions. Daemon commits with hallucinated messages (e.g. "chore: update project configuration") authored as "Unknown Author", and re-introduces stale files via broad `git add -A`. Fix lives in `larsartmann/buildflow` (pre-commit `golangci-lint config verify`, message templates derived from `git diff --stat`, `GOWORK=off`). Blocked on separate-repo work. Mitigated in this repo by `scripts/check-lint-config.sh` + `TestGolangciDisabledLinters`. |

---

## Deferred — v2.0 breaking changes

| #   | Task                                                            | Notes                                                                                                                                                                                  |
| --- | --------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 35  | Flip defaults: self-host htmx + semantic tokens → default       | Both shipped opt-in (v0.22.0). `HTMXSrc` opt-in; `templ-components-theme.css` opt-in. Default flip deferred to v2.0 (insufficient deprecation time). See ADR-0007, ADR-0008, ADR-0022. |
| 38  | Remove `AlertType` / `ToastType` backward-compat aliases        | Other aliases (`ModalSizeFull`, `DrawerFull`, `FamilyFromErrorFamily`, `FormProps.Inline`) removed in v1.0.0. These two remain as `type X = FeedbackType` aliases.                     |
| 39  | Compound component pattern (Trigger/Content/Close) for overlays | Current Modal/Drawer are monolithic. v2.0 design — ADR-0023 written.                                                                                                                   |

---

## Deferred — v1.0 follow-up

| #   | Task                                                  | Notes                                                                                                                                                      |
| --- | ----------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 33  | `Validate() error` methods on remaining props structs | `ErrorPageProps.Validate()` shipped v1.0.0. Other props use graceful `utils.Lookup` fallback — add `Validate` only where invalid states are representable. |
| 34  | Move test helpers to `internal/testutil/`             | 70+ test files depend on exported helpers. Large mechanical migration, deferred post-v1.0.                                                                 |
