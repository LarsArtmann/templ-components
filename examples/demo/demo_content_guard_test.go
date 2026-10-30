package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

// The demo-content guard (#C13): every demoSection anchor on the registry
// pages must be followed by non-whitespace CONTENT — not just its heading.
// A broken composition (a component renamed out from under the demo, a
// conditional that silently renders nothing) otherwise ships as a
// visible-but-empty section that only a human eyeballing the page would
// catch.
//
// Shipped-structure note (TODO row 485): the demo calls @demoSection(...)
// WITHOUT a children block, so each section renders its heading only and the
// section's content lands as SIBLING markup after </section> — visually
// identical in block flow, which is why this shipped unnoticed. This guard
// therefore asserts content over the span [section start, next section
// start), heading excluded — the contract the pages actually ship.

var (
	sectionAnchorRe  = regexp.MustCompile(`<section class="demo-section[^"]*" id="([^"]+)"`)
	sectionHeadingRe = regexp.MustCompile(`^<div class="flex items-center justify-between mb-6">.*?</div>`)
	tagRe            = regexp.MustCompile(`<[^>]+>`)
)

// contentSkipList exempts sections whose content is legitimately textless
// (pure icon/SVG specimens) from the non-whitespace-text assertion.
var contentSkipList = map[string]string{
	"feedback-spinners":  "spinner specimens are pure SVG — textless by design",
	"feedback-skeletons": "skeleton grid is pure shimmer boxes — textless by design",
}

// noSectionPages are registry routes that intentionally compose no
// demoSection blocks (appshell/card pages); the guard must not demand
// sections from them.
var noSectionPages = map[string]string{
	"/":       "home page — hero + cards, no demoSection blocks",
	"/recipes": "appshell recipe screens, no demoSection blocks",
	"/users":  "appshell table screen, no demoSection blocks",
}

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
				if reason, known := noSectionPages[page.Path]; known {
					t.Logf("no demo sections: %s", reason)

					return
				}

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
				// wrapper (one div, no nested div inside) so only the
				// section's and siblings' CONTENT text remains.
				if openEnd := strings.IndexByte(chunk, '>'); openEnd >= 0 {
					chunk = chunk[openEnd+1:]
				}

				chunk = sectionHeadingRe.ReplaceAllString(chunk, "")

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
