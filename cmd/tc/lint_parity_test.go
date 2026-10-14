package main

import (
	"bytes"
	"path/filepath"
	"testing"
)

// TestLintAllParityWithGuardTests pins the parity contract between `tc lint`
// and the library's guard tests (master-plan M2/F7): each scanner must return
// CLEAN on exactly the directory set its guard test sweeps, because the
// scanners ARE the guard engines (same exported funcs) and the guards are
// green on this repo. A seeded-violation proof per rule lives in the utils
// scanner tests; together the two directions prove verdict equality on this
// tree: guard-green ⇒ lint-clean, and lint fires where the guard fires.
func TestLintAllParityWithGuardTests(t *testing.T) {
	t.Parallel()

	// Repo-root-relative, mirroring the guard tests' dir lists (utils tests
	// run with CWD=utils and prefix each dir with "..").
	guardDirs := map[string][]string{
		"--rtl":    {"display", "feedback", "forms", "navigation", "errorpage", "layout", "htmx"},
		"--motion": {"display", "feedback", "forms", "navigation", "errorpage", "layout", "htmx"},
		"--dark":   {"display", "feedback", "forms", "navigation", "errorpage", "layout", "htmx", "examples/demo"},
		"--icons": {
			"display", "forms", "feedback", "layout", "navigation",
			"errorpage", "recipes", "htmx", "datastar", "charts/echarts",
			"icons", "utils",
		},
	}

	for flag, dirs := range guardDirs {
		t.Run(flag, func(t *testing.T) {
			t.Parallel()

			args := []string{flag}
			for _, dir := range dirs {
				args = append(args, filepath.Join("..", "..", dir))
			}

			stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}

			if code := runLint(args, stdout, stderr); code != 0 {
				t.Errorf(
					"tc lint %s over the guard dir set is NOT clean (exit %d) — the CLI and the guard tests disagree:\n%s%s",
					flag,
					code,
					stderr.String(),
					stdout.String(),
				)
			}
		})
	}
}

// TestParseLintArgsAll proves --all selects every scanner exactly once and
// that duplicate/repeated flags do not double-run a scanner.
func TestParseLintArgsAll(t *testing.T) {
	t.Parallel()

	selected, dirs, ok := parseLintArgs([]string{"--all", "--rtl", "."}, &bytes.Buffer{})
	if !ok {
		t.Fatal("parseLintArgs(--all --rtl .) reported !ok")
	}

	if len(selected) != len(lintScanners) {
		t.Errorf("--all selected %d scanners, want %d", len(selected), len(lintScanners))
	}

	seen := map[string]bool{}
	for _, s := range selected {
		if seen[s.flag] {
			t.Errorf("scanner %s selected twice", s.flag)
		}

		seen[s.flag] = true
	}

	if len(dirs) != 1 || dirs[0] != "." {
		t.Errorf("dirs = %v, want [.] with --all consuming no positional", dirs)
	}
}
