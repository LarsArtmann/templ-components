// Package md converts the docs markdown (ported from the Astro MDX sources)
// to HTML: YAML frontmatter, GFM tables, auto heading IDs for the TOC, and
// chroma class-based syntax highlighting (dual-theme CSS comes from the site
// build).
package md

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
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

var frontmatterRe = regexp.MustCompile(`(?s)\A---\n(.*?)\n---\n?`)
var frontmatterKeyRe = regexp.MustCompile(`(?m)^([a-zA-Z_]+):\s*(.*)$`)

// templLexer highlights templ fences as Go (chroma has no dedicated templ
// lexer; Go coloring covers keywords, strings, and struct literals well).
type templLexer struct {
	chroma.Lexer
}

func (l templLexer) Tokenise(opts *chroma.TokeniseConfig, source string) (chroma.Iterator, error) {
	return lexers.Get("go").Tokenise(opts, source)
}

func newMarkdown() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			extension.Linkify,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github-dark"),
				highlighting.WithCustomLexer(templLexer{Lexer: lexers.Get("go")}),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
				),
				highlighting.WithWrapperRenderer(func(w io.Writer, ctx highlighting.CodeBlockContext, entering bool) {
					if entering {
						fmt.Fprint(w, `<div class="code-block">`)
						fmt.Fprint(w, `<button type="button" class="code-copy" aria-label="Copy code to clipboard">Copy</button>`)
						return
					}

					fmt.Fprint(w, `</div>`)
				}),
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
	page := Page{}

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
	var headings []Heading

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
		if id, ok := heading.AttributeString("id"); ok {
			if raw, isBytes := id.([]byte); isBytes {
				idText = string(raw)
			}
		}

		var textParts []string
		for child := heading.FirstChild(); child != nil; child = child.NextSibling() {
			if codeSpan, isCode := child.(*ast.CodeSpan); isCode {
				var codeText strings.Builder
				for segment := range codeSpan.Segments {
					codeText.Write(segment.Value(source))
				}

				textParts = append(textParts, codeText.String())
				continue
			}

			if textNode, isText := child.(*ast.Text); isText {
				textParts = append(textParts, string(textNode.Value(source)))
			}
		}

		headings = append(headings, Heading{ID: idText, Text: strings.Join(textParts, ""), Level: heading.Level})

		return ast.WalkContinue, nil
	})
	if err != nil {
		return headings
	}

	return headings
}
