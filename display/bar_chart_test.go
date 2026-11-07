package display

import (
	"strconv"
	"testing"
	"unicode/utf8"

	"github.com/larsartmann/templ-components/utils"
)

func TestBarChartHorizontal(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, BarChart(BarChartProps{
		Bars: []BarChartBar{
			{Label: "general", Value: 100},
			{Label: "random", Value: 50},
		},
	}))
	utils.AssertContains(t, output, "general")
	utils.AssertContains(t, output, "random")
	utils.AssertContains(t, output, "100")
	utils.AssertContains(t, output, "50")
	utils.AssertContains(t, output, "bg-gray-100")
	utils.AssertContains(t, output, "space-y-2")
}

func TestBarChartVertical(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, BarChart(BarChartProps{
		Orient: BarVertical,
		Bars: []BarChartBar{
			{Label: "Mon", Value: 10},
			{Label: "Tue", Value: 20},
		},
	}))
	utils.AssertContains(t, output, "height:12rem")
	utils.AssertContains(t, output, "h-full")
	utils.AssertContains(t, output, "relative w-full flex-1")
	utils.AssertContainsAll(t, output, "absolute", "inset-x-0", "rounded-t")
	utils.AssertContains(t, output, "style=\"bottom:0.0%;height:")
	utils.AssertContains(t, output, "min-w-12")
}

func TestBarChartEmpty(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, BarChart(BarChartProps{}))
	utils.AssertContains(t, output, "No data")
}

func TestBarChartCustomEmptyMessage(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, BarChart(BarChartProps{
		EmptyMessage: "No channels with messages.",
	}))
	utils.AssertContains(t, output, "No channels with messages.")
}

func TestBarChartBarColor(t *testing.T) {
	t.Parallel()
	t.Run("default color", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, BarChart(BarChartProps{
			Bars: []BarChartBar{{Label: "A", Value: 10}},
		}))
		utils.AssertContains(t, output, "bg-accent-600")
	})
	t.Run("per-bar color override", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, BarChart(BarChartProps{
			Bars: []BarChartBar{
				{Label: "A", Value: 10, Color: "bg-emerald-600 dark:bg-emerald-500"},
			},
		}))
		utils.AssertContains(t, output, "bg-emerald-600")
	})
}

func TestBarChartHref(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, BarChart(BarChartProps{
		Bars: []BarChartBar{
			{Label: "general", Value: 10, Href: "/channels/123"},
		},
	}))
	utils.AssertContains(t, output, `href="/channels/123"`)
}

func TestBarChartHideValues(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, BarChart(BarChartProps{
		Bars:       []BarChartBar{{Label: "A", Value: 100}},
		ShowValues: false,
	}))
	utils.AssertNotContains(t, output, `>100<`)
}

func TestBarChartCustomValueFormat(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, BarChart(BarChartProps{
		Bars:       []BarChartBar{{Label: "A", Value: 1234}},
		ShowValues: true,
		ValueFormat: func(v float64) string {
			return strconv.FormatFloat(v/1000, 'f', 1, 64) + "k"
		},
	}))
	utils.AssertContains(t, output, "1.2k")
}

func TestBarChartMaxOverride(t *testing.T) {
	t.Parallel()
	t.Run("override affects width", func(t *testing.T) {
		t.Parallel()

		maxVal := 200.0
		value := 100.0

		got := barPercentWidth(value, maxVal)
		if got != "50.0%" {
			t.Errorf("expected 50.0%%, got %s", got)
		}
	})
}

func TestBarChartPercentWidth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value float64
		max   float64
		want  string
	}{
		{0, 100, "0.0%"},
		{50, 100, "50.0%"},
		{100, 100, "100.0%"},
		{150, 100, "100.0%"},
		{-10, 100, "0.0%"},
		{10, 0, "0.0%"},
	}
	for _, tt := range tests {
		got := barPercentWidth(tt.value, tt.max)
		if got != tt.want {
			t.Errorf("barPercentWidth(%f, %f) = %q, want %q", tt.value, tt.max, got, tt.want)
		}
	}
}

func TestBarOrientIsValid(t *testing.T) {
	t.Parallel()

	if !BarOrientIsValid(BarHorizontal) {
		t.Error("BarHorizontal should be valid")
	}

	if !BarOrientIsValid(BarVertical) {
		t.Error("BarVertical should be valid")
	}

	if BarOrientIsValid("diagonal") {
		t.Error("diagonal should be invalid")
	}
}

func TestBarChartDefaults(t *testing.T) {
	t.Parallel()

	props := DefaultBarChartProps()
	if props.Orient != BarHorizontal {
		t.Errorf("expected Orient=horizontal")
	}

	if props.BarColor != "bg-accent-600 dark:bg-accent-500" {
		t.Errorf("expected default BarColor")
	}

	if !props.ShowValues {
		t.Error("expected ShowValues=true")
	}

	if props.EmptyMessage != "No data" {
		t.Errorf("expected EmptyMessage=No data")
	}
}

func TestTruncateLabel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"short", 10, "short"},
		{"exactly 10", 10, "exactly 10"},
		{"this is too long", 10, "this is t…"},
		// Regression: byte-slicing split multibyte runes, emitting invalid
		// UTF-8 (mojibake) for international bar labels. Truncation must be
		// rune-safe and still end with the ellipsis.
		{"aÜÜÜÜÜÜ", 6, "aÜÜÜÜ…"},
		{"Änderungsübersicht", 10, "Änderungs…"},
	}
	for _, tt := range tests {
		got := truncateLabel(tt.input, tt.maxLen)
		if got != tt.want {
			t.Errorf("truncateLabel(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
		}

		if !utf8.ValidString(got) {
			t.Errorf("truncateLabel(%q, %d) = %q is not valid UTF-8", tt.input, tt.maxLen, got)
		}
	}
}

func TestBarChartSignedPlacement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		value        float64
		maxVal       float64
		minVal       float64
		wantZero     string
		wantTop      string
		wantHeight   string
		wantNegative bool
	}{
		{"all-positive pins the zero line to the bottom", 6, 8, 0, "0.0%", "100.0%", "75.0%", false},
		{"mixed series lifts the zero line", 6, 6, -2, "25.0%", "75.0%", "75.0%", false},
		{"negative hangs below the zero line", -1, 6, -2, "25.0%", "75.0%", "12.5%", true},
		{"all-negative pins the zero line to the top", -3, 0, -3, "100.0%", "0.0%", "100.0%", true},
		{"degenerate range collapses to the baseline", 0, 0, 0, "0.0%", "100.0%", "0.0%", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			zero, top, height, negative := barChartSignedPlacement(tt.value, tt.maxVal, tt.minVal)
			if zero != tt.wantZero || top != tt.wantTop || height != tt.wantHeight || negative != tt.wantNegative {
				t.Errorf(
					"barChartSignedPlacement(%v, %v, %v) = (%q, %q, %q, %v), want (%q, %q, %q, %v)",
					tt.value,
					tt.maxVal,
					tt.minVal,
					zero,
					top,
					height,
					negative,
					tt.wantZero,
					tt.wantTop,
					tt.wantHeight,
					tt.wantNegative,
				)
			}
		})
	}
}

func TestBarChartBoundsSigned(t *testing.T) {
	t.Parallel()

	bars := []BarChartBar{
		{Label: "a", Value: 420},
		{Label: "b", Value: -180},
		{Label: "c", Value: 90},
	}

	maxVal, minVal := barChartBounds(bars, 0)
	if maxVal != 420 || minVal != -180 {
		t.Errorf("barChartBounds() = (%v, %v), want (420, -180)", maxVal, minVal)
	}

	// The Max override replaces the positive bound only; the negative bound
	// always comes from the data.
	maxVal, minVal = barChartBounds(bars, 1000)
	if maxVal != 1000 || minVal != -180 {
		t.Errorf("barChartBounds(override) = (%v, %v), want (1000, -180)", maxVal, minVal)
	}
}

func TestBarChartVerticalSigned(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, BarChart(BarChartProps{
		Orient: BarVertical,
		Bars: []BarChartBar{
			{Label: "in", Value: 420},
			{Label: "out", Value: -180},
		},
		Height:     "10rem",
		ShowValues: true,
	}))

	// The zero baseline is mandatory for signed series.
	utils.AssertContains(t, output, "h-px")
	// Negatives take the danger family and hang below the zero line.
	utils.AssertContains(t, output, "bg-rose-600 dark:bg-rose-500")
	utils.AssertContains(t, output, "rounded-b")
	utils.AssertContains(t, output, "style=\"top:70.0%;height:30.0%;\"")
	// Positives still grow upward from the zero line.
	utils.AssertContains(t, output, "style=\"bottom:30.0%;height:70.0%;\"")
}

func TestBarChartVerticalMaxBarWidth(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, BarChart(BarChartProps{
		Orient:      BarVertical,
		Bars:        []BarChartBar{{Label: "a", Value: 5}, {Label: "b", Value: 3}},
		MaxBarWidth: "max-w-24",
	}))

	utils.AssertContains(t, output, "justify-center")
	utils.AssertContains(t, output, "max-w-24")
}

func TestBarChartVerticalBaselineOptIn(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, BarChart(BarChartProps{
		Orient:       BarVertical,
		Bars:         []BarChartBar{{Label: "a", Value: 5}},
		ShowBaseline: true,
	}))

	utils.AssertContains(t, output, "h-px")
}
