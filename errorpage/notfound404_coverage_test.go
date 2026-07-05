package errorpage

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestNotFound404Coverage(t *testing.T) {
	t.Parallel()

	t.Run("numeral defaults to 404 when empty", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{Numeral: ""}))
		utils.AssertContains(t, output, "404")
	})

	t.Run("title defaults to Page not found when empty", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{Title: ""}))
		utils.AssertContains(t, output, "Page not found")
	})

	t.Run("go home text defaults when empty", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{GoHomeHref: "/", GoHomeText: ""}))
		utils.AssertContains(t, output, "Go to homepage")
	})

	t.Run("attrs propagate to root element", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{
			BaseProps: utils.BaseProps{Attrs: templ.Attributes{"data-testid": "not-found-page"}},
		}))
		utils.AssertContains(t, output, `data-testid="not-found-page"`)
	})

	t.Run("search placeholder defaults when empty", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{
			SearchAction: "/search", SearchPlaceholder: "",
		}))
		utils.AssertContains(t, output, `placeholder="Search..."`)
		utils.AssertContains(t, output, `aria-label="Search..."`)
	})

	t.Run("go home with default href from DefaultNotFound404Props", func(t *testing.T) {
		t.Parallel()
		props := DefaultNotFound404Props()
		output := utils.Render(t, NotFound404(props))
		utils.AssertContains(t, output, `href="/"`)
		utils.AssertContains(t, output, "Go to homepage")
	})
}

func TestNotFound404DefaultPropsValues(t *testing.T) {
	t.Parallel()
	props := DefaultNotFound404Props()
	if props.Numeral != "404" {
		t.Errorf("Numeral = %q, want %q", props.Numeral, "404")
	}
	if props.Title != "Page not found" {
		t.Errorf("Title = %q, want %q", props.Title, "Page not found")
	}
	if props.Message == "" {
		t.Error("Message should not be empty")
	}
	if props.SearchInputName != "q" {
		t.Errorf("SearchInputName = %q, want %q", props.SearchInputName, "q")
	}
	if props.GoHomeHref != "/" {
		t.Errorf("GoHomeHref = %q, want %q", props.GoHomeHref, "/")
	}
	if props.GoHomeText != "Go to homepage" {
		t.Errorf("GoHomeText = %q, want %q", props.GoHomeText, "Go to homepage")
	}
	if !props.ShowGoBack {
		t.Error("ShowGoBack should be true")
	}
}
