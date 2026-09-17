# Recipe: Kanban Card Anatomy — Rich Card Content

**Audience:** Consumers building a `display.KanbanBoard` whose cards carry more
than a title — issue keys, priority, tags, assignees, a description preview.

**Problem:** `KanbanCard` is deliberately small (ID, Title, Href, Content
slot). Teams coming from Trello/Linear/Jira expect issue-key badges, truncated
tag rows, avatar stacks, and a one-line description preview. All of that
composes through the existing `Content` slot — but the composition has four
non-obvious details (tag overflow, avatar overlap, preview flattening, and
what moves mean on a sorted board).

**Outcome:** Copy-pasteable patterns for each card-content ingredient, plus
the move-handling contract for sorted views.

---

## The `Content` slot is the anatomy

Everything below renders inside one `KanbanCard.Content` slot:

```go
display.KanbanCard{
    ID:      task.ID,
    Title:   task.Title,
    Href:    "/tasks/" + task.ID,
    Content: taskCardContent(task),
}
```

`Content` is a `templ.Component`. Wrap your composition in a consumer-side
`templ` template and pass it through. Nothing in `display` needs to change.

## 1. Display ID + priority (the identity row)

A stable, human-readable key (PROJ-17, GH-402) plus a status or priority
badge is the cheapest scannability win:

```templ
templ taskIdentity(task Task) {
    <div class="flex flex-wrap items-center gap-1.5">
        @display.Badge(display.BadgeProps{
            Text: task.Key,
            Size: display.BadgeSizeSM,
        })
        @display.StatusBadge(task.Priority)
    </div>
}
```

`display.StatusBadge` maps ~20 common status strings (`"open"`, `"in
progress"`, `"blocked"`, …) onto badge colors automatically — use it for
anything status-shaped, and `Badge` with an explicit `Type` when you need
custom semantics (e.g. `BadgePurple` for "P1").

## 2. Tag overflow: render 2, then "+N"

Five tags never fit on a card. The convention (same as `CountBadge`'s "N+"
overflow): render at most two `Badge`s, then a muted "+N" remainder. Render
the remainder **once, at position 2**, instead of breaking out of the loop —
templ has no `break` inside loops:

```templ
templ taskTags(task Task) {
    if len(task.Tags) > 0 {
        <div class="mt-1.5 flex flex-wrap items-center gap-1">
            for i, tag := range task.Tags {
                if i == 0 || i == 1 {
                    @display.Badge(display.BadgeProps{
                        Text: tag,
                        Type: display.BadgeNeutral,
                        Size: display.BadgeSizeSM,
                    })
                } else if i == 2 {
                    <span class="text-xs font-medium text-gray-500 dark:text-gray-400">
                        +{ strconv.Itoa(len(task.Tags) - 2) }
                    </span>
                }
            }
        </div>
    }
}
```

Two or fewer tags render with no "+N" text; one tag renders alone. Do not
truncate server-side and lose the count — the "+N" tells the user there is
more, which is the whole point.

## 3. Assignee stack: overlapping avatars

Avatars overlap with a negative inset (physical `-space-x` is fine here —
overlap direction follows document order, not writing direction). Prefer
initials over images on cards; they render instantly and stay legible at XS:

```templ
templ taskAssignees(task Task) {
    if len(task.Assignees) > 0 {
        <div class="mt-1.5 flex -space-x-1.5">
            for _, a := range task.Assignees {
                @display.Avatar(display.AvatarProps{
                    Initials: a.Initials,
                    Alt:      a.Name,
                    Size:     display.AvatarSizeXS,
                })
            }
        </div>
    }
}
```

Cap the stack at ~3–4 avatars and fold the rest into a "+N" avatar
(`Initials: "+3"`) — same overflow logic as tags.

## 4. Description preview: markdown to one muted line

A raw markdown blob is unreadable on a card. Flatten it to one plain line:
attachments become labeled placeholders, inline syntax is stripped, whitespace
collapses, and the result truncates. This is consumer-side Go — no component
change. The compiling implementation lives in
`display/kanban_example_test.go` (`oneLinePreview` plus its
`cardAnatomy*Pattern` regexps) — **that file is the single source of truth**;
copy from there rather than from this doc so the two can never drift:

```go
// oneLinePreview flattens a markdown description into a single plain-text
// line for card display: attachments become [label] placeholders, link and
// emphasis syntax is stripped, whitespace collapses. Returns "" when nothing
// readable remains. Truncated previews end with an ellipsis.
func oneLinePreview(markdown string) string
```

(Rune cap: 90, then an ellipsis.)

```templ
templ taskPreview(task Task) {
    if preview := oneLinePreview(task.DescriptionMarkdown); preview != "" {
        <p class="mt-1.5 truncate text-xs text-gray-500 dark:text-gray-400">{ preview }</p>
    }
}
```

`truncate` (overflow ellipsis) plus the rune-cap are belt and suspenders: the
rune-cap keeps the DOM small, `truncate` handles long unbroken strings.

## 5. Hidden columns: navigation, not component scope

Boards with "Backlog"/"Cancelled" graveyards stay readable by hiding those
columns behind `display.Tabs` or a link, and rendering the board filtered.
Filtering the `Columns` slice is a consumer navigation concern — the component
intentionally has no filter/WIP state (see the 2026-09-16 scope decision,
TODO_LIST #191).

---

## Moves on a sorted view: the contract

`KanbanBoard`'s move exchange submits `(card, column, index)` — but **the
server owns the real ordering**. The index is advisory. That has one
consequence you must handle deliberately:

**If the view is server-sorted (by priority, created-at, …), a same-column
move is meaningless.** The client drops the card at index 3, the server
re-renders, the sort re-asserts itself and the card jumps back. The user sees
a broken board.

The component cannot know your sort — but your move handler does. Reject
same-column moves on sorted views instead of silently persisting them:

**Pseudo-code:** `columnOf`, `sortedByPriority`, `persistMove`, and
`renderBoard` stand in for your app's lookups and render path — the
library-relevant lines are `display.ParseKanbanMove` and the 422 rejection.

```go
func moveCard(w http.ResponseWriter, r *http.Request) {
    move, err := display.ParseKanbanMove(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // A priority-sorted view: within-column reorders are no-ops the sort
    // will undo — reject them explicitly instead of accepting a ghost move.
    if move.Column == columnOf(move.Card) && sortedByPriority {
        http.Error(w,
            "column is sorted by priority; within-column moves are not accepted",
            http.StatusUnprocessableEntity)
        return
    }

    persistMove(move)
    renderBoard(w, r) // re-render the SAME board element (the swap contract)
}
```

Cross-column moves (a status change) stay meaningful on sorted views — do not
reject those. If you need within-column reordering, render that board view
unsorted and let the index be authoritative.

The full exchange contract (hidden form, defaults, re-render response) is
documented on `KanbanBoardProps.Wire` and decoded by `ParseKanbanMove`.

---

## The optimistic pending register (and the reserved top-right corner)

A wired board does not freeze while a move is in flight (ADR-0041). The
moment a move is submitted — drop OR keyboard button, both transports — the
client moves the card optimistically and marks it with a **pending
register**:

- the card dims and a small **spinner ring spins in its top-right corner**
  (a `::after` pseudo-element positioned with `inset-inline-*`, so it
  mirrors in RTL; `prefers-reduced-motion` users get a static ring),
- the card carries `aria-busy="true"`,
- the target column's "No cards" placeholder hides and both columns' count
  badges/aria-labels update immediately.

When the server's re-render lands, the swap replaces the board DOM and the
pending state is forgotten. When the move FAILS (4xx/5xx, network error),
the card **snaps back to its exact original position**, counts are
recomputed from the DOM, the card flashes a red border for 4s, and an
sr-only `role="alert"` announces "Moving X to Y failed. The board was
restored." — never `aria-live="assertive"` (repo urgency policy).

**Design consequence for your `Content` slot: the card's top-right corner is
reserved.** The pending spinner ring is drawn there, so avoid persistent
affordances (menu buttons, checkboxes) in that corner of custom card
content — or accept that they are covered for the duration of a move. Put
such controls in the identity row's normal flow instead (as the demo board
does).

Concurrency semantics (overlapping moves, double failures, listener
double-bind) are pinned by `TestKanbanJSConcurrentMoves`; the pending and
failed visual states by `TestKanbanPendingRegisterVisualStates`.

---

## Complete worked example

The Go side of this recipe is compile-proven and render-exercised by
`display/kanban_example_test.go` (`taskCardContent` +
`ExampleKanbanBoard_cardAnatomy`) — treat that file as the source of truth
for the patterns above. The `templ` snippets here are consumer-side sketches
(they assume your `Task` model); when they and the compiled example ever
disagree, trust the example and fix this doc.

**Related — the column Action slot:** per-column affordances (e.g. an
add-card button) compose through `KanbanColumn.Action`. The consumer owns
the whole element, including its accessibility: the content must be
keyboard-reachable (a real button or link), and mutations must set
`wire.MethodPost` explicitly — an unspecified Method renders GET.

**Grounding:** this recipe transcribes the battle-tested card content design
of `BloopAI/vibe-kanban` (research: `docs/research/vibe-kanban-analysis.md`
§4.3, §4.4, §4.6, §4.7) onto `templ-components` primitives.
