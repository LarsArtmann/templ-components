# AGENTS.md — templ-components

**Updated:** 2026-05-20 | **Coverage:** 67.5% | **Tests:** 179 | **Packages:** 9+demo | **Generated files:** 32 `*_templ.go` committed

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

## CRITICAL: Generated `*_templ.go` Files MUST Be Committed

This is a **templ library**, not an application. The Go module proxy (proxy.golang.org) fetches
source from the Git tag — it does **not** run `templ generate`. Without committed `*_templ.go`
files, consumers get uncompilable code (`undefined` errors on every component function).

- The `.gitignore` uses `!*_templ.go` to override the global gitignore's `*_templ.go` entry
- After editing any `.templ` file, always run `templ generate ./...` and commit the updated `*_templ.go` files alongside the source
- Never add `*_templ.go` back to `.gitignore` — this is the standard pattern for publishable templ packages
- 32 generated files across 8 packages + examples/demo (display, feedback, forms, htmx, icons, internal/svg, layout, navigation)

**Why this matters:** The Go module proxy serves source as-is. Consumers who `go get` this package
will have their Go toolchain download the tagged commit. If `*_templ.go` is missing from that
commit, the package won't compile. Unlike applications (where you generate at build time), a
**library's generated code is part of its distributable artifact**.

## Architecture

- **Module:** `github.com/larsartmann/templ-components`
- **Go:** 1.26, **templ:** v0.3.x
- **No framework deps** — pure Go + templ + Tailwind v4 class strings
- **Import graph:** `utils ← all`, `internal/svg ← display,feedback,icons`, `icons ← display,feedback`, `feedback ← none (htmx decoupled)`
- **No circular imports** allowed

## Code Conventions

- All component props embed `utils.BaseProps` (exception: `layout.PageProps`)
- All root elements propagate `props.Class`, `props.Attrs`, and `props.ID` from BaseProps (25/25 components, including NavLink/MobileNavLink)
- Class attributes use `utils.Class()` for Tailwind conflict resolution (exception: `templ.KV` conditionals where comma-join is required)
- Style lookups use maps/structs, not switches (e.g., `badgeStyleMap`, `badgeSizeLookup`, `cardPaddingLookup`, `iconPathData`, `alertIconMap`, `toastIconMap`, `spinnerSizeLookup`, `progressHeightLookup`, `avatarSizeLookup`, `avatarDotSizeLookup`)
- String enums: `type XxxType string` + `const XxxDefault XxxType = "default"`
- Size constants: uppercase suffix pattern `[Component]Size[SM|MD|LG]` (e.g., `AvatarSizeSM`, `BadgeSizeSM`, `SpinnerSM`)
- Default constructors: `DefaultXxxProps()` for every component with non-zero defaults
- Private helpers: `xxxClass()` for Tailwind class mapping
- CSP: all inline scripts use `nonce={ props.Nonce }`
- Sub-templates: extract shared rendering to private `templ` functions
- Feedback styles: shared `feedbackStyleSet` struct + `lookupFeedbackStyle[T]()` generic
- FeedbackType: canonical `FeedbackType` enum (`FeedbackSuccess/Error/Warning/Info`); `AlertType` and `ToastType` are type aliases for backward compat
- Icons: `iconPathData` map with `|` separator for multi-path icons
- Form errors: `ErrorAttrs(id, errMsg, helpTextID)` helper returns `templ.Attributes` for aria-invalid/aria-describedby
- Card shell CSS: shared `cardShellClass` constant for consistent card styling
- HTMX loading: accepts `templ.Component` spinner parameter (decoupled from feedback package)
- Toast icons: generated from Go `iconPathData` via `icons.IconPathJS()` (single source of truth)
- TrendDirection: `TrendNone = "none"` (non-empty sentinel, not "")
- Layout: `Minimal(MinimalProps)` uses props struct like `Base(PageProps)`
- Modal: focus save/restore via `data-tc-prev-focus` attribute on open, restored on close
- NavLink/MobileNavLink: `NavLink` uses `utils.Class()` for merge; `MobileNavLink` appends `props.Class` to `templ.KV` chain
- InputType: validates via `inputType()` with `validInputTypes` map; panics on unknown, defaults empty to `"text"`
- Structural variants (TabsVariant, DropdownPosition, TrendDirection): use `if`-branch for DOM structure, not map lookup — map pattern is for pure class lookups only
- `forms.SanitizeID`: exported utility for library consumers, not used internally
- Enum validation: 2 panic-on-unknown (InputType, icons.Name), 10 map+fallback, structural variants use if-branch
- Modal/Dropdown/Accordion: ID validation at render time (`validateDropdownID`, `validateModalID`, `validateAccordionItems`) panic on empty
- JS patterns: Accordion + Dropdown use global singleton (`window.tc*Attached`), Modal uses per-instance IIFE (focus trap + focus restore), ThemeToggle uses IIFE-wrapped global guard
- Dismiss JS: Alert + Toast share `tcDismissAttached` handler using `[data-dismiss]` selector
- Toast JS: dismiss icon from `icons.IconPathJS()` via `tcToastIcons.dismiss`
- Table: row cells auto-padded/truncated to match header count
- HTMX retry: per-element `data-tc-retry` attribute (no shared counter)

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
- `LoadingOverlay(message, showProgress, progress)` → `LoadingOverlay(LoadingOverlayProps)` (now has BaseProps + Message + ShowProgress + Progress)
- `TrendNone` = `""` → `"none"` (non-empty sentinel)
- `DefaultStatCardProps()` now sets `Trend: TrendNone`
- `FillIcon(class, path)` → `FillIcon(class, path, rotate bool)` (no longer variadic)
- `AlertType`/`ToastType` → type aliases for `FeedbackType` (deprecated, use FeedbackType directly)
- `Breadcrumbs(items []BreadcrumbItem)` → `Breadcrumbs(BreadcrumbsProps)` (now has BaseProps)
- `utils.BoolString()` → removed, use `strconv.FormatBool` directly
- `utils.Deref/DerefOr/MergeAttrs` → removed (zero production callers)

## Lint Command

```bash
# Must lint specific packages — examples/ excluded via .golangci.yml
golangci-lint run ./display/... ./feedback/... ./forms/... ./htmx/... ./icons/... ./layout/... ./navigation/... ./utils/... ./internal/...
```
