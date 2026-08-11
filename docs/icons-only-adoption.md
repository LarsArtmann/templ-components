# Icons Package — API Reference

> **Note:** All LarsArtmann projects should use Tailwind CSS v4+ (see
> `docs/tailwind-v4-adoption-guide.md`). The `icons` package happens to be
> CSS-agnostic (pure SVG data), which makes it useful in any Go project — but
> this is a natural property of icons, not a portability strategy.

The `icons` package provides 102 named SVG icons (101 Heroicons v2 outline
path-icon constants + 1 animated Spinner; 5 discoverability aliases like
`Close`→`X` resolve to canonical paths). It depends only on
`github.com/a-h/templ` and the `utils/svg` path constants — no Tailwind, no
CSS framework.

## Standalone module (icons-only adoption)

Since ADR-0034, `icons` is a **separate Go module**. You can adopt it without
pulling in any UI components, Tailwind, or CSS:

```bash
go get github.com/larsartmann/templ-components/icons@latest
```

This is the lightest possible dependency — only `github.com/a-h/templ` and
`utils/svg` (pure SVG path data). No `display`, `forms`, `feedback`, or any
Tailwind-emitting package is imported.

## Three API levels

### 1. `icons.Icon(name, class)` — full templ component

Use when you render inside a `.templ` file and want the library to build the
`<svg>` element. The `class` parameter is yours — pass any CSS class:

```templ
import "github.com/larsartmann/templ-components/icons"

@icons.Icon(icons.Users, "my-nav-icon text-blue-600")
```

### 2. `icons.IconPathData(name) []string` — raw path d-strings

Use when you need **full control** over the `<svg>` wrapper — custom class,
stroke-width, `width`/`height`, or `aria-hidden`. Returns raw SVG path
`d`-attribute strings with no markup wrapper:

```go
import "github.com/larsartmann/templ-components/icons"

func iconSVG(name icons.Name) string {
    var inner strings.Builder
    for _, d := range icons.IconPathData(name) {
        inner.WriteString(`<path d="`)
        inner.WriteString(d)
        inner.WriteString(`"/>`)
    }
    return `<svg class="my-icon" width="18" height="18" viewBox="0 0 24 24"
        fill="none" stroke="currentColor" stroke-width="1.8"
        stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">` +
        inner.String() + `</svg>`
}
```

### 3. `icons.IconPathJS(name) string` — pre-wrapped path elements

Use in JavaScript that dynamically creates icons. Returns `<path>` elements
with a fixed `stroke-width="1.5"`:

```go
icons.IconPathJS(icons.Home)
// => `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M2.25 12l8.954...`
```

## Animated icons (heroicons-animated inspired)

`AnimatedIcon` and `AnimatedIconWithAnimation` render any icon with a
hover-triggered CSS animation — pure CSS, zero JavaScript, `prefers-reduced-motion`
support. Inspired by [heroicons-animated.com](https://www.heroicons-animated.com/).

```templ
@icons.AnimatedIcon(icons.Heart, "h-6 w-6 text-red-500")

@icons.AnimatedIconWithAnimation(icons.Bell, icons.AnimWiggle, "h-6 w-6")
```

Each icon has a default animation via `DefaultAnimation()`:

| Icon | Default | Icon | Default |
| --- | --- | --- | --- |
| Heart | `AnimPulse` | Home | `AnimJump` |
| Star | `AnimBeat` | Search | `AnimBounce` |
| Bell | `AnimWiggle` | ChevronDown | `AnimNod` |
| Settings | `AnimSpin` | Eye | `AnimBlink` |
| Beaker | `AnimWobble` | Bolt | `AnimDraw` |
| ExternalLink | `AnimShake` | Refresh | `AnimSpin` |

11 presets total: `AnimPulse`, `AnimBeat`, `AnimBounce`, `AnimWiggle`, `AnimSpin`,
`AnimJump`, `AnimNod`, `AnimShake`, `AnimBlink` (per-path), `AnimWobble`, `AnimDraw` (self-draw via `stroke-dashoffset`).

**CSS requirement:** copy the `.tc-anim-*` classes from `templates/custom.css`
into your stylesheet. The animations use CSS `@keyframes`, individual transform
properties (`scale`, `rotate`, `translate`), and `:hover` / `:focus-within` triggers.

**DOM structure:** `AnimatedIcon` wraps the SVG in a `<span>` element (the hover
trigger), unlike `Icon()` which renders a bare `<svg>`. This extra wrapper may
affect flex/grid layouts or CSS sibling combinators that assume the SVG is a
direct child. Use `Icon()` / `IconRTL()` if you need the bare SVG.

## Full icon catalog

102 icons (101 path-icon constants + 1 animated Spinner). Typed constants prevent typos:

```go
icons.Users          // multi-person
icons.BuildingOffice2 // building/tenant
icons.Key            // credential/key
icons.ArrowRightOnRectangle // logout
icons.Squares2x2     // dashboard grid
icons.Search         // magnifying glass
icons.Mail           // envelope
icons.Trash          // trash can
icons.Plus           // plus sign
icons.Clock          // clock face
// ...see icons/icon_names.go for the full list
```

Unknown names fall back to `icons.Question` — the UI never breaks.

## Adding new icons

Icons use [Heroicons v2 outline](https://heroicons.com) path data. To add one:

1. Add a `Name` constant in `icons/icon_names.go`
2. Add the path data in `icons/icon_paths.go` (use `|` to separate multiple paths)
3. Run `go test ./icons/` — the auto-generated name list test verifies sync

## Package CSS dependencies

All packages emit Tailwind v4+ utility classes except `icons` (pure SVG data).
The `AnimatedIcon` / `AnimatedIconWithAnimation` functions require the `.tc-anim-*`
classes from `templates/custom.css` — copy them into your stylesheet if using
animated icons without the full library CSS.
Tailwind v4+ is the standard — see `docs/tailwind-v4-adoption-guide.md`.
