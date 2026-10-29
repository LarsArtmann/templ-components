package display

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

func TestGoldenSweepAreaChart(t *testing.T) {
	t.Parallel()

	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "area_chart_single_series", HTML: utils.Render(t, AreaChart(AreaChartProps{
			AriaLabel: "Active users by month",
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
		AriaLabel:   "User growth",
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
		Class:  "max-w-xl",
		ID:     "area-chart",
		Series: []LineChartSeries{{Values: []float64{1, 2, 3}}},
	}))
	utils.AssertContains(t, html, "max-w-xl")
	utils.AssertContains(t, html, `id="area-chart"`)
}

// areaPathD extracts the d attribute of the filled area path (fill with
// opacity, no stroke) from rendered chart HTML.
func areaPathD(t *testing.T, html string) string {
	t.Helper()

	const marker = `fill="currentColor" fill-opacity="0.2" stroke="none"`

	before, _, found := strings.Cut(html, marker)
	if !found {
		t.Fatalf("no area fill path in output: %s", html)
	}

	const dPrefix = `<path d="`

	_, rest, found := strings.Cut(before, dPrefix)
	if !found {
		t.Fatalf("no path d attribute before area fill marker: %s", html)
	}

	d, _, _ := strings.Cut(rest, `"`)

	return d
}

// TestAreaChartTicksStayInsidePlotDomain pins the fix for ticks rendering
// OUTSIDE the plot area: ComputeNiceTicks intentionally covers the data
// domain (first tick <= min, last tick >= max — pinned by the geometry
// property test), but chartAxes projected those outer ticks using the DATA
// range, so e.g. data 120..340 drew a "100" gridline below the x-axis (in
// the x-label zone) and a "350" gridline above the plot top. Out-of-domain
// ticks must not render; in-domain ticks must.
func TestAreaChartTicksStayInsidePlotDomain(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, AreaChart(AreaChartProps{
		Series: []LineChartSeries{
			{Name: "Users", Values: []float64{120, 180, 250, 300, 280, 340}},
		},
	}))

	// ComputeNiceTicks(120, 340, 8) yields 100, 150, 200, 250, 300, 350:
	// 100 and 350 are outside the data domain and must not render.
	utils.AssertNotContains(t, html, ">100</text>")
	utils.AssertNotContains(t, html, ">350</text>")

	// In-domain ticks still render.
	utils.AssertContains(t, html, ">150</text>")
	utils.AssertContains(t, html, ">300</text>")
}

// TestAreaChartSmoothFillFollowsLine pins the fix for the filled area edge
// not following a smooth (curved) line: with Style smooth the stroked line
// used a Catmull-Rom spline while the fill below it stayed a straight
// polyline, so the curve visibly detached from its own fill. The area path
// must carry the same C (cubic) segments as the line.
func TestAreaChartSmoothFillFollowsLine(t *testing.T) {
	t.Parallel()

	html := utils.Render(t, AreaChart(AreaChartProps{
		Series: []LineChartSeries{
			{Name: "Trend", Values: []float64{120, 180, 250, 300, 280, 340}},
		},
		Style: LineChartStyleSmooth,
	}))

	if d := areaPathD(t, html); !strings.Contains(d, " C ") {
		t.Errorf("smooth area path has no curve segments (straight fill under a curved line): %s", d)
	}
}
