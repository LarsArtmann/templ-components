package utils

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func BenchmarkUtils(b *testing.B) {
	b.Run("Class merge 2 strings", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			_ = Class("px-4 py-2", "px-6")
		}
	})

	b.Run("Class merge 4 strings", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			_ = Class("px-4 py-2", "px-6", "bg-red-500", "bg-blue-500")
		}
	})

	b.Run("EnsureID with prefix", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			_ = EnsureID("tooltip", "")
		}
	})

	b.Run("Ternary", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			_ = Ternary(true, "yes", "no")
		}
	})

	b.Run("Lookup hit", func(b *testing.B) {
		m := map[string]string{"a": "alpha", "b": "beta"}

		b.ResetTimer()

		for b.Loop() {
			_ = Lookup(m, "a", "default")
		}
	})

	b.Run("Lookup miss", func(b *testing.B) {
		m := map[string]string{"a": "alpha"}

		b.ResetTimer()

		for b.Loop() {
			_ = Lookup(m, "z", "default")
		}
	})
}

// RenderToBuffer is a benchmark helper that renders a component to a buffer.
func RenderToBuffer(b *testing.B, c templ.Component) {
	b.Helper()

	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		b.Fatalf("render failed: %v", err)
	}

	_ = strings.TrimSpace(buf.String())
}
