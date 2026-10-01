package forms_test

import (
	"github.com/larsartmann/templ-components/forms"
)

func ExampleCheckbox() {
	_ = forms.Checkbox(forms.DefaultCheckboxProps())
	// Output:
}

func ExampleCombobox() {
	_ = forms.Combobox(forms.DefaultComboboxProps())
	// Output:
}

func ExampleDatePicker() {
	_ = forms.DatePicker(forms.DefaultDatePickerProps())
	// Output:
}

func ExampleDirtyGuard() {
	_ = forms.DirtyGuard(forms.DirtyGuardProps{})
	// Output:
}

func ExampleFileInput() {
	_ = forms.FileInput(forms.DefaultFileInputProps())
	// Output:
}

func ExampleInputGroup() {
	_ = forms.InputGroup(forms.DefaultInputGroupProps())
	// Output:
}

func ExampleFieldError() {
	_ = forms.FieldError("email", "A valid email is required")
	// Output:
}

func ExampleFormFieldWrapper() {
	_ = forms.FormFieldWrapper(forms.FormFieldProps{ID: "email", Label: "Email"})
	// Output:
}

func ExampleLabel() {
	_ = forms.Label("email", "Email address", true)
	// Output:
}

func ExampleRadio() {
	_ = forms.Radio(forms.DefaultRadioProps())
	// Output:
}

func ExampleRadioGroup() {
	_ = forms.RadioGroup(forms.DefaultRadioGroupProps())
	// Output:
}

func ExampleToggle() {
	_ = forms.Toggle(forms.DefaultToggleProps())
	// Output:
}

func ExampleValidationSummary() {
	_ = forms.ValidationSummary(forms.DefaultValidationSummaryProps())
	// Output:
}
