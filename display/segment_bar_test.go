package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

func TestSegmentBarRender(t *testing.T) {
	t.Parallel()

	segments := []SegmentBarSegment{
		{Label: "Go", Value: 62.3},
		{Label: "TypeScript", Value: 37.7},
	}

	t.Run("renders proportional flex-grow segments", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SegmentBar(SegmentBarProps{Segments: segments}))
		utils.AssertContains(t, output, "flex-grow:62.3")
		utils.AssertContains(t, output, "flex-grow:37.7")
		utils.AssertContains(t, output, `role="img"`)
	})

	t.Run("legend shows labels and percents", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SegmentBar(SegmentBarProps{Segments: segments}))
		utils.AssertContains(t, output, "Go")
		utils.AssertContains(t, output, "62.3%")
		utils.AssertContains(t, output, "37.7%")
	})

	t.Run("no legend opt-out", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, renderSegmentBarNoLegend(segments))
		utils.AssertNotContains(t, output, "<li")
	})

	t.Run("aria label composed from segments", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SegmentBar(SegmentBarProps{Segments: segments}))
		utils.AssertContains(t, output, `aria-label="Go 62.3%, TypeScript 37.7%"`)
	})

	t.Run("aria label override", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SegmentBar(SegmentBarProps{
			BaseProps: utils.BaseProps{AriaLabel: "Language shares"},
			Segments:  segments,
		}))
		utils.AssertContains(t, output, `aria-label="Language shares"`)
	})

	t.Run("empty or zero-total renders nothing", func(t *testing.T) {
		t.Parallel()
		if output := utils.Render(t, SegmentBar(SegmentBarProps{})); output != "" {
			t.Errorf("empty segments rendered %q", output)
		}

		zero := []SegmentBarSegment{{Label: "A", Value: 0}, {Label: "B", Value: -5}}
		if output := utils.Render(t, SegmentBar(SegmentBarProps{Segments: zero})); output != "" {
			t.Errorf("zero-total rendered %q", output)
		}
	})

	t.Run("negative values clamp in percents", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SegmentBar(SegmentBarProps{
			Segments: []SegmentBarSegment{{Label: "A", Value: 10}, {Label: "B", Value: -3}},
		}))
		utils.AssertContains(t, output, "100.0%")
		utils.AssertNotContains(t, output, "-3")
	})

	t.Run("explicit color override", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SegmentBar(SegmentBarProps{
			Segments: []SegmentBarSegment{{Label: "A", Value: 1, Color: "bg-emerald-600 dark:bg-emerald-400"}},
		}))
		utils.AssertContains(t, output, "bg-emerald-600")
	})

	t.Run("deterministic palette keyed by label", func(t *testing.T) {
		t.Parallel()
		first := utils.Render(t, SegmentBar(SegmentBarProps{
			Segments: []SegmentBarSegment{{Label: "Go", Value: 1}},
		}))
		second := utils.Render(t, SegmentBar(SegmentBarProps{
			Segments: []SegmentBarSegment{{Label: "Go", Value: 1}},
		}))
		if first != second {
			t.Error("same label renders different colors across renders")
		}
	})

	t.Run("propagates base props", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SegmentBar(SegmentBarProps{
			BaseProps: utils.BaseProps{ID: "bar", Class: "mb-2", Attrs: templ.Attributes{"data-x": "1"}},
			Segments:  segments,
		}))
		utils.AssertContains(t, output, `id="bar"`)
		utils.AssertContains(t, output, "mb-2")
		utils.AssertContains(t, output, `data-x="1"`)
	})
}

func TestSegmentBarHelpers(t *testing.T) {
	t.Parallel()

	t.Run("percent zero total", func(t *testing.T) {
		t.Parallel()
		if got := segmentBarPercent(5, 0); got != "0.0%" {
			t.Errorf("segmentBarPercent(5, 0) = %q", got)
		}
	})

	t.Run("color palette cycle is deterministic", func(t *testing.T) {
		t.Parallel()
		a := segmentBarColor(SegmentBarSegment{Label: "a"})
		b := segmentBarColor(SegmentBarSegment{Label: "a"})
		if a != b {
			t.Errorf("label 'a' resolved to two colors: %q vs %q", a, b)
		}
	})

	t.Run("total ignores negative values", func(t *testing.T) {
		t.Parallel()
		got := segmentBarTotal([]SegmentBarSegment{{Value: 5}, {Value: -2}, {Value: 5}})
		if got != 10 {
			t.Errorf("segmentBarTotal = %v, want 10", got)
		}
	})
}

func TestGoldenSweepSegmentBar(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "segment_bar_two", HTML: utils.Render(t, SegmentBar(SegmentBarProps{
			Segments: []SegmentBarSegment{
				{Label: "Go", Value: 62.3},
				{Label: "TypeScript", Value: 37.7},
			},
		}))},
		{Name: "segment_bar_many", HTML: utils.Render(t, SegmentBar(SegmentBarProps{
			Segments: []SegmentBarSegment{
				{Label: "Uploads", Value: 120},
				{Label: "Fetches", Value: 80},
				{Label: "Body", Value: 40},
				{Label: "Path", Value: 12},
			},
		}))},
		{Name: "segment_bar_no_legend", HTML: utils.Render(t, renderSegmentBarNoLegend([]SegmentBarSegment{
			{Label: "Go", Value: 62.3},
			{Label: "TypeScript", Value: 37.7},
		}))},
		{Name: "segment_bar_custom_color", HTML: utils.Render(t, SegmentBar(SegmentBarProps{
			Segments: []SegmentBarSegment{
				{Label: "Safe", Value: 90, Color: "bg-emerald-600 dark:bg-emerald-400"},
				{Label: "Risky", Value: 10, Color: "bg-rose-600 dark:bg-rose-400"},
			},
		}))},
	})
}

// SegmentBarProps2 renders a legend-less bar (test shorthand).
func renderSegmentBarNoLegend(segments []SegmentBarSegment) templ.Component {
	return SegmentBar(SegmentBarProps{Segments: segments, NoLegend: true})
}
