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
golangci-lint run ./display/... ./errorpage/... ./feedback/... ./forms/... ./htmx/... ./icons/... ./layout/... ./navigation/... ./utils/... ./svg/... ./internal/...
```

## Full verification (do this before every PR)

```bash
find . -name '*_templ.go' -print0 | xargs -0 rm && templ generate ./... && go build ./... && go test ./... && golangci-lint run ./display/... ./errorpage/... ./feedback/... ./forms/... ./htmx/... ./icons/... ./layout/... ./navigation/... ./utils/... ./svg/... ./internal/...
```

---

## Conventions

**Read [`AGENTS.md`](AGENTS.md) first** — it is the canonical reference for all
architecture, code conventions, and gotchas. The highlights:

| Convention               | Rule                                                                                                                                   |
| ------------------------ | -------------------------------------------------------------------------------------------------------------------------------------- |
| Props structs            | Every component's props embed `utils.BaseProps` (exception: `layout.PageProps`).                                                       |
| RTL / logical properties | Never use `ml-`/`mr-`/`pl-`/`pr-`/`left-`/`right-`. Use `ms-`/`me-`/`ps-`/`pe-`/`start-`/`end-`.                                       |
| Motion                   | Use shared constants (`transitionFast`, `transitionNormal`, `transitionColors`, `transitionTransform`). All include `motion-reduce:*`. |
| Typed enums              | `type XxxType string` + typed constants + `IsValid()` method + test. Lookup maps use typed keys. All lookups via `utils.Lookup`.       |
| CSP safety               | Every inline `<script>` uses `nonce={ props.Nonce }`. No exceptions.                                                                   |
| Class merging            | Always use `utils.Class()` for Tailwind conflict resolution (thread-safe via `sync.Mutex`).                                            |
| Style lookups            | Use maps + `utils.Lookup`, not switches. Structural variants use `if`-branch for DOM structure.                                        |
| Zero runtime panics      | Component code must never panic. Enum lookups fall back gracefully.                                                                    |

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

---

## Reporting issues

- Include Go version, templ version, and a minimal reproduction.
- For feature requests, describe the **use case** — not just the solution.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
