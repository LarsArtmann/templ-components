package display

import (
	"testing"

	"github.com/a-h/templ"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
	"github.com/larsartmann/templ-components/utils/wire"
)

// Golden sweep for the Kanban board: the wired board under both transport
// dialects, the read-only board, and CSRF/empty-column variations. EnsureID
// normalization keeps the auto-generated board id stable.

func TestGoldenSweepKanban(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "kanban_readonly", HTML: utils.Render(t, KanbanBoard(KanbanBoardProps{
			Columns: kanbanTestColumns(),
		}))},
		{Name: "kanban_wired_htmx", HTML: utils.Render(t, KanbanBoard(KanbanBoardProps{
			Columns:   kanbanTestColumns(),
			Wire:      &wire.Action{URL: "/api/kanban/move"},
			CSRFToken: "tok-1",
		}))},
		{Name: "kanban_wired_datastar", HTML: utils.Render(t, KanbanBoard(KanbanBoardProps{
			Columns: kanbanTestColumns(),
			Wire:    &wire.Action{Transport: wire.TransportDatastar, URL: "/api/kanban/move"},
		}))},
		{Name: "kanban_empty_board", HTML: utils.Render(t, KanbanBoard(DefaultKanbanBoardProps()))},
		{Name: "kanban_card_content", HTML: utils.Render(t, KanbanBoard(KanbanBoardProps{
			Columns: []KanbanColumn{{
				ID:    "todo",
				Title: "To do",
				Cards: []KanbanCard{{
					ID:      "c1",
					Title:   "With body",
					Content: templ.Raw(`<p class="mt-1 text-xs text-gray-500">meta line</p>`),
				}},
			}},
		}))},
	})
}
