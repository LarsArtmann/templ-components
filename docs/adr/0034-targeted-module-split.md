# ADR 0034: Targeted 7-Module Workspace Split

## Date

2026-08-10

## Status

**Accepted — executed.** Supersedes [ADR-0020](0020-per-package-modules-split.md)
(proposed a per-package ~12-module split, deferred until consumer demand).

## Context

ADR-0020 proposed splitting the library into per-package modules (one go.mod
per package: utils, svg, icons, htmx, layout, forms, display, feedback,
navigation, errorpage, recipes, root). The split was deferred until concrete
consumer demand. As of 2026-08-10, that demand remains unconfirmed (0 known
external consumers).

However, the maintainer chose to execute a **smaller, targeted** split for
three reasons:

1. **`go-error-family` isolation.** Only `errorpage` depends on
   `go-error-family`. Splitting `errorpage` into its own module means consumers
   who don't use error pages never see `go-error-family` in their dependency
   graph.
2. **Icons-only adoption.** The `icons` package is the most likely candidate
   for standalone consumption (documented in `docs/icons-only-adoption.md`).
   Splitting it into its own module enables `go get .../icons` without the full
   library.
3. **ECharts opt-in.** `charts/echarts` is already documented as an opt-in
   adapter (ADR-0031). Promoting it to its own module makes the opt-in boundary
   explicit at the dependency level.

A full per-package split (~12 modules) was rejected because:

- `utils/svg` is a universal dependency (every UI package uses it); splitting
  it adds overhead without composability payoff
- `utils/golden` is used by 8 packages' tests; per-package splitting would
  create a test-helpers dependency web
- ADR-0020's consumer-demand triggers were still unmet
- More modules = more go.mod files to maintain, tag, and keep in sync

## Decision

Split into **7 modules** connected by a `go.work` workspace for local dev and
`replace` directives for CI/consumer builds.

### Module boundaries

```
Layer 0 (leaf):  utils (utils, utils/svg, utils/cdn, utils/golden)
Layer 1:         icons, charts/echarts, htmx, datastar    [depend on utils]
Layer 2:         errorpage                                [depends on utils, icons]
Layer 3:         root (display, feedback, forms, layout, navigation,
                          recipes, integration, cmd/tc,
                          internal/contract, examples/demo)
```

| Module path                                              | Deps                                 | Purpose                                              |
| -------------------------------------------------------- | ------------------------------------ | ---------------------------------------------------- |
| `github.com/larsartmann/templ-components/utils`          | templ, tailwind-merge-go             | Leaf: BaseProps, Class(), EnsureID, svg, cdn, golden |
| `github.com/larsartmann/templ-components/icons`          | templ, utils                         | 106 named SVG icons; icons-only adoption             |
| `github.com/larsartmann/templ-components/errorpage`      | templ, go-error-family, icons, utils | Error pages + handler; isolates go-error-family      |
| `github.com/larsartmann/templ-components/charts/echarts` | templ, utils                         | Opt-in ECharts adapter (ADR-0031)                    |
| `github.com/larsartmann/templ-components/htmx`           | templ, utils                         | HTMX loading, error handling, OOB swaps              |
| `github.com/larsartmann/templ-components/datastar`       | templ, utils, go-datastar/static     | Datastar runtime + SSE LiveRegion (ADR-0030)         |
| `github.com/larsartmann/templ-components` (root)         | all above + testify                  | Core UI + recipes + integration + demo + CLI         |

### `internal/` package promotion

Go's `internal/` rule blocks cross-module access. Three `internal/` packages
were promoted to `utils/` sub-packages:

| Old path          | New path       | Reason                       |
| ----------------- | -------------- | ---------------------------- |
| `internal/svg`    | `utils/svg`    | Used by icons, display, etc. |
| `internal/cdn`    | `utils/cdn`    | Used by layout, htmx         |
| `internal/golden` | `utils/golden` | Used by 8 packages' tests    |

`internal/contract` stays in the root module (cross-cutting test package that
imports 10+ packages; works fine as a root-module internal).

### Dual strategy: `go.work` + `replace`

- **`go.work`** (gitignored): local development. Provides seamless cross-module
  refactoring. `use` directives for all 7 modules + visualtest.
- **`replace` directives** in each module's `go.mod`: CI and consumer builds.
  Each sub-module replaces its sibling deps with relative paths (`./utils`,
  `../icons`, etc.). The root module replaces all 6 sub-modules.

### Versioning: shared tag

All modules share a single version number (currently v1.8.1). At release time,
the root module gets tag `v<version>` and each sub-module gets tag
`<dir>/v<version>` (e.g., `utils/v2.0.0`, `icons/v2.0.0`). The Go module proxy
resolves sub-module versions via these directory-prefixed tags.

The `require` entries in each go.mod reference the current shared version
(v1.8.1). At release time, `scripts/release.sh` bumps all require entries to
the new version and creates all tags.

### Compatibility

Import paths are **unchanged**. A consumer who was importing
`github.com/larsartmann/templ-components/icons` continues to do so. The only
difference is that the import is now backed by a separate module. Consumers who
`go get` the root module get all sub-modules transitively via the require
directives.

The `internal/` packages that moved (`svg`, `cdn`, `golden`) were not importable
by external consumers anyway (Go's `internal/` rule), so no consumer is
affected by the rename.

## Consequences

### Positive

- **Dependency graph isolation.** Consumers who `go get .../icons` no longer
  pull `go-error-family`, testify, or the full UI package set.
- **Icons-only adoption** is now a first-class `go get` away.
- **Compile-time DAG enforcement.** The module boundaries enforce the
  dependency direction at the Go toolchain level, not just by convention.
- **Clean separation of concerns.** Each module has a clear, documented
  purpose and dependency set.

### Negative

- **Release complexity.** `scripts/release.sh` must tag 7 directories instead
  of 1. The version-sync guard must check all modules.
- **CI complexity.** golangci-lint does not support go.work; each module is
  linted independently. Per-module isolation tests verify standalone builds.
- **`go.work` is gitignored.** Contributors need to run `go work use` or rely
  on the `go.work` template. This is standard for Go monorepo workflows.
- **No compat shim for proxy consumers (yet).** The `v1.8.1` require entries
  with `replace` directives work for dev and repo clones. Proxy consumers need
  published sub-module tags. This is resolved at the next release (v2.0) where
  `scripts/release.sh` tags all sub-module directories.

## Verification

All 7 modules build, test, and lint clean:

- Workspace mode (`go.work`): `go build ./... && go test ./...`
- Standalone mode (`GOWORK=off`): per-module `go build ./... && go test ./...`
- Lint: per-module `golangci-lint run` (golangci-lint does not support go.work)

Strict acyclic DAG verified: no upward dependencies, no cycles.

## References

- [ADR-0020: Per-Package Go Modules Split](0020-per-package-modules-split.md) —
  the original proposal (deferred, superseded by this ADR)
- [ADR-0019: Recipes Package](0019-recipes-package.md) — established the
  top-of-DAG composition pattern
- [ADR-0031: Two-Tier Chart Architecture](0031-two-tier-chart-architecture.md) —
  charts/echarts as opt-in adapter
- [icons-only adoption guide](../icons-only-adoption.md) — consumer pattern
- [go-modularize skill](https://github.com/larsartmann/crush-skills) —
  7-phase workflow that drove the execution
