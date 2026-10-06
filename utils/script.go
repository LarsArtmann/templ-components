package utils

import "github.com/a-h/templ"

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
