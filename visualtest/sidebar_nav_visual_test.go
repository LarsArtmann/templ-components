package visualtest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/navigation"
)

// sidebarNavDemoItems is the canonical golden item set: ungrouped top entries
// followed by two section groups (exercises both render branches).
func sidebarNavDemoItems() []navigation.SidebarNavItem {
	return []navigation.SidebarNavItem{
		{Href: "/", Label: "Overview"},
		{Href: "/reports", Label: "Reports"},
		{Href: "/users", Label: "Users", Section: "Workspace"},
		{Href: "/settings", Label: "Settings", Section: "Workspace"},
		{Href: "/docs", Label: "Documentation", Section: "Resources"},
		{Href: "/support", Label: "Support", Section: "Resources"},
	}
}

// sidebarNavForTest builds a component-level SidebarNav (not the
// AppShell-incidental render) with an active first item.
func sidebarNavForTest(items []navigation.SidebarNavItem) templ.Component {
	props := navigation.DefaultSidebarNavProps()
	props.Items = items
	props.CurrentPath = "/"

	return navigation.SidebarNav(props)
}

// TestSidebarNavLightDark pins the theme-adaptive sidebar at component level:
// WHITE chrome in light mode, the classic near-black chrome in dark mode. The
// 2026 sidebar regression (invisible 2+ months) lived exactly between these
// two renders — AppShell captures show it incidentally; these are deliberate.
func TestSidebarNavLightDark(t *testing.T) {
	t.Parallel()

	AssertScreenshot(t, "sidebar_nav/light", sidebarNavForTest(sidebarNavDemoItems()),
		Options{Viewport: ViewportDesktop},
	)
	AssertScreenshot(t, "sidebar_nav/dark", sidebarNavForTest(sidebarNavDemoItems()),
		Options{Dark: new(true), Viewport: ViewportDesktop},
	)
}

// classicDarkSidebarCSS is the documented permanently-dark opt-out — the same
// block as templates/custom.css (the canonical copy consumers are pointed at).
const classicDarkSidebarCSS = `<style>
  :root {
    --tc-sidebar-bg: var(--color-gray-900);
    --tc-sidebar-border: transparent;
    --tc-sidebar-fg: var(--color-gray-300);
    --tc-sidebar-fg-hover: var(--color-white);
    --tc-sidebar-item-hover-bg: var(--color-gray-800);
    --tc-sidebar-muted: var(--color-gray-500);
    --tc-sidebar-muted-hover: var(--color-gray-300);
  }
</style>`

// TestSidebarNavClassicDarkOptOut browser-proves the pure-CSS opt-out: in
// LIGHT mode with the classic-dark token block injected, the sidebar root's
// computed background must be gray-900 (rgb(17, 24, 39)) — the exact classic
// admin chrome. Token wiring is CSS; only a real browser proves the
// custom-property override reaches the rendered element.
func TestSidebarNavClassicDarkOptOut(t *testing.T) {
	t.Parallel()

	optOut := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, classicDarkSidebarCSS); err != nil {
			return err
		}

		return sidebarNavForTest(sidebarNavDemoItems()).Render(ctx, w)
	})

	page, err := renderHTML(optOut, defaultOptions(Options{Viewport: ViewportDesktop}))
	if err != nil {
		t.Fatalf("build page: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, page)
	}))
	t.Cleanup(srv.Close)

	ctx, cancel := newTab(t)
	defer cancel()

	chromedp.Run(ctx, chromedp.Navigate(srv.URL))

	var bg string

	// Palette-independent assertion: the aside's computed background must
	// equal a live bg-gray-900 probe — whatever --color-gray-900 maps to in
	// the active theme (the demo palette remaps gray → stone values, so a
	// hard-coded rgb() would false-fail). The probe guards against BOTH
	// failure modes: override-not-reaching-element (transparent/white aside)
	// and a probe that never got its class (empty comparison).
	chromedp.Run(ctx,
		chromedp.WaitVisible("aside", chromedp.ByQuery),
		chromedp.Evaluate(`(() => {
			const probe = document.createElement('div');
			probe.className = 'bg-gray-900';
			document.body.appendChild(probe);
			const probeBg = getComputedStyle(probe).backgroundColor;
			const asideBg = getComputedStyle(document.querySelector('aside')).backgroundColor;
			probe.remove();
			return probeBg !== 'rgba(0, 0, 0, 0)' && probeBg === asideBg ? probeBg : '';
		})()`, &bg),
	)

	if bg == "" {
		t.Error("opt-out sidebar background does not match a bg-gray-900 probe — the classic-dark token override did not reach the element")
	}

	AssertScreenshot(t, "sidebar_nav/classic_dark_opt_out", optOut,
		Options{Viewport: ViewportDesktop},
	)
}
