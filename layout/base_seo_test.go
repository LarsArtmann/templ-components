package layout

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// seoRender renders Base with the given props for SEO assertions.
func seoRender(t *testing.T, props PageProps) string {
	t.Helper()

	return utils.Render(t, Base(props))
}

func TestBase_SEOMeta_ZeroValueOmitsAll(t *testing.T) {
	t.Parallel()

	props := DefaultPageProps()
	props.Title = "SEO zero"
	output := seoRender(t, props)

	for _, forbidden := range []string{
		`name="robots"`,
		`rel="canonical"`,
		`rel="alternate"`,
		`application/ld+json`,
	} {
		if strings.Contains(output, forbidden) {
			t.Errorf("zero SEOMeta should not emit %q", forbidden)
		}
	}
}

func TestBase_SEOMeta_NoIndex(t *testing.T) {
	t.Parallel()

	props := DefaultPageProps()
	props.Title = "Noindex page"
	props.SEO.NoIndex = true
	output := seoRender(t, props)

	if !strings.Contains(output, `<meta name="robots" content="noindex">`) {
		t.Error("NoIndex should emit the robots noindex meta tag")
	}
}

func TestBase_SEOMeta_Canonical(t *testing.T) {
	t.Parallel()

	props := DefaultPageProps()
	props.Title = "Canonical page"
	props.SEO.Canonical = "https://example.com/cv"
	output := seoRender(t, props)

	if !strings.Contains(output, `<link rel="canonical" href="https://example.com/cv">`) {
		t.Error("Canonical should emit the canonical link tag")
	}
}

func TestBase_SEOMeta_Alternates(t *testing.T) {
	t.Parallel()

	props := DefaultPageProps()
	props.Title = "Alternates page"
	props.SEO.Alternates = []SEOAlternate{
		{Lang: "en", URL: "https://example.com/cv"},
		{Lang: "de", URL: "https://example.com/de/cv"},
		{Lang: "x-default", URL: "https://example.com/cv"},
	}
	output := seoRender(t, props)

	enIdx := strings.Index(output, `<link rel="alternate" hreflang="en" href="https://example.com/cv">`)
	deIdx := strings.Index(output, `<link rel="alternate" hreflang="de" href="https://example.com/de/cv">`)
	xdIdx := strings.Index(output, `<link rel="alternate" hreflang="x-default" href="https://example.com/cv">`)

	if enIdx == -1 || deIdx == -1 || xdIdx == -1 {
		t.Fatalf("all three alternates should render, indices: en=%d de=%d x-default=%d", enIdx, deIdx, xdIdx)
	}

	if !(enIdx < deIdx && deIdx < xdIdx) {
		t.Error("alternates should render in declaration order")
	}
}

func TestBase_SEOMeta_JSONLD(t *testing.T) {
	t.Parallel()

	props := DefaultPageProps()
	props.Title = "JSON-LD page"
	props.SEO.JSONLD = `{"@context":"https://schema.org","@type":"Person","name":"Test"}`
	output := seoRender(t, props)

	if !strings.Contains(output, `<script type="application/ld+json">{"@context":"https://schema.org","@type":"Person","name":"Test"}</script>`) {
		t.Error("JSONLD should be embedded verbatim as an ld+json script")
	}
}

func TestBase_SEOMeta_ImplementationExample(t *testing.T) {
	// Documents the consumer migration this field set exists for: a page
	// that previously hand-rolled all four tags in HeadContent (see the
	// CV project's screenHeadContent) composes them declaratively instead.
	t.Parallel()

	props := DefaultPageProps()
	props.Title = "Lars Artmann"
	props.Description = "Software Architect"
	props.SEO = SEOMeta{
		NoIndex:   true,
		Canonical: "https://lars.software/cv",
		Alternates: []SEOAlternate{
			{Lang: "en", URL: "https://lars.software/cv"},
			{Lang: "de", URL: "https://lars.software/de/cv"},
		},
		JSONLD: `{"@type":"ProfilePage"}`,
	}
	output := seoRender(t, props)

	for _, want := range []string{
		`<meta name="robots" content="noindex">`,
		`<link rel="canonical" href="https://lars.software/cv">`,
		`hreflang="en"`,
		`hreflang="de"`,
		`application/ld+json`,
	} {
		if !strings.Contains(output, want) {
			t.Errorf("combined SEOMeta should contain %q", want)
		}
	}
}
