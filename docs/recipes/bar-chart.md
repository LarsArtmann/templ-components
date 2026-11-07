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

## Vertical needs a definite height (the 0px-bars trap)

Vertical bars are percentage heights, so the chart must have a definite
height to resolve them against. `Height` takes a **CSS value** (not a
Tailwind class — `h-40` is invalid here); when empty it defaults to `12rem`.
The old failure mode — bars silently rendering 0px — is now guarded by a
real-browser test, but keep the rule in mind when embedding the chart in a
collapsed container.

```go
display.BarChart(display.BarChartProps{
    Bars:   topChannels,
    Orient: display.BarVertical,
    Height: "8rem", // CSS value; empty = 12rem
})
```

## Signed series (positive and negative values)

Vertical charts plot negative values below an automatic zero line (drawn in
gray; all-positive series stay bottom-anchored). Negatives take the danger
family by default so the sign survives the render, and `ShowBaseline` adds
the zero line to all-positive charts too.

```go
display.BarChart(display.BarChartProps{
    Bars: []display.BarChartBar{
        {Label: "mar", Value: 420},
        {Label: "apr", Value: -180},
    },
    Orient:          display.BarVertical,
    Height:          "10rem",
    ShowBaseline:    true,
    NegativeBarColor: "bg-rose-700 dark:bg-rose-500", // optional override
})
```

## Capping bar width on small series

A two-bar series stretches each column to half the chart. Cap the columns
and the track centers them:

```go
display.BarChart(display.BarChartProps{
    Bars:        monthlyNet, // 2-6 bars
    Orient:      display.BarVertical,
    MaxBarWidth: "max-w-24",
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
