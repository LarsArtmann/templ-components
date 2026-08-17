# Modularization

This directory documents the multi-module workspace split executed on 2026-08-10.

## What was done

The single-module library was split into **7 Go modules** connected by a
`go.work` workspace (local dev) and `replace` directives (CI/consumers).

```
Layer 0 (leaf):  utils (utils, utils/svg, utils/cdn, utils/golden)
Layer 1:         icons, charts/echarts, datastar, htmx    [depend on utils]
Layer 2:         errorpage                [depends on utils, icons]
Layer 3:         root (display, feedback, forms, layout, navigation,
                          recipes, integration, cmd/tc,
                          internal/contract, examples/demo)
```

## Why 7 modules (not per-package)

A full per-package split (~12 modules) was proposed in ADR-0020 but rejected
because:

- `utils/svg` is a universal dependency; splitting adds overhead without payoff
- `utils/golden` is used by 8 packages' tests; per-package splitting creates a
  test-helpers dependency web
- Consumer demand triggers (ADR-0020) were still unmet
- The targeted split captures the real isolation wins: `go-error-family`
  (errorpage), icons-only adoption, ECharts opt-in

## Key decisions

| Decision             | Choice                               | Rationale                                                                |
| -------------------- | ------------------------------------ | ------------------------------------------------------------------------ |
| Scope                | Targeted 7-module                    | Real isolation wins without per-package overhead                         |
| Versioning           | Shared (one version for all modules) | Simplest for solo-maintained lib                                         |
| `internal/` packages | Promoted to `utils/` sub-packages    | Go's `internal/` rule blocks cross-module access                         |
| Compat               | Import paths unchanged               | `internal/*` was never externally importable; sub-module paths unchanged |

## Files

- **[ADR-0034](../adr/0034-targeted-module-split.md)** — the executed split decision
- **[ADR-0020](../adr/0020-per-package-modules-split.md)** — the original proposal (superseded)
- **`2026-07-18_10-50_assessment.html`** — historical assessment from the analysis phase

## Release process

At release time, `scripts/release.sh` tags all 7 modules:

```
v<version>              # root module
utils/v<version>        # utils
icons/v<version>        # icons
errorpage/v<version>    # errorpage
charts/echarts/v<version>  # charts/echarts
datastar/v<version>     # datastar
htmx/v<version>         # htmx
```

The Go module proxy resolves sub-module versions via these directory-prefixed
tags. The `replace` directives in each `go.mod` are for local dev only; proxy
consumers resolve via the published tags.

## Contributing

`go.work` is gitignored. Contributors create it automatically when they run
`nix develop` or `direnv allow` (the `.envrc` sets up the workspace). If
`go.work` is missing, run:

```bash
go work use . utils icons errorpage charts/echarts datastar htmx visualtest
```
