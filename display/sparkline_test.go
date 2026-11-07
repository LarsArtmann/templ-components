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
		AriaLabel: "Trend: increasing",
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

// TestSparklinePointsStayWithinCanvas pins the invariant that every projected
// x coordinate lies in [0, width]: the last point must land exactly on the
// right canvas edge, never past it. Step-rounding (i * round(step)) put the
// final point OUTSIDE the SVG for widths not divisible by the point count
// (e.g. 8 values at width 200 landed at x=203).
func TestSparklinePointsStayWithinCanvas(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		values []float64
		width  int
	}{
		{name: "8 values at width 200", values: []float64{1, 3, 2, 5, 4, 6, 3, 7}, width: 200},
		{name: "7 values at width 100", values: []float64{1, 3, 2, 5, 4, 6, 3}, width: 100},
		{name: "10 values at width 120", values: []float64{3, 5, 4, 7, 6, 9, 8, 12, 10, 14}, width: 120},
		{name: "3 values at width 37", values: []float64{1, 2, 3}, width: 37},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			coords, ok := sparklinePointCoords(tt.values, tt.width, 30, 0, 10)
			if !ok {
				t.Fatalf("sparklinePointCoords returned ok=false for %d values", len(tt.values))
			}

			for i, c := range coords {
				if c.x < 0 || c.x > tt.width {
					t.Errorf("point %d: x=%d outside canvas [0,%d]", i, c.x, tt.width)
				}
			}

			if last := coords[len(coords)-1]; last.x != tt.width {
				t.Errorf("last point x=%d, want exactly %d (right canvas edge)", last.x, tt.width)
			}
		})
	}
}

func TestSparklineBounds(t *testing.T) {
	t.Parallel()
	t.Run("auto from data", func(t *testing.T) {
		t.Parallel()

		minVal, maxVal := sparklineBounds(SparklineProps{Values: []float64{3, 1, 4, 1, 5, 9, 2, 6}})
		if minVal != 1 {
			t.Errorf("expected min=1, got %f", minVal)
		}

		if maxVal != 9 {
			t.Errorf("expected max=9, got %f", maxVal)
		}
	})
	t.Run("equal values bump max", func(t *testing.T) {
		t.Parallel()

		minVal, maxVal := sparklineBounds(SparklineProps{Values: []float64{5, 5, 5}})
		if minVal != 5 {
			t.Errorf("expected min=5, got %f", minVal)
		}

		if maxVal != 6 {
			t.Errorf("expected max=6 (bumped), got %f", maxVal)
		}
	})
	t.Run("zero min is respected not treated as auto", func(t *testing.T) {
		t.Parallel()

		zero := 0.0
		minVal, maxVal := sparklineBounds(SparklineProps{Values: []float64{3, 1, 4, 1, 5}, Min: &zero})

		if minVal != 0 {
			t.Errorf("expected min=0 (explicit), got %f", minVal)
		}

		if maxVal != 5 {
			t.Errorf("expected max=5 (auto from data), got %f", maxVal)
		}
	})
	t.Run("explicit min and max overrides", func(t *testing.T) {
		t.Parallel()

		minIn, maxIn := -10.0, 10.0
		minVal, maxVal := sparklineBounds(SparklineProps{Values: []float64{3, 1, 4, 1, 5}, Min: &minIn, Max: &maxIn})

		if minVal != -10 {
			t.Errorf("expected min=-10, got %f", minVal)
		}

		if maxVal != 10 {
			t.Errorf("expected max=10, got %f", maxVal)
		}
	})
}

func TestSparklineZeroAwareBounds(t *testing.T) {
	t.Parallel()
	t.Run("all-positive widens min to 0", func(t *testing.T) {
		t.Parallel()

		minVal, maxVal := sparklineBounds(SparklineProps{Values: []float64{3, 1, 5}, ZeroAware: true})
		if minVal != 0 {
			t.Errorf("expected min=0, got %f", minVal)
		}

		if maxVal != 5 {
			t.Errorf("expected max=5, got %f", maxVal)
		}
	})
	t.Run("all-negative widens max to 0", func(t *testing.T) {
		t.Parallel()

		minVal, maxVal := sparklineBounds(SparklineProps{Values: []float64{-3, -1, -5}, ZeroAware: true})
		if minVal != -5 {
			t.Errorf("expected min=-5, got %f", minVal)
		}

		if maxVal != 0 {
			t.Errorf("expected max=0, got %f", maxVal)
		}
	})
	t.Run("straddling series unchanged", func(t *testing.T) {
		t.Parallel()

		minVal, maxVal := sparklineBounds(SparklineProps{Values: []float64{-3, 5}, ZeroAware: true})
		if minVal != -3 {
			t.Errorf("expected min=-3, got %f", minVal)
		}

		if maxVal != 5 {
			t.Errorf("expected max=5, got %f", maxVal)
		}
	})
	t.Run("ShowZeroLine implies zero-aware", func(t *testing.T) {
		t.Parallel()

		minVal, maxVal := sparklineBounds(SparklineProps{Values: []float64{3, 1, 5}, ShowZeroLine: true})
		if minVal != 0 {
			t.Errorf("expected min=0 (implied), got %f", minVal)
		}

		if maxVal != 5 {
			t.Errorf("expected max=5, got %f", maxVal)
		}
	})
	t.Run("explicit overrides win over widening", func(t *testing.T) {
		t.Parallel()

		floor := 1.0

		minVal, maxVal := sparklineBounds(SparklineProps{Values: []float64{3, 1, 5}, ZeroAware: true, Min: &floor})
		if minVal != 1 {
			t.Errorf("expected min=1 (explicit), got %f", minVal)
		}

		if maxVal != 5 {
			t.Errorf("expected max=5, got %f", maxVal)
		}
	})
}

func TestSparklineZeroY(t *testing.T) {
	t.Parallel()
	// Domain [-10, 10] at height 40: zero is the vertical midpoint.
	if got := sparklineZeroY(-10, 10, 40); got != 20 {
		t.Errorf("expected zeroY=20 for straddling domain, got %d", got)
	}
	// Domain [0, 10]: zero sits on the bottom edge.
	if got := sparklineZeroY(0, 10, 40); got != 40 {
		t.Errorf("expected zeroY=40 for non-negative domain, got %d", got)
	}
	// Domain [-10, 0]: zero sits on the top edge.
	if got := sparklineZeroY(-10, 0, 40); got != 0 {
		t.Errorf("expected zeroY=0 for non-positive domain, got %d", got)
	}
}

func TestSparklineResponsiveMarkup(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Sparkline(SparklineProps{
		Values:     []float64{1, 2, 3},
		Responsive: true,
	}))
	utils.AssertContains(t, output, `viewBox="0 0 120 30"`)
	utils.AssertContains(t, output, `preserveAspectRatio="none"`)
	utils.AssertContainsAll(t, output, "w-full", "h-auto")
	utils.AssertContains(t, output, `vector-effect="non-scaling-stroke"`)
	utils.AssertNotContains(t, output, `width="120"`)
}

func TestSparklineFixedModeUnchanged(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Sparkline(SparklineProps{Values: []float64{1, 2, 3}}))
	utils.AssertContains(t, output, `width="120"`)
	utils.AssertContains(t, output, `height="30"`)
	utils.AssertNotContains(t, output, "viewBox")
	utils.AssertNotContains(t, output, "vector-effect")
	utils.AssertNotContains(t, output, "preserveAspectRatio")
}

func TestSparklineZeroLineMarkup(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Sparkline(SparklineProps{
		Values:       []float64{-5, 5},
		ShowZeroLine: true,
	}))
	utils.AssertContains(t, output, "<line")
	utils.AssertContains(t, output, `stroke-opacity="0.25"`)
	// Zero is the vertical midpoint of a straddling domain on the 30px canvas.
	utils.AssertContains(t, output, `y1="15"`)
	plain := utils.Render(t, Sparkline(SparklineProps{Values: []float64{-5, 5}}))
	utils.AssertNotContains(t, plain, "<line")
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
