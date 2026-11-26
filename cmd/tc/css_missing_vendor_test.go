package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

// TestCSSCheckNamesAMissingVendorTree covers the failure that ships unstyled
// pages silently: the entry names the library's vendor tree, but no tree
// exists where CSS is compiled (`go mod vendor` did not run). The check must
// say that, not ask for synced files.
func TestCSSCheckNamesAMissingVendorTree(t *testing.T) {
	t.Parallel()

	root := newCSSConsumer(t, libraryVersion())
	writeCSSFile(t, filepath.Join(root, "css", "app.css"), "@import \"tailwindcss\" source(none);\n"+
		"@source \"../vendor/"+libraryModulePath+"\";\n")

	var out bytes.Buffer

	if failures := cssCheck(root, defaultCSSOptions(), &out); failures != 2 {
		t.Errorf("failures = %d, want 2 (vendor tree missing, stylesheet not imported):\n%s", failures, out.String())
	}

	if !strings.Contains(out.String(), "a Tailwind entry scans the library's vendor tree") ||
		!strings.Contains(out.String(), "run `go mod vendor` before Tailwind") {
		t.Errorf("the missing vendor tree should be named with its fix:\n%s", out.String())
	}

	if strings.Contains(out.String(), corpusName+" (") {
		t.Errorf("an entry that names the vendor tree wants vendoring, not synced files:\n%s", out.String())
	}
}

// TestCSSCheckFromTheProjectRoot runs the check the way the command does:
// from the project root, with entry paths relative to it. The glob form of the
// vendor line is the one consumers write by hand.
func TestCSSCheckFromTheProjectRoot(t *testing.T) { // t.Chdir cannot run in parallel
	root := newCSSConsumer(t, libraryVersion())
	writeCSSFile(t, filepath.Join(root, "css", "app.css"), "@import \"tailwindcss\" source(none);\n"+
		"@source \"../vendor/"+libraryModulePath+"/**/*.templ\";\n")

	t.Chdir(root)

	var out bytes.Buffer

	if failures := cssCheck(".", defaultCSSOptions(), &out); failures != 2 {
		t.Errorf("failures = %d, want 2 (vendor tree missing, stylesheet not imported):\n%s", failures, out.String())
	}

	if !strings.Contains(out.String(), "a Tailwind entry scans the library's vendor tree") ||
		strings.Contains(out.String(), corpusName+" (") {
		t.Errorf("a relative entry path must still be recognized as the vendor tree:\n%s", out.String())
	}

	if isVendoredLibraryPath("/src/vendor/" + libraryModulePath + "-other/x.templ") {
		t.Errorf("a sibling module name that merely starts with the library path is not the library")
	}
}
