package main

import (
	"net/http"
	"net/mail"
	"net/url"
	"strings"

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

func (s recipeLoginState) validationErrors() []formsValidationError {
	errs := make([]formsValidationError, 0, 2)
	if s.ErrEmail != "" {
		errs = append(errs, formsValidationError{Field: "email", Message: s.ErrEmail})
	}

	if s.ErrPassword != "" {
		errs = append(errs, formsValidationError{Field: "password", Message: s.ErrPassword})
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
		st.ErrPassword = "Your password must be at least " + itoa(recipePasswordMinLen) + " characters."
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

func (s recipeSignupState) validationErrors() []formsValidationError {
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
		st.ErrPassword = "Your password must be at least " + itoa(recipePasswordMinLen) + " characters."
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
// back to the page with a saved confirmation for the submitted section.
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

func itoa(n int) string {
	return strings.TrimSpace(strings.Repeat("", 0) + fmtInt(n))
}

func fmtInt(n int) string {
	return strconv.Itoa(n)
}
