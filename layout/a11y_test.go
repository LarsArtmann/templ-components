package layout

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
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

// TestBodyPrimitivesDoNotEmitMain is a contract guard: layout body primitives
// (AppShell, Split, Stack, Container) must NEVER emit their own <main> element.
// Base owns the singleton <main> landmark (WCAG 2.4.1 Bypass Blocks); nested
// <main> is invalid HTML and breaks screen-reader navigation. This test would
// have caught the original Split <main> regression.
func TestBodyPrimitivesDoNotEmitMain(t *testing.T) {
	t.Parallel()

	props := DefaultAppShellProps()
	props.Content = Container(DefaultContainerProps())

	for _, tc := range []struct {
		name string
		comp templ.Component
	}{
		{"AppShell", AppShell(props)},
		{"Split", Split(DefaultSplitProps())},
		{"Stack", Stack(DefaultStackProps())},
		{"Container", Container(DefaultContainerProps())},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := utils.Render(t, tc.comp)
			if strings.Contains(output, "<main") {
				t.Errorf("%s must not emit <main> — Base owns the main landmark", tc.name)
			}
		})
	}
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

	t.Run("unknown version returns empty (no wrong hash)", func(t *testing.T) {
		t.Parallel()

		hash := htmxMainSRI("99.99.99")
		if hash != "" {
			t.Errorf("expected empty SRI for unknown version, got %q (wrong hash would block the script)", hash)
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

func TestScriptA11y(t *testing.T) {
	t.Parallel()

	t.Run("nonce is always present on script tag", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Script("my-nonce", "/app.js", nil))
		utils.AssertContains(t, output, `nonce="my-nonce"`)
	})

	t.Run("empty nonce still emits the attribute", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Script("", "/app.js", nil))
		utils.AssertContains(t, output, `nonce=""`)
	})
}
