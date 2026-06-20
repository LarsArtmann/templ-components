package forms

func ExampleForm() {
	_ = Form(FormProps{
		Action:    "/api/users",
		Method:    FormPost,
		CSRFToken: "example-token",
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
