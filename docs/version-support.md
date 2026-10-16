# Version Support Policy

How old can your toolchain be, and what breaks when floors move? This page
is the contract.

## Supported floors (v1.17.0)

| Dependency   | Floor                                     | Why                                                                                                                                                                                                                                                                                                                                                |
| ------------ | ----------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Go           | 1.27                                      | `encoding/json/v2` (via `errorpage`, `navigation/breadcrumbs`) is stable on Go 1.27 — no `GOEXPERIMENT` flag needed. The floor moved because dependencies (chromedp cdproto) required >= 1.27 (2026-10-05 sweep, all directives in lockstep); the toolchain floor tracks the latest security release (currently 1.27.1 via the flake's `go_1_27`). |
| templ        | v0.3.1020                                 | The generator version the committed `*_templ.go` files are built with. Older generators emit incompatible runtime calls; newer unpublished versions are not on the proxy.                                                                                                                                                                          |
| Tailwind CSS | v4.x (4.1+)                               | CSS-first config, `@theme`, `@custom-variant`, `@source` scanning of `.templ` files. v3 does not support the theming model.                                                                                                                                                                                                                        |
| HTMX         | 2.0.10 (self-hosted, embedded) or 2.x CDN | The library embeds 2.0.10 inline; `HTMXVersion` consumers need 2.x for response-handling defaults the components rely on.                                                                                                                                                                                                                          |

## What a floor bump looks like

1. The floor moves **only** when a dependency feature or a security fix is
   actually required — never speculatively.
2. The bump lands with: the `go.mod`/`go.work` change, a CHANGELOG
   **Changed** entry naming the new floor and the reason, and a note here.
3. Consumers on older toolchains get a clear build error (missing
   `encoding/json/v2`, unknown runtime symbols) — not silent misbehavior.
4. `tc doctor` reports the current floors from inside your project.
5. **Every toolchain pin moves in the SAME commit** (the go-1.27 sweep
   missed half of these on the first pass): `go.work` + ALL module
   `go.mod` directives (identical spelling — mixed forms break workspace
   builds), the flake `nixpkgs-go` pin, the golangci-lint pin (verify it
   supports the new Go — ≤ v2.13.2 panics on 1.27), CI `go-version` (or
   `go-version-file: go.mod`), and the demo Dockerfile `golang:` image.
6. **Regenerate BOTH golden families**: HTML goldens
   (`go test ./pkg/... -update`) and the website's page goldens
   (`cd website && go test ./internal/pages/ -update`) — the website suite
   was missed entirely in the go-1.27 sweep and went red only later.
7. **Sweep the docs floors**: prose claims pinned by `TestDocsCountDrift`
   (test-function counts move with every sweep), this page, and AGENTS.md's
   Go-floor line.
8. **Consumer-note probe**: after the release, smoke one consumer tag
   (`scripts/check-tag-compiles.sh <version>`) — the floor is only real if
   a tagged consumer build passes on the minimum toolchain.

## Deprecation and removal

- A deprecated API ships a CHANGELOG **Deprecated** entry and stays for at
  least one minor release before removal in the next major (see
  `docs/adr/0039-v2-module-path-timing.md` for the module-path mechanics).
- Breaking changes accumulate in TODO_LIST's _Deferred — v2.0 breaking
  changes_ section; nothing there ships in a v1 minor.
