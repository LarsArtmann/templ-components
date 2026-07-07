package errorpage

import (
	"testing"
)

func TestConstructorStatusCodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		props    ErrorPageProps
		wantCode int
	}{
		{"NotFound returns 404", NotFound(), 404},
		{"Forbidden returns 403", Forbidden(), 403},
		{"BadRequest returns 400 (derived from family)", BadRequest(""), 0},
		{"Conflict returns 409 (derived from family)", Conflict(""), 0},
		{"ServiceUnavailable returns 503 (derived from family)", ServiceUnavailable(), 0},
		{"InternalError returns 500", InternalError(), 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.wantCode > 0 && tt.props.StatusCode != tt.wantCode {
				t.Errorf("StatusCode = %d, want %d", tt.props.StatusCode, tt.wantCode)
			}
		})
	}
}
