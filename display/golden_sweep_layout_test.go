package display

import (
	"testing"

	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

// Golden sweep for the layout/visualization components added in v1.6.0:
// CollapsibleSection and Heatmap.

func TestGoldenSweepCollapsibleSection(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "collapsible_section_open", HTML: utils.Render(t, CollapsibleSection(CollapsibleSectionProps{
			Title: "Advanced Settings",
		}))},
		{Name: "collapsible_section_closed", HTML: utils.Render(t, CollapsibleSection(CollapsibleSectionProps{
			Title:     "Debug Info",
			Collapsed: true,
		}))},
		{Name: "collapsible_section_h2", HTML: utils.Render(t, CollapsibleSection(CollapsibleSectionProps{
			Title:    "Section Title",
			TitleTag: "h2",
		}))},
		{Name: "collapsible_section_storage_key", HTML: utils.Render(t, CollapsibleSection(CollapsibleSectionProps{
			Title:      "Persisted Section",
			StorageKey: "persisted-section",
		}))},
		{Name: "collapsible_section_custom_class", HTML: utils.Render(t, CollapsibleSection(CollapsibleSectionProps{
			Title:     "Styled Section",
			BaseProps: utils.BaseProps{Class: "border border-gray-200"},
		}))},
	})
}

func TestGoldenSweepHeatmap(t *testing.T) {
	t.Parallel()

	rows := []HeatmapRow{
		{Label: "Mon", Cells: []HeatmapCell{
			{Value: 5, Label: "Mon 00:00 — 5 msgs"},
			{Value: 12, Label: "Mon 06:00 — 12 msgs"},
			{Value: 0},
			{Value: 8, Label: "Mon 18:00 — 8 msgs"},
		}},
		{Label: "Tue", Cells: []HeatmapCell{
			{Value: 3, Label: "Tue 00:00 — 3 msgs"},
			{Value: 20, Label: "Tue 06:00 — 20 msgs"},
			{Value: 7, Label: "Tue 12:00 — 7 msgs"},
			{Value: 2, Label: "Tue 18:00 — 2 msgs"},
		}},
		{Label: "Wed", Cells: []HeatmapCell{
			{Value: 0},
			{Value: 0},
			{Value: 15, Label: "Wed 12:00 — 15 msgs"},
			{Value: 0},
		}},
	}

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "heatmap_basic", HTML: utils.Render(t, Heatmap(HeatmapProps{
			Rows:         rows,
			ColumnLabels: []string{"00:00", "06:00", "12:00", "18:00"},
		}))},
		{Name: "heatmap_peak", HTML: utils.Render(t, Heatmap(HeatmapProps{
			Rows:          rows,
			ColumnLabels:  []string{"00:00", "06:00", "12:00", "18:00"},
			HighlightPeak: true,
		}))},
		{Name: "heatmap_show_values", HTML: utils.Render(t, Heatmap(HeatmapProps{
			Rows:       rows,
			ShowValues: true,
		}))},
		{Name: "heatmap_empty", HTML: utils.Render(t, Heatmap(DefaultHeatmapProps()))},
		{Name: "heatmap_max_override", HTML: utils.Render(t, Heatmap(HeatmapProps{
			Rows: rows,
			Max:  100,
		}))},
	})
}
