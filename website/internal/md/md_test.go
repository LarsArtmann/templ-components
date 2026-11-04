package md

import (
	"strings"
	"testing"
)

// TestCodeBlockChrome verifies the fenced-code chrome: every block renders
// the chrome bar with a language label, and a goldmark brace attribute
// ({.title="..."}) surfaces as the title span. Both chrome strings are
// HTML-escaped fence metadata, never raw.
func TestCodeBlockChrome(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		md       string
		contains []string
		omits    []string
	}{
		{
			name: "language label renders",
			md:   "```go\npackage main\n```\n",
			contains: []string{
				`<div class="code-block-chrome">`,
				`<span class="code-block-lang">go</span>`,
				`class="code-copy" aria-label="Copy code to clipboard"`,
			},
			omits: []string{`code-block-title`},
		},
		{
			name: "brace-attribute title renders escaped",
			md:   "```bash {title=\"a<b.go\"}\ngo get example\n```\n",
			contains: []string{
				`<span class="code-block-title">a&lt;b.go</span>`,
				`<span class="code-block-lang">bash</span>`,
			},
		},
		{
			name: "fence without language renders chrome without lang span",
			md:   "```\nplain text\n```\n",
			contains: []string{
				`<div class="code-block-chrome"><span class="code-block-meta">`,
			},
			omits: []string{`code-block-lang`, `code-block-title`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			page, err := Parse(tt.md)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}

			for _, want := range tt.contains {
				if !strings.Contains(page.HTML, want) {
					t.Errorf("rendered HTML missing %q\n---\n%s", want, page.HTML)
				}
			}

			for _, no := range tt.omits {
				if strings.Contains(page.HTML, no) {
					t.Errorf("rendered HTML unexpectedly contains %q\n---\n%s", no, page.HTML)
				}
			}
		})
	}
}
