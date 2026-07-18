package icons

import "testing"

func TestNameIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		n    Name
		want bool
	}{
		{"Home icon", Home, true},
		{"Spinner", Spinner, true},
		{"Question icon", Question, true},
		{"ChevronDown icon", ChevronDown, true},
		{"invalid empty", Name(""), false},
		{"invalid bogus", Name("bogus-icon"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := NameIsValid(tt.n); got != tt.want {
				t.Errorf("NameIsValid(%q) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}
