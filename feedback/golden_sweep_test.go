package feedback

import (
	"testing"

	"github.com/larsartmann/templ-components/utils/golden"
	"github.com/larsartmann/templ-components/utils"
)

// Golden sweep for feedback components that previously lacked golden tests.

func TestGoldenSweepInlineError(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "inline_error", HTML: utils.Render(t, InlineError("Email is required"))},
	})
}

func TestGoldenSweepInlineSuccess(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "inline_success", HTML: utils.Render(t, InlineSuccess("Username is available"))},
	})
}

func TestGoldenSweepSkeletonGroup(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "skeleton_group", HTML: utils.Render(t, SkeletonGroup([]SkeletonVariant{
			SkeletonTableRow, SkeletonTableRow, SkeletonTableRow,
		}))},
	})
}

func TestGoldenSweepToastContainer(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "toast_container", HTML: utils.Render(t, ToastContainer("test-nonce"))},
	})
}
