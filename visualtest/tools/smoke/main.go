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
	"context"
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

// defaultTimeoutMS is the default per-request timeout.
const defaultTimeoutMS = 5000

// smokeConfig is the parsed invocation: one request plus its assertions.
type smokeConfig struct {
	verb            string
	url             string
	body            string
	wantStatus      int
	timeoutMS       int
	headers         multiFlag
	wantContains    multiFlag
	wantNotContains multiFlag
}

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	cfg := parseFlags(args)

	req, err := http.NewRequestWithContext(
		context.Background(),
		cfg.verb,
		cfg.url,
		strings.NewReader(cfg.body),
	)
	if err != nil {
		return fail(
			http.StatusBadRequest,
			"build request %s %s: %v — check the URL and method spelling",
			cfg.verb,
			cfg.url,
			err,
		)
	}

	for _, h := range cfg.headers {
		k, v, found := strings.Cut(h, ":")
		if !found {
			return fail(http.StatusBadRequest, "header %q is not k:v — use e.g. -H \"Sec-Fetch-Site: same-origin\"", h)
		}

		req.Header.Set(strings.TrimSpace(k), strings.TrimSpace(v))
	}

	client := &http.Client{ //nolint:exhaustruct_v5 // a smoke client wants exactly the zero-value transport/redirect/jar plus a timeout
		Timeout: time.Duration(cfg.timeoutMS) * time.Millisecond,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fail(
			http.StatusServiceUnavailable,
			"%s %s failed: %v — is the server running and the PORT correct?",
			cfg.verb,
			cfg.url,
			err,
		)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fail(http.StatusInternalServerError, "read response body: %v", err)
	}

	if code := assertResponse(cfg, resp.StatusCode, string(respBody)); code != 0 {
		return code
	}

	fmt.Fprintf(os.Stdout, "ok %d %s %s\n", resp.StatusCode, cfg.verb, cfg.url)

	return 0
}

// parseFlags builds the invocation from argv-style args.
func parseFlags(args []string) smokeConfig {
	var (
		verb            string
		body            string
		wantStatus      int
		timeoutMS       int
		headers         multiFlag
		wantContains    multiFlag
		wantNotContains multiFlag
	)

	flags := flag.NewFlagSet("smoke", flag.ContinueOnError)
	flags.StringVar(&verb, "X", http.MethodGet, "HTTP method")
	flags.StringVar(&body, "d", "", "request body")
	flags.IntVar(&wantStatus, "status", http.StatusOK, "expected status code")
	flags.IntVar(&timeoutMS, "timeout-ms", defaultTimeoutMS, "request timeout in milliseconds")
	flags.Var(&headers, "H", "request header as k:v (repeatable)")
	flags.Var(&wantContains, "contains", "marker that must appear in the response body (repeatable)")
	flags.Var(&wantNotContains, "not-contains", "marker that must NOT appear in the response body (repeatable)")

	if err := flags.Parse(args); err != nil {
		fail(http.StatusBadRequest, "parse flags: %v", err)
	}

	if flags.NArg() != 1 {
		fail(
			http.StatusBadRequest,
			"exactly one URL argument is required (got %d) — pass the target as the last argument",
			flags.NArg(),
		)
	}

	return smokeConfig{
		verb:            verb,
		url:             flags.Arg(0),
		body:            body,
		wantStatus:      wantStatus,
		timeoutMS:       timeoutMS,
		headers:         headers,
		wantContains:    wantContains,
		wantNotContains: wantNotContains,
	}
}

// assertResponse checks status and body markers, returning a nonzero exit
// code (the offending status) on the first mismatch.
func assertResponse(cfg smokeConfig, gotStatus int, body string) int {
	if gotStatus != cfg.wantStatus {
		return fail(
			gotStatus,
			"%s %s returned %d, want %d — body: %.200s",
			cfg.verb,
			cfg.url,
			gotStatus,
			cfg.wantStatus,
			body,
		)
	}

	for _, marker := range cfg.wantContains {
		if !strings.Contains(body, marker) {
			return fail(
				http.StatusInternalServerError,
				"response of %s %s is missing required marker %q — the page or endpoint changed, or the server is serving stale content",
				cfg.verb,
				cfg.url,
				marker,
			)
		}
	}

	for _, marker := range cfg.wantNotContains {
		if strings.Contains(body, marker) {
			return fail(
				http.StatusInternalServerError,
				"response of %s %s contains forbidden marker %q — the state change it implies did not happen",
				cfg.verb,
				cfg.url,
				marker,
			)
		}
	}

	return 0
}

// fail prints a cause-and-fix message and returns the process exit code.
func fail(code int, format string, args ...any) int {
	fmt.Fprintf(os.Stderr, "smoke: "+format+"\n", args...)

	return code
}
