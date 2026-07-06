package forms

import "testing"

// FuzzFormMethod verifies formMethod() never panics on arbitrary input.
func FuzzFormMethod(f *testing.F) {
	f.Add("GET")
	f.Add("POST")
	f.Add("")
	f.Add("DELETE")
	f.Add("invalid-method")
	f.Fuzz(func(t *testing.T, input string) {
		result := formMethod(FormMethod(input))
		if result == "" {
			t.Errorf("formMethod(%q) returned empty", input)
		}
	})
}
