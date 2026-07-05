package errorpage

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestNotFound404A11y(t *testing.T) {
	t.Parallel()

	t.Run("default aria-label when none provided", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(DefaultNotFound404Props()))
		utils.AssertContains(t, output, `aria-label="404 — Page not found"`)
	})

	t.Run("propagates custom aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{
			BaseProps: utils.BaseProps{AriaLabel: "Custom 404 label"},
		}))
		utils.AssertContains(t, output, `aria-label="Custom 404 label"`)
		utils.AssertNotContains(t, output, `aria-label="404 — Page not found"`)
	})

	t.Run("numeral is aria-hidden decorative", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(DefaultNotFound404Props()))
		utils.AssertContains(t, output, `aria-hidden="true"`)
	})

	t.Run("search form has role=search", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{SearchAction: "/search"}))
		utils.AssertContains(t, output, `role="search"`)
	})

	t.Run("search input has aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{
			SearchAction:      "/search",
			SearchPlaceholder: "Search the site...",
		}))
		utils.AssertContains(t, output, `aria-label="Search the site..."`)
	})

	t.Run("go back button has type=button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{ShowGoBack: true}))
		utils.AssertContains(t, output, `type="button"`)
		utils.AssertContains(t, output, `data-tc-go-back`)
	})

	t.Run("has dark mode classes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(DefaultNotFound404Props()))
		utils.AssertContains(t, output, "dark:")
	})

	t.Run("interactive elements have motion-reduce", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{
			SearchAction: "/search", ShowGoBack: true, Links: DefaultNotFoundLinks(),
		}))
		utils.AssertContains(t, output, "motion-reduce:transition-none")
		utils.AssertContains(t, output, "motion-reduce:duration-0")
	})

	t.Run("go home link has focus-visible ring", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{GoHomeHref: "/"}))
		utils.AssertContains(t, output, "focus-visible:ring-blue-500")
	})

	t.Run("link cards have focus-visible ring", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{
			Links: []NotFoundLink{{Text: "Home", Href: "/"}},
		}))
		utils.AssertContains(t, output, "focus-visible:ring-blue-500")
	})

	t.Run("nonce on script tag when ShowGoBack", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{
			ShowGoBack: true, BaseProps: utils.BaseProps{Nonce: "test-nonce-xyz"},
		}))
		utils.AssertContains(t, output, `nonce="test-nonce-xyz"`)
	})

	t.Run("no script tag when ShowGoBack is false", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, NotFound404(NotFound404Props{ShowGoBack: false}))
		utils.AssertNotContains(t, output, "<script")
	})
}
