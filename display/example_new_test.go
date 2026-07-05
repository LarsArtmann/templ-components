package display_test

import (
	"bytes"
	"context"
	"time"

	"github.com/larsartmann/templ-components/display"
)

func ExampleCopyButton() {
	props := display.DefaultCopyButtonProps()
	props.Text = "npm install my-package"
	props.Label = "Copy command"

	var buf bytes.Buffer
	_ = display.CopyButton(props).Render(context.Background(), &buf)
}

func ExampleRelativeTime() {
	props := display.RelativeTimeProps{
		Time: time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC),
	}

	var buf bytes.Buffer
	_ = display.RelativeTime(props).Render(context.Background(), &buf)
}

func ExampleCountBadge() {
	props := display.CountBadgeProps{Count: 12, Max: 99}

	var buf bytes.Buffer
	_ = display.CountBadge(props).Render(context.Background(), &buf)
}

func ExampleDefinitionGrid() {
	props := display.DefinitionGridProps{
		Cols: display.GridCols2,
		Items: []display.DefinitionItem{
			{Term: "CPU", Detail: "42%"},
			{Term: "Memory", Detail: "8.2 GB"},
		},
	}

	var buf bytes.Buffer
	_ = display.DefinitionGrid(props).Render(context.Background(), &buf)
}

func ExampleImage() {
	props := display.ImageProps{
		Src:    "/profile.jpg",
		Alt:    "Profile photo",
		Width:  128,
		Height: 128,
		Lazy:   true,
	}

	var buf bytes.Buffer
	_ = display.Image(props).Render(context.Background(), &buf)
}
