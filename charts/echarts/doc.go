// Package echarts provides a templ component adapter for [Apache ECharts] — a
// powerful, ~1 MB JavaScript charting library with 25+ chart types, interactive
// tooltips, zoom/pan, and legend toggles.
//
// This package is **opt-in** and follows the same zero-dependency pattern as the
// [datastar] package: it does NOT import go-echarts. The consumer builds their
// chart using [go-echarts], calls chart.RenderSnippet() to get
// {Element, Script, Option}, and passes those three strings to [EChart]. This
// keeps the core templ-components dependency graph clean.
//
// # Quick start
//
// Add go-echarts to your go.mod:
//
//	go get github.com/go-echarts/go-echarts/v2
//
// Inject the ECharts runtime once per page (in your layout):
//
//	@echarts.SDKScript(echarts.DefaultSDKScriptProps())
//
// Render a chart from go-echarts output:
//
//	bar := charts.NewBar()
//	bar.SetXAxis([]string{"Mon", "Tue", "Wed"}).
//	    AddSeries("Category", []int{120, 200, 150})
//	snippet, _ := bar.RenderSnippet()
//	@echarts.EChart(echarts.EChartsProps{
//	    Element: snippet.Element,
//	    Script:  snippet.Script,
//	    Nonce:   nonce,
//	})
//
// # When to choose Tier 2 (ECharts) vs Tier 1 (native SVG)
//
// Use the native SVG charts (display.LineChart, display.PieChart,
// display.AreaChart) when you need:
//   - Zero JavaScript
//   - Static or server-rendered charts
//   - CSP compliance without extra scripts
//
// Use ECharts (this package) when you need:
//   - Interactive tooltips on hover
//   - Zoom/pan/drag selection
//   - Legend toggle (show/hide series)
//   - Exotic chart types (radar, scatter, heatmap, geo, 3D)
//   - Animations and transitions
//
// See docs/adr/0031-two-tier-chart-architecture.md for the full rationale.
//
// # Dark mode
//
// The dark mode bridge script (auto-injected by [EChart]) watches the `.dark`
// class on `<html>` and applies the ECharts "dark" theme at runtime. This keeps
// ECharts charts in sync with your Tailwind dark mode toggle without any
// consumer-side wiring.
//
// [Apache ECharts]: https://echarts.apache.org/
// [go-echarts]: https://github.com/go-echarts/go-echarts
// [datastar]: https://data-star.dev/
package echarts
