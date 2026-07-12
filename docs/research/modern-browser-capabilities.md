# Modern Browser Capabilities for templ-components

> **Deep research synthesis** — How native browser APIs can eliminate custom
> JavaScript, improve accessibility, and reduce bundle size across the
> templ-components library.
>
> **Date:** 2026-07-12 · **Baseline target:** All features Baseline 2024+ (Chrome,
> Firefox, Safari, Edge).

---

## Executive Summary

The library ships **~690 lines of custom JavaScript** across 23 interactive
components. Modern browser APIs — many Baseline since 2022-2024 — can replace
the majority of this JS with **declarative HTML attributes** and **CSS-only
animations**:

| API                          | Baseline | Components Affected                        | JS Eliminated                   |
| ---------------------------- | -------- | ------------------------------------------ | ------------------------------- |
| `<dialog>` element           | 2022     | Modal, Drawer                              | ~150 lines                      |
| Popover API                  | 2024     | Dropdown, Popover, ContextMenu, MobileMenu | ~134 lines                      |
| `@starting-style`            | 2024     | All overlay/popup animations               | CSS replaces JS opacity toggles |
| `<details>`/`<summary>`      | Always   | Accordion                                  | ~26 lines                       |
| CSS scroll-snap              | 2023     | Carousel                                   | ~13 lines                       |
| View Transitions API         | 2023     | HTMX swaps, MPA navigations                | Enhancement                     |
| `content-visibility: auto`   | 2023     | Table, DataTable, long lists               | Perf boost                      |
| `:has()` selector            | 2023     | Toggle, Checkbox, RadioGroup, Tabs         | Simplification                  |
| `<search>` element           | 2023     | Nav search forms                           | Semantics                       |
| CSS scroll-driven animations | 2024     | ProgressBar (scroll-linked)                | Enhancement                     |
| `color-mix()`                | 2023     | Theming overrides                          | Flexibility                     |

**Potential JS reduction: ~320+ lines (46%)** while improving accessibility
(browser-native focus management, ARIA, and keyboard handling).

---

## 1. The `<dialog>` Element — Native Modal/Drawer

**Baseline:** 2022 (all major browsers). Universally safe.

### What it provides natively

| Feature                | Current custom JS                        | Native `<dialog>` equivalent            |
| ---------------------- | ---------------------------------------- | --------------------------------------- |
| Focus trap             | `overlayTrapJS()` — 28 lines per overlay | `showModal()` — browser-managed         |
| Backdrop dimming       | `<div class="bg-black/50">` + JS toggle  | `::backdrop` pseudo-element             |
| Escape key dismissal   | `keydown` listener in trap JS            | Built-in                                |
| Click-outside dismiss  | `[data-tc-close]` click delegation       | `<form method="dialog">` backdrop       |
| Top-layer rendering    | `z-50` + `pointer-events` toggling       | Automatic top-layer (immune to z-index) |
| `inert` on background  | Manual `setAttribute('inert', '')` in JS | Browser handles background inertness    |
| Focus restore on close | `data-tc-prev-focus` + manual restore    | Browser-managed                         |
| `aria-modal`           | Manual `aria-modal="true"`               | Implicit — no attribute needed          |

### How Modal would change

**Before** (current — ~75 lines JS):

```html
<div role="dialog" aria-modal="true" class="opacity-0 pointer-events-none ...">
  <div class="bg-black/50" data-tc-close="modal"></div>
  <div class="scale-95 opacity-0 ..." inert>...</div>
</div>
<script>
  /* 75 lines of focus trap + open/close/animation JS */
</script>
```

**After** (`<dialog>` — ~5 lines JS for HTMX compat):

```html
<dialog id="my-modal" class="backdrop:bg-black/50 ...">
  <form method="dialog"><button class="absolute top-2 end-2">X</button></form>
  ...content...
</dialog>
<script>
  // Only needed for HTMX-opened dialogs (declarative <dialog> opens via JS API)
  if (!window.tcDialogAttached) {
    window.tcDialogAttached = true;
    document.addEventListener("click", (e) => {
      const trigger = e.target.closest("[data-tc-dialog-open]");
      if (trigger) document.getElementById(trigger.dataset.tcDialogOpen)?.showModal();
    });
  }
</script>
```

### Benefits

- **~150 lines of focus-trap JS eliminated** (Modal + Drawer)
- **Better accessibility**: browser focus management is more robust than manual
  Tab-key trapping (handles Shadow DOM, iframes, custom elements)
- **No z-index wars**: top-layer rendering means overlays can never be stuck
  behind other content
- **Native inert**: background content is automatically non-interactive

### Migration Considerations

- `showModal()` requires JS call (no purely declarative open) — but the trigger
  JS is ~5 lines vs. current ~75
- `::backdrop` styling needs to go in CSS (app.css), not Tailwind classes
- Exit animations require `@starting-style` + `transition-behavior: allow-discrete`
- Drawer would use `<dialog>` with side-panel positioning (no native slide, but
  CSS transforms work on `<dialog>`)

### Recommendation

**Phase 2 migration** — high impact but changes rendered HTML structure. Keep the
Go API (`ModalProps`, `DrawerProps`) identical; only the `.templ` internals change.
Update golden tests. Consumer CSS targeting `role="dialog"` still works.

---

## 2. Popover API — Declarative Floating UI

**Baseline:** 2024 (popover attribute), 2025 (hint type). Safe for all modern
browsers.

### What it provides natively

| Feature                       | Current custom JS                              | Native Popover API                           |
| ----------------------------- | ---------------------------------------------- | -------------------------------------------- |
| Light dismiss (click-outside) | `document.addEventListener('click', ...)`      | `popover="auto"` — built-in                  |
| Escape key dismissal          | `keydown` listener checking for Escape         | Built-in for `popover="auto"`                |
| Top-layer rendering           | `z-20`/`z-50` + absolute positioning           | Automatic top-layer                          |
| Single-open stacking          | Manual close-all-others logic                  | `popover="auto"` auto-closes others          |
| Declarative trigger           | `data-dropdown-trigger` + JS listener          | `popovertarget="id"` — ZERO JS               |
| Backdrop                      | Manual overlay div                             | `::backdrop` pseudo-element                  |
| Open/close state              | `classList.toggle('hidden')` + `aria-expanded` | `:popover-open` pseudo-class + implicit ARIA |

### Components this replaces

| Component       | Current JS | Replacement                                             |
| --------------- | ---------- | ------------------------------------------------------- |
| **Dropdown**    | ~53 lines  | `popover="auto"` + `popovertarget` on button            |
| **Popover**     | ~40 lines  | `popover="auto"` + `popovertarget` on button            |
| **ContextMenu** | ~15 lines  | `popover="manual"` + JS only for positioning at cursor  |
| **MobileMenu**  | ~26 lines  | `popover="auto"` + `popovertarget` on hamburger         |
| **Tooltip**     | ~18 lines  | `popover="hint"` (CSS-only hover via `:hover` + anchor) |

### Dropdown Example

**Before** (~53 lines JS):

```html
<div class="relative inline-block">
  <button data-dropdown-trigger="dd1" aria-expanded="false">Actions</button>
  <div class="hidden absolute z-10 ..." role="menu" data-dropdown-menu="dd1">...</div>
</div>
<script>
  /* open/close + click-outside + Escape + arrow nav = 53 lines */
</script>
```

**After** (~0 lines JS for open/close, ~15 lines for arrow-key nav):

```html
<div class="relative inline-block">
  <button popovertarget="dd1-menu">Actions</button>
  <div popover="auto" id="dd1-menu" class="..." role="menu">...</div>
</div>
<script>
  // Only arrow-key navigation (everything else is native)
  if (!window.tcMenuNavAttached) {
    window.tcMenuNavAttached = true;
    document.addEventListener("keydown", (e) => {
      if (!e.target.closest('[role="menu"]')) return;
      // arrow up/down logic only — ~15 lines
    });
  }
</script>
```

### Benefits

- **~134 lines of JS eliminated** (Dropdown + Popover + ContextMenu + MobileMenu)
- **Declarative triggers**: `popovertarget` is pure HTML — no event listeners
- **Better stacking**: top-layer means popovers never conflict with `overflow: hidden`
  parents or z-index contexts
- **Light dismiss**: browser handles click-outside and Escape natively
- **Implicit ARIA**: `popovertarget` creates implicit `aria-expanded` + `aria-controls`

### Migration Considerations

- Arrow-key navigation still needs JS (WAI-ARIA APG requirement for menus)
- `popover="hint"` for tooltips is newer (may need `popover="manual"` fallback)
- CSS anchor positioning for smart placement is not yet Baseline (use absolute
  positioning + existing position maps for now)
- `::backdrop` styling needs CSS (can't use Tailwind utility classes on pseudo-elements)

### Recommendation

**Phase 1 migration** — straightforward, high impact. Dropdown and Popover
components are nearly identical to the Popover API's design. Keep Go API
identical; update `.templ` internals + golden tests.

---

## 3. `@starting-style` — CSS-Only Entry/Exit Animations

**Baseline:** 2024. Works in Chrome, Edge, Safari 17.5+, Firefox 129+.

### Problem it solves

CSS transitions don't fire when an element goes from `display: none` to visible,
or when it first appears in the DOM. This is why the library currently uses JS
to toggle opacity/transform classes — the animation can't be pure CSS.

`@starting-style` defines the "from" state for elements entering the DOM or
becoming visible, enabling **pure CSS entry/exit animations**.

### Application

Every overlay (Modal, Drawer, Dropdown, Popover, Tooltip) currently uses JS to:

1. Remove `hidden` class
2. Add `opacity-100` + remove `opacity-0`
3. Add `scale-100` + remove `scale-95`

With `@starting-style`, this becomes pure CSS:

```css
/* Entry animation */
[popover],
dialog {
  transition:
    opacity 200ms,
    transform 200ms,
    overlay 200ms allow-discrete,
    display 200ms allow-discrete;
}

[popover]:popover-open,
dialog:modal {
  opacity: 1;
  transform: scale(1);
}

/* Exit animation */
@starting-style {
  [popover]:popover-open,
  dialog:modal {
    opacity: 0;
    transform: scale(0.95);
  }
}

/* When closed (exit) */
[popover],
dialog {
  opacity: 0;
  transform: scale(0.95);
}
```

### Benefits

- **Eliminates JS opacity/transform toggling** for all overlay components
- **Smooth exit animations** (currently impossible with `display: none` removal)
- **Works with Popover API and `<dialog>`** natively

### Recommendation

**Phase 1** — add to `templates/app.css` as part of the CSS foundation. Works
even before component migrations (enhances current opacity-transition components).

---

## 4. `<details>`/`<summary>` — Native Accordion

**Baseline:** Always supported. Universally safe.

### What it provides natively

| Feature                | Current custom JS                        | Native `<details>`                            |
| ---------------------- | ---------------------------------------- | --------------------------------------------- |
| Toggle open/closed     | 26 lines of click delegation + classList | `open` attribute — zero JS                    |
| `aria-expanded`        | Manual set/remove in JS                  | Implicit via `open` attribute                 |
| Chevron rotation       | JS `classList.toggle('rotate-180')`      | CSS `:has([open])` or `[open] svg { rotate }` |
| Keyboard (Enter/Space) | Handled by `<button>`                    | Native `<summary>` handles it                 |
| Content animation      | JS grid-rows toggle                      | CSS `grid-rows` on `[open]`                   |

### How Accordion would change

**Before** (~26 lines JS):

```html
<div class="divide-y ...">
  <div>
    <button data-accordion-trigger="acc1" aria-expanded="true">Title</button>
    <div class="grid grid-rows-[1fr]" data-open>...</div>
  </div>
</div>
<script>
  /* toggle grid-rows + aria-expanded + rotate icon = 26 lines */
</script>
```

**After** (zero JS):

```html
<div class="divide-y ...">
  <details open>
    <summary class="flex items-center justify-between ...">
      <span>Title</span>
      <svg class="transition-transform motion-reduce:transition-none [[open]_&]:rotate-180" />
    </summary>
    <div class="grid grid-rows-[1fr] transition-all motion-reduce:transition-none">
      ...content...
    </div>
  </details>
</div>
```

### Benefits

- **26 lines of JS eliminated**
- **Better accessibility**: `<summary>` is natively focusable, keyboard-operable,
  and announced as a disclosure toggle by screen readers
- **Progressive enhancement**: works with zero JS (content accessible even if JS fails)
- **Browser-native state**: `open` attribute survives page reloads via bfcache

### Migration Considerations

- `<details>` renders a default disclosure triangle — hide with `summary { list-style: none; }`
  and `summary::-webkit-details-marker { display: none; }`
- The `grid-rows-[1fr]`/`grid-rows-[0fr]` animation technique still works — CSS
  `details[open] > div` selector replaces JS class toggle
- Tailwind v4 variant `[[open]_&]:rotate-180` targets the open state (needs the
  `open` attribute to be scannable — it is, as a bare attribute)

### Recommendation

**Phase 1** — simplest, safest migration. Zero breaking changes to Go API.

---

## 5. CSS Scroll-Snap — Native Carousel

**Baseline:** 2023. Universally safe.

### What it provides natively

| Feature          | Current custom JS                      | Native scroll-snap                       |
| ---------------- | -------------------------------------- | ---------------------------------------- |
| Slide navigation | 13 lines of translateX + dot update JS | `scroll-snap-type: x mandatory`          |
| Touch/drag swipe | Not supported                          | Native touch scrolling                   |
| Dot indicators   | JS class toggle                        | `scrollend` event + IntersectionObserver |
| Keyboard nav     | Not supported                          | Native tab + arrow keys on container     |

### How Carousel would change

**Before** (~13 lines JS):

```html
<div class="overflow-hidden">
  <div class="flex transition-transform" data-tc-carousel-track>
    <div class="w-full flex-shrink-0">Slide 1</div>
    <div class="w-full flex-shrink-0">Slide 2</div>
  </div>
  <button data-tc-carousel-prev>←</button>
  <button data-tc-carousel-next>→</button>
</div>
<script>
  /* translateX + dot classes = 13 lines */
</script>
```

**After** (~8 lines JS for dots sync only):

```html
<div class="overflow-x-auto snap-x snap-mandatory scroll-smooth">
  <div class="flex">
    <div class="w-full flex-shrink-0 snap-center">Slide 1</div>
    <div class="w-full flex-shrink-0 snap-center">Slide 2</div>
  </div>
</div>
<button onclick="carousel.scrollBy({left: width, behavior: 'smooth'})">→</button>
<script>
  // Dot sync via scrollend (8 lines)
  carousel.addEventListener("scrollend", () => {
    const idx = Math.round(carousel.scrollLeft / carousel.offsetWidth);
    dots.forEach((d, i) => d.classList.toggle("active-dot", i === idx));
  });
</script>
```

### Benefits

- **Native touch/drag support** (huge UX improvement on mobile)
- **Native momentum scrolling** on trackpads and touch devices
- **Keyboard accessible** (container is focusable with `tabindex="0"`)
- **Smooth animated scrolling** via CSS `scroll-behavior: smooth`

### Recommendation

**Phase 1** — significant UX improvement, especially for mobile/touch.

---

## 6. View Transitions API — Smooth HTMX Swaps

**Baseline:** SPA transitions 2023 (Chrome, Edge, Safari). Cross-document (MPA)
transitions 2024+ (Chrome only, Firefox/Safari in progress).

### Application for HTMX

HTMX swaps content via DOM replacement. Wrapping the swap in
`document.startViewTransition()` enables **smooth morphing animations** between
old and new states — cross-fades, shared-element transitions, custom animations.

```javascript
// Global HTMX hook — add once per page
document.body.addEventListener("htmx:beforeSwap", (e) => {
  if (!document.startViewTransition) return; // graceful degradation
  e.preventDefault();
  document.startViewTransition(() => {
    htmx.process(e.detail.target);
    e.detail.target.innerHTML = e.detail.xhr.responseText;
  });
});
```

### Shared-element transitions

Named elements morph between old and new states:

```css
/* A card that morphs during HTMX navigation */
.card-featured {
  view-transition-name: featured-card;
}
```

### Benefits

- **Perceived performance**: transitions mask network latency
- **Spatial continuity**: elements visually "move" to their new position
- **Zero library dependency**: native browser API
- **Graceful degradation**: older browsers just do instant swap (current behavior)

### Recommendation

**Phase 1** — additive helper, no breaking changes. Ship as `htmx.ViewTransitions(nonce)`
component that injects the swap hook.

---

## 7. `content-visibility: auto` — Rendering Performance

**Baseline:** 2023. Safe for all modern browsers.

### Problem it solves

Long tables, lists, and data-heavy pages render all content immediately, even
off-screen items. `content-visibility: auto` tells the browser to skip rendering
for off-screen elements until they approach the viewport.

### Application

```css
/* Table rows below the fold */
tbody tr {
  content-visibility: auto;
  contain-intrinsic-size: auto 48px; /* estimated row height */
}

/* Long card grids */
.grid-card {
  content-visibility: auto;
  contain-intrinsic-size: auto 200px;
}
```

### Benefits

- **2-5x faster initial render** on pages with 100+ items
- **Lower memory usage**: browser skips layout/paint for off-screen content
- **`contain-intrinsic-size: auto Xpx`**: learns actual size after first render,
  preventing scrollbar jitter

### Recommendation

**Phase 1** — add as utility classes or CSS foundation. Apply to Table, Grid,
Pagination output. No JS changes, no API changes.

---

## 8. `:has()` Selector — CSS-Only State Detection

**Baseline:** 2023 (all major browsers).

### Applications

```css
/* Toggle: style container based on checkbox state */
.toggle-container:has(input:checked) .toggle-thumb {
  transform: translateX(100%);
}

/* Form field: red border when it has an error */
.form-field:has(.field-error) input {
  border-color: rgb(220 38 38);
}

/* Card: hover effect when it contains a focused element */
.card:has(:focus-visible) {
  box-shadow: 0 0 0 2px rgb(59 130 246);
}

/* Table row: highlight when it contains a selected checkbox */
tr:has(input:checked) {
  background-color: rgb(239 246 255);
}

/* Tabs: show panel when its radio is checked */
.tab-group:has(#tab1:checked) [aria-labelledby="tab1"] {
  display: block;
}
```

### Benefits

- **Eliminates JS class toggling** for parent-state styling
- **Works without JS** (progressive enhancement)
- **Better than JS**: reacts to DOM state changes instantly, no listener needed

### Recommendation

**Phase 1** — add to CSS foundation. Some components (Toggle, Checkbox) can
gradually adopt `:has()` for state styling.

---

## 9. `<search>` Element — Semantic Search Landmark

**Baseline:** 2023. Universally safe.

```html
<!-- Before -->
<div class="relative">
  <input type="search" />
</div>

<!-- After -->
<search class="relative">
  <input type="search" />
</search>
```

### Benefits

- **Screen reader landmark**: users can navigate to search via landmark menu
- **No ARIA needed**: `<search>` implies `role="search"` natively
- **SEO**: search engines identify the search functionality

### Recommendation

**Phase 1** — wrap search inputs in `<search>` element. Minimal change.

---

## 10. CSS Scroll-Driven Animations

**Baseline:** 2024 (Chrome, Edge, Safari TP). Firefox in progress.

### Application: Reading Progress Bar

```css
.reading-progress-bar {
  animation: grow-progress linear;
  animation-timeline: scroll(root);
}

@keyframes grow-progress {
  from {
    transform: scaleX(0);
  }
  to {
    transform: scaleX(1);
  }
}
```

### Application: Reveal-on-Scroll

```css
.scroll-reveal {
  animation: reveal linear;
  animation-timeline: view();
  animation-range: entry 10% cover 30%;
}

@keyframes reveal {
  from {
    opacity: 0;
    transform: translateY(40px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
```

### Benefits

- **Zero JS** scroll-linked animations (currently requires IntersectionObserver)
- **Main-thread independent**: runs on compositor thread (smooth 60fps)
- **Graceful degradation**: unsupported browsers show static content

### Recommendation

**Phase 2** — add as utility classes for opt-in use. Not yet universal enough
for core components.

---

## 11. `color-mix()` — Dynamic Theming

**Baseline:** 2023. Universally safe.

```css
/* Dynamic hover color from a single base */
.button:hover {
  background-color: color-mix(in srgb, var(--color-primary) 85%, black);
}

/* Dynamic transparency from base color */
.overlay {
  background-color: color-mix(in srgb, var(--color-primary) 10%, transparent);
}
```

### Benefits

- **Eliminates pre-computed color variants** (blue-500, blue-600, blue-700)
- **Single source of truth**: define one primary color, derive shades dynamically
- **Works with CSS custom properties**: consumers set `--color-primary` and all
  shades derive automatically

### Recommendation

**Phase 1** — document in theming guide. Don't change existing components (they
use Tailwind's color system), but show consumers how to use `color-mix()` for
custom overrides.

---

## Implementation Roadmap

### Phase 1: Foundation + Low-Risk Wins (immediate, non-breaking)

| Task                              | Effort | Impact | Breaking?        |
| --------------------------------- | ------ | ------ | ---------------- |
| CSS foundation in `app.css`       | Small  | Medium | No               |
| Accordion → `<details>`           | Small  | Medium | No (same Go API) |
| Carousel → scroll-snap            | Medium | High   | No (same Go API) |
| View Transitions helper           | Small  | Medium | No (additive)    |
| `content-visibility` utilities    | Tiny   | Medium | No               |
| `<search>` element in Nav         | Tiny   | Low    | No               |
| Document `color-mix()` + `:has()` | Tiny   | Low    | No               |

### Phase 2: Popover API Migration (medium risk, changes HTML)

| Task                            | Effort | Impact | Breaking?                   |
| ------------------------------- | ------ | ------ | --------------------------- |
| Dropdown → Popover API          | Medium | High   | HTML changes (golden tests) |
| Popover component → Popover API | Medium | High   | HTML changes                |
| ContextMenu → Popover API       | Small  | Medium | HTML changes                |
| MobileMenu → Popover API        | Medium | Medium | HTML changes                |
| Tooltip → `popover="hint"`      | Medium | Medium | HTML changes                |

### Phase 3: Dialog Migration (higher risk, major refactor)

| Task                         | Effort | Impact | Breaking?    |
| ---------------------------- | ------ | ------ | ------------ |
| Modal → `<dialog>`           | Large  | High   | HTML changes |
| Drawer → `<dialog>`          | Large  | High   | HTML changes |
| Remove overlay focus-trap JS | Small  | High   | JS removal   |

### Phase 4: Advanced (when Baseline universal)

| Task                                | Waiting for                      |
| ----------------------------------- | -------------------------------- |
| CSS anchor positioning for popovers | Firefox support (Baseline 2025+) |
| Cross-document View Transitions     | Firefox/Safari support           |
| `popover="hint"` for tooltips       | Universal Baseline               |

---

## Browser Support Strategy

The library targets **Baseline 2024** — features that work across all major
browsers (Chrome, Firefox, Safari, Edge) in their latest versions.

All recommended features include **graceful degradation**:

- `<dialog>`: falls back to standard `display: none` div
- Popover API: `popovertarget` buttons do nothing (progressive enhancement)
- `@starting-style`: elements appear without animation (instant show)
- `<details>`: works in all browsers (Baseline since forever)
- scroll-snap: falls back to normal scrolling
- View Transitions: falls back to instant swap
- `content-visibility`: falls back to normal rendering
- `:has()`: falls back to unstyled state (ensure base styles are correct without it)

### Feature Detection

For components that need a JS fallback when a browser API is unavailable:

```javascript
// Popover API
if (Object.hasOwn(HTMLElement.prototype, 'popover')) { ... }

// View Transitions
if (document.startViewTransition) { ... }

// <dialog>
if (typeof HTMLDialogElement !== 'undefined') { ... }
```

---

## Impact Summary

| Metric                         | Current     | After Phase 1-3    | Change |
| ------------------------------ | ----------- | ------------------ | ------ |
| Custom JS lines                | ~690        | ~340               | -51%   |
| Components with custom JS      | 23          | ~12                | -48%   |
| Focus trap implementations     | 1 (manual)  | 0 (browser-native) | -100%  |
| Click-outside dismiss handlers | 5           | 0                  | -100%  |
| Escape key handlers            | 6           | 1 (Combobox)       | -83%   |
| Top-layer rendering            | 0 (z-index) | 8 components       | Native |
| Touch/drag carousel            | No          | Yes (scroll-snap)  | New    |
| View transitions on HTMX       | No          | Yes                | New    |
