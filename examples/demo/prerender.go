package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/layout"
)

type prerenderPage struct {
	filename string
	title    string
	desc     string
	render   func(layout.PageProps) templ.Component
}

func prerender(outputDir string) error {
	nonce := "demo-nonce"
	cssHead := demoFonts(nonce)

	pages := []prerenderPage{
		{"index.html", "templ-components Demo", "Showcase of all templ-components", demoPage},
		{"forms/index.html", "Forms Demo - templ-components", "Complete form showcase with validation", formsDemoPage},
	}

	ctx := context.Background()

	for _, page := range pages {
		props := layout.DefaultPageProps()
		props.Title = page.title
		props.Description = page.desc
		props.Nonce = nonce
		props.CSSPath = ""
		props.HeadContent = cssHead

		outPath := filepath.Join(outputDir, page.filename)
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return fmt.Errorf("create dir for %s: %w", page.filename, err)
		}

		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("create %s: %w", page.filename, err)
		}

		if err := page.render(props).Render(ctx, f); err != nil {
			f.Close()
			return fmt.Errorf("render %s: %w", page.filename, err)
		}
		f.Close()
		fmt.Printf("  wrote %s\n", outPath)
	}

	return nil
}
