package layout

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestBaseFullCoverage(t *testing.T) {
	t.Parallel()
	t.Run("with all props", func(t *testing.T) {
		t.Parallel()
		headContent := templ.Raw(`<meta name="custom" content="test">`)
		footer := templ.Raw(`<script>console.log("footer")</script>`)
		output := utils.Render(t, Base(PageProps{
			Title:          "My Page",
			Description:    "A test page",
			Locale:         "en",
			CSSPath:        "/css/app.css",
			Favicon:        "/favicon.ico",
			ThemeColor:     "#ffffff",
			DarkThemeColor: "#0a0a0a",
			BodyClass:      "bg-gray-50",
			Nonce:          "nonce-test",
			HeadContent:    headContent,
			Footer:         footer,
		}))
		utils.AssertContains(t, output, "My Page")
		utils.AssertContains(t, output, "A test page")
		utils.AssertContains(t, output, `/css/app.css`)
		utils.AssertContains(t, output, `content="test"`)
	})
}

func TestMinimalFullCoverage(t *testing.T) {
	t.Parallel()
	t.Run("with title and locale", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Minimal(MinimalProps{
			Title:  "Minimal Page",
			Locale: "fr",
		}))
		utils.AssertContains(t, output, "Minimal Page")
		utils.AssertContains(t, output, `lang="fr"`)
	})
}

func TestThemeScriptFullCoverage(t *testing.T) {
	t.Parallel()
	t.Run("renders with nonce", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ThemeScript("nonce-ts"))
		utils.AssertContains(t, output, `nonce="nonce-ts"`)
		utils.AssertContains(t, output, "dark")
	})
}

func TestThemeToggleFullCoverage(t *testing.T) {
	t.Parallel()
	t.Run("renders with nonce", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ThemeToggle("Toggle theme", "nonce-tt"))
		utils.AssertContains(t, output, `aria-label="Toggle theme"`)
		utils.AssertContains(t, output, `nonce="nonce-tt"`)
	})
}
