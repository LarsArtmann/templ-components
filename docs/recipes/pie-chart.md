# Recipe: Pie Chart

Pure-SVG pie and donut charts with labels and legend. Zero JavaScript.

## Basic Pie

```go
@display.PieChart(display.PieChartProps{
    Slices: []display.PieChartSlice{
        {Label: "Direct", Value: 45},
        {Label: "Organic", Value: 30},
        {Label: "Referral", Value: 25},
    },
})
```

## Donut Chart with Center Label

```go
@display.PieChart(display.PieChartProps{
    Slices: []display.PieChartSlice{
        {Label: "Used", Value: 80},
        {Label: "Free", Value: 48},
    },
    Donut:       true,
    CenterLabel: "128GB",
})
```

## Custom Colors

```go
@display.PieChart(display.PieChartProps{
    Slices: []display.PieChartSlice{
        {Label: "High", Value: 70, Color: "text-emerald-600 dark:text-emerald-400"},
        {Label: "Med", Value: 20, Color: "text-amber-600 dark:text-amber-400"},
        {Label: "Low", Value: 10, Color: "text-rose-600 dark:text-rose-400"},
    },
})
```

## Accessibility

```go
@display.PieChart(display.PieChartProps{
    BaseProps: utils.BaseProps{AriaLabel: "Traffic sources breakdown"},
    Slices:    slices,
})
```

## Props Reference

| Prop           | Type                 | Default | Description                     |
| -------------- | -------------------- | ------- | ------------------------------- |
| `Slices`       | `[]PieChartSlice`    | —       | Data (Label, Value, Color)      |
| `Width`        | `int`                | 400     | SVG canvas width                |
| `Height`       | `int`                | 300     | SVG canvas height               |
| `Donut`        | `bool`               | false   | Donut (ring) vs full pie        |
| `InnerRadius`  | `float64`            | 0.6     | Donut hole size (0.0–1.0)       |
| `ShowLabels`   | `bool`               | true    | External labels with percentages|
| `ShowLegend`   | `bool`               | true    | Color-swatch legend             |
| `CenterLabel`  | `string`             | —       | Center text (donut only)        |
