# Recipe: Hybrid rendering — templ components into a strings.Builder

**When to read:** your codebase renders HTML from plain Go (a
`strings.Builder`, `io.Writer`, or hand-rolled renderer) and you want to use
library components without converting handlers to `.templ` files. This is
the "hybrid path" — named after the cqrs-htmx dashboardui pattern that
adopted it.

## The one rule that bites

Children slots (`{ children... }`) only populate when a component is invoked
from **another templ template**. Rendered standalone from Go:

```go
var b strings.Builder
_ = Grid(GridProps{Cols: GridCols2}).Render(ctx, &b)
// b contains the grid div — and NO children. The slot is empty.
```

This applies to every children-slot component (`Grid`, `Card`, `Form`,
`htmx.PolledRegion`, …). It is templ runtime behavior, not a library bug:
children travel through the context, and only a templ caller wraps the
context for you. Pinned by `TestGridHybridChildrenEmptyWithoutWithChildren`.

## Three ways to get content in

### 1. Prefer props-driven components (zero ceremony)

Most of the catalogue carries its content in props and renders fully in the
hybrid path: `StatCard`, `Badge`, `ListNote`, `CopyButton`, `Table` (with
`Rows`), `Alert`, `EmptyState`, `Pagination`, `Button`, … If a props-driven
equivalent exists, use it.

### 2. Pass children programmatically: `templ.WithChildren`

The escape hatch for children-slot components — thread children through the
context exactly the way a templ caller would:

```go
child := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
	_, err := io.WriteString(w, "<p>rendered from Go</p>")
	return err
})

var b strings.Builder
ctx := templ.WithChildren(context.Background(), child)
_ = Grid(GridProps{Cols: GridCols2}).Render(ctx, &b)
// b contains the grid div WITH the child inside.
```

Multiple children compose by writing them all in one `ComponentFunc` (or by
composing `templ.Component` values). Pinned by
`TestGridHybridChildrenViaWithChildren`.

### 3. Wrap in a real templ file (recommended for polled regions)

For `htmx.PolledRegion` the children ARE the server-side partial that the
poll endpoint re-serves — render it from a `.templ` file whose children
match the endpoint's response, so the initial render and every poll produce
identical markup. `WithChildren` works here too, but a shared `.templ`
partial is the shape that cannot drift from its endpoint.

## Gotchas

- **Nested children:** `WithChildren` applies to the component you render.
  If that component invokes OTHER children-slot components internally, the
  context flows down — but explicit is better than implicit: compose each
  level yourself.
- **Golden tests use the templ path** (`utils.Render`); hybrid output is
  byte-identical for props-driven components, so goldens transfer.
- **CSP:** components that emit inline scripts (`CopyButton`, `Modal`, …)
  still need `Nonce` in `BaseProps` — the hybrid path changes nothing there.
