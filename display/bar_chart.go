package display

import (
	"fmt"
	"math"

	"github.com/larsartmann/templ-components/utils"
)

// BarOrient controls bar chart orientation.
type BarOrient string

const (
	// BarHorizontal renders bars left-to-right (default).
	BarHorizontal BarOrient = "horizontal"
	// BarVertical renders bars bottom-to-top (column chart).
	BarVertical BarOrient = "vertical"
)

// BarOrientIsValid reports whether v is one of the defined BarOrient constants.
func BarOrientIsValid(v BarOrient) bool {
	return v == BarHorizontal || v == BarVertical
}

// BarChartBar represents a single bar in a bar chart.
type BarChartBar struct {
	// Label is the bar's category name (e.g. "general", "Alice").
	Label string

	// Value is the bar's magnitude.
	Value float64

	// Color overrides the default bar color with a Tailwind bg-* class
	// (e.g. "bg-emerald-600 dark:bg-emerald-500"). Empty = use default.
	Color string

	// Href makes the bar label a clickable link.
	Href string
}

// BarChartProps configures a CSS bar chart.
type BarChartProps struct {
	utils.BaseProps

	// Bars are the data points to render.
	Bars []BarChartBar

	// Orient is the bar orientation. Default: BarHorizontal.
	Orient BarOrient

	// Max overrides the auto-computed maximum (0 = auto from data).
	Max float64

	// BarColor is the default Tailwind bg-* class for bars without a per-bar
	// Color override. Default: "bg-blue-600 dark:bg-blue-500".
	BarColor string

	// LabelWidth is the label column width for horizontal charts (Tailwind
	// width class, e.g. "w-32"). Default: "w-32".
	LabelWidth string

	// ShowValues renders the numeric value next to each bar.
	ShowValues bool

	// ValueFormat formats the value for display. Default: fmt.Sprintf("%.0f", v).
	ValueFormat func(float64) string

	// EmptyMessage is shown when Bars is empty. Default: "No data".
	EmptyMessage string
}

const (
	defaultBarColor   = "bg-blue-600 dark:bg-blue-500"
	defaultEmptyMsg   = "No data"
	defaultLabelWidth = "w-32"
	percentScale      = 100.0
	percentRound      = 10.0
	zeroPercent       = "0.0%"
)

// DefaultBarChartProps returns sensible defaults for a bar chart.
func DefaultBarChartProps() BarChartProps {
	return BarChartProps{ //nolint:exhaustruct // intentionally minimal defaults
		Orient:       BarHorizontal,
		BarColor:     defaultBarColor,
		LabelWidth:   defaultLabelWidth,
		ShowValues:   true,
		EmptyMessage: defaultEmptyMsg,
		ValueFormat:  func(v float64) string { return fmt.Sprintf("%.0f", v) },
	}
}

// barChartMax returns the effective max, auto-computing from data when 0.
func barChartMax(bars []BarChartBar, override float64) float64 {
	return chartMaxWithOverride(override, func() float64 {
		var maxVal float64

		for _, b := range bars {
			if b.Value > maxVal {
				maxVal = b.Value
			}
		}

		return maxVal
	})
}

// barPercentWidth returns the CSS width percentage for a bar value.
func barPercentWidth(value, maxVal float64) string {
	if maxVal <= 0 {
		return zeroPercent
	}

	pct := (value / maxVal) * percentScale

	if pct < 0 {
		pct = 0
	}

	if pct > percentScale {
		pct = percentScale
	}

	return fmt.Sprintf("%.1f%%", math.Round(pct*percentRound)/percentRound)
}
