package navigation

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// Coverage: exercise 0% default constructors and low-coverage functions.

func TestDefaultLoadMoreProps(t *testing.T) {
	t.Parallel()

	p := DefaultLoadMoreProps()
	if p.Label != "Load more" {
		t.Errorf("DefaultLoadMoreProps().Label = %q, want %q", p.Label, "Load more")
	}
}

func TestDefaultFooterProps(t *testing.T) {
	t.Parallel()

	p := DefaultFooterProps()
	if p.BrandText != "" {
		t.Error("DefaultFooterProps should have empty BrandText")
	}
}

func TestDefaultSidebarNavProps(t *testing.T) {
	t.Parallel()

	p := DefaultSidebarNavProps()
	if len(p.Items) != 0 {
		t.Error("DefaultSidebarNavProps should have empty Items")
	}
}

func TestEndOfListCoverageGaps(t *testing.T) {
	t.Parallel()

	t.Run("aria-label propagation", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EndOfList(EndOfListProps{
			BaseProps: utils.BaseProps{AriaLabel: "End of results"},
		}))
		utils.AssertContains(t, output, `aria-label="End of results"`)
	})

	t.Run("attrs propagation", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, EndOfList(EndOfListProps{
			BaseProps: utils.BaseProps{Attrs: templ.Attributes{"data-testid": "eol"}},
		}))
		utils.AssertContains(t, output, `data-testid="eol"`)
	})
}

func TestFooterCoverageGaps(t *testing.T) {
	t.Parallel()

	t.Run("renders with brand text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Footer(FooterProps{BrandText: "Acme Corp"}))
		utils.AssertContains(t, output, "Acme Corp")
		utils.AssertContains(t, output, "<footer")
	})

	t.Run("renders without brand text", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Footer(FooterProps{}))
		utils.AssertContains(t, output, "<footer")
	})

	t.Run("propagates class", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Footer(FooterProps{
			BaseProps: utils.BaseProps{Class: "mt-12"},
		}))
		utils.AssertContains(t, output, "mt-12")
	})

	t.Run("propagates aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Footer(FooterProps{
			BaseProps: utils.BaseProps{AriaLabel: "Site footer"},
		}))
		utils.AssertContains(t, output, `aria-label="Site footer"`)
	})
}

func TestResolveBreadcrumbURLCoverage(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name    string
		href    string
		baseURL string
		want    string
	}{
		{"empty href", "", "/base", ""},
		{"absolute URL", "https://example.com/page", "/base", "https://example.com/page"},
		{"protocol-relative", "//cdn.example.com/lib.js", "/base", "//cdn.example.com/lib.js"},
		{"relative with baseURL", "users", "/admin", "/admin/users"},
		{"relative without baseURL", "users", "", "users"},
		{"trailing slash in baseURL", "page", "/app/", "/app/page"},
		{"leading slash in href", "/page", "/app", "/app/page"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := resolveBreadcrumbURL(tt.href, tt.baseURL)
			if got != tt.want {
				t.Errorf("resolveBreadcrumbURL(%q, %q) = %q, want %q", tt.href, tt.baseURL, got, tt.want)
			}
		})
	}
}

func TestBreadcrumbSeparatorRender(t *testing.T) {
	t.Parallel()

	t.Run("custom separator renders", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, breadcrumbSeparator("→"))
		utils.AssertContains(t, output, "→")
	})

	t.Run("default separator renders", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, breadcrumbSeparator(""))
		utils.AssertNotContains(t, output, "→")
	})
}
