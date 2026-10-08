package utils

import (
	"regexp"
	"testing"
)

// htmlCommentRe matches HTML comments across line breaks. An unterminated
// comment is left in place, so its text stays scannable — a violation-shaped
// token inside broken markup still gets reported (fails safe).
var htmlCommentRe = regexp.MustCompile(`(?s)<!--.*?-->`)

// stripHTMLComments removes single- and multi-line HTML comments so the
// class-shaped compliance regexes judge markup, not prose. Both the
// motion-reduce and RTL guards false-positived on comment prose (2026-10-08),
// which previously forced component comments to be worded around them.
func stripHTMLComments(content string) string {
	return htmlCommentRe.ReplaceAllString(content, "")
}

func TestStripHTMLComments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "no comments unchanged",
			content: `<div class="flex items-center">x</div>`,
			want:    `<div class="flex items-center">x</div>`,
		},
		{
			name:    "single-line comment removed",
			content: `<!-- note about inset-s-0 usage --><div>ok</div>`,
			want:    `<div>ok</div>`,
		},
		{
			name:    "multi-line comment removed",
			content: "a\n<!-- spans\nseveral lines -->\nb",
			want:    "a\n\nb",
		},
		{
			name:    "class-shaped tokens inside comment do not survive",
			content: `<!-- never write transition-colors or end-of-track in prose -->`,
			want:    ``,
		},
		{
			name:    "unterminated comment fails safe",
			content: `<div>ok</div><!-- dangling`,
			want:    `<div>ok</div><!-- dangling`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := stripHTMLComments(tt.content); got != tt.want {
				t.Errorf("stripHTMLComments() = %q, want %q", got, tt.want)
			}
		})
	}
}
