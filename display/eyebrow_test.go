package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

func TestEyebrowGoldenSweep(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "eyebrow_default", HTML: utils.Render(t, Eyebrow(EyebrowProps{Text: "Deploy #142 · production"}))},
		{Name: "eyebrow_accent", HTML: utils.Render(t, Eyebrow(EyebrowProps{
			Text:       "DNS block · 12:47:03",
			BaseProps:  utils.BaseProps{Class: "text-red-600 dark:text-red-400"},
		}))},
		{Name: "eyebrow_empty", HTML: utils.Render(t, Eyebrow(DefaultEyebrowProps()))},
	})
}

func TestEyebrowBehavior(t *testing.T) {
	t.Parallel()

	t.Run("user sees the eyebrow text verbatim", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Eyebrow(EyebrowProps{Text: "Section label"}))
		utils.AssertContains(t, output, "Section label")
	})

	t.Run("eyebrow reads as an uppercase overline label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Eyebrow(EyebrowProps{Text: "overline"}))
		utils.AssertContainsAll(t, output, "uppercase", "font-semibold", "tracking-[0.18em]", "font-mono")
	})

	t.Run("consumer class override wins for accent theming", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Eyebrow(EyebrowProps{
			Text:      "accented",
			BaseProps: utils.BaseProps{Class: "text-blue-600 dark:text-blue-400"},
		}))
		utils.AssertContains(t, output, "text-blue-600")
		utils.AssertNotContains(t, output, "text-gray-500")
	})
}

func TestEyebrowA11y(t *testing.T) {
	t.Parallel()

	t.Run("propagates aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Eyebrow(EyebrowProps{
			Text:      "status",
			BaseProps: utils.BaseProps{AriaLabel: "Deployment status"},
		}))
		utils.AssertContains(t, output, `aria-label="Deployment status"`)
	})

	t.Run("neutral text has dark mode variant", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Eyebrow(EyebrowProps{Text: "x"}))
		utils.AssertContainsAll(t, output, "text-gray-500", "dark:text-gray-400")
	})
}

func TestEyebrowEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("empty text renders nothing", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Eyebrow(EyebrowProps{}))
		if output != "" {
			t.Errorf("expected empty output, got %q", output)
		}
	})

	t.Run("propagates id and attrs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Eyebrow(EyebrowProps{
			Text: "x",
			BaseProps: utils.BaseProps{
				ID:    "eyebrow-1",
				Attrs: map[string]any{"data-testid": "eyebrow"},
			},
		}))
		utils.AssertContainsAll(t, output, `id="eyebrow-1"`, `data-testid="eyebrow"`)
	})
}
