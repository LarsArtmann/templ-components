# ADR 0040: Component-Level Module Extraction (Kanban, Heatmap, …) — Trigger-Gated

**Status:** Accepted — no extraction; triggers recorded
**Date:** 2026-09-14
**Sources:** go-modularize 7-phase re-assessment · proposal
[docs/modularization/2026-09-14_component-module-boundaries.html](../modularization/2026-09-14_component-module-boundaries.html)

## Context

Recurring question: should large single components or component families —
`display.KanbanBoard`, `display.Heatmap`, the native SVG chart family — become
their own Go modules (like `charts/echarts`, `htmx`, `datastar`), instead of
living in the root module's `display` package?

The repo's module doctrine is already fixed by ADR-0020 (per-package split,
deferred on unmet demand) and ADR-0034 (targeted 7-module split, executed).
This ADR applies that doctrine at component granularity and records the
assessment so future sessions do not re-litigate it.

## Decision

**Do not extract.** KanbanBoard and Heatmap remain `display` package components.
No module, package, or import-path change. The decision is trigger-gated: the
moment any trigger below fires, execute an extraction following the ADR-0034
mechanics, timed with the v2 module-path migration (ADR-0039) because moving
`display.KanbanBoard` to a `kanban` module changes both the import path and the
identifier path — a breaking change.

## Evidence (2026-09-14 audit)

1. **No dependency isolation.** Kanban imports `utils`, `utils/svg`,
   `utils/wire`; Heatmap imports `utils` plus same-package
   `chart_geometry.go` (shared nice-max math with BarChart,
   `display/chart_geometry.go:49,61`). Dependency isolation is the only
   criterion that ever justified a sub-module here (go-error-family →
   errorpage, go-datastar/static → datastar, icons-only adoption → icons,
   ECharts opt-in → charts/echarts). A kanban/heatmap module would have the
   same dependency footprint consumers already have.
2. **No composability payoff.** Zero consumers have asked for kanban-only or
   heatmap-only adoption (ADR-0020 demand triggers unmet since 2026-07-21).
   Importing a module is cheap; importing the `display` _package_ is what Go
   compiles, and the linker prunes unreachable code.
3. **Concrete recurring cost.** A new module touches the release tag set and
   require-bump sweep, the per-module lint/test loops in
   `scripts/release.sh:322,332`, the CI lint matrix, GOWORK=off isolation
   tests, tidy, govulncheck, the `check-module-sync.sh` path table,
   `check-module-layers.sh`, `go.work`, and the require+replace pairs of
   `visualtest` and `website` — 12+ touchpoints, permanently.
4. **Breaking-change gating.** Import-path moves are v2-gated per ADR-0039
   (option 1: migrate at first real breakage).
5. **Cohesion.** 93 commits touched `display/` in Aug–Sep 2026; kanban in 10,
   heatmap in 6 — overwhelmingly introduction + repo-wide sweeps. They are
   frozen, not churn centers. Heatmap shares chart geometry with BarChart;
   extracting it would force `chart_geometry` into another shared module —
   the exact universal-dependency split ADR-0034 rejected for `utils/svg`.

## Triggers (execute an extraction when ANY fires)

| Trigger                                                                                   | Precedent                                                                          |
| ----------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| **T1** — the family needs a unique external dependency no other component uses            | datastar (go-datastar/static), errorpage (go-error-family)                         |
| **T2** — a documented consumer needs family-only adoption without the rest of the library | icons (icons-only adoption)                                                        |
| **T3** — the family ships its own build pipeline or runtime artifact                      | charts/echarts (opt-in adapter); a future WASM/server-driven runtime would qualify |

## Consequences

- The `display` god-package (43 components) stays an **accepted deviation**
  from per-package modularity for v1; revisit only at v2 (ADR-0039) or on a
  trigger above.
- Module count stays at 7 published modules; release, sync guards, and CI
  matrices are unchanged.
- ROADMAP's "Per-package modules split" row points here as the
  component-granularity record.

## References

- [ADR-0020: Per-Package Go Modules Split](0020-per-package-modules-split.md) —
  trigger doctrine origin
- [ADR-0034: Targeted 7-Module Workspace Split](0034-targeted-module-split.md) —
  executed mechanics and rejection rationale
- [ADR-0039: v2 Module-Path Migration Timing](0039-v2-module-path-timing.md) —
  breaking-change gating
- [Modularization proposal (this assessment)](../modularization/2026-09-14_component-module-boundaries.html)
