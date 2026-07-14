// Demo application showcasing all templ-components.
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/a-h/templ"
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

	// Health check endpoint for Cloud Run / container orchestration
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Serve pre-compiled CSS (embedded in binary, no CDN dependency)
	mux.HandleFunc("/css/app.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		_, _ = w.Write(embeddedCSS)
	})

	// Mock HTMX endpoints for interactive demo components
	mux.HandleFunc("/api/load-more", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<div class="rounded-lg border border-gray-200 dark:border-gray-700 p-4 bg-white dark:bg-gray-800">
			<p class="text-sm text-gray-600 dark:text-gray-400">Loaded via HTMX! This item was fetched from the server.</p>
		</div>`)
	})

	mux.HandleFunc("/api/delete", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `<div class="rounded-lg border border-green-200 dark:border-green-800 bg-green-50 dark:bg-green-900/20 p-4">
			<p class="text-sm text-green-800 dark:text-green-200">Item deleted successfully (mock endpoint).</p>
		</div>`)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/forms" {
			renderPage(w, r, "Forms Demo - templ-components", "Complete form showcase with validation", formsDemoPage)
			return
		}
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		renderPage(w, r, "templ-components Demo", "Showcase of all templ-components", demoPage)
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

func renderPage(w http.ResponseWriter, r *http.Request, title, description string, page func(layout.PageProps) templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	props := layout.DefaultPageProps()
	props.Title = title
	props.Description = description
	props.Nonce = "demo-nonce"
	props.CSSPath = "/css/app.css"
	props.HeadContent = demoFonts("demo-nonce")
	if err := page(props).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
