# Heatmap

Render a CSS-based grid heatmap where each cell's background opacity reflects its value relative to the maximum.

## Basic Usage

```templ
@display.Heatmap(display.HeatmapProps{
    Rows: []display.HeatmapRow{
        {Label: "Mon", Cells: []display.HeatmapCell{
            {Value: 5, Label: "Mon 00:00 — 5 msgs"},
            {Value: 12, Label: "Mon 06:00 — 12 msgs"},
            {Value: 0},
            {Value: 8, Label: "Mon 18:00 — 8 msgs"},
        }},
        {Label: "Tue", Cells: []display.HeatmapCell{
            {Value: 3},
            {Value: 20},
            {Value: 7},
            {Value: 2},
        }},
    },
    ColumnLabels: []string{"00:00", "06:00", "12:00", "18:00"},
})
```

## Peak Highlighting

Set `HighlightPeak: true` to add a ring around the cell with the highest value:

```templ
@display.Heatmap(display.HeatmapProps{
    Rows:          rows,
    ColumnLabels:  labels,
    HighlightPeak: true,
})
```

## Showing Values Inside Cells

```templ
@display.Heatmap(display.HeatmapProps{
    Rows:       rows,
    ShowValues: true,
})
```

## Max Override

By default the max is auto-computed from the data. Override it to normalize across multiple heatmaps:

```templ
@display.Heatmap(display.HeatmapProps{
    Rows: rows,
    Max:  100,
})
```

## Clickable Cells

Each cell can link to a detail page:

```templ
{Label: "Mon", Cells: []display.HeatmapCell{
    {Value: 5, Href: "/activity?day=mon&hour=0"},
}},
```

## Custom Color Variable

Cells use `rgba(var(--ds-brand-rgb), opacity)` by default. Override with any CSS custom property that has an `-rgb` companion:

```templ
@display.Heatmap(display.HeatmapProps{
    Rows:     rows,
    ColorVar: "--my-accent",
})
```

Requires `--my-accent-rgb: R, G, B` in your CSS.

## Accessibility

- Uses native `<table>` semantics for screen-reader compatibility
- Each cell has a `title` tooltip (set via `HeatmapCell.Label`)
- Supports `AriaLabel` on the container for overall context
- Empty state renders a configurable message
