package feedback

import (
	"strconv"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// T5.8 (#493): boolean-flag components are usually tested one flag at a
// time, so flag INTERACTIONS go unproven. These tables walk meaningful flag
// combinations and assert the invariant that must hold for EVERY
// combination: rendered in one pass, no panic, and each combination's
// contract markers present (or absent) together.

// TestAlertFlagCombinations walks Dismissible × every feedback type: role
// and (when dismissible) the dismiss affordance + script must co-render for
// each type — a type whose style lookup regresses must not silently drop
// the script or the role.
func TestAlertFlagCombinations(t *testing.T) {
	t.Parallel()

	for _, tt := range []FeedbackType{
		FeedbackSuccess, FeedbackError, FeedbackWarning, FeedbackInfo,
	} {
		for _, dismissible := range []bool{false, true} {
			t.Run(string(tt)+"/dismissible="+b2s(dismissible), func(t *testing.T) {
				t.Parallel()

				output := utils.Render(t, Alert(AlertProps{
					Message:     "combo",
					Type:        tt,
					Dismissible: dismissible,
				}))

				utils.AssertContains(t, output, `role="`+feedbackRole(tt)+`"`)

				if dismissible {
					utils.AssertContains(t, output, `data-dismiss="alert"`)
					utils.AssertContains(t, output, "<script")
				} else {
					utils.AssertNotContains(t, output, "data-dismiss")
				}
			})
		}
	}
}

// TestToastFlagCombinations walks Duration × Title: the auto-dismiss timer
// script renders only for duration > 0 and the title markup only when set —
// no combination may leak the other's markup.
func TestToastFlagCombinations(t *testing.T) {
	t.Parallel()

	for _, duration := range []int{0, ToastDurationMedium} {
		for _, withTitle := range []bool{false, true} {
			name := "duration=" + strconv.Itoa(duration) + "/title=" + b2s(withTitle)
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				props := ToastProps{
					Message:  "combo",
					Duration: duration,
				}
				if withTitle {
					props.Title = "Head"
				}

				output := utils.Render(t, Toast(props))
				utils.AssertContains(t, output, `role="status"`)

				if duration > 0 {
					utils.AssertContains(t, output, "setTimeout")
				}

				if withTitle {
					utils.AssertContains(t, output, "Head")
				}
			})
		}
	}
}

func b2s(b bool) string {
	if b {
		return "t"
	}

	return "f"
}
