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

	// Tooltip sets a per-bar title attribute (native browser tooltip).
	// Useful for dense charts where per-bar labels are hidden.
	Tooltip string

	// ValueLabel overrides the auto-formatted value display. When set,
	// this string is shown instead of ValueFormat(Value). Useful for
	// composite labels like "123 (45%)" or "1.2 GB".
	ValueLabel string
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

	// MinBarWidth sets the minimum width for vertical bars (Tailwind
	// width class, e.g. "min-w-1" for dense time-series). Default: "min-w-12".
	MinBarWidth string

	// Gap controls the spacing between bars (Tailwind gap class,
	// e.g. "gap-px" for dense charts). Default: "gap-2" (vertical),
	// "" (horizontal, uses space-y-2).
	Gap string

	// Height sets the chart container height (CSS value, e.g. "8rem").
	// Essential for vertical charts — percentage bar heights need a
	// definite parent height. No effect on horizontal.
	Height string
}

const (
	defaultBarColor    = "bg-blue-600 dark:bg-blue-500"
	defaultEmptyMsg    = "No data"
	defaultLabelWidth  = "w-32"
	defaultMinBarWidth = "min-w-12"
	defaultVerticalGap = "gap-2"
	percentScale       = 100.0
	percentRound       = 10.0
	zeroPercent        = "0.0%"
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

// barValueLabel returns the display string for a bar's value, preferring
// the per-bar ValueLabel override when set.
func barValueLabel(bar BarChartBar, fallback func(float64) string) string {
	if bar.ValueLabel != "" {
		return bar.ValueLabel
	}

	return fallback(bar.Value)
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
