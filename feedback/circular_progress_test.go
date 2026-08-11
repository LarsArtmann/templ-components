package feedback

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

func TestGoldenCircularProgress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		props CircularProgressProps
	}{
		{
			name: "md_75_with_label",
			props: CircularProgressProps{
				Value:     75,
				Size:      CircularProgressSizeMD,
				ShowLabel: true,
			},
		},
		{
			name: "sm_50_no_label",
			props: CircularProgressProps{
				Value: 50,
				Size:  CircularProgressSizeSM,
			},
		},
		{
			name: "lg_100_clamped",
			props: CircularProgressProps{
				Value: 150,
				Size:  CircularProgressSizeLG,
			},
		},
		{
			name:  "default_zero",
			props: CircularProgressProps{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, CircularProgress(tt.props))
			golden.Assert(t, "circular_progress_"+tt.name, output)
		})
	}
}
