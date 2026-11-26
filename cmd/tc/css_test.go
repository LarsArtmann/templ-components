package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const consumerEntry = `@import "tailwindcss" source(none);
@source "../**/*.templ";
@custom-variant dark (&:where(.dark, .dark *));
`

// wiringBlock is what a consumer adds to its entry after the first sync.
const wiringBlock = `@source "./templ-components-classes.txt";
@import "./templ-components.css";
`

// newCSSConsumer writes the smallest project `tc css` operates on: a go.mod
// pinning the library, and a Tailwind entry at css/app.css.
func newCSSConsumer(t *testing.T, pin string) string {
	t.Helper()

	root := t.TempDir()
	writeCSSFile(t, filepath.Join(root, "go.mod"), "module example.com/app\n\ngo 1.27\n\nrequire (\n\t"+
		libraryModulePath+" "+pin+"\n\tgithub.com/a-h/templ v0.3.1100 // indirect\n)\n")
	writeCSSFile(t, filepath.Join(root, "css", "app.css"), consumerEntry)

	return root
}

// syncedWiredConsumer is a project whose library files are synced and wired.
func syncedWiredConsumer(t *testing.T) (string, cssOptions) {
	t.Helper()

	root := newCSSConsumer(t, libraryVersion())
	opts := defaultCSSOptions()

	if err := cssSync(root, opts, &bytes.Buffer{}); err != nil {
		t.Fatalf("cssSync: %v", err)
	}

	entry := filepath.Join(root, "css", "app.css")
	writeCSSFile(t, entry, strings.Replace(readCSSFile(t, entry), "@source ", wiringBlock+"@source ", 1))

	return root, opts
}

func writeCSSFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func readCSSFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	return string(content)
}

func TestCSSSyncWritesStampedLibraryFiles(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer

	root := newCSSConsumer(t, libraryVersion())

	if err := cssSync(root, defaultCSSOptions(), &out); err != nil {
		t.Fatalf("cssSync: %v", err)
	}

	for _, file := range libraryFiles {
		content := readCSSFile(t, filepath.Join(root, "css", file.name))

		canonical, err := libraryContent(file.name)
		if err != nil {
			t.Fatalf("libraryContent(%s): %v", file.name, err)
		}

		version, body := splitStamp([]byte(content))
		if version != libraryVersion() {
			t.Errorf("%s stamped %q, want %q", file.name, version, libraryVersion())
		}

		if !bytes.Equal(body, canonical) {
			t.Errorf("%s body differs from the embedded canonical file", file.name)
		}
	}

	for _, want := range []string{
		`@source "./templ-components-classes.txt";`, `@import "./templ-components.css";`,
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("sync output lacks the wiring line %s; got:\n%s", want, out.String())
		}
	}
}

func TestCSSCheckPassesOnSyncedWiredConsumer(t *testing.T) {
	t.Parallel()

	root, opts := syncedWiredConsumer(t)

	var out bytes.Buffer

	if failures := cssCheck(root, opts, &out); failures != 0 {
		t.Fatalf("cssCheck reported %d failure(s) on a synced, wired consumer:\n%s", failures, out.String())
	}

	var syncOut bytes.Buffer

	if err := cssSync(root, opts, &syncOut); err != nil {
		t.Fatalf("second cssSync: %v", err)
	}

	if strings.Contains(syncOut.String(), "wrote") || strings.Contains(syncOut.String(), "refreshed") {
		t.Errorf("second sync rewrote up-to-date files:\n%s", syncOut.String())
	}

	if !strings.Contains(syncOut.String(), "already reaches every library file") {
		t.Errorf("second sync should confirm the wiring instead of hinting:\n%s", syncOut.String())
	}
}

func TestCSSCheckFailsWhileEntryIsNotWired(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer

	root := newCSSConsumer(t, libraryVersion())

	if err := cssSync(root, defaultCSSOptions(), &out); err != nil {
		t.Fatalf("cssSync: %v", err)
	}

	if failures := cssCheck(root, defaultCSSOptions(), &out); failures != 2 {
		t.Errorf("failures = %d, want 2 (classes not reached, stylesheet not imported):\n%s", failures, out.String())
	}

	for _, want := range []string{
		"a Tailwind entry reaches the library's classes",
		"a Tailwind entry imports the library stylesheet",
		`@source "./templ-components-classes.txt"; (after ` + "`tc css sync`" + `)`,
		"when `go mod vendor` runs before Tailwind",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("check output lacks %q:\n%s", want, out.String())
		}
	}
}

func TestCSSCheckFailsOnDriftAndMissingFiles(t *testing.T) {
	t.Parallel()

	root, opts := syncedWiredConsumer(t)
	custom := filepath.Join(root, "css", customName)
	corpus := filepath.Join(root, "css", corpusName)

	writeCSSFile(t, custom, readCSSFile(t, custom)+"\n.my-rule {}\n")

	var out bytes.Buffer

	// The entry still imports the synced stylesheet, so only the drifted file counts.
	failures := cssCheck(root, opts, &out)
	if failures != 1 || !strings.Contains(out.String(), "differs from "+libraryVersion()) {
		t.Errorf("drifted custom.css not reported (failures = %d):\n%s", failures, out.String())
	}

	if err := cssSync(root, opts, &bytes.Buffer{}); err != nil {
		t.Fatalf("cssSync repairs a stamped copy without --force: %v", err)
	}

	if err := os.Remove(corpus); err != nil {
		t.Fatalf("remove corpus: %v", err)
	}

	out.Reset()

	failures = cssCheck(root, opts, &out)
	if failures != 2 || !strings.Contains(out.String(), "templ-components-classes.txt ("+corpus+") present") ||
		!strings.Contains(out.String(), "a Tailwind entry reaches the library's classes") {
		t.Errorf("missing corpus not reported (failures = %d):\n%s", failures, out.String())
	}

	if err := cssSync(root, opts, &bytes.Buffer{}); err != nil {
		t.Fatalf("cssSync restores the corpus: %v", err)
	}

	out.Reset()

	if failures = cssCheck(root, opts, &out); failures != 0 {
		t.Errorf("check after repair reported %d failure(s):\n%s", failures, out.String())
	}
}

func TestCSSCheckReportsUnstampedDivergentCopy(t *testing.T) {
	t.Parallel()

	root, opts := syncedWiredConsumer(t)
	custom := filepath.Join(root, "css", customName)

	writeCSSFile(t, custom, "/* an older scaffold */\n.mine {}\n")

	var out bytes.Buffer

	failures := cssCheck(root, opts, &out)
	if failures != 1 || !strings.Contains(out.String(), "carries no tc css stamp") ||
		!strings.Contains(out.String(), "tc css sync --force") {
		t.Errorf("unstamped divergent copy not reported with the --force hint (failures = %d):\n%s", failures,
			out.String())
	}
}

func TestCSSSyncRefusesVersionMismatch(t *testing.T) {
	t.Parallel()

	root := newCSSConsumer(t, "v0.0.1")

	err := cssSync(root, defaultCSSOptions(), &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "go.mod pins v0.0.1 but this tc is "+libraryVersion()) ||
		!strings.Contains(err.Error(), syncCommand) {
		t.Fatalf("mismatched pin should refuse with the pinned-version command, got: %v", err)
	}

	for _, file := range libraryFiles {
		if _, statErr := os.Stat(filepath.Join(root, "css", file.name)); statErr == nil {
			t.Errorf("%s was written although the pin does not match this tc", file.name)
		}
	}
}

func TestCSSSyncRefusesUnstampedDivergentCopy(t *testing.T) {
	t.Parallel()

	root := newCSSConsumer(t, libraryVersion())
	custom := filepath.Join(root, "css", customName)

	writeCSSFile(t, custom, "/* mine */\n.mine {}\n")

	err := cssSync(root, defaultCSSOptions(), &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Fatalf("unstamped divergent copy should be refused, got: %v", err)
	}

	for _, file := range libraryFiles {
		if file.name == customName {
			continue
		}

		if _, statErr := os.Stat(filepath.Join(root, "css", file.name)); statErr == nil {
			t.Errorf("%s was written although the run was refused (sync must be all or nothing)", file.name)
		}
	}

	opts := defaultCSSOptions()
	opts.force = true

	if err := cssSync(root, opts, &bytes.Buffer{}); err != nil {
		t.Fatalf("--force cssSync: %v", err)
	}

	if !strings.HasPrefix(readCSSFile(t, custom), "/* "+stampMarker+libraryVersion()) {
		t.Errorf("--force did not replace the copy with the stamped canonical file")
	}

	opts = defaultCSSOptions()
	opts.scaffold = true

	writeCSSFile(t, custom, "/* mine again */\n.mine {}\n")

	if err := cssSync(root, opts, &bytes.Buffer{}); err != nil {
		t.Fatalf("scaffold cssSync: %v", err)
	}

	if readCSSFile(t, custom) != "/* mine again */\n.mine {}\n" {
		t.Errorf("scaffold mode must never replace an existing file")
	}
}

func TestCSSSyncScaffoldsWithoutGoMod(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	opts := defaultCSSOptions()
	opts.dir = filepath.Join(root, "styles")
	opts.scaffold = true

	var out bytes.Buffer

	if err := cssSync(root, opts, &out); err != nil {
		t.Fatalf("scaffold cssSync: %v", err)
	}

	if !strings.Contains(out.String(), "scaffolding anyway") {
		t.Errorf("scaffold mode should warn about the missing pin:\n%s", out.String())
	}

	for _, file := range libraryFiles {
		if _, err := os.Stat(filepath.Join(opts.dir, file.name)); err != nil {
			t.Errorf("%s missing after scaffold: %v", file.name, err)
		}
	}

	if err := cssSync(root, defaultCSSOptions(), &bytes.Buffer{}); err == nil {
		t.Errorf("a project without go.mod must not sync outside scaffold mode")
	}
}

func TestCSSDirectoryDerivesFromReferencedPaths(t *testing.T) {
	t.Parallel()

	root := newCSSConsumer(t, libraryVersion())
	writeCSSFile(
		t,
		filepath.Join(root, "css", "app.css"),
		"@import \"tailwindcss\" source(none);\n@import \"./base.css\";\n",
	)
	writeCSSFile(t, filepath.Join(root, "css", "base.css"),
		"@source \"../lib/tc/templ-components-classes.txt\";\n@import \"../lib/tc/custom.css\";\n"+
			"@import \"../lib/tc/templ-components-theme.css\";\n@import \"../lib/tc/safelist.css\";\n")

	if err := cssSync(root, defaultCSSOptions(), &bytes.Buffer{}); err != nil {
		t.Fatalf("cssSync: %v", err)
	}

	for _, file := range libraryFiles {
		if _, err := os.Stat(filepath.Join(root, "lib", "tc", file.name)); err != nil {
			t.Errorf("%s should be synced next to the paths the entry references: %v", file.name, err)
		}

		if _, err := os.Stat(filepath.Join(root, "css", file.name)); err == nil {
			t.Errorf("%s was synced beside the entry although the wiring points elsewhere", file.name)
		}
	}

	var out bytes.Buffer

	if failures := cssCheck(root, defaultCSSOptions(), &out); failures != 0 {
		t.Errorf("cssCheck reported %d failure(s):\n%s", failures, out.String())
	}
}

func TestCSSDiscoveryIgnoresCompiledOutputAndVendor(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeCSSFile(t, filepath.Join(root, "app.css"), consumerEntry)
	writeCSSFile(t, filepath.Join(root, "static", "app.css"), "@layer theme{:root{}}")
	writeCSSFile(t, filepath.Join(root, "vendor", "lib", "index.css"), consumerEntry)
	writeCSSFile(t, filepath.Join(root, "other", "site.css"), "@import 'tailwindcss';\n")
	writeCSSFile(t, filepath.Join(root, "commented", "x.css"), "/* @import \"tailwindcss\"; */\n.mine{}\n")

	entries, err := tailwindEntries(root, defaultCSSOptions())
	if err != nil {
		t.Fatalf("tailwindEntries: %v", err)
	}

	want := []string{filepath.Join(root, "app.css"), filepath.Join(root, "other", "site.css")}
	if len(entries) != len(want) {
		t.Fatalf("entries = %v, want %v", entries, want)
	}

	for i := range want {
		if entries[i] != want[i] {
			t.Errorf("entries[%d] = %s, want %s", i, entries[i], want[i])
		}
	}
}

func TestSplitStamp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		version string
		body    string
	}{
		{
			name:    "css stamp",
			content: "/* tc:sync templ-components v1.2.3 (x) */\n.a{}\n",
			version: "v1.2.3",
			body:    ".a{}\n",
		},
		{name: "text stamp", content: "# tc:sync templ-components v1.2.3 (x)\na\n", version: "v1.2.3", body: "a\n"},
		{name: "unstamped", content: ".a{}\n", version: "", body: ".a{}\n"},
		{
			name: "stamp not on first line", content: ".a{}\n/* tc:sync templ-components v1.2.3 */\n", version: "",
			body: ".a{}\n/* tc:sync templ-components v1.2.3 */\n",
		},
		{
			name: "stamp not at line start", content: ".a{} /* tc:sync templ-components v1.2.3 */\n", version: "",
			body: ".a{} /* tc:sync templ-components v1.2.3 */\n",
		},
		{name: "empty", content: "", version: "", body: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			version, body := splitStamp([]byte(test.content))
			if version != test.version || string(body) != test.body {
				t.Errorf("splitStamp() = (%q, %q), want (%q, %q)", version, body, test.version, test.body)
			}
		})
	}
}

func TestCSSDirectives(t *testing.T) {
	t.Parallel()

	text := `@import "tailwindcss" source(none);
@import "tailwindcss/theme.css";
/* @import "./commented.css"; */
@source not "../docs";
@source "./templ-components-classes.txt";
@source inline("[p-4]");
@source "../views/**/*.templ";
@import "../lib/custom.css" layer(components);
@import url("https://example.com/x.css");
@import 'single quoted.css';
@import "../noext";
.mine { content: "@import fake.css"; }
`
	directives := cssDirectives(text)

	want := []string{
		"source ./templ-components-classes.txt",
		"source ../views/**/*.templ",
		"import ../lib/custom.css",
		"import single quoted.css",
		"import ../noext",
	}

	if len(directives) != len(want) {
		t.Fatalf("cssDirectives() returned %d directive(s), want %d: %+v", len(directives), len(want), directives)
	}

	for i, expected := range want {
		kind := "source"
		if directives[i].isImport {
			kind = "import"
		}

		if got := kind + " " + directives[i].target; got != expected {
			t.Errorf("directive[%d] = %q, want %q", i, got, expected)
		}
	}
}

func TestRequireVersionFindsModuleInBothForms(t *testing.T) {
	t.Parallel()

	content := "module x\n\nrequire github.com/a-h/templ v0.3.1100\n\nrequire (\n\t" + libraryModulePath +
		" v1.27.2 // indirect\n\t" + libraryModulePath + "/utils v1.27.2\n)\n\nreplace " + libraryModulePath +
		" => ../templ-components\n"

	tests := []struct {
		module string
		want   string
		found  bool
	}{
		{module: libraryModulePath, want: "v1.27.2", found: true},
		{module: libraryModulePath + "/utils", want: "v1.27.2", found: true},
		{module: "github.com/a-h/templ", want: "v0.3.1100", found: true},
		{module: "example.com/absent", want: "", found: false},
	}

	for _, test := range tests {
		version, found := requireVersion(content, test.module)
		if version != test.want || found != test.found {
			t.Errorf("requireVersion(%q) = (%q, %v), want (%q, %v)", test.module, version, found, test.want, test.found)
		}
	}
}

func TestFindGoModWalksUpward(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	writeCSSFile(t, filepath.Join(root, "go.mod"), "module x\n")

	found, ok := findGoMod(filepath.Join(root, "a", "b"))
	if !ok || found != filepath.Join(root, "go.mod") {
		t.Errorf("findGoMod() = (%q, %v), want (%q, true)", found, ok, filepath.Join(root, "go.mod"))
	}

	if _, ok := findGoMod(t.TempDir()); ok {
		t.Errorf("findGoMod() found a go.mod above a directory that has none")
	}
}

func TestParseCSSArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args    []string
		mode    string
		opts    cssOptions
		wantErr bool
	}{
		{args: []string{"sync"}, mode: "sync", opts: defaultCSSOptions()},
		{
			args: []string{"sync", "--force", "--dir", "styles"}, mode: "sync",
			opts: cssOptions{dir: "styles", force: true, css: "", scaffold: false},
		},
		{
			args: []string{"check", "--css", "css/app.css"}, mode: "check",
			opts: cssOptions{dir: "", force: false, css: "css/app.css", scaffold: false},
		},
		{args: nil, wantErr: true},
		{args: []string{"build"}, wantErr: true},
		{args: []string{"sync", "--dir"}, wantErr: true},
		{
			args: []string{"check", "--force"}, mode: "check",
			opts: cssOptions{dir: "", force: true, css: "", scaffold: false},
		},
	}

	for _, test := range tests {
		mode, opts, err := parseCSSArgs(test.args)
		if (err != nil) != test.wantErr || mode != test.mode || opts != test.opts {
			t.Errorf("parseCSSArgs(%v) = (%q, %+v, %v), want (%q, %+v, err=%v)", test.args, mode, opts, err,
				test.mode, test.opts, test.wantErr)
		}
	}
}

// TestStarterFilesMatchTemplates keeps the embedded starter identical to
// templates/: the starter is excluded from the mirror guard ("curated, no
// twins"), which let its theme copy drift once. `tc css sync` writes these
// bytes, so drift here ships stale CSS to every consumer.
func TestStarterFilesMatchTemplates(t *testing.T) {
	t.Parallel()

	names := make([]string, 0, len(libraryFiles)+1)
	names = append(names, starterEntryName)

	for _, file := range libraryFiles {
		names = append(names, file.name)
	}

	for _, name := range names {
		starter, err := sourcesFS.ReadFile(filepath.Join("_sources", "starter", name))
		if err != nil {
			t.Errorf("starter does not embed %s: %v", name, err)

			continue
		}

		template, err := os.ReadFile(filepath.Join("..", "..", "templates", name))
		if err != nil {
			t.Errorf("templates/%s unreadable: %v", name, err)

			continue
		}

		if !bytes.Equal(starter, template) {
			t.Errorf("cmd/tc/_sources/starter/%s differs from templates/%s; re-copy it (tc css sync ships these bytes)",
				name, name)
		}
	}

	corpus, err := libraryContent(corpusName)
	if err != nil || !strings.Contains(string(corpus), "package display") {
		t.Errorf("the embedded class corpus should carry the library sources, got (%v, %d bytes)", err, len(corpus))
	}
}
