package forms

import (
	"testing"
)

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		name     string
		v        float64
		decimals int
		want     string
	}{
		{"zero returns empty", 0, 2, ""},
		{"integer", 42, 0, "42"},
		{"two decimals", 3.14159, 2, "3.14"},
		{"negative", -2.5, 1, "-2.5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatFloat(tt.v, tt.decimals)
			if got != tt.want {
				t.Errorf("FormatFloat(%v, %d) = %q, want %q", tt.v, tt.decimals, got, tt.want)
			}
		})
	}
}

func TestIsSelected(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		option string
		want   bool
	}{
		{"matches", "a", "a", true},
		{"no match", "a", "b", false},
		{"empty", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSelected(tt.value, tt.option)
			if got != tt.want {
				t.Errorf("IsSelected(%q, %q) = %v, want %v", tt.value, tt.option, got, tt.want)
			}
		})
	}
}

func TestIfNotNil(t *testing.T) {
	type item struct {
		Name string
	}

	t.Run("non-nil", func(t *testing.T) {
		obj := &item{Name: "hello"}
		got := IfNotNil(obj, func(i item) string { return i.Name })
		if got != "hello" {
			t.Errorf("IfNotNil() = %q, want %q", got, "hello")
		}
	})

	t.Run("nil", func(t *testing.T) {
		var obj *item
		got := IfNotNil(obj, func(i item) string { return i.Name })
		if got != "" {
			t.Errorf("IfNotNil() = %q, want empty", got)
		}
	})
}

func TestIfNotNilString(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		s := "hello"
		got := IfNotNilString(&s)
		if got != "hello" {
			t.Errorf("IfNotNilString() = %q, want %q", got, "hello")
		}
	})

	t.Run("nil", func(t *testing.T) {
		got := IfNotNilString(nil)
		if got != "" {
			t.Errorf("IfNotNilString() = %q, want empty", got)
		}
	})
}

func TestBoolToString(t *testing.T) {
	tests := []struct {
		name string
		b    bool
		want string
	}{
		{"true", true, "true"},
		{"false", false, "false"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BoolToString(tt.b)
			if got != tt.want {
				t.Errorf("BoolToString(%v) = %q, want %q", tt.b, got, tt.want)
			}
		})
	}
}
