# SidebarNav: theme-adaptive chrome (and the permanently-dark opt-out)

**Shipped:** shell theming tokens (`--tc-sidebar-*`, M22/F101) ·
**Browser-proof:** `visualtest/sidebar_nav_visual_test.go`
(`TestSidebarNavLightDark`, `TestSidebarNavClassicDarkOptOut`)

## What changed

`SidebarNav` used to be permanently dark (near-black chrome in both themes)
with hard-coded classes. It is now **theme-adaptive**, driven entirely by CSS
custom properties defined in `templates/custom.css`:

| Mode  | Sidebar background | Border      | Text       |
| ----- | ------------------ | ----------- | ---------- |
| Light | white              | `gray-200`  | `gray-700` |
| Dark  | `#000`             | transparent | `gray-300` |

Item hover, muted/section-label, and header tokens follow the same pattern.
No Go API changed — every visual difference is token-driven, so consumers who
rebrand via `@theme` get the sidebar rebrand for free.

## Opt-out: restore the classic permanently-dark sidebar

Set the same variables in `:root` (pure CSS, no library code touched). The
canonical, always-current snippet lives in `templates/custom.css` next to the
token definitions:

```css
:root {
  --tc-sidebar-bg: var(--color-gray-900);
  --tc-sidebar-border: transparent;
  --tc-sidebar-fg: var(--color-gray-300);
  --tc-sidebar-fg-hover: var(--color-white);
  --tc-sidebar-item-hover-bg: var(--color-gray-800);
  --tc-sidebar-muted: var(--color-gray-500);
  --tc-sidebar-muted-hover: var(--color-gray-300);
}
```

Dark-mode values mirror the original hard-coded classes exactly, so a
consumer with this block renders the classic dark admin chrome in BOTH modes
— pixel-equivalent to the pre-token sidebar (guarded by
`TestSidebarNavClassicDarkOptOut`, which DOM-asserts the override reaches the
rendered element and pins a golden of the result).

## Migration checklist

- No action needed to keep the new default (adaptive) behavior.
- If your admin UI depended on the always-dark sidebar, paste the opt-out
  block into your CSS entry point (after the library CSS).
- ADR-0011's documented exception ("permanently dark sidebar") is retired by
  this change; the compliance-test exemption it described is gone.
