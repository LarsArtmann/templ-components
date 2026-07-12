# ADR 0015: Stylable Select API (`appearance: base-select`)

## Date

2026-07-12

## Status

Accepted

## Context

The native `<select>` element has historically been impossible to style
consistently across browsers. The `<option>` elements, the dropdown picker,
and the arrow indicator all render with browser-default styling that cannot be
overridden with CSS.

Developers have worked around this for years with:

1. Custom JS dropdown components (select2, Choices.js, Combobox)
2. `appearance: none` hacks (hide native UI, overlay custom div)
3. Library-specific select replacements (Radix Select, MUI Select)

These approaches all sacrifice native accessibility, keyboard navigation, or
mobile UX to achieve visual consistency.

The CSS Working Group and browser vendors have now standardized the
**customizable select API** via `appearance: base-select` (CSS) + `<button>`
and `<selectedcontent>` elements (HTML). This enables full CSS styling of every
part of the select while preserving native semantics.

**Browser support (as of 2026-07):**

| Browser    | Version | Status        |
| ---------- | ------- | ------------- |
| Chrome     | 135+    | Supported     |
| Edge       | 135+    | Supported     |
| Safari     | 27+     | Supported     |
| Firefox    | —       | Not supported |
| iOS Safari | —       | Not supported |

This is **not yet Baseline** (Firefox and iOS Safari missing).

## Decision

Add `SelectProps.Stylable bool` as an **opt-in** field. When `true`:

1. The `<select>` gets `class="tc-select"` (activates `appearance: base-select`)
2. A `<button><selectedcontent></selectedcontent></button>` structure is emitted
   inside the select — the browser clones the selected option's content into
   `<selectedcontent>` for display in the closed state
3. CSS in `templates/app.css` (126 lines, guarded by
   `@supports (appearance: base-select)`) styles:
   - The button (background, border, padding, hover, focus, disabled)
   - The dropdown picker (`::picker(select)` — container, shadow, max-height)
   - The arrow icon (`::picker-icon` — rotates on open)
   - Options (padding, hover, checked states)
   - Full dark mode support

When `false` (default), the select renders as a normal native `<select>` with
no structural changes.

## Consequences

**Positive:**

- Full CSS control over select appearance without sacrificing accessibility
- Progressive enhancement: non-supporting browsers get native `<select>` (no
  broken UI)
- No JavaScript required (pure CSS + HTML)
- Matches the `<selectedcontent>` spec pattern — future-proof

**Negative:**

- Opt-in means most consumers won't discover it unless documented
- The 126-line CSS block is the largest per-component CSS in `app.css`
- Consumers must copy the `.tc-select` CSS from `app.css` to use it
- No Firefox/iOS Safari support means inconsistent appearance across browsers
  for the same consumer app

## Alternatives considered

### Make Stylable the default

Rejected. `appearance: base-select` is not Baseline (no Firefox/iOS Safari).
Making it default would change the appearance for consumers on those browsers
in unexpected ways. The opt-in approach lets consumers choose progressive
enhancement when their browser matrix allows it.

### Custom JS dropdown

Rejected. The library already has `Combobox` for autocomplete use cases.
A JS-based styled select would duplicate that complexity and sacrifice native
keyboard navigation and mobile UX.

### Wait for Baseline

Considered. Rejected because the API is stable in Chrome/Safari/Edge, and
progressive enhancement means non-supporting browsers get a working native
select. The opt-in nature means no consumer is forced to use it.
