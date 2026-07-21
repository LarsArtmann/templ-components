# Popover API + CSS Anchor Positioning — Baseline Research

**Researched:** 2026-07-21 · **For:** ADR-0017 (Popover API migration) · **Baseline today:** 2026-07-20

## Browser Support Matrix

### HTML `popover` attribute — ✅ Baseline 2024 (April 2024)

| Browser | First supported | Notes                    |
| ------- | --------------- | ------------------------ |
| Chrome  | 114             | May 2023                 |
| Edge    | 114             |                          |
| Safari  | 17.0            | Sep 2023 (desktop + iOS) |
| Firefox | 125             | Apr 2024                 |
| Samsung | 23              |                          |

**Verdict:** Safe to use as a baseline today. The full Popover API (event surface, `hint` state)
reached Baseline 2025 (January 2025). For our purposes (`auto` + `manual` + `popovertarget`),
2024 Baseline applies.

### `popovertarget` / `popovertargetaction` — ✅ same as `popover`

Same browser versions. Works only on `<button>` and `<input type="button">`.

### CSS Anchor Positioning — ✅ Baseline 2026 (January 2026)

| Property / Function | Chrome | Edge  | Firefox | Safari | Status                                                   |
| ------------------- | ------ | ----- | ------- | ------ | -------------------------------------------------------- |
| `anchor-name`       | 125    | 125   | 147     | 26     | Baseline 2026                                            |
| `position-area`     | 129    | 129   | 147     | 26     | Baseline 2026                                            |
| `position-try`      | 125    | 125   | 147     | 26     | Baseline 2026                                            |
| `anchor()`          | 125    | 125   | 147     | 26     | Baseline 2026                                            |
| `anchor-size()`     | 125    | 125   | 147     | 26     | Baseline 2026                                            |
| `position-anchor`   | 125\*  | 125\* | 147\*   | 26\*   | Limited (initial-value bugs; \* = full only in 151+/27+) |

**Verdict:** Use `position-area` + `anchor-name` (Baseline 2026); keep CSS `inset` fallback for older browser support inside the same file.

### `popover="hint"`

Newer third value, narrower support (Chrome 133+, Safari 26+, Firefox behind flag as of mid-2026). Ideal for tooltips (light-dismiss + doesn't close other `auto` popovers). Falls back to `manual` in unsupporting browsers — safe to use as progressive enhancement.

## Light-Dismiss Behavior

`popover="auto"` (and `hint`) get **free** light-dismiss:

- Click outside the popover
- `Escape` key
- Opening another `auto` popover closes the current one (unless nested)
- `popover="manual"` does NOT light-dismiss — must be closed explicitly

## Current-State Inventory (what JS actually exists today)

Honest audit of the 5 components named in the plan:

| Component   | Has singleton JS? | JS location                  | Approx JS LOC | Purpose of JS                                                     |
| ----------- | ----------------- | ---------------------------- | ------------- | ----------------------------------------------------------------- |
| Tooltip     | Yes               | `display/shared.go`          | ~20           | aria-describedby propagation; click-toggle for touch; Escape      |
| HoverCard   | **No**            | —                            | 0             | Already pure CSS (`:hover` / `:focus-within`)                     |
| ContextMenu | Yes (inline)      | `display/context_menu.templ` | ~12           | `contextmenu` event positioning; click-outside; Escape            |
| Popover     | Yes (inline)      | `display/popover.templ`      | ~40           | click-toggle; click-outside; Escape; open/close w/ focus          |
| Dropdown    | Yes (inline)      | `display/dropdown.templ`     | ~50           | click-toggle; click-outside; Escape; open/close; ArrowUp/Down nav |

**Total today:** ~120 lines of JS across 4 components (HoverCard is already native).

## Fallback Strategy Decision

**Decision: progressive enhancement with CSS `inset` fallback.**

- Popover API: ship as primary (Baseline 2024)
- Anchor Positioning: use where it simplifies positioning (Baseline 2026), but the existing CSS-class-based positioning (e.g. `bottom-full left-1/2 -translate-x-1/2`) already works without it. Keep CSS positioning; add Anchor as enhancement only if a clear win emerges.
- For browsers without `popover` support: element renders inline (no top-layer) — graceful degradation, not breakage.

**Rejected alternatives:**

1. Tiny JS positioner — would replicate what popover API does natively; no win.
2. Anchor-only positioning — limits browser support unnecessarily; current CSS positioning works fine.

## Migration Plan (revised from roadmap)

The roadmap assumed "5 migrations". The honest reality:

| Component   | Migration action                                                                                                                                                                                                                                                                  | LOC delta |
| ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------- |
| Tooltip     | **De-JS entirely.** Remove singleton script. Document touch limitation (tooltips are hover/focus progressive enhancement only). The `aria-describedby`-propagation problem becomes a documented consumer responsibility (set `aria-describedby` on the trigger element directly). | -20       |
| HoverCard   | **No change.** Already native CSS. Document as the reference implementation.                                                                                                                                                                                                      | 0         |
| ContextMenu | **Full migration.** Add `popover="auto"`; replace JS with ~6-line `contextmenu` event handler that calls `showPopover()` + sets `inset` to `clientX/clientY`. Native light-dismiss + Escape.                                                                                      | -6        |
| Popover     | **Full migration.** Use `popovertarget` on trigger button, `popover="auto"` on panel. **Zero JS.**                                                                                                                                                                                | -40       |
| Dropdown    | **Partial migration.** Use `popovertarget` on trigger button, `popover="auto"` on menu. Keep ~25 lines JS for ArrowUp/Down/Home/End keyboard nav. Native dismiss/Escape/click-outside.                                                                                            | -25       |

**Net:** ~120 lines JS → ~31 lines JS. The 5 inline scripts in `display/shared.go` and `*.templ` shrink to 2 (Dropdown keyboard nav + ContextMenu contextmenu handler).

## Reuse of existing `[popover]::backdrop` CSS

`templates/custom.css:26` has:

```css
[popover]::backdrop {
  background-color: transparent;
}
```

This is already correct for our use — we don't want a visible backdrop on dropdowns/menus. Reuse as-is. Add component-specific entrance animations via `@starting-style` + `allow-discrete` (same pattern as the `<dialog>` overlay system at `custom.css:30-200`).

## Per-Component `popover` mode selection

| Component   | `popover` value | Why                                                         |
| ----------- | --------------- | ----------------------------------------------------------- |
| Tooltip     | (none)          | Removed — pure CSS hover/focus is the right model           |
| HoverCard   | (none)          | Pure CSS hover/focus works (already shipped)                |
| ContextMenu | `auto`          | Light-dismiss on any click — correct for right-click menus  |
| Popover     | `auto`          | Click-triggered, dismiss-on-outside-click — textbook `auto` |
| Dropdown    | `auto`          | Same as Popover — action menu semantics                     |

No component uses `popover="manual"` — that value is for toasts / persistent notifications, which we already handle differently.

## Risks & Mitigations

| Risk                                                | Mitigation                                                          |
| --------------------------------------------------- | ------------------------------------------------------------------- |
| Keyboard nav regression in Dropdown                 | Preserve thin JS handler for ArrowUp/Down/Home/End only             |
| `popovertarget` only works on `<button>`            | All current triggers are already `<button>` — no change needed      |
| Touch devices lose tooltip hover                    | Documented; tooltips are non-critical progressive enhancement       |
| Anchor Positioning partial bugs (`position-anchor`) | We don't use `position-anchor`; CSS class positioning is unaffected |
| Old browsers render popover content inline          | Graceful degradation — content visible inline, not broken           |

## CSP Impact

- Net **reduction** in inline scripts: Tooltip + Popover scripts deleted entirely.
- CSP nonce test (`integration/csp_nonce_test.go`) currently does NOT cover Tooltip/Popover — add them post-migration as negative assertions (no `<script>` tag emitted).
- Dropdown + ContextMenu still emit scripts — nonce test continues to cover them.

## Sources

- <https://developer.mozilla.org/en-US/docs/Web/HTML/Global_attributes/popover>
- <https://developer.mozilla.org/en-US/docs/Web/API/Popover_API>
- <https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_anchor_positioning>
- <https://caniuse.com/mdn-html_global_attributes_popover>
- <https://caniuse.com/css-anchor-positioning>
