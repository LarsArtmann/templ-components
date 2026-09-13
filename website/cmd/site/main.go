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
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/larsartmann/templ-components/website/internal/build"
	"github.com/larsartmann/templ-components/website/internal/md"
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

	docsPages, err := renderDocs(ctx, renderer, repoRoot, nonce)
	if err != nil {
		return err
	}

	sitePages = append(sitePages, docsPages...)

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

	fmt.Fprintf(os.Stdout, "site: wrote %d page(s) + sitemaps to %s (components=%d icons=%d enums=%d modules=%d)\n",
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

// renderDocs renders every registered docs page from content/docs, wiring
// prev/next navigation and git last-updated dates.
func renderDocs(ctx context.Context, renderer *build.Renderer, repoRoot, nonce string) ([]build.Page, error) {
	all := pages.AllDocs()

	out := make([]build.Page, 0, len(all))

	for index, doc := range all {
		source, err := os.ReadFile(filepath.Join(repoRoot, "website", "content", "docs", doc.Slug+".md"))
		if err != nil {
			return nil, fmt.Errorf("read docs %s: %w", doc.Slug, err)
		}

		parsed, err := md.Parse(string(source))
		if err != nil {
			return nil, fmt.Errorf("parse docs %s: %w", doc.Slug, err)
		}

		var prev, next *pages.DocRef

		if index > 0 {
			previous := all[index-1]
			prev = &previous
		}

		if index < len(all)-1 {
			following := all[index+1]
			next = &following
		}

		out = append(out, build.Page{
			Path:      doc.Slug + ".html",
			Component: pages.DocsLayout(doc.Slug, parsed, prev, next, lastUpdated(repoRoot, doc.Slug), nonce),
		})
	}

	return out, nil
}

// lastUpdated asks git for the last commit date touching a docs source
// (YYYY-MM-DD); empty when unavailable (shallow clones, non-git runs).
func lastUpdated(repoRoot, slug string) string {
	cmd := exec.Command("git", "-C", repoRoot, "log", "-1", "--format=%cs", "--", "website/content/docs/"+slug+".md")

	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

// sitemapEntry is one URL in the sitemap.
type sitemapEntry struct {
	loc     string
	lastmod string
}

// writeSitemaps emits sitemap.xml (all pages, git lastmod where known) and
// sitemap-index.xml (the URL robots.txt already references).
func writeSitemaps(outDir, repoRoot string) error {
	entries := []sitemapEntry{{loc: pages.SiteURL + "/"}}

	for _, doc := range pages.AllDocs() {
		entries = append(entries, sitemapEntry{
			loc:     pages.SiteURL + "/" + doc.Slug,
			lastmod: lastUpdated(repoRoot, doc.Slug),
		})
	}

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sb.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")

	for _, entry := range entries {
		sb.WriteString("\t<url>\n\t\t<loc>" + entry.loc + "</loc>\n")

		if entry.lastmod != "" {
			sb.WriteString("\t\t<lastmod>" + entry.lastmod + "</lastmod>\n")
		}

		sb.WriteString("\t</url>\n")
	}

	sb.WriteString("</urlset>\n")

	sitemapPath := filepath.Join(outDir, "sitemap.xml")
	if err := os.WriteFile(sitemapPath, []byte(sb.String()), 0o644); err != nil {
		return fmt.Errorf("write sitemap: %w", err)
	}

	index := `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<sitemap>
		<loc>` + pages.SiteURL + `/sitemap.xml</loc>
	</sitemap>
</sitemapindex>
`
	indexPath := filepath.Join(outDir, "sitemap-index.xml")
	if err := os.WriteFile(indexPath, []byte(index), 0o644); err != nil {
		return fmt.Errorf("write sitemap index: %w", err)
	}

	return nil
}
