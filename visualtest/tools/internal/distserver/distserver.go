// Package distserver serves a built website dist/ directory with Firebase
// cleanUrls semantics. The handler is shared by the capture tools so the
// URL-resolution contract (/foo resolves to foo.html when the exact file is
// absent) lives in exactly one place.
package distserver

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Handler returns an http.Handler serving dist with Firebase cleanUrls
// semantics: /foo resolves to foo.html when the exact file is absent.
func Handler(dist string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clean, _, _ := strings.Cut(r.URL.Path, "?")
		if !strings.HasSuffix(clean, ".html") && !strings.Contains(clean, ".") {
			cleanPath := filepath.Join(dist, clean+".html")
			if _, err := os.Stat(cleanPath); err == nil { //nolint:gosec // CLI-controlled dist root
				http.ServeFile(w, r, cleanPath) //nolint:gosec // CLI-controlled dist root

				return
			}
		}

		http.FileServer(http.Dir(dist)).ServeHTTP(w, r)
	})

	return mux
}
