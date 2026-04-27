package utils

import (
	"testing"
	"time"
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
