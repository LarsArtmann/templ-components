# Consumer guide — templ-components

Everything a consumer needs beyond the component index in `skill/SKILL.md`: quick start,
transports, head-tag overrides, CSP, theming, compiled-CSS rules, the dark-mode checklist,
recipes, and site knobs. Versions and counts follow the repo; when a claim here and the
code disagree, the code wins — then fix this file.

## Quick start: adopting the library

```bash
go get github.com/larsartmann/templ-components@latest
```

### Minimal page

```templ
@layout.Base(layout.DefaultPageProps()) {
    @navigation.SimpleNav(navigation.SimpleNavProps{
        BrandText:   "MyApp",
        Links:       []navigation.NavLinkProps{{Href: "/", Text: "Home"}},
        RightItems:  layout.ThemeToggle("Toggle theme", ""),
        CurrentPath: "/",
    })
    <main class="mx-auto max-w-7xl px-4 py-8">
        @display.Grid(display.GridProps{Cols: display.GridCols3}) {
            for _, u := range users {
                @display.StatCard(display.StatCardProps{Label: u.Name, Value: u.Count, Href: "/users/" + u.ID})
            }
        }
    </main>
    @feedback.ToastContainer("")
    @htmx.GlobalErrorHandling(htmx.DefaultErrorHandlingConfig())
}
```

`RightItems` is a `templ.Component` field: pass the component function value
(`layout.ThemeToggle(...)`), never `@layout.ThemeToggle(...)` — `@` only works in template body.

### Dual-transport form (HTMX + Datastar)

One `wire.Action`, either runtime — submit, server-side validation, inline field errors
(`docs/recipes/server-side-validation.md`):

```templ
@forms.Form(forms.FormProps{
	Method:     forms.FormPost,
	CSRFToken:  csrf,
	NoValidate: true, // let server-side validation own the error UX
	Wire:       &wire.Action{Method: wire.MethodPost, URL: "/api/subscribe", Target: "#form-region"},
}) {
	@forms.Input(forms.InputProps{Name: "email", Type: forms.InputEmail, Label: "Email"})
	@display.Button(display.ButtonProps{Text: "Subscribe", Type: display.ButtonHTMLSubmit})
}
<div id="form-region" aria-live="polite"></div>
```

Render the form INSIDE its swap region and re-render it on validation errors (values
preserved); wrap the handler in `wire.Handler` and answer 200 OK even for errors. Under
Datastar, patch target and mode come ONLY from response headers — `Action.Target`/`Swap`
render for htmx only (see `docs/transport-wiring.md`).

### Suppressing auto-injected `<head>` tags

`DefaultPageProps()` auto-injects HTMX and the stylesheet. Override them:

```go
props := layout.DefaultPageProps()
props.HTMXVersion = "" // suppress the HTMX <script>
props.CSSPath = ""     // suppress <link rel="stylesheet" href="/app.css">
```

Self-hosting and pure-Datastar pages: `props.HTMXCDN = "https://unpkg.com"`, or
`props.HTMXSrc = layout.HTMXNone` for no htmx runtime at all.

### CSP-safe script tags

Use `layout.Script` for external scripts; the nonce can never be forgotten:

```templ
@layout.Script(nonce, "/static/app.js", templ.Attributes{"defer": true})
```

Inline component scripts MUST render through `utils.ScriptComponent(nonce, script, errLabel)`
— it owns the omit-empty nonce rule, nonce escaping, and error wrapping. The rule: an empty
nonce renders NO nonce attribute, never `nonce=""` (strict-CSP browsers reject the empty
value, silently killing the feature). A component's singleton script ALWAYS renders — never
skip the whole `<script>` when the nonce is empty, or non-CSP consumers get no JS.
templ's `<script>` context sanitizes interpolation, so JS with braces/template literals must
go through the ComponentFunc writer, not a templ `<script>` block.

### Theming

Components emit standard Tailwind classes (`bg-blue-600`). Override globally, no Go changes:

```css
@theme {
  --color-blue-600: #4f46e5; /* indigo instead of blue */
}
```

Dark mode: `@custom-variant dark (&:where(.dark, .dark *))` + `layout.ThemeScript()` +
`layout.ThemeToggle()`. `color-scheme: light` on `:root`, `color-scheme: dark` on `.dark`.

**Custom CSS utilities** (`templates/custom.css`): `.tc-*` classes for what Tailwind can't
express — overlays (`.tc-overlay`, `.tc-modal`, `.tc-drawer`), scroll-snap (`.tc-snap-*`),
stylable select (`.tc-select`), auto-growing textarea (`.tc-auto-grow`), content-visibility
(`.tc-content-auto`), squircle corners (`.tc-squircle`), and fluid typography
(`.tc-fluid-display`, `.tc-fluid-h1`–`h4`, `.tc-fluid-lead`, container-`cqi` based; see
`docs/recipes/fluid-typography.md`). Consumers vendor `templates/custom.css` via `app.css`.
Guarded by `TestCustomCSSUtilities`.

### Compiled-CSS artifact rule (fleet adoption lesson)

A consumer with a TRACKED compiled CSS artifact must regenerate it in the SAME commit as any
component adoption — an adopted component whose utilities never reached the stylesheet renders
unstyled outside local dev (string tests cannot see this).

1. Regenerate via the consumer's canonical pipeline (nix app or `regen-ui.sh`-style script),
   not ad-hoc tool invocations.
2. If classes enter the scan through a committed inventory, the scan must be REPLACE-AWARE:
   resolve the library with `go list -m -f '{{.Dir}}'` (honors go.mod replace) and cover
   EVERY imported package — a missing package in the scan list shipped a 404-page outage.
3. Grep the compiled artifact for the adopted component's classes before committing.

### Dark mode checklist (new components)

Every new component passes `TestDarkModeCompliance` + `TestDarkModeSemanticColors`:

1. **Neutral** — every `text-gray-*`, `bg-white`, `bg-gray-*`, `border-gray-*`, `ring-gray-*`
   carries a `dark:` variant on the same line.
2. **Semantic** — every `bg-blue-600`, `text-red-600`, … carries `dark:` (bg `-600`→`-500`,
   text `-600`→`-400`).
3. **Focus rings** — `focus:ring-blue-500` needs `dark:focus:ring-blue-400`.
4. **Ring offset** — `focus:ring-offset-2` needs `dark:focus:ring-offset-gray-900`.
5. **Shadows** — overlay/card shadows get `dark:shadow-black/20`.
6. **Palette** — `gray-*` exclusively; no `slate-*`/`zinc-*`/`neutral-*`/`stone-*`.

Full convention: `docs/adr/0011-dark-mode-convention.md`.

### Consumer `_templ.go` guidance

**Consumers gitignore `*_templ.go`** and generate at build time. **The library commits them**
because the module proxy serves tagged source as-is and never runs `templ generate` — missing
generated files break every consumer.

## Recipes

| Recipe | When to read |
| --- | --- |
| `docs/migration/play-cdn-to-tailwind-v4.md` | Migrating from Tailwind Play CDN to a CSS-first build |
| `docs/tailwind-v4-adoption-guide.md` | Full Tailwind v4 setup with `@source` scanning |
| `docs/recipes/server-rendered-htmx-error-feedback.md` | Wiring HTMX error feedback (toast/alert/page) |
| `docs/recipes/error-pages.md` | Family-aware error pages: status→family, per-code copy, HTMX branch, noindex |
| `docs/recipes/server-side-validation.md` | Dual-transport form validation round-trip |
| `docs/recipes/horizontal-filter-bar.md` | Horizontal HTMX filter bar vs `forms.Form` |
| `docs/recipes/sse-live-rows.md` | Append live `<tr>` over SSE — `Table.BodyID` + `hx-swap-oob` |
| `docs/recipes/bounding-tables.md` | Cap + `ListNote`, or Pagination/LoadMore — never an unbounded table |
| `docs/recipes/hybrid-strings-builder-rendering.md` | Rendering into a `strings.Builder` (`templ.WithChildren` escape hatch) |
| `docs/recipes/theme-bridge.md` / `scoped-theme-bridge.md` | Remap library colors, globally or scoped |
| `docs/icons-only-adoption.md` | Adopting just the `icons` package (CSS-agnostic) |
| `docs/transport-wiring.md` + `docs/datastar-runtime-facts.md` | The wire contract and pinned-runtime facts |

## How to know if a component already exists

1. **"By use case" table above** — start from your page type.
2. **`references/component-catalogue.md`** — every signature and one-liner.
3. **`README.md` "Component Catalog"** — grouped code examples.
4. **`pkg.go.dev/github.com/larsartmann/templ-components`** — full API reference.
5. **`rg -n "^templ [A-Z]" --glob "*.templ"`** — source of truth.

### Consumer tip: track adoption in your AGENTS.md

The #1 discoverability gap across consumers is not knowing which components are already
adopted vs hand-rolled. Keep a grep-able table in the consumer project's `AGENTS.md`:

```markdown
| Library component      | Status  | Where                   |
| ---------------------- | ------- | ----------------------- |
| `display.Grid`         | adopted | dashboard.templ         |
| `display.Table`        | custom  | filters.templ (thinner) |
| `forms.Form`           | custom  | filterForm helper       |
| `feedback.ProgressBar` | adopted | backfill.templ          |
```
