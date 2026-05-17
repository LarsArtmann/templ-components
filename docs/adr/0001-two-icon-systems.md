# ADR: Two Icon Systems — Stroke (24×24) and Filled (20×20)

**Date:** 2026-05-17
**Status:** Accepted

## Context

The codebase has two parallel SVG icon systems:

1. **`icons.Icon`** — 24×24 stroke-based icons via `iconPathData` map. Uses `stroke="currentColor"`, `stroke-width="1.5"`. Used for navigation, buttons, status indicators.

2. **`internal/svg.FillIcon`** — 20×20 filled icons via `fillIcon` templ component. Uses `fill="currentColor"`. Used for chevrons (accordion, dropdown), trend arrows (stat card), and status indicators.

## Decision

Keep both systems. They serve different visual purposes:

- **Stroke icons** — Line-style icons for UI controls, menus, actions. Standard HeroIcons outline style.
- **Filled icons** — Solid directional indicators where visual weight matters (chevrons, arrows, status dots).

## Rationale

- Merging would require a `variant` parameter on every icon call, adding complexity.
- The filled icons are small directional helpers (chevrons, arrows), not a parallel icon library.
- `internal/svg` is intentionally internal — consumers use `icons.Icon` for their own icons.

## Consequences

- Two small SVG systems coexist.
- `internal/svg.FillIcon` is not exported; only used internally by display/navigation packages.
- If filled icons grow beyond ~10, reconsider merging with a variant parameter.
