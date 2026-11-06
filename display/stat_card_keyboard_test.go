package display

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// htmx-driven StatCards are a click target on a <div>: without role, tabindex,
// and key activation a keyboard user cannot reach the tile at all (#36).
func TestStatCardHxGetKeyboardAccess(t *testing.T) {
	t.Parallel()

	t.Run("hx-get tile is a keyboard-operable button", func(t *testing.T) {
		t.Parallel()
		html := utils.Render(t, StatCard(StatCardProps{
			Label:    "Active users",
			Value:    "1,204",
			HxGet:    "/api/stats",
			HxTarget: "#stat-container",
		}))

		utils.AssertContainsAll(t, html,
			`role="button"`,
			`tabindex="0"`,
			"data-tc-hx-activate",
			`hx-get="/api/stats"`,
			"focus-visible:ring-2",
			"cursor-pointer",
		)
	})

	t.Run("activation script is nonce-guarded and precedes the tile", func(t *testing.T) {
		t.Parallel()
		html := utils.Render(t, StatCard(StatCardProps{
			Nonce: "stat-nonce",
			Label: "Active users",
			Value: "1,204",
			HxGet: "/api/stats",
		}))

		utils.AssertContains(t, html, `<script nonce="stat-nonce">`)
		utils.AssertContains(t, html, "tcHxTileActivateAttached")

		if scriptAt, tileAt := indexOfScript(
			html,
		), indexOfTile(
			html,
		); scriptAt == -1 || tileAt == -1 ||
			scriptAt > tileAt {
			t.Errorf("activation script must render before the tile (script=%d tile=%d):\n%s", scriptAt, tileAt, html)
		}
	})

	t.Run("plain card stays inert", func(t *testing.T) {
		t.Parallel()
		html := utils.Render(t, StatCard(StatCardProps{Label: "Active users", Value: "1,204"}))

		utils.AssertNotContains(t, html, `role="button"`)
		utils.AssertNotContains(t, html, "data-tc-hx-activate")
		utils.AssertNotContains(t, html, "tcHxTileActivateAttached")
	})

	t.Run("href card keeps anchor semantics", func(t *testing.T) {
		t.Parallel()
		html := utils.Render(t, StatCard(StatCardProps{
			Label: "Active users",
			Value: "1,204",
			Href:  "/users?activity=active",
		}))

		utils.AssertContains(t, html, `href="/users?activity=active"`)
		utils.AssertNotContains(t, html, `role="button"`)
		utils.AssertNotContains(t, html, "data-tc-hx-activate")
	})
}

func indexOfScript(html string) int {
	return strings.Index(html, "<script")
}

func indexOfTile(html string) int {
	return strings.Index(html, "data-tc-hx-activate")
}
