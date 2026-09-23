// Package distserver serves a built website dist/ directory with Firebase
// cleanUrls semantics. The handler is shared by the capture tools so the
// URL-resolution contract (/foo resolves to foo.html when the exact file is
// absent) lives in exactly one place.
package distserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
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

// Selftest verifies that a dist directory is servable through this package's
// handler: it binds a loopback listener, fetches the given page, and shuts
// down. Both dist-based capture tools call it from their -selftest flag —
// one home for the come-up check (backlog #305/#262).
func Selftest(dist, page string) error {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	server := &http.Server{Handler: Handler(dist)} //nolint:gosec // loopback-only dev server
	defer server.Close()

	go func() { _ = server.Serve(listener) }()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+listener.Addr().String()+"/"+page, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("fetch /%s: %w", page, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch /%s: status %d", page, resp.StatusCode)
	}

	return nil
}
