# Tailwind CSS v4+ Adoption Guide

**The standard for all LarsArtmann projects.** Tailwind v4+ (latest), CSS-first
config, no Node.js runtime, no DaisyUI. Small custom `.css` only where Tailwind
doesn't cover something.

---

## Why Tailwind v4+?

| Benefit                   | Details                                                                                                                                          |
| ------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| **No Node.js required**   | Build step is a single Go-callable or standalone binary (`@tailwindcss/cli`). No `package.json`, no `node_modules`.                              |
| **CSS-first config**      | No `tailwind.config.js`. Everything lives in CSS via `@import "tailwindcss"` and `@theme`.                                                       |
| **`@source` scanning**    | Point Tailwind at any directory — including vendored Go modules — and it extracts class names from `.templ`, `.go`, `.html` files automatically. |
| **`@theme` overrides**    | Recolor the entire design system by overriding `--color-*` variables. One line changes every `bg-blue-600` globally.                             |
| **Class-based dark mode** | `@custom-variant dark (&:where(.dark, .dark *))` enables JS-toggled dark mode without `prefers-color-scheme`.                                    |
| **Zero runtime JS**       | Pure CSS output. No Alpine.js, no DaisyUI JS, no hydration. Server-rendered HTML just works.                                                     |

---

## Setup (5 minutes)

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

If you use `templ-components`, vendor it so Tailwind can scan for class names:

```bash
go mod vendor
```

Then add the `@source` line pointing at the vendored path (shown above).

---

## Theming without touching component code

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

## FAQ

**Do I need Node.js?** No. The standalone Tailwind CLI binary works without Node.

**Can I use DaisyUI?** No. DaisyUI adds a component layer that conflicts with
templ-components. Use templ-components instead.

**Can I use arbitrary values?** Yes. `w-[248px]`, `grid-cols-[auto_1fr]`,
`text-[13px]` — all standard Tailwind v4.

**How do I handle dark mode?** Include `@custom-variant dark` in your CSS and
use `layout.ThemeScript()` + `layout.ThemeToggle()` for JS-toggled dark mode.

**What about CSP?** Tailwind produces pure CSS files — no inline styles, no
`eval()`. Include the built `.css` via `<link>` and you're CSP-compliant.

**Where does custom CSS go?** In your main CSS file, after `@import "tailwindcss"`.
Use `@layer utilities { ... }` or `@layer components { ... }` for organization.
