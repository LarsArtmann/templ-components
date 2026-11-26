package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestStarterEntryCompilesWithTailwind compiles the starter `app.css` exactly
// as `tc init` leaves it. A CSS comment in the starter once contained a glob
// with the comment terminator in it, which broke every project scaffolded
// from it; only the real compiler sees that. Skips without tailwindcss.
func TestStarterEntryCompilesWithTailwind(t *testing.T) {
	t.Parallel()

	tailwind, err := exec.LookPath("tailwindcss")
	if err != nil {
		t.Skipf("tailwindcss not on PATH: %v", err)
	}

	starter, err := os.ReadFile(filepath.Join("..", "..", "templates", starterEntryName))
	if err != nil {
		t.Fatalf("read the starter entry: %v", err)
	}

	root := t.TempDir()
	entry := filepath.Join(root, starterEntryName)
	writeCSSFile(t, entry, string(starter))

	opts := defaultCSSOptions()
	opts.dir = root
	opts.scaffold = true

	if err := cssSync(root, opts, &strings.Builder{}); err != nil {
		t.Fatalf("cssSync beside the starter: %v", err)
	}

	output := filepath.Join(t.TempDir(), "styles.css")

	buildErr := exec.CommandContext(t.Context(), tailwind, "--input", entry, "--output", output).Run()
	if buildErr != nil {
		t.Fatalf("the starter entry does not compile; `tc init` would scaffold a broken project: %v", buildErr)
	}

	if compiled := readCSSFile(t, output); !strings.Contains(compiled, ".tc-modal") {
		t.Errorf("the compiled starter lacks the library's component CSS")
	}
}
