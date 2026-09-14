package build

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// tagRemnantRe matches tag-opening shapes ("<h2", "</p") that must never
// survive PlainText (bare "<"/">" in decoded text content are legitimate).
var tagRemnantRe = regexp.MustCompile(`<[a-zA-Z/]`)

// writeFixtureRepo lays out a minimal fake library checkout for CountStats.
func writeFixtureRepo(t *testing.T) string {
	t.Helper()

	root := t.TempDir()

	files := map[string]string{
		"go.mod":                "module example.com/lib\n",
		"display/go.mod":        "module example.com/lib/display\n",
		"display/card.templ":    "templ Card(props CardProps) {\n}\n\ntempl Badge(text string) {\n}\n",
		"display/types.go":      "package display\n\nfunc BadgeTypeIsValid(v BadgeType) bool { return true }\n",
		"display/types_test.go": "package display\n\nfunc NeverIsValid(v string) bool { return false }\n",
		"feedback/go.mod":       "module example.com/lib/feedback\n",
		"feedback/alert.templ":  "templ Alert(message string) {\n}\n",
		"icons/icon_names.go":   "package icons\n\nconst (\n\tHome  Name = \"home\"\n\tClose Name = \"home\"\n\tX     Name = \"x\"\n)",
	}

	for path, content := range files {
		target := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			t.Fatalf("mkdir for %s: %v", path, err)
		}

		if err := os.WriteFile(target, []byte(content), 0o600); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	return root
}

func TestCountStats(t *testing.T) {
	t.Parallel()

	stats, err := CountStats(writeFixtureRepo(t))
	if err != nil {
		t.Fatalf("CountStats: %v", err)
	}

	if stats.Components != 3 {
		t.Errorf("Components = %d, want 3 (2 display + 1 feedback)", stats.Components)
	}

	if stats.Icons != 2 {
		t.Errorf("Icons = %d, want 2 (aliases deduplicated)", stats.Icons)
	}

	if stats.Enums != 1 {
		t.Errorf("Enums = %d, want 1 (_test.go excluded)", stats.Enums)
	}

	if stats.Modules != 3 {
		t.Errorf("Modules = %d, want 3 (root + display + feedback)", stats.Modules)
	}
}

func TestPlainText(t *testing.T) {
	t.Parallel()

	markup := `<h2 id="setup">Setup</h2>` +
		`<p>Install <code>go&nbsp;&gt;= 1.26</code> &amp; templ.</p>` +
		`<script>var tracked = "never indexed";</script>` +
		`<style>.x { color: red; }</style>` +
		`<pre><code class="language-bash">go get example.com/lib</code></pre>` +
		`<ul><li>First</li><li>Second</li></ul>`

	got := PlainText(markup)

	if strings.Contains(got, "tracked") || strings.Contains(got, "color: red") {
		t.Errorf("script/style content leaked: %q", got)
	}

	if !strings.Contains(got, "go >= 1.26") {
		t.Errorf("entity decoding lost: %q", got)
	}

	if !strings.Contains(got, "go get example.com/lib") {
		t.Errorf("code fences must stay searchable: %q", got)
	}

	if tagRemnantRe.MatchString(got) || strings.Contains(got, "&amp;") || strings.Contains(got, "&nbsp;") {
		t.Errorf("markup survived: %q", got)
	}

	if strings.Contains(got, "  ") {
		t.Errorf("whitespace not collapsed: %q", got)
	}
}

func TestWriteSearchIndex(t *testing.T) {
	t.Parallel()

	outDir := t.TempDir()

	docs := []SearchDoc{
		{URL: "/a", Title: "A", Body: "alpha", Sections: []SearchSection{{ID: "s", Text: "S"}}},
		{URL: "/b", Title: "B", Body: "beta"},
	}

	if err := WriteSearchIndex(outDir, docs); err != nil {
		t.Fatalf("WriteSearchIndex: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outDir, "search-index.json"))
	if err != nil {
		t.Fatalf("read index: %v", err)
	}

	if !strings.Contains(string(data), `"url":"/a"`) {
		t.Errorf("unexpected index content: %s", data)
	}
}
