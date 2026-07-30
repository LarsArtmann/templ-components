# Status Report: Keyboard Navigation Improvements — Round 3

**Date:** 2026-07-30 17:00
**Session goal:** Close the 4 known keyboard-navigation gaps left from rounds 1-2 (Rating arrow-key direction, Carousel focus ring, ContextMenu keyboard access, shared menu-nav extraction)

---

## A) FULLY DONE (verified: build + 15/15 tests + 0 lint issues)

### 1. ContextMenu keyboard accessibility
- **Shift+F10 and the ContextMenu key** now open the menu (positioned at the trigger via `getBoundingClientRect()`, not the cursor — appropriate for keyboard activation)
- Menuitems gained `tabindex="-1"` (roving tabindex for WAI-ARIA menu pattern)
- Disabled items expose `aria-disabled="true"` (was just `pointer-events-none` opacity, invisible to AT)
- Shares the extracted menu keyboard nav helper (ArrowUp/Down, Home/End, PageUp/PageDown, focus-first-on-open)
- Escape and click-outside remain native via Popover API `popover="auto"`
- **Files:** `display/context_menu.templ`, `display/context_menu_test.go` (4 new tests)

### 2. Shared menu keyboard-nav extraction (DRY)
- Extracted Dropdown's 35-line inline `<script>` into `menuKeyboardNavJS()` + `menuKeyboardNavScriptComponent()` in `display/shared.go`
- Singleton guard: `tcMenuKeyNavAttached`
- Both Dropdown and ContextMenu now inject the same script via `@menuKeyboardNavScriptComponent(props.Nonce)`
- Improved: the nav selector now skips `[aria-disabled="true"]` items too (was only `[disabled]`), so ContextMenu's `aria-disabled` items are correctly skipped
- **Files:** `display/shared.go`, `display/dropdown.templ`, `display/dropdown_test.go`

### 3. Carousel focus-visible outline
- The `tabindex="0"` carousel region now shows `focus-visible:ring-2 focus-visible:ring-blue-500` with ring-offset
- Keyboard users can now see where focus landed before pressing arrow keys
- **Files:** `display/carousel.templ`, `display/carousel_test.go`

### 4. Documentation + infrastructure
- CHANGELOG `[Unreleased]` updated with all changes (Added + Fixed)
- AGENTS.md updated: Carousel entry (focus ring), Dropdown/ContextMenu entry (shared nav), new Rating DOM-order convention, RTL keyboard mapping updated
- Demo CSS recompiled via `nix develop -c tailwindcss` (new `flex-row-reverse` class added to compiled output)
- All golden files updated: `rating_interactive.golden`, `dropdown_basic.golden`, `carousel_basic.golden`
- All `*_templ.go` regenerated (91 files)

---

## B) PARTIALLY DONE (implemented but has a known defect)

### 5. Rating arrow-key direction fix — CORRECT but introduced a FILL BUG

**What's correct:**
- DOM order changed from reverse (5→1) to forward (1→N) — ArrowDown/ArrowRight now increase the value per WAI-ARIA radiogroup pattern ✓
- Fill classes (`peer-checked:text-amber-400`) moved from the nested `<svg>` to the `<label>` — Tailwind's `peer-checked` `~` combinator now resolves (SVG was a descendant, not a sibling, so it never matched before) ✓
- `flex-row-reverse` added to the interactive container to reverse the visual layout back to 5-left/1-right (matching the original design's visual order) ✓
- Read-only mode is unchanged (no flex-row-reverse, normal flow) ✓

**What's BROKEN (see section D below):**
- The cumulative star fill is **INVERTED** for all values except the exact middle. Selecting 1 star fills ALL 5 stars; selecting 5 stars fills only 1 star. Only value=3 out of 5 happens to render correctly by coincidence.

---

## C) NOT STARTED

1. **Chromedp keyboard-event test harness** — all keyboard tests remain string-assertion-only. No test dispatches real key events. This is why the fill bug went undetected (string tests verify class presence, not CSS activation).
2. **Visual regression test for Rating** — no `visualtest` golden exists for the Rating component, so the inverted fill would not be caught by pixel comparison either.
3. **Rating JavaScript fill handler** — the correct fix for both arrow-key direction AND cumulative fill requires either JS (singleton that fills stars on `change`/`hover`) or a fundamentally different CSS approach. Not implemented.

---

## D) TOTALLY FUCKED UP

### The Rating cumulative fill is INVERTED — a self-inflicted regression

**Root cause:** I used forward DOM order (1→5) + `flex-row-reverse` + `peer-checked:text-amber-400` on each label. The CSS `.peer:checked ~ .peer-checked\:text-amber-400` matches ALL labels that are siblings AFTER the checked radio in DOM order. With forward DOM:

| Selected value | Labels that match `~` (after checked radio in DOM) | Visual with flex-row-reverse | Expected | Correct? |
|---|---|---|---|---|
| 1 | labels 1,2,3,4,5 (ALL after radio-1) | ★★★★★ | ★☆☆☆☆ | **NO** |
| 2 | labels 2,3,4,5 | ★★★★☆ | ★★☆☆☆ | **NO** |
| 3 | labels 3,4,5 | ★★★☆☆ | ★★★☆☆ | YES (coincidence) |
| 4 | labels 4,5 | ★★☆☆☆ | ★★★★☆ | **NO** |
| 5 | label 5 only | ★☆☆☆☆ | ★★★★★ | **NO** |

The `peer-checked` `~` combinator only matches **forward siblings** (elements AFTER the checked peer in DOM). For a correct left-fill rating (★★★☆☆), we need the LOWER-value stars to be AFTER the checked radio — which is **reverse DOM order**. But reverse DOM order breaks arrow-key direction. This is a fundamental CSS limitation: `~` cannot match preceding siblings.

**Why it went undetected:** All tests are string-presence assertions (`utils.AssertContains`). They verify the class string exists in the HTML, not that the CSS rule activates for the correct elements. The `TestRatingKeyboardOrder` test checks forward DOM order and flex-row-reverse presence — both pass — but no test verifies that the fill renders correctly for value=1 or value=5.

**The old code was also broken** (differently): `peer-checked` classes were on the nested `<svg>` (a descendant of `<label>`, not a sibling of `.peer`), so the `~` combinator never matched and NO star ever filled. My fix moved classes to the label (correct for resolution) but the fill DIRECTION is wrong.

**Correct approaches (not yet implemented):**
1. **Forward DOM + JS singleton** for cumulative fill on `change`/`hover` (most production rating widgets do this)
2. **Forward DOM + non-cumulative fill** (only the selected star is amber; less pretty but honest)
3. **Forward DOM + Tailwind `has-*` / arbitrary variant** — e.g. `[&:has(~.peer:checked)]:text-amber-400` checks if a FOLLOWING sibling is checked. But this fills labels BEFORE the checked radio, which is the WRONG direction for forward DOM + left-fill.
4. **Keep reverse DOM order** (correct fill via `~`) and accept inverted arrow keys (the original tradeoff, which the user explicitly asked to fix)

The fundamental tension: CSS `~` can only look forward in DOM. Correct arrow keys require forward DOM. Correct cumulative fill via `peer-checked` requires reverse DOM. Both cannot be achieved with pure CSS alone.

### Other mistakes this session (all caught and fixed, but shouldn't have happened):

1. **Inconsistent singleton guard name** — wrote `tcMenu_KEYNAV_GUARD` / `tcMENU_KEYNAV_GUARD` in the same string. Caught on self-review before testing. Sloppy copy-paste.
2. **Applied `flex-row-reverse` to the shared container** (both readonly and interactive branches). Reversed the readonly icons to ☆★★★★ instead of ★★★★☆. Caught by golden test failure. Fixed with `utils.Ternary` conditional class.
3. **Wrong golden test name pattern** — ran `-run TestGoldenSweep` for forms, but rating goldens are named `TestGoldenRating*`. Forms goldens were stale after first update pass. Wasted a round trip.
4. **Ordered class substring assertion** — asserted `opacity-50 pointer-events-none` in ContextMenu test, but tailwind-merge reorders tokens. AGENTS.md explicitly warns against this. Fixed with `AssertContainsAll`.
5. **wsl_v5 whitespace** — missing blank line before `if` in rating test. Should have anticipated given the codebase uses wsl_v5 + gofumpt.

---

## E) WHAT WE SHOULD IMPROVE

### Process improvements

1. **Reason through CSS selector logic BEFORE implementing.** I moved fill classes to the label and added flex-row-reverse without tracing the `.peer:checked ~ .peer-checked\:*` selector through all 5 possible values. A 5-minute whiteboard trace would have caught the inversion immediately.

2. **Test edge cases, not just the middle.** `TestRatingKeyboardOrder` uses `Value: 3` (the one value that coincidentally renders correctly). Had I tested `Value: 1` and `Value: 5`, the fill inversion would have been obvious.

3. **String-assertion tests give false confidence for CSS-dependent behavior.** The test says "the class is present" — but CSS activation depends on DOM structure and selector direction. We need at least one visual or JS-execution test for critical interactive components.

4. **Read the AGENTS.md warnings BEFORE writing tests.** The ordered-class-substring mistake is explicitly documented. I should have internalized the `AssertContainsAll` pattern before writing any test that checks multiple Tailwind tokens.

### Architectural improvements

5. **The Rating fill problem exposes a broader gap:** components that rely on CSS-only interactivity (peer-checked, peer-hover) cannot handle bidirectional state. The library should either accept JS for these cases or document the limitation clearly.

6. **The shared menu-keyboard nav extraction is good DRY** but the singleton name namespace is ad-hoc (`tcMenuKeyNavAttached`, `tcCtxMenuAttached`, `tcCarouselAttached`, etc.). Consider a naming convention document.

---

## F) Up to 50 Things We Should Get Done Next

### Critical (fix the regression)

1. **Fix the Rating cumulative fill** — implement JS singleton for fill on `change`/`focus`/`hover`, OR switch to non-cumulative fill, OR use `has()` for a pure-CSS solution. This is the #1 priority.
2. **Add edge-case Rating tests** — test fill behavior for value=1 and value=5, not just value=3.
3. **Add a visual regression golden for Rating** in `visualtest/` — the fill inversion would have been caught by pixel comparison.
4. **Verify the Rating visual order** — is 5-on-left / 1-on-right the right UX? Standard ratings put 1 on the left. The original code had this backwards too (reverse DOM = 5 on left), and my flex-row-reverse preserves it. May need to flip to 1-on-left.

### Keyboard testing infrastructure

5. **Build chromedp keyboard test harness** in `visualtest/` — `input.DispatchKeyEvent` helpers, one test per interactive component.
6. **Add keyboard event simulation tests** for Dropdown, ContextMenu, Carousel, MobileMenu, Tooltip, Tabs.
7. **CSP nonce propagation test for ContextMenu** — verify the Shift+F10 script has `nonce`.
8. **Integration test: ContextMenu Shift+F10 opens menu and first menuitem gets focus.**

### ContextMenu improvements

9. **ContextMenu: Tab key should close the menu** (WAI-ARIA menu pattern: Tab closes, doesn't cycle within menu).
10. **ContextMenu: verify `e.target.closest('[data-tc-ctxmenu-trigger]')` works when focus is on a nested element inside the trigger.**
11. **ContextMenu: position the menu more intelligently for keyboard activation** (trigger bottom-left may overflow viewport; reuse the popover positioner).
12. **ContextMenu: add `aria-haspopup="menu"` on the trigger container.**
13. **ContextMenu: Enter/Space on a menuitem should activate it (native behavior, but verify for `<span>` non-link items).**

### Rating improvements

14. **Rating: add hover-fill preview** (highlight stars up to the hovered one) — requires JS or a different CSS approach.
15. **Rating: add keyboard `Home`/`End` support** (jump to 1 star / max stars).
16. **Rating: the `required` attribute is on radio value=1 only** — verify this is sufficient for HTML constraint validation across browsers.
17. **Rating: consider half-star support** (common in review UIs).
18. **Rating: the `sr-only` span says "N star(s)" — verify screen reader announces the radiogroup correctly with the label.**

### Carousel improvements

19. **Carousel: add `aria-live="polite"`** to announce slide changes to screen readers.
20. **Carousel: pause auto-advance on hover/focus** (if auto-advance is ever added).
21. **Carousel: the focus ring uses `ring-offset` — verify the offset color works on dark backgrounds.**
22. **Carousel: arrow buttons should also have `focus-visible` styling** (they're absolutely positioned, may overlap the ring).

### Shared menu nav improvements

23. **Extract a `tcRtlArrowKeys` helper** — the RTL arrow-key ternary is duplicated in Tabs, the shared menu nav, and Carousel.
24. **Extract a `tcFirstFocusable` helper** — the "focus first menuitem on toggle" pattern could be generalized.
25. **The menu nav `pageSize` calculation** (`Math.floor(items.length / 4)`) is arbitrary — consider a fixed step or make it configurable.
26. **Menu nav: character-key jump** (pressing "E" jumps to first item starting with "E") — WAI-ARIA menu pattern recommends this.

### Dropdown improvements

27. **Dropdown: the golden file now uses the compact shared JS** — verify the diff is clean and no old assertions break in other tests.
28. **Dropdown: verify that `aria-disabled` items are keyboard-skippable but still visible** (not hidden from AT).
29. **Dropdown: the `tabindex="-1"` on the menu container** + roving tabindex on items — verify Tab key exits the menu correctly.

### Testing improvements

30. **Add fuzz test for ContextMenu** — arbitrary item counts, disabled states, empty items.
31. **Add contract test: all `role="menu"` containers** use the shared keyboard nav helper.
32. **Golden test for ContextMenu WITH nonce** — the current golden has empty nonce so the script isn't captured.
33. **Add `TestMotionReduceCompliance` exemption check for new focus-visible rings** (they use transitions).
34. **Benchmark the shared menu nav JS** — ensure the singleton guard is fast.

### Documentation improvements

35. **Document the CSS `~` limitation** in an ADR — "why Rating can't have both correct arrow keys and cumulative fill with pure CSS."
36. **Update `docs/javascript-guide.md`** with the Rating fill case study.
37. **Update the skill SKILL.md** with the shared menu nav pattern.
38. **AGENTS.md: add a "CSS selector direction" warning** — always trace `~` and `+` combinators through all possible states before implementing.

### Demo / integration

39. **Verify the demo Rating page** renders correctly with the new DOM order + flex-row-reverse.
40. **Add a demo page for ContextMenu keyboard usage** (Shift+F10 hint text).
41. **Verify the recompiled demo CSS includes all new classes** (flex-row-reverse confirmed, but check ring-offset, has-* if added).

### Code quality

42. **The `pluralStars` helper in rating.templ** could use a named constant or `strings.Plural` (Go 1.26 has `strings.Plural`? — verify).
43. **The ContextMenu `<span>` for non-link items** renders a non-interactive element with `role="menuitem"` — consider rendering a `<button>` instead for keyboard activation.
44. **Consistent disabled styling** — ContextMenu now uses `aria-disabled` but Dropdown uses `[disabled]` attribute on buttons. Unify the pattern.
45. **The `utils.Ternary` for readonly vs interactive container class** in rating.templ is a long inline expression — extract to a helper for readability.

### Release readiness

46. **The Rating fill bug blocks release** — must be fixed before the next version cut.
47. **Run `nix run .#visual`** to verify no visual regressions across all components (requires Chromium).
48. **Run `nix flake check`** to verify formatting passes.
49. **Verify `scripts/check-lint-config.sh`** still passes after all changes.
50. **Review all BuildFlow auto-commit messages** from this session — they're likely hallucinated and don't mention the actual changes. The `[Unreleased]` CHANGELOG section is the authoritative record.

---

## G) Questions I CANNOT Figure Out Myself

### 1. Rating fill approach: JS singleton or non-cumulative or accept the tradeoff?

The fundamental CSS limitation (`~` only looks forward in DOM) means we cannot have BOTH correct arrow-key direction (forward DOM) AND correct cumulative star fill (needs reverse DOM) with pure CSS alone. The three options are:

- **(A) JS singleton** (~30 lines): listen for `change`/`hover`/`focus` on the radiogroup, fill stars cumulatively. Most production rating widgets do this. Adds JS to a previously JS-free component.
- **(B) Non-cumulative fill**: only the selected star is amber. Less pretty but honest and zero JS. Users see "the 3rd star is highlighted" not "3 stars filled."
- **(C) Accept inverted arrow keys**: revert to reverse DOM order for correct fill, accept that ArrowDown decreases the value. The original tradeoff.

Which approach do you want?

### 2. Rating visual direction: 1-on-left or 5-on-left?

The original code (reverse DOM, no flex-row-reverse) rendered star 5 on the LEFT and star 1 on the RIGHT. My code (forward DOM + flex-row-reverse) preserves this same visual. Standard UX conventions put 1 on the left (you read left-to-right, fill left-to-right). Should I flip to 1-on-left (remove flex-row-reverse, use forward DOM + non-cumulative or JS fill), or keep 5-on-left (current design)?

### 3. Chromedp keyboard test harness: build it now or defer?

The fill bug went undetected because all tests are string-assertion-only. Building a real keyboard-event test harness in `visualtest/` (~200-300 lines: `input.DispatchKeyEvent` helpers + per-component tests) would catch CSS-activation regressions like this. But it's a significant effort in a separate module. Should I build it as part of this keyboard-navigation work, or defer it to a separate testing-infrastructure task?

---

## Session Metrics

| Metric | Value |
|---|---|
| Files modified (source) | 6 (.templ + .go) |
| Files modified (generated) | 91 (*_templ.go) + 3 golden + 1 demo CSS |
| Files created | 1 (context_menu_test.go) |
| Tests added | 7 (Carousel focus, Rating order, ContextMenu ×4, updated Dropdown) |
| Tests passing | 15/15 packages |
| Lint issues | 0 |
| Golden files updated | 3 (rating_interactive, dropdown_basic, carousel_basic) |
| Known bugs introduced | 1 (Rating fill inversion — see section D) |
| BuildFlow auto-commits | ~5 (messages are hallucinated per AGENTS.md) |
| Round trips wasted on silly mistakes | 5 (guard name, readonly container, golden name, class assertion, wsl whitespace) |
