package main

import (
	"encoding/json/v2"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/website/internal/build"
	"github.com/larsartmann/templ-components/website/internal/pages"
)

var scriptTagRe = regexp.MustCompile(`(?s)<script\b([^>]*)>`)

// TestSiteBuildIntegrity renders the whole site into a temp dir and verifies
// the deployment invariants: page count, internal links, script nonces, the
// search index, and the sitemap.
func TestSiteBuildIntegrity(t *testing.T) {
	t.Parallel()

	outDir := t.TempDir()
	repoRoot := "../../.."

	cfg := config{
		outDir:       outDir,
		repoRoot:     repoRoot,
		firebasePath: filepath.Join(repoRoot, "website", "firebase.json"),
		skipStars:    true,
		updateCSP:    false,
	}

	if err := run(cfg); err != nil {
		t.Fatalf("site build: %v", err)
	}

	rendered, assets := collectDist(t, outDir)

	// assets/app.css is compiled by build.sh's Tailwind step, outside the Go
	// SSG; declare it so stylesheet references validate in the test dist.
	assets = append(assets, "assets/app.css")

	// The expected page count is derived structurally: the same topLevelPages
	// source run() uses, plus the docs set — no hand-kept constants.
	stats, err := build.CountStats(repoRoot)
	if err != nil {
		t.Fatalf("library stats: %v", err)
	}

	wantPages := len(topLevelPages(stats, pages.StarsLabel(0), "")) + len(pages.AllDocs())
	if len(rendered) != wantPages {
		t.Fatalf("wrote %d HTML pages, want %d", len(rendered), wantPages)
	}

	if problems := build.CheckLinks(rendered, assets); len(problems) != 0 {
		t.Errorf("%d broken links:\n%s", len(problems), strings.Join(problems, "\n"))
	}

	assertScriptNonces(t, rendered)
	assertSearchIndex(t, outDir, rendered)
	assertSitemap(t, outDir, rendered)
	assertNoFrameworkScripts(t, rendered)
}

// collectDist reads every dist file back as rendered pages + asset paths.
func collectDist(t *testing.T, outDir string) ([]build.RenderedPage, []string) {
	t.Helper()

	var rendered []build.RenderedPage

	var assets []string

	err := filepath.WalkDir(outDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(outDir, path)
		if err != nil {
			return err
		}

		rel = filepath.ToSlash(rel)

		if !strings.HasSuffix(rel, ".html") {
			assets = append(assets, rel)

			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		rendered = append(rendered, build.RenderedPage{Path: rel, HTML: string(data)})

		return nil
	})
	if err != nil {
		t.Fatalf("walk dist: %v", err)
	}

	return rendered, assets
}

// assertScriptNonces enforces the site's script policy: every <script> is
// either external same-origin (src=) or carries a non-empty nonce. JSON-LD
// blocks are exempt — data blocks are never executed (CSP3 exempts them from
// script-src) and the hash-based CSP header covers them anyway.
func assertScriptNonces(t *testing.T, rendered []build.RenderedPage) {
	t.Helper()

	for _, page := range rendered {
		for _, match := range scriptTagRe.FindAllStringSubmatch(page.HTML, -1) {
			attrs := strings.ToLower(match[1])

			if strings.Contains(attrs, "src=") || strings.Contains(attrs, "application/ld+json") {
				continue
			}

			nonce := regexp.MustCompile(`nonce="([^"]*)"`).FindStringSubmatch(match[1])
			if len(nonce) < 2 || nonce[1] == "" {
				t.Errorf("%s: inline <script%s> has no nonce", page.Path, match[1])
			}
		}
	}
}

func assertSearchIndex(t *testing.T, outDir string, rendered []build.RenderedPage) {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(outDir, "search-index.json"))
	if err != nil {
		t.Fatalf("read search index: %v", err)
	}

	var docs []build.SearchDoc
	if err := json.Unmarshal(data, &docs); err != nil {
		t.Fatalf("parse search index: %v", err)
	}

	if len(docs) != len(pages.AllDocs()) {
		t.Fatalf("search index has %d docs, want %d", len(docs), len(pages.AllDocs()))
	}

	// Scope invariant (G1, 2026-09-19): search stays DOCS-ONLY. Marketing
	// pages (/, /sales) are persuasion surfaces — matching them for
	// documentation queries would bury the reference content users came
	// for. renderDocs is the only SearchDoc source; this assertion fails
	// loudly if a future change starts indexing other pages.
	docURLs := map[string]bool{}
	for _, ref := range pages.AllDocs() {
		if ref.Slug != "" {
			docURLs["/"+ref.Slug] = true
		}
	}

	for _, doc := range docs {
		if !docURLs[doc.URL] {
			t.Errorf("search index entry %s is outside the docs scope", doc.URL)
		}
	}

	pagePaths := map[string]bool{}
	for _, page := range rendered {
		pagePaths[page.Path] = true
	}

	for _, doc := range docs {
		if !pagePaths[strings.TrimPrefix(doc.URL, "/")+".html"] {
			t.Errorf("search index URL %s has no page", doc.URL)
		}

		if doc.Title == "" || doc.Body == "" {
			t.Errorf("search doc %s missing title/body", doc.URL)
		}
	}
}

func assertSitemap(t *testing.T, outDir string, rendered []build.RenderedPage) {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(outDir, "sitemap.xml"))
	if err != nil {
		t.Fatalf("read sitemap: %v", err)
	}

	for _, page := range rendered {
		if page.Path == "404.html" {
			continue
		}

		url := pages.SiteURL + "/"
		if page.Path != "index.html" {
			url = pages.SiteURL + "/" + strings.TrimSuffix(page.Path, ".html")
		}

		if !strings.Contains(string(data), "<loc>"+url+"</loc>") {
			t.Errorf("sitemap missing %s", url)
		}
	}
}

// externalScriptSrcRe matches a script tag loading from an absolute URL. The
// site's argument is "no framework, no CDN" — every script must be inline
// (nonce-hashed) or root-relative (self-hosted).
var externalScriptSrcRe = regexp.MustCompile(`(?i)<script[^>]+src=["'](?:https?:)?//`)

// scriptBodyRe extracts inline script bodies so framework identifiers are
// only banned in CODE. Prose mentions are legitimate: the landing comparison
// matrix and related-projects page name Alpine.js, React, and Vue as the
// alternatives the library replaces.
var scriptBodyRe = regexp.MustCompile(`(?is)<script\b[^>]*>(.*?)</script>`)

// frameworkMentionRe matches client-framework identifiers as whole words.
// Word boundaries keep prose safe ("revenue" must not trip "vue"); "next"
// and "nuxt" are deliberately absent because docs prose says "next page" —
// their runtimes are blocked by externalScriptSrcRe instead.
var frameworkMentionRe = regexp.MustCompile(`(?i)\b(react|vue|alpine|preact|svelte|angular|jquery)\b`)

// findFrameworkViolations returns every framework/CDN violation in an HTML
// document: an absolute-URL script source, or a framework identifier inside
// an inline script body. Empty slice = clean.
func findFrameworkViolations(html string) []string {
	var violations []string

	for _, match := range externalScriptSrcRe.FindAllString(html, -1) {
		violations = append(violations, "external script src: "+match)
	}

	for _, body := range scriptBodyRe.FindAllStringSubmatch(html, -1) {
		for _, hit := range frameworkMentionRe.FindAllString(body[1], -1) {
			violations = append(violations, "framework identifier in inline script: "+hit)
		}
	}

	return violations
}

// assertNoFrameworkScripts enforces the sales page's core claim — "view
// source and look for the framework" — across every site page. The deployed
// CSP (script-src 'self' + hashes) blocks violations in the browser; this
// test catches them at build time, before the claim ships.
func assertNoFrameworkScripts(t *testing.T, rendered []build.RenderedPage) {
	t.Helper()

	for _, page := range rendered {
		if violations := findFrameworkViolations(page.HTML); len(violations) != 0 {
			t.Errorf("%s: %d framework/CDN violations:\n%s",
				page.Path, len(violations), strings.Join(violations, "\n"))
		}
	}
}

// TestFindFrameworkViolations is the negative-control: the detector must
// flag a fake framework frame and stay quiet on legitimate prose.
func TestFindFrameworkViolations(t *testing.T) {
	t.Parallel()

	frames := map[string]string{
		"cdn script":        `<script src="https://cdn.example.com/react.js"></script>`,
		"scheme-relative":   `<script src="//unpkg.com/vue@3"></script>`,
		"inline identifier": `<script>React.createElement("div")</script>`,
		"script-body hit":   `<script>Alpine.start()</script>`,
	}

	for name, html := range frames {
		if violations := findFrameworkViolations(html); len(violations) == 0 {
			t.Errorf("%s: detector missed frame %q", name, html)
		}
	}

	clean := map[string]string{
		"prose substring":   `<p>Revenue grew and the value persisted.</p>`,
		"next-page link":    `<a href="/guides">Next page</a>`,
		"root-relative":     `<script src="/assets/js/theme-sync.js" defer></script>`,
		"inline nonce":      `<script nonce="abc">var x = 1;</script>`,
		"prose competitor":  `<span class="font-mono">Alpine.js</span>`,
		"jsonld data block": `<script type="application/ld+json">{"name":"templ-components"}</script>`,
	}

	for name, html := range clean {
		if violations := findFrameworkViolations(html); len(violations) != 0 {
			t.Errorf("%s: detector flagged clean HTML %q: %v", name, html, violations)
		}
	}
}
