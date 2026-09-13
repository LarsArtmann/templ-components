package display

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// benchRender times one component rendering into a fresh buffer.
func benchRender(b *testing.B, component templ.Component) {
	b.Helper()

	b.ResetTimer()

	for b.Loop() {
		var buf bytes.Buffer

		_ = component.Render(context.Background(), &buf)
	}
}

func BenchmarkHotPaths(b *testing.B) {
	b.Run("Class merge", func(b *testing.B) {
		for b.Loop() {
			utils.Class("px-4 py-2", "px-6", "bg-red-500", "bg-blue-500")
		}
	})

	b.Run("Badge render", func(b *testing.B) {
		props := DefaultBadgeProps()
		props.Text = activeBadgeText

		benchRender(b, Badge(props))
	})

	b.Run("Card render", func(b *testing.B) {
		props := DefaultCardProps()
		props.Title = "Users"

		benchRender(b, Card(props))
	})

	b.Run("Table render", func(b *testing.B) {
		props := TableProps{
			Headers: []string{"Name", "Email", "Role"},
			Rows: []TableRow{
				SimpleTableRow("Alice", "alice@example.com", "Admin"),
				SimpleTableRow("Bob", "bob@example.com", "User"),
			},
		}

		benchRender(b, Table(props))
	})

	b.Run("Modal render", func(b *testing.B) {
		props := DefaultModalProps()
		props.ID = "test-modal"
		props.Title = "Confirm"

		benchRender(b, Modal(props))
	})

	b.Run("Dropdown render", func(b *testing.B) {
		props := DropdownProps{
			BaseProps: utils.BaseProps{ID: "dd"},
			Label:     dropdownLabelActions,
			Items: []DropdownItem{
				{Text: dropdownItemEdit, Href: dropdownHrefEdit},
				{Text: "Delete", Href: "/delete"},
			},
		}

		benchRender(b, Dropdown(props))
	})

	b.Run("CopyButton render", func(b *testing.B) {
		props := DefaultCopyButtonProps()
		props.Text = "pnpm add foo"

		benchRender(b, CopyButton(props))
	})

	b.Run("CountBadge render", func(b *testing.B) {
		props := CountBadgeProps{Count: 42}

		benchRender(b, CountBadge(props))
	})

	b.Run("Image render", func(b *testing.B) {
		props := ImageProps{Src: "/photo.jpg", Alt: "Photo", Width: 128, Height: 128, FallbackSrc: "/placeholder.jpg"}

		benchRender(b, Image(props))
	})

	b.Run("RelativeTime render", func(b *testing.B) {
		props := RelativeTimeProps{Time: mustTime("2025-01-15T10:30:00Z")}

		benchRender(b, RelativeTime(props))
	})

	b.Run("Sparkline render", func(b *testing.B) {
		props := DefaultSparklineProps()
		props.Values = []float64{1, 3, 2, 5, 4, 6, 3, 7, 5, 8}

		benchRender(b, Sparkline(props))
	})

	b.Run("BarChart render", func(b *testing.B) {
		props := DefaultBarChartProps()
		props.Bars = []BarChartBar{
			{Label: "general", Value: 1200},
			{Label: "random", Value: 800},
			{Label: "dev", Value: 450},
		}

		benchRender(b, BarChart(props))
	})

	b.Run("ExternalLink render", func(b *testing.B) {
		props := DefaultExternalLinkProps()
		props.Href = "https://example.com"
		props.Text = "Open"

		benchRender(b, ExternalLink(props))
	})
}

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}

	return t
}

func BenchmarkHotPaths_CollapsibleSection_render(b *testing.B) {
	props := CollapsibleSectionProps{
		Title: "Advanced Settings",
		BaseProps: utils.BaseProps{
			Class: "border border-gray-200",
		},
	}

	b.ResetTimer()

	for b.Loop() {
		var buf bytes.Buffer

		_ = CollapsibleSection(props).Render(context.Background(), &buf)
	}
}

func BenchmarkHotPaths_Heatmap_render(b *testing.B) {
	rows := make([]HeatmapRow, 7)
	for i := range rows {
		cells := make([]HeatmapCell, 24)
		for j := range cells {
			cells[j] = HeatmapCell{Value: float64(i*24 + j)}
		}

		rows[i] = HeatmapRow{Label: "Day", Cells: cells}
	}

	props := HeatmapProps{
		Rows:          rows,
		HighlightPeak: true,
	}

	b.ResetTimer()

	for b.Loop() {
		var buf bytes.Buffer

		_ = Heatmap(props).Render(context.Background(), &buf)
	}
}

func BenchmarkHotPaths_KanbanBoard_render(b *testing.B) {
	columns := make([]KanbanColumn, 3)
	for ci := range columns {
		cards := make([]KanbanCard, 5)
		for i := range cards {
			cards[i] = KanbanCard{ID: "c", Title: "Card title"}
		}

		columns[ci] = KanbanColumn{ID: "col", Title: "Column", Cards: cards}
	}

	readonly := KanbanBoardProps{Columns: columns}

	wired := DefaultKanbanBoardProps()
	wired.BaseProps = utils.BaseProps{ID: "kb-bench"}
	wired.Columns = columns
	wired.Wire = &wire.Action{URL: "/api/kanban"}

	b.Run("readonly", func(b *testing.B) {
		for b.Loop() {
			var buf bytes.Buffer

			_ = KanbanBoard(readonly).Render(context.Background(), &buf)
		}
	})

	b.Run("wired", func(b *testing.B) {
		for b.Loop() {
			var buf bytes.Buffer

			_ = KanbanBoard(wired).Render(context.Background(), &buf)
		}
	})
}
