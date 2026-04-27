// Package layout provides tests for layout components like Base, Minimal, ThemeScript, and ThemeToggle.
package layout

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestThemeScriptRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ThemeScript("abc123"))
	utils.AssertContains(t, output, `nonce="abc123"`)
	utils.AssertContains(t, output, "localStorage")
	utils.AssertContains(t, output, "classList.add('dark')")
}

func TestThemeToggleRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ThemeToggle("Toggle theme", "nonce-xyz"))
	utils.AssertContains(t, output, `aria-label="Toggle theme"`)
	utils.AssertContains(t, output, `nonce="nonce-xyz"`)
	utils.AssertContains(t, output, "tcToggleTheme")
}

func TestBaseRender(t *testing.T) {
	t.Parallel()
	props := BaseProps{
		Title:          "Test Page",
		Description:    "",
		Locale:         "en",
		OGImage:        "",
		ThemeColor:     "",
		DarkThemeColor: "",
		CSSPath:        "/app.css",
		Favicon:        "",
		HTMXVersion:    "",
		HTMXSRI:        "",
		BodyClass:      "bg-white",
		Nonce:          "test-nonce",
		HeadContent:    nil,
	}
	output := utils.Render(t, Base(props))
	utils.AssertContains(t, output, "<title>Test Page</title>")
	utils.AssertContains(t, output, `lang="en"`)
	utils.AssertContains(t, output, `nonce="test-nonce"`)
	utils.AssertContains(t, output, `href="/app.css"`)
}

func TestMinimalRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Minimal("Simple", "en"))
	utils.AssertContains(t, output, "<title>Simple</title>")
	utils.AssertContains(t, output, `lang="en"`)
	utils.AssertNotContains(t, output, "htmx")
}
