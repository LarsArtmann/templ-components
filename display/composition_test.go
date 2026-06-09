package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestCompositionCardWithBadge(t *testing.T) {
	t.Parallel()
	t.Run("card header action renders badge", func(t *testing.T) {
		t.Parallel()
		badge := Badge(BadgeProps{
			Text: activeBadgeText,
			Type: BadgeSuccess,
		})
		props := CardProps{
			Title:        "Service Status",
			HeaderAction: badge,
			Padding:      CardPaddingMD,
		}
		output := utils.Render(t, Card(props))
		utils.AssertContains(t, output, "Service Status")
		utils.AssertContains(t, output, "Active")
		utils.AssertContains(t, output, "bg-green-")
	})
}

func TestCompositionTableWithContent(t *testing.T) {
	t.Parallel()
	t.Run("table renders templ components in cells", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Table(TableProps{
			Headers: []string{tableHeaderName, "Status"},
			Rows: []TableRow{
				{
					Cells: []TableCell{
						{Text: avatarAltAlice},
						{Content: templ.Component(nil)},
					},
				},
			},
		}))
		utils.AssertContains(t, output, "Alice")
		utils.AssertContains(t, output, "<table")
	})
}

func TestCompositionCardWithStatCards(t *testing.T) {
	t.Parallel()
	t.Run("card title and stat card render independently", func(t *testing.T) {
		t.Parallel()
		statProps := StatCardProps{
			Value:  "1,234",
			Label:  "Total Users",
			Change: "+12%",
			Trend:  TrendUp,
		}
		output := utils.Render(t, StatCard(statProps))
		utils.AssertContains(t, output, "1,234")
		utils.AssertContains(t, output, "Total Users")
		utils.AssertContains(t, output, "+12%")
		utils.AssertContains(t, output, "text-green-600")
	})
}

func TestCompositionAccordionWithMultipleItems(t *testing.T) {
	t.Parallel()
	t.Run("accordion items with proper aria controls", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Accordion(AccordionProps{
			Items: []AccordionItem{
				{ID: "item1", Title: "First", Open: true},
				{ID: "item2", Title: "Second"},
			},
		}))
		utils.AssertContains(t, output, `aria-controls="item1-panel"`)
		utils.AssertContains(t, output, `id="item1-trigger"`)
		utils.AssertContains(t, output, `aria-expanded="true"`)
		utils.AssertContains(t, output, `aria-expanded="false"`)
	})
}

func TestCompositionTabsWithActiveState(t *testing.T) {
	t.Parallel()
	t.Run("tabs render with correct active tab", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Tabs(TabsProps{
			ActiveTabID: "overview",
			Tabs: []Tab{
				{ID: "overview", Label: "Overview"},
				{ID: "settings", Label: "Settings"},
			},
		}))
		utils.AssertContains(t, output, `aria-selected="true"`)
		utils.AssertContains(t, output, "Overview")
		utils.AssertContains(t, output, "Settings")
	})
}

func TestCompositionDropdownWithMixedItems(t *testing.T) {
	t.Parallel()
	t.Run("dropdown with internal and external links", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dropdown(DropdownProps{
			BaseProps: utils.BaseProps{ID: "mixed-menu"},
			Label:     dropdownLabelMenu,
			Items: []DropdownItem{
				{Text: "Internal", Href: "/page"},
				{Text: "External", Href: "https://example.com", External: true},
			},
		}))
		utils.AssertContains(t, output, `href="/page"`)
		utils.AssertContains(t, output, `target="_blank"`)
		utils.AssertContains(t, output, `rel="noopener noreferrer"`)
	})
}
