package display

import (
	"testing"
	"time"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

func TestGoldenDateRange(t *testing.T) {
	t.Parallel()

	start := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		props DateRangeProps
	}{
		{
			name: "range",
			props: DateRangeProps{
				Start: &start,
				End:   &end,
			},
		},
		{
			name: "present",
			props: DateRangeProps{
				Start: &start,
				End:   nil,
			},
		},
		{
			name: "compact_format",
			props: DateRangeProps{
				Start:  &start,
				End:    &end,
				Layout: DateFormatJan2006,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, DateRange(tt.props))
			golden.Assert(t, "date_range_"+tt.name, output)
		})
	}
}
