package navigation

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
	"github.com/larsartmann/templ-components/utils/wire"
)

// TestNavLinkWire pins the NavLink-level Wire contract (#155): a wired link
// keeps its href (no-JS fallback) and carries the action's dialect
// attributes; nil Wire or an empty URL stays a plain link; the active link
// renders as a span with no wiring at all; the consumer's action is never
// mutated. The field rides NavLinkProps, so Nav, SimpleNav, MobileMenu, and
// MobileNavLink inherit it through the same anchor.
func TestNavLinkWire(t *testing.T) {
	t.Parallel()

	base := NavLinkProps{Href: "/inbox", Text: "Inbox"}

	t.Run("htmx dialect renders hx attributes alongside href", func(t *testing.T) {
		t.Parallel()

		props := base
		props.Wire = &wire.Action{ //nolint:exhaustruct // URL+Target are the wiring surface
			URL:    "/api/inbox-fragment",
			Target: "#content",
		}

		html := utils.Render(t, NavLink(props, "/"))
		utils.AssertContainsAll(t, html,
			`href="/inbox"`,
			`hx-get="/api/inbox-fragment"`,
			`hx-target="#content"`,
		)
	})

	t.Run("datastar dialect renders click binding alongside href", func(t *testing.T) {
		t.Parallel()

		props := base
		props.Wire = &wire.Action{ //nolint:exhaustruct // URL is the wiring surface
			Transport: wire.TransportDatastar,
			URL:       "/api/inbox-fragment",
		}

		html := utils.Render(t, NavLink(props, "/"))
		utils.AssertContainsAll(t, html,
			`href="/inbox"`,
			`data-on:click="@get(&#39;/api/inbox-fragment&#39;)"`,
		)
	})

	t.Run("nil Wire keeps a plain link", func(t *testing.T) {
		t.Parallel()

		props := base

		html := utils.Render(t, NavLink(props, "/"))
		utils.AssertContains(t, html, `href="/inbox"`)
		utils.AssertNotContains(t, html, "hx-get=")
		utils.AssertNotContains(t, html, "data-on:click=")
	})

	t.Run("empty URL Wire is inert", func(t *testing.T) {
		t.Parallel()

		props := base
		props.Wire = &wire.Action{Target: "#content"} //nolint:exhaustruct // empty URL = no wiring

		html := utils.Render(t, NavLink(props, "/"))
		utils.AssertContains(t, html, `href="/inbox"`)
		utils.AssertNotContains(t, html, "hx-get=")
	})

	t.Run("active link keeps the anchor with wiring for a region refresh", func(t *testing.T) {
		t.Parallel()

		props := base
		props.Wire = &wire.Action{
			URL:    "/api/inbox-fragment",
			Target: "#content",
		} //nolint:exhaustruct // URL+Target are the wiring surface

		html := utils.Render(t, NavLink(props, "/inbox"))
		utils.AssertContainsAll(t, html,
			`aria-current="page"`,
			`href="/inbox"`,
			`hx-get="/api/inbox-fragment"`,
		)
	})

	t.Run("MobileNavLink carries the wiring too", func(t *testing.T) {
		t.Parallel()

		props := base
		props.Wire = &wire.Action{
			URL:    "/api/inbox-fragment",
			Target: "#content",
		} //nolint:exhaustruct // URL+Target are the wiring surface

		html := utils.Render(t, MobileNavLink(props, "/"))
		utils.AssertContainsAll(t, html, `href="/inbox"`, `hx-get="/api/inbox-fragment"`)
	})

	t.Run("consumer action is never mutated", func(t *testing.T) {
		t.Parallel()

		props := base
		action := &wire.Action{URL: "/api/inbox-fragment"} //nolint:exhaustruct // URL is the wiring surface
		props.Wire = action

		_ = utils.Render(t, NavLink(props, "/"))

		if action.URL != "/api/inbox-fragment" {
			t.Errorf("Wire action URL mutated: %q — components render consumer specs, never rewrite them", action.URL)
		}
	})
}

// TestGoldenSweepNavLinkWire pins the wired NavLink markup for both dialects.
func TestGoldenSweepNavLinkWire(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "nav_link_wire_htmx", HTML: utils.Render(t, NavLink(NavLinkProps{
			Href: "/inbox",
			Text: "Inbox",
			Wire: &wire.Action{
				URL:    "/api/inbox-fragment",
				Target: "#content",
			}, //nolint:exhaustruct // URL+Target are the wiring surface
		}, "/"))},
		{Name: "nav_link_wire_datastar", HTML: utils.Render(t, NavLink(NavLinkProps{
			Href: "/inbox",
			Text: "Inbox",
			Wire: &wire.Action{
				Transport: wire.TransportDatastar,
				URL:       "/api/inbox-fragment",
			}, //nolint:exhaustruct // URL is the wiring surface
		}, "/"))},
	})
}
