package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestContainerQueryCompliance scans .templ files for structural viewport
// breakpoints (sm:/md:/lg:/xl: on grid-cols, flex, hidden, col-span, etc.)
// that lack a corresponding ContainerAware or ContainerResponsive flag.
//
// Container queries let a component adapt to its parent container width
// instead of the browser viewport — critical for embedded/nested layouts.
// If a component uses responsive grid/flex/display classes, it SHOULD offer
// a ContainerAware opt-in so consumers can switch to container-based
// responsiveness.
//
// Exceptions are listed below with explicit reasons.
func TestContainerQueryCompliance(t *testing.T) {
	t.Parallel()

	dirs := []string{
		"display", "feedback", "forms", "navigation",
		"errorpage", "layout", "htmx", "recipes",
	}

	structuralRe := regexp.MustCompile(
		`(sm|md|lg|xl):(grid-cols|grid-rows|flex-col|flex-row|hidden|block|col-span|col-start|col-end|flex-1|flex-auto|flex-none|basis-)`,
	)

	containerAwareRe := regexp.MustCompile(
		`ContainerAware|ContainerResponsive`,
	)

	scanContainerQuery(t, dirs, structuralRe, containerAwareRe)
}

// scanContainerQuery walks .templ files in the given dirs and reports
// structural viewport breakpoints without ContainerAware.
func scanContainerQuery(
	t *testing.T,
	dirs []string,
	structuralRe *regexp.Regexp,
	containerAwareRe *regexp.Regexp,
) {
	t.Helper()

	violations := 0

	for _, dir := range dirs {
		dirPath := filepath.Join("..", dir)

		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			return checkContainerQueryFile(
				t, path, info, err, structuralRe, containerAwareRe, &violations,
			)
		})
		if err != nil {
			t.Logf("walk error for %s: %v", dir, err)
		}
	}

	if violations > 0 {
		t.Errorf("found %d container-query compliance violations", violations)
	}
}

// checkContainerQueryFile is the filepath.Walk callback for a single .templ file.
func checkContainerQueryFile(
	t *testing.T,
	path string,
	info os.FileInfo,
	walkErr error,
	structuralRe *regexp.Regexp,
	containerAwareRe *regexp.Regexp,
	violations *int,
) error {
	t.Helper()

	if walkErr != nil || info.IsDir() {
		return walkErr
	}

	if !strings.HasSuffix(path, ".templ") {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	src := string(data)

	if !structuralRe.MatchString(src) {
		return nil
	}

	if containerAwareRe.MatchString(src) {
		return nil
	}

	relPath := strings.TrimPrefix(path, "../")

	if isContainerQueryException(relPath) {
		return nil
	}

	commentRe := regexp.MustCompile(`^\s*//`)

	for line := range strings.SplitSeq(src, "\n") {
		if commentRe.MatchString(line) {
			continue
		}

		if structuralRe.MatchString(line) {
			*violations++

			t.Errorf(
				"%s uses structural viewport breakpoints without ContainerAware:\n  %s",
				relPath,
				strings.TrimSpace(line),
			)
		}
	}

	return nil
}

// containerQueryException defines a file that legitimately uses viewport
// breakpoints without ContainerAware.
type containerQueryException struct {
	file   string
	reason string
}

// containerQueryExceptions lists files that use structural viewport
// breakpoints intentionally without ContainerAware. Each entry was verified
// (2026-07-28) to (a) actually contain a structural breakpoint the scanner
// would flag, and (b) be a full-page component where container queries do not
// apply — the component fills the viewport, so "container width" == "viewport
// width" and a ContainerAware flag would be meaningless.
//
// Do NOT add an entry unless the file BOTH trips the scanner AND is genuinely
// viewport-relative. Pre-emptive exemptions for files that have no structural
// breakpoint weaken the test (they hide future regressions); such entries were
// pruned in this pass.
var containerQueryExceptions = []containerQueryException{
	{"layout/appshell.templ", "full-page admin shell: sidebar visibility is viewport-gated (lg:block/lg:hidden)"},
	{"navigation/mobile_menu.templ", "mobile menu: shown/hidden by viewport width (sm:hidden)"},
	{"navigation/nav.templ", "responsive nav bar: grid layout switches at viewport md: breakpoint"},
	{"errorpage/notfound404.templ", "full-page 404: link grid columns switch at viewport sm:/lg: breakpoints"},
	{
		"recipes/auth_layout.templ",
		"full-page auth split-screen: branding panel visibility is viewport-gated (lg:flex/hidden)",
	},
}

func isContainerQueryException(path string) bool {
	for _, ex := range containerQueryExceptions {
		if strings.Contains(path, ex.file) {
			return true
		}
	}

	return false
}
