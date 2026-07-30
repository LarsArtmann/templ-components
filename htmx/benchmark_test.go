package htmx

import (
	"bytes"
	"context"
	"testing"
)

func BenchmarkHTMX(b *testing.B) {
	b.Run("LoadingIndicator render", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = LoadingIndicator(nil).Render(context.Background(), &buf)
		}
	})

	b.Run("CSRFToken render", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = CSRFToken("token-value-12345").Render(context.Background(), &buf)
		}
	})

	b.Run("SwapOOB render", func(b *testing.B) {
		props := SwapOOBProps{
			Selector:  "#target",
			SwapStyle: SwapOuterHTML,
		}

		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = SwapOOB(props).Render(context.Background(), &buf)
		}
	})

	b.Run("PolledRegion render", func(b *testing.B) {
		props := PolledRegionProps{
			URL:           "/api/live-stats",
			Every:         "10s",
			Eager:         true,
			ShowTimestamp: true,
		}

		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = PolledRegion(props).Render(context.Background(), &buf)
		}
	})
}
