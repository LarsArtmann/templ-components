package forms

import "testing"

// FuzzInputType verifies inputType() never panics on arbitrary input.
func FuzzInputType(f *testing.F) {
	f.Add("text")
	f.Add("email")
	f.Add("")
	f.Add("invalid-type")
	f.Fuzz(func(t *testing.T, input string) {
		result := inputType(InputType(input))
		if result == "" {
			t.Errorf("inputType(%q) returned empty string", input)
		}
	})
}
