# JavaScript in Templ Projects: The Complete Guide

> A deep-research synthesis across the official templ docs, templ-components'
> own ADR 0005, the T.A.H. Stack (templ + Alpine + HTMX), Datastar, React islands,
> CSP compliance, TypeScript workflows, and the View Transitions API.

---

## The Decision Ladder

Before writing ANY JavaScript, walk this ladder. Stop at the first rung that solves
your problem.

| Rung | Tool                         | When to use                                                                      | JS required |
| ---- | ---------------------------- | -------------------------------------------------------------------------------- | ----------- |
| 1    | Native HTML                  | `<details>`, `<form>`, `:checked`, `:target`                                     | Zero        |
| 2    | HTMX                         | Server round-trips, partial updates, form submits                                | Zero (attrs) |
| 3    | Inline `<script>` + singleton | Document-level event delegation (dropdowns, tabs, accordions)                   | Minimal, CSP-safe |
| 4    | Alpine.js                    | Client-side state (toggles, modals, filtering) without server round-trip         | ~8KB lib    |
| 5    | External JS bundle (esbuild) | Complex client-side logic, charts, editors                                       | Full build  |
| 6    | React/Vue islands            | Rich interactive widgets that can't be server-rendered                           | Framework + bundler |
| 7    | Datastar                     | Server-driven reactive UI with SSE streaming (replaces HTMX + Alpine)            | ~15KB lib   |

---

## Pattern 1: Singleton Guard + Event Delegation (Zero-Dependency)

**Used by:** templ-components (this repo), production Go apps that want zero JS
dependencies.

Every interactive component emits ONE inline `<script>` that registers a
document-level event listener, guarded by a global flag so HTMX re-renders are
idempotent:

```templ
templ Dropdown(props DropdownProps) {
	<div data-tc-dropdown={ id }>
		<button data-tc-dropdown-trigger={ id }>Open</button>
		<div data-tc-dropdown-menu={ id } class="hidden">...</div>
	</div>
	<script nonce={ props.Nonce }>
		if (!window.tcDropdownAttached) {
			window.tcDropdownAttached = true;
			document.addEventListener('click', function(e) {
				var trigger = e.target.closest('[data-tc-dropdown-trigger]');
				if (trigger) {
					var id = trigger.getAttribute('data-tc-dropdown-trigger');
					var menu = document.querySelector('[data-tc-dropdown-menu="' + id + '"]');
					menu.classList.toggle('hidden');
				}
			});
		}
	</script>
}
```

### Key principles

- `window.tc*Attached` guard prevents double-binding on HTMX swaps.
- `document.addEventListener` + `e.target.closest('[data-*]')` = event delegation
  (works for dynamically added elements).
- `nonce={ props.Nonce }` on every `<script>` for CSP compliance.
- Zero inline handlers (`onclick`, etc.) — fully `script-src-attr 'none'` compliant.

### When to choose this pattern

- You're building a **component library** (like templ-components) where adding a
  JS framework dependency is unacceptable.
- Your interactivity is **document-level event delegation** (dropdowns, accordions,
  tabs, copy buttons, dismiss buttons).
- You need **HTMX DOM-swap compatibility** — document-level listeners automatically
  handle newly swapped elements.

### Exceptions: per-instance IIFE

Modal and Drawer use per-instance IIFE (not document-level delegation) because they
need **per-instance state** (focus trap, previous-focus element). See
`overlayTrapJS()` in `display/shared.go`.

### Pros / Cons

| Pros                                      | Cons                                    |
| ----------------------------------------- | --------------------------------------- |
| Zero dependencies                         | Verbose for complex state               |
| HTMX-compatible                           | Manual focus management                 |
| CSP-safe (nonce on every script)          | No reactivity system                    |
| Works after DOM swaps                     | Each component repeats the boilerplate  |

**This repo uses 12 singleton guards** across all interactive components (see
ADR 0005).

---

## Pattern 2: HTMX-Only (Zero Custom JS)

**Used by:** Server-rendered apps where the server is the source of truth.

```templ
templ SearchResults(query string, results []Result) {
	<form hx-get="/search" hx-target="#results" hx-trigger="keyup changed delay:300ms">
		<input name="q" value={ query } />
	</form>
	<div id="results">
		for _, r := range results {
			<div>{ r.Title }</div>
		}
	</div>
}
```

### Key principle

The same templ component renders **both** the initial page AND the HTMX partial
update. No separate API endpoint or JSON serialization needed.

### Pros / Cons

| Pros                                  | Cons                                     |
| ------------------------------------- | ---------------------------------------- |
| Simplest possible interactivity       | Server round-trip for every interaction  |
| Progressive enhancement (works w/o JS)| No client-side state                     |
| Reuses same templ component           | Latency on every interaction             |

---

## Pattern 3: HTMX + Alpine.js (The T.A.H. Stack)

**Used by:** templUI, most community templates, the "GOAT Stack"

Alpine.js handles client-side state (toggles, dropdowns); HTMX handles server
communication.

```templ
templ EditableRow(item Item) {
	<tr x-data="{ editMode: false }">
		<td>
			<template x-if="!editMode">
				<span x-text={ templ.JSONString(item.Email) }></span>
			</template>
			<template x-if="editMode">
				<input type="text" x-model={ templ.JSONString(item.Email) } />
			</template>
		</td>
		<td>
			<button x-show="!editMode" @click="editMode = true">Edit</button>
			<button x-show="editMode" hx-post={ "/save/" + item.ID } hx-swap="outerHTML">Save</button>
		</td>
	</tr>
}
```

### Passing Go data to Alpine

```templ
<div x-data={ templ.JSONString(myGoStruct) }>
```

`templ.JSONString` HTML-encodes the JSON so it's safe in an attribute value.

### Common Alpine features in templ projects

| Feature               | Directive          | Use case                              |
| --------------------- | ------------------ | ------------------------------------- |
| Toggles/show-hide     | `x-show`, `@click` | Dropdowns, modals, confirmation       |
| Conditional rendering  | `x-if` (template)  | Edit mode toggling                    |
| Two-way data binding   | `x-model`          | Form inputs                           |
| Text binding           | `x-text`           | Displaying server data reactively     |
| List rendering         | `x-for` (template) | Client-side filtering                 |
| Attribute binding      | `x-bind` / `:attr` | Dynamic classes, styles               |

### Gotcha: string values in `x-text` / `x-model`

When templ renders a Go string directly into `x-text` or `x-model`, the string is
output **without quotes**, which Alpine interprets as a JS variable name rather than
a string literal. Use `templ.JSONString()` or wrap in a `x-data` initialization to
avoid this.

### Pros / Cons

| Pros                                | Cons                                        |
| ----------------------------------- | ------------------------------------------- |
| No build step                       | Another library (~8KB)                      |
| Tiny, declarative                   | `x-data` values are JS expressions (escape!) |
| Handles client-side state elegantly | Doesn't help with server communication      |

---

## Pattern 4: Datastar (Server-Driven Reactivity)

**Used by:** datastarui (shadcn/ui port in Go), growing adoption in Go community

Replaces **both** HTMX and Alpine.js with a single ~15KB library using SSE.

### Frontend (templ)

```templ
templ Counter(signals CounterSignals) {
	<div data-signals={ templ.JSONString(signals) }>
		<button data-on:click="@post('/increment')">+1</button>
		<span data-text="$count">0</span>
	</div>
}
```

### Backend (Go handler)

```go
sse := datastar.NewSSE(w, r)

// Option A: Patch HTML elements (morphed into DOM by ID)
sse.MergeFragmentTempl(counterTemplate(updatedSignals))

// Option B: Patch only signals (most efficient — no HTML sent)
sse.MarshalAndMergeSignals(patchJSON)
```

### Datastar vs HTMX

| Aspect          | HTMX                               | Datastar                                    |
| --------------- | ---------------------------------- | ------------------------------------------- |
| Architecture    | Frontend-driven (attrs on trigger) | Server-driven (server sends patches)        |
| Transport       | AJAX (fetch) request/response      | SSE streaming                               |
| DOM updates     | `innerHTML` swap (morph via ext)   | DOM morphing by default (preserves state)   |
| Client state    | None (pair with Alpine.js)         | Built-in reactive signals                   |
| Bundle size     | ~14KB                              | ~15KB                                       |
| Expressions     | `hx-*` attributes + hyperscript    | `data-*` attributes use plain JS + `$signal` |

### Pros / Cons

| Pros                                          | Cons                                    |
| --------------------------------------------- | --------------------------------------- |
| Single library (no Alpine needed)             | Steeper learning curve                  |
| SSE streaming for real-time                   | SSE lifecycle management on server      |
| DOM morphing preserves focus/scroll/input     | Newer ecosystem than HTMX               |
| Signal-only patches (no HTML round-trip)      | `data-*` expressions are JS (escaping!) |

---

## Pattern 5: React/Vue Islands

**Used by:** Apps with complex widgets (charts, editors, React Flow) that can't be
server-rendered.

```templ
templ ChartWidget(data ChartData) {
	<div id="chart" data-config={ templ.JSONString(data) }>
		<script nonce={ nonce }>
			bundle.renderChart(document.currentScript.closest('div'));
		</script>
	</div>
}
```

Build: `esbuild --bundle --global-name=bundle --outfile=assets/js/chart.js src/Chart.tsx`

The server renders a placeholder `<div>` with `data-*` attributes; the JS bundle
mounts the React component into it. The rest of the page is pure templ.

---

## Templ's Built-in JS Features (Reference)

| Feature                         | Purpose                                              | Example                                                  |
| ------------------------------- | ---------------------------------------------------- | -------------------------------------------------------- |
| `templ.NewOnceHandle()`         | Render a `<script>` block only once per page         | `@handle.Once() { <script>...</script> }`               |
| `templ.WithNonce(ctx, nonce)`   | Inject CSP nonce into context                        | All inline `<script>` tags get `nonce="..."`            |
| `templ.JSFuncCall("fn", data)`  | Call JS function with JSON-encoded Go data           | `onclick={ templ.JSFuncCall("alert", msg) }`            |
| `templ.JSExpression("event")`   | Pass raw JS expression (event objects, `this`)       | `onclick={ templ.JSFuncCall("fn", templ.JSExpression("event")) }` |
| `templ.JSONString(data)`        | Encode Go data as JSON in an HTML attribute          | `x-data={ templ.JSONString(data) }`                     |
| `templ.JSONScript("id", data)`  | Embed data in `<script type="application/json">`     | Large data payloads, chart configs                       |
| `{{ value }}` in `<script>`     | Interpolate Go values into JS (auto-escaped)         | `const msg = "{{ user.Name }}";`                        |

### `templ.OnceHandle` — deduplicate scripts

When a component is used multiple times on a page, use `OnceHandle` so its
`<script>` block renders only once:

```templ
var helloHandle = templ.NewOnceHandle()

templ hello(label, name string) {
	@helloHandle.Once() {
		<script>
			function hello(name) { alert('Hello, ' + name + '!'); }
		</script>
	}
	<button data-name={ name }>{ label }</button>
}
```

> **Warning:** Don't write `@templ.NewOnceHandle().Once()` — this creates a new handle
> each time, defeating the purpose. Always declare the handle as a package-level `var`.

### `templ.JSONScript` — safe data embedding

```templ
@templ.JSONScript("page-config", config)
```

Output: `<script id="page-config" type="application/json">{"theme":"dark"}</script>`

Client-side: `const config = JSON.parse(document.getElementById('page-config').textContent);`

---

## CSP Compliance Checklist

```
[x] Every <script> has nonce={ props.Nonce }
[x] Zero inline event handlers (onclick, onload, etc.)
[x] Use document.addEventListener with event delegation
[x] Use templ.WithNonce() middleware to inject nonce via context
[x] External scripts use layout.Script(nonce, src, attrs)
[ ] CSP header includes: script-src 'nonce-<random>'
[ ] If using Alpine.js hx-on:click: add 'unsafe-hashes' or 'unsafe-inline'
```

> **CSP gotcha:** Even with nonces, inline `onclick` / `hx-on:click` attributes
> require `'unsafe-hashes'` in your CSP. This is a browser limitation, not a templ
> issue. Avoid inline handlers entirely if your CSP is strict — use event delegation
> (Pattern 1) or `addEventListener` via IIFE scripts.

### Nonce middleware pattern

```go
func withNonce(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonce := securelyGenerateRandomString()
		w.Header().Set("Content-Security-Policy",
			fmt.Sprintf("script-src 'nonce-%s'", nonce))
		ctx := templ.WithNonce(r.Context(), nonce)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
```

---

## TypeScript Workflow

### Project structure

```
project/
├── assets/js/          # esbuild output (served by Go's http.FileServer)
├── src/ts/             # TypeScript source
│   └── main.ts
├── *.templ             # templ templates
└── go.mod
```

### Build command

```bash
esbuild --bundle --minify --outfile=assets/js/app.js src/ts/main.ts
```

### Reference in templ

```templ
@layout.Script(nonce, "/assets/js/app.js", templ.Attributes{"defer": true})
```

### Passing Go → TypeScript data

```templ
@templ.JSONScript("page-data", pageConfig)
```

```ts
const config = JSON.parse(document.getElementById('page-data').textContent);
```

### Hot reload

- [WGO](https://github.com/bokwoon95/wgo) — stable Go dev server with file watching
- [Templiér](https://github.com/romshark/templier) — hot-reload for templ dev mode
- esbuild's `--watch` flag for TS rebuilds

---

## View Transitions API (Zero-JS Animations)

HTMX 1.9.0+ supports native browser View Transitions — smooth page transitions with
zero JS:

```templ
<div hx-get="/page" hx-swap="innerHTML transition:true" hx-target="main">
	Load Page
</div>
```

```css
::view-transition-old(root) { animation: 180ms both fade-out; }
::view-transition-new(root) { animation: 420ms both fade-in; }
```

For multi-page apps (MPA), Chrome supports cross-document view transitions natively
— no JS at all. This eliminates the "ka-chunk" of full-page reloads.

---

## Event Delegation — Avoiding Inline Handlers

Mozilla considers inline event handlers (`onclick`, `onload`) bad practice. Templ
enables a clean alternative using IIFE scripts with `document.currentScript`:

```templ
templ hello(label, name string) {
	<div>
		<input type="button" value={ label } data-name={ name }/>
		<script nonce={ nonce }>
			(() => {
				let parent = document.currentScript.closest('div');
				let btn = parent.querySelector('input[data-name]');
				btn.addEventListener('click', function() {
					let name = btn.getAttribute('data-name');
					alert('Hello, ' + name + '!');
				});
			})()
		</script>
	</div>
}
```

The key technique: `document.currentScript` lets the inline script find its own
position in the DOM, then traverse to sibling elements. This avoids needing unique
IDs and works after HTMX swaps.

---

## What NOT to Do

| Anti-pattern                                        | Why                                              | Fix                                            |
| --------------------------------------------------- | ------------------------------------------------ | --------------------------------------------- |
| Inline `onclick` / `onload` handlers                | Breaks strict CSP, Mozilla bad practice         | Use `addEventListener` or event delegation    |
| `@templ.NewOnceHandle().Once()` inline              | Creates new handle each render                   | Declare as package-level `var`                 |
| `templ.JSExpression` with user data                 | Bypasses escaping, XSS risk                      | Use `templ.JSONString` or `templ.JSFuncCall`   |
| Per-element `addEventListener` without delegation   | Lost on HTMX DOM swap                            | Use `document.addEventListener` + `closest()`  |
| Missing `nonce` on `<script>`                       | Blocked by CSP                                   | Use `layout.Script(nonce, src, attrs)`         |
| `x-text={ goString }` without quotes (Alpine.js)    | Alpine interprets as JS variable, not string     | Use `templ.JSONString(goString)`               |

---

## Recommendation by Project Type

| Project type                         | Recommended stack          | Why                                                |
| ----------------------------------- | ------------------------- | -------------------------------------------------- |
| Content site / blog                  | HTMX only                 | Server renders everything, minimal interactivity    |
| Admin dashboard                      | HTMX + singleton-guard JS | Type-safe, CSP-safe, zero deps (this repo's pattern)|
| SaaS app (forms heavy)               | HTMX + Alpine.js          | Client-side form state, validation, toggles         |
| Real-time app (chat, live data)      | Datastar                  | SSE streaming, signal patching, built-in reactivity |
| App with complex widgets             | HTMX + React islands      | Server-render most pages, islands for widgets       |

---

## Sources

- [templ official docs: JavaScript](https://templ.guide/syntax-and-usage/script-templates/)
- [templ official docs: CSP](https://templ.guide/security/content-security-policy/)
- [templ official docs: HTMX](https://templ.guide/server-side-rendering/htmx/)
- [templ official docs: Datastar](https://templ.guide/server-side-rendering/datastar/)
- [templ official docs: React islands](https://templ.guide/syntax-and-usage/using-react-with-templ/)
- [ADR 0005: JavaScript Attachment Patterns](adr/0005-js-attachment-patterns.md)
- [GitHub Issue #220: CSP safe templ](https://github.com/a-h/templ/issues/220)
- [templ-csp-example](https://github.com/leonyork/templ-csp-example)
- [templUI (templ + Alpine.js component library)](https://github.com/axzilla/templui)
- [DatastarUI (Go/templ port of shadcn/ui)](https://github.com/CoreyCole/datastarui)
- [datastar-templ (type-safe Datastar helpers)](https://github.com/Yacobolo/datastar-templ)
- [Go Templ + Alpine.js + HTMX guide](https://sachinsharma.dev/blogs/go-templ-alpine-js-interactive-htmx)
- [HTMX View Transitions essay](https://htmx.org/essays/view-transitions/)
- [surreal (JS helper library)](https://github.com/gnat/surreal)
