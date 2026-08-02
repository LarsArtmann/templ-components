package datastar

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
)

func BenchmarkDatastar(b *testing.B) {
	b.Run("SDKScript render", func(b *testing.B) {
		props := DefaultSDKScriptProps()

		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = SDKScript(props).Render(context.Background(), &buf)
		}
	})

	b.Run("SDKScript self-hosted render", func(b *testing.B) {
		props := SDKScriptProps{Src: "/assets/datastar.js"}

		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = SDKScript(props).Render(context.Background(), &buf)
		}
	})

	b.Run("LiveRegion render", func(b *testing.B) {
		props := LiveRegionProps{
			URL:       "/stream/metrics",
			AutoStart: true,
		}

		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = LiveRegion(props).Render(context.Background(), &buf)
		}
	})

	b.Run("Indicator render", func(b *testing.B) {
		props := IndicatorProps{Signal: "loading"}

		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = Indicator(props).Render(context.Background(), &buf)
		}
	})

	b.Run("Indicator with custom spinner render", func(b *testing.B) {
		props := IndicatorProps{
			Signal:  "saving",
			Spinner: templ.NopComponent,
		}

		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = Indicator(props).Render(context.Background(), &buf)
		}
	})
}
