package visualtest

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

// TestDemoKanbanHTTPContracts pins the demo's kanban endpoint security
// semantics over plain HTTP — the layer a browser e2e cannot isolate:
//
//   - the move form's CSRF token is mandatory (missing or wrong → 403),
//   - the bodyless add/reset POSTs enforce same-origin (no Origin → 403),
//   - an add to an unknown column 404s instead of 200-nooping,
//   - reset removes the cards an add created.
//
// The valid-token leg follows the browser's real flow: GET /, harvest the
// hidden csrf_token input from the rendered board, submit it back.
func TestDemoKanbanHTTPContracts(t *testing.T) {
	server := StartDemoServer(t)
	base := server.BaseURL()

	// Test fixture: the zero-value transport/redirect/jar behavior is exactly
	// what these contract probes want.
	client := &http.Client{
		Timeout: demoHTTPTimeout,
	}

	get := func(path string) string {
		t.Helper()

		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, base+path, nil)
		if err != nil {
			t.Fatalf("visualtest[demo]: build GET %s: %v", path, err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("visualtest[demo]: GET %s: %v", path, err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("visualtest[demo]: GET %s body: %v", path, err)
		}

		return string(body)
	}

	post := func(path string, headers map[string]string, body string) int {
		t.Helper()

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, base+path, strings.NewReader(body))
		if err != nil {
			t.Fatalf("visualtest[demo]: build POST %s: %v", path, err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("visualtest[demo]: POST %s: %v", path, err)
		}
		defer resp.Body.Close()

		_, _ = io.Copy(io.Discard, resp.Body)

		return resp.StatusCode
	}

	const (
		noHeaders       = ""
		originHeader    = "Origin"
		fetchSiteHeader = "Sec-Fetch-Site"

		movedCard   = "kb-1"
		movedColumn = "progress"
		addedCardID = "kb-new-0"
	)

	// The browser flow: the CSRF token arrives inside the rendered page.
	// Harvest scoped to the htmx board's own form — the demo index renders
	// other components first and one of them could carry the same input
	// name with a different value.
	page := get("/")

	boardStart := strings.Index(page, `id="kanban-demo-htmx"`)
	if boardStart < 0 {
		t.Fatal("visualtest[demo]: rendered page has no htmx kanban board")
	}

	token := kanbanDemoCSRFTokenFromHTML(t, page[boardStart:])
	if token == "" {
		t.Fatal("visualtest[demo]: harvested csrf_token is empty")
	}

	sameOrigin := map[string]string{fetchSiteHeader: "same-origin"}

	// Add: same-origin is mandatory; unknown columns 404; known columns 200.
	if got := post("/api/kanban/htmx/add/backlog", nil, ""); got != http.StatusForbidden {
		t.Fatalf("visualtest[demo]: add without same-origin proof = %d, want 403", got)
	}

	if got := post("/api/kanban/htmx/add/nope", sameOrigin, ""); got != http.StatusNotFound {
		t.Fatalf("visualtest[demo]: add to unknown column = %d, want 404", got)
	}

	if got := post("/api/kanban/htmx/add/backlog", sameOrigin, ""); got != http.StatusOK {
		t.Fatalf("visualtest[demo]: same-origin add = %d, want 200", got)
	}

	if page := get("/"); !strings.Contains(page, `data-tc-kanban-card="`+addedCardID+`"`) {
		t.Fatal("visualtest[demo]: added card missing from re-rendered board")
	}

	// Move: the CSRF token is mandatory — missing and wrong both 403; the
	// harvested token passes.
	moveBody := "card=" + movedCard + "&column=" + movedColumn + "&index=0"

	if got := post("/api/kanban/htmx", sameOrigin, moveBody); got != http.StatusForbidden {
		t.Fatalf("visualtest[demo]: move without CSRF token = %d, want 403", got)
	}

	if got := post("/api/kanban/htmx", sameOrigin, moveBody+"&csrf_token=wrong"); got != http.StatusForbidden {
		t.Fatalf("visualtest[demo]: move with wrong CSRF token = %d, want 403", got)
	}

	if got := post("/api/kanban/htmx", sameOrigin, moveBody+"&csrf_token="+token); got != http.StatusOK {
		t.Fatalf("visualtest[demo]: move with harvested CSRF token = %d, want 200", got)
	}

	// Reset (Origin-header branch this time) removes the added card again.
	if got := post("/api/kanban/htmx/reset", map[string]string{originHeader: base}, ""); got != http.StatusOK {
		t.Fatalf("visualtest[demo]: reset via Origin header = %d, want 200", got)
	}

	if page := get("/"); strings.Contains(page, `data-tc-kanban-card="`+addedCardID+`"`) {
		t.Fatal("visualtest[demo]: reset did not remove the added card")
	}

	// The Datastar-wrapped add goes through the same same-origin gate.
	if got := post("/api/kanban/datastar/add/backlog", nil, ""); got != http.StatusForbidden {
		t.Fatalf("visualtest[demo]: datastar add without same-origin proof = %d, want 403", got)
	}

	if got := post("/api/kanban/datastar/add/backlog", sameOrigin, ""); got != http.StatusOK {
		t.Fatalf("visualtest[demo]: datastar same-origin add = %d, want 200", got)
	}

	server.FailIfServerErrors(t)
}

// kanbanDemoCSRFTokenFromHTML extracts the move form's hidden csrf_token
// value from a rendered demo page — exactly what a browser submits back.
func kanbanDemoCSRFTokenFromHTML(t *testing.T, html string) string {
	t.Helper()

	const marker = `name="csrf_token" value="`

	_, rest, found := strings.Cut(html, marker)
	if !found {
		t.Fatal("visualtest[demo]: rendered page has no csrf_token hidden input")
	}

	value, _, terminated := strings.Cut(rest, `"`)
	if !terminated {
		t.Fatal("visualtest[demo]: csrf_token input value unterminated")
	}

	return value
}
