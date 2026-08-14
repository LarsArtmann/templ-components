package navigation

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

func TestSidebarNavGroupedCoverage(t *testing.T) {
	t.Parallel()

	groupedItems := []SidebarNavItem{
		{Label: "Overview", Href: "/"},
		{Label: "Dashboard", Href: "/dash", Section: "Main"},
		{Label: "Reports", Href: "/reports", Section: "Main"},
		{Label: "Profile", Href: "/settings/profile", Section: "Account"},
		{Label: "Security", Href: "/settings/security", Section: "Account"},
	}

	t.Run("sections render as grouped details", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SidebarNav(SidebarNavProps{Items: groupedItems}))
		utils.AssertContains(t, output, "Main")
		utils.AssertContains(t, output, "Account")
		utils.AssertContains(t, output, "<details")
		utils.AssertContains(t, output, "<summary")
		utils.AssertContains(t, output, "Overview")
	})

	t.Run("group with active item is open", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SidebarNav(SidebarNavProps{
			Items:       groupedItems,
			CurrentPath: "/reports",
		}))
		utils.AssertContains(t, output, "open")
		utils.AssertContains(t, output, `aria-current="page"`)
	})

	t.Run("explicit active flag inside group", func(t *testing.T) {
		t.Parallel()

		items := []SidebarNavItem{
			{Label: "A", Href: "/a", Section: "S1", Active: true},
			{Label: "B", Href: "/b", Section: "S1"},
		}

		output := utils.Render(t, SidebarNav(SidebarNavProps{Items: items}))
		utils.AssertContains(t, output, `aria-current="page"`)
		utils.AssertContains(t, output, "open")
	})

	t.Run("brand header footer slots and icons", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SidebarNav(SidebarNavProps{
			BaseProps: utils.BaseProps{
				ID:        "sb-1",
				Class:     "w-72",
				AriaLabel: "Admin sidebar",
				Attrs:     templ.Attributes{"data-ctx": "admin"},
			},
			Brand:  templ.Raw(`<strong data-testid="brand">Acme</strong>`),
			Header: templ.Raw(`<input data-testid="filter" />`),
			Footer: templ.Raw(`<div data-testid="foot">v1</div>`),
			Items: []SidebarNavItem{
				{Label: "Home", Href: "/", Icon: icons.Home},
				{Label: "Users", Href: "/users", Icon: icons.Users},
			},
		}))
		utils.AssertContains(t, output, `id="sb-1"`)
		utils.AssertContains(t, output, `aria-label="Admin sidebar"`)
		utils.AssertContains(t, output, `data-ctx="admin"`)
		utils.AssertContains(t, output, `data-testid="brand"`)
		utils.AssertContains(t, output, `data-testid="filter"`)
		utils.AssertContains(t, output, `data-testid="foot"`)
		utils.AssertContains(t, output, "<svg")
	})

	t.Run("empty items renders shell", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SidebarNav(SidebarNavProps{}))
		utils.AssertContains(t, output, `aria-label="Sidebar"`)
		utils.AssertContains(t, output, "<nav")
	})
}
