package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// TestRTLRendering verifies that components render correctly with dir="rtl"
// — logical CSS properties (ms-/me-/start-0/end-0) resolve identically
// regardless of direction, so we assert the logical properties are present.
func TestRTLRendering(t *testing.T) {
	t.Parallel()

	t.Run("drawer uses logical positioning via data-side", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Drawer(DrawerProps{
			BaseProps: utils.BaseProps{ID: "rtl-drawer"},
			Open:      true,
			Side:      DrawerLeft,
		}))
		utils.AssertContains(t, output, `data-side="left"`)
	})

	t.Run("statcard renders without physical margins", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatCard(StatCardProps{
			Value:  "42",
			Label:  "Active",
			Change: "+5%",
			Trend:  TrendUp,
		}))
		// ml-2 should have been converted to ms-2 by logical property migration
		utils.AssertContains(t, output, "ms-2")
		utils.AssertNotContains(t, output, "ml-2")
	})

	t.Run("table uses logical text alignment, not physical", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers:     []string{"Name"},
			Rows:        []TableRow{SimpleTableRow("Alice")},
			CellPadding: TableCellPaddingCompact,
		}))
		utils.AssertContains(t, output, "text-start")
		utils.AssertNotContains(t, output, "text-left")
	})
}
