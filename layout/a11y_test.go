package layout

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

const (
	testNonce   = "test-nonce"
	testCSSPath = "/app.css"
	testPage    = "Test Page"
)

func TestSecurityHeaders(t *testing.T) {
	t.Run("security headers rendered when enabled", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		props.SecurityHeaders = true
		props.Nonce = testNonce
		output := utils.Render(t, Base(props))
		utils.AssertContains(t, output, `http-equiv="X-Content-Type-Options"`)
		utils.AssertContains(t, output, `content="nosniff"`)
		utils.AssertContains(t, output, `http-equiv="Referrer-Policy"`)
		utils.AssertContains(t, output, `content="strict-origin-when-cross-origin"`)
	})

	t.Run("security headers not rendered when disabled", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		props.SecurityHeaders = false
		props.Nonce = testNonce
		output := utils.Render(t, Base(props))
		utils.AssertNotContains(t, output, `http-equiv="X-Content-Type-Options"`)
	})

	t.Run("skip link is present for accessibility", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		props.Nonce = "n"
		output := utils.Render(t, Base(props))
		utils.AssertContains(t, output, `href="#main-content"`)
		utils.AssertContains(t, output, "Skip to main content")
		utils.AssertContains(t, output, "sr-only")
	})

	t.Run("main landmark exists", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		props.Nonce = "n"
		output := utils.Render(t, Base(props))
		utils.AssertContains(t, output, `id="main-content"`)
	})
}

func TestDefaultPageProps(t *testing.T) {
	t.Parallel()

	t.Run("defaults are sensible", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		if props.Locale != "en" {
			t.Errorf("Locale = %q, want %q", props.Locale, "en")
		}
		if props.HTMXUseSRI != true {
			t.Error("HTMXUseSRI should be true by default")
		}
		if props.CSSPath != testCSSPath {
			t.Errorf("CSSPath = %q, want %q", props.CSSPath, testCSSPath)
		}
		if props.BodyClass == "" {
			t.Error("BodyClass should not be empty")
		}
	})
}

func TestHTMXSRI(t *testing.T) {
	t.Parallel()

	t.Run("known version returns SRI hash", func(t *testing.T) {
		t.Parallel()
		hash := htmxSRI(htmxVersion206, "main")
		if hash == "" {
			t.Error("expected non-empty SRI hash for known version")
		}
		utils.AssertContains(t, hash, "sha384-")
	})

	t.Run("unknown version returns empty string", func(t *testing.T) {
		t.Parallel()
		hash := htmxSRI("99.99.99", "main")
		if hash != "" {
			t.Errorf("expected empty hash for unknown version, got %q", hash)
		}
	})

	t.Run("response-targets extension returns hash", func(t *testing.T) {
		t.Parallel()
		hash := htmxSRI(htmxVersion206, "response-targets")
		if hash == "" {
			t.Error("expected non-empty SRI hash for response-targets")
		}
	})
}

func TestBaseDarkMode(t *testing.T) {
	t.Parallel()

	t.Run("body has dark mode classes", func(t *testing.T) {
		t.Parallel()
		props := DefaultPageProps()
		props.Nonce = "n"
		output := utils.Render(t, Base(props))
		utils.AssertContains(t, output, "dark:bg-gray-950")
		utils.AssertContains(t, output, "dark:text-gray-100")
	})
}
