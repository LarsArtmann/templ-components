package forms

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// renderWrapperWithChildren renders FormFieldWrapper with raw markup as its
// children, mirroring how Input/Select route their control through it.
func renderWrapperWithChildren(t *testing.T, props FormFieldProps, child string) string {
	t.Helper()

	var buf bytes.Buffer

	wrapper := FormFieldWrapper(props)

	ctx := templ.WithChildren(context.Background(), templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := w.Write([]byte(child))

		return err
	}))

	if err := wrapper.Render(ctx, &buf); err != nil {
		t.Fatalf("failed to render FormFieldWrapper: %v", err)
	}

	return strings.TrimSpace(buf.String())
}

func TestFormFieldWrapperStackedSingleChild(t *testing.T) {
	t.Parallel()

	out := renderWrapperWithChildren(t, FormFieldProps{
		ID: "email", Label: "Email", Error: "Required", HelpText: "Use work email", Stacked: true,
	}, `<input type="email" name="email" />`)

	if !strings.HasPrefix(out, "<div>") || !strings.HasSuffix(out, "</div>") {
		t.Fatalf("expected one wrapping <div> around all field parts, got: %s", out)
	}

	utils.AssertContainsAll(t, out, "<label", "<input", `role="alert"`, "Use work email")

	if strings.Index(out, "<label") > strings.Index(out, "<input") {
		t.Errorf("label must precede the control inside the container: %s", out)
	}
}

func TestFormFieldWrapperStackedFalseKeepsSiblings(t *testing.T) {
	t.Parallel()

	out := renderWrapperWithChildren(t, FormFieldProps{ID: "email", Label: "Email"}, `<input id="email" />`)

	if strings.HasPrefix(out, "<div>") {
		t.Fatalf("default path must keep the parts as unwrapped siblings, got: %s", out)
	}

	if strings.Index(out, "<label") > strings.Index(out, "<input") {
		t.Errorf("label must precede the control: %s", out)
	}
}

func TestFormFieldWrapperStackedPropagatesFromConsumers(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		out  string
	}{
		{"input", utils.Render(t, Input(InputProps{Name: "email", Label: "Email", Stacked: true}))},
		{"select", utils.Render(t, Select(SelectProps{Name: "plan", Label: "Plan", Stacked: true}))},
		{"textarea", utils.Render(t, Textarea(TextareaProps{Name: "bio", Label: "Bio", Stacked: true}))},
		{"file_input", utils.Render(t, FileInput(FileInputProps{Name: "doc", Label: "Document", Stacked: true}))},
		{"combobox", utils.Render(t, Combobox(ComboboxProps{Name: "region", Label: "Region", Stacked: true}))},
		{"date_picker", utils.Render(t, DatePicker(DatePickerProps{Name: "start", Label: "Start", Stacked: true}))},
	} {
		if !strings.HasPrefix(tc.out, "<div>") {
			t.Errorf("%s: Stacked must wrap the field in one container, got: %s", tc.name, tc.out)
		}
	}
}
