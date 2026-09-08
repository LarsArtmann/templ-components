package visualtest_test

import (
	"context"
	"io"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/navigation"
	"github.com/larsartmann/templ-components/visualtest"
)

// The responsive sweeps close #159 (mobile 375px) and #160 (RTL): components
// whose BREAKPOINT behavior and MIRRORING are the risky part — Nav's
// hamburger collapse, MobileMenu, form stacking, table overflow, Split's
// single-column collapse, and the logical-property mirroring of Nav, Split,
// Carousel, Drawer, and an open Dropdown.

func sweepNavProps() navigation.NavProps {
	props := navigation.DefaultNavProps()
	props.Sticky = false
	props.Brand = templ.Raw(`<span class="text-base font-bold text-gray-900 dark:text-white">MyApp</span>`)
	props.Links = []navigation.NavLinkProps{
		{Href: "/", Text: "Home"},
		{Href: "/reports", Text: "Reports"},
		{Href: "/settings", Text: "Settings"},
	}
	props.CurrentPath = "/"

	return props
}

func sweepTableProps() display.TableProps {
	props := display.DefaultTableProps()
	props.Headers = []string{"ID", "Customer", "Plan", "MRR", "Status"}
	rows := [][]string{
		{"1", "Acme Corporation", "Enterprise", "$1,200", "Active"},
		{"2", "Globex", "Pro", "$99", "Trial"},
		{"3", "Initech", "Starter", "$29", "Active"},
	}
	for _, cells := range rows {
		row := display.TableRow{}
		for _, text := range cells {
			row.Cells = append(row.Cells, display.TableCell{Text: text})
		}

		props.Rows = append(props.Rows, row)
	}

	return props
}

func sweepSplitProps() layout.SplitProps {
	return layout.SplitProps{
		Main: templ.Raw(
			`<article class="prose max-w-none"><h2 class="text-lg font-semibold text-gray-900 dark:text-white">Article title</h2>` +
				`<p class="mt-2 text-sm text-gray-600 dark:text-gray-400">Main content column with a paragraph of body text that wraps naturally at narrow widths.</p></article>`,
		),
		Aside: templ.Raw(
			`<nav class="rounded-lg border border-gray-200 dark:border-gray-700 p-3 text-sm text-gray-600 dark:text-gray-400">` +
				`<p class="font-medium text-gray-900 dark:text-white">On this page</p>` +
				`<ul class="mt-2 space-y-1"><li>Section A</li><li>Section B</li></ul></nav>`,
		),
	}
}

// TestMobileSweep captures the 375px (iPhone SE) behavior of the components
// whose desktop rendering hides mobile-specific risk.
func TestMobileSweep(t *testing.T) {
	t.Parallel()

	mobile := visualtest.Options{Viewport: visualtest.ViewportMobile}

	t.Run("nav collapses to hamburger", func(t *testing.T) {
		t.Parallel()
		visualtest.AssertScreenshot(t, "nav/mobile_light", navigation.Nav(sweepNavProps()), mobile)
	})

	t.Run("mobile menu renders closed", func(t *testing.T) {
		t.Parallel()

		links := []navigation.NavLinkProps{
			{Href: "/", Text: "Home"},
			{Href: "/reports", Text: "Reports"},
		}

		menu := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			if err := navigation.MobileMenuToggle(false, "sweep-menu", false).Render(ctx, w); err != nil {
				return err
			}

			return navigation.MobileMenu(links, "/", "test-nonce", "sweep-menu", false).Render(ctx, w)
		})
		visualtest.AssertScreenshot(t, "mobilemenu/closed_light", menu, mobile)
	})

	t.Run("stacked form fields", func(t *testing.T) {
		t.Parallel()

		page := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			return forms.Form(forms.FormProps{Action: "/save", Method: forms.FormPost}).
				Render(templ.WithChildren(ctx, templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
					email := forms.DefaultInputProps()
					email.Name = "email"
					email.Label = "Email"
					if err := forms.Input(email).Render(ctx, w); err != nil {
						return err
					}

					return nil
				})), w)
		})
		visualtest.AssertScreenshot(t, "form/stack_mobile_light", page, mobile)
	})

	t.Run("table overflows horizontally", func(t *testing.T) {
		t.Parallel()
		visualtest.AssertScreenshot(t, "table/mobile_light", display.Table(sweepTableProps()), mobile)
	})

	t.Run("split stacks to one column", func(t *testing.T) {
		t.Parallel()
		visualtest.AssertScreenshot(t, "split/stacked_mobile_light", layout.Split(sweepSplitProps()), mobile)
	})
}

// TestRTLSweep captures the dir="rtl" mirroring of components whose layouts
// depend on logical properties actually flipping (audit f13). The goldens are
// the machine-checkable half; a human eyeball pass remains TODO #80 family.
func TestRTLSweep(t *testing.T) {
	t.Parallel()

	rtl := visualtest.Options{RTL: visualtest.Bool(true)}

	t.Run("nav", func(t *testing.T) {
		t.Parallel()
		visualtest.AssertScreenshot(t, "nav/rtl_light", navigation.Nav(sweepNavProps()), rtl)
	})

	t.Run("split", func(t *testing.T) {
		t.Parallel()
		visualtest.AssertScreenshot(t, "split/rtl_light", layout.Split(sweepSplitProps()),
			visualtest.Options{RTL: visualtest.Bool(true), Viewport: visualtest.Viewport{Width: 800, Height: 400}},
		)
	})

	t.Run("carousel", func(t *testing.T) {
		t.Parallel()

		carousel := display.DefaultCarouselProps()
		carousel.Slides = []display.CarouselSlide{
			{Content: templ.Raw(`<div class="flex h-28 items-center justify-center rounded-xl bg-blue-600 text-xl font-bold text-white">1</div>`)},
			{Content: templ.Raw(`<div class="flex h-28 items-center justify-center rounded-xl bg-emerald-600 text-xl font-bold text-white">2</div>`)},
		}
		carousel.ShowArrows = true
		carousel.ShowIndicators = true
		carousel.Nonce = "test-nonce"
		visualtest.AssertScreenshot(t, "carousel/rtl_light", display.Carousel(carousel),
			visualtest.Options{RTL: visualtest.Bool(true), Viewport: visualtest.Viewport{Width: 480, Height: 280}},
		)
	})

	t.Run("drawer right", func(t *testing.T) {
		t.Parallel()

		drawer := display.DefaultDrawerProps()
		drawer.Title = "Filters"
		drawer.Open = true
		drawer.Side = display.DrawerRight
		visualtest.AssertScreenshot(t, "drawer/right_rtl", display.Drawer(drawer),
			visualtest.Options{
				RTL:          visualtest.Bool(true),
				WaitSelector: "dialog",
				FullViewport: true,
				Viewport:     visualtest.Viewport{Width: 480, Height: 400},
				MaxMismatch:  0.01,
			},
		)
	})

	t.Run("dropdown open", func(t *testing.T) {
		t.Parallel()

		dropdown := display.DefaultDropdownProps()
		dropdown.Label = "Options"
		dropdown.Nonce = "test-nonce"
		dropdown.Items = []display.DropdownItem{
			{Text: "Edit", Href: "/edit"},
			{Text: "Archive", Href: "/archive"},
		}
		visualtest.AssertScreenshot(t, "dropdown/open_rtl", display.Dropdown(dropdown),
			visualtest.Options{
				RTL:          visualtest.Bool(true),
				State:        visualtest.StateClick,
				WaitSelector: "[popover]",
				FullViewport: true,
				Viewport:     visualtest.Viewport{Width: 480, Height: 360},
				MaxMismatch:  0.01,
			},
		)
	})
}
