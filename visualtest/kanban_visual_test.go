package visualtest_test

import (
	"context"
	"io"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
	"github.com/larsartmann/templ-components/visualtest"
)

// kanbanBoardComponent renders a wired three-column board — the richest
// visual state: populated + empty columns, card counts, hidden move buttons,
// drag affordances — for pixel-level regression coverage. The move pipeline
// itself is browser-proven by kanban_e2e_test.go; these goldens pin layout
// and colors only.
//
// #80-family caveat: the PNGs below are agent-captured; a human should
// eyeball them once (column proportions, dark-mode surfaces, RTL mirroring).
func kanbanBoardComponent() templ.Component {
	props := display.DefaultKanbanBoardProps()
	props.BaseProps = utils.BaseProps{ID: "kb-visual"}
	props.Columns = []display.KanbanColumn{
		{
			ID:    "todo",
			Title: "To do",
			Cards: []display.KanbanCard{
				{ID: "c1", Title: "Draft release notes"},
				{ID: "c2", Title: "Review open PRs"},
			},
		},
		{ID: "doing", Title: "In progress", Cards: []display.KanbanCard{{ID: "c3", Title: "Kanban polish"}}},
		{ID: "done", Title: "Done", Cards: []display.KanbanCard{{ID: "c4", Title: "Pareto plan"}}},
	}
	props.Wire = &wire.Action{URL: "/api/kanban/visual"}

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<div class="w-[64rem]">`); err != nil {
			return err
		}

		if err := display.KanbanBoard(props).Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w, `</div>`)

		return err
	})
}

// TestKanbanBoard pins the wired board's presentation in light, dark, and
// RTL. The empty-column text variant is covered by the HTML goldens; dark
// column surfaces and RTL mirroring are only provable at the pixel level.
func TestKanbanBoard(t *testing.T) {
	t.Parallel()

	visualtest.AssertScreenshot(t, "kanban/board_light", kanbanBoardComponent())
	visualtest.AssertScreenshot(
		t,
		"kanban/board_dark",
		kanbanBoardComponent(),
		visualtest.Options{Dark: new(true)},
	)
	visualtest.AssertScreenshot(
		t,
		"kanban/board_rtl",
		kanbanBoardComponent(),
		visualtest.Options{RTL: new(true)},
	)
}
