# ADR: CSS-Variable Portability Layer

**Status:** Proposed
**Date:** 2026-06-27
**Decision:** Adopt semantic CSS-variable tokens + color-bridge CSS file

## Context

`templ-components` emits hardcoded Tailwind utility classes (`bg-white dark:bg-gray-800`,
`text-gray-900`, etc.). This couples consumers to Tailwind CSS. Projects like
`cqrs-htmx/adminui` that use custom CSS (not Tailwind) cannot adopt any tc
component except `icons`.

**Goal:** Let non-Tailwind projects use tc components by including a standalone
CSS file — no Tailwind build, no Node.js.

## Analysis

### Class inventory

tc components use **178 unique class tokens**:

- **88 color-related** (`bg-*`, `text-*`, `border-*`, `ring-*`, etc.)
- **90 layout/spacing** (`flex`, `grid`, `px-4`, `items-center`, etc.)
- Many use `dark:` prefix variants (e.g., `dark:bg-gray-900`)

### Options evaluated

| Option                                                                              | Effort    | Non-Tailwind works? | Verdict                                                                               |
| ----------------------------------------------------------------------------------- | --------- | ------------------- | ------------------------------------------------------------------------------------- |
| **A. CSS bridge file** — map all 178 classes to plain CSS                           | Medium    | Yes (if complete)   | Fragile — must track every new class. Layout classes change across Tailwind versions. |
| **B. Semantic class refactor** — replace `bg-white` with `bg-tc-surface` everywhere | Very High | Yes                 | Clean long-term, but 73-component refactor is a multi-day effort.                     |
| **C. Inline style mode** — components emit `style="background:var(--tc-surface)"`   | High      | Yes                 | Breaks `utils.Class()` merge, loses Tailwind dark: variant, ugly HTML.                |
| **D. Color-bridge + document Tailwind requirement for layout**                      | Low       | Partially           | Pragmatic. Color classes become CSS vars; layout classes need Tailwind or manual CSS. |
| **E. Do nothing**                                                                   | Zero      | No                  | Status quo. adminui uses icons-only.                                                  |

### Recommendation: **Option D (pragmatic bridge) → Option B (long-term)**

**Phase 1 (now):** Create `templ-components-colors.css` that:

1. Defines `--tc-*` CSS variables for all colors (light + `.dark` overrides)
2. Maps each color class tc uses to its variable (using escaped `dark\:` selectors)
3. Consumers override `--tc-*` variables to retheme — no Go code changes

This immediately enables non-Tailwind projects to get **correct colors** from tc
components. Layout (flex, padding, grid) still requires either Tailwind or manual
CSS definitions — documented as a known limitation.

**Phase 2 (future):** Progressively migrate components from hardcoded Tailwind
colors to semantic classes (`bg-tc-surface`, `text-tc-primary`). Provide a
`templ-components.css` that defines these semantic classes in plain CSS. This
makes tc fully Tailwind-free.

### Why not Option A (full bridge)?

A full 178-class bridge would essentially reimplement a Tailwind subset in plain
CSS. It would break every time a new class is added, and layout classes
(`flex`, `grid-cols-3`, `sm:px-6`) are tightly coupled to Tailwind's engine
(responsive prefixes, arbitrary values). Maintaining this would be worse than
the problem it solves.

### Why not Option C (inline styles)?

Inline styles bypass `utils.Class()` (which uses tailwind-merge-go for conflict
resolution). They also can't express `dark:` variants without JS. And they make
the HTML output verbose and hard to override.

## Consequences

- **Immediate:** Non-Tailwind consumers include `templ-components-colors.css`
  and get correct light/dark colors from all components. They define their own
  layout CSS or accept unstyled layout.
- **Future:** Semantic class migration makes tc progressively more portable.
  Each component migrated reduces the Tailwind coupling by one unit.
- **Risk:** Low. The color bridge is additive — existing Tailwind consumers
  are unaffected.
