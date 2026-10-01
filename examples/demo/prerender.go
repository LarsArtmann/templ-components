package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/layout"
)

type prerenderPage struct {
	filename string
	title    string
	desc     string
	render   func(layout.PageProps) templ.Component
}

// prerenderRequest builds the synthetic request the registry Content funcs
// receive during static export: no query (the wire transport defaults to
// both) and a build-scoped CSRF token in context — a served-static demo page
// cannot pass live CSRF validation (bind moves to the serving server).
func prerenderRequest() *http.Request {
	req, err := http.NewRequest(http.MethodGet, demoURL("/"), nil)
	if err != nil {
		panic(fmt.Sprintf("prerender: build synthetic request: %v", err))
	}
	*req = *req.WithContext(context.WithValue(req.Context(), demoSessionCtxKey{}, newDemoSessionToken()))

	return req
}

func prerender(outputDir string) error {
	nonce := demoNonceConst
	cssHead := demoFonts(nonce)
	req := prerenderRequest()

	pages := []prerenderPage{
		{
			"index.html",
			"templ-components Demo - live component showcase",
			homePageMeta.Short,
			func(_ layout.PageProps) templ.Component {
				return demoShell(homePageMeta, nonce, demoHomeContent(nonce))
			},
		},
	}

	for _, page := range demoPages() {
		meta := page
		pages = append(pages, prerenderPage{
			filename: strings.TrimPrefix(meta.Path, "/") + "/index.html",
			title:    meta.Title + " - templ-components Demo",
			desc:     meta.Short,
			render: func(_ layout.PageProps) templ.Component {
				return demoShell(meta, nonce, meta.Content(req))
			},
		})
	}

	pages = append(
		pages,
		prerenderPage{
			"recipes/dashboard.html",
			"Dashboard Recipe - templ-components",
			"Dashboard recipe demo",
			recipesDashboardPage,
		},
		prerenderPage{
			"recipes/settings.html",
			"Settings Recipe - templ-components",
			"Settings recipe demo",
			recipesSettingsPage,
		},
		prerenderPage{
			"recipes/login.html",
			"Login Recipe - templ-components",
			"Login card recipe demo",
			recipesLoginPage,
		},
		prerenderPage{
			"recipes/auth.html",
			"Auth Layout Recipe - templ-components",
			"Auth layout recipe demo",
			recipesAuthPage,
		},
	)

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
			return fmt.Errorf(
				"create dir for %s (outPath=%s, outputDir=%s): %w",
				page.filename,
				outPath,
				outputDir,
				err,
			)
		}

		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("create outPath=%s: %w", outPath, err)
		}

		if err := page.render(props).Render(ctx, f); err != nil {
			f.Close()
			return fmt.Errorf("render outPath=%s: %w", outPath, err)
		}
		f.Close()
		fmt.Printf("  wrote %s\n", outPath)
	}

	return nil
}
