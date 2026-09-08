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
	"sort"
	"strings"

	"github.com/larsartmann/templ-components/utils"
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
			pkgDisplay, pkgFeedback, pkgForms, pkgLayout,
			pkgNav, pkgHTMX, pkgDatastar, pkgErrorpage, pkgRecipes,
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

const enumsGoFile = "enums_go.go"

// Package-name constants: the literals repeat across the registry list, the
// packageDeps/packageImports maps, and the datastar copy special-casing
// (goconst keeps them single-sourced).
const (
	pkgDisplay   = "display"
	pkgFeedback  = "feedback"
	pkgForms     = "forms"
	pkgLayout    = "layout"
	pkgNav       = "navigation"
	pkgHTMX      = "htmx"
	pkgDatastar  = "datastar"
	pkgErrorpage = "errorpage"
	pkgRecipes   = "recipes"
)

// Import-path constants for the packageImports checklist: the same literals
// repeat across entries (goconst), and centralizing them keeps the map
// diffable against the drift test.
const (
	impTempl        = "github.com/a-h/templ"
	impUtils        = "github.com/larsartmann/templ-components/utils"
	impUtilsSVG     = impUtils + "/svg"
	impUtilsWire    = impUtils + "/wire"
	impUtilsCDN     = impUtils + "/cdn"
	impIcons        = "github.com/larsartmann/templ-components/icons"
	impDisplay      = "github.com/larsartmann/templ-components/display"
	impLayout       = "github.com/larsartmann/templ-components/layout"
	impDatastarStat = "github.com/larsartmann/go-datastar/static"
	impErrorFamily  = "github.com/larsartmann/go-error-family"
	impHTMX         = "github.com/larsartmann/templ-components/htmx"
)

// datastarBumpProtocolDoc is copied alongside every datastar component add:
// the package integrates with the pinned go-datastar/static runtime bundle,
// so vendored copies inherit the bump/re-audit checklist by construction.
const datastarBumpProtocolDoc = "DATASTAR-BUMP-PROTOCOL.md"

// packageImports lists the module-level imports each package's non-test,
// non-generated sources use. 'tc add --list-deps' prints them as the
// go.mod checklist a vendoring consumer needs. Guarded against drift by
// TestPackageImportsMatchSources.
var packageImports = map[string][]string{
	pkgDisplay: {
		impTempl,
		impHTMX,
		impIcons,
		impUtils,
		impUtilsSVG,
		impUtilsWire,
	},

	pkgFeedback: {
		impIcons,
		impUtils,
		impUtilsSVG,
	},
	pkgForms: {
		impTempl,
		impIcons,
		impUtils,
		impUtilsWire,
	},
	pkgLayout: {
		impTempl,
		impIcons,
		impUtils,
		impUtilsCDN,
	},
	pkgNav: {
		impIcons,
		impUtils,
		impUtilsSVG,
		impUtilsWire,
	},
	pkgHTMX: {
		impTempl,
		impUtils,
	},

	pkgDatastar: {
		impTempl,
		impDatastarStat,
		impUtils,
		impUtilsCDN,
	},

	pkgErrorpage: {
		impTempl,
		impErrorFamily,
		impIcons,
		impUtils,
	},
	pkgRecipes: {
		impTempl,
		impDisplay,
		impIcons,
		impLayout,
		impUtils,
	},
}

// packageDeps lists the non-test, non-generated, non-types .go files in each
// package. These are the sibling files that 'tc add' does NOT copy but that the
// .templ source references (class lookups, enums, shared helpers, etc.).
// Use 'tc add <component> --list-deps' to print them.
var packageDeps = map[string][]string{
	pkgDisplay: {
		"bar_chart.go", "button_go.go", "collapsible_section.go",
		"drawer_go.go", enumsGoFile, "external_link.go",
		"heatmap.go", "modal_go.go", "shared.go", "sparkline.go",
	},
	pkgFeedback: {enumsGoFile, "styles.go"},
	pkgForms: {
		"aria.go", enumsGoFile, "ids.go",
		"input_classes.go", "radio.go",
	},
	pkgLayout: {
		"appshell_types.go", "container_types.go",
		"sri.go", "split_types.go", "stack_types.go",
	},

	pkgNav: {},
	pkgHTMX: {
		enumsGoFile, "polled_region.go",
		"view_transitions.go",
	},

	pkgDatastar: {
		"cancellation.go", enumsGoFile, "indicator.go",
		"live_region.go", "retry.go", "sdk_script.go",
		// version.go additionally pulls in go-datastar/static (embedded
		// runtime bundle) + utils/cdn — SDKScript consumers only.
		"version.go",
	},

	pkgErrorpage: {
		"constructors.go", "fromerror.go", "handler.go",
		"notfound404_types.go", "styles.go",
	},
	pkgRecipes: {
		"dashboard_types.go", "login_card_types.go",
		"settings_layout_types.go",
	},
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
	case "version", "-v", "--version":
		fmt.Fprintln(os.Stdout, utils.Version)
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
  tc add <component> --list-deps List sibling .go files the component depends on.
  tc version                Print the library version.

Examples:
  tc init
  tc add button
  tc add dropdown --out ./src/components
  tc add modal --list-deps

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
	out, listDeps, positional := parseAddArgs(args)
	if len(positional) < 1 {
		_, _ = fmt.Fprintln(os.Stderr, "tc add <component> [--out DIR] [--list-deps]")

		return
	}

	name := strings.ToLower(strings.TrimSpace(positional[0]))
	files, ok := r.files[name]
	if !ok {
		failf("unknown component %q. Run 'tc ls' for the list.", name)
	}

	pkg := r.pkg[name]

	if listDeps {
		deps := packageDeps[pkg]
		sort.Strings(deps)
		fmt.Fprintf(os.Stdout, "# sibling .go files in package '%s' needed by '%s':\n", pkg, name)
		for _, dep := range deps {
			fmt.Fprintf(os.Stdout, "  %s/%s\n", pkg, dep)
		}
		imports := packageImports[pkg]
		sort.Strings(imports)
		fmt.Fprintf(os.Stdout, "\n# imports the vendored '%s' package pulls into go.mod:\n", pkg)
		for _, imp := range imports {
			fmt.Fprintf(os.Stdout, "  %s\n", imp)
		}
		if pkg == pkgDatastar {
			fmt.Fprintf(os.Stderr, "\ntc: datastar vendors an external runtime bundle — the add copies\n")
			fmt.Fprintf(os.Stderr, "      %s with the re-audit checklist. Follow it on every\n",
				datastarBumpProtocolDoc)
			fmt.Fprintf(os.Stderr, "      go-datastar/static bump.\n")
		}
		fmt.Fprintf(os.Stderr, "\ntc: %d file(s). These are NOT copied by 'tc add'.\n", len(deps))
		fmt.Fprintf(os.Stderr, "tc: to get a working component, vendor the full package:\n")
		fmt.Fprintf(os.Stderr, "      go get github.com/larsartmann/templ-components/%s\n", pkg)

		return
	}

	if err := os.MkdirAll(out, 0o750); err != nil {
		failf("create output dir: %v", err)
	}

	for _, src := range files {
		copyFile(src, filepath.Join(out, filepath.Base(src)))
	}

	// datastar components vendor an external runtime integration: always
	// drop the bump/re-audit checklist next to them (post-audit patterns by
	// construction).
	if pkg == pkgDatastar {
		copyFile(filepath.Join("_sources", pkgDatastar, datastarBumpProtocolDoc),
			filepath.Join(out, datastarBumpProtocolDoc))
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

	// Warn: a .templ file references package-level helpers (class lookups,
	// enums, utils, sub-templates) defined in sibling .go files that are NOT
	// embedded or copied. The output will not compile standalone.
	fmt.Fprintf(os.Stderr,
		"tc: note: '%s.templ' is part of the '%s' package and references helpers\n"+
			"    (class lookups, enums, sub-templates) defined in sibling .go files\n"+
			"    that were NOT copied. The output will not compile standalone.\n"+
			"    For a working component, vendor the full package:\n"+
			"      go get github.com/larsartmann/templ-components/%s\n",
		name, pkg, pkg)
}

func parseAddArgs(args []string) (out string, listDeps bool, positional []string) {
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
		case a == "--list-deps":
			listDeps = true
			i++
		case a == "-h" || a == "--help":
			_, _ = fmt.Fprintln(os.Stderr, "tc add <component> [--out DIR] [--list-deps]")

			return out, listDeps, positional
		default:
			positional = append(positional, a)
			i++
		}
	}

	return out, listDeps, positional
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
