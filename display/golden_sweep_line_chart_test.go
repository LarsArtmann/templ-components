package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils/golden"
	"github.com/larsartmann/templ-components/utils"
)

func TestGoldenSweepLineChart(t *testing.T) {
	t.Parallel()

	weekdays := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "line_chart_single_series", HTML: utils.Render(t, LineChart(LineChartProps{
			BaseProps: utils.BaseProps{AriaLabel: "Revenue by day"},
			Series: []LineChartSeries{
				{Name: "Revenue", Values: []float64{10, 25, 40, 35, 60, 55, 80}},
			},
			XAxisLabels: weekdays,
		}))},
		{Name: "line_chart_multi_series", HTML: utils.Render(t, LineChart(LineChartProps{
			Series: []LineChartSeries{
				{Name: "Revenue", Values: []float64{10, 25, 40, 35, 60, 55, 80}},
				{Name: "Expenses", Values: []float64{8, 15, 20, 30, 35, 40, 50}},
			},
			XAxisLabels: weekdays,
		}))},
		{Name: "line_chart_smooth", HTML: utils.Render(t, LineChart(LineChartProps{
			Series: []LineChartSeries{
				{Name: "Trend", Values: []float64{10, 25, 40, 35, 60, 55, 80}},
			},
			XAxisLabels: weekdays,
			Style:       LineChartStyleSmooth,
		}))},
		{Name: "line_chart_no_grid", HTML: utils.Render(t, LineChart(LineChartProps{
			Series: []LineChartSeries{
				{Name: "Data", Values: []float64{10, 25, 40, 35, 60, 55, 80}},
			},
			XAxisLabels: weekdays,
			ShowGrid:    false,
		}))},
		{Name: "line_chart_no_dots", HTML: utils.Render(t, LineChart(LineChartProps{
			Series: []LineChartSeries{
				{Name: "Data", Values: []float64{10, 25, 40, 35, 60, 55, 80}},
			},
			XAxisLabels: weekdays,
			ShowDots:    false,
		}))},
		{Name: "line_chart_custom_color", HTML: utils.Render(t, LineChart(LineChartProps{
			Series: []LineChartSeries{
				{
					Name:   "Growth",
					Values: []float64{5, 10, 20, 35, 55, 80, 110},
					Color:  "text-emerald-600 dark:text-emerald-400",
				},
			},
			XAxisLabels: weekdays,
		}))},
		{Name: "line_chart_dashed", HTML: utils.Render(t, LineChart(LineChartProps{
			Series: []LineChartSeries{
				{Name: "Target", Values: []float64{50, 50, 50, 50, 50, 50, 50}, Dashed: true},
				{Name: "Actual", Values: []float64{30, 40, 55, 45, 60, 70, 65}},
			},
			XAxisLabels: weekdays,
		}))},
		{Name: "line_chart_empty", HTML: utils.Render(t, LineChart(DefaultLineChartProps()))},
		{Name: "line_chart_no_labels", HTML: utils.Render(t, LineChart(LineChartProps{
			Series: []LineChartSeries{
				{Name: "Data", Values: []float64{10, 25, 40, 35, 60, 55, 80}},
			},
		}))},
		{Name: "line_chart_decorative", HTML: utils.Render(t, LineChart(LineChartProps{
			Series: []LineChartSeries{
				{Name: "Data", Values: []float64{10, 25, 40, 35, 60, 55, 80}},
			},
		}))},
	})
}

func TestLineChartAriaLabel(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, LineChart(LineChartProps{
		BaseProps:   utils.BaseProps{AriaLabel: "Revenue trend"},
		Series:      []LineChartSeries{{Name: "Rev", Values: []float64{1, 2, 3}}},
		XAxisLabels: []string{"A", "B", "C"},
	}))
	utils.AssertContains(t, html, `aria-label="Revenue trend"`)
	utils.AssertContains(t, html, `role="img"`)
}

func TestLineChartDecorativeAriaHidden(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, LineChart(LineChartProps{
		Series:      []LineChartSeries{{Name: "Rev", Values: []float64{1, 2, 3}}},
		XAxisLabels: []string{"A", "B", "C"},
	}))
	utils.AssertContains(t, html, `aria-hidden="true"`)
}

func TestLineChartEmptyState(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, LineChart(DefaultLineChartProps()))
	utils.AssertContains(t, html, "No data")
	utils.AssertNotContains(t, html, "<path")
}

func TestLineChartMinOverride(t *testing.T) {
	t.Parallel()

	minVal := 0.0
	html := utils.Render(t, LineChart(LineChartProps{
		Series: []LineChartSeries{{Values: []float64{10, 20, 30}}},
		Min:    &minVal,
	}))
	utils.AssertContains(t, html, ">0<")
}

func TestLineChartCustomValueFormat(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, LineChart(LineChartProps{
		Series: []LineChartSeries{{Values: []float64{10, 20, 30}}},
		ValueFormat: func(v float64) string {
			return "$" + FormatTickValue(v)
		},
	}))
	utils.AssertContains(t, html, "$")
}

func TestLineChartBasePropsPropagation(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, LineChart(LineChartProps{
		BaseProps: utils.BaseProps{
			Class: "max-w-2xl",
			ID:    "my-chart",
		},
		Series: []LineChartSeries{{Values: []float64{1, 2, 3}}},
	}))
	utils.AssertContains(t, html, "max-w-2xl")
	utils.AssertContains(t, html, `id="my-chart"`)
}
