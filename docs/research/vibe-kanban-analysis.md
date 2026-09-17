# Research: What templ-components' KanbanBoard can learn from BloopAI/vibe-kanban

_Date: 2026-09-16_
_Subject: `display.KanbanBoard` vs `BloopAI/vibe-kanban` — adopt / avoid analysis_
_Method: read the vibe-kanban source on GitHub `main` (commit at time of reading), not the marketing site. Files read: `packages/ui/src/components/KanbanBoard.tsx`, `KanbanCardContent.tsx`, `KanbanFilterBar.tsx`, `packages/web-core/src/features/kanban/ui/KanbanContainer.tsx` (the `handleDragEnd` owner), `packages/web-core/src/features/kanban/model/hooks/useKanbanFilters.ts`, `packages/web-core/src/pages/kanban/ProjectKanban.tsx`, `docs/cloud/kanban-board.mdx`, `docs/configuration-customisation/keyboard-shortcuts.mdx`, repo tree + README._

## 1. What vibe-kanban is

A coding-agent orchestration product ("Get 10X more out of Claude Code, Codex…") whose
planning surface is a kanban board. Stack:

- **Frontend:** React SPA, `@hello-pangea/dnd` (the maintained `react-beautiful-dnd` fork) for
  drag-and-drop, Zustand stores for selection state, TanStack Query for data, Electric-style
  sync for realtime board updates.
- **Backend:** Rust (crates: `server`, `db` (sqlx/SQLite), `relay-webrtc`, `relay-tunnel`,
  `mcp`, `worktree-manager`, …) behind a pnpm workspace (`packages/ui`, `web-core`,
  `local-web`, `remote-web`).
- **Status: SUNSETTING.** The README banner links a shutdown announcement, and
  `ProjectKanban.tsx` on `main` is gutted to a `ProjectSunsetPage` stub. The kanban
  components still exist in `packages/ui`, which is what we studied.

Their kanban is one feature of a project-management suite (issues, sub-issues, tags,
assignees, priorities, PR badges, workspaces where agents execute code, diff review,
preview browser). That scope fact matters for the "what NOT to learn" section.

## 2. How their board works (the parts that matter)

### 2.1 Composition (`packages/ui/src/components/KanbanBoard.tsx`)

- `KanbanProvider` — `DragDropContext` + an `inline-grid grid-flow-col
  auto-cols-[minmax(200px,400px)] divide-x` layout (CSS grid with fixed-ish column widths).
- `KanbanCards` — the `Droppable` column body, renders `{children}{provided.placeholder}`
  (the dnd library's space keeper during a drag).
- `KanbanCard` — a `Draggable` wrapper around their `Card` primitive. Notable props:
  `dragDisabled`, `isSelected` (accent ring), `isOpen`, `isMobile`.
- `KanbanHeader` — sticky `top-0`, status **color dot** (`hsl(var(--color))`), column name,
  and a **`+` add-task button** (tooltip'd, per column).

### 2.2 Mobile strategy (`KanbanCard`)

On mobile the card is NOT whole-card draggable. It renders a dedicated **drag handle**
(`DotsSixVerticalIcon`, `cursor-grab`) on the leading edge; the rest of the card stays
tappable for opening the issue. On desktop the whole card is the drag surface
(`provided.dragHandleProps` on the root), and card-open is wired via `onMouseUp`
with `e.button === 0 && !snapshot.isDragging` — **not** `onClick`, because the dnd
runtime eats clicks.

### 2.3 The move pipeline (`KanbanContainer.tsx`, `handleDragEnd`)

1. Bail if no `destination`, or source == destination (column AND index) — no-op drop.
2. **Sort-mode gate:** within-column reorder is blocked unless
   `kanbanFilters.sortField === 'sort_order'` ("Manual"). Cross-column moves are always
   allowed (they're status changes; the server sorts within the column).
3. **Optimistic local state:** `setItems(prev => …)` splices the card into place in a
   client-side `Record<columnId, issueId[]>`.
4. **Persist as a bulk rewrite:** every issue in the destination column (and the source
   column on cross-column moves) gets `sort_order = 1000 * columnIndex + (issueIndex + 1)`
   (both 1-based) and `status_id`, then one `bulkUpdateIssues(updates)` call.
5. **Echo suppression:** an `isSyncingRef` flag is set during the bulk update and cleared
   only after a `setTimeout(…, 500)` "to let Electric sync complete" — so the realtime
   sync of your own write doesn't clobber the optimistic state you already applied.

### 2.4 Card content (`KanbanCardContent.tsx`)

Card anatomy: `displayId` (e.g. `ISS-1`), title, description preview, priority icon,
assignee avatars, tags, PR badges, relationship badges, sub-issue marker, loading dots.
Two details worth stealing conceptually:

- **Tag overflow:** `tags.slice(0, 2)` rendered, then `+{tags.length - 2}` when truncated.
- **Description preview normalization:** `formatKanbanDescriptionPreview` collapses a
  markdown body to one plain-text line, replacing code blocks / images / attachments with
  localized placeholder labels ("code block", "image", "file") instead of leaking syntax.

### 2.5 Filtering and views (`KanbanFilterBar.tsx`, `useKanbanFilters.ts`)

Filter bar above the board: text search, priority / assignee / tag filters, sort field +
direction (manual, priority, created, …), show-sub-issues, show-workspaces, hide-blocked.
Status tabs above the board (Active / All / Backlog / Cancelled) with hidden-status columns
reachable as list views. Hidden-by-default statuses (Backlog, Cancelled) keep the board clean.

### 2.6 Selection model (docs + stores)

Multi-select with `Cmd/Ctrl+Click`, `Shift+Click` range, `X` toggle,
`Shift+J`/`Shift+K` extend, `Cmd+A` select all, `Esc` clear. While 2+ cards are selected,
**drag-and-drop is disabled** and a bulk action bar appears (status / priority / assignee /
delete for all selected).

### 2.7 Keyboard shortcuts (docs)

Vim-flavored: `J/K/H/L` navigation, `?` help, `C` create, `D` delete, `/` search, plus
**chord sequences** (`G S` go-to-settings, `W D` duplicate workspace, `X P` create PR, …)
with a documented 500ms inter-key window.

### 2.8 Concurrency

Realtime last-write-wins, documented as such: "If you and a teammate move the same issue
at the same time, the last action wins."

## 3. What our KanbanBoard already is (grounding)

- `display.KanbanBoard` (`display/kanban.templ` + `display/kanban.go`): server-rendered
  columns/cards; the move exchange is ONE hidden form (`card` / `column` / `index`) carrying
  `wire.Action` attributes, submitted via `requestSubmit()` by a CSP-safe singleton script.
  Works identically under htmx (`hx-post` + component-owned `hx-target="#<board>"
  hx-swap="outerHTML"`) and Datastar (`@post` + `wire.Handler` `PatchModeOuter`).
- HTML5 drag-and-drop computing the insertion index from pointer position, adjusting for
  the dragged card's own removal within a column, skipping no-op drops, ignoring cross-board
  drops (`CSS.escape`-guarded selectors).
- **Per-card keyboard move buttons** (prev/next column), hover-revealed, always visible on
  coarse pointers (`pointer:coarse` CSS rule) — the board is fully keyboard operable.
- `role="status"` live region announcing "Moving X to Y" / "Moved X to Y" around the swap.
- Cards/columns without stable consumer IDs render but don't participate in moves
  (invalid states unrepresentable); an unwired board is read-only with zero JS.
- Server contract: `display.ParseKanbanMove(r) → KanbanMove{Card, Column, Index}`; the
  server owns ALL ordering — we send only the target insertion index.
- Proven by `visualtest/kanban_e2e_test.go` (keyboard + synthetic HTML5 DnD, both
  transports, coarse pointer, cross-board guard, live-region announcements) and demo
  RTL/mobile e2e.

## 4. What we CAN learn (adopt / adapt)

### 4.1 Per-column "add card" affordance — real gap, cheap

Their column header `+` button. Ours has only `EmptyColumnText`. Proposal shape:
`KanbanColumn.Action templ.Component` slot (or `AddHref string`) rendered in the column
header next to the count badge. Server-rendered link/button; zero JS. Full component-
anatomy treatment (golden, a11y, e2e, CHANGELOG) if implemented.

### 4.2 Column status color dot — real gap, cheap

Their header shows a colored status indicator per column. Proposal shape:
`KanbanColumn.Tone` typed enum reusing the existing Badge/status palette + map lookup
(our standard pattern), rendering a small dot before the title. Gives consumers instant
column semantics (todo=blue, done=green…) without any new CSS mechanism.

### 4.3 Tag overflow "+N" pattern — recipe, not API

`tags.slice(0, 2)` + `+{n}` matches our CountBadge "N+" overflow convention exactly.
No component change needed — document as a recipe composing `KanbanCard.Content` with
`Badge`/`Avatar`/`StatusBadge` (see 4.4).

### 4.4 A card-anatomy recipe — their strongest deliverable

Their card content design (display ID, title, priority, truncated tags, assignees,
plain-text description preview) is battle-tested UX. All of it already composes through
our `KanbanCard.Content` slot. Worth writing `docs/recipes/kanban-card-anatomy.md` with a
worked example (ID badge + title + tags-with-overflow + avatar row). Their markdown→
one-line preview with attachment placeholders is a nice consumer-side snippet to include.

### 4.5 Drag-handle separation on touch — consider, don't rush

On touch they drag via a dedicated handle so taps keep meaning "open". We currently make
the whole card draggable AND keep move buttons visible on coarse pointers; taps vs drags
are disambiguated by the browser (a tap without movement isn't a drag). If we ever add
whole-card click-to-open, adopt their separation: an explicit `data-tc-kanban-handle`
region for dragging, card body stays clickable. Do NOT preemptively add it — it costs
visible UI chrome and our e2e proves the current path operable.

### 4.6 Sort-mode gating as a DOCUMENTED contract — no code change

Their insight: **draggability and view sort interact.** A board sorted by priority/created
shouldn't accept within-column reorders (they're meaningless; the sort re-asserts itself).
We already have the mechanics: a consumer rendering a sorted view can omit `Wire`
(read-only) — but that also kills cross-column moves, which ARE meaningful (status
change). Options if we ever want parity: document that consumers' `ParseKanbanMove`
handlers should reject same-column moves when the view is server-sorted (the server owns
ordering; the index is advisory), or add `KanbanBoardProps.Reorderable bool`. For now: a
godoc + recipe note. Cheapest correct step.

### 4.7 Hidden-by-default statuses — product thought, consumer-side

Backlog/Cancelled columns hidden behind tabs keeps boards clean. For us that's a consumer
navigation concern (`Tabs` + filtered board renders), worth one sentence in a recipe.

### 4.8 Validation of our architecture (the quiet lesson)

Their move pipeline needs optimistic client state + a bulk-write endpoint + an
echo-suppression ref + a 500ms timer to avoid the realtime sync clobbering their own
optimism. Our pipeline is: one form POST, server re-renders the board, HTML lands. No
client/server state can drift because the client holds no state. Their `handleDragEnd`
plus `isSyncingRef` is a museum piece of the bug class HATEOAS eliminates — cite this doc
next time someone proposes optimistic kanban moves.

## 5. What we should NOT learn (reject, with reasons)

### 5.1 Integer `sort_order = 1000 * columnIndex + index`

The worst idea in the file. Consequences: (a) one drag rewrites **every card in the
destination column AND the source column** — O(column) writes per move; (b) the encoding
couples two independent orders (column position × card position) into one number, so
reordering columns invalidates every card's stored order; (c) 1000-gap sparse ordering
sounds nice but is immediately consumed by big columns. Our contract — send
`card`/`column`/`index`, server owns persistence — is strictly better. If we ever document
ordering guidance for consumers: fractional indexing (LexoRank-style) or server-computed
densities, never derived integers.

### 5.2 Client-side state as the source of truth

`setItems` optimistic splice + `bulkUpdateIssues` + sync echo fights. Every divergence
between client items and server truth is a bug we'd have to write. Our singleton script
holds drag-transient state only between `dragstart` and `drop`, then submits and forgets.

### 5.3 The SPA dependency stack for a board

`@hello-pangea/dnd` + Zustand + TanStack Query + a realtime sync engine to move cards.
Our equivalent interaction ships ~2KB of dependency-free, CSP-safe script and one form.
The dnd library is also the reason they can't use plain `onClick` (5.5) and why their
DnD doesn't work on touch without a bespoke `isMobile` variant tree (2.2). ADR-0033
(no Web Components / no SPA runtime) and ADR-0030/0036 (transport-agnostic wiring) already
encode this; vibe-kanban is the counterexample that proves it.

### 5.4 Vim chords and single-letter shortcuts

`J/K/H/L`, `X`, `Shift+J/K`, `G S`, `W D`, `X P` with a 500ms window. Reasons to reject
for a component library: (a) undiscoverable — their own docs need a `?` cheat-sheet dialog;
(b) a library component cannot know page context (is the user typing in an input? does the
page already bind `X`?); (c) chords are hostile to screen-reader users and to
international keyboards; (d) it's app-level chrome, not board semantics. Our per-card
prev/next buttons (real buttons, real labels, tabbable) follow WAI-ARIA APG instead.

### 5.5 `onMouseUp`-as-click hack

They fire card-open on `onMouseUp` with `button === 0 && !isDragging` because the dnd
runtime swallows clicks. A symptom, not a technique. Our card titles are native
`<a href>` elements — the browser handles click-vs-drag for free.

### 5.6 Multi-select + bulk actions + "selection disables drag"

A selection store (`useIssueSelectionStore`), a multi-select hook, a bulk action bar, and
a modal rule that selection freezes DnD. Heavy state machine, and the freeze rule is a
usability trap (select, then try to drag all — blocked). If a consumer ever needs bulk
moves, the server-rendered answer is checkboxes + one form POST of checked IDs + target
column — no client store. Until a consumer asks: YAGNI.

### 5.7 Whole-card draggable with no visible affordance (desktop)

On desktop their cards give no cue they're draggable until grabbed. Ours show `cursor-grab`

- hover border + move buttons on hover/focus/coarse-pointer. Discoverability is a feature;
  keep ours.

### 5.8 Last-write-wins as the entire concurrency story

Fine for their product, but they document zero conflict detection. Our contract makes
concurrency a server concern and stays out of it; if it ever needs documenting, the note
belongs in a transport-wiring recipe (consumers can add a version field / ETag check in
their move handler), not in the component.

### 5.9 The scope that killed it (the meta-lesson)

The kanban is the planning surface of a suite that grew workspaces, agent runners, diff
review, a preview browser with devtools, relay tunnels, WebRTC, and an MCP server — and
the product sunset. For us the boundary is already drawn and correct:
`KanbanBoard` is a display + move-exchange component. Persistence, filtering, sorting,
realtime, selection, and issue-panel UX are consumer concerns (composable via our
`Content` slots and `Wire`). Keep saying no.

## 6. Follow-up candidates (ranked — 1–4 shipped 2026-09-16, #5 deferred)

| # | Item                                                                                                                                         | Cost    | Value                                          |
| - | -------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ---------------------------------------------- |
| 1 | ✅ DONE 2026-09-16 — `docs/recipes/kanban-card-anatomy.md` (card anatomy + tag overflow + description preview)                               | trivial | high — unlocks their best UX with existing API |
| 2 | ✅ DONE 2026-09-16 — `KanbanColumn.Action` slot (per-column add-card), browser-proven e2e                                                    | small   | high — real gap                                |
| 3 | ✅ DONE 2026-09-16 — `KanbanColumn.Tone` status dot                                                                                          | small   | medium                                         |
| 4 | ✅ DONE 2026-09-16 — Godoc/recipe note: sorted views ↔ read-only or server-rejected reorder (`ParseKanbanMove` godoc + recipe move contract) | trivial | medium — prevents consumer confusion           |
| 5 | Drag-handle variant — ONLY if whole-card click actions ever land (still deferred, TODO_LIST #190)                                            | medium  | low today                                      |

## 7. Sources

- `BloopAI/vibe-kanban` GitHub, `main` branch at time of reading (sunset banner present).
  Key files listed in the header. Local copies were pulled via the GitHub API for reading;
  nothing was committed to this repo.
- Our side: `display/kanban.templ`, `display/kanban.go`, `visualtest/kanban_e2e_test.go`,
  `examples/demo/kanban_demo.templ`, `utils/wire` (ADR-0036).
