package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// --- Snapshot tests for FilterDropdown, Slider, Rating ---

func TestFilterDropdownSnapshot(t *testing.T) {
	t.Parallel()

	t.Run("full-featured filter dropdown", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterDropdown(FilterDropdownProps{
			Name:        "status",
			Label:       "Status",
			Value:       "active",
			HxGet:       "/api/users",
			HxTarget:    "#user-list",
			HxSwap:      "innerHTML",
			HxTrigger:   "change delay:300ms",
			HxInclude:   "closest form",
			HxIndicator: "#loading",
			HelpText:    "Filter by account status",
			Options: []SelectOption{
				{Value: "all", Label: "All Users"},
				{Value: "active", Label: "Active"},
				{Value: "inactive", Label: "Inactive"},
				{Value: "banned", Label: "Banned", Disabled: true},
			},
		}))
		utils.AssertContainsAll(t, output,
			`name="status"`,
			`hx-get="/api/users"`,
			`hx-target="#user-list"`,
			`hx-swap="innerHTML"`,
			`hx-trigger="change delay:300ms"`,
			`hx-include="closest form"`,
			`hx-indicator="#loading"`,
			"Filter by account status",
			"All Users",
			"Active",
			"Inactive",
			"Banned",
			"disabled",
			"selected",
		)
	})
}

func TestSliderSnapshot(t *testing.T) {
	t.Parallel()

	t.Run("full-featured slider with all options", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Slider(SliderProps{
			BaseProps: utils.BaseProps{ID: "price-range", Class: "mb-4"},
			Name:      "max_price",
			Label:     "Maximum Price",
			Min:       0,
			Max:       1000,
			Value:     250,
			Step:      50,
			Required:  true,
			ShowValue: true,
			HelpText:  "Drag to adjust your budget",
		}))
		utils.AssertContainsAll(t, output,
			`type="range"`,
			`id="price-range"`,
			`name="max_price"`,
			`min="0"`,
			`max="1000"`,
			`value="250"`,
			`step="50"`,
			"required",
			"250",
			"Drag to adjust your budget",
			"mb-4",
		)
	})
}

func TestRatingSnapshot(t *testing.T) {
	t.Parallel()

	t.Run("interactive rating with all features", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			BaseProps: utils.BaseProps{ID: "product-rating", Class: "mb-4"},
			Name:      "product_score",
			Value:     4,
			Max:       5,
			Size:      RatingSizeLG,
			Label:     "Rate this product",
			Required:  true,
			HelpText:  "How would you rate your experience?",
		}))
		utils.AssertContainsAll(t, output,
			`role="radiogroup"`,
			"Rate this product",
			"required",
			"How would you rate your experience?",
			"h-6 w-6",
		)
	})

	t.Run("readonly rating with custom max", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			BaseProps: utils.BaseProps{AriaLabel: "Average customer rating"},
			Value:     7,
			Max:       10,
			Size:      RatingSizeSM,
			ReadOnly:  true,
		}))
		utils.AssertContainsAll(t, output,
			`role="img"`,
			`aria-label="Average customer rating"`,
			"7 out of 10",
			"h-4 w-4",
		)
	})
}
