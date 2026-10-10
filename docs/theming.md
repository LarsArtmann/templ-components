# Theming templ-components

The library uses standard Tailwind CSS v4 color utilities (`bg-blue-600`,
`text-red-600`, etc.) so you can re-skin everything from a single CSS file
without touching any Go code.

There are **three** ways to override colors, in order of recommended use.

---

## 1. Semantic token layer (recommended)

Import `templ-components-theme.css` from your `app.css`. This file aliases
every Tailwind palette color used by the library to a semantic name — since
2026-10-10 over the **full 50–950 ramp** per family
(`--color-tc-primary-600`, `--color-tc-danger-700`, …), with `yellow-*`
bound to the warning family too (badge tints use yellow, feedback amber).

```css
/* app.css */
@import "tailwindcss" source(none);
@import "./templ-components-theme.css";
@import "./custom.css";

@theme {
  /* Override semantic shades — every component updates */
  --color-tc-primary-600: #4f46e5; /* indigo buttons/links */
  --color-tc-primary-700: #4338ca; /* hover */
  --color-tc-primary-100: #e0e7ff; /* badge tints */
  --color-tc-primary-800: #3730a3; /* badge tint text */
  /* … */
}
```

Once you've imported the theme file, every `bg-blue-600` in the library
silently becomes your `--color-tc-primary-600`. One override re-skins the
whole library — buttons, links, focus rings, active states, badge tints,
toasts, progress bars, all of it. The defaults are Tailwind's palette values
as literals, so an unmodified import renders identically to stock Tailwind.

**Available semantic tokens** (see `templates/templ-components-theme.css`
for the full list):

| Token family              | Default palette | Used by                                    |
| ------------------------- | --------------- | ------------------------------------------ |
| `--color-tc-primary-*`    | `blue-*` (50–950) | Buttons, links, active states, focus rings, badge primary tints |
| `--color-tc-danger-*`     | `red-*` (50–950)  | Destructive buttons, errors, validation    |
| `--color-tc-success-*`    | `green-*` (50–950) | Positive feedback, success toasts        |
| `--color-tc-warning-*`    | `amber-*` + `yellow-*` (50–950) | Caution, warning badges/tints |

Legacy named tokens (`--color-tc-primary`, `--color-tc-danger`, …) remain
and resolve through the shade tokens (e.g. `tc-primary` ≡ `tc-primary-600`).

Dark-mode equivalents: override the shade tokens inside a `.dark` scope or
use Tailwind's `dark:` variants in your override CSS.

Need the remap scoped to LIBRARY components only (your own blue stays blue)?
See [recipes/scoped-theme-bridge.md](recipes/scoped-theme-bridge.md) — and
ADR-0044's `--color-accent-*` family, which makes scoping structural in v2.

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

Four starter presets ship in `templates/presets/`:

| Preset    | File                            | Style                                  |
| --------- | ------------------------------- | -------------------------------------- |
| `default` | `templates/presets/default.css` | The library defaults (blue + gray)     |
| `minimal` | `templates/presets/minimal.css` | Reduced palette, more whitespace       |
| `glass`   | `templates/presets/glass.css`   | Frosted-glass surfaces, blurred panels |
| `emerald` | `templates/presets/emerald.css` | Emerald-green brand palette            |

Import the one you want from your `app.css`:

```css
@import "./presets/glass.css";
```

Each preset file overrides only the `--color-tc-*` tokens — it does not
re-define the entire Tailwind palette.

---

## Print-safe vs screen-only components

Most components render fine on paper, but a few are interactive-only.
The table below documents which components include `print:` variants for
graceful degradation when the page is printed.

| Component          | Print behaviour                               | Notes                                                   |
| ------------------ | --------------------------------------------- | ------------------------------------------------------- |
| `Card`             | Borders, shadows, and backgrounds removed     | `print:shadow-none print:border-0 print:bg-transparent` |
| `Modal`            | Hidden                                        | `print:hidden`                                          |
| `Drawer`           | Hidden                                        | `print:hidden`                                          |
| `PageHeader`       | Title heading avoids page break               | `break-after-avoid` on `<h1>`                           |
| `SectionHeading`   | Heading avoids page break                     | `break-after-avoid`                                     |
| `StatCard`         | No special print rules                        | Renders as-is                                           |
| `ProgressBar`      | No special print rules                        | Animations are static when printed                      |
| `CircularProgress` | No special print rules                        | SVG renders correctly in print                          |
| `Spinner`          | No special print rules                        | Animation is static when printed                        |
| `Alert`            | No special print rules                        | Renders as-is                                           |
| `Tabs`             | All panels visible (no JS to toggle in print) | Consider hiding tabs in print CSS if needed             |
| `Tooltip`          | No special print rules                        | Invisible unless hovered (which cannot happen in print) |
| `Dropdown`         | No special print rules                        | Closed state is the default                             |
| All others         | No special print rules                        | Render as-is                                            |
