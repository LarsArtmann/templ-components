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
        Label: "Active Users",
        Value: "—",
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

        // Patch the StatCard's value by targeting its element ID
        if err := resp.PatchElementsTempl(statCard(metrics),
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

### Why this is better than PolledRegion

| Aspect       | PolledRegion (HTMX)            | LiveRegion (Datastar SSE)           |
| ------------ | ------------------------------ | ----------------------------------- |
| Requests     | N per minute (every interval)  | 1 long-lived connection             |
| Latency      | Up to `interval` seconds stale | Sub-second (push on change)         |
| Idle traffic | Full HTML fragment each poll   | Zero (server pushes only on change) |
| Bandwidth    | Full HTML every poll           | Signal-only patches possible        |
| Reconnection | Automatic (next poll)          | Automatic (Last-Event-ID resume)    |

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
script-src 'self' https://cdn.jsdelivr.net;
```

Or self-host and use `script-src 'self'` only. No `'unsafe-inline'` or
`'unsafe-eval'` needed — Datastar's `data-*` expressions are sandboxed via
`Function()` constructors, not `eval()`.

---

## Further reading

- [Datastar official docs](https://data-star.dev/)
- [go-datastar](https://github.com/LarsArtmann/go-datastar) — Go protocol library (patches as values)
- [go-datastar/static](https://pkg.go.dev/github.com/larsartmann/go-datastar/static) — embedded JS bundle (zero deps)
- `docs/research/datastar-integration-analysis.md` — full deep-research analysis
- `docs/adr/0030-datastar-integration-strategy.md` — architectural decision
- `docs/javascript-guide.md` — JS decision ladder (Datastar is rung 7)
