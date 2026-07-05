package errorpage

import (
	"strconv"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestNotFound404EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("empty props renders without panic", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{}))
		utils.AssertContains(t, output, "min-h-screen")
		utils.AssertContains(t, output, "404")
		utils.AssertContains(t, output, "Page not found")
	})

	t.Run("no search action omits search form", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{}))
		utils.AssertNotContains(t, output, `role="search"`)
		utils.AssertNotContains(t, output, "<form")
	})

	t.Run("no links omits links section", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{}))
		utils.AssertNotContains(t, output, "Popular pages")
	})

	t.Run("no go home href omits go home button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{GoHomeHref: ""}))
		utils.AssertNotContains(t, output, "Go to homepage")
	})

	t.Run("show go back false omits go back button and script", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{ShowGoBack: false}))
		utils.AssertNotContains(t, output, "Go back")
		utils.AssertNotContains(t, output, "data-tc-go-back")
		utils.AssertNotContains(t, output, "<script")
	})

	t.Run("empty message omits message paragraph", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{Message: ""}))
		utils.AssertNotContains(t, output, "doesn't exist")
	})

	t.Run("empty link icon renders without icon span", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{
			Links: []NotFoundLink{{Text: "No Icon", Href: "/no-icon"}},
		}))
		utils.AssertContains(t, output, "No Icon")
		utils.AssertContains(t, output, "/no-icon")
	})

	t.Run("custom numeral renders correctly", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{Numeral: "ERROR"}))
		utils.AssertContains(t, output, "ERROR")
	})

	t.Run("search with default input name when empty", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{SearchAction: "/search"}))
		utils.AssertContains(t, output, `name="q"`)
		utils.AssertContains(t, output, `placeholder="Search..."`)
	})

	t.Run("many links render in grid", func(t *testing.T) {
		t.Parallel()
		links := make([]NotFoundLink, 6)
		for i := range links {
			links[i] = NotFoundLink{Text: "Link " + strconv.Itoa(i+1), Href: "/" + strconv.Itoa(i+1)}
		}
		output := utils.Render(t, NotFound404(NotFound404Props{Links: links}))
		for i := 1; i <= 6; i++ {
			utils.AssertContains(t, output, "Link "+strconv.Itoa(i))
		}
	})

	t.Run("custom ID and class propagate", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{
			BaseProps: utils.BaseProps{ID: "my-404", Class: "custom-bg"},
		}))
		utils.AssertContains(t, output, `id="my-404"`)
		utils.AssertContains(t, output, "custom-bg")
	})
}

func TestNotFound404_SearchInputNameDefault(t *testing.T) {
	t.Parallel()
	t.Run("explicit input name overrides default", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{
			SearchAction: "/search", SearchInputName: "query",
		}))
		utils.AssertContains(t, output, `name="query"`)
		utils.AssertNotContains(t, output, `name="q"`)
	})
}

func TestDefaultNotFoundLinksReturnsValues(t *testing.T) {
	t.Parallel()
	links := DefaultNotFoundLinks()
	if len(links) < 2 {
		t.Fatalf("expected at least 2 default links, got %d", len(links))
	}
	for _, link := range links {
		if link.Text == "" || link.Href == "" {
			t.Errorf("default link has empty Text or Href: %+v", link)
		}
	}
}
