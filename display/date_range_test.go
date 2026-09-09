package display

import (
	"github.com/a-h/templ"

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

// TestGoldenDateRangeAdjacent pins the block-vs-inline semantics: DateRange
// renders an INLINE <time> element, so two adjacent ranges flow on the same
// line separated only by natural whitespace — the component injects no
// separator or line break between instances. The golden documents the exact
// adjacency output consumers stacking two ranges (e.g. parallel roles) get.
func TestGoldenDateRangeAdjacent(t *testing.T) {
	t.Parallel()

	start := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	first := utils.Render(t, DateRange(DateRangeProps{Start: &start, End: &end}))
	second := utils.Render(t, DateRange(DateRangeProps{
		BaseProps: utils.BaseProps{Attrs: templ.Attributes{"data-test": "second"}},
		Start:     &start,
		End:       &end,
	}))

	golden.Assert(t, "date_range_adjacent", first+"\n"+second)
}
