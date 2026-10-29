package display

import (
	"testing"
)

// TestToneAdapterDefaults cover the unknown-value fallbacks of every Tone/Hue
// adapter: an out-of-set enum must map to the neutral slot, never panic and
// never return an empty value (consumers pass user data through these).
func TestToneAdapterDefaults(t *testing.T) {
	t.Parallel()

	unknown := "not-a-real-value"

	if got := StatusTone(unknown).Tone(); got == "" {
		t.Error("StatusTone(unknown).Tone() returned empty")
	}

	if got := StatusTone(unknown).DotHue(); got == "" {
		t.Error("StatusTone(unknown).DotHue() returned empty")
	}

	if got := KanbanTone(unknown).Hue(); got == "" {
		t.Error("KanbanTone(unknown).Hue() returned empty")
	}

	if got := StatTone(unknown).Hue(); got == "" {
		t.Error("StatTone(unknown).Hue() returned empty")
	}

	if got := TrendDirection(unknown).Tone(); got == "" {
		t.Error("TrendDirection(unknown).Tone() returned empty")
	}

	if got := ScrollbackTone(unknown).Tone(); got == "" {
		t.Error("ScrollbackTone(unknown).Tone() returned empty")
	}

	if got := BadgeType(unknown).Tone(); got == "" {
		t.Error("BadgeType(unknown).Tone() returned empty")
	}
}

// TestDefaultValueFormatters pins the documented default formatters that no
// golden reaches (goldens always pass their own formatter).
func TestDefaultValueFormatters(t *testing.T) {
	t.Parallel()

	barFormat := DefaultBarChartProps().ValueFormat
	if barFormat == nil {
		t.Fatal("DefaultBarChartProps().ValueFormat is nil")
	}

	if got := barFormat(42.6); got != "43" {
		t.Errorf("default bar ValueFormat(42.6) = %q, want %q", got, "43")
	}

	heatmapFormat := DefaultHeatmapProps().ValueFormat
	if heatmapFormat == nil {
		t.Fatal("DefaultHeatmapProps().ValueFormat is nil")
	}

	if got := heatmapFormat(3.4); got != "3" {
		t.Errorf("default heatmap ValueFormat(3.4) = %q, want %q", got, "3")
	}
}

// TestBarValueLabelPrecedence: a per-bar ValueLabel wins; otherwise the shared
// formatter is used.
func TestBarValueLabelPrecedence(t *testing.T) {
	t.Parallel()

	fallback := func(float64) string { return "fmt" }

	if got := barValueLabel(BarChartBar{Value: 7}, fallback); got != "fmt" {
		t.Errorf("barValueLabel without ValueLabel = %q, want %q", got, "fmt")
	}

	if got := barValueLabel(BarChartBar{Value: 7, ValueLabel: "custom"}, fallback); got != "custom" {
		t.Errorf("barValueLabel with ValueLabel = %q, want %q", got, "custom")
	}
}

// TestChartLegendSkipsUnnamedSeries: a series without a name is data, not a
// legend row, and must not consume horizontal space.
func TestChartLegendSkipsUnnamedSeries(t *testing.T) {
	t.Parallel()

	items := buildChartLegend([]LineChartSeries{
		{Name: "", Values: []float64{1, 2}},
		{Name: "Named", Values: []float64{1, 2}},
	}, 12)
	if len(items) != 1 || items[0].Name != "Named" {
		t.Fatalf("legend items = %+v, want exactly the named series", items)
	}

	if items[0].X != 12 {
		t.Errorf("first legend X = %d, want startX 12", items[0].X)
	}
}

// TestNiceStepBuckets pins the human-readable tick ladder, including the 5
// bucket that the default fixtures never reach.
func TestNiceStepBuckets(t *testing.T) {
	t.Parallel()

	if got := niceStepForNormalized(4); got != niceStepFive {
		t.Errorf("niceStepForNormalized(4) = %v, want %v", got, niceStepFive)
	}

	if got := niceStepForNormalized(42); got != niceStepTen {
		t.Errorf("niceStepForNormalized(42) = %v, want %v", got, niceStepTen)
	}
}

// TestLineChartBoundsEdges: caller overrides win over the data range, and a
// flat series is padded so the renderer never divides by zero.
func TestLineChartBoundsEdges(t *testing.T) {
	t.Parallel()

	series := []LineChartSeries{{Values: []float64{3, 9}}}
	minOverride, maxOverride := 0.0, 100.0

	minVal, maxVal := lineChartBounds(series, &minOverride, &maxOverride)
	if minVal != minOverride || maxVal != maxOverride {
		t.Errorf("lineChartBounds overrides = (%v, %v), want (0, 100)", minVal, maxVal)
	}

	minVal, maxVal = lineChartBounds([]LineChartSeries{{Values: []float64{5, 5}}}, nil, nil)
	if maxVal <= minVal {
		t.Errorf("flat range = (%v, %v), want max > min", minVal, maxVal)
	}
}

// TestLineChartGeometryEdges: positive stroke widths pass through, and a
// single tick label sits at the plot origin.
func TestLineChartGeometryEdges(t *testing.T) {
	t.Parallel()

	if got := lineChartStrokeWidth(3.5); got != 3.5 {
		t.Errorf("lineChartStrokeWidth(3.5) = %v, want 3.5", got)
	}

	if got := lineChartStrokeWidth(0); got != lineChartDefaultStroke {
		t.Errorf("lineChartStrokeWidth(0) = %v, want default %v", got, lineChartDefaultStroke)
	}

	if got := lineChartTickLabelX(200, 1, 0); got != 0 {
		t.Errorf("lineChartTickLabelX(n=1) = %v, want 0", got)
	}

	if got := lineChartTickLabelX(200, 2, 1); got != 200 {
		t.Errorf("lineChartTickLabelX(n=2, i=1) = %v, want 200", got)
	}
}

// TestPieChartFormattingEdges: a zero total yields no angle, and whole
// percentages drop the decimal.
func TestPieChartFormattingEdges(t *testing.T) {
	t.Parallel()

	if got := pieChartSliceAngle(4, 0); got != 0 {
		t.Errorf("pieChartSliceAngle(total=0) = %v, want 0", got)
	}

	if got := pieChartFormatPercent(50); got != "50%" {
		t.Errorf("pieChartFormatPercent(50) = %q, want %q", got, "50%")
	}

	if got := pieChartFormatPercent(12.5); got != "12.5%" {
		t.Errorf("pieChartFormatPercent(12.5) = %q, want %q", got, "12.5%")
	}
}

// TestSparklineDegenerateInputs: fewer than two points render nothing, and a
// flat range is padded so the area path stays finite.
func TestSparklineDegenerateInputs(t *testing.T) {
	t.Parallel()

	if got := sparklineAreaPath([]float64{1}, 20, 10, 0, 1); got != "" {
		t.Errorf("sparklineAreaPath(single point) = %q, want empty", got)
	}

	span, ok := sparklineGeometry([]float64{2, 2}, 2, 2)
	if !ok {
		t.Fatal("sparklineGeometry(two points) reported ok=false")
	}

	if span <= 0 {
		t.Errorf("sparklineGeometry(flat range) span = %v, want > 0", span)
	}

	if _, ok := sparklineGeometry([]float64{1}, 0, 1); ok {
		t.Error("sparklineGeometry(single point) reported ok=true")
	}
}

// TestOverlayKindComponentName: unknown overlay kinds fall back to the generic
// JS-name segment instead of emitting an empty function name.
func TestOverlayKindComponentName(t *testing.T) {
	t.Parallel()

	if got := OverlayKind("bogus").componentName(); got != "Overlay" {
		t.Errorf("OverlayKind(bogus).componentName() = %q, want %q", got, "Overlay")
	}

	if got := OverlayModal.componentName(); got != "Modal" {
		t.Errorf("OverlayModal.componentName() = %q, want %q", got, "Modal")
	}
}
