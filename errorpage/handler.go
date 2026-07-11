package errorpage

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"fmt"
	"html"
	"io"
	"log/slog"
	"net/http"

	"encoding/json/jsontext"

	"github.com/a-h/templ"
)

// ErrorHandlerConfig controls how ErrorHandler renders errors.
type ErrorHandlerConfig struct {
	// Nonce is used for CSP-compliant inline scripts.
	Nonce string

	// Override allows per-error customization of the ErrorPageProps
	// before rendering. When the returned pointer is non-nil, its values
	// replace the derived props. When nil, the original derived props are used.
	Override func(err error, props ErrorPageProps) *ErrorPageProps

	// HTMLShell wraps the error page in a minimal HTML document with
	// DOCTYPE, html, head, title, and body tags. Use when the error page
	// is served as a standalone HTTP response (not embedded in an existing layout).
	HTMLShell bool

	// JSON renders a JSON error response instead of HTML.
	// The response includes family, code, message, title, why, and fix fields.
	// Use for API endpoints or HTMX error handling.
	JSON bool

	// Lang sets the <html lang="..."> attribute when HTMLShell is true.
	// Defaults to "en" when empty.
	Lang string
}

// ErrorHandler returns an http.Handler that renders a go-error-family
// aware error page. Use it in your HTTP error handling:
//
//	http.HandleFunc("/api/...", func(w http.ResponseWriter, r *http.Request) {
//	    if err := doSomething(); err != nil {
//	        errorpage.ErrorHandler(err, errorpage.ErrorHandlerConfig{Nonce: nonce}).ServeHTTP(w, r)
//	        return
//	    }
//	    w.WriteHeader(http.StatusOK)
//	})
func ErrorHandler(err error, cfg ErrorHandlerConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		props := FromError(err)
		props.ShowTimestamp = true

		if cfg.Override != nil {
			if overridden := cfg.Override(err, props); overridden != nil {
				props = *overridden
			}
		}

		if cfg.Nonce != "" {
			props.Nonce = cfg.Nonce
		}

		statusCode := props.StatusCode
		if statusCode == 0 {
			statusCode = FamilyStatusCode(props.Family)
		}

		if cfg.JSON {
			writeJSONError(w, statusCode, props)
			return
		}

		if cfg.HTMLShell {
			title := props.Title
			if title == "" {
				title = fmt.Sprintf("Error %d", statusCode)
			}
			lang := cfg.Lang
			if lang == "" {
				lang = "en"
			}
			html, renderErr := renderShellToBuffer(r.Context(), title, lang, props)
			if renderErr != nil {
				slog.Error("error page render failed", "error", renderErr, "original_error", err)
				writeFallbackError(w, statusCode)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(statusCode)
			_, _ = w.Write(html)
			return
		}

		buf, renderErr := renderToBuffer(r.Context(), ErrorPage(props)) //nolint:contextcheck // intentional passthrough
		if renderErr != nil {
			slog.Error("error page render failed", "error", renderErr, "original_error", err)
			writeFallbackError(w, statusCode)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(statusCode)
		_, _ = w.Write(buf)
	})
}

// WriteError writes an error page to an http.ResponseWriter.
// Convenience wrapper around ErrorHandler for simpler usage.
func WriteError(w http.ResponseWriter, r *http.Request, err error, nonce string) {
	ErrorHandler(err, ErrorHandlerConfig{Nonce: nonce}).ServeHTTP(w, r) //nolint:exhaustruct // minimal config
}

// WriteErrorPage writes a pre-configured error page with the given HTTP status code.
// If statusCode is 0, the status code is derived from props.Family via FamilyStatusCode.
// Use with the pre-built constructors:
//
//	errorpage.WriteErrorPage(w, r, 0, errorpage.NotFound(), "")
func WriteErrorPage(w http.ResponseWriter, r *http.Request, statusCode int, props ErrorPageProps, nonce string) {
	if statusCode == 0 {
		statusCode = FamilyStatusCode(props.Family)
	}
	if nonce != "" {
		props.Nonce = nonce
	}
	buf, renderErr := renderToBuffer(r.Context(), ErrorPage(props)) //nolint:contextcheck // intentional passthrough
	if renderErr != nil {
		slog.Error("error page render failed", "error", renderErr, "status_code", statusCode)
		writeFallbackError(w, statusCode)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write(buf)
}

// WriteNotFound404 writes a NotFound404 page to an http.ResponseWriter with a
// 404 status code. Convenience wrapper for common 404 handler usage.
//
//	func main() {
//	    mux := http.NewServeMux()
//	    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//	        errorpage.WriteNotFound404(w, r, errorpage.DefaultNotFound404Props(), "")
//	    })
//	}
func WriteNotFound404(w http.ResponseWriter, r *http.Request, props NotFound404Props, nonce string) {
	if nonce != "" {
		props.Nonce = nonce
	}
	buf, renderErr := renderToBuffer(r.Context(), NotFound404(props)) //nolint:contextcheck // intentional passthrough
	if renderErr != nil {
		slog.Error("not found 404 page render failed", "error", renderErr)
		writeFallbackError(w, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write(buf)
}

func writeJSONError(w http.ResponseWriter, statusCode int, props ErrorPageProps) {
	resp := errorResponse{ //nolint:exhaustruct // Context set conditionally below
		Family:  string(props.Family),
		Code:    string(props.Code),
		Message: props.Message,
		Title:   props.Title,
		Why:     props.Why,
		Fix:     props.Fix,
	}
	if len(props.Context) > 0 {
		ctx := make(map[string]string, len(props.Context))
		for _, p := range props.Context {
			ctx[p.Key] = p.Value
		}
		resp.Context = ctx
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	enc := jsontext.NewEncoder(w)
	enc.SetEscapeHTML(true)
	if err := json.MarshalEncode(enc, resp); err != nil {
		slog.Error("JSON error response encode failed", "error", err)
	}
}

// renderToBuffer renders a templ component to a byte slice.
// This ensures the full document renders successfully before writing
// any bytes to the http.ResponseWriter, preventing truncated HTML
// at the wrong status code if rendering fails mid-stream.
func renderToBuffer(ctx context.Context, comp templ.Component) ([]byte, error) {
	var buf bytes.Buffer
	if err := comp.Render(ctx, &buf); err != nil {
		return nil, fmt.Errorf("render component: %w", err)
	}
	return buf.Bytes(), nil
}

// renderShellToBuffer renders the error page wrapped in a minimal HTML
// document to a byte slice, ensuring the full document renders before
// writing to the ResponseWriter.
func renderShellToBuffer(ctx context.Context, title, lang string, props ErrorPageProps) ([]byte, error) {
	shell := templ.ComponentFunc(func(_ context.Context, bw io.Writer) error {
		_, _ = fmt.Fprint(bw, `<!DOCTYPE html>`)
		_, _ = fmt.Fprintf(bw, `<html lang="%s"><head>`, html.EscapeString(lang))
		_, _ = fmt.Fprint(bw, `<meta charset="UTF-8">`)
		_, _ = fmt.Fprint(bw, `<meta name="viewport" content="width=device-width, initial-scale=1.0">`)
		_, _ = fmt.Fprintf(bw, `<title>%s</title>`, html.EscapeString(title))
		_, _ = fmt.Fprint(bw, `</head><body>`)
		renderErr := ErrorPage(props).Render(ctx, bw) //nolint:contextcheck // intentional passthrough
		if renderErr != nil {
			return fmt.Errorf("render error page: %w", renderErr)
		}
		_, _ = fmt.Fprint(bw, `</body></html>`)
		return nil
	})
	var buf bytes.Buffer
	if err := shell.Render(ctx, &buf); err != nil {
		return nil, fmt.Errorf("render error page shell (title=%q): %w", title, err)
	}
	return buf.Bytes(), nil
}

// writeFallbackError writes a minimal plain-text error response when
// the templ error page itself fails to render. This ensures the client
// always receives a response with the correct status code.
func writeFallbackError(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// If headers haven't been written yet, WriteHeader will succeed.
	// If they have (e.g. superedge case), this is a no-op.
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprintf(w, "Error %d\n", statusCode)
}

// Verify interface compliance.
var _ http.Handler = ErrorHandler(nil, ErrorHandlerConfig{}) //nolint:exhaustruct // type check only
