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
			Title:     "Audit",
			ID:        "page-header-test",
			Class:     "custom-header-class",
			AriaLabel: "Page header",
		}))
		utils.AssertContains(t, output, `id="page-header-test"`)
		utils.AssertContains(t, output, "custom-header-class")
		utils.AssertContains(t, output, `aria-label="Page header"`)
	})

	t.Run("TitleComponent renders inside the h1", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PageHeader(PageHeaderProps{
			Title:          "Ignored string",
			TitleComponent: templ.Raw(`<code>orders.conflict</code>`),
		}))
		utils.AssertContains(t, output, "<h1")
		utils.AssertContains(t, output, "<code>orders.conflict</code>")
		utils.AssertNotContains(t, output, "Ignored string")
	})

	t.Run("SubtitleComponent renders inside the subtitle paragraph", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PageHeader(PageHeaderProps{
			Title:             "Deployments",
			Subtitle:          "Ignored string",
			SubtitleComponent: templ.Raw(`<span class="wp-pill">3 running</span>`),
		}))
		utils.AssertContains(t, output, "<p")
		utils.AssertContains(t, output, `wp-pill">3 running</span>`)
		utils.AssertNotContains(t, output, "Ignored string")
	})

	t.Run("string fields still render without components", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, PageHeader(PageHeaderProps{
			Title:    "Plain",
			Subtitle: "Plain subtitle",
		}))
		utils.AssertContains(t, output, ">Plain</h1>")
		utils.AssertContains(t, output, "Plain subtitle")
	})
}
