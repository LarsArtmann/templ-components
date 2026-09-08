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
	// the grid should not constrain it. CAUTION with SidebarNav specifically:
	// SidebarNav hardcodes `w-64` (16rem). Under SidebarWidthAuto the grid
	// track sizes to that element, which is fine — but under SidebarWidthSM
	// (12rem) the fixed `w-64` sidebar OVERFLOWS its 12rem track (the demo hit
	// this: it now uses the MD default). Either widen SidebarWidth, override
	// SidebarNav's width via its Class, or use Auto and let `w-64` win.
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

// shellClassFor picks the AppShell outer wrapper class. The two-track grid
// template is only correct when a sidebar slot exists; with Sidebar nil the
// single content column would land in the sidebar-width track and collapse
// the shell, so the no-sidebar class is used instead.
func shellClassFor(sidebar templ.Component) string {
	if sidebar == nil {
		return appshellShellNoSidebarClass
	}

	return appshellShellClass
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
//
// # Empty-slot contract
//
// Every slot is optional; nil renders predictably:
//
//	Sidebar    nil → no sidebar wrapper, NO --tc-sidebar-w style, single-column
//	             shell (the two-track grid is only emitted when a sidebar
//	             exists — emitting it anyway collapses the shell; fixed
//	             2026-09-08 after the dashboard recipe rendered ~16rem wide)
//	Header     nil → no <header> element
//	MobileNav  nil → nothing rendered below lg
//	Footer     nil → no <footer> element
//	Content    nil → empty content column (the shell frame still renders)
//
// # Embedded contexts and min-h-dvh
//
// The shell stretches to the full viewport height (`min-h-dvh`) because it
// is designed as the direct child of `layout.Base`'s <main>. When embedding
// an AppShell inside a card or other non-page container (like the demo
// does), override the height via Class: `utils.Class` merges consumer
// classes last, so `Class: "min-h-0"` (or an explicit height) wins over the
// shell default.
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
//nolint:exhaustruct_v5 // constructor intentionally sets only non-zero defaults
func DefaultAppShellProps() AppShellProps {
	return AppShellProps{
		SidebarWidth:   SidebarWidthDefault,
		StickyHeader:   true,
		Container:      true,
		ContainerWidth: ContainerWidthDefault,
	}
}
