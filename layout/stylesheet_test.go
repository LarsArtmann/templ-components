package layout

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestStylesheet(t *testing.T) {
	t.Parallel()
	t.Run("basic", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Stylesheet("/app.css", nil))
		utils.AssertContains(t, output, `rel="stylesheet"`)
		utils.AssertContains(t, output, `href="/app.css"`)
	})
	t.Run("with attrs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Stylesheet("/print.css", templ.Attributes{
			"media":       "print",
			"crossorigin": "anonymous",
		}))
		utils.AssertContains(t, output, `media="print"`)
		utils.AssertContains(t, output, `crossorigin="anonymous"`)
	})
}

func TestScriptCoverage(t *testing.T) {
	t.Parallel()
	t.Run("with attrs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Script("test-nonce", "/app.js", templ.Attributes{
			"defer":     true,
			"integrity": "sha384-abc",
		}))
		utils.AssertContains(t, output, `nonce="test-nonce"`)
		utils.AssertContains(t, output, `src="/app.js"`)
		utils.AssertContains(t, output, "defer")
	})
	t.Run("nil attrs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Script("nonce123", "/lib.js", nil))
		utils.AssertContains(t, output, `nonce="nonce123"`)
	})
}
