package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// TestAddSmoke exercises the real `tc add` binary end-to-end in a throwaway
// directory (TODO #283c, 2026-09-23): the unit tests pin the registry, but
// only an exec'd binary proves the os.Args dispatch, the source-embed read,
// and the file writes actually work together. eyebrow is a single-file
// display component; auth_layout is a recipe shipping a *_types.go companion
// — the two shapes `tc add` must handle.
func TestAddSmoke(t *testing.T) {
	t.Parallel()

	bin := filepath.Join(t.TempDir(), "tc-build")

	buildCtx, cancelBuild := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancelBuild)

	build := exec.CommandContext(buildCtx, "go", "build", "-o", bin, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build tc: %v\n%s", err, out)
	}

	cases := []struct {
		name      string
		wantFiles []string
	}{
		{name: "eyebrow", wantFiles: []string{"eyebrow.templ"}},
		{name: "auth_layout", wantFiles: []string{"auth_layout.templ", "auth_layout_types.go"}},
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			add := exec.CommandContext(ctx, bin, "add", tcase.name)
			add.Dir = dir

			out, err := add.CombinedOutput()
			if err != nil {
				t.Fatalf("tc add %s: %v\n%s", tcase.name, err, out)
			}

			for _, want := range tcase.wantFiles {
				got := filepath.Join(dir, "components", want)

				data, err := os.ReadFile(got)
				if err != nil {
					t.Errorf("tc add %s did not produce %s: %v", tcase.name, got, err)

					continue
				}

				if len(data) == 0 {
					t.Errorf("scaffolded %s is empty", got)
				}
			}
		})
	}
}
