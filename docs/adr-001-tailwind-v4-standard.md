# ADR-001: Tailwind CSS v4+ as the Standard

**Status:** Accepted
**Date:** 2026-06-27
**Decision:** Full Tailwind CSS v4+ (latest) for all projects. No CSS-variable
portability layer. Small custom CSS only where Tailwind doesn't cover something.

## Context

`templ-components` emits Tailwind utility classes. A cross-project analysis of
`cqrs-htmx/adminui` revealed that adminui uses a hand-rolled 706-line CSS design
system (CSS custom properties, no Tailwind) and therefore cannot adopt tc
components beyond the `icons` package.

A "CSS-variable portability layer" was proposed — mapping tc's Tailwind classes
to `--tc-*` CSS variables so non-Tailwind projects could use components without
a Tailwind build. A 1388-line color bridge CSS file was generated.

## Decision

**Reject the portability layer. Adopt Tailwind CSS v4+ as the standard for all
projects.**

### Rationale

1. **Tailwind v4 is CSS-first and build-free.** v4's `@import "tailwindcss"` and
   `@source` directives work with zero JavaScript at runtime. The build step
   (`@tailwindcss/cli` or Lightning CSS) is a single binary, no Node.js required.

2. **Maintaining two systems is worse than migrating.** A portability layer
   means tracking every new class in parallel plain-CSS definitions. It breaks
   on every responsive prefix, arbitrary value, or plugin class. The maintenance
   cost grows with every component added.

3. **Semantic class migration is a multi-day dead end.** Replacing
   `bg-white` with `bg-tc-surface` across 73 components doesn't add features —
   it adds indirection. Tailwind v4's `@theme` block already solves theming
   cleanly by overriding `--color-*` variables.

4. **The right answer for adminui is migration, not accommodation.** adminui
   should adopt Tailwind v4+ (CSS-first, no Node.js) and progressively replace
   its hand-rolled CSS. This eliminates the maintenance of a 706-line custom
   design system and unlocks every templ-components component for free.

5. **Small custom CSS is fine.** Where Tailwind doesn't cover something (rare
   animations, complex selectors, third-party widget styling), a small custom
   `.css` file is the right tool. This is not a portability layer — it's
   pragmatic augmentation.

## Consequences

- **templ-components stays Tailwind-only.** Components emit standard Tailwind
  utility classes. Consumers override via `@theme { --color-*: ... }`.
- **No `templ-components-colors.css` bridge file.** Removed.
- **`icons` package remains CSS-agnostic** (pure SVG data). This is a natural
  property of icons, not a portability strategy.
- **Recommendation for all projects:** adopt Tailwind v4+. CSS-first config,
  `@source` scanning, `@theme` for overrides. No DaisyUI, no Node.js runtime.
- **adminui path forward:** migrate from hand-rolled CSS to Tailwind v4+ in a
  future session. Start by wrapping `admin.css` tokens in a `@theme` block,
  then progressively replace custom classes with utilities.
