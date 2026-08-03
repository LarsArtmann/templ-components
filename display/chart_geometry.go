package display

import (
	"fmt"
	"math"
	"strconv"
	"strings"
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
	n := len(values)
	if n == 0 {
		return nil
	}

	if maxVal <= minVal {
		maxVal = minVal + scaleRangePad
	}

	points := make([]Point, n) //nolint:makezero // pre-allocated with exact size, filled by index
	rangeVal := maxVal - minVal

	for i, v := range values {
		var x float64
		if n > 1 {
			x = float64(i) * float64(width) / float64(n-1)
		}

		normalized := (v - minVal) / rangeVal
		if normalized < normalizedMin {
			normalized = normalizedMin
		} else if normalized > normalizedMax {
			normalized = normalizedMax
		}

		y := float64(height) - normalized*float64(height)

		points[i] = Point{X: x, Y: y}
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
	n := len(points)
	if n < smoothPathMinPoints {
		return BuildPolylinePath(points)
	}

	var b strings.Builder

	fmt.Fprintf(&b, "M %s %s", formatCoord(points[0].X), formatCoord(points[0].Y))

	tension := 1.0 / splineTension

	for i := range n - 1 {
		p0 := points[max(i-1, 0)]
		p1 := points[i]
		p2 := points[i+1]
		p3 := points[min(i+2, n-1)]

		cp1x := p1.X + (p2.X-p0.X)*tension
		cp1y := p1.Y + (p2.Y-p0.Y)*tension
		cp2x := p2.X - (p3.X-p1.X)*tension
		cp2y := p2.Y - (p3.Y-p1.Y)*tension

		fmt.Fprintf(&b, " C %s %s %s %s %s %s",
			formatCoord(cp1x), formatCoord(cp1y),
			formatCoord(cp2x), formatCoord(cp2y),
			formatCoord(p2.X), formatCoord(p2.Y))
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
