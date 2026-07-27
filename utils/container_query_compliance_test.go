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
// breakpoints intentionally without ContainerAware.
var containerQueryExceptions = []containerQueryException{
	// AppShell: full-page layout — always viewport-relative.
	{"layout/appshell.templ", "full-page shell always viewport-relative"},
	{"navigation/sidebar_nav.templ", "sidebar is viewport-relative by design"},
	{"navigation/mobile_menu.templ", "mobile menu is viewport-gated by design"},
	{"navigation/nav.templ", "responsive nav bar is viewport-relative"},
	// Error pages — full-page layouts, viewport-relative.
	{"errorpage/notfound404.templ", "full-page 404 is viewport-relative"},
	// Demo/recipe pages — viewport-relative by design.
	{"recipes/dashboard.templ", "dashboard recipe is viewport-relative"},
	{"recipes/settings_layout.templ", "settings layout is viewport-relative"},
}

func isContainerQueryException(path string) bool {
	for _, ex := range containerQueryExceptions {
		if strings.Contains(path, ex.file) {
			return true
		}
	}

	return false
}
