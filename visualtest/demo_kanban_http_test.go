package visualtest

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

// kanbanContractFixture parameterizes the shared kanban-endpoint contract
// prober per implementation (TODO #227 anti-drift tie).
type kanbanContractFixture struct {
	// name prefixes failure messages and marks which implementation failed.
	name string
	// boardID is the htmx board's root element id on the index page; the CSRF
	// token is harvested from that board's own move form.
	boardID string
	// addColumn is a known column id for the add probes.
	addColumn string
	// moveCard/moveColumn identify a valid move for the CSRF probes.
	moveCard, moveColumn string
	// addedCardID is the card id the first successful add creates.
	addedCardID string
}

// runKanbanContractProbes is the CROSS-BINDING anti-drift tie between the
// demo implementation (examples/demo/kanban_demo.go — kanbanDemoState) and
// the e2e implementation (kanban_e2e_test.go — kanbanE2EServer): the two
// CANNOT share a builder (visualtest boots the demo as an external binary),
// so they share THIS probe table instead. Both implementations must keep
// passing every probe — when one changes behavior, this table fails for the
// divergent side and names it. Contract pinned per docs/visual-testing.md's
// demo contract-marker cheat-sheet:
//
//   - the move form's CSRF token is mandatory (missing or wrong → 403),
//   - the bodyless add/reset POSTs enforce same-origin (no Origin → 403),
//   - an add to an unknown column 404s instead of 200-nooping,
//   - reset removes the cards an add created (both dialects).
//
// The valid-token leg follows the browser's real flow: GET /, harvest the
// hidden csrf_token input from the rendered board, submit it back.
func runKanbanContractProbes(t *testing.T, base string, fx kanbanContractFixture) {
	t.Helper()

	// Test fixture: the zero-value transport/redirect/jar behavior is exactly
	// what these contract probes want.
	client := &http.Client{
		Timeout: demoHTTPTimeout,
	}

	get := func(path string) string {
		t.Helper()

		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, base+path, nil)
		if err != nil {
			t.Fatalf("visualtest[%s]: build GET %s: %v", fx.name, path, err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("visualtest[%s]: GET %s: %v", fx.name, path, err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("visualtest[%s]: GET %s body: %v", fx.name, path, err)
		}

		return string(body)
	}

	post := func(path string, headers map[string]string, body string) int {
		t.Helper()

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, base+path, strings.NewReader(body))
		if err != nil {
			t.Fatalf("visualtest[%s]: build POST %s: %v", fx.name, path, err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("visualtest[%s]: POST %s: %v", fx.name, path, err)
		}
		defer resp.Body.Close()

		_, _ = io.Copy(io.Discard, resp.Body)

		return resp.StatusCode
	}

	const (
		originHeader    = "Origin"
		fetchSiteHeader = "Sec-Fetch-Site"
	)

	// The browser flow: the CSRF token arrives inside the rendered page.
	// Harvest scoped to the htmx board's own form — an index page renders
	// other components first and one of them could carry the same input
	// name with a different value.
	page := get("/")

	boardStart := strings.Index(page, `id="`+fx.boardID+`"`)
	if boardStart < 0 {
		t.Fatalf("visualtest[%s]: rendered page has no %q kanban board", fx.name, fx.boardID)
	}

	token := kanbanCSRFTokenFromHTML(t, fx.name, page[boardStart:])
	if token == "" {
		t.Fatalf("visualtest[%s]: harvested csrf_token is empty", fx.name)
	}

	sameOrigin := map[string]string{fetchSiteHeader: "same-origin"}

	// Add: same-origin is mandatory; unknown columns 404; known columns 200.
	if got := post("/api/kanban/htmx/add/"+fx.addColumn, nil, ""); got != http.StatusForbidden {
		t.Fatalf("visualtest[%s]: add without same-origin proof = %d, want 403", fx.name, got)
	}

	if got := post("/api/kanban/htmx/add/nope", sameOrigin, ""); got != http.StatusNotFound {
		t.Fatalf("visualtest[%s]: add to unknown column = %d, want 404", fx.name, got)
	}

	if got := post("/api/kanban/htmx/add/"+fx.addColumn, sameOrigin, ""); got != http.StatusOK {
		t.Fatalf("visualtest[%s]: same-origin add = %d, want 200", fx.name, got)
	}

	if page := get("/"); !strings.Contains(page, `data-tc-kanban-card="`+fx.addedCardID+`"`) {
		t.Fatalf("visualtest[%s]: added card %s missing from re-rendered board", fx.name, fx.addedCardID)
	}

	// Move: the CSRF token is mandatory — missing and wrong both 403; the
	// harvested token passes.
	moveBody := "card=" + fx.moveCard + "&column=" + fx.moveColumn + "&index=0"

	if got := post("/api/kanban/htmx", sameOrigin, moveBody); got != http.StatusForbidden {
		t.Fatalf("visualtest[%s]: move without CSRF token = %d, want 403", fx.name, got)
	}

	if got := post("/api/kanban/htmx", sameOrigin, moveBody+"&csrf_token=wrong"); got != http.StatusForbidden {
		t.Fatalf("visualtest[%s]: move with wrong CSRF token = %d, want 403", fx.name, got)
	}

	if got := post("/api/kanban/htmx", sameOrigin, moveBody+"&csrf_token="+token); got != http.StatusOK {
		t.Fatalf("visualtest[%s]: move with harvested CSRF token = %d, want 200", fx.name, got)
	}

	// Reset (Origin-header branch this time) removes the added card again.
	if got := post("/api/kanban/htmx/reset", map[string]string{originHeader: base}, ""); got != http.StatusOK {
		t.Fatalf("visualtest[%s]: reset via Origin header = %d, want 200", fx.name, got)
	}

	if page := get("/"); strings.Contains(page, `data-tc-kanban-card="`+fx.addedCardID+`"`) {
		t.Fatalf("visualtest[%s]: reset did not remove the added card", fx.name)
	}

	// The Datastar-wrapped add goes through the same same-origin gate, and
	// the datastar board has a reset of its own (leaves shared package-global
	// boards in their initial layout for whichever test runs next).
	if got := post("/api/kanban/datastar/add/"+fx.addColumn, nil, ""); got != http.StatusForbidden {
		t.Fatalf("visualtest[%s]: datastar add without same-origin proof = %d, want 403", fx.name, got)
	}

	if got := post("/api/kanban/datastar/add/"+fx.addColumn, sameOrigin, ""); got != http.StatusOK {
		t.Fatalf("visualtest[%s]: datastar same-origin add = %d, want 200", fx.name, got)
	}

	if got := post("/api/kanban/datastar/reset", sameOrigin, ""); got != http.StatusOK {
		t.Fatalf("visualtest[%s]: datastar reset = %d, want 200", fx.name, got)
	}
}

// TestDemoKanbanHTTPContracts pins the DEMO's kanban endpoint contract over
// plain HTTP — the layer a browser e2e cannot isolate. The e2e parity twin
// is TestKanbanE2EHTTPContractParity; the two share runKanbanContractProbes
// so the implementations cannot silently diverge (TODO #227).
func TestDemoKanbanHTTPContracts(t *testing.T) {
	server := StartDemoServer(t)

	runKanbanContractProbes(t, server.BaseURL(), kanbanContractFixture{
		name:        "demo",
		boardID:     "kanban-demo-htmx",
		addColumn:   "backlog",
		moveCard:    "kb-1",
		moveColumn:  "progress",
		addedCardID: "kb-new-0",
	})

	server.FailIfServerErrors(t)
}

// TestKanbanE2EHTTPContractParity runs the SAME contract probes against the
// e2e harness server (kanbanE2EServer) that TestDemoKanbanHTTPContracts runs
// against the real demo binary. One probe table, two implementations —
// divergence fails here naming the side that drifted (TODO #227).
func TestKanbanE2EHTTPContractParity(t *testing.T) {
	srv := kanbanE2EServer(t)
	defer srv.Close()

	runKanbanContractProbes(t, srv.URL, kanbanContractFixture{
		name:        "e2e-parity",
		boardID:     "kb-htmx",
		addColumn:   "todo",
		moveCard:    "e1",
		moveColumn:  "doing",
		addedCardID: "ea0",
	})
}

// kanbanCSRFTokenFromHTML extracts the move form's hidden csrf_token
// value from a rendered board — exactly what a browser submits back.
func kanbanCSRFTokenFromHTML(t *testing.T, impl, html string) string {
	t.Helper()

	const marker = `name="csrf_token" value="`

	_, rest, found := strings.Cut(html, marker)
	if !found {
		t.Fatalf("visualtest[%s]: rendered board has no csrf_token hidden input", impl)
	}

	value, _, terminated := strings.Cut(rest, `"`)
	if !terminated {
		t.Fatalf("visualtest[%s]: csrf_token input value unterminated", impl)
	}

	return value
}
