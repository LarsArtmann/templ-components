package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/icons"
)

// demoPageMeta describes one page of the multi-page demo. The registry is
// the single source of truth for the routes in newMux, the sidebar items,
// the home-page cards, and the prerender page list — adding a page means
// adding ONE entry here.
type demoPageMeta struct {
	// Path is the mux-relative route ("/display"). The canonical URL is
	// demoURL(Path); sidebar links and CurrentPath use the canonical form.
	Path string
	// Title is the sidebar label, card title, and page <h1>.
	Title string
	// Short is the home-card description.
	Short string
	// Icon renders in the sidebar and on the home card.
	Icon icons.Name
	// Section groups the sidebar item (empty = ungrouped, always visible).
	Section string
	// ShowHeader renders the PageHeader block above the content. Home is the
	// only page that sets it false (its hero replaces the header).
	ShowHeader bool
	// Content renders the page body. Request-scoped pages (wire transport
	// query, users sort/pagination query) read their inputs from r.
	Content func(r *http.Request) templ.Component
}

// canonicalPath is the URL visitors see: base path + route.
func (m demoPageMeta) canonicalPath() string {
	return demoURL(m.Path)
}

// homePageMeta describes the landing page (not part of the sidebar/cards).
var homePageMeta = demoPageMeta{
	Path:    "/",
	Title:   "templ-components Demo",
	Short:   "Showcase of all templ-components",
	Icon:    icons.Home,
	Section: "",
}

// demoPages lists every section page in sidebar order. Home ("/") lives in
// homePageMeta; the standalone recipe screens (dashboard, settings, login,
// auth), the users page, and the /errors/* routes are linked FROM pages
// (recipes, users) and keep their standalone shells.
func demoPages() []demoPageMeta {
	return []demoPageMeta{
		{
			Path:       "/layout",
			Title:      "Layout",
			Short:      "Page shell primitives: Base, Container, AppShell, Stack, theme.",
			Icon:       icons.Squares2x2,
			Section:    "Components",
			ShowHeader: true,
			Content:    func(*http.Request) templ.Component { return layoutDemo() },
		},
		{
			Path:       "/display",
			Title:      "Display",
			Short:      "Cards, tables, badges, avatars, charts, kanban, carousel.",
			Icon:       icons.Cube,
			Section:    "Components",
			ShowHeader: true,
			Content:    func(*http.Request) templ.Component { return displayDemo() },
		},
		{
			Path:       "/feedback",
			Title:      "Feedback",
			Short:      "Alerts, toasts, spinners, skeletons, progress, steps.",
			Icon:       icons.Bell,
			Section:    "Components",
			ShowHeader: true,
			Content:    func(*http.Request) templ.Component { return feedbackDemo() },
		},
		{
			Path:       "/forms",
			Title:      "Forms",
			Short:      "Inputs, selects, toggles, validation, filter bar, calendar.",
			Icon:       icons.Edit,
			Section:    "Components",
			ShowHeader: true,
			Content:    func(*http.Request) templ.Component { return formsDemoContent() },
		},
		{
			Path:       "/navigation",
			Title:      "Navigation",
			Short:      "Nav bars, breadcrumbs, pagination, load-more, footer.",
			Icon:       icons.MapPin,
			Section:    "Components",
			ShowHeader: true,
			Content:    func(*http.Request) templ.Component { return navigationDemo() },
		},
		{
			Path:       "/icons",
			Title:      "Icons",
			Short:      "102 SVG icons with hover animations, zero JavaScript.",
			Icon:       icons.Star,
			Section:    "Components",
			ShowHeader: true,
			Content:    func(*http.Request) templ.Component { return iconsDemo() },
		},
		{
			Path:       "/htmx",
			Title:      "HTMX",
			Short:      "Loading states, error handling, polling, OOB swaps.",
			Icon:       icons.ArrowPath,
			Section:    "Interactive",
			ShowHeader: true,
			Content:    func(*http.Request) templ.Component { return htmxDemo() },
		},
		{
			Path:       "/datastar",
			Title:      "Datastar",
			Short:      "SSE LiveRegion, loading indicators, SDK script.",
			Icon:       icons.Bolt,
			Section:    "Interactive",
			ShowHeader: true,
			Content:    func(*http.Request) templ.Component { return datastarDemo() },
		},
		{
			Path:       "/wire",
			Title:      "Wire (htmx + Datastar)",
			Short:      "One wiring spec, two transports — forms, filters, wizard.",
			Icon:       icons.Link,
			Section:    "Interactive",
			ShowHeader: true,
			Content: func(r *http.Request) templ.Component {
				return wireDemo(parseDemoTransport(r.URL.Query().Get("transport")))
			},
		},
		{
			Path:       "/kanban",
			Title:      "Kanban",
			Short:      "Drag-and-drop board with optimistic moves over both transports.",
			Icon:       icons.QueueList,
			Section:    "Interactive",
			ShowHeader: true,
			Content: func(r *http.Request) templ.Component {
				return kanbanDemo(demoSessionCSRF(r))
			},
		},
		{
			Path:       "/echarts",
			Title:      "ECharts adapter",
			Short:      "Opt-in interactive charts via go-echarts + SDK script.",
			Icon:       icons.Chart,
			Section:    "Interactive",
			ShowHeader: true,
			Content:    func(*http.Request) templ.Component { return echartsDemo() },
		},
		{
			Path:       "/recipes",
			Title:      "Recipes & full pages",
			Short:      "Screen-level compositions: dashboard, settings, login, users.",
			Icon:       icons.Book,
			Section:    "Full pages",
			ShowHeader: true,
			Content:    func(*http.Request) templ.Component { return recipesDemo() },
		},
		{
			Path:       "/users",
			Title:      "Users list",
			Short:      "Server-driven DataTable: sort + paginate round-trips.",
			Icon:       icons.Users,
			Section:    "Full pages",
			ShowHeader: true,
			Content:    usersDemoContent,
		},
		{
			Path:       "/error-pages",
			Title:      "Error pages",
			Short:      "Family-aware ErrorPage, NotFound404, live playground.",
			Icon:       icons.ExclamationTriangle,
			Section:    "Full pages",
			ShowHeader: true,
			Content:    func(*http.Request) templ.Component { return errorpageDemo() },
		},
	}
}

// demoPageByPath finds a registry entry by mux-relative path.
func demoPageByPath(path string) (demoPageMeta, bool) {
	for _, page := range demoPages() {
		if page.Path == path {
			return page, true
		}
	}

	return demoPageMeta{}, false
}
