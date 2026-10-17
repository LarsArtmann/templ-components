package utils

import (
	"context"
	"html"
	"strings"
	"testing"
)

// FuzzScriptComponent pins the writer's contract edges against hostile
// inputs (#360): the nonce attribute is HTML-escaped and never empty when a
// nonce resolves, the script body is written verbatim (templ's own script
// context handles JS-level sanitization upstream of this writer), and no
// input shape panics.
func FuzzScriptComponent(f *testing.F) {
	f.Add(`alert("x");`, "nonce-123")
	f.Add(`</script><script>alert(1)`, "")
	f.Add(`quotes " ' angles < > ampersand &`, `n"ce<&>`)
	f.Add("unicode ünïcödé 日本語 \U0001F642", "n")
	f.Add("", "n")
	f.Add("var tc = `template ${literal}`;", "")

	f.Fuzz(func(t *testing.T, js, nonce string) {
		var sb strings.Builder

		if err := ScriptComponent(nonce, js, "fuzz").Render(context.Background(), &sb); err != nil {
			t.Fatalf("render error: %v", err)
		}

		out := sb.String()

		if strings.Contains(out, `nonce=""`) {
			t.Fatalf("nonce=%q rendered an empty nonce attribute: %q", nonce, out)
		}

		if nonce != "" {
			want := `nonce="` + html.EscapeString(nonce) + `"`
			if !strings.Contains(out, want) {
				t.Fatalf("nonce=%q missing escaped attribute %q in %q", nonce, want, out)
			}
		}

		if !strings.HasSuffix(out, "\n"+js+"</script>\n") {
			t.Fatalf("script body not written verbatim (nonce=%q): %q", nonce, out)
		}

		if !strings.HasPrefix(out, "<script") {
			t.Fatalf("output does not open with <script: %q", out)
		}
	})
}
