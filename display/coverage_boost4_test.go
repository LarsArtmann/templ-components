package display

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

func TestCardVariantCoverage(t *testing.T) {
	t.Parallel()

	slot := templ.Raw(`<span data-testid="slot">slot</span>`)

	t.Run("title tags h1 through h6", func(t *testing.T) {
		t.Parallel()

		for _, tag := range []string{"h1", "h2", "h4", "h5", "h6"} {
			output := utils.Render(t, Card(CardProps{Title: "Tagged", TitleTag: tag}))
			utils.AssertContains(t, output, "<"+tag)
			utils.AssertContains(t, output, "</"+tag+">")
			utils.AssertContains(t, output, "Tagged")
		}
	})

	t.Run("subtitle and header action", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Card(CardProps{
			Title:        "Users",
			Subtitle:     "Manage your team",
			HeaderAction: slot,
		}))
		utils.AssertContains(t, output, "Manage your team")
		utils.AssertContains(t, output, `data-testid="slot"`)
	})

	t.Run("custom header slot replaces default header", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Card(CardProps{Title: "Ignored", Header: slot}))
		utils.AssertContains(t, output, `data-testid="slot"`)
		utils.AssertNotContains(t, output, ">Ignored<")
	})

	t.Run("body slot renders", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Card(CardProps{
			Title: "Slot body",
			Body:  slot,
		}))
		utils.AssertContains(t, output, `data-testid="slot"`)
	})

	t.Run("padding sm and lg", func(t *testing.T) {
		t.Parallel()
		utils.AssertContains(t, utils.Render(t, Card(CardProps{Padding: CardPaddingSM})), "px-3 py-3")
		utils.AssertContains(t, utils.Render(t, Card(CardProps{Padding: CardPaddingLG})), "px-6 py-6")
	})

	t.Run("container aware wrapper and variants", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Card(CardProps{
			Title:          "CA",
			Padding:        CardPaddingMD,
			ContainerAware: true,
			Footer:         slot,
		}))
		utils.AssertContains(t, output, "@container")
		utils.AssertContains(t, output, "@sm:px-6")
	})

	t.Run("base props propagate", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Card(CardProps{
			BaseProps: utils.BaseProps{
				ID:        "card-x",
				Class:     "w-96",
				AriaLabel: "Stats card",
				Attrs:     templ.Attributes{"data-ctx": "dash"},
			},
			Title: "T",
		}))
		utils.AssertContains(t, output, `id="card-x"`)
		utils.AssertContains(t, output, `aria-label="Stats card"`)
		utils.AssertContains(t, output, `data-ctx="dash"`)
	})

	t.Run("simple card with body slot", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SimpleCard(SimpleCardProps{
			BaseProps: utils.BaseProps{ID: "simple-1"},
			Padding:   CardPaddingLG,
			Body:      slot,
		}))
		utils.AssertContains(t, output, `id="simple-1"`)
		utils.AssertContains(t, output, `data-testid="slot"`)
		utils.AssertContains(t, output, "px-6 py-6")
	})
}

func TestStatCardVariantCoverage(t *testing.T) {
	t.Parallel()

	t.Run("all trend directions", func(t *testing.T) {
		t.Parallel()

		for _, trend := range []TrendDirection{TrendUp, TrendDown, TrendWarn, TrendNone} {
			output := utils.Render(t, StatCard(StatCardProps{
				Value:  "1,234",
				Label:  "Revenue",
				Change: "+12%",
				Trend:  trend,
			}))
			utils.AssertContains(t, output, "1,234")
			utils.AssertContains(t, output, "Revenue")
			utils.AssertContains(t, output, "+12%")
		}
	})

	t.Run("invalid trend falls back to none", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatCard(StatCardProps{
			Value:  "10",
			Label:  "L",
			Change: "c",
			Trend:  TrendDirection("bogus"),
		}))
		utils.AssertContains(t, output, "10")
	})

	t.Run("icon tile renders", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatCard(StatCardProps{
			Value: "99",
			Label: "Users",
			Icon:  icons.Users,
		}))
		utils.AssertContains(t, output, "<svg")
	})

	t.Run("href wraps card in anchor", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatCard(StatCardProps{
			Value: "5",
			Label: "Open tickets",
			Href:  "/tickets",
		}))
		utils.AssertContains(t, output, `href="/tickets"`)
	})

	t.Run("htmx attrs render", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatCard(StatCardProps{
			Value:    "5",
			Label:    "Live",
			HxGet:    "/api/stats",
			HxTarget: "#stat",
			HxSwap:   "innerHTML",
		}))
		utils.AssertContains(t, output, `hx-get="/api/stats"`)
		utils.AssertContains(t, output, `hx-target="#stat"`)
		utils.AssertContains(t, output, `hx-swap="innerHTML"`)
	})
}

func TestSectionHeadingVariantCoverage(t *testing.T) {
	t.Parallel()

	t.Run("all heading levels", func(t *testing.T) {
		t.Parallel()

		levels := map[HeadingLevel]string{
			HeadingLevelH1: "h1", HeadingLevelH3: "h3", HeadingLevelH4: "h4",
			HeadingLevelH5: "h5", HeadingLevelH6: "h6",
		}
		for level, tag := range levels {
			output := utils.Render(t, SectionHeading(SectionHeadingProps{Title: "Sec", Level: level}))
			utils.AssertContains(t, output, "<"+tag)
			utils.AssertContains(t, output, "Sec")
		}

		utils.AssertContains(t, utils.Render(t, SectionHeading(SectionHeadingProps{Title: "D"})), "<h2")
	})

	t.Run("alignments", func(t *testing.T) {
		t.Parallel()
		utils.AssertContains(t, utils.Render(t, SectionHeading(SectionHeadingProps{
			Title: "A", Align: TextAlignCenter,
		})), "text-center")
		utils.AssertContains(t, utils.Render(t, SectionHeading(SectionHeadingProps{
			Title: "A", Align: TextAlignRight,
		})), "text-end")
		utils.AssertContains(t, utils.Render(t, SectionHeading(SectionHeadingProps{
			Title: "A", Align: TextAlignLeft,
		})), "text-start")
	})

	t.Run("subtitle and base props", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SectionHeading(SectionHeadingProps{
			BaseProps: utils.BaseProps{
				ID:        "sec-1",
				AriaLabel: "Section",
				Attrs:     templ.Attributes{"data-k": "v"},
			},
			Title:    "T",
			SubTitle: "Sub",
		}))
		utils.AssertContains(t, output, `id="sec-1"`)
		utils.AssertContains(t, output, `aria-label="Section"`)
		utils.AssertContains(t, output, `data-k="v"`)
		utils.AssertContains(t, output, "Sub")
	})
}

func TestPieChartLabelAndLegendCoverage(t *testing.T) {
	t.Parallel()

	t.Run("external labels with both text anchors", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PieChart(PieChartProps{
			Slices: []PieChartSlice{
				{Label: "Right", Value: 30},
				{Label: "Left", Value: 40},
				{Label: "Back", Value: 30},
			},
			ShowLabels: true,
			LabelMode:  PieChartLabelExternal,
		}))
		utils.AssertContains(t, output, "Right")
		utils.AssertContains(t, output, "Left")
		utils.AssertContains(t, output, `text-anchor="end"`)
		utils.AssertContains(t, output, `text-anchor="start"`)
	})

	t.Run("labels suppressed with non-external mode", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PieChart(PieChartProps{
			Slices:     []PieChartSlice{{Label: "Hidden", Value: 50}},
			ShowLabels: true,
			LabelMode:  PieChartLabelNone,
		}))
		utils.AssertNotContains(t, output, "Hidden")
	})

	t.Run("legend skips zero values and falls back for empty labels", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PieChart(PieChartProps{
			Slices: []PieChartSlice{
				{Label: "Visible", Value: 60},
				{Label: "", Value: 30},
				{Label: "Zero", Value: 0},
			},
			ShowLegend: true,
		}))
		utils.AssertContains(t, output, "Visible")
		utils.AssertContains(t, output, "\u2014")
		utils.AssertNotContains(t, output, ">Zero<")
	})

	t.Run("custom size and inner radius", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PieChart(PieChartProps{
			Slices: []PieChartSlice{{Label: "A", Value: 1}},
			Width:  400, Height: 300,
			Donut:       true,
			InnerRadius: 0.5,
		}))
		utils.AssertContains(t, output, `viewBox="0 0 400 300"`)
	})

	t.Run("custom empty message", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PieChart(PieChartProps{EmptyMessage: "Nothing here"}))
		utils.AssertContains(t, output, "Nothing here")
	})
}

func TestChartLegendCoverage(t *testing.T) {
	t.Parallel()

	series := []LineChartSeries{
		{Name: "Revenue", Values: []float64{10, 25, 40, 35, 60, 55, 80}},
		{Name: "Expenses", Values: []float64{8, 15, 20, 30, 35, 40, 50}},
		{Name: "Profit", Values: []float64{2, 10, 20, 5, 25, 15, 30}},
	}
	labels := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

	t.Run("line chart legend", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LineChart(LineChartProps{
			Series:      series,
			XAxisLabels: labels,
			ShowLegend:  true,
		}))
		utils.AssertContains(t, output, "Revenue")
		utils.AssertContains(t, output, "Expenses")
		utils.AssertContains(t, output, "Profit")
	})

	t.Run("area chart legend", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, AreaChart(AreaChartProps{
			Series:      series,
			XAxisLabels: labels,
			ShowLegend:  true,
		}))
		utils.AssertContains(t, output, "Revenue")
		utils.AssertContains(t, output, "Profit")
	})
}

func TestBarChartVariantCoverage(t *testing.T) {
	t.Parallel()

	bars := []BarChartBar{
		{Label: "Alpha", Value: 40},
		{Label: "Beta", Value: 60, Color: "bg-emerald-600 dark:bg-emerald-500", Href: "/beta"},
		{Label: "Gamma", Value: 20},
	}

	t.Run("vertical orientation with values and format", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, BarChart(BarChartProps{
			Bars:        bars,
			Orient:      BarVertical,
			Height:      "8rem",
			ShowValues:  true,
			ValueFormat: func(v float64) string { return strings.Repeat("*", int(v/10)) + "!" },
			MinBarWidth: "min-w-2",
			Gap:         "gap-1",
		}))
		utils.AssertContains(t, output, "Alpha")
		utils.AssertContains(t, output, "****!")
		utils.AssertContains(t, output, "8rem")
	})

	t.Run("horizontal with max override and custom label width", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, BarChart(BarChartProps{
			Bars:       bars,
			Max:        100,
			ShowValues: true,
			LabelWidth: "w-40",
		}))
		utils.AssertContains(t, output, `href="/beta"`)
		utils.AssertContains(t, output, "w-40")
	})

	t.Run("custom empty message", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, BarChart(BarChartProps{EmptyMessage: "No bars"}))
		utils.AssertContains(t, output, "No bars")
	})
}

func TestHeatmapVariantCoverage(t *testing.T) {
	t.Parallel()

	rows := []HeatmapRow{
		{Label: "Mon", Cells: []HeatmapCell{{Value: 3, Href: "/zero"}, {Value: 7, Label: "Seven"}, {Value: 0}}},
		{Label: "Tue", Cells: []HeatmapCell{{Value: 5}, {Value: 2}, {Value: 9}}},
	}

	t.Run("values format max colorvar cellsize", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Heatmap(HeatmapProps{
			Rows:         rows,
			ColumnLabels: []string{"A", "B", "C"},
			Max:          10,
			ColorVar:     "--tc-accent",
			CellSize:     "h-6",
			ShowValues:   true,
			ValueFormat:  func(v float64) string { return strings.Repeat("#", int(v)) },
		}))
		utils.AssertContains(t, output, "Mon")
		utils.AssertContains(t, output, "######")
		utils.AssertContains(t, output, "--tc-accent")
		utils.AssertContains(t, output, "h-6")
		utils.AssertContains(t, output, `href="/zero"`)
	})
}

func TestCollapsibleSectionVariantCoverage(t *testing.T) {
	t.Parallel()

	t.Run("title tags", func(t *testing.T) {
		t.Parallel()

		for _, tag := range []string{"h1", "h2", "h4", "h5", "h6"} {
			output := utils.Render(t, CollapsibleSection(CollapsibleSectionProps{
				Title:    "CS",
				TitleTag: tag,
			}))
			utils.AssertContains(t, output, "<"+tag)
			utils.AssertContains(t, output, "CS")
		}
	})

	t.Run("collapsed storage key and custom icon", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CollapsibleSection(CollapsibleSectionProps{
			Title:      "Advanced",
			Collapsed:  true,
			StorageKey: "adv-section",
			Icon:       icons.Plus,
		}))
		utils.AssertContains(t, output, `data-collapsible="adv-section"`)
		utils.AssertNotContains(t, output, "<details open")
	})
}

func TestEmptyStateVariantCoverage(t *testing.T) {
	t.Parallel()

	t.Run("title tag variants", func(t *testing.T) {
		t.Parallel()

		for _, tag := range []string{"h1", "h2", "h4", "h5", "h6"} {
			output := utils.Render(t, EmptyState(EmptyStateProps{
				Title:    "Nothing",
				TitleTag: tag,
			}))
			utils.AssertContains(t, output, "<"+tag)
		}
	})

	t.Run("minimal without description or action", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EmptyState(EmptyStateProps{Title: "Bare"}))
		utils.AssertContains(t, output, "Bare")
		utils.AssertNotContains(t, output, "<a ")
		utils.AssertNotContains(t, output, "<button")
	})
}
