# ADR 0020: Per-Package Go Modules Split

## Date

2026-07-21

## Status

**Proposed — deferred until consumer demand.** A prototype was developed on the
`modularize/strategic-split` branch (never merged). This ADR documents the
design; the execution awaits concrete consumer requests for independent
importability.

## Context

The library is a single Go module (`github.com/larsartmann/templ-components`)
with 15 packages. Consumers who only need icons (the `icons` package) currently
`go get` the entire module — including display, forms, feedback, layout,
navigation, errorpage, recipes, htmx. The Go module proxy dedupes shared
dependencies at the module-graph level, so this is not a build-performance
problem in practice. But it is a **dependency-graph** problem: a consumer who
wants only icons ends up with templ + tailwind-merge-go in their go.sum because
`icons` depends on `internal/svg` + `utils` which pulls tailwind-merge-go.

For most consumers this is a non-issue. For libraries that want to embed just
the icons (the "icons-only adoption" pattern documented at
`docs/icons-only-adoption.md`), per-package modules would let them `go get
github.com/larsartmann/templ-components/icons` without dragging in the rest.

## Decision (when executed)

Split the single module into per-package modules connected by a `go.work`
workspace for local development. Each leaf module is independently versioned
and taggable.

### Module boundaries

| Module path                                            | Depends on                                                  |
| ------------------------------------------------------ | ----------------------------------------------------------- |
| `github.com/larsartmann/templ-components/utils`        | (stdlib only + tailwind-merge-go + templ)                   |
| `github.com/larsartmann/templ-components/internal/svg` | (stdlib only)                                               |
| `github.com/larsartmann/templ-components/icons`        | utils, internal/svg                                         |
| `github.com/larsartmann/templ-components/htmx`         | utils                                                       |
| `github.com/larsartmann/templ-components/layout`       | utils, icons                                                |
| `github.com/larsartmann/templ-components/forms`        | utils, icons                                                |
| `github.com/larsartmann/templ-components/display`      | utils, icons, internal/svg, htmx                            |
| `github.com/larsartmann/templ-components/feedback`     | utils, icons, internal/svg                                  |
| `github.com/larsartmann/templ-components/navigation`   | utils, icons, internal/svg                                  |
| `github.com/larsartmann/templ-components/errorpage`    | utils, icons                                                |
| `github.com/larsartmann/templ-components/recipes`      | utils, icons, layout, navigation, forms, display            |
| `github.com/larsartmann/templ-components` (root)       | re-exports everything for backward compat (one minor cycle) |

### Compatibility strategy

The root module `github.com/larsartmann/templ-components` becomes a compat
shim: it re-exports every package via type aliases and re-emits every
component function as a thin wrapper. Consumers who don't update their import
paths keep working. The compat shim is removed in v2.1 (one minor cycle after
v2.0 ships the split).

### Versioning

Each sub-module gets its own tagged release. The root compat module stays at
`v2.0.x`. Sub-modules start at `v1.0.0` (they have no prior independent
versioning history).

### CI matrix

CI builds each module independently to catch "works in the workspace, breaks
standalone" regressions. The `go.work` file is for local dev only; CI
simulates the consumer experience by `go get`-ing the dependency modules from
their tagged versions.

### `internal/` packages

`internal/svg`, `internal/golden`, `internal/contract`, `internal/testutil`
(proposed in ADR-0014) cannot be imported externally by Go's `internal` rule.
They become modules in their own right but the `internal/` path prefix is
dropped — `svg`, `golden`, `contract`, `testutil` become regular packages
under the module paths `.../svg`, etc.

## Why deferred

- **No consumer demand yet.** The single-module model has worked for 22
  releases (v0.1 → v1.0). Icons-only adopters currently vendor the library
  and accept the go.sum bloat.
- **CI complexity cost.** Per-module CI is real ongoing work — every PR
  touches the workspace, but the matrix build catches per-module breakage
  separately. This is worth doing when there are 3+ active consumers; today
  there are 0 confirmed.
- **Compat shim is throwaway code.** It exists for exactly one minor cycle
  and then is deleted. Writing it before consumers ask for the split is
  pay-now-benefit-maybe-later.
- **`go.work` workspace UX.** Local development with `go.work` works well
  but is a new pattern for contributors. Adding it as part of a deferred
  split lets us learn the workflow before committing.

## Trigger criteria

Execute this ADR when **any one** of:

1. A consumer publicly requests `go get .../icons` (or any single-package
   import) without the rest.
2. The library's `go.sum` adds >50 entries from a single new dependency that
   only one package actually uses.
3. A downstream library wants to embed `templ-components` icons but cannot
   afford the `templ` + `tailwind-merge-go` dependency closure.

Until then, the single-module model stands.

## What v2.0 ships instead

v2.0 ships the CLI scaffolding tool (`tc add <component>`, Phase 5.3) and the
headless-variants evaluation (Phase 5.4). Both are non-breaking additions.
The modules split awaits its trigger.

## References

- `modularize/strategic-split` branch — prototype of this split (unmerged)
- [ADR-0019: recipes package](0019-recipes-package.md) — established the
  top-of-DAG composition pattern that the modules split preserves
- [icons-only adoption guide](../icons-only-adoption.md) — the consumer
  pattern most likely to trigger this split
