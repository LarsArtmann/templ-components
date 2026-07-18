package htmx

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestViewTransitionsRender(t *testing.T) {
	t.Parallel()

	t.Run("global enables config", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ViewTransitions(ViewTransitionsProps{
			Global:    true,
			BaseProps: utils.BaseProps{Nonce: "test-nonce"},
		}))
		utils.AssertContains(t, output, "globalViewTransitions")
		utils.AssertContains(t, output, `nonce="test-nonce"`)
	})

	t.Run("non-global omits script", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ViewTransitions(ViewTransitionsProps{
			Global:    false,
			BaseProps: utils.BaseProps{Nonce: "n"},
		}))
		utils.AssertNotContains(t, output, "globalViewTransitions")
		utils.AssertNotContains(t, output, "<script")
	})

	t.Run("always renders CSS", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ViewTransitions(ViewTransitionsProps{
			BaseProps: utils.BaseProps{Nonce: "n"},
		}))
		utils.AssertContains(t, output, "view-transition-old")
		utils.AssertContains(t, output, "view-transition-new")
		utils.AssertContains(t, output, "prefers-reduced-motion")
	})

	t.Run("default props global is true", func(t *testing.T) {
		t.Parallel()

		props := DefaultViewTransitionsProps()
		if !props.Global {
			t.Error("expected Global=true by default")
		}
	})
}
