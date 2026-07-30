package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestSparklineRender(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Sparkline(SparklineProps{
		Values: []float64{1, 3, 2, 5, 4},
	}))
	utils.AssertContains(t, output, "<svg")
	utils.AssertContains(t, output, "<polyline")
	utils.AssertContains(t, output, "stroke=\"currentColor\"")
	utils.AssertContains(t, output, `aria-hidden="true"`)
}

func TestSparklineFilled(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Sparkline(SparklineProps{
		Values: []float64{1, 3, 2, 5, 4},
		Filled: true,
	}))
	utils.AssertContains(t, output, "<path")
	utils.AssertContains(t, output, "fill=\"currentColor\"")
}

func TestSparklineAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Sparkline(SparklineProps{
		Values:    []float64{1, 2, 3},
		BaseProps: utils.BaseProps{AriaLabel: "Trend: increasing"},
	}))
	utils.AssertContains(t, output, `aria-label="Trend: increasing"`)
	utils.AssertNotContains(t, output, `aria-hidden`)
}

func TestSparklineCustomDimensions(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Sparkline(SparklineProps{
		Values:      []float64{1, 2, 3},
		Width:       200,
		Height:      40,
		StrokeWidth: 2,
	}))
	utils.AssertContains(t, output, `width="200"`)
	utils.AssertContains(t, output, `height="40"`)
	utils.AssertContains(t, output, `stroke-width="2"`)
}

func TestSparklineSingleValueRendersNothing(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Sparkline(SparklineProps{
		Values: []float64{42},
	}))
	utils.AssertNotContains(t, output, "<svg")
}

func TestSparklineEmptyRendersNothing(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Sparkline(SparklineProps{
		Values: []float64{},
	}))
	utils.AssertNotContains(t, output, "<svg")
}

func TestSparklinePointsComputes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		values []float64
		width  int
		height int
		min    float64
		max    float64
		want   string
	}{
		{
			name:   "ascending",
			values: []float64{0, 1, 2},
			width:  100,
			height: 20,
			min:    0,
			max:    2,
			want:   "0,20 50,10 100,0",
		},
		{
			name:   "descending",
			values: []float64{2, 1, 0},
			width:  100,
			height: 20,
			min:    0,
			max:    2,
			want:   "0,0 50,10 100,20",
		},
		{
			name:   "flat line",
			values: []float64{1, 1, 1},
			width:  100,
			height: 20,
			min:    0,
			max:    1,
			want:   "0,0 50,0 100,0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := sparklinePoints(tt.values, tt.width, tt.height, tt.min, tt.max)
			if got != tt.want {
				t.Errorf("sparklinePoints() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSparklineBounds(t *testing.T) {
	t.Parallel()
	t.Run("auto from data", func(t *testing.T) {
		t.Parallel()
		min, max := sparklineBounds(SparklineProps{Values: []float64{3, 1, 4, 1, 5, 9, 2, 6}})
		if min != 1 {
			t.Errorf("expected min=1, got %f", min)
		}
		if max != 9 {
			t.Errorf("expected max=9, got %f", max)
		}
	})
	t.Run("equal values bump max", func(t *testing.T) {
		t.Parallel()
		min, max := sparklineBounds(SparklineProps{Values: []float64{5, 5, 5}})
		if min != 5 {
			t.Errorf("expected min=5, got %f", min)
		}
		if max != 6 {
			t.Errorf("expected max=6 (bumped), got %f", max)
		}
	})
}

func TestSparklineDefaults(t *testing.T) {
	t.Parallel()
	props := DefaultSparklineProps()
	if props.Width != 120 {
		t.Errorf("expected Width=120, got %d", props.Width)
	}
	if props.Height != 30 {
		t.Errorf("expected Height=30, got %d", props.Height)
	}
	if props.StrokeWidth != 1.5 {
		t.Errorf("expected StrokeWidth=1.5, got %f", props.StrokeWidth)
	}
}
