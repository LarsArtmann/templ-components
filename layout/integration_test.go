package layout

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestIntegrationFullPageRender(t *testing.T) {
	t.Parallel()

	t.Run("base renders complete HTML document", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		props.Title = testPage
		props.Description = "A test page"
		output := utils.Render(t, Base(props))

		if !strings.HasPrefix(strings.ToLower(output), "<!doctype html>") {
			t.Error("output should start with <!DOCTYPE html>")
		}
		for _, check := range []struct {
			contains string
			message  string
		}{
			{`<html lang="en"`, "html element with lang"},
			{"<title>Test Page</title>", "title"},
			{`content="A test page"`, "description meta"},
			{`<main id="main-content"`, "main element"},
			{"Skip to main content", "skip link"},
		} {
			if !strings.Contains(output, check.contains) {
				t.Errorf("output should contain %s", check.message)
			}
		}
	})

	t.Run("base with security headers renders meta tags", func(t *testing.T) {
		t.Parallel()
		assertSecurityHeadersPresent(t, true)
	})

	t.Run("base without security headers omits meta tags", func(t *testing.T) {
		t.Parallel()
		assertSecurityHeadersPresent(t, false)
	})

	t.Run("base with HTMX SRI renders integrity attributes", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		props.HTMXUseSRI = true
		output := utils.Render(t, Base(props))

		utils.AssertContains(t, output, `integrity="`)
		utils.AssertContains(t, output, `crossorigin="anonymous"`)
	})

	t.Run("base without HTMX SRI omits integrity", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		props.HTMXUseSRI = false
		output := utils.Render(t, Base(props))

		utils.AssertNotContains(t, output, `integrity="`)
	})

	t.Run("base with CSS path renders stylesheet", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		props.CSSPath = "/custom.css"
		output := utils.Render(t, Base(props))

		utils.AssertContains(t, output, `href="/custom.css"`)
		utils.AssertContains(t, output, `rel="stylesheet"`)
	})

	t.Run("base with OG image renders meta", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		props.OGImage = "/og.png"
		output := utils.Render(t, Base(props))

		utils.AssertContains(t, output, `property="og:image" content="/og.png"`)
	})

	t.Run("minimal renders stripped HTML", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Minimal(MinimalProps{Title: "Simple", Locale: "de"}))

		utils.AssertContains(t, strings.ToLower(output), "<!doctype html>")
		utils.AssertContains(t, output, `<html lang="de">`)
		utils.AssertContains(t, output, "<title>Simple</title>")
		utils.AssertNotContains(t, output, "htmx")
	})
}
