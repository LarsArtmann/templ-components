package display

import (
	"context"
	"io"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

// proseMarkdownSample renders a representative longform-content body —
// every element family the .tc-prose CSS covers (headings, prose, links,
// lists, blockquote, inline + block code, table, image, rule).
func proseMarkdownSample() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `
<h2>Release notes</h2>
<p>Longform content renders <strong>bold leads</strong>, <em>emphasis</em>, and
<a href="https://templ.guide">links</a> without per-element classes.</p>
<h3>Highlights</h3>
<ul>
	<li>Typed props for every component</li>
	<li>Dark mode with
		<ol>
			<li>class strategy</li>
			<li>color-scheme</li>
		</ol>
	</li>
</ul>
<blockquote>Cite your sources; style them once.</blockquote>
<p>Inline <code>utils.Render</code> references and a block:</p>
<pre><code class="language-go">err := Prose(props).Render(ctx, w)</code></pre>
<hr>
<table>
	<thead>
		<tr><th>Package</th><th>Components</th></tr>
	</thead>
	<tbody>
		<tr><td>display</td><td>52</td></tr>
		<tr><td>forms</td><td>26</td></tr>
	</tbody>
</table>
<img src="/screenshots/prose.png" alt="Prose rendering sample" width="640" height="360">
`)

		return err
	})
}

// TestGoldenSweepProse snapshots the Prose wrapper: the three sizes, BaseProps
// propagation, and the full markdown sample the .tc-prose CSS targets.
func TestGoldenSweepProse(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "prose_default", HTML: utils.Render(t, Prose(ProseProps{Children: proseMarkdownSample()}))},
		{
			Name: "prose_small",
			HTML: utils.Render(t, Prose(ProseProps{
				Children: proseMarkdownSample(),
				Size:     ProseSizeSM,
			})),
		},
		{
			Name: "prose_large",
			HTML: utils.Render(t, Prose(ProseProps{
				Children: proseMarkdownSample(),
				Size:     ProseSizeLG,
			})),
		},
		{
			Name: "prose_props",
			HTML: utils.Render(t, Prose(ProseProps{
				Children: proseMarkdownSample(),
				Size:     ProseSizeMD,
				BaseProps: utils.BaseProps{
					ID:    "release-notes",
					Class: "max-w-3xl mx-auto",
				},
			})),
		},
	})
}
