# Recipe: Migrating an HTMX-Only Page to Dual Transport

You have a page wired with htmx attributes (`hx-get`, `hx-target`, …) and you
want the same page — or a section of it — to work under Datastar without
rewriting every component call. This is exactly what the `utils/wire`
contract (ADR-0036) is for. The whole migration is three moves.

## Before: hand-written htmx attributes

```templ
@display.Button(display.ButtonProps{
    BaseProps: utils.BaseProps{Attrs: templ.Attributes{
        "hx-get":    "/api/items",
        "hx-target": "#items",
    }},
    Text: "Load more",
})

<div id="items"></div>
```

Two problems: the wiring lives in untyped string maps, and switching the page
to Datastar means touching every component call.

## Move 1: Replace attribute literals with a `wire.Action`

```templ
@display.Button(display.ButtonProps{
    BaseProps: utils.BaseProps{Attrs: wire.Action{
        URL:    "/api/items",
        Target: "#items",
    }.Attributes()},
    Text: "Load more",
})
```

Rendered output is identical (`hx-get` + `hx-target`; the zero-value Transport
is htmx per ADR-0030). This is a pure refactor — your tests should not change.

Components that opted in take the action directly, which reads even better:

```templ
@display.Button(display.ButtonProps{
    Wire: &wire.Action{URL: "/api/items", Target: "#items"},
    Text: "Load more",
})
```

## Move 2: Make the endpoint transport-agnostic with `wire.Handler`

Wrap your fragment handler once. Datastar callers get response-header
targeting; htmx and plain browser callers pass through untouched.

```go
mux.Handle("/api/items", wire.Handler(wire.PatchTarget{
    Selector: "#items",      // where Datastar patches; htmx ignores this
    Mode:     wire.PatchModeInner,
}, http.HandlerFunc(itemsFragment)))
```

Under htmx the region is still chosen client-side by `hx-target`. Under
Datastar the response decides — the pinned v1.0.2 runtime has no client-side
target option. If you relied on self-replacement (`hx-target="this"`), give
the element an explicit `ID` instead: a Datastar fragment whose root carries
the same id patches it automatically (id matching).

## Move 3: Flip the transport per page (or per section)

```go
// one field, whole section switches
transport := wire.TransportDatastar
action := wire.Action{Transport: transport, URL: "/api/items", Target: "#items"}
```

In practice you drive this from a query param, a feature flag, or the page
model — the demo does it with `?transport=htmx|datastar` and a segmented
control (`examples/demo/wire_demo.templ`).

## What does NOT migrate (by design)

| htmx feature                         | Datastar equivalent                          | Verdict |
| ------------------------------------ | -------------------------------------------- | ------- |
| `hx-trigger="every 2s"` polling      | `data-on-interval`                           | Stay in modules: `htmx.PolledRegion` / datastar actions. See the [signaling notes](../transport-wiring.md). |
| `hx-trigger="revealed"`              | `data-on-intersect`                          | Same — `navigation.LoadMore`'s `InfiniteScroll` stays htmx-only under `Wire`. |
| `hx-swap-oob` out-of-band swaps      | SSE patch events with selectors              | `htmx.SwapOOB` / `datastar.LiveRegion`. |
| `hx-confirm`                         | none in the pinned runtime                   | `htmx.ConfirmDelete` stays htmx-only. |
| Carrying field values on `change`    | bound signals + interpolated expressions     | htmx: typed contract works. Datastar: `Attrs` escape hatch (demo's `/api/wire/validate` block shows both). |

The rule (ADR-0036): the wire contract covers only what both dialects express
the same way. Everything else keeps living in its transport module, where the
accessible components already exist.

## Verifying the migration

1. String level: golden tests assert the dialect attributes; remember templ
   entity-encodes attribute values — assert `data-on:click="@get(&#39;/api/…&#39;)"`
   (the encoded form) in Go string tests.
2. Browser level: `visualtest/wire_e2e_test.go` is the executable template —
   it clicks the same `Action` under both runtimes and asserts the fragment
   lands in the right region.
