package errorpage

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/internal/golden"
	"github.com/larsartmann/templ-components/utils"
)

// TestErrorHeaderSubTemplate verifies the shared errorHeader sub-template renders
// correctly when composed into both ErrorPage and ErrorDetail.
func TestErrorHeaderSubTemplate(t *testing.T) {
	t.Parallel()

	// ErrorPage renders errorHeader internally
	props := ErrorPageProps{
		BaseProps: utils.BaseProps{AriaLabel: "Error"},
		Title:     "Test Error",
		Message:   "Something broke",
	}
	result := utils.Render(t, ErrorPage(props))
	if !strings.Contains(result, "Test Error") {
		t.Error("errorHeader: title not rendered")
	}
	if !strings.Contains(result, "Something broke") {
		t.Error("errorHeader: message not rendered")
	}
}

// TestActionLinkBodySubTemplate verifies the shared actionLinkBody renders text + arrow icon.
func TestActionLinkBodySubTemplate(t *testing.T) {
	t.Parallel()

	props := ErrorPageProps{
		WayOut:     "Go home",
		WayOutHref: "/",
	}
	result := utils.Render(t, ErrorPage(props))
	if !strings.Contains(result, "Go home") {
		t.Error("actionLinkBody: text not rendered")
	}
	if !strings.Contains(result, "<svg") {
		t.Error("actionLinkBody: arrow icon SVG not rendered")
	}
}

// TestGoBackScriptSubTemplate verifies the shared goBackScript renders a nonce'd script.
func TestGoBackScriptSubTemplate(t *testing.T) {
	t.Parallel()

	props := ErrorPageProps{
		BaseProps: utils.BaseProps{Nonce: "test-nonce-123"},
		WayOut:    "Go back",
	}
	result := utils.Render(t, ErrorPage(props))
	if !strings.Contains(result, "test-nonce-123") {
		t.Error("goBackScript: nonce not propagated")
	}
	if !strings.Contains(result, "history.back") {
		t.Error("goBackScript: history.back() not present")
	}
}

// TestNotFound404GoBackScript verifies goBackScript is also used by NotFound404.
func TestNotFound404GoBackScript(t *testing.T) {
	t.Parallel()

	props := NotFound404Props{
		ShowGoBack: true,
		BaseProps:  utils.BaseProps{Nonce: "404-nonce"},
	}
	result := utils.Render(t, NotFound404(props))
	if !strings.Contains(result, "404-nonce") {
		t.Error("NotFound404: nonce not rendered in goBackScript")
	}
	if !strings.Contains(result, "history.back") {
		t.Error("NotFound404: history.back() not present")
	}
}

// TestGoldenErrorHeaderConsistency ensures errorHeader output is stable.
func TestGoldenErrorHeaderConsistency(t *testing.T) {
	t.Parallel()

	props := ErrorPageProps{
		Title:   "Service Unavailable",
		Message: "Try again in a moment",
	}
	got := utils.Render(t, ErrorPage(props))
	golden.Assert(t, "error_header_consistency", got)
}
