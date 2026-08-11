package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

func TestGoldenSectionHeading(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		props SectionHeadingProps
	}{
		{
			name: "h2_center",
			props: SectionHeadingProps{
				Title: "Experience",
				Level: HeadingLevelH2,
				Align: TextAlignCenter,
			},
		},
		{
			name: "h3_left_subtitle",
			props: SectionHeadingProps{
				Title:    "Skills",
				Level:    HeadingLevelH3,
				Align:    TextAlignLeft,
				SubTitle: "Programming languages and tools",
			},
		},
		{
			name: "default_level",
			props: SectionHeadingProps{
				Title: "About",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := utils.Render(t, SectionHeading(tt.props))
			golden.Assert(t, "section_heading_"+tt.name, output)
		})
	}
}
