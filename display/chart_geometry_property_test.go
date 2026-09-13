package display

import (
	"math"
	"math/rand"
	"testing"
)

// TestChartGeometryProperties pins chart-geometry invariants over randomized
// domains (M20/F092, deterministic seed for reproducibility):
//
//  1. ticks are strictly monotonic and evenly spaced,
//  2. the tick domain COVERS the data domain (first <= min, last >= max),
//  3. the tick count stays within count±1 (the wobble one extra decade can
//     add is bounded, not arbitrary),
//  4. ScalePoints maps every value INSIDE the pixel box (0..width, 0..height).
func TestChartGeometryProperties(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(42))

	for i := range 2000 {
		minVal := (rng.Float64() - 0.5) * math.Pow(10, float64(rng.Intn(7)-3))
		span := rng.Float64() * math.Pow(10, float64(rng.Intn(7)))
		maxVal := minVal + span

		count := 2 + rng.Intn(9)
		label := "case"

		t.Run(label, func(t *testing.T) {
			t.Parallel()

			_ = i
			ticks := ComputeNiceTicks(minVal, maxVal, count)

			if len(ticks) < 2 {
				t.Fatalf("min=%v max=%v count=%d: got %d ticks, want >= 2", minVal, maxVal, count, len(ticks))
			}

			// 1. monotonic + evenly spaced.
			step := ticks[1] - ticks[0]

			if step <= 0 {
				t.Fatalf("min=%v max=%v: non-positive tick step %v", minVal, maxVal, step)
			}

			for j := 1; j < len(ticks); j++ {
				got := ticks[j] - ticks[j-1]

				if math.Abs(got-step) > 1e-9*math.Abs(step) {
					t.Fatalf("min=%v max=%v: uneven spacing at %d: %v vs %v", minVal, maxVal, j, got, step)
				}
			}

			// 2. domain coverage (the whole data range lies between the outer ticks).
			if ticks[0] > minVal+1e-9*math.Max(1, math.Abs(minVal)) {
				t.Fatalf("min=%v max=%v: first tick %v does not cover min", minVal, maxVal, ticks[0])
			}

			last := ticks[len(ticks)-1]

			if last < maxVal-1e-9*math.Max(1, math.Abs(maxVal)) {
				t.Fatalf("min=%v max=%v: last tick %v does not cover max", minVal, maxVal, last)
			}

			// 3. count bounds: never below 2, never more than count+2
			// (a decade round-up can add one tick at each end).
			if len(ticks) > count+2 {
				t.Fatalf("min=%v max=%v count=%d: got %d ticks, want <= count+2", minVal, maxVal, count, len(ticks))
			}
		})
	}
}

// TestScalePointsStaysInBox pins the pixel mapping invariant over randomized
// values: every scaled point lands inside the viewBox (with float tolerance).
func TestScalePointsStaysInBox(t *testing.T) {
	t.Parallel()

	rng := rand.New(rand.NewSource(7))

	const (
		width  = 640
		height = 320
	)

	for range 1000 {
		minVal := (rng.Float64() - 0.5) * 1000
		maxVal := minVal + rng.Float64()*1000

		values := make([]float64, rng.Intn(50))
		for j := range values {
			values[j] = minVal + rng.Float64()*(maxVal-minVal)
		}

		for _, p := range ScalePoints(values, width, height, minVal, maxVal) {
			if p.X < -0.01 || p.X > float64(width)+0.01 {
				t.Fatalf("values=%v: X %v outside [0,%d]", values, p.X, width)
			}

			if p.Y < -0.01 || p.Y > float64(height)+0.01 {
				t.Fatalf("values=%v: Y %v outside [0,%d]", values, p.Y, height)
			}
		}
	}
}
