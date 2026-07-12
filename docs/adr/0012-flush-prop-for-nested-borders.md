# ADR 0012: Flush prop for nested component borders

## Date

2026-07-12

## Status

Accepted

## Context

When `display.Table` is nested inside `display.Card(CardPaddingNone)`, both
components render their own `border border-gray-200 dark:border-gray-700
rounded-lg`. The result is a visible double-border defect — two concentric
1px borders with slightly mismatched corners.

This nesting pattern (Table inside Card with no padding) is the standard
admin dashboard layout. It appears on virtually every data-listing page:
users, tenants, audit logs, orders, etc.

The consumer (cqrs-htmx adminui) had a fragile CSS workaround:

```css
.overflow-hidden > .overflow-x-auto {
  border: 0 !important;
  border-radius: 0 !important;
}
```

This selector depends on Tailwind utility class names that could change in
a future version. It also couples the consumer's CSS to the library's
internal class structure.

## Decision

Add a `Flush bool` field to `TableProps`. When `true`, the table wrapper div
suppresses its own `border` and `rounded-lg` classes, keeping only
`overflow-x-auto`. The consumer opts in explicitly:

```go
@display.Card(display.CardProps{Title: "Users", Padding: display.CardPaddingNone}) {
    @display.Table(display.TableProps{Rows: rows, Flush: true})
}
```

### Why `Flush` over alternatives

| Option                                | Verdict      | Reason                                                                                                                             |
| ------------------------------------- | ------------ | ---------------------------------------------------------------------------------------------------------------------------------- |
| `Flush bool`                          | **Accepted** | Narrow scope — only suppresses the wrapper border. Matches CSS terminology. Concise.                                               |
| `Borderless bool`                     | Rejected     | Too broad — implies removing ALL borders including the `Bordered` table element's cell borders. Ambiguous.                         |
| Auto-detect parent context            | Rejected     | Components should not inspect their parent/children's type. Magic behavior, hard to debug, couples components to nesting patterns. |
| Split wrapper into separate component | Rejected     | Breaks the simple API for standalone tables. Over-engineering for a common pattern.                                                |

### Why `Flush` is Table-specific (not a shared concept)

`Flush` stays on `TableProps` only. A shared `Flushable` interface or
`BaseProps` flag would be premature — Table-in-Card is the only known case.
If a second component needs the same pattern (e.g., `Form` inside
`CardPaddingNone`), promote `Flush` to a shared concept at that time.
This follows the YAGNI principle and the library's promotion trigger
convention (see [ADR 0010](0010-sub-template-extraction-pattern.md)).

### Wrapper div bypasses `utils.Class()`

The `tableWrapperClass()` result is used directly in the template
(`class={ tableWrapperClass(props.Flush) }`), not through `utils.Class()`.
This is intentional:

1. The wrapper is an internal element — consumer `props.Class` is applied to
   the `<table>` element, not the wrapper.
2. `utils.Class()` uses tailwind-merge-go which reorders classes, breaking
   existing tests that assert `"border border-gray-200"` as a substring.

## Consequences

- **Positive:** Eliminates the double-border defect with a single boolean.
  Consumers can delete fragile CSS workarounds. Type-safe, explicit, no magic.
- **Positive:** The pattern is documented for future components that may need
  the same opt-out (see [table-in-card recipe](../recipes/table-in-card.md)).
- **Negative:** Consumers must know to set `Flush: true` — it's not automatic.
  Mitigated by godoc, recipe docs, and the demo.
- **Neutral:** `Flush` interacts with `CardPaddingNone` — both are needed for
  the correct table-in-card layout. Neither alone is sufficient.

## Related

- [Recipe: Table Inside Card](../recipes/table-in-card.md)
- [Bug report: Table-in-Card Double Border](../feedback/2026-07-12_table-in-card-double-border.md)
