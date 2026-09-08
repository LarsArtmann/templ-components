package forms

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestZZProbeInlineCtx(t *testing.T) {
	input := DefaultInputProps()
	input.Name = "q"
	input.Label = "q"

	var buf bytes.Buffer

	err := Input(input).Render(context.Background(), &buf)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("baseline group div present: %v", strings.Contains(buf.String(), "min-w-40"))

	buf.Reset()

	marked := context.WithValue(context.Background(), formLayoutContextKey{}, FormLayoutInline)
	err = Input(input).Render(marked, &buf)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("marked group div present: %v", strings.Contains(buf.String(), "min-w-40"))

	buf.Reset()

	page := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		t.Logf("inner children ctx has marker: %v", ctx.Value(formLayoutContextKey{}) != nil)

		return Input(input).Render(ctx, w)
	})

	_ = Form(FormProps{Action: "/f", Layout: FormLayoutInline}).Render(
		templ.WithChildren(context.Background(), page), &buf,
	)

	t.Logf("form-inline group div present: %v", strings.Contains(buf.String(), "min-w-40"))
}
