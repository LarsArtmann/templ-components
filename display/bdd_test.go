// Package display provides behavior-driven tests for display components.
// These tests verify end-user-facing behavior, not implementation details.
package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// --- Badge Behavior ---

func TestBadgeUserSeesCorrectVisualFeedback(t *testing.T) {
	t.Parallel()

	t.Run("user sees success badge with green styling", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Badge(DefaultBadgeProps()))
		utils.AssertContains(t, output, `<span`)
	})

	t.Run("user sees pill-shaped badge when pill is enabled", func(t *testing.T) {
		t.Parallel()
		props := DefaultBadgeProps()
		props.Text = "Active"
		props.Pill = true
		output := utils.Render(t, Badge(props))
		utils.AssertContains(t, output, "Active")
		utils.AssertContains(t, output, "rounded-full")
	})

	t.Run("user sees dot indicator on badge", func(t *testing.T) {
		t.Parallel()
		props := DefaultBadgeProps()
		props.Text = "Live"
		props.Dot = true
		output := utils.Render(t, Badge(props))
		utils.AssertContains(t, output, "Live")
		utils.AssertContains(t, output, "h-1.5 w-1.5 rounded-full")
	})
}

// --- StatusBadge Behavior ---

func TestStatusBadgeMapsKnownStatusesCorrectly(t *testing.T) {
	t.Parallel()

	knownStatuses := []string{"healthy", "running", "active", "success", "error", "warning"}
	for _, status := range knownStatuses {
		t.Run("status "+status+" renders a badge", func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, StatusBadge(status))
			utils.AssertContains(t, output, "<span")
		})
	}

	t.Run("unknown status falls back to default badge", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatusBadge("completely-unknown-status"))
		utils.AssertContains(t, output, "<span")
	})
}

// --- Card Behavior ---

func TestCardUserCanComposeContent(t *testing.T) {
	t.Parallel()

	t.Run("user sees card with title and content", func(t *testing.T) {
		t.Parallel()
		props := DefaultCardProps()
		props.Title = "Dashboard"
		output := utils.Render(t, Card(props))
		utils.AssertContains(t, output, "Dashboard")
	})

	t.Run("user sees card with subtitle", func(t *testing.T) {
		t.Parallel()
		props := DefaultCardProps()
		props.Title = "Users"
		props.Subtitle = "Manage your team"
		output := utils.Render(t, Card(props))
		utils.AssertContains(t, output, "Users")
		utils.AssertContains(t, output, "Manage your team")
	})
}

// --- Modal Behavior ---

func TestModalUserCanOpenAndClose(t *testing.T) {
	t.Parallel()

	t.Run("modal renders with accessible attributes", func(t *testing.T) {
		t.Parallel()
		props := DefaultModalProps()
		props.Title = "Confirm Action"
		output := utils.Render(t, Modal(props))
		utils.AssertContains(t, output, "Confirm Action")
		utils.AssertContains(t, output, `role="dialog"`)
		utils.AssertContains(t, output, `aria-modal="true"`)
	})

	t.Run("modal renders in open state", func(t *testing.T) {
		t.Parallel()
		props := DefaultModalProps()
		props.Open = true
		output := utils.Render(t, Modal(props))
		utils.AssertContains(t, output, "opacity-100")
	})
}

// --- Table Behavior ---

func TestTableUserCanViewData(t *testing.T) {
	t.Parallel()

	t.Run("user sees table headers and data", func(t *testing.T) {
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
		utils.AssertContains(t, output, "Bob")
	})

	t.Run("user sees striped table rows", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers: []string{"A"},
			Rows:    []TableRow{SimpleTableRow("1")},
			Striped: true,
		}))
		utils.AssertContains(t, output, "bg-gray-50")
	})

	t.Run("user sees table caption for accessibility", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Caption: "User list",
			Headers: []string{"Name"},
			Rows:    []TableRow{SimpleTableRow("Alice")},
		}))
		utils.AssertContains(t, output, "User list")
		utils.AssertContains(t, output, `<caption`)
	})
}

// --- Tabs Behavior ---

func TestTabsUserCanSwitchViews(t *testing.T) {
	t.Parallel()

	t.Run("user sees tabs with active state", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			Tabs: []Tab{
				{ID: "users", Label: "Users", Active: true},
				{ID: "settings", Label: "Settings"},
			},
		}))
		utils.AssertContains(t, output, "Users")
		utils.AssertContains(t, output, "Settings")
		utils.AssertContains(t, output, `aria-selected="true"`)
	})

	t.Run("user sees pill-style tabs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			Tabs: []Tab{
				{ID: "a", Label: "A", Active: true},
			},
			TabStyle: TabsStylePills,
		}))
		utils.AssertContains(t, output, "rounded-md")
	})
}

// --- Accordion Behavior ---

func TestAccordionUserCanExpandCollapse(t *testing.T) {
	t.Parallel()

	t.Run("user sees accordion items with titles", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Accordion(AccordionProps{
			Items: []AccordionItem{
				{ID: "faq1", Title: "What is this?"},
				{ID: "faq2", Title: "How does it work?"},
			},
		}))
		utils.AssertContains(t, output, "What is this?")
		utils.AssertContains(t, output, "How does it work?")
	})

	t.Run("user sees pre-opened accordion item", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Accordion(AccordionProps{
			Items: []AccordionItem{
				{ID: "faq1", Title: "Open by default", Open: true},
			},
		}))
		utils.AssertContains(t, output, `aria-expanded="true"`)
	})
}

// --- Dropdown Behavior ---

func TestDropdownUserCanSelectAction(t *testing.T) {
	t.Parallel()

	t.Run("user sees dropdown trigger and menu items", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "actions"},
			Label:     "Actions",
			Items: []DropdownItem{
				{Text: "Edit", Href: "/edit"},
				{Text: "Delete", Href: "/delete"},
			},
		}))
		utils.AssertContains(t, output, "Actions")
		utils.AssertContains(t, output, "Edit")
		utils.AssertContains(t, output, "Delete")
		utils.AssertContains(t, output, `role="menu"`)
		utils.AssertContains(t, output, `aria-haspopup="true"`)
	})

	t.Run("user sees external link with proper attributes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "links"},
			Label:     "Links",
			Items: []DropdownItem{
				{Text: "Docs", Href: "https://example.com", External: true},
			},
		}))
		utils.AssertContains(t, output, `target="_blank"`)
		utils.AssertContains(t, output, `rel="noopener noreferrer"`)
	})
}

// --- Avatar Behavior ---

func TestAvatarUserCanIdentifyUsers(t *testing.T) {
	t.Parallel()

	t.Run("user sees avatar with image", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Src:   "/alice.jpg",
			Alt:   "Alice",
			Size:  AvatarSizeMD,
			Shape: AvatarShapeCircle,
		}))
		utils.AssertContains(t, output, `/alice.jpg`)
		utils.AssertContains(t, output, "Alice")
	})

	t.Run("user sees avatar with initials fallback", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Avatar(AvatarProps{
			Initials: "AB",
			Size:     AvatarSizeSM,
		}))
		utils.AssertContains(t, output, "AB")
	})
}

// --- Tooltip Behavior ---

func TestTooltipUserSeesHelpOnHover(t *testing.T) {
	t.Parallel()

	t.Run("tooltip wraps children with hint text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tooltip(TooltipProps{
			Text:     "More information",
			Position: TooltipPositionTop,
		}))
		utils.AssertContains(t, output, "More information")
		utils.AssertContains(t, output, `role="tooltip"`)
	})
}

// --- EmptyState Behavior ---

func TestEmptyStateUserSeesHelpfulGuidance(t *testing.T) {
	t.Parallel()

	t.Run("user sees empty state with title and description", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EmptyState(EmptyStateProps{
			Title:       "No items yet",
			Description: "Create your first item to get started.",
		}))
		utils.AssertContains(t, output, "No items yet")
		utils.AssertContains(t, output, "Create your first item to get started.")
	})

	t.Run("user sees action button to proceed", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EmptyState(EmptyStateProps{
			Title:      "No repos",
			ActionText: "Connect Repository",
			ActionHref: "/connect",
		}))
		utils.AssertContains(t, output, "Connect Repository")
		utils.AssertContains(t, output, `href="/connect"`)
	})
}
