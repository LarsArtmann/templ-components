package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

// Golden sweep for the chart/visualization components: Sparkline, BarChart,
// ExternalLink. These previously lacked golden snapshot tests.

func TestGoldenSweepSparkline(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "sparkline_basic", HTML: utils.Render(t, Sparkline(SparklineProps{
			Values: []float64{1, 3, 2, 5, 4, 6, 3, 7},
		}))},
		{Name: "sparkline_filled", HTML: utils.Render(t, Sparkline(SparklineProps{
			Values: []float64{1, 3, 2, 5, 4, 6, 3, 7},
			Filled: true,
		}))},
		{Name: "sparkline_aria_label", HTML: utils.Render(t, Sparkline(SparklineProps{
			Values:    []float64{1, 3, 2, 5, 4, 6, 3, 7},
			AriaLabel: "Message rate over 8 hours",
		}))},
		{Name: "sparkline_custom_dims", HTML: utils.Render(t, Sparkline(SparklineProps{
			Values:      []float64{1, 3, 2, 5, 4, 6, 3, 7},
			Width:       200,
			Height:      40,
			StrokeWidth: 2,
		}))},
		{Name: "sparkline_empty", HTML: utils.Render(t, Sparkline(SparklineProps{}))},
	})
}

func TestGoldenSweepBarChart(t *testing.T) {
	t.Parallel()

	bars := []BarChartBar{
		{Label: "general", Value: 1200},
		{Label: "random", Value: 800},
		{Label: "dev", Value: 450},
		{Label: "gaming", Value: 200},
	}

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "bar_chart_horizontal", HTML: utils.Render(t, BarChart(BarChartProps{
			Bars: bars,
		}))},
		{Name: "bar_chart_vertical", HTML: utils.Render(t, BarChart(BarChartProps{
			Bars:   bars,
			Orient: BarVertical,
		}))},
		{Name: "bar_chart_empty", HTML: utils.Render(t, BarChart(DefaultBarChartProps()))},
		{Name: "bar_chart_custom_color", HTML: utils.Render(t, BarChart(BarChartProps{
			Bars: []BarChartBar{
				{Label: "high", Value: 900, Color: "bg-emerald-600 dark:bg-emerald-500"},
				{Label: "med", Value: 500, Color: "bg-amber-600 dark:bg-amber-500"},
				{Label: "low", Value: 100, Color: "bg-red-600 dark:bg-red-500"},
			},
		}))},
		{Name: "bar_chart_with_href", HTML: utils.Render(t, BarChart(BarChartProps{
			Bars: []BarChartBar{
				{Label: "general", Value: 1200, Href: "/channels/general"},
				{Label: "random", Value: 800, Href: "/channels/random"},
			},
		}))},
		{Name: "bar_chart_vertical_signed", HTML: utils.Render(t, BarChart(BarChartProps{
			Bars: []BarChartBar{
				{Label: "mar", Value: 420},
				{Label: "apr", Value: -180},
				{Label: "may", Value: 90},
				{Label: "jun", Value: -60},
			},
			Orient:     BarVertical,
			Height:     "10rem",
			ShowValues: true,
		}))},
		{Name: "bar_chart_vertical_capped", HTML: utils.Render(t, BarChart(BarChartProps{
			Bars:        bars[:2],
			Orient:      BarVertical,
			MaxBarWidth: "max-w-24",
		}))},
		{Name: "bar_chart_vertical_baseline", HTML: utils.Render(t, BarChart(BarChartProps{
			Bars:         bars[:3],
			Orient:       BarVertical,
			ShowBaseline: true,
		}))},
	})
}

func TestGoldenSweepExternalLink(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "external_link_text", HTML: utils.Render(t, ExternalLink(ExternalLinkProps{
			Href: "https://discord.com",
			Text: "Open in Discord",
		}))},
		{Name: "external_link_no_icon", HTML: utils.Render(t, ExternalLink(ExternalLinkProps{
			Href:     "https://example.com",
			Text:     "Example",
			ShowIcon: false,
		}))},
		{Name: "external_link_custom_class", HTML: utils.Render(t, ExternalLink(ExternalLinkProps{
			Href:  "https://docs.example.com",
			Text:  "Documentation",
			Class: "text-blue-600 hover:text-blue-500 dark:text-blue-400",
		}))},
	})
}
