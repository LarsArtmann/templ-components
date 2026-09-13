---
title: Dark Mode
description: Built-in dark mode with FOUC prevention and regression tests.
---

## Setup

1. Include `ThemeScript` before your content to prevent FOUC:

```templ
@layout.Base(layout.DefaultPageProps()) {
    @layout.ThemeScript("")
    {@content}
}
```

2. Add the `@custom-variant` directive in your CSS:

```css
@import "tailwindcss";
@custom-variant dark (&:where(.dark, .dark *));
```

3. Optionally add a `ThemeToggle`:

```templ
@layout.ThemeToggle("Toggle theme", "")
```

## How It Works

- `ThemeScript` runs before first paint, reads `localStorage`, and adds the `.dark` class to `<html>`.
- `color-scheme: light` on `:root`, `color-scheme: dark` on `.dark` — native form controls render correctly in both modes.
- `ThemeToggle` syncs across all instances on the page via `querySelectorAll`.

## Color Convention

| Context             | Light           | Dark                 |
| ------------------- | --------------- | -------------------- |
| Semantic background | `bg-blue-600`   | `dark:bg-blue-500`   |
| Semantic text       | `text-blue-600` | `dark:text-blue-400` |
| Neutral text        | `text-gray-500` | `dark:text-gray-400` |

Every neutral and semantic color class has a `dark:` variant — enforced by regression tests (`TestDarkModeCompliance`, `TestDarkModeSemanticColors`).

## LocalStorage Guard

`ThemeScript` and `ThemeToggle` both wrap `localStorage` access in `try/catch` to handle Safari private mode (`QuotaExceededError`).
