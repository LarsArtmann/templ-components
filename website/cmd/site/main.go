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
		return err
	}
	stats, err := build.CountStats(repoRoot)
	if err != nil {
		return err
	}

	stars := 0
	if !skipStars {
		stars = fetchStars()
	}

	renderer := build.NewRenderer(nonce)
	ctx := context.Background()
	pages := []build.Page{
		{Path: "index.html", Component: pages.Landing(stats, pages.StarsLabel(stars), nonce)},
	}
	if err := renderer.WritePages(ctx, out, pages); err != nil {
		return err
	}
	fmt.Printf("site: wrote %d page(s) to %s (components=%d icons=%d enums=%d modules=%d)\n",
		len(pages), out, stats.Components, stats.Icons, stats.Enums, stats.Modules)
	return nil
}

// fetchStars mirrors the Astro hero's build-time star lookup: fail-soft, any
// error means the static "Star on GitHub" fallback label.
func fetchStars() int {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodGet, gitHubAPIURL, nil)
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
		StargazersCount int `json:"stargazers_count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0
	}
	return payload.StargazersCount
}
