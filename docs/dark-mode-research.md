# Dark Mode & Theming Research

**Date:** 2026-07-10
**Status:** Research — informs documentation and packaging, not a decision record

Triggered by consumer feedback from cqrs-htmx ([feedback](feedback/2026-07-10_cqrs-htmx-consumer-feedback.md)),
which reported friction between the library's class-based dark mode strategy and
adminui's `prefers-color-scheme` + CSS-variable design system. This document
researches all available dark mode mechanisms in Tailwind v4 and modern CSS,
reasons from first principles about which is best for a component library, and
proposes concrete documentation/packaging changes.

Related: [ADR 0011 (dark mode convention)](adr/0011-dark-mode-convention.md),
[ADR 0008 (semantic tokens)](adr/0008-semantic-tokens.md),
[ADR 001 (Tailwind v4 standard)](adr-001-tailwind-v4-standard.md),
[Tailwind v4 Adoption Guide](tailwind-v4-adoption-guide.md).

---

## The three CSS-level mechanisms

Modern CSS offers three independent mechanisms for dark mode. Each solves a
different part of the problem:

| Mechanism               | What it does                                                                                                            | Baseline | Tailwind v4?                           |
| ----------------------- | ----------------------------------------------------------------------------------------------------------------------- | -------- | -------------------------------------- |
| `prefers-color-scheme`  | Media query: "the user's OS wants dark"                                                                                 | Jan 2020 | **Tailwind's default `dark:` variant** |
| `color-scheme` property | Tells the browser which schemes the page supports; drives native UI rendering (scrollbars, form controls, canvas color) | Jan 2022 | Consumer CSS                           |
| `light-dark(a, b)`      | Returns `a` or `b` based on the _computed_ `color-scheme` — no media query, no variant, no class                        | May 2024 | Not expressible as a utility class     |

### `prefers-color-scheme` (media query)

The oldest mechanism. The browser exposes the user's OS-level preference:

```css
@media (prefers-color-scheme: dark) {
  .card {
    background: #1e293b;
  }
}
```

This is a **read-only signal** — the page can _detect_ the preference but cannot
_change_ it. It works with zero JavaScript and zero configuration. The user gets
dark mode if and only if their OS is set to dark.

### `color-scheme` (CSS property)

This is the foundation that the other two build on. It tells the browser which
color schemes the page (or element) can comfortably render in:

```css
:root {
  color-scheme: light dark;
}
```

Setting this makes the browser automatically adjust native UI: scrollbars,
checkboxes, date pickers, `<input>` backgrounds, selection colors, and the canvas
(default page background). **Without this, native form controls render light even
in dark mode** — the most common cause of "dark mode looks broken on form fields."

It can be scoped per-element:

```css
header {
  color-scheme: only light;
} /* always-light header */
main {
  color-scheme: light dark;
} /* follows the page mode */
```

### `light-dark()` (CSS function)

The newest option. Instead of writing two media-query blocks or two class
variants, you write a single declaration:

```css
:root {
  color-scheme: light dark;
}

.card {
  background: light-dark(white, oklch(0.279 0.041 260.031));
  color: light-dark(oklch(0.208 0.042 265.755), white);
}
```

`light-dark()` reads the **computed** `color-scheme` (influenced by both
`prefers-color-scheme` and explicit `color-scheme` declarations) and returns the
first or second value. It works with **any** dark mode strategy because it reads
`color-scheme`, which both the `.dark` class and `prefers-color-scheme` set.

**Limitation for this library:** `light-dark()` is not expressible as a Tailwind
utility class. Adopting it would require writing custom CSS for every color in
every component, abandoning the utility-class approach entirely. Baseline only
since May 2024.

Sources: [MDN: light-dark()](https://developer.mozilla.org/en-US/docs/Web/CSS/color_value/light-dark),
[MDN: color-scheme](https://developer.mozilla.org/en-US/docs/Web/CSS/color-scheme)

---

## The two application strategies

On top of those CSS mechanisms, there are two ways to _apply_ mode-switching in a
Tailwind v4 component library:

### Strategy A — Explicit `dark:` variants (what templ-components does today)

Every color class gets its dark counterpart written explicitly:

```html
<div class="bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700"></div>
```

- The `dark:` variant maps to `prefers-color-scheme` **by default** in Tailwind v4
- A consumer can override it to a `.dark` class via `@custom-variant`
- Requires ~200 `dark:` variants across 26+ components (enforced by compliance tests)
- Each component is mode-aware — the component author picks the dark shade

### Strategy B — Token indirection via `@theme`

Components reference semantic tokens, and the token's _value_ changes per mode:

```css
/* Consumer's CSS */
@theme inline {
  --color-surface: var(--my-surface);
}

:root {
  --my-surface: white;
}
.dark {
  --my-surface: oklch(0.279 0.041 260.031);
}
/* or: */
@media (prefers-color-scheme: dark) {
  :root {
    --my-surface: oklch(0.279 0.041 260.031);
  }
}
```

```html
<!-- Component HTML — mode-blind -->
<div class="bg-surface border-border"></div>
```

- The component is completely unaware of dark mode — the token handles it
- No `dark:` variants needed
- Requires rewriting all component classes from `bg-blue-600` to `bg-tc-primary` etc.

### How `@theme` vs `@theme inline` works

This distinction is critical and widely misunderstood:

| Form                                         | Compiled CSS                           | Cascade-overridable?   | Global `:root` var emitted? |
| -------------------------------------------- | -------------------------------------- | ---------------------- | --------------------------- |
| `@theme { --color-x: #fff }`                 | `.bg-x { background: var(--color-x) }` | **Yes** — cascade wins | **Yes**                     |
| `@theme inline { --color-x: var(--my-var) }` | `.bg-x { background: var(--my-var) }`  | **No** — value inlined | **No**                      |

The `inline` form is essential for CSS-variable-based theming: it lets a utility
class resolve to a consumer-controlled variable that changes per mode.

Source: [Tailwind v4 theme docs](https://tailwindcss.com/docs/theme)

---

## Critical discovery: Tailwind v4's default IS `prefers-color-scheme`

This inverts the framing of the consumer feedback. From the
[official Tailwind v4 docs](https://tailwindcss.com/docs/dark-mode):

> **By default**, the `dark` variant is driven by the browser's
> **`prefers-color-scheme`** CSS media feature — meaning dark mode utilities
> activate automatically based on the user's OS/system setting, with no
> configuration required.

The `.dark` class strategy is an **opt-in override** via `@custom-variant`. But
`templ-components-theme.css` ships this override unconditionally (line 24):

```css
@custom-variant dark (&:where(.dark, .dark *));
```

This means every consumer who imports the theme file gets class-based dark mode
— even if they wanted OS-following. A consumer who uses `prefers-color-scheme`
and does NOT import the theme file gets dark mode **for free** from every `dark:`
variant in the library.

---

## First-principles analysis: what's best for a component library?

### The library's constraints

1. **Must work with any consumer dark mode strategy** — `.dark` class, `prefers-color-scheme`, or CSS variables
2. **Must be Tailwind-native** — ADR-001 mandates standard Tailwind utility classes
3. **Must be testable** — the compliance test suite enforces dark mode coverage
4. **Must not require JavaScript for basic dark mode** — some consumers want zero JS
5. **Must support custom color systems** — consumers with existing design tokens

### Evaluation matrix

| Criterion                             | Explicit `dark:` variants (current)           | Semantic tokens (`bg-tc-*`)         | `light-dark()`               |
| ------------------------------------- | --------------------------------------------- | ----------------------------------- | ---------------------------- |
| **Tailwind-native**                   | Yes — every Tailwind user understands `dark:` | No — custom `bg-tc-*` classes       | No — requires custom CSS     |
| **Strategy-agnostic**                 | Yes — `prefers-color-scheme` or `.dark` class | Yes — token abstracts mode          | Yes — reads `color-scheme`   |
| **Works with CSS-variable consumers** | Yes — via `@theme` palette override           | Yes — override the token            | Yes                          |
| **Works without JS**                  | Yes (with `prefers-color-scheme` default)     | Yes                                 | Yes                          |
| **Compliance-testable**               | Yes — existing `TestDarkModeCompliance`       | Not needed (no variants)            | Not applicable               |
| **Broad browser support**             | Baseline 2020                                 | Baseline 2017                       | Baseline 2024                |
| **ADR-001 alignment**                 | Aligns (standard Tailwind classes)            | Conflicts (rejected in ADR-001)     | Conflicts (non-Tailwind)     |
| **Maintenance cost**                  | Low (already built + tested)                  | High (token layer to maintain)      | High (rewrite everything)    |
| **Consumer mental model**             | "Tailwind classes, I already know this"       | "New `bg-tc-*` vocabulary to learn" | "Custom CSS for every color" |

### Conclusion: keep the architecture, fix the packaging

The current approach — **explicit `dark:` variants on standard Tailwind palette
classes** — is the correct design. It wins on every axis that matters for a
library. The consumer feedback pain points are real but their proposed solution
(surface token redesign) solves a documentation problem with an architecture
change.

ADR-001 already rejected the portability/semantic-token layer. ADR-0008 proposes
it as opt-in but deferred to v1.0. This research confirms that decision: the
`@theme` palette override already solves the CSS-variable consumer case without
requiring any component code changes.

---

## The three first-class consumer paths

These should be documented as equal options, not buried or conflated:

### Path 1: OS-following (zero config)

The simplest path. Works out of the box with Tailwind v4's default `dark:`
variant:

```css
/* app.css — that's it */
@import "tailwindcss";
@source "../vendor/github.com/larsartmann/templ-components";
```

Every `dark:bg-gray-800` in every library component activates automatically when
the user's OS is set to dark mode. No JavaScript, no `@custom-variant`, no
`.dark` class. Native form controls render correctly if the consumer also sets:

```css
:root {
  color-scheme: light dark;
}
```

**Best for:** blogs, docs sites, internal tools where the user never needs to
override their OS preference.

### Path 2: Toggle (user-controlled)

For apps that want a dark/light toggle button (using the library's
`layout.ThemeScript()` + `layout.ThemeToggle()`):

```css
@import "tailwindcss";
@source "../vendor/github.com/larsartmann/templ-components";

/* Override the dark variant to use a class instead of media query */
@custom-variant dark (&:where(.dark, .dark *));

:root {
  color-scheme: light;
}
.dark {
  color-scheme: dark;
}
```

```go
// server.go — renders <head> with FOUC-prevention script
layout.Base(layout.PageProps{
    // ThemeScript injects an inline <script> that reads localStorage
    // and adds the .dark class before first paint
})
```

**Best for:** SaaS apps, dashboards, any site where the user wants to pick.

### Path 3: CSS-variable design system (the cqrs-htmx/adminui case)

For consumers with an existing CSS-variable-based design system (e.g.,
`var(--surface)`, `var(--text)`, `var(--border)`). Instead of a fragile bridge
like `.bg-white { background: var(--surface) }`, map the library's palette to
consumer variables via `@theme`:

```css
@import "tailwindcss";
@source "../vendor/github.com/larsartmann/templ-components";

/* Map Tailwind palette slots to consumer design tokens.
 * Every bg-white, text-gray-*, border-gray-*, dark:bg-gray-* in
 * every library component now resolves to these variables. */
@theme {
  --color-white: var(--surface);
  --color-gray-50: var(--surface-2);
  --color-gray-100: var(--border);
  --color-gray-200: var(--border);
  --color-gray-400: var(--text-muted);
  --color-gray-500: var(--text-muted);
  --color-gray-700: var(--border-dark);
  --color-gray-800: var(--surface-dark);
  --color-gray-900: var(--bg);
  --color-gray-100: var(--text-dark);
}

:root {
  color-scheme: light dark;
}
```

Now both `bg-white` and `dark:bg-gray-800` resolve to consumer-controlled
variables that already switch per mode. The `dark:` variants become irrelevant
because both sides of the pair resolve to the same consumer variable. No bridge
hacks, no `.dark` class needed, no strategy mismatch.

**Best for:** projects with an existing CSS custom property design system that
want to adopt library components without rewriting their theming layer.

---

## Issues in the current packaging

### 1. `templ-components-theme.css` forces `.dark` class strategy

The theme file ships `@custom-variant dark (&:where(.dark, .dark *))`
unconditionally. Consumers who import it for the `@theme` color overrides also
get the class-based dark mode override, silently disabling OS-following.

**Fix:** Split the file. The `@custom-variant` line belongs in a separate
"toggle dark mode" snippet, not bundled with color overrides.

### 2. `color-scheme: light` on `:root` is wrong for dark mode

Line 27 of `templ-components-theme.css`:

```css
:root {
  color-scheme: light;
}
```

This tells the browser the page only supports light mode, so native form controls
and scrollbars render light even in dark mode. Should be:

```css
:root {
  color-scheme: light dark;
}
```

This lets the browser pick the correct scheme based on the active mode (whether
driven by `prefers-color-scheme` or the `.dark` class).

### 3. Adoption guide presents only one strategy

`tailwind-v4-adoption-guide.md` presents the `.dark` class strategy as the only
option. The `@theme` palette override pattern (Path 3 above) is buried in the
"Phase 1" migration section, not presented as a first-class theming option.

---

## What NOT to do

### Do not build a mandatory semantic token layer (`bg-tc-surface` everywhere)

ADR-001 already rejected this correctly. ADR-0008 proposes it as opt-in but
deferred. This research confirms: the token layer adds indirection, breaks the
"standard Tailwind classes" promise, and the problem it solves (CSS-variable
consumer compatibility) is already handled by the `@theme` palette override.

### Do not adopt `light-dark()` for component styling

Baseline only since May 2024, not expressible as a Tailwind utility class, and
would require rewriting every component's color approach. Revisit in 2027+ when
browser support is universal.

### Do not redesign surface tokens (`--tc-surface` as a mandatory abstraction)

The `@theme` palette override (`--color-white: var(--surface)`) already solves
this more elegantly. A mandatory abstraction layer forces every consumer through
an indirection that most don't need.

---

## Action items

1. **Split `templ-components-theme.css`** — separate the `@custom-variant dark`
   opt-in from the `@theme` color overrides, so importing color overrides doesn't
   silently change the dark mode strategy
2. **Fix `color-scheme: light` → `color-scheme: light dark`** on `:root`
3. **Rewrite the adoption guide's dark mode section** — present the three consumer
   paths as equal first-class options, with OS-following as the zero-config default
4. **Add the `@theme` palette override pattern** to the adoption guide as a
   top-level "Theming" section (currently buried in the migration section)
5. **Add a note to `ThemeScript`/`ThemeToggle` godoc** clarifying they're only
   needed for the toggle strategy, not for OS-following

---

## References

- [Tailwind CSS v4 — Dark Mode](https://tailwindcss.com/docs/dark-mode)
- [Tailwind CSS v4 — Theme Variables](https://tailwindcss.com/docs/theme)
- [MDN — light-dark() CSS function](https://developer.mozilla.org/en-US/docs/Web/CSS/color_value/light-dark)
- [MDN — color-scheme property](https://developer.mozilla.org/en-US/docs/Web/CSS/color-scheme)
- [ADR 0011: Dark Mode Color Convention](adr/0011-dark-mode-convention.md)
- [ADR 0008: Semantic Token Layer](adr/0008-semantic-tokens.md)
- [ADR 001: Tailwind CSS v4+ Standard](adr-001-tailwind-v4-standard.md)
- [Tailwind v4 Adoption Guide](tailwind-v4-adoption-guide.md)
- [cqrs-htmx Consumer Feedback](feedback/2026-07-10_cqrs-htmx-consumer-feedback.md)
