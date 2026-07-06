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
	data, err := os.ReadFile(goldenPath) //nolint:gosec // test temp dir
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

// TestLineAtOutOfRange verifies lineAt returns "" for out-of-range index.
func TestLineAtOutOfRange(t *testing.T) {
	t.Parallel()

	if lineAt([]string{"a\n", "b\n"}, 5) != "" {
		t.Error("lineAt(5) should return '' for out-of-range")
	}
}

//go:fix inline
func boolPtr(b bool) *bool { return new(b) }
