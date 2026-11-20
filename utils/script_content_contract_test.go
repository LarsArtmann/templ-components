package utils

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// TestScriptComponentSingleCloseTag pins the script-content contract: the JS
// string a caller passes MUST NOT contain a literal `</script>` — the
// component wraps the content verbatim, so a close tag inside the content
// breaks out of the element and the trailing close tag becomes a stray
// second close (the v1.27.0 toast bug class: a raw tag tail kept inside a
// string now fed to ScriptComponent). The assertion documents the contract
// mechanically rather than by comment alone: for a clean script the render
// carries exactly ONE close tag, at the tail.
func TestScriptComponentSingleCloseTag(t *testing.T) {
	t.Parallel()

	t.Run("clean script renders exactly one close tag at the tail", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer

		if err := ScriptComponent("nonce-123", "console.log(1);", "contract test").
			Render(t.Context(), &buf); err != nil {
			t.Fatalf("render failed: %v", err)
		}

		out := buf.String()

		if got := strings.Count(out, "</script>"); got != 1 {
			t.Fatalf("expected exactly 1 close tag, got %d: %q", got, out)
		}

		if !strings.HasSuffix(out, "</script>\n") {
			t.Fatalf("close tag is not at the tail: %q", out)
		}
	})

	t.Run(
		"content carrying </script> refuses at render time (the caller-error class closes here)",
		func(t *testing.T) {
			t.Parallel()

			// Since v1.27.2 the verbatim-wrap hazard is a render-time refusal,
			// not a documented hazard: the component fails loudly and names the
			// errLabel, so a caller can never ship the element break-out (the
			// v1.27.0 toast class). The escape form <\/script> is the sanctioned
			// way to carry a close tag inside a JS string.
			var buf bytes.Buffer

			hazard := "var s = \"</script>\";"

			err := ScriptComponent("", hazard, "hazard demo").Render(t.Context(), &buf)
			if err == nil {
				t.Fatalf("the hazard render must refuse, got output %q", buf.String())
			}

			if !strings.Contains(err.Error(), "hazard demo") {
				t.Fatalf("refusal must name the errLabel, got: %v", err)
			}

			if !errors.Is(err, ErrScriptCloseTag) {
				t.Fatalf("refusal must wrap the sentinel ErrScriptCloseTag, got: %v", err)
			}

			if strings.Contains(buf.String(), "</script>") {
				t.Fatalf("a refused render must not emit any close tag, got %q", buf.String())
			}
		},
	)
}
