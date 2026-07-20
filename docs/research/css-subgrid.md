# Research Note: CSS Subgrid — Status and templ-components Implications

**Created:** 2026-07-20
**Status:** Tracking — not yet actionable. Revisit when Baseline status changes.

---

## What is CSS subgrid?

`grid-template-columns: subgrid` (and `grid-template-rows: subgrid`) lets a
nested grid container inherit its parent's grid tracks instead of defining
new ones. This enables **true 2D alignment across nested components** —
something that currently requires either flattening the DOM or duplicating
column widths.

```css
.parent {
  display: grid;
  grid-template-columns: 12rem minmax(0, 1fr);
}
.child {
  display: grid;
  grid-template-columns: subgrid;
  grid-column: span 2;
}
/* .child's two columns now exactly match .parent's — no width duplication */
```

## Current Baseline status (as of 2026-07-20)

- **Chrome/Edge:** shipped (v117, 2023-09)
- **Safari:** shipped (v16, 2022-09)
- **Firefox:** shipped (v71, 2019-12)
- **Baseline:** Widely available (all three engines shipped >2 years ago)

Subgrid IS Baseline. The holdup is not browser support.

## What it would unlock in templ-components

### 1. Card header/body/footer alignment

Currently `Card` has a header, body, and footer that each manage their own
padding and layout independently. With subgrid, the Card could establish a
grid (e.g. `[icon] [title] [action]`) and the header/body/footer would
inherit those tracks — so an icon in the header aligns with content in the
body and an action in the footer.

### 2. DefinitionList term/detail alignment

`DefinitionList` currently uses `grid-cols-[auto_1fr]` which works for a
flat list. With subgrid, a nested component (e.g. `Card` inside a detail
cell) could inherit the parent's term column width, so all term/detail
pairs across multiple cards stay aligned.

### 3. Form field alignment

`Form` Layout: Grid uses `sm:grid-cols-[auto_minmax(0,1fr)]` for label/value
alignment. With subgrid, a `FormFieldWrapper` could inherit the form's
columns, so nested input groups (e.g. `InputGroup` with prefix/suffix) keep
the label column aligned with sibling fields.

### 4. Table-like layouts without `<table>`

Subgrid enables CSS-only "tables" where a header row establishes column
tracks and each data row inherits them — without the HTML `<table>` element
(which carries semantic baggage and has limited styling flexibility).

## Why we haven't adopted it yet

1. **No clear consumer demand.** The current `minmax(0, 1fr)` + `auto`
   patterns handle the 80% case. Subgrid's value is in cross-component
   alignment, which is a nice-to-have, not a blocker.
2. **API complexity.** Exposing subgrid via a typed API (e.g.
   `Card.Subgrid: true`) adds surface area. The raw CSS property is simpler
   but harder to reason about across the component boundary.
3. **Testing gap.** Golden-file tests compare rendered HTML, not visual
   alignment. Subgrid's benefit is visual; we'd need screenshot testing to
   verify it works. That's a separate initiative.

## When to revisit

- A consumer requests cross-component column alignment that subgrid would solve
- The library adds screenshot/visual regression testing
- A new component (e.g. a settings form with mixed input types) would benefit
  from label alignment across nested wrappers

## Decision

**Track only.** No action until a concrete use case surfaces. The primitives
shipped in this plan (`AppShell`, `Container`, `Split`, `Stack`, multi-col
`Footer`, `Form` Grid layout) do NOT need subgrid — `minmax(0, 1fr)` is
sufficient for their 2D layouts.
