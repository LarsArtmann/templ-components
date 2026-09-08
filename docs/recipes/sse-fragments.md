# Recipe: SSE fragment streams (named events → server-rendered HTML)

Pattern for live dashboards under a strict CSP (`script-src 'self'`, no
`unsafe-eval`): a plain `EventSource` client receives **named SSE events**
whose payloads are server-rendered HTML fragments, and swaps them into
fixed containers. Proven in production by the CV career-pipeline dashboard.

## Why not Datastar/Alpine here

Clients that compile expressions at runtime (`new Function()`) are dead
under CSP — see `docs/datastar-runtime-facts.md` (`'unsafe-eval'` is never
granted). The fix is named events + a nonce'd vanilla-JS controller, not a
CSP exception.

## Server: one event per region, HTML as the wire format

Render each region as a **fragment component** and broadcast under a
region name. Keep the fragment component the ONLY place the markup lives —
initial server render and every broadcast reuse it, so the JS never
re-implements HTML.

```templ
templ applicationsFragment(apps []Application) {
    if len(apps) == 0 {
        @display.EmptyState(display.EmptyStateProps{Icon: icons.Inbox, Title: "No applications yet"})
    } else {
        <ul id="applications-list">{ for _, a := range apps { ... } }</ul>
    }
}
```

Broadcast per store notification:

```go
func (h *Handlers) broadcast(ctx context.Context) {
    events := []sse.Event{
        {Event: "dashboard", Data: marshal(h.signals())},                    // JSON scalars
        {Event: "applications", Data: render(applicationsFragment(h.apps()))},
        {Event: "recent-events", Data: render(recentEventsFragment(h.events()))},
    }
    h.broadcaster.BroadcastMany(events...)
}
```

## Wire contract

- **JSON scalars** (`dashboard` event): counters and flags the client
  writes into `#stat-*` nodes by id.
- **HTML fragments** (`applications`, `approvals`, `interviews`,
  `recent-events`): full region markup, swapped with `innerHTML`.
- One connect sends the full initial state (scalars + fragments) before
  broadcasts — first paint needs no extra round-trip. **Tests: drain ALL
  initial events before asserting on broadcasts**; a partial drain lets a
  stale initial fragment interleave into the broadcast window (a real,
  load-dependent flake).

## Client: ~40 lines, CSP-safe

```js
var targets = { applications: "applications-container",
                "recent-events": "recent-events-container" };
var source = new EventSource("/api/pipeline/events");
source.addEventListener("dashboard", function (e) {
  var s = JSON.parse(e.data);
  Object.keys(s).forEach(function (k) { setText("stat-" + k, String(s[k])); });
});
Object.keys(targets).forEach(function (name) {
  source.addEventListener(name, function (e) {
    document.getElementById(targets[name]).innerHTML = e.data;
  });
});
source.onopen  = function () { setConnected(true); };
source.onerror = function () { setConnected(false); };  // banner + auto-reconnect
```

**Scripts inside swapped fragments never execute** — the HTML spec skips
`<script>` content inserted via `innerHTML`. Anything dynamic in a
fragment must be data (attributes/text), not scripts: timestamps ship
pre-rendered (`display.RelativeTime` with `AutoRefresh: false`); state
flags ride `data-*` attributes read by the page controller.

## Reconnect + staleness UX

`EventSource` reconnects automatically; surface the gap: an `onerror`
handler unhides a "Live updates interrupted" banner (`role="status"`,
`aria-live="polite"`) and `onopen` hides it again.

## Checklist

- [ ] One fragment component per region (single markup source)
- [ ] JSON only for scalars; HTML for anything structured
- [ ] Client sets `textContent` (never `innerHTML` with untrusted strings; fragment HTML is trusted server output)
- [ ] No `<script>` inside fragments (it will not run)
- [ ] Reconnect banner with `role="status"`
- [ ] Tests drain the full initial state before broadcast assertions
