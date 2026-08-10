package forms

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils/golden"
	"github.com/larsartmann/templ-components/utils"
)

// Golden sweep for forms components that previously lacked golden tests.
// Combobox uses EnsureID — safe for golden testing via ID normalization.

func TestGoldenSweepCheckbox(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "checkbox_default", HTML: utils.Render(t, Checkbox(CheckboxProps{
			Name: "terms", Label: "I agree to the terms of service", Checked: true,
		}))},
		{Name: "checkbox_error", HTML: utils.Render(t, Checkbox(CheckboxProps{
			Name: "newsletter", Label: "Subscribe to newsletter", Error: "Please select at least one option",
		}))},
		{Name: "checkbox_disabled", HTML: utils.Render(t, Checkbox(CheckboxProps{
			Name: "locked", Label: "Cannot change this", Disabled: true, HelpText: "Managed by admin",
		}))},
	})
}

func TestGoldenSweepToggle(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "toggle_on", HTML: utils.Render(t, Toggle(ToggleProps{
			Name: "notifications", Label: "Enable notifications", Checked: true,
		}))},
		{Name: "toggle_off", HTML: utils.Render(t, Toggle(ToggleProps{
			Name: "darkmode", Label: "Dark mode",
		}))},
		{Name: "toggle_disabled", HTML: utils.Render(t, Toggle(ToggleProps{
			Name: "readonly", Label: "Read-only setting", Disabled: true, HelpText: "Locked by policy",
		}))},
	})
}

func TestGoldenSweepRadioGroup(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "radiogroup_stacked", HTML: utils.Render(t, RadioGroup(RadioGroupProps{
			Name: "plan", Label: "Select a plan",
			Options: []RadioOption{
				{Value: "free", Label: "Free tier"},
				{Value: "pro", Label: "Pro tier", Checked: true},
				{Value: "enterprise", Label: "Enterprise"},
			},
		}))},
		{Name: "radiogroup_inline", HTML: utils.Render(t, RadioGroup(RadioGroupProps{
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
		{Name: "combobox_basic", HTML: utils.Render(t, Combobox(ComboboxProps{
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
		{Name: "input_group_both", HTML: utils.Render(t, InputGroup(InputGroupProps{
			LeftAddon:  templ.Raw(`<span>$</span>`),
			RightAddon: templ.Raw(`<span>.00</span>`),
		}))},
	})
}

func TestGoldenSweepForm(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "form_basic", HTML: utils.Render(t, Form(FormProps{
			Action: "/submit", Method: FormPost,
		}))},
	})
}

func TestGoldenSweepValidationSummary(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "validation_summary", HTML: utils.Render(t, ValidationSummary(ValidationSummaryProps{
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
		{Name: "file_input_basic", HTML: utils.Render(t, FileInput(FileInputProps{
			Name: "avatar", Label: "Upload avatar", Accept: "image/*",
		}))},
		{Name: "file_input_multiple", HTML: utils.Render(t, FileInput(FileInputProps{
			Name: "documents", Label: "Upload documents", Multiple: true, HelpText: "PDF, DOCX up to 10MB each",
		}))},
	})
}

func TestGoldenSweepDatePicker(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "date_picker_basic", HTML: utils.Render(t, DatePicker(DatePickerProps{
			Name: "dob", Label: "Date of birth",
		}))},
		{Name: "date_picker_with_range", HTML: utils.Render(t, DatePicker(DatePickerProps{
			Name: "appointment", Label: "Appointment date", Min: "2026-01-01", Max: "2026-12-31", Value: "2026-03-15",
		}))},
	})
}
