package layout

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// ---------------------------------------------------------------------------
// Base: full props coverage — Title, Description, OG, Theme, Security, HeadContent
// ---------------------------------------------------------------------------

func TestBaseFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Base(PageProps{
		Title:           "My Page",
		Description:     "A description",
		Locale:          "fr",
		OGImage:         "https://example.com/og.png",
		ThemeColor:      "#ffffff",
		DarkThemeColor:  "#000000",
		CSSPath:         "/styles.css",
		Favicon:         "/icon.png",
		BodyClass:       "min-h-screen",
		Nonce:           "nonce-abc",
		SecurityHeaders: true,
		HeadContent:     templ.Raw("<meta name=\"custom\" content=\"data\">"),
	}))
	utils.AssertContainsAll(t, output,
		"<title>My Page</title>",
		"A description",
		`lang="fr"`,
		"https://example.com/og.png",
		"/styles.css",
		"/icon.png",
		"min-h-screen",
		`name="custom"`,
	)
}

func TestBaseNoCSSPath(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Base(PageProps{
		Title:   "No CSS",
		CSSPath: "",
	}))
	utils.AssertNotContains(t, output, `rel="stylesheet"`)
}

func TestBaseNoHTMX(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Base(PageProps{
		Title:       "No HTMX",
		HTMXVersion: "", // disable auto-injection
	}))
	utils.AssertNotContains(t, output, "htmx.org")
}

func TestBaseWithHTMXCustomCDN(t *testing.T) {
	t.Parallel()

	props := DefaultPageProps()
	props.Title = "Custom CDN"
	props.HTMXCDN = "https://unpkg.com"
	output := utils.Render(t, Base(props))
	utils.AssertContains(t, output, "unpkg.com")
}

func TestBaseWithHTMXNoResponseTargets(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Base(PageProps{
		Title:               "No Response Targets",
		HTMXResponseTargets: false,
	}))
	utils.AssertNotContains(t, output, "response-targets")
}

func TestBaseWithHTMXSRI(t *testing.T) {
	t.Parallel()

	props := DefaultPageProps()
	props.Title = "SRI"
	props.HTMXUseSRI = true
	output := utils.Render(t, Base(props))
	utils.AssertContains(t, output, "sha384-")
}

func TestBaseWithFooterSlot(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Base(PageProps{
		Title:  "With Footer",
		Footer: templ.Raw("<footer>Footer content</footer>"),
	}))
	utils.AssertContains(t, output, "Footer content")
}

func TestBaseDefaultProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Base(DefaultPageProps()))
	utils.AssertContainsAll(t, output,
		`lang="en"`,
		"/app.css",
		"htmx.org",
	)
}

func TestBaseWithSecurityHeadersDisabled(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Base(PageProps{
		Title:           "No Security",
		SecurityHeaders: false,
	}))
	// Without security headers, CSP meta tag should be absent
	utils.AssertNotContains(t, output, "Content-Security-Policy")
}

// ---------------------------------------------------------------------------
// Minimal: locale fallback, title, children
// ---------------------------------------------------------------------------

func TestMinimalFullProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Minimal(MinimalProps{
		Title:  "Minimal Page",
		Locale: "de",
	}))
	utils.AssertContainsAll(t, output,
		"<!doctype html>",
		`lang="de"`,
		"<title>Minimal Page</title>",
	)
}

func TestMinimalDefaultLocale(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Minimal(MinimalProps{
		Title: "Default",
	}))
	utils.AssertContains(t, output, `lang="en"`)
}

func TestMinimalDefaultProps(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Minimal(DefaultMinimalProps()))
	utils.AssertContainsAll(t, output, "<!doctype html>", `lang="en"`)
}

func TestMinimalEmptyTitle(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Minimal(MinimalProps{}))
	utils.AssertContains(t, output, "<title></title>")
}

// ---------------------------------------------------------------------------
// ThemeScript
// ---------------------------------------------------------------------------

func TestThemeScriptWithNonce(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ThemeScript("nonce-123"))
	utils.AssertContainsAll(t, output,
		`nonce="nonce-123"`,
		"<script",
		"dark",
	)
}

func TestThemeScriptEmptyNonce(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ThemeScript(""))
	utils.AssertContains(t, output, "<script")
}

// ---------------------------------------------------------------------------
// ThemeToggle
// ---------------------------------------------------------------------------

func TestThemeToggleWithAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ThemeToggle("Toggle dark mode", "nonce-abc"))
	utils.AssertContainsAll(t, output,
		`aria-label="Toggle dark mode"`,
		`nonce="nonce-abc"`,
	)
}

func TestThemeToggleDefaultAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ThemeToggle("", "n"))
	utils.AssertContains(t, output, "<button")
}

// ---------------------------------------------------------------------------
// Script helper
// ---------------------------------------------------------------------------

func TestScriptWithNonce(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Script("nonce-xyz", "/app.js", nil))
	utils.AssertContainsAll(t, output,
		`src="/app.js"`,
		`nonce="nonce-xyz"`,
	)
}

func TestScriptWithAttrs(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Script("n", "/lib.js", templ.Attributes{
		"defer":   true,
		"async":   true,
		"data-cf": "none",
	}))
	utils.AssertContainsAll(t, output,
		`src="/lib.js"`,
		"defer",
		"async",
		`data-cf="none"`,
	)
}

func TestScriptEmptySrc(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Script("n", "", nil))
	utils.AssertContains(t, output, "<script")
}

// ---------------------------------------------------------------------------
// Stylesheet helper
// ---------------------------------------------------------------------------

func TestStylesheetBasic(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Stylesheet("/styles.css", nil))
	utils.AssertContainsAll(t, output,
		`rel="stylesheet"`,
		`href="/styles.css"`,
	)
}

func TestStylesheetWithAttrs(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Stylesheet("/print.css", templ.Attributes{
		"media": "print",
	}))
	utils.AssertContainsAll(t, output,
		`href="/print.css"`,
		`media="print"`,
	)
}

func TestStylesheetEmptyHref(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Stylesheet("", nil))
	_ = output // should not panic
}

// ---------------------------------------------------------------------------
// HTMX CDN URL builder coverage
// ---------------------------------------------------------------------------

func TestBaseHTMXCDNOverrides(t *testing.T) {
	t.Parallel()

	props := DefaultPageProps()
	props.Title = "Self-hosted"
	props.HTMXCDN = "https://cdn.example.com"
	output := utils.Render(t, Base(props))
	utils.AssertContains(t, output, "cdn.example.com")
}

// TestBaseHTMXSrcSelfHost verifies the HTMXSrc field switches the htmx main
// script from the CDN URL to a self-hosted path, emits no SRI (same-origin),
// skips the CDN preconnect, and suppresses the CDN response-targets extension
// (self-host implies you manage extensions). See ADR-0007.
func TestBaseHTMXSrcSelfHost(t *testing.T) {
	t.Parallel()

	t.Run("HTMXSrc emits self-hosted script, skips CDN + response-targets", func(t *testing.T) {
		t.Parallel()

		props := DefaultPageProps()
		props.HTMXSrc = "/static/htmx.min.js"
		// HTMXResponseTargets left at default (true) — must still be suppressed.
		output := utils.Render(t, Base(props))
		utils.AssertContains(t, output, `src="/static/htmx.min.js"`)
		utils.AssertNotContains(t, output, "cdn.jsdelivr.net")
		utils.AssertNotContains(t, output, "preconnect")
		utils.AssertNotContains(t, output, "integrity=")
		utils.AssertNotContains(t, output, "response-targets")
	})

	t.Run("HTMXSrc with empty Version still emits script", func(t *testing.T) {
		t.Parallel()

		props := DefaultPageProps()
		props.HTMXVersion = ""
		props.HTMXSrc = "/vendor/htmx.js"
		output := utils.Render(t, Base(props))
		utils.AssertContains(t, output, `src="/vendor/htmx.js"`)
	})
}
