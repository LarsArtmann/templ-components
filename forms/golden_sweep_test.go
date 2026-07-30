package forms

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

// Golden sweep for forms components that previously lacked golden tests.
// Combobox uses EnsureID — safe for golden testing via ID normalization.

func TestGoldenSweepCheckbox(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{"checkbox_default", utils.Render(t, Checkbox(CheckboxProps{
			Name: "terms", Label: "I agree to the terms of service", Checked: true,
		}))},
		{"checkbox_error", utils.Render(t, Checkbox(CheckboxProps{
			Name: "newsletter", Label: "Subscribe to newsletter", Error: "Please select at least one option",
		}))},
		{"checkbox_disabled", utils.Render(t, Checkbox(CheckboxProps{
			Name: "locked", Label: "Cannot change this", Disabled: true, HelpText: "Managed by admin",
		}))},
	})
}

func TestGoldenSweepToggle(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{"toggle_on", utils.Render(t, Toggle(ToggleProps{
			Name: "notifications", Label: "Enable notifications", Checked: true,
		}))},
		{"toggle_off", utils.Render(t, Toggle(ToggleProps{
			Name: "darkmode", Label: "Dark mode",
		}))},
		{"toggle_disabled", utils.Render(t, Toggle(ToggleProps{
			Name: "readonly", Label: "Read-only setting", Disabled: true, HelpText: "Locked by policy",
		}))},
	})
}

func TestGoldenSweepRadioGroup(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{"radiogroup_stacked", utils.Render(t, RadioGroup(RadioGroupProps{
			Name: "plan", Label: "Select a plan",
			Options: []RadioOption{
				{Value: "free", Label: "Free tier"},
				{Value: "pro", Label: "Pro tier", Checked: true},
				{Value: "enterprise", Label: "Enterprise"},
			},
		}))},
		{"radiogroup_inline", utils.Render(t, RadioGroup(RadioGroupProps{
			Name: "size", Label: "Size", Inline: true,
			Options: []RadioOption{
				{Value: "sm", Label: "Small"},
				{Value: "md", Label: "Medium", Checked: true},
				{Value: "lg", Label: "Large"},
			},
		}))},
	})
}

func TestGoldenSweepCombobox(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{"combobox_basic", utils.Render(t, Combobox(ComboboxProps{
			Name: "country", Label: "Country", Placeholder: "Search countries...",
			Options: []ComboboxOption{
				{Value: "de", Label: "Germany"},
				{Value: "at", Label: "Austria"},
				{Value: "ch", Label: "Switzerland"},
			},
		}))},
	})
}

func TestGoldenSweepInputGroup(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{"input_group_both", utils.Render(t, InputGroup(InputGroupProps{
			LeftAddon:  templ.Raw(`<span>$</span>`),
			RightAddon: templ.Raw(`<span>.00</span>`),
		}))},
	})
}

func TestGoldenSweepForm(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{"form_basic", utils.Render(t, Form(FormProps{
			Action: "/submit", Method: FormPost,
		}))},
	})
}

func TestGoldenSweepValidationSummary(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{"validation_summary", utils.Render(t, ValidationSummary(ValidationSummaryProps{
			Errors: []ValidationError{
				{Field: "email", Message: "Email is required"},
				{Field: "password", Message: "Password must be at least 8 characters"},
			},
		}))},
	})
}

func TestGoldenSweepFileInput(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{"file_input_basic", utils.Render(t, FileInput(FileInputProps{
			Name: "avatar", Label: "Upload avatar", Accept: "image/*",
		}))},
		{"file_input_multiple", utils.Render(t, FileInput(FileInputProps{
			Name: "documents", Label: "Upload documents", Multiple: true, HelpText: "PDF, DOCX up to 10MB each",
		}))},
	})
}

func TestGoldenSweepDatePicker(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{"date_picker_basic", utils.Render(t, DatePicker(DatePickerProps{
			Name: "dob", Label: "Date of birth",
		}))},
		{"date_picker_with_range", utils.Render(t, DatePicker(DatePickerProps{
			Name: "appointment", Label: "Appointment date", Min: "2026-01-01", Max: "2026-12-31", Value: "2026-03-15",
		}))},
	})
}
