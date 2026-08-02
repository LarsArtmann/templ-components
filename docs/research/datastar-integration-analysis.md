# Datastar Integration Analysis — How templ-components Could Benefit

> **Deep research into [Datastar](https://data-star.dev/) (v1.0.2) and its fit
> with this library's architecture, philosophy, and component surface.**
>
> **Date:** 2026-08-02
> **Status:** Recommendation — adopt as opt-in complement, not HTMX replacement

---

## Executive Summary

Datastar is a ~12 KiB, zero-dependency frontend framework that unifies HTMX's
backend reactivity with Alpine.js's frontend reactivity into a single library,
built on Server-Sent Events (SSE) and reactive signals. It is the natural
evolution of the hypermedia approach this library already champions.

This analysis identifies **four tiers of opportunity**, ranked by value and
risk:

| Tier | Opportunity                                              | Risk     | Verdict      |
| ---- | -------------------------------------------------------- | -------- | ------------ |
| A    | Consumer coexistence docs + recipe                       | None     | **Do now**   |
| B    | Parallel opt-in `datastar` package (zero new deps)       | Low      | **Do now**   |
| C    | New SSE/signal-powered capabilities (LiveRegion, etc.)   | Medium   | **Prototype** |
| D    | Full HTMX → Datastar migration                           | **High** | **Reject**   |

**Core recommendation:** Build a `datastar` package that mirrors the existing
`htmx` package — emitting `data-*` attributes and injecting the runtime with
**zero new dependencies** at the library level. Add SSE-powered components
(`LiveRegion`) that unlock capabilities HTMX cannot provide. Keep HTMX as the
default; let consumers opt into Datastar per-page or app-wide. Document the
coexistence path so the 104 existing components work unchanged in Datastar apps.

---

## 1. What Datastar Adds That HTMX Does Not

### 1.1. Reactive Signals (Client-Side State)

HTMX has **no client-side state**. Every interaction that needs state (dropdown
open/close, tab selection, combobox filter text, tag list) requires either a
server round-trip or a separate library (Alpine.js). This library currently
solves this with **15 singleton-guard inline scripts** — hand-written event
delegation that is correct but verbose and重复 per component.

Datastar provides **reactive signals** natively:

```html
<!-- Datastar: one attribute, reactive, no JS -->
<div data-signals="{ open: false }">
  <button data-on:click="$open = !$open">Toggle</button>
  <div data-show="$open">Content</div>
</div>
```

Two-way binding on inputs, computed signals (`data-computed`), class toggling
(`data-class`), and text binding (`data-text`) are all declarative.

### 1.2. SSE-First Streaming (Real-Time Without Polling)

HTMX's real-time story is **polling** (`hx-trigger="every 5s"`). The
`PolledRegion` component uses exactly this. Polling has three costs:

1. **Server load** — N clients × (60 / interval) requests per minute, even when
   nothing changed.
2. **Latency** — updates arrive at most every `interval` seconds. A 10s poll
   means up to 10s stale data.
3. **Wasted bandwidth** — each poll re-sends the full HTML fragment even if the
   data didn't change.

Datastar uses **Server-Sent Events** — a single long-lived connection per client
where the server pushes patches only when data changes:

```go
sse := datastar.NewSSE(w, r)
for {
    sse.PatchElementTempl(liveStats(currentMetrics()))
    time.Sleep(1 * time.Second)
}
```

Benefits: zero idle requests, sub-second latency, bandwidth proportional to
actual changes. Datastar is proven at scale (1 billion checkboxes demo, 800k
DOM updates/second).

### 1.3. DOM Morphing (State Preservation)

HTMX swaps DOM via `innerHTML`/`outerHTML` by default — the old element is
destroyed and replaced. Focus, scroll position, input state, and transitions are
lost unless the morph extension is added.

Datastar **morphs by default** (idiomorph-style, matching by element ID). Only
changed DOM nodes are updated; everything else is preserved. Focus stays where
it was, scroll is maintained, CSS transitions complete naturally.

### 1.4. Server-Driven Updates (Server Knows Best)

HTMX scatters update logic across the trigger element — the button must declare
`hx-get`, `hx-target`, `hx-swap`, `hx-select`. This couples the frontend to the
server's response shape.

Datastar keeps update logic **server-side** — the server decides what to patch
and sends targeted element/signal patches. The frontend trigger is a single
attribute: `data-on:click="@get('/endpoint')"`. The response (SSE stream) tells
the client exactly what changed.

### 1.5. Built-in Request Lifecycle

HTMX requires this library's `GlobalErrorHandling` component (~130 lines of
inline JS) for retry, error-to-toast mapping, and error history.

Datastar has this **built in**: automatic request cancellation (deduplicates
in-flight requests), retry with exponential backoff, and lifecycle events
(`started`, `finished`, `error`, `retrying`, `retries-failed`). The `data-indicator`
attribute automatically toggles a boolean signal during in-flight requests —
no custom loading-state JS needed.

---

## 2. The Architectural Tension

This library's core philosophy (ADR 0005, the JS decision ladder) is:

1. **HATEOAS-first** — HTML is the source of truth; prefer native HTML over scripts.
2. **Zero JS runtime dependencies** — the singleton-guard pattern exists precisely
   to avoid adding Alpine.js or any framework.
3. **Progressive enhancement** — components work without JS; JS enhances.
4. **CSP-safe by construction** — every inline `<script>` carries a nonce.

Datastar is a **~12 KiB JS runtime dependency**. Adopting it as the default
would:

- **Break every consumer** who relies on the zero-JS-dependency guarantee
- **Betray the HATEOAS-first philosophy** (Datastar's `data-*` expressions are
  JavaScript, not pure hypermedia controls)
- **Add SSE lifecycle complexity** to the server side
- **Conflict with HTMX** — both libraries intercept `data-*`/form behavior

**Therefore: Datastar is not a replacement for HTMX in this library.** It is a
**complement** — an opt-in upgrade path for consumers who want real-time
streaming, reactive client state, or who have already chosen Datastar as their
app's interactivity layer.

This is the same relationship the library already has with Alpine.js: documented
as an option (JS guide rung 4), never a dependency.

---

## 3. The Four Tiers of Opportunity

### Tier A: Consumer Coexistence (Documentation Only)

**Cost:** Zero code changes. **Value:** High — unblocks every Datastar-using consumer.

The library's 104 components emit standard HTML with Tailwind classes. Most work
**unchanged** in a Datastar app because they are pure server-rendered HTML. The
only friction points are components with HTMX-specific attributes:

| Component               | HTMX coupling                        | Datastar equivalent                                  |
| ----------------------- | ------------------------------------ | ---------------------------------------------------- |
| `PolledRegion`          | `hx-get`, `hx-trigger="every Ns"`    | `data-on:load="@get('/stream')"` (SSE)               |
| `LoadMore`              | `hx-get`, `hx-swap`, `hx-trigger`    | `data-on:click="@get('/more')"` + server patches     |
| `ConfirmDelete`         | `hx-delete`, `hx-confirm`            | `data-on:click="@delete('/item')"` + confirm signal  |
| `LoadingButton`         | `htmx:balloonBeforeRequest` classes  | `data-indicator` signal                              |
| `FilterDropdown`        | `hx-get`, `hx-trigger="change"`      | `data-on:change="@get('/filter')"`                   |
| `GlobalErrorHandling`   | `htmx:responseError` listeners       | Datastar's built-in retry + `data-on-interval`       |
| `SwapOOB`               | `hx-swap-oob`                        | Server patches target element by ID directly         |
| `InlineLoadingOverlay`  | `hx-indicator`                       | `data-indicator` signal + `data-show`                |

**Deliverable:** A recipe (`docs/recipes/datastar-integration.md`) documenting
each mapping with code examples.

### Tier B: Parallel Opt-In `datastar` Package (Zero New Dependencies)

**Cost:** New package, no go.mod changes. **Value:** High — first-class Datastar support.

The existing `htmx` package proves the pattern: it emits `hx-*` attributes and
injects the runtime `<script>` **without importing any HTMX Go SDK**. A
`datastar` package does the same:

```
datastar/
├── sdk_script.go        # SDKScript() — inject Datastar runtime <script>
├── live_region.go       # LiveRegion() — SSE-powered PolledRegion alternative
├── live_region.templ    # LiveRegion template
├── attributes.go        # Typed data-* attribute builders
├── indicator.go         # Indicator() — loading-state component
├── enums_go.go          # Typed enums (DatastarVersion, IndicatorSignal, etc.)
└── *_test.go            # Golden + unit + contract tests
```

**Key insight:** The library stays dependency-free. Consumers who want SSE
streaming add `github.com/starfederation/datastar-go` to **their** go.mod. The
library's `datastar` package only handles the templ/rendering side (attributes +
runtime injection), exactly like the `htmx` package handles `hx-*` without an
HTMX SDK.

#### SDKScript — Runtime Injection

```go
// Mirrors how layout.Base injects HTMX. CSP-safe via layout.Script.
@datastar.SDKScript(datastar.SDKScriptProps{Nonce: nonce})

// Renders:
// <script type="module" nonce="..." src="https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.2/bundles/datastar.js"></script>
```

Supports `SelfHosted` path (like HTMXCDN override) and version pinning.

#### Attribute Helpers — Typed `data-*` Builders

```go
// Ergonomic, typed, CSP-safe attribute construction
datastar.OnClick("@get('/api/search')")     // → data-on:click="@get('...')"
datastar.Text("$resultCount")               // → data-text="$resultCount"
datastar.Bind("searchQuery")                // → data-bind:searchQuery
datastar.Show("$loading")                   // → data-show="$loading"
datastar.Signals(store)                     // → data-signals={JSON}
datastar.Indicator("searching")             // → data-indicator:searching
```

These return `templ.Attributes`, matching the library's existing
`utils.Class()` / `templ.Attributes` ergonomic pattern.

### Tier C: New Capabilities Datastar Unlocks

These are things the library **cannot do today** — genuinely new value, not
replacements.

#### C1. LiveRegion — SSE-Powered Real-Time Updates

`PolledRegion` polls on an interval. `LiveRegion` establishes an SSE stream:

```templ
@datastar.LiveRegion(datastar.LiveRegionProps{
    URL:       "/stream/metrics",
    AutoStart: true,
    Live:      datastar.LivePolite,
}) {
    @display.StatCard(display.StatCardProps{Label: "Active Users", Value: "—"})
}
```

The server endpoint streams patches:

```go
sse := datastar.NewSSE(w, r)
sse.PatchElementTempl(metricsCard(currentData()))
```

**Why this matters:** Real-time dashboards, live activity feeds, collaborative
multi-cursors, chat — all currently impossible without polling. SSE is the
browser-native real-time transport (no WebSocket complexity, automatic
reconnection with `Last-Event-ID`).

#### C2. Reactive Combobox (No Singleton JS)

The current `forms.Combobox` uses a ~60-line singleton-guard script
(`tcComboboxAttached`) for input filtering, keyboard navigation, and selection.
A Datastar-native variant could do this declaratively:

```html
<div data-signals="{ query: '', open: false, selected: '' }">
  <input data-bind:query
         data-on:input="$open = true"
         data-on:focus="$open = true"
         data-on:click__outside="$open = false" />
  <div data-show="$open">
    for _, opt := range filteredOptions {
      <button data-on:click="$selected = '...'; $open = false">{ opt.Label }</button>
    }
  </div>
</div>
```

Zero custom JS. The filtering could even be server-side via
`data-on:input__debounce.300ms="@get('/filter?q=' + $query)"`.

#### C3. Optimistic UI

HTMX can't do optimistic updates — the UI waits for the server response before
changing. Datastar signals enable instant feedback:

```html
<button data-on:click="$liked = true; @post('/like')">
  <span data-text="$liked ? 'Liked' : 'Like'"></span>
</button>
```

The UI updates instantly; the server confirms or reconciles via SSE patch.

#### C4. Multi-Step Forms Without Round-Trips

Currently, multi-step forms require either HTMX round-trips per step or custom
JS. Datastar signals make step state declarative:

```html
<div data-signals="{ step: 1 }">
  <div data-show="$step === 1">...step 1 fields...</div>
  <div data-show="$step === 2">...step 2 fields...</div>
  <button data-on:click="$step = $step + 1">Next</button>
</div>
```

### Tier D: Full HTMX → Datastar Migration (REJECTED)

Replacing HTMX as the default would:

1. **Break the zero-dependency guarantee** — the library's #1 selling point
2. **Break every `hx-*` consumer** — all 8 `htmx` package components + every
   consumer's `hx-*` attributes
3. **Add ~12 KiB JS to every page** even for static content
4. **Require SSE infrastructure** — not every deployment can hold long-lived
   connections (serverless functions, some CDNs)
5. **Betray HATEOAS-first** — Datastar's `data-*` expressions are JS expressions,
   not pure hypermedia

**Verdict: Do not migrate.** HTMX remains the default; Datastar is opt-in.

---

## 4. Component-by-Component Analysis

### Components That Benefit Most from Datastar

| Component             | Current approach                    | Datastar improvement                          | Impact |
| --------------------- | ----------------------------------- | --------------------------------------------- | ------ |
| `PolledRegion`        | HTMX polling (idle waste, latency)  | SSE `LiveRegion` (push, zero idle)            | **High** |
| `Combobox`            | ~60-line singleton JS               | Signals for filter/open/selected              | **High** |
| `TagsInput`           | ~50-line singleton JS               | Signals for tag array                         | **High** |
| `Tabs`                | ~40-line singleton JS               | Signals for active tab                        | Medium |
| `Carousel`            | ~50-line singleton JS               | Signals for slide index                       | Medium |
| `Dropdown`/`Popover`  | Popover API + position JS           | Signals for open state                        | Low (already native) |
| `GlobalErrorHandling` | ~130-line inline JS                 | Built-in retry/indicator lifecycle            | **High** |
| `LoadingButton`       | `.htmx-request` CSS classes         | `data-indicator` signal                       | Medium |
| `FilterDropdown`      | `hx-get` on change                  | `data-on:change="@get()"`                     | Low |
| `LoadMore`            | `hx-get` + `hx-trigger="revealed"`  | `data-on-intersect="@get()"`                  | Low |

### Components Unaffected by Datastar

All purely presentational components (`Card`, `Badge`, `Avatar`, `Table`,
`Grid`, `StatCard`, icons, forms inputs without interactivity, error pages,
layout shell, navigation, feedback components) work **identically** in HTMX and
Datastar apps. They emit HTML + Tailwind classes — framework-agnostic.

**~80% of the component surface is already Datastar-compatible.**

---

## 5. Case Study: PolledRegion vs LiveRegion

This is the strongest single opportunity. Here is the concrete difference:

### PolledRegion (current — HTMX polling)

```
Client                          Server
  |--- GET /stats (every 10s) --->|
  |<--- 200 <div>full HTML</div>--|  (even if data unchanged)
  |--- GET /stats (every 10s) --->|
  |<--- 200 <div>full HTML</div>--|
  ... (repeats forever, 6 req/min idle)
```

- 6 requests/minute per client, even when idle
- Up to 10s latency on updates
- Full HTML re-rendered every poll
- 100 clients = 600 requests/minute baseline load

### LiveRegion (proposed — Datastar SSE)

```
Client                          Server
  |--- GET /stream (one connection) -->|
  |<--- SSE: patch <div>stats</div>---->|  (only when data changes)
  |<--- SSE: patch <div>stats</div>---->|
  |<--- SSE: patch signals {...}------>|  (signal-only, no HTML!)
  ... (connection stays open, server pushes on change)
```

- 1 connection per client, held open
- Updates pushed instantly when data changes (sub-second)
- Can patch signals only (no HTML round-trip) for pure data updates
- 100 clients = 100 open connections, near-zero idle traffic
- Automatic reconnection with `Last-Event-ID` on disconnect

---

## 6. The Singleton-JS Reduction Opportunity

The library currently has **15 singleton-guard scripts** across interactive
components. These represent ~500+ lines of hand-maintained JavaScript that must
be kept CSP-safe, idempotent, and RTL-aware.

A Datastar-native variant of each component would replace this with declarative
`data-*` attributes — zero custom JS, zero singleton guards, zero nonce
management for that component.

**However**, this should NOT replace the existing components. It should be an
**alternative** for consumers who have opted into Datastar. The existing
singleton-guard components remain the zero-dependency default.

| Approach          | JS dependency | Custom JS lines | CSP nonce needed |
| ----------------- | ------------- | --------------- | ---------------- |
| Singleton-guard   | None          | ~500 (15 scripts) | Yes (per script) |
| Datastar-native   | ~12 KiB       | 0               | No (module script) |
| Native HTML only  | None          | 0               | No               |

---

## 7. Concrete Recommendations & Roadmap

### Phase 1: Foundation (this PR)

1. **`datastar` package** — `SDKScript`, attribute helpers, `LiveRegion`,
   `Indicator`. Zero new dependencies.
2. **Golden + unit tests** — same testing strategy as all other packages.
3. **Recipe** — `docs/recipes/datastar-integration.md` documenting the
   HTMX-to-Datastar attribute mapping and SSE handler pattern.
4. **This research doc + ADR 0030** — document the positioning decision.

### Phase 2: Reactive Component Variants (future)

- `datastar.Combobox` — signal-driven, zero custom JS
- `datastar.TagsInput` — signal-driven
- `datastar.MultiStepForm` — step state via signals
- `datastar.LiveActivityFeed` — SSE-powered infinite scroll feed

Each as a **separate component** in the `datastar` package, not a replacement
for the `forms.*` equivalents.

### Phase 3: Demo Integration (future)

- Add Datastar-powered routes to `examples/demo` showcasing `LiveRegion` with
  a mock SSE endpoint.

---

## 8. Tradeoffs and Risks

| Concern                          | Assessment                                                                 |
| -------------------------------- | -------------------------------------------------------------------------- |
| **Dependency bloat**             | Mitigated: library adds zero deps. Consumer opts in by adding datastar-go. |
| **SSE infrastructure**           | Real: SSE needs long-lived connections. Not all hosts support it (serverless cold starts). Documented as a consumer responsibility. |
| **Learning curve**               | Real: Datastar's signal model is new to HTMX users. Mitigated by recipe + attribute helpers. |
| **Two interactivity models**     | Real: consumers must choose HTMX or Datastar per page. Mitigated by clear decision ladder (JS guide rung 2 vs 7). |
| **`data-*` attribute collision** | Low: Datastar uses `data-signals`, `data-on`, `data-bind` — distinct from `data-tc-*`. |
| **CSP**                          | Non-issue: Datastar loads as `<script type="module" src>` — CSP-safe with `script-src`. No inline handlers needed if using `data-on`. |
| **Bundle size**                  | ~12 KiB (smaller than HTMX's ~14 KiB). Loads as ES module (deferred by default). |
| **Maturity**                     | v1.0.2 stable, 501(c)(3) nonprofit governance, 4.8k GitHub stars, 14+ SDKs. Go SDK maintained by the creator. |

---

## 9. Decision

**Adopt Datastar as an opt-in complement, not a replacement for HTMX.**

- **HTMX remains the default** interactivity layer (zero-dependency, HATEOAS-first).
- **A new `datastar` package** provides first-class Datastar support with zero
  new library dependencies.
- **SSE-powered components** (`LiveRegion`) unlock real-time capabilities the
  library cannot offer today.
- **Consumers choose** per-page or app-wide whether to use HTMX, Datastar, both,
  or neither (native HTML only).

This positions the library as the **framework-agnostic Go component layer** that
works with whichever hypermedia/reactivity approach the consumer prefers — the
strongest competitive position for a server-rendered UI library.

---

## Sources

- https://data-star.dev/ — official site
- https://data-star.dev/guide/getting_started — CDN script tag, v1.0.2
- https://data-star.dev/reference/sdks — SDK overview
- https://github.com/starfederation/datastar-go — Go SDK source
- https://data-star.dev/reference/sse_events — SSE event wire format
- `docs/javascript-guide.md` — this library's JS decision ladder (rung 7: Datastar)
- `docs/adr/0005-js-attachment-patterns.md` — singleton-guard pattern rationale
- `docs/adr/0007-self-host-htmx-default.md` — HTMX as default interactivity
- `docs/adr/0024-polled-region-design.md` — PolledRegion design decisions
