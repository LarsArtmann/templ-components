# ADR 0011: Dark Mode Color Convention

**Date:** 2026-07-08
**Status:** Accepted

## Context

The templ-components library uses Tailwind CSS v4's class-based dark mode (`@custom-variant dark (&:where(.dark, .dark *))`) toggled by `layout.ThemeScript()` + `layout.ThemeToggle()`. Every component must look correct in both light and dark mode.

A comprehensive audit (2026-07-08) found 30+ instances where components had light-mode color classes without corresponding `dark:` variants, causing invisible text, washed-out icons, and inconsistent focus rings in dark mode.

## Decision

### Color Shade Convention

| Context               | Light Mode          | Dark Mode               | Example                                        |
| --------------------- | ------------------- | ----------------------- | ---------------------------------------------- |
| Primary background    | `-600`              | `-500`                  | `bg-blue-600 dark:bg-blue-500`                 |
| Primary text          | `-600`              | `-400`                  | `text-blue-600 dark:text-blue-400`             |
| Danger background     | `-600`              | `-500`                  | `bg-red-600 dark:bg-red-500`                   |
| Danger text           | `-600`              | `-400`                  | `text-red-600 dark:text-red-400`               |
| Neutral text (high)   | `gray-900`          | `gray-100`              | `text-gray-900 dark:text-gray-100`             |
| Neutral text (medium) | `gray-700`          | `gray-200`              | `text-gray-700 dark:text-gray-200`             |
| Neutral text (low)    | `gray-500`          | `gray-400`              | `text-gray-500 dark:text-gray-400`             |
| Neutral text (muted)  | `gray-400`          | `gray-500`              | `text-gray-400 dark:text-gray-500`             |
| Surface               | `white` / `gray-50` | `gray-800` / `gray-900` | `bg-white dark:bg-gray-800`                    |
| Border                | `gray-200`          | `gray-700`              | `border-gray-200 dark:border-gray-700`         |
| Focus ring            | `blue-500`          | `blue-400`              | `focus:ring-blue-500 dark:focus:ring-blue-400` |
| Ring offset           | `white` (default)   | `gray-900`              | `dark:focus:ring-offset-gray-900`              |

### Palette Rule

All components use `gray-*` exclusively for neutral colors. No mixing of `slate-*`, `zinc-*`, `neutral-*`, or `stone-*`. This ensures consistent grayscale ramp across all components.

### Enforcement

Two failing tests block CI:

1. `utils.TestDarkModeCompliance` — scans all `.templ`/`.go` source files for neutral color classes (`text-gray-*`, `bg-white`, `bg-gray-*`, `border-gray-*`, `ring-gray-*`) without `dark:` variants.
2. `utils.TestDarkModeSemanticColors` — scans for semantic color classes (`bg-blue-600`, `text-red-600`, etc.) without `dark:` variants.

### Documented Exceptions

- **Toggle thumb** (`bg-white` in both modes): The track changes color instead; the thumb stays white.
- **SidebarNav** (permanently dark sidebar): Uses `bg-gray-900` as base; `hover:bg-gray-800` is a lighter shade on the dark sidebar.
- **Avatar silhouette icon** (`text-blue-200`): Decorative SVG inside a blue background; light-on-blue in both modes.

## Consequences

- Every new component MUST include `dark:` variants for all neutral and semantic colors.
- The compliance tests prevent regressions — adding a `text-gray-500` without `dark:text-gray-400` will fail CI.
- Consumers can override colors via Tailwind v4 `@theme` variables without touching component code.
- `color-scheme: light` on `:root` and `color-scheme: dark` on `.dark` ensures native form controls (scrollbars, checkboxes, date pickers) render correctly in both modes.

> **Note:** The `color-scheme: light` value on `:root` should be `color-scheme: light dark`
> to support consumers who use `prefers-color-scheme` instead of the `.dark` class.
> See [Dark Mode & Theming Research](../dark-mode-research.md) for the full analysis
> and the three first-class consumer dark mode paths.
