package forms

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// --- FilterDropdown Behavior (BDD-style) ---

func TestFilterDropdownUserCanFilterResults(t *testing.T) {
	t.Parallel()

	t.Run("user sees a select with filter options", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterDropdown(FilterDropdownProps{
			Name:  "status",
			Label: "Status",
			Options: []SelectOption{
				{Value: "all", Label: "All"},
				{Value: "active", Label: "Active"},
				{Value: "inactive", Label: "Inactive"},
			},
			HxGet:    "/api/users",
			HxTarget: "#user-list",
		}))
		utils.AssertContains(t, output, "All")
		utils.AssertContains(t, output, "Active")
		utils.AssertContains(t, output, "Inactive")
	})

	t.Run("user sees their previous selection preserved", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterDropdown(FilterDropdownProps{
			Name:  "status",
			Value: "active",
			Options: []SelectOption{
				{Value: "all", Label: "All"},
				{Value: "active", Label: "Active"},
				{Value: "inactive", Label: "Inactive"},
			},
			HxGet:    "/api/users",
			HxTarget: "#user-list",
		}))
		utils.AssertContains(t, output, "selected")
		utils.AssertContains(t, output, `value="active"`)
	})

	t.Run("user sees help text explaining the filter", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterDropdown(FilterDropdownProps{
			Name:     "date",
			Label:    "Date Range",
			HelpText: "Filter by creation date",
			Options:  []SelectOption{{Value: "today", Label: "Today"}},
			HxGet:    "/api",
			HxTarget: "#results",
		}))
		utils.AssertContains(t, output, "Filter by creation date")
	})
}

// --- Slider Behavior ---

func TestSliderUserCanAdjustRange(t *testing.T) {
	t.Parallel()

	t.Run("user sees a range input with min max and step", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Slider(SliderProps{
			Name:  "volume",
			Label: "Volume",
			Min:   0,
			Max:   100,
			Value: 50,
			Step:  5,
		}))
		utils.AssertContains(t, output, `type="range"`)
		utils.AssertContains(t, output, `min="0"`)
		utils.AssertContains(t, output, `max="100"`)
		utils.AssertContains(t, output, `step="5"`)
	})

	t.Run("user sees current value displayed", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Slider(SliderProps{
			Name:      "volume",
			Label:     "Volume",
			Value:     42,
			ShowValue: true,
		}))
		utils.AssertContains(t, output, "42")
	})

	t.Run("user sees error message when invalid", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Slider(SliderProps{
			BaseProps: utils.BaseProps{ID: "vol"},
			Name:      "vol",
			Value:     50,
			Error:     "Value exceeds maximum",
		}))
		utils.AssertContains(t, output, "Value exceeds maximum")
		utils.AssertContains(t, output, `aria-invalid="true"`)
	})
}

// --- Rating Behavior ---

func TestRatingUserCanSelectStars(t *testing.T) {
	t.Parallel()

	t.Run("user sees 5 star options for a 5-star rating", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			Name: "quality",
			Max:  5,
		}))
		count := substringCount(output, `type="radio"`)
		if count != 5 {
			t.Errorf("expected 5 radio inputs, got %d", count)
		}
	})

	t.Run("user sees current rating pre-selected", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			Name:  "quality",
			Value: 4,
			Max:   5,
		}))
		utils.AssertContains(t, output, `value="4"`)
		utils.AssertContains(t, output, "checked")
	})

	t.Run("user sees read-only display of rating", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			Value:    3,
			Max:      5,
			ReadOnly: true,
		}))
		utils.AssertContains(t, output, "3 out of 5")
		utils.AssertNotContains(t, output, `type="radio"`)
	})
}

// --- Edge Cases: FilterDropdown ---

func TestFilterDropdownEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("preselect does not mutate caller slice", func(t *testing.T) {
		t.Parallel()
		opts := []SelectOption{
			{Value: "a", Label: "A"},
			{Value: "b", Label: "B", Selected: true},
		}
		_ = filterDropdownPreselect(opts, "a")
		// Caller's slice must be unchanged
		if !opts[1].Selected {
			t.Error("preselect mutated caller slice: B.Selected should still be true")
		}
		if opts[0].Selected {
			t.Error("preselect mutated caller slice: A.Selected should still be false")
		}
	})

	t.Run("preselect with no matching value clears all selections", func(t *testing.T) {
		t.Parallel()
		opts := []SelectOption{
			{Value: "a", Label: "A"},
			{Value: "b", Label: "B", Selected: true},
		}
		result := filterDropdownPreselect(opts, "nonexistent")
		if result[0].Selected {
			t.Error("A should not be selected")
		}
		if result[1].Selected {
			t.Error("B should not be selected after non-matching preselect")
		}
	})

	t.Run("empty value returns same slice identity", func(t *testing.T) {
		t.Parallel()
		opts := []SelectOption{
			{Value: "a", Label: "A"},
		}
		result := filterDropdownPreselect(opts, "")
		// When value is empty, the same slice is returned (not a copy)
		if len(result) != len(opts) {
			t.Errorf("expected same length, got result=%d opts=%d", len(result), len(opts))
		}
		if result[0].Selected != opts[0].Selected {
			t.Error("values should match when value is empty")
		}
	})

	t.Run("no options renders empty select", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterDropdown(FilterDropdownProps{
			Name:     "empty",
			HxGet:    "/api",
			HxTarget: "#results",
		}))
		utils.AssertContains(t, output, "<select")
		utils.AssertContains(t, output, "</select>")
	})

	t.Run("all optional HTMX attrs omitted when empty", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, FilterDropdown(FilterDropdownProps{
			Name:     "status",
			Options:  []SelectOption{{Value: "all", Label: "All"}},
			HxGet:    "/api",
			HxTarget: "#results",
		}))
		utils.AssertNotContains(t, output, "hx-swap")
		utils.AssertNotContains(t, output, "hx-include")
		utils.AssertNotContains(t, output, "hx-indicator")
	})
}

// --- Edge Cases: Slider ---

func TestSliderEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("negative min and max range", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Slider(SliderProps{
			Name:  "temp",
			Min:   -50,
			Max:   50,
			Value: 0,
		}))
		utils.AssertContains(t, output, `min="-50"`)
		utils.AssertContains(t, output, `max="50"`)
		utils.AssertContains(t, output, `value="0"`)
	})

	t.Run("decimal step", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Slider(SliderProps{
			Name:  "price",
			Min:   0,
			Max:   10,
			Value: 5.5,
			Step:  0.5,
		}))
		utils.AssertContains(t, output, `step="0.5"`)
		utils.AssertContains(t, output, `value="5.5"`)
	})

	t.Run("zero step defaults to 1", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Slider(SliderProps{
			Name:  "x",
			Value: 50,
			Step:  0,
		}))
		utils.AssertContains(t, output, `step="1"`)
	})

	t.Run("no label renders bare input", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Slider(SliderProps{
			Name:  "vol",
			Value: 50,
		}))
		utils.AssertNotContains(t, output, "<label")
	})

	t.Run("custom class propagated", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Slider(SliderProps{
			BaseProps: utils.BaseProps{Class: "my-slider"},
			Name:      "vol",
			Value:     50,
		}))
		utils.AssertContains(t, output, "my-slider")
	})
}

// --- Edge Cases: Rating ---

func TestRatingEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("value exceeds max clips display", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			Value:    10,
			Max:      5,
			ReadOnly: true,
		}))
		// Value > Max: "10 out of 5" in sr-only text, but only 5 stars rendered
		utils.AssertContains(t, output, "10 out of 5")
	})

	t.Run("zero value in interactive mode renders no checked", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			Name:  "score",
			Value: 0,
			Max:   5,
		}))
		// "checked" also appears in CSS class names (peer-checked:),
		// so check for the boolean attribute pattern specifically
		checkedAttrCount := strings.Count(output, `" checked>`)
		if checkedAttrCount != 0 {
			t.Errorf("expected 0 checked boolean attributes, got %d", checkedAttrCount)
		}
	})

	t.Run("negative value in readonly mode", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			Value:    -1,
			Max:      5,
			ReadOnly: true,
		}))
		utils.AssertContains(t, output, "-1 out of 5")
	})

	t.Run("max stars of 1", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			Name:  "binary",
			Value: 1,
			Max:   1,
		}))
		count := substringCount(output, `type="radio"`)
		if count != 1 {
			t.Errorf("expected 1 radio input, got %d", count)
		}
	})

	t.Run("invalid size falls back to MD", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			Value:    3,
			Max:      5,
			Size:     RatingSize("invalid"),
			ReadOnly: true,
		}))
		utils.AssertContains(t, output, "h-5 w-5")
	})

	t.Run("custom class propagated", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			BaseProps: utils.BaseProps{Class: "my-rating"},
			Value:     3,
			Max:       5,
		}))
		utils.AssertContains(t, output, "my-rating")
	})

	t.Run("help text rendered", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Rating(RatingProps{
			Name:     "q",
			Value:    3,
			HelpText: "Click a star to rate",
		}))
		utils.AssertContains(t, output, "Click a star to rate")
	})
}

// --- Rating Coverage (private helpers) ---

func TestRatingSizeClass(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		size RatingSize
		want string
	}{
		{RatingSizeSM, "h-4 w-4"},
		{RatingSizeMD, "h-5 w-5"},
		{RatingSizeLG, "h-6 w-6"},
		{RatingSize("invalid"), "h-5 w-5"},
	} {
		got := ratingSizeClass(tt.size)
		if got != tt.want {
			t.Errorf("ratingSizeClass(%q) = %q, want %q", tt.size, got, tt.want)
		}
	}
}

func TestPluralStars(t *testing.T) {
	t.Parallel()
	if got := pluralStars(1); got != "" {
		t.Errorf("pluralStars(1) = %q, want empty", got)
	}
	if got := pluralStars(2); got != "s" {
		t.Errorf("pluralStars(2) = %q, want %q", got, "s")
	}
}
