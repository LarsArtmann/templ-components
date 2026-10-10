package visualtest_test

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/visualtest"
)

// TestSegmentBar pins the proportion bar's pixel rendering — flex-grow
// proportions, palette colors, and the legend — in light and dark (plan C7f:
// proportions are pixel-relevant, not just HTML-golden relevant).
func TestSegmentBar(t *testing.T) {
	t.Parallel()

	visualtest.AssertScreenshot(t, "segment_bar/light", segmentBarComponent())
	visualtest.AssertScreenshot(
		t,
		"segment_bar/dark",
		segmentBarComponent(),
		visualtest.Options{Dark: new(true)},
	)
}

func segmentBarComponent() templ.Component {
	return display.SegmentBar(display.SegmentBarProps{
		Segments: []display.SegmentBarSegment{
			{Label: "Go", Value: 62.3},
			{Label: "TypeScript", Value: 27.1},
			{Label: "Shell", Value: 8.4},
			{Label: "Other", Value: 2.2},
		},
	})
}
