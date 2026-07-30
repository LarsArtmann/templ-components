package display

import (
	"math"
	"strconv"
	"strings"

	"github.com/larsartmann/templ-components/utils"
)

// SparklineProps configures a tiny inline SVG line chart for trend visualization.
type SparklineProps struct {
	utils.BaseProps

	// Values are the data points to plot. Fewer than 2 points renders nothing.
	Values []float64

	// Width is the SVG canvas width in pixels. Default: 120.
	Width int

	// Height is the SVG canvas height in pixels. Default: sparklineDefaultHeight.
	Height int

	// StrokeWidth controls the line thickness. Default: 1.5.
	StrokeWidth float64

	// Filled renders a filled area beneath the line (default: false).
	Filled bool

	// Min overrides the auto-computed minimum value (nil = auto from data).
	Min *float64

	// Max overrides the auto-computed maximum value (nil = auto from data).
	Max *float64
}

// Sparkline defaults.
const (
	sparklineDefaultWidth  = 120
	sparklineDefaultHeight = 30
	sparklineDefaultStroke = 1.5
)

// DefaultSparklineProps returns sensible defaults for a sparkline.
func DefaultSparklineProps() SparklineProps {
	return SparklineProps{ //nolint:exhaustruct // intentionally minimal defaults
		Width:       sparklineDefaultWidth,
		Height:      sparklineDefaultHeight,
		StrokeWidth: sparklineDefaultStroke,
	}
}

// sparklinePoints computes the SVG polyline points for the given values.
func sparklinePoints(values []float64, width, height int, minVal, maxVal float64) string {
	if len(values) < 2 {
		return ""
	}

	if maxVal <= minVal {
		maxVal = minVal + 1
	}

	stepX := float64(width) / float64(len(values)-1)
	rangeVal := maxVal - minVal
	points := make([]string, 0, len(values))

	for i, v := range values {
		xCoord := i * int(math.Round(stepX))

		normalized := (v - minVal) / rangeVal

		yCoord := max(height-int(math.Round(normalized*float64(height))), 0)
		yCoord = min(yCoord, height)

		points = append(points, strconv.Itoa(xCoord)+","+strconv.Itoa(yCoord))
	}

	return strings.Join(points, " ")
}

// sparklineAreaPath builds a closed SVG path for the filled area beneath the line.
func sparklineAreaPath(values []float64, width, height int, minVal, maxVal float64) string {
	if len(values) < 2 {
		return ""
	}

	if maxVal <= minVal {
		maxVal = minVal + 1
	}

	stepX := float64(width) / float64(len(values)-1)
	rangeVal := maxVal - minVal

	var b strings.Builder

	for i, v := range values {
		xCoord := i * int(math.Round(stepX))

		normalized := (v - minVal) / rangeVal

		yCoord := max(height-int(math.Round(normalized*float64(height))), 0)
		yCoord = min(yCoord, height)

		if i == 0 {
			b.WriteString("M ")
		} else {
			b.WriteString(" L ")
		}

		b.WriteString(strconv.Itoa(xCoord))
		b.WriteString(" ")
		b.WriteString(strconv.Itoa(yCoord))
	}

	b.WriteString(" L ")
	b.WriteString(strconv.Itoa(width))
	b.WriteString(" ")
	b.WriteString(strconv.Itoa(height))
	b.WriteString(" L 0 ")
	b.WriteString(strconv.Itoa(height))
	b.WriteString(" Z")

	return b.String()
}

// sparklineBounds returns the effective min/max for the given values,
// respecting any caller-provided overrides.
func sparklineBounds(props SparklineProps) (float64, float64) {
	if len(props.Values) == 0 {
		return 0, 1
	}

	minVal, maxVal := props.Values[0], props.Values[0]

	for _, v := range props.Values[1:] {
		if v < minVal {
			minVal = v
		}

		if v > maxVal {
			maxVal = v
		}
	}

	// Apply caller overrides.
	if props.Min != nil {
		minVal = *props.Min
	}

	if props.Max != nil {
		maxVal = *props.Max
	}

	if maxVal == minVal {
		maxVal = minVal + 1
	}

	return minVal, maxVal
}
