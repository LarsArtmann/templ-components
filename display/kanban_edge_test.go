package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// Edge cases: degenerate boards and identities.

// TestKanbanEdgeNoColumns renders an empty board without panicking.
func TestKanbanEdgeNoColumns(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, KanbanBoard(KanbanBoardProps{
		Wire: &wire.Action{URL: "/api/kanban/move"},
	}))
	utils.AssertContains(t, html, `role="region"`)
	utils.AssertNotContains(t, html, "<section")
}

// TestKanbanEdgeCardsWithoutIDs verifies cards lacking a stable ID render
// (statically) instead of emitting an empty move identity — the server
// could never route such a move.
func TestKanbanEdgeCardsWithoutIDs(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, KanbanBoard(KanbanBoardProps{
		Columns: []KanbanColumn{{
			ID:    "todo",
			Title: "To do",
			Cards: []KanbanCard{{Title: "Anonymous"}, {ID: "c2", Title: "Known"}},
		}},
		Wire: &wire.Action{URL: "/api/kanban/move"},
	}))

	utils.AssertNotContains(t, html, `data-tc-kanban-card=""`)
	utils.AssertContains(t, html, `data-tc-kanban-card="c2"`)
	utils.AssertContains(t, html, "Anonymous")
}

// TestKanbanEdgeColumnsWithoutIDs verifies columns lacking an ID render but
// are not wired as drop targets (empty column-body identity).
func TestKanbanEdgeColumnsWithoutIDs(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, KanbanBoard(KanbanBoardProps{
		Columns: []KanbanColumn{
			{Title: "No id", Cards: []KanbanCard{{ID: "c1", Title: "Card"}}},
			{ID: "known", Title: "Known"},
		},
		Wire: &wire.Action{URL: "/api/kanban/move"},
	}))

	utils.AssertNotContains(t, html, `data-tc-kanban-column-body=""`)
	utils.AssertContains(t, html, `data-tc-kanban-column-body="known"`)

	// The card in the id-less column is not draggable either — its column
	// cannot receive or send a well-formed move.
	utils.AssertNotContains(t, html, `draggable="true"`)
}

// TestKanbanEdgeColumnIDsWithSpaces verifies heading ids stay valid HTML
// (whitespace collapsed) and the raw identity survives in the body
// attribute where the move contract needs it.
func TestKanbanEdgeColumnIDsWithSpaces(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, KanbanBoard(KanbanBoardProps{
		Columns: []KanbanColumn{{ID: "in progress", Title: "In progress"}},
	}))

	utils.AssertContains(t, html, `data-tc-kanban-column-body="in progress"`)
	utils.AssertNotContains(t, html, `id="kb`+"-col-in progress"+`-title"`)
}

// TestKanbanEdgeMoveFormIsLastResortOnly verifies the hidden move form does
// not leak into the visual layout (display hidden) and carries the no-JS
// POST action for plain form semantics.
func TestKanbanEdgeMoveFormIsHidden(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, KanbanBoard(KanbanBoardProps{
		Columns: kanbanTestColumns(),
		Wire:    &wire.Action{URL: "/api/kanban/move"},
	}))

	utils.AssertContainsAll(t, html,
		`<form`,
		`method="POST"`,
		`action="/api/kanban/move"`,
		`class="hidden"`,
		`type="hidden"`,
	)
}

// TestKanbanEdgeDefaultProps verifies DefaultKanbanBoardProps matches the
// documented defaults.
func TestKanbanEdgeDefaultProps(t *testing.T) {
	t.Parallel()

	props := DefaultKanbanBoardProps()
	if props.CSRFTokenName != "csrf_token" {
		t.Errorf("DefaultKanbanBoardProps CSRFTokenName = %q, want csrf_token", props.CSRFTokenName)
	}

	if props.Wire != nil {
		t.Error("DefaultKanbanBoardProps Wire should be nil (read-only by default)")
	}
}
