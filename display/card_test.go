package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestCardRender(t *testing.T) {
	t.Run("basic card with title", func(t *testing.T) {
		output := utils.Render(t, Card(CardProps{Title: "Users", Padding: "md"}))
		utils.AssertContains(t, output, "Users")
		utils.AssertContains(t, output, "bg-white")
		utils.AssertContains(t, output, "rounded-lg")
	})

	t.Run("card with custom class and id", func(t *testing.T) {
		props := CardProps{
			BaseProps: utils.BaseProps{ID: "my-card", Class: "mt-4"},
			Title:     "Test",
			Padding:   "md",
		}
		output := utils.Render(t, Card(props))
		utils.AssertContains(t, output, `id="my-card"`)
		utils.AssertContains(t, output, "mt-4")
	})

	t.Run("simple card", func(t *testing.T) {
		output := utils.Render(t, SimpleCard())
		utils.AssertContains(t, output, "bg-white")
		utils.AssertContains(t, output, "rounded-lg")
	})
}

func TestEmptyStateRender(t *testing.T) {
	t.Run("with action link", func(t *testing.T) {
		props := EmptyStateProps{
			Title:      "No repos",
			Icon:       "folder",
			ActionText: "Add Repo",
			ActionHref: "/repos/new",
		}
		output := utils.Render(t, EmptyState(props))
		utils.AssertContains(t, output, "No repos")
		utils.AssertContains(t, output, "Add Repo")
		utils.AssertContains(t, output, `href="/repos/new"`)
	})

	t.Run("simple empty state", func(t *testing.T) {
		output := utils.Render(t, SimpleEmptyState("Nothing here"))
		utils.AssertContains(t, output, "Nothing here")
	})
}
