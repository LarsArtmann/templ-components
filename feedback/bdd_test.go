// Package feedback provides behavior-driven tests for feedback components.
// These tests verify end-user-facing behavior: what the user sees and experiences.
package feedback

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func renderLoadingOverlayWithProgress(t *testing.T, message string, progress int) string {
	t.Helper()
	return utils.Render(t, LoadingOverlay(LoadingOverlayProps{
		Message:      message,
		ShowProgress: true,
		Progress:     progress,
	}))
}

func assertLoadingOverlayProgress(t *testing.T, message string, progress int, wantPercent string) {
	t.Helper()
	output := renderLoadingOverlayWithProgress(t, message, progress)
	utils.AssertContains(t, output, message)
	utils.AssertContains(t, output, wantPercent)
}

// --- Alert Behavior ---

func TestAlertUserReceivesImportantMessages(t *testing.T) {
	t.Parallel()

	t.Run("user sees alert with title and message", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Alert(AlertProps{
			Title:   "Warning",
			Message: "This action cannot be undone.",
			Type:    AlertWarning,
		}))
		utils.AssertContains(t, output, "Warning")
		utils.AssertContains(t, output, "This action cannot be undone.")
		utils.AssertContains(t, output, `role="alert"`)
	})

	t.Run("user sees success alert with green styling", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Alert(AlertProps{
			Title:   "Success",
			Message: "Your changes have been saved.",
			Type:    AlertSuccess,
		}))
		utils.AssertContains(t, output, "Success")
		utils.AssertContains(t, output, "green")
	})

	t.Run("user can dismiss a dismissible alert", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Alert(AlertProps{
			Title:       "Info",
			Message:     "New features available.",
			Type:        AlertInfo,
			Dismissible: true,
		}))
		utils.AssertContains(t, output, "New features available.")
		utils.AssertContains(t, output, `data-dismiss="alert"`)
	})
}

// --- Toast Behavior ---

func TestToastUserGetsNonIntrusiveNotifications(t *testing.T) {
	t.Parallel()

	t.Run("user sees toast with message", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Toast(ToastProps{
			Message: "Item saved successfully!",
			Type:    ToastSuccess,
		}))
		utils.AssertContains(t, output, "Item saved successfully!")
	})

	t.Run("user sees toast with title and auto-dismiss timing", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Toast(ToastProps{
			Title:    "Saved",
			Message:  "Your changes are live.",
			Type:     ToastSuccess,
			Duration: ToastDurationMedium,
		}))
		utils.AssertContains(t, output, "Saved")
		utils.AssertContains(t, output, "Your changes are live.")
	})
}

// --- Spinner / Loading Behavior ---

func TestSpinnerUserSeesLoadingProgress(t *testing.T) {
	t.Parallel()

	t.Run("user sees animated spinner", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Spinner(SpinnerProps{Size: SpinnerMD, Color: "text-blue-600"}))
		utils.AssertContains(t, output, "animate-spin")
		utils.AssertContains(t, output, "text-blue-600")
	})

	t.Run("user sees full-screen loading overlay", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadingOverlay(LoadingOverlayProps{Message: "Loading data..."}))
		utils.AssertContains(t, output, "Loading data...")
		utils.AssertContains(t, output, "animate-spin")
	})

	t.Run("user sees loading overlay with progress", func(t *testing.T) {
		t.Parallel()
		assertLoadingOverlayProgress(t, "Uploading...", 65, "65%")
	})
}

// --- Skeleton Behavior ---

func TestSkeletonUserSeesContentPlaceholder(t *testing.T) {
	t.Parallel()

	t.Run("user sees pulsing skeleton placeholder", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Skeleton(SkeletonText))
		utils.AssertContains(t, output, "animate-pulse")
	})

	t.Run("user sees skeleton group with multiple variants", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SkeletonGroup([]SkeletonVariant{
			SkeletonTitle,
			SkeletonText,
			SkeletonAvatar,
		}))
		utils.AssertContains(t, output, "animate-pulse")
	})
}

// --- Progress Bar Behavior ---

func TestProgressBarUserSeesCompletion(t *testing.T) {
	t.Parallel()

	t.Run("user sees progress at 50 percent", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Current: 50,
			Total:   100,
		}))
		utils.AssertContains(t, output, "50%")
		utils.AssertContains(t, output, `role="progressbar"`)
	})

	t.Run("user sees progress with label", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Current:   3,
			Total:     5,
			Label:     "Steps completed",
			ShowLabel: true,
		}))
		utils.AssertContains(t, output, "Steps completed")
	})

	t.Run("progress bar clamps overflow to 100 percent", func(t *testing.T) {
		t.Parallel()
		for _, tt := range []struct {
			name    string
			current int
			total   int
			extra   string
		}{
			{"150/100", 150, 100, ""},
			{"200/100 with aria-valuenow", 200, 100, `aria-valuenow="200"`},
		} {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				output := utils.Render(t, ProgressBar(ProgressBarProps{
					Current: tt.current,
					Total:   tt.total,
				}))
				utils.AssertContains(t, output, "100%")
				utils.AssertNotContains(t, output, "150%")
				if tt.extra != "" {
					utils.AssertContains(t, output, tt.extra)
				}
			})
		}
	})
}

// --- Inline Messages Behavior ---

func TestInlineMessagesUserSeesFieldFeedback(t *testing.T) {
	t.Parallel()

	t.Run("user sees inline error message", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, InlineError("Email is required"))
		utils.AssertContains(t, output, "Email is required")
		utils.AssertContains(t, output, `role="alert"`)
		utils.AssertContains(t, output, "text-red-600")
	})

	t.Run("user sees inline success message", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, InlineSuccess("Username is available"))
		utils.AssertContains(t, output, "Username is available")
	})
}

// --- SkeletonCardGrid Behavior ---

func TestSkeletonCardGridUserSeesLoadingState(t *testing.T) {
	t.Parallel()

	t.Run("user sees 6 skeleton cards while loading", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SkeletonCardGrid(6))
		utils.AssertContains(t, output, `role="status"`)
		utils.AssertContains(t, output, "lg:grid-cols-3")
	})

	t.Run("zero count falls back to single placeholder", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, SkeletonCardGrid(0))
		if got := strings.Count(output, "h-48"); got != 1 {
			t.Errorf("expected 1 fallback card, got %d", got)
		}
	})
}
