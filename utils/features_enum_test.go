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

	features := string(readDoc(t, "FEATURES.md"))

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

// TestFeaturesEnumValuesExhaustive keeps the VALUES in FEATURES.md's enum
// tables honest: for every documented enum type whose Go constants are
// declared as `Name Type = "literal"` on a single line, a constant that no
// documented value covers fails (the Family/Orchestration class of drift —
// the type was named, a value was silently left out). Conservative by
// design: a row is only enforced when EVERY documented value resolves to a
// constant name ending in that value; rows using a different naming
// convention are skipped rather than mis-judged, and `*Default`/`*Unspecified`
// convenience constants are exempt.
func TestFeaturesEnumValuesExhaustive(t *testing.T) {
	t.Parallel()

	root := ".."

	// Matches single-line typed const declarations inside const blocks:
	// \tFamilyRejection Family = "rejection"
	constDecl := regexp.MustCompile(`(?m)^\t([A-Z][A-Za-z0-9]*)\s+([A-Z][A-Za-z0-9]*)\s*=\s*"([^"]*)"`)

	// Matches a FEATURES.md enum table row: | `Type` | A, B, C |
	tableRow := regexp.MustCompile(`(?m)^\|\s*` + "`([A-Za-z0-9]+)`" + `\s*\|\s*([^|]+)\|`)

	consts := map[string]map[string]bool{} // type name -> const name set

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

		for _, m := range constDecl.FindAllSubmatch(data, -1) {
			typ := string(m[2])
			if consts[typ] == nil {
				consts[typ] = map[string]bool{}
			}

			consts[typ][string(m[1])] = true
		}

		return nil
	})
	if err != nil {
		t.Fatalf("walk packages: %v", err)
	}

	features := string(readDoc(t, "FEATURES.md"))

	var problems []string

	for _, m := range tableRow.FindAllStringSubmatch(features, -1) {
		typ, values := m[1], m[2]

		cset := consts[typ]
		if len(cset) == 0 {
			continue
		}

		documented := []string{}
		for _, v := range strings.Split(values, ",") {
			v = strings.TrimSpace(v)
			if v != "" {
				documented = append(documented, v)
			}
		}

		// All-or-nothing: only enforce rows where EVERY documented value
		// resolves to a constant name ending in that value. Rows using a
		// different naming convention are skipped rather than mis-judged.
		resolvable := true

		for _, v := range documented {
			matched := false
			for c := range cset {
				if strings.HasSuffix(strings.ToLower(c), strings.ToLower(v)) {
					matched = true

					break
				}
			}

			if !matched {
				resolvable = false

				break
			}
		}

		if !resolvable {
			continue
		}

		for c := range cset {
			lower := strings.ToLower(c)
			if strings.HasSuffix(lower, "default") || strings.HasSuffix(lower, "unspecified") {
				continue
			}

			listed := false
			for _, v := range documented {
				if strings.HasSuffix(lower, strings.ToLower(v)) {
					listed = true

					break
				}
			}

			if !listed {
				problems = append(problems, fmt.Sprintf(
					"enum %s has constant %s that no documented value in FEATURES.md covers",
					typ, c,
				))
			}
		}
	}

	if len(problems) > 0 {
		sort.Strings(problems)
		t.Fatalf("FEATURES.md enum tables are missing %d constants:\n%s",
			len(problems), strings.Join(problems, "\n"))
	}
}
