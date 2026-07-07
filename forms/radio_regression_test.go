package forms

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestRadioGroupRequiredPropagatesToInputs(t *testing.T) {
	t.Parallel()

	props := RadioGroupProps{
		Name:     "plan",
		Label:    "Select a plan",
		Required: true,
		Options: []RadioOption{
			{Value: "free", Label: "Free"},
			{Value: "pro", Label: "Pro"},
		},
	}

	var buf strings.Builder
	_ = RadioGroup(props).Render(t.Context(), &buf)
	html := buf.String()

	if !strings.Contains(html, "required") {
		t.Error("RadioGroup with Required=true should render 'required' on radio inputs")
	}
}

func TestRadioOptionCheckedPreSelects(t *testing.T) {
	t.Parallel()

	props := RadioGroupProps{
		Name:  "plan",
		Label: "Select a plan",
		Options: []RadioOption{
			{Value: "free", Label: "Free"},
			{Value: "pro", Label: "Pro", Checked: true},
		},
	}

	var buf strings.Builder
	_ = RadioGroup(props).Render(t.Context(), &buf)
	html := buf.String()

	if !strings.Contains(html, "checked") {
		t.Error("RadioOption with Checked=true should render 'checked' on the input")
	}
}

func TestInputGroupRightAddonHasPointerEventsNone(t *testing.T) {
	t.Parallel()

	props := InputGroupProps{
		RightAddon: templ.NopComponent,
	}

	var buf strings.Builder
	_ = InputGroup(props).Render(t.Context(), &buf)
	html := buf.String()

	if !strings.Contains(html, "pointer-events-none") {
		t.Error("InputGroup right addon div should have pointer-events-none class")
	}
}

func TestFieldErrorHasRoleAlert(t *testing.T) {
	t.Parallel()

	var buf strings.Builder
	_ = FieldError("email", "Required").Render(t.Context(), &buf)
	html := buf.String()

	if !strings.Contains(html, `role="alert"`) {
		t.Error("FieldError should have role=\"alert\" for screen reader announcement")
	}
}

func TestFieldErrorEmptyMessageRendersNothing(t *testing.T) {
	t.Parallel()

	var buf strings.Builder
	_ = FieldError("email", "").Render(t.Context(), &buf)
	html := buf.String()

	if strings.Contains(html, "<p") {
		t.Error("FieldError with empty message should not render a <p> element")
	}
}
