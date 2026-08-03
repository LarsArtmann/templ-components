package display

import (
	"github.com/larsartmann/templ-components/utils"
)

// LineChartStyle controls how series lines are drawn.
type LineChartStyle string

const (
	// LineChartStyleLinear connects data points with straight line segments.
	LineChartStyleLinear LineChartStyle = "linear"
	// LineChartStyleSmooth connects data points with a Catmull-Rom spline curve.
	LineChartStyleSmooth LineChartStyle = "smooth"
)

// LineChartStyleIsValid reports whether v is one of the defined LineChartStyle constants.
func LineChartStyleIsValid(v LineChartStyle) bool {
	return v == LineChartStyleLinear || v == LineChartStyleSmooth
}

// ChartPadding defines the inset spacing around the plot area in an SVG chart.
// Left holds Y-axis labels, Bottom holds X-axis labels, Top/Right give
// breathing room.
type ChartPadding struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

// DefaultChartPadding returns sensible defaults for a 600×300 chart.
func DefaultChartPadding() ChartPadding {
	return ChartPadding{
		Top:    chartPaddingTop,
		Right:  chartPaddingRight,
		Bottom: chartPaddingBottom,
		Left:   chartPaddingLeft,
	}
}

// LineChartSeries is a single data series in a line chart.
type LineChartSeries struct {
	// Name is the series label shown in the legend.
	Name string

	// Values are the Y-axis data points.
	Values []float64

	// Color overrides the palette color with a Tailwind text-* class
	// (e.g. "text-emerald-600 dark:text-emerald-400"). Empty = palette by index.
	Color string

	// StrokeWidth is the line thickness. Default: 2.
	StrokeWidth float64

	// Dashed renders the line with a dashed pattern.
	Dashed bool
}

// LineChartProps configures a pure-SVG line chart with axes, gridlines,
// multi-series support, and a legend. Zero JavaScript — all rendering is
// server-side SVG. Dark-mode aware via Tailwind dark: variants on SVG elements.
type LineChartProps struct {
	utils.BaseProps

	// Series is the data to plot. Each series becomes one line.
	Series []LineChartSeries

	// XAxisLabels are category labels along the X-axis (e.g. month names).
	// If empty, no X-axis labels are rendered.
	XAxisLabels []string

	// Width is the SVG canvas width in pixels. Default: 600.
	Width int

	// Height is the SVG canvas height in pixels. Default: 300.
	Height int

	// Padding controls the inset around the plot area. Default: DefaultChartPadding().
	Padding ChartPadding

	// Min overrides the auto-computed Y-axis minimum (nil = auto from data).
	Min *float64

	// Max overrides the auto-computed Y-axis maximum (nil = auto from data).
	Max *float64

	// ShowGrid renders dashed horizontal gridlines at each Y tick. Default: true.
	ShowGrid bool

	// ShowDots renders a circle at each data point. Default: true.
	ShowDots bool

	// ShowLegend renders a color-swatch legend above the chart when there
	// are 2+ series. Default: true.
	ShowLegend bool

	// Style controls whether lines are straight (Linear) or curved (Smooth).
	// Default: LineChartStyleLinear.
	Style LineChartStyle

	// ValueFormat formats Y-axis tick labels. Default: FormatTickValue.
	ValueFormat func(float64) string

	// EmptyMessage is shown when Series is empty. Default: "No data".
	EmptyMessage string
}

// LineChart defaults.
const (
	lineChartDefaultWidth  = 600
	lineChartDefaultHeight = 300
	lineChartDefaultStroke = 2.0

	chartPaddingTop    = 20
	chartPaddingRight  = 20
	chartPaddingBottom = 30
	chartPaddingLeft   = 40

	lineChartDotRadius = 3
	lineChartTickLen   = 5
	lineChartMaxTicks  = 8

	lineChartDefaultEmptyMsg = "No data"

	lineChartLabelYOffset  = 4.0
	lineChartFontSize      = "11"
	lineChartEmptyFontSize = "14"
	lineChartXLabelOffset  = 18
	lineChartLegendCharW   = 7
	lineChartLegendGap     = 28
	lineChartLegendY       = 12
)

// Shared Tailwind color class constants for chart palettes.
const (
	chartColorBlue    = "text-blue-600 dark:text-blue-400"
	chartColorEmerald = "text-emerald-600 dark:text-emerald-400"
	chartColorAmber   = "text-amber-600 dark:text-amber-400"
	chartColorRose    = "text-rose-600 dark:text-rose-400"
	chartColorViolet  = "text-violet-600 dark:text-violet-400"
	chartColorCyan    = "text-cyan-600 dark:text-cyan-400"
	chartColorOrange  = "text-orange-600 dark:text-orange-400"
	chartColorPink    = "text-pink-600 dark:text-pink-400"
)

// lineChartPalette is the default color cycle for series without an explicit Color.
// Uses Tailwind text-* classes so SVG stroke/fill inherit via currentColor.
//
//nolint:gochecknoglobals // Package-level lookup table for chart series colors
var lineChartPalette = []string{
	chartColorBlue,
	chartColorEmerald,
	chartColorAmber,
	chartColorRose,
	chartColorViolet,
	chartColorCyan,
}

// DefaultLineChartProps returns sensible defaults for a line chart.
func DefaultLineChartProps() LineChartProps {
	return LineChartProps{ //nolint:exhaustruct // intentionally minimal defaults
		Width:        lineChartDefaultWidth,
		Height:       lineChartDefaultHeight,
		Padding:      DefaultChartPadding(),
		ShowGrid:     true,
		ShowDots:     true,
		ShowLegend:   true,
		Style:        LineChartStyleLinear,
		ValueFormat:  FormatTickValue,
		EmptyMessage: lineChartDefaultEmptyMsg,
	}
}

// lineChartBounds returns the effective min/max across all series values,
// respecting any caller-provided overrides. If all series are empty, returns
// (0, 1).
func lineChartBounds(series []LineChartSeries, minOverride, maxOverride *float64) (float64, float64) {
	var minVal, maxVal float64

	var hasData bool

	for _, s := range series {
		for _, v := range s.Values {
			if !hasData {
				minVal, maxVal = v, v
				hasData = true

				continue
			}

			if v < minVal {
				minVal = v
			}

			if v > maxVal {
				maxVal = v
			}
		}
	}

	if !hasData {
		return 0, 1
	}

	if minOverride != nil {
		minVal = *minOverride
	}

	if maxOverride != nil {
		maxVal = *maxOverride
	}

	if maxVal == minVal {
		maxVal = minVal + 1
	}

	return minVal, maxVal
}

// lineChartHasData reports whether at least one series has 2+ values.
func lineChartHasData(series []LineChartSeries) bool {
	for _, s := range series {
		if len(s.Values) >= 2 {
			return true
		}
	}

	return false
}

// lineChartTickFormat delegates to the custom format function or the default.
func lineChartTickFormat(format func(float64) string, v float64) string {
	if format != nil {
		return format(v)
	}

	return FormatTickValue(v)
}

// lineChartColor returns the Tailwind class for a series, using the palette
// when the series does not specify a custom color.
func lineChartColor(seriesIdx int, custom string) string {
	if custom != "" {
		return custom
	}

	return lineChartPalette[seriesIdx%len(lineChartPalette)]
}

// lineChartStrokeWidth returns the effective stroke width, defaulting to 2.
func lineChartStrokeWidth(w float64) float64 {
	if w > 0 {
		return w
	}

	return lineChartDefaultStroke
}

// lineChartTickLabelPositions distributes labels evenly across the plot width.
// Returns the X coordinate for label index i.
func lineChartTickLabelX(plotW, n, i int) float64 {
	if n <= 1 {
		return 0
	}

	return float64(i) * float64(plotW) / float64(n-1)
}
