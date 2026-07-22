# TODO List — templ-components

**Updated:** 2026-07-22 | **Version:** 1.1.0

> Only open, actionable items. Completed work is tracked in CHANGELOG.md.
> Statuses: ⬜ deferred, ⚫ blocked (needs external resources).

---

## Blocked — External dependencies

| #   | Task                                               | Blocker                                      |
| --- | -------------------------------------------------- | -------------------------------------------- |
| 28  | `awesome-templ` PR submission                      | Needs maintainer approval                    |
| 29  | `templ.guide` listing submission                   | Needs maintainer approval                    |
| 30  | Configure SSH tag signing (`allowedSignersFile`)   | Requires user's local git config + SSH key   |
| 13  | Visual regression testing (Playwright screenshots) | Requires npm/playwright — no Node.js in repo |

---

## Deferred — v1.0 breaking changes

| #   | Task                                                      | Notes                                                                            |
| --- | --------------------------------------------------------- | -------------------------------------------------------------------------------- |
| 33  | `Validate() error` methods on all props structs           | `ErrorPageProps.Validate()` already shipped; design decision needed for the rest |
| 34  | Move test helpers to `internal/testutil/`                 | 70+ test files depend on exported helpers                                        |
| 35  | Self-host htmx as default, CDN opt-in                     | ADR 0007 written; `layout/sri.go` CDN still default                              |
| 36  | Semantic token layer (`bg-tc-primary` etc.)               | ADR 0008 written; 256 hardcoded color refs remain                                |
| 38  | Remove deprecated type aliases (`AlertType`, `ToastType`) | `FamilyFromErrorFamily` already removed                                          |

---

## Deferred — v2.0 architectural

| #   | Task                                                            | Notes                               |
| --- | --------------------------------------------------------------- | ----------------------------------- |
| 39  | Compound component pattern (Trigger/Content/Close) for overlays | Current Modal/Drawer are monolithic |

---

## Deferred — Tooling

| #   | Task                                              | Notes                                                                                                            |
| --- | ------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| 67  | Switch treefmt `gofmt` → `gofumpt` in `flake.nix` | Latent conflict with `.golangci.yml` `gofumpt` linter; deferred to avoid formatting churn across entire codebase |
