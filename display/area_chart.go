package display

import "github.com/larsartmann/templ-components/utils"

// AreaChartProps configures a pure-SVG area chart — a line chart with filled
// areas beneath each series. Like LineChart, it supports axes, gridlines,
// multi-series, dots, and a legend. Zero JavaScript.
//
// Each series renders both a line (top edge) and a semi-transparent fill
// from the line down to the X-axis baseline. Uses currentColor for stroke
// and fill so Tailwind text-* classes control per-series coloring.
type AreaChartProps struct {
	utils.BaseProps

	// Series is the data to plot.
	Series []LineChartSeries

	// XAxisLabels are category labels along the X-axis.
	XAxisLabels []string

	// Width is the SVG canvas width in pixels. Default: 600.
	Width int

	// Height is the SVG canvas height in pixels. Default: 300.
	Height int

	// Padding controls the inset around the plot area.
	Padding ChartPadding

	// Min overrides the auto-computed Y-axis minimum.
	Min *float64

	// Max overrides the auto-computed Y-axis maximum.
	Max *float64

	// ShowGrid renders dashed horizontal gridlines at each Y tick. Default: true.
	ShowGrid bool

	// ShowDots renders a circle at each data point. Default: false (cleaner area look).
	ShowDots bool

	// ShowLegend renders a color-swatch legend above the chart when 2+ series.
	ShowLegend bool

	// Style controls whether lines are straight (Linear) or curved (Smooth).
	Style LineChartStyle

	// FillOpacity controls the transparency of the area fill (0.0–1.0).
	// Default: 0.2.
	FillOpacity float64

	// ValueFormat formats Y-axis tick labels.
	ValueFormat func(float64) string

	// EmptyMessage is shown when Series is empty.
	EmptyMessage string
}

// AreaChart defaults.
const (
	areaChartDefaultWidth   = 600
	areaChartDefaultHeight  = 300
	areaChartDefaultFill  = 0.2
	areaChartDefaultEmpty = "No data"
)

// DefaultAreaChartProps returns sensible defaults for an area chart.
func DefaultAreaChartProps() AreaChartProps {
	return AreaChartProps{ //nolint:exhaustruct // intentionally minimal defaults
		Width:        areaChartDefaultWidth,
		Height:       areaChartDefaultHeight,
		Padding:      DefaultChartPadding(),
		ShowGrid:     true,
		ShowDots:     false,
		ShowLegend:   true,
		Style:        LineChartStyleLinear,
		FillOpacity:  areaChartDefaultFill,
		ValueFormat:  FormatTickValue,
		EmptyMessage: areaChartDefaultEmpty,
	}
}

// areaChartFillOpacityStr converts a 0–1 float to an SVG fill-opacity attribute string.
func areaChartFillOpacityStr(opacity float64) string {
	if opacity <= 0 || opacity > 1 {
		return "0.2"
	}

	return formatCoord(opacity)
}
