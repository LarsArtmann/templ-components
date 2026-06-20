package errorpage

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestErrorPageFullCoverage(t *testing.T) {
	t.Parallel()

	t.Run("all fields populated with WayOutHref and ShowTimestamp", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			BaseProps: utils.BaseProps{
				ID:        "err-page",
				Class:     "custom-class",
				AriaLabel: "Error page",
				Nonce:     "nonce-123",
				Attrs:     templ.Attributes{"data-test": "true"},
			},
			Family:        FamilyConflict,
			Code:          "ERR_CONFLICT_409",
			Title:         "Resource Conflict",
			Message:       "The resource already exists",
			Why:           "Duplicate key constraint violated",
			Fix:           "Use a unique identifier",
			WayOut:        "Back to list",
			WayOutHref:    "/dashboard",
			ShowTimestamp: true,
			Timestamp:     "2024-06-20T12:00:00Z",
			Context: []ContextPair{
				{Key: "resource", Value: "users/42"},
				{Key: "field", Value: "email"},
			},
			CauseChain: []CauseItem{
				{Message: "unique constraint failed", Code: "DB_001"},
				{Message: "insert blocked", Code: ""},
			},
		}))
		utils.AssertContains(t, output, `id="err-page"`)
		utils.AssertContains(t, output, "custom-class")
		utils.AssertContains(t, output, `aria-label="Error page"`)
		utils.AssertContains(t, output, `data-test="true"`)
		utils.AssertContains(t, output, `href="/dashboard"`)
		utils.AssertContains(t, output, "Back to list")
		utils.AssertContains(t, output, "2024-06-20T12:00:00Z")
		utils.AssertContains(t, output, "DB_001")
		utils.AssertContains(t, output, "Cause chain")
	})

	t.Run("WayOut without WayOutHref renders button and script", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			BaseProps: utils.BaseProps{Nonce: "n1"},
			Family:    FamilyRejection,
			Title:     "Rejected",
			WayOut:    "Return to list",
		}))
		utils.AssertContains(t, output, "data-tc-go-back")
		utils.AssertContains(t, output, "Return to list")
		utils.AssertContains(t, output, "history.back()")
	})

	t.Run("WayOutHref without WayOut renders link with default label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:     FamilyInfrastructure,
			Title:      "Infra Error",
			WayOutHref: "/home",
		}))
		utils.AssertContains(t, output, `href="/home"`)
		utils.AssertContains(t, output, "back")
	})

	t.Run("ShowTimestamp without timestamp value", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:        FamilyCorruption,
			Title:         "Corrupted",
			ShowTimestamp: true,
		}))
		utils.AssertContains(t, output, "Corrupted")
	})

	t.Run("all families render without error", func(t *testing.T) {
		t.Parallel()
		for _, family := range []Family{
			FamilyRejection,
			FamilyConflict,
			FamilyTransient,
			FamilyCorruption,
			FamilyInfrastructure,
		} {
			output := utils.Render(t, ErrorPage(ErrorPageProps{
				Family: family,
				Title:  string(family),
				Code:   "TEST",
			}))
			utils.AssertContains(t, output, string(family))
			utils.AssertContains(t, output, "TEST")
		}
	})
}

func TestErrorDetailFullCoverage(t *testing.T) {
	t.Parallel()

	t.Run("with all fields including BaseProps and cause codes", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorDetail(ErrorDetailProps{
			BaseProps: utils.BaseProps{
				ID:        "detail-1",
				Class:     "detail-class",
				AriaLabel: "Error detail",
				Nonce:     "nonce-x",
			},
			Family:  FamilyCorruption,
			Code:    "CORRUPT_500",
			Title:   "Data Corruption",
			Message: "Records are damaged",
			Fix:     "Restore from backup",
			Context: []ContextPair{
				{Key: "table", Value: "orders"},
				{Key: "rows", Value: "1234"},
			},
			CauseChain: []CauseItem{
				{Message: "checksum mismatch", Code: "CHK001"},
				{Message: "disk read error", Code: "DISK002"},
			},
			Timestamp: "2024-06-20T08:00:00Z",
		}))
		utils.AssertContains(t, output, `id="detail-1"`)
		utils.AssertContains(t, output, "detail-class")
		utils.AssertContains(t, output, `aria-label="Error detail"`)
		utils.AssertContains(t, output, "CORRUPT_500")
		utils.AssertContains(t, output, "Restore from backup")
		utils.AssertContains(t, output, "CHK001")
		utils.AssertContains(t, output, "DISK002")
		utils.AssertContains(t, output, "2024-06-20T08:00:00Z")
	})

	t.Run("without optional fields", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorDetail(ErrorDetailProps{
			Family: FamilyTransient,
			Title:  "Minimal",
		}))
		utils.AssertContains(t, output, "Minimal")
		utils.AssertContains(t, output, "transient")
	})
}

func TestErrorAlertFullCoverage(t *testing.T) {
	t.Parallel()

	t.Run("with BaseProps", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorAlert(ErrorAlertProps{
			BaseProps: utils.BaseProps{
				ID:        "alert-1",
				Class:     "alert-class",
				AriaLabel: "Error alert",
				Attrs:     templ.Attributes{"data-severity": "high"},
			},
			Family:      FamilyCorruption,
			Title:       "Critical Alert",
			Message:     "Immediate action required",
			Fix:         "Restart the service",
			Dismissible: true,
		}))
		utils.AssertContains(t, output, `id="alert-1"`)
		utils.AssertContains(t, output, "alert-class")
		utils.AssertContains(t, output, `aria-label="Error alert"`)
		utils.AssertContains(t, output, `data-severity="high"`)
		utils.AssertContains(t, output, "data-dismiss")
		utils.AssertContains(t, output, "Restart the service")
	})
}
