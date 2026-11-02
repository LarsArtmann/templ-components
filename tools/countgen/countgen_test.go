package countgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// wrongNumber seeds every count-bearing doc site with a recognizably wrong
// value so Apply has something to fix; the site patterns locate the lines by
// structure, not by the number.
const wrongNumber = "99"

// syntheticTree writes a minimal tree exercising every counted axis:
// exported templ funcs (counted + ignored), icon entries (+Spinner),
// IsValid declarations (in-scope, test-file, .git, _sources, website),
// generated files, HTML goldens (in-scope + website), visual goldens
// (real + dotted-dir + .fail.png), and errorpage constructors.
func syntheticTree(t *testing.T) string {
	t.Helper()

	root := t.TempDir()

	write := func(rel, content string) {
		t.Helper()

		path := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", rel, err)
		}

		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
	}

	write("display/alpha.templ", "templ Card() {}\n\ntempl Badge() {}\n\ntempl helper() {}\n")
	write("display/alpha_templ.go", "// generated\n")
	write("display/enums.go", "package display\n\nfunc StatusIsValid(s Status) bool { return true }\n")
	write("display/enums_test.go", "package display\n\nfunc OtherIsValid(s Other) bool { return true }\n")
	write("forms/beta.templ", "templ Input() {}\n")
	write("recipes/gamma.templ", "templ Dashboard() {}\n")

	write("icons/icon_paths.go", "package icons\n\nvar iconPathData = map[string]string{\n"+
		"\tHome:  \"M1 1\",\n"+
		"\tClose: \"M2 2\",\n"+
		"}\n")

	write(
		"errorpage/constructors.go",
		"package errorpage\n\nfunc NotFound() *Error { return nil }\n\nfunc InternalError() *Error { return nil }\n",
	)

	write("display/snapshot.golden", "<div>golden</div>\n")
	write("website/excluded.golden", "<div>excluded</div>\n")
	write("website/excluded_templ.go", "// generated, excluded\n")

	write("visualtest/testdata/routes/dashboard.png", "png")
	write("visualtest/testdata/.fail/route.png", "png")
	write("visualtest/testdata/routes/flaky.fail.png", "png")

	write(".git/ignored.go", "func GitIsValid(s string) bool { return true }\n")
	write("cmd/tc/_sources/ignored.go", "func MirrorIsValid(s string) bool { return true }\n")

	write(featuresDoc, featuresFixture())
	write(skillDoc, skillFixture())

	return root
}

func featuresFixture() string {
	rows := &strings.Builder{}

	for _, pkg := range countedPackages() {
		rows.WriteString("| `" + pkg + "` | " + wrongNumber + " |\n")
	}

	return "# Features\n\n" +
		"**Totals:** " + wrongNumber + " templ components (primitives, drift-guard verified) + " +
		wrongNumber + " recipe screens = " + wrongNumber + ", " + wrongNumber + " icon names, 74 typed enums (" +
		wrongNumber + " with `IsValid()`, " + wrongNumber + " generated `*_templ.go` files\n\n" +
		"| Package | Components |\n| - | - |\n" +
		"| `icons` | " + wrongNumber + " (" + wrongNumber + " icons) |\n" +
		"| `recipes` | " + wrongNumber + " screens |\n" +
		rows.String() + "\n" +
		"- " + wrongNumber + " typed string enums (" + wrongNumber + " with `IsValid()` methods + tests)\n" +
		"- HTML snapshots (" + wrongNumber + " in 5 packages, " + wrongNumber + " baselines)\n" +
		"- Visual regression (" + wrongNumber + "; " + wrongNumber + " goldens incl. charts)\n"
}

func skillFixture() string {
	headings := &strings.Builder{}

	for _, pkg := range []string{"display", "forms", "feedback", "layout", "navigation", "htmx", "datastar"} {
		headings.WriteString("#### `" + pkg + "` — " + wrongNumber + " components\n\n")
	}

	return "# Skill\n\n" +
		wrongNumber + " components across " + wrongNumber + " packages (bogus counted templ packages, plus `utils` and `icons`)\n\n" +
		"- " + wrongNumber + " templ components + " + wrongNumber + " icons. Full signatures live elsewhere.\n\n" +
		"#### `errorpage` — " + wrongNumber + " components + " + wrongNumber + " constructors + handler\n\n" +
		"#### `charts/echarts` — " + wrongNumber + " components (opt-in module)\n\n" +
		headings.String()
}

func TestDeriveSyntheticTree(t *testing.T) {
	t.Parallel()

	counts, err := Derive(syntheticTree(t))
	if err != nil {
		t.Fatalf("Derive: %v", err)
	}

	want := Counts{
		Components:            3, // display 2 + forms 1
		PackageCounts:         map[string]int{"display": 2, "forms": 1, "recipes": 1, "icons": 0},
		Recipes:               1,
		ComponentsPlusRecipes: 4,
		Icons:                 3, // 2 entries + Spinner
		IsValid:               1, // enums.go only: test/.git/_sources excluded
		GeneratedFiles:        1, // website's excluded
		HTMLGoldens:           1, // website's excluded
		VisualGoldens:         1, // dotted dir + .fail.png excluded
		ErrorpageConstructors: 2,
	}

	for _, pkg := range countedPackages() {
		if _, ok := want.PackageCounts[pkg]; !ok && pkg != "recipes" && pkg != "icons" {
			want.PackageCounts[pkg] = 0
		}
	}

	if counts.Components != want.Components ||
		counts.Recipes != want.Recipes ||
		counts.ComponentsPlusRecipes != want.ComponentsPlusRecipes ||
		counts.Icons != want.Icons ||
		counts.IsValid != want.IsValid ||
		counts.GeneratedFiles != want.GeneratedFiles ||
		counts.HTMLGoldens != want.HTMLGoldens ||
		counts.VisualGoldens != want.VisualGoldens ||
		counts.ErrorpageConstructors != want.ErrorpageConstructors {
		t.Fatalf("Derive = %+v, want %+v", counts, want)
	}

	for pkg, wantCount := range want.PackageCounts {
		if got := counts.PackageCounts[pkg]; got != wantCount {
			t.Errorf("PackageCounts[%s] = %d, want %d", pkg, got, wantCount)
		}
	}
}

func TestCheckDriftApplyFixesIdempotence(t *testing.T) {
	t.Parallel()

	root := syntheticTree(t)

	err := Check(root)
	if err == nil {
		t.Fatal("Check on drifted docs: want error, got nil")
	}

	for _, site := range []string{"Totals line", "icons table row", "heading display"} {
		if !strings.Contains(err.Error(), site) {
			t.Errorf("Check error missing site %q:\n%v", site, err)
		}
	}

	if !strings.Contains(err.Error(), "go run ./tools/countgen/cmd/countgen -w") {
		t.Errorf("Check error missing the fix command:\n%v", err)
	}

	changed, err := Apply(root)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}

	if len(changed) != 2 || changed[0] != featuresDoc || changed[1] != skillDoc {
		t.Fatalf("Apply changed = %v, want [%s %s]", changed, featuresDoc, skillDoc)
	}

	if err := Check(root); err != nil {
		t.Fatalf("Check after Apply: %v", err)
	}

	if changed, err := Apply(root); err != nil || len(changed) != 0 {
		t.Fatalf("second Apply = %v, %v; want no writes (idempotence)", changed, err)
	}

	assertDocContains(
		t,
		root,
		featuresDoc,
		"**Totals:** 3 templ components (primitives, drift-guard verified) + 1 recipe screens = 4, 3 icon names, 74 typed enums (1 with `IsValid()`, 1 generated `*_templ.go` files",
	)
	assertDocContains(t, root, featuresDoc, "| `display` | 2 |")
	assertDocContains(t, root, featuresDoc, "| `icons` | 0 (3 icons)")
	assertDocContains(t, root, skillDoc,
		"3 components across 11 packages (nine counted templ packages, plus `utils` and `icons`)")
	assertDocContains(t, root, skillDoc, "#### `errorpage` — 0 components + 2 constructors + handler")
	assertDocContains(t, root, skillDoc, "#### `display` — 2 components")
}

func TestApplyPreservesHandTypedEnumsTotal(t *testing.T) {
	t.Parallel()

	root := syntheticTree(t)

	if _, err := Apply(root); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	features, err := os.ReadFile(filepath.Join(root, featuresDoc))
	if err != nil {
		t.Fatalf("read %s: %v", featuresDoc, err)
	}

	if !strings.Contains(string(features), "74 typed enums (1 with") {
		t.Error("Apply rewrote the hand-typed enums total (must stay passthrough)")
	}
}

func TestCheckStructuralDriftMultiMatch(t *testing.T) {
	t.Parallel()

	root := syntheticTree(t)

	features, err := os.ReadFile(filepath.Join(root, featuresDoc))
	if err != nil {
		t.Fatalf("read %s: %v", featuresDoc, err)
	}

	totals := strings.Split(string(features), "\n")
	doubled := make([]string, 0, len(totals)+1)

	for _, line := range totals {
		doubled = append(doubled, line)
		if strings.HasPrefix(line, "**Totals:**") {
			doubled = append(doubled, line)
		}
	}

	if err := os.WriteFile(filepath.Join(root, featuresDoc), []byte(strings.Join(doubled, "\n")), 0o644); err != nil {
		t.Fatalf("write %s: %v", featuresDoc, err)
	}

	err = Check(root)
	if err == nil {
		t.Fatal("Check with a duplicated totals line: want structural-drift error, got nil")
	}

	if !strings.Contains(err.Error(), `matched 2 times (want exactly 1)`) {
		t.Errorf("error missing multi-match diagnosis:\n%v", err)
	}
}

func TestDeriveMissingRequiredFile(t *testing.T) {
	t.Parallel()

	_, err := Derive(t.TempDir())
	if err == nil {
		t.Fatal("Derive on an empty root: want error, got nil")
	}

	if !strings.Contains(err.Error(), "icons/icon_paths.go") {
		t.Errorf("error should name the missing file, got: %v", err)
	}
}

func TestNumberWordBounds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		total int
		want  string
	}{
		{total: -1, want: "-1"},
		{total: 0, want: "zero"},
		{total: 9, want: "nine"},
		{total: 20, want: "twenty"},
		{total: 21, want: "21"},
		{total: 100, want: "100"},
	}

	for _, tt := range tests {
		if got := numberWord(tt.total); got != tt.want {
			t.Errorf("numberWord(%d) = %q, want %q", tt.total, got, tt.want)
		}
	}
}

func assertDocContains(t *testing.T, root, name, want string) {
	t.Helper()

	content, err := os.ReadFile(filepath.Join(root, name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}

	if !strings.Contains(string(content), want) {
		t.Errorf("%s missing %q", name, want)
	}
}
