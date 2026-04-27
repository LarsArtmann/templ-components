package utils

import (
	"testing"
	"time"

	"github.com/a-h/templ"
)

func TestCurrentYear(t *testing.T) {
	got := CurrentYear()
	want := time.Now().Format("2006")
	if got != want {
		t.Errorf("CurrentYear() = %q, want %q", got, want)
	}
}

func TestTernary(t *testing.T) {
	tests := []struct {
		name      string
		condition bool
		a, b      string
		want      string
	}{
		{"true returns a", true, "yes", "no", "yes"},
		{"false returns b", false, "yes", "no", "no"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Ternary(tt.condition, tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Ternary() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPtr(t *testing.T) {
	v := "hello"
	p := Ptr(v)
	if p == nil {
		t.Fatal("Ptr() returned nil")
	}
	if *p != v {
		t.Errorf("*Ptr() = %q, want %q", *p, v)
	}
}

func TestDeref(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		v := 42
		got := Deref(&v)
		if got != 42 {
			t.Errorf("Deref() = %d, want 42", got)
		}
	})
	t.Run("nil", func(t *testing.T) {
		var p *int
		got := Deref(p)
		if got != 0 {
			t.Errorf("Deref() = %d, want 0", got)
		}
	})
}

func TestDerefOr(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		v := 42
		got := DerefOr(&v, 99)
		if got != 42 {
			t.Errorf("DerefOr() = %d, want 42", got)
		}
	})
	t.Run("nil", func(t *testing.T) {
		var p *int
		got := DerefOr(p, 99)
		if got != 99 {
			t.Errorf("DerefOr() = %d, want 99", got)
		}
	})
}

func TestClass(t *testing.T) {
	tests := []struct {
		name      string
		defaults  string
		overrides string
		want      string
	}{
		{"empty overrides", "a b c", "", "a b c"},
		{"with overrides", "a b c", "d e", "d e"},
		{"trim spaces", "  a b  ", "", "a b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Class(tt.defaults, tt.overrides)
			if got != tt.want {
				t.Errorf("Class(%q, %q) = %q, want %q", tt.defaults, tt.overrides, got, tt.want)
			}
		})
	}
}

func TestMergeAttrs(t *testing.T) {
	t.Run("merge two maps", func(t *testing.T) {
		a := templ.Attributes{"data-foo": "1"}
		b := templ.Attributes{"data-bar": "2"}
		got := MergeAttrs(a, b)
		if got["data-foo"] != "1" || got["data-bar"] != "2" {
			t.Errorf("MergeAttrs() = %v", got)
		}
	})
	t.Run("later overrides earlier", func(t *testing.T) {
		a := templ.Attributes{"data-foo": "1"}
		b := templ.Attributes{"data-foo": "2"}
		got := MergeAttrs(a, b)
		if got["data-foo"] != "2" {
			t.Errorf("MergeAttrs() data-foo = %q, want 2", got["data-foo"])
		}
	})
	t.Run("empty maps", func(t *testing.T) {
		got := MergeAttrs()
		if len(got) != 0 {
			t.Errorf("MergeAttrs() len = %d, want 0", len(got))
		}
	})
}
