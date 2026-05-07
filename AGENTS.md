# AGENTS.md — templ-components

**Updated:** 2026-05-07

## Build & Test Commands

```bash
# Full build (required before go build after .templ changes)
find . -name '*_templ.go' -print0 | xargs -0 rm && templ generate ./... && go build ./...

# Tests
go test ./...

# Lint
golangci-lint run ./...

# All-in-one verification
find . -name '*_templ.go' -print0 | xargs -0 rm && templ generate ./... && go build ./... && go test ./... && golangci-lint run ./...
```

## Architecture

- **Module:** `github.com/larsartmann/templ-components`
- **Go:** 1.26, **templ:** v0.3.1001
- **No framework deps** — pure Go + templ + Tailwind class strings
- **Import graph:** `utils ← all`, `internal/svg ← display,feedback,icons`, `icons ← display`
- **No circular imports** allowed

## Code Conventions

- All component props embed `utils.BaseProps` (exception: `layout.PageProps`)
- All root elements propagate `props.Class`, `props.Attrs`, and `props.ID` from BaseProps
- Class attributes use `utils.Class()` for Tailwind conflict resolution (exception: `templ.KV` conditionals where comma-join is required)
- Style lookups use maps, not switches (e.g., `badgeColorMap`, `iconPathData`)
- String enums: `type XxxType string` + `const XxxDefault XxxType = "default"`
- Size constants: uppercase suffix pattern `[Component]Size[SM|MD|LG]` (e.g., `AvatarSizeSM`, `BadgeSizeSM`, `SpinnerSM`)
- Default constructors: `DefaultXxxProps()` for every component with non-zero defaults
- Private helpers: `xxxClass()` for Tailwind class mapping
- CSP: all inline scripts use `nonce={ props.Nonce }`
- Sub-templates: extract shared rendering to private `templ` functions
- Feedback styles: shared `feedbackStyleSet` struct + `lookupFeedbackStyle[T]()` generic
- Icons: `iconPathData` map with `|` separator for multi-path icons
- Form errors: `ErrorAttrs(id, errMsg)` helper returns `templ.Attributes` for aria-invalid/aria-describedby

## Breaking Changes (v0.1 → v0.2)

- `AvatarProps.Online/Offline bool` → `AvatarStatus` enum
- `StatCard(value, label, change, positive)` → `StatCardProps` struct with `TrendDirection` enum
- `PageProps.HTMXSRI string` → `HTMXUseSRI bool`
- `SecurityHeaders` defaults to `true` (was implicit `false`)
- `DropdownItem.Icon string` → `icons.Name`
- `EmptyStateProps.Icon string` → `icons.Name`
- `BadgeSizeSm/Md/Lg` → `BadgeSizeSM/MD/LG` (uppercase suffix)
- `SpinnerSmall/Medium/Large` → `SpinnerSM/MD/LG` (uppercase suffix)
- `TabsStyle`/`TabsStylePills` → `TabsVariant`/`TabsPills`, field `TabStyle` → `Variant`

## Lint Command

```bash
# Must lint specific packages — examples/ has 23 issues that can't be excluded via config
golangci-lint run ./display/... ./feedback/... ./forms/... ./htmx/... ./icons/... ./layout/... ./navigation/... ./utils/... ./internal/...
```
