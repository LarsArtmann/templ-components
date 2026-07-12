package navigation

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestEndOfList(t *testing.T) {
	t.Parallel()

	t.Run("default message", func(t *testing.T) {
		t.Parallel()
		html := utils.Render(t, EndOfList(DefaultEndOfListProps()))
		utils.AssertContains(t, html, "reached the end")
		utils.AssertContains(t, html, `role="status"`)
		utils.AssertContainsAll(t, html, "text-gray-500", "dark:text-gray-400")
	})

	t.Run("custom message", func(t *testing.T) {
		t.Parallel()
		html := utils.Render(t, EndOfList(EndOfListProps{Message: "No more battles"}))
		utils.AssertContains(t, html, "No more battles")
	})

	t.Run("propagates class and id", func(t *testing.T) {
		t.Parallel()
		html := utils.Render(t, EndOfList(EndOfListProps{
			BaseProps: utils.BaseProps{ID: "eol", Class: "border-t pt-4"},
		}))
		utils.AssertContains(t, html, `id="eol"`)
		utils.AssertContains(t, html, "border-t")
		utils.AssertContains(t, html, "pt-4")
	})

	t.Run("empty message falls back to default", func(t *testing.T) {
		t.Parallel()
		html := utils.Render(t, EndOfList(EndOfListProps{}))
		utils.AssertContains(t, html, "reached the end")
	})
}
