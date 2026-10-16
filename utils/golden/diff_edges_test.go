package golden

import (
	"strings"
	"testing"
)

// TestDiffLcsEdges pins the LCS diff's structural edges (M20/F094): pure
// insertion, pure deletion, mixed multi-hunk changes, identical input, and
// the empty-vs-content extremes. The golden files only exercise the
// "one replaced line" path; these cases are what a contributor hits when a
// golden legitimately grows or shrinks.
func TestDiffLcsEdges(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		got        string
		want       string
		wantHas    []string
		wantHasNot []string
	}{
		{
			name:       "identical inputs produce no diff",
			got:        "a\nb\nc\n",
			want:       "a\nb\nc\n",
			wantHasNot: []string{"---", "+++"},
		},
		{
			name:    "pure insertion marks only added lines",
			got:     "a\nb\nc\n",
			want:    "a\nb\n",
			wantHas: []string{"--- [3] c"}, // got has c, want does not
		},
		{
			name:    "pure deletion marks only removed lines",
			got:     "a\n",
			want:    "a\nx\ny\n",
			wantHas: []string{"+++ [2] x", "+++ [3] y"},
		},
		{
			name:    "multiple hunks are all reported",
			got:     "a\nb\nc\nd\ne\n",
			want:    "a\nX\nc\nd\nY\n",
			wantHas: []string{"--- [2] b", "+++ [2] X", "--- [5] e", "+++ [5] Y"},
		},
		{
			name:    "empty got against content is all additions",
			got:     "",
			want:    "one\ntwo\n",
			wantHas: []string{"+++ [1] one", "+++ [2] two"},
		},
		{
			name:    "content against empty want is all removals",
			got:     "one\ntwo\n",
			want:    "",
			wantHas: []string{"--- [1] one", "--- [2] two"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := diff(tt.got, tt.want)

			for _, needle := range tt.wantHas {
				if !strings.Contains(got, needle) {
					t.Errorf("diff must contain %q:\n%s", needle, got)
				}
			}

			for _, needle := range tt.wantHasNot {
				if strings.Contains(got, needle) {
					t.Errorf("diff must NOT contain %q:\n%s", needle, got)
				}
			}
		})
	}
}

// TestDiffLongLineReadability pins the #339 triage form: a replacement pair
// whose lines exceed the threshold is windowed around the first difference
// and followed by a word diff that isolates the changed tokens — instead of
// dumping two multi-kilobyte lines.
func TestDiffLongLineReadability(t *testing.T) {
	t.Parallel()

	pad := strings.Repeat("shared-token ", 40)
	wantLine := "<div class=\"" + pad + "old-class\" id=\"x\">"
	gotLine := "<div class=\"" + pad + "new-class\" id=\"x\">"

	got := diff(gotLine+"\n", wantLine+"\n")

	for _, needle := range []string{
		"~~~ word diff:",
		"-old-class ",
		"+new-class ",
	} {
		if !strings.Contains(got, needle) {
			t.Errorf("long-line diff must contain %q:\n%s", needle, got)
		}
	}

	for _, banned := range []string{"(+", "shared-token"} {
		if strings.Contains(got, banned) {
			t.Errorf("long-line diff must be windowed (found %q):\n%s", banned, got)
		}
	}
}

// TestDiffShortLinesStayVerbatim pins the historical format: short
// replacement pairs print the full lines with no word-diff section.
func TestDiffShortLinesStayVerbatim(t *testing.T) {
	t.Parallel()

	got := diff("<p>b</p>\n", "<p>a</p>\n")

	for _, needle := range []string{"--- [1] <p>a</p>", "+++ [1] <p>b</p>"} {
		if !strings.Contains(got, needle) {
			t.Errorf("short-line diff must contain %q:\n%s", needle, got)
		}
	}

	if strings.Contains(got, "word diff") {
		t.Errorf("short-line diff must not carry a word-diff section:\n%s", got)
	}
}
