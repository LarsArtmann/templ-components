// Demo application showcasing all templ-components.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/larsartmann/templ-components/layout"
)

func main() {
	prerenderDir := flag.String(
		"prerender",
		"",
		"Pre-render demo pages to static HTML in the given directory instead of starting a server",
	)
	flag.Parse()

	if *prerenderDir != "" {
		if err := prerender(*prerenderDir); err != nil {
			fmt.Fprintf(os.Stderr, "Pre-render error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	mux := http.NewServeMux()

	// Mock HTMX endpoints for interactive demo components
	mux.HandleFunc("/api/load-more", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<div class="rounded-lg border border-gray-200 dark:border-gray-700 p-4 bg-white dark:bg-gray-800">
			<p class="text-sm text-gray-600 dark:text-gray-400">Loaded via HTMX! This item was fetched from the server.</p>
		</div>`)
	})

	mux.HandleFunc("/api/delete", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `<div class="rounded-lg border border-green-200 dark:border-green-800 bg-green-50 dark:bg-green-900/20 p-4">
			<p class="text-sm text-green-800 dark:text-green-200">Item deleted successfully (mock endpoint).</p>
		</div>`)
	})

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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if _, err := strconv.Atoi(port); err != nil {
		port = "8080"
	}
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Demo running at http://localhost:%s\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
	}
}
