# Transport Wiring Guide — one `wire.Action`, htmx & Datastar

`utils/wire` is a transport-agnostic wiring contract (ADR-0036, extended by
[ADR-0038](adr/0038-common-subset-extensions.md) with `ContentType` and
`DebounceMS`): you describe
a client-initiated hypermedia exchange **once**, and the package renders it as
the attribute dialect of the runtime your page loaded. It composes with every
component in this library, because every component spreads `BaseProps.Attrs`
— and `display.Button`, `navigation.LoadMore`, and `forms.Form` additionally
accept a `Wire` field directly.

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
| `Transport` | `""` → htmx (library default)                                                             |
| `Method`    | `""` → GET                                                                                |
| `Event`     | `""` → htmx: attribute omitted (element defaults: click/submit/change); Datastar: `click` |
| `URL`       | `""` → renders nothing (inert)                                                            |
| `ContentType` | `""` → Datastar signals as JSON (runtime default); `ContentTypeForm` serializes the enclosing form's fields; htmx ignores it |
| unknowns    | `TransportIsValid`/`MethodIsValid`/`EventIsValid`/`ContentTypeIsValid` exist; rendering falls back to defaults |

## Dialect mapping

| `wire.Action`        | htmx rendering           | Datastar rendering                      |
| -------------------- | ------------------------ | --------------------------------------- |
| `Method` + `URL`     | `hx-get="/api/fragment"` | `data-on:click="@get('/api/fragment')"` |
| `Event: EventSubmit` | `hx-trigger="submit"`    | event key: `data-on:submit="…"`         |
| `Target: "#out"`     | `hx-target="#out"`       | _not rendered_ — see below              |
| `ContentType: ContentTypeForm` | _not rendered_ (native form serialization) | `{contentType: 'form'}` appended — serializes the enclosing form |
| `URL: ""`            | nothing                  | nothing                                 |

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

| Constant                      | Direction | Meaning                                           |
| ----------------------------- | --------- | ------------------------------------------------- |
| `wire.HeaderHXRequest`        | request   | set by htmx on every AJAX request                 |
| `wire.HeaderDatastarRequest`  | request   | set by Datastar on every fetch action             |
| `wire.HeaderDatastarSelector` | response  | names the patch region for non-SSE HTML responses |
| `wire.HeaderDatastarMode`     | response  | merge mode (`inner`, `outer`, …)                  |

The demo implements this end-to-end: `examples/demo/wire_demo.templ` renders
the same Action under both transports, `/api/wire/fragment` serves both, and
`examples/demo/wire_demo_test.go` pins the branching contract.

## Busy, polling, and reveal signaling (research notes, 2026-09-04)

How the two runtimes communicate "work in progress" — verified against the
pinned bundles, so you don't have to:

- **htmx 2.0.10 (the self-hosted bundle here) sets no `aria-busy`.** The
  embedded bundle contains no aria attributes at all. Loading state is CSS:
  the `htmx-request` class on the requesting element plus
  `htmx-indicator`-style visibility rules. The `htmx` module's components
  (`LoadingButton`, `InlineLoadingOverlay` with `role="status"`,
  `LoadingIndicator`) are the accessible surface — use them, they announce.
- **Datastar signals busy through its engine.** `datastar.Indicator` wires a
  fetch-in-progress signal to any spinner, and `datastar.LiveRegion` emits
  `aria-busy` on its live region while streaming (its busy script).
- **Decision: docs, not a helper.** The signaling models (CSS classes vs
  signals) are too different for a common `wire` helper without inventing a
  third abstraction. Each transport module already ships the accessible
  component for its model.

**Interval and intersect triggers are symmetric concepts that stay out of the
contract for now.** The pinned Datastar v1.0.2 bundle includes
`data-on-interval` and `data-on-intersect` (IntersectionObserver is in the
bundle), and htmx has `hx-trigger="every 2s"` / `revealed` — but the trigger
_syntax_ is dialect-specific (durations, options, filters). Extending
`wire.Event` would mean modeling a mini trigger language; that is a future
ADR-sized decision, deliberately not smuggled into the current common subset.

**Form-submit parity, updated 2026-09-07.** `wire.EventSubmit` renders
`hx-trigger="submit"` / `data-on:submit` fine, and whole-**form** submission
is now symmetric: the pinned Datastar v1.0.3 runtime accepts
`{contentType: 'form'}` on fetch actions, which serializes the enclosing
form's fields (with an HTML5 validation gate, the submitter button's
name/value, and enctype-aware bodies — see
`docs/datastar-runtime-facts.md`). `wire.Action.ContentType` models that
option, and `forms.FormProps.Wire` applies it by default. What stays
asymmetric is **per-field value binding**: htmx includes the triggering
element's (or form's) fields natively, while a Datastar action with a static
URL needs bound signals
(`data-bind:value` + an interpolated expression like
`@get('/api/validate?value=' + encodeURIComponent($value || ''))`). The demo's
server-validation block (`/api/wire/validate`) shows both: the typed contract
under htmx, the `Attrs` escape hatch under Datastar.

## Dual-transport forms

`forms.Form` ships a `Wire` field (the TODO #153 survey candidate): the same
form submits under either runtime, with the form's fields — including the
CSRF hidden input — traveling in both dialects.

```go
@forms.Form(forms.FormProps{
    Method: forms.FormPost, // no-JS fallback
    Wire:   &wire.Action{Transport: wire.TransportDatastar, Method: wire.MethodPost, URL: "/api/save"},
}) {
    @forms.Input(forms.InputProps{Name: "email", Type: forms.InputEmail, Label: "Email"})
    @display.Button(display.ButtonProps{Text: "Save", Type: display.ButtonHTMLSubmit})
}
```

| Aspect             | htmx dialect                                                 | Datastar dialect (`ContentTypeForm`, applied by Form)          |
| ------------------ | ------------------------------------------------------------ | -------------------------------------------------------------- |
| Rendering          | `hx-post="/api/save" hx-trigger="submit"` (implicit trigger) | `data-on:submit="@post('/api/save', {contentType: 'form'})"`  |
| Field values       | native serialization (urlencoded; multipart with `enctype`)  | same — FormData → urlencoded (or multipart with `enctype`)      |
| HTML5 validation   | `Validate: true` adds `hx-validate="true"`; `NoValidate: true` renders `novalidate` | automatic (`checkValidity` gate); `NoValidate: true` renders `novalidate`, which skips the gate |
| Submitter button   | name/value included                                          | name/value appended by the runtime                             |
| Response targeting | `Wire.Target` → `hx-target` (default swaps into the form)    | response-driven — wrap the handler in `wire.Handler`            |

The form-level defaults: an unspecified `Event` becomes `submit`, an
unspecified `ContentType` becomes `ContentTypeForm` (set `ContentTypeJSON`
explicitly to submit signals instead — then no form fields travel). `Method`
follows wire's zero-value contract (GET) — set `MethodPost` for mutating
submissions. The demo implements this end-to-end in
`examples/demo/wire_demo.templ` (`/api/wire/form`).

### Server-side validation round-trip

The #1 forms use case — submit → validate → inline errors — is proven
dual-transport in Chromium (`visualtest/wire_form_e2e_test.go`) and
documented as a copy-pasteable recipe:
**[`docs/recipes/server-side-validation.md`](recipes/server-side-validation.md)**.
The essentials: render the form INSIDE its swap region, wrap the handler in
`wire.Handler`, re-render the form with `ValidationSummary` + `Input.Error`
(values preserved) on invalid input, and answer **200 OK even for errors** —
htmx 2.0.10's default `responseHandling` does not swap 4xx/5xx, and the
pinned Datastar bundle dispatches a fetch error event at status >= 400
(both verified against the runtime bundles).

One timing note for test automation: htmx wires swapped-in nodes during its
~20ms settle phase; a client that clicks a freshly swapped submit button
inside that window falls through to a native submit. Humans cannot click
that fast — e2e drivers must wait for the swap to settle first.

### Auto-submit filter components (FilterInput, FilterDropdown)

`forms.FilterInput` (debounced search) and `FilterDropdown.Wire`
(dual-transport select) share one design decision, made here in writing:
**when the component is wired, it wraps its field in a GET form rather than
binding Datastar signals.** The wrapper form wins because:

1. The pinned bundle's `contentType: 'form'` path serializes
   `closest("form")` — without an enclosing form a standalone field cannot
   use form encoding, and the signals default would force servers to parse
   `?datastar=<json>` instead of a clean `?q=<value>`.
2. htmx includes the closest form's values for free, so multi-field filter
   bars serialize all controls under both transports with zero extra
   attributes (the same job `HxInclude` does by hand).
3. The no-JS fallback is a plain GET form — `FilterInput` submits natively
   on Enter; `FilterDropdown` renders a `<noscript>` Apply button.

The alternative (binding the value to a signal and string-building the URL
in the expression) was rejected: expression-concatenated URLs are an
injection surface and cannot be validated as complete literals.

### Busy state on wired submits

Perceived performance: acknowledge the submit instantly. Each runtime has a
native busy mechanism; both compose with `wire` because busy styling is
orthogonal to the wiring attributes.

- **htmx** — wrap `htmx.LoadingButton(defaultText, loadingText, spinner)` in
  the button element that carries the `hx-*` attributes. While the request
  runs, htmx puts `.htmx-request` on that element, hiding the default label
  (`[.htmx-request_&]:hidden`) and revealing the loading label + spinner
  (`.htmx-indicator`). The wired `forms.Form` case is the same: the form
  carries `.htmx-request`, so a `LoadingButton` inside the form reacts
  automatically.
- **Datastar** — set an indicator signal on the trigger and toggle it
  elsewhere: `data-indicator:saving` on the wired `display.Button` (via
  `BaseProps.Attrs` — it composes with `Wire`), plus
  `@datastar.Indicator(datastar.IndicatorProps{Signal: "saving"})` anywhere
  on the page. The Indicator renders `role="status" aria-live="polite"`, so
  screen readers announce the busy state.

Both patterns are live in the demo's "Busy state" card (the endpoint sleeps
800ms so the state is visible): `examples/demo/wire_demo.templ`.

### GET search forms

Search-as-a-GET needs no special API: `forms.Form` with
`Method: forms.FormGet` and a wired submit. Under htmx the fields ride the
query string natively; under Datastar a GET fetch with
`contentType: 'form'` carries the fields as **query parameters too**
(bundle-verified) — identical URLs, one endpoint:

```templ
@forms.Form(forms.FormProps{
	Method: forms.FormGet,
	Action: "/search", // no-JS fallback
	Wire:   &wire.Action{Method: wire.MethodGet, URL: "/api/wire/search", Target: "#results"},
}) {
	@forms.Input(forms.InputProps{Name: "q", Type: forms.InputSearch, Label: "Search"})
	@display.Button(display.ButtonProps{Text: "Search", Type: display.ButtonHTMLSubmit})
}
```

Both runtimes produce `/api/wire/search?q=<value>` — pinned by the demo's
"GET search form" card and its tests (`TestWireSearchEndpoint`,
`TestWireDemoUploadAndSearchCards`). File uploads follow the same
form-encoding path with multipart — see
**[`docs/recipes/file-upload.md`](recipes/file-upload.md)**.

## Deliberate scope boundaries

`wire` covers only the dialects' common subset. Transport-specific machinery
stays in its module, where it already exists:

| Need                         | Use instead                                                       |
| ---------------------------- | ----------------------------------------------------------------- |
| Polling / reveal / lazy load | `htmx.PolledRegion`, `navigation.LoadMore`                        |
| Out-of-band swaps            | `htmx.SwapOOB`                                                    |
| Confirm dialogs              | `htmx.ConfirmDelete` (`hx-confirm`)                               |
| Loading indicators           | `htmx.InlineLoadingOverlay`, `datastar.Indicator`                 |
| SSE streams / signals        | `datastar.LiveRegion`, `datastar.Get/Post/...` with retry options |
| View transitions             | `htmx.ViewTransitions`                                            |

If you need htmx trigger-engine power beyond a plain event (`hx-trigger="click
delay:1s"`, `from:`, `once:`), pass raw attributes via `Attrs` — that is the
designed escape hatch, not a wire gap.

## Web Components: the consumer-side recipe (ADR-0033 stands)

The library ships **no** custom elements — Shadow DOM would break Tailwind
theming, and the zero-JS identity is a feature (ADR-0033). But if _you_ want
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
  attributes on elements _inside_ the custom element fire normally.
- **Datastar-compatible** — the runtime observes `data-*` attributes in the
  whole document (including inside custom elements) via MutationObserver.

Constraints to respect (same list as ADR-0033, minus the ones light DOM fixes):

- You own the JS: custom elements upgrade via script, so these islands are
  never zero-JS. Serve the definition with a CSP nonce or as an external file.
- `innerHTML` from attribute strings is your responsibility — escape
  consumer-provided strings before interpolation.
- Anything needing styling isolation must not reach for Shadow DOM; use
  scoped class names or container queries (`@container`) instead.

If demand ever justifies a wrapper _module_ (ADR-0033's narrow exception), it
would look exactly like this recipe packaged as Go — light-DOM hosts only,
never Shadow DOM — and would need its own superseding ADR.

## Testing notes

- `utils/wire` ships table tests for every enum and dialect path, a render
  test through `templ.RenderAttributes` (proves `data-on:<event>` colon keys
  survive the attribute writer), invariant tests pinning the cross-dialect
  guarantees (empty URL inert, URL present in both dialects, no target under
  Datastar, no script emission), and a fuzz target (`FuzzAction`) asserting
  no panic and no empty-valued attributes on arbitrary input.
- Attribute values are HTML-entity encoded in source (`'` → `&#39;`) and
  decode back in the DOM — assert on the decoded meaning in DOM-level tests,
  or on the encoded form in string tests (see `TestActionAttributesRender`).
- Browser-level proof lives in `visualtest/wire_e2e_test.go`: real Chromium,
  both runtimes, one `wire.Handler` endpoint. Copy it when wiring new
  components.
- Migrating an existing htmx page? `docs/recipes/transport-migration.md` is
  the step-by-step recipe, including the what-does-NOT-migrate table.
