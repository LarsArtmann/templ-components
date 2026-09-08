package display

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// TestKanbanA11y verifies the accessibility contract of the board: region
// landmark, labelled columns, list semantics, keyboard move buttons with
// descriptive labels, a polite live region, and nonce-carrying script.
func TestKanbanA11y(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, KanbanBoard(KanbanBoardProps{
		BaseProps: utils.BaseProps{ID: "kb-a11y", Nonce: "n-a11y"},
		Columns:   kanbanTestColumns(),
		Wire:      &wire.Action{URL: "/api/kanban/move"},
	}))

	utils.AssertContainsAll(t, html,
		`role="region"`,
		`aria-label="Kanban board"`,
		`aria-labelledby="kb-a11y-col-todo-title"`,
		`aria-labelledby="kb-a11y-col-doing-title"`,
		`<h3 id="kb-a11y-col-todo-title">`,
		`aria-live="polite"`,
		`role="status"`,
		`aria-label="Move Write docs to next column"`,
		`aria-label="Move Kanban board to previous column"`,
		`aria-label="To do: 2 cards"`,
		`aria-label="Done: no cards"`,
		`nonce="n-a11y"`,
	)
}

// TestKanbanA11yAriaLabelOverride verifies BaseProps.AriaLabel wins over
// the default region label.
func TestKanbanA11yAriaLabelOverride(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, KanbanBoard(KanbanBoardProps{
		BaseProps: utils.BaseProps{AriaLabel: "Sprint board"},
		Columns:   kanbanTestColumns(),
	}))
	utils.AssertContains(t, html, `aria-label="Sprint board"`)
	utils.AssertNotContains(t, html, `aria-label="Kanban board"`)
}

// TestKanbanA11yKeyboardButtonsOnlyWhereMovesExist verifies first-column
// cards get no "previous" button and last-column cards no "next" button —
// keyboard affordances match the possible moves.
func TestKanbanA11yKeyboardButtonsOnlyWhereMovesExist(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, KanbanBoard(KanbanBoardProps{
		Columns: kanbanTestColumns(),
		Wire:    &wire.Action{URL: "/api/kanban/move"},
	}))

	firstCard := html[strings.Index(html, `data-tc-kanban-card="c1"`):strings.Index(html, `data-tc-kanban-card="c2"`)]
	if !strings.Contains(firstCard, `data-tc-kanban-move="next"`) {
		t.Error("first-column card lacks a next button")
	}

	if strings.Contains(firstCard, `data-tc-kanban-move="prev"`) {
		t.Error("first-column card has a prev button but has no previous column")
	}

	lastCard := html[strings.Index(html, `data-tc-kanban-card="c3"`):]
	if strings.Contains(lastCard, `data-tc-kanban-move="next"`) {
		t.Error("last-column card has a next button but has no next column")
	}

	if !strings.Contains(lastCard, `data-tc-kanban-move="prev"`) {
		t.Error("last-column card lacks a prev button")
	}
}

// TestKanbanA11yReadonlyBoardIsSemantic verifies an unwired board keeps the
// region/column/list semantics without interactive affordances.
func TestKanbanA11yReadonlyBoardIsSemantic(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, KanbanBoard(KanbanBoardProps{
		Columns: kanbanTestColumns(),
	}))

	utils.AssertContainsAll(t, html,
		`role="region"`,
		`<ul`,
		`<li`,
	)
	utils.AssertNotContains(t, html, "draggable")
	utils.AssertNotContains(t, html, "data-tc-kanban-move")
	utils.AssertNotContains(t, html, "aria-live")
}

// TestKanbanMotionReduce verifies the reveal transition on the move buttons
// carries a motion-reduce fallback (shared constants guarantee it).
func TestKanbanMotionReduce(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, KanbanBoard(KanbanBoardProps{
		Columns: kanbanTestColumns(),
		Wire:    &wire.Action{URL: "/api/kanban/move"},
	}))
	utils.AssertContains(t, html, "motion-reduce:transition-none")
}
