# Theme Bridge: Custom Semantic Palettes

## Problem

templ-components emits standard Tailwind color classes (`bg-blue-600`, `bg-white`,
`text-gray-900`). If your app uses a **custom semantic palette** (e.g. `bg-surface`,
`bg-accent`, `text-primary`), the library's hardcoded color names seem like a blocker.

They aren't. Tailwind v4's `@theme` directive lets you **remap** standard color tokens
to your custom values. The library components keep working unchanged.

## The pattern: remap, don't replace

In your `app.css`, override the specific Tailwind color variables that the library
uses. Point them at your semantic tokens:

```css
@import "tailwindcss";

/* Your custom semantic palette */
@theme {
  --color-surface: #1a1b23;
  --color-elevated: #24252e;
  --color-accent: #6366f1;
  --color-subtle-border: #2d2e38;
}

/* Bridge: redirect library colors to your palette */
@theme {
  /* Library cards/tables use bg-white → your surface color */
  --color-white: var(--color-surface);

  /* Library primary buttons/badges use bg-blue-600 → your accent */
  --color-blue-600: var(--color-accent);
  --color-blue-500: var(--color-accent);
  --color-blue-400: color-mix(in srgb, var(--color-accent) 80%, white);
  --color-blue-50: color-mix(in srgb, var(--color-accent) 10%, var(--color-surface));
  --color-blue-900: color-mix(in srgb, var(--color-accent) 30%, black);

  /* Library text colors */
  --color-gray-900: #f3f4f6; /* primary text (dark bg) */
  --color-gray-700: #d1d5db; /* secondary text */
  --color-gray-500: #9ca3af; /* muted text */
  --color-gray-300: #4b5563; /* borders */
  --color-gray-200: #374151; /* subtle borders */
  --color-gray-100: var(--color-elevated);
  --color-gray-50: var(--color-elevated);

  /* Dark mode variants (same as above — both modes use the dark palette) */
}
```

### Dark mode with dual palettes

If your app supports both light and dark mode, scope the bridge overrides:

```css
@layer base {
  :root {
    /* Light mode: let library defaults work */
  }

  .dark {
    /* Dark mode: redirect library colors to your dark palette */
    --color-white: var(--color-surface);
    --color-gray-900: #f3f4f6;
    --color-gray-700: #d1d5db;
    --color-gray-500: #9ca3af;
    --color-gray-200: #374151;
    --color-gray-50: var(--color-elevated);
    --color-blue-600: var(--color-accent);
  }
}
```

## Which colors does the library use?

The most common hardcoded colors across all components:

| Tailwind token                             | Used by                              | Maps to (example)       |
| ------------------------------------------ | ------------------------------------ | ----------------------- |
| `bg-white` / `dark:bg-gray-900`            | Card backgrounds, Table body, Nav    | Your surface color      |
| `bg-gray-50` / `dark:bg-gray-800`          | Table headers, secondary backgrounds | Your elevated color     |
| `bg-blue-600` / `dark:bg-blue-500`         | Primary buttons, badges, links       | Your accent color       |
| `text-gray-900` / `dark:text-white`        | Primary text                         | Your primary text color |
| `text-gray-500` / `dark:text-gray-400`     | Muted/secondary text                 | Your muted text color   |
| `border-gray-200` / `dark:border-gray-700` | Card/table borders                   | Your border color       |
| `text-blue-600` / `dark:text-blue-400`     | Links, active states                 | Your accent color       |
| `bg-green-600` / `bg-red-600`              | Success/error states                 | Keep or remap           |
| `text-amber-600`                           | TrendWarn                            | Keep or remap           |

## Quick adoption checklist

If your app currently reimplements library components with custom class helpers:

1. **Don't fork the components.** Use `@theme` to bridge colors.
2. **Audit which tokens you need.** Grep your current CSS for the colors you use.
3. **Add the bridge block** to your `app.css` `@theme` section.
4. **Adopt library components incrementally** — start with non-interactive ones
   (Card, Badge, Grid), then forms (Input, Select), then interactive ones (Modal,
   Dropdown).
5. **Verify dark mode** — render with and without `.dark` class.

## Why this works

Tailwind v4 generates CSS from `@theme` tokens. When a component emits
`bg-white`, Tailwind looks up `--color-white`. If you've set
`--color-white: var(--color-surface)`, every `bg-white` in the library now renders
your surface color — with zero Go code changes.
