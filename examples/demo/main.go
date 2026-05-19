// Demo application showcasing all templ-components.
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/larsartmann/templ-components/layout"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		nonce := "demo-nonce"
		props := layout.DefaultPageProps()
		props.Title = "templ-components Demo"
		props.Description = "Showcase of all templ-components"
		props.Nonce = nonce
		props.CSSPath = ""
		props.HeadContent = tailwindV4CDN(nonce)
		props.HTMXVersion = ""
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
