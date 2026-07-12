package errorpage

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// TestErrorPageWithClassAndAriaLabel exercises the BaseProps.Class and
// BaseProps.AriaLabel branches in the templ templates — previously untested.
func TestErrorPageWithClassAndAriaLabel(t *testing.T) {
	t.Parallel()
	html := utils.Render(t, ErrorPage(ErrorPageProps{
		BaseProps: utils.BaseProps{Class: "my-error", AriaLabel: "Error dialog", Nonce: "n"},
		Family:    FamilyInfrastructure,
		Title:     "Oops",
	}))
	if !strings.Contains(html, "my-error") {
		t.Error("Class should propagate to root element")
	}
}

func TestErrorDetailWithNonceAndClass(t *testing.T) {
	t.Parallel()
	html := utils.Render(t, ErrorDetail(ErrorDetailProps{
		BaseProps: utils.BaseProps{Class: "detail-cls", Nonce: "x", AriaLabel: "Detail"},
		Family:    FamilyTransient,
	}))
	if !strings.Contains(html, "detail-cls") {
		t.Error("Class should propagate to ErrorDetail root")
	}
}

func TestErrorAlertDismissibleFalse(t *testing.T) {
	t.Parallel()
	html := utils.Render(t, ErrorAlert(ErrorAlertProps{
		Family:      FamilyCorruption,
		Title:       "Bug",
		Dismissible: false,
	}))
	if strings.Contains(html, "dismiss") {
		t.Error("Dismissible=false should not show dismiss button")
	}
}

func TestNotFound404NoSearchInput(t *testing.T) {
	t.Parallel()
	props := DefaultNotFound404Props()
	props.SearchAction = ""
	html := utils.Render(t, NotFound404(props))
	if strings.Contains(html, "<form") {
		t.Error("Empty SearchAction should not render search form")
	}
}

func TestNotFound404WithSearchPlaceholder(t *testing.T) {
	t.Parallel()
	props := DefaultNotFound404Props()
	props.SearchAction = "/api/search"
	props.SearchPlaceholder = "Type to search..."
	props.SearchInputName = "query"
	html := utils.Render(t, NotFound404(props))
	if !strings.Contains(html, "Type to search...") {
		t.Error("SearchPlaceholder should render")
	}
	if !strings.Contains(html, `name="query"`) {
		t.Error("SearchInputName should set the input name")
	}
}

func TestErrorHandlerWithOverrideNonNil(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(t.Context(), "GET", "/", nil)
	h := ErrorHandler(errors.New("boom"), ErrorHandlerConfig{
		Nonce: "n",
		Override: func(err error, p ErrorPageProps) *ErrorPageProps {
			p.Title = "Overridden"
			return &p
		},
	})
	h.ServeHTTP(w, r)
	if !strings.Contains(w.Body.String(), "Overridden") {
		t.Error("Override should replace title")
	}
}

func TestErrorHandlerWithOverrideNilResult(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(t.Context(), "GET", "/", nil)
	h := ErrorHandler(errors.New("boom"), ErrorHandlerConfig{
		Nonce: "n",
		Override: func(err error, p ErrorPageProps) *ErrorPageProps {
			return nil
		},
	})
	h.ServeHTTP(w, r)
	if w.Code < 400 {
		t.Error("Should still return error status")
	}
}

func TestWriteErrorPageDerivesStatusCodeFromFamily(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(t.Context(), "GET", "/", nil)
	props := DefaultErrorPageProps()
	props.Family = FamilyRejection
	props.StatusCode = 0
	WriteErrorPage(w, r, 0, props, "nonce")
	if w.Code != http.StatusBadRequest {
		t.Errorf("FamilyRejection should derive 400, got %d", w.Code)
	}
}

func TestParseFamilyCaseInsensitive(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  Family
	}{
		{"REJECTION", FamilyRejection},
		{"Conflict", FamilyConflict},
		{"TRANSIENT", FamilyTransient},
		{"Corruption", FamilyCorruption},
		{"INFRASTRUCTURE", FamilyInfrastructure},
		{"garbage", FamilyTransient},
		{"", FamilyTransient},
	}
	for _, tt := range tests {
		got := ParseFamily(tt.input)
		if got != tt.want {
			t.Errorf("ParseFamily(%q) = %s, want %s", tt.input, got, tt.want)
		}
	}
}

func TestExtractCauseChainNilError(t *testing.T) {
	t.Parallel()
	chain := ExtractCauseChain(nil, 5)
	if len(chain) != 0 {
		t.Errorf("Expected 0 causes for nil error, got %d", len(chain))
	}
}

func TestAllConstructorsRenderWithCoverage(t *testing.T) {
	t.Parallel()
	constructors := []struct {
		name  string
		props ErrorPageProps
	}{
		{"NotFound", NotFound()},
		{"Forbidden", Forbidden()},
		{"BadRequest", BadRequest("")},
		{"BadRequestMsg", BadRequest("custom bad")},
		{"Conflict", Conflict("")},
		{"ConflictMsg", Conflict("custom conflict")},
		{"ServiceUnavailable", ServiceUnavailable()},
		{"InternalError", InternalError()},
	}
	for _, c := range constructors {
		html := utils.Render(t, ErrorPage(c.props))
		if html == "" {
			t.Errorf("Constructor %s rendered empty", c.name)
		}
	}
}
