package feedback

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// Alert's dismiss script must never emit an EMPTY nonce attribute: under a
// strict CSP, `nonce=""` matches no source, so the browser blocks the script
// and the dismiss button silently dies. When the caller omits Nonce the
// component falls back to templ.GetNonce(ctx); when neither is set the
// attribute is omitted entirely (no-CSP sites keep the same behavior minus
// the useless empty attribute). Consumer evidence: crush-daily's
// collectErrorBanner forgot the nonce and its alert dismiss script was
// blocked by the site's script-src 'nonce-…' policy.
func TestAlertNonceFallback(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		ctxNonce  string
		propNonce string
		want      string
		notWant   string
	}{
		{
			name:      "prop nonce wins",
			ctxNonce:  "ctx-nonce",
			propNonce: "prop-nonce",
			want:      `nonce="prop-nonce"`,
			notWant:   `nonce=""`,
		},
		{
			name:      "falls back to context nonce",
			ctxNonce:  "ctx-nonce",
			propNonce: "",
			want:      `nonce="ctx-nonce"`,
			notWant:   `nonce=""`,
		},
		{
			name:      "no nonce anywhere omits the attribute",
			ctxNonce:  "",
			propNonce: "",
			want:      "<script>",
			notWant:   `nonce=`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			if tt.ctxNonce != "" {
				ctx = templ.WithNonce(ctx, tt.ctxNonce)
			}

			var buf strings.Builder

			err := Alert(AlertProps{
				Message:     "nonce probe",
				Dismissible: true,
				BaseProps:   utils.BaseProps{Nonce: tt.propNonce},
			}).Render(ctx, &buf)
			if err != nil {
				t.Fatalf("render: %v", err)
			}

			out := buf.String()

			if tt.want != "" && !strings.Contains(out, tt.want) {
				t.Errorf("output missing %q:\n%s", tt.want, out)
			}

			if tt.notWant != "" && strings.Contains(out, tt.notWant) {
				t.Errorf("output contains %q:\n%s", tt.notWant, out)
			}
		})
	}
}
