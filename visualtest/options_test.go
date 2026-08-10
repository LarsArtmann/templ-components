package visualtest

import (
	"testing"
)

func TestInteractionStateString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		state InteractionState
		want  string
	}{
		{StateRest, "rest"},
		{StateHover, "hover"},
		{StateFocus, "focus"},
		{StateClick, "click"},
		{StateContext, "context"},
		{InteractionState(99), "unknown(99)"},
	}

	for _, tt := range tests {
		if got := tt.state.String(); got != tt.want {
			t.Errorf("InteractionState(%d).String() = %q, want %q", int(tt.state), got, tt.want)
		}
	}
}

func TestBoolHelper(t *testing.T) {
	t.Parallel()

	b := new(true)
	if b == nil || !*b {
		t.Fatal("Bool(true) should return non-nil pointer to true")
	}

	f := new(false)
	if f == nil || *f {
		t.Fatal("Bool(false) should return non-nil pointer to false")
	}
}

func TestViewportPresets(t *testing.T) {
	t.Parallel()

	if ViewportMobile.Width != 375 || ViewportMobile.Height != 667 {
		t.Errorf("ViewportMobile = %dx%d, want 375x667", ViewportMobile.Width, ViewportMobile.Height)
	}

	if ViewportTablet.Width != 768 || ViewportTablet.Height != 1024 {
		t.Errorf("ViewportTablet = %dx%d, want 768x1024", ViewportTablet.Width, ViewportTablet.Height)
	}

	if ViewportDesktop.Width != 1280 || ViewportDesktop.Height != 800 {
		t.Errorf("ViewportDesktop = %dx%d, want 1280x800", ViewportDesktop.Width, ViewportDesktop.Height)
	}
}

func TestResolveOptionsTriState(t *testing.T) {
	t.Parallel()

	// nil Dark → default light.
	o := resolveOptions(nil)
	if isDark(o) {
		t.Error("default should be light mode")
	}

	// Explicit dark.
	o = resolveOptions([]Options{{Dark: new(true)}})
	if !isDark(o) {
		t.Error("Dark=Bool(true) should be dark mode")
	}

	// Explicit light (tri-state: distinguish from unset).
	o = resolveOptions([]Options{{Dark: new(false)}})
	if isDark(o) {
		t.Error("Dark=Bool(false) should be light mode")
	}

	if o.Dark == nil {
		t.Error("Dark=Bool(false) should be non-nil (explicitly set)")
	}

	// Later option overrides earlier.
	o = resolveOptions([]Options{{Dark: new(true)}, {Dark: new(false)}})
	if isDark(o) {
		t.Error("later Dark=Bool(false) should override earlier Bool(true)")
	}
}
