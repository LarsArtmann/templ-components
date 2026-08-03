# ADR-0031: Two-Tier Chart Architecture (Native SVG + Opt-in ECharts)

Date: 2026-08-03

## Status

Accepted

## Context

templ-components had three basic visualization components — Sparkline (SVG
polyline), BarChart (CSS flexbox), and Heatmap (CSS grid) — all zero-JS per
ADR-0025. However, the **single biggest gap** for real-world dashboards was the
lack of:

1. **Line charts with axes** — every metrics/monitoring dashboard needs time-series
2. **Pie/donut charts** — proportions are the #2 chart need
3. **Interactive charts** — tooltips, zoom, legend toggle

The `recipes.Dashboard` component has `Charts []templ.Component` slots, but
consumers had no chart components to fill them. Two options were considered:

**Option A: Only native SVG charts.** Zero-JS, dependency-free. Covers ~80% of
use cases (static dashboards, PDF reports, CSP-strict environments). Leaves power
users to build their own ECharts integration from scratch — the exact
duplication the library exists to prevent.

**Option B: Only an ECharts adapter.** Full interactivity, 25+ chart types.
Forces a ~1MB JS dependency on everyone, even users who just want a simple
static line chart.

## Decision

Implement a **two-tier chart architecture**:

- **Tier 1: Native SVG charts** (`display` package) — LineChart, PieChart,
  AreaChart, DonutChart. Zero JavaScript. Server-rendered SVG with Tailwind
  dark: variants. Covers ~80% of dashboard chart needs.

- **Tier 2: Opt-in ECharts adapter** (`charts/echarts` package) — a thin
  wrapper component that accepts go-echarts `RenderSnippet()` output as three
  strings (Element, Script, Option). Does NOT import go-echarts. Covers the
  remaining ~20% requiring interactivity (tooltips, zoom, exotic chart types).

This mirrors the library's existing pattern: HTMX by default (zero JS),
Datastar opt-in for power users.

## Rationale

1. **Consumer choice:** Most dashboards need simple static charts. Forcing a
   1MB JS bundle for a 7-day revenue line chart is overkill. But power users
   who need interactive charts shouldn't have to build their own CSP-safe,
   dark-mode-aware ECharts wrapper from scratch.

2. **The `datastar` precedent:** The library already has an opt-in package
   (`datastar/`) that wraps a JS library without importing the SDK. The
   `charts/echarts` package follows this exact pattern.

3. **Shared geometry module:** The `chart_geometry.go` helpers (ScalePoints,
   ComputeNiceTicks, BuildPolylinePath, BuildSmoothPath) are the "1% that
   delivers 51%" — every Tier 1 chart type composes these primitives. This
   keeps the SVG chart code DRY and consistent.

4. **PieChart uses SVG arc paths** (`<path d="M...A...">`) rather than the
   CSS `stroke-dasharray` trick on circles. Arc paths support donut holes,
   label positioning at arc midpoints, and non-integer totals.

## Consequences

- The `display` package gains three new components (LineChart, PieChart,
  AreaChart) plus shared geometry helpers. The component count increases.
- The `charts/echarts` package is a new opt-in package. It does NOT add any
  dependency to the core library's `go.mod`.
- Consumers needing interactive charts must `go get go-echarts` separately.
- The dark mode bridge script in the ECharts adapter syncs ECharts themes
  with the Tailwind `.dark` class via MutationObserver.
- Stacked/grouped BarChart is deferred to Tier 2 (ECharts) — too complex for
  the native SVG approach.
