package errorpage

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
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

// TestGoldenSweepErrorFamilyMatrix pins the full five-family matrix —
// every go-error-family renders through ErrorAlert (rejection, conflict,
// transient, corruption, infrastructure) so a family color/icon/text
// regression is caught for ANY family, not just the two that previously
// had goldens. ErrorDetail gains the missing infrastructure family too.
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
