package display

import (
	"math"
	"testing"
)

// FuzzScalePoints verifies ScalePoints never panics on adversarial inputs
// (NaN, Inf, negative, very large). Uses fixed value slices since the Go fuzz
// corpus only supports scalar types.
func FuzzScalePoints(f *testing.F) {
	valuesSets := [][]float64{
		{0, 50, 100},
		{},
		{math.NaN()},
		{math.Inf(1)},
		{-1e20, 1e20},
	}

	f.Add(600, 300, 0.0, 100.0)
	f.Add(-100, -100, math.Inf(-1), math.Inf(1))
	f.Add(0, 0, math.NaN(), math.NaN())

	f.Fuzz(func(t *testing.T, width int, height int, minVal float64, maxVal float64) {
		for _, values := range valuesSets {
			// Should never panic regardless of input.
			ScalePoints(values, width, height, minVal, maxVal)
		}
	})
}

// FuzzComputeNiceTicks verifies ComputeNiceTicks never panics on adversarial
// inputs (NaN, Inf, equal min/max, negative count, very large range).
func FuzzComputeNiceTicks(f *testing.F) {
	f.Add(0.0, 100.0, 5)
	f.Add(math.NaN(), 100.0, 5)
	f.Add(0.0, math.NaN(), 5)
	f.Add(math.Inf(-1), math.Inf(1), -1)
	f.Add(50.0, 50.0, 0)
	f.Add(-1e20, 1e20, 1000)

	f.Fuzz(func(t *testing.T, minVal float64, maxVal float64, count int) {
		// Should never panic regardless of input.
		ticks := ComputeNiceTicks(minVal, maxVal, count)
		// All ticks must be finite (no NaN/Inf in output).
		for _, tick := range ticks {
			if math.IsNaN(tick) || math.IsInf(tick, 0) {
				t.Errorf("ComputeNiceTicks(%f, %f, %d) produced non-finite tick: %f", minVal, maxVal, count, tick)
			}
		}
	})
}

// FuzzComputeArcPath verifies computeArcPath never panics on adversarial
// inputs (NaN, Inf, negative radius, zero inner, extreme angles).
func FuzzComputeArcPath(f *testing.F) {
	f.Add(200.0, 150.0, 100.0, 60.0, 0.0, 90.0)
	f.Add(200.0, 150.0, 0.0, 0.0, 0.0, 360.0)
	f.Add(200.0, 150.0, -10.0, -5.0, -720.0, 720.0)
	f.Add(math.NaN(), math.NaN(), math.NaN(), math.NaN(), math.NaN(), math.NaN())
	f.Add(200.0, 150.0, 100.0, 200.0, 0.0, 360.0)

	f.Fuzz(
		func(t *testing.T, cx float64, cy float64, radius float64, innerRadius float64, startAngle float64, endAngle float64) {
			// Should never panic regardless of input.
			_ = computeArcPath(cx, cy, radius, innerRadius, startAngle, endAngle)
		},
	)
}

// FuzzBuildSmoothPath verifies BuildSmoothPath never panics on adversarial
// inputs (NaN, Inf, empty, single point, collinear, extreme coordinates).
func FuzzBuildSmoothPath(f *testing.F) {
	pointSets := [][]Point{
		{},
		{{X: 0, Y: 0}},
		{{X: 0, Y: 0}, {X: 100, Y: 50}},
		{{X: math.NaN(), Y: 0}, {X: 50, Y: math.Inf(1)}},
		{{X: -1e20, Y: 1e20}, {X: 1e20, Y: -1e20}, {X: 0, Y: 0}},
	}

	f.Add(0.0, 0.0, 100.0, 50.0)
	f.Add(math.NaN(), math.Inf(-1), math.Inf(1), math.NaN())
	f.Add(-1e20, 1e20, 1e20, -1e20)

	f.Fuzz(func(t *testing.T, x1 float64, y1 float64, x2 float64, y2 float64) {
		dynamicSet := []Point{{X: x1, Y: y1}, {X: x2, Y: y2}, {X: x1 + x2, Y: y1 + y2}}

		for _, points := range pointSets {
			// Should never panic regardless of input.
			_ = BuildSmoothPath(points)
		}

		// Should never panic regardless of input.
		_ = BuildSmoothPath(dynamicSet)
	})
}

// FuzzBuildAreaPath verifies BuildAreaPath never panics on adversarial inputs
// (NaN, Inf, empty, negative height, extreme coordinates).
func FuzzBuildAreaPath(f *testing.F) {
	pointSets := [][]Point{
		{},
		{{X: 0, Y: 0}},
		{{X: 0, Y: 0}, {X: 100, Y: 50}},
		{{X: math.NaN(), Y: 0}, {X: 50, Y: math.Inf(1)}},
	}

	f.Add(0.0, 0.0, 100.0, 50.0, 300)
	f.Add(math.NaN(), math.Inf(-1), 0.0, 0.0, -100)
	f.Add(-1e20, 1e20, 1e20, -1e20, 0)

	f.Fuzz(func(t *testing.T, x1 float64, y1 float64, x2 float64, y2 float64, height int) {
		dynamicSet := []Point{{X: x1, Y: y1}, {X: x2, Y: y2}}

		for _, points := range pointSets {
			// Should never panic regardless of input.
			_ = BuildAreaPath(points, height)
		}

		// Should never panic regardless of input.
		_ = BuildAreaPath(dynamicSet, height)
	})
}
