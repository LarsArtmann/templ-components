package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// vendoredConsumer is a project whose Tailwind entry reads the library from a
// vendor tree, the way `go mod vendor` ships it. The vendored library is a
// symlink to this repository, so the test reads the real sources.
func vendoredConsumer(t *testing.T, vendoredVersion string) string {
	t.Helper()

	root := newCSSConsumer(t, libraryVersion())
	library := filepath.Join(root, "vendor", libraryModulePath)

	if err := os.MkdirAll(filepath.Dir(library), 0o750); err != nil {
		t.Fatalf("mkdir vendor: %v", err)
	}

	repository, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolve repository: %v", err)
	}

	if err := os.Symlink(repository, library); err != nil {
		t.Fatalf("symlink the vendored library: %v", err)
	}

	writeCSSFile(t, filepath.Join(root, "vendor", "modules.txt"),
		"# "+libraryModulePath+" "+vendoredVersion+"\n\n## explicit; go 1.27\n"+libraryModulePath+"/layout\n")
	writeCSSFile(t, filepath.Join(root, "css", "app.css"), "@import \"tailwindcss\" source(none);\n"+
		"@source \"../vendor/"+libraryModulePath+"\";\n"+
		"@import \"../vendor/"+libraryModulePath+"/"+vendoredStylesheetPath+"\";\n")

	return root
}

func TestCSSCheckPassesOnVendoredLayout(t *testing.T) {
	t.Parallel()

	root := vendoredConsumer(t, libraryVersion())

	var out bytes.Buffer

	if failures := cssCheck(root, defaultCSSOptions(), &out); failures != 0 {
		t.Fatalf("cssCheck reported %d failure(s) on a vendored layout:\n%s", failures, out.String())
	}

	for _, want := range []string{
		"a Tailwind entry scans the library sources", "vendor tree matches the go.mod pin",
		"a Tailwind entry imports the library stylesheet",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("check output lacks %q:\n%s", want, out.String())
		}
	}

	if strings.Contains(out.String(), "templ-components-classes.txt") {
		t.Errorf("a vendored layout reads no synced file, so none should be checked:\n%s", out.String())
	}
}

func TestCSSCheckFailsOnStaleVendorTree(t *testing.T) {
	t.Parallel()

	root := vendoredConsumer(t, "v0.0.1")

	var out bytes.Buffer

	if failures := cssCheck(root, defaultCSSOptions(), &out); failures != 1 {
		t.Errorf("stale vendor tree should fail exactly once, got %d:\n%s", failures, out.String())
	}

	if !strings.Contains(out.String(), "vendor tree matches the go.mod pin (v0.0.1 vs "+libraryVersion()+")") ||
		!strings.Contains(out.String(), "run `go mod vendor`") {
		t.Errorf("the stale vendor tree should be named with its fix:\n%s", out.String())
	}
}

func TestCSSCheckNamesVendoredLinesWhenTheImportIsMissing(t *testing.T) {
	t.Parallel()

	root := vendoredConsumer(t, libraryVersion())
	entry := filepath.Join(root, "css", "app.css")

	writeCSSFile(t, entry, strings.Replace(readCSSFile(t, entry),
		"@import \"../vendor/"+libraryModulePath+"/"+vendoredStylesheetPath+"\";", "", 1))

	var out bytes.Buffer

	if failures := cssCheck(root, defaultCSSOptions(), &out); failures != 1 {
		t.Errorf("a vendored layout without the stylesheet import should fail exactly once, got %d:\n%s", failures,
			out.String())
	}

	want := "@import \"../vendor/" + libraryModulePath + "/" + vendoredStylesheetPath + "\"; (when `go mod vendor` runs before Tailwind)"
	if !strings.Contains(out.String(), want) {
		t.Errorf("the fix should name the vendored import relative to the entry, want %s:\n%s", want, out.String())
	}
}

func TestVendoredLibraryRelativeFindsTheVendorTree(t *testing.T) {
	t.Parallel()

	root := vendoredConsumer(t, libraryVersion())

	for _, test := range []struct {
		entryDirectory string
		want           string
	}{
		{entryDirectory: filepath.Join(root, "css"), want: "../vendor/" + libraryModulePath},
		{entryDirectory: filepath.Join(root, "a", "b", "c"), want: "../../../vendor/" + libraryModulePath},
	} {
		if got := vendoredLibraryRelative(test.entryDirectory); got != test.want {
			t.Errorf("vendoredLibraryRelative(%s) = %q, want %q", test.entryDirectory, got, test.want)
		}
	}

	if got := vendoredLibraryRelative(
		filepath.Join(newCSSConsumer(t, libraryVersion()), "css"),
	); got != vendoredLibraryPath {
		t.Errorf("without a vendor tree the conventional path applies, got %q", got)
	}
}

// TestCSSCheckNamesAHandCuratedVendorList covers the layout that ships the
// next adopted component unstyled: an entry that lists single files of the
// vendor tree. The check names the directory line that replaces the list.
func TestCSSCheckNamesAHandCuratedVendorList(t *testing.T) {
	t.Parallel()

	root := vendoredConsumer(t, libraryVersion())
	entry := filepath.Join(root, "css", "app.css")

	writeCSSFile(t, entry, "@import \"tailwindcss\" source(none);\n"+
		"@source \"../vendor/"+libraryModulePath+"/layout/base.templ\";\n"+
		"@source \"../vendor/"+libraryModulePath+"/layout/theme.templ\";\n"+
		"@import \"../vendor/"+libraryModulePath+"/"+vendoredStylesheetPath+"\";\n")

	var out bytes.Buffer

	if failures := cssCheck(root, defaultCSSOptions(), &out); failures != 1 {
		t.Errorf("a curated vendor list should fail exactly once, got %d:\n%s", failures, out.String())
	}

	if !strings.Contains(out.String(), "2 single library files are listed by hand") ||
		!strings.Contains(out.String(), vendoredSourceLine(filepath.Join(root, "css"))) {
		t.Errorf("the fix should replace the list with the directory line:\n%s", out.String())
	}
}
