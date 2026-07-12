package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

func TestGoldenFilterDropdown(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FilterDropdown(FilterDropdownProps{
		Name:  "status",
		Label: "Status",
		Value: "active",
		Options: []SelectOption{
			{Value: "all", Label: "All"},
			{Value: "active", Label: "Active"},
			{Value: "inactive", Label: "Inactive"},
		},
		HxGet:    "/api/users",
		HxTarget: "#user-list",
	}))
	golden.Assert(t, "filter_dropdown_basic", output)
}

func TestGoldenSlider(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Slider(SliderProps{
		Name:      "volume",
		Label:     "Volume",
		Min:       0,
		Max:       100,
		Value:     50,
		Step:      5,
		ShowValue: true,
	}))
	golden.Assert(t, "slider_basic", output)
}

func TestGoldenRatingInteractive(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Rating(RatingProps{
		Name:  "quality",
		Value: 3,
		Max:   5,
		Label: "Quality",
	}))
	golden.Assert(t, "rating_interactive", output)
}

func TestGoldenRatingReadOnly(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Rating(RatingProps{
		Value:    4,
		Max:      5,
		ReadOnly: true,
	}))
	golden.Assert(t, "rating_readonly", output)
}
