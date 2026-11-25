package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCSSSyncCompilesWithTailwind is the compile proof behind `tc css sync`:
// string tests cannot show that Tailwind turns the synced files into CSS. A
// consumer that syncs and wires its entry gets the library's utility classes
// (from Go lookup maps and from .templ sources), its component CSS and the
// safelist in the compiled stylesheet; the same consumer without the library
// files gets none of them. Skips when tailwindcss is not on PATH.
func TestCSSSyncCompilesWithTailwind(t *testing.T) {
	t.Parallel()

	tailwind, err := exec.LookPath("tailwindcss")
	if err != nil {
		t.Skipf("tailwindcss not on PATH: %v", err)
	}

	root := newCSSConsumer(t, libraryVersion())
	entry := filepath.Join(root, "css", "app.css")
	writeCSSFile(t, filepath.Join(root, "internal", "views", "page.templ"), `<div class="text-center">own</div>`)

	compile := func() string {
		output := filepath.Join(t.TempDir(), "styles.css")

		buildErr := exec.CommandContext(t.Context(), tailwind, "--input", entry, "--output", output).Run()
		if buildErr != nil {
			t.Fatalf("tailwindcss: %v", buildErr)
		}

		return readCSSFile(t, output)
	}

	proofs := []struct {
		name     string
		selector string
	}{
		{name: "utility from a Go lookup map", selector: ".border-amber-200"},
		{name: "utility from a .templ source", selector: `.focus\:ring-inset`},
		{name: "component CSS", selector: ".tc-modal"},
		{name: "safelisted Grid class", selector: "minmax(160px,1fr)"},
	}

	before := compile()
	if !strings.Contains(before, ".text-center") {
		t.Fatalf("the consumer's own class is missing before the sync, so the compile proof is blind:\n%s", before)
	}

	for _, proof := range proofs {
		if strings.Contains(before, proof.selector) {
			t.Errorf("%s (%s) is present before the sync; the proof cannot show the sync added it", proof.selector,
				proof.name)
		}
	}

	if err := cssSync(root, defaultCSSOptions(), &strings.Builder{}); err != nil {
		t.Fatalf("cssSync: %v", err)
	}

	writeCSSFile(t, entry, strings.Replace(readCSSFile(t, entry), "@source ", wiringBlock+"@source ", 1))

	if failures := cssCheck(root, defaultCSSOptions(), &strings.Builder{}); failures != 0 {
		t.Fatalf("cssCheck reported %d failure(s) on the synced, wired consumer", failures)
	}

	after := compile()

	for _, proof := range proofs {
		if !strings.Contains(after, proof.selector) {
			t.Errorf("%s (%s) is missing from the compiled stylesheet after the sync", proof.selector, proof.name)
		}
	}
}
