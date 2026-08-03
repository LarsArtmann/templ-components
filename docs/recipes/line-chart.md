# Recipe: Line Chart

Pure-SVG line chart with axes, gridlines, multi-series support, and a legend.
Zero JavaScript.

## Basic Usage

```go
@display.LineChart(display.LineChartProps{
    Series: []display.LineChartSeries{
        {Name: "Revenue", Values: []float64{10, 25, 40, 35, 60, 55, 80}},
    },
    XAxisLabels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
})
```

## Multi-Series

```go
@display.LineChart(display.LineChartProps{
    Series: []display.LineChartSeries{
        {Name: "Revenue", Values: []float64{10, 25, 40, 35, 60, 55, 80}},
        {Name: "Expenses", Values: []float64{8, 15, 20, 30, 35, 40, 50}},
    },
    XAxisLabels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
})
```

## Smooth Curves

```go
@display.LineChart(display.LineChartProps{
    Series: []display.LineChartSeries{
        {Name: "Trend", Values: []float64{10, 25, 40, 35, 60, 55, 80}},
    },
    XAxisLabels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
    Style:       display.LineChartStyleSmooth,
})
```

## Custom Colors

Each series accepts a Tailwind `text-*` class. The line uses `currentColor`.

```go
@display.LineChart(display.LineChartProps{
    Series: []display.LineChartSeries{
        {Name: "Growth", Values: []float64{5, 10, 20, 35, 55, 80, 110},
         Color: "text-emerald-600 dark:text-emerald-400"},
    },
})
```

## Accessibility

Set `AriaLabel` for screen readers. Without it, the chart is `aria-hidden`
(decorative).

```go
@display.LineChart(display.LineChartProps{
    BaseProps:   utils.BaseProps{AriaLabel: "Revenue trend from January to June"},
    Series:      []display.LineChartSeries{{Values: monthlyRevenue}},
})
```

## Props Reference

| Prop          | Type                          | Default     | Description                        |
| ------------- | ----------------------------- | ----------- | ---------------------------------- |
| `Series`      | `[]LineChartSeries`           | —           | Data series (Name, Values, Color)  |
| `XAxisLabels` | `[]string`                    | —           | Category labels for the X-axis     |
| `Width`       | `int`                         | 600         | SVG canvas width                   |
| `Height`      | `int`                         | 300         | SVG canvas height                  |
| `ShowGrid`    | `bool`                        | true        | Dashed gridlines                   |
| `ShowDots`    | `bool`                        | true        | Data point circles                 |
| `ShowLegend`  | `bool`                        | true        | Color-swatch legend (2+ series)    |
| `Style`       | `LineChartStyle`              | Linear      | Linear or Smooth curves            |
| `Min`/`Max`   | `*float64`                    | auto        | Y-axis range override              |
| `ValueFormat` | `func(float64) string`        | FormatTick  | Y-axis tick label formatter        |
