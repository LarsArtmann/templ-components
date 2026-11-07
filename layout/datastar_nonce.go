package layout

import (
	"github.com/a-h/templ"
)

// DatastarNonceAttrs returns the <html>-level attributes that activate the
// Datastar runtime's CSP mode (docs/datastar-runtime-facts.md §CSP): the
// runtime reads data-nonce off document.documentElement once at init and,
// when it is present and nonempty, executes every expression as a nonce'd
// inline script instead of Function() eval — so a nonce-based script-src CSP
// needs no 'unsafe-eval'. Spread the result into PageProps.HTMLDataAttrs and
// pass the SAME nonce to datastar.SDKScriptProps.Nonce; one without the other
// is broken (script without flag → eval attempt → CSP violation; flag without
// nonced script → the runtime itself is blocked).
//
// An empty nonce returns nil so <html data-nonce=""> can never render: the
// runtime treats an empty attribute as a misconfiguration and THROWS at init
// ("Datastar CSP requires a nonempty html data-nonce."), which would take
// every data-* expression on the page down with it. Same omit-empty contract
// as utils.ScriptAttrs/ScriptComponent for an empty nonce.
func DatastarNonceAttrs(nonce string) templ.Attributes {
	if nonce == "" {
		return nil
	}

	return templ.Attributes{"data-nonce": nonce}
}
