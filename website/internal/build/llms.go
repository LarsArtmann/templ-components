package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// WriteLLMS emits dist/llms.txt — the lean machine-readable site index
// (llmstxt.org shape: one H1, a blockquote summary, one link per line under
// section headings). AI agents and crawlers fetch it to orient before
// reading pages; the full corpus variant (llms-full.txt) stays a separate
// future surface. Docs entries reuse the search index's parsed metadata so
// the two indexes can never disagree about titles or URLs.
func WriteLLMS(outDir string, stats Stats, docs []SearchDoc) error {
	content := LLMSIndex(stats, docs)

	target := filepath.Join(outDir, "llms.txt")
	//nolint:gosec // public site asset must be world-readable
	if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write llms.txt: %w", err)
	}

	return nil
}

// LLMSIndex renders the llms.txt document. Deterministic: fixed section
// order, docs in their registered reading order, no timestamps.
func LLMSIndex(stats Stats, docs []SearchDoc) string {
	var sb strings.Builder

	sb.WriteString("# templ-components\n\n")
	sb.WriteString("> Server-rendered Go UI components for templ: ")
	sb.WriteString(strconv.Itoa(stats.Components))
	sb.WriteString(" components, ")
	sb.WriteString(strconv.Itoa(stats.Icons))
	sb.WriteString(" SVG icons, ")
	sb.WriteString(strconv.Itoa(stats.Enums))
	sb.WriteString(" typed enums, ")
	sb.WriteString(strconv.Itoa(stats.Modules))
	sb.WriteString(" opt-in Go modules. HTMX-native with opt-in Datastar,")
	sb.WriteString(" Tailwind CSS v4, CSP-safe by construction, dark mode and RTL")
	sb.WriteString(" built in. MIT licensed, no client framework required.\n\n")
	sb.WriteString("Typed props make invalid states compile errors; every component")
	sb.WriteString(" ships HTML golden tests, an axe-core accessibility sweep, and")
	sb.WriteString(" motion-reduce fallbacks. Version ")
	sb.WriteString(stats.LibraryVersion)
	sb.WriteString(", requires Go ")
	sb.WriteString(stats.GoVersion)
	sb.WriteString(".\n\n")

	sb.WriteString("## Product\n\n")
	sb.WriteString("- [Home](" + siteURLPath("/") + "): the library, live-rendered\n")
	sb.WriteString("- [Sales pitch](" + siteURLPath("/sales") + "): every adoption argument on one page\n\n")

	if len(docs) > 0 {
		sb.WriteString("## Documentation\n\n")

		for _, doc := range docs {
			sb.WriteString("- [" + llmsLine(doc.Title) + "](" + siteURLPath(doc.URL) + ")")

			if description := llmsLine(doc.Description); description != "" {
				sb.WriteString(": " + description)
			}

			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// siteURLPath joins the canonical site origin with a root-relative path.
func siteURLPath(path string) string {
	return "https://templcomponents.lars.software" + path
}

// llmsLine flattens a title or description to one spec-friendly line.
func llmsLine(text string) string {
	return strings.Join(strings.Fields(text), " ")
}

// WriteLLMSFull emits dist/llms-full.txt — the llmstxt.org "full corpus"
// companion to llms.txt: the same H1/blockquote header followed by every
// docs page's full plain-text body under a heading carrying its canonical
// URL. Deterministic: docs in their registered reading order, no timestamps.
func WriteLLMSFull(outDir string, stats Stats, docs []SearchDoc) error {
	content := LLMSFullIndex(stats, docs)

	target := filepath.Join(outDir, "llms-full.txt")
	//nolint:gosec // public site asset must be world-readable
	if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write llms-full.txt: %w", err)
	}

	return nil
}

// LLMSFullIndex renders the llms-full.txt document from the same SearchDoc
// metadata as the search index and llms.txt — one source, no drift.
func LLMSFullIndex(stats Stats, docs []SearchDoc) string {
	var sb strings.Builder

	sb.WriteString(LLMSIndex(stats, docs))

	if len(docs) > 0 {
		sb.WriteString("\n## Full documentation\n\n")

		for _, doc := range docs {
			sb.WriteString("\n---\n\n# " + llmsLine(doc.Title) +
				" (" + siteURLPath(doc.URL) + ")\n\n")
			sb.WriteString(strings.TrimSpace(doc.Body))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
