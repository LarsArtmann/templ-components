package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestDefaultRatingProps(t *testing.T) {
	t.Parallel()
	p := DefaultRatingProps()
	if p.Max != 5 || p.Size != RatingSizeMD {
		t.Error("expected default Max=5, Size=MD")
	}
}

func TestRatingSizeIsValid(t *testing.T) {
	t.Parallel()
	if !RatingSizeIsValid(RatingSizeSM) {
		t.Error("SM should be valid")
	}
	if !RatingSizeIsValid(RatingSizeMD) {
		t.Error("MD should be valid")
	}
	if !RatingSizeIsValid(RatingSizeLG) {
		t.Error("LG should be valid")
	}
	if RatingSizeIsValid(RatingSize("xl")) {
		t.Error("XL should be invalid")
	}
}

func TestRatingInteractiveRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Rating(RatingProps{
		Name:  "quality",
		Value: 3,
		Max:   5,
	}))
	utils.AssertContains(t, output, `type="radio"`)
	utils.AssertContains(t, output, `name="quality"`)
	utils.AssertContains(t, output, `value="3"`)
	utils.AssertContains(t, output, `checked`)
	utils.AssertContains(t, output, `role="radiogroup"`)
}

func TestRatingReadOnlyRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Rating(RatingProps{
		Value:    4,
		Max:      5,
		ReadOnly: true,
	}))
	utils.AssertContains(t, output, `role="img"`)
	utils.AssertContains(t, output, "4 out of 5")
	utils.AssertNotContains(t, output, `type="radio"`)
}

func TestRatingMaxStars(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Rating(RatingProps{
		Name:  "score",
		Value: 7,
		Max:   10,
	}))
	// Should have radio inputs for 10 stars
	radioCount := substringCount(output, `type="radio"`)
	if radioCount != 10 {
		t.Errorf("expected 10 radio inputs, got %d", radioCount)
	}
}

func TestRatingDefaultValueWhenMax0(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Rating(RatingProps{
		Name:  "score",
		Value: 3,
		Max:   0,
	}))
	// Max defaults to 5
	radioCount := substringCount(output, `type="radio"`)
	if radioCount != 5 {
		t.Errorf("expected 5 radio inputs when Max=0, got %d", radioCount)
	}
}

func TestRatingHelpText(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Rating(RatingProps{
		Name:     "quality",
		Value:    3,
		HelpText: "Rate the overall quality",
	}))
	utils.AssertContains(t, output, "Rate the overall quality")
}

func TestRatingLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Rating(RatingProps{
		Name:  "quality",
		Value: 3,
		Label: "Product Quality",
	}))
	utils.AssertContains(t, output, "Product Quality")
}

func TestRatingDarkMode(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Rating(RatingProps{
		Name:  "q",
		Value: 3,
	}))
	utils.AssertContains(t, output, "dark:text-gray-600")
}

func TestRatingReadOnlyZeroValue(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Rating(RatingProps{
		Value:    0,
		Max:      5,
		ReadOnly: true,
	}))
	utils.AssertContains(t, output, "0 out of 5")
}

func TestRatingSizes(t *testing.T) {
	t.Parallel()
	for _, size := range []RatingSize{RatingSizeSM, RatingSizeMD, RatingSizeLG} {
		output := utils.Render(t, Rating(RatingProps{
			Value:    3,
			Max:      5,
			Size:     size,
			ReadOnly: true,
		}))
		expected := ratingSizeLookup[size]
		utils.AssertContains(t, output, expected)
	}
}

//nolint:unparam // substr is always type="radio" but kept for flexibility
func substringCount(s, substr string) int {
	count := 0
	for {
		idx := indexOf(s, substr)
		if idx == -1 {
			break
		}
		count++
		s = s[idx+len(substr):]
	}
	return count
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
