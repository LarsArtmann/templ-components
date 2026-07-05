package errorpage

import (
	"testing"

	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

func TestNotFound404Golden(t *testing.T) {
	t.Parallel()

	t.Run("full 404 page with all features", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{
			Numeral:           "404",
			Title:             "Page not found",
			Message:           "The page you're looking for doesn't exist or has been moved.",
			SearchAction:      "/search",
			SearchPlaceholder: "Search pages...",
			SearchInputName:   "q",
			GoHomeHref:        "/",
			GoHomeText:        "Go to homepage",
			ShowGoBack:        true,
			Links: []NotFoundLink{
				{Text: "Home", Href: "/", Icon: "home"},
				{Text: "Documentation", Href: "/docs", Icon: "document"},
				{Text: "Blog", Href: "/blog", Icon: "bookmark"},
			},
		}))
		golden.Assert(t, "notfound404_full", output)
	})

	t.Run("minimal 404 page with defaults", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(DefaultNotFound404Props()))
		golden.Assert(t, "notfound404_minimal", output)
	})
}
