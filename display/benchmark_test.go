package display

import (
	"bytes"
	"context"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func BenchmarkHotPaths(b *testing.B) {
	b.Run("Class merge", func(b *testing.B) {
		for b.Loop() {
			utils.Class("px-4 py-2", "px-6", "bg-red-500", "bg-blue-500")
		}
	})

	b.Run("Badge render", func(b *testing.B) {
		props := DefaultBadgeProps()
		props.Text = activeBadgeText
		b.ResetTimer()
		for b.Loop() {
			var buf bytes.Buffer
			_ = Badge(props).Render(context.Background(), &buf)
		}
	})

	b.Run("Card render", func(b *testing.B) {
		props := DefaultCardProps()
		props.Title = "Users"
		b.ResetTimer()
		for b.Loop() {
			var buf bytes.Buffer
			_ = Card(props).Render(context.Background(), &buf)
		}
	})

	b.Run("Table render", func(b *testing.B) {
		props := TableProps{
			Headers: []string{"Name", "Email", "Role"},
			Rows: []TableRow{
				SimpleTableRow("Alice", "alice@example.com", "Admin"),
				SimpleTableRow("Bob", "bob@example.com", "User"),
			},
		}
		b.ResetTimer()
		for b.Loop() {
			var buf bytes.Buffer
			_ = Table(props).Render(context.Background(), &buf)
		}
	})

	b.Run("Modal render", func(b *testing.B) {
		props := DefaultModalProps()
		props.ID = "test-modal"
		props.Title = "Confirm"
		b.ResetTimer()
		for b.Loop() {
			var buf bytes.Buffer
			_ = Modal(props).Render(context.Background(), &buf)
		}
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
		b.ResetTimer()
		for b.Loop() {
			var buf bytes.Buffer
			_ = Dropdown(props).Render(context.Background(), &buf)
		}
	})
}
