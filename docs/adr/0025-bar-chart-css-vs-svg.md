# ADR-0025: BarChart CSS vs SVG

**Status:** Accepted (2026-07-30)
**Decider:** Lars Artmann

## Context

The library needed a bar chart component. Two rendering approaches:

1. **SVG** — `<rect>` elements positioned via calculated coordinates
2. **Pure CSS** — `<div>` elements with `width:` or `height:` percentages

SVG offers pixel-perfect control and is the standard for data visualization
libraries (D3, Recharts, Chart.js). CSS is simpler, inherently responsive,
and requires no coordinate math.

## Decision

**Use pure CSS** for BarChart. Reserve SVG for Sparkline (where precise
polyline rendering at small sizes is essential).

Rationale:
- The primary use case is dashboard stat breakdowns (top-N categories), not
  scientific data visualization.
- CSS bar charts are inherently responsive — `width: 75%` scales with the
  container. SVG requires `viewBox` and `preserveAspectRatio` juggling.
- Dark mode works via Tailwind `dark:` variants on the bar's `bg-*` class.
  SVG would need stroke/fill color management.
- The CSS approach matches what DiscordSync (the reference consumer) was
  already hand-rolling — this component is a direct extraction.

## Consequences

- Vertical (column) charts use `height:` percentages inside a flex container.
  This works well but may need `min-height` on the container for very short bars.
- Per-bar colors use Tailwind `bg-*` classes, not SVG `fill` attributes.
- No axes, gridlines, or tooltips (out of scope — use a charting library for that).
- The chart is accessible via `role="img"` + `aria-label`.
