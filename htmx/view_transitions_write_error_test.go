package htmx

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type failingStyleWriter struct{}

func (failingStyleWriter) Write([]byte) (int, error) {
	return 0, errors.New("connection closed")
}

// TestStyleComponentPropagatesWriteErrors proves the CSP-safe <style> writer
// reports a failed write instead of dropping it.
func TestStyleComponentPropagatesWriteErrors(t *testing.T) {
	t.Parallel()

	err := styleComponent("", "::view-transition-old(root){}").
		Render(context.Background(), failingStyleWriter{})
	if err == nil {
		t.Fatal("styleComponent swallowed a write error")
	}

	if !strings.Contains(err.Error(), "view transitions style") {
		t.Errorf("error %q lacks the writer label", err)
	}
}
