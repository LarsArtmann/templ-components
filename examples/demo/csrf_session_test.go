package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// postWithSession builds a mutating request carrying a session cookie and a
// csrf_token form field, the way the kanban endpoints receive them.
func postWithSession(t *testing.T, cookie, formToken string) *http.Request {
	t.Helper()

	form := url.Values{}
	if formToken != "" {
		form.Set(kanbanCSRFFieldName, formToken)
	}

	r := httptest.NewRequest(http.MethodPost, "/api/kanban/htmx", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if cookie != "" {
		r.Header.Set("Cookie", demoSessionCookie+"="+cookie)
	}

	return r
}

// TestDemoSessionCSRF pins the session-scoped CSRF contract (#229): a first
// visit mints the SameSite=Lax cookie and exposes the token to handlers;
// the submitted form token validates only against THAT session's value — a
// token from another session, or no cookie at all, fails closed.
func TestDemoSessionCSRF(t *testing.T) {
	t.Parallel()

	var handlerToken string
	firstVisit := withDemoSession(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerToken = demoSessionCSRF(r)
	}))

	rec := httptest.NewRecorder()
	firstVisit.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if handlerToken == "" {
		t.Fatal("first visit exposed an empty token to the handler")
	}

	for _, c := range rec.Result().Cookies() {
		if c.Name != demoSessionCookie {
			continue
		}

		if c.Value != handlerToken {
			t.Errorf("cookie value %q != handler token %q", c.Value, handlerToken)
		}

		if c.SameSite != http.SameSiteLaxMode {
			t.Errorf("cookie SameSite = %v, want Lax", c.SameSite)
		}
	}

	t.Run("matching session token validates", func(t *testing.T) {
		t.Parallel()

		var valid bool
		withDemoSession(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			valid = demoSessionTokenValid(r)
		})).ServeHTTP(httptest.NewRecorder(), postWithSession(t, handlerToken, handlerToken))

		if !valid {
			t.Error("the session's own token must validate")
		}
	})

	t.Run("foreign session cookie fails closed", func(t *testing.T) {
		t.Parallel()

		var valid bool
		withDemoSession(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			valid = demoSessionTokenValid(r)
		})).ServeHTTP(httptest.NewRecorder(), postWithSession(t, "some-other-session", handlerToken))

		if valid {
			t.Error("a foreign session must not validate another session's token")
		}
	})

	t.Run("no cookie fails closed", func(t *testing.T) {
		t.Parallel()

		var valid bool
		withDemoSession(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			valid = demoSessionTokenValid(r)
		})).ServeHTTP(httptest.NewRecorder(), postWithSession(t, "", handlerToken))

		if valid {
			t.Error("no session cookie must fail validation even with a token in the form")
		}
	})
}
