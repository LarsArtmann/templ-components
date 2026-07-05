# Container Queries

> Use `@container` when a component should adapt to its **parent's** width, not the viewport.

## When to use container queries

Use container queries when a component is placed in variable-width containers
(sidebar, split pane, card body) and must reflow based on available space — not
the browser window.

Use viewport queries (`sm:`, `lg:`) when the component spans the full page width.

## Grid with ContainerResponsive

The `Grid` component supports an opt-in container-query mode:

```go
@display.Grid(display.GridProps{
    Cols:               display.GridCols3,
    ContainerResponsive: true,
    Class:              "max-w-2xl",
}) {
    // cards...
}
```

When `ContainerResponsive` is `true`:

- The grid wraps in a `<div class="@container">` element
- Column counts use `@sm:` / `@md:` / `@lg:` variants instead of `sm:` / `lg:`
- The grid responds to the wrapper's width, not the viewport

When `false` (default), the grid uses standard viewport breakpoints — backward
compatible with all existing consumers.

## Manual container queries in custom components

```html
<div class="@container">
  <div class="grid grid-cols-1 @sm:grid-cols-2 @lg:grid-cols-3">
    <!-- columns reflow based on parent width -->
  </div>
</div>
```

Tailwind v4 supports `@container` natively — no plugin or config needed.

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
