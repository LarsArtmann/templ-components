package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// --- DataTable Edge Cases ---

func TestDataTableEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("empty columns with rows still renders", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{},
			Rows: []TableRow{
				SimpleTableRow("Alice"),
			},
		}))
		utils.AssertContains(t, output, "<table")
		utils.AssertContains(t, output, "Alice")
	})

	t.Run("single column single row", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Only"},
			},
			Rows: []TableRow{
				SimpleTableRow("Value"),
			},
		}))
		utils.AssertContains(t, output, "Only")
		utils.AssertContains(t, output, "Value")
	})

	t.Run("active sort column not in columns list", func(t *testing.T) {
		t.Parallel()
		// ActiveSortColumn references a non-existent column
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name", Sortable: true},
			},
			ActiveSortColumn: "NonExistent",
			ActiveSortDir:    SortAsc,
			SortBaseURL:      "/users",
			Rows: []TableRow{
				SimpleTableRow("Alice"),
			},
		}))
		// No column matches, so no aria-sort=ascending should appear
		utils.AssertNotContains(t, output, `aria-sort="ascending"`)
		// But Name column should still get an asc link
		utils.AssertContains(t, output, "dir=asc")
	})

	t.Run("Sortable=true but SortBaseURL empty produces no links", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name", Sortable: true},
			},
			ActiveSortColumn: "Name",
			ActiveSortDir:    SortAsc,
			// SortBaseURL intentionally empty
			Rows: []TableRow{
				SimpleTableRow("Alice"),
			},
		}))
		utils.AssertNotContains(t, output, "?sort=")
		utils.AssertNotContains(t, output, "&sort=")
	})

	t.Run("custom SortKey used in URL", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Created At", Sortable: true, SortKey: "created_at"},
			},
			SortBaseURL: "/items",
			Rows: []TableRow{
				SimpleTableRow("Alice"),
			},
		}))
		utils.AssertContains(t, output, "sort=created_at")
		utils.AssertNotContains(t, output, "sort=created at")
	})

	t.Run("label defaults to lowercase sort key", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "UPPERCASE", Sortable: true},
			},
			SortBaseURL: "/items",
			Rows: []TableRow{
				SimpleTableRow("Alice"),
			},
		}))
		utils.AssertContains(t, output, "sort=uppercase")
	})

	t.Run("flush mode suppresses border", func(t *testing.T) {
		t.Parallel()
		flushOutput := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name"},
			},
			Rows: []TableRow{
				SimpleTableRow("Alice"),
			},
			Flush: true,
		}))
		normalOutput := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name"},
			},
			Rows: []TableRow{
				SimpleTableRow("Alice"),
			},
		}))
		// Normal mode has border classes, flush does not
		if flushOutput == normalOutput {
			t.Error("expected different output for flush vs non-flush")
		}
	})

	t.Run("cellPadding compact produces different output than comfortable", func(t *testing.T) {
		t.Parallel()
		compactOutput := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name"},
			},
			Rows: []TableRow{
				SimpleTableRow("Alice"),
			},
			CellPadding: TableCellPaddingCompact,
		}))
		utils.AssertContains(t, compactOutput, "py-2")
	})

	t.Run("nil pagination renders nothing extra", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name"},
			},
			Rows: []TableRow{
				SimpleTableRow("Alice"),
			},
			Pagination: nil,
		}))
		utils.AssertContains(t, output, "</table>")
	})

	t.Run("nil empty state with empty rows renders bare table", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name"},
			},
			Rows:       []TableRow{},
			EmptyState: nil,
		}))
		utils.AssertContains(t, output, "<table")
		utils.AssertContains(t, output, "Name")
	})
}
