package display_test

import (
	"bytes"
	"context"

	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/utils/wire"
)

func ExampleKanbanBoard() {
	props := display.DefaultKanbanBoardProps()
	props.Columns = []display.KanbanColumn{
		{
			ID:    "todo",
			Title: "To do",
			Cards: []display.KanbanCard{
				{ID: "c1", Title: "Write docs"},
				{ID: "c2", Title: "Cut release", Href: "/issues/2"},
			},
		},
		{ID: "doing", Title: "In progress"},
		{ID: "done", Title: "Done"},
	}
	props.Wire = &wire.Action{URL: "/api/kanban/move"}

	var buf bytes.Buffer

	_ = display.KanbanBoard(props).Render(context.Background(), &buf)
}

func ExampleKanbanBoard_readOnly() {
	// Without Wire the board renders as a static, read-only view.
	var buf bytes.Buffer

	_ = display.KanbanBoard(display.KanbanBoardProps{
		Columns: []display.KanbanColumn{
			{ID: "todo", Title: "To do", Cards: []display.KanbanCard{{ID: "c1", Title: "Write docs"}}},
		},
	}).Render(context.Background(), &buf)
}
