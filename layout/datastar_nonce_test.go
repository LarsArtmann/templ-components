package layout

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

func TestDatastarNonceAttrs(t *testing.T) {
	t.Parallel()

	t.Run("empty nonce degrades to nil (no data-nonce attribute)", func(t *testing.T) {
		t.Parallel()

		if attrs := DatastarNonceAttrs(""); attrs != nil {
			t.Errorf("empty nonce must return nil, got %#v (data-nonce=\"\" throws the runtime at init)", attrs)
		}
	})

	t.Run("nonempty nonce maps to the activation attribute", func(t *testing.T) {
		t.Parallel()

		attrs := DatastarNonceAttrs("abc123")
		if attrs == nil || attrs["data-nonce"] != "abc123" {
			t.Errorf(`want templ.Attributes{"data-nonce": "abc123"}, got %#v`, attrs)
		}
	})

	t.Run("rendered Base carries data-nonce on html and omits it when unset", func(t *testing.T) {
		t.Parallel()

		withNonce := DefaultPageProps()
		withNonce.HTMLDataAttrs = DatastarNonceAttrs("n0nce-value")

		rendered := utils.Render(t, Base(withNonce))
		utils.AssertContains(t, rendered, `data-nonce="n0nce-value"`)

		withoutNonce := DefaultPageProps()
		bare := utils.Render(t, Base(withoutNonce))
		utils.AssertNotContains(t, bare, `data-nonce=""`)

		emptyAttrs := DefaultPageProps()
		emptyAttrs.HTMLDataAttrs = DatastarNonceAttrs("")

		if utils.Render(t, Base(emptyAttrs)) != bare {
			t.Error("empty nonce must render byte-identically to no attributes at all")
		}
	})
}

func TestGoldenSweepBaseDatastarNonce(t *testing.T) {
	t.Parallel()

	props := DefaultPageProps()
	props.HTMLDataAttrs = DatastarNonceAttrs("csp-nonce-1")

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "base_datastar_nonce_set", HTML: utils.Render(t, Base(props))},
	})
}
