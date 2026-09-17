---
title: Kanban Board
description: Wire a dual-transport KanbanBoard — drag-and-drop and keyboard moves, the server move contract, sorted views, and the optimistic pending register.
---

`display.KanbanBoard` renders a complete kanban board (columns, cards, count
badges, a hidden move form) and wires drag-and-drop AND per-card keyboard
move buttons through **one hidden form carrying one `wire.Action`** — the
same markup serves htmx and Datastar, so a wired board works under either
runtime with zero component changes.

## Quick start

```templ
import (
    "github.com/larsartmann/templ-components/display"
    "github.com/larsartmann/templ-components/utils"
    "github.com/larsartmann/templ-components/utils/wire"
)

@display.KanbanBoard(display.KanbanBoardProps{
    BaseProps: utils.BaseProps{ID: "board"},
    Wire: &wire.Action{
        Transport: dialect,        // htmx (zero value) or wire.TransportDatastar
        Method:    wire.MethodPost, // moves are mutations — always POST
        URL:       "/api/kanban/move",
        Target:    "#" + "board",   // htmx: swap the board
    },
    Columns: columns, // []display.KanbanColumn
    CSRFToken: csrfToken,
})
```

Render rules that make the board moveable:

- Columns and cards need **stable consumer IDs** (`KanbanColumn.ID`,
  `KanbanCard.ID`). Cards or columns without IDs still render, but they do
  not participate in moves.
- Unwired boards (nil `Wire`) render read-only: no script, no move form.
- Mutations MUST set `Method: wire.MethodPost` — an unspecified method
  renders GET in both dialects.

## The server move contract

Every move submits one urlencoded body — `card`, `column`, `index` — and the
library decodes it for you:

```go
move, err := display.ParseKanbanMove(r)
if err != nil {
    http.Error(w, "invalid move", http.StatusBadRequest)
    return
}
```

Apply the move to your state, then **re-render the SAME board ID**. The
client's optimistic pending register resolves from that re-render, so the
board id you render back must match the id you rendered first. `index` is
**advisory**: the server owns the real ordering, so server-side validation
(unknown column → 404, wrong-state card → 422) is expected and the client
reverts cleanly on those responses.

## Sorted views: advisory index

If you render a board from a server-sorted slice (top-N, filtered, grouped),
the client-computed index may not match your stored order. The documented
contract: accept the move, apply your own ordering rule, and re-render the
board from your sorted view. Never trust the client index as authoritative —
it is a hint about where the card was dropped.

## The optimistic pending register

A wired board never freezes while a move is in flight. The moment a move is
submitted (drop OR keyboard button, both transports):

- the card moves in the DOM and wears `tc-kanban-pending` (dim + spinner
  ring in the card's top-right corner, `aria-busy="true"`),
- the target column's "No cards" placeholder hides,
- both columns' count badges and `aria-label`s update immediately.

On success the server's re-render lands and pending state is forgotten. On
failure (4xx/5xx, network error) the card **snaps back to its exact original
position**, counts are recomputed from the DOM, the card flashes a red
border for 4s, and an sr-only `role="alert"` announces the revert — never
`aria-live="assertive"` (the repo's urgency policy).

Design consequence: the card's top-right corner is reserved for the pending
ring. Keep persistent affordances out of custom card content's corner — the
[card anatomy recipe](https://github.com/larsartmann/templ-components/blob/master/docs/recipes/kanban-card-anatomy.md)
shows the pattern.

Concurrency semantics (overlapping moves, double failures, listener
double-bind) are pinned by `TestKanbanJSConcurrentMoves`; the pending and
failed states also have deterministic PNG goldens.

## Keyboard support

Every card renders hover/keyboard move buttons (↑/↓ per WAI-ARIA direction
with RTL-aware arrows). They submit the same hidden move form as
drag-and-drop — one server contract serves pointer and keyboard users. The
buttons stay visible under `pointer: coarse` (touch), and drag-and-drop
remains pointer-only by design.

## Further reading

- [Transport Wiring](/guides/transport-wiring) — the `wire.Action` contract,
  `wire.Handler` response targeting, and the GlobalErrorHandling interplay.
- [Card anatomy recipe](https://github.com/larsartmann/templ-components/blob/master/docs/recipes/kanban-card-anatomy.md)
  — the `Content` slot patterns (identity row, tag overflow, assignees).
- [ADR-0041](https://github.com/larsartmann/templ-components/blob/master/docs/adr/0041-kanban-optimistic-move-pending-register.md)
  — the pending register decision.
