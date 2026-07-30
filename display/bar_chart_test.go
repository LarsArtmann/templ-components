package display

import (
	"strconv"
	"testing"

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
	utils.AssertContains(t, output, "items-end")
	utils.AssertContains(t, output, "rounded-t")
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
		utils.AssertContains(t, output, "bg-blue-600")
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

	if props.BarColor != "bg-blue-600 dark:bg-blue-500" {
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
	}
	for _, tt := range tests {
		got := truncateLabel(tt.input, tt.maxLen)
		if got != tt.want {
			t.Errorf("truncateLabel(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
		}
	}
}
