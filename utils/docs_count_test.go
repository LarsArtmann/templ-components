package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestDocsCountDrift(t *testing.T) {
	t.Parallel()

	root := ".."
	actualComponents := countExportedTemplFunctions(t, root)
	actualGenerated := countGeneratedFiles(t, root)
	actualIsValid := countIsValidMethods(t, root)

	features := readDoc(t, "FEATURES.md")
	assertCount(t, features, `(\d+)\s+templ components`, "FEATURES.md templ components", actualComponents)
	assertCount(t, features, `(\d+)\s+generated\s+.*_templ\.go.*files`, "FEATURES.md generated files", actualGenerated)

	agents := readDoc(t, "AGENTS.md")
	assertCount(t, agents, `(\d+)\s+generated files across all packages`, "AGENTS.md generated files", actualGenerated)

	skill := readDoc(t, "skill", "SKILL.md")
	assertCount(t, skill, `(\d+)\s+components across 9 packages`, "SKILL.md components", actualComponents)

	sections := readDoc(t, "website", "src", "data", "sections.ts")
	assertCount(t, sections, `(\d+)\s+components across 9 packages`, "website sections.ts components", actualComponents)
	assertCount(t, sections, `(\d+)\s+typed string enums`, "website sections.ts typed string enums", actualIsValid)
}

func countExportedTemplFunctions(t *testing.T, root string) int {
	t.Helper()
	templFuncRe := regexp.MustCompile(`^templ\s+([A-Z][A-Za-z0-9]*)\s*\(`)
	count := 0
	packages := []string{"display", "feedback", "forms", "navigation", "errorpage", "layout", "htmx"}
	for _, pkg := range packages {
		files, err := filepath.Glob(filepath.Join(root, pkg, "*.templ"))
		if err != nil {
			t.Fatalf("glob error for %s: %v", pkg, err)
		}
		for _, file := range files {
			//nolint:gosec // test reads known package files
			data, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("read %s: %v", file, err)
			}
			for line := range strings.SplitSeq(string(data), "\n") {
				if templFuncRe.MatchString(line) {
					count++
				}
			}
		}
	}
	return count
}

func countGeneratedFiles(t *testing.T, root string) int {
	t.Helper()
	count := 0
	skipDirs := map[string]bool{
		".git":    true,
		"website": true,
	}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if skipDirs[filepath.Base(path)] {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, "_templ.go") {
			count++
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk generated files: %v", err)
	}
	return count
}

func countIsValidMethods(t *testing.T, root string) int {
	t.Helper()
	count := 0
	isValidRe := regexp.MustCompile(`func\s+[A-Z][A-Za-z0-9]*IsValid\s*\(`)
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
		//nolint:gosec // test reads known package files
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		if matches := isValidRe.FindAll(data, -1); matches != nil {
			count += len(matches)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk IsValid methods: %v", err)
	}
	return count
}

func readDoc(t *testing.T, parts ...string) []byte {
	t.Helper()
	path := filepath.Join(parts...)
	//nolint:gosec // test reads known documentation files
	data, err := os.ReadFile(filepath.Join("..", path))
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func assertCount(t *testing.T, doc []byte, pattern, label string, want int) {
	t.Helper()
	re := regexp.MustCompile(pattern)
	m := re.FindSubmatch(doc)
	if m == nil {
		t.Errorf("%s missing count", label)
		return
	}
	got, err := strconv.Atoi(string(m[1]))
	if err != nil {
		t.Errorf("%s has non-numeric count %q", label, string(m[1]))
		return
	}
	if got != want {
		t.Errorf("%s says %d; actual is %d", label, got, want)
	}
}
