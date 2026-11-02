package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// The dependency budget is closed (templ + tailwind-merge-go +
// go-error-family) and every coupling travels through a known set of go.mods.
// The 2026-11-01 audit caught exactly this class SHIPPED: errorpage required
// go-error-family v0.11.0 while root/visualtest/website sat at v0.12.0 —
// legal under workspace MVS, invisible to every workspace build, and a live
// consumer-facing version skew until the next proxy fetch. These guards turn
// that class into a <50ms test failure: the pin SURFACE (which go.mods may
// mention a dependency at all) and the pin TRAIN (every mention carries the
// same version) are both asserted per dependency.
//
// templ is deliberately NOT re-guarded here: TestTemplVersionPin already
// anchors every go.mod's templ require to the generator constant in
// cmd/tc/doctor.go — a stronger form of the same train (single source of
// truth beats pairwise equality).

// pinSpec describes one dependency's pin train: the module path, the module
// whose go.mod carries the direct require line, and every go.mod allowed to
// mention the dependency (the direct module included).
type pinSpec struct {
	dep          string
	directModule string
	allowedMods  []string
}

var (
	goErrorFamilyPin = pinSpec{
		dep:          "github.com/larsartmann/go-error-family",
		directModule: "errorpage",
		allowedMods:  []string{".", "errorpage", "visualtest", "website"},
	}

	tailwindMergePin = pinSpec{
		dep:          "github.com/Oudwins/tailwind-merge-go",
		directModule: "utils",
		// All 9 modules mention it: utils directly, the rest indirectly
		// through utils' class-merging call site (charts/echarts included —
		// the guard itself caught the initial inventory missing it).
		allowedMods: []string{
			".", "utils", "icons", "htmx", "datastar", "errorpage",
			"visualtest", "website", "charts/echarts",
		},
	}
)

// pinLineRe matches a require line for dep: module path followed by a
// v-prefixed version. Anchored at line start so `replace` lines and comments
// never match.
func pinLineRe(dep string) *regexp.Regexp {
	return regexp.MustCompile(`(?m)^[ \t]*(?:require[ \t]+)?` + regexp.QuoteMeta(dep) + `[ \t]+(v\S+)`)
}

// moduleLabel renders a module dir as its go.mod path for failure messages.
func moduleLabel(modDir string) string {
	if modDir == "." {
		return "go.mod (root)"
	}

	return modDir + "/go.mod"
}

// workspaceGoMods returns every go.mod under root, skipping VCS/vendor/build
// output directories (hidden dirs, node_modules, testdata, result* symlinks,
// the cmd/tc self-healing mirror).
func workspaceGoMods(tb testing.TB, root string) []string {
	tb.Helper()

	var mods []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			name := d.Name()
			skip := strings.HasPrefix(name, ".") ||
				name == "node_modules" ||
				name == "testdata" ||
				name == "_sources" ||
				strings.HasPrefix(name, "result")

			if skip && path != root {
				return filepath.SkipDir
			}

			return nil
		}

		if d.Name() == "go.mod" {
			mods = append(mods, path)
		}

		return nil
	})
	if err != nil {
		tb.Fatalf("walk %s for go.mods: %v", root, err)
	}

	sort.Strings(mods)

	return mods
}

// checkPinTrain asserts one dependency's pin surface and version train:
//
//  1. every go.mod mentioning the dependency is in spec.allowedMods (no
//     module silently gained a dependency), and every allowed module
//     mentions it (no lockstep desync);
//  2. the direct module's version is the anchor and every other mention
//     matches it (no MVS-masked skew like the 2026-11-01 go-error-family
//     incident).
//
// Surface membership is what is pinned — a module's direct/indirect marker
// is free to flip under `go mod tidy` and is not asserted.
func checkPinTrain(tb testing.TB, repoRoot string, spec pinSpec) {
	tb.Helper()

	versions := collectPinVersions(tb, repoRoot, spec.dep)
	checkPinSurface(tb, spec, versions)
	checkTrainVersions(tb, spec, versions)
}

// collectPinVersions maps each module dir to the versions it requires dep
// with (scanning the workspace's go.mods).
func collectPinVersions(tb testing.TB, repoRoot, dep string) map[string][]string {
	tb.Helper()

	versions := map[string][]string{}

	for _, mod := range workspaceGoMods(tb, repoRoot) {
		content, err := os.ReadFile(mod)
		if err != nil {
			tb.Fatalf("read %s: %v", mod, err)
		}

		rel, err := filepath.Rel(repoRoot, mod)
		if err != nil {
			tb.Fatalf("rel %s: %v", mod, err)
		}

		dir := filepath.ToSlash(filepath.Dir(rel))

		for _, m := range pinLineRe(dep).FindAllStringSubmatch(string(content), -1) {
			versions[dir] = append(versions[dir], m[1])
		}
	}

	return versions
}

// checkPinSurface asserts exactly the allowed modules mention the dependency.
func checkPinSurface(tb testing.TB, spec pinSpec, versions map[string][]string) {
	tb.Helper()

	allowed := make(map[string]bool, len(spec.allowedMods))

	for _, m := range spec.allowedMods {
		allowed[m] = true
	}

	for dir := range versions {
		if !allowed[dir] {
			tb.Errorf(
				"%s mentions %s but is not in the allowed pin surface %v — a module silently gained a dependency (consumer-graph decision, not a mechanical bump)",
				moduleLabel(dir),
				spec.dep,
				spec.allowedMods,
			)
		}
	}

	for _, m := range spec.allowedMods {
		if _, ok := versions[m]; !ok {
			tb.Errorf(
				"%s does not mention %s but the pin surface requires it — the release lockstep will desync this module's go.sum",
				moduleLabel(m),
				spec.dep,
			)
		}
	}
}

// checkTrainVersions asserts the version train: the direct module's version
// anchors, every other mention matches, and no module carries two versions.
func checkTrainVersions(tb testing.TB, spec pinSpec, versions map[string][]string) {
	tb.Helper()

	direct := versions[spec.directModule]
	if len(direct) == 0 {
		tb.Fatalf(
			"%s carries no require for %s — the pin train has no anchor (direct module renamed?)",
			moduleLabel(spec.directModule), spec.dep,
		)
	}

	want := direct[0]
	for _, v := range direct[1:] {
		if v != want {
			tb.Errorf(
				"%s requires %s at both %s and %s — one module, two pins",
				moduleLabel(spec.directModule), spec.dep, want, v,
			)
		}
	}

	for dir, vs := range versions {
		if dir == spec.directModule {
			continue
		}

		for _, v := range vs {
			if v != want {
				tb.Errorf(
					"%s requires %s %s but the pin train is anchored at %s (direct in %s) — align the versions in one commit",
					moduleLabel(dir),
					spec.dep,
					v,
					want,
					moduleLabel(spec.directModule),
				)
			}
		}
	}
}

func TestGoErrorFamilyPinAlignment(t *testing.T) {
	t.Parallel()

	checkPinTrain(t, "..", goErrorFamilyPin)
}

func TestTailwindMergePinAlignment(t *testing.T) {
	t.Parallel()

	checkPinTrain(t, "..", tailwindMergePin)
}

// recordingTB captures failures instead of failing the outer test, so the
// guard's detection logic can be negative-tested against synthetic trees.
// The embedded nil testing.TB satisfies the (Go 1.27, ArtifactDir-carrying)
// interface; only Helper/Errorf/Fatalf are ever invoked.
type recordingTB struct {
	testing.TB

	failures []string
}

func (r *recordingTB) Helper() {}

func (r *recordingTB) Errorf(format string, args ...any) {
	r.failures = append(r.failures, fmt.Sprintf(format, args...))
}

type guardFatal struct{}

func (r *recordingTB) Fatalf(format string, args ...any) {
	r.Errorf(format, args...)
	panic(guardFatal{})
}

// runGuardOnTree runs checkPinTrain against a synthetic tree, mirroring
// testing's Fatal semantics (panic stops the check) while collecting every
// failure message.
func runGuardOnTree(root string, spec pinSpec) []string {
	fake := &recordingTB{}

	func() {
		defer func() {
			if r := recover(); r != nil {
				if _, ok := r.(guardFatal); !ok {
					panic(r)
				}
			}
		}()

		checkPinTrain(fake, root, spec)
	}()

	return fake.failures
}

// writeTempGoMod materializes one module's go.mod in a synthetic tree.
func writeTempGoMod(t *testing.T, root, dir, requireLines string) {
	t.Helper()

	full := filepath.Join(root, filepath.FromSlash(dir))
	if err := os.MkdirAll(full, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", full, err)
	}

	body := "module example.com/" + strings.ReplaceAll(
		dir,
		"/",
		"-",
	) + "\n\ngo 1.27\n\nrequire (\n" + requireLines + ")\n"
	if err := os.WriteFile(filepath.Join(full, "go.mod"), []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", filepath.Join(full, "go.mod"), err)
	}
}

func TestCheckPinTrainDetectsDrift(t *testing.T) {
	t.Parallel()

	// The REAL spec: the skew that shipped 2026-11-01 (errorpage v0.11.0
	// while every indirect consumer sat at v0.12.0) must fail this guard.
	spec := goErrorFamilyPin
	dep := spec.dep

	t.Run("version skew fails and names both versions", func(t *testing.T) {
		t.Parallel()

		root := t.TempDir()
		writeTempGoMod(t, root, ".", "\t"+dep+" v0.12.0 // indirect\n")
		writeTempGoMod(t, root, "errorpage", "\t"+dep+" v0.11.0\n")
		writeTempGoMod(t, root, "visualtest", "\t"+dep+" v0.12.0 // indirect\n")
		writeTempGoMod(t, root, "website", "\t"+dep+" v0.12.0 // indirect\n")

		joined := strings.Join(runGuardOnTree(root, spec), "\n")
		if !strings.Contains(joined, "v0.11.0") || !strings.Contains(joined, "v0.12.0") {
			t.Errorf("expected skew failure naming both versions, got:\n%s", joined)
		}
	})

	t.Run("surface drift fails", func(t *testing.T) {
		t.Parallel()

		root := t.TempDir()
		for _, m := range spec.allowedMods {
			writeTempGoMod(t, root, m, "\t"+dep+" v0.12.0 // indirect\n")
		}

		writeTempGoMod(t, root, "charts/echarts", "\t"+dep+" v0.12.0 // indirect\n")

		joined := strings.Join(runGuardOnTree(root, spec), "\n")
		if !strings.Contains(joined, "charts/echarts/go.mod") || !strings.Contains(joined, "pin surface") {
			t.Errorf("expected surface-drift failure naming charts/echarts, got:\n%s", joined)
		}
	})

	t.Run("missing module fails", func(t *testing.T) {
		t.Parallel()

		root := t.TempDir()
		writeTempGoMod(t, root, ".", "\t"+dep+" v0.12.0 // indirect\n")
		writeTempGoMod(t, root, "errorpage", "\t"+dep+" v0.12.0\n")
		writeTempGoMod(t, root, "visualtest", "\t"+dep+" v0.12.0 // indirect\n")

		joined := strings.Join(runGuardOnTree(root, spec), "\n")
		if !strings.Contains(joined, "website") {
			t.Errorf("expected missing-surface failure naming website, got:\n%s", joined)
		}
	})

	t.Run("aligned train passes", func(t *testing.T) {
		t.Parallel()

		root := t.TempDir()
		for _, m := range spec.allowedMods {
			writeTempGoMod(t, root, m, "\t"+dep+" v0.12.0 // indirect\n")
		}

		if failures := runGuardOnTree(root, spec); len(failures) != 0 {
			t.Errorf("aligned train must pass cleanly, got:\n%s", strings.Join(failures, "\n"))
		}
	})

	t.Run("replace lines never count as pins", func(t *testing.T) {
		t.Parallel()

		root := t.TempDir()
		writeTempGoMod(t, root, ".", "\t"+dep+" v0.12.0 // indirect\n")
		writeTempGoMod(t, root, "errorpage", "\t"+dep+" v0.12.0\n\treplace "+dep+" => ../dep\n")
		writeTempGoMod(t, root, "visualtest", "\t"+dep+" v0.12.0 // indirect\n")
		writeTempGoMod(t, root, "website", "\t"+dep+" v0.12.0 // indirect\n")

		if failures := runGuardOnTree(root, spec); len(failures) != 0 {
			t.Errorf("replace directives must not register as pin mentions, got:\n%s", strings.Join(failures, "\n"))
		}
	})
}
