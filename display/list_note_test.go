package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestListNoteRender(t *testing.T) {
	t.Parallel()

	t.Run("renders notice when total exceeds shown", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ListNote(ListNoteProps{Shown: 50, Total: 127}))
		utils.AssertContains(t, output, "Showing 50 of 127")
		utils.AssertContains(t, output, "Narrow your search")
	})

	t.Run("renders nothing when all items fit", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ListNote(ListNoteProps{Shown: 100, Total: 100}))
		utils.AssertNotContains(t, output, "Showing")
	})

	t.Run("renders nothing when shown exceeds total", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ListNote(ListNoteProps{Shown: 10, Total: 5}))
		utils.AssertNotContains(t, output, "Showing")
	})

	t.Run("zero total renders nothing", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ListNote(ListNoteProps{Shown: 0, Total: 0}))
		utils.AssertNotContains(t, output, "Showing")
	})

	t.Run("has role=status for screen readers", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ListNote(ListNoteProps{Shown: 1, Total: 2}))
		utils.AssertContains(t, output, `role="status"`)
	})
}
