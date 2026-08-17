# Container Queries

> Use `@container` when a component should adapt to its **parent's** width, not the viewport.

**See also:** [Container Query Leveraging Strategy](../container-query-strategy.md) | [Fluid Typography](fluid-typography.md) | [ADR-0018](../adr/0018-container-query-native-contract.md)

## When to use container queries

Use container queries when a component is placed in variable-width containers
(sidebar, split pane, card body) and must reflow based on available space — not
the browser window.

Use viewport queries (`sm:`, `lg:`) when the component spans the full page width.

## Grid with ContainerAware

The `Grid` component defaults to container-query mode (`ContainerAware: true` since v2.0):

```go
@display.Grid(display.GridProps{
    Cols:            display.GridCols3,
    ContainerAware:  true, // default since v2.0; set false for viewport breakpoints
    Class:           "max-w-2xl",
}) {
    // cards...
}
```

When `ContainerAware` is `true`:

- The grid wraps in a `<div class="@container">` element
- Column counts use `@sm:` / `@md:` / `@lg:` variants instead of `sm:` / `lg:`
- The grid responds to the wrapper's width, not the viewport

When `false`, the grid uses standard viewport breakpoints.

## Manual container queries in custom components

```html
<div class="@container">
  <div class="grid grid-cols-1 @sm:grid-cols-2 @lg:grid-cols-3">
    <!-- columns reflow based on parent width -->
  </div>
</div>
```

Tailwind v4 supports `@container` natively — no plugin or config needed.

## All container-aware components

Every component with a `ContainerAware` flag follows the
same contract (ADR-0018): emits `@container` wrapper or root
class, swaps viewport breakpoints (`sm:`/`md:`/`lg:`) for container variants
(`@sm:`/`@md:`/`@lg:`). Since v2.0, `Grid`, `Card`, and `Split` default `true`;
the other 5 default `false`.

| Component                   | Flag             | What adapts                                         | Default |
| --------------------------- | ---------------- | --------------------------------------------------- | ------- |
| `display.Grid`              | `ContainerAware` | Column count (1→2→3→N)                              | `true`  |
| `display.Card`              | `ContainerAware` | Padding (compact below `@sm:`)                      | `true`  |
| `display.DefinitionGrid`    | `ContainerAware` | Term-detail card grid column count                  | `false` |
| `navigation.Nav`            | `ContainerAware` | Collapse to hamburger below `@sm:`                  | `false` |
| `navigation.Pagination`     | `ContainerAware` | Mobile prev/next vs full page numbers               | `false` |
| `layout.Split`              | `ContainerAware` | 2-col main+aside collapses to stacked below `@md:`  | `true`  |
| `forms.Form`                | `ContainerAware` | Grid layout label/value columns below `@sm:`        | `false` |
| `feedback.SkeletonCardGrid` | `ContainerAware` | Loading skeleton grid matches `Grid.ContainerAware` | `false` |

### Split — article+sidebar in a constrained container

```go
@layout.Split(layout.SplitProps{
    Main:           articleBody,
    Aside:          tocWidget,
    ContainerAware: true,
})
```

### Form — settings form in a modal or drawer

```go
@forms.Form(forms.FormProps{
    Layout:         forms.FormLayoutGrid,
    ContainerAware: true,
}) {
    @forms.Input(forms.InputProps{Name: "email"})
}
```

### Pagination — in a card footer

```go
@navigation.Pagination(navigation.PaginationProps{
    CurrentPage:   3,
    TotalPages:    10,
    BaseURL:       "/users",
    ContainerAware: true,
})
```

### SkeletonCardGrid — loading state for a container-aware grid

```go
@feedback.SkeletonCardGrid(feedback.SkeletonCardGridProps{
    Count:          6,
    ContainerAware: true,
})
```

## Container query size reference

| Variant | Min width     | Equivalent viewport |
| ------- | ------------- | ------------------- |
| `@sm:`  | 24rem (384px) | `sm:` (640px)       |
| `@md:`  | 28rem (448px) | —                   |
| `@lg:`  | 32rem (512px) | `md:` (768px)       |
| `@xl:`  | 36rem (576px) | `lg:` (1024px)      |
| `@2xl:` | 42rem (672px) | `xl:` (1280px)      |
| `@3xl:` | 48rem (768px) | —                   |
| `@4xl:` | 56rem (896px) | —                   |

Container breakpoints are **smaller** than viewport breakpoints because containers
are typically narrower than the full window.
