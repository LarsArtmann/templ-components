package icons

import (
	"bytes"
	"context"
	"testing"
)

func BenchmarkIcons(b *testing.B) {
	b.Run("Icon render", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = Icon(Home, "h-5 w-5").Render(context.Background(), &buf)
		}
	})

	b.Run("IconWithStrokeWidth render", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = IconWithStrokeWidth(Home, "h-5 w-5", 2.0).Render(context.Background(), &buf)
		}
	})

	b.Run("IconPathData lookup", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			_ = IconPathData(Home)
		}
	})

	b.Run("IconPathJS lookup", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			_ = IconPathJS(Home)
		}
	})
}
