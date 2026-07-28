# ADR-0022: v2.0 Default-Flip Migration Plan

**Status:** Draft (2026-07-28)
**Decider:** Lars Artmann

## Context

Three opt-in features have proven stable and valuable since their introduction:

1. **Self-host HTMX** (`PageProps.HTMXSrc`) — shipped v0.22.0 as opt-in.
   Currently defaults to the jsDelivr CDN.
2. **Semantic token layer** (`templ-components-theme.css`) — shipped v0.22.0
   as opt-in. Currently consumers must explicitly import the theme file.
3. **Container-aware mode** (`ContainerAware`/`ContainerResponsive`) — shipped
   v0.21.0–[Unreleased] as opt-in on 8 components. Currently defaults to
   viewport breakpoints.

Each opt-in adds a cognitive step for new consumers. v2.0 flips these
defaults so the "better" behavior is the out-of-the-box experience.

## Decision

Flip all three defaults in a single major version:

### 1. HTMX Self-Host by Default

- `PageProps.HTMXCDN` → removed
- `PageProps.HTMXSrc` → defaults to `"self"` (embeds HTMX via Go's `embed`)
- CDN becomes opt-in: `PageProps.HTMXCDN = "https://cdn.jsdelivr.net/npm"`
- **Migration:** consumers who relied on the CDN set `HTMXCDN` explicitly

### 2. Semantic Tokens by Default

- `templates/app.css` includes `@import "./templ-components-theme.css"` by default
- Raw Tailwind color classes (`bg-blue-600`) become aliases for semantic tokens
- **Migration:** consumers who overrode colors via `@theme { --color-blue-600 }`
  keep working (the alias layer resolves to the same value)

### 3. Container-Aware by Default (selective)

Only components commonly placed in constrained containers flip:
- `Grid.ContainerResponsive` → default `true`
- `Card.ContainerAware` → default `true`
- `Split.ContainerAware` → default `true`

Components where viewport is usually correct stay viewport-default:
- `Nav`, `Pagination`, `Form`, `DefinitionGrid`, `SkeletonCardGrid` — stay opt-in

- **Migration:** consumers who need viewport behavior set `ContainerAware: false`

### 4. Deprecated Alias Removal

- `AlertType`/`ToastType` aliases for `FeedbackType` — removed
- `Grid.ContainerResponsive` renamed to `Grid.ContainerAware` for consistency

## Migration Guide Structure

```
docs/migration/v1-to-v2.md
├── HTMX: self-host → CDN opt-in
├── Semantic tokens: automatic → override if needed
├── ContainerAware: 3 components flip → set false for viewport
├── Aliases removed: AlertType → FeedbackType, etc.
└── Deprecation timeline: v1.3 warnings → v2.0 removal
```

## Deprecation Timeline

| Version | Action                                              |
|---------|-----------------------------------------------------|
| v1.3.0  | Add `// Deprecated` comments + `HTMXCDN` log warning |
| v1.4.0  | Default `HTMXSrc = "self"` with CDN opt-in (breaking for CDN users) |
| v2.0.0  | Full default flip + alias removal                   |

## Consequences

- **Breaking changes** require major version bump
- **Migration guide** must be comprehensive (3 features × migration paths)
- **Semantic tokens** may produce visual diffs if consumers had custom `@theme` overrides
- **Container-aware default** changes responsive behavior for Grid/Card/Split
