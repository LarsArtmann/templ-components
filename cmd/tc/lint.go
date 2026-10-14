//nolint:gochecknoglobals,nonamedreturns // CLI scanner registry + arg parser; mirrors main.go conventions
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/larsartmann/templ-components/utils"
)

// lintScanner is one compliance check selectable via `tc lint --<flag>`.
type lintScanner struct {
	flag string
	help string
	scan func(dir string) ([]utils.Finding, error)
}

// lintScanners lists every `tc lint` scanner in flag order. The CLI stays in
// lockstep with the library's guard tests: each entry's scan function IS the
// engine behind a utils guard test, so a consumer running `tc lint` enforces
// the same rules the library enforces on itself.
var lintScanners = []lintScanner{
	{
		flag: "rtl",
		help: "physical direction classes (ml-/mr-/pl-/pr-/text-left/text-right/border-l-/border-r-) and deprecated start-*/end-* inset aliases in .templ files",
		scan: utils.ScanRTL,
	},
	{
		flag: "dark",
		help: "neutral + semantic color classes missing a same-line dark: variant (.templ/.go sources)",
		scan: utils.ScanDarkMode,
	},
}

// cmdLint runs the compliance scanners and exits with 1 when findings exist,
// so CI and pre-commit hooks can gate on it.
func cmdLint(_ *registry, args []string) {
	os.Exit(runLint(args, os.Stdout, os.Stderr))
}

// runLint parses `tc lint` arguments, runs the selected scanners over the
// given directories, prints findings, and returns the process exit code:
// 0 = clean, 1 = findings or scan failure, 2 = usage error.
func runLint(args []string, stdout, stderr io.Writer) int {
	selected, dirs, ok := parseLintArgs(args, stderr)
	if !ok {
		return 2
	}

	flags := make([]string, 0, len(selected))
	for _, s := range selected {
		flags = append(flags, "--"+s.flag)
	}

	findings := 0

	for _, dir := range dirs {
		if _, err := os.Stat(dir); err != nil {
			fmt.Fprintf(stderr, "tc lint: no such directory: %s\n", dir)

			return 1
		}

		for _, s := range selected {
			found, err := s.scan(dir)
			if err != nil {
				fmt.Fprintf(stderr, "tc lint: --%s scan failed on %s: %v\n", s.flag, dir, err)

				return 1
			}

			for _, f := range found {
				findings++

				fmt.Fprintf(stdout, "%s:%d: %s: %s\n    %s\n", f.File, f.Line, f.Rule, f.Message, f.Snippet)
			}
		}
	}

	pluralDirs := "directory"
	if len(dirs) != 1 {
		pluralDirs = "directories"
	}

	if findings > 0 {
		fmt.Fprintf(
			stdout,
			"tc lint: %d finding(s) (%s) in %d %s\n",
			findings,
			strings.Join(flags, " "),
			len(dirs),
			pluralDirs,
		)

		return 1
	}

	fmt.Fprintf(stdout, "tc lint: clean — 0 findings (%s) in %d %s\n", strings.Join(flags, " "), len(dirs), pluralDirs)

	return 0
}

// parseLintArgs splits `tc lint` arguments into selected scanners and target
// directories. It prints usage to stderr and reports !ok for usage errors.
func parseLintArgs(args []string, stderr io.Writer) (selected []lintScanner, dirs []string, ok bool) {
	i := 0

	for i < len(args) {
		a := args[i]

		switch {
		case a == "-h" || a == "--help":
			fmt.Fprint(stderr, lintUsage())

			return nil, nil, false
		case strings.HasPrefix(a, "--"):
			name := strings.TrimPrefix(a, "--")
			matched := false

			for _, s := range lintScanners {
				if name == s.flag {
					selected = append(selected, s)
					matched = true

					break
				}
			}

			if !matched {
				fmt.Fprintf(stderr, "tc lint: unknown flag %q\n\n%s", a, lintUsage())

				return nil, nil, false
			}

			i++
		default:
			dirs = append(dirs, a)
			i++
		}
	}

	if len(selected) == 0 {
		fmt.Fprintf(stderr, "tc lint: select at least one scanner\n\n%s", lintUsage())

		return nil, nil, false
	}

	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	return selected, dirs, true
}

func lintUsage() string {
	var b strings.Builder

	b.WriteString("tc lint [--rtl] [DIR ...]\n\n")
	b.WriteString("Scans source files for the library's compliance rules. Without DIR,\n")
	b.WriteString("scans the current directory recursively.\n\nScanners:\n")

	for _, s := range lintScanners {
		fmt.Fprintf(&b, "  --%-8s %s\n", s.flag, s.help)
	}

	b.WriteString("\nExit codes: 0 clean, 1 findings or scan failure, 2 usage error.\n")

	return b.String()
}
