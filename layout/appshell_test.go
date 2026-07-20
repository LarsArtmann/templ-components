package layout

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestAppShell(t *testing.T) {
	t.Parallel()

	t.Run("default props produce grid + minmax + dvh", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, AppShell(AppShellProps{
			Content: templ.Raw(`<p>body</p>`),
		}))
		utils.AssertContainsAll(t, output,
			"lg:grid",
			"lg:grid-cols-[var(--tc-sidebar-w)_minmax(0,1fr)]",
			"min-h-dvh",
		)
	})

	t.Run("sidebar width CSS var emitted", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			width SidebarWidth
			want  string
		}{
			{SidebarWidthSM, "12rem"},
			{SidebarWidthMD, "16rem"},
			{SidebarWidthLG, "20rem"},
			{SidebarWidthAuto, "auto"},
			{SidebarWidthDefault, "16rem"},
		}
		for _, tc := range cases {
			output := utils.Render(t, AppShell(AppShellProps{
				Content:      templ.Raw(`<p>x</p>`),
				SidebarWidth: tc.width,
			}))

			want := "--tc-sidebar-w: " + tc.want
			if !strings.Contains(output, want) {
				t.Errorf("width %q: output missing %q", tc.width, want)
			}
		}
	})

	t.Run("unknown sidebar width falls back to default (16rem)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, AppShell(AppShellProps{
			Content:      templ.Raw(`<p>x</p>`),
			SidebarWidth: SidebarWidth("bogus"),
		}))
		utils.AssertContains(t, output, "--tc-sidebar-w: 16rem")
	})

	t.Run("sidebar slot renders inside hidden lg:block wrapper", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, AppShell(AppShellProps{
			Sidebar: templ.Raw(`<nav data-test="sidebar">SidebarNav</nav>`),
			Content: templ.Raw(`<p>x</p>`),
		}))
		utils.AssertContainsAll(t, output,
			`class="`+appshellSidebarWrapperClass+`"`,
			`data-test="sidebar"`,
		)
	})

	t.Run("MobileNav slot renders inside lg:hidden wrapper", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, AppShell(AppShellProps{
			MobileNav: templ.Raw(`<button data-test="mobile">Menu</button>`),
			Content:   templ.Raw(`<p>x</p>`),
		}))
		utils.AssertContainsAll(t, output,
			`class="lg:hidden"`,
			`data-test="mobile"`,
		)
	})

	t.Run("Header sticky when StickyHeader=true", func(t *testing.T) {
		t.Parallel()

		props := DefaultAppShellProps()
		props.Header = templ.Raw(`<div data-test="hdr">Header</div>`)
		props.Content = templ.Raw(`<p>x</p>`)
		output := utils.Render(t, AppShell(props))
		utils.AssertContainsAll(t, output,
			"sticky",
			"top-0",
			"z-40",
			"border-b",
			"dark:border-gray-800",
		)
	})

	t.Run("StickyHeader=false omits sticky classes", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, AppShell(AppShellProps{
			Header:       templ.Raw(`<div>Header</div>`),
			Content:      templ.Raw(`<p>x</p>`),
			StickyHeader: false,
		}))
		if strings.Contains(output, "sticky") {
			t.Errorf("StickyHeader=false should omit sticky class")
		}
	})

	t.Run("Content wrapped in Container when Container=true (max-w-7xl + pad)", func(t *testing.T) {
		t.Parallel()

		props := DefaultAppShellProps()
		props.Content = templ.Raw(`<p data-test="body">body</p>`)
		output := utils.Render(t, AppShell(props))
		utils.AssertContainsAll(t, output,
			"max-w-7xl",
			"px-4",
			"sm:px-6",
			"lg:px-8",
			`data-test="body"`,
		)
	})

	t.Run("Container=false leaves Content unwrapped", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, AppShell(AppShellProps{
			Content:   templ.Raw(`<p data-test="raw">raw</p>`),
			Container: false,
		}))
		utils.AssertContains(t, output, `data-test="raw"`)

		if strings.Contains(output, "max-w-7xl") {
			t.Errorf("Container=false should not wrap Content in max-w-7xl")
		}
	})

	t.Run("ContainerWidth propagates (Prose)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, AppShell(AppShellProps{
			Content:        templ.Raw(`<p>x</p>`),
			Container:      true,
			ContainerWidth: ContainerWidthProse,
		}))
		utils.AssertContains(t, output, "max-w-prose")
	})

	t.Run("Footer slot renders in content column with mt-auto", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, AppShell(AppShellProps{
			Content: templ.Raw(`<p>x</p>`),
			Footer:  templ.Raw(`<div data-test="footer">Footer</div>`),
		}))
		utils.AssertContainsAll(t, output,
			`<footer`,
			"mt-auto",
			`data-test="footer"`,
		)
	})

	t.Run("BaseProps propagate (ID, Class, AriaLabel, Attrs)", func(t *testing.T) {
		t.Parallel()

		props := AppShellProps{
			BaseProps: utils.BaseProps{
				ID:        "app",
				Class:     "data-tc-test",
				AriaLabel: "Application",
				Attrs:     templ.Attributes{"data-testid": "shell"},
			},
			Content: templ.Raw(`<p>x</p>`),
		}
		output := utils.Render(t, AppShell(props))
		utils.AssertContainsAll(t, output,
			`id="app"`,
			"data-tc-test",
			`aria-label="Application"`,
			`data-testid="shell"`,
		)
	})

	t.Run("nil Content does not panic and renders empty content area", func(t *testing.T) {
		t.Parallel()
		// Content is the only required-ish slot; nil must not panic.
		output := utils.Render(t, AppShell(DefaultAppShellProps()))
		utils.AssertContainsAll(t, output, "min-h-dvh", "lg:grid", "minmax(0,1fr)")
	})

	t.Run("dark mode: header has dark:border-gray-800 + dark:bg-gray-900", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, AppShell(AppShellProps{
			Header:  templ.Raw(`<div>x</div>`),
			Content: templ.Raw(`<p>x</p>`),
		}))
		utils.AssertContainsAll(t, output,
			"dark:border-gray-800",
			"dark:bg-gray-900",
		)
	})

	t.Run("minmax(0,1fr) present (grid-blowout guard)", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, AppShell(AppShellProps{
			Content: templ.Raw(`<p>x</p>`),
		}))
		if !strings.Contains(output, "minmax(0,1fr)") {
			t.Errorf("AppShell must use minmax(0,1fr) not bare 1fr (ADR-0016)")
		}
	})
}

func TestSidebarWidthIsValid(t *testing.T) {
	t.Parallel()

	for _, w := range []SidebarWidth{SidebarWidthSM, SidebarWidthMD, SidebarWidthLG, SidebarWidthAuto} {
		if !SidebarWidthIsValid(w) {
			t.Errorf("SidebarWidthIsValid(%q) = false; want true", w)
		}
	}

	if SidebarWidthIsValid(SidebarWidth("bogus")) {
		t.Errorf("SidebarWidthIsValid(\"bogus\") = true; want false")
	}
}
