package utils

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils/svg"
)

// TestScriptAttrsOmitEmpty pins the omit-empty rule DIRECTLY in the owning
// package (the integration sweeps cover it component-side; this is the
// unit-level contract): an empty nonce renders NO attribute — never
// nonce="" — and a set nonce carries its value verbatim.
func TestScriptAttrsOmitEmpty(t *testing.T) {
	t.Parallel()

	if attrs := ScriptAttrs(""); attrs != nil {
		t.Errorf("ScriptAttrs(\"\") = %v, want nil", attrs)
	}

	attrs := ScriptAttrs("n1")
	if attrs == nil {
		t.Fatal("ScriptAttrs(n1) = nil, want the nonce attribute")
	}

	if got := attrs["nonce"]; got != "n1" {
		t.Errorf("ScriptAttrs(n1)[nonce] = %q, want n1", got)
	}
}

// TestDismissButtonRender covers the shared dismiss control's render
// contract: data-dismiss hook, aria-label, the X-mark icon path, and the
// caller-owned tone classes landing on the button.
func TestDismissButtonRender(t *testing.T) {
	t.Parallel()

	output := Render(t, DismissButton("bg-red-100", "text-red-800"))

	for _, want := range []string{
		`type="button"`,
		`data-dismiss="alert"`,
		`aria-label="Dismiss"`,
		"bg-red-100",
		"text-red-800",
		svg.PathXMark,
	} {
		if !strings.Contains(output, want) {
			t.Errorf("DismissButton output missing %q", want)
		}
	}
}
