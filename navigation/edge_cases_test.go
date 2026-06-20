package navigation

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestNavLinkEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name   string
		props  NavLinkProps
		active string
		want   []string
	}{
		{"inactive link", NavLinkProps{Href: "/about", Text: "About"}, "/", []string{"About"}},
		{"external link", NavLinkProps{Href: "https://example.com", Text: "Ext", External: true}, "", []string{`target="_blank"`, `rel="noopener noreferrer"`}},
		{"custom id/class", NavLinkProps{BaseProps: utils.BaseProps{ID: "nl", Class: "mt-2"}, Href: "/", Text: "Home"}, "/", []string{`id="nl"`, "mt-2"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, NavLink(tt.props, tt.active))
			utils.AssertContainsAll(t, output, tt.want...)
		})
	}
}

func TestMobileNavLinkEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name   string
		props  NavLinkProps
		active string
		want   []string
	}{
		{"inactive", NavLinkProps{Href: "/about", Text: "About"}, "/", []string{"About"}},
		{"active", NavLinkProps{Href: "/", Text: "Home"}, "/", []string{"Home", "bg-blue-50"}},
		{"custom id/class", NavLinkProps{BaseProps: utils.BaseProps{ID: "mnl", Class: "mt-2"}, Href: "/", Text: "Home"}, "/", []string{`id="mnl"`, "mt-2"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, MobileNavLink(tt.props, tt.active))
			utils.AssertContainsAll(t, output, tt.want...)
		})
	}
}

func TestPaginationEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		props PaginationProps
		want  []string
	}{
		{"two pages", PaginationProps{CurrentPage: 1, TotalPages: 2, BaseURL: "/items"}, []string{"href="}},
		{"custom id/class", PaginationProps{BaseProps: utils.BaseProps{ID: "pg", Class: "mt-2"}, CurrentPage: 1, TotalPages: 3, BaseURL: "/items"}, []string{`id="pg"`, "mt-2"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Pagination(tt.props))
			utils.AssertContainsAll(t, output, tt.want...)
		})
	}
}

func TestNavEdgeCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		props NavProps
		want  []string
	}{
		{"no brand", NavProps{Links: testNavLinks}, []string{navItemHome}},
		{"sticky", NavProps{Links: testNavLinks, Sticky: true}, []string{"sticky"}},
		{"custom id/class", NavProps{BaseProps: utils.BaseProps{ID: "nv", Class: "mt-2"}, Links: testNavLinks}, []string{`id="nv"`, "mt-2"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, Nav(tt.props))
			utils.AssertContainsAll(t, output, tt.want...)
		})
	}
}
