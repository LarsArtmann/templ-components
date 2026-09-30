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

	t.Run("zero variant equals truncated variant", func(t *testing.T) {
		t.Parallel()

		explicit := utils.Render(t, ListNote(ListNoteProps{Shown: 50, Total: 127, Variant: ListNoteTruncated}))

		implicit := utils.Render(t, ListNote(ListNoteProps{Shown: 50, Total: 127}))
		if explicit != implicit {
			t.Fatalf("zero Variant must behave as ListNoteTruncated:\nexplicit: %s\nimplicit: %s", explicit, implicit)
		}
	})

	t.Run("unknown variant degrades to truncated semantics", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ListNote(ListNoteProps{Shown: 1, Total: 2, Variant: ListNoteVariant("bogus")}))
		utils.AssertContains(t, output, "Showing 1 of 2")
	})

	t.Run("count variant renders count-only message", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ListNote(ListNoteProps{Shown: 42, Variant: ListNoteCount}))
		utils.AssertContains(t, output, "Showing 42 items.")
		utils.AssertNotContains(t, output, "Narrow your search")
	})

	t.Run("count variant pluralizes singular", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ListNote(ListNoteProps{Shown: 1, Variant: ListNoteCount}))
		utils.AssertContains(t, output, "Showing 1 item.")
	})

	t.Run("count variant renders zero for an empty-but-loaded range", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ListNote(ListNoteProps{Shown: 0, Variant: ListNoteCount}))
		utils.AssertContains(t, output, "Showing 0 items.")
	})

	t.Run("count variant carries role=status for screen readers", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ListNote(ListNoteProps{Shown: 3, Variant: ListNoteCount}))
		utils.AssertContains(t, output, `role="status"`)
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

	t.Run("range variant renders the X–Y of Z position hint", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ListNote(ListNoteProps{
			RangeFrom: 193,
			RangeTo:   216,
			Total:     237,
			Variant:   ListNoteRange,
		}))
		utils.AssertContains(t, output, "Showing 193–216 of 237.")
		utils.AssertNotContains(t, output, "Narrow your search")
	})

	t.Run("range variant renders on the last page too", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ListNote(ListNoteProps{
			RangeFrom: 1,
			RangeTo:   237,
			Total:     237,
			Variant:   ListNoteRange,
		}))
		utils.AssertContains(t, output, "Showing 1–237 of 237.")
	})

	t.Run("range variant degrades to the count message for an empty page", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ListNote(ListNoteProps{
			RangeFrom: 0,
			RangeTo:   0,
			Shown:     0,
			Total:     237,
			Variant:   ListNoteRange,
		}))
		utils.AssertContains(t, output, "Showing 0 items.")
		utils.AssertNotContains(t, output, "–")
	})

	t.Run("range variant carries role=status for screen readers", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ListNote(ListNoteProps{
			RangeFrom: 1,
			RangeTo:   10,
			Total:     20,
			Variant:   ListNoteRange,
		}))
		utils.AssertContains(t, output, `role="status"`)
	})
}
