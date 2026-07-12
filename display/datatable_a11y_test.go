package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// --- DataTable Accessibility ---

func TestDataTableA11y(t *testing.T) {
	t.Parallel()

	t.Run("table has implicit role from <table> element", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name"},
			},
			Rows: []TableRow{
				SimpleTableRow("Alice"),
			},
		}))
		utils.AssertContains(t, output, "<table")
	})

	t.Run("sortable active column has aria-sort", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name", Sortable: true},
			},
			ActiveSortColumn: "Name",
			ActiveSortDir:    SortAsc,
			SortBaseURL:      "/users",
			Rows: []TableRow{
				SimpleTableRow("Alice"),
			},
		}))
		utils.AssertContains(t, output, `aria-sort="ascending"`)
	})

	t.Run("sortable descending column has aria-sort", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name", Sortable: true},
			},
			ActiveSortColumn: "Name",
			ActiveSortDir:    SortDesc,
			SortBaseURL:      "/users",
			Rows: []TableRow{
				SimpleTableRow("Alice"),
			},
		}))
		utils.AssertContains(t, output, `aria-sort="descending"`)
	})

	t.Run("sortable inactive column has aria-sort=none", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name", Sortable: true},
				{Label: "Email", Sortable: true},
				{Label: "Role"},
			},
			ActiveSortColumn: "Name",
			ActiveSortDir:    SortAsc,
			SortBaseURL:      "/users",
			Rows: []TableRow{
				SimpleTableRow("Alice", "alice@example.com", "Admin"),
			},
		}))
		// Email is sortable but not active → aria-sort="none"
		utils.AssertContains(t, output, `aria-sort="none"`)
	})

	t.Run("sort header link has accessible text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name", Sortable: true},
			},
			SortBaseURL: "/users",
			Rows: []TableRow{
				SimpleTableRow("Alice"),
			},
		}))
		utils.AssertContains(t, output, ">Name<")
	})

	t.Run("caption renders as visually hidden", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name"},
			},
			Rows: []TableRow{
				SimpleTableRow("Alice"),
			},
			Caption: "User directory",
		}))
		utils.AssertContains(t, output, "User directory")
		utils.AssertContains(t, output, "<caption")
	})

	t.Run("propagates aria-label from BaseProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			BaseProps: utils.BaseProps{AriaLabel: "Users table"},
			Columns: []DataTableColumn{
				{Label: "Name"},
			},
			Rows: []TableRow{
				SimpleTableRow("Alice"),
			},
		}))
		utils.AssertContains(t, output, `aria-label="Users table"`)
	})

	t.Run("propagates custom class from BaseProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			BaseProps: utils.BaseProps{Class: "my-table"},
			Columns: []DataTableColumn{
				{Label: "Name"},
			},
			Rows: []TableRow{
				SimpleTableRow("Alice"),
			},
		}))
		utils.AssertContains(t, output, "my-table")
	})
}
