package utils

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// TestScriptComponentContextNonce is the ADR-0047 pilot proof: the canonical
// script writer resolves the nonce from the render context when the prop is
// empty, an explicit prop always wins, and the omit-empty rule survives.
func TestScriptComponentContextNonce(t *testing.T) {
	t.Parallel()

	const js = "console.log('x');"

	render := func(t *testing.T, ctx context.Context, nonce string) string {
		t.Helper()

		var buf bytes.Buffer

		if err := ScriptComponent(nonce, js, "test script").Render(ctx, &buf); err != nil {
			t.Fatalf("render: %v", err)
		}

		return buf.String()
	}

	t.Run("props nonce unchanged", func(t *testing.T) {
		t.Parallel()

		got := render(t, context.Background(), "prop-nonce")
		if !strings.Contains(got, `<script nonce="prop-nonce">`) {
			t.Errorf("props nonce not rendered: %s", got)
		}
	})

	t.Run("context fills empty prop", func(t *testing.T) {
		t.Parallel()

		got := render(t, ContextWithNonce(context.Background(), "ctx-nonce"), "")
		if !strings.Contains(got, `<script nonce="ctx-nonce">`) {
			t.Errorf("context nonce not rendered: %s", got)
		}
	})

	t.Run("props and context agree byte-identically", func(t *testing.T) {
		t.Parallel()

		viaProps := render(t, context.Background(), "same-nonce")

		viaContext := render(t, ContextWithNonce(context.Background(), "same-nonce"), "")

		if viaProps != viaContext {
			t.Errorf("context path diverges from props path:\nprops:    %s\ncontext:  %s", viaProps, viaContext)
		}
	})

	t.Run("disagreement is caught: props win", func(t *testing.T) {
		t.Parallel()

		got := render(t, ContextWithNonce(context.Background(), "ctx-nonce"), "prop-nonce")
		if !strings.Contains(got, `nonce="prop-nonce"`) || strings.Contains(got, "ctx-nonce") {
			t.Errorf("explicit prop must win over context: %s", got)
		}
	})

	t.Run("both empty omits the attribute", func(t *testing.T) {
		t.Parallel()

		got := render(t, context.Background(), "")
		if strings.Contains(got, "nonce=") {
			t.Errorf("empty nonce must render no nonce attribute (omit-empty rule): %s", got)
		}

		if !strings.Contains(got, "<script>\n"+js+"</script>") {
			t.Errorf("script body must still render: %s", got)
		}
	})
}
