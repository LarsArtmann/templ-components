# AGENTS.md — templ-components

**Updated:** 2026-05-18

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
- **No framework deps** — pure Go + templ + Tailwind v4 class strings
- **Import graph:** `utils ← all`, `internal/svg ← display,feedback,icons`, `icons ← display,feedback`, `feedback ← none (htmx decoupled)`
- **No circular imports** allowed

## Code Conventions

- All component props embed `utils.BaseProps` (exception: `layout.PageProps`)
- All root elements propagate `props.Class`, `props.Attrs`, and `props.ID` from BaseProps
- Class attributes use `utils.Class()` for Tailwind conflict resolution (exception: `templ.KV` conditionals where comma-join is required)
- Style lookups use maps/structs, not switches (e.g., `badgeStyleMap`, `badgeSizeLookup`, `cardPaddingLookup`, `iconPathData`)
- String enums: `type XxxType string` + `const XxxDefault XxxType = "default"`
- Size constants: uppercase suffix pattern `[Component]Size[SM|MD|LG]` (e.g., `AvatarSizeSM`, `BadgeSizeSM`, `SpinnerSM`)
- Default constructors: `DefaultXxxProps()` for every component with non-zero defaults
- Private helpers: `xxxClass()` for Tailwind class mapping
- CSP: all inline scripts use `nonce={ props.Nonce }`
- Sub-templates: extract shared rendering to private `templ` functions
- Feedback styles: shared `feedbackStyleSet` struct + `lookupFeedbackStyle[T]()` generic
- Icons: `iconPathData` map with `|` separator for multi-path icons
- Form errors: `ErrorAttrs(id, errMsg, helpTextID)` helper returns `templ.Attributes` for aria-invalid/aria-describedby
- Card shell CSS: shared `cardShellClass` constant for consistent card styling
- HTMX loading: accepts `templ.Component` spinner parameter (decoupled from feedback package)
- Toast icons: generated from Go `iconPathData` via `icons.IconPathJS()` (single source of truth)
- TrendDirection: `TrendNone = "none"` (non-empty sentinel, not "")
- Layout: `Minimal(MinimalProps)` uses props struct like `Base(PageProps)`
- Modal/Dropdown: JS IDs escaped via `strconv.Quote()` (XSS prevention)
- Table: row cells auto-padded/truncated to match header count

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
- `Tab.Active bool` → `TabsProps.ActiveTabID string` (prevents zero/multiple active)
- `BadgeDefault` → removed, use `BadgeNeutral`
- `ErrorAttrs(id, errMsg)` → `ErrorAttrs(id, errMsg, helpTextID)`
- `Minimal(title, locale)` → `Minimal(MinimalProps)`
- `LoadingIndicator()` → `LoadingIndicator(spinner templ.Component)` (decoupled from feedback)
- `InlineLoadingOverlay(id)` → `InlineLoadingOverlay(id, spinner templ.Component)`
- `LoadingButton(default, loading)` → `LoadingButton(default, loading, spinner templ.Component)`
- `SimpleCard()` → `SimpleCard(SimpleCardProps)` (now has BaseProps + CardPadding)
- `TrendNone` = `""` → `"none"` (non-empty sentinel)
- `DefaultStatCardProps()` now sets `Trend: TrendNone`

## Lint Command

```bash
# Must lint specific packages — examples/ excluded via .golangci.yml
golangci-lint run ./display/... ./feedback/... ./forms/... ./htmx/... ./icons/... ./layout/... ./navigation/... ./utils/... ./internal/...
```
