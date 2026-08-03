# Recipe: ECharts Adapter (Tier 2 Interactive Charts)

Opt-in wrapper for [Apache ECharts] — 25+ interactive chart types, tooltips,
zoom/pan, and legend toggle. Requires ~1MB JS from CDN.

## When to Choose Tier 1 vs Tier 2

| Use Tier 1 (Native SVG) when... | Use Tier 2 (ECharts) when...                    |
| ------------------------------- | ----------------------------------------------- |
| You need zero JavaScript        | You need interactive tooltips on hover          |
| CSP-strict environment          | You need zoom/pan/drag selection                |
| Static or PDF-rendered charts   | You need legend toggle (show/hide series)       |
| Simple line/pie/area charts     | You need exotic types (radar, scatter, geo, 3D) |
| Minimal bundle size is critical | You need animations and transitions             |

## Setup

### 1. Add go-echarts to your project

```bash
go get github.com/go-echarts/go-echarts/v2
```

### 2. Inject the ECharts runtime once per page

```go
@echarts.SDKScript(echarts.DefaultSDKScriptProps())
```

### 3. Render a chart

```go
import charts "github.com/go-echarts/go-echarts/v2/charts"

bar := charts.NewBar()
bar.SetXAxis([]string{"Mon", "Tue", "Wed"}).
    AddSeries("Sales", []int{120, 200, 150})
snippet, _ := bar.RenderSnippet()

@echarts.EChart(echarts.EChartsProps{
    Element: snippet.Element,
    Script:  snippet.Script,
    Nonce:   nonce,
})
```

## Self-Hosting ECharts

For air-gapped or CSP-strict environments, self-host `echarts.min.js`:

```go
@echarts.SDKScript(echarts.SDKScriptProps{
    Src:   "/static/vendor/echarts.min.js",
    Nonce: nonce,
})
```

## Dark Mode

The dark mode bridge is auto-injected by `EChart`. It watches the `.dark` class
on `<html>` and applies dark/light color overrides to all ECharts instances via
`setOption({merge: true})`. No consumer-side wiring needed.

To disable: `EChartsProps{DarkModeBridge: false}`.

## Custom CDN

```go
@echarts.SDKScript(echarts.SDKScriptProps{
    Version: "5.4.3",
    CDN:     "https://unpkg.com",
    Nonce:   nonce,
})
```

[Apache ECharts]: https://echarts.apache.org/
