// Package layout provides tests for layout components like Base, Minimal, ThemeScript, and ThemeToggle.
package layout

import (
	"testing"

	"github.com/a-h/templ"
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
	utils.AssertContains(t, output, "data-theme-toggle")
	utils.AssertNotContains(t, output, "onclick=")
}

func TestBaseRender(t *testing.T) {
	t.Parallel()
	props := PageProps{
		Title:          testPage,
		Description:    "",
		Locale:         "en",
		OGImage:        "",
		ThemeColor:     "",
		DarkThemeColor: "",
		CSSPath:        testCSSPath,
		Favicon:        "",
		HTMXVersion:    "",
		HTMXUseSRI:     false,
		BodyClass:      "bg-white",
		Nonce:          testNonce,
		HeadContent:    nil,
	}
	output := utils.Render(t, Base(props))
	utils.AssertContains(t, output, "<title>Test Page</title>")
	utils.AssertContains(t, output, `lang="en"`)
	utils.AssertContains(t, output, `nonce="test-nonce"`)
	utils.AssertContains(t, output, `href="/app.css"`)
	utils.AssertContains(t, output, `id="main-content"`)
	utils.AssertContains(t, output, "Skip to main content")
}

func TestBaseRenderFullProps(t *testing.T) {
	t.Parallel()
	props := PageProps{
		Title:          "Full Page",
		Description:    "A test page",
		Locale:         "de",
		OGImage:        "/og.png",
		ThemeColor:     "#ffffff",
		DarkThemeColor: "#000000",
		CSSPath:        "/style.css",
		Favicon:        "/fav.ico",
		HTMXVersion:    htmxVersion206,
		HTMXUseSRI:     true,
		BodyClass:      "bg-gray-50",
		Nonce:          "abc",
		HeadContent:    nil,
	}
	output := utils.Render(t, Base(props))
	utils.AssertContains(t, output, "<title>Full Page</title>")
	utils.AssertContains(t, output, `name="description" content="A test page"`)
	utils.AssertContains(t, output, `property="og:image" content="/og.png"`)
	utils.AssertContains(t, output, `name="theme-color" content="#ffffff"`)
	utils.AssertContains(t, output, `name="theme-color" content="#000000"`)
	utils.AssertContains(t, output, `href="/style.css"`)
	utils.AssertContains(t, output, `href="/fav.ico"`)
	utils.AssertContains(t, output, `nonce="abc"`)
	utils.AssertContains(t, output, "htmx.org@2.0.6")
	utils.AssertContains(t, output, "response-targets.js")
}

func TestMinimalRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Minimal(MinimalProps{Title: "Simple", Locale: "en"}))
	utils.AssertContains(t, output, "<title>Simple</title>")
	utils.AssertNotContains(t, output, "htmx")
}

func TestBaseWithFooter(t *testing.T) {
	t.Parallel()
	props := PageProps{
		Title:     testPage,
		BodyClass: "bg-white",
		Footer:    footerTestComponent(),
	}
	output := utils.Render(t, Base(props))
	utils.AssertContains(t, output, "<footer")
	utils.AssertContains(t, output, "Test Footer")
}

func footerTestComponent() templ.Component {
	return templ.Raw(`<footer>Test Footer</footer>`)
}

func TestDefaultMinimalProps(t *testing.T) {
	t.Parallel()
	props := DefaultMinimalProps()
	if props.Locale != "en" {
		t.Errorf("DefaultMinimalProps().Locale = %q, want %q", props.Locale, "en")
	}
}
