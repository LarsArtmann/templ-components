# Bounding Unbounded Tables

## When to use this

Your data set has no natural page size — an event catalog, a message list, a
backfill log. Rendering every row produces a 500-row DOM with no truncation
notice, no paging, and a page that gets slower as the data grows. A consumer
audit (go-cqrs-lite docserver, 2026-10-01) flagged exactly this shape: four
tables rendering every row with "no cap, note, or pagination — `ListNote` and
`Pagination` sit unused".

The library ships the whole answer. This recipe is the decision tree plus the
wiring.

## The decision tree

Answer two questions about your page:

1. **Must the user be able to reach row 400 without a search?**
   - **No** → cap + notice (the common case). Render the first N rows and put
     a `ListNote` under the table telling the user what happened and what to
     do about it.
   - **Yes** → paginate. Pick by data shape:
     - Bounded / stable ordering → `navigation.Pagination` (page numbers).
     - Unbounded / concurrent inserts → `navigation.LoadMore` (cursor; see
       [Cursor Pagination](cursor-pagination.md)).

2. **Is the table inside an HTMX swap?** Then the cap/pager re-renders with
   each response — keep the cap constant per request and echo it into
   `ListNote`.

## Pattern A: cap + truncation notice (the default choice)

```templ
const catalogCap = 50

@display.Table(display.TableProps{
    Headers: []string{"Message", "Producers", "Consumers"},
    Rows:    rows, // already sliced to catalogCap
})
@display.ListNote(display.ListNoteProps{
    Shown: len(rows),
    Total: totalMessages, // the full match count from the same query
})
```

`ListNote` renders **nothing** when `Total <= Shown` — the common small-table
case needs no conditional on your side. When the cap bites, it renders
"Showing 50 of 500. Narrow your search to see more."

Variant semantics (pick per page meaning, see `ListNoteVariant` docs):

| Variant             | Renders                            | Use when                                                 |
| ------------------- | ---------------------------------- | -------------------------------------------------------- |
| `ListNoteTruncated` | "Showing N of M. Narrow your…"     | A search/filter narrows results; hidden rows are noise    |
| `ListNoteCount`     | "Showing N items." (always)        | List semantics: the window IS the answer (log between ts) |
| `ListNoteRange`     | "Showing X–Y of Z." (always)       | Cursor-paginated position: pairs with `LoadMore` batches  |

## Pattern B: page-number pagination

For bounded sets with stable ordering, slice server-side and let
`Pagination` own the links (SEO-friendly `rel` attributes included):

```templ
@display.Table(display.TableProps{Headers: headers, Rows: pageRows})
@navigation.Pagination(navigation.PaginationProps{
    CurrentPage: page,
    TotalPages:  totalPages,
    BaseURL:     "/messages",
    QueryParam:  "page", // default
})
```

Render `ListNote{Variant: ListNoteRange, Shown: len(pageRows), RangeFrom: offset + 1, RangeTo: offset + len(pageRows), Total: total}` when the position
inside the whole set matters more than the page number.

## Pattern C: cursor batches (unbounded / shifting data)

Compose `LoadMore` + `ListNoteRange` — each response appends a batch and
updates the position notice:

```templ
@display.Table(display.TableProps{Headers: headers, Rows: accumulated})
@display.ListNote(display.ListNoteProps{
    Variant:   display.ListNoteRange,
    Shown:     len(accumulated),
    RangeFrom: 1,
    RangeTo:   len(accumulated),
    Total:     totalMatches,
})
@navigation.LoadMore(navigation.LoadMoreProps{
    Label:    "Load more",
    Endpoint: "/messages",
    Cursor:   nextCursor,
})
```

End the list with `navigation.EndOfList` once `nextCursor` is empty.

## What NOT to do

- **Don't render every row "because it's simpler".** A 500-row DOM is a
  performance and accessibility regression that grows silently with data.
- **Don't hand-roll "Showing X of Y" text.** `ListNote` owns the phrasing,
  pluralization, dark mode, and the empty-case suppression — a local
  `fmt.Sprintf` copy drifts.
- **Don't paginate inside a card without `Table.Flush` + `CellPadding`** —
  see [Table Inside Card](table-in-card.md) for the compact composition.

## Related

- [Cursor Pagination](cursor-pagination.md) — the LoadMore round trip in full
- [Custom Table Rows](custom-table-rows.md) — sortable headers for the same table
- `Table.LazyRows` — `content-visibility: auto` for the rare case you truly
  must render hundreds of rows
