package layout

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// TestMinimal_SEO verifies Minimal's SEOMeta support: zero value emits
// nothing (Minimal stays dependency-free), and every populated field emits
// the same tags Base emits — the two shells must not drift.
func TestMinimal_SEO(t *testing.T) {
	t.Parallel()

	t.Run("zero value emits no seo tags", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, Minimal(DefaultMinimalProps()))

		utils.AssertNotContains(t, output, `name="robots"`)
		utils.AssertNotContains(t, output, `rel="canonical"`)
		utils.AssertNotContains(t, output, `hreflang=`)
		utils.AssertNotContains(t, output, "application/ld+json")
	})

	t.Run("full seo block", func(t *testing.T) {
		t.Parallel()

		props := MinimalProps{
			Title: "Printable CV",
			SEO: SEOMeta{
				NoIndex:   true,
				Canonical: "https://example.com/cv",
				Alternates: []SEOAlternate{
					{Lang: "en", URL: "https://example.com/cv"},
					{Lang: "de", URL: "https://example.com/de/cv"},
				},
				JSONLD: `{"@type":"Person"}`,
			},
		}

		output := utils.Render(t, Minimal(props))

		utils.AssertContainsAll(t, output,
			`<meta name="robots" content="noindex">`,
			`<link rel="canonical" href="https://example.com/cv">`,
			`<link rel="alternate" hreflang="en" href="https://example.com/cv">`,
			`<link rel="alternate" hreflang="de" href="https://example.com/de/cv">`,
			`<script type="application/ld+json">{"@type":"Person"}</script>`,
		)
	})

	t.Run("base and minimal emit identical tag sets", func(t *testing.T) {
		t.Parallel()

		shared := SEOMeta{
			NoIndex:   true,
			Canonical: "https://example.com/cv",
			Alternates: []SEOAlternate{
				{Lang: "en", URL: "https://example.com/cv"},
				{Lang: "x-default", URL: "https://example.com/cv"},
			},
			JSONLD: `{"@type":"WebPage"}`,
		}

		minimalOut := utils.Render(t, Minimal(MinimalProps{Title: "T", SEO: shared}))

		pageProps := DefaultPageProps()
		pageProps.Title = "T"
		pageProps.SEO = shared
		baseOut := utils.Render(t, Base(pageProps))

		for _, tag := range []string{
			`<meta name="robots" content="noindex">`,
			`<link rel="canonical" href="https://example.com/cv">`,
			`<link rel="alternate" hreflang="x-default" href="https://example.com/cv">`,
			`<script type="application/ld+json">{"@type":"WebPage"}</script>`,
		} {
			if !strings.Contains(minimalOut, tag) || !strings.Contains(baseOut, tag) {
				t.Errorf("tag %q missing from Minimal or Base output", tag)
			}
		}
	})
}
