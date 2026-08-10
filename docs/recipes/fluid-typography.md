# Fluid Typography via Container Query Units

> Size text relative to its **container's** width, not the viewport. A heading in a
> 320px sidebar shrinks smoothly; the same heading in a full-width hero grows.

## Why container query units?

Standard responsive typography uses viewport units (`vw`) or breakpoint font sizes
(`text-2xl md:text-4xl`). Both respond to the **browser window**. When a component lives
inside a constrained container — a sidebar, a grid cell, a card body, a drawer — viewport
units produce the wrong size: a sidebar heading renders huge because the window is wide,
even though the sidebar is narrow.

CSS **Container Query Length Units** (`cqi`, `cqw`, `cqb`, `cqh`, `cqmin`, `cqmax`) solve
this. They resolve against the nearest `@container` ancestor's size:

- `1cqi` = 1% of the container's **inline size** (width in horizontal writing mode) — the
  most useful unit for fluid type.
- `1cqw` = 1% of container width (same as `cqi` in LTR/RTL horizontal scripts).

Pair `cqi` with `clamp()` for a smooth scale with a sane floor and ceiling:

```css
font-size: clamp(1.5rem, 3.5cqi + 0.5rem, 2rem);
/*  → 24px minimum, grows with container width, never exceeds 32px  */
```

Baseline: container query units are **Baseline 2023** (same as container queries
themselves). No fallback strategy needed for any browser receiving security updates.

## The `.tc-fluid-*` utility classes

The library ships ready-made fluid heading and lead-text classes in `templates/custom.css`.
They compose directly inside any container-aware component (Card, Split, Nav, Form, Grid,
DefinitionGrid, Pagination, SkeletonCardGrid) or any `<div class="@container">`:

| Class              | Min (floor) | Max (ceiling) | Scale formula          |
| ------------------ | ----------- | ------------- | ---------------------- |
| `.tc-fluid-display` | 2rem (32px) | 3.5rem (56px) | `5cqi + 1rem`         |
| `.tc-fluid-h1`      | 1.75rem     | 2.5rem        | `4cqi + 0.75rem`      |
| `.tc-fluid-h2`      | 1.5rem      | 2rem          | `3.5cqi + 0.5rem`     |
| `.tc-fluid-h3`      | 1.25rem     | 1.75rem       | `2.5cqi + 0.5rem`     |
| `.tc-fluid-h4`      | 1.125rem    | 1.5rem        | `2cqi + 0.375rem`     |
| `.tc-fluid-lead`    | 1.125rem    | 1.375rem      | `1.75cqi + 0.5rem`    |

> **Requirement:** the element must have a query-container ancestor. Tailwind's
> `@container` class sets `container-type: inline-size`. Without it, `cqi` resolves
> against the **small viewport** (a graceful but viewport-driven fallback).

## Usage

### Inside a container-aware Card

```go
@display.Card(display.CardProps{
    Title:          "Quarterly Revenue",
    ContainerAware: true,  // emits <div class="@container">
}) {
    <p class="tc-fluid-lead text-gray-600 dark:text-gray-400">
        Revenue grew across all regions this quarter.
    </p>
    <span class="tc-fluid-display font-bold text-blue-600 dark:text-blue-400">
        $4.2M
    </span>
}
```

The metric scales down in a 4-up grid and up when the card spans full width — with zero
JavaScript and zero Go logic per breakpoint.

### Inside a Split (article + sidebar)

```go
@layout.Split(layout.SplitProps{
    Main:           articleBody,
    Aside:          tocWidget,
    ContainerAware: true,
})
```

A `.tc-fluid-h2` inside `Aside` renders smaller than the same class inside `Main`, because
each column establishes (or lives within) its own container context.

### Standalone — add `@container` yourself

```html
<div class="@container max-w-md">
    <h2 class="tc-fluid-h2 font-semibold">Heading</h2>
    <p class="tc-fluid-lead">Lead paragraph that scales with the 28rem container.</p>
</div>
```

### Composing with Tailwind utilities

The `.tc-fluid-*` classes only set `font-size` and `line-height`. Compose freely with
Tailwind for weight, color, tracking, and spacing:

```html
<h1 class="tc-fluid-h1 font-bold tracking-tight text-gray-900 dark:text-white">
    Dashboard
</h1>
```

## Rolling your own scale

For a custom fluid size, use `clamp()` + `cqi` directly in your consumer CSS or a Tailwind
arbitrary value:

```html
<!-- Inline arbitrary value -->
<span class="text-[clamp(0.875rem,2cqi+0.25rem,1.25rem)]">Responsive label</span>
```

```css
/* In your custom.css */
.metric-value {
  font-size: clamp(2rem, 8cqi, 5rem);
  font-weight: 800;
  line-height: 1;
}
```

## Container query unit reference

| Unit    | Resolves against                         |
| ------- | ---------------------------------------- |
| `cqw`   | 1% of container width                    |
| `cqh`   | 1% of container height                   |
| `cqi`   | 1% of container inline size (= width in LTR/RTL horizontal text) |
| `cqb`   | 1% of container block size (= height in horizontal text) |
| `cqmin` | 1% of the smaller of container width/height |
| `cqmax` | 1% of the larger of container width/height  |

For typography, **always use `cqi`** (inline size): it mirrors correctly in vertical
writing modes and is the canonical choice for horizontal text scaling.

## Accessibility

- Fluid `cqi` scaling with `clamp()` floors never produces text smaller than the floor, so
  minimum readable sizes are guaranteed.
- These classes set **no animation** — text reflows instantly on resize. There is no
  `motion-reduce` concern (nothing is animated).
- The library's global `text-wrap: balance` (on all headings) composes automatically with
  these classes for clean line breaks.

## See also

- [Container Queries](container-queries.md) — the 8 container-aware components and the
  `@container` wrapper contract (ADR-0018)
- [ADR-0018: Container-Query-Native Contract](../adr/0018-container-query-native-contract.md)
- [Tailwind v4 Container Queries](https://tailwindcss.com/docs/container-queries)
