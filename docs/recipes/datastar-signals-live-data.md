# Recipe: Datastar Signals Live Data

Push live values into a page over Server-Sent Events without sending ANY
HTML — the server streams `datastar-patch-signals` events, the client's
`data-text` bindings re-render automatically. This is the cheap channel:
where the htmx twin (`sse-live-rows.md`) ships rendered `<tr>` fragments,
this recipe ships a few dozen JSON bytes per tick and lets the runtime's
reactive graph update every bound node.

Every wire-format claim below is verified against the pinned runtime bundle
(`go-datastar/static` v0.6.1 = Datastar 1.0.3, sha256 `5d6b7794…`,
`datastar.TestPinnedRuntimeBundleContract` guards the tokens) and
re-proven in a browser under a strict nonce-only CSP by
`visualtest.TestDatastarCSPModeBrowser`.

## The shape (client)

Declare the initial signals once, bind nodes to dotted signal paths:

```templ
<div data-signals="{'metrics': {users: 42, revenue: 1280}}">
	<p><span data-text="$metrics.users"></span> users online</p>
	<p>$<span data-text="$metrics.revenue"></span> MRR</p>
</div>
```

- `data-text` re-renders whenever the referenced signal changes — there is
  no per-binding subscription to write.
- **Signal names in attribute suffixes are LOWERCASE** (`data-signals-users`,
  `data-indicator:busy`): the HTML parser lowercases attribute NAMES, so a
  camelCase suffix silently declares a different signal than the expression
  reads. camelCase keys are fine inside VALUE-borne syntax (`data-signals`
  JSON above) and inside expressions (`$metrics.users`).

## The endpoint (server)

One SSE stream per page region; each tick is one `datastar-patch-signals`
event carrying a `signals` dataline with a JSON **merge patch** (RFC 7396
semantics — the bundle's handler is literally `mergePatch`, nested objects
deep-merge, `null` deletes):

```go
// writeDatastarSignals writes one datastar-patch-signals SSE event. The
// format is pinned by examples/demo/sse_test.go and guarded by the bundle
// contract: event name, keyed datalines, blank-line terminator.
func writeDatastarSignals(w io.Writer, patch map[string]any) error {
	buf, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	var b strings.Builder

	b.WriteString("event: datastar-patch-signals\n")
	b.WriteString("data: signals " + string(buf) + "\n")
	b.WriteString("\n")

	_, err = io.WriteString(w, b.String())

	return err
}

func sseMetrics(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)

		return
	}

	_ = http.NewResponseController(w).SetWriteDeadline(time.Time{})
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			if err := writeDatastarSignals(w, map[string]any{
				"metrics": map[string]any{"users": nextCount(), "revenue": nextRevenue()},
			}); err != nil {
				return
			}

			flusher.Flush()
		}
	}
}
```

- `Content-Type: text/event-stream` + `Cache-Control: no-cache`, one flush
  per event, and a write deadline of zero (`SetWriteDeadline(time.Time{})`)
  so a slow tick never kills the stream — identical hygiene to the
  patch-elements twin in `examples/demo/main.go` (`writeDatastarPatch`).
- A real newline terminates each dataline and a BLANK line terminates the
  event. A literal `\n` (backslash-n text) means the browser never
  dispatches the event — this exact bug shipped once (see
  `docs/datastar-runtime-facts.md`).
- The demo endpoint skeleton this recipe references lives at
  `/api/datastar/signals` (examples/demo) — its wire format is pinned by
  `examples/demo/sse_test.go`.

### `onlyIfMissing`

The runtime accepts a second dataline, `onlyIfMissing`: when truthy, keys
the client ALREADY holds are left untouched (the patch becomes
create-only). Emit it only when you mean it:

```
event: datastar-patch-signals
data: signals {"session": {"startedAt": "…"}}
data: onlyIfMissing true

```

## Reacting to patches (client listener)

`data-on-signal-patch` runs an expression on every signal patch; the patch
payload arrives as the `patch` argument (`argNames:["patch"]` in the pinned
bundle):

```templ
<div data-on-signal-patch="$patch.metrics && $patch.metrics.users > 90 && toast('high load')"></div>
```

A companion attribute, `data-on-signal-patch-filter`, narrows WHICH patches
reach the handler (any other key on the attribute throws `KeyNotAllowed` in
the pinned bundle). Its exact value grammar is only partially decoded from
the minified runtime — it parses as an object applied over the patch detail
— so treat the precise shape as unverified (evidence tier D) and prefer the
`$patch`-guard expression above until decoded at the next re-audit.

## Local vs backend signals

By DEFAULT the runtime excludes signals whose names contain a `_` segment
from request bodies (the fetch-options destructure defaults
`filterSignals: {include: /.*/, exclude: /(^|\.)_/}`). The convention:
`_`-prefixed signals are LOCAL client state (transient UI flags) and never
travel to the server; everything else is backend state and rides the next
action request. Name accordingly — a `_draft` signal stays client-side, a
`metrics` signal shows up in request payloads.

## CSP

The whole channel works under a strict nonce-only CSP: bind and listener
expressions execute through the same CSP-mode script mechanism as every
other `data-*` attribute (browser-proven with zero violations, 2026-11-07 —
`docs/datastar-runtime-facts.md` §CSP). Pair the page nonce into
`datastar.SDKScriptProps.Nonce` AND `<html data-nonce>` via
`layout.DatastarNonceAttrs` as usual.

## When to prefer which channel

| Need | Reach for |
| ---- | --------- |
| New/changed ROWS or fragments (structure) | `sse-live-rows.md` (patch-elements) |
| Numbers, counters, toggles — values only | THIS recipe (patch-signals) |
| Both from one stream | Both events on one SSE response are legal; keep each tick single-purpose |

The library stays transport-agnostic on purpose: `utils/wire` models the
request half (ADR-0036), and server streaming stays with your handler code —
the htmx/Datastar decision ladder in `docs/javascript-guide.md` is unchanged.
