# Recipe: Family-Aware Error Pages

**Audience:** Consumers wiring styled error pages for failed actions, broken
routes, and degraded dependencies — especially admin/dashboard surfaces that
must keep the HTTP status honest and never leak internals.

**Problem:** The `errorpage` package ships the components (`ErrorPage`,
`NotFound404`, `ErrorAlert`, `ErrorDetail`) and handler helpers
(`WriteError`, `ErrorHandler`, `FromError`), but not the _policy_ around
them: how to map your domain errors to visual families, where per-code copy
lives, how to branch HTMX swaps vs navigations, and what to do when your
normal layout cannot render on the error path.

**Outcome:** You have the two patterns proven in production consumers
(cqrs-htmx's `adminui` and `dashboardui`), copy-pasteable, with the
decision rules that keep error UX honest.

---

## Tier 1: One-call escape hatch — `WriteError` / `ErrorHandler`

When an error IS the response (a route registered specifically to handle a
failure), skip composition entirely:

```go
// Any handler: one call writes a family-aware, status-correct page.
// FromError(err) prefers err's Public() user-safe message and derives
// family/title/code/why/fix/trace from go-error-family semantics.
errorpage.WriteError(w, r, err, nonce)
```

```go
// A dedicated route for a known failure mode:
mux.Handle("/payments/declined", errorpage.ErrorHandler(err, errorpage.ErrorHandlerConfig{
    Nonce:     nonce,
    HTMLShell: true, // standalone response, not embedded in a layout
}))
```

Use these when you have an `error` value. Reach for Tier 2 when you have a
_status code and a decision to make_ instead.

## Tier 2: Map statuses to families, codes to copy

`Family` selects the visual treatment (accent bar, tone, icon). Map HTTP
statuses once, in one function:

```go
func errorFamilyFor(status int) errorpage.Family {
    switch status {
    case http.StatusConflict:
        return errorpage.FamilyConflict
    case http.StatusServiceUnavailable:
        return errorpage.FamilyTransient
    case http.StatusInternalServerError:
        return errorpage.FamilyCorruption
    }
    if status >= 500 {
        return errorpage.FamilyInfrastructure
    }
    return errorpage.FamilyRejection
}
```

Per-code copy beats per-family copy: a table of known domain codes →
user-safe guidance, with a family-level fallback for everything else. This
is the `adminui` shape (provenance: cqrs-htmx `adminui/errorpage.go`):

```go
var actionErrorMessages = map[string]string{
    "tenant.already_exists": "A tenant with this ID already exists. Pick a different identifier.",
    "user_not_found":        "That user no longer exists. Refresh the list and try again.",
    // ...
}

func actionErrorFallback(f errorfamily.Family) string {
    switch f {
    case errorfamily.Conflict:
        return "The record changed since you loaded it. Refresh and try again."
    case errorfamily.Transient:
        return "The service is temporarily unavailable. Try again in a moment."
    // ...
    default:
        return "The request failed. Try again."
    }
}
```

Then one writer serves every failure path. The two rules it bakes in:

1. **Branch HTMX vs navigation.** An HTMX swap must receive the bare card
   (it lands inside an existing document); a navigation gets the card inside
   your layout shell.
2. **Keep the status and `no-store` honest.** Error pages are never
   cacheable, and the status code is what operators and monitors see.

```go
func writeErrorPage(w http.ResponseWriter, r *http.Request, status int, title, message string) {
    props := errorpage.ErrorPageProps{
        Family:     errorFamilyFor(status),
        StatusCode: status,
        Title:      title,
        Message:    message,
        WayOut:     "Back to dashboard",
        WayOutHref: "/",
    }
    var b strings.Builder
    if err := errorpage.ErrorPage(props).Render(r.Context(), &b); err != nil {
        http.Error(w, title, status) // library render failed; stay plain
        return
    }
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Header().Set("Cache-Control", "no-store")
    w.WriteHeader(status)
    if isHTMX(r) { // e.g. wire.IsHTMX(r) or your framework's HX-Request check
        _, _ = w.Write([]byte(b.String()))
        return
    }
    // navigation: render the card inside your normal layout
    _ = layoutShell(page(title), templ.Raw(b.String())).Render(r.Context(), w)
}
```

For failed _actions_ (form submits, mutations), pair the page with immediate
feedback: fire a toast via `HX-Trigger` AND render the error page, so the
user sees why it failed in both the interaction loop and the navigation.

## Tier 3: The noindex error-shell — when your layout can't render

Some layouts need request-scoped data (session, pageData, nav state) that an
error path may not have — or must not trust (the failure might BE that
data). The `dashboardui` pattern (provenance: cqrs-htmx) renders the library
card inside a **minimal document** that loads only the stylesheets:

```templ
templ errorShell(title, basePath, nonce string, inner templ.Component) {
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <meta charset="utf-8"/>
            <meta name="viewport" content="width=device-width, initial-scale=1"/>
            <title>{ title }</title>
            <meta name="robots" content="noindex"/>
            @layout.ThemeScript(nonce)
            <link rel="stylesheet" href={ basePath + "/-/app.css" }/>
        </head>
        <body>
            <div class="error-shell">
                @inner
            </div>
        </body>
    </html>
}
```

A one-line CSS rule (`.error-shell { min-height: 100vh; display: grid; place-items: center; padding: 1rem; }`)
keeps the card centered without duplicating the real layout. The shell is
deliberately static: no nav, no session reads, `noindex` so search engines
never learn your error URLs.

Choose Tier 3 over "render the normal layout anyway" whenever the error
originated in the layout's own dependencies.

## The way out: primary action, ghost action, code chip, width

The card's recovery affordances are four props on `ErrorPageProps`:

- **`WayOutAction` (typed bundle)** — `Text` + `Href`. When `Text` is set it
  wins entirely over the loose `WayOut`/`WayOutHref` strings (no merging).
  Text with an empty href renders as a history-back button.
- **`SecondaryWayOut` / `SecondaryWayOutHref`** — a ghost-styled second
  action beside the primary (e.g. a status page or docs link). Same
  empty-href = history-back rule.
- **`CopyCode`** — renders a clipboard button next to the error-code chip
  (only when a `Code` is set; without a code there is no chip to copy).
  The button reuses the library's `data-tc-copy` singleton script and needs
  a `Nonce` like every inline script.
- **`MaxWidth`** — closed-set enum (`ErrorMaxWidthLG`/`XL`/`2XL`/`4XL`);
  empty or unknown falls back to XL. Useful when the error shell is narrow
  (mobile-first Tier 3 shells read well at `LG`).

```go
props := errorpage.ErrorPageProps{
    // ...family/status/title as in Tier 2...
    WayOutAction: errorpage.WayOutAction{Text: "Back to dashboard", Href: "/dashboard"},
    SecondaryWayOut:     "Status page",
    SecondaryWayOutHref: "https://status.example.com",
    CopyCode:            true,
    MaxWidth:            errorpage.ErrorMaxWidthLG,
}
```

One behavioral default to know: in the `ErrorHandler` pipeline, if you set
NONE of the way-out fields and the error implements `IsRetryable() bool`,
the handler fills in a same-path "Retry" link (`applyRetrySuggestion`).
Explicit caller way outs always win — there is no merge to fight.

## Rules that apply at every tier

- **Never render raw `err.Error()` to users.** It leaks stream IDs, storage
  internals, and provider text that helps no one. Log the full error
  server-side (with its code), render the table/fallback copy.
- **404s are their own component.** Use `errorpage.NotFound404` with
  `DefaultNotFound404Props()` — see [Custom 404 Page](custom-404-page.md).
- **Inline HTMX failures are a different surface.** A failed swap that
  should recover in place (toast + keep the page alive) belongs to the
  [HTMX Error Feedback](server-rendered-htmx-error-feedback.md) pipeline;
  full error pages are for failures that END the interaction.
- **Sanitization is `FromError`'s job when you have an error value** — it
  prefers the error's `Public()` message. Hand-rolled paths (Tier 2/3) must
  do their own table lookup.

## Where both patterns live

The provenance implementations are cqrs-htmx's `adminui` (status→family
mapping, per-code copy table, toast+page action errors, layout-wrapped
navigations) and `dashboardui` (noindex error-shell, nil-safe
`renderError`, plain-text fallback). They legitimately share no code — the
needs differ — which is why this recipe, not a helper, is the deliverable.
