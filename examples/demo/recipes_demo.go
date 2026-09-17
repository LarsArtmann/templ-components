package main

import (
	"net/http"
	"net/mail"
	"net/url"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/layout"
)

// The recipe demo pages are interactive: login, signup, and settings forms
// POST back to their own routes and re-render with server-side validation
// errors or a success state. These types carry that state between the HTTP
// handlers and the templ templates.

const recipePasswordMinLen = 8

// recipeLoginState is the server-rendered state of the /recipes/login form.
type recipeLoginState struct {
	// Email is echoed back into the email field after a failed submit.
	Email string
	// Success marks a completed (demo) sign-in: renders the success banner.
	Success bool
	// ErrEmail / ErrPassword are per-field error messages.
	ErrEmail    string
	ErrPassword string
}

func (s recipeLoginState) hasErrors() bool {
	return s.ErrEmail != "" || s.ErrPassword != ""
}

func (s recipeLoginState) validationErrors() []forms.ValidationError {
	errs := make([]forms.ValidationError, 0, 2)
	if s.ErrEmail != "" {
		errs = append(errs, forms.ValidationError{Field: "email", Message: s.ErrEmail})
	}

	if s.ErrPassword != "" {
		errs = append(errs, forms.ValidationError{Field: "password", Message: s.ErrPassword})
	}

	return errs
}

// parseRecipeLogin validates a login POST. Any well-formed email plus an
// 8+ character password "succeeds" — the demo persists nothing.
func parseRecipeLogin(r *http.Request) recipeLoginState {
	st := recipeLoginState{Email: strings.TrimSpace(r.PostFormValue("email"))}

	if st.Email == "" {
		st.ErrEmail = "Enter your email address."
	} else if _, err := mail.ParseAddress(st.Email); err != nil {
		st.ErrEmail = "Enter a valid email address, e.g. you@example.com."
	}

	if password := r.PostFormValue("password"); len(password) < recipePasswordMinLen {
		st.ErrPassword = "Your password must be at least " + strconv.Itoa(recipePasswordMinLen) + " characters."
	}

	if !st.hasErrors() {
		st.Success = true
	}

	return st
}

// recipeSignupState is the server-rendered state of the /recipes/auth form.
type recipeSignupState struct {
	Name  string
	Email string
	// Success marks a completed (demo) account creation.
	Success bool
	// ErrName / ErrEmail / ErrPassword / ErrTerms are per-field errors.
	ErrName     string
	ErrEmail    string
	ErrPassword string
	ErrTerms    string
}

func (s recipeSignupState) hasErrors() bool {
	return s.ErrName != "" || s.ErrEmail != "" || s.ErrPassword != "" || s.ErrTerms != ""
}

func (s recipeSignupState) validationErrors() []forms.ValidationError {
	errs := make([]forms.ValidationError, 0, 4)
	if s.ErrName != "" {
		errs = append(errs, forms.ValidationError{Field: "full_name", Message: s.ErrName})
	}

	if s.ErrEmail != "" {
		errs = append(errs, forms.ValidationError{Field: "email", Message: s.ErrEmail})
	}

	if s.ErrPassword != "" {
		errs = append(errs, forms.ValidationError{Field: "password", Message: s.ErrPassword})
	}

	if s.ErrTerms != "" {
		errs = append(errs, forms.ValidationError{Field: "terms", Message: s.ErrTerms})
	}

	return errs
}

func parseRecipeSignup(r *http.Request) recipeSignupState {
	st := recipeSignupState{
		Name:  strings.TrimSpace(r.PostFormValue("full_name")),
		Email: strings.TrimSpace(r.PostFormValue("email")),
	}

	if st.Name == "" {
		st.ErrName = "Enter your name."
	}

	if st.Email == "" {
		st.ErrEmail = "Enter your work email."
	} else if _, err := mail.ParseAddress(st.Email); err != nil {
		st.ErrEmail = "Enter a valid email address, e.g. you@company.com."
	}

	if password := r.PostFormValue("password"); len(password) < recipePasswordMinLen {
		st.ErrPassword = "Your password must be at least " + strconv.Itoa(recipePasswordMinLen) + " characters."
	}

	if r.PostFormValue("terms") == "" {
		st.ErrTerms = "Please accept the terms of service to continue."
	}

	if !st.hasErrors() {
		st.Success = true
	}

	return st
}

// recipeSettingsSections are the known settings sections; the section hidden
// input must carry one of these IDs.
var recipeSettingsSections = map[string]bool{
	"profile":       true,
	"security":      true,
	"notifications": true,
}

// handleRecipeSettingsSave accepts a settings section form POST and redirects
// back to the page with a saved confirmation for the submitted section. The
// form itself persists nothing — the redirect loop exists to demonstrate a
// full POST/redirect/GET round trip against library form components.
func handleRecipeSettingsSave(w http.ResponseWriter, r *http.Request) {
	section := r.PostFormValue("section")
	if !recipeSettingsSections[section] {
		http.Error(w, "unknown settings section", http.StatusBadRequest)

		return
	}

	http.Redirect(w, r, "/recipes/settings?saved="+url.QueryEscape(section)+"#"+section, http.StatusSeeOther)
}

// recipeSavedSection resolves the ?saved= query parameter to a known section
// ID, or "" when absent/unknown.
func recipeSavedSection(raw string) string {
	if recipeSettingsSections[raw] {
		return raw
	}

	return ""
}

// renderRecipeSettings serves GET (page, optionally with the saved banner)
// and POST (save redirect) for the settings recipe demo.
func renderRecipeSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		handleRecipeSettingsSave(w, r)

		return
	}

	renderRecipePageState(w, r, "Settings Recipe - templ-components", "Settings recipe demo",
		func(props layout.PageProps) templ.Component {
			return recipesSettingsPageState(props, recipeSavedSection(r.URL.Query().Get("saved")))
		})
}

// renderRecipeLogin serves GET (empty form) and POST (validated form) for
// the login recipe demo.
func renderRecipeLogin(w http.ResponseWriter, r *http.Request) {
	st := recipeLoginState{}
	if r.Method == http.MethodPost {
		st = parseRecipeLogin(r)
	}

	renderRecipePageState(w, r, "Login Recipe - templ-components", "Login card recipe demo",
		func(props layout.PageProps) templ.Component { return recipesLoginPageState(props, st) })
}

// renderRecipeAuth serves GET (empty form) and POST (validated form) for the
// auth layout recipe demo.
func renderRecipeAuth(w http.ResponseWriter, r *http.Request) {
	st := recipeSignupState{}
	if r.Method == http.MethodPost {
		st = parseRecipeSignup(r)
	}

	renderRecipePageState(w, r, "Auth Layout Recipe - templ-components", "Auth layout recipe demo",
		func(props layout.PageProps) templ.Component { return recipesAuthPageState(props, st) })
}

func renderRecipePageState(
	w http.ResponseWriter,
	r *http.Request,
	title, description string,
	page func(layout.PageProps) templ.Component,
) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := page(demoPageProps(title, description)).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
