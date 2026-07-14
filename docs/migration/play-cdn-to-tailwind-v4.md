# Recipe: Migrate from Tailwind Play CDN to Tailwind v4 CSS-first

**Audience:** Consumers who started with the Tailwind Play CDN
(`<script src="/static/tailwind.min.js">`) and want to move to the CSS-first
build for faster page loads, smaller payloads, and a strict Content-Security-Policy.

**Trigger:** dashboard feels slow (~300 KB JS download + client-side compile
on every page load), or your CSP needs `style-src 'unsafe-inline'` weakened.

**Outcome:** a single ~20–40 KB pre-compiled `app.css`, no runtime JS, and a
CSP that allows only `style-src 'self'`.

---

## Why migrate?

| Play CDN (before)                      | CSS-first build (after)              |
| -------------------------------------- | ------------------------------------ |
| ~300 KB JS downloaded on every page    | ~20–40 KB CSS, cacheable             |
| Compiles CSS client-side on each load  | Pre-compiled at build time           |
| Requires `style-src 'unsafe-inline'`   | Works with `style-src 'self'`        |
| Flash of unstyled content possible     | Styles available immediately         |
| No `@source` scanning of vendored deps | Scans templ-components automatically |

---

## Quick migration

Instead of the manual 7-step process below, use one of these shortcuts:

- **BuildFlow:** If using BuildFlow, the `tailwind-build` provider does everything automatically.
- **Starter template:** Copy [`templates/app.css`](../../templates/app.css) and [`templates/custom.css`](../../templates/custom.css) as your entry point (both files needed — `app.css` imports `custom.css`), run `go mod vendor`, then build with the Tailwind CLI.

---

## Step 1 — Install the Tailwind v4 CLI

Pick one option. The CLI is the same binary in all three cases.

### Option A: Standalone binary (no Node.js)

```bash
curl -sSLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
chmod +x tailwindcss-linux-x64
mv tailwindcss-linux-x64 ~/.local/bin/tailwindcss
```

### Option B: Nix flake (recommended for LarsArtmann projects)

Add `pkgs.tailwindcss_4` to your `devShells.default.packages`:

```nix
devShells.default = pkgs.mkShellNoCC {
  packages = with pkgs; [ go_1_26 templ tailwindcss_4 ];
};
```

### Option C: npm

```bash
npm install -D @tailwindcss/cli
npx @tailwindcss/cli -i app.css -o static/app.css --minify
```

## Step 2 — Create your CSS entry point

Create `app.css` at the project root:

```css
/* app.css */
@import "tailwindcss";

/* Scan templ-components for class names (vendored module path) */
@source "../vendor/github.com/larsartmann/templ-components";

/* Class-based dark mode toggle (required for templ-components dark mode) */
@custom-variant dark (&:where(.dark, .dark *));

/* Optional: override brand colors globally */
@theme {
  --color-blue-600: #4f46e5; /* indigo-600 */
  --color-blue-500: #6366f1;
}

/* Optional: small custom CSS where Tailwind doesn't cover something */
@layer utilities {
  .scrollbar-thin {
    scrollbar-width: thin;
  }
}
```

> **Where is templ-components vendored?** Run `go mod vendor && ls vendor/github.com/larsartmann/templ-components`. If you don't vendor, use `@source "path/to/local/checkout"` during development, or rely on the `@source` of the templ-components module in your GOPATH.

## Step 3 — Build the CSS

```bash
tailwindcss -i app.css -o static/app.css --minify
```

For development with auto-rebuild:

```bash
tailwindcss -i app.css -o static/app.css --watch
```

Add this to your `flake.nix` as a `build-css` app so the build is reproducible:

```nix
apps.build-css = {
  type = "app";
  program = pkgs.writeShellApplication {
    name = "build-css";
    runtimeInputs = [ pkgs.tailwindcss_4 ];
    text = ''
      tailwindcss -i app.css -o static/app.css --minify
    '';
  };
};
```

## Step 4 — Wire it into `layout.PageProps`

templ-components auto-injects `<link rel="stylesheet" href="/app.css">` by default
(that's the `DefaultPageProps()` CSSPath). You usually don't need to change
anything — just serve the generated `static/app.css` at the path `/app.css`.

If you serve it at a different path, set `CSSPath`:

```go
props := layout.DefaultPageProps()
props.CSSPath = "/static/app.css"
```

If you serve styles another way (inline, CDN), suppress the auto-inject:

```go
props.CSSPath = "" // no <link rel="stylesheet"> emitted
```

## Step 5 — Remove the Play CDN script

Delete this line from your base template / head content:

```html
<script src="/static/tailwind.min.js"></script>
```

If you were loading it via `HeadContent`, replace with nothing —
the compiled `app.css` replaces it.

## Step 6 — Tighten the CSP

Your `style-src` directive can now drop `'unsafe-inline'`:

```
# Before (Play CDN required unsafe-inline)
Content-Security-Policy: style-src 'self' 'unsafe-inline'

# After (CSS-first build)
Content-Security-Policy: style-src 'self'
```

> **Note:** if you use `layout.ThemeScript` (dark mode toggle), that script tag
> carries a CSP nonce via `nonce={ nonce }` — it does not need `'unsafe-inline'`.

## Step 7 — Verify

1. Load your page. No 404 on `/app.css`.
2. Inspect: the network tab should show one CSS file (~20–40 KB), no
   `tailwind.min.js` request.
3. Toggle dark mode — the `dark` class should still flip colors.
4. Check the browser console for CSP violations. There should be none.

---

## Common pitfalls

- **`@source` path wrong** → Tailwind compiles but utility classes from
  templ-components are missing. Verify the vendor path with `ls`.
- **Forgot `@custom-variant dark`** → dark mode toggle does nothing.
- **Cached `tailwind.min.js`** → browser serves the old Play CDN script from
  cache. Hard-reload or bust the cache.
- **`@theme` overrides ignored** → make sure they're inside the `@theme {}`
  block, not top-level CSS.

See `docs/tailwind-v4-adoption-guide.md` for the full v4 setup reference.
