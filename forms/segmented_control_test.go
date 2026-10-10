package forms

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

func TestSegmentedControlRender(t *testing.T) {
	t.Parallel()

	radioItems := []SegmentedControlItem{
		{Label: "15m", Value: "15m"},
		{Label: "1h", Value: "60m"},
		{Label: "1d", Value: "1440m"},
	}

	linkItems := []SegmentedControlItem{
		{Label: "Strict", Value: "30", Href: "/settings?threshold=30"},
		{Label: "Balanced", Value: "50", Href: "/settings?threshold=50"},
		{Label: "Lenient", Value: "70", Href: "/settings?threshold=70"},
	}

	t.Run("radio mode renders radiogroup with forward order", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SegmentedControl(SegmentedControlProps{
			Name:        "duration",
			Items:       radioItems,
			ActiveValue: "15m",
		}))
		utils.AssertContains(t, output, `role="radiogroup"`)
		utils.AssertContains(t, output, `type="radio"`)
		utils.AssertContains(t, output, `name="duration"`)
		utils.AssertContains(t, output, `value="15m"`)
		utils.AssertContains(t, output, "peer-checked:bg-blue-600")
		// Forward DOM order: the first radio in the HTML is the first item.
		first := strings.Index(output, `value="15m"`)
		second := strings.Index(output, `value="60m"`)
		if first > second || first < 0 {
			t.Errorf("radios not in forward DOM order (15m at %d, 60m at %d)", first, second)
		}
	})

	t.Run("checked radio matches active value exactly once", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SegmentedControl(SegmentedControlProps{
			Name:        "duration",
			Items:       radioItems,
			ActiveValue: "60m",
		}))
		if got, want := strings.Count(output, ` checked`), 1; got != want {
			t.Errorf("checked count = %d, want %d", got, want)
		}
	})

	t.Run("link mode renders aria-current links without radios", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SegmentedControl(SegmentedControlProps{
			Items:       linkItems,
			ActiveValue: "50",
		}))
		utils.AssertContains(t, output, `role="group"`)
		utils.AssertContains(t, output, `href="/settings?threshold=50"`)
		utils.AssertContains(t, output, `aria-current="true"`)
		utils.AssertNotContains(t, output, `type="radio"`)
		utils.AssertNotContains(t, output, "<script")
	})

	t.Run("radiogroup fallback label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SegmentedControl(SegmentedControlProps{
			Name:  "duration",
			Items: radioItems,
		}))
		utils.AssertContains(t, output, `aria-label="Options"`)
	})

	t.Run("aria label override", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SegmentedControl(SegmentedControlProps{
			BaseProps:   utils.BaseProps{AriaLabel: "Threshold presets"},
			Name:        "threshold",
			Items:       radioItems,
			ActiveValue: "15m",
		}))
		utils.AssertContains(t, output, `aria-label="Threshold presets"`)
	})

	t.Run("propagates class and attrs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SegmentedControl(SegmentedControlProps{
			BaseProps: utils.BaseProps{Class: "mb-2", Attrs: templ.Attributes{"data-x": "1"}},
			Name:      "duration",
			Items:     radioItems,
		}))
		utils.AssertContains(t, output, "mb-2")
		utils.AssertContains(t, output, `data-x="1"`)
	})

	t.Run("empty items render empty group", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SegmentedControl(SegmentedControlProps{Name: "x"}))
		utils.AssertContains(t, output, `role="radiogroup"`)
	})

	t.Run("logical properties only (RTL safe)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SegmentedControl(SegmentedControlProps{
			Name:  "duration",
			Items: radioItems,
		}))
		for _, banned := range []string{"rounded-l-", "rounded-r-", " ml-", " mr-", " pl-", " pr-", "border-l-", "border-r-", " translate-x"} {
			if strings.Contains(output, banned) {
				t.Errorf("physical property %q present — use logical utilities", banned)
			}
		}
	})
}

func TestSegmentedControlIsLinkMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		items []SegmentedControlItem
		want  bool
	}{
		{"no items", nil, false},
		{"no hrefs", []SegmentedControlItem{{Label: "A", Value: "a"}}, false},
		{"one href switches", []SegmentedControlItem{{Label: "A", Value: "a", Href: "/x"}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := segmentedControlIsLinkMode(tt.items); got != tt.want {
				t.Errorf("segmentedControlIsLinkMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoldenSweepSegmentedControl(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "segmented_control_radio", HTML: utils.Render(t, SegmentedControl(SegmentedControlProps{
			Name: "duration",
			Items: []SegmentedControlItem{
				{Label: "15m", Value: "15m"},
				{Label: "1h", Value: "60m"},
				{Label: "1d", Value: "1440m"},
			},
			ActiveValue: "60m",
		}))},
		{Name: "segmented_control_link", HTML: utils.Render(t, SegmentedControl(SegmentedControlProps{
			Items: []SegmentedControlItem{
				{Label: "Strict", Value: "30", Href: "/settings?threshold=30"},
				{Label: "Balanced", Value: "50", Href: "/settings?threshold=50"},
				{Label: "Lenient", Value: "70", Href: "/settings?threshold=70"},
			},
			ActiveValue: "50",
		}))},
		{Name: "segmented_control_unselected", HTML: utils.Render(t, SegmentedControl(SegmentedControlProps{
			Name: "threshold",
			Items: []SegmentedControlItem{
				{Label: "Low", Value: "low"},
				{Label: "High", Value: "high"},
			},
		}))},
	})
}
