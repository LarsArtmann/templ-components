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
	"encoding/json/v2"
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
	firebasePath := flag.String("firebase-config", "firebase.json", "path to firebase.json for the CSP header guard")
	updateCSP := flag.Bool("update-csp", false, "rewrite the firebase.json CSP header instead of only checking it")

	flag.Parse()

	cfg := config{
		outDir:       *out,
		repoRoot:     *repoRoot,
		firebasePath: *firebasePath,
		skipStars:    *skipStars,
		updateCSP:    *updateCSP,
	}

	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

// config is one site-build invocation.
type config struct {
	outDir       string
	repoRoot     string
	firebasePath string
	skipStars    bool
	updateCSP    bool
}

// topLevelPages returns the non-docs HTML pages the site renders (index,
// sales, 404) — the single structural source for the page set and its count.
func topLevelPages(stats build.Stats, starsLabel, nonce string) []build.Page {
	return []build.Page{
		{Path: "index.html", Component: pages.Landing(stats, starsLabel, nonce)},
		{Path: "sales.html", Component: pages.Sales(stats, starsLabel, nonce)},
		{Path: "404.html", Component: pages.NotFound(nonce)},
	}
}

func run(cfg config) error {
	nonce, err := build.Nonce()
	if err != nil {
		return fmt.Errorf("build nonce: %w", err)
	}

	stats, err := build.CountStats(cfg.repoRoot)
	if err != nil {
		return fmt.Errorf("library stats: %w", err)
	}

	stars := 0
	if !cfg.skipStars {
		stars = fetchStars()
	}

	renderer := build.NewRenderer(nonce)

	docsPages, searchDocs, err := renderDocs(cfg.repoRoot, nonce)
	if err != nil {
		return err
	}

	ctx := context.Background()

	topPages := topLevelPages(stats, pages.StarsLabel(stars), nonce)
	sitePages := make([]build.Page, 0, len(topPages)+len(docsPages))
	sitePages = append(sitePages, topPages...)
	sitePages = append(sitePages, docsPages...)

	rendered, err := renderer.RenderPages(ctx, sitePages)
	if err != nil {
		return fmt.Errorf("render pages: %w", err)
	}

	if err := syncCSP(cfg, rendered); err != nil {
		return err
	}

	if err := build.WriteRendered(cfg.outDir, rendered); err != nil {
		return fmt.Errorf("write pages: %w", err)
	}

	if err := writeSitemaps(cfg.outDir, cfg.repoRoot); err != nil {
		return fmt.Errorf("sitemaps: %w", err)
	}

	if err := build.WriteSearchIndex(cfg.outDir, searchDocs); err != nil {
		return fmt.Errorf("search index: %w", err)
	}

	if err := writeAssets(cfg.outDir, cfg.repoRoot); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "site: wrote %d page(s) + sitemaps to %s (components=%d icons=%d enums=%d modules=%d)\n",
		len(sitePages), cfg.outDir, stats.Components, stats.Icons, stats.Enums, stats.Modules)

	return nil
}

// syncCSP pins the rendered pages' inline scripts into the firebase.json CSP
// header: --update-csp rewrites the config, the default verifies it so a
// stale committed header fails the build instead of deploying broken pages.
func syncCSP(cfg config, rendered []build.RenderedPage) error {
	header := build.CSPHeader(build.InlineScriptHashes(rendered))

	if cfg.updateCSP {
		changed, err := build.SyncFirebaseCSP(cfg.firebasePath, header)
		if err != nil {
			return fmt.Errorf("sync firebase CSP: %w", err)
		}

		if changed {
			fmt.Fprintf(os.Stdout, "site: updated %s CSP header\n", cfg.firebasePath)
		}

		return nil
	}

	if err := build.CheckFirebaseCSP(cfg.firebasePath, header); err != nil {
		return fmt.Errorf("CSP guard: %w", err)
	}

	return nil
}

// writeAssets emits the generated chroma stylesheet and copies the static
// asset trees (assets/, public/) into the dist directory. Source paths are
// resolved under repoRoot so the build does not depend on the working
// directory.
func writeAssets(out, repoRoot string) error {
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

	websiteRoot := filepath.Join(repoRoot, "website")

	for _, tree := range []struct{ src, dst string }{
		{filepath.Join(websiteRoot, "assets"), filepath.Join(out, "assets")},
		{filepath.Join(websiteRoot, "public"), out},
	} {
		if _, err := os.Stat(tree.src); err == nil {
			if err := build.CopyTree(tree.src, tree.dst); err != nil {
				return fmt.Errorf("copy %s: %w", tree.src, err)
			}
		}
	}

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
	if err := json.UnmarshalRead(resp.Body, &payload); err != nil {
		return 0
	}

	return payload.StargazersCount
}

// renderDocs renders every registered docs page from content/docs, wiring
// prev/next navigation and git last-updated dates, and collects the
// client-side search index entries.
func renderDocs(repoRoot, nonce string) ([]build.Page, []build.SearchDoc, error) {
	all := pages.AllDocs()

	out := make([]build.Page, 0, len(all))
	searchDocs := make([]build.SearchDoc, 0, len(all))

	for index, doc := range all {
		source, err := os.ReadFile(filepath.Join(repoRoot, "website", "content", "docs", doc.Slug+".md"))
		if err != nil {
			return nil, nil, fmt.Errorf("read docs %s: %w", doc.Slug, err)
		}

		parsed, err := md.Parse(string(source))
		if err != nil {
			return nil, nil, fmt.Errorf("parse docs %s: %w", doc.Slug, err)
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

		sections := make([]build.SearchSection, 0, len(parsed.Headings))
		for _, heading := range parsed.Headings {
			sections = append(sections, build.SearchSection{ID: heading.ID, Text: heading.Text})
		}

		searchDocs = append(searchDocs, build.SearchDoc{
			URL:         "/" + doc.Slug,
			Title:       parsed.Title,
			Description: parsed.Description,
			Sections:    sections,
			Body:        build.PlainText(parsed.HTML),
		})

		out = append(out, build.Page{
			Path: doc.Slug + ".html",

			Component: pages.DocsLayout(doc.Slug, parsed, prev, next, lastUpdated(repoRoot, doc.Slug), nonce),
		})
	}

	return out, searchDocs, nil
}

// lastUpdated asks git for the last commit date touching a docs source
// (YYYY-MM-DD); empty when unavailable (shallow clones, non-git runs).
// lastUpdated returns the committer date (YYYY-MM-DD) of the newest commit
// touching any of the given repo-relative paths — the sitemap lastmod source.
// Empty when git is unavailable or the paths have no history.
func lastUpdated(repoRoot string, paths ...string) string {
	if len(paths) == 0 {
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []string{"-C", repoRoot, "log", "-1", "--format=%cs", "--"}
	args = append(args, paths...)

	//nolint:gosec // repository-controlled paths, never user input
	cmd := exec.CommandContext(ctx, "git", args...)

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
// sitemap-index.xml (the URL robots.txt already references). Non-docs pages
// take their lastmod from the page's primary .templ source(s); docs pages
// from their content markdown. CONVENTION: a new top-level page must be
// added to topLevelPages AND get a lastmod source here — one entry per
// rendered top-level page, never a hand-kept count.
func writeSitemaps(outDir, repoRoot string) error {
	entries := make([]sitemapEntry, 0, 2+len(pages.AllDocs()))
	entries = append(
		entries,
		sitemapEntry{loc: pages.SiteURL + "/", lastmod: lastUpdated(repoRoot, "website/internal/pages/landing.templ", "website/internal/pages/hero.templ")},
		sitemapEntry{loc: pages.SiteURL + "/sales", lastmod: lastUpdated(repoRoot, "website/internal/pages/sales.templ")},
	)

	for _, doc := range pages.AllDocs() {
		entries = append(entries, sitemapEntry{
			loc:     pages.SiteURL + "/" + doc.Slug,
			lastmod: lastUpdated(repoRoot, "website/content/docs/"+doc.Slug+".md"),
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
	//nolint:gosec // public site asset must be world-readable
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
	//nolint:gosec // public site asset must be world-readable
	if err := os.WriteFile(indexPath, []byte(index), 0o644); err != nil {
		return fmt.Errorf("write sitemap index: %w", err)
	}

	return nil
}
