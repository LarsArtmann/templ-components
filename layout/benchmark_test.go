package layout

import (
	"bytes"
	"context"
	"testing"
)

func BenchmarkLayout(b *testing.B) {
	b.Run("ThemeScript render", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = ThemeScript("nonce-123").Render(context.Background(), &buf)
		}
	})

	b.Run("ThemeToggle render", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = ThemeToggle("Toggle theme", "nonce-123").Render(context.Background(), &buf)
		}
	})

	b.Run("Script render", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = Script("nonce-123", "/app.js", nil).Render(context.Background(), &buf)
		}
	})

	b.Run("Minimal render", func(b *testing.B) {
		props := DefaultMinimalProps()

		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = Minimal(props).Render(context.Background(), &buf)
		}
	})
}
