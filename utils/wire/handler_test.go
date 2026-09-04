package wire

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPatchModeIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value PatchMode
		valid bool
	}{
		{name: "unspecified is the defined zero value", value: PatchModeUnspecified, valid: true},
		{name: "inner", value: PatchModeInner, valid: true},
		{name: "outer", value: PatchModeOuter, valid: true},
		{name: "prepend", value: PatchModePrepend, valid: true},
		{name: "append", value: PatchModeAppend, valid: true},
		{name: "before", value: PatchModeBefore, valid: true},
		{name: "after", value: PatchModeAfter, valid: true},
		{name: "replace", value: PatchModeReplace, valid: true},
		{name: "unknown", value: PatchMode("upsert"), valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := PatchModeIsValid(tt.value); got != tt.valid {
				t.Errorf("PatchModeIsValid(%q) = %v, want %v", tt.value, got, tt.valid)
			}
		})
	}
}

func TestRequestPredicates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		headers      map[string]string
		wantDatastar bool
		wantHTMX     bool
	}{
		{
			name:         "plain browser request has no transport markers",
			headers:      nil,
			wantDatastar: false,
			wantHTMX:     false,
		},
		{
			name:         "datastar request marker",
			headers:      map[string]string{HeaderDatastarRequest: "true"},
			wantDatastar: true,
			wantHTMX:     false,
		},
		{
			name:         "htmx request marker",
			headers:      map[string]string{HeaderHXRequest: "true"},
			wantDatastar: false,
			wantHTMX:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/api/fragment", nil)
			for k, v := range tt.headers {
				r.Header.Set(k, v)
			}

			if got := IsDatastar(r); got != tt.wantDatastar {
				t.Errorf("IsDatastar() = %v, want %v", got, tt.wantDatastar)
			}

			if got := IsHTMX(r); got != tt.wantHTMX {
				t.Errorf("IsHTMX() = %v, want %v", got, tt.wantHTMX)
			}
		})
	}
}

func TestHandler(t *testing.T) {
	t.Parallel()

	const body = "<p>fragment</p>"
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, body)
	})

	tests := []struct {
		name         string
		target       PatchTarget
		datastarCall bool
		wantSelector string
		wantMode     string
	}{
		{
			name:         "datastar caller gets default inner mode",
			target:       PatchTarget{Selector: "#out"},
			datastarCall: true,
			wantSelector: "#out",
			wantMode:     "inner",
		},
		{
			name:         "datastar caller gets explicit mode",
			target:       PatchTarget{Selector: "#out", Mode: PatchModeAppend},
			datastarCall: true,
			wantSelector: "#out",
			wantMode:     "append",
		},
		{
			name:         "htmx caller gets no routing headers",
			target:       PatchTarget{Selector: "#out"},
			datastarCall: false,
			wantSelector: "",
			wantMode:     "",
		},
		{
			name:         "plain browser caller gets no routing headers",
			target:       PatchTarget{Selector: "#out"},
			datastarCall: false,
			wantSelector: "",
			wantMode:     "",
		},
		{
			name:         "empty selector degrades to id-matched patching",
			target:       PatchTarget{Mode: PatchModeOuter},
			datastarCall: true,
			wantSelector: "",
			wantMode:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(Handler(tt.target, next))
			t.Cleanup(srv.Close)

			req, err := http.NewRequest(http.MethodGet, srv.URL+"/api/fragment", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.datastarCall {
				req.Header.Set(HeaderDatastarRequest, "true")
			}

			resp, err := srv.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = resp.Body.Close() })

			if resp.StatusCode != http.StatusOK {
				t.Errorf("status = %d, want 200", resp.StatusCode)
			}

			if got := resp.Header.Get(HeaderDatastarSelector); got != tt.wantSelector {
				t.Errorf("Datastar-Selector = %q, want %q", got, tt.wantSelector)
			}

			if got := resp.Header.Get(HeaderDatastarMode); got != tt.wantMode {
				t.Errorf("Datastar-Mode = %q, want %q", got, tt.wantMode)
			}

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if string(respBody) != body {
				t.Errorf("body = %q, want %q", respBody, body)
			}
		})
	}
}
