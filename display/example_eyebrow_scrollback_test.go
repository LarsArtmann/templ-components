package display_test

import (
	"bytes"
	"context"

	"github.com/larsartmann/templ-components/display"
)

func ExampleEyebrow() {
	var buf bytes.Buffer

	_ = display.Eyebrow(display.EyebrowProps{
		Text: "Deploy #142 · production",
	}).Render(context.Background(), &buf)
}

func ExampleScrollback() {
	props := display.DefaultScrollbackProps()
	props.Lines = []display.ScrollbackLine{
		{
			Timestamp: "12:47:03.184", Tag: "query",
			Text: "ads.example.com A", Tone: display.ScrollbackToneInfo,
		},
		{
			Timestamp: "12:47:03.184", Tag: "match",
			Text: "blocklist: StevenBlack/hosts", Tone: display.ScrollbackToneDanger,
		},
		{
			Timestamp: "12:47:03.185", Tag: "action",
			Text: "NXDOMAIN", Tone: display.ScrollbackToneWarning,
		},
	}

	var buf bytes.Buffer

	_ = display.Scrollback(props).Render(context.Background(), &buf)
}
