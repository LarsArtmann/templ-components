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

func TestGoldenStylableSelect(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Select(SelectProps{
		BaseProps: utils.BaseProps{ID: "country"},
		Name:      "country",
		Label:     "Country",
		Stylable:  true,
		Options: []SelectOption{
			{Value: "de", Label: "Germany"},
			{Value: "at", Label: "Austria", Selected: true},
		},
	}))
	golden.Assert(t, "stylable_select", output)
}

func TestGoldenAutoGrowTextarea(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Textarea(TextareaProps{
		BaseProps: utils.BaseProps{ID: "bio"},
		Name:      "bio",
		Label:     "Bio",
		AutoGrow:  true,
	}))
	golden.Assert(t, "textarea_autogrow", output)
}

func TestGoldenSearchInput(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Input(InputProps{
		BaseProps:   utils.BaseProps{ID: "q"},
		Type:        InputSearch,
		Name:        "q",
		Placeholder: "Search...",
	}))
	golden.Assert(t, "search_input", output)
}
