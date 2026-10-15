package tone

// Hue is the categorical hue family: which Tailwind color palette a class
// draws from, independent of meaning. Thirteen hues are in active use across
// the library's lookup tables (verified 2026-10-15); indigo is deliberately
// absent (nothing renders it).
//
// Hue is for decisions that are genuinely categorical (chart series, kanban
// lanes, stat tiles): "the third series is violet" carries no semantics.
// When a color MEANS something, use Tone and let the tone table pick the
// hue.
type Hue string

const (
	// Chart palette hues first, in canonical series order (line/pie/segment
	// charts cycle this order — chart_geometry.go, pie_chart.go).
	HueBlue    Hue = "blue"
	HueEmerald Hue = "emerald"
	HueAmber   Hue = "amber"
	HueRose    Hue = "rose"
	HueViolet  Hue = "violet"
	HueCyan    Hue = "cyan"
	HueOrange  Hue = "orange"
	HuePink    Hue = "pink"
	// Semantic-dialect hues (the feedback/badge/kanban/stat families).
	HueGray   Hue = "gray"
	HueGreen  Hue = "green"
	HueYellow Hue = "yellow"
	HueRed    Hue = "red"
	HuePurple Hue = "purple"
)

// HueDefault is the canonical default (Gray).
const HueDefault = HueGray

// HueIsValid reports whether h is a recognized Hue.
func HueIsValid(h Hue) bool {
	switch h {
	case HueBlue, HueEmerald, HueAmber, HueRose, HueViolet, HueCyan,
		HueOrange, HuePink, HueGray, HueGreen, HueYellow, HueRed, HuePurple:
		return true
	default:
		return false
	}
}

// Hues returns every Hue, chart palette first (canonical series order),
// then the semantic-dialect hues.
func Hues() []Hue {
	return []Hue{
		HueBlue, HueEmerald, HueAmber, HueRose, HueViolet, HueCyan,
		HueOrange, HuePink, HueGray, HueGreen, HueYellow, HueRed, HuePurple,
	}
}

// ChartHues returns the eight chart-series hues in canonical cycle order.
// Line/pie/segment palettes are views of this slice (see
// display/chart_geometry.go for the pre-tone constants).
func ChartHues() []Hue {
	return []Hue{HueBlue, HueEmerald, HueAmber, HueRose, HueViolet, HueCyan, HueOrange, HuePink}
}
