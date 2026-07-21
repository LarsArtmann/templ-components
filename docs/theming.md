# Theming templ-components

The library uses standard Tailwind CSS v4 color utilities (`bg-blue-600`,
`text-red-600`, etc.) so you can re-skin everything from a single CSS file
without touching any Go code.

There are **three** ways to override colors, in order of recommended use.

---

## 1. Semantic token layer (recommended)

Import `templ-components-theme.css` from your `app.css`. This file aliases
every Tailwind palette color used by the library to a semantic name
(`--color-tc-primary`, `--color-tc-danger`, `--color-tc-success`, etc.).

```css
/* app.css */
@import "tailwindcss" source(none);
@import "./templ-components-theme.css";
@import "./custom.css";

@theme {
  /* Override semantic names — every component updates */
  --color-tc-primary: #4f46e5; /* indigo */
  --color-tc-primary-hover: #4338ca;
  --color-tc-danger: #dc2626;
  /* … */
}
```

Once you've imported the theme file, every `bg-blue-600` in the library
silently becomes your `--color-tc-primary`. One override re-skins the
whole library — buttons, links, focus rings, active states, toasts,
progress bars, all of it.

**Available semantic tokens** (see `templates/templ-components-theme.css`
for the full list):

| Token                      | Default     | Used by                                    |
| -------------------------- | ----------- | ------------------------------------------ |
| `--color-tc-primary`       | `blue-600`  | Buttons, links, active states, focus rings |
| `--color-tc-primary-hover` | `blue-700`  | Button hover                               |
| `--color-tc-danger`        | `red-600`   | Destructive buttons, errors, validation    |
| `--color-tc-success`       | `green-600` | Positive feedback, success toasts          |
| `--color-tc-warning`       | `amber-500` | Caution, "holding" trends                  |
| `--color-tc-info`          | `blue-500`  | Info toasts, informational feedback        |

Dark-mode equivalents: override `--color-tc-primary` etc. inside a
`.dark` scope or use Tailwind's `dark:` variants in your override CSS.

See [ADR-0008](adr/0008-semantic-tokens.md) for the design rationale.

---

## 2. Direct Tailwind palette override

If you don't want the semantic indirection, override Tailwind palette
colors directly in `@theme`:

```css
@theme {
  --color-blue-600: #4f46e5; /* now bg-blue-600 = indigo */
  --color-blue-500: #6366f1;
}
```

This works today without the theme file but couples your theme to
Tailwind's palette names. You'd need to override `--color-red-600`,
`--color-green-600`, etc. individually for each semantic intent.

---

## 3. Component-level `Class` override

Every component accepts `BaseProps.Class` which is merged via
`tailwind-merge-go`. Override one component instance:

```go
display.Button(display.ButtonProps{
    BaseProps: utils.BaseProps{Class: "bg-green-600 hover:bg-green-700"},
    Text:      "Confirm",
})
```

Use this for one-off styling. For library-wide theming, use option 1.

---

## Dark mode

Dark mode is class-based. `layout.ThemeScript()` and
`layout.ThemeToggle()` add/remove the `dark` class on `<html>`. To
override a dark-mode color, target `.dark`:

```css
.dark {
  --color-tc-primary: #818cf8; /* lighter indigo for dark backgrounds */
}
```

See [ADR-0011](adr/0011-dark-mode-convention.md) for the dark-mode
color convention (`-500` shade in dark mode, `-600` in light).

---

## Theme presets

Three starter presets ship in `templates/presets/`:

| Preset    | File                            | Style                                  |
| --------- | ------------------------------- | -------------------------------------- |
| `default` | `templates/presets/default.css` | The library defaults (blue + gray)     |
| `minimal` | `templates/presets/minimal.css` | Reduced palette, more whitespace       |
| `glass`   | `templates/presets/glass.css`   | Frosted-glass surfaces, blurred panels |

Import the one you want from your `app.css`:

```css
@import "./presets/glass.css";
```

Each preset file overrides only the `--color-tc-*` tokens — it does not
re-define the entire Tailwind palette.
