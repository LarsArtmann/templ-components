// Package display provides tests for display components.
package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestTableRender(t *testing.T) {
	t.Parallel()
	t.Run("basic table", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers: []string{tableHeaderName, "Email"},
			Rows: []TableRow{
				SimpleTableRow("Alice", "alice@example.com"),
				SimpleTableRow("Bob", "bob@example.com"),
			},
		}))
		utils.AssertContains(t, output, tableHeaderName)
		utils.AssertContains(t, output, "Email")
		utils.AssertContains(t, output, "Alice")
		utils.AssertContains(t, output, "bob@example.com")
		utils.AssertContains(t, output, "scope=\"col\"")
	})

	t.Run("striped rows", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers: []string{"A"},
			Rows:    []TableRow{SimpleTableRow("1"), SimpleTableRow("2")},
			Striped: true,
		}))
		utils.AssertContains(t, output, "bg-gray-50")
	})

	t.Run("with caption", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Caption: "User list",
			Headers: []string{tableHeaderName},
			Rows:    []TableRow{SimpleTableRow("Alice")},
		}))
		utils.AssertContains(t, output, "User list")
		utils.AssertContains(t, output, "<caption")
	})

	t.Run("with custom id", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			BaseProps: utils.BaseProps{ID: "users-table"},
			Headers:   []string{"Name"},
			Rows:      []TableRow{SimpleTableRow("Alice")},
		}))
		utils.AssertContains(t, output, `id="users-table"`)
	})

	t.Run("row with fewer cells than headers is padded", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers: []string{"A", "B", "C"},
			Rows:    []TableRow{SimpleTableRow("only-one")},
		}))
		utils.AssertContains(t, output, "only-one")
	})

	t.Run("row with more cells than headers is truncated", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers: []string{"A"},
			Rows:    []TableRow{SimpleTableRow("1", "2", "3")},
		}))
		utils.AssertContains(t, output, "1")
		utils.AssertNotContains(t, output, ">2<")
		utils.AssertNotContains(t, output, ">3<")
	})

	t.Run("no headers renders all cells", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Rows: []TableRow{SimpleTableRow("x", "y", "z")},
		}))
		utils.AssertContains(t, output, "x")
		utils.AssertContains(t, output, "y")
		utils.AssertContains(t, output, "z")
	})

	t.Run("bordered table", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers:  []string{"A"},
			Rows:     []TableRow{SimpleTableRow("1")},
			Bordered: true,
		}))
		utils.AssertContains(t, output, "border border-gray-200")
	})

	t.Run("hover table", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers: []string{"A"},
			Rows:    []TableRow{SimpleTableRow("1")},
			Hover:   true,
		}))
		utils.AssertContains(t, output, "hover:bg-gray-100")
	})

	t.Run("default props", func(t *testing.T) {
		t.Parallel()

		props := DefaultTableProps()
		if props.Striped != true {
			t.Error("DefaultTableProps().Striped should be true")
		}
	})

	t.Run("Body slot overrides Rows", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers: []string{"Name", "Status"},
			Rows:    []TableRow{SimpleTableRow("should-not-appear")},
			Body: templ.Raw(
				"<tr><td>Alice</td><td><span class=\"badge\">Admin</span></td></tr><tr><td>Bob</td><td>User</td></tr>",
			),
		}))
		utils.AssertContains(t, output, "Alice")
		utils.AssertContains(t, output, "Bob")
		utils.AssertContains(t, output, "<span class=\"badge\">Admin</span>")
		utils.AssertNotContains(t, output, "should-not-appear")
		utils.AssertContains(t, output, "<tbody")
	})

	t.Run("Body slot renders headers and tbody wrapper", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers: []string{"Name"},
			Body:    templ.Raw("<tr><td>Custom row</td></tr>"),
		}))
		utils.AssertContains(t, output, "Name")
		utils.AssertContains(t, output, "Custom row")
		utils.AssertContains(t, output, "<thead")
		utils.AssertContains(t, output, "<tbody")
		utils.AssertContains(t, output, "divide-y")
	})
}

func TestTableFlush(t *testing.T) {
	t.Parallel()

	t.Run("default table has wrapper border", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers: []string{"A"},
			Rows:    []TableRow{SimpleTableRow("1")},
		}))
		utils.AssertContains(t, output, "overflow-x-auto")
		utils.AssertContains(t, output, "rounded-lg")
		utils.AssertContains(t, output, "border border-gray-200 dark:border-gray-700")
	})

	t.Run("flush table suppresses wrapper border", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers: []string{"A"},
			Rows:    []TableRow{SimpleTableRow("1")},
			Flush:   true,
		}))
		utils.AssertContains(t, output, "overflow-x-auto")
		utils.AssertNotContains(t, output, "rounded-lg")
		utils.AssertNotContains(t, output, "border border-gray-200")
	})
}

func TestTableCellPaddingOption(t *testing.T) {
	t.Parallel()

	t.Run("default uses comfortable padding", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers: []string{"A"},
			Rows:    []TableRow{SimpleTableRow("1")},
		}))
		utils.AssertContainsAll(t, output, "px-4", "py-3")
		utils.AssertNotContains(t, output, "py-2")
	})

	t.Run("compact uses py-2 padding", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers:     []string{"A"},
			Rows:        []TableRow{SimpleTableRow("1")},
			CellPadding: TableCellPaddingCompact,
		}))
		utils.AssertContainsAll(t, output, "px-4", "py-2")
		utils.AssertNotContains(t, output, "py-3")
	})

	t.Run("comfortable explicitly uses py-3", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers:     []string{"A"},
			Rows:        []TableRow{SimpleTableRow("1")},
			CellPadding: TableCellPaddingComfortable,
		}))
		utils.AssertContainsAll(t, output, "px-4", "py-3")
	})

	t.Run("invalid cell padding falls back to comfortable", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers:     []string{"A"},
			Rows:        []TableRow{SimpleTableRow("1")},
			CellPadding: TableCellPadding("bogus"),
		}))
		utils.AssertContainsAll(t, output, "px-4", "py-3")
		utils.AssertNotContains(t, output, "py-2")
	})

	t.Run("compact applies to typed headers too", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			TypedHeaders: []TableHeader{{Label: "A"}},
			Rows:         []TableRow{SimpleTableRow("1")},
			CellPadding:  TableCellPaddingCompact,
		}))
		utils.AssertContainsAll(t, output, "px-4", "py-2")
		utils.AssertNotContains(t, output, "py-3")
	})
}

func TestTableTypedHeaders(t *testing.T) {
	t.Parallel()
	t.Run("sortable headers with aria-sort", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			TypedHeaders: []TableHeader{
				{Label: "Name", Sortable: true, SortDirection: SortAsc, Href: "/sort?col=name&dir=asc"},
				{Label: "Date", Sortable: true, SortDirection: SortNone, Href: "/sort?col=date&dir=asc"},
				{Label: "Status"},
			},
			Rows: []TableRow{
				SimpleTableRow("Alice", "2024-01-01", "Active"),
			},
		}))
		utils.AssertContains(t, output, `aria-sort="ascending"`)
		utils.AssertContains(t, output, `aria-sort="none"`)
		utils.AssertContains(t, output, `/sort?col=name`)
		utils.AssertContains(t, output, "Name")
		utils.AssertContains(t, output, "↑")
		utils.AssertNotContains(t, output, "Status</a>")
	})
	t.Run("non-sortable header has no aria-sort", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			TypedHeaders: []TableHeader{
				{Label: "Name"},
			},
		}))
		utils.AssertNotContains(t, output, "aria-sort")
	})
}
