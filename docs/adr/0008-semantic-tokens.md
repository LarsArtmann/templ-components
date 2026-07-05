# ADR 0008: Semantic Token Layer

## Status

Proposed — deferred to v1.0

## Context

The library currently uses 256+ hardcoded Tailwind color references (`bg-blue-600`,
`text-red-500`, `border-gray-300`, etc.) across all components. Consumers override
colors by setting `@theme { --color-blue-600: #...; }` in their CSS — this works but
couples the consumer's theme to Tailwind's color palette names.

The problem: a consumer wanting a "primary" brand color must re-define `--color-blue-600`
to their brand color, which is semantically wrong (they're not making blue look
different — they're saying "my primary color is blue, but could be green"). There's no
indirection between the **semantic intent** (primary, danger, success) and the **color
value** (blue-600, red-500, green-500).

## Decision

Introduce an **opt-in semantic token layer** that aliases Tailwind colors to semantic
names:

```css
/* templ-components-theme.css (consumer-side, opt-in) */
@theme {
    --color-tc-primary: var(--color-blue-600);
    --color-tc-primary-hover: var(--color-blue-500);
    --color-tc-danger: var(--color-red-600);
    --color-tc-danger-hover: var(--color-red-500);
    --color-tc-success: var(--color-green-600);
    --color-tc-warning: var(--color-yellow-500);
    --color-tc-info: var(--color-blue-500);
}
```

### Migration path (opt-in, then flip)

1. **Phase 1 (v0.9.0):** Document the semantic token pattern. `templ-components-theme.css`
   already has examples. No code changes.
2. **Phase 2 (v1.0):** Add a build-time flag or CSS layer so components emit
   `bg-tc-primary` instead of `bg-blue-600` when the semantic layer is active.
3. **Phase 3 (post-v1.0):** Make `bg-tc-primary` the default. `bg-blue-600` becomes
   the override.

### Why not do it now?

- 256 color references across all `.templ` files = massive golden file churn
- Each component needs review to determine the correct semantic name
- Risk of breaking consumer themes that already override via `--color-blue-600`
- No consumer has requested this yet (the CSS-variable override model works)

## Alternatives Considered

### A. CSS custom properties directly (`--tc-primary`)

Use `style="background: var(--tc-primary)"` instead of Tailwind classes.

**Rejected:** Loses Tailwind's utility-class ergonomics (hover states, dark mode,
responsive variants all require explicit CSS).

### B. Component-level `Variant`/`Color` props

Add a `Color string` field to every component.

**Rejected:** Defeats the CSS-variable model. 86 components × N variants = massive
API surface for something CSS already handles better.

### C. Tailwind v4 `@theme` aliases only (current approach)

Consumers define `--color-blue-600` in their `@theme` block.

**Accepted as interim.** Works today, no code changes, but couples semantic intent
to Tailwind color names. The semantic token layer is the evolution of this approach.

## Consequences

- **Deferred to v1.0:** No immediate code changes. The audit is documented for planning.
- **Consumer impact:** Zero until v1.0. Existing `@theme` overrides continue to work.
- **Migration effort:** ~4 hours of mechanical find-and-replace + golden file updates
  when the time comes.
