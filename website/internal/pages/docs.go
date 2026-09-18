package pages

// DocRef is one docs page in the sidebar registry.
type DocRef struct {
	// Slug is the URL path without extension ("getting-started/installation").
	Slug string
	// Title is the sidebar label (and <h1> fallback).
	Title string
}

// DocGroup is one sidebar section.
type DocGroup struct {
	Label string
	Docs  []DocRef
}

// pkgGoDevDoc is the sidebar entry that links off-site (mirrors the Astro
// sidebar's "Full API on pkg.go.dev").
//
//nolint:gochecknoglobals // site information architecture table
var pkgGoDevDoc = DocRef{Slug: "", Title: "Full API on pkg.go.dev"}

// docSidebar mirrors the Astro Starlight sidebar configuration
// (website/astro.config.mjs) so information architecture is unchanged.
//
//nolint:gochecknoglobals // site information architecture table
var docSidebar = []DocGroup{
	{
		Label: "Getting Started",
		Docs: []DocRef{
			{Slug: "getting-started/installation", Title: "Installation"},
			{Slug: "getting-started/quick-start", Title: "Quick Start"},
		},
	},
	{
		Label: "Guides",
		Docs: []DocRef{
			{Slug: "guides/theming", Title: "Theming"},
			{Slug: "guides/dark-mode", Title: "Dark Mode"},
			{Slug: "guides/htmx-integration", Title: "HTMX Integration"},
			{Slug: "guides/transport-wiring", Title: "Transport Wiring"},
			{Slug: "guides/kanban-board", Title: "Kanban Board"},
			{Slug: "guides/error-pages", Title: "Error Pages"},
			{Slug: "guides/accessibility", Title: "Accessibility"},
			{Slug: "guides/csp-compliance", Title: "CSP Compliance"},
			{Slug: "guides/invariants", Title: "Guarantees"},
			{Slug: "guides/version-support", Title: "Version Support"},
		},
	},
	{
		Label: "API Reference",
		Docs: []DocRef{
			{Slug: "api-reference", Title: "Public API"},
			pkgGoDevDoc,
		},
	},
	{
		Label: "Community",
		Docs: []DocRef{
			{Slug: "changelog", Title: "Changelog"},
			{Slug: "contributing", Title: "Contributing"},
			{Slug: "related-projects", Title: "Related Projects"},
		},
	},
}

// AllDocs returns the sidebar pages in reading order (used for prev/next).
func AllDocs() []DocRef {
	var all []DocRef

	for _, group := range docSidebar {
		for _, doc := range group.Docs {
			if doc.Slug == "" {
				continue
			}

			all = append(all, doc)
		}
	}

	return all
}

// Sidebar returns the sidebar groups.
func Sidebar() []DocGroup {
	return docSidebar
}
