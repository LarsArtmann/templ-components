package display

import (
	"fmt"
	"math"

	"github.com/larsartmann/templ-components/utils"
)

// PieChartLabelMode controls where slice labels are rendered.
type PieChartLabelMode string

const (
	// PieChartLabelExternal renders labels outside the pie with leader lines.
	PieChartLabelExternal PieChartLabelMode = "external"
	// PieChartLabelNone suppresses labels (use the legend instead).
	PieChartLabelNone PieChartLabelMode = "none"
)

// PieChartLabelModeIsValid reports whether v is one of the defined PieChartLabelMode constants.
func PieChartLabelModeIsValid(v PieChartLabelMode) bool {
	return v == PieChartLabelExternal || v == PieChartLabelNone
}

// PieChartSlice represents a single slice in a pie or donut chart.
type PieChartSlice struct {
	// Label is the slice's category name.
	Label string

	// Value is the slice's magnitude. Must be non-negative.
	Value float64

	// Color overrides the palette color with a Tailwind text-* class
	// (e.g. "text-emerald-600 dark:text-emerald-400"). Empty = palette by index.
	Color string
}

// PieChartProps configures a pure-SVG pie or donut chart. Zero JavaScript —
// all rendering is server-side SVG using arc paths. Dark-mode aware via
// Tailwind dark: variants on SVG elements.
type PieChartProps struct {
	utils.BaseProps

	// Slices is the data to plot.
	Slices []PieChartSlice

	// Width is the SVG canvas width in pixels. Default: 400.
	Width int

	// Height is the SVG canvas height in pixels. Default: 300.
	Height int

	// Donut renders a donut (ring) chart with an inner hole. Default: false (full pie).
	Donut bool

	// InnerRadius controls the donut hole size as a fraction of the radius
	// (0.0–1.0). Default: 0.6.
	InnerRadius float64

	// ShowLabels renders slice labels outside the pie. Default: true.
	ShowLabels bool

	// LabelMode controls label positioning. Default: PieChartLabelExternal.
	LabelMode PieChartLabelMode

	// ShowLegend renders a color-swatch legend below the chart. Default: true.
	ShowLegend bool

	// CenterLabel is text rendered in the center of a donut chart (e.g. "128GB").
	CenterLabel string

	// EmptyMessage is shown when Slices is empty. Default: "No data".
	EmptyMessage string
}

// PieChart defaults.
const (
	pieChartDefaultWidth   = 400
	pieChartDefaultHeight  = 300
	pieChartDefaultInner   = 0.6
	pieChartDefaultEmpty   = "No data"
	pieChartMinValue       = 0.0
	pieChartDegToRad       = math.Pi / 180.0
	pieChartFullCircle     = 360.0
	pieChartHalfCircle     = 180.0
	pieChartQuarterCircle  = 90.0
	pieChartFullCircleEps  = 0.01
	pieChartPercentScale   = 100.0
	pieChartLabelOffset    = 1.25
	pieChartFontSize       = "11"
	pieChartCenterFontSize = "18"
	pieChartLegendYGap     = 18
	pieChartLegendCharW    = 7
	pieChartLegendSwatch   = 12
	pieChartLegendGap      = 8
	pieChartLeftPad        = 20
)

// pieChartPalette is the default color cycle for slices.
//
//nolint:gochecknoglobals // Package-level lookup table for chart slice colors
var pieChartPalette = []string{
	chartColorBlue,
	chartColorEmerald,
	chartColorAmber,
	chartColorRose,
	chartColorViolet,
	chartColorCyan,
	chartColorOrange,
	chartColorPink,
}

// DefaultPieChartProps returns sensible defaults for a pie chart.
func DefaultPieChartProps() PieChartProps {
	return PieChartProps{ //nolint:exhaustruct // intentionally minimal defaults
		Width:        pieChartDefaultWidth,
		Height:       pieChartDefaultHeight,
		InnerRadius:  pieChartDefaultInner,
		ShowLabels:   true,
		LabelMode:    PieChartLabelExternal,
		ShowLegend:   true,
		EmptyMessage: pieChartDefaultEmpty,
	}
}

// DefaultDonutChartProps returns sensible defaults for a donut chart.
func DefaultDonutChartProps() PieChartProps {
	props := DefaultPieChartProps()
	props.Donut = true

	return props
}

// pieChartTotal returns the sum of all slice values, ignoring negative values.
func pieChartTotal(slices []PieChartSlice) float64 {
	var total float64

	for _, s := range slices {
		if s.Value > pieChartMinValue {
			total += s.Value
		}
	}

	return total
}

// pieChartColor returns the Tailwind class for a slice, using the palette
// when the slice does not specify a custom color.
func pieChartColor(sliceIdx int, custom string) string {
	if custom != "" {
		return custom
	}

	return pieChartPalette[sliceIdx%len(pieChartPalette)]
}

// pieChartSliceAngle converts a value to an angle in degrees.
func pieChartSliceAngle(value, total float64) float64 {
	if total <= 0 {
		return 0
	}

	return (value / total) * pieChartFullCircle
}

// sliceAngleResult holds a computed arc for a single slice.
type sliceAngleResult struct {
	sliceIdx   int     // original index in the Slices slice
	startAngle float64 // degrees, measured clockwise from 12 o'clock
	endAngle   float64
	midAngle   float64
	value      float64
	percent    float64
}

// computeSliceAngles converts slice values to arc angles.
// Angles are in degrees, measured clockwise from the top (12 o'clock position).
func computeSliceAngles(slices []PieChartSlice) []sliceAngleResult {
	total := pieChartTotal(slices)
	if total <= 0 {
		return nil
	}

	results := make([]sliceAngleResult, 0, len(slices))
	var startAngle float64

	for idx, s := range slices {
		angle := pieChartSliceAngle(s.Value, total)
		if angle <= 0 {
			continue
		}

		endAngle := startAngle + angle
		midAngle := startAngle + angle/2

		results = append(results, sliceAngleResult{
			sliceIdx:   idx,
			startAngle: startAngle,
			endAngle:   endAngle,
			midAngle:   midAngle,
			value:      s.Value,
			percent:    (s.Value / total) * pieChartPercentScale,
		})

		startAngle = endAngle
	}

	return results
}

// polarToCartesian converts polar coordinates (angle in degrees from top,
// radius) to Cartesian coordinates centered at (cx, cy).
func polarToCartesian(cx, cy, radius, angleDeg float64) (float64, float64) {
	rad := (angleDeg - pieChartQuarterCircle) * pieChartDegToRad

	return cx + radius*math.Cos(rad), cy + radius*math.Sin(rad)
}

// computeArcPath builds an SVG path string for a pie or donut arc.
// For a full pie: innerRadius = 0. For a donut: innerRadius > 0.
// When the arc spans 360 degrees (single slice), it renders two half-circle
// arcs because SVG arc commands cannot represent a full circle in one segment.
func computeArcPath(cx, cy, radius, innerRadius, startAngle, endAngle float64) string {
	if radius <= 0 {
		return ""
	}

	sweep := endAngle - startAngle
	fullCircle := sweep >= pieChartFullCircle-pieChartFullCircleEps

	if fullCircle {
		return computeFullCirclePath(cx, cy, radius, innerRadius)
	}

	x1, y1 := polarToCartesian(cx, cy, radius, startAngle)
	x2, y2 := polarToCartesian(cx, cy, radius, endAngle)

	largeArc := 0
	if sweep > pieChartHalfCircle {
		largeArc = 1
	}

	if innerRadius <= 0 {
		return fmt.Sprintf("M %g %g A %g %g 0 %d 1 %g %g L %g %g Z",
			x1, y1, radius, radius, largeArc, x2, y2, cx, cy)
	}

	ix1, iy1 := polarToCartesian(cx, cy, innerRadius, endAngle)
	ix2, iy2 := polarToCartesian(cx, cy, innerRadius, startAngle)

	return fmt.Sprintf("M %g %g A %g %g 0 %d 1 %g %g L %g %g A %g %g 0 %d 0 %g %g Z",
		x1, y1, radius, radius, largeArc, x2, y2,
		ix1, iy1, innerRadius, innerRadius, largeArc, ix2, iy2)
}

// computeFullCirclePath builds an SVG path for a full circle (pie or donut).
// SVG arcs cannot represent a 360° arc in one segment, so we split it into
// two semicircles.
func computeFullCirclePath(cx, cy, radius, innerRadius float64) string {
	x1, y1 := polarToCartesian(cx, cy, radius, 0)
	x2, y2 := polarToCartesian(cx, cy, radius, pieChartHalfCircle)

	if innerRadius <= 0 {
		return fmt.Sprintf("M %g %g A %g %g 0 1 1 %g %g A %g %g 0 1 1 %g %g Z",
			x1, y1, radius, radius, x2, y2, radius, radius, x1, y1)
	}

	ix1, iy1 := polarToCartesian(cx, cy, innerRadius, 0)
	ix2, iy2 := polarToCartesian(cx, cy, innerRadius, pieChartHalfCircle)

	return fmt.Sprintf("M %g %g A %g %g 0 1 1 %g %g A %g %g 0 1 1 %g %g "+
		"L %g %g A %g %g 0 1 0 %g %g A %g %g 0 1 0 %g %g Z",
		x1, y1, radius, radius, x2, y2, radius, radius, x1, y1,
		ix1, iy1, innerRadius, innerRadius, ix2, iy2, innerRadius, innerRadius, ix1, iy1)
}

// pieChartHasData reports whether at least one slice has a positive value.
func pieChartHasData(slices []PieChartSlice) bool {
	return pieChartTotal(slices) > 0
}

// pieChartFormatPercent formats a percentage value with no trailing zeros.
func pieChartFormatPercent(v float64) string {
	if v == math.Trunc(v) {
		return fmt.Sprintf("%g%%", v)
	}

	return fmt.Sprintf("%.1f%%", v)
}
