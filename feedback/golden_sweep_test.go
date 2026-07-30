package feedback

import (
	"testing"

	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

// Golden sweep for feedback components that previously lacked golden tests.

func TestGoldenSweepInlineError(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{"inline_error", utils.Render(t, InlineError("Email is required"))},
	})
}

func TestGoldenSweepInlineSuccess(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{"inline_success", utils.Render(t, InlineSuccess("Username is available"))},
	})
}

func TestGoldenSweepSkeletonGroup(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{"skeleton_group", utils.Render(t, SkeletonGroup([]SkeletonVariant{
			SkeletonTableRow, SkeletonTableRow, SkeletonTableRow,
		}))},
	})
}

func TestGoldenSweepToastContainer(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{"toast_container", utils.Render(t, ToastContainer("test-nonce"))},
	})
}
