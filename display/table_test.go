// Package display provides tests for display components.
package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestTableRender(t *testing.T) {
	t.Parallel()
	t.Run("basic table", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers: []string{"Name", "Email"},
			Rows: []TableRow{
				SimpleTableRow("Alice", "alice@example.com"),
				SimpleTableRow("Bob", "bob@example.com"),
			},
		}))
		utils.AssertContains(t, output, "Name")
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
			Headers: []string{"Name"},
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
}
