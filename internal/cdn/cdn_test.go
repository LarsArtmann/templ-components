package cdn

import "testing"

func TestResolveBase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		cdn         string
		defaultBase string
		want        string
	}{
		{"empty uses default", "", "https://cdn.jsdelivr.net/npm", "https://cdn.jsdelivr.net/npm"},
		{"trims trailing slash", "https://unpkg.com/", "https://cdn.jsdelivr.net/npm", "https://unpkg.com"},
		{"passes through unchanged", "https://unpkg.com", "https://cdn.jsdelivr.net/npm", "https://unpkg.com"},
		{"keeps path", "https://unpkg.com/path", "https://cdn.jsdelivr.net/npm", "https://unpkg.com/path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ResolveBase(tt.cdn, tt.defaultBase)
			if got != tt.want {
				t.Errorf("ResolveBase(%q, %q) = %q, want %q", tt.cdn, tt.defaultBase, got, tt.want)
			}
		})
	}
}

func TestOrigin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		cdn         string
		defaultBase string
		want        string
	}{
		{"empty uses default", "", "https://cdn.jsdelivr.net/gh", "https://cdn.jsdelivr.net"},
		{"unpkg clean", "https://unpkg.com", "https://cdn.jsdelivr.net/gh", "https://unpkg.com"},
		{"unpkg with path", "https://unpkg.com/path", "https://cdn.jsdelivr.net/gh", "https://unpkg.com"},
		{"relative path returns empty", "/assets", "https://cdn.jsdelivr.net/gh", ""},
		{"http scheme", "http://example.com/foo", "https://cdn.jsdelivr.net/gh", "http://example.com"},
		{"trims slash then origin", "https://unpkg.com/", "https://cdn.jsdelivr.net/gh", "https://unpkg.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Origin(tt.cdn, tt.defaultBase)
			if got != tt.want {
				t.Errorf("Origin(%q, %q) = %q, want %q", tt.cdn, tt.defaultBase, got, tt.want)
			}
		})
	}
}
