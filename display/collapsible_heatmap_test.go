package display

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollapsibleSectionDefaults(t *testing.T) {
	t.Parallel()

	p := resolveCollapsibleDefaults(CollapsibleSectionProps{})
	assert.Equal(t, "h3", p.TitleTag)
	assert.False(t, p.Collapsed, "default should be expanded")
}

func TestCollapsibleSectionRespectsOverrides(t *testing.T) {
	t.Parallel()

	p := resolveCollapsibleDefaults(CollapsibleSectionProps{
		TitleTag:  "h2",
		Collapsed: true,
	})
	assert.Equal(t, "h2", p.TitleTag)
	assert.True(t, p.Collapsed)
}

func TestIsValidHeadingTag(t *testing.T) {
	t.Parallel()

	for _, tag := range []string{"h1", "h2", "h3", "h4", "h5", "h6"} {
		assert.True(t, isValidHeadingTag(tag), "%s should be valid", tag)
	}

	for _, tag := range []string{"", "h0", "h7", "div", "span"} {
		assert.False(t, isValidHeadingTag(tag), "%q should be invalid", tag)
	}
}

func TestHeatmapMax(t *testing.T) {
	t.Parallel()

	rows := []HeatmapRow{
		{Cells: []HeatmapCell{{Value: 5}, {Value: 12}}},
		{Cells: []HeatmapCell{{Value: 3}, {Value: 20}}},
	}

	assert.InDelta(t, 20.0, heatmapMax(rows, 0), 0.001)
	assert.InDelta(t, 100.0, heatmapMax(rows, 100), 0.001)
	assert.InDelta(t, 1.0, heatmapMax(nil, 0), 0.001)
}

func TestHeatmapOpacity(t *testing.T) {
	t.Parallel()

	assert.InDelta(t, 0.0, heatmapOpacity(0, 10), 0.001)
	assert.InDelta(t, 0.0, heatmapOpacity(5, 0), 0.001)
	assert.InDelta(t, 1.0, heatmapOpacity(10, 10), 0.001)
	assert.InDelta(t, 0.5, heatmapOpacity(5, 10), 0.001)
	assert.InDelta(t, 0.05, heatmapOpacity(1, 100), 0.001, "should clamp to minimum")
}

func TestHeatmapPeakLocation(t *testing.T) {
	t.Parallel()

	rows := []HeatmapRow{
		{Cells: []HeatmapCell{{Value: 5}, {Value: 12}}},
		{Cells: []HeatmapCell{{Value: 3}, {Value: 20}}},
	}

	ri, ci := heatmapPeakLocation(rows)
	assert.Equal(t, 1, ri)
	assert.Equal(t, 1, ci)

	ri2, ci2 := heatmapPeakLocation(nil)
	assert.Equal(t, -1, ri2)
	assert.Equal(t, -1, ci2)
}

func TestHeatmapCellBackground(t *testing.T) {
	t.Parallel()

	bg := heatmapCellBackground(5, 10, "--ds-brand")
	assert.Contains(t, bg, "rgba(var(--ds-brand-rgb)")
	assert.Contains(t, bg, "0.50")
}
