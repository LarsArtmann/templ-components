package main

import (
	"encoding/json"
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

	wantPages := 2 + len(pages.AllDocs())
	if len(rendered) != wantPages {
		t.Fatalf("wrote %d HTML pages, want %d", len(rendered), wantPages)
	}

	if problems := build.CheckLinks(rendered, assets); len(problems) != 0 {
		t.Errorf("%d broken links:\n%s", len(problems), strings.Join(problems, "\n"))
	}

	assertScriptNonces(t, rendered)
	assertSearchIndex(t, outDir, rendered)
	assertSitemap(t, outDir, rendered)
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

		data, err := os.ReadFile(path) //nolint:gosec // trusted test-only dist read
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
