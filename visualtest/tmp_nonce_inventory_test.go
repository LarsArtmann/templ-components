package visualtest

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"testing"
)

func TestTmpNonceInventory(t *testing.T) {
	server := StartDemoServer(t)

	client := &http.Client{}
	re := regexp.MustCompile(`(?s)<script nonce="">(.{0,90})`)

	pages := []string{"/", "/layout", "/display", "/feedback", "/forms", "/navigation", "/icons", "/htmx", "/datastar", "/wire", "/kanban", "/echarts", "/recipes", "/users", "/error-pages", "/recipes/dashboard", "/errors/404"}
	for _, p := range pages {
		resp, err := client.Get(server.BaseURL() + p)
		if err != nil {
			t.Fatalf("GET %s: %v", p, err)
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		matches := re.FindAllStringSubmatch(string(body), -1)
		if len(matches) == 0 {
			continue
		}
		seen := map[string]bool{}
		out := ""
		for _, m := range matches {
			snippet := regexp.MustCompile(`\s+`).ReplaceAllString(m[1], " ")
			if !seen[snippet] {
				seen[snippet] = true
				out += fmt.Sprintf("\n      [%s]", snippet)
			}
		}
		t.Logf("PAGE %s: %d empty-nonce scripts%s", p, len(matches), out)
	}
}
