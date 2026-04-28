// Package display provides tests for display components.
package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestTooltipRender(t *testing.T) {
	t.Parallel()

	t.Run("top tooltip", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tooltip(TooltipProps{
			Text:     "More info",
			Position: TooltipPositionTop,
		}))
		utils.AssertContains(t, output, "More info")
		utils.AssertContains(t, output, `role="tooltip"`)
		utils.AssertContains(t, output, "group-hover:block")
		utils.AssertContains(t, output, "bottom-full")
	})

	t.Run("bottom tooltip", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tooltip(TooltipProps{
			Text:     "Help",
			Position: TooltipPositionBottom,
		}))
		utils.AssertContains(t, output, "Help")
		utils.AssertContains(t, output, "top-full")
	})

	t.Run("right tooltip", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tooltip(TooltipProps{
			Text:     "Details",
			Position: TooltipPositionRight,
		}))
		utils.AssertContains(t, output, "left-full")
	})
}
