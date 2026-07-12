package display

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestDefaultDataTableProps(t *testing.T) {
	t.Parallel()
	p := DefaultDataTableProps()
	if !p.Striped {
		t.Error("expected Striped=true by default")
	}
	if p.SortParam != "sort" {
		t.Error("expected SortParam='sort' by default")
	}
	if p.DirParam != "dir" {
		t.Error("expected DirParam='dir' by default")
	}
}

func TestDataTableBasicRender(t *testing.T) {
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
		Striped: true,
	}))
	utils.AssertContains(t, output, "Alice")
	utils.AssertContains(t, output, "bob@example.com")
	utils.AssertContains(t, output, "<table")
	utils.AssertContains(t, output, "<thead")
	utils.AssertContains(t, output, "Name")
	utils.AssertContains(t, output, "Email")
}

func TestDataTableSortURLsGenerated(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, DataTable(DataTableProps{
		Columns: []DataTableColumn{
			{Label: "Name", Sortable: true},
			{Label: "Email", Sortable: true, SortKey: "email_address"},
			{Label: "Role"},
		},
		ActiveSortColumn: "Name",
		ActiveSortDir:    SortAsc,
		SortBaseURL:      "/users",
		Rows: []TableRow{
			SimpleTableRow("Alice", "alice@example.com", "Admin"),
		},
	}))
	// Active column (Name, asc) should link to toggle (desc)
	utils.AssertContains(t, output, "dir=desc")
	utils.AssertContains(t, output, "sort=name")
	// Other sortable column should link to ascending
	utils.AssertContains(t, output, "sort=email_address")
	utils.AssertContains(t, output, "dir=asc")
	// aria-sort should be present
	utils.AssertContains(t, output, `aria-sort="ascending"`)
	// Non-sortable column should NOT have aria-sort
	nameStart := strings.Index(output, "Role")
	if nameStart == -1 {
		t.Fatal("expected Role in output")
	}
}

func TestDataTableSortDirectionToggle(t *testing.T) {
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
	// Active column (desc) should toggle to asc
	utils.AssertContains(t, output, "dir=asc")
	utils.AssertContains(t, output, "sort=name")
	utils.AssertContains(t, output, `aria-sort="descending"`)
}

func TestDataTableNoSortBaseURL(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, DataTable(DataTableProps{
		Columns: []DataTableColumn{
			{Label: "Name", Sortable: true},
		},
		Rows: []TableRow{
			SimpleTableRow("Alice"),
		},
	}))
	// Without SortBaseURL, no sort query links generated (aria-sort is OK)
	utils.AssertNotContains(t, output, "?sort=")
	utils.AssertNotContains(t, output, "&sort=")
}

func TestDataTableEmptyStateRendered(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, DataTable(DataTableProps{
		Columns: []DataTableColumn{
			{Label: "Name"},
		},
		Rows:       []TableRow{},
		EmptyState: SimpleEmptyState("No data available"),
	}))
	utils.AssertContains(t, output, "No data available")
	utils.AssertNotContains(t, output, "<table")
}

func TestDataTableEmptyRowsNoEmptyState(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, DataTable(DataTableProps{
		Columns: []DataTableColumn{
			{Label: "Name"},
		},
		Rows: []TableRow{},
	}))
	// Without EmptyState, still renders the table (with headers only)
	utils.AssertContains(t, output, "<table")
	utils.AssertContains(t, output, "Name")
}

func TestDataTablePaginationRendered(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, DataTable(DataTableProps{
		Columns: []DataTableColumn{
			{Label: "Name"},
		},
		Rows: []TableRow{
			SimpleTableRow("Alice"),
		},
		Pagination: SimpleEmptyState("Page 1 of 5"),
	}))
	utils.AssertContains(t, output, "Page 1 of 5")
}

func TestDataTableNoPaginationWhenEmpty(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, DataTable(DataTableProps{
		Columns: []DataTableColumn{
			{Label: "Name"},
		},
		Rows:       []TableRow{},
		EmptyState: SimpleEmptyState("No data"),
		Pagination: SimpleEmptyState("Page 1"),
	}))
	// When empty, pagination should NOT render
	utils.AssertNotContains(t, output, "Page 1")
}

func TestDataTableCustomSortParam(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, DataTable(DataTableProps{
		Columns: []DataTableColumn{
			{Label: "Name", Sortable: true},
		},
		ActiveSortColumn: "Name",
		ActiveSortDir:    SortAsc,
		SortBaseURL:      "/items",
		SortParam:        "orderBy",
		DirParam:         "order",
		Rows: []TableRow{
			SimpleTableRow("Alice"),
		},
	}))
	utils.AssertContains(t, output, "orderBy=name")
	utils.AssertContains(t, output, "order=desc")
}

func TestDataTableFlushPassedToTable(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, DataTable(DataTableProps{
		Columns: []DataTableColumn{
			{Label: "Name"},
		},
		Rows: []TableRow{
			SimpleTableRow("Alice"),
		},
		Flush: true,
	}))
	// Flush suppresses the wrapper border — should NOT have rounded-lg border
	if strings.Contains(output, "rounded-lg border") {
		t.Error("Flush=true should suppress the table wrapper border")
	}
}

func TestDataTableCaptionRendered(t *testing.T) {
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
	utils.AssertContains(t, output, "sr-only")
}

func TestToggleSortDir(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input SortDirection
		want  SortDirection
	}{
		{SortAsc, SortDesc},
		{SortDesc, SortAsc},
		{SortNone, SortAsc},
	}
	for _, tt := range tests {
		got := toggleSortDir(tt.input)
		if got != tt.want {
			t.Errorf("toggleSortDir(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestColumnSortKey(t *testing.T) {
	t.Parallel()
	// Explicit SortKey
	col := DataTableColumn{Label: "Name", SortKey: "full_name"}
	if got := columnSortKey(col); got != "full_name" {
		t.Errorf("expected full_name, got %s", got)
	}
	// Default to lowercase Label
	col2 := DataTableColumn{Label: "CreatedAt"}
	if got := columnSortKey(col2); got != "createdat" {
		t.Errorf("expected createdat, got %s", got)
	}
}
