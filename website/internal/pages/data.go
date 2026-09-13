// Package pages contains the templ components and typed content data for the
// templ-components marketing site. Every page renders through the library's
// own components (layout.Base, display.Button, icons, ...) — the site is the
// library's largest dogfood consumer.
package pages

import "github.com/larsartmann/templ-components/icons"

// Site identity. Single source of truth for every page (head tags, footer,
// JSON-LD). Mirrors the values the Astro site shipped in src/data/config.ts.
const (
	SiteName        = "templ-components"
	SiteTitle       = "templ-components — Server-Rendered UI Components for Go"
	SiteDescription = "A pure Tailwind CSS v4 component library for Go's templ engine with first-class HTMX integration. No DaisyUI, no Node.js, no framework lock-in."
	SiteURL         = "https://templcomponents.lars.software"
	DemoURL         = "https://templcomponents-demo-132045829579.us-central1.run.app"
	GitHubURL       = "https://github.com/larsartmann/templ-components"
	AuthorName      = "LarsArtmann"
	AuthorURL       = "https://larsartmann.com/"
	PkgGoDevURL     = "https://pkg.go.dev/github.com/larsartmann/templ-components"
)

// Feature is a landing-page feature card.
type Feature struct {
	Icon  icons.Name
	Title string
	Desc  string
}

// Features powers the FeatureGrid section. Content ported from the Astro
// src/data/features.ts (icon keys now map to the real icons module).
var Features = []Feature{
	{
		Icon:  icons.ShieldCheck,
		Title: "Type-Safe Props",
		Desc:  "Typed string enums make invalid states unrepresentable. Every props struct embeds BaseProps for consistent ID, class, ARIA, and CSP nonce propagation.",
	},
	{
		Icon:  icons.Bolt,
		Title: "Zero Node.js",
		Desc:  "Pure Go + templ + Tailwind CSS v4. No build pipeline beyond templ generate. No pnpm, no bundlers, no SPA framework. CSS-first config, class-based dark mode.",
	},
	{
		Icon:  icons.CodeBracket,
		Title: "Server-Rendered",
		Desc:  "Every component renders HTML on the server. HATEOAS-aligned: JavaScript enhances rather than replaces HTML. Interactive features use minimal vanilla JS with CSP nonces.",
	},
	{
		Icon:  icons.Moon,
		Title: "Built-in Dark Mode",
		Desc:  "Every component has dark: variants for all neutral and semantic colors, enforced by regression tests. ThemeScript prevents FOUC. color-scheme for native form controls.",
	},
	{
		Icon:  icons.CheckCircle,
		Title: "CSP-Ready",
		Desc:  "All inline scripts use nonce attributes. No eval(), no inline event handlers. Integration test suite verifies nonce compliance on every inline script across all components.",
	},
	{
		Icon:  icons.Bolt,
		Title: "HTMX Integration",
		Desc:  "Dedicated htmx package with loading indicators, error handling, CSRF protection, out-of-band swaps, and View Transitions. First-class server-rendered HTMX patterns.",
	},
}

// StepColor is the accent color of a HowItWorks step card.
type StepColor string

const (
	StepAccent StepColor = "accent"
	StepAmber  StepColor = "amber"
)

// IsValid reports whether the step color is a known value.
func (c StepColor) IsValid() bool {
	switch c {
	case StepAccent, StepAmber:
		return true
	default:
		return false
	}
}

// Step is one card of the HowItWorks section.
type Step struct {
	Step      string
	StepColor StepColor
	Title     string
	Desc      string
	Code      string
}

// Steps powers the HowItWorks section (ported from src/data/sections.ts).
var Steps = []Step{
	{
		Step:      "1",
		StepColor: StepAccent,
		Title:     "Import",
		Desc:      "Pick only the packages you need. No monolithic bundle — pay for what you use.",
		Code:      `import "github.com/larsartmann/templ-components/display"`,
	},
	{
		Step:      "2",
		StepColor: StepAccent,
		Title:     "Render",
		Desc:      "Pass typed props structs. Every invalid state is a compile-time error, not a runtime check.",
		Code:      `@display.Card(display.CardProps{Title: "Hello"}) { ... }`,
	},
	{
		Step:      "3",
		StepColor: StepAmber,
		Title:     "Generate",
		Desc:      "templ generate compiles .templ files to Go. The generated code is committed for library consumers.",
		Code:      "templ generate ./...",
	},
	{
		Step:      "4",
		StepColor: StepAmber,
		Title:     "Ship",
		Desc:      "Server renders HTML. HTMX enhances where needed. No client-side framework required.",
		Code:      "// HTML over the wire — done",
	},
}

// MatrixCell is one comparison-table cell. Yes/No render as glyphs, Partial as
// a tilde, and any other text renders verbatim (e.g. "Tailwind v4 CSS-first").
type MatrixCell string

const (
	MatrixYes     MatrixCell = "yes"
	MatrixNo      MatrixCell = "no"
	MatrixPartial MatrixCell = "partial"
)

// ComparisonColumn is a competitor/library column of the matrix. The last
// column (templ-components) is highlighted.
var ComparisonColumns = []string{"templUI", "goshipit", SiteName}

// ComparisonMatrixRow is one feature row of the comparison table.
type ComparisonMatrixRow struct {
	Feature string
	Values  []MatrixCell // aligned with ComparisonColumns
}

// ComparisonMatrix powers the comparison table (ported and refreshed from
// src/data/sections.ts — enum counts now live in the hero, not here).
var ComparisonMatrix = []ComparisonMatrixRow{
	{Feature: "CSS approach", Values: []MatrixCell{"Tailwind + vars", "Tailwind + DaisyUI", "Tailwind v4 CSS-first"}},
	{Feature: "JavaScript", Values: []MatrixCell{"Alpine.js", "DaisyUI JS", "HATEOAS (enhances HTML)"}},
	{Feature: "Requires Node.js", Values: []MatrixCell{MatrixNo, MatrixYes, MatrixNo}},
	{Feature: "Typed props enums", Values: []MatrixCell{MatrixNo, MatrixNo, "58 (tested IsValid)"}},
	{Feature: "CSP nonce support", Values: []MatrixCell{MatrixYes, MatrixNo, MatrixYes}},
	{Feature: "Dark mode", Values: []MatrixCell{"CSS vars", "DaisyUI", "Tailwind dark: (tested)"}},
	{Feature: "HTMX integration", Values: []MatrixCell{MatrixNo, MatrixNo, MatrixYes}},
	{Feature: "Native SVG charts", Values: []MatrixCell{MatrixNo, MatrixNo, "yes (zero-JS)"}},
	{Feature: "ECharts adapter", Values: []MatrixCell{MatrixNo, MatrixNo, "yes (opt-in)"}},
	{Feature: "Standalone library", Values: []MatrixCell{MatrixNo, MatrixNo, MatrixYes}},
}

// UseCase is one card of the UseCases section.
type UseCase struct {
	Icon  icons.Name
	Title string
	Desc  string
}

// UseCases powers the UseCases section (ported from src/data/sections.ts).
var UseCases = []UseCase{
	{
		Icon:  icons.Squares2x2,
		Title: "Admin Dashboards",
		Desc:  "Tables, stat cards, charts, and sidebars. Build data-dense panels in minutes, not days.",
	},
	{
		Icon:  icons.Clipboard,
		Title: "Forms & CRUD",
		Desc:  "22 form components with validation, comboboxes, date pickers, and accessible error handling.",
	},
	{
		Icon:  icons.Bars3,
		Title: "SaaS Interfaces",
		Desc:  "Navigation, modals, drawers, tabs, and breadcrumbs. Full app chrome without a frontend framework.",
	},
}

// heroCode is the templ sample rendered (and highlighted) in the hero code
// window. Ported from src/data/hero-code.ts.
const heroCode = `package main

import (
    "github.com/larsartmann/templ-components/display"
    "github.com/larsartmann/templ-components/feedback"
    "github.com/larsartmann/templ-components/layout"
)

templ Page() {
    @layout.Base(layout.DefaultPageProps()) {
        @layout.ThemeScript("")
        @display.PageHeader(display.PageHeaderProps{
            Title: "Dashboard",
            Subtitle: "Welcome back",
        })
        @display.Grid(display.GridProps{
            Cols: display.GridCols3,
        }) {
            @display.StatCard(display.StatCardProps{
                Label: "Revenue",
                Value: "$42,189",
                Trend: display.TrendUp,
            })
        }
        @feedback.Toast(feedback.ToastProps{
            Message: "Data refreshed!",
            Type:   feedback.ToastSuccess,
        })
    }
}`
