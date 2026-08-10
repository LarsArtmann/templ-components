package navigation

import (
	"testing"

	"github.com/larsartmann/templ-components/utils/golden"
	"github.com/larsartmann/templ-components/utils"
)

// Golden sweep for navigation components that previously lacked golden tests.
// Nav uses EnsureID for its mobile menu — safe for golden testing via ID normalization.

func TestGoldenSweepNav(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "nav_basic", HTML: utils.Render(t, Nav(NavProps{
			CurrentPath: "/dashboard",
			Links: []NavLinkProps{
				{Href: "/dashboard", Text: "Dashboard"},
				{Href: "/projects", Text: "Projects"},
				{Href: "/settings", Text: "Settings"},
			},
		}))},
	})
}

func TestGoldenSweepSimpleNav(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "simple_nav", HTML: utils.Render(t, SimpleNav(SimpleNavProps{
			BrandText:   "Demo App",
			BrandHref:   "/",
			CurrentPath: "/",
			Links: []NavLinkProps{
				{Href: "/", Text: "Home"},
				{Href: "/docs", Text: "Docs"},
				{Href: "/pricing", Text: "Pricing"},
			},
		}))},
	})
}

func TestGoldenSweepFooter(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "footer_basic", HTML: utils.Render(t, Footer(FooterProps{
			BrandText: "MyApp",
		}))},
	})
}

func TestGoldenSweepEndOfList(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "end_of_list_default", HTML: utils.Render(t, EndOfList(DefaultEndOfListProps()))},
		{Name: "end_of_list_custom", HTML: utils.Render(t, EndOfList(EndOfListProps{
			Message: "No more items to load",
		}))},
	})
}
