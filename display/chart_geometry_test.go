package display

import (
	"math"
	"strings"
	"testing"
)

func TestScalePoints(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		values     []float64
		width      int
		height     int
		min        float64
		max        float64
		wantLen    int
		wantFirstX float64
		wantLastX  float64
	}{
		{
			name:       "basic 4 points",
			values:     []float64{0, 50, 100},
			width:      100,
			height:     100,
			min:        0,
			max:        100,
			wantLen:    3,
			wantFirstX: 0,
			wantLastX:  100,
		},
		{
			name:       "single point",
			values:     []float64{42},
			width:      200,
			height:     100,
			min:        0,
			max:        100,
			wantLen:    1,
			wantFirstX: 0,
			wantLastX:  0,
		},
		{
			name:       "empty",
			values:     []float64{},
			width:      100,
			height:     100,
			min:        0,
			max:        100,
			wantLen:    0,
			wantFirstX: 0,
			wantLastX:  0,
		},
		{
			name:       "min equals max (range padded)",
			values:     []float64{5, 5, 5},
			width:      100,
			height:     100,
			min:        5,
			max:        5,
			wantLen:    3,
			wantFirstX: 0,
			wantLastX:  100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pts := ScalePoints(tt.values, tt.width, tt.height, tt.min, tt.max)
			if len(pts) != tt.wantLen {
				t.Fatalf("got %d points, want %d", len(pts), tt.wantLen)
			}

			if len(pts) > 0 {
				if pts[0].X != tt.wantFirstX {
					t.Errorf("first X: got %v, want %v", pts[0].X, tt.wantFirstX)
				}

				if pts[len(pts)-1].X != tt.wantLastX {
					t.Errorf("last X: got %v, want %v", pts[len(pts)-1].X, tt.wantLastX)
				}
			}
		})
	}
}

func TestScalePointsYInversion(t *testing.T) {
	t.Parallel()

	// Value at max should map to Y=0 (top), value at min should map to Y=height (bottom).
	pts := ScalePoints([]float64{0, 100}, 100, 100, 0, 100)
	if pts[0].Y != 100 {
		t.Errorf("min value Y: got %v, want 100 (bottom)", pts[0].Y)
	}

	if pts[1].Y != 0 {
		t.Errorf("max value Y: got %v, want 0 (top)", pts[1].Y)
	}
}

func TestScalePointsClamps(t *testing.T) {
	t.Parallel()

	// Values outside [min, max] should be clamped.
	pts := ScalePoints([]float64{-50, 150}, 100, 100, 0, 100)
	if pts[0].Y != 100 {
		t.Errorf("below-min value Y: got %v, want 100 (clamped to bottom)", pts[0].Y)
	}

	if pts[1].Y != 0 {
		t.Errorf("above-max value Y: got %v, want 0 (clamped to top)", pts[1].Y)
	}
}

func TestBuildPolylinePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		points []Point
		want   string
	}{
		{"empty", []Point{}, ""},
		{"single", []Point{{X: 10, Y: 20}}, "M 10 20"},
		{"two points", []Point{{X: 0, Y: 0}, {X: 10, Y: 20}}, "M 0 0 L 10 20"},
		{
			"three points",
			[]Point{{X: 0, Y: 100}, {X: 50, Y: 50}, {X: 100, Y: 0}},
			"M 0 100 L 50 50 L 100 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := BuildPolylinePath(tt.points)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildSmoothPath(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		if got := BuildSmoothPath(nil); got != "" {
			t.Errorf("got %q, want empty", got)
		}
	})

	t.Run("two points falls back to polyline", func(t *testing.T) {
		t.Parallel()
		pts := []Point{{X: 0, Y: 0}, {X: 10, Y: 20}}
		got := BuildSmoothPath(pts)
		want := BuildPolylinePath(pts)
		if got != want {
			t.Errorf("got %q, want %q (polyline fallback)", got, want)
		}
	})

	t.Run("three points produces Bezier", func(t *testing.T) {
		t.Parallel()
		pts := []Point{{X: 0, Y: 0}, {X: 50, Y: 50}, {X: 100, Y: 0}}
		got := BuildSmoothPath(pts)
		if !strings.HasPrefix(got, "M 0 0 C ") {
			t.Errorf("expected Bezier curve, got %q", got)
		}
		if !strings.Contains(got, "100 0") {
			t.Errorf("expected path to end at (100, 0), got %q", got)
		}
	})
}

func TestBuildAreaPath(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		if got := BuildAreaPath(nil, 100); got != "" {
			t.Errorf("got %q, want empty", got)
		}
	})

	t.Run("closes to baseline", func(t *testing.T) {
		t.Parallel()
		pts := []Point{{X: 0, Y: 0}, {X: 50, Y: 50}, {X: 100, Y: 0}}
		got := BuildAreaPath(pts, 100)
		if !strings.HasSuffix(got, "Z") {
			t.Errorf("expected path to end with Z (closed), got %q", got)
		}

		if !strings.Contains(got, "100 100") {
			t.Errorf("expected baseline Y=100 in path, got %q", got)
		}
	})
}

func TestComputeNiceTicks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		min       float64
		max       float64
		count     int
		wantFirst float64
		wantLast  float64
		wantSpan  float64
	}{
		{"0 to 100", 0, 100, 5, 0, 100, 100},
		{"0 to 92 rounds up", 0, 92, 5, 0, 100, 100},
		{"10 to 90", 10, 90, 4, 0, 100, 100},
		{"negative range", -50, 50, 5, -60, 60, 120},
		{"small range", 0, 1, 5, 0, 1, 1},
		{"single value range", 5, 5, 5, 5, 6, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ticks := ComputeNiceTicks(tt.min, tt.max, tt.count)
			if len(ticks) < 2 {
				t.Fatalf("got %d ticks, want at least 2", len(ticks))
			}

			if ticks[0] != tt.wantFirst {
				t.Errorf("first tick: got %v, want %v", ticks[0], tt.wantFirst)
			}

			lastIdx := len(ticks) - 1
			if ticks[lastIdx] != tt.wantLast {
				t.Errorf("last tick: got %v, want %v", ticks[lastIdx], tt.wantLast)
			}

			span := ticks[lastIdx] - ticks[0]
			if span != tt.wantSpan {
				t.Errorf("span: got %v, want %v", span, tt.wantSpan)
			}

			// All ticks should be evenly spaced.
			step := ticks[1] - ticks[0]
			for i := 1; i < len(ticks); i++ {
				diff := ticks[i] - ticks[i-1]
				if math.Abs(diff-step) > 1e-9 {
					t.Errorf("tick %d-%d spacing: got %v, want %v", i-1, i, diff, step)
				}
			}
		})
	}
}

func TestFormatTickValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v    float64
		want string
	}{
		{"zero", 0, "0"},
		{"whole number", 42, "42"},
		{"negative whole", -17, "-17"},
		{"one decimal", 0.5, "0.5"},
		{"thousands", 15000, "15K"},
		{"thousands with decimal", 12500, "12.5K"},
		{"millions", 2000000, "2M"},
		{"millions with decimal", 1500000, "1.5M"},
		{"negative thousands", -50000, "-50K"},
		{"small number", 3.14, "3.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := FormatTickValue(tt.v)
			if got != tt.want {
				t.Errorf("FormatTickValue(%v): got %q, want %q", tt.v, got, tt.want)
			}
		})
	}
}

func BenchmarkScalePoints(b *testing.B) {
	values := make([]float64, 1000)
	for i := range values {
		values[i] = float64(i)
	}

	b.ResetTimer()

	for range b.N {
		ScalePoints(values, 600, 300, 0, 999)
	}
}

func BenchmarkBuildPolylinePath(b *testing.B) {
	points := make([]Point, 1000)
	for i := range points {
		points[i] = Point{X: float64(i), Y: float64(i % 100)}
	}

	b.ResetTimer()

	for range b.N {
		BuildPolylinePath(points)
	}
}
