package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

// Golden sweep for StatusDot and LivePill (TODO_LIST #392).

func TestGoldenSweepStatusIndicators(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "status_dot_static", HTML: utils.Render(t, StatusDot(StatusDotProps{Tone: StatusToneDanger}))},
		{Name: "status_dot_pulse_labeled", HTML: utils.Render(t, StatusDot(StatusDotProps{
			Tone:  StatusToneSuccess,
			Pulse: true,
			Label: "Live",
		}))},
		{Name: "status_dot_warning", HTML: utils.Render(t, StatusDot(StatusDotProps{Tone: StatusToneWarning, Label: "Reconnecting"}))},
		{Name: "live_pill_live", HTML: utils.Render(t, LivePill(LivePillProps{Tone: StatusToneSuccess, Text: "Live", Pulse: true}))},
		{Name: "live_pill_reconnecting", HTML: utils.Render(t, LivePill(LivePillProps{Tone: StatusToneWarning, Text: "Reconnecting…"}))},
		{Name: "live_pill_paused", HTML: utils.Render(t, LivePill(LivePillProps{Tone: StatusToneNeutral, Text: "Paused"}))},
	})
}
