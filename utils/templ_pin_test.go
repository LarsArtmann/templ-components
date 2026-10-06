package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// TestTemplVersionPin pins the github.com/a-h/templ version required by every
// module's go.mod to the templGeneratedWith constant in cmd/tc/doctor.go.
//
// The 2026-10-06 incident: an auto-commit daemon re-applied the templ
// v0.3.1070 bump across 8 go.mod files while the sanctioned generator
// (nixpkgs templ in the dev shell) was still v0.3.1020. That exact bump was
// deliberately rolled back on 2026-10-05 (see AGENTS.md "sweep note"): the
// 0.3.1070 generator changes OUTPUT (icon+label spacing), so the bump needs a
// flake-level generator pin, a golden + pixel re-baseline, and a source-trim
// review first (TODO #335). Until that plan lands, generator and go.mod MUST
// stay in lockstep. This guard fails on any partial bump — go.mod moving
// without doctor.go's constant, or the constant moving without every module.
func TestTemplVersionPin(t *testing.T) {
	t.Parallel()

	const doctorFile = "../cmd/tc/doctor.go"

	doctorSrc, err := os.ReadFile(filepath.Clean(doctorFile))
	if err != nil {
		t.Fatalf("read %s: %v", doctorFile, err)
	}

	doctorRe := regexp.MustCompile(`templGeneratedWith = "(v[^"]+)"`)

	m := doctorRe.FindSubmatch(doctorSrc)
	if m == nil {
		t.Fatalf("templGeneratedWith constant not found in %s — constant renamed?", doctorFile)
	}

	expected := string(m[1])

	requireRe := regexp.MustCompile(`github\.com/a-h/templ\s+(v\S+)`)

	for _, mod := range []string{
		".", "utils", "icons", "errorpage", "charts/echarts", "htmx",
		"datastar", "website", "visualtest",
	} {
		t.Run(mod, func(t *testing.T) {
			t.Parallel()

			goMod := filepath.Join("..", mod, "go.mod")

			src, err := os.ReadFile(filepath.Clean(goMod))
			if err != nil {
				t.Fatalf("read %s: %v", goMod, err)
			}

			got := requireRe.FindSubmatch(src)
			if got == nil {
				t.Fatalf("%s has no github.com/a-h/templ require line", goMod)
			}

			if string(got[1]) != expected {
				t.Errorf(
					"%s requires templ %s but the generator pin is %s — "+
						"revert the go.mod bump or execute the TODO #335 migration plan "+
						"(flake generator pin + golden re-baseline) in the same commit",
					goMod, string(got[1]), expected,
				)
			}
		})
	}
}
