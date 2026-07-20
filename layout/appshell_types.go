package layout

import (
	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// SidebarWidth is a typed enum for the desktop sidebar column width of an
// AppShell. Unknown values fall back to SidebarWidthMD.
type SidebarWidth string

const (
	// SidebarWidthSM is 12rem (w-48) — compact icon-heavy sidebars.
	SidebarWidthSM SidebarWidth = "sm"
	// SidebarWidthMD is 16rem (w-64) — the default, matches SidebarNav.
	SidebarWidthMD SidebarWidth = "md"
	// SidebarWidthLG is 20rem (w-80) — wide sidebars with labels + descriptions.
	SidebarWidthLG SidebarWidth = "lg"
	// SidebarWidthAuto sizes the sidebar to its content via `auto` in the grid
	// template. Use when the sidebar (e.g. SidebarNav) sets its own width and
	// the grid should not constrain it.
	SidebarWidthAuto SidebarWidth = "auto"
	// SidebarWidthDefault is the canonical default (MD).
	SidebarWidthDefault SidebarWidth = SidebarWidthMD
)

// SidebarWidthIsValid reports whether w is a recognized SidebarWidth.
func SidebarWidthIsValid(w SidebarWidth) bool {
	_, ok := sidebarWidthLookup[w]

	return ok
}

// sidebarWidthMD is the default sidebar width in CSS units. Extracted as a
// named constant to satisfy goconst (this value appears in the lookup table,
// the default fallback, and the godoc examples).
const sidebarWidthMDValue = "16rem"

//nolint:gochecknoglobals // Package-level lookup table for sidebar widths
var sidebarWidthLookup = map[SidebarWidth]string{
	SidebarWidthSM:   "12rem",
	SidebarWidthMD:   sidebarWidthMDValue,
	SidebarWidthLG:   "20rem",
	SidebarWidthAuto: "auto",
}

// sidebarWidthValue returns the CSS grid-template-columns value (a length or
// "auto") for the sidebar column. Falls back to SidebarWidthDefault.
func sidebarWidthValue(w SidebarWidth) string {
	return utils.Lookup(sidebarWidthLookup, w, sidebarWidthLookup[SidebarWidthDefault])
}

// AppShellProps configures a sidebar + content application shell. This is the
// 2D layout primitive every admin dashboard rebuilds by hand: a fixed-width
// sidebar on desktop, a sticky optional header, and a flexible main column.
//
// AppShell renders INSIDE the `<main>` landmark provided by `layout.Base` —
// it does NOT emit its own `<main>` or skip-link (Base owns those). Use it
// as the body content of Base:
//
//	@layout.Base(props) {
//	   @layout.AppShell(layout.AppShellProps{
//	      Sidebar:  navigation.SidebarNav(...),
//	      Header:   navigation.Nav(...),
//	      Content:  dashboardContent(),
//	   })
//	}
//
// # Mobile
//
// The desktop sidebar is `hidden lg:block` — invisible below the lg
// breakpoint. For mobile navigation, render a separate `display.Drawer` (or
// any other mobile pattern) and pass it as the MobileNav slot. AppShell will
// render that slot only below lg. This keeps layout free of the display
// import and gives the consumer full control over mobile UX.
type AppShellProps struct {
	utils.BaseProps

	// Sidebar renders in the left grid column on desktop (lg+). Hidden on
	// mobile. Typically navigation.SidebarNav, but any templ.Component works.
	Sidebar templ.Component
	// MobileNav renders below lg only (lg:hidden). Optional. Use it to host a
	// display.Drawer or other mobile navigation pattern. Nil = no mobile nav.
	MobileNav templ.Component
	// Header renders inside the content column, sticky to the top when
	// StickyHeader is true. Typically navigation.Nav. Optional.
	Header templ.Component
	// Footer renders inside the content column, after Content. Optional.
	Footer templ.Component
	// Content is the main body of the shell. Required.
	Content templ.Component
	// SidebarWidth controls the desktop sidebar column width via a CSS custom
	// property (--tc-sidebar-w). Default SidebarWidthMD (16rem). The main
	// column is always minmax(0, 1fr) — never bare 1fr (grid-blowout guard).
	SidebarWidth SidebarWidth
	// StickyHeader pins the Header to the top of the content column on scroll.
	// Default true. Set to false for static headers.
	StickyHeader bool
	// Container, when true (default), wraps Content in a layout.Container at
	// ContainerWidth. Set to false when Content manages its own width
	// (e.g. an edge-to-edge table or map).
	Container bool
	// ContainerWidth is the max-width applied to the Content wrapper when
	// Container is true. Default ContainerWidthLG.
	ContainerWidth ContainerWidth
}

// DefaultAppShellProps returns sensible defaults: MD sidebar, sticky header,
// Content wrapped in a Container at LG width.
//
//nolint:exhaustruct // constructor intentionally sets only non-zero defaults
func DefaultAppShellProps() AppShellProps {
	return AppShellProps{
		SidebarWidth:   SidebarWidthDefault,
		StickyHeader:   true,
		Container:      true,
		ContainerWidth: ContainerWidthDefault,
	}
}
