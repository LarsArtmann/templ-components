# Recipe: Table Inside Card

The standard admin dashboard pattern: a `Table` nested inside a `Card` with
`CardPaddingNone`. Without special handling this produces a double border
(the card border + the table wrapper border). The `Flush` prop solves this.

## Basic: Flush + compact padding

```go
@display.Card(display.CardProps{
    Title:   "Users",
    Padding: display.CardPaddingNone,
}) {
    @display.Table(display.TableProps{
        Headers:     []string{"Name", "Email", "Role"},
        Rows:        rows,
        Flush:       true,
        CellPadding: display.TableCellPaddingCompact,
        Hover:       true,
        Striped:     true,
    })
}
```

- **`Flush: true`** suppresses the table wrapper's `border`, `rounded-lg`, keeping
  only `overflow-x-auto`. The card's own border provides the visual boundary.
- **`CardPaddingNone`** renders the table as a direct child of the card shell,
  with no padding wrapper. This is required for the table to fill the card edge-to-edge.
- **`CellPadding: TableCellPaddingCompact`** uses `px-4 py-2` (vs the default
  `px-4 py-3`). Ideal for data-heavy admin panels where rows should be scannable.

## Without Flush (standalone table)

A standalone table NOT inside a card should NOT use `Flush` — it needs its own
border and rounded corners:

```go
@display.Table(display.TableProps{
    Headers: []string{"Name", "Email", "Role"},
    Rows:    rows,
    Hover:   true,
})
```

## Combining with sortable headers

Flush and CellPadding compose with all other Table features:

```go
@display.Card(display.CardProps{
    Title:   "Audit Log",
    Padding: display.CardPaddingNone,
}) {
    @display.Table(display.TableProps{
        TypedHeaders: []display.TableHeader{
            {Label: "Timestamp", Sortable: true, SortDirection: sortDir, Href: sortURL("ts", sortDir)},
            {Label: "Actor"},
            {Label: "Action"},
        },
        Rows:        auditRows,
        Flush:       true,
        CellPadding: display.TableCellPaddingCompact,
    })
}
```

## Combining with clickable rows

Clickable rows (`TableRow.Href`) work inside flushed tables:

```go
@display.Card(display.CardProps{
    Title:   "Tenants",
    Padding: display.CardPaddingNone,
}) {
    @display.Table(display.TableProps{
        Headers: []string{"Name", "Plan", "Status"},
        Rows: []display.TableRow{
            {Cells: []string{"Acme Corp", "Pro", "Active"}, Href: "/tenants/acme"},
            {Cells: []string{"Globex", "Basic", "Trial"}, Href: "/tenants/globex"},
        },
        Flush:       true,
        CellPadding: display.TableCellPaddingCompact,
        Hover:       true,
    })
}
```

## Migration from CSS workaround

If you previously used a CSS hack to suppress the double border:

```css
/* OLD — remove this */
.overflow-hidden > .overflow-x-auto {
  border: 0 !important;
  border-radius: 0 !important;
}
```

Replace with `Flush: true` on the `TableProps`:

```go
// BEFORE — relies on fragile CSS selector
@display.Card(display.CardProps{Padding: display.CardPaddingNone}) {
    @display.Table(display.TableProps{Rows: rows})
}

// AFTER — explicit, type-safe
@display.Card(display.CardProps{Padding: display.CardPaddingNone}) {
    @display.Table(display.TableProps{Rows: rows, Flush: true})
}
```

Then delete the CSS workaround from your stylesheet.

## Available cell padding levels

| Value                         | Padding     | Use case                                          |
| ----------------------------- | ----------- | ------------------------------------------------- |
| `TableCellPaddingComfortable` | `px-4 py-3` | Default. Spacious, good for public-facing tables. |
| `TableCellPaddingCompact`     | `px-4 py-2` | Dense data grids, admin dashboards, audit logs.   |

## Why not auto-detect?

The `Flush` prop is explicit rather than auto-detected from the parent context.
Components should not inspect their parent or children's type — this is "magic
behavior" that makes debugging difficult and couples components to specific
nesting patterns. The consumer opts in with `Flush: true`.
