package forms

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// --- FilterDropdown Accessibility ---

func TestFilterDropdownA11y(t *testing.T) {
	t.Parallel()

	t.Run("select has correct name attribute", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterDropdown(FilterDropdownProps{
			Name:     "status",
			Options:  []SelectOption{{Value: "all", Label: "All"}},
			HxGet:    "/api",
			HxTarget: "#results",
		}))
		utils.AssertContains(t, output, `name="status"`)
	})

	t.Run("propagates aria-label from BaseProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterDropdown(FilterDropdownProps{
			BaseProps: utils.BaseProps{AriaLabel: "Filter by status"},
			Name:      "status",
			Options:   []SelectOption{{Value: "all", Label: "All"}},
			HxGet:     "/api",
			HxTarget:  "#results",
		}))
		utils.AssertContains(t, output, `aria-label="Filter by status"`)
	})

	t.Run("label is associated with select when ID is set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterDropdown(FilterDropdownProps{
			BaseProps: utils.BaseProps{ID: "status-filter"},
			Name:      "status",
			Label:     "Status",
			Options:   []SelectOption{{Value: "all", Label: "All"}},
			HxGet:     "/api",
			HxTarget:  "#results",
		}))
		utils.AssertContains(t, output, `id="status-filter"`)
	})

	t.Run("dark mode classes present on select", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterDropdown(FilterDropdownProps{
			Name:     "status",
			Options:  []SelectOption{{Value: "all", Label: "All"}},
			HxGet:    "/api",
			HxTarget: "#results",
		}))
		utils.AssertContains(t, output, "dark:bg-gray-800")
		utils.AssertContains(t, output, "dark:text-white")
	})
}

// --- Slider Accessibility ---

func TestSliderA11y(t *testing.T) {
	t.Parallel()

	t.Run("range input has correct aria attributes when error", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Slider(SliderProps{
			BaseProps: utils.BaseProps{ID: "vol"},
			Name:      "vol",
			Value:     50,
			Error:     "Too loud",
		}))
		utils.AssertContains(t, output, `aria-invalid="true"`)
	})

	t.Run("range input has aria-label when set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Slider(SliderProps{
			BaseProps: utils.BaseProps{AriaLabel: "Volume control"},
			Name:      "vol",
			Value:     50,
		}))
		utils.AssertContains(t, output, `aria-label="Volume control"`)
	})

	t.Run("label associated via for when ID set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Slider(SliderProps{
			BaseProps: utils.BaseProps{ID: "vol-slider"},
			Name:      "vol",
			Label:     "Volume",
			Value:     50,
		}))
		utils.AssertContains(t, output, `for="vol-slider"`)
	})

	t.Run("disabled slider has disabled attribute", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Slider(SliderProps{
			Name:     "vol",
			Value:    50,
			Disabled: true,
		}))
		utils.AssertContains(t, output, "disabled")
	})

	t.Run("required slider shows asterisk", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Slider(SliderProps{
			BaseProps: utils.BaseProps{ID: "vol"},
			Name:      "vol",
			Label:     "Volume",
			Value:     50,
			Required:  true,
		}))
		utils.AssertContains(t, output, "*")
		utils.AssertContains(t, output, "required")
	})

	t.Run("showValue displays current value", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Slider(SliderProps{
			Name:      "vol",
			Label:     "Volume",
			Value:     75,
			ShowValue: true,
		}))
		utils.AssertContains(t, output, "75")
		utils.AssertContains(t, output, "data-slider-value")
	})
}

// --- Rating Accessibility ---

func TestRatingA11y(t *testing.T) {
	t.Parallel()

	t.Run("interactive rating has radiogroup role", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			Name:  "quality",
			Value: 3,
			Max:   5,
		}))
		utils.AssertContains(t, output, `role="radiogroup"`)
	})

	t.Run("readonly rating has img role", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			Value:    4,
			Max:      5,
			ReadOnly: true,
		}))
		utils.AssertContains(t, output, `role="img"`)
	})

	t.Run("interactive rating has aria-required when required", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			Name:     "quality",
			Value:    3,
			Max:      5,
			Required: true,
		}))
		utils.AssertContains(t, output, `aria-required="true"`)
	})

	t.Run("readonly rating does NOT have aria-required even when required", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			Value:    3,
			Max:      5,
			ReadOnly: true,
			Required: true,
		}))
		utils.AssertNotContains(t, output, `aria-required`)
	})

	t.Run("screen reader text shows value out of max", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			Value:    3,
			Max:      5,
			ReadOnly: true,
		}))
		utils.AssertContains(t, output, "3 out of 5")
		utils.AssertContains(t, output, "sr-only")
	})

	t.Run("each radio has sr-only label with star text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			Name:  "quality",
			Value: 3,
			Max:   5,
		}))
		utils.AssertContains(t, output, "1 star")
		utils.AssertContains(t, output, "2 stars")
		utils.AssertContains(t, output, "3 stars")
		utils.AssertContains(t, output, "4 stars")
		utils.AssertContains(t, output, "5 stars")
	})

	t.Run("aria-label defaults to Label when set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			Name:  "quality",
			Value: 3,
			Max:   5,
			Label: "Product Quality",
		}))
		utils.AssertContains(t, output, `aria-label="Product Quality"`)
	})

	t.Run("aria-label defaults to 'Rating' when no label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			Name:  "quality",
			Value: 3,
			Max:   5,
		}))
		utils.AssertContains(t, output, `aria-label="Rating"`)
	})

	t.Run("aria-label overrides Label when both set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			BaseProps: utils.BaseProps{AriaLabel: "Custom Rating"},
			Name:      "quality",
			Value:     3,
			Max:       5,
			Label:     "Product Quality",
		}))
		utils.AssertContains(t, output, `aria-label="Custom Rating"`)
	})

	t.Run("only first radio has required boolean attribute", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			BaseProps: utils.BaseProps{ID: "rate"},
			Name:      "quality",
			Value:     3,
			Max:       5,
			Required:  true,
		}))
		// Count only the boolean required attribute on inputs (not aria-required on div)
		requiredAttrCount := strings.Count(output, " required>")
		if requiredAttrCount != 1 {
			t.Errorf("expected 1 required boolean attribute, got %d", requiredAttrCount)
		}
	})
}
