package layout

import (
	"context"
	_ "embed"
	"fmt"
	"io"

	"github.com/a-h/templ"
)

//go:embed static/htmx.min.js
var htmxSource string

// HTMXSelfHost is the sentinel value for PageProps.HTMXSrc that triggers
// inline embedding of the HTMX source (v2.0 default). When HTMXSrc is set
// to HTMXSelfHost, the HTMX script is rendered inline via a <script> tag
// with the embedded source, avoiding any external CDN request.
const HTMXSelfHost = "self"

// HTMXNone is the sentinel value for PageProps.HTMXSrc that renders NO htmx
// runtime at all (backlog #282). Choose it when a page uses no htmx-powered
// components (or wires everything through Datastar): Base then emits no
// htmx <script>, no CDN preconnect/dns-prefetch hints, and no SRI — smaller
// pages, zero framework bytes. Every library component degrades gracefully
// without the runtime: wired anchors/forms fall back to plain navigation or
// no-op, and unwired markup is unaffected.
const HTMXNone = "none"

// htmxSelfHostComponent returns a templ.Component that renders the embedded
// HTMX source inline in a nonce-protected <script> tag. Used when
// PageProps.HTMXSrc == HTMXSelfHost. Writing raw HTML via ComponentFunc is
// required because templ's <script> context sanitizes interpolation.
func htmxSelfHostComponent(nonce string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := fmt.Fprintf(w, `<script nonce="%s">%s</script>`, nonce, htmxSource)
		if err != nil {
			return fmt.Errorf("write inline htmx script: %w", err)
		}

		return nil
	})
}
