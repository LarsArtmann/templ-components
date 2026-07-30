# htmx.PolledRegion

Auto-refreshing HTMX region that fetches new content on an interval. The region
replaces itself (`hx-swap="outerHTML"`) on each poll, re-rendering server-side.

## When to use

- Live dashboards (stats, activity feeds, health checks)
- Any content that should update without page reload
- Progressive enhancement — content renders server-side first, then polls

## Basic usage

```templ
@htmx.PolledRegion(htmx.PolledRegionProps{URL: "/partials/stats", Every: "10s"}) {
    @display.StatCard(display.StatCardProps{Label: "Messages", Value: msgCount})
}
```

## Eager loading

`Eager: true` fires the first fetch immediately on page load, then continues
polling at the interval. Use when you want the initial render to show a
placeholder that's replaced with live data ASAP:

```templ
@htmx.PolledRegion(htmx.PolledRegionProps{
    URL:   "/api/live-feed",
    Every: "5s",
    Eager: true,
}) {
    <p>Loading live data...</p>
}
```

## Timestamp footer

`ShowTimestamp: true` renders an "Updated HH:MM:SS" footer so operators can
verify polling is active — if the timestamp freezes, polling has stalled:

```templ
@htmx.PolledRegion(htmx.DefaultPolledRegionProps()) {
    <!-- children -->
}
```

Custom time format:

```go
props := htmx.DefaultPolledRegionProps()
props.URL = "/api/heartbeat"
props.TimeFormat = "2006-01-02 15:04:05" // full date+time instead of time-only
```

## Screen reader announcements

`aria-live` controls how screen readers announce updates:

- `PolledLivePolite` (default) — announces after current speech
- `PolledLiveAssertive` — interrupts current speech (use for alerts)
- `PolledLiveOff` — no announcements

```templ
@htmx.PolledRegion(htmx.PolledRegionProps{
    URL:   "/api/alerts",
    Every: "1s",
    Live:  htmx.PolledLiveAssertive,
}) {
    @display.Alert(display.AlertProps{Type: display.AlertError, Message: "High error rate!"})
}
```

## Server-side handler

The endpoint returns an HTML fragment that replaces the entire region:

```go
func statsHandler(w http.ResponseWriter, r *http.Request) {
    stats := getStats()
    // Return just the region HTML (not a full page)
    component := statsRegion(stats)
    component.Render(r.Context(), w)
}
```

## Swap styles

Default is `outerHTML` (the region div replaces itself). You can use other
HTMX swap styles:

```go
htmx.PolledRegionProps{
    URL:   "/api/counter",
    Every: "2s",
    Swap:  htmx.SwapInnerHTML, // replace children only, keep the wrapper div
}
```
