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
	t.Parallel()
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

	t.Run("skip link and main landmark present", func(t *testing.T) {
		t.Parallel()
		output := renderBaseWithNonce(t, "n")
		utils.AssertContains(t, output, `href="#main-content"`)
		utils.AssertContains(t, output, "Skip to main content")
		utils.AssertContains(t, output, "sr-only")
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

func TestHTMXMainSRI(t *testing.T) {
	t.Parallel()

	t.Run("known version returns SRI hash", func(t *testing.T) {
		t.Parallel()
		hash := htmxMainSRI(defaultHTMXVersion)
		if hash == "" {
			t.Error("expected non-empty SRI hash for known version")
		}
		utils.AssertContains(t, hash, "sha384-")
	})

	t.Run("unknown version falls back to default SRI", func(t *testing.T) {
		t.Parallel()
		hash := htmxMainSRI("99.99.99")
		if hash == "" {
			t.Error("expected non-empty fallback SRI hash for unknown version (SRI must not be silently dropped)")
		}
	})

	t.Run("response-targets extension SRI is pinned", func(t *testing.T) {
		t.Parallel()
		if sriResponseTargets == "" {
			t.Error("expected non-empty SRI hash for response-targets")
		}
		utils.AssertContains(t, sriResponseTargets, "sha384-")
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
