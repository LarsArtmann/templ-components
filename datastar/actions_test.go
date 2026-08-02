package datastar

import (
	"testing"
)

func TestActionExpressions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		got  string
		want string
	}{
		{"Get", Get("/api/search"), "@get('/api/search')"},
		{"Post", Post("/api/save"), "@post('/api/save')"},
		{"Put", Put("/api/update/1"), "@put('/api/update/1')"},
		{"Patch", Patch("/api/patch/1"), "@patch('/api/patch/1')"},
		{"Delete", Delete("/api/delete/1"), "@delete('/api/delete/1')"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.got != tt.want {
				t.Errorf("%s(%q) = %q, want %q", tt.name, tt.want, tt.got, tt.want)
			}
		})
	}
}

func TestActionExpressionEscapesSingleQuotes(t *testing.T) {
	t.Parallel()

	got := Get("/api/search?q=it's")
	want := "@get('/api/search?q=it\\'s')"
	if got != want {
		t.Errorf("Get with single quote = %q, want %q", got, want)
	}
}

func TestIndicatorSignalExpr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		signal string
		want   string
	}{
		{"saving", "$saving"},
		{"fetching", "$fetching"},
		{"", "$"},
	}
	for _, tt := range tests {
		got := indicatorSignalExpr(tt.signal)
		if got != tt.want {
			t.Errorf("indicatorSignalExpr(%q) = %q, want %q", tt.signal, got, tt.want)
		}
	}
}

func TestDatastarScriptURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		version DatastarVersion
		cdn     string
		want    string
	}{
		{
			name:    "default CDN and version",
			version: DatastarVersion1_0_2,
			cdn:     "",
			want:    "https://cdn.jsdelivr.net/gh/starfederation/[email protected]/bundles/datastar.js",
		},
		{
			name:    "custom CDN",
			version: DatastarVersion1_0_2,
			cdn:     "https://unpkg.com",
			want:    "https://unpkg.com/starfederation/[email protected]/bundles/datastar.js",
		},
		{
			name:    "custom CDN with trailing slash",
			version: DatastarVersion1_0_2,
			cdn:     "https://unpkg.com/",
			want:    "https://unpkg.com/starfederation/[email protected]/bundles/datastar.js",
		},
		{
			name:    "empty version defaults to pinned",
			version: "",
			cdn:     "",
			want:    "https://cdn.jsdelivr.net/gh/starfederation/[email protected]/bundles/datastar.js",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := datastarScriptURL(tt.version, tt.cdn)
			if got != tt.want {
				t.Errorf("datastarScriptURL(%q, %q) = %q, want %q", tt.version, tt.cdn, got, tt.want)
			}
		})
	}
}

func TestDatastarCDNOrigin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		cdn  string
		want string
	}{
		{"", "https://cdn.jsdelivr.net"},
		{"https://unpkg.com", "https://unpkg.com"},
		{"https://unpkg.com/path", "https://unpkg.com"},
		{"/assets", ""},
	}
	for _, tt := range tests {
		got := datastarCDNOrigin(tt.cdn)
		if got != tt.want {
			t.Errorf("datastarCDNOrigin(%q) = %q, want %q", tt.cdn, got, tt.want)
		}
	}
}
