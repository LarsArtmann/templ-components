package visualtest

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/navigation"
)

// appShellTestContent is the main column of the AppShell goldens.
func appShellTestContent() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<div class="p-6">
			<h2 class="text-lg font-semibold text-gray-900">Dashboard</h2>
			<p class="mt-1 text-sm text-gray-500">Shell content column</p>
			<div class="mt-4 grid grid-cols-2 gap-4">
				<div class="rounded-lg border border-gray-200 p-4"><p class="text-sm font-medium">Metric A</p><p class="text-2xl font-bold">1,204</p></div>
				<div class="rounded-lg border border-gray-200 p-4"><p class="text-sm font-medium">Metric B</p><p class="text-2xl font-bold">98%</p></div>
			</div>
		</div>`)

		return err
	})
}

// appShellTestHeader is a simple shell header bar.
func appShellTestHeader() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<div class="flex items-center justify-between px-6 py-3"><span class="text-sm font-medium">Header</span><span class="text-xs text-gray-500">v1</span></div>`)

		return err
	})
}

// appShellSidebar is the shell's sidebar slot (SidebarNav-shaped).
func appShellSidebar() templ.Component {
	props := navigation.DefaultSidebarNavProps()
	props.Items = []navigation.SidebarNavItem{
		{Href: "/", Text: "Overview"},
		{Href: "/reports", Text: "Reports"},
		{Href: "/settings", Text: "Settings"},
	}
	props.CurrentPath = "/"

	return navigation.SidebarNav(props)
}

// appShellForTest builds a full shell: SidebarNav sidebar + header + content.
func appShellForTest() templ.Component {
	props := layout.DefaultAppShellProps()
	props.Sidebar = appShellSidebar()
	props.Header = appShellTestHeader()
	props.Content = appShellTestContent()

	return layout.AppShell(props)
}

// TestAppShellFull is AppShell's first visual golden (#164): the shell with a
// SidebarNav sidebar, header, and content at a desktop (lg+) viewport so the
// two-track grid is active. The no-sidebar collapse regression is unit-guarded
// (layout/appshell_test.go); this pins the WITH-sidebar layout at pixel level.
func TestAppShellFull(t *testing.T) {
	t.Parallel()

	visualtest.AssertScreenshot(t, "appshell/light", appShellForTest(),
		visualtest.Options{Viewport: visualtest.Viewport{Width: 1280, Height: 800}},
	)
	visualtest.AssertScreenshot(t, "appshell/dark", appShellForTest(),
		visualtest.Options{
			Dark:     visualtest.Bool(true),
			Viewport: visualtest.Viewport{Width: 1280, Height: 800},
		},
	)
}

// TestAppShellNoSidebar pins the sidebar-less shell at pixel level: a single
// full-width column (the 2026-09-08 collapse bug squeezed it to the sidebar
// track width — unit tests assert the class, this asserts the rendered
// geometry).
func TestAppShellNoSidebar(t *testing.T) {
	t.Parallel()

	props := layout.DefaultAppShellProps()
	props.Header = appShellTestHeader()
	props.Content = appShellTestContent()

	visualtest.AssertScreenshot(t, "appshell/no_sidebar_light", layout.AppShell(props),
		visualtest.Options{Viewport: visualtest.Viewport{Width: 1280, Height: 600}},
	)
}

// TestAppShellSidebarFitsTrack DOM-measures the documented SidebarWidth ×
// SidebarNav interaction (audit f18, replacing pixel estimates): at the MD
// default the SidebarNav (w-64, 16rem) must fit inside its 16rem grid track.
// Overflow here is the "SM track vs w-64 sidebar" bug class — measurable in
// the DOM, so no golden needed.
func TestAppShellSidebarFitsTrack(t *testing.T) {
	t.Parallel()

	page, err := renderHTML(appShellForTest(), defaultOptions(Options{
		Viewport: Viewport{Width: 1280, Height: 800},
	}))
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

	var overflow float64
	err = chromedp.Run(ctx,
		chromedp.EmulateViewport(1280, 800),
		chromedp.Navigate(srv.URL),
		chromedp.WaitVisible("#tc-root", chromedp.ByQuery),
		chromedp.Evaluate(`(() => {
			const wrapper = document.querySelector('.hidden.lg\\:block');
			const aside = wrapper.querySelector('aside') || wrapper.firstElementChild;
			if (!wrapper || !aside) return -1;
			return aside.getBoundingClientRect().width - wrapper.getBoundingClientRect().width;
		})()`, &overflow),
	)
	if err != nil {
		t.Fatalf("chromedp: %v", err)
	}

	if overflow < 0 {
		t.Fatal("could not locate the sidebar wrapper/aside in the rendered shell")
	}

	// Sub-pixel tolerance only: any real overflow (like a w-64 sidebar in a
	// 12rem SM track) is >= tens of pixels.
	if overflow > 1 {
		t.Errorf("sidebar overflows its grid track by %.1fpx — SidebarWidth must be >= the sidebar's own width (SidebarNav is w-64/16rem; the MD default fits exactly)", overflow)
	}
}
