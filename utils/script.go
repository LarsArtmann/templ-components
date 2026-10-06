package utils

import (
	"context"
	"fmt"
	"html"
	"io"

	"github.com/a-h/templ"
)

// ScriptAttrs returns the nonce attribute for an inline <script> (or
// <style>) tag, or nil when nonce is empty.
//
// Rendering `nonce=""` for an unset nonce silently kills the script on any
// strict-CSP page (a nonce attribute with the wrong value never matches the
// CSP header) while non-CSP pages gain nothing from it. Omitting the
// attribute on empty keeps the script working for non-CSP consumers and
// changes nothing for strict-CSP ones — an inline script without a valid
// nonce is blocked there either way. This is the same omit-empty rule the
// theme scripts and the Datastar live-region busy script follow
// (issues #7, #9, #24).
func ScriptAttrs(nonce string) templ.Attributes {
	if nonce == "" {
		return nil
	}

	return templ.Attributes{"nonce": nonce}
}

// ScriptComponent renders a CSP-safe `<script>` tag wrapping the given JS
// verbatim — THE canonical writer for component singleton scripts. New
// script-emitting components MUST use it instead of hand-rolling an
// Fprintf: it owns the omit-empty nonce rule (an empty nonce renders NO
// nonce attribute, never `nonce=""` — the dead-script class behind the
// 2026-10-01 demo outage), HTML-escapes the nonce value, and wraps write
// errors with the caller's errLabel for debugging. A nonce-carrying render
// is byte-identical to `<script nonce="…">\n<js></script>\n`.
func ScriptComponent(nonce, script, errLabel string) templ.Component {
	nonceAttr := ""
	if nonce != "" {
		nonceAttr = fmt.Sprintf(" nonce=\"%s\"", html.EscapeString(nonce))
	}

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := fmt.Fprintf(w, "<script%s>\n%s</script>\n", nonceAttr, script); err != nil {
			return fmt.Errorf("write %s: %w", errLabel, err)
		}

		return nil
	})
}
