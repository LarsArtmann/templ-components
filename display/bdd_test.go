// Package display provides behavior-driven tests for display components.
// These tests verify end-user-facing behavior, not implementation details.
package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

const (
	cardTitleUsers   = "Users"
	statusActiveText = "active"
	accordionFAQ1    = "faq1"
	statusErrorText  = "error"
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
		props.Text = activeBadgeText
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

	knownStatuses := []string{
		"healthy",
		"running",
		statusActiveText,
		"success",
		statusErrorText,
		"warning",
	}
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
		props.Title = cardTitleUsers
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
		props.ID = "bdd-modal"
		props.Title = "Confirm Action"
		output := utils.Render(t, Modal(props))
		utils.AssertContains(t, output, "Confirm Action")
		utils.AssertContains(t, output, `role="dialog"`)
		utils.AssertContains(t, output, `aria-modal="true"`)
	})

	t.Run("modal renders in open state", func(t *testing.T) {
		t.Parallel()
		props := DefaultModalProps()
		props.ID = "bdd-modal-open"
		props.Open = true
		output := utils.Render(t, Modal(props))
		utils.AssertContains(t, output, "opacity-100")
	})
}

// --- Table Behavior ---

func TestTableUserCanViewData(t *testing.T) {
	t.Parallel()

	t.Run("user sees striped table rows", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers: []string{"A"},
			Rows:    []TableRow{SimpleTableRow("1")},
			Striped: true,
		}))
		utils.AssertContains(t, output, "bg-gray-50")
	})
}

// --- Accordion Behavior ---

func TestAccordionUserCanExpandCollapse(t *testing.T) {
	t.Parallel()

	t.Run("user sees accordion items with titles", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Accordion(AccordionProps{
			Items: []AccordionItem{
				{ID: accordionFAQ1, Title: "What is this?"},
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

func TestStatCardUserCanClickToFilter(t *testing.T) {
	t.Parallel()

	t.Run("user sees a clickable stat card with href", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatCard(StatCardProps{
			Label: "Active",
			Value: "42",
			Href:  "/?activity=active",
		}))
		utils.AssertContains(t, output, `href="/?activity=active"`)
		utils.AssertContains(t, output, "42")
		utils.AssertContains(t, output, "Active")
	})

	t.Run("user can keyboard-focus the linked stat card", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatCard(StatCardProps{
			Value: "10",
			Label: "Total",
			Href:  "/all",
		}))
		utils.AssertContains(t, output, "focus-visible:ring")
	})
}

func TestGridUserSeesResponsiveLayout(t *testing.T) {
	t.Parallel()

	t.Run("user sees a grid that stacks on mobile", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Grid(GridProps{Cols: GridCols3}))
		utils.AssertContains(t, output, "grid-cols-1")
		utils.AssertContains(t, output, "lg:grid-cols-3")
	})

	t.Run("grid falls back to 3 cols for unknown value", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Grid(GridProps{Cols: GridCols("bogus")}))
		utils.AssertContains(t, output, "lg:grid-cols-3")
	})
}
