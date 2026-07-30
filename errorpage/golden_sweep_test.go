package errorpage

import (
	"testing"

	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

// Golden sweep for errorpage components that previously lacked golden tests.

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
	})
}
