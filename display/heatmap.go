package display

import (
	"fmt"
	"math"

	"github.com/larsartmann/templ-components/utils"
)

// HeatmapCell represents a single cell in a heatmap grid.
type HeatmapCell struct {
	// Value is the cell's magnitude.
	Value float64

	// Label is the tooltip text for the cell.
	Label string

	// Href makes the cell a clickable link.
	Href string
}

// HeatmapRow represents a single row in a heatmap grid.
type HeatmapRow struct {
	// Label is the row heading (e.g. "Mon", "general").
	Label string

	// Cells are the per-column values for this row.
	Cells []HeatmapCell
}

// HeatmapProps configures a CSS-based heatmap grid.
type HeatmapProps struct {
	utils.BaseProps

	// Rows are the data rows. Each row's Cells slice should have the same
	// length as ColumnLabels.
	Rows []HeatmapRow

	// ColumnLabels are the labels for each column (e.g. hour labels).
	ColumnLabels []string

	// Max overrides the auto-computed maximum value (0 = auto from data).
	Max float64

	// ColorVar is the CSS custom property used for cell background color.
	// Default: "--ds-brand".
	ColorVar string

	// CellSize is the Tailwind height class for each cell.
	// Default: "h-5".
	CellSize string

	// ShowValues renders the numeric value inside each cell.
	ShowValues bool

	// ValueFormat formats the value for display. Default: fmt.Sprintf("%.0f", v).
	ValueFormat func(float64) string

	// EmptyMessage is shown when Rows is empty. Default: "No data".
	EmptyMessage string

	// HighlightPeak adds a ring to the cell with the highest value.
	HighlightPeak bool
}

const (
	defaultHeatmapColorVar = "--ds-brand"
	defaultHeatmapCellSize = "h-5"
	minHeatmapOpacity      = 0.05
)

// DefaultHeatmapProps returns sensible defaults for a heatmap.
func DefaultHeatmapProps() HeatmapProps {
	return HeatmapProps{ //nolint:exhaustruct // intentionally minimal defaults
		ColorVar:     defaultHeatmapColorVar,
		CellSize:     defaultHeatmapCellSize,
		ShowValues:   false,
		EmptyMessage: defaultEmptyMsg,
		ValueFormat:  func(v float64) string { return fmt.Sprintf("%.0f", v) },
	}
}

// heatmapMax returns the effective max value, auto-computing from data when 0.
func heatmapMax(rows []HeatmapRow, override float64) float64 {
	if override > 0 {
		return override
	}

	var maxVal float64

	for _, row := range rows {
		for _, cell := range row.Cells {
			if cell.Value > maxVal {
				maxVal = cell.Value
			}
		}
	}

	return chartMaxFloor(maxVal)
}

// heatmapOpacity returns an opacity value (0.05–1.0) for a cell based on
// its value relative to the maximum. Returns 0 when max is zero.
func heatmapOpacity(value, maxVal float64) float64 {
	if maxVal <= 0 || value <= 0 {
		return 0
	}

	opacity := value / maxVal
	if opacity < minHeatmapOpacity {
		opacity = minHeatmapOpacity
	}

	return opacity
}

// heatmapPeakLocation returns the row and column indices of the peak cell.
// Returns -1, -1 when there is no data.
func heatmapPeakLocation(rows []HeatmapRow) (int, int) {
	peakRow, peakCol := -1, -1

	var peakVal float64

	for ri, row := range rows {
		for ci, cell := range row.Cells {
			if cell.Value > peakVal {
				peakVal = cell.Value
				peakRow = ri
				peakCol = ci
			}
		}
	}

	return peakRow, peakCol
}

// heatmapCellBackground returns the inline style for a heatmap cell.
func heatmapCellBackground(value, maxVal float64, colorVar string) string {
	opacity := heatmapOpacity(value, maxVal)

	return fmt.Sprintf("background-color: rgba(var(%s-rgb), %s)", colorVar, formatHeatmapOpacity(opacity))
}

const opacityScaleFactor = 100.0

// formatHeatmapOpacity formats an opacity float to 2 decimal places.
func formatHeatmapOpacity(f float64) string {
	return fmt.Sprintf("%.2f", math.Round(f*opacityScaleFactor)/opacityScaleFactor)
}
