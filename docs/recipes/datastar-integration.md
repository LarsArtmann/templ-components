# Recipe: Using templ-components with Datastar

> How to use this library in a [Datastar](https://data-star.dev/) app, or add
> Datastar-powered real-time components alongside HTMX.

---

## When to choose this recipe

- Your app already uses Datastar (or wants to) for reactive client-side state
- You need **real-time SSE streaming** (live dashboards, activity feeds, chat)
- You want **reactive forms** without server round-trips for every keystroke
- You're migrating from HTMX to Datastar and want to keep your components

**Not sure?** Read `docs/research/datastar-integration-analysis.md` first.

---

## Step 1: Inject the Datastar runtime

Add `datastar.SDKScript` to your layout — alongside or instead of the HTMX
script:

```go
@layout.Base(layout.PageProps{
    HTMXVersion: "", // suppress HTMX if going Datastar-only
}) {
    @datastar.SDKScript(datastar.DefaultSDKScriptProps())
    // ... page content ...
}
```

Self-hosting via [go-datastar/static](https://pkg.go.dev/github.com/larsartmann/go-datastar/static)
(zero-dependency module — its go.mod requires nothing):

```bash
go get github.com/larsartmann/go-datastar/static
```

```go
// Serve the embedded bundle with ETag + Cache-Control
mux.Handle("GET /datastar.js", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
    _, _ = w.Write(static.Bytes())
}))

// Point SDKScript at your self-hosted path
@datastar.SDKScript(datastar.SDKScriptProps{Src: "/datastar.js"})
```

The version pinned by `datastar.DatastarVersion1_0_2` is derived from
`static.Version`, so the CDN URL and the embedded bundle can never drift.

---

## Step 2: Add go-datastar to your project

```bash
go get github.com/larsartmann/go-datastar
```

[go-datastar](https://github.com/LarsArtmann/go-datastar) provides the
server-side SSE protocol vocabulary where every patch is a first-class value.
The library itself does **not** depend on it — only your application does.

---

## Step 3: SSE-powered live region (replaces PolledRegion)

### Frontend (templ)

```templ
@datastar.LiveRegion(datastar.LiveRegionProps{
    URL:       "/stream/metrics",
    AutoStart: true,
    Live:      datastar.LivePolite,
}) {
    @display.StatCard(display.StatCardProps{
        BaseProps: utils.BaseProps{ID: "metrics"},
        Label:     "Active Users",
        Value:     "—",
    })
}
```

### Backend (Go handler)

```go
import (
    "github.com/larsartmann/go-datastar"
    "github.com/larsartmann/go-sse"
)

func streamMetrics(w http.ResponseWriter, r *http.Request) {
    stream := sse.NewStream(w, r)
    defer func() { _ = stream.Close() }()

    resp := datastar.NewResponse(stream)

    for {
        metrics := fetchCurrentMetrics()

        // Patch the StatCard's content. WithSelector is REQUIRED with
        // WithModeInner — the runtime rejects non-default modes without a
        // target. (The default outer mode matches incoming root elements
        // by their id instead, so fragments must carry an id that already
        // exists in the DOM.)
        if err := resp.PatchElementsTempl(statCardContent(metrics),
            datastar.WithSelector("#metrics"),
            datastar.WithModeInner(),
        ); err != nil {
            return
        }

        // Or patch signals only (no HTML round-trip):
        // _ = resp.MarshalAndPatchSignals(map[string]any{"activeUsers": metrics.Users})

        select {
        case <-r.Context().Done():
            return
        case <-time.After(time.Second):
        }
    }
}
```

### One-liner variant and proxy keep-alive

`NewResponseFromHTTP` collapses the two setup lines into one:

```go
func streamMetrics(w http.ResponseWriter, r *http.Request) {
    resp := datastar.NewResponseFromHTTP(w, r)
    // NOTE: this variant does not hand you the *sse.Stream handle. Use the
    // explicit NewResponse form above whenever you need one (heartbeats,
    // Last-Event-ID, onDisconnect hooks).
}
```

Production streams behind a reverse proxy (Nginx, Cloudflare, AWS ALB) should
also run a heartbeat goroutine — proxies kill "idle" connections after 30–60s
of silence, and a metrics stream can legitimately be quiet that long:

```go
stream := sse.NewStream(w, r)
go stream.Heartbeat(stream.Context(), 15*time.Second)
```

The heartbeat is a standard SSE comment frame (`: heartbeat`) — browsers ignore
it, proxies reset their idle timer. Datastar's runtime parses datalines keyed
terminated by a blank line, so comment frames are safely ignored on the wire
(verified against the pinned bundle; see `docs/datastar-runtime-facts.md`).

### Why this is better than PolledRegion

| Aspect       | PolledRegion (HTMX)            | LiveRegion (Datastar SSE)                                                            |
| ------------ | ------------------------------ | ------------------------------------------------------------------------------------ |
| Requests     | N per minute (every interval)  | 1 long-lived connection                                                              |
| Latency      | Up to `interval` seconds stale | Sub-second (push on change)                                                          |
| Idle traffic | Full HTML fragment each poll   | Zero (server pushes only on change)                                                  |
| Bandwidth    | Full HTML every poll           | Signal-only patches possible                                                         |
| Reconnection | Automatic (next poll)          | Automatic (retry w/ backoff; Last-Event-ID resume when the server assigns event ids) |

---

## Step 4: Reactive loading indicators

Datastar's `data-indicator` system replaces HTMX's `hx-indicator`:

```templ
<button
    data-on:click={ datastar.Post("/api/save") }
    data-indicator:saving
>
    Save
</button>

@datastar.Indicator(datastar.IndicatorProps{
    Signal: "saving",
    Spinner: feedback.Spinner(feedback.SpinnerProps{
        Size: feedback.SpinnerSM,
    }),
})
```

The indicator automatically shows/hides based on the in-flight request state —
no custom JS needed.

---

## HTMX to Datastar attribute mapping

Most components work unchanged. Only components with `hx-*` attributes need
Datastar equivalents:

| HTMX                      | Datastar equivalent                                     |
| ------------------------- | ------------------------------------------------------- |
| `hx-get="/url"`           | `data-on:click={ datastar.Get("/url") }`                |
| `hx-post="/url"`          | `data-on:click={ datastar.Post("/url") }`               |
| `hx-delete="/url"`        | `data-on:click={ datastar.Delete("/url") }`             |
| `hx-trigger="every 10s"`  | `data-init={ datastar.Get("/stream") }` (SSE)           |
| `hx-target="#id"`         | Server patches element by ID (no client-side target)    |
| `hx-swap="outerHTML"`     | Server chooses morph mode per patch                     |
| `hx-indicator=".loading"` | `data-indicator:loading` + `data-show="$loading"`       |
| `hx-confirm="Sure?"`      | `data-on:click="if (confirm('Sure?')) @delete('/url')"` |
| `hx-swap-oob`             | Server patches multiple elements in one SSE event       |

### Action helper functions

```go
datastar.Get("/api/search")      // → "@get('/api/search')"
datastar.Post("/api/save")       // → "@post('/api/save')"
datastar.Put("/api/update/1")    // → "@put('/api/update/1')"
datastar.Patch("/api/patch/1")   // → "@patch('/api/update/1')"
datastar.Delete("/api/del/1")    // → "@delete('/api/del/1')"
```

---

## Lifecycle observability (analytics hooks)

The runtime emits a document-level `datastar-fetch` CustomEvent for every
action (`@get`, `@post`, ...) and SSE stream. Hook it for analytics, loading
counters, or telemetry:

```js
document.addEventListener('datastar-fetch', (e) => {
    const { type, el, argsRaw } = e.detail;
    switch (type) {
        case 'started':  analytics.track('stream:start', { url: argsRaw.url }); break;
        case 'finished': analytics.track('stream:end',   { status: argsRaw.status }); break;
        case 'error':    analytics.track('stream:error', { status: argsRaw.status }); break;
        // 'retrying' and 'retries-failed' also exist (reconnect exhaustion)
    }
});
```

Facts (pinned by `datastar.TestPinnedRuntimeBundleContract`, so renames fail CI):

- `type` ∈ `started`, `finished`, `error` (HTTP ≥ 400; status in
  `argsRaw.status`), `retrying`, `retries-failed`
- There is **no `datastar-sse-error` event** — listening for it is dead code
- Reconnection: clean stream EOF reconnects only with `retry: 'always'`
  (`LiveRegionProps.Retry: datastar.RetryAlways`); defaults are 10 retries,
  1s ×2 exponential backoff, 30s cap

---

## Using existing components in a Datastar app

~80% of this library's components are pure server-rendered HTML + Tailwind.
They work identically in HTMX and Datastar apps:

- All `display` components (Card, Badge, Avatar, Table, Grid, etc.)
- All `forms` components (Input, Select, Textarea, Toggle, etc.)
- All `feedback` components (Alert, Toast, Spinner, Skeleton, etc.)
- All `layout` components (Base, AppShell, Container, Split, Stack)
- All `navigation` components (Nav, Breadcrumbs, Pagination, etc.)
- All `errorpage` components
- All `icons`

Only the `htmx` package's 8 components are HTMX-specific. Replace them with
their Datastar equivalents as shown above.

---

## Coexistence: HTMX + Datastar on the same page

HTMX and Datastar use different attribute namespaces (`hx-*` vs `data-*`), so
they can coexist on the same page without conflict:

```go
@layout.Base(layout.PageProps{HTMXVersion: layout.HTMXVersion2_0_10}) {
    // Datastar runtime
    @datastar.SDKScript(datastar.DefaultSDKScriptProps())

    // HTMX-powered form submission
    @forms.Form(forms.FormProps{Action: "/save", Method: "POST"}) {
        @forms.Input(forms.InputProps{Label: "Name"})
    }

    // Datastar-powered live region
    @datastar.LiveRegion(datastar.LiveRegionProps{
        URL:       "/stream/updates",
        AutoStart: true,
    }) {
        @display.Card(display.CardProps{Title: "Live Updates"}) {
            <p>Waiting for updates...</p>
        }
    }
}
```

---

## CSP considerations

Datastar loads as `<script type="module" src="...">`. Under a strict CSP:

```
script-src 'self' https://cdn.jsdelivr.net 'unsafe-eval';
```

**`'unsafe-eval'` is required**: the Datastar runtime compiles every `data-*`
expression with the `Function()` constructor, and CSP classifies that as
`eval` — without the directive all expression attributes silently fail.
`'unsafe-inline'` is still NOT needed: the runtime script is external, and
every inline script this library emits carries a nonce.

If `'unsafe-eval'` is unacceptable for your threat model, self-host the
bundle and avoid expression attributes — or use the HTMX path, which needs
no eval.

---

## Further reading

- `docs/datastar-runtime-facts.md` — audited runtime facts (SSE wire format,
  lifecycle events, CSP, reconnect matrix); the source of truth for everything
  summarized above
- [Datastar official docs](https://data-star.dev/)
- [go-datastar](https://github.com/LarsArtmann/go-datastar) — Go protocol library (patches as values)
- [go-datastar/static](https://pkg.go.dev/github.com/larsartmann/go-datastar/static) — embedded JS bundle (zero deps)
- `docs/research/datastar-integration-analysis.md` — full deep-research analysis
- `docs/adr/0030-datastar-integration-strategy.md` — architectural decision
- `docs/javascript-guide.md` — JS decision ladder (Datastar is rung 7)
