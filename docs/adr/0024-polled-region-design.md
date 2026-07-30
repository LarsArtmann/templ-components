# ADR-0024: PolledRegion Design

**Status:** Accepted (2026-07-30)
**Decider:** Lars Artmann

## Context

Dashboards need auto-refreshing content regions. HTMX provides polling via
`hx-trigger="every Ns"`, but every consumer hand-rolls the same wrapper:
a div with `hx-get`, `hx-trigger`, `hx-swap`, `aria-live`, and optional
timestamp.

The key design questions were:

1. Default swap style: `outerHTML` or `innerHTML`?
2. Should the first poll fire on load (eager) or after the first interval?
3. How to make polling visible to operators?
4. How to handle screen-reader announcements?

## Decision

1. **Default swap: `outerHTML`** — the region replaces itself on each poll.
   This means the entire region re-renders server-side. Consumers who want to
   preserve the wrapper div can set `Swap: SwapInnerHTML`.

2. **Eager is opt-in** — `Eager: false` by default. The first fetch happens
   after the first interval. Consumers who want immediate data set `Eager: true`
   and render a placeholder as children.

3. **ShowTimestamp defaults to true** — renders "Updated HH:MM:SS" in a
   `<time>` element. The timestamp ticks forward on each successful poll
   because the whole region re-renders. If it freezes, polling has stalled.
   `TimeFormat` is configurable (default: `"15:04:05"`).

4. **aria-live defaults to `polite`** — announcements happen after current
   speech finishes. `assertive` is available for alert feeds.

## Consequences

- Server-side handler must return the complete region HTML (not a partial).
- `outerHTML` swap means the children re-render on every poll — efficient for
  small regions, potentially expensive for large ones.
- The timestamp uses `time.Now()` server-side, so it reflects the server's
  clock, not the client's.
