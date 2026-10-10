package format_test

import (
	"math"
	"testing"
	"time"

	"github.com/larsartmann/templ-components/utils/format"
)

func BenchmarkFormat(b *testing.B) {
	b.Run("Bytes plain", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			_ = format.Bytes(512)
		}
	})

	b.Run("Bytes fractional", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			_ = format.Bytes(44_040_192)
		}
	})

	b.Run("CompactDuration", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			_ = format.CompactDuration(26*time.Hour + 3*time.Minute + 4*time.Second)
		}
	})

	b.Run("ClockDuration", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			_ = format.ClockDuration(3723 * time.Second)
		}
	})

	b.Run("Percent", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			_ = format.Percent(0.873)
		}
	})

	b.Run("CompactCount", func(b *testing.B) {
		b.ResetTimer()

		for b.Loop() {
			_ = format.CompactCount(1_200_000)
		}
	})
}

func ExampleBytes() {
	_ = format.Bytes(1536)       // "1.5 KiB"
	_ = format.Bytes(44_040_192) // "42 MiB"
	_ = format.Bytes(-512)       // "-512 B"
	// Output:
}

func ExampleCompactDuration() {
	_ = format.CompactDuration(42 * time.Second)              // "42s"
	_ = format.CompactDuration(90 * time.Second)              // "1m30s"
	_ = format.CompactDuration(12*time.Hour + 30*time.Minute) // "12h30m"
	// Output:
}

func ExampleClockDuration() {
	_ = format.ClockDuration(225 * time.Second)  // "3:45"
	_ = format.ClockDuration(3723 * time.Second) // "1:02:03"
	_ = format.ClockDuration(53 * time.Hour)     // "2d 5h"
	// Output:
}

func ExamplePercent() {
	_ = format.Percent(0.873)      // "87.3%"
	_ = format.Percent(1.5)        // "150.0%"
	_ = format.Percent(math.NaN()) // "—" placeholder
	// Output:
}
