# TODO List — templ-components

**Updated:** 2026-08-09 | **Version:** 1.8.0

> Only open, actionable items. Completed work is tracked in [`CHANGELOG.md`](CHANGELOG.md).
> Statuses: ⬜ deferred, ⚫ blocked (needs external resources).

---

## Blocked — External dependencies

| #   | Task                                        | Blocker                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| --- | ------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 80  | Human-eyeball the AI-generated overlay PNGs | `dropdown` (light/dark), `popover`, `contextmenu`, `modal`, and `drawer` PNGs were captured by an automated agent. AI cannot read PNGs, so a human must confirm no enshrined rendering bug (e.g. wrong top-layer position). Run `nix run .#visual`, then inspect `visualtest/testdata/{dropdown,popover,contextmenu,modal,drawer}/`. The regenerated `dropdown/open_dark.png` (6781→4790 bytes) makes this slightly more urgent.                                  |
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
