# Status Report: Keyboard Navigation Fixes & Remaining Gaps

**Date:** 2026-07-30 16:22
**Session focus:** Fixing the critical Dropdown HTMX regression, reverting unnecessary changes, and auditing what remains

---

## a) FULLY DONE

1. **Dropdown Enter/Space HTMX regression FIXED.** The custom JS handler at `display/dropdown.templ` used `window.location.href = item.href` for `<a>` menuitems, which triggered a full page load and bypassed HTMX's AJAX swap. The entire Enter/Space block was removed. Native browser behavior activates links/buttons correctly, and HTMX intercepts native events. Test (`display/dropdown_test.go:TestDropdownKeyboardEnhancements`) now asserts `window.location.href` is absent via `AssertNotContains`.

2. **Tooltip querySelector syntax error FIXED.** The original `display/shared.go` had `[tabindex="-1"]` inside a double-quoted JS string (`querySelector("...[tabindex="-1"]...")`) — a latent syntax error that was never caught because tests use string-contains assertions, not execution. Fixed to single-quoted JS string with double-quoted CSS attribute: `'...[tabindex="-1"]...'`, matching the pattern already used in `navigation/mobile_menu.templ:81`.

3. **Tooltip Escape-to-dismiss (kept as-is).** Evaluated the `e.target.blur()` alternative and decided the attribute-based approach (`data-tc-tooltip-dismissed` + CSS rule) is correct: `blur()` would move focus to `<body>`, losing the user's tab position. The attribute approach hides the tooltip while keeping focus on the trigger.

4. **Golden files regenerated** for dropdown_basic, tooltip_top, tooltip_bottom after the source fixes.

5. **Full test suite passes** — all 16 packages, 0 failures.

6. **Lint passes** — golangci-lint with 67 enabled linters, 0 issues.

7. **CHANGELOG.md updated** — `[Unreleased]` section has Added (Carousel keyboard, MobileMenu keyboard, Dropdown keyboard enhancements, Tooltip Escape) and Fixed (Dropdown HTMX regression).

8. **AGENTS.md updated** — Carousel, Dropdown, Tooltip, and MobileMenu entries now document keyboard behavior. RTL keyboard mapping entry updated to include Carousel.

9. **Previous session work confirmed intact:**
   - Carousel: `tabindex="0"` on region, ArrowLeft/Right/Home/End keydown handler, RTL-aware
   - MobileMenu: Escape closes + focus to button, open focuses first child, shared `tcMobileMenuSet` helper
   - Dropdown: Home/End/PageUp/PageDown, `:not([disabled])` selector, first-item focus on `toggle` event
   - Tooltip: Escape sets `data-tc-tooltip-dismissed`, mouseenter/focusin clears it

---

## b) PARTIALLY DONE

1. **Keyboard test coverage is assertion-only, not execution-based.** All 4 new keyboard tests (`TestCarouselKeyboardNavigation`, `TestMobileMenuKeyboardNavigation`, `TestDropdownKeyboardEnhancements`, `TestTooltipEscapeDismiss`) verify that specific JS strings are present in the rendered HTML. **None of them actually dispatch keyboard events.** This is exactly why the Dropdown HTMX regression went undetected — the test asserted `e.key === 'Enter'` was present, but never tested what happened when Enter was actually pressed. The `visualtest` module (chromedp + headless Chromium) has the infrastructure to dispatch real keyboard events but no keyboard interaction harness exists.

2. **Carousel golden test does not capture keyboard JS.** The golden sweep test (`TestGoldenSweepCarousel`) renders CarouselProps without a Nonce, so the `if props.Nonce != ""` guard prevents the `<script>` tag from rendering. The golden file (`carousel_basic.golden`) correctly has `tabindex="0"` on the div but has zero JS. This means golden tests cannot catch regressions in the keyboard handler JS. Same issue applies to all component goldens that gate scripts behind Nonce.

3. **AGENTS.md documentation is spread across multiple bullet points.** Keyboard behavior for each component is documented in separate entries (Carousel, Dropdown, Tooltip, MobileMenu) rather than in a consolidated "Keyboard Navigation" section. This makes it harder to understand the overall keyboard story.

---

## c) NOT STARTED

1. **ContextMenu keyboard accessibility.** `display/context_menu.templ` only listens to the `contextmenu` mouse event (right-click). There is zero keyboard support: no Shift+F10 to trigger, no arrow key navigation, no Enter to activate. This component is completely inaccessible to keyboard-only users. WAI-ARIA APG requires Shift+F10 + full menu keyboard nav.

2. **Rating component reversed DOM order.** `forms/rating.templ:104` renders radio inputs in reverse order (`for i := maxStars; i >= 1; i--`) so CSS can show filled stars via `~` selector. But this means ArrowLeft/Right follow DOM order (5→4→3→2→1), breaking the visual left-to-right expectation. ArrowRight should go to the next higher rating, not lower. This is a pre-existing accessibility bug.

3. **Keyboard interaction test harness.** No infrastructure exists to simulate keyboard events in tests. The `visualtest` module could provide this (chromedp has `input.DispatchKeyEvent`), but no keyboard-specific helper has been written.

4. **Focus visible styles audit.** No systematic check that all interactive components have visible focus indicators (`:focus-visible` outlines). Components that use `tabindex="0"` (Carousel) or are focusable by default should have clear focus rings.

5. **Skip links / bypass blocks.** No audit of whether Page/AppShell provides a "Skip to main content" link for screen reader and keyboard users (WCAG 2.4.1).

6. **Tab order audit.** No systematic check that tab order matches visual/semantic order across all components, especially in composite widgets.

---

## d) TOTALLY FUCKED UP

1. **The Dropdown Enter/Space handler was a regression I introduced and then had to fix.** In the previous session, I added a custom Enter/Space handler that used `window.location.href = item.href` for `<a>` elements. This was fundamentally broken for HTMX — it bypassed AJAX and caused full page loads. The fix (removing it entirely) was correct, but the bug should never have been introduced. The root cause was implementing keyboard activation without understanding that native browser behavior already handles it and HTMX hooks into native events.

2. **The tooltip quote-style change was initially misdiagnosed.** I initially thought the `'-1'` → `"-1"` change was an unnecessary side effect. Then I "fixed" it by escaping with `\"`, then realized that was ugly, then finally settled on the correct pattern (single-quoted JS string). This took 3 iterations for what should have been a 1-step fix. The original code was actually a latent syntax error (double quotes inside double quotes), so the previous session's change to single quotes was actually a *fix*, not a regression — I was wrong about it being wrong.

3. **BuildFlow daemon commit messages continue to be terrible.** The daemon committed my source fixes with messages like `"ANGELOG.md documentation"` (typo, should be CHANGELOG) and `"refactor(display): consolidate shared component rendering logic"` (hallucinated summary that doesn't mention the HTMX regression fix). These commits are invisible to `git log --grep` for "HTMX" or "Enter/Space". This is a known issue documented in AGENTS.md (T13) but remains unfixed.

4. **No keyboard event simulation in any test.** This is the deepest structural gap. The Dropdown HTMX regression went undetected because tests only check string presence. Until we have real keyboard event dispatch (even via chromedp), keyboard regressions will keep slipping through.

---

## e) WHAT WE SHOULD IMPROVE

### Architecture / Design

1. **Consolidate keyboard navigation into a shared helper.** Each component implements its own `addEventListener('keydown', ...)` handler with its own key matching, RTL detection, and focus management. A shared `tcMenuKeyboard(handler)` or `tcRtlArrowKeys(e, onNext, onPrev)` utility would reduce duplication and ensure consistent behavior.

2. **Nonce-gated scripts break golden test coverage.** When `Nonce` is empty, the `<script>` tag is not rendered, so golden tests never capture JS. Consider always rendering scripts (with empty nonce when CSP is not enforced) or adding a test mode that renders scripts regardless.

3. **Tooltip dismiss state machine complexity.** The `data-tc-tooltip-dismissed` attribute + CSS rule + 3 event listeners (keydown, mouseenter, focusin) is a lot of machinery for "hide tooltip on Escape". Consider whether a simpler approach exists, or whether the complexity is justified by the focus-preservation requirement.

4. **Singleton guard pattern repetition.** Every component's JS starts with `if(!window.tcXxxAttached){window.tcXxxAttached=true;...}`. This is a shared pattern that could be extracted, though the templ raw-string format makes extraction awkward.

### Testing

5. **Add chromedp-based keyboard event tests.** The `visualtest` module already has headless Chromium. A `keyboard_test.go` could dispatch real keydown events and assert DOM state changes (focus moved, class toggled, etc.).

6. **Test HTMX interaction in Dropdown.** Specifically: render a Dropdown with an `<a hx-get="/data">` menuitem, dispatch Enter, and assert that HTMX's AJAX swap fires (not `window.location.href`).

7. **Golden tests with Nonce set.** Add golden variants where `Nonce: "test-nonce"` to capture the JS in golden files.

### Accessibility

8. **Focus-visible styles.** Audit all components for `:focus-visible` outlines. Carousel's `tabindex="0"` div needs a focus ring.

9. **ARIA live regions for dynamic changes.** Carousel slide changes should be announced to screen readers (currently silent).

10. **Reduced motion + keyboard.** Verify that keyboard navigation works when CSS animations are disabled (`motion-reduce`).

### Code Quality

11. **The Carousel `e.preventDefault()` is called before checking if the key is actually used.** It filters keys first (`if(e.key!=='ArrowLeft'&&...)`), then calls `preventDefault()`. This is correct but fragile — if the filter list changes, preventDefault might swallow unintended keys.

12. **Dropdown pageSize calculation.** `Math.floor(items.length / 4)` means a 3-item menu jumps by 0 (PageDown does nothing useful). Consider `Math.max(1, ...)` — oh wait, it already has that. But for 3 items, pageSize=1, so PageDown moves by 1, same as ArrowDown. Consider whether PageUp/PageDown add value for small menus.

13. **MobileMenu focus selector is duplicated.** `navigation/mobile_menu.templ:81` and `display/shared.go:328` both have `a[href],button:not([disabled]),input,select,textarea,[tabindex]:not([tabindex="-1"])`. Extract to a shared constant.

---

## f) Up to 50 Things We Should Get Done Next

### Critical (accessibility blockers)

1. ContextMenu: Add Shift+F10 keyboard trigger
2. ContextMenu: Add arrow key navigation (reuse Dropdown menu keyboard pattern)
3. ContextMenu: Add Enter/Space activation (native, no custom handler)
4. Rating: Fix reversed DOM order breaking arrow key direction
5. Carousel: Add `:focus-visible` outline for keyboard users
6. Audit all components for `:focus-visible` styles
7. Add "Skip to main content" link to Page/AppShell (WCAG 2.4.1)

### High priority (testing infrastructure)

8. Build chromedp keyboard event test harness in visualtest
9. Add keyboard interaction test for Carousel (ArrowLeft/Right/Home/End)
10. Add keyboard interaction test for Dropdown (Home/End/PageUp/PageDown/disabled-skip)
11. Add keyboard interaction test for MobileMenu (Escape closes, focus management)
12. Add keyboard interaction test for Tooltip (Escape dismisses)
13. Add HTMX interaction test for Dropdown (Enter on `<a hx-get>` triggers AJAX, not navigation)
14. Add golden test variants with Nonce set to capture JS output
15. Add ContextMenu keyboard golden tests after implementing keyboard support

### Medium priority (code quality)

16. Extract shared keyboard nav utility (RTL arrow detection, focus cycling)
17. Extract shared "first focusable element" selector constant
18. Consolidate keyboard documentation into a single AGENTS.md section
19. Add `aria-live` region for Carousel slide changes
20. Review PageUp/PageDown value for small Dropdown menus
21. Add Tab key trap prevention audit for all overlay components (Modal, Drawer, Dropdown)
22. Verify focus restoration after Modal/Drawer close works with keyboard activation
23. Audit Combobox keyboard nav against WAI-ARIA combobox pattern
24. Audit Tabs keyboard nav against WAI-ARIA tabs pattern (already good, verify)
25. Audit TagsInput keyboard nav (Enter to add, Backspace to remove)

### Lower priority (polish)

26. Add keyboard shortcut help text to demo page
27. Document keyboard shortcuts in component docs
28. Add `accesskey` support to Button/Link components
29. Add inert attribute support for background content during modal open
30. Audit scroll-behavior: smooth respects motion-reduce
31. Add focus trap utility for custom modal-like components
32. Verify all form controls work with Enter key submission
33. Add keyboard navigation to Pagination (already has arrow buttons, verify keyboard)
34. Add keyboard navigation to Accordion (native `<details>` — verify Enter/Space)
35. Audit Sidebar nav for keyboard accessibility
36. Add keyboard test for CopyButton (Enter/Space to copy)
37. Add keyboard test for LoadMore button
38. Verify Table sortable headers are keyboard accessible
39. Add keyboard support for TableRow href navigation (Enter on focused row)
40. Verify Nav menu keyboard accessibility (tab order, Enter activation)
41. Add visual focus indicators to all interactive elements in dark mode
42. Audit that `prefers-reduced-motion` doesn't break keyboard navigation
43. Add keyboard shortcut to toggle theme (Ctrl+Shift+L or similar)
44. Test all components in screen reader (NVDA/VoiceOver) for correct announcements
45. Add `role="application"` audit for Carousel (currently `role="region"`)
46. Consider `aria-activedescendant` pattern for Combobox instead of `aria-owns`
47. Add keyboard test for Stylable Select (Enter to open, arrow keys, Enter to select)
48. Add keyboard test for Calendar/DatePicker if it exists
49. Verify Slider keyboard support (arrow keys, Home, End)
50. Add keyboard navigation end-to-end test in demo application

---

## g) Questions (3)

### Q1: Should we build a chromedp keyboard event test harness now, or wait?

Building a real keyboard event test harness in `visualtest/` (chromedp `input.DispatchKeyEvent`) would catch regressions like the Dropdown HTMX bug. But it's a significant infrastructure investment (~200-300 lines + a test per component). Alternatively, we could use a lighter approach: extract the JS handlers to pure functions and unit-test the logic in Go. The JS-in-raw-strings format makes extraction hard though. **Which approach do you prefer: chromedp integration tests, or accept that keyboard JS is tested via string assertions + manual testing?**

### Q2: Should ContextMenu get full keyboard support (Shift+F10 + arrow nav)?

ContextMenu is currently mouse-only (right-click). Adding Shift+F10 + arrow key navigation would make it WAI-ARIA compliant, but it would reuse the Dropdown menu keyboard pattern (which we just stabilized). The alternative is to document ContextMenu as "mouse/touch only — use Dropdown for keyboard-accessible menus." **Should we invest in ContextMenu keyboard support, or steer users to Dropdown for keyboard scenarios?**

### Q3: Should we fix the Rating reversed-DOM-order bug?

`forms/rating.templ` renders radio inputs in reverse DOM order (5→1) so CSS `~` sibling selector can show filled stars. This breaks ArrowLeft/Right direction. Fixing it requires either: (a) CSS Grid to reverse visual order while keeping DOM order, (b) `flex-row-reverse`, or (c) JS to handle arrow key direction reversal. Option (b) is simplest but has browser quirks. **Should we fix this now, or is it a known tradeoff we accept (the visual star-fill behavior is the priority)?**

---

## Session Metrics

| Metric | Value |
|--------|-------|
| Files changed (source) | 4 (.templ + .go) |
| Files changed (generated) | 3 (_templ.go) |
| Files changed (golden) | 4 (.golden) |
| Files changed (docs) | 2 (CHANGELOG, AGENTS) |
| Files changed (tests) | 2 (dropdown_test, tooltip_test) |
| Tests passing | 16/16 packages |
| Lint issues | 0 |
| Commits this session (BuildFlow) | 4 (2 source + 1 docs + 1 status report) |
| Critical bugs fixed | 2 (HTMX regression, querySelector syntax error) |
| Known issues remaining | 2 (ContextMenu keyboard, Rating DOM order) |
| BuildFlow commit message quality | Poor (hallucinated summaries, typo "ANGELOG") |
