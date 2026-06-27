package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestPageHeaderRender(t *testing.T) {
	t.Parallel()

	t.Run("title only", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PageHeader(PageHeaderProps{
			Title: "Dashboard",
		}))
		utils.AssertContains(t, output, "Dashboard")
		utils.AssertContains(t, output, "<h1")
	})

	t.Run("title and subtitle", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PageHeader(PageHeaderProps{
			Title:    "Users",
			Subtitle: "Manage user accounts",
		}))
		utils.AssertContains(t, output, "Users")
		utils.AssertContains(t, output, "Manage user accounts")
		utils.AssertContains(t, output, "<p")
	})

	t.Run("action slot renders on right", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PageHeader(PageHeaderProps{
			Title:  "Tenants",
			Action: templ.Raw(`<a href="/tenants/new">New tenant</a>`),
		}))
		utils.AssertContains(t, output, "Tenants")
		utils.AssertContains(t, output, "New tenant")
		utils.AssertContains(t, output, "flex-shrink-0")
	})

	t.Run("breadcrumb slot renders above title", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PageHeader(PageHeaderProps{
			Title:      "User detail",
			Breadcrumb: templ.Raw(`<nav>Home / Users</nav>`),
		}))
		utils.AssertContains(t, output, "User detail")
		utils.AssertContains(t, output, "<nav>Home / Users</nav>")
		utils.AssertContains(t, output, "mb-3")
	})

	t.Run("no subtitle hides paragraph", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PageHeader(PageHeaderProps{
			Title: "Settings",
		}))
		utils.AssertNotContains(t, output, "text-gray-500")
	})

	t.Run("no action hides action slot", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PageHeader(PageHeaderProps{
			Title: "Settings",
		}))
		utils.AssertNotContains(t, output, "flex-shrink-0")
	})

	t.Run("propagates BaseProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PageHeader(PageHeaderProps{
			Title: "Audit",
			BaseProps: utils.BaseProps{
				ID:        "page-header-test",
				Class:     "custom-header-class",
				AriaLabel: "Page header",
			},
		}))
		utils.AssertContains(t, output, `id="page-header-test"`)
		utils.AssertContains(t, output, "custom-header-class")
		utils.AssertContains(t, output, `aria-label="Page header"`)
	})
}
