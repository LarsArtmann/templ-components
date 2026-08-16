# Contributing to templ-components

Thanks for your interest. This project is pre-v1.0; APIs may change and
contributions that improve consistency are especially welcome.

---

## Prerequisites

| Tool      | Version           | Notes                                                                                  |
| --------- | ----------------- | -------------------------------------------------------------------------------------- |
| Go        | 1.26+             | Pinned in `go.mod` and `flake.nix`                                                     |
| Nix       | any (recommended) | `nix develop` provides Go, `golangci-lint`, and **templ v0.3.1020** (matches `go.mod`) |
| templ CLI | v0.3.1020         | The dev shell pins this; do **not** use a system binary that may be v0.3.1036+         |

> **Why Nix?** The system `templ` binary may be an unreleased upstream build.
> Always use `nix develop` before generating. See [`AGENTS.md`](AGENTS.md).

---

## Build

```bash
nix develop

# Regenerate all *_templ.go from .templ sources, then build
find . -name '*_templ.go' -print0 | xargs -0 rm && templ generate ./... && go build ./...
```

## Test

```bash
go test ./...
```

## Lint

```bash
golangci-lint run ./...
```

## Full verification (do this before every PR)

```bash
find . -name '*_templ.go' -print0 | xargs -0 rm && templ generate ./... && go build ./... && go test ./... && golangci-lint run ./...
```

---

## Conventions

**Read [`AGENTS.md`](AGENTS.md) first** — it is the canonical reference for all
architecture, code conventions, and gotchas. The highlights:

| Convention               | Rule                                                                                                                                                                                                                                                                      |
| ------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Props structs            | Every component's props embed `utils.BaseProps` (exception: `layout.PageProps`).                                                                                                                                                                                          |
| RTL / logical properties | Never use `ml-`/`mr-`/`pl-`/`pr-`/`left-`/`right-`. Use `ms-`/`me-`/`ps-`/`pe-`/`start-`/`end-`.                                                                                                                                                                          |
| Motion                   | Use shared constants (`transitionFast`, `transitionNormal`, `transitionColors`, `transitionTransform`). All include `motion-reduce:*`.                                                                                                                                    |
| Dark mode                | Every neutral and semantic color MUST have a `dark:` variant. `-600` light → `-500` dark for backgrounds, `-400` for text. Use `gray-*` only. Enforced by `TestDarkModeCompliance` + `TestDarkModeSemanticColors`. See [ADR 0011](docs/adr/0011-dark-mode-convention.md). |
| Typed enums              | `type XxxType string` + typed constants + `IsValid()` method + test. Lookup maps use typed keys. All lookups via `utils.Lookup`.                                                                                                                                          |
| CSP safety               | Every inline `<script>` uses `nonce={ props.Nonce }`. No exceptions.                                                                                                                                                                                                      |
| Class merging            | Always use `utils.Class()` for Tailwind conflict resolution (thread-safe via `sync.Mutex`).                                                                                                                                                                               |
| Style lookups            | Use maps + `utils.Lookup`, not switches. Structural variants use `if`-branch for DOM structure.                                                                                                                                                                           |
| Zero runtime panics      | Component code must never panic. Enum lookups fall back gracefully.                                                                                                                                                                                                       |
| Nested borders           | Components nested inside `Card(CardPaddingNone)` must provide a `Flush`-style option to suppress their own border. See [ADR 0012](docs/adr/0012-flush-prop-for-nested-borders.md) and [Table-in-Card recipe](docs/recipes/table-in-card.md).                              |

### Commit messages

Follow [Conventional Commits](https://www.conventionalcommits.org):

```
feat(display): add carousel component
fix(feedback): correct toast dismiss timing
docs: update quick start example
```

---

## Generated files: `*_templ.go` MUST be committed

This is a **templ library**, not an application. The Go module proxy
(`proxy.golang.org`) fetches source from the Git tag — it does not run
`templ generate`. Without committed `*_templ.go` files, consumers who `go get`
this package get uncompilable code.

After editing any `.templ` file, run `templ generate ./...` and commit the
updated `*_templ.go` files alongside the source change. See [`AGENTS.md`](AGENTS.md).

---

## Release

```bash
scripts/release.sh <new-version> "<release-summary>"
```

One-commit convention. SSH-signed tags. House rule: **never push automatically**.
See [`AGENTS.md`](AGENTS.md) § Release Convention for details.

### Release checklist (lockstep tagging)

Every root tag must ship with a `<sub-module>/v<version>` tag for **every**
published sub-module — `utils/`, `icons/`, `errorpage/`, `htmx/`, `datastar/`,
`charts/echarts/` — all pointing at the same release commit. Consumers pin
sub-modules in go.mod; a root-only release leaves those pins unresolvable on
the Go module proxy (v1.8.3 shipped root-only and broke every dependent
build).

Before pushing:

1. `scripts/release.sh` created the tags (the sub-module set is derived from
   the root go.mod's replace directives — adding a published module is a
   replace-directive edit, nothing else).
2. `scripts/check-release-tags.sh` passes (names any missing or diverged
   sub-module tag).
3. Push with `git push origin master --follow-tags` so the sibling tags ride
   along.

The `release-tags` CI job runs the guard on every `v*` tag push and fails the
release if the set is incomplete.

---

## Reporting issues

- Include Go version, templ version, and a minimal reproduction.
- For feature requests, describe the **use case** — not just the solution.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
