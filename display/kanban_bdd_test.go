package display

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// BDD-style behaviour specs for the Kanban board: user-visible behaviour,
// not markup mechanics. The drag-and-drop mechanics are browser-proven in
// visualtest/kanban_e2e_test.go; these specs pin the render contract the
// script depends on.

// Spec: a consumer wiring a board gets every affordance needed to move
// cards — draggable cards, a hidden move form carrying the contract fields,
// keyboard buttons, and the singleton script.
func TestKanbanBehaviourWiredBoardExposesMovePipeline(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, KanbanBoard(KanbanBoardProps{
		Columns: kanbanTestColumns(),
		Wire:    &wire.Action{URL: "/api/kanban/move"},
	}))

	utils.AssertContainsAll(t, html,
		`draggable="true"`,
		`data-tc-kanban-form`,
		`data-tc-kanban-f-card`,
		`data-tc-kanban-f-column`,
		`data-tc-kanban-f-index`,
		`tcKanbanAttached`,
	)
}

// Spec: dropping a card on a column submits that column's identity — every
// column list (including empty ones) carries its column ID and title so the
// script can route the move and announce the target.
func TestKanbanBehaviourColumnsCarryMoveIdentity(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, KanbanBoard(KanbanBoardProps{
		Columns: kanbanTestColumns(),
		Wire:    &wire.Action{URL: "/api/kanban/move"},
	}))

	utils.AssertContainsAll(t, html,
		`data-tc-kanban-column-body="todo"`,
		`data-tc-kanban-column-body="doing"`,
		`data-tc-kanban-column-body="done"`,
		`data-tc-kanban-col-title="In progress"`,
	)
}

// Spec: an empty column tells the user it is empty and stays a valid drop
// target (the placeholder lives inside the same list element).
func TestKanbanBehaviourEmptyColumnPlaceholder(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, KanbanBoard(KanbanBoardProps{
		Columns:        kanbanTestColumns(),
		EmptyColumnText: "Nothing here yet",
	}))

	done := html[strings.Index(html, `data-tc-kanban-column-body="done"`):]
	if !strings.Contains(done, "Nothing here yet") {
		t.Error("empty column does not show the custom placeholder")
	}

	if strings.Contains(done, "No cards") {
		t.Error("empty column shows the default placeholder despite a custom one")
	}
}

// Spec: cards link to detail pages when Href is set, without losing the
// move identity the script needs.
func TestKanbanBehaviourCardLinksKeepMoveIdentity(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, KanbanBoard(KanbanBoardProps{
		Columns: kanbanTestColumns(),
		Wire:    &wire.Action{URL: "/api/kanban/move"},
	}))

	utils.AssertContainsAll(t, html,
		`href="/issues/2"`,
		`data-tc-kanban-card="c2"`,
	)
}

// Spec: a single-column board is a plain list — no move buttons, since no
// adjacent column exists.
func TestKanbanBehaviourSingleColumnHasNoMoveButtons(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, KanbanBoard(KanbanBoardProps{
		Columns: []KanbanColumn{{ID: "only", Title: "Only", Cards: []KanbanCard{{ID: "c1", Title: "Solo"}}}},
		Wire:    &wire.Action{URL: "/api/kanban/move"},
	}))

	utils.AssertNotContains(t, html, `data-tc-kanban-move=`)
	utils.AssertContains(t, html, `draggable="true"`)
}

// Spec: BaseProps propagate to the board root (ID, class, extra attrs).
func TestKanbanBehaviourBasePropsPropagate(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, KanbanBoard(KanbanBoardProps{
		BaseProps: utils.BaseProps{
			ID:    "custom-board",
			Class: "my-board",
			Attrs: map[string]any{"data-testid": "kanban"},
		},
		Columns: kanbanTestColumns(),
	}))

	utils.AssertContainsAll(t, html,
		`id="custom-board"`,
		`my-board`,
		`data-testid="kanban"`,
	)
}
