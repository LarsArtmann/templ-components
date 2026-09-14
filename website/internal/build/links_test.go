package build

import (
	"strings"
	"testing"
)

func TestCheckLinksClean(t *testing.T) {
	t.Parallel()

	pages := []RenderedPage{
		{Path: "index.html", HTML: `<a href="/guides/theming">Theming</a>` +
			`<a href="/guides/theming#tokens">Tokens</a>` +
			`<a href="#intro">Intro</a>` +
			`<script src="/assets/js/app.js" defer></script>` +
			`<link rel="stylesheet" href="/assets/app.css"/>` +
			`<a href="https://pkg.go.dev/example">API</a>` +
			`<a href="mailto:hi@example.com">Mail</a>` +
			`<div id="intro"></div>`},
		{Path: "guides/theming.html", HTML: `<a href="/">Home</a>` +
			`<img src="/og/home.png" alt=""/>` +
			`<h2 id="tokens">Tokens</h2>`},
	}

	assets := []string{"assets/js/app.js", "assets/app.css", "og/home.png", "favicon.svg", "index.html"}

	if problems := CheckLinks(pages, assets); len(problems) != 0 {
		t.Errorf("expected clean site, got %v", problems)
	}
}

func TestCheckLinksProblems(t *testing.T) {
	t.Parallel()

	pages := []RenderedPage{
		{Path: "index.html", HTML: `<a href="/missing-page">Dead</a>` +
			`<a href="/guides/theming#no-such-id">Bad anchor</a>` +
			`<a href="#nope">Bad self anchor</a>` +
			`<script src="/assets/js/gone.js"></script>` +
			`<a href="relative/page">Relative</a>`},
		{Path: "guides/theming.html", HTML: `<h2 id="tokens">Tokens</h2>`},
	}

	problems := CheckLinks(pages, nil)
	if len(problems) != 5 {
		t.Fatalf("expected 5 problems, got %d: %v", len(problems), problems)
	}

	for _, problem := range problems {
		if !strings.HasPrefix(problem, "index.html: ") {
			t.Errorf("problem not attributed to source page: %q", problem)
		}
	}
}

func TestResolveTarget(t *testing.T) {
	t.Parallel()

	index := &linkIndex{assets: map[string]bool{"assets/app.css": true, "manifest.json": true}}

	for _, tc := range []struct {
		path string
		want string
	}{
		{path: "/", want: "index.html"},
		{path: "", want: "index.html"},
		{path: "/guides/theming", want: "guides/theming.html"},
		{path: "/404.html", want: "404.html"},
		{path: "/assets/app.css", want: "assets/app.css"},
		{path: "/manifest.json", want: "manifest.json"},
	} {
		if got := index.resolveTarget(tc.path); got != tc.want {
			t.Errorf("resolveTarget(%q) = %q, want %q", tc.path, got, tc.want)
		}
	}
}
