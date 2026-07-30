package forms

import (
	"strings"
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

func TestRatingKeyboardOrder(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Rating(RatingProps{
		Name:  "quality",
		Value: 3,
		Max:   5,
	}))

	// Forward DOM order: radio value=1 must precede value=2 so radiogroup arrow
	// keys (ArrowDown/ArrowRight) move toward increasing value per WAI-ARIA.
	// The old reversed order (5..1) made arrows decrease the value.
	idx1 := strings.Index(output, `value="1"`)
	idx2 := strings.Index(output, `value="2"`)
	if idx1 < 0 || idx2 < 0 {
		t.Fatalf("expected value=1 and value=2 radios; got indices %d, %d", idx1, idx2)
	}

	if idx1 >= idx2 {
		t.Errorf("expected value=1 before value=2 (forward DOM order); got 1@%d, 2@%d", idx1, idx2)
	}

	// flex-row-reverse maps the forward DOM back to a correct visual so the
	// peer-checked ~ fill still lights the leftmost N stars.
	utils.AssertContains(t, output, "flex-row-reverse")

	// Fill classes live on the <label> (sibling of the .peer radio), not the
	// nested <svg>, so peer-checked resolves and cumulatively fills the checked
	// star and every lower-value star.
	utils.AssertContains(t, output, "peer-checked:text-amber-400")
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
