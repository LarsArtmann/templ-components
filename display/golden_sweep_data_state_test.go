package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

// Golden sweep for DataState (TODO_LIST #393): the four-rung honesty ladder.

func TestGoldenSweepDataState(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "data_state_empty", HTML: utils.Render(t, DataState(DataStateProps{
			State:       DataStateEmpty,
			Title:       "No tracked device activity",
			Description: "No requests landed in the last 30 days.",
		}))},
		{Name: "data_state_unavailable", HTML: utils.Render(t, DataState(DataStateProps{
			State:       DataStateUnavailable,
			Title:       "Unavailable — retrying",
			Description: "The query failed; this region keeps retrying every 30s.",
		}))},
		{Name: "data_state_disabled", HTML: utils.Render(t, DataState(DataStateProps{
			State:       DataStateDisabled,
			Title:       "Tracking is disabled",
			Description: "Enable tracking in settings to see device activity.",
		}))},
		{Name: "data_state_custom_icon", HTML: utils.Render(t, DataState(DataStateProps{
			State: DataStateEmpty,
			Title: "Nothing scheduled",
			Icon:  "calendar",
		}))},
	})
}
