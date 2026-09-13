// Command site is the static-site generator for the templ-components
// website. It renders every page through the library's own templ components
// and writes a deployable dist/ directory for Firebase Hosting.
//
// Usage (from the website/ module root):
//
//	GOEXPERIMENT=jsonv2 go run ./cmd/site --out dist --repo-root ..
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/larsartmann/templ-components/website/internal/build"
	"github.com/larsartmann/templ-components/website/internal/pages"
)

const gitHubAPIURL = "https://api.github.com/repos/larsartmann/templ-components"

func main() {
	out := flag.String("out", "dist", "output directory (relative to the website module)")
	repoRoot := flag.String("repo-root", "..", "repository root for library-fact statistics")
	skipStars := flag.Bool("skip-stars", false, "skip the GitHub star lookup (offline builds)")

	flag.Parse()

	if err := run(*out, *repoRoot, *skipStars); err != nil {
		log.Fatal(err)
	}
}

func run(out, repoRoot string, skipStars bool) error {
	nonce, err := build.Nonce()
	if err != nil {
		return fmt.Errorf("build nonce: %w", err)
	}

	stats, err := build.CountStats(repoRoot)
	if err != nil {
		return fmt.Errorf("library stats: %w", err)
	}

	stars := 0
	if !skipStars {
		stars = fetchStars()
	}

	renderer := build.NewRenderer(nonce)
	ctx := context.Background()

	sitePages := []build.Page{
		{Path: "index.html", Component: pages.Landing(stats, pages.StarsLabel(stars), nonce)},
	}
	if err := renderer.WritePages(ctx, out, sitePages); err != nil {
		return fmt.Errorf("write pages: %w", err)
	}

	chromaCSS, err := build.ChromaCSS()
	if err != nil {
		return fmt.Errorf("chroma css: %w", err)
	}

	chromaPath := filepath.Join(out, "assets", "css", "chroma.css")
	if err := os.MkdirAll(filepath.Dir(chromaPath), 0o755); err != nil { //nolint:gosec // public site asset directory
		return fmt.Errorf("create chroma css dir: %w", err)
	}

	if err := os.WriteFile(chromaPath, []byte(chromaCSS), 0o644); err != nil { //nolint:gosec // public site asset
		return fmt.Errorf("write chroma css: %w", err)
	}

	for _, tree := range []struct{ src, dst string }{
		{"assets", filepath.Join(out, "assets")},
		{"public", out},
	} {
		if _, err := os.Stat(tree.src); err == nil {
			if err := build.CopyTree(tree.src, tree.dst); err != nil {
				return fmt.Errorf("copy %s: %w", tree.src, err)
			}
		}
	}

	fmt.Fprintf(os.Stdout, "site: wrote %d page(s) to %s (components=%d icons=%d enums=%d modules=%d)\n",
		len(sitePages), out, stats.Components, stats.Icons, stats.Enums, stats.Modules)

	return nil
}

// fetchStars mirrors the Astro hero's build-time star lookup: fail-soft, any
// error means the static "Star on GitHub" fallback label.
func fetchStars() int {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       5 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, gitHubAPIURL, nil)
	if err != nil {
		return 0
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0
	}

	var payload struct {
		StargazersCount int `json:"stargazers_count"` //nolint:tagliatelle // GitHub API wire format
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0
	}

	return payload.StargazersCount
}
