// Package md converts the docs markdown (ported from the Astro MDX sources)
// to HTML: YAML frontmatter, GFM tables, auto heading IDs for the TOC, and
// chroma class-based syntax highlighting (dual-theme CSS comes from the site
// build).
package md

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Heading is one TOC entry (h2/h3) collected from the document.
type Heading struct {
	ID    string
	Text  string
	Level int
}

// Page is one parsed docs page.
type Page struct {
	Title       string
	Description string
	HTML        string
	Headings    []Heading
}

var (
	frontmatterRe    = regexp.MustCompile(`(?s)\A---\n(.*?)\n---\n?`)
	frontmatterKeyRe = regexp.MustCompile(`(?m)^([a-zA-Z_]+):\s*(.*)$`)
)

// templLexer highlights templ fences as Go (chroma has no dedicated templ
// lexer; Go coloring covers keywords, strings, and struct literals well).
type templLexer struct {
	registry *chroma.LexerRegistry
}

func (l templLexer) Config() *chroma.Config {
	//nolint:exhaustruct_v5 // chroma's optional lexer-config fields are fine at zero value
	return &chroma.Config{Name: "templ", Aliases: []string{"templ"}, Filenames: []string{"*.templ"}}
}

func (l templLexer) SetRegistry(registry *chroma.LexerRegistry) chroma.Lexer {
	l.registry = registry

	return l
}

func (l templLexer) Tokenise(options *chroma.TokeniseOptions, source string) (chroma.Iterator, error) {
	goLexer := lexers.Get("go")
	if l.registry != nil {
		if found := l.registry.Get("go"); found != nil {
			goLexer = found
		}
	}

	return goLexer.Tokenise(options, source)
}

func (templLexer) AnalyseText(string) float32 {
	return 0
}

func (l templLexer) SetAnalyser(func(string) float32) chroma.Lexer {
	return l
}

func init() {
	lexers.Register(templLexer{})
}

func newMarkdown() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			extension.Linkify,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github-dark"),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
				),
				highlighting.WithWrapperRenderer(
					func(w util.BufWriter, ctx highlighting.CodeBlockContext, entering bool) {
						if entering {
							fmt.Fprint(w, `<div class="code-block">`)
							fmt.Fprint(
								w,
								`<button type="button" class="code-copy" aria-label="Copy code to clipboard">Copy</button>`,
							)
							return
						}

						fmt.Fprint(w, `</div>`)
					},
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
}

// Parse converts one markdown document (with optional frontmatter) to a
// rendered Page.
func Parse(source string) (Page, error) {
	page := Page{} //nolint:exhaustruct_v5 // fields are populated incrementally below

	remainder := source
	if match := frontmatterRe.FindStringSubmatch(source); match != nil {
		for _, keyMatch := range frontmatterKeyRe.FindAllStringSubmatch(match[1], -1) {
			value := strings.Trim(strings.TrimSpace(keyMatch[2]), `"`)
			switch strings.ToLower(keyMatch[1]) {
			case "title":
				page.Title = value
			case "description":
				page.Description = value
			}
		}

		remainder = source[len(match[0]):]
	}

	var buf bytes.Buffer
	if err := newMarkdown().Convert([]byte(remainder), &buf); err != nil {
		return page, fmt.Errorf("convert markdown: %w", err)
	}

	page.HTML = buf.String()
	page.Headings = collectHeadings([]byte(remainder))

	return page, nil
}

func collectHeadings(source []byte) []Heading {
	headings := []Heading{}

	parsed := newMarkdown().Parser().Parse(text.NewReader(source))

	err := ast.Walk(parsed, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		heading, ok := node.(*ast.Heading)
		if !ok || heading.Level < 2 || heading.Level > 3 {
			return ast.WalkContinue, nil
		}

		idText := ""
		if id, idOK := heading.AttributeString("id"); idOK {
			if raw, isBytes := id.([]byte); isBytes {
				idText = string(raw)
			}
		}

		var textParts []string

		segments := heading.Lines()

		for i := range segments.Len() {
			segment := segments.At(i)

			textParts = append(textParts, string(segment.Value(source)))
		}

		headings = append(headings, Heading{ID: idText, Text: strings.Join(textParts, ""), Level: heading.Level})

		return ast.WalkContinue, nil
	})
	if err != nil {
		return headings
	}

	return headings
}
