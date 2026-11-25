---
title: Theming
description: Customize colors and styles without touching component code.
---

## How Theming Works

Components emit standard Tailwind utility classes (`bg-blue-600`, `text-gray-900`). To customize colors, override the underlying CSS variables in your `@theme` block:

```css
/* app.css */
@import "tailwindcss";
@source "./templ-components-classes.txt";

@theme {
  /* Override the primary brand color globally */
  --color-blue-600: #4f46e5; /* changes ALL bg-blue-600, text-blue-600, etc. */
  --color-blue-500: #6366f1;
  --color-blue-400: #818cf8;
}
```

No Go code changes needed. Every component that uses `bg-blue-600` picks up the new color.

## Semantic Token Layer

For a higher-level theming API, `tc css sync` writes the semantic theme next to your Tailwind entry, at the version your `go.mod` pins:

```bash
GOFLAGS=-mod=mod go run github.com/larsartmann/templ-components/cmd/tc css sync
```

Import it in your CSS:

```css
@import "tailwindcss";
@source "./templ-components-classes.txt";
@import "./templ-components-theme.css";
```

The theme file provides semantic tokens like `bg-tc-primary-600`, `text-tc-danger-600` that you can use in your own CSS or Tailwind arbitrary values. Override them in your `@theme` block to re-skin every component:

```css
@theme {
  --color-tc-primary-600: #7c3aed;
}
```

## Dark Mode

Dark mode is built-in. Components use Tailwind `dark:` variants for all colors. Include `ThemeScript` to prevent flash of unstyled content:

```templ
@layout.ThemeScript("")
```

See the [Dark Mode guide](/guides/dark-mode) for the full setup.

### The warm-dark `dark:`-pair override

To retune dark mode per component (the pattern this site uses for its warm
dark surfaces), override the color pair EXPLICITLY — both modes, complete
literal classes:

```templ
@display.Card(display.CardProps{
    Class: "bg-white dark:bg-[#1a1712]",
    // ...
})
```

Two scanner rules make or break it:

- Variant prefixes must be COMPLETE literals in source. `"dark:bg-" + shade`
  assembled at runtime is invisible to the Tailwind scanner — the CSS never
  gets generated. Runtime-assembled classes need the shipped safelist
  (`templates/safelist.css`).
- Child targets use the arbitrary-child variant, also literal:
  `[&>details]:dark:bg-gray-900`.

For a library-wide dark retune, override the shade PAIR in `@theme` instead
(one block, every component — the pair-awareness gotcha is documented in the
Tailwind adoption guide).

## RTL / Internationalization

All components use CSS logical properties (`ms-`, `me-`, `ps-`, `pe-`, `start-`, `end-`, `text-start`, `border-s-`, `border-e`). They automatically mirror in RTL contexts when you set `dir="rtl"` on a parent element.
