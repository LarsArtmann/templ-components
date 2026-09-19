// Package build contains the static-site generator plumbing for the
// templ-components website: page rendering, build-scoped CSP nonces, and
// library-fact statistics derived from the repository checkout.
package build

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/a-h/templ"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/larsartmann/templ-components/utils"
)

// Nonce returns a fresh CSP nonce (32 hex chars) for one site build. Every
// inline script across every generated page shares it; firebase.json pins the
// same nonce via its header config at deploy time (regenerated each build).
func Nonce() (string, error) {
	const nonceBytes = 16

	buf := make([]byte, nonceBytes) //nolint:makezero // fixed-size random buffer for crypto/rand
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate CSP nonce: %w", err)
	}

	return hex.EncodeToString(buf), nil
}

// Page is one rendered output file of the site.
type Page struct {
	// Path is the dist-relative file path, e.g. "index.html" or
	// "guides/theming.html" (Firebase cleanUrls serves it as /guides/theming).
	Path string
	// Component renders the full document.
	Component templ.Component
}

// Renderer renders templ components and writes pages to the dist directory.
type Renderer struct {
	Nonce string
}

// NewRenderer creates a Renderer bound to a build-scoped nonce.
func NewRenderer(nonce string) *Renderer {
	return &Renderer{Nonce: nonce}
}

// RenderToString renders a component to an HTML string.
func (r *Renderer) RenderToString(ctx context.Context, c templ.Component) (string, error) {
	var sb strings.Builder
	if err := c.Render(ctx, &sb); err != nil {
		return "", fmt.Errorf("render component: %w", err)
	}

	return sb.String(), nil
}

// RenderPages renders every page into memory, preserving input order.
func (r *Renderer) RenderPages(ctx context.Context, pages []Page) ([]RenderedPage, error) {
	rendered := make([]RenderedPage, 0, len(pages))

	for _, page := range pages {
		pageHTML, err := r.RenderToString(ctx, page.Component)
		if err != nil {
			return nil, fmt.Errorf("render %s: %w", page.Path, err)
		}

		rendered = append(rendered, RenderedPage{Path: page.Path, HTML: pageHTML})
	}

	return rendered, nil
}

// WriteRendered writes rendered pages under outDir, creating parent
// directories as needed.
func WriteRendered(outDir string, pages []RenderedPage) error {
	for _, page := range pages {
		target := filepath.Join(outDir, filepath.FromSlash(page.Path))
		//nolint:gosec // public site asset directory and file must be world-readable
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("create dir for %s: %w", page.Path, err)
		}

		//nolint:gosec // public site asset must be world-readable
		if err := os.WriteFile(target, []byte(page.HTML), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", page.Path, err)
		}
	}

	return nil
}

var (
	templComponentRe = regexp.MustCompile(`(?m)^templ ([A-Z][A-Za-z0-9]*)\(`)
	iconNameRe       = regexp.MustCompile(`Name\s*=\s*("[a-z0-9-]+")`)
	isValidRe        = regexp.MustCompile(`(?m)^func ([A-Za-z0-9]+)IsValid\(`)
)

// Stats are the library facts shown on the site (hero metric strip). Derived
// from the actual checkout at build time so the site cannot drift from
// reality the way the hand-typed Astro numbers did.
type Stats struct {
	Components     int    // templ components across the published library packages
	Icons          int    // distinct SVG icons in the icons module
	Enums          int    // closed-set enums shipping an IsValid() method
	Modules        int    // published Go modules in the library workspace
	GoVersion      string // root go.mod go directive, patch-trimmed for display ("1.26")
	LibraryVersion string // utils.Version verbatim ("1.18.0") — the sold library's own semver
}

// statsDirs are the published library packages scanned for templ components.
//
//nolint:gochecknoglobals // static scan configuration table
var statsDirs = []string{
	"display", "feedback", "forms", "layout", "navigation", "recipes",
	"htmx", "datastar", "errorpage",
}

// excludedModules are repo-local Go modules that are not part of the
// published library (site generator + visual regression harness).
//
//nolint:gochecknoglobals // static scan configuration table
var excludedModules = map[string]bool{
	"website":    true,
	"visualtest": true,
}

// CountStats derives the library facts shown on the site from the repo
// checkout rooted at repoRoot: templ components, icons, IsValid enums, and
// published modules. Fail-soft per metric: a missed file counts nothing.
func CountStats(repoRoot string) (Stats, error) {
	var stats Stats

	for _, dir := range statsDirs {
		root := filepath.Join(repoRoot, filepath.FromSlash(dir))
		if _, err := os.Stat(root); err != nil {
			continue
		}

		count, err := countMatches(root, ".templ", templComponentRe)
		if err != nil {
			return stats, fmt.Errorf("count components in %s: %w", dir, err)
		}

		stats.Components += count
	}

	stats.Icons = countUniqueIcons(filepath.Join(repoRoot, "icons", "icon_names.go"))

	enumCount, err := countMatches(repoRoot, ".go", isValidRe)
	if err != nil {
		return stats, fmt.Errorf("count enums: %w", err)
	}

	stats.Enums = enumCount

	stats.Modules, err = countModules(repoRoot)
	if err != nil {
		return stats, fmt.Errorf("count modules: %w", err)
	}

	stats.GoVersion = goDisplayVersion(filepath.Join(repoRoot, "go.mod"))
	stats.LibraryVersion = utils.Version

	return stats, nil
}

// goDisplayVersion reads the root go.mod's go directive ("go 1.26.7") and
// trims the patch segment for display ("1.26"). Empty string on any read or
// parse failure — fail-soft per metric, like every other stat.
func goDisplayVersion(goModPath string) string {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return ""
	}

	for line := range strings.SplitSeq(string(data), "\n") {
		directive, found := strings.CutPrefix(strings.TrimSpace(line), "go ")
		if !found {
			continue
		}

		parts := strings.Split(strings.TrimSpace(directive), ".")
		if len(parts) < 2 {
			return ""
		}

		return parts[0] + "." + parts[1]
	}

	return ""
}

// countUniqueIcons counts distinct icon name constants (aliases like Close/X
// point at the same name value, so the raw constant count overstates).
func countUniqueIcons(iconNamesFile string) int {
	data, err := os.ReadFile(iconNamesFile)
	if err != nil {
		return 0
	}

	seen := map[string]bool{}
	for _, match := range iconNameRe.FindAllStringSubmatch(string(data), -1) {
		seen[match[1]] = true
	}

	return len(seen)
}

func countMatches(root, suffix string, regex *regexp.Regexp) (int, error) {
	count := 0

	err := filepath.WalkDir(
		root,
		func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if entry.IsDir() {
				switch entry.Name() {
				case "node_modules", "dist", ".git", "testdata", "static":
					return filepath.SkipDir
				}

				return nil
			}

			if !strings.HasSuffix(path, suffix) {
				return nil
			}
			// IsValid lives in production code only; tests never define it.
			if strings.HasSuffix(path, "_test.go") {
				return nil
			}

			data, readErr := os.ReadFile(path) //nolint:gosec // trusted local repository scan
			if readErr != nil {
				return fmt.Errorf("read %s: %w", path, readErr)
			}

			count += len(regex.FindAllSubmatch(data, -1))

			return nil
		},
	)
	if err != nil {
		return count, fmt.Errorf("walk %s: %w", root, err)
	}

	return count, nil
}

func countModules(repoRoot string) (int, error) {
	count := 0

	err := filepath.WalkDir(
		repoRoot,
		func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if entry.IsDir() {
				switch entry.Name() {
				case "node_modules", "dist", ".git", "testdata":
					return filepath.SkipDir
				}

				rel, relErr := filepath.Rel(repoRoot, path)
				if relErr == nil && excludedModules[filepath.ToSlash(rel)] {
					return filepath.SkipDir
				}

				return nil
			}

			if entry.Name() == "go.mod" {
				count++
			}

			return nil
		},
	)
	if err != nil {
		return count, fmt.Errorf("walk %s: %w", repoRoot, err)
	}

	return count, nil
}

// ChromaCSS returns the class-based syntax-highlighting stylesheet: the
// github-light theme scoped under `html .chroma`, the github-dark theme
// under `html.dark .chroma`. Regenerated every site build so highlighting
// can never drift from the chroma dependency.
//
// Light theme is "github" (chroma v2.27's "github-light" ships inverted
// colors), and PreWrapper backgrounds are overridden to transparent — the
// site stylesheet owns the code-block surface via --color-bg-code.
func ChromaCSS() (string, error) {
	lightBuf, darkBuf := &bytes.Buffer{}, &bytes.Buffer{}
	light := chromahtml.New(chromahtml.WithClasses(true))
	dark := chromahtml.New(chromahtml.WithClasses(true))

	if err := light.WriteCSS(lightBuf, styles.Get("github")); err != nil {
		return "", fmt.Errorf("chroma light css: %w", err)
	}

	if err := dark.WriteCSS(darkBuf, styles.Get("github-dark")); err != nil {
		return "", fmt.Errorf("chroma dark css: %w", err)
	}

	lightScoped := strings.ReplaceAll(lightBuf.String(), ".chroma", "html .chroma")
	darkScoped := strings.ReplaceAll(darkBuf.String(), ".chroma", "html.dark .chroma")

	transparent := "\n/* Site-owned code surface: theme tokens pick the background. */\n" +
		"html .chroma { background-color: transparent; color: inherit; }\n" +
		"html.dark .chroma { background-color: transparent; color: inherit; }\n"

	return lightScoped + "\n" + darkScoped + transparent, nil
}

// CopyTree copies every regular file under src into dst, preserving the
// relative structure (used for assets/ and public/ → dist).
func CopyTree(src, dst string) error {
	err := filepath.WalkDir(
		src,
		func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if entry.IsDir() {
				return nil
			}

			rel, err := filepath.Rel(src, path)
			if err != nil {
				return fmt.Errorf("relativize %s: %w", path, err)
			}

			target := filepath.Join(dst, rel)
			//nolint:gosec // public site asset directory must be traversable
			if err := os.MkdirAll(
				filepath.Dir(target),
				0o755,
			); err != nil {
				return fmt.Errorf("create dir %s: %w", filepath.Dir(target), err)
			}

			data, err := os.ReadFile(path) //nolint:gosec // trusted local asset tree
			if err != nil {
				return fmt.Errorf("read %s: %w", path, err)
			}

			//nolint:gosec // public site asset must be world-readable
			if err := os.WriteFile(target, data, 0o644); err != nil {
				return fmt.Errorf("write %s: %w", target, err)
			}

			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("copy tree %s: %w", src, err)
	}

	return nil
}
