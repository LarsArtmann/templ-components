package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

// The demo-content guard (#C13): every demoSection on every registry page
// must render non-whitespace CONTENT — not just its heading. A broken
// composition (children passed where the component wants a Content slot, a
// component renamed out from under the demo, a conditional that silently
// renders nothing) otherwise ships as a visible-but-empty section that only
// a human eyeballing the page would catch.

var (
	sectionAnchorRe  = regexp.MustCompile(`<section class="demo-section[^"]*" id="([^"]+)"`)
	sectionHeadingRe = regexp.MustCompile(`^<div class="flex items-center justify-between mb-6">.*?</div>`)
	tagRe            = regexp.MustCompile(`<[^>]+>`)
)

// contentSkipList exempts sections whose content is legitimately textless
// (pure icon/SVG specimens) from the non-whitespace-text assertion.
var contentSkipList = map[string]string{} // filled if icon-only sections surface

func stripTags(s string) string { return strings.TrimSpace(tagRe.ReplaceAllString(s, "")) }

func TestDemoSectionsRenderContent(t *testing.T) {
	srv := httptest.NewServer(newDemoHandler())
	t.Cleanup(srv.Close)

	client := srv.Client()

	for _, page := range demoPages() {
		t.Run(page.Path, func(t *testing.T) {
			resp, err := client.Get(srv.URL + page.Path) //nolint:noctx // test-local server, bounded response
			if err != nil {
				t.Fatalf("fetch %s: %v", page.Path, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("fetch %s: status %d", page.Path, resp.StatusCode)
			}

			html, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("read %s: %v", page.Path, err)
			}

			pageHTML := string(html)

			matches := sectionAnchorRe.FindAllStringSubmatchIndex(pageHTML, -1)
			if len(matches) == 0 {
				t.Fatalf("page %s renders NO demo sections — did the page composition change?", page.Path)
			}

			for i, m := range matches {
				anchor := pageHTML[m[2]:m[3]]

				if reason, skip := contentSkipList[anchor]; skip {
					t.Logf("skip %s: %s", anchor, reason)

					continue
				}

				end := len(pageHTML)
				if i+1 < len(matches) {
					end = matches[i+1][0]
				}

				chunk := pageHTML[m[0]:end]

				// Cut past the section's own open tag, then drop the heading
				// wrapper (one div, no nested div inside) so only CONTENT
				// text remains for the assertion.
				if openEnd := strings.IndexByte(chunk, '>'); openEnd >= 0 {
					chunk = chunk[openEnd+1:]
				}

				chunk = sectionHeadingRe.ReplaceAllString(chunk, "")

				if tail := strings.LastIndex(chunk, "</section>"); tail >= 0 {
					chunk = chunk[:tail]
				}

				if text := stripTags(chunk); text == "" {
					t.Errorf(
						"section #%s on %s renders NO content (heading only) — a demo composition broke",
						anchor, page.Path,
					)
				}
			}
		})
	}
}
