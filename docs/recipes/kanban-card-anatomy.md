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
change:

```go
var (
    attachmentPattern  = regexp.MustCompile(`!\[([^\]]*)\]\([^)]*\)`)
    markdownSyntaxPattern = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)|\*\*([^*]+)\*\*|\*([^*]+)\*|` + "`([^`]+)`")
)

// oneLinePreview flattens a markdown description into a single plain-text
// line for card display: attachments become [label] placeholders, link and
// emphasis syntax is stripped, whitespace collapses. Returns "" when nothing
// readable remains. Truncated previews end with an ellipsis.
func oneLinePreview(markdown string) string {
    line := attachmentPattern.ReplaceAllString(markdown, "[$1]")
    line = markdownSyntaxPattern.ReplaceAllString(line, "$1$2$3$4")
    line = strings.Join(strings.Fields(line), " ")

    const maxPreviewRunes = 90
    runes := []rune(line)
    if len(runes) > maxPreviewRunes {
        line = string(runes[:maxPreviewRunes]) + "…"
    }

    return line
}
```

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

## Complete worked example

The library's own test suite proves this composition compiles and renders:
see `ExampleKanbanBoard_cardAnatomy` in `display/kanban_example_test.go`
(source of the patterns above).

**Grounding:** this recipe transcribes the battle-tested card content design
of `BloopAI/vibe-kanban` (research: `docs/research/vibe-kanban-analysis.md`
§4.3, §4.4, §4.6, §4.7) onto `templ-components` primitives.
