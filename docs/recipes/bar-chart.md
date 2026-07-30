# display.BarChart

CSS-based bar chart for categorical data. Horizontal or vertical orientation,
per-bar colors, clickable labels, value formatting. No SVG, no JavaScript.

## When to use

- Top-N breakdowns (channels by message count, authors by activity)
- Category comparisons
- Compact data visualizations in dashboard cards

## Basic usage

```go
display.BarChart(display.BarChartProps{
    Bars: []display.BarChartBar{
        {Label: "general", Value: 1200},
        {Label: "random",  Value: 800},
        {Label: "dev",     Value: 450},
    },
})
```

## Vertical (column chart)

```go
display.BarChart(display.BarChartProps{
    Bars:   topChannels,
    Orient: display.BarVertical,
})
```

## Per-bar colors

```go
display.BarChart(display.BarChartProps{
    Bars: []display.BarChartBar{
        {Label: "healthy", Value: 900, Color: "bg-emerald-600 dark:bg-emerald-500"},
        {Label: "warning", Value: 500, Color: "bg-amber-600 dark:bg-amber-500"},
        {Label: "critical", Value: 100, Color: "bg-red-600 dark:bg-red-500"},
    },
})
```

## Clickable labels

```go
display.BarChart(display.BarChartProps{
    Bars: []display.BarChartBar{
        {Label: "general", Value: 1200, Href: "/channels/general"},
        {Label: "random",  Value: 800,  Href: "/channels/random"},
    },
})
```

## Custom value formatting

```go
display.BarChart(display.BarChartProps{
    Bars: data,
    ValueFormat: func(v float64) string {
        return humanize.FormatFloat("#,###", v)
    },
})
```

## Empty state

```go
display.BarChart(display.BarChartProps{
    Bars:         emptySlice,
    EmptyMessage: "No messages in this period",
})
```

## Accessibility

The chart container has `role="img"`. Provide `AriaLabel` for screen readers:

```go
display.BarChart(display.BarChartProps{
    Bars:      channels,
    AriaLabel: "Messages by channel, last 7 days",
})
```
