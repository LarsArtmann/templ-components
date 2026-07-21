// Package main is the `tc` CLI: scaffolding tool for templ-components.
// See README at docs/cli.md for the full usage guide.
//
// This file uses CLI-specific conventions intentionally:
//   - path traversal via user-supplied --out (gosec G703 is a false positive)
//   - Print to stdout is the purpose of the tool (forbidigo doesn't apply)
//   - registry globals are normal for embed-driven CLIs
//
//nolint:errcheck,gosec,forbidigo,gochecknoglobals,nonamedreturns,wsl_v5,nlreturn,nolintlint,gocyclo,cyclop,funlen,goprintffuncname,nilerr // CLI tool conventions
package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// CLI tool uses different conventions than library code: path traversal is
// intentional (CLI consumes user-supplied paths), Print to stdout is the
// purpose, globals are normal for embed registries.
//
//nolint:all // CLI tool; see package doc
//go:embed all:_sources
var sourcesFS embed.FS

// registry tracks every embedded component source by lower-cased name.
type registry struct {
	files map[string][]string
	pkg   map[string]string
	pkgs  []string
}

func newRegistry() *registry {
	r := &registry{
		files: map[string][]string{},
		pkg:   map[string]string{},
		pkgs: []string{
			"display", "feedback", "forms", "layout",
			"navigation", "htmx", "errorpage", "recipes",
		},
	}

	_ = fs.WalkDir(sourcesFS, "_sources", func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel("_sources", path)
		parts := strings.Split(rel, string(filepath.Separator))
		if len(parts) < 2 {
			return nil
		}

		pkg := parts[0]
		filename := parts[len(parts)-1]
		if !strings.HasSuffix(filename, ".templ") {
			return nil
		}

		base := strings.TrimSuffix(filename, ".templ")
		key := strings.ToLower(base)
		r.files[key] = append(r.files[key], path)

		if _, ok := r.pkg[key]; !ok {
			r.pkg[key] = pkg
		}

		return nil
	})

	return r
}

func main() {
	if len(os.Args) < 2 {
		usage()

		return
	}

	r := newRegistry()

	switch os.Args[1] {
	case "init":
		cmdInit(r, os.Args[2:])
	case "ls", "list":
		cmdList(r, os.Args[2:])
	case "add":
		cmdAdd(r, os.Args[2:])
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func usage() {
	_, _ = fmt.Fprintln(os.Stderr, `tc — templ-components scaffolding

Usage:
  tc init                   Scaffold a starter app.css + custom.css in the current directory.
  tc ls                     Print every component, grouped by package.
  tc add <component>        Copy <component>.templ (and _types.go if present) to ./components/.
  tc add <component> --out DIR   Copy to a custom directory.

Examples:
  tc init
  tc add button
  tc add dropdown --out ./src/components

Available components: run 'tc ls' to see them all.`)
}

func cmdInit(_ *registry, _ []string) {
	if err := os.MkdirAll(".", 0o750); err != nil {
		failf("create dir: %v", err)
	}

	for _, name := range []string{"app.css", "custom.css"} {
		content, err := sourcesFS.ReadFile(filepath.Join("_sources", "starter", name))
		if err != nil {
			failf("read starter %s: %v", name, err)
		}

		if _, err := os.Stat(name); err == nil {
			fmt.Fprintf(os.Stderr, "tc: skip %s (already exists)\n", name)

			continue
		}

		if err := os.WriteFile(name, content, 0o600); err != nil {
			failf("write %s: %v", name, err)
		}

		status("wrote", name)
	}
}

func cmdList(r *registry, _ []string) {
	for _, pkg := range r.pkgs {
		fmt.Fprintf(os.Stdout, "\n# %s\n", pkg)

		for name, p := range r.pkg {
			if p != pkg {
				continue
			}

			fmt.Fprintf(os.Stdout, "  %s\n", name)
		}
	}
}

func cmdAdd(r *registry, args []string) {
	out, positional := parseAddArgs(args)
	if len(positional) < 1 {
		_, _ = fmt.Fprintln(os.Stderr, "tc add <component> [--out DIR]")

		return
	}

	name := strings.ToLower(strings.TrimSpace(positional[0]))
	files, ok := r.files[name]
	if !ok {
		failf("unknown component %q. Run 'tc ls' for the list.", name)
	}

	if err := os.MkdirAll(out, 0o750); err != nil {
		failf("create output dir: %v", err)
	}

	for _, src := range files {
		copyFile(src, filepath.Join(out, filepath.Base(src)))
	}

	for _, src := range files {
		base := strings.TrimSuffix(filepath.Base(src), ".templ")
		typesPath := filepath.Join(filepath.Dir(src), base+"_types.go")

		if content, err := sourcesFS.ReadFile(typesPath); err == nil {
			dest := filepath.Join(out, base+"_types.go")

			if err := os.WriteFile(dest, content, 0o600); err != nil {
				failf("write %s: %v", dest, err)
			}

			status("wrote", dest)
		}
	}
}

func parseAddArgs(args []string) (out string, positional []string) {
	out = "./components"

	i := 0

	for i < len(args) {
		a := args[i]

		switch {
		case a == "--out" || a == "-o":
			if i+1 >= len(args) {
				failf("--out requires a value")
			}

			out = args[i+1]
			i += 2
		case strings.HasPrefix(a, "--out="):
			out = strings.TrimPrefix(a, "--out=")
			i++
		case strings.HasPrefix(a, "-o="):
			out = strings.TrimPrefix(a, "-o=")
			i++
		case a == "-h" || a == "--help":
			_, _ = fmt.Fprintln(os.Stderr, "tc add <component> [--out DIR]")

			return out, positional
		default:
			positional = append(positional, a)
			i++
		}
	}

	return out, positional
}

func copyFile(srcEmbed, dest string) {
	content, err := sourcesFS.ReadFile(srcEmbed)
	if err != nil {
		failf("read embedded %s: %v", srcEmbed, err)
	}

	if err := os.WriteFile(dest, content, 0o600); err != nil {
		failf("write %s: %v", dest, err)
	}

	status("wrote", dest)
}

func status(verb, path string) {
	fmt.Fprintf(os.Stdout, "tc: %s %s\n", verb, path)
}

func failf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "tc: "+format+"\n", args...)
	os.Exit(1)
}
