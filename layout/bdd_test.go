// Package layout provides behavior-driven tests for layout components.
// These tests verify end-user experience: page structure, dark mode, security headers.
package layout

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func renderBaseWithSecurity(t *testing.T, enabled bool) string {
	t.Helper()

	props := DefaultPageProps()
	props.SecurityHeaders = enabled

	return utils.Render(t, Base(props))
}

func assertSecurityHeadersPresent(t *testing.T, enabled bool) {
	t.Helper()

	output := renderBaseWithSecurity(t, enabled)
	if enabled {
		utils.AssertContains(t, output, `http-equiv="X-Content-Type-Options"`)
		utils.AssertContains(t, output, `http-equiv="Referrer-Policy"`)
	} else {
		utils.AssertNotContains(t, output, `http-equiv="X-Content-Type-Options"`)
		utils.AssertNotContains(t, output, `http-equiv="Referrer-Policy"`)
	}
}

func renderBaseWithNonce(t *testing.T, nonce string) string {
	t.Helper()

	props := DefaultPageProps()
	props.Nonce = nonce

	return utils.Render(t, Base(props))
}

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
		props.HTMXVersion = defaultHTMXVersion
		output := utils.Render(t, Base(props))
		utils.AssertContains(t, output, "htmx.org@2.0.10")
	})

	t.Run("user sees HTMX loaded with SRI when enabled", func(t *testing.T) {
		t.Parallel()

		props := DefaultPageProps()
		props.HTMXUseSRI = true
		output := utils.Render(t, Base(props))
		utils.AssertContains(t, output, `integrity="sha384-`)
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
		output := renderBaseWithNonce(t, "my-nonce-123")
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

	output := utils.Render(t, ThemeScript("nonce-abc"))
	for _, tt := range []struct {
		name string
		want string
	}{
		{"reads localStorage preference", `localStorage.getItem('theme')`},
		{"has CSP nonce", `nonce="nonce-abc"`},
		{"checks system preference", "prefers-color-scheme"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			utils.AssertContains(t, output, tt.want)
		})
	}
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

		if props.HTMXVersion != defaultHTMXVersion {
			t.Errorf("expected HTMXVersion %q, got %q", defaultHTMXVersion, props.HTMXVersion)
		}

		if !props.SecurityHeaders {
			t.Error("expected SecurityHeaders to be true")
		}

		if !props.HTMXUseSRI {
			t.Error("expected HTMXUseSRI to be true")
		}
	})
}

// --- Script Behavior ---

func TestScriptUserGetsCSPSafeScriptTag(t *testing.T) {
	t.Parallel()

	t.Run("user includes a script with nonce", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Script("abc123", "/static/app.js", nil))
		utils.AssertContains(t, output, `src="/static/app.js"`)
		utils.AssertContains(t, output, `nonce="abc123"`)
	})

	t.Run("user can add defer and module type", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Script("n1", "/app.js", templ.Attributes{
			"defer": true,
			"type":  "module",
		}))
		utils.AssertContains(t, output, "defer")
		utils.AssertContains(t, output, `type="module"`)
	})
}
