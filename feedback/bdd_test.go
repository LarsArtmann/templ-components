// Package feedback provides behavior-driven tests for feedback components.
// These tests verify end-user-facing behavior: what the user sees and experiences.
package feedback

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

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
		output := utils.Render(t, Spinner(SpinnerMD, "text-blue-600"))
		utils.AssertContains(t, output, "animate-spin")
		utils.AssertContains(t, output, "text-blue-600")
	})

	t.Run("user sees full-screen loading overlay", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadingOverlay("Loading data...", false, 0))
		utils.AssertContains(t, output, "Loading data...")
		utils.AssertContains(t, output, "animate-spin")
	})

	t.Run("user sees loading overlay with progress", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, LoadingOverlay("Uploading...", true, 65))
		utils.AssertContains(t, output, "Uploading...")
		utils.AssertContains(t, output, "65%")
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
		output := utils.Render(t, ProgressBar(ProgressBarProps{
			Current: 150,
			Total:   100,
		}))
		utils.AssertContains(t, output, "100%")
		utils.AssertNotContains(t, output, "150%")
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
	})

	t.Run("user sees inline success message", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, InlineSuccess("Username is available"))
		utils.AssertContains(t, output, "Username is available")
	})
}
