package main

import (
	"path/filepath"
	"strings"
	"testing"
)

// The doctor's two false-FAIL fixes (#35): the templ pin check must read
// both go.mod require forms (the line-prefix matcher read only the
// single-line form, so every block-form consumer — the gofmt default —
// read as "templ missing"), and the Tailwind entry check must discover the
// entry by content (the fixed filename probe missed layouts like
// assets/css/input.css).

func TestTemplRequireVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    string
		wantOK  bool
	}{
		{
			name: "block require form resolves",
			content: `module example.com/consumer

go 1.27

require (
	github.com/a-h/templ v0.3.1100
	github.com/larsartmann/templ-components v1.27.0
)

require (
	github.com/Oudwins/tailwind-merge-go v0.2.3 // indirect
)
`,
			want:   "0.3.1100",
			wantOK: true,
		},
		{
			name: "single-line require form resolves",
			content: `module example.com/consumer

go 1.27

require github.com/a-h/templ v0.3.1020
`,
			want:   "0.3.1020",
			wantOK: true,
		},
		{
			name:    "single-line indirect resolves",
			content: "require github.com/a-h/templ v0.3.1100 // indirect\n",
			want:    "0.3.1100",
			wantOK:  true,
		},
		{
			name: "trailing comment stripped",
			content: `require (
	github.com/a-h/templ v0.3.1100 // pinned to the generator
)
`,
			want:   "0.3.1100",
			wantOK: true,
		},
		{
			name: "replace block does not masquerade",
			content: `require github.com/a-h/templ v0.3.1100

replace github.com/a-h/templ => github.com/fork/templ v9.9.9
`,
			want:   "0.3.1100",
			wantOK: true,
		},
		{
			name: "closed and reopened require block resolves",
			content: `require (
	github.com/a v1.0.0
)

exclude github.com/a-h/templ v0.3.900

require (
	github.com/a-h/templ v0.3.1100
)
`,
			want:   "0.3.1100",
			wantOK: true,
		},
		{
			name:    "missing dependency reports not found",
			content: "module example.com/consumer\n\ngo 1.27\n\nrequire github.com/some/other v1.0.0\n",
			want:    "",
			wantOK:  false,
		},
		{
			name:    "empty go.mod reports not found",
			content: "",
			want:    "",
			wantOK:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := templRequireVersion(tt.content)
			if ok != tt.wantOK || got != tt.want {
				t.Errorf("templRequireVersion() = (%q, %v), want (%q, %v)", got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

func TestCheckTemplVersionPinnedAtVerdicts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		goMod    string
		wantPASS bool
	}{
		{
			name:     "block-form current pin passes",
			goMod:    "module c\n\nrequire (\n\tgithub.com/a-h/templ v0.3.1100\n)\n",
			wantPASS: true,
		},
		{
			name:     "block-form old pin fails",
			goMod:    "module c\n\nrequire (\n\tgithub.com/a-h/templ v0.3.1020\n)\n",
			wantPASS: false,
		},
		{name: "missing templ fails", goMod: "module c\n\ngo 1.27\n", wantPASS: false},
		{name: "unreadable path fails", goMod: "", wantPASS: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := filepath.Join(t.TempDir(), "go.mod")

			if tt.goMod != "" {
				writeFile(t, path, tt.goMod)
			}

			if (checkTemplVersionPinnedAt(path) == 0) != tt.wantPASS {
				t.Fatalf("checkTemplVersionPinnedAt(%q) verdict != pass=%v", path, tt.wantPASS)
			}
		})
	}
}

func TestDiscoverCSSEntry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		cssPath    string
		cssContent string
		want       string
	}{
		{
			name:       "nested entry discovered",
			cssPath:    "assets/css/input.css",
			cssContent: "@import \"tailwindcss\";\n@source \"../../../**/*.templ\";\n",
			want:       "assets/css/input.css",
		},
		{
			name:       "templ-less entry returned as hint",
			cssPath:    "assets/css/input.css",
			cssContent: "@import \"tailwindcss\";\n@source \"../src\";\n",
			want:       "assets/css/input.css",
		},
		{
			name:    "no css at all",
			cssPath: "",
			want:    "",
		},
		{
			name:       "css without @source ignored",
			cssPath:    "styles/theme.css",
			cssContent: "body { color: red; }\n",
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()

			if tt.cssPath != "" {
				writeFile(t, filepath.Join(dir, tt.cssPath), tt.cssContent)
			}

			got, found := discoverCSSEntry(dir)
			if tt.want == "" {
				if found {
					t.Fatalf("discoverCSSEntry() = (%q, true), want not found", got)
				}

				return
			}

			if !found || got != filepath.Join(dir, tt.want) {
				t.Errorf("discoverCSSEntry() = (%q, %v), want (%q, true)", got, found, filepath.Join(dir, tt.want))
			}
		})
	}
}

func TestDiscoverCSSEntryPrefersScanningEntry(t *testing.T) {
	t.Parallel()

	// When several @source-carrying css files exist, the one that actually
	// scans .templ wins over an earlier templ-less hit (lexical walk order
	// must not doom a valid entry).
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "aaa/inert.css"), "@source \"./other\";\n")
	writeFile(t, filepath.Join(dir, "zzz/input.css"), "@source \"../../**/*.templ\";\n")

	got, found := discoverCSSEntry(dir)
	if !found || got != filepath.Join(dir, "zzz/input.css") {
		t.Errorf("discoverCSSEntry() = (%q, %v), want the scanning entry", got, found)
	}
}

func TestDiscoverCSSEntryDepthBound(t *testing.T) {
	t.Parallel()

	// The walk must not reach arbitrarily deep trees (huge consumer repos).
	dir := t.TempDir()

	deep := dir + "/level" + strings.Repeat("/level", cssDiscoveryMaxDepth+4)
	writeFile(t, deep+"/late.css", "@source \"**/*.templ\";\n")

	if got, found := discoverCSSEntry(dir); found {
		t.Errorf("walked past the depth bound (found %q)", got)
	}
}

func TestCssFastPathEntry(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "static/app.css"), "@source \"**/*.templ\";\n")

	if got := cssFastPathEntry(dir); got != filepath.Join(dir, "static/app.css") {
		t.Errorf("cssFastPathEntry() = %q, want the static/app.css candidate", got)
	}

	if got := cssFastPathEntry(t.TempDir()); got != "" {
		t.Errorf("cssFastPathEntry() = %q in an empty dir, want \"\"", got)
	}
}

func TestDoctorGeneratedFileNotBroken(t *testing.T) {
	t.Parallel()

	// The doctor must hand out advice anchored to the same generator version
	// the library pins (kept in lockstep by utils.TestTemplVersionPin on the
	// library side); this asserts the doctor constant is version-shaped.
	if !isVersionLiteral(templGeneratedWith) {
		t.Errorf("templGeneratedWith = %q, not a version literal", templGeneratedWith)
	}
}
