# ADR-0041: Kanban Optimistic Move with an Honest Pending Register

Date: 2026-09-17
Status: Accepted
Deciders: Lars Artmann, Crush

## Context

A wired `display.KanbanBoard` shipped a perception gap: a move (drop or keyboard
button) filled the hidden form and called `requestSubmit()`, but **nothing
visibly happened until the server response swapped the board**. The drag ghost
snapped back to the old position and the card sat still through the whole
round-trip. On localhost that is a blink; on a real network (hundreds of
milliseconds, or a hung server) the board looks broken or dead. Screen-reader
users heard "Moving X to Y." and then silence until the swap landed.

The failure mode is worse than latency: if the request failed (500, offline),
the board kept showing the pre-move state with **no explanation** — the client
never registered that an action was in flight, so it could not honestly report
its outcome.

The ask: register client actions that are **not yet on the server** — an
optimistic update, but *transparent*: the UI moves immediately AND clearly
marks the move as pending until the server confirms, and undoes it visibly when
the server refuses.

## Decision

Wired boards apply moves **optimistically in the DOM** and carry an explicit
pending register until the server responds. Default-on for every wired board;
no props flag (YAGNI — there is no meaningful "off" story: the previous
behavior was the bug).

### 1. Optimistic placement (transport-agnostic)

`tcKbSubmit` physically moves the card's DOM node to the target column at the
same index the server receives, hides the target column's empty placeholder
(`data-tc-kanban-empty`), and syncs both columns' count badges
(`data-tc-kanban-count`) and `aria-label`s from the DOM. This happens inside
the shared submit pipeline, so it composes with htmx and Datastar alike —
the transports only take over for the exchange itself.

### 2. Pending register on the moved card

The moved card gets the `tc-kanban-pending` class (dim + a CSS-only spinner
ring via `::after`, reduced-motion safe) and `aria-busy="true"`. The live
region keeps its existing two-step register: "Moving X to Y." at submit,
"Moved X to Y." once the swap lands.

A visible board-level status chip was **rejected**: it duplicates the signal
the card already carries at the exact point of action, costs layout space on
every board, and adds a third thing to keep in sync. Consumers who want a
global "saving" indicator can compose the transport events (htmx indicators,
`datastar-fetch`) — the library stays out of the page chrome.

### 3. Success — pending can never stick

Three independent clearers, any one of which is sufficient:

- **htmx**: `htmx:afterRequest` with `detail.successful === true` on the move
  form (verified in the vendored htmx 2.0.10 source: `afterRequest` fires from
  `onload` after the swap, with `successful = !isError`).
- **Datastar**: `datastar-fetch` with `type: 'finished'` on the move form
  (verified in the pinned v1.0.3 bundle: `finished` dispatches in `finally`
  after every action; on the ≥400 path `error` fires first, so the
  success-clear that follows is a no-op).
- **The existing live-region swap poll** (kanbanAnnounceJS, now polling 30s
  instead of 5s): strips pending alongside the post-swap announcement.

Clearing is idempotent and keyed by board id, so a fast first swap clearing a
second in-flight move's indicator is harmless — the swap re-renders server
truth, which is the state the second request will confirm or refute.

### 4. Failure — revert honestly

- **htmx**: `htmx:afterRequest` with `successful !== true` — per the vendored
  source this covers 4xx/5xx AND network errors, aborts, and timeouts
  (`successful` is undefined for those, never `true`). `htmx:responseError` /
  `htmx:sendError` are wired as belt-and-braces.
- **Datastar**: `datastar-fetch` `type: 'error'` (HTTP ≥ 400, immediate) or
  `'retries-failed'` (network errors retry with the runtime's default backoff
  first — up to ~2 minutes of honest "still trying" before the revert).

Revert restores the card to its exact original position (guarded against a
detached original sibling), unhides placeholders, recomputes counts/aria-labels
from the DOM, removes the pending state, and:

- flashes `tc-kanban-move-failed` (red border) on the card for 4s, and
- announces "Moving X to Y failed. The board was restored." into a new
  `role="alert"` live region (the repo's urgency policy: `role="alert"`,
  never `aria-live="assertive"`).

Every listener filters on `data-tc-kanban-form` so unrelated wired elements
(add/reset buttons inside the board) cannot trigger a revert.

## Consequences

- Wired-board goldens change: `data-tc-kanban-count` (all boards),
  `data-tc-kanban-empty` (empty columns), and the wired-only
  `role="alert"` region. Unwired boards gain only the two inert data
  attributes.
- No props/API change; no new dependencies; JS remains a CSP-safe
  nonce'd singleton; CSS additions live in `templates/custom.css` with dark
  and reduced-motion variants.
- The demo move endpoints sleep 800ms on purpose so the pending register is
  perceivable — real apps have latency; the demo should show the state the
  library ships.
- Browser proof lives in `visualtest/kanban_e2e_test.go` (flaky board pair:
  one delayed card proving optimistic placement + pending, one 500 card
  proving revert) under both runtimes.
