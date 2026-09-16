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
		`<h3 id="kb-a11y-col-todo-title"`,
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
// cards get no "previous" button and middle cards get both — keyboard
// affordances match the possible moves (the last column is empty here, so
// the middle column's card is the deepest with buttons).
func TestKanbanA11yKeyboardButtonsOnlyWhereMovesExist(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, KanbanBoard(KanbanBoardProps{
		Columns: kanbanTestColumns(),
		Wire:    &wire.Action{URL: "/api/kanban/move"},
	}))

	firstCard := html[mustIndex(t, html, `data-tc-kanban-card="c1"`):mustIndex(t, html, `data-tc-kanban-card="c2"`)]
	if !strings.Contains(firstCard, `data-tc-kanban-move="next"`) {
		t.Error("first-column card lacks a next button")
	}

	if strings.Contains(firstCard, `data-tc-kanban-move="prev"`) {
		t.Error("first-column card has a prev button but has no previous column")
	}

	middleCard := html[mustIndex(t, html, `data-tc-kanban-card="c3"`):]
	if !strings.Contains(middleCard, `data-tc-kanban-move="next"`) {
		t.Error("middle-column card lacks a next button")
	}

	if !strings.Contains(middleCard, `data-tc-kanban-move="prev"`) {
		t.Error("middle-column card lacks a prev button")
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

// TestKanbanA11yActionSlot verifies the per-column Action slot lands inside
// the column header (after the count, before the card list) and renders on
// read-only boards too — it is presentation, not part of the move exchange.
func TestKanbanA11yActionSlot(t *testing.T) {
	t.Parallel()

	props := KanbanBoardProps{
		Columns: []KanbanColumn{{
			ID:    "todo",
			Title: "To do",
			Cards: []KanbanCard{{ID: "c1", Title: "Write docs"}},
			Action: Button(ButtonProps{
				Text:    "Add card",
				Variant: ButtonSecondary,
				Size:    ButtonSizeSM,
			}),
		}},
	}

	wired := props
	wired.Wire = &wire.Action{URL: "/api/kanban/move"}
	wiredHTML := utils.Render(t, KanbanBoard(wired))
	actionIdx := mustIndex(t, wiredHTML, ">Add card</button>")
	headerIdx := mustIndex(t, wiredHTML, ">To do</h3>")
	bodyIdx := mustIndex(t, wiredHTML, `data-tc-kanban-column-body="todo"`)

	if actionIdx < headerIdx || actionIdx > bodyIdx {
		t.Errorf(
			"action button not between the column header and the card list (action %d, header %d, body %d)",
			actionIdx, headerIdx, bodyIdx,
		)
	}

	utils.AssertContains(t, wiredHTML, `aria-label="To do: 1 card"`)
	utils.AssertContains(t, wiredHTML, ">Add card</button>")

	readonly := utils.Render(t, KanbanBoard(props))
	utils.AssertContains(t, readonly, ">Add card</button>")
}

// TestKanbanToneDotRendering verifies the column status dot: rendered
// before the title when Tone is set (aria-hidden, decorative), and absent
// for the zero value and unknown tones (graceful degradation — no invisible
// element).
func TestKanbanToneDotRendering(t *testing.T) {
	t.Parallel()

	toned := utils.Render(t, KanbanBoard(KanbanBoardProps{
		Columns: []KanbanColumn{{ID: "done", Title: "Done", Tone: KanbanToneGreen}},
	}))
	utils.AssertContainsAll(t, toned,
		"rounded-full", "h-2", "w-2", "bg-green-500", "dark:bg-green-400", `aria-hidden="true"`,
	)

	for name, tone := range map[string]KanbanTone{"zero": "", "unknown": KanbanTone("bogus")} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			html := utils.Render(t, KanbanBoard(KanbanBoardProps{
				Columns: []KanbanColumn{{ID: "done", Title: "Done", Tone: tone}},
			}))
			utils.AssertNotContains(t, html, "rounded-full")
		})
	}
}

// mustIndex returns strings.Index(s, substr) or fails the test when the
// substring is absent (a -1 index would slice out of range).
func mustIndex(t *testing.T, s, substr string) int {
	t.Helper()

	idx := strings.Index(s, substr)
	if idx < 0 {
		t.Fatalf("substring %q not found", substr)
	}

	return idx
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

// TestKanbanTouchButtonsHook verifies the move buttons live under the
// tc-kanban-buttons hook — the selector templates/custom.css targets with a
// pointer:coarse rule so touch devices (no hover) still see them. Browser
// proof of the computed style lives in visualtest/kanban_e2e_test.go.
func TestKanbanTouchButtonsHook(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, KanbanBoard(KanbanBoardProps{
		Columns: kanbanTestColumns(),
		Wire:    &wire.Action{URL: "/api/kanban/move"},
	}))
	utils.AssertContainsAll(t, html,
		`class="tc-kanban-buttons`,
		`opacity-0 group-hover:opacity-100 group-focus-within:opacity-100`,
	)

	// Read-only boards render no button wrapper at all.
	readonly := utils.Render(t, KanbanBoard(KanbanBoardProps{Columns: kanbanTestColumns()}))
	utils.AssertNotContains(t, readonly, "tc-kanban-buttons")
}
