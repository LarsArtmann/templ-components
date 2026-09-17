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

// kanbanSectionAddButton mirrors the demo's per-column add affordance.
func kanbanSectionAddButton(url, columnTitle string) templ.Component {
	return display.Button(display.ButtonProps{
		BaseProps: utils.BaseProps{AriaLabel: "Add card to " + columnTitle},
		Text:      "+ Add",
		Variant:   display.ButtonGhost,
		Size:      display.ButtonSizeSM,
		Wire: &wire.Action{
			Method: wire.MethodPost,
			Event:  wire.EventClick,
			URL:    url,
			Target: "#kb-visual-action",
		},
	})
}

// kanbanSectionComponent mirrors the demo's kanban section — the visual
// surface the Action/Tone phase added: a per-board Reset button in a header
// row, columns with Tone status dots, and a per-column Action add-card
// button. Three columns so every affordance sits inside the capture frame
// (the board's column row scrolls horizontally by design); the six tone
// variants are pinned by the kanban_column_tone HTML golden. Behavior is
// browser-proven by kanban_e2e_test.go; these goldens pin the new pixels,
// which no below-the-fold route golden ever covered.
func kanbanSectionComponent() templ.Component {
	board := display.DefaultKanbanBoardProps()
	board.BaseProps = utils.BaseProps{ID: "kb-visual-action"}
	board.Columns = []display.KanbanColumn{
		{
			ID:    "backlog",
			Title: "Backlog",
			Tone:  display.KanbanToneGray,
			Cards: []display.KanbanCard{
				{ID: "c1", Title: "Audit CSP headers"},
				{ID: "c2", Title: "Dark mode pass"},
			},
			Action: kanbanSectionAddButton("/api/kanban/visual/add/backlog", "Backlog"),
		},
		{
			ID:     "progress",
			Title:  "In progress",
			Tone:   display.KanbanToneBlue,
			Action: kanbanSectionAddButton("/api/kanban/visual/add/progress", "In progress"),
		},
		{
			ID:     "done",
			Title:  "Done",
			Tone:   display.KanbanToneGreen,
			Cards:  []display.KanbanCard{{ID: "c3", Title: "Set up CI"}},
			Action: kanbanSectionAddButton("/api/kanban/visual/add/done", "Done"),
		},
	}
	board.Wire = &wire.Action{URL: "/api/kanban/visual"}

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(
			w,
			`<div class="w-[64rem]"><div class="mb-2 flex items-center justify-between">`,
		); err != nil {
			return err
		}

		reset := display.Button(display.ButtonProps{
			BaseProps: utils.BaseProps{AriaLabel: "Reset demo board"},
			Text:      "Reset",
			Variant:   display.ButtonGhost,
			Size:      display.ButtonSizeSM,
			Wire: &wire.Action{
				Method: wire.MethodPost,
				URL:    "/api/kanban/visual/reset",
				Target: "#kb-visual-action",
			},
		})
		if err := reset.Render(ctx, w); err != nil {
			return err
		}

		if _, err := io.WriteString(w, `</div>`); err != nil {
			return err
		}

		if err := display.KanbanBoard(board).Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w, `</div>`)

		return err
	})
}

// TestKanbanSection pins the Action/Tone surface in both modes: tone-dot
// colors, the Action button in the column header, and the board-header
// reset affordance.
func TestKanbanSection(t *testing.T) {
	t.Parallel()

	visualtest.AssertScreenshot(t, "kanban/section_action_tone_light", kanbanSectionComponent())
	visualtest.AssertScreenshot(
		t,
		"kanban/section_action_tone_dark",
		kanbanSectionComponent(),
		visualtest.Options{Dark: new(true)},
	)
}
