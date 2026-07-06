package display

import "testing"

// FuzzButtonHTMLType verifies buttonHTMLType() never panics on arbitrary input.
func FuzzButtonHTMLType(f *testing.F) {
	f.Add("button")
	f.Add("submit")
	f.Add("reset")
	f.Add("")
	f.Add("invalid")
	f.Fuzz(func(t *testing.T, input string) {
		result := buttonHTMLType(ButtonHTMLType(input))
		// Should always return a valid HTML type or default
		if result == "" {
			t.Errorf("buttonHTMLType(%q) returned empty", input)
		}
	})
}
