# TODO List — templ-components

**Updated:** 2026-08-05 | **Version:** 1.7.0

> Only open, actionable items. Completed work is tracked in [`CHANGELOG.md`](CHANGELOG.md).
> Statuses: ⬜ deferred, ⚫ blocked (needs external resources).

---

## Open — actionable

### Testing gaps (highest impact)

| #   | Task                                                      | Evidence                                                                                                                                                                                                       |
| --- | --------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 95  | Add visual regression tests for chart components          | No `visualtest/` entries for LineChart, PieChart, AreaChart (`visualtest/testdata/` — verified absent). Biggest testing gap per 3 status reports. Golden tests verify HTML structure but not visual rendering. |
| 96  | Add dark-mode visual variants for 6 newer components      | Combobox, Tooltip, Carousel, Skeleton, ErrorPage, NotFound404 are light-only in `visualtest/`. Every pre-existing component family has both light + dark goldens.                                              |
| 97  | Add visual tests for v1.5–v1.6 components without goldens | CollapsibleSection, Heatmap, Sparkline, BarChart, ExternalLink, PolledRegion, DataTable have no visual goldens (`visualtest/testdata/` — verified absent).                                                     |
| 98  | Add fuzz tests for chart geometry math                    | `ScalePoints`, `ComputeNiceTicks`, `computeArcPath` untested with NaN/Inf/negative/very-large inputs. Pure math functions — perfect fuzz targets. Source: `display/chart_geometry.go`, `display/pie_chart.go`. |
| 99  | Add `waitAnimationSettled` unit test                      | `visualtest/harness.go` — the helper has no dedicated test. Exercised indirectly by every overlay visual test but polling logic, empty-animations path, and timeout path are untested in isolation.            |

### Drift prevention (process hardening)

| #   | Task                                                 | Evidence                                                                                                                                                                                                                                                                                                      |
| --- | ---------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 103 | Write `scripts/check-templ-sync.sh` pre-commit guard | `TestTemplGeneratedInSync` exists but only fires in CI (BuildFlow daemon has 60s budget, no `go test`). A <100ms shell script mirroring `scripts/check-lint-config.sh` would catch `*_templ.go` drift at commit time. Source: `docs/status/2026-08-03_00-29_templ-sync-drift-root-cause-and-process-gaps.md`. |
| 104 | Add CSS freshness CI check                           | Compile demo CSS in CI, diff against committed `examples/demo/static/app.css`, fail if different. `TestCSSFreshness` only warns locally. The v1.7.0 release shipped stale CSS because this check wasn't enforced.                                                                                             |

### Architecture / DRY

| #   | Task                                            | Evidence                                                                                                                                                                                                                                          |
| --- | ----------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 108 | Extract shared LineChart/AreaChart sub-template | `display/line_chart.templ` + `display/area_chart.templ` — ~80% template duplication (axes, gridlines, X-axis labels, legend, empty state identical). Deferred from ADR-0010's 8+ parameter guidance, but the duplication is a maintenance burden. |
| 109 | Add benchmarks for chart geometry helpers       | `BenchmarkScalePoints` + `BenchmarkBuildPolylinePath` exist in `chart_geometry_test.go`, but no benchmarks for PieChart arc computation (`computeSliceAngles` + `computeArcPath` for 100 slices) or full LineChart render.                        |

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
