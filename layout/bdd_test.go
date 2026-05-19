// Package layout provides behavior-driven tests for layout components.
// These tests verify end-user experience: page structure, dark mode, security headers.
package layout

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// --- Base Layout Behavior ---

func TestBaseUserGetsCompleteHTMLPage(t *testing.T) {
	t.Parallel()

	t.Run("user sees page title in HTML head", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		props.Title = "My Application"
		output := utils.Render(t, Base(props))
		utils.AssertContains(t, output, "<title>My Application</title>")
	})

	t.Run("user sees meta description for SEO", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		props.Description = "A beautiful web application"
		output := utils.Render(t, Base(props))
		utils.AssertContains(t, output, `name="description"`)
		utils.AssertContains(t, output, "A beautiful web application")
	})

	t.Run("user sees HTMX script with correct version", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		props.HTMXVersion = htmxVersion206
		output := utils.Render(t, Base(props))
		utils.AssertContains(t, output, "htmx.org@2.0.6")
	})

	t.Run("user sees HTMX loaded with SRI when enabled", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		props.HTMXUseSRI = true
		output := utils.Render(t, Base(props))
		utils.AssertContains(t, output, `integrity="sha384-`)
	})

	t.Run("user gets security headers when enabled", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		props.SecurityHeaders = true
		output := utils.Render(t, Base(props))
		utils.AssertContains(t, output, "X-Content-Type-Options")
		utils.AssertContains(t, output, "Referrer-Policy")
	})

	t.Run("user sees OG meta tags when OG image is set", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		props.OGImage = "https://example.com/og.png"
		output := utils.Render(t, Base(props))
		utils.AssertContains(t, output, `property="og:image"`)
		utils.AssertContains(t, output, "https://example.com/og.png")
	})

	t.Run("user sees custom CSS path", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		props.CSSPath = "/custom.css"
		output := utils.Render(t, Base(props))
		utils.AssertContains(t, output, `href="/custom.css"`)
	})

	t.Run("user sees CSP nonce on inline scripts", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		props.Nonce = "my-nonce-123"
		output := utils.Render(t, Base(props))
		utils.AssertContains(t, output, `nonce="my-nonce-123"`)
	})
}

// --- Minimal Layout Behavior ---

func TestMinimalUserGetsCleanHTMLDocument(t *testing.T) {
	t.Parallel()

	t.Run("user sees minimal HTML without external dependencies", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Minimal(MinimalProps{Title: testPage, Locale: "en"}))
		utils.AssertContains(t, output, "<title>Test Page</title>")
		utils.AssertContains(t, output, `lang="en"`)
	})
}

// --- ThemeScript Behavior ---

func TestThemeScriptUserGetsDarkModeWithoutFOUC(t *testing.T) {
	t.Parallel()

	t.Run("theme script reads localStorage preference", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ThemeScript("nonce-abc"))
		utils.AssertContains(t, output, `localStorage.getItem('theme')`)
	})

	t.Run("theme script has CSP nonce", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ThemeScript("nonce-abc"))
		utils.AssertContains(t, output, `nonce="nonce-abc"`)
	})

	t.Run("theme script checks system preference", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ThemeScript("n"))
		utils.AssertContains(t, output, "prefers-color-scheme")
	})
}

// --- ThemeToggle Behavior ---

func TestThemeToggleUserCanSwitchTheme(t *testing.T) {
	t.Parallel()

	t.Run("user sees theme toggle button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ThemeToggle("Toggle theme", "nonce"))
		utils.AssertContains(t, output, "<button")
		utils.AssertContains(t, output, "Toggle theme")
	})

	t.Run("theme toggle has accessible label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ThemeToggle("Switch appearance", "nonce"))
		utils.AssertContains(t, output, `aria-label="Switch appearance"`)
	})
}

// --- DefaultPageProps Behavior ---

func TestDefaultPagePropsProvidesSensibleDefaults(t *testing.T) {
	t.Parallel()

	t.Run("defaults include locale, HTMX version, and security", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		if props.Locale != "en" {
			t.Errorf("expected Locale 'en', got %q", props.Locale)
		}
		if props.HTMXVersion != htmxVersion206 {
			t.Errorf("expected HTMXVersion %q, got %q", htmxVersion206, props.HTMXVersion)
		}
		if !props.SecurityHeaders {
			t.Error("expected SecurityHeaders to be true")
		}
		if !props.HTMXUseSRI {
			t.Error("expected HTMXUseSRI to be true")
		}
	})
}
