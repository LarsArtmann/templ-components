# Migrating from v1.x to v2.0

v2.0 is a major release with breaking changes across four areas: module
structure, HTMX defaults, container-query defaults, and type aliases.
This guide walks through each change and what you need to do.

## Quick summary

| Change | Impact | Migration effort |
|--------|--------|-----------------|
| 7-module workspace split | Import paths unchanged; `internal/*` → `utils/*` | None (if you didn't import `internal/`) |
| HTMX self-host by default | HTMX embedded inline, no CDN request | Set `HTMXSrc: ""` to keep CDN |
| Container-aware by default | Grid, Card, Split use container queries | Set `ContainerAware: false` for viewport |
| Alias removal | `AlertType`/`ToastType` removed | Rename to `FeedbackType` |

---

## 1. Module structure: 7-module workspace

The library is now 7 Go modules coordinated by `go.work` (local dev) and
`replace` directives (CI/consumers). See ADR-0034.

### What changed

- `internal/svg`, `internal/cdn`, `internal/golden` promoted to `utils/svg`,
  `utils/cdn`, `utils/golden` (Go's `internal/` rule blocks cross-module access)
- **All externally-importable package paths are unchanged** — `display`,
  `feedback`, `forms`, `layout`, `navigation`, `htmx`, `datastar`, `recipes`,
  `errorpage`, `icons`, `charts/echarts` all have the same import paths

### What you need to do

**If you never imported `internal/*` packages** (most consumers): nothing.
Your imports continue to work as-is.

**If you imported `internal/svg`, `internal/cdn`, or `internal/golden`**: update
your imports:

```go
// Before (v1.x)
import "github.com/larsartmann/templ-components/internal/svg"
import "github.com/larsartmann/templ-components/internal/golden"

// After (v2.0)
import "github.com/larsartmann/templ-components/utils/svg"
import "github.com/larsartmann/templ-components/utils/golden"
```

### Individual module adoption

You can now adopt individual modules without the full library:

```bash
go get github.com/larsartmann/templ-components/icons@latest        # icons only
go get github.com/larsartmann/templ-components/errorpage@latest    # error pages only
go get github.com/larsartmann/templ-components/charts/echarts@latest  # ECharts adapter only
```

---

## 2. HTMX self-host by default

HTMX is now **embedded inline** in the page by default (v2.0). No external CDN
request is made. The embedded version is HTMX 2.0.10.

### What changed

- `DefaultPageProps()` now sets `HTMXSrc: "self"` (embeds HTMX inline)
- The CDN path (`HTMXVersion` + `HTMXCDN` + SRI) is still available but opt-in
- `HTMXSelfHost` is a sentinel constant (`"self"`) that triggers inline embedding

### What you need to do

**If you want to keep using the CDN** (e.g., for caching benefits):

```go
props := layout.DefaultPageProps()
props.HTMXSrc = "" // clear self-host sentinel to use CDN
// HTMXVersion is already set to "2.0.10" in DefaultPageProps
```

**If you want self-host (new default)**: nothing. It just works.

> **CSP requirement:** The embedded HTMX renders as `<script nonce="...">`. Your
> Content Security Policy must allow `script-src 'nonce-...'` (or `'unsafe-inline'`).
> If your CSP only allows `script-src https://cdn.jsdelivr.net`, the inline HTMX
> will be blocked. Either add `'nonce-...'` to your CSP or switch back to CDN mode
> (`HTMXSrc: ""`).

**If you were already self-hosting via a custom path**: keep your existing
`HTMXSrc` value (e.g., `"/static/htmx.min.js"`). The `"self"` sentinel only
activates when `HTMXSrc == "self"`.

### Trade-offs

| Self-host (default) | CDN (opt-in) |
|---------------------|-------------|
| No external request | Browser caching across sites |
| ~50KB inline per page | Shared cache via CDN |
| No SRI needed (same-origin) | SRI verification via integrity attr |
| Works offline / air-gapped | Requires network access |

---

## 3. Container-aware by default

Three components now default to container-query-based responsiveness instead of
viewport breakpoints. See ADR-0018.

### What changed

| Component | Field | v1.x default | v2.0 default |
|-----------|-------|-------------|-------------|
| `Grid` | `ContainerAware` | `false` (viewport) | `true` (container) |
| `Card` | `ContainerAware` | `false` (viewport) | `true` (container) |
| `Split` | `ContainerAware` | `false` (viewport) | `true` (container) |

Additionally, `Grid.ContainerResponsive` has been **renamed** to
`Grid.ContainerAware` for consistency with all other components.

### What you need to do

**If you used `ContainerResponsive: true` on Grid**: rename to `ContainerAware`.
The default is now `true`, so you can omit it entirely.

```go
// Before (v1.x)
display.GridProps{Cols: display.GridCols3, ContainerResponsive: true}

// After (v2.0) — ContainerAware is now the default, just omit it
display.GridProps{Cols: display.GridCols3}

// Or if you want viewport-based breakpoints (the old default):
display.GridProps{Cols: display.GridCols3, ContainerAware: false}
```

**Components that stayed viewport-default**: `Nav`, `Pagination`, `Form`,
`DefinitionGrid`, `SkeletonCardGrid` — their `ContainerAware` flag remains
opt-in (default `false`).

---

## 4. Deprecated alias removal

Type aliases and constants that were deprecated since v0.x have been removed.

### What changed

| Removed | Replacement |
|---------|------------|
| `feedback.AlertType` | `feedback.FeedbackType` |
| `feedback.ToastType` | `feedback.FeedbackType` |
| `feedback.AlertSuccess` | `feedback.FeedbackSuccess` |
| `feedback.AlertError` | `feedback.FeedbackError` |
| `feedback.AlertWarning` | `feedback.FeedbackWarning` |
| `feedback.AlertInfo` | `feedback.FeedbackInfo` |
| `feedback.ToastSuccess` | `feedback.FeedbackSuccess` |
| `feedback.ToastError` | `feedback.FeedbackError` |
| `feedback.ToastWarning` | `feedback.FeedbackWarning` |
| `feedback.ToastInfo` | `feedback.FeedbackInfo` |
| `Grid.ContainerResponsive` | `Grid.ContainerAware` |

### What you need to do

Find-and-replace:

```go
// Before (v1.x)
Type: feedback.AlertError
Type: feedback.ToastSuccess

// After (v2.0)
Type: feedback.FeedbackError
Type: feedback.FeedbackSuccess
```

Struct field types changed too, but since `AlertType` and `ToastType` were
aliases for `FeedbackType`, the underlying type is the same. The `AlertProps`
and `ToastProps` structs still exist — only the type of the `Type` field
changed from `AlertType`/`ToastType` to `FeedbackType`.

---

## Semantic tokens (non-breaking)

`templates/app.css` now imports `templ-components-theme.css` by default (it was
previously marked "optional"). This means semantic color tokens
(`--color-tc-primary`, `--color-tc-danger`, etc.) are available out of the box.

**No action required.** If you had a custom `@import "./templ-components-theme.css"`
in your CSS, you can remove it (the starter `app.css` now includes it). If you
were overriding colors via `@theme { --color-blue-600: #custom; }`, that
continues to work.

---

## Checklist

- [ ] Update imports: `internal/svg` → `utils/svg`, `internal/golden` → `utils/golden`, `internal/cdn` → `utils/cdn` (if applicable)
- [ ] Rename `AlertType` → `FeedbackType`, `ToastType` → `FeedbackType` (if used)
- [ ] Rename `AlertSuccess` → `FeedbackSuccess`, etc. (if used)
- [ ] Rename `ToastSuccess` → `FeedbackSuccess`, etc. (if used)
- [ ] Rename `ContainerResponsive` → `ContainerAware` on `GridProps` (if used)
- [ ] Set `ContainerAware: false` on Grid/Card/Split if you need viewport breakpoints
- [ ] Set `HTMXSrc: ""` if you want to keep using the HTMX CDN
- [ ] Update CSP: allow `script-src 'nonce-...'` for inline HTMX (or keep CDN mode)
- [ ] Run `go build ./...` and fix any compile errors from the above changes
