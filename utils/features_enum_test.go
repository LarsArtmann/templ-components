package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// TestFeaturesEnumTableExhaustive keeps FEATURES.md's enum tables honest:
// every exported type that ships an IsValid method must be named somewhere in
// FEATURES.md. The tables document values per package; a new enum that skips
// the table (the StatTone/ScrollbackTone class of drift) fails here instead of
// silently staying undocumented.
func TestFeaturesEnumTableExhaustive(t *testing.T) {
	t.Parallel()

	root := ".."

	isValidDecl := regexp.MustCompile(`func\s+([A-Z][A-Za-z0-9]*)IsValid\s*\(`)

	names := map[string]bool{}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			base := filepath.Base(path)
			if base == ".git" || base == "website" {
				return filepath.SkipDir
			}

			return nil
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		for _, m := range isValidDecl.FindAllSubmatch(data, -1) {
			names[string(m[1])] = true
		}

		return nil
	})
	if err != nil {
		t.Fatalf("walk packages: %v", err)
	}

	features := readDoc(t, "FEATURES.md")

	var missing []string

	for name := range names {
		if !strings.Contains(features, "`"+name+"`") {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		sort.Strings(missing)
		t.Fatalf(
			"FEATURES.md does not name %d enum types that ship IsValid: %v — add them to the package's enum table (or, for inline-documented packages, name the type)",
			len(missing),
			missing,
		)
	}
}
