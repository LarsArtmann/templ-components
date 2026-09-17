package pages

import (
	"fmt"
	"sync"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
	"github.com/larsartmann/templ-components/website/internal/build"
	"github.com/larsartmann/templ-components/website/internal/md"
)

// goldenStats derives the library facts for golden rendering from the actual
// codebase (build.CountStats) instead of hand-typed literals — the literals
// drifted (58-vs-60 enums) and shipped on the live site until the derived
// comparison-matrix fix. The path is relative to THIS package's directory
// (website/internal/pages) up to the repo root; the value is cached so the
// walk happens once per test binary.
var goldenStats = sync.OnceValue(func() build.Stats {
	stats, err := build.CountStats("../../..")
	if err != nil {
		panic(fmt.Sprintf("golden_test: derive library stats: %v", err))
	}
	return stats
})

// goldenNonce is a fixed nonce so rendered scripts are byte-stable.
const goldenNonce = "golden-test-nonce"

// goldenDocsPage is a minimal docs fixture independent of content/ edits.
var goldenDocsPage = md.Page{
	Title:       "Installation",
	Description: "How to install the library",
	HTML: `<h2 id="requirements">Requirements</h2>` +
		`<p>Go 1.26+ with the templ CLI.</p>` +
		`<div class="code-block"><button type="button" class="code-copy">Copy</button>` +
		`<pre><code class="chroma">go get github.com/larsartmann/templ-components</code></pre></div>` +
		`<h3 id="tools">Tools</h3>` +
		`<p>Add the <code>templ</code> binary.</p>`,
	Headings: []md.Heading{
		{ID: "requirements", Text: "Requirements", Level: 2},
		{ID: "tools", Text: "Tools", Level: 3},
	},
}

func render(t *testing.T, component templ.Component) string {
	t.Helper()

	return utils.Render(t, component)
}

func TestGoldenSweepPages(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "landing", HTML: render(t, Landing(goldenStats(), StarsLabel(1024), goldenNonce))},
		{Name: "landing-no-stars", HTML: render(t, Landing(goldenStats(), StarsLabel(0), goldenNonce))},
		{Name: "notfound", HTML: render(t, NotFound(goldenNonce))},
		{
			Name: "docs-layout",
			HTML: render(
				t,
				DocsLayout("getting-started/installation", goldenDocsPage, nil, nil, "2026-09-13", goldenNonce),
			),
		},
	})
}
