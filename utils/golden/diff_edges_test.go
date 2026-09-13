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
