package golden

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAssertMatchesGoldenFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	content := "<div class=\"bg-blue-600\">Hello</div>\n"
	name := "test_assert"
	goldenPath := filepath.Join(dir, name+".golden")

	if err := os.WriteFile(goldenPath, []byte(content), 0o600); err != nil {
		t.Fatalf("setup write: %v", err)
	}

	assertInDir(t, dir, name, content)
}

func TestAssertRejectsMismatch(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	name := "test_mismatch"
	goldenPath := filepath.Join(dir, name+".golden")

	if err := os.WriteFile(goldenPath, []byte("old content\n"), 0o600); err != nil {
		t.Fatalf("setup write: %v", err)
	}

	mockT := &testing.T{}
	assertInDir(mockT, dir, name, "new content\n")

	if !mockT.Failed() {
		t.Error("expected Assert to fail on mismatch")
	}
}

func TestDiffOutput(t *testing.T) {
	t.Parallel()

	got := diff("line1\nline2\n", "line1\nchanged\n")
	if got == "" {
		t.Error("expected non-empty diff for different inputs")
	}

	if !strings.Contains(got, "--- [2] line2") || !strings.Contains(got, "+++ [2] changed") {
		t.Errorf("unexpected diff output:\n%s", got)
	}
}

func TestNormalizeClasses(t *testing.T) {
	t.Parallel()

	input := `class="c-b c-a c-c"`
	want := `class="c-a c-b c-c"`

	got := normalizeClasses(input)
	if got != want {
		t.Errorf("normalizeClasses(%q) = %q, want %q", input, got, want)
	}
}

func TestNormalizeClassesMultiple(t *testing.T) {
	t.Parallel()

	input := `<div class="z-50 flex items-center">text</div>`

	got := normalizeClasses(input)
	if got != `<div class="flex items-center z-50">text</div>` {
		t.Errorf("unexpected: %s", got)
	}
}

// TestNormalizeClassesSkipsLongerAttributeNames pins the attribute-name
// boundary: a data-class value is a JavaScript OBJECT expression whose token
// order is its meaning, so only the real class attribute may be sorted.
func TestNormalizeClassesSkipsLongerAttributeNames(t *testing.T) {
	t.Parallel()

	input := `<button data-class="{ 'b a': $x==='y', 'd c': $x!=='y' }" class="z y" data-attr:data-class="q p">`

	want := `<button data-class="{ 'b a': $x==='y', 'd c': $x!=='y' }" class="y z" data-attr:data-class="q p">`

	if got := normalizeClasses(input); got != want {
		t.Errorf("normalizeClasses:\n got: %s\nwant: %s", got, want)
	}
}

// TestNormalizeClassesKeepsBoundary pins that the matched boundary character
// survives the rewrite (a dropped space or newline would glue attributes).
func TestNormalizeClassesKeepsBoundary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		html string
		want string
	}{
		{"space boundary", `<div id="a" class="b a">`, `<div id="a" class="a b">`},
		{"newline boundary", "<div\nclass=\"b a\">", "<div\nclass=\"a b\">"},
		{"line start", `class="b a"`, `class="a b"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := normalizeClasses(tt.html); got != tt.want {
				t.Errorf("normalizeClasses(%q) = %q, want %q", tt.html, got, tt.want)
			}
		})
	}
}
