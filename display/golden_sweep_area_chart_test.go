package display

import (
	"testing"

	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

func TestGoldenSweepAreaChart(t *testing.T) {
	t.Parallel()

	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "area_chart_single_series", HTML: utils.Render(t, AreaChart(AreaChartProps{
			BaseProps: utils.BaseProps{AriaLabel: "Active users by month"},
			Series: []LineChartSeries{
				{Name: "Users", Values: []float64{120, 180, 250, 300, 280, 340}},
			},
			XAxisLabels: months,
		}))},
		{Name: "area_chart_multi_series", HTML: utils.Render(t, AreaChart(AreaChartProps{
			Series: []LineChartSeries{
				{Name: "Active", Values: []float64{120, 180, 250, 300, 280, 340}},
				{Name: "Churned", Values: []float64{20, 30, 25, 40, 35, 50}},
			},
			XAxisLabels: months,
		}))},
		{Name: "area_chart_smooth", HTML: utils.Render(t, AreaChart(AreaChartProps{
			Series: []LineChartSeries{
				{Name: "Trend", Values: []float64{120, 180, 250, 300, 280, 340}},
			},
			XAxisLabels: months,
			Style:       LineChartStyleSmooth,
		}))},
		{Name: "area_chart_with_dots", HTML: utils.Render(t, AreaChart(AreaChartProps{
			Series: []LineChartSeries{
				{Name: "Data", Values: []float64{120, 180, 250, 300, 280, 340}},
			},
			XAxisLabels: months,
			ShowDots:    true,
		}))},
		{Name: "area_chart_custom_opacity", HTML: utils.Render(t, AreaChart(AreaChartProps{
			Series: []LineChartSeries{
				{Name: "Data", Values: []float64{120, 180, 250, 300, 280, 340}},
			},
			XAxisLabels: months,
			FillOpacity: 0.4,
		}))},
		{Name: "area_chart_custom_color", HTML: utils.Render(t, AreaChart(AreaChartProps{
			Series: []LineChartSeries{
				{
					Name:   "Growth",
					Values: []float64{120, 180, 250, 300, 280, 340},
					Color:  "text-emerald-600 dark:text-emerald-400",
				},
			},
			XAxisLabels: months,
		}))},
		{Name: "area_chart_empty", HTML: utils.Render(t, AreaChart(DefaultAreaChartProps()))},
	})
}

func TestAreaChartAriaLabel(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, AreaChart(AreaChartProps{
		BaseProps:   utils.BaseProps{AriaLabel: "User growth"},
		Series:      []LineChartSeries{{Name: "Users", Values: []float64{1, 2, 3}}},
		XAxisLabels: []string{"A", "B", "C"},
	}))
	utils.AssertContains(t, html, `aria-label="User growth"`)
	utils.AssertContains(t, html, `role="img"`)
}

func TestAreaChartDecorativeAriaHidden(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, AreaChart(AreaChartProps{
		Series: []LineChartSeries{{Name: "Data", Values: []float64{1, 2, 3}}},
	}))
	utils.AssertContains(t, html, `aria-hidden="true"`)
}

func TestAreaChartEmptyState(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, AreaChart(DefaultAreaChartProps()))
	utils.AssertContains(t, html, "No data")
	utils.AssertNotContains(t, html, "<path")
}

func TestAreaChartBasePropsPropagation(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, AreaChart(AreaChartProps{
		BaseProps: utils.BaseProps{
			Class: "max-w-xl",
			ID:    "area-chart",
		},
		Series: []LineChartSeries{{Values: []float64{1, 2, 3}}},
	}))
	utils.AssertContains(t, html, "max-w-xl")
	utils.AssertContains(t, html, `id="area-chart"`)
}
