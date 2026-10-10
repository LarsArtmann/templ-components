package format_test

import (
	"math"
	"testing"
	"time"

	"github.com/larsartmann/templ-components/utils/format"
)

func TestBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		b    int64
		want string
	}{
		{"zero", 0, "0 B"},
		{"one byte", 1, "1 B"},
		{"plain bytes", 512, "512 B"},
		{"just below KiB", 1023, "1023 B"},
		{"exact KiB", 1024, "1.0 KiB"},
		{"fractional KiB", 1536, "1.5 KiB"},
		{"ten KiB integer form", 10240, "10 KiB"},
		{"integer MiB form at 42", 44_040_192, "42 MiB"}, // 42 * 1024 * 1024
		{"integer MiB", 100 * 1024 * 1024, "100 MiB"},
		{"GiB", 5 * 1024 * 1024 * 1024, "5.0 GiB"},
		{"TiB", 3 * 1024 * 1024 * 1024 * 1024, "3.0 TiB"},
		{"PiB", 2 * 1024 * 1024 * 1024 * 1024 * 1024, "2.0 PiB"},
		{"EiB", 1 << 60, "1.0 EiB"},
		{"max int64", math.MaxInt64, "8.0 EiB"},
		{"negative fractional", -1536, "-1.5 KiB"},
		{"negative plain", -512, "-512 B"},
		{"min int64 no overflow", math.MinInt64, "-8.0 EiB"},
		{"rounding boundary bumps unit", 1048575, "1.0 MiB"}, // 1023.999 KiB rounds to 1024
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := format.Bytes(tt.b); got != tt.want {
				t.Errorf("Bytes(%d) = %q, want %q", tt.b, got, tt.want)
			}
		})
	}
}

func TestCompactDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    time.Duration
		want string
	}{
		{"zero", 0, "0s"},
		{"seconds only", 42 * time.Second, "42s"},
		{"minute and seconds", 90 * time.Second, "1m30s"},
		{"exact hour drops trailing zero", time.Hour, "1h"},
		{"hours and minutes", 12*time.Hour + 30*time.Minute, "12h30m"},
		{"all units", 26*time.Hour + 3*time.Minute + 4*time.Second, "26h3m4s"},
		{"sub-second dropped", 1500 * time.Millisecond, "1s"},
		{"sub-second only rounds to zero", 500 * time.Millisecond, "0s"},
		{"negative keeps sign", -90 * time.Second, "-1m30s"},
		{"negative seconds", -42 * time.Second, "-42s"},
		{"min duration no overflow", math.MinInt64, "-2562047h47m16s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := format.CompactDuration(tt.d); got != tt.want {
				t.Errorf("CompactDuration(%v) = %q, want %q", tt.d, got, tt.want)
			}
		})
	}
}

func TestClockDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		d    time.Duration
		want string
	}{
		{"zero", 0, "0:00"},
		{"seconds", 42 * time.Second, "0:42"},
		{"minutes and seconds", 225 * time.Second, "3:45"},
		{"leading zero seconds", 62 * time.Second, "1:02"},
		{"hour with clock format", 3723 * time.Second, "1:02:03"},
		{"under a day", 23*time.Hour + 59*time.Minute + 59*time.Second, "23:59:59"},
		{"a day and an hour", 25 * time.Hour, "1d 1h"},
		{"multiple days", 53 * time.Hour, "2d 5h"},
		{"negative clamps to zero", -time.Minute, "0:00"},
		{"sub-second dropped", 90*time.Second + 500*time.Millisecond, "1:30"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := format.ClockDuration(tt.d); got != tt.want {
				t.Errorf("ClockDuration(%v) = %q, want %q", tt.d, got, tt.want)
			}
		})
	}
}

func TestPercent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		ratio float64
		want  string
	}{
		{"zero", 0, "0.0%"},
		{"typical ratio", 0.873, "87.3%"},
		{"full", 1, "100.0%"},
		{"above one hundred", 1.5, "150.0%"},
		{"negative", -0.125, "-12.5%"},
		{"rounding to one decimal", 0.33333, "33.3%"},
		{"NaN degrades to placeholder", math.NaN(), "—"},
		{"positive infinity", math.Inf(1), "—"},
		{"negative infinity", math.Inf(-1), "—"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := format.Percent(tt.ratio); got != tt.want {
				t.Errorf("Percent(%v) = %q, want %q", tt.ratio, got, tt.want)
			}
		})
	}
}

func TestStringOrDash(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty becomes dash", "", "—"},
		{"value passes through", "hello", "hello"},
		{"zero value is not empty", "0", "0"},
		{"whitespace is a value", " ", " "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := format.StringOrDash(tt.in); got != tt.want {
				t.Errorf("StringOrDash(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestCompactCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		n    int64
		want string
	}{
		{"zero", 0, "0"},
		{"plain below thousand", 999, "999"},
		{"exactly one thousand", 1000, "1.0k"},
		{"fractional k", 1200, "1.2k"},
		{"integer k", 42000, "42k"},
		{"integer form near unit top", 999_499, "999k"},
		{"fractional M", 1_200_000, "1.2M"},
		{"integer M", 42_000_000, "42M"},
		{"billions", 5_600_000_000, "5.6B"},
		{"trillions", 7_800_000_000_000, "7.8T"},
		{"max int64 capped at T", math.MaxInt64, "9223372T"},
		{"negative", -2500, "-2.5k"},
		{"negative plain", -999, "-999"},
		{"rounding boundary bumps unit", 999_999, "1.0M"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := format.CompactCount(tt.n); got != tt.want {
				t.Errorf("CompactCount(%d) = %q, want %q", tt.n, got, tt.want)
			}
		})
	}
}
