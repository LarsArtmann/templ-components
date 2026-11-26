package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCSSVendoredLayoutCompilesWithTailwind is the compile proof for the
// vendored layout, the one that needs no tooling: `go mod vendor` ships the
// library sources and its stylesheet, and two native Tailwind lines reach
// both. The entry scans three vendored packages and imports the stylesheet
// from the vendor tree; the compiled output must carry a utility from a Go
// lookup map, one from a .templ source, the component CSS and the safelist.
// Skips when tailwindcss is not on PATH.
func TestCSSVendoredLayoutCompilesWithTailwind(t *testing.T) {
	t.Parallel()

	tailwind, err := exec.LookPath("tailwindcss")
	if err != nil {
		t.Skipf("tailwindcss not on PATH: %v", err)
	}

	root := vendoredConsumer(t, libraryVersion())
	entry := filepath.Join(root, "css", "app.css")
	writeCSSFile(t, entry, "@import \"tailwindcss\" source(none);\n"+
		"@source \"../vendor/"+libraryModulePath+"/{display,feedback,layout}\";\n"+
		"@import \"../vendor/"+libraryModulePath+"/"+vendoredStylesheetPath+"\";\n")

	output := filepath.Join(t.TempDir(), "styles.css")

	buildErr := exec.CommandContext(t.Context(), tailwind, "--input", entry, "--output", output).Run()
	if buildErr != nil {
		t.Fatalf("tailwindcss: %v", buildErr)
	}

	compiled := readCSSFile(t, output)

	for _, proof := range []struct {
		name     string
		selector string
	}{
		{name: "utility from a Go lookup map", selector: ".border-amber-200"},
		{name: "utility from a .templ source", selector: `.focus\:ring-inset`},
		{name: "component CSS", selector: ".tc-modal"},
		{name: "safelisted Grid class", selector: "minmax(160px,1fr)"},
	} {
		if !strings.Contains(compiled, proof.selector) {
			t.Errorf(
				"%s (%s) is missing from the compiled stylesheet of the vendored layout",
				proof.selector,
				proof.name,
			)
		}
	}
}
