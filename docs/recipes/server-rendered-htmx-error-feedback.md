# Recipe: Server-Rendered HTMX Error Feedback Loop

**Audience:** Consumers building server-rendered (non-SPA) HTMX apps who want
accessible, family-aware error feedback — toast on failure, alert on partial
swap, full error page on navigation — without hand-rolling a switch statement.

**Trigger:** your handler returns `http.Error(w, "query error", 500)` and the
user sees a blank white page. You want a toast or inline alert instead, but
the integration path from "handler error → visible UI" is unclear.

**Outcome:** one error type (`go-error-family` Classified) flows through three
render modes (full page, partial swap, toast), all family-colored and accessible.

---

## The three render modes

| Mode                  | When to use                                   | Component                                     |
| --------------------- | --------------------------------------------- | --------------------------------------------- |
| Full error page       | Navigation request (full page load) fails     | `errorpage.ErrorPage` / `ErrorHandler`        |
| Inline alert fragment | HTMX partial swap returns an error fragment   | `errorpage.ErrorAlert` / `feedback.Alert`     |
| Toast (auto-dismiss)  | Background HTMX request fails (no navigation) | `htmx.GlobalErrorHandling` + `feedback.Toast` |

The **same error** can produce any of these — the choice depends on whether the
request is a navigation, a partial swap, or a background poll.

---

## Setup: install the error-feedback plumbing

### 1. Define your errors with `go-error-family`

```go
import "github.com/larsartmann/go-error-family"

var ErrUserNotFound = errorfamily.New(
    errorfamily.WithCode("USER_NOT_FOUND"),
    errorfamily.WithFamily(errorfamily.FamilyRejection),
    errorfamily.WithWhy("No user with that ID exists."),
    errorfamily.WithFix("Check the ID and try again, or create a new user."),
)
```

The `Family` (Rejection / Transient / Infrastructure / Corruption) drives the
color, icon, and HTTP status code automatically.

### 2. Include the global HTMX error handler once in your base layout

This wires HTMX's `htmx:responseError` event to render a family-aware toast
in the `feedback.ToastContainer`. Include it once near the end of `<body>`:

```templ
@layout.Base(props) {
    <main>{ children... }</main>
    @feedback.ToastContainer("")
    @htmx.GlobalErrorHandling(htmx.DefaultErrorHandlingConfig())
}
```

The handler:

- Listens for `htmx:responseError` events.
- Parses the JSON error body (produced by the server's `ErrorHandlerConfig{JSON: true}`).
- Maps the `family` field to a toast type (Rejection → warning, Transient → info,
  Corruption/Infrastructure → error).
- Announces the error via `aria-live="polite"` for screen readers.

### 3. Wire the server-side error handler

For HTMX-targeted endpoints, return JSON errors so the client-side handler can
read the family:

```go
func handleSave(w http.ResponseWriter, r *http.Request) {
    if err := svc.Save(r.Context(), r.FormValue("id")); err != nil {
        // Returns JSON: {"family":"rejection","code":"USER_NOT_FOUND",...}
        errorpage.ErrorHandler(err, errorpage.ErrorHandlerConfig{
            JSON:  true,
            Nonce: nonce,
        }).ServeHTTP(w, r)
        return
    }
    // success: render the updated fragment
    w.Header().Set("Content-Type", "text/html")
    _ = myFragment().Render(r.Context(), w)
}
```

For non-HTMX navigation endpoints, return an HTML error page:

```go
func handlePage(w http.ResponseWriter, r *http.Request) {
    if err := svc.Load(r.Context(), r.URL.Path); err != nil {
        errorpage.ErrorHandler(err, errorpage.ErrorHandlerConfig{
            HTMLShell: true,
            Nonce:     nonce,
        }).ServeHTTP(w, r)
        return
    }
    // success: render the page
}
```

---

## Flow 1: Background poll fails → toast

User scenario: a dashboard polls `/health` every 5 seconds. The server
returns 503 (Infrastructure). The user sees a red toast that auto-dismisses.

```
[Browser]                    [Server]
    |                            |
    | HTMX GET /health           |
    |--------------------------->|
    |                            | svc.Health() → ErrDBDown
    |                            |   family=infrastructure
    |                            |
    |<---------------------------|
    | 503                        |
    | JSON {family:..., code:...}|
    |                            |
    | htmx:responseError fires   |
    | GlobalErrorHandling parses |
    | family→toast type=error    |
    | Toast rendered in container|
    | aria-live announces        |
    | auto-dismiss after 5s      |
```

No additional code needed — the `GlobalErrorHandling` script handles it.

## Flow 2: Partial swap fails → inline alert

User scenario: user clicks "Save" on a form. HTMX POSTs to `/save`. The server
returns a 409 (Conflict). Instead of a toast, you want the error visible inline
above the form.

Render the error as an HTML fragment (NOT JSON) and swap it into a target div:

```go
func handleSave(w http.ResponseWriter, r *http.Request) {
    if err := svc.Save(r.Context(), r.FormValue("id")); err != nil {
        // Render the error as an inline alert fragment
        w.WriteHeader(http.StatusConflict)
        w.Header().Set("Content-Type", "text/html")
        props := errorpage.FromError(err)
        _ = errorpage.ErrorAlert(props).Render(r.Context(), w)
        return
    }
    // success: render the updated view
}
```

On the client:

```html
<form hx-post="/save" hx-target="#result" hx-swap="innerHTML">
  <!-- form fields -->
</form>
<div id="result"></div>
```

HTMX swaps the returned alert into `#result`. The `ErrorAlert` component is
family-colored, has `role="alert"`, and includes the Why/Fix text from
`go-error-family`.

## Flow 3: Navigation fails → full error page

User scenario: user clicks a link to `/users/123`. The handler returns 404.
The browser navigates to a full error page.

```go
func handleUserPage(w http.ResponseWriter, r *http.Request) {
    user, err := svc.FindUser(r.URL.Query().Get("id"))
    if err != nil {
        errorpage.ErrorHandler(err, errorpage.ErrorHandlerConfig{
            HTMLShell: true,
            Nonce:     nonce,
        }).ServeHTTP(w, r)
        return
    }
    // render the user page
}
```

The `ErrorPage` component renders a full HTML document with the family color,
icon, title, Why/Fix card, and an optional cause chain. `HTMLShell: true` wraps
it in a valid `<!DOCTYPE html>` document so it works as a standalone response.

---

## Why not just `http.Error`?

`http.Error(w, "query error", 500)` produces:

- A blank white page with plain text.
- No family awareness (color, icon, retry hint).
- No accessibility (no `role="alert"`, no `aria-live`).
- No structured JSON for HTMX clients.
- No cause chain or fix suggestion.

The `errorpage.ErrorHandler` family replaces this with one line that adapts to
the request type (API/HTMX vs browser) and error family.

---

## Common pitfalls

- **Forgot `ToastContainer`** → the GlobalErrorHandling script logs to console
  but renders nothing. Always include `@feedback.ToastContainer("")` in the base layout.
- **Mixed JSON and HTML** → an HTMX-targeted endpoint that returns HTML instead
  of JSON confuses the toast handler. Use `JSON: true` for background/poll
  endpoints, HTML for partial-swap endpoints.
- **No nonce on GlobalErrorHandling** → the script is blocked by CSP. Pass
  `ErrorHandlingConfig{Nonce: nonce}`.
- **Forgot `aria-live`** → `GlobalErrorHandling` includes its own
  `tc-error-announcer` div with `aria-live="polite"`, but if you remove it,
  screen-reader users won't hear the error.

See `docs/feedback/2026-07-05_browser-history.md` for the consumer report that
motivated this recipe, and `errorpage/handler.go` for the full handler API.
