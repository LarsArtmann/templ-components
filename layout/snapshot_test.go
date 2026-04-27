package layout

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestThemeScriptRender(t *testing.T) {
	output := utils.Render(t, ThemeScript("abc123"))
	utils.AssertContains(t, output, `nonce="abc123"`)
	utils.AssertContains(t, output, "localStorage")
	utils.AssertContains(t, output, "classList.add('dark')")
}

func TestThemeToggleRender(t *testing.T) {
	output := utils.Render(t, ThemeToggle("Toggle theme", "nonce-xyz"))
	utils.AssertContains(t, output, `aria-label="Toggle theme"`)
	utils.AssertContains(t, output, `nonce="nonce-xyz"`)
	utils.AssertContains(t, output, "tcToggleTheme")
}

func TestBaseRender(t *testing.T) {
	props := BaseProps{
		Title:    "Test Page",
		Locale:   "en",
		CSSPath:  "/app.css",
		Nonce:    "test-nonce",
		BodyClass: "bg-white",
	}
	output := utils.Render(t, Base(props))
	utils.AssertContains(t, output, "<title>Test Page</title>")
	utils.AssertContains(t, output, `lang="en"`)
	utils.AssertContains(t, output, `nonce="test-nonce"`)
	utils.AssertContains(t, output, `href="/app.css"`)
}

func TestMinimalRender(t *testing.T) {
	output := utils.Render(t, Minimal("Simple", "en"))
	utils.AssertContains(t, output, "<title>Simple</title>")
	utils.AssertContains(t, output, `lang="en"`)
	utils.AssertNotContains(t, output, "htmx")
}
