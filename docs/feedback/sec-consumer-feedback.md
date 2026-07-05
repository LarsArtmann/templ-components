# templ-components — SDK Feedback from SEC

**Consumer:** [SEC](https://github.com/larsartmann/sec) — dice-based game (Templ + HTMX + TailwindCSS v4)
**Date:** 2026-07-05
**Version used:** Not imported (evaluated, deferred)
**Session:** Evaluated for SEC's neon/dark game UI, deferred due to design-system mismatch

---

## Context: Why SEC evaluated templ-components

SEC uses Templ + HTMX + TailwindCSS v4 — the exact stack templ-components targets. The project has ~20 template files rendering game UI, dice animations, countdowns, finished states, history, analytics, simulation panels, and an islands mini-game map. We evaluated templ-components to replace hand-written components (buttons, cards, panels, form inputs, navigation).

---

## What's excellent about the library

### 1. `utils.BaseProps` pattern

The `BaseProps` struct (ID, Class, Style, Data attributes) embedded into every component's props is the right abstraction. It lets consumers override styling without the library guessing every possible CSS class. Clean composition.

### 2. Typed enum system

The `EnumTable[T]` pattern for typed enums (button variants, sizes, colors) with `String()`, `MarshalJSON()`, `UnmarshalJSON()`, `Display()`, `Emoji()`, `Label()` — all from a single `[]EnumEntry[T]` declaration — is excellent. SEC adopted this exact pattern for its own enums after seeing it in the skill docs.

### 3. CSP-safe icon paths

Icons as inline SVG path data (not external files, not inline `<script>`) is correct for strict CSP. The `IconPath(name)` function returns a path string that templates embed safely.

### 4. HTMX helper integration

Components accept HTMX attributes via props (`hxGet`, `hxPost`, `hxTarget`, `hxSwap`, etc.) rather than raw `hx-` HTML attributes. Type-safe, IDE-autocompletable.

### 5. Dark-mode awareness

Components have explicit dark-mode variants via Tailwind classes. The `dark:` prefix pattern is standard Tailwind, but the library applies it consistently across all components.

---

## Why SEC deferred adoption

### 1. Design system mismatch — neon dark aesthetic

SEC has a highly custom visual identity:

- Neon-on-black color scheme with CSS custom properties (`--neon-cyan`, `--neon-amber`, `--neon-red`, `--neon-pink`)
- Glassmorphism panels with backdrop-blur and semi-transparent backgrounds
- Custom glow effects (`--shadow-neon-amber`, `--shadow-neon-red`)
- Pixel-art game UI elements (dice faces, round dots, tension bars, streak calendars)
- Animated CSS-only effects (dice cup shake, countdown pulse, multiplier escalation)

templ-components uses standard Tailwind v4 utility classes with semantic color tokens (`bg-surface`, `text-primary`, `border-border`). The mental model is "clean SaaS dashboard" not "neon game UI". Adopting templ-components would mean fighting the default styling on every component, or wrapping each one in a custom override layer.

**Impact:** The cost of overriding exceeds the cost of maintaining hand-written components, because SEC's components are deeply customized (dice animations, tension bars, streak heatmaps) and don't have templ-components equivalents anyway.

### 2. No game-UI components

SEC needs: dice cup with shake animation, round-progress dots, tension bar with gradient levels, streak calendar heatmap, multiplier display with escalation colors, countdown timer, game state transitions. templ-components has: buttons, cards, inputs, modals, tables, navigation, badges, alerts. Zero overlap with SEC's actual component needs.

### 3. Tailwind v4 CSS-variable model vs SEC's custom variables

SEC defines its own CSS variables in `input.css` under `@theme` and uses them throughout. templ-components also uses Tailwind v4's `@theme` tokens. But the token names and values are different (`--color-surface` vs `--neon-cyan`). Adopting templ-components would require either:

- Renaming all SEC variables to match templ-components' token names
- Creating a mapping layer (fragile)
- Overriding every component's classes (defeats the purpose)

---

## What would make SEC adopt templ-components

### 1. Headless component variants

If templ-components offered "headless" variants — components that provide behavior and accessibility (ARIA, keyboard nav, focus trapping) without any styling — SEC could apply its own neon classes while getting the structural/semantic value. This is the Radix UI / Headless UI model.

**Example:** A headless modal that handles focus trapping, escape-to-close, backdrop click, and ARIA attributes, but renders zero CSS classes. SEC wraps it with glassmorphism + neon styling.

### 2. Game-UI component pack

A separate package (`templ-components-game`?) with:

- Progress dots (filled/unfilled/error states)
- Heatmap calendar (GitHub-style contribution graph)
- Stat display (label + value + trend indicator)
- Badge/chip with icon
- Animated counter
- Timer/countdown display

These are generic enough for any game UI, not just SEC.

### 3. CSS-variable aliasing layer

A config step where consumers map their existing CSS variable names to templ-components' expected tokens:

```go
templcomponents.ThemeMap{
    templcomponents.Surface: "neon-panel",
    templcomponents.Primary: "neon-cyan",
    templcomponents.Border:  "neon-border",
}
```

This would let SEC adopt without renaming its variables.

### 4. Slot/render-prop pattern for custom internals

Components that accept a `templ.Component` for their inner content (not just text/props) would let SEC inject custom game UI into a structured shell. E.g., a `Card` that accepts a `Header`, `Body`, and `Footer` as templ components rather than strings.

---

## Ideas for improvement (general)

### 1. Component count vs component depth

The library has many components, each with moderate customization. Consider whether fewer components with deeper customization (render props, slots, composition) would serve more use cases. The current model works for SaaS dashboards but not for visually distinctive apps.

### 2. Template literal helpers

SEC's `game_helpers.go` has functions like `tensionLevel()`, `outcomeColor()`, `roundDotClass()` — domain-specific presentation logic. templ-components could provide a pattern (not necessarily code) for how to structure domain-specific helpers that integrate with the component system.

### 3. Animation utilities

SEC has CSS-only animations (dice shake, countdown pulse, multiplier glow). A set of animation utility classes or a templ `<Animation>` component would be useful for any interactive UI.

---

## Overall verdict

templ-components is well-built — the `BaseProps` pattern, typed enums, CSP-safe icons, and HTMX integration are all correct decisions. For a SaaS/dashboard/admin-panel project using Tailwind v4, I'd recommend it without hesitation.

For visually distinctive projects (games, creative portfolios, branded marketing sites) with custom design systems, the library's opinionated styling creates more friction than value. The gap isn't quality — it's target audience. The library targets "clean, consistent, professional UI" but SEC needs "distinctive, branded, animated game UI."

The path to adoption for projects like SEC is either headless components (behavior without styling) or a CSS-variable aliasing layer. Either would significantly expand the library's addressable audience without changing its value proposition for existing consumers.

SEC did adopt one pattern from templ-components: the `EnumTable[T]` generic enum system. We saw it in the skill docs and implemented the same pattern for our 10 enum types. That alone was worth reading the skill docs.
