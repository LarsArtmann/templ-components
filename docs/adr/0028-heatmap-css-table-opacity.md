# ADR-0028: Heatmap Uses CSS Table with Opacity

**Status:** Accepted (2026-07-30)
**Decider:** Lars Artmann

## Context

Activity visualizations need a grid heatmap (e.g. 7-day × 24-hour activity grid, channel × author matrix). The heatmap must be responsive, dark-mode aware, and accessible.

Options considered:

1. **SVG `<rect>` grid** — full control over rendering, but complex for tooltips and responsive sizing
2. **CSS Grid (`div` cells)** — flexible but loses table semantics
3. **HTML `<table>` with opacity backgrounds** — native semantics, easy tooltips via `title`, responsive via `overflow-x-auto`

## Decision

Use HTML `<table>` with inline `background-color: rgba(var(--css-var-rgb), opacity)` styles. This provides:

- Native table semantics for screen readers
- Easy cell tooltips via `title` attribute
- Responsive scrolling via `overflow-x-auto` wrapper
- Consistent with the BarChart ADR (0025) decision to prefer CSS over SVG

Opacity is computed as `value / max`, clamped to a minimum of 0.05 so low-activity cells remain visible. Peak highlighting adds a ring to the max-value cell.

## Consequences

- Depends on a CSS custom property with `-rgb` suffix (e.g. `--ds-brand-rgb`)
- Inline styles required for dynamic opacity (cannot use Tailwind classes for computed values)
- Table layout naturally handles responsive column widths
- `HighlightPeak` must scan all cells on every render (O(rows × cols))
