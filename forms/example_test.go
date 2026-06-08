package forms

import (
	"github.com/a-h/templ"
)

func ExampleForm() {
	_ = Form(FormProps{
		Action:    "/api/users",
		Method:    FormPost,
		CSRFToken: "example-token",
		Content: templ.Join(
			Input(InputProps{Name: "name", Type: InputText, Label: "Name"}),
			Input(InputProps{Name: "email", Type: InputEmail, Label: "Email"}),
			Checkbox(CheckboxProps{Name: "terms", Label: "I agree to the terms"}),
		),
	})
	// Output:
}

func ExampleInput() {
	_ = Input(InputProps{
		Name:        "email",
		Type:        InputEmail,
		Label:       "Email address",
		Placeholder: "you@example.com",
		HelpText:    "We'll never share your email.",
	})
	// Output:
}

func ExampleSelect() {
	_ = Select(SelectProps{
		Name:  "country",
		Label: "Country",
		Options: []SelectOption{
			{Value: "us", Label: "United States"},
			{Value: "de", Label: "Germany"},
		},
	})
	// Output:
}

func ExampleTextarea() {
	_ = Textarea(TextareaProps{
		Name:        "bio",
		Label:       "Biography",
		Placeholder: "Tell us about yourself...",
		Rows:        4,
		MaxLength:   500,
	})
	// Output:
}
