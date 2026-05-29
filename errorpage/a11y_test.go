package errorpage

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestErrorPageA11y(t *testing.T) {
	t.Parallel()

	t.Run("error page has no role attribute (full-page, not an alert)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:  FamilyRejection,
			Title:   "Not Found",
			Message: "Page does not exist.",
		}))
		utils.AssertNotContains(t, output, `role="alert"`)
	})

	t.Run("error page propagates aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:    FamilyRejection,
			Title:     "Error",
			BaseProps: utils.BaseProps{AriaLabel: "404 Error"},
		}))
		utils.AssertContains(t, output, `aria-label="404 Error"`)
	})

	t.Run("error page propagates custom ID", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:    FamilyTransient,
			BaseProps: utils.BaseProps{ID: "error-page"},
		}))
		utils.AssertContains(t, output, `id="error-page"`)
	})
}

func TestErrorDetailA11y(t *testing.T) {
	t.Parallel()

	t.Run("error detail has role=alert", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorDetail(ErrorDetailProps{
			Family:  FamilyCorruption,
			Title:   "Data Error",
			Message: "Data is corrupted.",
		}))
		utils.AssertContains(t, output, `role="alert"`)
	})

	t.Run("error detail propagates aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorDetail(ErrorDetailProps{
			Family:    FamilyCorruption,
			Title:     "Error",
			BaseProps: utils.BaseProps{AriaLabel: "Detailed error information"},
		}))
		utils.AssertContains(t, output, `aria-label="Detailed error information"`)
	})
}

func TestErrorAlertA11y(t *testing.T) {
	t.Parallel()

	t.Run("error alert has role=alert", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorAlert(ErrorAlertProps{
			Family:  FamilyTransient,
			Title:   "Temporary Error",
			Message: "Try again.",
		}))
		utils.AssertContains(t, output, `role="alert"`)
	})

	t.Run("error alert has dismiss button with aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorAlert(ErrorAlertProps{
			Family:      FamilyInfrastructure,
			Title:       "Service Down",
			Message:     "Try later.",
			Dismissible: true,
		}))
		utils.AssertContains(t, output, `aria-label="Dismiss"`)
		utils.AssertContains(t, output, `data-dismiss="alert"`)
	})

	t.Run("dismissible alert uses nonce on script", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorAlert(ErrorAlertProps{
			Family:      FamilyTransient,
			Title:       "Error",
			Dismissible: true,
			BaseProps:   utils.BaseProps{Nonce: "test-nonce-abc"},
		}))
		utils.AssertContains(t, output, `nonce="test-nonce-abc"`)
	})

	t.Run("non-dismissible alert has no script tag", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorAlert(ErrorAlertProps{
			Family: FamilyTransient,
			Title:  "Error",
		}))
		utils.AssertNotContains(t, output, "<script")
	})
}
