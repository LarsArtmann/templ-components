package golden

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAssertMatchesGoldenFile(t *testing.T) {
	t.Parallel()

	content := "<div class=\"bg-blue-600\">Hello</div>\n"
	name := "test_assert"
	goldenPath := filepath.Join("testdata", name+".golden")

	if err := os.MkdirAll("testdata", 0o750); err != nil {
		t.Fatalf("setup: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll("testdata")
	})

	if err := os.WriteFile(goldenPath, []byte(content), 0o600); err != nil {
		t.Fatalf("setup write: %v", err)
	}

	Assert(t, name, content)
}

func TestAssertRejectsMismatch(t *testing.T) {
	t.Parallel()

	name := "test_mismatch"
	goldenPath := filepath.Join("testdata", name+".golden")

	if err := os.MkdirAll("testdata", 0o750); err != nil {
		t.Fatalf("setup: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll("testdata")
	})

	if err := os.WriteFile(goldenPath, []byte("old content\n"), 0o600); err != nil {
		t.Fatalf("setup write: %v", err)
	}

	mockT := &testing.T{}
	Assert(mockT, name, "new content\n")

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
	if !strings.Contains(got, "--- line2") || !strings.Contains(got, "+++ changed") {
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
