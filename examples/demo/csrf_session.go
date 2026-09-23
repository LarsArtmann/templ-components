package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// Session-scoped CSRF for the demo's mutating endpoints (backlog #229,
// 2026-09-23). A real app issues one token per session and validates it
// server-side on every mutating request; the demo used to mint ONE random
// token per process, which protected nothing beyond the process boundary.
// Now every visitor gets a random token in a SameSite=Lax cookie on their
// first response, and mutating endpoints compare the submitted csrf_token
// against that cookie — a token scraped from yesterday's page render no
// longer validates.
//
// Still demo-grade by design: no HMAC binding to the session, no expiry, in
// a public demo holding no real data. The shape matches what a real app
// should grow into.

const demoSessionCookie = "tc_demo_csrf"

type demoSessionCtxKey struct{}

// withDemoSession middleware: mint the CSRF cookie on first response and
// expose the token to handlers via the request context.
func withDemoSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie(demoSessionCookie); err == nil && c.Value != "" {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), demoSessionCtxKey{}, c.Value)))

			return
		}

		token := newDemoSessionToken()
		http.SetCookie(w, &http.Cookie{
			Name:     demoSessionCookie,
			Value:    token,
			Path:     "/",
			MaxAge:   86400,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   true,
		})
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), demoSessionCtxKey{}, token)))
	})
}

// demoSessionCSRF returns the visitor's session token (empty when the
// middleware has not run or the visitor sent no cookie yet — a mutating
// endpoint must treat empty as a failed validation, since the form could not
// have legitimately carried a token).
func demoSessionCSRF(r *http.Request) string {
	token, _ := r.Context().Value(demoSessionCtxKey{}).(string)

	return token
}

// demoSessionTokenValid reports whether the submitted form token matches the
// visitor's session cookie token.
func demoSessionTokenValid(r *http.Request) bool {
	session := demoSessionCSRF(r)

	return session != "" && r.FormValue(kanbanCSRFFieldName) == session
}

// newDemoSessionToken mints one random session token.
func newDemoSessionToken() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic("demo session: generate csrf token: " + err.Error())
	}

	return hex.EncodeToString(b[:])
}
