package errorpage

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestNotFound404UserSeesFullPage(t *testing.T) {
	t.Parallel()

	t.Run("user sees large 404 numeral as hero", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(DefaultNotFound404Props()))
		utils.AssertContains(t, output, "404")
	})

	t.Run("user sees title and message", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{
			Title: "Oops! Lost in space", Message: "This page drifted into the void.",
		}))
		utils.AssertContains(t, output, "Oops! Lost in space")
		utils.AssertContains(t, output, "This page drifted into the void.")
	})

	t.Run("user sees search box to find content", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{
			SearchAction: "/api/search", SearchPlaceholder: "Search articles...",
		}))
		utils.AssertContains(t, output, `action="/api/search"`)
		utils.AssertContains(t, output, `placeholder="Search articles..."`)
		utils.AssertContains(t, output, `type="search"`)
		utils.AssertContains(t, output, `method="get"`)
	})

	t.Run("user sees go home button linking to homepage", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{
			GoHomeHref: "/dashboard", GoHomeText: "Back to dashboard",
		}))
		utils.AssertContains(t, output, `href="/dashboard"`)
		utils.AssertContains(t, output, "Back to dashboard")
	})

	t.Run("user sees go back button that navigates history", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{ShowGoBack: true}))
		utils.AssertContains(t, output, "Go back")
		utils.AssertContains(t, output, "history.back()")
	})

	t.Run("user sees suggested links as cards", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{
			Links: []NotFoundLink{
				{Text: "Getting Started", Href: "/docs/getting-started"},
				{Text: "API Reference", Href: "/docs/api"},
			},
		}))
		utils.AssertContains(t, output, "Getting Started")
		utils.AssertContains(t, output, "/docs/getting-started")
		utils.AssertContains(t, output, "API Reference")
		utils.AssertContains(t, output, "Popular pages")
	})

	t.Run("user sees custom numeral for different error codes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{
			Numeral: "418", Title: "Teapot Error",
		}))
		utils.AssertContains(t, output, "418")
		utils.AssertContains(t, output, "Teapot Error")
	})

	t.Run("user sees defaults when props are empty", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{}))
		utils.AssertContains(t, output, "404")
		utils.AssertContains(t, output, "Page not found")
	})
}
