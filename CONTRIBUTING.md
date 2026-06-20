# Contributing to templ-components

Thanks for your interest. This project is in early development (pre-release), so APIs may change and contributions that improve consistency are especially welcome.

## Setup

Requirements:

- Go 1.26+
- [templ CLI](https://templ.guide/quick-start/installation)

```bash
git clone https://github.com/larsartmann/templ-components.git
cd templ-components
templ generate ./...
go build ./...
go test ./...
```

## Development Workflow

1. Create a feature branch from `master`
2. Make your changes
3. Run the full verification:

```bash
find . -name '*_templ.go' -print0 | xargs -0 rm && templ generate ./... && go build ./... && go test ./... && golangci-lint run ./display/... ./errorpage/... ./feedback/... ./forms/... ./htmx/... ./icons/... ./layout/... ./navigation/... ./utils/... ./internal/...
```

4. Open a pull request

## Code Conventions

- All component props embed `utils.BaseProps` (exception: `layout.PageProps`)
- All root elements propagate `props.Class`, `props.Attrs`, and `props.ID`
- Style lookups use maps, not switches
- String enums: `type XxxType string` + constants
- Size constants: uppercase suffix `[Component]Size[SM|MD|LG]`
- Default constructors: `DefaultXxxProps()` with meaningful non-zero defaults
- Class merging: always use `utils.Class()` for Tailwind conflict resolution. Thread-safe via `sync.Mutex` protecting the full Merge() call sequence.
- CSP: all inline scripts use `nonce={ props.Nonce }`
- JS IDs: escape with `strconv.Quote()` to prevent XSS (see `validateDropdownID`)
- No external dependencies beyond `templ`, `tailwind-merge-go`, and `go-error-family` (errorpage only)

## Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org):

```
feat(display): add carousel component
fix(feedback): correct toast dismiss timing
docs: update quick start example
```

## Reporting Issues

- Include Go version, templ version, and a minimal reproduction
- For feature requests, describe the use case — not just the solution

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
