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

## visualtest sibling-pin policy

`visualtest/` is a separate Go module but is **not published** (no proxy tags).
It has its own pin policy, distinct from the seven published modules:

- **Pins:** `visualtest/go.mod` requires its sibling modules (`utils`, `icons`,
  `errorpage`, `htmx`) at the **latest released version** — the same version as
  `utils.Version`. `scripts/release.sh` bumps these pins in step 5b, so they
  move in the release commit itself.
- **Local replaces:** because visualtest is never fetched from the proxy, its
  `go.mod` permanently carries relative `replace` directives for the root
  module and all siblings (like root's local-dev replaces). This keeps
  `GOWORK=off go build/test/tidy` working locally and in CI **before the
  release tags are pushed** — stale pins can no longer leave master red (the
  v1.11.0 incident: 9 days red because pins were bumped by hand after push).
- **Enforcement:** `scripts/check-version-sync.sh` (pre-commit + CI) fails if
  any go.mod — visualtest included — pins a templ-components sibling at a
  version other than `utils.Version`. `scripts/check-module-sync.sh` verifies
  module paths and relative replace paths.

After every release cut, the post-release "re-add replace directives" commit
re-tidies all modules (release.sh step 10b/10c) and asserts the second tidy is
a no-op — the same invariant CI enforces via "Verify no untracked changes".

## Contributing

`go.work` is gitignored. Contributors create it automatically when they run
`nix develop` or `direnv allow` (the `.envrc` sets up the workspace). If
`go.work` is missing, run:

```bash
go work use . utils icons errorpage charts/echarts datastar htmx visualtest
```
