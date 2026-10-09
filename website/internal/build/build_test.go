package build

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// tagRemnantRe matches tag-opening shapes ("<h2", "</p") that must never
// survive PlainText (bare "<"/">" in decoded text content are legitimate).
var tagRemnantRe = regexp.MustCompile(`<[a-zA-Z/]`)

// writeFixtureRepo lays out a minimal fake library checkout for CountStats.
func writeFixtureRepo(t *testing.T) string {
	t.Helper()

	root := t.TempDir()

	files := map[string]string{
		"go.mod":                "module example.com/lib\n\ngo 1.26.7\n",
		"display/go.mod":        "module example.com/lib/display\n",
		"display/card.templ":    "templ Card(props CardProps) {\n}\n\ntempl Badge(text string) {\n}\n",
		"display/types.go":      "package display\n\nfunc BadgeTypeIsValid(v BadgeType) bool { return true }\n",
		"display/types_test.go": "package display\n\nfunc NeverIsValid(v string) bool { return false }\n",
		"feedback/go.mod":       "module example.com/lib/feedback\n",
		"feedback/alert.templ":  "templ Alert(message string) {\n}\n",
		"icons/icon_paths.go":   "package icons\n\nvar iconPathData = map[Name]string{\n\tHome: \"M0 0\",\n\tX:    svg.PathXMark,\n}\n",
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

	if stats.Icons != 3 {
		t.Errorf("Icons = %d, want 3 (2 path entries + Spinner)", stats.Icons)
	}

	if stats.Enums != 1 {
		t.Errorf("Enums = %d, want 1 (_test.go excluded)", stats.Enums)
	}

	if stats.Modules != 3 {
		t.Errorf("Modules = %d, want 3 (root + display + feedback)", stats.Modules)
	}

	if stats.GoVersion != "1.26" {
		t.Errorf("GoVersion = %q, want \"1.26\" (patch segment trimmed)", stats.GoVersion)
	}

	if stats.LibraryVersion != utils.Version {
		t.Errorf("LibraryVersion = %q, want utils.Version (%q)", stats.LibraryVersion, utils.Version)
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

func TestLLMSIndex(t *testing.T) {
	t.Parallel()

	stats := Stats{Components: 123, Icons: 102, Enums: 63, Modules: 7, GoVersion: "1.27", LibraryVersion: "1.21.0"}
	docs := []SearchDoc{
		{URL: "/install", Title: "Installation", Description: "go get and templ generate", Body: "ignored"},
		{URL: "/guide", Title: "Guide\nwith newline", Body: "ignored"},
	}

	got := LLMSIndex(stats, docs)

	if !strings.HasPrefix(got, "# templ-components\n\n> ") {
		t.Errorf("missing llms.txt H1 + blockquote shape: %q", got[:60])
	}

	for _, want := range []string{
		"123 components, 102 SVG icons, 63 typed enums, 7 opt-in Go modules",
		"## Product",
		"](" + "https://templcomponents.lars.software/sales)",
		"## Documentation",
		"- [Installation](https://templcomponents.lars.software/install): go get and templ generate",
		"- [Guide with newline](https://templcomponents.lars.software/guide)",
		"Version 1.21.0",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("llms.txt missing %q\ngot:\n%s", want, got)
		}
	}

	if strings.Contains(got, "\n\n\n") {
		t.Errorf("blank-line runs in llms.txt:\n%s", got)
	}
}

// TestStatsDirsAreCanonical pins statsDirs to the canonical 9-package
// component-scan set, order-sensitive: any add, removal, rename, or reorder
// changes what the site's hero "N components" counts, so it must be a
// deliberate two-place update (here and utils.TestDocsCountDrift's canonical
// package set).
func TestStatsDirsAreCanonical(t *testing.T) {
	canonical := []string{
		"display", "feedback", "forms", "layout", "navigation",
		"charts/echarts", "htmx", "datastar", "errorpage",
	}

	if !slices.Equal(statsDirs, canonical) {
		t.Errorf(
			"statsDirs drifted from the canonical set — update BOTH statsDirs and utils.TestDocsCountDrift's package set together, then re-pin here.\ngot:  %q\nwant: %q",
			statsDirs,
			canonical,
		)
	}
}
