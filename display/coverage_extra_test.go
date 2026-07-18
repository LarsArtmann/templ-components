package display

import (
	"testing"
	"time"

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
}

func TestCountBadgeMaxOverflow(t *testing.T) {
	t.Parallel()
	t.Run("shows N+ when count exceeds max", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CountBadge(CountBadgeProps{Count: 150, Max: 99}))
		utils.AssertContains(t, output, "99+")
	})
	t.Run("shows exact count when under max", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CountBadge(CountBadgeProps{Count: 5}))
		utils.AssertContains(t, output, "5")
		utils.AssertNotContains(t, output, "+")
	})
}

func TestFormatRelativeTimeBoundaries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		from time.Time
		want string
	}{
		{"just now (under 60s)", time.Now().Add(-59 * time.Second), "just now"},
		{"1 minute ago", time.Now().Add(-60 * time.Second), "1 minute ago"},
		{"59 minutes ago", time.Now().Add(-59 * time.Minute), "59 minutes ago"},
		{"1 hour ago", time.Now().Add(-60 * time.Minute), "1 hour ago"},
		{"23 hours ago", time.Now().Add(-23 * time.Hour), "23 hours ago"},
		{"1 day ago", time.Now().Add(-24 * time.Hour), "1 day ago"},
		{"6 days ago", time.Now().Add(-144 * time.Hour), "6 days ago"},
		{"7 days ago (1 week)", time.Now().Add(-168 * time.Hour), "1 week ago"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := formatRelativeTime(tt.from, time.Now())
			if got != tt.want {
				t.Errorf("formatRelativeTime() = %q, want %q", got, tt.want)
			}
		})
	}
}
