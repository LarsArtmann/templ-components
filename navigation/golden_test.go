package navigation

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

func TestGoldenSidebarNav(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, SidebarNav(SidebarNavProps{
		Brand: templ.Raw(`<span class="font-bold text-white">MyApp</span>`),
		Items: []SidebarNavItem{
			{Label: "Dashboard", Href: "/", Icon: icons.Squares2x2, Active: true},
			{Label: "Users", Href: "/users", Icon: icons.Users},
		},
	}))
	golden.Assert(t, "sidebar_nav", output)
}
