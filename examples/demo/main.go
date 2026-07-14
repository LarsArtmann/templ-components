// Demo application showcasing all templ-components.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/larsartmann/templ-components/layout"
)

func main() {
	prerenderDir := flag.String("prerender", "", "Pre-render demo pages to static HTML in the given directory instead of starting a server")
	flag.Parse()

	if *prerenderDir != "" {
		if err := prerender(*prerenderDir); err != nil {
			fmt.Fprintf(os.Stderr, "Pre-render error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/forms" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			nonce := "demo-nonce"
			props := layout.DefaultPageProps()
			props.Title = "Forms Demo - templ-components"
			props.Description = "Complete form showcase with validation"
			props.Nonce = nonce
			props.CSSPath = ""
			props.HeadContent = tailwindV4CDN(nonce)
			if err := formsDemoPage(props).Render(r.Context(), w); err != nil {
				http.Error(w, err.Error(), 500)
			}
			return
		}
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		nonce := "demo-nonce"
		props := layout.DefaultPageProps()
		props.Title = "templ-components Demo"
		props.Description = "Showcase of all templ-components"
		props.Nonce = nonce
		props.CSSPath = ""
		props.HeadContent = tailwindV4CDN(nonce)
		if err := demoPage(props).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	//nolint:exhaustruct // Demo code - HTTP server for demonstration only
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Println("Demo running at http://localhost:8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
	}
}
