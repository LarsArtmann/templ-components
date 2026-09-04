# Transport Wiring Guide — one `wire.Action`, htmx & Datastar

`utils/wire` is a transport-agnostic wiring contract (ADR-0036): you describe
a client-initiated hypermedia exchange **once**, and the package renders it as
the attribute dialect of the runtime your page loaded. It composes with every
component in this library, because every component spreads `BaseProps.Attrs`
— and `display.Button` additionally accepts a `Wire` field directly.

## The model: two orthogonal axes

```
Component (templ, pure HTML)
  ├─ transport axis:   htmx (default) │ Datastar      ← wire.Action
  └─ element axis:     native element │ custom element ← consumer's choice
```

- **Transport** = which client runtime performs the exchange. Mutually
  exclusive per element, chosen per wiring.
- **Element model** = what hosts the markup. Orthogonal: both runtimes work
  inside native elements and inside light-DOM custom elements alike.
- Shadow DOM is not an option here at all (ADR-0033) — it cannot receive the
  consumer's Tailwind classes.

## Quick start

```go
import "github.com/larsartmann/templ-components/utils/wire"
```

Same props, either runtime — only the Transport differs:

```go
// htmx (Transport zero value resolves to htmx per ADR-0030)
@display.Button(display.ButtonProps{
    Text: "Load",
    Wire: &wire.Action{URL: "/api/fragment", Target: "#out"},
})

// Datastar
@display.Button(display.ButtonProps{
    Text: "Load",
    Wire: &wire.Action{Transport: wire.TransportDatastar, URL: "/api/fragment"},
})
```

Any component, even without a `Wire` field — spread the attributes yourself:

```go
@forms.Input(forms.InputProps{
    Name:  "email",
    Attrs: wire.Action{URL: "/api/validate", Event: wire.EventChange}.Attributes(),
})
```

### Zero values and validation

| Field       | Zero value behavior                                                                       |
| ----------- | ----------------------------------------------------------------------------------------- |
| `Transport` | `""` → htmx (library default)                                                              |
| `Method`    | `""` → GET                                                                                 |
| `Event`     | `""` → htmx: attribute omitted (element defaults: click/submit/change); Datastar: `click` |
| `URL`       | `""` → renders nothing (inert)                                                             |
| unknowns    | `TransportIsValid`/`MethodIsValid`/`EventIsValid` exist; rendering falls back to defaults  |

## Dialect mapping

| `wire.Action`        | htmx rendering                       | Datastar rendering                          |
| -------------------- | ------------------------------------ | ------------------------------------------- |
| `Method` + `URL`     | `hx-get="/api/fragment"`             | `data-on:click="@get('/api/fragment')"`     |
| `Event: EventSubmit` | `hx-trigger="submit"`                | event key: `data-on:submit="…"`             |
| `Target: "#out"`     | `hx-target="#out"`                   | *not rendered* — see below                  |
| `URL: ""`            | nothing                              | nothing                                     |

### Why Target is htmx-only (the #1 FAQ)

The pinned Datastar v1.0.2 runtime (`go-datastar/static` v0.4.0) accepts **no
`target` option** on `@get`/`@post`/… — verified against the embedded bundle
(see `docs/datastar-runtime-facts.md`). Datastar targeting is response-driven:

1. **Non-SSE HTML responses**: the runtime reads `Datastar-Selector` and
   `Datastar-Mode` **response headers** and patches that region.
2. **SSE responses**: the `datastar-patch-elements` event datalines carry the
   selector/mode.
3. **No selector at all**: default `outer` mode matches the fragment root by
   its `id` — give the fragment root the target's id and nothing else is
   needed.

So the handler owns the target under Datastar. `wire` gives you the constants
to do that cleanly:

## One handler, both transports

The library packages the branching for you: wrap your fragment handler in
`wire.Handler` and it becomes a both-transports endpoint.

```go
mux.Handle("/api/fragment", wire.Handler(wire.PatchTarget{
    Selector: "#out",
    Mode:     wire.PatchModeInner, // zero value resolves to inner
}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    noStore(w)
    componentOr500(w, r, fragment(time.Now().Format("15:04:05")))
})))
```

`wire.Handler` inspects the request headers: a Datastar caller (marked by
`Datastar-Request`) receives `Datastar-Selector`/`Datastar-Mode` response
headers targeting the `PatchTarget` — Datastar decides where to patch from the
response. htmx and plain callers pass through untouched: htmx targets
client-side via `hx-target`, and a plain navigation just renders the fragment.
An empty `PatchTarget.Selector` degrades to Datastar's default id-matched
patching. For custom branching, the building blocks are exported too:

```go
if wire.IsDatastar(r) {
    w.Header().Set(wire.HeaderDatastarSelector, "#out")
    w.Header().Set(wire.HeaderDatastarMode, "inner")
}
```

| Constant                 | Direction | Meaning                                             |
| ------------------------ | --------- | --------------------------------------------------- |
| `wire.HeaderHXRequest`   | request   | set by htmx on every AJAX request                    |
| `wire.HeaderDatastarRequest` | request | set by Datastar on every fetch action              |
| `wire.HeaderDatastarSelector` | response | names the patch region for non-SSE HTML responses |
| `wire.HeaderDatastarMode` | response  | merge mode (`inner`, `outer`, …)                    |

The demo implements this end-to-end: `examples/demo/wire_demo.templ` renders
the same Action under both transports, `/api/wire/fragment` serves both, and
`examples/demo/wire_demo_test.go` pins the branching contract.

## Deliberate scope boundaries

`wire` covers only the dialects' common subset. Transport-specific machinery
stays in its module, where it already exists:

| Need                        | Use instead                            |
| --------------------------- | -------------------------------------- |
| Polling / reveal / lazy load | `htmx.PolledRegion`, `navigation.LoadMore` |
| Out-of-band swaps           | `htmx.SwapOOB`                          |
| Confirm dialogs             | `htmx.ConfirmDelete` (`hx-confirm`)     |
| Loading indicators          | `htmx.InlineLoadingOverlay`, `datastar.Indicator` |
| SSE streams / signals       | `datastar.LiveRegion`, `datastar.Get/Post/...` with retry options |
| View transitions            | `htmx.ViewTransitions`                  |

If you need htmx trigger-engine power beyond a plain event (`hx-trigger="click
delay:1s"`, `from:`, `once:`), pass raw attributes via `Attrs` — that is the
designed escape hatch, not a wire gap.

## Web Components: the consumer-side recipe (ADR-0033 stands)

The library ships **no** custom elements — Shadow DOM would break Tailwind
theming, and the zero-JS identity is a feature (ADR-0033). But if *you* want
custom elements in your app, light-DOM custom elements compose with
everything in this library with zero library changes:

```html
<tc-stat-card label="Requests" value="1,204" trend="up"></tc-stat-card>
```

```js
// CSP: serve from an external module or inline with a nonce.
class TcStatCard extends HTMLElement {
  static get observedAttributes() { return ["label", "value", "trend"]; }
  attributeChangedCallback() { this.#render(); }
  connectedCallback() { this.#render(); }
  #render() {
    // Light DOM: children and rendered markup stay in the document —
    // Tailwind classes apply, htmx swaps and Datastar patches inside work.
    this.innerHTML = `
      <div class="rounded-lg border p-4">
        <p class="text-sm text-gray-500">${this.getAttribute("label")}</p>
        <p class="text-2xl font-bold">${this.getAttribute("value")}</p>
      </div>`;
  }
}
customElements.define("tc-stat-card", TcStatCard);
```

Why this works when "Web Components" usually would not:

- **Light DOM, no `attachShadow`** — the element's content participates in the
  document stylesheet, so your compiled Tailwind (including `@theme` token
  overrides) styles it like any other markup.
- **htmx-compatible** — swaps into or out of the element's light-DOM children
  do not desync anything (there is no parallel shadow tree); `wire.Action`
  attributes on elements *inside* the custom element fire normally.
- **Datastar-compatible** — the runtime observes `data-*` attributes in the
  whole document (including inside custom elements) via MutationObserver.

Constraints to respect (same list as ADR-0033, minus the ones light DOM fixes):

- You own the JS: custom elements upgrade via script, so these islands are
  never zero-JS. Serve the definition with a CSP nonce or as an external file.
- `innerHTML` from attribute strings is your responsibility — escape
  consumer-provided strings before interpolation.
- Anything needing styling isolation must not reach for Shadow DOM; use
  scoped class names or container queries (`@container`) instead.

If demand ever justifies a wrapper *module* (ADR-0033's narrow exception), it
would look exactly like this recipe packaged as Go — light-DOM hosts only,
never Shadow DOM — and would need its own superseding ADR.

## Testing notes

- `utils/wire` ships table tests for every enum and dialect path, a render
  test through `templ.RenderAttributes` (proves `data-on:<event>` colon keys
  survive the attribute writer), and a fuzz target (`FuzzAction`) asserting
  no panic and no empty-valued attributes on arbitrary input.
- Attribute values are HTML-entity encoded in source (`'` → `&#39;`) and
  decode back in the DOM — assert on the decoded meaning in DOM-level tests,
  or on the encoded form in string tests (see `TestActionAttributesRender`).
