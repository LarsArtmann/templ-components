package forms

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
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

// TestRadioGroupErrorPropagatesAriaToInputs verifies that aria-invalid and
// aria-describedby are propagated to individual radio <input> elements, not
// just the <fieldset>. Screen readers need these attributes on the focusable
// form controls to announce the error state.
func TestRadioGroupErrorPropagatesAriaToInputs(t *testing.T) {
	t.Parallel()

	props := RadioGroupProps{
		Name:  "plan",
		Label: "Select a plan",
		BaseProps: utils.BaseProps{
			ID: "plan-group",
		},
		Error: "Please select a plan",
		Options: []RadioOption{
			{Value: "free", Label: "Free"},
			{Value: "pro", Label: "Pro"},
		},
	}

	var buf strings.Builder
	_ = RadioGroup(props).Render(t.Context(), &buf)
	html := buf.String()

	// Count aria-invalid occurrences — should be on each input
	ariaInvalidCount := strings.Count(html, `aria-invalid="true"`)
	if ariaInvalidCount < 2 {
		t.Errorf("expected aria-invalid on each radio input (at least 2), got %d", ariaInvalidCount)
	}

	// Each input should also have aria-describedby
	ariaDescribedByCount := strings.Count(html, `aria-describedby`)
	if ariaDescribedByCount < 2 {
		t.Errorf("expected aria-describedby on each radio input (at least 2), got %d", ariaDescribedByCount)
	}
}

// TestRadioGroupHelpTextPropagatesAriaToInputs verifies that aria-describedby
// for help text is also propagated to individual inputs.
func TestRadioGroupHelpTextPropagatesAriaToInputs(t *testing.T) {
	t.Parallel()

	props := RadioGroupProps{
		Name:  "plan",
		Label: "Select a plan",
		BaseProps: utils.BaseProps{
			ID: "plan-group",
		},
		HelpText: "Choose your subscription level",
		Options: []RadioOption{
			{Value: "free", Label: "Free"},
		},
	}

	var buf strings.Builder
	_ = RadioGroup(props).Render(t.Context(), &buf)
	html := buf.String()

	if !strings.Contains(html, `aria-describedby`) {
		t.Error("expected aria-describedby on radio inputs when HelpText is set")
	}
}

// TestRadioGroupNoErrorNoAriaOnInputs verifies that without error or help text,
// no extra aria attributes are added to inputs.
func TestRadioGroupNoErrorNoAriaOnInputs(t *testing.T) {
	t.Parallel()

	props := RadioGroupProps{
		Name:  "plan",
		Label: "Select a plan",
		BaseProps: utils.BaseProps{
			ID: "plan-group",
		},
		Options: []RadioOption{
			{Value: "free", Label: "Free"},
		},
	}

	var buf strings.Builder
	_ = RadioGroup(props).Render(t.Context(), &buf)
	html := buf.String()

	if strings.Contains(html, `aria-invalid`) {
		t.Error("should not have aria-invalid when no error is set")
	}
}
