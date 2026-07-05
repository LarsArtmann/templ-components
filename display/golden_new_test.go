package display

import (
	"testing"
	"time"

	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

func TestGoldenCopyButton(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, CopyButton(CopyButtonProps{
		Text:      "npm install foo",
		Label:     "Copy",
		Icon:      true,
		BaseProps: utils.BaseProps{Nonce: "abc123"},
	}))
	golden.Assert(t, "copy_button", output)
}

func TestGoldenRelativeTime(t *testing.T) {
	t.Parallel()
	ts := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	output := utils.Render(t, RelativeTime(RelativeTimeProps{
		Time: ts,
	}))
	golden.Assert(t, "relative_time", output)
}

func TestGoldenCountBadge(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, CountBadge(CountBadgeProps{Count: 5}))
	golden.Assert(t, "count_badge", output)
}

func TestGoldenDefinitionGrid(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, DefinitionGrid(DefinitionGridProps{
		Cols: GridCols2,
		Items: []DefinitionItem{
			{Term: "CPU", Detail: "42%"},
			{Term: "Memory", Detail: "8.2 GB"},
		},
	}))
	golden.Assert(t, "definition_grid", output)
}

func TestGoldenImage(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Image(ImageProps{
		Src:    "/photo.jpg",
		Alt:    "Profile photo",
		Width:  128,
		Height: 128,
	}))
	golden.Assert(t, "image", output)
}
