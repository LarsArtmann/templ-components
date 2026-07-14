# Tailwind CSS v4+ Adoption Guide

**The standard for all LarsArtmann projects.** Tailwind v4+ (latest), CSS-first
config, no Node.js runtime, no DaisyUI. Small custom `.css` only where Tailwind
doesn't cover something.

---

## Why Tailwind v4+?

| Benefit                 | Details                                                                                                                                          |
| ----------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| **No Node.js required** | Build step is a single Go-callable or standalone binary (`@tailwindcss/cli`). No `package.json`, no `node_modules`.                              |
| **CSS-first config**    | No `tailwind.config.js`. Everything lives in CSS via `@import "tailwindcss"` and `@theme`.                                                       |
| **`@source` scanning**  | Point Tailwind at any directory — including vendored Go modules — and it extracts class names from `.templ`, `.go`, `.html` files automatically. |
| **`@theme` overrides**  | Recolor the entire design system by overriding `--color-*` variables. One line changes every `bg-blue-600` globally.                             |
| **Dark mode**           | `dark:` variants follow `prefers-color-scheme` by default (zero-config). Optional `@custom-variant` for toggle strategy.                         |
| **Zero runtime JS**     | Pure CSS output. No Alpine.js, no DaisyUI JS, no hydration. Server-rendered HTML just works.                                                     |

---

## Setup (recommended)

### Option A: BuildFlow

If your project uses [BuildFlow](https://github.com/larsartmann/buildflow), the `tailwind-build` provider handles CSS automatically as part of its DAG:

```
go-mod-vendor → templ-generate → tailwind-build
```

Just ensure a CSS file with `@import "tailwindcss"` exists. No `go:generate` directive needed.

### Option B: Starter template

Copy [`templates/app.css`](../templates/app.css) and [`templates/custom.css`](../templates/custom.css) into your project as a ready-to-use entry point, then build:

```bash
go mod vendor
tailwindcss -i app.css -o styles.css --minify
```

`app.css` contains the Tailwind directives (`@import`, `@source`, `@theme`, `@custom-variant`). `custom.css` contains component-specific styles (dialog animations, stylable select, auto-grow textarea, etc.) and is imported by `app.css` via `@import "./custom.css"`. Copy both files to the same directory.

---

## Manual setup (5 minutes)

### 1. Install the Tailwind CLI

```bash
# Option A: Standalone binary (no Node.js)
curl -sSLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
chmod +x tailwindcss-linux-x64
mv tailwindcss-linux-x64 ~/.local/bin/tailwindcss

# Option B: npx (if you already have Node.js)
# Replace `tailwindcss` with `npx @tailwindcss/cli` in all commands below:
npx @tailwindcss/cli -i app.css -o static/app.css --minify

# Option C: Nix flake (recommended for LarsArtmann projects)
# Add pkgs.tailwindcss_4 to your devShell, then use `tailwindcss` directly.
```

### 2. Create your CSS entry point

```css
/* app.css */
@import "tailwindcss";

/* Scan templ-components (vendored) for class names */
@source "../vendor/github.com/larsartmann/templ-components";

/* Class-based dark mode toggle */
@custom-variant dark (&:where(.dark, .dark *));

/* Override the primary brand color */
@theme {
  --color-blue-600: #4f46e5; /* indigo-600 instead of blue-600 */
  --color-blue-500: #6366f1; /* indigo-500 */
}

/* Small custom CSS where Tailwind doesn't cover something */
@layer utilities {
  .scrollbar-thin {
    scrollbar-width: thin;
    scrollbar-color: var(--color-gray-400) transparent;
  }
}
```

### 3. Build

```bash
# One-shot build
tailwindcss -i app.css -o static/app.css --minify

# Watch mode (during development)
tailwindcss -i app.css -o static/app.css --watch
```

### 4. Include in your HTML

```html
<link rel="stylesheet" href="/static/app.css" />
```

---

## Vendoring templ-components

If you use `templ-components`, Tailwind needs to scan the library's source files
for class names. There are two approaches:

### Vendored (recommended for full control)

```bash
go mod vendor
```

Then add the `@source` line pointing at the vendored path (shown above).

### Go module cache (no vendoring)

If you don't vendor, Tailwind can scan the module cache directly:

```css
/* Use go list to find the exact path */
@source "$(go env GOMODCACHE)/github.com/larsartmann/templ-components@v0.13.0";
```

To generate this path automatically:

```bash
# Print the @source line for your current version
echo "@source \"$(go list -m -f '{{.Dir}}' github.com/larsartmann/templ-components)\";"
```

> **Troubleshooting:** If components render unstyled (no colors, no spacing), the
> `@source` path is wrong and Tailwind can't find the class names. Verify the path
> with `ls "$(go env GOMODCACHE)/github.com/larsartmann/templ-components@"*`.

---

## Setting Class on components

All component props embed `utils.BaseProps`, which provides the `Class`, `Attrs`,
`ID`, and `AriaLabel` fields. In Go, **promoted fields cannot be set in struct
literals** — you must use the embedded struct name:

```go
// This does NOT compile — Class is promoted from BaseProps:
// display.GridProps{Cols: display.GridCols3, Class: "gap-8"}

// Correct — set Class through the embedded BaseProps:
display.GridProps{
    BaseProps: utils.BaseProps{Class: "gap-8"},
    Cols:      display.GridCols3,
}
```

This is a Go language limitation, not a library design choice. Once you know the
pattern, it's consistent across every component.

---

## Theming

### Override the palette with `@theme`

All templ-components emit standard Tailwind classes (`bg-blue-600`, `text-gray-900`).
Override colors globally via `@theme`:

```css
@theme {
  /* Your brand color replaces blue everywhere */
  --color-blue-600: #0ea5e9; /* sky-500 */
  --color-blue-500: #38bdf8; /* sky-400 */

  /* Dark mode surface adjustments */
  --color-gray-900: #0f172a; /* slate-900 */
  --color-gray-800: #1e293b; /* slate-800 */
}
```

One `@theme` block changes every `bg-blue-600` across every component — no Go code
changes needed.

### Map to your CSS-variable design system

If your project uses CSS custom properties (e.g., `var(--surface)`, `var(--text)`),
map the library's palette to your tokens:

```css
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
}
```

Now every `bg-white`, `text-gray-*`, `border-gray-*`, `dark:bg-gray-*` in every
library component resolves to your variables. Both sides of each `dark:` pair
resolve to the same consumer variable, so mode-switching is automatic — no `.bg-white`
bridge hacks, no strategy mismatch. See [Dark Mode & Theming Research](dark-mode-research.md)
"Path 3" for the full worked example.

### Semantic aliases

For semantic aliases (`bg-tc-primary`, `text-tc-danger`), use the included
`templ-components-theme.css`:

```css
@import "tailwindcss";
@import "./templ-components-theme.css";
```

For a step-by-step recipe on mapping custom semantic palettes to library tokens,
see [Theme Bridge Recipe](recipes/theme-bridge.md).

All templ-components emit standard Tailwind classes (`bg-blue-600`, `text-gray-900`).
Override colors globally via `@theme`:

```css
@theme {
  /* Your brand color replaces blue everywhere */
  --color-blue-600: #0ea5e9; /* sky-500 */
  --color-blue-500: #38bdf8; /* sky-400 */

  /* Dark mode surface adjustments */
  --color-gray-900: #0f172a; /* slate-900 */
  --color-gray-800: #1e293b; /* slate-800 */
}
```

For semantic aliases (`bg-tc-primary`, `text-tc-danger`), use the included
`templ-components-theme.css`:

```css
@import "tailwindcss";
@import "./templ-components-theme.css";
```

---

## Migrating from a custom CSS design system

If your project has a hand-rolled CSS design system (like cqrs-htmx/adminui's
706-line `admin.css`), here's the migration path:

### Phase 1: Wrap tokens in `@theme` (1 hour)

Map your CSS custom properties to Tailwind's `--color-*` namespace:

```css
@theme {
  /* Map admin.css --accent to Tailwind's blue-600 */
  --color-blue-600: var(--accent);
  /* Map --bg, --surface, --border to gray scale */
  --color-gray-50: var(--surface-2);
  --color-gray-100: var(--border);
  --color-gray-900: var(--text);
}
```

### Phase 2: Replace utility classes progressively (ongoing)

Swap custom classes (`card`, `btn`, `badge`) for templ-components or raw Tailwind
utilities. Each component migrated deletes a chunk of custom CSS.

### Phase 3: Delete the custom CSS (when done)

Once all components use Tailwind, the hand-rolled design system is dead weight.
Delete it.

---

## Dark mode strategies

templ-components components emit `dark:` variants on all color classes. Tailwind
v4's **default** `dark:` variant follows `prefers-color-scheme` (zero-config).
There are three ways to handle dark mode:

### Path 1: OS-following (zero config)

The simplest option. Tailwind v4's default `dark:` variant follows the user's OS
setting. No configuration, no JavaScript:

```css
/* app.css — that's it */
@import "tailwindcss";
@source "../vendor/github.com/larsartmann/templ-components";

:root {
  color-scheme: light dark;
}
```

Every `dark:bg-gray-800` in every library component activates automatically.
Add `color-scheme: light dark` on `:root` so native form controls render correctly.

### Path 2: Toggle (user-controlled)

For apps that want a dark/light toggle button:

```css
@import "tailwindcss";
@source "../vendor/github.com/larsartmann/templ-components";

/* Override the dark variant to use a class instead of media query */
@custom-variant dark (&:where(.dark, .dark *));

:root {
  color-scheme: light dark;
}
.dark {
  color-scheme: dark;
}
```

Use `layout.ThemeScript()` (FOUC-prevention inline script) + `layout.ThemeToggle()`
(the toggle button). Both are only needed for this path.

### Path 3: CSS-variable design system

For projects with an existing CSS custom property design system. Map the library's
palette to your variables via `@theme` (see [Theming](#theming) section above).
Both sides of each `dark:` pair resolve to your variables, so mode-switching is
automatic regardless of strategy.

### Comparison

| Path                       | Config needed                                     | JavaScript           | Best for                        |
| -------------------------- | ------------------------------------------------- | -------------------- | ------------------------------- |
| OS-following               | `@import` only                                    | None                 | Blogs, docs, internal tools     |
| Toggle                     | `@custom-variant` + `ThemeScript` + `ThemeToggle` | ThemeScript (inline) | SaaS, dashboards, user-pickable |
| CSS-variable design system | `@theme` palette override                         | None (or your own)   | Existing design systems         |

> **`prefers-color-scheme` consumers:** Without a `.dark` class in the DOM,
> component-internal `dark:` variants that reference non-surface colors (e.g.,
> `dark:text-gray-300`) won't activate. Surface colors can be bridged via the
> `@theme` palette override (Path 3). Full `dark:` activation requires either
> the `.dark` class (Path 2) or OS-following (Path 1, which is the Tailwind
> default — just don't add `@custom-variant`).

For a deep analysis of all dark mode mechanisms in Tailwind v4 + modern CSS,
see [Dark Mode & Theming Research](dark-mode-research.md).

---

## Default constructors

Most components provide a `DefaultXxxProps()` constructor for the "I just want
the defaults with one field changed" pattern:

```go
props := display.DefaultGridProps()
props.Cols = display.GridCols4
props.BaseProps.Class = "gap-8"
```

Available constructors include `display.DefaultGridProps()`,
`display.DefaultCardProps()`, `display.DefaultEmptyStateProps()`,
`feedback.DefaultSpinnerProps()`, and more. Check each component's package for
the full list.

---

## FAQ

**Do I need Node.js?** No. The standalone Tailwind CLI binary works without Node.

**Can I use DaisyUI?** No. DaisyUI adds a component layer that conflicts with
templ-components. Use templ-components instead.

**Can I use arbitrary values?** Yes. `w-[248px]`, `grid-cols-[auto_1fr]`,
`text-[13px]` — all standard Tailwind v4.

**How do I handle dark mode?** See the [Dark mode strategies](#dark-mode-strategies)
section above. The simplest option is zero-config (`prefers-color-scheme`), or use
`@custom-variant dark` + `layout.ThemeScript()` + `layout.ThemeToggle()` for a
user-controlled toggle.

**What about CSP?** Tailwind produces pure CSS files — no inline styles, no
`eval()`. Include the built `.css` via `<link>` and you're CSP-compliant.

**Where does custom CSS go?** In your main CSS file, after `@import "tailwindcss"`.
Use `@layer utilities { ... }` or `@layer components { ... }` for organization.
