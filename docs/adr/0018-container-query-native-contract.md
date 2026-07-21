# ADR 0018: Container-Query-Native Contract

## Date

2026-07-21

## Status

Accepted

## Context

The library has one container-query component (`display.Grid` with `ContainerResponsive: true`,
shipped in v0.16.0). Every other responsive component uses viewport breakpoints (`sm:`, `md:`,
`lg:`, `xl:`). This is the right default — most components live at the top level of the page and
the viewport is the relevant size — but it breaks down for components placed inside constrained
containers: a `Nav` rendered inside a 320px sidebar shouldn't wait for `lg:viewport` to collapse,
and a `Card` inside a 4-up grid shouldn't pad as if it were full-width.

CSS Container Queries (Baseline 2023) solve this: a component can ask "how wide is my parent?"
and respond accordingly, independent of the viewport. The pattern is established in
`display/grid.templ` via the `ContainerResponsive bool` flag and the `@container` wrapper
emission.

This ADR codifies the contract for extending that pattern to other components.

## Decision

Establish a uniform `ContainerAware bool` opt-in flag pattern on components whose responsive
behavior is plausibly container-dependent.

### The `ContainerAware` flag contract

1. **Opt-in, default false.** Existing viewport-based behavior is byte-identical when the flag
   is off. No golden files change unless the consumer opts in.
2. **Flag name:** `ContainerAware` (NOT `ContainerResponsive`) — shorter, reads better at the
   call site (`CardProps{ContainerAware: true}`). The Grid precedent (`ContainerResponsive`)
   stays as-is for backward compat; we won't rename it pre-v1.0.
3. **Wrapper emission:** when true, the component emits its root inside a `<div class="@container">`
   wrapper. The wrapper is **not** the component's root element (so the consumer's `props.ID`,
   `props.Class`, `props.Attrs` still land on the actual component root, not the wrapper).
4. **Variant class swap:** when true, viewport-keyed classes (`sm:`, `lg:`) inside the component
   swap to container-keyed classes (`@sm:`, `@lg:`). Implementation: a small lookup map
   (e.g. `navContainerClass` alongside `navClass`) selected via `utils.Ternary`.
5. **No layout regression when off:** the lookup map for the container variant is unused code
   when the flag is false. The viewport-based classes are the default and remain unchanged.

### When to add `ContainerAware`

Add it to a component when **all** of:

- The component has viewport-based responsive behavior today (`sm:`, `lg:` etc.).
- The component is plausibly placed inside a constrained container (sidebar, card body, grid
  cell, drawer).
- A clear, named behavior change exists at container-width boundaries (collapse-to-burger,
  compact-padding, hide-secondary-actions).

Do NOT add it when:

- The component is always full-width in practice (e.g. `layout.Page` shell).
- The component has no responsive behavior today (nothing to container-ize).
- The behavior change would be cosmetic only (just colors or font sizes).

### Test approach

Container queries are CSS-only — they don't affect the SSR HTML structure beyond the wrapper
emission and the class swap. Both are unit-testable without a browser:

- Assert `<div class="@container">` wrapper is emitted when `ContainerAware: true`.
- Assert the wrapper is absent when false.
- Assert the container variant class (e.g. `@lg:`) is present when true.
- Assert the viewport variant class (e.g. `lg:`) is present when false.
- Wrap in a fixed-width parent for visual regression if needed (deferred — no Playwright today).

### RTL behavior

Container queries are direction-agnostic: `@container` measures inline-size (width), not
block-size, and logical CSS properties (`ms-`, `me-`, `start-`, `end-`) continue to mirror
automatically under `dir="rtl"`. No special handling required.

### Interaction with `AppShell`

`layout.AppShell` (v0.19.0) creates a 2D grid with a sidebar column and a main column. The
sidebar is `hidden lg:block`; the main column is `minmax(0,1fr)`. Children of either column
that opt into `ContainerAware` respond to their column's width, not the viewport. This is the
primary intended use case.

### Initial candidates

| Component | Container-aware behavior                                           |
| --------- | ------------------------------------------------------------------ |
| `Nav`     | Collapse to burger below `@sm:` instead of `sm:`                   |
| `Card`    | Compact padding (px-3 py-2) below `@sm:` instead of normal padding |

Both ship in this phase (v0.21.0). `Form.Inline` is a possible future candidate (decision
deferred until a clear container-driven behavior emerges).

## Consequences

**Positive:**

- Components become placement-aware — the same component works correctly in a sidebar, a
  card body, or a full-width page without consumer-side overrides.
- The pattern is uniform across the library: one flag name, one wrapper pattern, one test
  approach.
- Container Queries are Baseline 2023 — no fallback strategy needed for any browser still
  receiving security updates.
- No JS, no runtime cost — pure CSS via Tailwind v4 `@container` variants.

**Negative:**

- Each container-aware component carries two parallel class lookup maps (viewport +
  container). This is ~6 extra lines per component — acceptable cost for the capability.
- Consumers who don't opt in get no benefit but pay a tiny cognitive cost (one more bool flag
  in the props struct). Mitigated by clear godoc and the "default off" rule.

## References

- [ADR-0016: Grid-first for 2D layouts](0016-grid-first-for-2d-layouts.md)
- [CSS Container Queries on MDN](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_containment/Container_queries)
- [Tailwind v4 Container Queries](https://tailwindcss.com/docs/container-queries)
- [Baseline: Container Queries](https://web.dev/articles/container-query) — 2023
