package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// TestCopyButtonNoScriptSkipsScript pins the #63 fragment contract: a
// NoScript button renders markup only — no singleton script element — so it
// is safe inside htmx-swapped fragments where script elements inserted via
// innerHTML never execute.
func TestCopyButtonNoScriptSkipsScript(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, CopyButton(CopyButtonProps{
		Text:     "copy me",
		NoScript: true,
	}))
	utils.AssertContains(t, output, `data-tc-copy="copy me"`)
	utils.AssertNotContains(t, output, "<script")
	utils.AssertNotContains(t, output, "tcCopyAttached")

	// Default stays non-breaking: the script still renders.
	def := utils.Render(t, CopyButton(CopyButtonProps{Text: "copy me"}))
	utils.AssertContains(t, def, "tcCopyAttached")
}

// TestCopyButtonScriptMount pins the page-shell mount: CopyButtonScript
// renders the delegation script standalone, with the nonce only when set
// (empty nonce omits the attribute — strict-CSP browsers reject nonce="").
func TestCopyButtonScriptMount(t *testing.T) {
	t.Parallel()
	nonced := utils.Render(t, CopyButtonScript("abc123"))
	utils.AssertContains(t, nonced, "<script")
	utils.AssertContains(t, nonced, `nonce="abc123"`)
	utils.AssertContains(t, nonced, "tcCopyAttached")

	unnonced := utils.Render(t, CopyButtonScript(""))
	utils.AssertContains(t, unnonced, "<script")
	utils.AssertNotContains(t, unnonced, `nonce=""`)
}
