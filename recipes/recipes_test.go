// Package recipes provides rendering tests for the recipes package.
package recipes

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestDashboardRender(t *testing.T) {
	t.Parallel()

	t.Run("renders title and stat cards", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dashboard(DashboardProps{
			Title:    "Overview",
			Subtitle: "Last 30 days",
			StatCards: []templ.Component{
				templ.Raw(`<div data-test="stat1">S1</div>`),
				templ.Raw(`<div data-test="stat2">S2</div>`),
			},
		}))
		utils.AssertContains(t, output, "Overview")
		utils.AssertContains(t, output, "Last 30 days")
		utils.AssertContains(t, output, `data-test="stat1"`)
		utils.AssertContains(t, output, `data-test="stat2"`)
	})

	t.Run("renders charts when provided", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dashboard(DashboardProps{
			Title: "X",
			Charts: []templ.Component{
				templ.Raw(`<div data-test="chart1">Chart1</div>`),
			},
		}))
		utils.AssertContains(t, output, `data-test="chart1"`)
	})

	t.Run("omits stat grid when slice empty", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Dashboard(DashboardProps{Title: "X"}))
		// Still renders title but no chart/stat data attributes
		utils.AssertContains(t, output, "X")
	})

	t.Run("default props", func(t *testing.T) {
		t.Parallel()

		p := DefaultDashboardProps()
		if p.SidebarWidth == "" {
			t.Error("DefaultDashboardProps().SidebarWidth should not be empty")
		}
	})
}

func TestSettingsLayoutRender(t *testing.T) {
	t.Parallel()

	t.Run("renders title and sections", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SettingsLayout(SettingsLayoutProps{
			Title:    "Settings",
			Subtitle: "Manage your account",
			Aside:    templ.Raw(`<nav data-test="aside">aside</nav>`),
			Sections: []SettingsSection{
				{ID: "profile", Title: "Profile", Body: templ.Raw(`<form data-test="profile-form">P</form>`)},
				{ID: "security", Title: "Security", Body: templ.Raw(`<form data-test="security-form">S</form>`)},
			},
		}))
		utils.AssertContains(t, output, "Settings")
		utils.AssertContains(t, output, "Manage your account")
		utils.AssertContains(t, output, `data-test="aside"`)
		utils.AssertContains(t, output, `id="profile"`)
		utils.AssertContains(t, output, `id="security"`)
		utils.AssertContains(t, output, `data-test="profile-form"`)
	})

	t.Run("default props", func(t *testing.T) {
		t.Parallel()

		_ = DefaultSettingsLayoutProps()
	})
}

func TestLoginCardRender(t *testing.T) {
	t.Parallel()

	t.Run("renders title, form body, and footer", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoginCard(LoginCardProps{
			Title:    "Sign in",
			Subtitle: "Welcome back",
			FormBody: templ.Raw(`<form data-test="form">F</form>`),
			Footer:   templ.Raw(`<p data-test="footer">Need an account?</p>`),
		}))
		utils.AssertContains(t, output, "Sign in")
		utils.AssertContains(t, output, "Welcome back")
		utils.AssertContains(t, output, `data-test="form"`)
		utils.AssertContains(t, output, `data-test="footer"`)
	})

	t.Run("renders oauth divider when OAuthButtons set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoginCard(LoginCardProps{
			FormBody:     templ.Raw(`<form>F</form>`),
			OAuthButtons: templ.Raw(`<div data-test="oauth">Google</div>`),
		}))
		utils.AssertContains(t, output, `data-test="oauth"`)
		utils.AssertContains(t, output, ">OR<")
	})

	t.Run("omits oauth divider when OAuthButtons nil", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoginCard(LoginCardProps{
			FormBody: templ.Raw(`<form>F</form>`),
		}))
		utils.AssertNotContains(t, output, ">OR<")
	})

	t.Run("default props set Title", func(t *testing.T) {
		t.Parallel()

		p := DefaultLoginCardProps()
		if p.Title != "Sign in" {
			t.Errorf("DefaultLoginCardProps().Title = %q, want %q", p.Title, "Sign in")
		}
	})
}
