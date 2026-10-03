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

The extreme case is `display.KanbanBoard`: an entire drag-and-drop board wired through one `Wire` action, with optimistic moves (pending register + failure revert) working identically under both runtimes. See [the full guide](https://github.com/larsartmann/templ-components/blob/master/docs/transport-wiring.md#dual-transport-kanban-board).

## One Endpoint, Both Transports

The two runtimes mark their requests differently, and they pick the patch region differently:

| Caller   | Request marker     | Region chosen by                                        |
| -------- | ------------------ | ------------------------------------------------------- |
| htmx     | `HX-Request`       | `hx-target` (client-side)                               |
| Datastar | `Datastar-Request` | response headers (`Datastar-Selector`, `Datastar-Mode`) |

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

One deliberate asymmetry: `Action.Target` renders only for htmx. Datastar targeting is **response-driven only** (`wire.Handler` echoes `Datastar-Selector` back for you): the pinned runtime ignores client fetch options for patch targeting. `Action.Selector` renders `{selector: '#region'}` under Datastar form-encoding only — it picks which `<form>` serializes; it never targets patches.

The region-merge style is htmx client-side via `Action.Swap` (a `wire.PatchMode` rendering `hx-swap`, `innerHTML`/`outerHTML`/…). Datastar's merge mode is owned by the response: `wire.Handler(wire.PatchTarget{Selector: "#out", Mode: wire.PatchModeOuter}, …)` stamps `Datastar-Mode` for you (zero value = `inner`). Without routing headers the runtime falls back to id-matched patching — a fragment root whose `id` equals the caller element's id. All Datastar fetch options (selector, contentType, retry, …) travel in ONE object literal; the runtime reads exactly one options argument.

Beyond DOM events, the typed trigger language (ADR-0043) renders in both dialects: `Interval` polls (`"10s"`, `"2m"` — minutes/hours are normalized to seconds because both runtimes misparse them as milliseconds), `Reveal` lazy-loads on viewport entry (`Reveal.Exit` is Datastar-only), `ThrottleMS` rate-limits, and `PreventDefault` keeps wired anchors from both patching and navigating (the Datastar runtime auto-prevents only form+submit). Two htmx features have NO Datastar twin: confirm dialogs (`htmx.ConfirmDelete` — the pinned runtime ships no confirm machinery) and global view transitions (`htmx.ViewTransitions`); per-patch view transitions under Datastar ride `PatchTarget.UseViewTransitions`.

## Web Components

The library ships no custom elements (ADR-0033): Shadow DOM breaks Tailwind theming. Light-DOM custom elements defined by consumers compose perfectly with everything here — both runtimes traverse light DOM unchanged. See [the consumer recipe](https://github.com/larsartmann/templ-components/blob/master/docs/transport-wiring.md#web-components-the-consumer-side-recipe-adr-0033-stands).

## Further Reading

- [Transport wiring guide (full)](https://github.com/larsartmann/templ-components/blob/master/docs/transport-wiring.md)
- [ADR-0036: the wiring contract decision](https://github.com/larsartmann/templ-components/blob/master/docs/adr/0036-transport-wiring-contract.md)
- [Datastar runtime facts](https://github.com/larsartmann/templ-components/blob/master/docs/datastar-runtime-facts.md)
