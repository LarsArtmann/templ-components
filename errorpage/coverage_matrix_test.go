package errorpage

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRenderCoverageMatrix renders every component across every family and
// the optional-branch combinations, exercising the generated template
// branches for the coverage gate. Renders only — shape assertions live in
// the golden sweep; this test exists to close uncovered template paths
// (compact vs full headers, empty sections, optional chips).
func TestRenderCoverageMatrix(t *testing.T) {
	t.Parallel()

	families := []Family{
		FamilyRejection,
		FamilyConflict,
		FamilyTransient,
		FamilyCorruption,
		FamilyInfrastructure,
		FamilyOrchestration,
	}

	for _, family := range families {
		style := lookupFamilyStyle(family)

		// ErrorAlert across families × dismissible.
		for _, dismissible := range []bool{false, true} {
			if err := ErrorAlert(ErrorAlertProps{
				Family:      family,
				Title:       "Title " + string(family),
				Message:     "Message body.",
				Fix:         "Fix suggestion.",
				Dismissible: dismissible,
			}).Render(context.Background(), io.Discard); err != nil {
				t.Errorf("ErrorAlert(%s, dismissible=%v) render failed: %v", family, dismissible, err)
			}
		}

		// ErrorDetail compact header across families, bare and full.
		if err := ErrorDetail(ErrorDetailProps{Family: family}).Render(context.Background(), io.Discard); err != nil {
			t.Errorf("ErrorDetail(%s) bare render failed: %v", family, err)
		}

		if err := ErrorDetail(ErrorDetailProps{
			Family:     family,
			Code:       CodeInternalError,
			Title:      "Title",
			Message:    "Message body.",
			Fix:        "Fix suggestion.",
			Context:    []ContextPair{{Key: "k", Value: "v"}},
			CauseChain: []CauseItem{{Message: "cause", Code: "c1"}},
			Timestamp:  "2026-09-17T12:00:00Z",
			Trace:      "trc_1",
		}).Render(context.Background(), io.Discard); err != nil {
			t.Errorf("ErrorDetail(%s) full render failed: %v", family, err)
		}

		// The compact errorBody header (title/message/chip combinations).
		for _, title := range []string{"", "Title only"} {
			for _, message := range []string{"", "Message only"} {
				if err := errorBody(style, family, CodeConflict, title, "h3", "mt-1.5 text-sm font-semibold text-gray-900 dark:text-white", message, "mt-1 text-sm").Render(context.Background(), io.Discard); err != nil {
					t.Errorf("errorBody(%s) render failed: %v", family, err)
				}
			}
		}
	}
}

// TestNotFound404RenderVariants exercises the NotFound404 optional branches:
// search form present/absent, links present/absent, go-back off.
func TestNotFound404RenderVariants(t *testing.T) {
	t.Parallel()

	variants := []struct {
		name   string
		mutate func(props NotFound404Props) NotFound404Props
	}{
		{"minimal", func(p NotFound404Props) NotFound404Props {
			p.SearchAction = ""
			p.Links = nil
			p.ShowGoBack = false

			return p
		}},
		{"search-only", func(p NotFound404Props) NotFound404Props {
			p.Links = nil
			p.ShowGoBack = false

			return p
		}},
		{"links-only", func(p NotFound404Props) NotFound404Props {
			p.SearchAction = ""
			p.ShowGoBack = false

			return p
		}},
		{"full", func(p NotFound404Props) NotFound404Props { return p }},
	}

	for _, variant := range variants {
		props := variant.mutate(DefaultNotFound404Props())
		if err := NotFound404(props).Render(context.Background(), io.Discard); err != nil {
			t.Errorf("NotFound404(%s) render failed: %v", variant.name, err)
		}
	}
}

// TestWriteFallbackErrorPath drives the HTMLShell render-failure fallback: a
// canceled request context fails renderShellToBuffer, and the handler must
// still write a minimal body with the right status instead of hanging or
// panicking.
func TestWriteFallbackErrorPath(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)

	handler := ErrorHandler(errTestFallback{}, ErrorHandlerConfig{HTMLShell: true})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

type errTestFallback struct{}

func (errTestFallback) Error() string { return "boom" }

// TestErrorPageBranchCombos samples the ErrorPage optional-branch space
// (chips, diagnostics, action pair, footer, copy) so the generated template
// branches count toward the coverage gate.
func TestErrorPageBranchCombos(t *testing.T) {
	t.Parallel()

	base := func() ErrorPageProps {
		return ErrorPageProps{
			Family:     FamilyRejection,
			StatusCode: 400,
			Code:       CodeBadRequest,
			Title:      "Bad request",
			Message:    "The request could not be understood.",
		}
	}

	combos := []func(p ErrorPageProps) ErrorPageProps{
		func(p ErrorPageProps) ErrorPageProps { return p },
		func(p ErrorPageProps) ErrorPageProps {
			p.Why = "Why text"

			return p
		},
		func(p ErrorPageProps) ErrorPageProps {
			p.Fix = "Fix text"

			return p
		},
		func(p ErrorPageProps) ErrorPageProps {
			p.Context = []ContextPair{{Key: "k", Value: "v"}}

			return p
		},
		func(p ErrorPageProps) ErrorPageProps {
			p.CauseChain = []CauseItem{{Message: "cause"}}

			return p
		},
		func(p ErrorPageProps) ErrorPageProps {
			p.Trace = "trc_1"

			return p
		},
		func(p ErrorPageProps) ErrorPageProps {
			p.ShowTimestamp = true

			return p
		},
		func(p ErrorPageProps) ErrorPageProps {
			p.CopyCode = true

			return p
		},
		func(p ErrorPageProps) ErrorPageProps {
			p.WayOut = "Go back"

			return p
		},
		func(p ErrorPageProps) ErrorPageProps {
			p.WayOut = "Retry"

			p.WayOutHref = "/"

			return p
		},
		func(p ErrorPageProps) ErrorPageProps {
			p.SecondaryWayOut = "Contact support"

			p.SecondaryWayOutHref = "mailto:s@example.com"

			return p
		},
		func(p ErrorPageProps) ErrorPageProps {
			p.WayOutAction = WayOutAction{Text: "Status page", Href: "/status"}

			return p
		},
		func(p ErrorPageProps) ErrorPageProps {
			p.MaxWidth = ErrorMaxWidth2XL

			return p
		},
		func(p ErrorPageProps) ErrorPageProps {
			p.Code = ""

			return p
		},
		func(p ErrorPageProps) ErrorPageProps {
			p.StatusCode = 0

			return p
		},
		func(p ErrorPageProps) ErrorPageProps {
			p.Why = "Why"
			p.Fix = "Fix"
			p.Context = []ContextPair{{Key: "k", Value: "v"}}
			p.CauseChain = []CauseItem{{Message: "cause", Code: "c"}}
			p.Trace = "trc_2"
			p.ShowTimestamp = true
			p.CopyCode = true
			p.WayOut = "Retry"
			p.WayOutHref = "/"
			p.SecondaryWayOut = "Support"
			p.SecondaryWayOutHref = "mailto:s@example.com"
			p.MaxWidth = ErrorMaxWidth4XL

			return p
		},
	}

	for i, mutate := range combos {
		props := mutate(base())
		if err := ErrorPage(props).Render(context.Background(), io.Discard); err != nil {
			t.Errorf("combo %d render failed: %v", i, err)
		}
	}
}
