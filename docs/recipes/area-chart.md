# Recipe: Area Chart

Pure-SVG area chart — a line chart with semi-transparent filled areas. Zero JavaScript.

## Basic Usage

```go
@display.AreaChart(display.AreaChartProps{
    Series: []display.LineChartSeries{
        {Name: "Active Users", Values: []float64{120, 180, 250, 300, 280, 340}},
    },
    XAxisLabels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"},
})
```

## Multi-Series

```go
@display.AreaChart(display.AreaChartProps{
    Series: []display.LineChartSeries{
        {Name: "Active", Values: []float64{120, 180, 250, 300, 280, 340}},
        {Name: "Churned", Values: []float64{20, 30, 25, 40, 35, 50}},
    },
    XAxisLabels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"},
})
```

## Custom Fill Opacity

```go
@display.AreaChart(display.AreaChartProps{
    Series: []display.LineChartSeries{
        {Name: "Data", Values: data},
    },
    FillOpacity: 0.4, // 0.0–1.0, default 0.2
})
```

## Smooth Curves

```go
@display.AreaChart(display.AreaChartProps{
    Series: []display.LineChartSeries{
        {Name: "Trend", Values: data},
    },
    Style: display.LineChartStyleSmooth,
})
```

## Props Reference

| Prop             | Type                   | Default | Description                  |
| ---------------- | ---------------------- | ------- | ---------------------------- |
| `Series`         | `[]LineChartSeries`    | —       | Data series                  |
| `XAxisLabels`    | `[]string`             | —       | Category labels              |
| `Width`/`Height` | `int`                  | 600/300 | SVG canvas dimensions        |
| `ShowGrid`       | `bool`                 | true    | Dashed gridlines             |
| `ShowDots`       | `bool`                 | false   | Data point circles           |
| `ShowLegend`     | `bool`                 | true    | Color-swatch legend          |
| `Style`          | `LineChartStyle`       | Linear  | Linear or Smooth             |
| `FillOpacity`    | `float64`              | 0.2     | Area fill transparency (0–1) |
| `ValueFormat`    | `func(float64) string` | auto    | Y-axis tick formatter        |
