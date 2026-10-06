---
title: Error Pages
description: Render family-aware, CSP-safe error pages from any handler — go-error-family integration, the FromError pipeline, JSON responses, and the demo routes.
---

The `errorpage` package turns a classified error into a finished page: a
neutral card with a family-colored accent bar, an `HTTP nnn` chip row, the
suggested fix, context, the cause chain, and a trace footer. It integrates
with [go-error-family](https://pkg.go.dev/github.com/larsartmann/go-error-family)
and promotes samber/oops user-safe messages and trace IDs through the
bridge.

## Quick start

Wrap any handler with `ErrorHandler` — it derives everything from the error
and writes the real HTTP status:

```go
import (
    "github.com/larsartmann/templ-components/errorpage"
)

mux.Handle("GET /dashboard", errorpage.ErrorHandler(
    err,
    errorpage.ErrorHandlerConfig{Nonce: cspNonce()},
))
```

Prefer full control? Build the props yourself and render inside your own
layout shell (the [demo](https://templcomponents.lars.software/demo) uses exactly this pattern for its standalone
`/errors/*` routes):

```templ
props := errorpage.FromError(err) // family, title, code, why/fix, trace, status
props.ShowTimestamp = true
@errorpage.ErrorPage(props)
```

## Families

Every error resolves to one of six families, which drives the color, icon,
default title, and HTTP status:

| Family           | Status | Default title                  |
| ---------------- | ------ | ------------------------------ |
| `Rejection`      | 400    | Request could not be completed |
| `Conflict`       | 409    | Conflict detected              |
| `Transient`      | 503    | Temporary error                |
| `Corruption`     | 500    | Data integrity error           |
| `Infrastructure` | 503    | Service unavailable            |
| `Orchestration`  | 500    | Orchestration failure          |

`FromError` prefers the error's own `ErrorTitle()`; when absent it falls
back to the family default, so a page never renders headingless.

## JSON mode

API consumers get the same facts the page shows — set `JSON: true`:

```go
errorpage.ErrorHandler(err, errorpage.ErrorHandlerConfig{
    JSON: true, // application/json error body
})
```

The body carries `family`, `code`, `message`, `title`, `why`, `fix`,
`trace`, and `context`; untraced errors omit the `trace` key entirely.

## Recovery actions and the code chip

`ErrorPage` carries a typed primary action, an optional secondary ghost
action, a copy-to-clipboard button on the error-code chip, and a card-width
enum:

```go
props := errorpage.FromError(err)
props.WayOutAction = errorpage.WayOutAction{ // typed bundle; Text wins over the loose WayOut/WayOutHref strings
    Text: "Back to dashboard",
    Href: "/dashboard",
}
props.SecondaryWayOut = "Status page"                // ghost link beside the primary action
props.SecondaryWayOutHref = "https://status.example.com"
props.CopyCode = true                                // clipboard button next to the code chip
props.MaxWidth = errorpage.ErrorMaxWidthLG           // lg | xl (default) | 2xl | 4xl
```

An action with text but no href renders as a history-back button. When the
`ErrorHandler` pipeline finds no way out at all and the error implements
`IsRetryable() bool`, it fills in a same-path **Retry** link — an explicit
caller way out always wins.

## Try it live

The demo ships a stateless playground that renders a real `ErrorPage` from
query params — family, status, code, title, message, and card width:

**[Open the error-page playground →](https://templcomponents.lars.software/demo/errors/playground)**

## Components

- `ErrorPage` — full-page display (`<main>` landmark, action pair, width
  enum, copy-to-clipboard code chip).
- `NotFound404` — dedicated 404 with hero numeral, search, and quick links.
- `ErrorDetail` — inline card for dashboards and panels (`Variant: Tinted`
  or `Neutral`).
- `ErrorAlert` — family-aware banner.

All render family-safe colors in both themes and honor
`prefers-reduced-motion`. See the [API reference](/api-reference) for
the complete props model.
