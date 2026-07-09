// Package main implements the tc-css CLI tool for templ-components consumers.
//
// tc-css automates the Tailwind CSS v4 build pipeline:
//  1. Runs `go mod vendor` so vendored .templ files are available for @source scanning
//  2. Auto-generates an input CSS entry-point if none exists
//  3. Runs `tailwindcss` to compile the final CSS
//
// # Usage with go:generate
//
// Add this directive to any .go file in your project:
//
//	//go:generate go run github.com/larsartmann/templ-components/cmd/tc-css -input app.css -output styles.css
//
// Then run:
//
//	go generate ./...
//
// # Usage standalone
//
//	go run github.com/larsartmann/templ-components/cmd/tc-css -input internal/web/static/input.css -output internal/web/static/styles.css
//
// # Usage with BuildFlow
//
// BuildFlow's tailwind-build provider runs this automatically as part of its DAG.
// No go:generate directive is needed when using BuildFlow.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	input := flag.String("input", "app.css", "Input CSS file path (auto-generated if missing)")
	output := flag.String("output", "styles.css", "Output CSS file path")
	minify := flag.Bool("minify", true, "Minify output CSS")
	noVendor := flag.Bool("no-vendor", false, "Skip running 'go mod vendor' before building")
	flag.Parse()

	projectRoot := findProjectRoot()

	if !*noVendor {
		if err := runVendor(projectRoot); err != nil {
			fail("go mod vendor failed: %v", err)
		}
	}

	inputPath := absPath(projectRoot, *input)
	if !fileExists(inputPath) {
		fmt.Fprintf(os.Stderr, "tc-css: input CSS not found at %s — generating a starter file\n", inputPath)
		if err := generateInputCSS(inputPath, projectRoot); err != nil {
			fail("failed to generate input CSS: %v", err)
		}
	}

	outputPath := absPath(projectRoot, *output)

	args := []string{"-i", inputPath, "-o", outputPath}
	if *minify {
		args = append(args, "--minify")
	}

	fmt.Fprintf(os.Stderr, "tc-css: tailwindcss %s\n", strings.Join(args, " "))

	cmd := exec.CommandContext(
		context.Background(),
		"tailwindcss",
		args...) //nolint:gosec // CLI tool intentionally runs subprocesses
	cmd.Dir = projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fail("tailwindcss failed: %v", err)
	}

	fmt.Fprintln(os.Stderr, "tc-css: CSS built successfully.")
}

func runVendor(root string) error {
	cmd := exec.CommandContext(
		context.Background(),
		"go",
		"mod",
		"vendor",
	) //nolint:gosec // CLI tool intentionally runs subprocesses
	cmd.Dir = root
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run() //nolint:wrapcheck // CLI tool, error propagates to main and exits
}

func generateInputCSS(path, projectRoot string) error {
	if err := os.MkdirAll(
		filepath.Dir(path),
		0o750,
	); err != nil { //nolint:wrapcheck // CLI tool, error propagates to main and exits
		return err
	}

	cssDir := filepath.Dir(path)
	sources := findTemplVendoredSources(cssDir, projectRoot)

	var b strings.Builder
	b.WriteString(`@import "tailwindcss" source(none);`)
	b.WriteByte('\n')

	b.WriteString("\n/* Scan this project's templ files */\n")
	relProject, err := filepath.Rel(cssDir, projectRoot)
	if err != nil {
		relProject = "."
	}
	fmt.Fprintf(&b, "@source %q;\n", filepath.Join(relProject, "**/*.templ"))

	if len(sources) > 0 {
		b.WriteString("\n/* Scan vendored templ-components */\n")
		for _, src := range sources {
			fmt.Fprintf(&b, "@source %q;\n", src)
		}
	}

	b.WriteString("\n/* Enable class-based dark mode (required for ThemeScript/ThemeToggle) */\n")
	b.WriteString("@custom-variant dark (&:where(.dark, .dark *));\n")

	b.WriteString("\n/* Override colors without touching component code:\n")
	b.WriteString("@theme {\n")
	b.WriteString("  --color-blue-600: #4f46e5;\n")
	b.WriteString("  --color-blue-500: #6366f1;\n")
	b.WriteString("}\n")
	b.WriteString("*/\n")

	b.WriteString("\n/* For semantic aliases (bg-tc-primary, text-tc-danger), import the theme file:\n")
	b.WriteString("   @import \"./templ-components-theme.css\";\n")
	b.WriteString("*/\n")

	return os.WriteFile(path, []byte(b.String()), 0o600) //nolint:wrapcheck // CLI tool
}

// findTemplVendoredSources scans vendor/ for directories containing .templ
// files and returns their paths relative to cssDir.
func findTemplVendoredSources(cssDir, projectRoot string) []string {
	vendorDir := filepath.Join(projectRoot, "vendor")
	if !dirExists(vendorDir) {
		return nil
	}

	seen := make(map[string]bool)
	var sources []string

	_ = filepath.Walk(
		vendorDir,
		func(path string, info os.FileInfo, _ error) error { //nolint:nilerr // intentionally skip unreadable files
			if info.IsDir() {
				return nil
			}

			if !strings.HasSuffix(path, ".templ") {
				return nil
			}

			dir := filepath.Dir(path)
			rel, err := filepath.Rel(cssDir, dir)
			if err != nil {
				return nil
			}

			if !seen[rel] {
				seen[rel] = true
				sources = append(sources, filepath.Join(rel, "**/*.templ"))
			}

			return nil //nolint:nilerr // walk callback
		},
	)

	sort.Strings(sources)

	return sources
}

func findProjectRoot() string {
	dir, _ := os.Getwd()
	for {
		if fileExists(filepath.Join(dir, "go.mod")) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			cwd, _ := os.Getwd()
			return cwd
		}
		dir = parent
	}
}

func absPath(root, p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(root, p)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "tc-css: "+format+"\n", args...)
	os.Exit(1)
}
