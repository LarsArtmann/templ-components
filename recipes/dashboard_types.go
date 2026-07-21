package recipes

import (
	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/utils"
)

// DashboardProps configures a recipes.Dashboard screen. A dashboard is an
// AppShell with a stat-card grid and arbitrary chart/content slots. Every
// slot is a templ.Component so the consumer controls the actual content —
// recipes.Dashboard provides the layout scaffold only.
//
// Pass slots via templ.Component (Card.Body, navigation.Nav, display.Grid, etc.).
// The Sidebar slot is typically navigation.SidebarNav; the Header slot is
// typically navigation.Nav or navigation.SimpleNav. StatCards is a slice of
// display.StatCard components rendered inside a display.Grid. Charts is a
// slice of arbitrary chart components (display.Card-wrapped) rendered below
// the stat grid.
type DashboardProps struct {
	utils.BaseProps

	// Title and Subtitle render in a display.PageHeader at the top of the
	// content column.
	Title    string
	Subtitle string
	// Breadcrumb is an optional templ.Component slot (typically
	// navigation.Breadcrumbs) rendered above the title in the PageHeader.
	Breadcrumb templ.Component
	// Sidebar is the desktop sidebar slot (typically navigation.SidebarNav).
	// Required for the dashboard layout; pass nil for a sidebar-less shell.
	Sidebar templ.Component
	// MobileNav is the below-lg mobile navigation slot (typically
	// display.Drawer). Optional.
	MobileNav templ.Component
	// Header is the sticky top header slot (typically navigation.Nav). Optional.
	Header templ.Component
	// HeaderActions is rendered on the right side of the PageHeader (above
	// the stat grid). Use it for "Add" buttons, filters, or time-range
	// pickers. Optional.
	HeaderActions templ.Component
	// StatCards is the slice of display.StatCard components rendered in a
	// responsive 4-up grid at the top of the dashboard body.
	StatCards []templ.Component
	// Charts is the slice of chart/visualization components rendered below
	// the stat grid in a responsive 2-up grid. Each chart should be wrapped
	// in a display.Card by the consumer.
	Charts []templ.Component
	// SidebarWidth controls the AppShell sidebar width. Default
	// layout.SidebarWidthMD (16rem).
	SidebarWidth layout.SidebarWidth
}

// DefaultDashboardProps returns sensible defaults (MD sidebar, no slots).
//
//nolint:exhaustruct // constructor intentionally sets only non-zero defaults
func DefaultDashboardProps() DashboardProps {
	return DashboardProps{
		SidebarWidth: layout.SidebarWidthDefault,
	}
}
