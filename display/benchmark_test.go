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
}
