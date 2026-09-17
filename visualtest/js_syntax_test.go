package visualtest

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestDemoInlineScriptsAreSyntaxValid runs `node --check` over every inline
// <script> body served by the live demo (TODO #259). templ writes component
// JS as raw Go strings — a stray brace or a templ-context escape ships
// silently and only explodes in a browser console. The demo pages carry
// every component script, so parsing them all is a complete, cheap gate.
//
// Skips when node is not on PATH (CI's Build job has no node; the website
// workflow and dev machines do, so the gate still runs somewhere on every
// push that touches JS-bearing components).
func TestDemoInlineScriptsAreSyntaxValid(t *testing.T) {
	t.Parallel()

	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node not on PATH — JS syntax gate skipped (run on a machine with node, or the website workflow)")
	}

	server := StartDemoServer(t)

	routes := []string{
		"/", "/forms", "/users",
		"/recipes/dashboard", "/recipes/settings", "/recipes/login", "/recipes/auth",
		"/errors/full", "/errors/404", "/errors/404-page", "/errors/400",
		"/errors/403", "/errors/409",
	}

	scriptRe := regexp.MustCompile(`(?s)<script[^>]*>(.*?)</script>`)
	nonceAttr := regexp.MustCompile(`<script[^>]*>`)

	checked := 0

	for _, route := range routes {
		resp, err := http.Get(server.BaseURL() + route) //nolint:noctx // test fixture
		if err != nil {
			t.Fatalf("GET %s: %v", route, err)
		}

		body, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if readErr != nil {
			t.Fatalf("read %s: %v", route, readErr)
		}

		// No status assertion: the /errors/* routes deliberately serve real
		// error status codes (403/400/409/503) with full HTML error pages —
		// exactly the markup this gate wants to parse.

		for _, match := range scriptRe.FindAllStringSubmatch(string(body), -1) {
			js := match[1]
			if strings.TrimSpace(js) == "" {
				continue
			}

			checked++

			tmp := filepath.Join(t.TempDir(), "inline.js")
			if writeErr := os.WriteFile(tmp, []byte(js), 0o600); writeErr != nil {
				t.Fatalf("write temp script: %v", writeErr)
			}

			out, cmdErr := exec.Command(node, "--check", tmp).CombinedOutput() //nolint:gosec // fixed tmp path, node from PATH
			if cmdErr != nil {
				snippet := nonceAttr.FindString(string(body))
				t.Errorf("route %s: inline script fails `node --check`: %s\nscript head: %s", route, strings.TrimSpace(string(out)), firstN(snippet, 120))
			}
		}
	}

	if checked == 0 {
		t.Fatal("no inline scripts found on any demo route — the extraction broke")
	}

	t.Logf("node --check passed for %d inline script(s) across %d routes", checked, len(routes))
}

// firstN returns at most n characters of s (newlines collapsed) for failure messages.
func firstN(s string, n int) string {
	s = strings.ReplaceAll(s, "\n", " ")

	if len(s) > n {
		return s[:n] + "..."
	}

	return s
}
