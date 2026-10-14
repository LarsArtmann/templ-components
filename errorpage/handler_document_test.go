package errorpage

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestWriteNotFound404Document(t *testing.T) {
	t.Parallel()

	t.Run("writes a valid standalone HTML document with 404", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		r := httptest.NewRequestWithContext(context.Background(), "GET", "/missing", nil)
		WriteNotFound404Document(w, r, DefaultNotFound404Props(), "test-nonce", "")

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
		}

		body := w.Body.String()
		for _, want := range []string{
			"<!DOCTYPE html>",
			`<html lang="en">`,
			`<meta name="robots" content="noindex">`,
			"<title>Page not found</title>",
			">404<",
			"Go to homepage",
			"</body></html>",
		} {
			if !strings.Contains(body, want) {
				t.Errorf("document body missing %q", want)
			}
		}
	})

	t.Run("lang override reaches the html element", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		r := httptest.NewRequestWithContext(context.Background(), "GET", "/missing", nil)
		WriteNotFound404Document(w, r, DefaultNotFound404Props(), "", "de-DE")

		utils.AssertContains(t, w.Body.String(), `<html lang="de-DE">`)
	})

	t.Run("empty title falls back to the 404 default", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		r := httptest.NewRequestWithContext(context.Background(), "GET", "/missing", nil)
		WriteNotFound404Document(w, r, NotFound404Props{ShowGoBack: false}, "", "en")

		utils.AssertContains(t, w.Body.String(), "<title>Page not found</title>")
	})

	t.Run("nonce flows into the go-back script", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		r := httptest.NewRequestWithContext(context.Background(), "GET", "/missing", nil)
		WriteNotFound404Document(w, r, DefaultNotFound404Props(), "doc-nonce", "en")

		body := w.Body.String()
		utils.AssertContains(t, body, `nonce="doc-nonce"`)
		utils.AssertContains(t, body, "history.back()")
	})
}

// TestXRobotsTagHeaderInvariant pins the module-wide rule: every HTML error
// response carries X-Robots-Tag: noindex (soft-404 defense), while the JSON
// error path stays header-free.
func TestXRobotsTagHeaderInvariant(t *testing.T) {
	t.Parallel()

	newRequest := func() *http.Request {
		return httptest.NewRequestWithContext(context.Background(), "GET", "/missing", nil)
	}

	tests := []struct {
		name   string
		write  func(w http.ResponseWriter, r *http.Request)
		noindx bool
	}{
		{"WriteErrorPage", func(w http.ResponseWriter, r *http.Request) {
			WriteErrorPage(w, r, http.StatusNotFound, ErrorPageProps{Family: FamilyTransient}, "")
		}, true},
		{"WriteNotFound404", func(w http.ResponseWriter, r *http.Request) {
			WriteNotFound404(w, r, DefaultNotFound404Props(), "")
		}, true},
		{"WriteNotFound404Document", func(w http.ResponseWriter, r *http.Request) {
			WriteNotFound404Document(w, r, DefaultNotFound404Props(), "", "")
		}, true},
		{"ErrorHandler HTML", func(w http.ResponseWriter, r *http.Request) {
			ErrorHandler(errTest{}, ErrorHandlerConfig{}).ServeHTTP(w, r)
		}, true},
		{"ErrorHandler HTMLShell", func(w http.ResponseWriter, r *http.Request) {
			ErrorHandler(errTest{}, ErrorHandlerConfig{HTMLShell: true}).ServeHTTP(w, r)
		}, true},
		{"ErrorHandler JSON", func(w http.ResponseWriter, r *http.Request) {
			ErrorHandler(errTest{}, ErrorHandlerConfig{JSON: true}).ServeHTTP(w, r)
		}, false},
		{"fallback error", func(w http.ResponseWriter, r *http.Request) {
			writeFallbackError(w, http.StatusInternalServerError)
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			tt.write(w, newRequest())

			got := w.Header().Get("X-Robots-Tag")
			if tt.noindx && got != "noindex" {
				t.Errorf("X-Robots-Tag = %q, want %q", got, "noindex")
			}

			if !tt.noindx && got != "" {
				t.Errorf("X-Robots-Tag = %q, want empty (JSON path stays header-free)", got)
			}
		})
	}
}

type errTest struct{}

func (errTest) Error() string { return "boom" }
