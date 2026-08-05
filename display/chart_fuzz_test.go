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

	f.Fuzz(func(t *testing.T, cx float64, cy float64, radius float64, innerRadius float64, startAngle float64, endAngle float64) {
		// Should never panic regardless of input.
		_ = computeArcPath(cx, cy, radius, innerRadius, startAngle, endAngle)
	})
}
