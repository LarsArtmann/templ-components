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
