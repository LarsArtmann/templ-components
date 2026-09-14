package build

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var (
	linkAnchorRe = regexp.MustCompile(`(?i)<a\b[^>]*\bhref="([^"]*)"`)
	linkAssetRe  = regexp.MustCompile(`(?i)<(?:script|img|link|source)\b[^>]*\b(?:src|href)="([^"]*)"`)
	elementIDRe  = regexp.MustCompile(`(?i)\bid="([^"]*)"`)
)

// linkIndex is the lookup state for validating one site's references.
type linkIndex struct {
	pagePaths map[string]bool
	assets    map[string]bool
	pageIDs   map[string]map[string]bool
}

func newLinkIndex(pages []RenderedPage, assetFiles []string) *linkIndex {
	index := &linkIndex{
		pagePaths: make(map[string]bool, len(pages)),
		assets:    make(map[string]bool, len(assetFiles)),
		pageIDs:   make(map[string]map[string]bool, len(pages)),
	}

	for _, file := range assetFiles {
		index.assets[file] = true
	}

	for _, page := range pages {
		index.pagePaths[page.Path] = true

		ids := map[string]bool{}
		for _, match := range elementIDRe.FindAllStringSubmatch(page.HTML, -1) {
			ids[match[1]] = true
		}

		index.pageIDs[page.Path] = ids
	}

	return index
}

// CheckLinks validates every internal reference across the rendered pages:
// anchors must exist in their target page, page references must resolve to a
// rendered page, and asset references must resolve to a dist file. External
// URLs are skipped. It returns one problem per line (empty means clean), so
// a site build can fail on dead links before deployment.
func CheckLinks(pages []RenderedPage, assetFiles []string) []string {
	index := newLinkIndex(pages, assetFiles)

	var problems []string

	for _, page := range pages {
		for _, ref := range pageReferences(page.HTML) {
			if problem := index.checkRef(page.Path, ref); problem != "" {
				problems = append(problems, problem)
			}
		}
	}

	sort.Strings(problems)

	return problems
}

// pageReferences extracts the unique href/src references of one document.
func pageReferences(html string) []string {
	refs := map[string]bool{}

	for _, re := range []*regexp.Regexp{linkAnchorRe, linkAssetRe} {
		for _, match := range re.FindAllStringSubmatch(html, -1) {
			refs[match[1]] = true
		}
	}

	return sortedKeys(refs)
}

// checkRef validates one reference and returns a problem description (empty
// when the reference resolves).
func (index *linkIndex) checkRef(pagePath, ref string) string {
	if externalReference(ref) {
		return ""
	}

	if !strings.HasPrefix(ref, "/") && !strings.HasPrefix(ref, "#") {
		return fmt.Sprintf("%s: internal references must be absolute, got %q", pagePath, ref)
	}

	path, fragment, _ := strings.Cut(ref, "#")

	if path == "" {
		if !index.pageIDs[pagePath][fragment] {
			return fmt.Sprintf("%s: broken same-page anchor #%s", pagePath, fragment)
		}

		return ""
	}

	target := index.resolveTarget(path)
	if !index.pagePaths[target] && !index.assets[target] {
		return fmt.Sprintf("%s: broken link %q (no page or asset %s)", pagePath, ref, target)
	}

	if fragment != "" && index.pagePaths[target] && !index.pageIDs[target][fragment] {
		return fmt.Sprintf("%s: broken anchor %q (no id %q in %s)", pagePath, ref, fragment, target)
	}

	return ""
}

// externalReference reports whether a reference points off-site (or is a
// scheme the browser handles without a fetch).
func externalReference(ref string) bool {
	return ref == "" ||
		strings.HasPrefix(ref, "http://") ||
		strings.HasPrefix(ref, "https://") ||
		strings.HasPrefix(ref, "//") ||
		strings.HasPrefix(ref, "mailto:") ||
		strings.HasPrefix(ref, "tel:")
}

// resolveTarget maps a URL path to a dist file: clean URLs (no extension)
// resolve to "<path>.html", asset paths pass through when they exist as-is.
func (index *linkIndex) resolveTarget(path string) string {
	trimmed := strings.TrimPrefix(path, "/")
	if trimmed == "" {
		return "index.html"
	}

	if strings.HasSuffix(trimmed, ".html") || index.assets[trimmed] {
		return trimmed
	}

	return trimmed + ".html"
}

func sortedKeys(set map[string]bool) []string {
	keys := make([]string, 0, len(set))
	for key := range set {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}
