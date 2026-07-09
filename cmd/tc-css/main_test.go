package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindProjectRoot(t *testing.T) {
	root := findProjectRoot()
	if root == "" {
		t.Fatal("findProjectRoot returned empty string")
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatalf("findProjectRoot returned %q but go.mod not found there: %v", root, err)
	}
}

func TestFileExists(t *testing.T) {
	if !fileExists("main.go") {
		t.Error("fileExists(main.go) = false, want true")
	}
	if fileExists("nonexistent.go") {
		t.Error("fileExists(nonexistent.go) = true, want false")
	}
	if fileExists(".") {
		t.Error("fileExists(.) = true, want false (directory)")
	}
}

func TestDirExists(t *testing.T) {
	if !dirExists(".") {
		t.Error("dirExists(.) = false, want true")
	}
	if dirExists("main.go") {
		t.Error("dirExists(main.go) = true, want false")
	}
}

func TestAbsPath(t *testing.T) {
	got := absPath("/root", "relative/path.css")
	want := "/root/relative/path.css"
	if got != want {
		t.Errorf("absPath(/root, relative/path.css) = %q, want %q", got, want)
	}

	got = absPath("/root", "/abs/path.css")
	want = "/abs/path.css"
	if got != want {
		t.Errorf("absPath(/root, /abs/path.css) = %q, want %q", got, want)
	}
}

func TestGenerateInputCSS(t *testing.T) {
	tmp := t.TempDir()
	projectRoot := tmp
	cssDir := filepath.Join(projectRoot, "css")
	cssPath := filepath.Join(cssDir, "app.css")

	if err := generateInputCSS(cssPath, projectRoot); err != nil {
		t.Fatalf("generateInputCSS failed: %v", err)
	}

	content, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read generated CSS: %v", err)
	}

	cssStr := string(content)

	checks := []struct {
		name    string
		content string
	}{
		{"tailwind import", `@import "tailwindcss"`},
		{"source none", "source(none)"},
		{"templ scan", `**/*.templ`},
		{"dark variant", "@custom-variant dark"},
		{"theme comment", "@theme"},
	}

	for _, check := range checks {
		if !strings.Contains(cssStr, check.content) {
			t.Errorf("generated CSS missing %s: expected to contain %q", check.name, check.content)
		}
	}
}

func TestGenerateInputCSSWithVendor(t *testing.T) {
	tmp := t.TempDir()
	projectRoot := tmp
	cssDir := filepath.Join(projectRoot, "static")
	cssPath := filepath.Join(cssDir, "input.css")

	// Create a fake vendored .templ file
	vendorTempl := filepath.Join(projectRoot, "vendor", "github.com", "larsartmann", "templ-components", "display")
	if err := os.MkdirAll(vendorTempl, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vendorTempl, "card.templ"), []byte("package display"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := generateInputCSS(cssPath, projectRoot); err != nil {
		t.Fatalf("generateInputCSS failed: %v", err)
	}

	content, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read generated CSS: %v", err)
	}

	cssStr := string(content)

	// Should contain @source for the vendored templ-components
	if !strings.Contains(cssStr, "vendor") {
		t.Error("generated CSS should contain @source pointing to vendored templ-components")
	}
	if !strings.Contains(cssStr, "templ-components") {
		t.Error("generated CSS should reference templ-components in @source path")
	}
}

func TestFindTemplVendoredSources(t *testing.T) {
	tmp := t.TempDir()

	// No vendor dir → empty
	sources := findTemplVendoredSources(tmp, tmp)
	if len(sources) != 0 {
		t.Errorf("findTemplVendoredSources with no vendor/ should return empty, got %v", sources)
	}

	// Create vendor structure
	vendorTempl := filepath.Join(tmp, "vendor", "github.com", "example", "lib")
	if err := os.MkdirAll(vendorTempl, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vendorTempl, "comp.templ"), []byte("test"), 0o600); err != nil {
		t.Fatal(err)
	}

	cssDir := tmp
	sources = findTemplVendoredSources(cssDir, tmp)
	if len(sources) == 0 {
		t.Fatal("findTemplVendoredSources should find vendored .templ files")
	}

	found := false
	for _, src := range sources {
		if strings.Contains(src, "vendor") && strings.Contains(src, "**/*.templ") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("findTemplVendoredSources should include vendored .templ path, got %v", sources)
	}
}
