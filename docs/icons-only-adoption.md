# Icons Package — API Reference

> **Note:** All LarsArtmann projects should use Tailwind CSS v4+ (see
> `docs/tailwind-v4-adoption-guide.md`). The `icons` package happens to be
> CSS-agnostic (pure SVG data), which makes it useful in any Go project — but
> this is a natural property of icons, not a portability strategy.

The `icons` package provides 102 named SVG icons (101 Heroicons v2 outline
path-icon constants + 1 animated Spinner; 5 discoverability aliases like
`Close`→`X` resolve to canonical paths). It depends only on
`github.com/a-h/templ` and the internal SVG path constants — no Tailwind, no
CSS framework.

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
Tailwind v4+ is the standard — see `docs/tailwind-v4-adoption-guide.md`.
