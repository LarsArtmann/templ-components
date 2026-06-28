// Package display provides tests for display components.
package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestTooltipRender(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		text     string
		position TooltipPosition
		wantAll  []string
		wantAny  []string
		wantNone []string
	}{
		{
			name:     "top tooltip",
			text:     "More info",
			position: TooltipPositionTop,
			wantAll:  []string{"More info", `role="tooltip"`, "group-hover:block", "bottom-full"},
		},
		{
			name:     "bottom tooltip",
			text:     "Help",
			position: TooltipPositionBottom,
			wantAll:  []string{"Help", "top-full"},
		},
		{
			name:     "right tooltip",
			text:     "Details",
			position: TooltipPositionRight,
			wantAll:  []string{"left-full"},
		},
		{
			name:     "left tooltip",
			text:     "Sidebar",
			position: TooltipPositionLeft,
			wantAll:  []string{"right-full"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Tooltip(TooltipProps{
				Text:     tt.text,
				Position: tt.position,
			}))
			for _, want := range tt.wantAll {
				utils.AssertContains(t, output, want)
			}
		})
	}
}

func TestDefaultTooltipProps(t *testing.T) {
	t.Parallel()
	props := DefaultTooltipProps()
	if props.Position != TooltipPositionTop {
		t.Errorf("DefaultTooltipProps().Position = %q, want %q", props.Position, TooltipPositionTop)
	}
}

func TestTooltipA11yLinkage(t *testing.T) {
	t.Parallel()
	t.Run("tooltip has role=tooltip", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tooltip(TooltipProps{
			Text:     "Help text",
			Position: TooltipPositionTop,
		}))
		utils.AssertContains(t, output, `role="tooltip"`)
	})

	t.Run("tooltip with custom class and ID", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tooltip(TooltipProps{
			BaseProps: utils.BaseProps{ID: "tip-1", Class: "ml-2"},
			Text:      "Hint",
			Position:  TooltipPositionTop,
		}))
		utils.AssertContains(t, output, `id="tip-1"`)
		utils.AssertContains(t, output, "ml-2")
	})
}

func TestTooltipPositionEdgeCases(t *testing.T) {
	t.Parallel()
	t.Run("tooltip without ID auto-generates id and aria-describedby", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tooltip(TooltipProps{
			Text:     "Tip",
			Position: TooltipPositionTop,
		}))
		utils.AssertContains(t, output, "aria-describedby=")
		utils.AssertContains(t, output, `role="tooltip"`)
	})
	t.Run("tooltip with ID sets aria-describedby linkage", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tooltip(TooltipProps{
			BaseProps: utils.BaseProps{ID: "my-tip"},
			Text:      "Tip",
			Position:  TooltipPositionTop,
		}))
		utils.AssertContains(t, output, `aria-describedby="my-tip-tooltip"`)
		utils.AssertContains(t, output, `id="my-tip-tooltip"`)
	})
}
