package forms

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestDefaultRadioProps(t *testing.T) {
	t.Parallel()
	props := DefaultRadioProps()
	if props.Value != "" {
		t.Error("expected empty default value")
	}
}

func TestDefaultRadioGroupProps(t *testing.T) {
	t.Parallel()
	props := DefaultRadioGroupProps()
	if props.Name != "" {
		t.Error("expected empty default name")
	}
}

// --- Input ---

func TestInputReadOnly(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Input(InputProps{
		Label:    "Readonly",
		Name:     "ro",
		ReadOnly: true,
	}))
	utils.AssertContains(t, output, `readonly`)
}

func TestInputMaxLength(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Input(InputProps{
		Label:     "Limited",
		Name:      "limited",
		MaxLength: 100,
	}))
	utils.AssertContains(t, output, `maxlength="100"`)
}

func TestInputWithAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Input(InputProps{
		BaseProps: utils.BaseProps{AriaLabel: "Search field"},
		Name:      "search",
	}))
	utils.AssertContains(t, output, `aria-label="Search field"`)
}

func TestInputWithoutLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Input(InputProps{
		Name: "nolabel",
	}))
	utils.AssertNotContains(t, output, "<label")
}

// --- Checkbox ---

func TestCheckboxWithValue(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Checkbox(CheckboxProps{
		Label: "Accept",
		Name:  "accept",
		Value: "yes",
	}))
	utils.AssertContains(t, output, `value="yes"`)
}

func TestCheckboxRequired(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Checkbox(CheckboxProps{
		Label:    "Required",
		Name:     "req",
		Required: true,
	}))
	utils.AssertContains(t, output, `required`)
}

func TestCheckboxDisabled(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Checkbox(CheckboxProps{
		Label:    "Disabled",
		Name:     "dis",
		Disabled: true,
	}))
	utils.AssertContains(t, output, `disabled`)
}

func TestCheckboxWithAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Checkbox(CheckboxProps{
		BaseProps: utils.BaseProps{AriaLabel: "Accept terms"},
		Name:      "al",
	}))
	utils.AssertContains(t, output, `aria-label="Accept terms"`)
}

func TestCheckboxWithError(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Checkbox(CheckboxProps{
		Label: "Err",
		Name:  "err",
		Error: "Must accept",
	}))
	utils.AssertContains(t, output, "Must accept")
}

func TestCheckboxWithoutLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Checkbox(CheckboxProps{
		Name: "nolabel",
	}))
	utils.AssertContains(t, output, `type="checkbox"`)
}

// --- Form ---

func TestFormWithAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Form(FormProps{
		Action:    "/submit",
		BaseProps: utils.BaseProps{AriaLabel: "Login form"},
	}))
	utils.AssertContains(t, output, `aria-label="Login form"`)
}

func TestFormWithNilContent(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Form(FormProps{
		Action: "/submit",
	}))
	utils.AssertContains(t, output, `<form`)
}

func TestFormWithClass(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Form(FormProps{
		Action:    "/submit",
		BaseProps: utils.BaseProps{Class: "space-y-4"},
	}))
	utils.AssertContains(t, output, "space-y-4")
}

func TestFormWithAttrs(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Form(FormProps{
		Action:    "/submit",
		BaseProps: utils.BaseProps{Attrs: templ.Attributes{"data-testid": "form"}},
	}))
	utils.AssertContains(t, output, `data-testid="form"`)
}

// --- FileInput ---

func TestFileInputRequired(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FileInput(FileInputProps{
		Label:    "Upload",
		Name:     "file",
		Required: true,
	}))
	utils.AssertContains(t, output, `required`)
}

func TestFileInputDisabled(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FileInput(FileInputProps{
		Label:    "Upload",
		Name:     "file",
		Disabled: true,
	}))
	utils.AssertContains(t, output, `disabled`)
}

func TestFileInputWithAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FileInput(FileInputProps{
		BaseProps: utils.BaseProps{AriaLabel: "Upload file"},
		Name:      "upload",
	}))
	utils.AssertContains(t, output, `aria-label="Upload file"`)
}

func TestFileInputWithError(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FileInput(FileInputProps{
		Label: "Upload",
		Name:  "file",
		Error: "File too large",
	}))
	utils.AssertContains(t, output, "File too large")
}

func TestFileInputWithHelpText(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FileInput(FileInputProps{
		Label:    "Upload",
		Name:     "file",
		HelpText: "Max 5MB",
	}))
	utils.AssertContains(t, output, "Max 5MB")
}

// --- Radio ---

func TestRadioWithoutID(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Radio(RadioProps{
		Label: "Option",
		Name:  "opt",
	}))
	utils.AssertContains(t, output, `type="radio"`)
}

func TestRadioChecked(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Radio(RadioProps{
		BaseProps: utils.BaseProps{ID: "opt-1"},
		Label:     "Selected",
		Name:      "opt",
		Checked:   true,
	}))
	utils.AssertContains(t, output, `checked`)
}

func TestRadioDisabled(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Radio(RadioProps{
		Label:    "Disabled",
		Name:     "opt",
		Disabled: true,
	}))
	utils.AssertContains(t, output, `disabled`)
}

func TestRadioWithAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Radio(RadioProps{
		BaseProps: utils.BaseProps{AriaLabel: "Plan option"},
		Name:      "opt",
	}))
	utils.AssertContains(t, output, `aria-label="Plan option"`)
}

func TestRadioWithoutLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Radio(RadioProps{
		Name: "nolabel",
	}))
	utils.AssertContains(t, output, `type="radio"`)
}

// --- RadioGroup ---

func TestRadioGroupRequired(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, RadioGroup(RadioGroupProps{
		Name:     "plan",
		Required: true,
		Options:  []RadioOption{{Label: "A", Value: "a"}},
	}))
	utils.AssertContains(t, output, `aria-required="true"`)
}

func TestRadioGroupWithHelpText(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, RadioGroup(RadioGroupProps{
		Name:     "plan",
		HelpText: "Choose one",
		Options:  []RadioOption{{Label: "A", Value: "a"}},
	}))
	utils.AssertContains(t, output, "Choose one")
}

func TestRadioGroupWithoutLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, RadioGroup(RadioGroupProps{
		Name:    "plan",
		Options: []RadioOption{{Label: "A", Value: "a"}},
	}))
	utils.AssertContains(t, output, `type="radio"`)
}

// --- Toggle ---

func TestToggleDisabled(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Toggle(ToggleProps{
		Label:    "Setting",
		Name:     "setting",
		Disabled: true,
	}))
	utils.AssertContains(t, output, `disabled`)
}

func TestToggleWithAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Toggle(ToggleProps{
		BaseProps: utils.BaseProps{AriaLabel: "Toggle setting"},
		Name:      "toggle",
	}))
	utils.AssertContains(t, output, `aria-label="Toggle setting"`)
}

func TestToggleWithoutLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, Toggle(ToggleProps{
		Name: "nolabel",
	}))
	utils.AssertContains(t, output, `type="checkbox"`)
}

// --- InputGroup ---

func TestInputGroupWithoutLeftAddon(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, InputGroup(InputGroupProps{}))
	utils.AssertContains(t, output, "relative")
}

func TestInputGroupWithRightAddon(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, InputGroup(InputGroupProps{
		RightAddon: templ.Raw("<span>.00</span>"),
	}))
	utils.AssertContains(t, output, ".00")
}

func TestInputGroupWithBothAddons(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, InputGroup(InputGroupProps{
		LeftAddon:  templ.Raw("<span>$</span>"),
		RightAddon: templ.Raw("<span>.00</span>"),
	}))
	utils.AssertContains(t, output, "$")
	utils.AssertContains(t, output, ".00")
}

// --- FormFieldWrapper ---

func TestFormFieldWrapperEmptyLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FormFieldWrapper(FormFieldProps{}))
	if output != "" {
		t.Errorf("expected empty output, got %q", output)
	}
}

func TestFormFieldWrapperNoErrorNoHelp(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FormFieldWrapper(FormFieldProps{ID: "name", Label: "Name"}))
	utils.AssertContains(t, output, "Name")
}

func TestFormFieldWrapperWithError(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FormFieldWrapper(FormFieldProps{ID: "email", Label: "Email", Error: "Required"}))
	utils.AssertContains(t, output, "Required")
}

func TestFormFieldWrapperWithHelpText(t *testing.T) {
	t.Parallel()
	output := utils.Render(
		t,
		FormFieldWrapper(FormFieldProps{ID: "email", Label: "Email", HelpText: "Enter your email"}),
	)
	utils.AssertContains(t, output, "Enter your email")
}

func TestFormFieldWrapperRequired(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, FormFieldWrapper(FormFieldProps{ID: "name", Label: "Name", Required: true}))
	utils.AssertContains(t, output, "Name")
}

// --- ValidationSummary ---

func TestValidationSummaryWithErrors(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ValidationSummary(ValidationSummaryProps{
		Errors: []ValidationError{
			{Field: "email", Message: "Email is required"},
			{Field: "name", Message: "Name too short"},
		},
	}))
	utils.AssertContains(t, output, `role="alert"`)
	utils.AssertContains(t, output, "Email is required")
	utils.AssertContains(t, output, "Name too short")
	utils.AssertContains(t, output, `href="#email"`)
	utils.AssertContains(t, output, "2 errors found")
}

func TestValidationSummaryEmpty(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ValidationSummary(ValidationSummaryProps{
		Errors: []ValidationError{},
	}))
	if output != "" {
		t.Errorf("expected empty output, got: %s", output)
	}
}

func TestValidationSummaryNoField(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ValidationSummary(ValidationSummaryProps{
		Errors: []ValidationError{
			{Message: "Generic error"},
		},
	}))
	utils.AssertContains(t, output, "Generic error")
	utils.AssertNotContains(t, output, `href=`)
}

func TestValidationSummaryWithID(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ValidationSummary(ValidationSummaryProps{
		BaseProps: utils.BaseProps{ID: "val-summary"},
		Errors:    []ValidationError{{Field: "x", Message: "err"}},
	}))
	utils.AssertContains(t, output, `id="val-summary"`)
}

func TestValidationSummaryWithAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ValidationSummary(ValidationSummaryProps{
		BaseProps: utils.BaseProps{AriaLabel: "Form errors"},
		Errors:    []ValidationError{{Field: "x", Message: "err"}},
	}))
	utils.AssertContains(t, output, `aria-label="Form errors"`)
}

func TestRadioGroupAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, RadioGroup(RadioGroupProps{
		Name:      "plan",
		Label:     "Select Plan",
		BaseProps: utils.BaseProps{AriaLabel: "Plan selection group"},
		Options:   []RadioOption{{Value: "free", Label: "Free"}},
	}))
	utils.AssertContains(t, output, `aria-label="Plan selection group"`)
}

func TestRadioGroupWithoutAriaLabel(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, RadioGroup(RadioGroupProps{
		Name:    "plan",
		Label:   "Select Plan",
		Options: []RadioOption{{Value: "free", Label: "Free"}},
	}))
	utils.AssertNotContains(t, output, `aria-label=`)
}
