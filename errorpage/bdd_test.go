package errorpage

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// --- ErrorPage Behavior ---

func TestErrorPageUserSeesFullPageError(t *testing.T) {
	t.Parallel()

	t.Run("user sees error page with title and message", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:  FamilyRejection,
			Code:    "page.not_found",
			Title:   "Page not found",
			Message: "The page does not exist.",
		}))
		utils.AssertContains(t, output, "Page not found")
		utils.AssertContains(t, output, "The page does not exist.")
		utils.AssertContains(t, output, "page.not_found")
		utils.AssertContains(t, output, "amber")
	})

	t.Run("user sees transient error with blue styling", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:  FamilyTransient,
			Title:   "Temporary Error",
			Message: "Please try again in a moment.",
		}))
		utils.AssertContains(t, output, "Temporary Error")
		utils.AssertContains(t, output, "blue")
	})

	t.Run("user sees why explanation and fix suggestion", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:  FamilyCorruption,
			Title:   "Data Corrupted",
			Message: "Configuration file is unreadable.",
			Why:     "Some data appears to be damaged. This requires attention.",
			Fix:     "Check the YAML syntax and fix the indentation.",
		}))
		utils.AssertContains(t, output, "Some data appears to be damaged")
		utils.AssertContains(t, output, "Suggested fix:")
		utils.AssertContains(t, output, "Check the YAML syntax")
		utils.AssertContains(t, output, "red")
	})

	t.Run("user sees way out action link", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:     FamilyInfrastructure,
			Title:      "Service Unavailable",
			Message:    "The service is currently down.",
			WayOut:     "Go to Dashboard",
			WayOutHref: "/dashboard",
		}))
		utils.AssertContains(t, output, "Go to Dashboard")
		utils.AssertContains(t, output, "/dashboard")
		utils.AssertContains(t, output, "href=")
	})

	t.Run("user sees context details", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:  FamilyConflict,
			Title:   "Version Conflict",
			Message: "Someone else modified this resource.",
			Context: []ContextPair{
				{Key: "resource", Value: "user/42"},
				{Key: "expected_version", Value: "3"},
				{Key: "actual_version", Value: "5"},
			},
		}))
		utils.AssertContains(t, output, "resource")
		utils.AssertContains(t, output, "user/42")
		utils.AssertContains(t, output, "expected_version")
		utils.AssertContains(t, output, "3")
		utils.AssertContains(t, output, "orange")
	})

	t.Run("user sees cause chain", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:  FamilyTransient,
			Title:   "Database Timeout",
			Message: "Query took too long.",
			CauseChain: []CauseItem{
				{Message: "connection pool exhausted", Code: "db.pool"},
				{Message: "timeout after 30s", Code: "db.timeout"},
			},
		}))
		utils.AssertContains(t, output, "Cause chain")
		utils.AssertContains(t, output, "connection pool exhausted")
		utils.AssertContains(t, output, "db.timeout")
	})

	t.Run("user sees timestamp when enabled", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:        FamilyRejection,
			Title:         "Error",
			Timestamp:     "2026-05-29T12:00:00Z",
			ShowTimestamp: true,
		}))
		utils.AssertContains(t, output, "2026-05-29T12:00:00Z")
	})

	t.Run("user sees minimal error page with defaults", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(DefaultErrorPageProps()))
		utils.AssertContains(t, output, "min-h-screen")
	})

	t.Run("user sees custom ID and aria-label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:    FamilyRejection,
			Title:     "Not Found",
			BaseProps: utils.BaseProps{ID: "error-404", AriaLabel: "404 Error Page"},
		}))
		utils.AssertContains(t, output, `id="error-404"`)
		utils.AssertContains(t, output, `aria-label="404 Error Page"`)
	})
}

// --- ErrorDetail Behavior ---

func TestErrorDetailUserSeesErrorCard(t *testing.T) {
	t.Parallel()

	t.Run("user sees error detail with code and family badge", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorDetail(ErrorDetailProps{
			Family:  FamilyCorruption,
			Code:    "data.parse_failed",
			Title:   "Parse Failed",
			Message: "Invalid YAML at line 42.",
		}))
		utils.AssertContains(t, output, "data.parse_failed")
		utils.AssertContains(t, output, "corruption")
		utils.AssertContains(t, output, "Parse Failed")
		utils.AssertContains(t, output, "Invalid YAML at line 42.")
		utils.AssertContains(t, output, `role="alert"`)
	})

	t.Run("user sees suggested fix in detail card", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorDetail(ErrorDetailProps{
			Family: FamilyRejection,
			Title:  "Invalid Input",
			Fix:    "Check your email format and try again.",
		}))
		utils.AssertContains(t, output, "Suggested fix:")
		utils.AssertContains(t, output, "Check your email format")
	})

	t.Run("user sees context table in detail card", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorDetail(ErrorDetailProps{
			Family: FamilyTransient,
			Title:  "DB Error",
			Context: []ContextPair{
				{Key: "host", Value: "db.internal"},
				{Key: "port", Value: "5432"},
			},
		}))
		utils.AssertContains(t, output, "host")
		utils.AssertContains(t, output, "db.internal")
		utils.AssertContains(t, output, "port")
		utils.AssertContains(t, output, "5432")
	})

	t.Run("user sees cause chain in detail card", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorDetail(ErrorDetailProps{
			Family: FamilyInfrastructure,
			Title:  "Service Down",
			CauseChain: []CauseItem{
				{Message: "dns resolution failed"},
				{Message: "connection refused", Code: "net.refused"},
			},
		}))
		utils.AssertContains(t, output, "Caused by")
		utils.AssertContains(t, output, "dns resolution failed")
		utils.AssertContains(t, output, "net.refused")
	})

	t.Run("user sees timestamp in detail card footer", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorDetail(ErrorDetailProps{
			Family:    FamilyConflict,
			Title:     "Conflict",
			Timestamp: "2026-05-29T14:30:00Z",
		}))
		utils.AssertContains(t, output, "2026-05-29T14:30:00Z")
	})

	t.Run("minimal detail card renders without optional fields", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorDetail(DefaultErrorDetailProps()))
		utils.AssertContains(t, output, `role="alert"`)
	})
}

// --- ErrorAlert Behavior ---

func TestErrorAlertUserSeesFamilyAlert(t *testing.T) {
	t.Parallel()

	t.Run("user sees alert with family badge", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorAlert(ErrorAlertProps{
			Family:  FamilyTransient,
			Title:   "Temporary Error",
			Message: "Please try again shortly.",
		}))
		utils.AssertContains(t, output, "Temporary Error")
		utils.AssertContains(t, output, "Please try again shortly.")
		utils.AssertContains(t, output, "transient")
		utils.AssertContains(t, output, `role="alert"`)
	})

	t.Run("user sees fix suggestion in alert", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorAlert(ErrorAlertProps{
			Family:  FamilyRejection,
			Title:   "Bad Request",
			Message: "Invalid email format.",
			Fix:     "Enter a valid email address.",
		}))
		utils.AssertContains(t, output, "Enter a valid email address.")
	})

	t.Run("user can dismiss dismissible alert", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ErrorAlert(ErrorAlertProps{
			Family:      FamilyInfrastructure,
			Title:       "Service Down",
			Message:     "Try again later.",
			Dismissible: true,
		}))
		utils.AssertContains(t, output, `aria-label="Dismiss"`)
		utils.AssertContains(t, output, `data-dismiss="alert"`)
	})

	t.Run("each family has distinct styling", func(t *testing.T) {
		t.Parallel()

		for _, tt := range []struct {
			family    Family
			wantColor string
		}{
			{FamilyRejection, "amber"},
			{FamilyConflict, "orange"},
			{FamilyTransient, "blue"},
			{FamilyCorruption, "red"},
			{FamilyInfrastructure, "gray"},
		} {
			t.Run(string(tt.family), func(t *testing.T) {
				t.Parallel()
				output := utils.Render(t, ErrorAlert(ErrorAlertProps{
					Family:  tt.family,
					Title:   "Test",
					Message: "Test message",
				}))
				utils.AssertContains(t, output, tt.wantColor)
			})
		}
	})
}

// --- Family Style Logic ---

func TestFamilyStyleLogic(t *testing.T) {
	t.Parallel()

	t.Run("all five families are valid", func(t *testing.T) {
		t.Parallel()

		families := []Family{FamilyRejection, FamilyConflict, FamilyTransient, FamilyCorruption, FamilyInfrastructure}
		for _, f := range families {
			if !FamilyIsValid(f) {
				t.Errorf("expected %q to be valid", f)
			}
		}
	})

	t.Run("unknown family is not valid", func(t *testing.T) {
		t.Parallel()

		if FamilyIsValid("unknown") {
			t.Error("expected unknown family to be invalid")
		}
	})

	t.Run("unknown family gets default icon", func(t *testing.T) {
		t.Parallel()

		icon := FamilyIcon("nonexistent")
		if icon == "" {
			t.Error("expected default icon for unknown family")
		}
	})

	t.Run("each family has a unique color scheme", func(t *testing.T) {
		t.Parallel()

		seen := map[string]Family{}

		for _, f := range []Family{FamilyRejection, FamilyConflict, FamilyTransient, FamilyCorruption, FamilyInfrastructure} {
			style := lookupFamilyStyle(f)
			if existing, ok := seen[style.BG]; ok {
				t.Errorf("families %q and %q share background color %q", existing, f, style.BG)
			}

			seen[style.BG] = f
		}
	})
}
