package golden

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGoldenUpdateFlag verifies the -update flag writes golden files.
//
//nolint:paralleltest // NOT parallel — modifies the global `update` flag
func TestGoldenUpdateFlag(t *testing.T) {
	// NOT parallel — modifies the global `update` flag
	dir := t.TempDir()

	origUpdate := update
	defer func() { update = origUpdate }()

	update = new(true)

	content := "<div class=\"b a\">Hello</div>"
	assertInDir(t, dir, "test_update", content)

	goldenPath := filepath.Join(dir, "test_update.golden")

	data, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("golden file not written: %v", err)
	}

	if !strings.Contains(string(data), "a b") {
		t.Errorf("golden file should have sorted classes, got: %s", string(data))
	}
}

// TestGoldenUpdateFlagMkdir verifies -update creates directories as needed.
//
//nolint:paralleltest // NOT parallel — modifies the global `update` flag
func TestGoldenUpdateFlagMkdir(t *testing.T) {
	// NOT parallel — modifies the global `update` flag
	dir := t.TempDir()

	origUpdate := update
	defer func() { update = origUpdate }()

	update = new(true)

	nestedDir := filepath.Join(dir, "nested", "deep")
	content := "<p>test</p>"
	assertInDir(t, nestedDir, "nested_test", content)

	goldenPath := filepath.Join(nestedDir, "nested_test.golden")
	if _, err := os.Stat(goldenPath); err != nil {
		t.Errorf("golden file should exist at %s: %v", goldenPath, err)
	}
}

// TestNormalizeClassesEdgeCases verifies CSS class normalization edge cases.
func TestNormalizeClassesEdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		html string
		want string
	}{
		{
			name: "no classes",
			html: `<p>no classes here</p>`,
			want: `<p>no classes here</p>`,
		},
		{
			name: "multiple class attrs",
			html: `<span class="c b a"><em class="z y x">nested</em></span>`,
			want: `<span class="a b c"><em class="x y z">nested</em></span>`,
		},
		{
			name: "empty class",
			html: `<div class="">empty</div>`,
			want: `<div class="">empty</div>`,
		},
		{
			name: "single class",
			html: `<div class="solo">text</div>`,
			want: `<div class="solo">text</div>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := normalizeClasses(tt.html)
			if got != tt.want {
				t.Errorf("normalizeClasses:\n got: %s\nwant: %s", got, tt.want)
			}
		})
	}
}

// TestGoldenDiffIdentical verifies diff returns empty for identical strings.
func TestGoldenDiffIdentical(t *testing.T) {
	t.Parallel()

	got := diff("same\n", "same\n")
	if got != "" {
		t.Errorf("diff of identical strings should be empty, got: %s", got)
	}
}

// TestGoldenDiffMultiLine verifies diff output for multi-line differences.
func TestGoldenDiffMultiLine(t *testing.T) {
	t.Parallel()

	got := diff("a\nb\nc\n", "a\nx\nc\n")
	if !strings.Contains(got, "b") || !strings.Contains(got, "x") {
		t.Errorf("diff should show changed lines, got: %s", got)
	}
}

// TestGoldenDiffLCSAlignment verifies the LCS-based diff correctly handles
// insertions without cascading changes to subsequent lines.
func TestGoldenDiffLCSAlignment(t *testing.T) {
	t.Parallel()

	// Inserting a line in the middle — only the insertion should show, not every line after.
	got := diff("a\nc\n", "a\nb\nc\n")
	if !strings.Contains(got, "+++ [2] b") {
		t.Errorf("diff should show inserted line at position 2, got:\n%s", got)
	}

	if strings.Contains(got, "---") {
		t.Errorf("diff should have no deletions for pure insertion, got:\n%s", got)
	}
}

// TestNormalizeIDs verifies auto-generated EnsureID values are replaced with stable placeholders.
func TestNormalizeIDs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "primary hex format",
			input: `id="tc-modal-a1b2c3d4e5f6a7b8"`,
			want:  `id="tc-modal-NORMALIZED"`,
		},
		{
			name:  "fallback timestamp format",
			input: `id="tc-modal-1753897654321234567-1"`,
			want:  `id="tc-modal-NORMALIZED"`,
		},
		{
			name:  "cross-reference in popovertarget",
			input: `popovertarget="tc-popover-abcdef0123456789"`,
			want:  `popovertarget="tc-popover-NORMALIZED"`,
		},
		{
			name:  "explicit ID preserved",
			input: `id="my-custom-id"`,
			want:  `id="my-custom-id"`,
		},
		{
			name:  "short hex not matched (too short)",
			input: `id="tc-modal-abc"`,
			want:  `id="tc-modal-abc"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := normalizeIDs(tt.input)
			if got != tt.want {
				t.Errorf("normalizeIDs(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
