---
title: Transport Wiring
description: Wire components once and switch between HTMX and Datastar with one field — the transport-agnostic wiring contract.
---

## One Action, Two Dialects

HTMX and Datastar are dialects of the same operation: a method on a URL, triggered by a DOM event, patching a target region. The `utils/wire` package encodes that operation as one typed struct and renders it in whichever dialect you configure:

```go
import "github.com/larsartmann/templ-components/utils/wire"

// htmx (the default — zero value resolves to htmx per ADR-0030)
wire.Action{URL: "/api/items", Target: "#items"}
// → hx-get="/api/items" hx-target="#items"

// datastar — one field switches the transport
wire.Action{Transport: wire.TransportDatastar, URL: "/api/items"}
// → data-on:click="@get('/api/items')"
```

## Using It With Components

Every component accepts the rendered attributes through `BaseProps.Attrs`:

```templ
@display.Button(display.ButtonProps{
    BaseProps: utils.BaseProps{Attrs: action.Attributes()},
    Text:      "Load more",
})
```

Components that opt in take the action directly — `display.Button` has a typed `Wire` field:

```templ
@display.Button(display.ButtonProps{
    Text: "Load more",
    Wire: &wire.Action{URL: "/api/items", Target: "#items"},
})
```

## One Endpoint, Both Transports

The two runtimes mark their requests differently, and they pick the patch region differently:

| Caller  | Request marker   | Region chosen by                    |
| ------- | ---------------- | ----------------------------------- |
| htmx    | `HX-Request`     | `hx-target` (client-side)           |
| Datastar| `Datastar-Request` | response headers (`Datastar-Selector`, `Datastar-Mode`) |

`wire.Handler` wraps your fragment handler so one endpoint serves both — Datastar callers get the response-header targeting, everyone else passes through:

```go
mux.Handle("/api/items", wire.Handler(wire.PatchTarget{
    Selector: "#items",
    Mode:     wire.PatchModeInner,
}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    // render the fragment component...
})))
```

`wire.IsDatastar(r)` / `wire.IsHTMX(r)` are available for custom branching.

## Scope Boundaries

The contract covers the dialects' **common subset** only (ADR-0036). Transport-specific machinery stays in its module: polling (`htmx.PolledRegion`), out-of-band swaps (`htmx.SwapOOB`), confirm dialogs (`htmx.ConfirmDelete`), loading indicators (`htmx.InlineLoadingOverlay`, `datastar.Indicator`), and SSE streams (`datastar.LiveRegion`).

One deliberate asymmetry: `Action.Target` renders only for htmx. The pinned Datastar v1.0.2 runtime accepts no target option on fetch actions — targeting is response-driven, which is exactly what `wire.Handler` does for you.

## Web Components

The library ships no custom elements (ADR-0033): Shadow DOM breaks Tailwind theming. Light-DOM custom elements defined by consumers compose perfectly with everything here — both runtimes traverse light DOM unchanged. See [the consumer recipe](https://github.com/larsartmann/templ-components/blob/master/docs/transport-wiring.md#web-components-the-consumer-side-recipe-adr-0033-stands).

## Further Reading

- [Transport wiring guide (full)](https://github.com/larsartmann/templ-components/blob/master/docs/transport-wiring.md)
- [ADR-0036: the wiring contract decision](https://github.com/larsartmann/templ-components/blob/master/docs/adr/0036-transport-wiring-contract.md)
- [Datastar runtime facts](https://github.com/larsartmann/templ-components/blob/master/docs/datastar-runtime-facts.md)
