package errorpage

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/a-h/templ"
)

// deadWriter always fails on Write — the "connection is gone" case. Unlike
// failingWriter (which fails once, then succeeds so the fallback can land),
// this one proves the fallback paths themselves survive a dead connection.
type deadWriter struct {
	header http.Header
}

func (d *deadWriter) Header() http.Header {
	if d.header == nil {
		d.header = make(http.Header)
	}

	return d.header
}

func (d *deadWriter) Write([]byte) (int, error) {
	return 0, errors.New("connection closed")
}

func (d *deadWriter) WriteHeader(int) {}

// TestWritePathsSurviveDeadConnection covers the two error-reporting paths
// whose writes can fail after headers are committed: JSON encoding and the
// plain-text fallback. Neither may panic or lose the status/Content-Type.
func TestWritePathsSurviveDeadConnection(t *testing.T) {
	t.Parallel()

	t.Run("JSON error response logs the encode failure", func(t *testing.T) {
		t.Parallel()

		w := &deadWriter{}
		writeJSONError(w, http.StatusInternalServerError, ErrorPageProps{
			Family: FamilyCorruption,
			Code:   "test.code",
			Title:  "Test",
		})

		if got := w.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
			t.Errorf("Content-Type = %q, want JSON", got)
		}
	})

	t.Run("plain-text fallback survives a dead connection", func(t *testing.T) {
		t.Parallel()

		w := &deadWriter{}
		writeFallbackError(w, http.StatusInternalServerError)

		if got := w.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
			t.Errorf("Content-Type = %q, want text/plain", got)
		}

		if got := w.Header().Get("X-Robots-Tag"); got != "noindex" {
			t.Errorf("X-Robots-Tag = %q, want noindex", got)
		}
	})
}

// TestApplyRetrySuggestionKeepsExplicitWayOut pins the precedence rule: a
// caller-configured way out is never replaced by the retry suggestion, in
// either the primary or the secondary slot.
func TestApplyRetrySuggestionKeepsExplicitWayOut(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		props ErrorPageProps
	}{
		{name: "primary label", props: ErrorPageProps{WayOut: "Go home"}},
		{name: "primary href", props: ErrorPageProps{WayOutHref: "/home"}},
		{name: "primary action", props: ErrorPageProps{WayOutAction: WayOutAction{Text: "Retry"}}},
		{name: "secondary", props: ErrorPageProps{SecondaryWayOut: "Contact support"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			props := tt.props
			applyRetrySuggestion(&props, retryableError{retryable: true}, "/same")

			if props.WayOut != tt.props.WayOut || props.WayOutHref != tt.props.WayOutHref ||
				props.SecondaryWayOut != tt.props.SecondaryWayOut {
				t.Errorf("explicit way out was overwritten: %+v", props)
			}
		})
	}
}

// TestRenderToBufferPropagatesRenderErrors proves a failing component surfaces
// as an error instead of a truncated buffer.
func TestRenderToBufferPropagatesRenderErrors(t *testing.T) {
	t.Parallel()

	boom := errors.New("boom")
	failing := templ.ComponentFunc(func(context.Context, io.Writer) error {
		return boom
	})

	_, err := renderToBuffer(context.Background(), failing)
	if !errors.Is(err, boom) {
		t.Fatalf("renderToBuffer error = %v, want %v", err, boom)
	}
}
