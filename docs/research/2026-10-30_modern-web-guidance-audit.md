# Modern Web Guidance audit — how modern are we, really?

> **Question:** how well does templ-components support modern web APIs?
> **Method:** the 192-guide Modern Web Guidance library (CLI v0.0.196, retrieved
> 2026-10-30, 44 guides pulled into context) checked against the actual component
> sources — not against our own docs. Every verdict below cites the file that
> decides it.
> **Supersedes:** `docs/research/modern-browser-capabilities.md` (2026-07-12),
> which proposed the migration. That migration **shipped**. This document grades
> the result and lists what the next tier of platform work is.

---

## 1. Verdict

**Grade: B+.** The library is well ahead of a typical component library and
ahead of its own documentation. Two numbers, because they answer different
questions:

| Axis | Score | What it measures |
| --- | --- | --- |
| **Correctness of what ships** | **~92%** | Where we use a native API, do we use it the way the platform (and the guide) requires? |
| **Breadth of platform adoption** | **~65%** | How much of what the platform now offers do we actually exploit? |

- **Zero legacy anti-patterns in the high-traffic components.** Modal, Drawer,
  CommandPalette, Lightbox, Popover, Dropdown, ContextMenu, Accordion, Tabs,
  Carousel, Table all ride native primitives (`<dialog>`, Popover API,
  `<details>`, scroll-snap, `content-visibility`) rather than div-and-JS
  reimplementations.
- **One legacy holdout remains:** `feedback.LoadingOverlay` is still a
  `div role="dialog" aria-modal="true"` (`feedback/loading.templ:107`) — no focus
  containment, no inert page, no Escape, page still scrollable and interactive
  behind it. It is the exact pattern the `dialog` guide exists to kill.
- **The gaps are almost all "opt-in enhancement not taken", not "wrong API"** —
  `:has()` (zero uses), scroll-driven animation (zero), scroll-state queries
  (zero), `scroll-initial-target` (zero), scrollbar theming (zero),
  `light-dark()`/`contrast-color()` (zero), Trusted Types (absent).

Across the 44 guides retrieved: **20 surfaces now sit in the adopted column,
12 are partial, 10 are still gaps**, and 2 are not gradeable against this
library: `menu` is a CLI stub with no prescription, and
`accessible-generated-content` has almost nothing to bite on here (our only
meaning-carrying generated content is the print-only external-link URL in
`.tc-prose`). Partial dominates because we usually get
the hard 80% right and skip the final guide-mandated detail — that is where the
remaining risk lives.

---

## 2. What is genuinely best-in-class

Worth recording so this is not re-litigated. Each row is verified in source.

| Guide | Evidence |
| --- | --- |
| `dialog` (core semantics) | `display/shared.templ:17` renders a real `<dialog>`; `display/shared.go:78` opens with `showModal()` and closes with `dialog.close()`. No `role="dialog"`, no `aria-modal`, no focus trap, no `inert` juggling, no Esc handler, no manual focus restore — all the things the guide marks DON'T. |
| `dialog` (dismiss) | `closedby="any"` on every overlay `<dialog>` (Modal, Drawer, CommandPalette, Lightbox) and pointer-position backdrop detection shared through `dialogBackdropOutsideJS`/`backdropDismissJS` (`display/shared.go`). Press **and** release must land outside, so a selection drag ending on the backdrop does not dismiss, and clicks on the dialog's own padding no longer close it. (Added 2026-10-30.) |
| `light-dismiss-a-dialog` | Same `closedby="any"`; the declarative form where the browser supports it, the pointer path everywhere else, with `close()` on a closed dialog being a no-op. |
| `animate-to-intrinsic-sizes` | `templates/custom.css` — `details[data-tc-accordion]::details-content` gets `block-size: 0` → `auto` under `@supports selector(::details-content)` plus the `:root` `interpolate-size` opt-in, with `content-visibility` transitioned `allow-discrete` and the whole block disabled under `prefers-reduced-motion`. The accordion panel animates instead of snapping. (Added 2026-10-30.) |
| `cross-document-transitions` (MPA surfaces) | `examples/demo/demo.css` and `website/site.css` declare `@view-transition { navigation: auto }` with a reduced-motion opt-out. Both are MPAs where every page opts in, which is the requirement. Deliberately not in `templates/custom.css`, so consumers do not inherit it. (Added 2026-10-30.) |
| `prevent-page-scroll-under-modal-dialog` | `templates/custom.css:104-111` — `overscroll-behavior: contain` on `dialog:modal` **and** `overflow: hidden` + `contain` on `dialog::backdrop`. (Added 2026-10-30; this was the one MANDATORY we were missing.) |
| `animate-to-from-top-layer` | All three MANDATORYs: `overlay` in the transition list, `display … allow-discrete`, and `@starting-style` (`templates/custom.css:113-131`), with `prefers-reduced-motion` handling. |
| `size-aware-styling` / `responsive-design` | 55 `@container` sites; the `ContainerAware` flag swaps every breakpoint map viewport↔container (ADR-0018), and `Card` deliberately opts out with a visual-regression proof. |
| `form-fields-automatically-fit-contents` | `forms/textarea.templ` — `field-sizing-content` with `min-h/max-h`, on by default. |
| `branded-select-styling` | `templates/custom.css:369` gates `appearance: base-select` + `::picker(select)` behind `@supports`; non-supporting browsers get the untouched native `<select>` — the guide's required degradation. |
| `accessible-error-announcement` + `validate-input-after-interaction` + `required-field-feedback` | `templates/custom.css:349-357` (`:user-invalid`/`:user-valid`, Widely available since 2023), plus `forms.ErrorAttrs` wiring `aria-invalid`/`aria-describedby` and `role="alert"` on FieldError/ValidationSummary. |
| `carousel-snap-highlights` / `scroll-snap-state-sync` | `display/carousel_go.go:56` feature-detects `scrollend` before syncing dots; snap lives in CSS (`snap-x snap-mandatory`), not in JS. |
| `navigation-drawer` | `<dialog>` modal shell with `data-side` driving logical-property transforms (`custom.css` `dialog.tc-drawer`), RTL-correct. |
| `responsive-table` (semantics half) | `display/table.templ` — `role="region"` + `tabindex="0"` + accessible name on the scroll wrapper, `scope="col"`, `aria-sort`, `sr-only` `<caption>`, sticky `thead`. |
| `optimize-image-priority` | `display/image.templ:105,132` — eager images get `fetchpriority="high"` (LCP), lazy images get **no** `fetchpriority` (the guide forbids pairing `low` with `loading="lazy"`), `decoding="async"`, explicit `width`/`height`, blur-up placeholder. |
| `ime-safe-enter-submit` | `forms/combobox.templ:227`, `forms/tags_input.templ:54`, `display/command_palette.go:131` — Enter/`,` ignored while `isComposing` (or keyCode 229) is set. Added 2026-10-30. |
| `progress-ring` neighbours | `feedback/circular_progress.templ` — `role="progressbar"` + clamped `aria-valuenow`. |
| `switch` | `forms/toggle.templ` — `role="switch"` + `aria-checked` on the native checkbox: switch semantics with native form submission and keyboard behaviour. `aria-checked` is required once the role is overridden (the HTML validator enforces it, and a switch whose state is only implicit reads as unchecked), so it is rendered from `Checked` and re-synced by a singleton change listener that also repairs bfcache/history restoration. The native `switch` attribute is Safari-only today, so this is the portable form. (Added 2026-10-30.) |
| Dark-mode plumbing | `layout/theme.templ` — FOUC-prevention inline script, `localStorage` in try/catch, `prefers-color-scheme` fallback, `role="switch"` + synchronised `aria-checked`; `color-scheme` on `:root`/`.dark`; `prefers-contrast`, `forced-colors`, `prefers-reduced-transparency` blocks (`custom.css:555-620`). |
| Motion discipline | `motion-reduce:*` on every transition/animation, guarded repo-wide by `utils.TestMotionReduceCompliance`. |

Also adopted and easy to forget: `content-visibility: auto` for `Table.LazyRows`
(`custom.css:214-222`), `text-wrap: balance`/`pretty` (`custom.css:316-330`),
`accent-color` on every native control (`custom.css:333-344`),
`interpolate-size: allow-keywords` (`custom.css:494`), `@supports` gating for
every not-yet-Baseline feature, View Transitions via
`htmx.config.globalViewTransitions` (`htmx/view_transitions.go:51`), and
`corner-shape: squircle` as an opt-in.

---

## 3. Partial — the guide's final mile is missing

Ordered by user-visible impact.

| # | Guide | What we do | What the guide additionally requires |
| --- | --- | --- | --- |
| 1 | `declarative-dialog-popover-control` | `data-tc-overlay-trigger` click delegation (`display/shared.go:92`) | `command="show-modal"` + `commandfor` — no JS at all (Baseline 2025-12-12). Our path is the documented polyfill-equivalent, so this is a payload/behaviour simplification, not a bug. |
| 2 | `viewport-aware-dropdowns` + `menu` | Popover API + a ~60-line JS positioner with flip and viewport clamping (`display/shared.go:265`) | CSS anchor positioning (`position-area` + `position-try-fallbacks`). **Not a defect:** the guide itself reports anchor positioning as *Limited availability* and its own fallback is exactly `getBoundingClientRect()` math. The right move is an additive `@supports (anchor-name: --x)` layer, not a replacement. Note our flip resolves one axis per open and cannot express "flip both axes". |
| 3 | `position-aware-tooltips` | Tooltip is `absolute` inside its trigger's parent (`display/tooltip.templ:84`) | Edge-aware placement. Today a tooltip in an `overflow-hidden` ancestor is clipped and a near-edge tooltip overflows the viewport — the popover panels handle this, tooltips do not. |
| 4 | `interest-triggered-tooltips` | Pure CSS `:hover`/`:focus-within` + `aria-describedby` propagation to the focusable child | `interestfor` + `popover="hint"` (Limited availability) gives the 0.5s interest delay and browser-owned AT wiring. Our JS propagation is a reasonable stand-in. |
| 5 | `hovercard` | `role="tooltip"` on rich, interactive content, `absolute`-positioned (`display/hover_card.templ:57`) | `popover="hint"` + `interestfor`; the browser assigns `aria-details` (not `aria-describedby`) when the card is interactive, and the guide says **do not** set `role="tooltip"` yourself. Our current combination advertises a plain tooltip while shipping interactive content, and z-index/overflow applies. This is the weakest overlay in the library. |
| 6 | `responsive-table` (stacking half) | Sticky `thead` only | Sticky row headers, `@container` (not `@media`) breakpoints, and the stacked-card layout with `::before` label injection on narrow containers. |
| 7 | `persistent-toast-notifications` | `aria-live` container + `role="status"` per toast (`feedback/toast.templ:92,108`) | `popover="manual"` (MANDATORY in the guide) so toasts survive unrelated interactions and escape z-index/overflow stacking; `popovertargetaction="hide"` close buttons; `sibling-index()` spacing. Deviation is defensible (a live region is what announces reliably), but stacking/clipping inside `overflow-hidden` layouts is the price. |
| 8 | `progress` | `div role="progressbar"` + clamped `aria-valuenow` (`feedback/progressbar.templ`) | Native `<progress>` with `appearance: none` + `::-webkit-progress-value`/`::-moz-progress-bar` (never comma-joined). Native gets implicit semantics and platform behaviour for free; ours is the hand-roll the guide steers away from. |
| 9 | `slider` | Native `type="range"` + vendor pseudo-element classes | Explicit WebKit filled-track (`--control-value` synced on `input`), and optionally `progress()`/`attr()` where supported. Cosmetic, not behavioural. |
| 10 | `dark-mode-toggle` | `.dark` class + `localStorage` + `role="switch"` | The guide's data model is `<meta name="color-scheme">` mapped back to CSS with `:root:has(> head > meta[…])`, three states (system/light/dark), and a `storage` event listener. Ours has two states and **no OS fallback without JS** — if the theme script fails, a dark-OS user gets the light theme. |
| 11 | `typography` | `text-wrap: balance` on headings, `pretty` on `p`, rem sizes | `clamp()` fluid type is only in `.tc-fluid-*` helpers, not the components; no `text-box-trim` (heading spacing is padding-based). |
| 12 | `fluid-scaling` | `.tc-fluid-*` clamp utilities + container breakpoints | Container-query **units** (`cqi`/`cqw`) — zero uses; fluid type is still viewport-driven inside a container-aware component. |

---

## 4. Gap — modern platform features we do not use at all

| Guide | Status in repo | Cost of the gap |
| --- | --- | --- |
| `style-parent-with-has` / `child-state-based-styling` / `content-based-styling` | **zero** `:has()` (`0` hits repo-wide; `:is()` used) | Everything is done with `peer` + sibling variants. That works only when the thing you style is a *following sibling* — which is why the Toggle thumb needs `peer-checked:*:inset-s-N` gymnastics and why FieldError/help-text reveal needs extra wrapper glue. `:has()` (Widely available 2023-11-02) would express several of these directly. |
| `scroll-progress-indicator` / `scrollytelling` / `scroll-entry-exit-effects` | **zero** `animation-timeline` | Reading-progress and reveal effects still require JS scroll listeners (we simply do not ship them). Additive opportunity with `@supports (animation-timeline: scroll())`. |
| `state-aware-sticky-headers` / `scrollability-affordance-hints` | **zero** container scroll-state queries | Sticky headers/table wrappers cannot style themselves when stuck, and horizontal-overflow tables have no fade hint. Limited availability (Chrome/Edge only) — opt-in only. |
| `scroll-target-on-load` | **zero** `scroll-initial-target` | Calendar/DateRange "start on today" needs JS today. |
| `customize-scrollbar-color-and-thickness` / `adapt-scrollbar-to-contrast-preferences` | **zero** (only `scrollbar-none`) | Dark mode cannot theme scrollbars; `color-scheme` covers most of it implicitly. |
| `contrast-color` / `color` (modern color spaces) | **zero** `contrast-color()`, no `light-dark()` | Badges/avatars over user-injected background colors must hardcode a text color. `--alpha()` (deliberate, ADR-0008) compiles to `color-mix()` inside `@supports`, so the color-mix half is effectively covered. |
| `trusted-types` | **zero**; CSP is nonce-based only | Nonce CSP blocks injected `<script>`, but not `innerHTML` sinks. We build DOM with `createElement` in the tag/copy scripts, so the exposure is small; adding `require-trusted-types-for 'script'` would need an audit of the demo/website. |
| `respect-os-text-scale` | rem everywhere (good) | No `<meta name="text-scale">` support; fine for most consumers. |
| `search-hidden-content` | `<details>` yes, `hidden="until-found"` no | Find-in-page cannot reveal collapsed content. Applies if we ever ship a collapsible search surface. |
| `visually-stable-font-fallbacks` | none | No `font-size-adjust`/`ascent-override`; only matters if consumers swap webfonts. |

**Deliberate non-adoption, do not reopen:** Web Components / Shadow DOM
(ADR-0033), and component-level module extraction (ADR-0040). Canvas/WebGL,
passkeys, built-in AI, WebMCP, WASM, PWA origin migration, Temporal-based date
math, and local-network-access have no surface in this library (the charts are
pure SVG). `menu` turned out to be a **stub in the CLI** — a header comment and
a TODO list, no prescription — so dropdown-menu grading rests on
`viewport-aware-dropdowns` and the `dialog` guide's Rule 3.

---

## 5. Two data-quality warnings for the next session

1. **The CLI's per-browser tables are not trustworthy in isolation.** The
   `viewport-aware-dropdowns` guide reports anchor positioning as
   "Supported by: Safari 27 (Sep 2026). Unsupported in: Chrome, Edge, and
   Firefox", which contradicts Chrome having shipped it years earlier. Treat the
   *status line* ("Limited availability" ⇒ keep a fallback) as load-bearing and
   re-verify browser data at the source (MDN/caniuse) before rewriting anything
   on the strength of the table. Same caution for `<input type="checkbox"
   switch>`, which the guide reports as Safari-only.
2. **Guide Baseline claims age.** `interpolate-size`, `field-sizing`,
   `closedby`, and Invoker Commands are all "Newly available" class features
   that we already gate behind `@supports` or feature detection — keep gating
   them; do not let a future "Baseline 2026" line remove the graceful
   degradation.

---

## 6. Ranked next actions (Pareto order)

| # | Action | Guide | Effort | Why first |
| --- | --- | --- | --- | --- |
| 1 | `LoadingOverlay`: decide what it IS before converting it. It is currently a `div role="dialog" aria-modal` with consumer-overridable positioning (the demo pins it with `Class: "absolute"` inside a frame), so a naive `<dialog>` swap would break that specimen and force a page-blocking modal on render. The likely split: a static `LoadingOverlay` (no dialog semantics, `aria-live`/`aria-busy` on the region) plus a modal `LoadingDialog` with an `Open` prop on `overlayShell`. | `dialog` | M | The last legacy modal in the library — no focus containment, page still interactive behind it, no Esc. It needs a design decision + a demo migration, not a mechanical swap, which is why it is still open. |
| 2 | HoverCard → `popover="hint"` + `interestfor`, drop `role="tooltip"`; keep the CSS hover path as the non-supporting fallback | `hovercard`, `interest-triggered-tooltips` | M | Fixes the wrong-role/wrong-association problem AND the overflow-clipping problem in one move. |
| 3 | `@supports (anchor-name: --x)` enhancement layer over `popoverPositionJS` | `viewport-aware-dropdowns` | M | Keeps the JS fallback (still required) while giving supporting browsers declarative two-axis flip. |
| 4 | `popover="manual"` toasts (+ `popovertargetaction="hide"`) | `persistent-toast-notifications` | M | Escapes overflow/z-index clipping; keeps the live region for announcement. |
| 5 | `<progress>` behind `ProgressBar` with `appearance: none` styling | `progress` | M | Deletes the most ARIA-heavy hand-roll in `feedback`. |
| 6 | `:has()` sweep: Toggle/Checkbox/RadioGroup/FieldError glue → parent-state selectors | `style-parent-with-has` | M | Removes variant gymnastics that the scanner-safety rules made necessary. |
| 7 | Meta `color-scheme` data model + `storage` event + OS fallback when JS fails | `dark-mode` / `dark-mode-toggle` | M | The only *robustness* gap in the theming layer. |

---

## 7. Shipped by this audit (2026-10-30)

1. **Scroll containment** for every `<dialog>` overlay —
   `overscroll-behavior: contain` on `dialog:modal`, `overflow: hidden` +
   `contain` on `dialog::backdrop` (`templates/custom.css:89-111`). Both
   MANDATORY in `dialog` and `prevent-page-scroll-under-modal-dialog`.
2. **`fetchpriority` correction** in `display.Image` — `low` is no longer
   emitted for lazy images; eager images keep `high` (`display/image.templ:105,132`).
3. **IME-safe Enter** in the combobox, tags input, and command palette
   (`forms/combobox.templ:227`, `forms/tags_input.templ:54`,
   `display/command_palette.go:131`).
4. **`role="switch"`** on the `forms.Toggle` input (`forms/toggle.templ:92`) —
   screen readers announce a switch, not a checkbox.
5. **Backdrop dismissal by pointer position** for every overlay `<dialog>` —
   `display/shared.go` (`dialogBackdropOutsideJS` + `backdropDismissJS`), used by
   Modal/Drawer, CommandPalette, and Lightbox, with `closedby="any"` on all
   four. Fixes clicks on the dialog's own padding closing it, and stops a
   selection drag that ends on the backdrop from dismissing.
6. **Accordion panel animation** — a `details[data-tc-accordion]::details-content`
   block-size transition under `@supports selector(::details-content)`, with a
   reduced-motion opt-out (`templates/custom.css`). Scoped to the component so a
   consumer's own `<details>` is untouched.
7. **`forms.InputProps.Autocomplete`** + `AutocompleteType` (23 tokens across
   credentials, address, payment) that drops unknown tokens instead of emitting
   them — autofill was previously reachable only through `BaseProps.Attrs`.
8. **Cross-document view transitions** for the two MPA surfaces (demo and
   website) — `@view-transition { navigation: auto }`, reduced-motion gated.
9. **Browser proof for the two overlay/accordion changes** (the wired-component
   e2e rule): `visualtest/overlay_backdrop_e2e_test.go` drives a real Modal with
   real mouse input — an inside click (including a synthetic `pointerdown`/
   `pointerup` pair targeted at the `<dialog>` with inside coordinates, the exact
   shape a backdrop click has) keeps it open, a backdrop press+release closes it,
   the close button closes it, and neither press-inside/release-outside nor
   press-outside/release-inside dismisses it.
   `visualtest/accordion_animation_e2e_test.go` samples the `<details>` height
   every animation frame: an intermediate height exists while opening, and under
   `prefers-reduced-motion: reduce` the very first frame is already the settled
   height.
10. **The demo's accordion answers are rendered** — the demo passed them as
   templ children, but `Accordion` renders `AccordionItem.Content`, so every
   panel opened empty; the answers now ride `Content`. Without that fix the new
   animation e2e would have been measuring an empty panel.
11. **Deterministic visual goldens** — the demo links Google-hosted webfonts
   with `font-display: swap`, which made every capture race the font swap: route
   goldens flaked at ~1% pixel mismatch on a different subset of pages each run
   (and the "load flake" note in AGENTS.md only covered part of it). The capture
   browser now fails requests to `fonts.googleapis.com`/`fonts.gstatic.com` and
   waits for `document.fonts` in every pixel-comparison path
   (`visualtest/external_fonts.go`). All 53 demo route goldens were re-recorded;
   two consecutive verify runs are green.

Goldens (HTML + website pages), the ADR-0045 class inventory, the demo
`app.css`, `templates/styles.css`/`theme.out.css`, the starter CSS copy, the
`cmd/tc/_sources` mirror, and the enum counts in six docs were regenerated in
the same change.
