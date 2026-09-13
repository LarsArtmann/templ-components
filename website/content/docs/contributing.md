---
title: Contributing
description: Setup, conventions, and release flow.
---

Contributions are welcome! See [CONTRIBUTING.md](https://github.com/larsartmann/templ-components/blob/master/CONTRIBUTING.md) in the repository for the full guide.

## Quick Setup

```bash
git clone https://github.com/larsartmann/templ-components.git
cd templ-components
nix develop  # provides Go, templ, golangci-lint
```

## Build & Test

```bash
# Full build (regenerate templ first)
find . -name '*_templ.go' -print0 | xargs -0 rm && templ generate ./... && go build ./...

# Tests
go test ./...

# Lint
golangci-lint run ./...
```

## Key Conventions

- All props structs embed `utils.BaseProps`
- All root elements propagate `props.Class`, `props.Attrs`, `props.ID`, `props.AriaLabel`
- Class attributes use `utils.Class()` for Tailwind conflict resolution
- Every closed-set enum ships an `IsValid()` method + test
- All inline scripts use `nonce` attributes
- All components have `dark:` variants — enforced by regression tests
