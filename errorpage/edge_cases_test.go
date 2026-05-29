package errorpage

import (
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
