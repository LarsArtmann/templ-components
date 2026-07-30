# display.Sparkline

Tiny inline SVG line chart for trend visualization. No JavaScript — pure SVG,
uses `currentColor` so it inherits text color.

## When to use

- Dashboard stat cards with trend lines
- Compact data summaries where a full chart is overkill
- Anywhere you need a 120x30px visual trend indicator

## Basic usage

```go
display.Sparkline(display.SparklineProps{
    Values: []float64{1, 3, 2, 5, 4, 6, 3, 7},
})
```

## Filled area

```go
display.Sparkline(display.SparklineProps{
    Values: messageRates,
    Filled: true,
    Class:  "text-green-500 dark:text-green-400",
})
```

## Accessibility

By default the sparkline is `aria-hidden` (it's decorative). When the trend
carries meaningful information, provide `AriaLabel`:

```go
display.Sparkline(display.SparklineProps{
    Values:    dailyActiveUsers,
    AriaLabel: "Daily active users over the last 7 days",
})
```

## Custom dimensions

```go
display.Sparkline(display.SparklineProps{
    Values:      data,
    Width:       200,
    Height:      40,
    StrokeWidth: 2,
})
```

## Min/Max overrides

Min and Max accept `*float64` — `nil` means auto-compute from data. Use this
when you want consistent y-axis scaling across multiple sparklines:

```go
min := 0.0
max := 100.0
display.Sparkline(display.SparklineProps{
    Values: values,
    Min:    &min,
    Max:    &max,
})
```

## With PolledRegion (live-updating)

```go
@htmx.PolledRegion(htmx.PolledRegionProps{URL: "/api/spark", Every: "10s", Eager: true}) {
    @display.Sparkline(display.SparklineProps{
        Values: liveData,
        Filled: true,
    })
}
```
