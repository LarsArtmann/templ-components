# TODO List — templ-components

**Updated:** 2026-08-04 | **Version:** 1.6.0

> Only open, actionable items. Completed work is tracked in [`CHANGELOG.md`](CHANGELOG.md).
> Statuses: ⬜ deferred, ⚫ blocked (needs external resources).

> **2026-08-04 audit:** Items 67, 73, 79, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, and 94
> were all verified as **complete** against the code and removed from this list (their entries
> already live in `[Unreleased]` in CHANGELOG.md). #82's Chromium calibration is now done too
> (Chromium was available via `nix run .#visual` after all). The only remaining item requiring
> resources unavailable to an automated agent is human review of the overlay PNGs (#80).

---

## Open — actionable

_No actionable items. The actionable backlog is clear — all tasks are either complete
(see CHANGELOG `[Unreleased]`), blocked on external resources, or deferred to a future
major version below._

---

## Blocked — External dependencies

| #   | Task                                        | Blocker                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| --- | ------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 80  | Human-eyeball the AI-generated overlay PNGs | `dropdown` (light/dark), `popover`, `contextmenu`, `modal`, and `drawer` PNGs were captured by an automated agent. AI cannot read PNGs, so a human must confirm no enshrined rendering bug (e.g. wrong top-layer position). Run `nix run .#visual`, then inspect `visualtest/testdata/{dropdown,popover,contextmenu,modal,drawer}/`.                                                                                                                              |
| 28  | `awesome-templ` PR submission               | Needs upstream maintainer approval.                                                                                                                                                                                                                                                                                                                                                                                                                               |
| 29  | `templ.guide` listing submission            | Needs upstream maintainer approval.                                                                                                                                                                                                                                                                                                                                                                                                                               |
| 93  | BuildFlow daemon: honest commit messages    | 6+ sessions. Daemon commits with hallucinated messages (e.g. "chore: update project configuration") authored as "Unknown Author", and re-introduces stale files via broad `git add -A`. Fix lives in `larsartmann/buildflow` (pre-commit `golangci-lint config verify`, message templates derived from `git diff --stat`, `GOWORK=off`). Blocked on separate-repo work. Mitigated in this repo by `scripts/check-lint-config.sh` + `TestGolangciDisabledLinters`. |

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
