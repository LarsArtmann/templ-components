// smoke is a tiny HTTP smoke client for the demo (and any local server):
// issue one request, assert the status code and body markers, exit 0/1.
// It replaces the throwaway /tmp smoke clients — raw curl/wget are banned
// in this environment, and ad-hoc one-offs cannot be reused or linted.
//
// Examples:
//
//	smoke -status 403 -X POST http://localhost:8901/api/kanban/htmx/add/backlog
//	smoke -X POST -H "Sec-Fetch-Site: same-origin" \
//	    -contains 'data-tc-kanban-card' \
//	    http://localhost:8901/api/kanban/htmx/add/backlog
//	smoke -contains '<title>Demo' http://localhost:8901/
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// multiFlag collects repeated -flag values into a slice.
type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, ", ") }

func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)

	return nil
}

func main() {
	var (
		verb            = flag.String("X", http.MethodGet, "HTTP method")
		body            = flag.String("d", "", "request body")
		wantStatus      = flag.Int("status", http.StatusOK, "expected status code")
		timeoutMS       = flag.Int("timeout-ms", 5000, "request timeout in milliseconds")
		headers         multiFlag
		wantContains    multiFlag
		wantNotContains multiFlag
	)
	flag.Var(&headers, "H", "request header as k:v (repeatable)")
	flag.Var(&wantContains, "contains", "marker that must appear in the response body (repeatable)")
	flag.Var(&wantNotContains, "not-contains", "marker that must NOT appear in the response body (repeatable)")
	flag.Parse()

	if flag.NArg() != 1 {
		fail(http.StatusBadRequest, "exactly one URL argument is required (got %d) — pass the target as the last argument", flag.NArg())
	}

	url := flag.Arg(0)

	req, err := http.NewRequest(*verb, url, strings.NewReader(*body))
	if err != nil {
		fail(http.StatusBadRequest, "build request %s %s: %v — check the URL and method spelling", *verb, url, err)
	}

	for _, h := range headers {
		k, v, found := strings.Cut(h, ":")
		if !found {
			fail(http.StatusBadRequest, "header %q is not k:v — use e.g. -H \"Sec-Fetch-Site: same-origin\"", h)
		}

		req.Header.Set(strings.TrimSpace(k), strings.TrimSpace(v))
	}

	client := &http.Client{Timeout: time.Duration(*timeoutMS) * time.Millisecond}

	resp, err := client.Do(req)
	if err != nil {
		fail(http.StatusServiceUnavailable, "%s %s failed: %v — is the server running and the PORT correct?", *verb, url, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fail(http.StatusInternalServerError, "read response body: %v", err)
	}

	if resp.StatusCode != *wantStatus {
		fail(resp.StatusCode, "%s %s returned %d, want %d — body: %.200s", *verb, url, resp.StatusCode, *wantStatus, string(respBody))
	}

	text := string(respBody)

	for _, marker := range wantContains {
		if !strings.Contains(text, marker) {
			fail(http.StatusInternalServerError, "response of %s %s is missing required marker %q — the page or endpoint changed, or the server is serving stale content", *verb, url, marker)
		}
	}

	for _, marker := range wantNotContains {
		if strings.Contains(text, marker) {
			fail(http.StatusInternalServerError, "response of %s %s contains forbidden marker %q — the state change it implies did not happen", *verb, url, marker)
		}
	}

	fmt.Printf("ok %d %s %s\n", resp.StatusCode, *verb, url)
}

// fail prints a cause-and-fix message and exits with the given code.
func fail(code int, format string, args ...any) {
	fmt.Fprintf(os.Stderr, "smoke: "+format+"\n", args...)
	os.Exit(code)
}
