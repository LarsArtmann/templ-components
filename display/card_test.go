// Package display provides tests for display components like Badge, Card, Modal, and EmptyState.
package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

func TestCardRender(t *testing.T) {
	t.Parallel()
	t.Run("basic card with title", func(t *testing.T) {
		t.Parallel()
		props := CardProps{
			Title:   cardTitleUsers,
			Padding: CardPaddingMD,
		}
		output := utils.Render(t, Card(props))
		utils.AssertContains(t, output, cardTitleUsers)
		utils.AssertContains(t, output, "bg-white")
		utils.AssertContains(t, output, "rounded-lg")
	})

	t.Run("card with custom class and id", func(t *testing.T) {
		t.Parallel()
		props := CardProps{
			BaseProps: utils.BaseProps{
				ID:    "my-card",
				Class: cssClassMt4,
			},
			Title:   "Test",
			Padding: CardPaddingMD,
		}
		output := utils.Render(t, Card(props))
		utils.AssertContains(t, output, `id="my-card"`)
		utils.AssertContains(t, output, cssClassMt4)
	})

	t.Run("simple card", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SimpleCard(DefaultSimpleCardProps()))
		utils.AssertContains(t, output, "bg-white")
		utils.AssertContains(t, output, "rounded-lg")
	})

	t.Run("simple card with custom class and id", func(t *testing.T) {
		t.Parallel()
		props := SimpleCardProps{
			BaseProps: utils.BaseProps{
				ID:    "simple-card",
				Class: cssClassMt4,
			},
		}
		output := utils.Render(t, SimpleCard(props))
		utils.AssertContains(t, output, `id="simple-card"`)
		utils.AssertContains(t, output, cssClassMt4)
	})
}

func TestStatCardRender(t *testing.T) {
	t.Parallel()

	t.Run("trend up shows Increased by", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatCard(StatCardProps{
			Value:  "1234",
			Label:  "Users",
			Change: "12%",
			Trend:  TrendUp,
		}))
		utils.AssertContains(t, output, "1234")
		utils.AssertContains(t, output, "text-green-600")
		utils.AssertContains(t, output, "Increased by")
		utils.AssertNotContains(t, output, "Decreased by")
	})

	t.Run("trend down shows Decreased by", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatCard(StatCardProps{
			Value:  "500",
			Label:  "Errors",
			Change: "5%",
			Trend:  TrendDown,
		}))
		utils.AssertContains(t, output, "text-red-600")
		utils.AssertContains(t, output, "Decreased by")
		utils.AssertNotContains(t, output, "Increased by")
	})

	t.Run("trend none shows no direction text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatCard(StatCardProps{
			Value:  "99.9%",
			Label:  "Uptime",
			Change: "stats",
			Trend:  TrendNone,
		}))
		utils.AssertNotContains(t, output, "Increased by")
		utils.AssertNotContains(t, output, "Decreased by")
		utils.AssertContains(t, output, "99.9%")
	})

	t.Run("no change hides trend indicator", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatCard(StatCardProps{
			Value: "100",
			Label: "Total",
		}))
		utils.AssertContains(t, output, "100")
		utils.AssertNotContains(t, output, "sr-only")
	})

	t.Run("default props has TrendNone", func(t *testing.T) {
		t.Parallel()
		props := DefaultStatCardProps()
		if props.Trend != TrendNone {
			t.Errorf("DefaultStatCardProps().Trend = %q, want %q", props.Trend, TrendNone)
		}
	})

	t.Run("icon renders a leading icon tile", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatCard(StatCardProps{
			Label: "Users",
			Value: "1,204",
			Icon:  icons.Users,
		}))
		utils.AssertContains(t, output, "1,204")
		utils.AssertContains(t, output, "<svg")
		utils.AssertContains(t, output, "bg-blue-50")
		utils.AssertContains(t, output, "flex items-center gap-4")
	})

	t.Run("no icon keeps plain layout without icon tile", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, StatCard(StatCardProps{
			Label: "Users",
			Value: "10",
		}))
		utils.AssertNotContains(t, output, "bg-blue-50")
		utils.AssertNotContains(t, output, "<svg")
	})
}

func TestEmptyStateRender(t *testing.T) {
	t.Parallel()
	t.Run("with action link", func(t *testing.T) {
		t.Parallel()
		props := EmptyStateProps{
			BaseProps:   utils.BaseProps{},
			Title:       "No repos",
			Description: "Create your first repo",
			Icon:        icons.Folder,
			ActionText:  "Add Repo",
			ActionHref:  "/repos/new",
		}
		output := utils.Render(t, EmptyState(props))
		utils.AssertContains(t, output, "No repos")
		utils.AssertContains(t, output, "Create your first repo")
		utils.AssertContains(t, output, "Add Repo")
		utils.AssertContains(t, output, `href="/repos/new"`)
	})

	t.Run("with custom icon", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EmptyState(EmptyStateProps{
			Title: "Inbox empty",
			Icon:  icons.Inbox,
		}))
		utils.AssertContains(t, output, "Inbox empty")
	})

	t.Run("with action attrs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EmptyState(EmptyStateProps{
			Title:       "No data",
			Icon:        icons.Folder,
			ActionText:  "Add",
			ActionHref:  "/add",
			ActionAttrs: templ.Attributes{"data-testid": "add-btn"},
		}))
		utils.AssertContains(t, output, `data-testid="add-btn"`)
	})

	t.Run("simple empty state", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SimpleEmptyState("Nothing here"))
		utils.AssertContains(t, output, "Nothing here")
	})

	t.Run("default props", func(t *testing.T) {
		t.Parallel()
		props := DefaultEmptyStateProps()
		if props.Icon != icons.Inbox {
			t.Errorf("DefaultEmptyStateProps().Icon = %q, want %q", props.Icon, icons.Inbox)
		}
	})

	t.Run("without icon", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EmptyState(EmptyStateProps{
			Title: "No data",
		}))
		utils.AssertContains(t, output, "No data")
	})

	t.Run("empty state with button action (no href)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EmptyState(EmptyStateProps{
			Title:      "No items",
			Icon:       icons.Folder,
			ActionText: "Add Item",
		}))
		utils.AssertContains(t, output, "Add Item")
		utils.AssertContains(t, output, "<button")
		utils.AssertContains(t, output, `type="button"`)
		utils.AssertNotContains(t, output, "<a ")
	})

	t.Run("without action", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EmptyState(EmptyStateProps{
			Title:       "Empty",
			Description: "Nothing here yet",
			Icon:        icons.Folder,
		}))
		utils.AssertContains(t, output, "Nothing here yet")
		utils.AssertNotContains(t, output, "<a ")
	})
}

func TestNormalizeTrend(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input TrendDirection
		want  TrendDirection
	}{
		{TrendUp, TrendUp},
		{TrendDown, TrendDown},
		{TrendNone, TrendNone},
		{TrendDirection("invalid"), TrendNone},
		{TrendDirection(""), TrendNone},
	}
	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			t.Parallel()
			got := normalizeTrend(tt.input)
			if got != tt.want {
				t.Errorf("normalizeTrend(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestCardFeatures(t *testing.T) {
	t.Parallel()

	t.Run("card with subtitle", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Card(CardProps{
			Title:    "Users",
			Subtitle: "Manage your team",
			Padding:  CardPaddingMD,
		}))
		utils.AssertContains(t, output, "Users")
		utils.AssertContains(t, output, "Manage your team")
	})

	t.Run("card with footer", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Card(CardProps{
			Title:   "Settings",
			Padding: CardPaddingMD,
			Footer:  templ.Raw("<div>Footer content</div>"),
		}))
		utils.AssertContains(t, output, "Footer content")
	})

	t.Run("card with header action", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Card(CardProps{
			Title:        "Projects",
			Padding:      CardPaddingMD,
			HeaderAction: templ.Raw("<button>Add</button>"),
		}))
		utils.AssertContains(t, output, "Add")
	})

	t.Run("card without title", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Card(CardProps{
			Padding: CardPaddingMD,
		}))
		utils.AssertContains(t, output, "bg-white")
		utils.AssertNotContains(t, output, "font-semibold")
	})
}
