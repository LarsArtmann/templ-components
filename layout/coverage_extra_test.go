package layout

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestBaseExtraCoverage(t *testing.T) {
	t.Parallel()

	t.Run("with all optional fields", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Base(PageProps{
			Title:       "Test Page",
			Description: "A test page",
			Locale:      "en",
			ThemeColor:  "#000000",
			CSSPath:     "/style.css",
		}))
		utils.AssertContains(t, output, "Test Page")
		utils.AssertContains(t, output, "A test page")
		utils.AssertContains(t, output, "/style.css")
		utils.AssertContains(t, output, "#000000")
	})

	t.Run("with dark theme color", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Base(PageProps{
			Title:          "Dark",
			DarkThemeColor: "#1a1a1a",
		}))
		utils.AssertContains(t, output, "#1a1a1a")
	})

	t.Run("with favicon and OG image", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Base(PageProps{
			Title:   "Social",
			Favicon: "/favicon.ico",
			OGImage: "/og.png",
		}))
		utils.AssertContains(t, output, "/favicon.ico")
		utils.AssertContains(t, output, "/og.png")
	})

	t.Run("with nonce", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Base(PageProps{
			Title: "Nonce Test",
			Nonce: "test-nonce",
		}))
		utils.AssertContains(t, output, `nonce="test-nonce"`)
	})

	t.Run("default props", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Base(DefaultPageProps()))
		utils.AssertContains(t, output, "<html")
	})
}

func TestMinimalExtraCoverage(t *testing.T) {
	t.Parallel()

	t.Run("default props", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Minimal(DefaultMinimalProps()))
		utils.AssertContains(t, output, "<html")
	})
}

func TestThemeToggleExtraCoverage(t *testing.T) {
	t.Parallel()

	t.Run("renders with custom aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ThemeToggle("Switch appearance", ""))
		utils.AssertContains(t, output, "Switch appearance")
	})

	t.Run("default aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ThemeToggle("", ""))
		utils.AssertContains(t, output, "Toggle theme")
	})
}
