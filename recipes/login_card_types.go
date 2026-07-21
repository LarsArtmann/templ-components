package recipes

import (
	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// LoginCardProps configures a recipes.LoginCard screen — the canonical
// centered Card with email + password form, optional OAuth slots, and
// "remember me" / "forgot password" affordances.
//
// The form body and OAuth buttons are consumer-supplied templ.Component
// slots so the consumer controls the actual input fields, validation, and
// HTMX wiring. LoginCard provides the layout: a centered max-w card with
// a title header and a vertical form stack body.
type LoginCardProps struct {
	utils.BaseProps

	// Title renders as the card header. Default: "Sign in".
	Title string
	// Subtitle renders as muted text under the Title. Optional.
	Subtitle string
	// FormBody is the main form content (email + password + submit button).
	// Typically a forms.Form with Layout: forms.FormLayoutStack. Required.
	FormBody templ.Component
	// OAuthButtons renders below the form body, separated by a divider.
	// Use it for "Continue with Google/GitHub" buttons. Optional.
	OAuthButtons templ.Component
	// Footer renders below everything (typically "Don't have an account?
	// Sign up"). Optional.
	Footer templ.Component
}

// DefaultLoginCardProps returns sensible defaults (Title: "Sign in").
//
//nolint:exhaustruct // constructor intentionally sets only non-zero defaults
func DefaultLoginCardProps() LoginCardProps {
	return LoginCardProps{
		Title: loginDefaultTitle,
	}
}

const loginDefaultTitle = "Sign in"
