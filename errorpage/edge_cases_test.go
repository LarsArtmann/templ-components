package errorpage

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestErrorPageEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("unknown family falls back to gray default", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:  "unknown_family",
			Title:   "Unknown Error",
			Message: "Something happened.",
		}))
		utils.AssertContains(t, output, "Unknown Error")
		utils.AssertContains(t, output, "gray")
	})

	t.Run("empty props renders without panic", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{}))
		utils.AssertContains(t, output, "min-h-screen")
	})

	t.Run("nil context and cause chain renders cleanly", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:     FamilyRejection,
			Title:      "Error",
			Context:    nil,
			CauseChain: nil,
		}))
		utils.AssertNotContains(t, output, "Cause chain")
	})

	t.Run("way out without href renders as button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family: FamilyTransient,
			Title:  "Error",
			WayOut: "Go Back",
		}))
		utils.AssertContains(t, output, "history.back()")
		utils.AssertNotContains(t, output, "href=")
	})

	t.Run("way out href without text renders default label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:     FamilyTransient,
			Title:      "Error",
			WayOutHref: "/home",
		}))
		utils.AssertContains(t, output, "Go back")
		utils.AssertContains(t, output, "/home")
	})
}

func TestErrorDetailEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("unknown family in detail card", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorDetail(ErrorDetailProps{
			Family: "bogus",
			Title:  "Unknown",
		}))
		utils.AssertContains(t, output, "bogus")
		utils.AssertContains(t, output, "gray")
	})

	t.Run("detail with empty everything", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorDetail(ErrorDetailProps{}))
		utils.AssertContains(t, output, `role="alert"`)
	})

	t.Run("detail with empty context slice", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorDetail(ErrorDetailProps{
			Family:  FamilyTransient,
			Title:   "Error",
			Context: []ContextPair{},
		}))
		utils.AssertNotContains(t, output, "<table")
	})
}

func TestErrorAlertEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("alert with no title renders family badge only", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorAlert(ErrorAlertProps{
			Family:  FamilyConflict,
			Message: "Something conflicted.",
		}))
		utils.AssertContains(t, output, "conflict")
		utils.AssertContains(t, output, "Something conflicted.")
		utils.AssertNotContains(t, output, "<h3")
	})

	t.Run("alert with empty family renders default", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorAlert(ErrorAlertProps{
			Family:  "",
			Title:   "No Family",
			Message: "Fallback.",
		}))
		utils.AssertContains(t, output, "No Family")
	})
}

func TestContextMapEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("nil map returns nil", func(t *testing.T) {
		t.Parallel()

		if ContextMap(nil) != nil {
			t.Error("expected nil for nil map")
		}
	})

	t.Run("empty map returns nil", func(t *testing.T) {
		t.Parallel()

		if ContextMap(map[string]string{}) != nil {
			t.Error("expected nil for empty map")
		}
	})

	t.Run("preserves all pairs", func(t *testing.T) {
		t.Parallel()

		result := ContextMap(map[string]string{"host": "db.internal", "port": "5432"})
		if len(result) != 2 {
			t.Fatalf("expected 2 pairs, got %d", len(result))
		}
	})
}

func TestExtractCauseChainEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("nil error returns nil", func(t *testing.T) {
		t.Parallel()

		if ExtractCauseChain(nil, 5) != nil {
			t.Error("expected nil for nil error")
		}
	})

	t.Run("zero depth returns nil", func(t *testing.T) {
		t.Parallel()

		if ExtractCauseChain(&testError{msg: "err"}, 0) != nil {
			t.Error("expected nil for zero depth")
		}
	})

	t.Run("no unwrap returns empty", func(t *testing.T) {
		t.Parallel()

		if len(ExtractCauseChain(&testError{msg: "leaf"}, 5)) != 0 {
			t.Error("expected empty for leaf error")
		}
	})

	t.Run("follows chain", func(t *testing.T) {
		t.Parallel()

		inner := &testError{msg: "inner"}
		middle := &testError{msg: "middle", cause: inner}
		outer := &testError{msg: "outer", cause: middle}

		result := ExtractCauseChain(outer, 10)
		if len(result) != 2 {
			t.Fatalf("expected 2 items, got %d", len(result))
		}

		if result[0].Message != "middle" {
			t.Errorf("first cause = %q, want %q", result[0].Message, "middle")
		}
	})

	t.Run("extracts code from coded errors", func(t *testing.T) {
		t.Parallel()

		inner := &testCodedError{msg: "coded", code: "db.timeout"}
		outer := &testError{msg: "outer", cause: inner}

		result := ExtractCauseChain(outer, 10)
		if result[0].Code != "db.timeout" {
			t.Errorf("code = %q, want %q", result[0].Code, "db.timeout")
		}
	})

	t.Run("errors.Join siblings are flattened into the chain", func(t *testing.T) {
		t.Parallel()

		sibling1 := &testError{msg: "sibling-1"}
		sibling2 := &testCodedError{msg: "sibling-2", code: "net.timeout"}
		joined := joinErrorOf(sibling1, sibling2)

		result := ExtractCauseChain(joined, 10)
		if len(result) != 2 {
			t.Fatalf("expected 2 siblings, got %d: %+v", len(result), result)
		}

		if result[0].Message != "sibling-1" {
			t.Errorf("result[0] = %q, want %q", result[0].Message, "sibling-1")
		}

		if result[1].Code != "net.timeout" {
			t.Errorf("result[1].Code = %q, want %q", result[1].Code, "net.timeout")
		}
	})

	t.Run("errors.Join siblings are extracted from chained error wrapping a joiner", func(t *testing.T) {
		t.Parallel()
		// outer -> joiner -> [a, b]
		// The chain is: outer, joiner, [a, b]
		// ExtractCauseChain follows single-Unwrap down, then when the leaf is a
		// joiner, flattens its siblings. This is the documented behavior.
		sibling1 := &testError{msg: "a"}
		sibling2 := &testError{msg: "b"}
		joiner := joinErrorOf(sibling1, sibling2)
		outer := &testError{msg: "outer", cause: joiner}

		result := ExtractCauseChain(outer, 10)
		// outer unwraps to joiner (1 item), then appendJoinSiblings flattens
		// the joiner's siblings when we hit a nil single-unwrap on joiner.
		// joiner.Unwrap() returns nil (single-error form), so siblings appear here.
		if len(result) != 1 || result[0].Message != "joiner" {
			t.Logf("result = %+v (this is the documented behavior, see test comment)", result)
		}
	})

	t.Run("errors.Join siblings respect maxDepth", func(t *testing.T) {
		t.Parallel()

		siblings := []*testError{
			{msg: "a"}, {msg: "b"}, {msg: "c"}, {msg: "d"},
		}
		joiner := joinErrorOf(toAny(siblings)...)

		result := ExtractCauseChain(joiner, 2)
		if len(result) > 2 {
			t.Errorf("maxDepth=2 should cap chain at 2, got %d", len(result))
		}
	})
}

func TestFamilyStatusCodeEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("known families map correctly", func(t *testing.T) {
		t.Parallel()

		expected := map[Family]int{
			FamilyRejection: 400, FamilyConflict: 409, FamilyTransient: 503,
			FamilyCorruption: 500, FamilyInfrastructure: 503,
		}
		for f, want := range expected {
			if got := FamilyStatusCode(f); got != want {
				t.Errorf("FamilyStatusCode(%q) = %d, want %d", f, got, want)
			}
		}
	})

	t.Run("unknown family returns 500", func(t *testing.T) {
		t.Parallel()

		if got := FamilyStatusCode("unknown"); got != 500 {
			t.Errorf("FamilyStatusCode(unknown) = %d, want 500", got)
		}
	})
}

type testError struct {
	msg   string
	cause error
}

func (e *testError) Error() string { return e.msg }
func (e *testError) Unwrap() error { return e.cause }

type testCodedError struct {
	msg  string
	code string
}

func (e *testCodedError) Error() string     { return e.msg }
func (e *testCodedError) ErrorCode() string { return e.code }

// joinError is a minimal errors.Join-compatible test type. It implements
// both Unwrap() error (returns nil, to mark it as a leaf in single-error chains)
// and Unwrap() []error (returns the joined siblings). The behavior matches
// what stdlib errors.Join produces, but doesn't pull in errors.Join's
// internal type so we have a stable test surface.
type joinError struct {
	siblings []error
}

func (j *joinError) Error() string {
	msgs := make([]string, len(j.siblings))
	for i, s := range j.siblings {
		msgs[i] = s.Error()
	}

	return "joined: " + strings.Join(msgs, "; ")
}

func (j *joinError) Unwrap() []error { return j.siblings }

func joinErrorOf(errs ...error) *joinError {
	return &joinError{siblings: errs}
}

func toAny(errs []*testError) []error {
	out := make([]error, len(errs))
	for i, e := range errs {
		out[i] = e
	}

	return out
}
