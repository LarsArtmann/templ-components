# Context — templ-components

**Updated:** 2026-05-07

## What

A Go component library built on [templ](https://templ.guide) and [Tailwind CSS](https://tailwindcss.com) for building server-rendered web applications with HTMX.

## Tech Stack

| Layer         | Technology                           |
| ------------- | ------------------------------------ |
| Language      | Go 1.26                              |
| Templates     | templ v0.3.1001                      |
| Styling       | Tailwind CSS (via class strings)     |
| Class merging | tailwind-merge-go v0.2.1             |
| Interactivity | HTMX 2.0.6 + vanilla JS              |
| Build         | `templ generate` + `go build`        |
| CI            | GitHub Actions (lint + build + test) |

## Package Layout

```
templ-components/
├── utils/           # Base types, Tailwind class merging, generic helpers (BoolString, MapEnum)
├── internal/svg/    # Shared SVG primitives (fillIcon, spinner)
├── display/         # UI: card, badge, modal, table, tabs, avatar, tooltip, accordion, dropdown
├── feedback/        # User feedback: alert, toast, spinner, progress, skeleton (shared feedbackStyleSet)
├── forms/           # Form controls: input, select, textarea, checkbox, label
├── htmx/            # HTMX helpers: loading, error handling, CSRF, OOB swap
├── icons/           # Named SVG icons (42 constants, map-driven rendering)
├── layout/          # Page layout: base HTML, theme toggle, dark mode
└── navigation/      # Nav: navbar, breadcrumbs, pagination, mobile menu
```

### Import Graph

```
utils          ← all packages
internal/svg   ← display, feedback, icons
icons          ← display (empty_state)
```

No circular imports. `internal/svg` is not importable by consumers (Go `internal/` convention).

## Key Patterns

### Props Embedding

All component Props structs embed `utils.BaseProps`:

```go
type CardProps struct {
    utils.BaseProps      // ID, Class, Attrs, AriaLabel, Nonce
    Title string
    // ...
}
```

Exception: `layout.PageProps` (page metadata, not component props).

### Style Lookups

Style maps (not switches) for all visual variant lookups:

```go
var badgeColorMap = map[BadgeType]string{...}
func badgeColorClass(t BadgeType) string { ... }
```

Feedback styles share a common struct with a generic lookup:

```go
type feedbackStyleSet struct { Border, BG, Text, Icon string }
func lookupFeedbackStyle[T ~string](m map[T]feedbackStyleSet, def feedbackStyleSet, t T) feedbackStyleSet
```

### Icon Rendering

Icons use a map-driven approach instead of a switch:

```go
var iconPathData = map[Name]string{ Home: "M2.25 12l8.954...", ... }
// Multi-path icons use "|" separator
var iconPathData = map[Name]string{ Eye: "M2.036...|M15 12...", ... }
```

### Enum Types (Impossible States Unrepresentable)

| Type             | Values                | Replaces                    |
| ---------------- | --------------------- | --------------------------- |
| `AvatarStatus`   | Online, Offline, None | Two bool fields (both true) |
| `TrendDirection` | Up, Down, None        | `positive bool` on StatCard |

### CSP Compliance

All inline `<script>` tags use `nonce={ nonce }` or `nonce={ props.Nonce }`.

### Sub-templates

Complex components extract shared rendering logic into private sub-templates:

```go
templ fillIcon(class, path string, rotate ...bool) { ... }  // display/
templ strokeIcon(class string, paths []string) { ... }        // icons/
templ paginationArrow(enabled, href, ...) { ... }             // navigation/
templ inlineMessage(message, colorClass, ...) { ... }        // feedback/
```

## Dependencies

- `github.com/a-h/templ` — template engine
- `github.com/Oudwins/tailwind-merge-go` — Tailwind class conflict resolution

No other runtime dependencies.

## Naming Conventions

| Pattern             | Example              | Purpose                                   |
| ------------------- | -------------------- | ----------------------------------------- |
| `XxxProps`          | `CardProps`          | Component configuration struct            |
| `XxxType`           | `AlertType`          | String enum for visual variants           |
| `XxxSize`           | `BadgeSize`          | String enum for size variants             |
| `XxxPosition`       | `TooltipPosition`    | String enum for positional variants       |
| `XxxStatus`         | `AvatarStatus`       | String enum for state variants            |
| `XxxDirection`      | `TrendDirection`     | String enum for directional variants      |
| `DefaultXxxProps()` | `DefaultCardProps()` | Constructor with sensible defaults        |
| `xxxClass()`        | `badgeColorClass()`  | Unexported: enum → Tailwind class mapping |

## Architecture Decisions

1. **`utils.BaseProps` over per-component fields** — Shared ID/Class/Attrs/AriaLabel/Nonce across all components
2. **`internal/svg` package** — Centralized SVG primitives to avoid cross-package import issues
3. **Map-based style lookups** — Data-driven, extensible, consistent across packages
4. **`layout.PageProps` (not BaseProps)** — Different purpose, different name to avoid confusion
5. **String enums** — Type-safe without code generation; `type XxxType string` + constants
6. **No framework dependencies** — Pure Go + templ; Tailwind classes are strings
7. **`feedbackStyleSet` + generic lookup** — Shared style struct with `lookupFeedbackStyle[T]()` eliminates per-component duplicate types
8. **`iconPathData` map** — Data-driven icon rendering replaces switch statements; multi-path icons use `|` separator
9. **`AvatarStatus` / `TrendDirection` enums** — Impossible states unrepresentable; boolean pairs eliminated
10. **`utils.BoolString()`** — Standardized boolean-to-string conversion replacing local implementations
