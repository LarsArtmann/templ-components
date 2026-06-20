package display

import (
	"testing"

	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

func TestEmptyStateCoverage(t *testing.T) {
	t.Parallel()

	t.Run("with icon and action link", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EmptyState(EmptyStateProps{
			Icon:        icons.Users,
			Title:       "No users found",
			Description: "Add your first user to get started",
			ActionText:  "Add user",
			ActionHref:  "/users/new",
		}))
		utils.AssertContains(t, output, "No users found")
		utils.AssertContains(t, output, "Add user")
		utils.AssertContains(t, output, "/users/new")
	})

	t.Run("with action as button (no href)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EmptyState(EmptyStateProps{
			Title:      "Empty",
			ActionText: "Click me",
		}))
		utils.AssertContains(t, output, "Click me")
	})

	t.Run("with ID and class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EmptyState(EmptyStateProps{
			Title:     "Test",
			BaseProps: utils.BaseProps{ID: "empty", Class: "py-12"},
		}))
		utils.AssertContains(t, output, `id="empty"`)
		utils.AssertContains(t, output, "py-12")
	})

	t.Run("minimal with just title", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EmptyState(EmptyStateProps{
			Title: "Nothing here",
		}))
		utils.AssertContains(t, output, "Nothing here")
	})
}

func TestStatusBadgeExtraCoverage(t *testing.T) {
	t.Parallel()

	for _, status := range []string{"active", "pending", "inactive", "archived", "draft", "unknown-status"} {
		t.Run(status, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, StatusBadge(status))
			utils.AssertContains(t, output, status)
		})
	}
}

func TestTableExtraCoverage(t *testing.T) {
	t.Parallel()

	t.Run("with caption", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers: []string{"Name", "Email"},
			Caption: "User list",
			Rows: []TableRow{
				SimpleTableRow("Alice", "alice@example.com"),
				SimpleTableRow("Bob", "bob@example.com"),
			},
		}))
		utils.AssertContains(t, output, "User list")
		utils.AssertContains(t, output, "Alice")
		utils.AssertContains(t, output, "alice@example.com")
	})

	t.Run("with ID", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers:   []string{"A"},
			Rows:      []TableRow{SimpleTableRow("1")},
			BaseProps: utils.BaseProps{ID: "tbl"},
		}))
		utils.AssertContains(t, output, `id="tbl"`)
	})

	t.Run("rows shorter than headers get padded", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers: []string{"A", "B", "C"},
			Rows:    []TableRow{SimpleTableRow("1")},
		}))
		utils.AssertContains(t, output, "1")
	})
}
