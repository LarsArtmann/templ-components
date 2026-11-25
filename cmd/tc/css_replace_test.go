package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

// TestCSSSyncAcceptsLocalReplace covers the development layout: a `replace`
// directive points the library at a checkout, so go.mod carries no version for
// it. The tc built from that tree ships that tree's files, so the version
// comparison must not refuse.
func TestCSSSyncAcceptsLocalReplace(t *testing.T) {
	t.Parallel()

	root := newCSSConsumer(t, "v0.0.0-00010101000000-000000000000")
	goModPath := filepath.Join(root, "go.mod")

	var out bytes.Buffer

	err := cssSync(root, defaultCSSOptions(), &out)
	if err == nil || !strings.Contains(err.Error(), "refuses") && !strings.Contains(err.Error(), "go.mod pins") {
		t.Fatalf("the zero pseudo-version without a replace should be refused, got: %v", err)
	}

	writeCSSFile(t, goModPath, readCSSFile(t, goModPath)+"\nreplace "+libraryModulePath+" => ../templ-components\n")

	out.Reset()

	if err := cssSync(root, defaultCSSOptions(), &out); err != nil {
		t.Fatalf("cssSync with a local replace: %v", err)
	}

	if !strings.Contains(out.String(), "version check is skipped") {
		t.Errorf("sync should name the skipped version check:\n%s", out.String())
	}

	out.Reset()

	if failures := cssCheck(root, defaultCSSOptions(), &out); failures != 2 {
		t.Errorf("check with a local replace = %d failure(s), want 2 (only the unwired entry):\n%s", failures,
			out.String())
	}

	if !strings.Contains(out.String(), "replaces templ-components with a local directory") {
		t.Errorf("check should report the local replace:\n%s", out.String())
	}

	writeCSSFile(
		t,
		goModPath,
		strings.Replace(readCSSFile(t, goModPath), " => ../templ-components", " => "+libraryModulePath+" v0.0.1", 1),
	)

	if libraryReplacedLocally(root) {
		t.Errorf("a replace with a module version is not a local replace")
	}
}
