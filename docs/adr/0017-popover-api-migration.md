# ADR 0017: Popover API migration for Dropdown, Popover, ContextMenu

## Date

2026-07-21

## Status

Accepted (revised 2026-07-21 — see Revision below)

## Revision (2026-07-21)

The original ADR claimed that CSS class-based positioning (`top-full left-1/2`,
etc.) "continues to work because the popover element remains a DOM descendant
of the relatively-positioned trigger wrapper." **This was wrong.** The Popover
API promotes the panel to the **top layer**, where the UA stylesheet forces
`position: fixed; inset: 0` — the panel is detached from its trigger's DOM
subtree, so CSS classes resolve against the viewport, not the trigger.

The fix: a shared singleton script (`popoverPositionJS` in `display/shared.go`)
reads the trigger's `getBoundingClientRect()` on every `toggle` open and sets
`style.left/top` with viewport clamping. This mirrors the proven
`ContextMenu` cursor-positioning pattern. Popover and Dropdown now emit this
positioner (Popover: ~30 lines shared; Dropdown: positioner + ~25 lines
keyboard nav).

Additionally, the `Tooltip` `aria-describedby` propagation was restored via a
small singleton script (`tooltipAriaJS`) rather than left as consumer
responsibility — the regression was not acceptable for existing consumers.

## Context

Three overlay components — `Dropdown`, `Popover`, `ContextMenu` — currently ship with inline
`<script>` blocks that reimplement browser-native popover behavior in JavaScript:

- Click-triggered toggle on the trigger button
- Click-outside-to-close
- `Escape`-to-close
- Open/close state via `classList.add('hidden')` / `setAttribute('aria-expanded', ...)`
- Manual focus management on open

Combined, these three components carry ~100 lines of singleton JS that the browser now provides
for free via the HTML `popover` attribute (Baseline 2024, April 2024).

A fourth component, `Tooltip`, carries a smaller singleton (~20 lines) in `display/shared.go`
for aria-describedby propagation, touch click-toggle, and Escape-dismiss. A fifth, `HoverCard`,
already uses pure CSS (`:hover` / `:focus-within`) with **no JavaScript** — it is the reference
implementation for the CSS-only approach.

The full browser-support matrix and fallback strategy is documented in
[`docs/research/popover-api.md`](../research/popover-api.md).

## Decision

Migrate `Dropdown`, `Popover`, and `ContextMenu` to the native Popover API. Remove the
`Tooltip` singleton script entirely. Leave `HoverCard` untouched.

### Per-component mode

| Component   | `popover` value          | Trigger mechanism                                     | JS remaining                                                            |
| ----------- | ------------------------ | ----------------------------------------------------- | ----------------------------------------------------------------------- |
| Tooltip     | _(none — CSS show/hide)_ | `:hover` / `:focus-within` (existing CSS)             | ~10 (`tooltipAriaJS` — aria-describedby propagation to focusable child) |
| HoverCard   | _(none — unchanged)_     | `:hover` / `:focus-within` (existing CSS)             | 0                                                                       |
| ContextMenu | `popover="auto"`         | `contextmenu` event → `showPopover()` + `style.inset` | ~6 (cursor positioning, was ~12)                                        |
| Popover     | `popover="auto"`         | `<button popovertarget="id">` (native click-toggle)   | ~30 (`popoverPositionJS` — top-layer positioning, shared singleton)     |
| Dropdown    | `popover="auto"`         | `<button popovertarget="id">` (native click-toggle)   | ~25 keyboard nav + shared positioner (was ~50)                          |

**Why `popover="auto"` (not `manual`) for all three:** light-dismiss (click-outside, Escape)
is the correct UX for menus and click-triggered floating panels. `popover="manual"` is for
persistent notifications / toasts, which these are not.

### Trigger pattern

For `Popover` and `Dropdown`, the trigger is a `<button>` — the native `popovertarget` attribute
turns it into a popover invoker with no JS:

```html
<button popovertarget="my-popover">Open</button>
<div id="my-popover" popover="auto">…content…</div>
```

The browser handles click-to-toggle, `aria-expanded` semantics (via `:popover-open`), light-dismiss,
and Escape.

For `ContextMenu`, there is no declarative right-click trigger, so a minimal per-instance IIFE
calls `showPopover()` and sets the menu's `inset` to `event.clientX/clientY`:

```html
<div data-tc-ctxmenu-trigger="cm1">…right-click target…</div>
<div id="cm1" popover="auto" role="menu">…items…</div>
<script nonce="…">
  (function () {
    var trigger = document.currentScript.previousElementSibling.previousElementSibling;
    // …or use document.querySelector; per-instance ~6 lines
  })();
</script>
```

(Singleton pattern preserved — one `tcCtxMenuAttached` guard, multiple triggers via attribute
delegation. See implementation in `display/context_menu.templ`.)

### Positioning strategy

The native Popover API promotes the panel to the **top layer**, where the UA
stylesheet forces `position: fixed; inset: 0`. The panel is therefore detached
from its trigger's DOM subtree — CSS classes like `top-full left-1/2` resolve
against the **viewport**, not the trigger, producing a panel at the wrong
location.

Three positioning approaches were considered:

1. **CSS Anchor Positioning** (`anchor-name` + `position-area`): the declarative
   ideal, but Baseline 2026 (Chrome-only as of 2026). Too new for a library
   targeting all evergreen browsers.
2. **JavaScript `getBoundingClientRect()` positioning**: read the trigger rect on
   `toggle` open, set `style.left/top`, clamp to viewport. Works in every
   browser that supports the Popover API (Baseline 2024). **Adopted.**
3. **Hybrid (Anchor + `@supports` JS fallback)**: forward-looking but doubles
   the maintenance surface. Deferred until Anchor Positioning reaches
   Baseline-wide.

The shared positioner (`popoverPositionJS`) is a singleton in
`display/shared.go`. It handles four positions (top/bottom/left/right) and
three alignments (start/center/end) with viewport clamping and scroll/resize
repositioning. `ContextMenu` uses a separate cursor-positioning path
(`clientX/clientY`) because it anchors to the mouse, not an element.

### Fallback for older browsers

`popover="auto"` gracefully degrades: in browsers without Popover API support, the element
renders inline (visible, no top-layer). This is a worse UX but is not broken. We accept this
trade-off because the Baseline (April 2024) is now ~2 years old.

### CSP implications

- Dropdown and Popover emit a shared positioner `<script nonce>` (singleton,
  injected once per page).
- Tooltip emits an aria-describedby propagation `<script nonce>` (singleton).
- ContextMenu still emits its cursor-positioning script.
- The `integration/csp_nonce_test.go` tripwire asserts every `<script>` tag
  carries the nonce — verified for all five components.

## Consequences

**Positive:**

- ~40 lines of JS deleted net (the old show/hide/state logic is gone; replaced
  by smaller positioning + aria scripts).
- Native light-dismiss, Escape, focus, and top-layer rendering replace fragile
  custom show/hide JS.
- `Tooltip` show/hide is pure CSS; only the small aria-propagation script
  remains.
- Matches the `<dialog>` migration precedent (ADR-0014): platform over JS.

**Negative:**

- Popover and Dropdown require a positioning script (the top layer detaches
  them from the trigger). This is unavoidable until CSS Anchor Positioning is
  Baseline-wide.
- `Tooltip` no longer supports click-toggle on touch devices (show/hide is
  CSS-only). Trade-off accepted: tooltips are progressive enhancement.
- Older browsers (pre-2023) render popover content inline — visible-but-
  unstyled. Accepted given Baseline 2024.

**Migration path for consumers:**

- No API changes to any props struct.
- `TooltipProps`, `PopoverProps`, `DropdownProps`, `ContextMenuProps` keep the same fields.
- The only behavioral change: `Tooltip` no longer shows on touch-tap. Consumers serving
  touch-heavy audiences should provide a `<button>`-triggered `Popover` instead.

## Rollback plan

Per-component rollback is straightforward: revert the `.templ` file and restore the deleted
singleton JS in `shared.go` (Tooltip) or inline `<script>` (Popover/Dropdown/ContextMenu).
No data migrations, no API breakage.

## References

- [Popover API research note](../research/popover-api.md)
- [ADR-0014: Native `<dialog>` element for Modal and Drawer](0014-dialog-migration.md) — the
  precedent for platform-over-JS migrations
- [MDN: Popover API](https://developer.mozilla.org/en-US/docs/Web/API/Popover_API)
- [Baseline: `popover` attribute](https://web.dev/articles/popover-api) — April 2024
