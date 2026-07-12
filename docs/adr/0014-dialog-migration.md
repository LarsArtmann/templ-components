# ADR 0014: Native `<dialog>` element for Modal and Drawer

## Date

2026-07-12

## Status

Accepted

## Context

Modal and Drawer components used a custom JavaScript overlay system that
required ~200 lines of JS per overlay instance to handle:

- Focus trapping (Tab cycling between first/last focusable element)
- Escape key dismissal
- Focus save and restore (return focus to trigger element)
- Top-layer rendering (z-index management)
- Backdrop dimming (custom `<div>` overlay)
- `inert` / `aria-hidden` toggling on background content
- Animation class toggling for open/close transitions

This JS was complex, fragile, and required careful synchronization between
server-rendered state and client-side DOM manipulation. The `shared.go` file
contained `overlayTrapJS`, `overlayOpenJS`, `overlayCloseJS`,
`overlayPanelConfig`, `focusableSelector`, and `jsClassArgs` — all to replicate
behavior that browsers now handle natively via the `<dialog>` element and the
`showModal()` / `close()` API.

## Decision

Migrate Modal and Drawer to the native HTML `<dialog>` element.

**What changed:**

- `<div role="dialog">` replaced with `<dialog>` element
- Custom backdrop `<div>` removed (now uses `::backdrop` pseudo-element)
- `showModal()` / `close()` replace all custom open/close JS
- CSS `@starting-style` + `transition-behavior: allow-discrete` replace JS
  animation class toggling
- `overlayDialogJS` (~15 lines) replaces the 200-line JS generation system:
  singleton guard + thin `tcOpen`/`tcClose` wrappers + per-instance IIFE for
  auto-open and click delegation

**What the browser now handles natively (previously custom JS):**

| Capability              | Before                                       | After                       |
| ----------------------- | -------------------------------------------- | --------------------------- |
| Focus trap              | Custom Tab/Shift+Tab cycling                 | Native to `showModal()`     |
| Escape dismissal        | Custom `keydown` listener                    | Native to `<dialog>`        |
| Focus restore           | Manual save/restore via `data-tc-prev-focus` | Native to `close()`         |
| Top-layer rendering     | z-index management (z-50, z-[100])           | Native top layer            |
| Backdrop dimming        | Custom `<div>` with opacity                  | `::backdrop` pseudo-element |
| Background `inert`      | Manual toggling                              | Native to `showModal()`     |
| `aria-hidden` lifecycle | Manual toggling                              | Native to `<dialog>`        |

## Consequences

**Positive:**

- ~200 lines of JS eliminated. The `shared.go` file shrank from 194+ lines of
  JS generation to ~15 lines.
- Browser-native behavior is more robust than custom JS — no edge cases with
  nested modals, Shadow DOM, or third-party iframes.
- CSS animations work even with JS disabled (progressive enhancement via
  `@starting-style`).
- Simpler CSP compliance — less inline JS to nonce.

**Negative:**

- Closed `<dialog>` elements are `display: none` (previously, content was in the
  DOM with `opacity-0`). This is a behavior change for consumers who query
  closed modal content via JS. Documented in AGENTS.md.
- `::backdrop` styling requires CSS — consumers must copy the `.tc-overlay`,
  `.tc-modal`, `.tc-drawer` CSS rules from `templates/app.css`. Without them,
  dialogs open with no animation and no backdrop dimming.
- `allow-discrete` and `@starting-style` are Chrome 117+ / Safari 17.4+ / Firefox
  128+. Browsers without support snap instantly (graceful degradation).

**CSS dependency chain:**

The `.tc-overlay`/`.tc-modal`/`.tc-drawer` classes in `templates/app.css` are
required for dialog animations and backdrop styling. This matches the existing
`.tc-*` CSS pattern (e.g., `.tc-auto-grow`, `.tc-select`). Consumers who vendor
`app.css` get these automatically.

## Alternatives considered

### Popover API

The Popover API (`popover` attribute) was considered but rejected because:

1. CSS Anchor Positioning (required for trigger-relative positioning) is not yet
   Baseline (Firefox unsupported).
2. Modal/Drawer are viewport-centered, not trigger-relative — `<dialog>` is the
   correct semantic element for modal dialogs.
3. Popover API does not provide focus trapping by default (only `popover="auto"`
   provides light dismiss, not focus trap).

### Keep custom JS

Considered maintaining the existing JS system. Rejected because:

1. The JS was the single largest source of complexity in the `display` package.
2. Native `<dialog>` is Baseline since 2023 (all major browsers).
3. The migration eliminated more code than it added — net reduction of ~180
   lines.
