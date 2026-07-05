# Recipe: Custom Table Rows

How to render custom `<tr>` elements with full control over cell content,
styling, and HTMX attributes — using the `Table.Body` slot and `TypedHeaders`.

## Basic: Body slot for custom rows

When you need per-row logic (conditional styling, components inside cells,
HTMX attributes), set `Body` to a `templ.Component` that renders `<tr>` elements:

```go
@display.Table(display.TableProps{
    Headers: []string{"User", "Role", "Actions"},
    Body:    userRows(users),
})

templ userRows(users []User) {
    for _, u := range users {
        <tr class="hover:bg-gray-50 dark:hover:bg-gray-800">
            <td class="px-4 py-3">{ u.Name }</td>
            <td class="px-4 py-3">
                @display.Badge(display.BadgeProps{Type: display.BadgePrimary, Children: u.Role})
            </td>
            <td class="px-4 py-3">
                <button hx-delete={ "/users/" + u.ID } hx-target="closest tr" hx-swap="outerHTML">
                    Delete
                </button>
            </td>
        </tr>
    }
}
```

## Sortable columns with TypedHeaders

Use `TypedHeaders` instead of `Headers` for columns that support server-side sorting:

```go
@display.Table(display.TableProps{
    TypedHeaders: []display.TableHeader{
        {Label: "Name", Sortable: true, SortDirection: sortDir, Href: toggleSortURL("name", sortDir)},
        {Label: "Email"},
        {Label: "Joined", Sortable: true, SortDirection: SortNone, Href: "/users?sort=joined"},
    },
    Rows: rows,
})
```

- Sortable headers get `aria-sort="ascending|descending|none"` for screen readers
- When `Href` is set, the label renders as a clickable `<a>` link
- Visual indicators (↑/↓) appear next to the active sort column
- Non-sortable headers render as plain text (no aria-sort)

## Server-side sort handler pattern

```go
func usersHandler(w http.ResponseWriter, r *http.Request) {
    sortCol := r.URL.Query().Get("sort")
    sortDir := display.SortAsc
    if r.URL.Query().Get("dir") == "desc" {
        sortDir = display.SortDesc
    }
    users := db.QueryUsers(r.Context(), sortCol, string(sortDir))
    // render table with TypedHeaders...
}
```
