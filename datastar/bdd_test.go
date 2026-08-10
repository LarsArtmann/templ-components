package datastar

import (
	"context"
	"io"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// --- SDKScript Behavior ---

func TestSDKScriptUserGetsDatastarRuntime(t *testing.T) {
	t.Parallel()

	t.Run("user sees ES module script tag with CDN URL", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SDKScript(DefaultSDKScriptProps()))
		utils.AssertContains(t, output, `<script type="module"`)
		utils.AssertContains(t, output, `src="`)
		utils.AssertContains(t, output, "starfederation/datastar")
	})

	t.Run("user gets self-hosted script when Src is set", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SDKScript(SDKScriptProps{
			Src: "/assets/datastar.js",
		}))
		utils.AssertContains(t, output, `src="/assets/datastar.js"`)
		utils.AssertNotContains(t, output, "cdn.jsdelivr.net")
	})

	t.Run("script always carries CSP nonce for security", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SDKScript(SDKScriptProps{
			BaseProps: utils.BaseProps{Nonce: "nonce-abc-789"},
		}))
		utils.AssertContains(t, output, `nonce="nonce-abc-789"`)
	})
}

// --- LiveRegion Behavior ---

func TestLiveRegionUserSeesAutoStreamingContent(t *testing.T) {
	t.Parallel()

	t.Run("region auto-connects to SSE on page load", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, LiveRegion(LiveRegionProps{
			URL:       "/stream/metrics",
			AutoStart: true,
		}))
		utils.AssertContainsAll(t, output, "data-init", "@get(", "/stream/metrics")
	})

	t.Run("manual start omits data-init attribute", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, LiveRegion(LiveRegionProps{
			URL:       "/stream/data",
			AutoStart: false,
		}))
		utils.AssertNotContains(t, output, "data-init")
	})

	t.Run("aria-live polite by default for screen readers", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, LiveRegion(DefaultLiveRegionProps()))
		utils.AssertContains(t, output, `aria-live="polite"`)
	})

	t.Run("invalid politeness falls back to polite gracefully", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, LiveRegion(LiveRegionProps{
			Live: "invalid-value",
		}))
		utils.AssertContains(t, output, `aria-live="polite"`)
	})

	t.Run("consumer can set assertive for critical alerts", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, LiveRegion(LiveRegionProps{
			Live: LiveAssertive,
		}))
		utils.AssertContains(t, output, `aria-live="assertive"`)
	})
}

// --- Indicator Behavior ---

func TestIndicatorUserSeesLoadingFeedback(t *testing.T) {
	t.Parallel()

	t.Run("indicator shows when signal is active", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Indicator(IndicatorProps{
			Signal: "fetching",
		}))
		utils.AssertContains(t, output, `data-show="$fetching"`)
		utils.AssertContains(t, output, `role="status"`)
		utils.AssertContains(t, output, `aria-live="polite"`)
	})

	t.Run("default spinner respects motion-reduce", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Indicator(IndicatorProps{
			Signal: "loading",
		}))
		utils.AssertContains(t, output, "motion-reduce:animate-none")
	})

	t.Run("consumer can inject custom spinner component", func(t *testing.T) {
		t.Parallel()

		customSpinner := templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
			_, err := io.WriteString(w, `<div class="my-custom-spinner"></div>`)

			return err //nolint:wrapcheck // test helper, direct passthrough
		})

		output := utils.Render(t, Indicator(IndicatorProps{
			Signal:  "saving",
			Spinner: customSpinner,
		}))
		utils.AssertContains(t, output, `data-show="$saving"`)
		utils.AssertContains(t, output, "my-custom-spinner")
		utils.AssertNotContains(t, output, "defaultIndicatorSpinner")
	})

	t.Run("nil spinner renders CSS pulse fallback", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Indicator(IndicatorProps{
			Signal:  "syncing",
			Spinner: templ.NopComponent,
		}))
		utils.AssertContains(t, output, `data-show="$syncing"`)
		utils.AssertNotContains(t, output, "animate-spin")
	})
}

// --- Action Expressions Behavior ---

func TestActionExpressionsUserTriggersBackendActions(t *testing.T) {
	t.Parallel()

	t.Run("get action builds correct expression", func(t *testing.T) {
		t.Parallel()

		if got, want := Get("/api/search"), "@get('/api/search')"; got != want {
			t.Errorf("Get = %q, want %q", got, want)
		}
	})

	t.Run("post action builds correct expression", func(t *testing.T) {
		t.Parallel()

		if got, want := Post("/api/save"), "@post('/api/save')"; got != want {
			t.Errorf("Post = %q, want %q", got, want)
		}
	})

	t.Run("delete action builds correct expression", func(t *testing.T) {
		t.Parallel()

		if got, want := Delete("/api/items/42"), "@delete('/api/items/42')"; got != want {
			t.Errorf("Delete = %q, want %q", got, want)
		}
	})

	t.Run("single quotes in URLs are escaped for safety", func(t *testing.T) {
		t.Parallel()

		got := Get("/api/search?q=it's")
		utils.AssertContains(t, got, "\\'")
		utils.AssertNotContains(t, got, "it's")
	})
}

// --- Cross-Cutting Concerns ---

func TestDatastarComponentsRenderValidHTML(t *testing.T) {
	t.Parallel()

	t.Run("all datastar components render without errors", func(t *testing.T) {
		t.Parallel()

		components := []struct {
			name string
			comp func() templ.Component
		}{
			{"SDKScript", func() templ.Component { return SDKScript(DefaultSDKScriptProps()) }},
			{"LiveRegion", func() templ.Component {
				return LiveRegion(LiveRegionProps{URL: "/stream", AutoStart: true})
			}},
			{"Indicator", func() templ.Component { return Indicator(IndicatorProps{Signal: "loading"}) }},
		}

		for _, c := range components {
			t.Run(c.name, func(t *testing.T) {
				t.Parallel()

				output := utils.Render(t, c.comp())
				if output == "" {
					t.Errorf("%s rendered empty output", c.name)
				}
			})
		}
	})
}
