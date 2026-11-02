// Command countgen keeps the count-bearing tables in FEATURES.md and
// skill/SKILL.md derived from the tree instead of hand-typed. Check mode (default)
// exits non-zero listing every drift; -w rewrites the docs.
//
//	go run ./tools/countgen/cmd/countgen          # check
//	go run ./tools/countgen/cmd/countgen -w      # write
//
// utils.TestDocsCountDrift is the independent verifier — it re-derives every
// count and asserts the docs without this package's code (see package docs).
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/larsartmann/templ-components/tools/countgen"
)

func main() {
	write := flag.Bool("w", false, "rewrite FEATURES.md/SKILL.md to match the tree")
	flag.Parse()

	root := "."

	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	if *write {
		changed, err := countgen.Apply(root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "countgen: %v\n", err)
			os.Exit(1)
		}

		if len(changed) == 0 {
			fmt.Fprintln(os.Stderr, "countgen: docs already match the tree")

			return
		}

		for _, name := range changed {
			fmt.Fprintf(os.Stderr, "countgen: wrote %s\n", name)
		}

		return
	}

	if err := countgen.Check(root); err != nil {
		fmt.Fprintf(os.Stderr, "countgen: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "countgen: OK — every count-bearing site matches the tree")
}
