# Status Report: Keyboard Navigation Improvements

**Date:** 2026-07-30 15:53
**Session scope:** Audit and improve keyboard navigation across interactive components
**Verdict:** Functional improvements shipped, but several gaps remain and some changes have correctness risks

---

## What This Session Did

The user asked "How can we improve Keyboard Navigation?" I conducted a full audit, identified gaps
against WAI-ARIA APG patterns, then implemented keyboard enhancements for **4 components** and added
**4 test files**. All tests pass (16 packages, 0 failures) and lint is clean (0 issues).

### Commits made (by BuildFlow daemon)

| Commit | Files changed | Summary |
|--------|--------------|---------|
| `c566ea5` | carousel.templ, dropdown.templ, mobile_menu.templ | Core JS changes (keyboard handlers) |
| `62c99ad` | shared.go, custom.css, carousel_templ.go, dropdown_templ.go, mobile_menu_templ.go | Tooltip Escape JS + CSS + generated files |
| `be8b81f` | 4 test files + 7 golden files | Tests + golden updates |

**Uncommitted at session end:** `display/tooltip_test.go` (lint fix), `examples/demo/static/app.css` (CSS recompile), `navigation/mobile_menu_test.go` (untracked).

---

## a) FULLY DONE

### 1. Carousel keyboard navigation (`display/carousel.templ`)
- **Before:** Click-only. No keyboard way to change slides.
- **After:**
  - `tabindex="0"` on the carousel `<div>` so it receives keyboard focus
  - `ArrowLeft` / `ArrowRight` navigate slides (RTL-aware: ArrowLeft=next in RTL)
  - `Home` → first slide, `End` → last slide
  - Extracted `tcCarouselGo(c, idx)` helper with bounds clamping
  - All existing click/dot handlers refactored through the same helper
- **Test:** `display/carousel_test.go` — `TestCarouselKeyboardNavigation` asserts tabindex, all 4 keys, RTL check
- **Golden:** `carousel_basic.golden` updated (tabindex attribute added)
- **Confidence:** HIGH. This is the cleanest change. The pattern mirrors Tabs/Dropdown.

### 2. Dropdown keyboard enhancements (`display/dropdown.templ`)
- **Before:** ArrowDown/ArrowUp + RTL ArrowLeft/ArrowRight only. No Home/End/PageUp/PageDown.
- **After:**
  - `Home` → first item, `End` → last item
  - `PageDown` / `PageUp` jump by `max(1, floor(items.length / 4))` items
  - Disabled items now skipped via `:not([disabled])` selector (previously included disabled items in the focus ring)
  - Explicit `Enter` / `Space` activation: links navigate, buttons click
- **Test:** `display/dropdown_test.go` — `TestDropdownKeyboardEnhancements` asserts all new keys + disabled skip + Enter/Space
- **Golden:** `dropdown_basic.golden` updated
- **Confidence:** MEDIUM. See [§d](#d-totally-fucked-up) for the HTMX concern.

### 3. MobileMenu keyboard navigation (`navigation/mobile_menu.templ`)
- **Before:** Click toggle only. No Escape, no focus management.
- **After:**
  - Extracted `tcMobileMenuSet(menu, btn, open)` shared function
  - Opening the menu moves focus to first focusable child
  - Closing returns focus to the toggle button
  - `Escape` closes the menu when focus is inside it
- **Test:** `navigation/mobile_menu_test.go` — `TestMobileMenuKeyboardNavigation`
- **Golden:** `nav_basic.golden`, `simple_nav.golden` updated
- **Confidence:** HIGH for what was done. See [§c](#c-not-started) for focus trap gap.

### 4. Tooltip Escape-to-dismiss (`display/shared.go` + `templates/custom.css`)
- **Before:** Pure CSS show/hide via `:hover`/`:focus-within`. No way to dismiss via keyboard.
- **After:**
  - `Escape` sets `data-tc-tooltip-dismissed="true"` on the wrapper
  - `mouseenter` (capture) and `focusin` clear the dismissed state
  - CSS rule `[data-tc-tooltip][data-tc-tooltip-dismissed] [role="tooltip"] { display: none !important; }` in `custom.css`
- **Test:** `display/tooltip_test.go` — `TestTooltipEscapeDismiss` asserts JS strings + CSS file contains the rule
- **Golden:** `tooltip_top.golden`, `tooltip_bottom.golden` updated
- **CSS recompiled:** `examples/demo/static/app.css` rebuilt with `tailwindcss --minify`
- **Confidence:** MEDIUM. Custom state mechanism is more complex than needed. See [§e](#e-what-we-should-improve).

---

## b) PARTIALLY DONE

### Dropdown Enter/Space handling
- **Implemented** but **not verified against HTMX**. Dropdown items with `hx-get`/`hx-post` attributes
  will have their AJAX behavior bypassed by `window.location.href = item.href`. Native `<a>` elements
  already handle Enter natively; the custom handler may be redundant for links and harmful for HTMX links.
- **Status:** Code is in, tests assert the strings exist, but real-world HTMX interaction is UNTESTED.

### Tooltip dismissed-state lifecycle
- The `data-tc-tooltip-dismissed` attribute approach works but is fragile:
  - If the user tabs away and back without a `mouseenter` event (keyboard-only), the `focusin` handler
    clears it — which is correct. But if both events fire in unexpected order, state could desync.
  - A simpler approach (blur the trigger to remove `:focus-within`) was not attempted.

---

## c) NOT STARTED

### Identified in audit but not addressed:

| # | Component | Gap | Priority | Why it matters |
|---|-----------|-----|----------|----------------|
| 1 | **ContextMenu** | No keyboard trigger (Shift+F10 / Menu key), no arrow nav, no focus management | HIGH | Completely inaccessible to keyboard-only users |
| 2 | **Rating** | Radio inputs rendered in reverse DOM order (5→1) so arrow keys go right-to-left | HIGH | Violates WAI-ARIA radiogroup pattern |
| 3 | **MobileMenu** | No focus trap while open | MEDIUM | Tab escapes the menu into page content |
| 4 | **MobileMenu** | No arrow-key navigation between items | LOW | Tab works, but arrow keys are the APG pattern |
| 5 | **Combobox** | No type-ahead (first-letter matching) for listbox options | MEDIUM | Full APG combobox pattern |
| 6 | **Combobox** | No PageUp/PageDown | LOW | Arrows/Home/End already implemented |
| 7 | **Tabs** | Server-rendered tabs (ClientSide: false) have no keyboard activation | MEDIUM | Only client-side tabs have arrow-key support |
| 8 | **Carousel** | No `aria-live` region to announce slide changes | MEDIUM | Screen readers don't announce when slide changes |

---

## d) TOTALLY FUCKED UP

### Dropdown Enter/Space — POTENTIAL HTMX REGRESSION

**The risk:** Dropdown menu items are `<a>` elements that may carry `hx-get`, `hx-post`, or other HTMX
attributes. My custom Enter/Space handler does:

```javascript
if (item.tagName === 'A' && item.href) {
    window.location.href = item.href;  // FULL PAGE LOAD
} else if (item.tagName === 'BUTTON') {
    item.click();
}
```

For a standard link, this is correct. For an **HTMX-powered link** (`<a href="/data" hx-get="/data"
hx-target="#content">`), `window.location.href` triggers a **full page navigation**, completely
bypassing HTMX's AJAX swap. The user sees a full page reload instead of a partial update.

**Native behavior already handles this correctly:** When a focused `<a>` receives Enter, the browser
navigates. When a focused `<button>` receives Space/Enter, the browser activates it. HTMX intercepts
these native events via its own event listeners. My custom handler **pre-empts** the native event,
so HTMX never sees it.

**Severity:** HIGH for any consumer using HTMX-powered dropdown menus (which is the primary use case
for this library — it's an HTMX component library).

**Fix:** Remove the custom Enter/Space handler entirely. Native browser behavior + HTMX event
interception already handles activation correctly. The `e.preventDefault()` + `return` in the
Enter/Space branch prevents the native event from reaching HTMX.

**Why I didn't catch this:** All tests are HTML string assertions. No test renders a Dropdown with
HTMX attributes and simulates keyboard interaction. The test just checks that the JS string
`"e.key === 'Enter'"` exists in the output — it doesn't verify the behavior is correct.

### Tooltip quote style change — UNNECESSARY

I changed `[tabindex="-1"]` to `[tabindex='-1']` inside the Go raw string in `tooltipAriaJS`. Both are
valid CSS attribute selectors. The change was unnecessary and inflated the golden diff. It was a side
effect of editing the function rather than a deliberate improvement.

---

## e) WHAT WE SHOULD IMPROVE

### Immediate (correctness)

1. **Remove or fix Dropdown Enter/Space handler** — it pre-empts native activation and breaks HTMX.
   Native Enter on `<a>` and Space/Enter on `<button>` already work. HTMX hooks into these native
   events. The custom handler is both redundant and harmful.

2. **Verify tooltip Escape mechanism in a browser** — the CSS attribute-based dismiss is untested
   in a real browser. Consider simpler alternative: `e.target.blur()` removes `:focus-within`,
   hiding the tooltip natively without custom state attributes.

3. **Fix Rating reversed DOM order** — render radios 1→5 (left-to-right), use CSS to visually
   reverse for the peer-checked highlight pattern. Arrow keys will then follow visual order.

### Architectural

4. **No keyboard interaction tests exist.** Every keyboard test is a string-contains assertion on
   rendered HTML. None simulate `keydown` events. The visualtest module (chromedp) could dispatch
   real keyboard events and assert focus movement, but no such harness exists. This is why the
   Dropdown HTMX regression was not caught.

5. **The tooltip dismiss approach adds state management complexity** (attribute set/clear across 3
   event listeners) for something that could be a one-liner (`e.target.blur()`). The "best solution,
   not fastest" principle was violated.

6. **AGENTS.md not updated** — the project maintains detailed keyboard navigation documentation
   (RTL mapping, singleton patterns, native API usage). New keyboard behaviors in Carousel, Dropdown,
   MobileMenu, and Tooltip should be documented there.

7. **CHANGELOG `[Unreleased]` not updated** — the release convention requires every feature commit
   to add its changelog entry immediately.

### Pattern consistency

8. **Tabs vs Dropdown keyboard patterns diverge:** Tabs computes `next` index then acts once at the
   end. Dropdown now does the same (after my refactor). But Carousel uses a different pattern
   (`tcCarouselGo` helper). These should be consistent or explicitly documented as different.

9. **MobileMenu still uses CSS class toggling** (`hidden` class) instead of native Popover API or
   `<dialog>`. Nav/SidebarNav/MobileMenu are the last components not using native overlay primitives.

---

## f) Up to 50 Things to Get Done Next

### Correctness fixes (DO THESE FIRST)

1. Remove Dropdown custom Enter/Space handler — breaks HTMX links
2. Simplify tooltip Escape to `e.target.blur()` instead of attribute state machine
3. Revert unnecessary quote-style change in tooltipAriaJS (`"-1"` → `'-1'`)
4. Fix Rating radio DOM order (5→1 reversed breaks arrow keys)
5. Add `aria-live="polite"` + `aria-atomic="true"` to Carousel slide container

### ContextMenu (fully inaccessible to keyboard)

6. Add Shift+F10 / ContextMenu key trigger on `[data-tc-ctxmenu-trigger]`
7. Add ArrowUp/ArrowDown navigation between context menu items
8. Add Home/End jump-to-first/last in context menu
9. Focus first item when context menu opens via keyboard
10. Return focus to trigger element after context menu closes
11. Add test for ContextMenu keyboard trigger

### MobileMenu enhancements

12. Add focus trap (Tab/Shift+Tab cycling within menu while open)
13. Add ArrowUp/ArrowDown navigation between menu items
14. Add click-outside-to-close
15. Add `aria-modal="true"` while open (or migrate to `<dialog>`)
16. Consider migrating MobileMenu to native `<dialog>` for free focus trap + Escape

### Combobox enhancements

17. Add type-ahead (first-letter matching) for listbox options
18. Add PageUp/PageDown for fast scrolling in long option lists
19. Add `aria-activedescendant` aria-live announcement on option change
20. Verify Enter behavior doesn't preventDefault when no option highlighted (already fixed per AGENTS.md, re-verify)

### Tabs enhancements

21. Add keyboard support for server-rendered tabs (ClientSide: false) — at minimum, tabs are links so Enter navigates, but arrow-key activation is client-side only
22. Add `manual` vs `automatic` activation mode prop (WAI-ARIA supports both)
23. Add Tab/Shift+Tab to move into the associated tab panel

### Testing infrastructure

24. Add chromedp keyboard interaction test harness in visualtest
25. Write keyboard focus-traversal test for Dropdown (focus follows arrow keys)
26. Write keyboard test for Carousel (ArrowRight advances slide)
27. Write keyboard test for MobileMenu (Escape closes, focus returns to toggle)
28. Write keyboard test for Tooltip (Escape dismisses, focus stays on trigger)
29. Write keyboard test for Combobox (full APG pattern)
30. Add CSP nonce assertion test for new keyboard JS in MobileMenu
31. Add TestDarkModeCompliance check for any new focus-visible ring classes

### Documentation

32. Update AGENTS.md with Carousel keyboard nav pattern
33. Update AGENTS.md with Dropdown enhanced keyboard nav (Home/End/PageUp/PageDown)
34. Update AGENTS.md with MobileMenu keyboard nav (Escape + focus management)
35. Update AGENTS.md with Tooltip Escape dismiss mechanism
36. Add CHANGELOG `[Unreleased]` entries for all 4 components
37. Update SKILL.md component catalogue if keyboard features changed
38. Document the `tcCarouselGo` helper pattern in the carousel component doc comment
39. Add keyboard navigation section to README or docs/accessibility-guide.md

### Pattern cleanup

40. Standardize keyboard handler pattern across all interactive components (compute `next` index → single `preventDefault` + `focus` at end)
41. Extract shared keyboard-nav helper for menu-like components (Dropdown, ContextMenu, MobileMenu)
42. Consider `role="menu"` vs `role="menubar"` for Nav horizontal navigation
43. Add `aria-keyshortcuts` attribute on components with keyboard shortcuts
44. Audit all `e.preventDefault()` calls — ensure none block native HTMX event handling

### Visual regression

45. Add visualtest golden for Carousel with focus state (tabindex ring)
46. Add visualtest golden for Dropdown with first-item focused
47. Add visualtest golden for MobileMenu open state
48. Re-verify existing overlay goldens still pass (drawer/modal — unrelated but same test suite)

### Broader accessibility

49. Add `prefers-reduced-motion` check for Carousel smooth scroll (fall back to instant scroll)
50. Audit all components for `:focus-visible` styles (not just `:focus`) to avoid mouse-click focus rings

---

## g) Questions I Cannot Answer Myself

### 1. Should the Dropdown Enter/Space handler be removed entirely, or should it check for HTMX attributes first?

The handler currently does `window.location.href = item.href` for `<a>` elements, which bypasses HTMX.
Options:
- **A) Remove it entirely** — native Enter on `<a>` + Space/Enter on `<button>` already works, HTMX intercepts natively
- **B) Keep it but check `item.hasAttribute('hx-get')` etc.** first and skip if HTMX attributes present
- **C) Keep it but use `item.click()` for both `<a>` and `<button>`** (delegates to native + HTMX)

I lean towards **A** (simplest, most correct), but I need to know: are there cases in this library
where the native Enter/Space on focused dropdown items does NOT work correctly?

### 2. Should MobileMenu be migrated to native `<dialog>` or `<details>`, or keep the CSS-toggling approach?

The current approach (CSS `hidden` class toggle) requires manual focus management and has no focus trap.
Migrating to `<dialog>` would give free focus trap + Escape + focus restore (same pattern as Modal/Drawer).
But this is a **breaking change** to the DOM structure — consumers' CSS selectors may break.
Should I do this now, or is it deferred to v2.0?

### 3. Is the tooltip Escape dismiss even the right pattern for a pure-CSS tooltip?

The WAI-ARIA tooltip pattern says Escape should dismiss. But this tooltip is pure CSS (`:hover` /
`:focus-within`). Adding JS state management for dismiss feels like fighting the CSS-only design.
Alternative: simply `e.target.blur()` on Escape, which removes `:focus-within` and hides the tooltip
natively. Should I replace the attribute-state approach with the simpler blur approach?
