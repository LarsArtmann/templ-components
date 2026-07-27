# TODO List — templ-components

**Updated:** 2026-07-27 | **Version:** 1.2.0

> Only open, actionable items. Completed work is tracked in CHANGELOG.md.
> Statuses: ⬜ deferred, ⚫ blocked (needs external resources).

---

## Open — actionable

| #   | Task                                                                         | Notes                                                                                                                                                                                                                                                           |
| --- | ---------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 70  | Set `GOWORK=off` in `flake.nix` devShell `shellHook`                         | Parent `/home/lars/projects/go.work` breaks `go generate ./...` and BuildFlow pre-commit in this repo. Recurring across 3 sessions.                                                                                                                             |
| 71  | Investigate GitHub Dependabot alert                                          | "1 moderate vulnerability" reported but never investigated. Check `https://github.com/LarsArtmann/templ-components/security/dependabot`.                                                                                                                        |
| 72  | Add demo CSS rebuild to `scripts/release.sh` (or document Docker handles it) | Release script runs `templ generate` but not the Tailwind CSS build. `examples/demo/static/app.css` may be stale after component class changes. Docker 3-stage build overwrites it, but local `go run` does not.                                                |
| 73  | Convert assertion-based tests to golden files (navigation, feedback, forms)  | Recurring from Pareto plan. Per-package, bounded. Golden files improve diff readability and reduce brittle substring checks. Start with `navigation` (highest value).                                                                                           |
| 74  | Add `utils.TestContainerQueryCompliance` scanner                             | Scan all `.templ` for structural viewport breakpoints (`sm:`/`md:`/`lg:` on grid/flex/hidden/col-span) with no corresponding `ContainerAware` flag. Mirrors `TestDarkModeCompliance`/`TestMotionReduceCompliance`. Source: container-query report (2026-07-27). |
| 75  | Visual regression tests for highest-risk components                          | Cover Modal, Drawer, Dropdown, Input, Select. Currently 4/15 packages have goldens (Button/Alert/Badge/Card). Interactive overlays are the highest regression risk (JS + top-layer + animations). Source: visual-testing report (2026-07-27).                   |
| 76  | Share one Chromium process across visual tests (tabs, not processes)         | Current harness launches a browser per test (~1s startup × N). Scales poorly: 100 tests = 100 launches. First design tried this but hit a context-cancellation bug. Source: visual-testing report e.2.                                                          |
| 77  | Add the first RTL visual test                                                | `Options.RTL` exists and works but has zero users — logical-property mirroring (`ms-`/`me-`/`start-`/`end-`) is completely unverified at the pixel level. Tiny, high-value. Source: visual-testing report b/P0#4.                                               |
| 78  | Lint test: Tailwind-class lookup maps must live in `.templ` files            | Tailwind v4's scanner only reads `*.templ`, not `*.go`. A `map[X]string` of class strings in a `.go` file produces silently-missing CSS (caught this session: Split's container classes compiled to nothing). Source: container-query report d/e.1.             |

---

## Blocked — External dependencies

| #   | Task                             | Blocker                   |
| --- | -------------------------------- | ------------------------- |
| 28  | `awesome-templ` PR submission    | Needs maintainer approval |
| 29  | `templ.guide` listing submission | Needs maintainer approval |

---

## Deferred — v2.0 breaking changes

| #   | Task                                                            | Notes                                                                                                                                                                        |
| --- | --------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 35  | Flip defaults: self-host htmx + semantic tokens → default       | Both shipped opt-in (v0.22.0). `HTMXSrc` opt-in; `templ-components-theme.css` opt-in. Default flip deferred to v2.0 (insufficient deprecation time). See ADR-0007, ADR-0008. |
| 38  | Remove `AlertType` / `ToastType` backward-compat aliases        | Other aliases (`ModalSizeFull`, `DrawerFull`, `FamilyFromErrorFamily`, `FormProps.Inline`) removed in v1.0.0. These two remain as `type X = FeedbackType` aliases.           |
| 39  | Compound component pattern (Trigger/Content/Close) for overlays | Current Modal/Drawer are monolithic. v2.0 design.                                                                                                                            |

---

## Deferred — v1.0 follow-up

| #   | Task                                                  | Notes                                                                                                                                                      |
| --- | ----------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 33  | `Validate() error` methods on remaining props structs | `ErrorPageProps.Validate()` shipped v1.0.0. Other props use graceful `utils.Lookup` fallback — add `Validate` only where invalid states are representable. |
| 34  | Move test helpers to `internal/testutil/`             | 70+ test files depend on exported helpers. Large mechanical migration, deferred post-v1.0.                                                                 |

---

## Deferred — Tooling

| #   | Task                                              | Notes                                                                                                            |
| --- | ------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| 67  | Switch treefmt `gofmt` → `gofumpt` in `flake.nix` | Latent conflict with `.golangci.yml` `gofumpt` linter; deferred to avoid formatting churn across entire codebase |
