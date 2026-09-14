# Version Support Policy

How old can your toolchain be, and what breaks when floors move? This page
is the contract.

## Supported floors (v1.17.0)

| Dependency   | Floor                                     | Why                                                                                                                                                                                                          |
| ------------ | ----------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Go           | 1.26                                      | `encoding/json/v2` (via `errorpage`, `navigation/breadcrumbs`) needs `GOEXPERIMENT=jsonv2` below 1.27; the toolchain floor tracks the latest security release (currently 1.26.7 for GO-2026-5972/6089/6090). |
| templ        | v0.3.1020                                 | The generator version the committed `*_templ.go` files are built with. Older generators emit incompatible runtime calls; newer unpublished versions are not on the proxy.                                    |
| Tailwind CSS | v4.x (4.1+)                               | CSS-first config, `@theme`, `@custom-variant`, `@source` scanning of `.templ` files. v3 does not support the theming model.                                                                                  |
| HTMX         | 2.0.10 (self-hosted, embedded) or 2.x CDN | The library embeds 2.0.10 inline; `HTMXVersion` consumers need 2.x for response-handling defaults the components rely on.                                                                                    |

## What a floor bump looks like

1. The floor moves **only** when a dependency feature or a security fix is
   actually required — never speculatively.
2. The bump lands with: the `go.mod`/`go.work` change, a CHANGELOG
   **Changed** entry naming the new floor and the reason, and a note here.
3. Consumers on older toolchains get a clear build error (missing
   `encoding/json/v2`, unknown runtime symbols) — not silent misbehavior.
4. `tc doctor` reports the current floors from inside your project.

## Deprecation and removal

- A deprecated API ships a CHANGELOG **Deprecated** entry and stays for at
  least one minor release before removal in the next major (see
  `docs/adr/0039-v2-module-path-timing.md` for the module-path mechanics).
- Breaking changes accumulate in TODO_LIST's _Deferred — v2.0 breaking
  changes_ section; nothing there ships in a v1 minor.
