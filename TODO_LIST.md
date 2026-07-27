# TODO List — templ-components

**Updated:** 2026-07-27 | **Version:** 1.2.0

> Only open, actionable items. Completed work is tracked in CHANGELOG.md.
> Statuses: ⬜ deferred, ⚫ blocked (needs external resources).

---

## Open — actionable

| #   | Task                                                                   | Notes                                                                                                                  |
| --- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| 70  | Set `GOWORK=off` in `flake.nix` devShell `shellHook`                   | Parent `/home/lars/projects/go.work` breaks `go generate ./...` and BuildFlow pre-commit in this repo. Recurring across 3 sessions. |
| 71  | Investigate GitHub Dependabot alert                                    | "1 moderate vulnerability" reported but never investigated. Check `https://github.com/LarsArtmann/templ-components/security/dependabot`. |
| 72  | Add demo CSS rebuild to `scripts/release.sh` (or document Docker handles it) | Release script runs `templ generate` but not the Tailwind CSS build. `examples/demo/static/app.css` may be stale after component class changes. Docker 3-stage build overwrites it, but local `go run` does not. |
| 73  | Convert assertion-based tests to golden files (navigation, feedback, forms) | Recurring from Pareto plan. Per-package, bounded. Golden files improve diff readability and reduce brittle substring checks. Start with `navigation` (highest value). |

---

## Blocked — External dependencies

| #   | Task                                               | Blocker                                      |
| --- | -------------------------------------------------- | -------------------------------------------- |
| 28  | `awesome-templ` PR submission                      | Needs maintainer approval                    |
| 29  | `templ.guide` listing submission                   | Needs maintainer approval                    |
| 13  | Visual regression testing (Playwright screenshots) | Requires npm/playwright — no Node.js in repo |

---

## Deferred — v2.0 breaking changes

| #   | Task                                                            | Notes                                                                                                                                            |
| --- | --------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| 35  | Flip defaults: self-host htmx + semantic tokens → default       | Both shipped opt-in (v0.22.0). `HTMXSrc` opt-in; `templ-components-theme.css` opt-in. Default flip deferred to v2.0 (insufficient deprecation time). See ADR-0007, ADR-0008. |
| 38  | Remove `AlertType` / `ToastType` backward-compat aliases        | Other aliases (`ModalSizeFull`, `DrawerFull`, `FamilyFromErrorFamily`, `FormProps.Inline`) removed in v1.0.0. These two remain as `type X = FeedbackType` aliases. |
| 39  | Compound component pattern (Trigger/Content/Close) for overlays | Current Modal/Drawer are monolithic. v2.0 design.                                                                                                 |

---

## Deferred — v1.0 follow-up

| #   | Task                                                      | Notes                                                                            |
| --- | --------------------------------------------------------- | -------------------------------------------------------------------------------- |
| 33  | `Validate() error` methods on remaining props structs     | `ErrorPageProps.Validate()` shipped v1.0.0. Other props use graceful `utils.Lookup` fallback — add `Validate` only where invalid states are representable. |
| 34  | Move test helpers to `internal/testutil/`                 | 70+ test files depend on exported helpers. Large mechanical migration, deferred post-v1.0. |

---

## Deferred — Tooling

| #   | Task                                              | Notes                                                                                                            |
| --- | ------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| 67  | Switch treefmt `gofmt` → `gofumpt` in `flake.nix` | Latent conflict with `.golangci.yml` `gofumpt` linter; deferred to avoid formatting churn across entire codebase |
