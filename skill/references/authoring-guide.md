# Authoring depth — templ-components

The long-form companion to `skill/SKILL.md` Part 2. Load this when you are writing or
reviewing component code and need the shared-helper inventory, the full per-component test
checklist, or a decision tree. Step-by-step process, build commands, invariants, and
anti-patterns stay in SKILL.md because they are consulted on every change.

## Shared internal helpers — reuse these

Re-inlining these is exactly what the art-dupl gate flags.

| Helper                                                                                             | Where                                           | Reuse for                                                                                            |
| -------------------------------------------------------------------------------------------------- | ----------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| `headingTag`                                                                                       | `display/shared.templ`                          | Title as h1–h6 (typed `HeadingTagType`, unknown → h3). Used by Card, EmptyState, CollapsibleSection. |
| `chartSeriesGroup`                                                                                 | `display/chart_shared.templ`                    | One series' `<g>` for a new SVG chart — compose `chart_geometry.go`, never reimplement the math.     |
| `sparklinePointCoords` / `sparklineGeometry`                                                       | `display/sparkline.go`                          | Value→pixel projection and min/max bounds (BarChart and Sparkline share it).                         |
| `calendarMonthNavLink` + `…Action`/`…URL`/`…Href`                                                  | `forms/calendar.templ`, `forms/calendar_nav.go` | Month navigation: href fallback + Wire cloning with `{year}`/`{month}` substitution.                 |
| `browser` / `distserver`                                                                           | `visualtest/internal/`                          | Capture tools (`ogshot`, `shots`, `siteshots`) — pinned Chromium + static-dist server.               |
| `feedbackStyleMap` / `feedbackStyle` / `feedbackIcon` / `feedbackRole`                             | `feedback/styles.go`                            | Any new alert/toast-like surface (shared `FeedbackType` styling).                                    |
| `DismissButton`                                                                                    | `utils/dismiss.templ`                           | Shared dismiss control (Alert, ErrorAlert); `data-dismiss="alert"`.                                  |
| `cardShellClass` / `mutedTextClass`                                                                | `display/card.templ` / `display/shared.go`      | Card shells and secondary text.                                                                      |
| Motion constants (`transitionFast`, `transitionNormal`, `transitionColors`, `transitionTransform`) | `display/shared.go`                             | Any transition — never inline timing strings.                                                        |

## Per-component testing checklist

“Assertion tests only” is not acceptable — that is how `layout.Script` became the
worst-tested component in the library.

```
[ ] golden_test.go      — exact rendered HTML vs .golden snapshot
[ ] a11y_test.go        — ARIA, roles, keyboard, motion-reduce, screen-reader text
[ ] rtl_check           — no physical properties (covered by TestRTLLogicalProperties)
[ ] bdd_test.go         — user-visible behaviour, not markup
[ ] edge_cases_test.go  — empty inputs, unknown enum values, ID collisions
[ ] example_test.go     — godoc ExampleXxx() compiles and renders
[ ] snapshot_test.go    — broader composition snapshot
[ ] coverage_*_test.go  — private helpers and branches
```

Golden tests use `utils/golden.Assert(t, name, got)` with class sorting + `EnsureID`
normalisation; `-update` regenerates after an intentional change (review the diff).

## Visual regression (`visualtest/`, separate module)

Chromedp pixel checks in their own module so browsers never enter the consumer dependency
graph. Supports `State` (`Hover`/`Focus`/`Click`/`Context`), `FullViewport` (top-layer
overlays), and `WaitSelector`. Regenerate goldens ONLY via `nix run .#visual -- -update`
(a bare `go test -update` there silently skips without the nix browser env). Load-flake
signature: light AND dark site goldens ~100% while demo routes pass → rerun, never `-update`.

## Decision trees

### Map lookup vs if-branch for a variant

- **Pure class/style data** → map + `utils.Lookup(map, key, fallback)`
  (`badgeStyleMap`, `cardPaddingLookup`, `spinnerSizeLookup`, `gridColsLookup`, …).
- **Structural DOM differences** (TabsVariant, DropdownPosition, TrendDirection) → `if`
  branch inside the template. Maps are for data, not for choosing which markup to emit.
- **Enum validation** → map + fallback (graceful); `InputType`, `ButtonHTMLType`, `FormMethod`
  fall back to HTML-spec defaults; only the icons path-data check may panic.

### `utils.Class()` vs `utils.Lookup()`

- Merging class strings → `utils.Class("a", b)` (tailwind-merge resolves conflicts).
- Selecting a style variant → `utils.Lookup(map, key, fallback)`.
- Both → look up first, then merge.

### When to introduce a typed enum

Any field with a fixed set of visual/behavioural variants becomes `type X string` + consts +
a `Default` constant + `IsValid()` + a test. Never accept a raw string for a closed set —
that reopens the invalid state you are trying to prevent.

### When to add an icon

Add the path to `iconPathData` in `icons/icon_paths.go` (multi-path icons use `|`;
`iconPaths()` panics on a stray separator). `allIconNames()` derives automatically — never
hand-maintain a second list. Update the icon count in the docs in the same commit.

### When interactivity needs JavaScript

1. Native HTML first (`<details>`, `<form>`, `:target`/`:checked`). If yes, stop.
2. One inline script through `utils.ScriptComponent`, guarded by a global singleton flag
   (`window.tcXxxAttached`), event delegation on `document`, Escape/focus restore/click-outside.
3. Escape interpolated IDs with `strconv.Quote()` (see `validateDropdownID`).

Full ladder: `docs/javascript-guide.md`; attachment rationale: ADR-0005.

### When to generate an ID

`utils.EnsureID("prefix", props.ID)` — returns the consumer's id unchanged, else
`tc-<prefix>-<16 hex>` from `crypto/rand`. Never `time.Now()` alone.

### Theming decisions

- Emit standard Tailwind utilities (`bg-blue-600`, `text-gray-900`).
- Consumers override via `@theme { --color-blue-600: #...; }` — do NOT add `Variant`/`Color`
  props for theming; that defeats the CSS-variable model.
- Semantic aliases live in `templ-components-theme.css` (`bg-tc-primary`, `text-tc-danger`).

## Website SSG page registration

The site (`website/`, a Go SSG rendering through this library) has three registration points
that must stay in sync when adding pages:

- **`topLevelPages`** (`website/cmd/site/main.go`) — non-docs HTML pages; a new entry needs
  the nav/header templates updated too.
- **Sitemap lastmod sources** — git history over declared repo-relative paths; a page without
  a source emits an undated sitemap entry (`writeSitemaps`).
- **Search index** — `build.PlainText` indexes markdown docs only; top-level pages are found
  via nav, not search (deliberate).

Visual changes are pinned by `TestSiteRouteGoldens` (`docs/visual-testing.md`).
