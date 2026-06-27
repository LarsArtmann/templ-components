package navigation

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

func TestSidebarNavRender(t *testing.T) {
	t.Parallel()

	t.Run("renders items with labels", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SidebarNav(SidebarNavProps{
			Items: []SidebarNavItem{
				{Label: "Dashboard", Href: "/"},
				{Label: "Users", Href: "/users"},
			},
		}))
		utils.AssertContains(t, output, "Dashboard")
		utils.AssertContains(t, output, "Users")
		utils.AssertContains(t, output, `href="/"`)
		utils.AssertContains(t, output, `href="/users"`)
	})

	t.Run("active item gets active styling and aria-current", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SidebarNav(SidebarNavProps{
			Items: []SidebarNavItem{
				{Label: "Dashboard", Href: "/", Active: true},
				{Label: "Users", Href: "/users"},
			},
		}))
		utils.AssertContains(t, output, "bg-blue-600")
		utils.AssertContains(t, output, `aria-current="page"`)
	})

	t.Run("currentpath auto-detects active item", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SidebarNav(SidebarNavProps{
			CurrentPath: "/users",
			Items: []SidebarNavItem{
				{Label: "Dashboard", Href: "/"},
				{Label: "Users", Href: "/users"},
			},
		}))
		utils.AssertContains(t, output, "bg-blue-600")
		utils.AssertContains(t, output, `aria-current="page"`)
	})

	t.Run("icon renders when set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SidebarNav(SidebarNavProps{
			Items: []SidebarNavItem{
				{Label: "Users", Href: "/users", Icon: icons.Users},
			},
		}))
		utils.AssertContains(t, output, "<svg")
		utils.AssertContains(t, output, "Users")
	})

	t.Run("no icon omits svg", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SidebarNav(SidebarNavProps{
			Items: []SidebarNavItem{
				{Label: "Settings", Href: "/settings"},
			},
		}))
		utils.AssertNotContains(t, output, "<svg")
	})

	t.Run("brand slot renders at top", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SidebarNav(SidebarNavProps{
			Brand: templ.Raw(`<span class="font-bold text-white">MyApp</span>`),
			Items: []SidebarNavItem{{Label: "Home", Href: "/"}},
		}))
		utils.AssertContains(t, output, "MyApp")
	})

	t.Run("footer slot renders at bottom", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SidebarNav(SidebarNavProps{
			Items:  []SidebarNavItem{{Label: "Home", Href: "/"}},
			Footer: templ.Raw(`<span>v1.0.0</span>`),
		}))
		utils.AssertContains(t, output, "v1.0.0")
	})

	t.Run("default aria-label is Sidebar", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SidebarNav(SidebarNavProps{
			Items: []SidebarNavItem{{Label: "Home", Href: "/"}},
		}))
		utils.AssertContains(t, output, `aria-label="Sidebar"`)
	})

	t.Run("propagates BaseProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SidebarNav(SidebarNavProps{
			Items: []SidebarNavItem{{Label: "Home", Href: "/"}},
			BaseProps: utils.BaseProps{
				ID:        "sidebar",
				Class:     "w-72",
				AriaLabel: "Main navigation",
			},
		}))
		utils.AssertContains(t, output, `id="sidebar"`)
		utils.AssertContains(t, output, "w-72")
		utils.AssertContains(t, output, `aria-label="Main navigation"`)
	})
}
