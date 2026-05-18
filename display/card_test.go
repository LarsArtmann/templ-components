// Package display provides tests for display components like Badge, Card, Modal, and EmptyState.
package display

import (
	"testing"

	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

func TestCardRender(t *testing.T) {
	t.Parallel()
	t.Run("basic card with title", func(t *testing.T) {
		t.Parallel()
		props := CardProps{
			Title:   "Users",
			Padding: CardPaddingMD,
		}
		output := utils.Render(t, Card(props))
		utils.AssertContains(t, output, "Users")
		utils.AssertContains(t, output, "bg-white")
		utils.AssertContains(t, output, "rounded-lg")
	})

	t.Run("card with custom class and id", func(t *testing.T) {
		t.Parallel()
		props := CardProps{
			BaseProps: utils.BaseProps{
				ID:    "my-card",
				Class: "mt-4",
			},
			Title:   "Test",
			Padding: CardPaddingMD,
		}
		output := utils.Render(t, Card(props))
		utils.AssertContains(t, output, `id="my-card"`)
		utils.AssertContains(t, output, "mt-4")
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
				Class: "mt-4",
			},
		}
		output := utils.Render(t, SimpleCard(props))
		utils.AssertContains(t, output, `id="simple-card"`)
		utils.AssertContains(t, output, "mt-4")
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
}

func TestEmptyStateRender(t *testing.T) {
	t.Parallel()
	t.Run("with action link", func(t *testing.T) {
		t.Parallel()
		props := EmptyStateProps{
			BaseProps:   utils.BaseProps{},
			Title:       "No repos",
			Description: "",
			Icon:        icons.Folder,
			ActionText:  "Add Repo",
			ActionHref:  "/repos/new",
			ActionAttrs: nil,
		}
		output := utils.Render(t, EmptyState(props))
		utils.AssertContains(t, output, "No repos")
		utils.AssertContains(t, output, "Add Repo")
		utils.AssertContains(t, output, `href="/repos/new"`)
	})

	t.Run("simple empty state", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SimpleEmptyState("Nothing here"))
		utils.AssertContains(t, output, "Nothing here")
	})
}
