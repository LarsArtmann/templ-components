package display

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

// Chart geometry constants — shared across all native SVG chart components.
const (
	niceStepOne           = 1.0
	niceStepTwo           = 2.0
	niceStepTwoPointFive  = 2.5
	niceStepFive          = 5.0
	niceStepTen           = 10.0
	defaultTickCount      = 5
	roundPrecisionDigits  = 10
	splineTension         = 6.0
	scaleRangePad         = 1.0
	normalizedMin         = 0.0
	normalizedMax         = 1.0
	millionThreshold      = 1_000_000.0
	tenKThreshold         = 10_000.0
	thousandDivisor       = 1_000.0
	millionDivisor        = 1_000_000.0
	decimalPlacesFallback = 1
	float64Bits           = 64
	smoothPathMinPoints   = 3
)

// Shared Tailwind color class constants for chart palettes. Used by LineChart,
// AreaChart, and PieChart series/slice coloring via currentColor inheritance.
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

// chartMaxFloor returns 1 when the computed max is zero so percentage-based
// chart math (bar widths, heatmap opacity) has a non-zero denominator.
// Used by BarChart and Heatmap when no caller-supplied max is provided.
func chartMaxFloor(maxVal float64) float64 {
	if maxVal == 0 {
		return 1
	}

	return maxVal
}

// chartMaxWithOverride returns override when positive, otherwise calls compute
// to derive the max from the data. The returned value is then floored via
// chartMaxFloor so the consumer can use it as a percentage denominator
// without a zero-division guard. Shared by BarChart and Heatmap.
func chartMaxWithOverride(override float64, compute func() float64) float64 {
	if override > 0 {
		return override
	}

	return chartMaxFloor(compute())
}

// Sanitize clamps all padding fields to non-negative values. Negative padding
// produces negative plot dimensions (width - left - right < 0), which corrupt
// SVG path math. Called by LineChart and AreaChart before rendering.
func (p ChartPadding) Sanitize() ChartPadding {
	return ChartPadding{
		Top:    max(p.Top, 0),
		Right:  max(p.Right, 0),
		Bottom: max(p.Bottom, 0),
		Left:   max(p.Left, 0),
	}
}

// ChartLegendItem holds a computed legend entry for a chart series.
type ChartLegendItem struct {
	Name  string
	Color string
	X     int
}

// ChartRenderData bundles the pre-computed values shared between LineChart
// and AreaChart rendering. Both charts compute identical setup (bounds, ticks,
// padding, legend positions) — this struct eliminates the duplicated logic.
type ChartRenderData struct {
	Width       int
	Height      int
	Padding     ChartPadding
	PlotW       int
	PlotH       int
	MinVal      float64
	MaxVal      float64
	RangeVal    float64
	Ticks       []float64
	HasData     bool
	LabelCount  int
	XAxisLabels []string
	ShowGrid    bool
	ValueFormat func(float64) string
	LegendItems []ChartLegendItem
	EmptyMsg    string
	Class       string
	AriaLabel   string
	ID          string
	Attrs       templ.Attributes
}

// computeChartRenderData calculates the shared rendering parameters for
// LineChart and AreaChart. Both charts share identical axis/gridline/legend
// layout — only the per-series path rendering differs.
func computeChartRenderData(
	width, height int,
	padding ChartPadding,
	series []LineChartSeries,
	xAxisLabels []string,
	minOverride, maxOverride *float64,
	showGrid, showLegend bool,
	valueFormat func(float64) string,
	emptyMsg, class, ariaLabel, id string,
	attrs templ.Attributes,
) ChartRenderData {
	padding = padding.Sanitize()
	if padding == (ChartPadding{}) { //nolint:exhaustruct_v5 // zero-value sentinel detects unset padding
		padding = DefaultChartPadding()
	}

	minVal, maxVal := lineChartBounds(series, minOverride, maxOverride)
	plotW := width - padding.Left - padding.Right
	plotH := height - padding.Top - padding.Bottom
	ticks := ComputeNiceTicks(minVal, maxVal, lineChartMaxTicks)
	hasData := lineChartHasData(series)
	rangeVal := maxVal - minVal

	maxSeriesLen := 0
	for _, s := range series {
		if len(s.Values) > maxSeriesLen {
			maxSeriesLen = len(s.Values)
		}
	}

	labelCount := min(len(xAxisLabels), maxSeriesLen)

	var legendItems []ChartLegendItem
	if showLegend && len(series) > 1 {
		legendItems = buildChartLegend(series, padding.Left)
	}

	return ChartRenderData{
		Width:       width,
		Height:      height,
		Padding:     padding,
		PlotW:       plotW,
		PlotH:       plotH,
		MinVal:      minVal,
		MaxVal:      maxVal,
		RangeVal:    rangeVal,
		Ticks:       ticks,
		HasData:     hasData,
		LabelCount:  labelCount,
		XAxisLabels: xAxisLabels,
		ShowGrid:    showGrid,
		ValueFormat: valueFormat,
		LegendItems: legendItems,
		EmptyMsg:    emptyMsg,
		Class:       class,
		AriaLabel:   ariaLabel,
		ID:          id,
		Attrs:       attrs,
	}
}

// buildChartLegend lays out legend items left-to-right from startX, spacing
// each item by its label width. Shared by LineChart and AreaChart.
func buildChartLegend(series []LineChartSeries, startX int) []ChartLegendItem {
	var items []ChartLegendItem

	legendX := startX

	for i, s := range series {
		if s.Name == "" {
			continue
		}

		items = append(items, ChartLegendItem{
			Name:  s.Name,
			Color: lineChartColor(i, s.Color),
			X:     legendX,
		})
		legendX += len(s.Name)*lineChartLegendCharW + lineChartLegendGap
	}

	return items
}

// Point is a coordinate pair in SVG user space.
type Point struct {
	X, Y float64
}

// ScalePoints maps data values to SVG point coordinates within a plot area of
// the given width and height. The X coordinate is distributed evenly across the
// value count; the Y coordinate is scaled from minVal..maxVal and inverted (SVG Y
// increases downward). If maxVal <= minVal the range is padded by 1 to avoid
// division by zero. The returned points are relative to the plot area origin
// (0,0); the caller offsets them by the chart padding.
func ScalePoints(values []float64, width, height int, minVal, maxVal float64) []Point {
	valueCount := len(values)
	if valueCount == 0 {
		return nil
	}

	if maxVal <= minVal {
		maxVal = minVal + scaleRangePad
	}

	points := make([]Point, valueCount) //nolint:makezero // pre-allocated with exact size, filled by index
	rangeVal := maxVal - minVal

	for i, v := range values {
		var posX float64
		if valueCount > 1 {
			posX = float64(i) * float64(width) / float64(valueCount-1)
		}

		normalized := (v - minVal) / rangeVal
		if normalized < normalizedMin {
			normalized = normalizedMin
		} else if normalized > normalizedMax {
			normalized = normalizedMax
		}

		y := float64(height) - normalized*float64(height)

		points[i] = Point{X: posX, Y: y}
	}

	return points
}

// BuildPolylinePath builds an SVG path string connecting the given points with
// straight line segments: "M x,y L x,y ...". Returns "" for an empty slice.
func BuildPolylinePath(points []Point) string {
	if len(points) == 0 {
		return ""
	}

	var b strings.Builder

	for i, p := range points {
		if i == 0 {
			fmt.Fprintf(&b, "M %s %s", formatCoord(p.X), formatCoord(p.Y))
		} else {
			fmt.Fprintf(&b, " L %s %s", formatCoord(p.X), formatCoord(p.Y))
		}
	}

	return b.String()
}

// BuildSmoothPath builds an SVG path string using cubic Bezier curves derived
// from a Catmull-Rom spline through the points. This produces visually smooth
// lines. Falls back to a straight polyline for fewer than 3 points.
func BuildSmoothPath(points []Point) string {
	pointCount := len(points)
	if pointCount < smoothPathMinPoints {
		return BuildPolylinePath(points)
	}

	var b strings.Builder

	fmt.Fprintf(&b, "M %s %s", formatCoord(points[0].X), formatCoord(points[0].Y))

	tension := 1.0 / splineTension

	for i := range pointCount - 1 {
		prev := points[max(i-1, 0)]
		curr := points[i]
		next := points[i+1]
		afterNext := points[min(i+2, pointCount-1)]

		cp1x := curr.X + (next.X-prev.X)*tension
		cp1y := curr.Y + (next.Y-prev.Y)*tension
		cp2x := next.X - (afterNext.X-curr.X)*tension
		cp2y := next.Y - (afterNext.Y-curr.Y)*tension

		fmt.Fprintf(&b, " C %s %s %s %s %s %s",
			formatCoord(cp1x), formatCoord(cp1y),
			formatCoord(cp2x), formatCoord(cp2y),
			formatCoord(next.X), formatCoord(next.Y))
	}

	return b.String()
}

// BuildAreaPath builds a closed SVG path for a filled area: the polyline through
// the points, then down to the baseline and back to the start. The baseline is
// the bottom of the plot area (y = height).
func BuildAreaPath(points []Point, height int) string {
	if len(points) == 0 {
		return ""
	}

	path := BuildPolylinePath(points)
	baseline := float64(height)

	last := points[len(points)-1]
	first := points[0]

	return fmt.Sprintf("%s L %s %s L %s %s Z",
		path,
		formatCoord(last.X), formatCoord(baseline),
		formatCoord(first.X), formatCoord(baseline))
}

// ComputeNiceTicks produces human-readable axis tick values spanning [minVal, maxVal].
// The tick count is approximate — the actual count may differ to produce round
// numbers (e.g., 0, 25, 50, 75, 100 instead of 0, 23, 46, 69, 92).
func ComputeNiceTicks(minVal, maxVal float64, count int) []float64 {
	if count <= 0 {
		count = defaultTickCount
	}

	if maxVal <= minVal {
		maxVal = minVal + scaleRangePad
	}

	dataRange := maxVal - minVal
	rawStep := dataRange / float64(count)
	magnitude := math.Pow(niceStepTen, math.Floor(math.Log10(rawStep)))
	normalized := rawStep / magnitude

	niceStep := niceStepForNormalized(normalized)
	niceStep *= magnitude

	niceMin := math.Floor(minVal/niceStep) * niceStep
	niceMax := math.Ceil(maxVal/niceStep) * niceStep

	var ticks []float64

	// Iterate with a small epsilon to avoid floating-point drift on the last tick.
	for v := niceMin; v <= niceMax+niceStep/2; v += niceStep {
		ticks = append(ticks, roundToPrecision(v, roundPrecisionDigits))
	}

	return ticks
}

// niceStepForNormalized maps a normalized raw step to the nearest "nice" step
// value (1, 2, 2.5, 5, or 10) to produce human-readable tick intervals.
func niceStepForNormalized(normalized float64) float64 {
	switch {
	case normalized <= niceStepOne:
		return niceStepOne
	case normalized <= niceStepTwo:
		return niceStepTwo
	case normalized <= niceStepTwoPointFive:
		return niceStepTwoPointFive
	case normalized <= niceStepFive:
		return niceStepFive
	default:
		return niceStepTen
	}
}

// FormatTickValue formats a numeric value for axis tick labels. Whole numbers
// are shown without decimals; thousands and millions get K/M suffixes; small
// fractions use up to 1 decimal place.
func FormatTickValue(v float64) string {
	absVal := math.Abs(v)

	switch {
	case absVal >= millionThreshold:
		return formatSuffixed(v, millionDivisor, "M")
	case absVal >= tenKThreshold:
		return formatSuffixed(v, thousandDivisor, "K")
	default:
		if v == math.Trunc(v) {
			return strconv.FormatFloat(v, 'f', -1, float64Bits)
		}

		return strconv.FormatFloat(v, 'f', decimalPlacesFallback, float64Bits)
	}
}

// formatSuffixed divides v by divisor and appends the suffix, trimming trailing
// ".0" for clean output (e.g., "1M" not "1.0M").
func formatSuffixed(v float64, divisor float64, suffix string) string {
	scaled := v / divisor

	if scaled == math.Trunc(scaled) {
		return strconv.FormatFloat(scaled, 'f', -1, float64Bits) + suffix
	}

	return strconv.FormatFloat(scaled, 'f', decimalPlacesFallback, float64Bits) + suffix
}

// formatCoord formats a float64 SVG coordinate, trimming trailing zeros.
func formatCoord(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, float64Bits)
}

// roundToPrecision rounds v to the given number of decimal places, returning
// a value with no floating-point representation artifacts.
func roundToPrecision(v float64, precision int) float64 {
	pow := math.Pow(niceStepTen, float64(precision))

	return math.Round(v*pow) / pow
}
