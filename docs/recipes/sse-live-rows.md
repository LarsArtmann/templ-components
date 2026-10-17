# Recipe: SSE Live Rows

Append `<tr>` fragments to a live table over Server-Sent Events — the
tail-a-log, live-feed, job-progress pattern — without re-rendering the whole
table. Built from three shipped pieces: the Table's `BodyID` (a stable id on
the `<tbody>`), the `Body` slot (custom rows), and htmx's SSE extension.

## The shape

```go
@display.Table(display.TableProps{
    Headers: []string{"Time", "Level", "Message"},
    BodyID:  "live-rows", // the append target
    Body:    logRows(entries), // server-rendered starting rows
})

templ logRows(entries []LogEntry) {
    for _, e := range entries {
        @logRow(e)
    }
}

templ logRow(e LogEntry) {
    <tr class="hover:bg-gray-50 dark:hover:bg-gray-800">
        <td class="px-4 py-2 font-mono text-xs">{ e.Time.Format(time.Kitchen) }</td>
        <td class="px-4 py-2">{ e.Level }</td>
        <td class="px-4 py-2">{ e.Message }</td>
    </tr>
}
```

`BodyID` lands on the `<tbody>`; `Body` renders the rows it starts with.

## The listener

One plain `<div>` owns the EventSource connection; a child does the swap.
The SSE extension is NOT part of `htmx.min.js` — load `htmx-ext-sse` after
the core (CSP: serve it via `layout.Script` so the nonce is applied):

```go
<div hx-ext="sse" sse-connect="/events">
    <div sse-swap="row" hx-target="#live-rows" hx-swap="beforeend"></div>
</div>
@layout.Script(nonce, "/static/htmx-ext-sse.js")
```

- `sse-connect` + `hx-ext="sse"` must sit on the SAME element (it owns the
  `EventSource`).
- `sse-swap="row"` matches the server's `event: row` name exactly — the
  browser discards events nobody is listening for (unnamed messages are
  listened for as `sse-swap="message"`; there is no catch-all).
- `hx-target="#live-rows"` + `hx-swap="beforeend"` redirect the message into
  the tbody as an append — the listener div itself never renders content.

## The endpoint

Emit one named event per row, `data:` carrying the rendered `<tr>`:

```go
func sseLogs(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-store")

    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "streaming unsupported", http.StatusInternalServerError)
        return
    }

    for entry := range tailLogs(r.Context()) {
        fragment, err := renderFragment(logRow(entry)) // strings.Builder — see the
        if err != nil {                                // hybrid-strings-builder recipe
            continue
        }
        fmt.Fprintf(w, "event: row\ndata: %s\n\n", fragment)
        flusher.Flush()
    }
}
```

Rules that keep the stream alive through proxies: flush after every event,
disable proxy buffering (`X-Accel-Buffering: no` for nginx), and stop when
`r.Context()` is done — the browser reconnects automatically and `sse-connect`
resumes the stream (send `Last-Event-ID`-aware data if gaps matter).

## Appending safely

- Give every rendered row a stable `id` (`id={ "log-" + e.ID }`) so later
  updates can patch a row in place: send
  `<tr id="log-42" hx-swap-oob="outerHTML:#log-42">…</tr>` inside any SSE
  message — out-of-band directives inside a message never touch the
  `sse-swap` target.
- Append-only feeds eventually grow unbounded. Cap the DOM: when the row
  count crosses your limit, send a full `Body`-slot re-render as an
  `hx-swap-oob="outerHTML:#live-rows"` message (the bounded-tables recipe's
  cap discipline applies to live tables too).
- Rows arriving while the user is mid-interaction do not steal focus —
  pure appends are the least disruptive live update. Avoid `hx-swap-oob`
  into anything the user is typing in.

## Datastar equivalent

The Datastar runtime consumes SSE natively (datalines carry signals or
fragments — see `docs/datastar-runtime-facts.md` and the datastar
integration recipe). A fragment patch targeting the tbody's `id` with mode
`beforeend`-equivalent is the twin of the pattern above; prefer it when the
page is already Datastar-wired.

## See also

- `custom-table-rows.md` — the `Body` slot and `TypedHeaders` basics
- `bounding-tables.md` — cap + notice discipline for unbounded data
- `hybrid-strings-builder-rendering.md` — rendering a templ fragment to a
  string for the `data:` payload
