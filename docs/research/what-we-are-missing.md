# What We're Missing — Hidden Gaps & Unknown Unknowns

> **Deep research into what templ-components does NOT yet leverage.**
> Covers browser APIs we're unaware of, codebase gaps discovered via audit,
> competitive landscape analysis, and emerging standards.
>
> **Date:** 2026-07-12

---

## Executive Summary

A systematic audit revealed **11 concrete gaps** in the current codebase —
features that are Baseline (widely supported) but not yet used. Additionally,
**6 emerging APIs** are documented as future opportunities. The gaps fall into
four categories:

| Category        | Gaps Found                                                                         | Impact                |
| --------------- | ---------------------------------------------------------------------------------- | --------------------- |
| **Performance** | `decoding="async"`, `fetchpriority`, `<link rel="preconnect">`, Speculation Rules  | Core Web Vitals       |
| **Mobile UX**   | `enterkeyhint`, `inputmode`, `field-sizing: content`                               | Mobile usability      |
| **Typography**  | `text-wrap: balance`, `text-wrap: pretty`                                          | Visual polish         |
| **CSS Power**   | `:user-valid`/`:user-invalid`, `interpolate-size`, `light-dark()`, `forced-colors` | Code reduction + a11y |

---

## Part 1: Codebase Gaps (Things We Should Already Have)

### 1.1. Image: Missing `decoding="async"` + `fetchpriority`

**Status:** Baseline (widely supported since 2021)

The Image component at `display/image.templ` emits `loading="lazy"` but has
**no `decoding` attribute and no `fetchpriority`**.

```html
<!-- Current -->
<img src="..." loading="lazy" />

<!-- Should be -->
<img src="..." loading="lazy" decoding="async" fetchpriority="low" />

<!-- For above-fold images (Lazy=false) -->
<img src="..." loading="eager" decoding="async" fetchpriority="high" />
```

**Why it matters:**

- `decoding="async"` prevents image decoding from blocking the main thread
  during initial render — images paint without jank
- `fetchpriority="high"` on above-fold images tells the browser to prioritize
  the LCP (Largest Contentful Paint) element
- `fetchpriority="low"` on below-fold images prevents them from competing
  with critical resources

**Impact:** Direct improvement to Core Web Vitals (LCP, FCP).

---

### 1.2. Forms: Missing `enterkeyhint` and `inputmode`

**Status:** Baseline (widely supported since 2021)

No form component emits `enterkeyhint` or `inputmode`. These are critical for
mobile UX:

```html
<!-- Search input should have: -->
<input type="search" enterkeyhint="search" inputmode="search" />

<!-- Email input should have: -->
<input type="email" enterkeyhint="next" inputmode="email" />

<!-- Tel input should have: -->
<input type="tel" enterkeyhint="next" inputmode="tel" />

<!-- URL input should have: -->
<input type="url" enterkeyhint="done" inputmode="url" />

<!-- Number input should have: -->
<input type="number" enterkeyhint="done" inputmode="decimal" />

<!-- Textarea in chat should have: -->
<textarea enterkeyhint="send"></textarea>
```

**Why it matters:**

- `enterkeyhint` changes the mobile keyboard's Enter key label/icon — "Search"
  on search inputs, "Next" in multi-field forms, "Send" in chat, "Done" on
  last field
- `inputmode` triggers the correct mobile keyboard layout — numeric keypad for
  numbers, email keyboard with @ key, tel keyboard with dial pad
- Both are **zero-cost progressive enhancement** — no effect on desktop

---

### 1.3. Missing `text-wrap: balance` and `text-wrap: pretty`

**Status:** Baseline since March 2024

No component or CSS file uses `text-wrap`. This is a one-line CSS improvement
for typography:

```css
/* Balanced headlines — no more orphaned words */
h1,
h2,
h3,
h4,
h5,
h6 {
  text-wrap: balance;
}

/* Clean paragraphs — no orphans at line end */
p {
  text-wrap: pretty;
}
```

**Why it matters:**

- `balance` makes headlines visually symmetric (binary-search algorithm finds
  the optimal line break points)
- `pretty` prevents single-word orphans at the end of paragraphs
- Zero performance cost for `balance` (limited to 6 lines)

---

### 1.4. Missing `field-sizing: content` for Textarea

**Status:** Baseline 2026 (newly available)

The Textarea uses a static `rows` attribute. `field-sizing: content` enables
**auto-growing textareas** — no JavaScript needed:

```css
textarea {
  field-sizing: content;
  min-height: 2.5rem;
  max-height: 20rem;
}
```

**Why it matters:**

- Textarea grows as the user types, shrinks as they delete — native, smooth,
  no JS measurement
- Eliminates the need for auto-resize JS libraries
- `max-height` ensures scrolling kicks in for very long content

---

### 1.5. Missing `:user-valid` / `:user-invalid` CSS

**Status:** Baseline since November 2023

Forms show error styling only via server-side `Error` prop + `aria-invalid`.
The `:user-valid` / `:user-invalid` pseudo-classes enable **CSS-only validation
feedback after user interaction**:

```css
/* Only style after user interacts — no premature errors */
input:user-invalid {
  border-color: rgb(220 38 38); /* red-600 */
}
input:user-valid {
  border-color: rgb(22 163 74); /* green-600 */
}
```

**Why it matters:**

- `:invalid` fires immediately on page load (shows red on empty required fields
  before user touches them) — terrible UX
- `:user-invalid` waits until the user blurs the field or submits the form
- This is **progressive enhancement** — works alongside the server-side error
  system, not replacing it

---

### 1.6. Missing `<link rel="preconnect">` for CDN

The Base layout injects HTMX from `cdn.jsdelivr.net` but never preconnects:

```html
<!-- Should be in <head> before the script tag -->
<link rel="preconnect" href="https://cdn.jsdelivr.net" crossorigin />
<link rel="dns-prefetch" href="https://cdn.jsdelivr.net" />
```

**Why it matters:**

- `preconnect` performs DNS lookup + TLS handshake + TCP connection in parallel
  with page render — saves 100-300ms on the first HTMX load
- `dns-prefetch` is the fallback for browsers that don't support `preconnect`

---

### 1.7. Missing `<meta name="color-scheme">`

The layout sets `color-scheme` via CSS but doesn't emit the meta tag. The meta
tag enables the browser to apply the correct color scheme **before CSS loads**:

```html
<meta name="color-scheme" content="light dark" />
```

**Why it matters:**

- Prevents flash of wrong color scheme before CSS parses
- Enables `light-dark()` CSS function to work
- Browser can render native form controls (scrollbars, date pickers) in the
  correct scheme immediately

---

### 1.8. Missing `interpolate-size: allow-keywords`

**Status:** Chrome 129+, Firefox/Safari in progress

This one-liner on `:root` enables **animating to/from `height: auto`** —
eliminating the `grid-rows-[1fr]`/`grid-rows-[0fr]` hack entirely:

```css
:root {
  interpolate-size: allow-keywords;
}

/* Now this just works */
details {
  height: 2.5rem;
  transition: height 0.3s ease;
  overflow: clip;
}
details[open] {
  height: auto;
}
```

**Why it matters:**

- Direct `height: auto` animation — no grid hack, no max-height hack
- Graceful degradation: unsupported browsers just snap (same as current)

---

### 1.9. Missing `forced-colors` Media Query Support

**Status:** Baseline since September 2022

Windows High Contrast mode forces a limited color palette and **nullifies
`box-shadow`**. Components that rely on shadows for visual definition (cards,
buttons, inputs) lose their boundaries:

```css
@media (forced-colors: active) {
  /* Add borders where shadows were removed */
  .card {
    border: 1px solid CanvasText;
  }
  /* Ensure focus rings use system colors */
  *:focus-visible {
    outline: 2px solid Highlight;
  }
}
```

**Why it matters:**

- Windows High Contrast users see invisible boundaries without this
- The browser forces colors automatically, but `box-shadow: none` means
  visual separation is lost
- **Accessibility compliance** — WCAG 1.4.3 Contrast (Minimum)

---

### 1.10. Missing `prefers-reduced-transparency` Support

**Status:** Baseline 2024 (Chrome 118+, Safari, Firefox)

Components use `backdrop-blur-xs` and `bg-black/50` for overlays. Users who
enable "Reduce Transparency" (macOS, iOS, Windows) need solid backgrounds:

```css
@media (prefers-reduced-transparency: reduce) {
  [data-tc-close] {
    backdrop-filter: none;
    background-color: rgb(0 0 0 / 0.8); /* more opaque */
  }
}
```

---

### 1.11. Missing `srcset` / `<picture>` First-Class Support

The Image component has no `SrcSet` or `Sizes` field — consumers must use
`Attrs` manually. Responsive images are a fundamental web performance feature:

```go
type ImageProps struct {
    // ... existing fields ...
    SrcSet string  // e.g. "img-480w.webp 480w, img-800w.webp 800w"
    Sizes  string  // e.g. "(max-width: 600px) 480px, 800px"
}
```

---

## Part 2: Emerging APIs (Future Opportunities)

### 2.1. CSS Anchor Positioning — Baseline 2025

**Browser support:** Chrome 125+, Firefox 147+, Safari 26+

Eliminates ALL JavaScript for tooltip/dropdown/popover positioning:

```css
.trigger {
  anchor-name: --my-anchor;
}
.tooltip {
  position-anchor: --my-anchor;
  position-area: top;
  position-try-fallbacks: flip-block, flip-inline;
}
```

**Current approach:** Absolute positioning with hardcoded `top-full`/`bottom-full`
classes. No auto-flip when viewport edge is near.

**Anchor positioning provides:**

- Automatic repositioning when near viewport edges (`position-try-fallbacks`)
- Sizing based on anchor dimensions (`anchor-size()`)
- Centering relative to anchor (`justify-self: anchor-center`)
- Works with Popover API in the top layer

---

### 2.2. `light-dark()` CSS Function — Baseline 2024

**Browser support:** Chrome 123+, Firefox 120+, Safari 17.5+

Could dramatically simplify dark mode. Instead of every component emitting
`text-gray-900 dark:text-gray-100`, the theme file could define:

```css
:root {
  color-scheme: light dark;
  --text-primary: light-dark(#111827, #f3f4f6);
  --text-muted: light-dark(#6b7280, #9ca3af);
  --bg-surface: light-dark(#ffffff, #111827);
}
```

**Current approach:** Every color has a `dark:` variant. The compliance test
(`TestDarkModeCompliance`) enforces this manually. `light-dark()` could
eventually replace the dual-class system for components that use semantic
CSS variables.

**Caveat:** Would require migrating from Tailwind utility classes to semantic
CSS variables — a large architectural shift. Better suited for v2.0.

---

### 2.3. Speculation Rules API — Chrome/Edge only (progressive enhancement)

**Browser support:** Chrome 121+, Edge 121+. Firefox/Safari ignore gracefully.

Enables **prerendering** of pages for near-instant navigation:

```html
<script type="speculationrules">
  {
    "prerender": [
      {
        "where": { "href_matches": "/*" },
        "eagerness": "moderate"
      }
    ]
  }
</script>
```

**Application for templ-components:** A `layout.SpeculationRules()` component
that consumers include once per page. Links are prerendered on hover (200ms),
making navigation feel instant. Unsupported browsers just navigate normally.

---

### 2.4. Declarative Shadow DOM — Baseline 2024

**Browser support:** Chrome 111+, Firefox 123+, Safari 16.4+

Enables **scoped CSS per component** without JavaScript:

```html
<my-card>
  <template shadowrootmode="open">
    <style>
      /* scoped to this component only */
    </style>
    <div class="card"><slot></slot></div>
  </template>
  Content here
</my-card>
```

**Relevance:** templ-components currently uses global Tailwind classes. DSD
would enable true component encapsulation — styles can't leak in or out.
However, this conflicts with the current Tailwind utility-class approach.
Better suited for a "scoped components" opt-in mode in v2.0+.

---

### 2.5. `commandfor` / `command` Attributes — Chrome only (experimental)

The successor to `popovertarget`. One button can control multiple targets:

```html
<button commandfor="popover1" command="toggle-popover">Toggle</button>
<button commandfor="dialog1" command="show-modal">Open Dialog</button>
```

**Advantage over `popovertarget`:** Works with `<dialog showModal()>`,
not just popovers. Still experimental — wait for Baseline.

---

### 2.6. `URL.canParse()` — Baseline 2024

Currently the library uses `url.Parse()` which returns an error. `URL.canParse()`
is a static method that returns a boolean — simpler validation:

```javascript
URL.canParse("https://example.com"); // true
URL.canParse("not a url"); // false
```

**Relevance:** Would simplify the Breadcrumb URL resolver. But this is
client-side JS — our URL parsing is server-side Go.

---

## Part 3: Competitive Landscape

### What shadcn/ui / Radix UI does that we don't

| Feature                  | shadcn/ui                       | templ-components              | Gap?                      |
| ------------------------ | ------------------------------- | ----------------------------- | ------------------------- |
| Form validation styling  | CSS `:invalid` + `aria-invalid` | Server-side `Error` prop only | Yes — add `:user-invalid` |
| Auto-growing textarea    | `field-sizing: content`         | Static `rows`                 | Yes                       |
| `enterkeyhint` on inputs | Not present                     | Not present                   | Both miss this            |
| Collapsible without JS   | Radix uses JS                   | We use `<details>` (better!)  | We're ahead               |
| `cmdk` command palette   | Full component                  | Not present                   | Consider adding           |
| Toast via `popover` API  | Uses Radix portal               | Custom `position: fixed`      | Could use `popover`       |

### What HTMX ecosystem does that we don't

| Feature          | HTMX ecosystem                         | templ-components | Gap?                   |
| ---------------- | -------------------------------------- | ---------------- | ---------------------- |
| `hx-preserve`    | Preserves element across swaps         | Not documented   | Document               |
| `hx-validate`    | HTML5 validation before submit         | Not used         | Add to Form            |
| `hx-disable-elt` | Disables elements during request       | Not used         | Consider               |
| `hx-include`     | Include extra form data                | Not documented   | Document               |
| `hx-params`      | Filter request parameters              | Not documented   | Document               |
| `hx-encoding`    | Multipart form uploads                 | Not documented   | Document for FileInput |
| ` boosted` links | `hx-boost` for progressive enhancement | Not documented   | Document               |
| `hx-history`     | History management                     | Not documented   | Document               |

---

## Part 4: Implementation Priority

### Immediate (trivial, high impact)

| Task                                          | Effort | Lines Changed |
| --------------------------------------------- | ------ | ------------- |
| `decoding="async"` + `fetchpriority` on Image | Tiny   | 2 lines       |
| `text-wrap: balance/pretty` in CSS            | Tiny   | 6 lines       |
| `<meta name="color-scheme">` in Base          | Tiny   | 1 line        |
| `<link rel="preconnect">` for CDN             | Tiny   | 2 lines       |
| `enterkeyhint` on Input types                 | Small  | ~15 lines     |
| `inputmode` on Input types                    | Small  | ~10 lines     |
| `field-sizing: content` utility class         | Tiny   | 3 lines       |
| `:user-valid`/`:user-invalid` CSS             | Small  | ~15 lines     |
| `interpolate-size: allow-keywords`            | Tiny   | 1 line        |
| `forced-colors` media query                   | Small  | ~10 lines     |
| `prefers-reduced-transparency`                | Tiny   | ~5 lines      |

### Medium term (component additions)

| Task                                                   | Effort |
| ------------------------------------------------------ | ------ |
| `SrcSet`/`Sizes` fields on ImageProps                  | Small  |
| `EnterKeyHint` field on InputProps/TextareaProps       | Small  |
| `InputMode` field on InputProps                        | Small  |
| Speculation Rules component                            | Medium |
| `hx-validate` on Form component                        | Tiny   |
| Document HTMX attributes (hx-preserve, hx-boost, etc.) | Medium |

### Long term (architectural)

| Task                                                   | Impact                        |
| ------------------------------------------------------ | ----------------------------- |
| Anchor Positioning for Tooltip/Dropdown/Popover        | Eliminates JS positioning     |
| `light-dark()` migration                               | Simplifies dark mode system   |
| Declarative Shadow DOM opt-in                          | Component style encapsulation |
| Popover API migration (Dropdown, Popover, ContextMenu) | ~134 lines JS eliminated      |
| `<dialog>` migration (Modal, Drawer)                   | ~150 lines JS eliminated      |
