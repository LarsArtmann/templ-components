package utils

import (
	"bytes"
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
		"content carrying </script> breaks out; contract documents why callers must not pass one",
		func(t *testing.T) {
			t.Parallel()

			// This is the documented hazard, not sanctioned behavior: the break-out
			// render is what the contract above forbids. Pinned so a future change
			// to the verbatim-wrap rule is a deliberate, visible decision.
			var buf bytes.Buffer

			hazard := "var s = \"</script>\";"

			if err := ScriptComponent("", hazard, "hazard demo").Render(t.Context(), &buf); err != nil {
				t.Fatalf("render failed: %v", err)
			}

			if got := strings.Count(buf.String(), "</script>"); got != 2 {
				t.Fatalf("expected the hazard render to carry 2 close tags (the documented break-out), got %d", got)
			}
		},
	)
}
