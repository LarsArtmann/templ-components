# ADR 0019: `recipes/` package — composition screens, not widgets

## Date

2026-07-21

## Status

Accepted

## Context

The library ships 98 individual components (display, forms, feedback, navigation, htmx, layout,
errorpage, icons). Consumers compose these into screens manually — every dashboard, settings
page, and login form rebuilds the same composition of AppShell + Grid + StatCard + Card slots
from scratch. This is the highest customer-value-per-line gap in the library: users don't want
widgets, they want screens.

TODO #31 (deferred since v0.12) asked for "blocks/composition examples." This ADR formalizes
the resolution: a new `recipes/` package that ships composition screens built from existing
components — no new primitives.

## Decision

Create a new top-level `recipes/` package that imports display/forms/layout/navigation
downward and exposes screen-level components.

### Package boundaries

- **Path:** `github.com/larsartmann/templ-components/recipes`
- **Depends on:** `display`, `forms`, `layout`, `navigation`, `utils`, `icons` (downward only)
- **Depended on by:** nothing in the library (top of the import DAG)
- **Import cycle risk:** zero — recipes never imports from a package that imports recipes.

### API model: `templ.Component` slots, not config bags

Each recipe takes a props struct with `templ.Component` slots for the variable parts
(charts, action panels, sidebar contents) and typed primitives for the fixed parts (titles,
breadcrumbs). This matches the existing Card.Body / PageHeader.Breadcrumb / AppShell.Sidebar
precedent — composition via slots, not inheritance.

```go
type DashboardProps struct {
    Title       string
    Subtitle    string
    Sidebar     templ.Component  // layout.AppShell sidebar slot
    HeaderBar   templ.Component  // optional right-side header items
    StatCards   []templ.Component // display.StatCard instances
    Charts      []templ.Component // user-supplied chart/visualization slots
    Footer      templ.Component
}
```

### Initial recipes (Phase 2)

| Recipe                   | Composes                                                                 |
| ------------------------ | ------------------------------------------------------------------------ |
| `recipes.Dashboard`      | `layout.AppShell` + `display.Grid` + `display.StatCard` + chart slots    |
| `recipes.SettingsLayout` | `layout.Split` + `navigation` (section nav) + `forms.Form` (per-section) |
| `recipes.LoginCard`      | `display.Card` + `forms.Form` (stack layout) + OAuth button slots        |

### Naming

- Package name: `recipes` (not `layouts`, not `screens`). Matches the shcnn/ui precedent
  ("recipes" = composition of primitives) and avoids collision with the existing `layout`
  package.
- Functions: `recipes.Dashboard`, `recipes.SettingsLayout`, `recipes.LoginCard` — single-word
  or camelCase screen names, no `Recipe` suffix (the package name is the suffix).

### What `recipes/` is NOT

- **Not a stylesheet bundle.** Recipes use the same Tailwind classes as the rest of the
  library; consumers still bring their own CSS build.
- **Not a fork-able starter.** Recipes live in the library and import the library. Consumers
  who want a starting point to fork should copy the `.templ` source files directly.
- **Not exhaustive.** Three recipes ship initially; more can be added when clear patterns
  emerge. We resist adding a recipe for every conceivable screen.

### Demo strategy

The demo site gets dedicated routes `/recipes/dashboard`, `/recipes/settings`, `/recipes/login`
showing the recipes rendered with mock content. Docs get `docs/recipes/*.md` with copy-paste
examples.

## Consequences

**Positive:**

- Resolves TODO #31 with a clean, downward-only composition layer.
- Consumers get high-value screens out of the box without abandoning the library's primitives.
- The recipes exercise the container-aware components (Phase 2.2-2.3) in real compositions,
  validating the container-query pattern end-to-end.
- New demo routes showcase the library's composition story.

**Negative:**

- One more package in the module (10 → 11 packages, counting examples/demo as separate).
- The `recipes.Dashboard` API surface is necessarily opinionated — consumers with very
  different dashboard layouts will still build their own. The slot model mitigates this but
  doesn't eliminate it.
- Documentation burden: each recipe needs its own `docs/recipes/*.md`.

## References

- [ADR-0018: Container-Query-Native Contract](0018-container-query-native-contract.md)
- [ADR-0016: Grid-first for 2D layouts](0016-grid-first-for-2d-layouts.md)
- TODO #31: "Blocks/composition examples" (resolved by this ADR)
