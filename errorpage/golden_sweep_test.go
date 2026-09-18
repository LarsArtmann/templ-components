package errorpage

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

// Golden sweep for errorpage components that previously lacked golden tests.

func TestGoldenSweepErrorPage(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "error_page_full", HTML: utils.Render(t, ErrorPage(ErrorPageProps{
			Family:     FamilyTransient,
			StatusCode: 503,
			Code:       CodeUnavailable,
			Title:      "Service temporarily unavailable",
			Message:    "We're performing maintenance or experiencing high traffic.",
			Why:        "This is a temporary issue. No data was lost.",
			Fix:        "Wait a moment and refresh the page.",
			WayOut:     "Retry",
			WayOutHref: "/",
			Context: []ContextPair{
				{Key: "region", Value: "eu-central-1"},
				{Key: "request_id", Value: "req_8fk2m1"},
			},
			CauseChain:    []CauseItem{{Message: "connection pool exhausted", Code: "db.pool"}},
			Timestamp:     "2026-09-17T12:00:00Z",
			Trace:         "trc_9f3a1c2d",
			ShowTimestamp: true,
		}))},
		{Name: "error_page_minimal", HTML: utils.Render(t, ErrorPage(ErrorPageProps{
			Family:  FamilyInfrastructure,
			Code:    CodeInternalError,
			Title:   titleInternalError,
			Message: msgInternalUnexpected,
			WayOut:  msgGoBack,
		}))},
		{Name: "error_page_secondary_action", HTML: utils.Render(t, ErrorPage(ErrorPageProps{
			Family:               FamilyTransient,
			StatusCode:           503,
			Title:                "Service temporarily unavailable",
			Message:              "We're performing maintenance or experiencing high traffic.",
			WayOut:               "Retry",
			WayOutHref:           "/",
			SecondaryWayOut:      "Contact support",
			SecondaryWayOutHref:  "mailto:support@example.com",
			ShowTimestamp:        true,
		}))},
		{Name: "error_page_secondary_go_back", HTML: utils.Render(t, ErrorPage(ErrorPageProps{
			Family:          FamilyInfrastructure,
			Title:           "Something went wrong",
			Message:         "An unexpected error occurred.",
			WayOut:          "Go home",
			WayOutHref:      "/",
			SecondaryWayOut: "Go back",
		}))},
	})
}

func TestGoldenSweepErrorAlert(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "error_alert_rejection", HTML: utils.Render(t, ErrorAlert(ErrorAlertProps{
			Family:  FamilyRejection,
			Title:   "Invalid Input",
			Message: "The email address format is invalid.",
			Fix:     "Enter a valid email address like name@example.com.",
		}))},
		{Name: "error_alert_transient", HTML: utils.Render(t, ErrorAlert(ErrorAlertProps{
			Family:  FamilyTransient,
			Title:   "Service Unavailable",
			Message: "The database is temporarily unreachable.",
		}))},
	})
}

func TestGoldenSweepErrorDetail(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "error_detail_full", HTML: utils.Render(t, ErrorDetail(ErrorDetailProps{
			Family:  FamilyCorruption,
			Code:    "config.parse_failed",
			Title:   "Configuration Parse Error",
			Message: "config.yaml has invalid syntax at line 42.",
			Fix:     "Check the YAML syntax — the indentation appears incorrect.",
			Context: []ContextPair{
				{Key: "file", Value: "config.yaml"},
				{Key: "line", Value: "42"},
				{Key: "column", Value: "8"},
			},
			CauseChain: []CauseItem{
				{Message: "yaml: found character that cannot start any token", Code: "yaml_scanner_error"},
			},
			Timestamp: "2026-07-30T12:00:00Z",
		}))},
		{Name: "error_detail_minimal", HTML: utils.Render(t, ErrorDetail(ErrorDetailProps{
			Family:  FamilyConflict,
			Title:   "Version Conflict",
			Message: "The record was modified by another user.",
		}))},
		{Name: "error_detail_neutral", HTML: utils.Render(t, ErrorDetail(ErrorDetailProps{
			Family:  FamilyCorruption,
			Code:    "config.parse_failed",
			Title:   "Configuration Parse Error",
			Message: "config.yaml has invalid syntax at line 42.",
			Fix:     "Check the YAML syntax — the indentation appears incorrect.",
			Context: []ContextPair{
				{Key: "file", Value: "config.yaml"},
				{Key: "line", Value: "42"},
				{Key: "column", Value: "8"},
			},
			Timestamp: "2026-07-30T12:00:00Z",
			Variant:   ErrorDetailNeutral,
		}))},
	})
}

// TestGoldenHandlerHTMLShell pins the ErrorHandler HTMLShell document path —
// the wire format a consumer's http.Server actually emits when configured
// with HTMLShell: true (full <!doctype html> document, lang, title, embedded
// error page). The timestamp is pinned via Override so the golden is
// deterministic.
func TestGoldenHandlerHTMLShell(t *testing.T) {
	t.Parallel()

	handler := ErrorHandler(errors.New("boom"), ErrorHandlerConfig{
		HTMLShell: true,
		Nonce:     "test-nonce",
		Override: func(_ error, props ErrorPageProps) *ErrorPageProps {
			props.Timestamp = "2026-09-17T12:00:00Z"
			props.Title = "Something went wrong"
			props.Message = "An unexpected error occurred."
			props.Code = CodeInternalError

			return &props
		},
	})

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/test", nil))

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "handler_htmlshell", HTML: rec.Body.String()},
	})
}

// TestGoldenSweepErrorFamilyMatrix pins the full six-family matrix —
// every go-error-family renders through ErrorAlert (rejection, conflict,
// transient, corruption, infrastructure, orchestration) so a family
// color/icon/text regression is caught for ANY family, not just the two that
// previously had goldens. ErrorDetail gains the missing infrastructure
// family too.
func TestGoldenSweepErrorFamilyMatrix(t *testing.T) {
	t.Parallel()

	families := []struct {
		name   string
		family Family
		title  string
	}{
		{"rejection", FamilyRejection, "Invalid Input"},
		{"conflict", FamilyConflict, "Version Conflict"},
		{"transient", FamilyTransient, "Service Unavailable"},
		{"corruption", FamilyCorruption, "Data Corruption"},
		{"infrastructure", FamilyInfrastructure, "Infrastructure Failure"},
		{"orchestration", FamilyOrchestration, "Orchestration Failure"},
	}

	alerts := make([]golden.Snapshot, 0, len(families))

	for _, f := range families {
		family := f

		alerts = append(alerts, golden.Snapshot{
			Name: "error_alert_family_" + family.name,
			HTML: utils.Render(t, ErrorAlert(ErrorAlertProps{
				Family:  family.family,
				Title:   family.title,
				Message: "Matrix message for " + string(family.family) + ".",
				Fix:     "Matrix fix for " + string(family.family) + ".",
			})),
		})
	}

	golden.AssertSnapshots(t, alerts)

	golden.Assert(t, "error_detail_infrastructure", utils.Render(t, ErrorDetail(ErrorDetailProps{
		Family:    FamilyInfrastructure,
		Code:      Code("infra.upstream_unreachable"),
		Title:     "Upstream Unreachable",
		Message:   "The upstream dependency did not respond in time.",
		Fix:       "Retry shortly; if it persists, check the status page.",
		Timestamp: "2026-09-09T12:00:00Z",
	})))
}
