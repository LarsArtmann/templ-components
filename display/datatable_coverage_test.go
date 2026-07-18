package display

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// --- DataTable Behavior (BDD-style) ---

func TestDataTableUserCanViewData(t *testing.T) {
	t.Parallel()

	t.Run("user sees column headers and row data", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name"},
				{Label: "Email"},
			},
			Rows: []TableRow{
				SimpleTableRow("Alice", "alice@example.com"),
				SimpleTableRow("Bob", "bob@example.com"),
			},
		}))
		utils.AssertContains(t, output, "Name")
		utils.AssertContains(t, output, "Email")
		utils.AssertContains(t, output, "Alice")
		utils.AssertContains(t, output, "bob@example.com")
	})
}

func TestDataTableUserCanSortColumns(t *testing.T) {
	t.Parallel()

	t.Run("user clicks active ascending column to sort descending", func(t *testing.T) {
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
		// The toggle link should switch to desc
		utils.AssertContains(t, output, "dir=desc")
	})

	t.Run("user clicks active descending column to sort ascending", func(t *testing.T) {
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
		utils.AssertContains(t, output, "dir=asc")
	})

	t.Run("user clicks inactive sortable column to start ascending", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name", Sortable: true},
				{Label: "Email", Sortable: true, SortKey: "email_address"},
			},
			ActiveSortColumn: "Name",
			ActiveSortDir:    SortAsc,
			SortBaseURL:      "/users",
			Rows: []TableRow{
				SimpleTableRow("Alice", "alice@example.com"),
			},
		}))
		// Email column should link to asc
		utils.AssertContains(t, output, "sort=email_address")
		// Both asc links should be present: email's asc + name's toggle to desc
		ascCount := strings.Count(output, "dir=asc")
		if ascCount < 1 {
			t.Error("expected at least one dir=asc link for inactive sortable column")
		}
	})
}

func TestDataTableUserSeesEmptyState(t *testing.T) {
	t.Parallel()

	t.Run("user sees custom empty state when no rows", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name"},
			},
			Rows:       []TableRow{},
			EmptyState: SimpleEmptyState("No users found"),
		}))
		utils.AssertContains(t, output, "No users found")
		utils.AssertNotContains(t, output, "<table")
	})
}

func TestDataTableUserCanPaginate(t *testing.T) {
	t.Parallel()

	t.Run("pagination renders below table when rows exist", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name"},
			},
			Rows: []TableRow{
				SimpleTableRow("Alice"),
			},
			Pagination: SimpleEmptyState("Page 1 of 3"),
		}))
		utils.AssertContains(t, output, "Page 1 of 3")
	})
}

// --- DataTable Snapshot ---

func TestDataTableSnapshot(t *testing.T) {
	t.Parallel()

	t.Run("sortable table with custom params and flush", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name", Sortable: true},
				{Label: "Status", Sortable: true, SortKey: "is_active"},
				{Label: "Actions"},
			},
			ActiveSortColumn: "Name",
			ActiveSortDir:    SortDesc,
			SortBaseURL:      "/admin/users",
			SortParam:        "orderBy",
			DirParam:         "order",
			Rows: []TableRow{
				SimpleTableRow("Alice", "Active", "Edit"),
				SimpleTableRow("Bob", "Inactive", "Edit"),
			},
			Flush:       true,
			Hover:       true,
			CellPadding: TableCellPaddingCompact,
		}))
		utils.AssertContainsAll(t, output, "orderBy=name", "order=asc", "orderBy=is_active", "Alice", "Bob")
		utils.AssertContains(t, output, `aria-sort="descending"`)
	})

	t.Run("table with pagination and empty state both set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, DataTable(DataTableProps{
			Columns: []DataTableColumn{
				{Label: "Name"},
			},
			Rows:       []TableRow{},
			EmptyState: SimpleEmptyState("Nothing here"),
			Pagination: SimpleEmptyState("Page 2"),
		}))
		// Empty state renders, pagination does not (no rows)
		utils.AssertContains(t, output, "Nothing here")
		utils.AssertNotContains(t, output, "Page 2")
	})
}

// --- DataTable Coverage (private helpers) ---

func TestDataTableSortURL(t *testing.T) {
	t.Parallel()
	// Normal URL
	got := dataTableSortURL("/users", "sort", "name", SortAsc, "dir")
	utils.AssertContains(t, got, "/users?")
	utils.AssertContains(t, got, "sort=name")
	utils.AssertContains(t, got, "dir=asc")

	// URL with existing query params
	got2 := dataTableSortURL("/users?page=1", "sort", "name", SortDesc, "dir")
	utils.AssertContains(t, got2, "/users?page=1&")
	utils.AssertContains(t, got2, "dir=desc")

	// Empty baseURL returns empty
	got3 := dataTableSortURL("", "sort", "name", SortAsc, "dir")
	if got3 != "" {
		t.Errorf("dataTableSortURL with empty baseURL = %q, want empty", got3)
	}
}

func TestDataTableTypedHeadersDefaults(t *testing.T) {
	t.Parallel()
	// When SortParam and DirParam are empty, defaults are applied
	props := DataTableProps{
		Columns: []DataTableColumn{
			{Label: "Name", Sortable: true},
		},
		ActiveSortColumn: "Name",
		ActiveSortDir:    SortAsc,
		SortBaseURL:      "/items",
		// SortParam and DirParam intentionally empty
		Rows: []TableRow{
			SimpleTableRow("Alice"),
		},
	}

	headers := dataTableTypedHeaders(props)
	if len(headers) != 1 {
		t.Fatalf("expected 1 header, got %d", len(headers))
	}

	if !headers[0].Sortable {
		t.Error("expected header to be sortable")
	}

	if headers[0].SortDirection != SortAsc {
		t.Error("expected SortDirection=SortAsc")
	}

	utils.AssertContains(t, headers[0].Href, "sort=name")
	utils.AssertContains(t, headers[0].Href, "dir=desc")
}
