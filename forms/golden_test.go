package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
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

// TestGoldenInputBasic and TestGoldenInputError snapshot the most-used Input
// render paths: a labelled text input (with help text) and an errored input
// (aria-invalid + field error). search_input above covers the search variant;
// together they guard the label/help/error wiring that assertion tests check
// only piecewise.
func TestGoldenInputBasic(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Input(InputProps{
		BaseProps:   utils.BaseProps{ID: "email"},
		Name:        "email",
		Label:       "Email address",
		Value:       "ada@example.com",
		Placeholder: "you@example.com",
		HelpText:    "We will never share your email.",
	}))
	golden.Assert(t, "input_basic", output)
}

func TestGoldenInputError(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Input(InputProps{
		BaseProps: utils.BaseProps{ID: "email"},
		Name:      "email",
		Label:     "Email address",
		Value:     "not-an-email",
		Error:     "Please enter a valid email address.",
	}))
	golden.Assert(t, "input_error", output)
}
