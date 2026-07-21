# ADR 0017: Popover API migration for Dropdown, Popover, ContextMenu

## Date

2026-07-21

## Status

Accepted

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

| Component   | `popover` value      | Trigger mechanism                                     | JS remaining                                               |
| ----------- | -------------------- | ----------------------------------------------------- | ---------------------------------------------------------- |
| Tooltip     | _(none — CSS only)_  | `:hover` / `:focus-within` (existing CSS)             | **0** (was ~20)                                            |
| HoverCard   | _(none — unchanged)_ | `:hover` / `:focus-within` (existing CSS)             | 0                                                          |
| ContextMenu | `popover="auto"`     | `contextmenu` event → `showPopover()` + `style.inset` | ~6 (was ~12)                                               |
| Popover     | `popover="auto"`     | `<button popovertarget="id">` (native click-toggle)   | **0** (was ~40)                                            |
| Dropdown    | `popover="auto"`     | `<button popovertarget="id">` (native click-toggle)   | ~25 (was ~50, kept for ArrowUp/Down/Home/End keyboard nav) |

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

We do **not** adopt CSS Anchor Positioning in this migration. The current CSS class-based
positioning (`bottom-full left-1/2 -translate-x-1/2`, etc.) continues to work for `Popover` and
`Dropdown` because the popover element remains a DOM descendant of the relatively-positioned
trigger wrapper. Anchor Positioning is Baseline 2026 but offers no clear win over the existing
class-based positioning for these components. We may revisit Anchor for `ContextMenu` (where
cursor-relative `inset` is required) in a future enhancement.

### Fallback for older browsers

`popover="auto"` gracefully degrades: in browsers without Popover API support, the element
renders inline (visible, no top-layer). This is a worse UX but is not broken. We accept this
trade-off because the Baseline (April 2024) is now ~2 years old.

### CSP implications

- Popover and Dropdown lose their inline `<script>` tags entirely → **CSP surface shrinks**.
- The `integration/csp_nonce_test.go` tripwire still passes (ContextMenu and Dropdown still
  emit scripts for keyboard / right-click handlers).
- Add new CSP test entries for `Tooltip` and `Popover` as **negative** assertions (verify no
  `<script>` tag is emitted).

## Consequences

**Positive:**

- ~70 lines of JS deleted across 3 components.
- `Popover` becomes the first component with **zero** JS — entirely native.
- `Tooltip` becomes pure CSS — no more singleton in `shared.go`.
- Native light-dismiss, Escape, focus, and top-layer rendering replace fragile custom JS.
- CSP audit story improves (fewer inline scripts).
- Matches the `<dialog>` migration precedent (ADR-0014): platform over JS.

**Negative:**

- `Tooltip` no longer supports click-toggle on touch devices. Trade-off accepted: tooltips
  are non-critical progressive enhancement; the hover/focus CSS still works on desktop.
- `aria-describedby` is no longer auto-propagated from the Tooltip wrapper to its first
  focusable child. **Consumer responsibility:** set `aria-describedby` directly on the
  focusable trigger element. Documented in godoc on `TooltipProps`.
- Older browsers (pre-2023) render popover content inline — visible-but-unstyled popovers.
  Accepted given Baseline 2024.

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
