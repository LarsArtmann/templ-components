# Recipe: Accessible Clickable Rows

`TableRow.Href` makes an entire row navigate on click or Enter — but the row
is a **pointer/keyboard convenience, not the accessible link**. A `<tr>`
cannot carry `role="link"` (ARIA in HTML forbids roles on rows inside a real
table, and the W3C validator rejects every page that tries), so a
screen-reader user never hears "link" from the row itself. Pair every
clickable row with a real `<a>` in a cell; keep the row-level `Href` as the
convenience layer on top.

## The contract in one line

| Layer                       | Who it serves                       | Mechanism                                       |
| --------------------------- | ----------------------------------- | ----------------------------------------------- |
| Cell link (`<a>` in a cell) | Screen readers, keyboard-only users | Native link semantics, in cell order            |
| Row (`TableRow.Href`)       | Pointer users, Enter-on-focused-row | `data-tc-row-href` + CSP-safe delegation script |

The library enforces its half: a row with `Href` gets `tabindex="0"`, a
`focus-visible` outline (browsers draw no default ring on a `<tr>`), and a
delegation script that follows the row target on click/Enter — while never
hijacking clicks on real links or buttons inside the row.

## Recommended: visible link in the first cell

Give the identifying cell a plain `<a>` through `TableCell.Content`. Screen
readers announce it in cell order ("link, INV-1042"), keyboard users reach it
with Tab, and pointer users can still click anywhere on the row:

```go
@display.Table(display.TableProps{
    TypedHeaders: []display.TableHeader{
        {Label: "Invoice"},
        {Label: "Amount"},
        {Label: "Status"},
    },
    Rows: []display.TableRow{
        {
            Href: "/invoices/1042",
            Cells: []display.TableCell{
                {Content: linkCell("/invoices/1042", "INV-1042")},
                {Text: "€1,204.00"},
                {Text: "Open"},
            },
        },
    },
})

templ linkCell(href, label string) {
    <a href={ href } class="font-medium text-accent-600 underline-offset-2 hover:underline dark:text-accent-400">
        { label }
    </a>
}
```

The link doubles as the skip-free keyboard path: a screen-reader user
browsing links hears every row's destination, without ever meeting the
row-level behavior.

## Variation: stretched link (rows without other controls)

If a row contains no other interactive element, one absolutely-positioned
link can cover the whole row so the native link — semantics included — is
what you actually click. The row's own `Href` becomes redundant (the browser
follows the `<a>` before the delegation script sees the click), but keep it:
users who zoom or override styles may get the row back as a plain target.

```go
templ stretchedLinkCell(href, label string) {
    <td class="relative">
        <a
            href={ href }
            class="after:absolute after:inset-0 after:z-10"
        >
            <span class="sr-only">Open invoice </span>{ label }
        </a>
    </td>
}
```

The containing `<td>` needs `relative` (positioning against the `<tr>`
itself is unreliable — table-row boxes are not positioning contexts in every
engine). The stretched overlay comes with real costs; use it only when all
of these hold:

- **One per row.** A second link or a button under the overlay becomes
  unclickable — hit the plain first-cell pattern instead.
- **Text selection dies** across the covered cells. Fine for navigation
  lists; wrong for data people copy (amounts, emails, IDs).
- **Hover styling moves to the link.** Per-cell hover effects stop firing
  because the pointer is over the overlay, not the cells.

When in doubt, ship the visible first-cell link: it costs nothing and has no
edge cases.

## Anti-patterns

- **`role="link"` on the `<tr>`.** Invalid per ARIA in HTML; the W3C
  validator fails the page. Early releases of this library shipped exactly
  that; it is removed on master and the fix ships in the next release (#76).
- **The row as the only interactive thing.** A clickable row without any
  cell link is invisible to screen-reader link lists and unreachable for
  keyboard-only users who skip focusable rows.
- **`href="#"` rows completed by JS.** Use the real destination in both the
  cell link and `TableRow.Href`; the delegation script navigates natively.

## Related

- [Custom Table Rows](custom-table-rows.md) — `Table.Body` slot for fully
  hand-rolled rows (apply the same contract: real link in a cell, row as
  convenience)
- [Table Inside Card](table-in-card.md) — the flushed-table layout the
  clickable-row focus outline was tuned in
