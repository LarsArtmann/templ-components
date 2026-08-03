package display

import (
	"math"
	"testing"

	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

func TestGoldenSweepPieChart(t *testing.T) {
	t.Parallel()

	slices := []PieChartSlice{
		{Label: "Direct", Value: 45},
		{Label: "Organic", Value: 30},
		{Label: "Referral", Value: 25},
	}

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "pie_chart_basic", HTML: utils.Render(t, PieChart(PieChartProps{
			BaseProps: utils.BaseProps{AriaLabel: "Traffic sources"},
			Slices:    slices,
		}))},
		{Name: "pie_chart_donut", HTML: utils.Render(t, PieChart(PieChartProps{
			Slices: slices,
			Donut:  true,
		}))},
		{Name: "pie_chart_donut_center", HTML: utils.Render(t, PieChart(PieChartProps{
			Slices:      slices,
			Donut:       true,
			CenterLabel: "100%",
		}))},
		{Name: "pie_chart_no_labels", HTML: utils.Render(t, PieChart(PieChartProps{
			Slices:     slices,
			ShowLabels: false,
		}))},
		{Name: "pie_chart_no_legend", HTML: utils.Render(t, PieChart(PieChartProps{
			Slices:     slices,
			ShowLegend: false,
		}))},
		{Name: "pie_chart_custom_colors", HTML: utils.Render(t, PieChart(PieChartProps{
			Slices: []PieChartSlice{
				{Label: "High", Value: 70, Color: "text-emerald-600 dark:text-emerald-400"},
				{Label: "Med", Value: 20, Color: "text-amber-600 dark:text-amber-400"},
				{Label: "Low", Value: 10, Color: "text-rose-600 dark:text-rose-400"},
			},
		}))},
		{Name: "pie_chart_empty", HTML: utils.Render(t, PieChart(DefaultPieChartProps()))},
		{Name: "pie_chart_single_slice", HTML: utils.Render(t, PieChart(PieChartProps{
			Slices: []PieChartSlice{
				{Label: "All", Value: 100},
			},
		}))},
	})
}

func TestPieChartAriaLabel(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, PieChart(PieChartProps{
		BaseProps: utils.BaseProps{AriaLabel: "Traffic sources"},
		Slices:    []PieChartSlice{{Label: "A", Value: 1}},
	}))
	utils.AssertContains(t, html, `aria-label="Traffic sources"`)
	utils.AssertContains(t, html, `role="img"`)
}

func TestPieChartDecorativeAriaHidden(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, PieChart(PieChartProps{
		Slices: []PieChartSlice{{Label: "A", Value: 1}},
	}))
	utils.AssertContains(t, html, `aria-hidden="true"`)
}

func TestPieChartEmptyState(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, PieChart(DefaultPieChartProps()))
	utils.AssertContains(t, html, "No data")
	utils.AssertNotContains(t, html, "<path")
}

func TestPieChartBasePropsPropagation(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, PieChart(PieChartProps{
		BaseProps: utils.BaseProps{
			Class: "max-w-sm",
			ID:    "my-pie",
		},
		Slices: []PieChartSlice{{Label: "A", Value: 1}},
	}))
	utils.AssertContains(t, html, "max-w-sm")
	utils.AssertContains(t, html, `id="my-pie"`)
}

func TestComputeSliceAngles(t *testing.T) {
	t.Parallel()

	t.Run("normal distribution", func(t *testing.T) {
		t.Parallel()

		slices := []PieChartSlice{
			{Label: "A", Value: 25},
			{Label: "B", Value: 25},
			{Label: "C", Value: 50},
		}
		angles := computeSliceAngles(slices)

		if len(angles) != 3 {
			t.Fatalf("got %d angles, want 3", len(angles))
		}

		if angles[0].startAngle != 0 {
			t.Errorf("first start angle: got %v, want 0", angles[0].startAngle)
		}

		if math.Abs(angles[2].endAngle-360) > 0.01 {
			t.Errorf("last end angle: got %v, want 360", angles[2].endAngle)
		}

		// A=25% → 90°, B=25% → 90°, C=50% → 180°
		if math.Abs(angles[0].endAngle-90) > 0.01 {
			t.Errorf("A end angle: got %v, want 90", angles[0].endAngle)
		}
	})

	t.Run("zero values skipped", func(t *testing.T) {
		t.Parallel()

		slices := []PieChartSlice{
			{Label: "A", Value: 0},
			{Label: "B", Value: 100},
		}
		angles := computeSliceAngles(slices)

		if len(angles) != 1 {
			t.Fatalf("got %d angles, want 1 (zero skipped)", len(angles))
		}

		if angles[0].sliceIdx != 1 {
			t.Errorf("sliceIdx: got %d, want 1", angles[0].sliceIdx)
		}
	})

	t.Run("all zeros returns nil", func(t *testing.T) {
		t.Parallel()

		angles := computeSliceAngles([]PieChartSlice{
			{Label: "A", Value: 0},
		})
		if angles != nil {
			t.Errorf("got %v, want nil", angles)
		}
	})
}

func TestComputeArcPath(t *testing.T) {
	t.Parallel()

	t.Run("non-zero radius", func(t *testing.T) {
		t.Parallel()

		path := computeArcPath(100, 100, 50, 0, 0, 90)
		if path == "" {
			t.Error("expected non-empty path")
		}

		if path[:2] != "M " {
			t.Errorf("expected path to start with 'M ', got %q", path[:2])
		}
	})

	t.Run("zero radius returns empty", func(t *testing.T) {
		t.Parallel()

		path := computeArcPath(100, 100, 0, 0, 0, 90)
		if path != "" {
			t.Errorf("expected empty path, got %q", path)
		}
	})

	t.Run("full circle single slice", func(t *testing.T) {
		t.Parallel()

		path := computeArcPath(100, 100, 50, 0, 0, 360)
		if path == "" {
			t.Error("expected non-empty path for full circle")
		}

		// Full circle uses two semicircle arcs
		if countingStr(path, "A ") < 2 {
			t.Errorf("expected at least 2 arc segments for full circle, path: %q", path)
		}
	})
}

func TestPolarToCartesian(t *testing.T) {
	t.Parallel()

	// At 0 degrees (top of circle), point should be above center
	x, y := polarToCartesian(100, 100, 50, 0)
	if math.Abs(x-100) > 0.01 {
		t.Errorf("x at 0°: got %v, want 100", x)
	}

	if y >= 100 {
		t.Errorf("y at 0°: got %v, want < 100 (above center)", y)
	}

	// At 180 degrees (bottom), point should be below center
	x, y = polarToCartesian(100, 100, 50, 180)
	if math.Abs(x-100) > 0.01 {
		t.Errorf("x at 180°: got %v, want 100", x)
	}

	if y <= 100 {
		t.Errorf("y at 180°: got %v, want > 100 (below center)", y)
	}
}

func countingStr(s, sub string) int {
	count := 0

	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			count++
		}
	}

	return count
}
