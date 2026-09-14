package build

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"regexp"
	"strings"
)

// SearchSection is one heading anchor inside an indexed document.
type SearchSection struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// SearchDoc is one document in the client-side search index
// (dist/search-index.json, fetched lazily by /assets/js/search.js).
type SearchDoc struct {
	URL         string          `json:"url"`
	Title       string          `json:"title"`
	Description string          `json:"description,omitempty"`
	Sections    []SearchSection `json:"sections,omitempty"`
	Body        string          `json:"body"`
}

var (
	searchScriptRe = regexp.MustCompile(`(?is)<script\b[^>]*>.*?</script>`)
	searchStyleRe  = regexp.MustCompile(`(?is)<style\b[^>]*>.*?</style>`)
	searchBlockRe  = regexp.MustCompile(
		`(?i)</?(p|div|li|tr|h[1-6]|pre|blockquote|section|article|table|ul|ol|figure|figcaption|dt|dd)\b[^>]*>`,
	)
	searchBrRe  = regexp.MustCompile(`(?i)<br\s*/?>`)
	searchTagRe = regexp.MustCompile(`<[^>]*>`)
)

// PlainText reduces rendered HTML to searchable text: script/style blocks are
// dropped, block boundaries become spaces, tags are stripped, entities are
// decoded, and whitespace is collapsed. Code fences survive as text on
// purpose — searching code snippets is a feature.
func PlainText(markup string) string {
	stripped := searchScriptRe.ReplaceAllString(markup, " ")
	stripped = searchStyleRe.ReplaceAllString(stripped, " ")
	stripped = searchBrRe.ReplaceAllString(stripped, " ")
	stripped = searchBlockRe.ReplaceAllString(stripped, " ")
	stripped = searchTagRe.ReplaceAllString(stripped, " ")
	stripped = html.UnescapeString(stripped)

	return strings.Join(strings.Fields(stripped), " ")
}

// WriteSearchIndex writes the client-side search index as compact JSON.
func WriteSearchIndex(outDir string, docs []SearchDoc) error {
	data, err := json.Marshal(docs)
	if err != nil {
		return fmt.Errorf("encode search index: %w", err)
	}

	target := outDir + "/search-index.json"
	//nolint:gosec // public site asset must be world-readable
	if err := os.WriteFile(target, data, 0o644); err != nil {
		return fmt.Errorf("write search index: %w", err)
	}

	return nil
}
