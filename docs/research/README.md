# Research Index

Primary-source research and deep-dives that inform templ-components
decisions. Reports are point-in-time snapshots: verify claims against the
current tree before editing based on them (the 2026-10-14 report's D1–D3
corrections are the case study).

## Deep-dives

| Report | Date | Subject | Follow-through |
| --- | --- | --- | --- |
| [2026-10-14 Tailwind CSS deep-dive](2026-10-14_tailwindcss-deep-dive.html) | 2026-10-14 | Tailwind v4.3.x audit: the remaining composability gaps (`@variant dark`, `@utility`, `@theme` animation tokens, `@plugin`, `aria-*` state variants) | ACTIVE — executed by [the composability Pareto plan](../planning/2026-10-14_20-45_TAILWIND-V4-COMPOSABILITY-PARETO-PLAN.md) |
| [2026-10-08 Tailwind CSS deep-dive](2026-10-08_tailwindcss-deep-dive.html) | 2026-10-08 | Tailwind v4 audit; its fixes are verified in-tree | Superseded by the 2026-10-14 report for the remaining gap list |
| [go-sse-go / Datastar deep-dive](2026-08-21_go-sse-go-datastar-deep-dive.html) | 2026-08-21 | SSE transport landscape for the Datastar integration | — |
| [UI library design research](ui-library-design-research.md) | — | How to build a superb web UI component library | — |

## Notes

- [css-subgrid.md](css-subgrid.md) — CSS subgrid status and templ-components implications.
- [datastar-integration-analysis.md](datastar-integration-analysis.md) — how templ-components could benefit from Datastar.
- [direction-context-propagation.md](direction-context-propagation.md) — shadcn-templ's DirectionProvider pattern comparison (RTL).
- [modern-browser-capabilities.md](modern-browser-capabilities.md) — modern browser capabilities inventory.
- [popover-api.md](popover-api.md) — Popover API + CSS anchor positioning baseline.
- [vibe-kanban-analysis.md](vibe-kanban-analysis.md) — what KanbanBoard can learn from BloopAI/vibe-kanban.
- [what-we-are-missing.md](what-we-are-missing.md) — hidden gaps and unknown unknowns.

## Conventions

- Findings promoted to work go through a Pareto plan in `docs/planning/`;
  the plan, not the report, is the execution record.
- Secondhand claims inside a report get primary-source verification before
  anyone edits based on them; corrections are annotated in place with the
  verifying command and date.
