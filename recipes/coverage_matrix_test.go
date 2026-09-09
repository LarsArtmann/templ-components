package recipes

import (
	"context"
	"errors"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// errWriter always fails, exercising the generated templates' write-error
// propagation branches (every WriteString error check in *_templ.go).
type errWriter struct{}

func (errWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }

// TestRenderErrorPropagation proves every recipe propagates writer/context
// failures instead of silently rendering partial output — and covers the
// runtime error branches the happy-path tests cannot reach.
func TestRenderErrorPropagation(t *testing.T) {
	t.Parallel()

	components := map[string]templ.Component{
		"dashboard": Dashboard(DashboardProps{Title: "X", Sidebar: templ.Raw("<nav>")}),
		"login":     LoginCard(LoginCardProps{Title: "X", FormBody: templ.Raw("<form>")}),
		"auth":      AuthLayout(AuthLayoutProps{Card: templ.Raw("<form>")}),
		"settings": SettingsLayout(SettingsLayoutProps{
			Title:    "X",
			Sections: []SettingsSection{{ID: "s", Title: "S", Body: templ.Raw("<form>")}},
		}),
	}

	for name, component := range components {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if err := component.Render(context.Background(), errWriter{}); err == nil {
				t.Errorf("%s rendered without error through a failing writer", name)
			}

			cancelled, cancel := context.WithCancel(context.Background())
			cancel()

			if err := component.Render(cancelled, errWriter{}); err == nil {
				t.Errorf("%s rendered without error on a cancelled context", name)
			}
		})
	}
}

// TestCoverageMatrix renders maximal and minimal prop combinations for every
// recipe so branch-heavy generated templates (slot presence checks, defaults)
// are exercised beyond the happy-path rendering tests. Assertions stay light:
// the point is branch execution, verified by a marker string per slot.
func TestCoverageMatrix(t *testing.T) {
	t.Parallel()

	marker := templ.Raw(`<div data-test="slot">slot</div>`)

	t.Run("dashboard maximal", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Dashboard(DashboardProps{
			BaseProps: utils.BaseProps{
				ID:        "dash-id",
				Class:     "dash-class",
				AriaLabel: "dash-label",
				Attrs:     templ.Attributes{"data-test-root": "dashboard"},
			},
			Title:         "Overview",
			Subtitle:      "Last 30 days",
			Breadcrumb:    marker,
			Sidebar:       marker,
			MobileNav:     marker,
			Header:        marker,
			HeaderActions: marker,
			StatCards:     []templ.Component{marker},
			Charts:        []templ.Component{marker},
		}))
		utils.AssertContainsAll(t, output,
			`data-test="slot"`, `data-test-root="dashboard"`, "dash-id", "dash-class", "dash-label")
	})

	t.Run("dashboard sidebar without header", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Dashboard(DashboardProps{
			Title:   "Shell",
			Sidebar: marker,
		}))
		utils.AssertContains(t, output, `data-test="slot"`)
	})

	t.Run("dashboard shell only", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Dashboard(DashboardProps{Title: "Bare"}))
		utils.AssertContains(t, output, "Bare")
	})

	t.Run("login card full", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, LoginCard(LoginCardProps{
			Title:        "Create account",
			Subtitle:     "Start free",
			FormBody:     marker,
			OAuthButtons: marker,
			Footer:       marker,
		}))
		utils.AssertContains(t, output, "Create account")
	})

	t.Run("auth layout panel footer", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, AuthLayout(AuthLayoutProps{
			BaseProps: utils.BaseProps{
				ID:    "auth-id",
				Class: "auth-class",
				Attrs: templ.Attributes{"data-test-root": "auth"},
			},
			Card:          marker,
			PanelTitle:    "Acme",
			PanelText:     "Build better products.",
			PanelFeatures: []string{"SOC 2", "Uptime"},
			PanelFooter:   marker,
			Reverse:       true,
		}))
		utils.AssertContainsAll(t, output, `data-test-root="auth"`, "auth-id", "auth-class", "Acme")
	})

	t.Run("auth layout card only", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, AuthLayout(AuthLayoutProps{Card: marker}))
		utils.AssertContains(t, output, `data-test="slot"`)
	})

	t.Run("settings sections with subtitles", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SettingsLayout(SettingsLayoutProps{
			Title:    "User Settings",
			Subtitle: "Manage your account",
			Aside:    marker,
			Sections: []SettingsSection{
				{ID: "profile", Title: "Profile", Subtitle: "Your details", Body: marker},
				{ID: "security", Title: "Security", Body: marker},
			},
		}))
		utils.AssertContainsAll(t, output, "Your details", `id="security"`)
	})

	t.Run("settings single section no aside", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SettingsLayout(SettingsLayoutProps{
			Title:    "Solo",
			Sections: []SettingsSection{{ID: "only", Title: "Only", Body: marker}},
		}))
		utils.AssertContains(t, output, `id="only"`)
	})
}
