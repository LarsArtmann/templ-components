# Cursor-Based Pagination with HTMX

## When to use this

Page-number pagination (`navigation.Pagination`) works for bounded result sets.
**Cursor pagination** is better when:

- The dataset is unbounded or frequently changing (new items insert at the top)
- You need stable ordering under concurrent writes (page numbers shift)
- You're building an infinite-scroll / "load more" UX

## Pattern

```
[item] [item] [item] ... [item]
         [ Load more → ]
```

The `navigation.LoadMore` button issues an `hx-get` to your endpoint with the
cursor. The server responds with the next batch of items **plus a new LoadMore
button** (with the updated cursor). `hx-swap="outerHTML"` replaces the old
button with the new items + button in one round trip.

## Server handler

```go
func itemsHandler(w http.ResponseWriter, r *http.Request) {
    cursor := r.URL.Query().Get("cursor")
    items, nextCursor := fetchItems(cursor, 20)

    // Render items + new LoadMore button
    templ.Handler(itemList(items, nextCursor)).Component().Render(r.Context(), w)
}
```

## Templ template

```templ
templ itemList(items []Item, nextCursor string) {
    for _, item := range items {
        @itemCard(item)
    }
    if nextCursor != "" {
        @navigation.LoadMore(navigation.LoadMoreProps{
            Endpoint: "/api/items",
            Cursor:   nextCursor,
        })
    }
}
```

## How it works

1. First page renders with the initial `LoadMore` button (cursor from first batch)
2. User clicks "Load more" → `hx-get="/api/items?cursor=<cursor>"`
3. Server returns next items + new `LoadMore` button with updated cursor
4. `hx-swap="outerHTML"` replaces old button → new items appear above the new button
5. When `nextCursor == ""` (no more items), no button is rendered — the list ends naturally

## Cursor design

The cursor is opaque to the client. Common approaches:

| Strategy     | Cursor format       | Use case                           |
| ------------ | ------------------- | ---------------------------------- |
| Offset       | `base64(offset=20)` | Simple, stable data                |
| Timestamp+ID | `base64(ts,id)`     | Time-ordered feeds (no page drift) |
| Token        | server-issued token | Encrypted/signed cursors           |

**Never expose raw database IDs or offsets** — base64-encode or encrypt the cursor.

## Infinite scroll variant

Replace the button with an intersection observer sentinel:

```templ
<div id="infinite-sentinel" hx-get={ "/api/items?cursor=" + cursor } hx-trigger="revealed" hx-swap="outerHTML">
</div>
```

HTMX's `revealed` event fires when the sentinel enters the viewport, triggering
the fetch automatically.

## Combining with `display.Grid`

```templ
@display.Grid(display.GridProps{Cols: display.GridCols3}) {
    for _, item := range items {
        @display.Card(display.CardProps{Title: item.Title}) {
            <p>{ item.Description }</p>
        }
    }
    if nextCursor != "" {
        <div class="col-span-full">
            @navigation.LoadMore(navigation.LoadMoreProps{
                Endpoint: "/api/items",
                Cursor:   nextCursor,
            })
        </div>
    }
}
```

## Error handling

Wire `htmx.GlobalErrorHandling` to show a toast when the load fails. The user
can retry by clicking the button again — the cursor hasn't changed.
