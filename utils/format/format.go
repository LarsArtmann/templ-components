package format

import (
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	// placeholderDash is the em-dash rendered for absent values (see
	// StringOrDash) and non-finite percentages (see Percent).
	placeholderDash = "—"

	bytesPerStep  = 1024 // IEC binary step between byte units
	countsPerStep = 1000 // decimal step between count units

	secondsPerMinute = 60
	secondsPerHour   = 3600
	secondsPerDay    = 86400

	// decimalThreshold is the value below which magnitudes render one
	// decimal place ("1.5") and from which they render integers ("42").
	decimalThreshold = 10
)

var (
	//nolint:gochecknoglobals // Package-level lookup tables for the unit ramps
	byteUnits = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
	//nolint:gochecknoglobals // Package-level lookup table for the count unit suffixes
	countUnits = []string{"", "k", "M", "B", "T"}
)

// Bytes formats a byte count using IEC binary units (B, KiB, MiB, …, EiB):
// one decimal below ten of a unit ("1.5 KiB"), integers above ("42 MiB"),
// plain bytes below 1024 ("512 B"), and "0 B" for zero. The sign is
// preserved for negative counts ("-1.5 KiB").
func Bytes(b int64) string {
	value := float64(b)

	sign := ""

	if value < 0 {
		sign = "-"
		value = -value
	}

	unit := 0
	for value >= bytesPerStep && unit < len(byteUnits)-1 {
		value /= bytesPerStep
		unit++
	}

	if unit == 0 {
		return sign + strconv.FormatFloat(value, 'f', 0, 64) + " " + byteUnits[0]
	}

	text := formatMagnitude(value)
	if text == "1024" && unit < len(byteUnits)-1 {
		// Rounding crossed the unit boundary (1023.995 KiB → "1024 KiB"):
		// bump to the next unit and reformat ("1.0 MiB").
		value /= bytesPerStep
		unit++
		text = strconv.FormatFloat(value, 'f', 1, 64)
	}

	return sign + text + " " + byteUnits[unit]
}

// CompactDuration renders a duration the way operators write it in env vars:
// "42s", "1m30s", "12h30m", "1h" — Go's h/m/s units without the padded
// trailing "0s" that time.Duration.String produces for round values. Zero
// renders "0s", sub-second components are dropped (truncation toward zero),
// and negative durations keep their sign ("-1m30s").
func CompactDuration(d time.Duration) string {
	seconds := int64(d / time.Second)

	sign := ""

	if seconds < 0 {
		sign = "-"
		seconds = -seconds
	}

	hours := seconds / secondsPerHour
	minutes := seconds % secondsPerHour / secondsPerMinute
	secs := seconds % secondsPerMinute

	var parts []string

	if hours > 0 {
		parts = append(parts, strconv.FormatInt(hours, 10)+"h")
	}

	if minutes > 0 {
		parts = append(parts, strconv.FormatInt(minutes, 10)+"m")
	}

	if secs > 0 || len(parts) == 0 {
		parts = append(parts, strconv.FormatInt(secs, 10)+"s")
	}

	return sign + strings.Join(parts, "")
}

// ClockDuration renders a duration for display: "M:SS" under an hour
// ("3:45"), "H:MM:SS" under a day ("1:02:03"), and "Nd Hh" from a day up
// ("2d 5h" — dashboard uptime convention). Negative durations clamp to
// "0:00": elapsed-time UIs have no meaningful negative. Sub-second
// components are dropped.
func ClockDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}

	seconds := int64(d / time.Second)

	switch {
	case seconds < secondsPerHour:
		return strconv.FormatInt(seconds/secondsPerMinute, 10) + ":" + twoDigits(seconds%secondsPerMinute)
	case seconds < secondsPerDay:
		return strconv.FormatInt(seconds/secondsPerHour, 10) + ":" +
			twoDigits(seconds%secondsPerHour/secondsPerMinute) + ":" + twoDigits(seconds%secondsPerMinute)
	default:
		return strconv.FormatInt(seconds/secondsPerDay, 10) + "d " +
			strconv.FormatInt(seconds%secondsPerDay/secondsPerHour, 10) + "h"
	}
}

// Percent formats a 0–1 ratio as a percentage with one decimal ("87.3%").
// Values above 1 render honestly ("150.0%" — capacity overcommit is real
// information) and negatives keep their sign. NaN and infinities degrade to
// the em-dash placeholder instead of rendering "NaN%" in a UI.
func Percent(ratio float64) string {
	if math.IsNaN(ratio) || math.IsInf(ratio, 0) {
		return placeholderDash
	}

	return strconv.FormatFloat(ratio*100, 'f', 1, 64) + "%"
}

// StringOrDash renders s, or the em-dash placeholder when s is empty — the
// table and detail-view convention for absent values.
func StringOrDash(s string) string {
	if s == "" {
		return placeholderDash
	}

	return s
}

// CompactCount formats a count with magnitude suffixes: plain below 1,000
// ("999"), then "k", "M", "B", "T" with one decimal below ten of a unit
// ("1.2k", "3.4M") and integers above ("42k", "1.2M"). The sign is
// preserved ("-2.5k").
func CompactCount(n int64) string {
	value := float64(n)

	sign := ""

	if value < 0 {
		sign = "-"
		value = -value
	}

	unit := 0
	for value >= countsPerStep && unit < len(countUnits)-1 {
		value /= countsPerStep
		unit++
	}

	if unit == 0 {
		return sign + strconv.FormatFloat(value, 'f', 0, 64)
	}

	text := formatMagnitude(value)
	if text == "1000" && unit < len(countUnits)-1 {
		// Rounding crossed the unit boundary (999,999 → "1000k"):
		// bump to the next unit ("1.0M").
		value /= countsPerStep
		unit++
		text = strconv.FormatFloat(value, 'f', 1, 64)
	}

	return sign + text + countUnits[unit]
}

// formatMagnitude renders a value ≥ 1 with one decimal below ten and no
// decimals from ten up.
func formatMagnitude(value float64) string {
	if value < decimalThreshold {
		return strconv.FormatFloat(value, 'f', 1, 64)
	}

	return strconv.FormatFloat(value, 'f', 0, 64)
}

// twoDigits renders 0–59 with a leading zero.
func twoDigits(n int64) string {
	if n < decimalThreshold {
		return "0" + strconv.FormatInt(n, 10)
	}

	return strconv.FormatInt(n, 10)
}
