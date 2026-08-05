package display

import (
	"bytes"
	"context"
	"testing"
)

// BenchmarkComputeSliceAngles measures the performance of slice angle
// computation for a pie chart with 100 slices.
func BenchmarkComputeSliceAngles(b *testing.B) {
	slices := make([]PieChartSlice, 100)
	for i := range slices {
		slices[i] = PieChartSlice{Label: "S", Value: float64(i + 1)}
	}

	b.ResetTimer()

	for b.Loop() {
		computeSliceAngles(slices)
	}
}

// BenchmarkComputeArcPath measures the performance of SVG arc path computation
// for 100 pie slices.
func BenchmarkComputeArcPath(b *testing.B) {
	slices := make([]PieChartSlice, 100)
	for i := range slices {
		slices[i] = PieChartSlice{Label: "S", Value: float64(i + 1)}
	}
	angles := computeSliceAngles(slices)

	b.ResetTimer()

	for b.Loop() {
		for _, a := range angles {
			computeArcPath(200, 150, 100, 60, a.startAngle, a.endAngle)
		}
	}
}

// BenchmarkLineChartRender measures the performance of rendering a full
// LineChart component to a buffer with two series of 50 data points each.
func BenchmarkLineChartRender(b *testing.B) {
	props := DefaultLineChartProps()
	props.Series = []LineChartSeries{
		{Name: "Revenue", Values: make([]float64, 50)},
		{Name: "Costs", Values: make([]float64, 50)},
	}
	for i := range props.Series[0].Values {
		props.Series[0].Values[i] = float64(i * 10)
		props.Series[1].Values[i] = float64(i * 5)
	}

	component := LineChart(props)
	ctx := context.Background()

	b.ResetTimer()

	for b.Loop() {
		var buf bytes.Buffer
		_ = component.Render(ctx, &buf)
	}
}
