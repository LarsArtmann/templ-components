package pages

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
	"github.com/larsartmann/templ-components/website/internal/build"
	"github.com/larsartmann/templ-components/website/internal/md"
)

// goldenStats are deterministic library facts for golden rendering.
var goldenStats = build.Stats{Components: 123, Icons: 105, Enums: 58, Modules: 7}

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
		{Name: "landing", HTML: render(t, Landing(goldenStats, StarsLabel(1024), goldenNonce))},
		{Name: "landing-no-stars", HTML: render(t, Landing(goldenStats, StarsLabel(0), goldenNonce))},
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
