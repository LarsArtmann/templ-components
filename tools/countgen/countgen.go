// Package countgen derives the library's count-bearing numbers from the tree
// and writes them into the two count-table docs (FEATURES.md, skill/SKILL.md).
//
// WHY THIS EXISTS: every component addition forced a hand-typed same-commit
// count bump in FEATURES/SKILL (TODO #996) — the tax this generator deletes.
// Run `go run ./tools/countgen/cmd/countgen -w` after any change that moves a
// count; commit the diff in the same commit as the source change.
//
// INDEPENDENT VERIFIER: utils.TestDocsCountDrift re-derives every count with
// its own implementation and asserts the docs directly — it is deliberately
// NOT this package's code, so the two layers cross-check each other. If they
// ever disagree, CI goes red and the reconciliation is the signal.
//
// DELIBERATELY NOT DERIVED (no mechanical source; semantics owned by
// docs/count-definitions.md):
//   - the enums TOTAL ("74 typed enums") — "exported closed-set string types"
//     is a semantic judgment; only the IsValid count is mechanical;
//   - the FEATURES.md `utils` row ("0"): utils owns no components by
//     definition, yet utils/dismiss.templ exports a helper templ function —
//     a mechanical count would print 1 against a deliberate 0;
//   - "~31,500 lines of Go/templ source" and "in 5 packages" prose counts;
//   - the "8 components accept ContainerAware" figure (ADR-0018 lookup-map
//     semantics).
package countgen

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	featuresDoc = "FEATURES.md"
	skillDoc    = "skill/SKILL.md"
)

// countedPackages are the packages whose exported templ functions compose the
// "N templ components" total (recipes screens counted separately).
var countedPackages = []string{
	"display", "feedback", "forms", "navigation", "errorpage",
	"layout", "htmx", "datastar", "charts/echarts",
}

// tablePackages are the packages whose exported templ functions countgen
// derives (utils excluded: its FEATURES row is the semantic "owns no
// components", not a templ-function count — utils/dismiss.templ exports
// one helper function and must keep printing 0).
var tablePackages = append([]string{"icons", "recipes"}, countedPackages...)

// Counts holds every number countgen derives from the tree.
type Counts struct {
	Components            int            // exported templ funcs across the 9 counted packages
	PackageCounts         map[string]int // per package (counted 9 + icons + recipes)
	Recipes               int            // recipes/*.templ exported screens
	ComponentsPlusRecipes int            // Components + Recipes
	Icons                 int            // iconPathData entries + Spinner
	IsValid               int            // <Name>IsValid( declarations, library source
	GeneratedFiles        int            // *_templ.go files (website excluded)
	HTMLGoldens           int            // committed .golden baselines (visualtest/website excluded)
	VisualGoldens         int            // visualtest/testdata .png goldens
	ErrorpageConstructors int            // exported constructors in errorpage/constructors.go
}

var (
	templFuncRe   = regexp.MustCompile(`^templ\s+([A-Z][A-Za-z0-9]*)\s*\(`)
	iconEntryRe   = regexp.MustCompile(`(?m)^\t[A-Z][A-Za-z0-9]*:\s+("[^"]*"|svg\.Path[A-Za-z]*)`)
	isValidRe     = regexp.MustCompile(`func\s+[A-Z][A-Za-z0-9]*IsValid\s*\(`)
	constructorRe = regexp.MustCompile(`(?m)^func [A-Z]`)
)

// Derive computes every count from the tree rooted at root.
func Derive(root string) (Counts, error) {
	counts := Counts{PackageCounts: map[string]int{}}

	for _, pkg := range tablePackages {
		total, err := countPackageTemplFuncs(root, pkg)
		if err != nil {
			return counts, err
		}

		counts.PackageCounts[pkg] = total
	}

	for _, pkg := range countedPackages {
		counts.Components += counts.PackageCounts[pkg]
	}

	counts.Recipes = counts.PackageCounts["recipes"]
	counts.ComponentsPlusRecipes = counts.Components + counts.Recipes

	if err := deriveWalkedCounts(root, &counts); err != nil {
		return counts, err
	}

	constructorSource, err := os.ReadFile(filepath.Join(root, "errorpage", "constructors.go"))
	if err != nil {
		return counts, fmt.Errorf("read errorpage/constructors.go: %w", err)
	}

	counts.ErrorpageConstructors = len(constructorRe.FindAll(constructorSource, -1))

	return counts, nil
}

func deriveWalkedCounts(root string, counts *Counts) error {
	icons, err := countIcons(root)
	if err != nil {
		return err
	}

	counts.Icons = icons

	isValid, err := countMatchesWalk(root, isValidRe, "IsValid")
	if err != nil {
		return err
	}

	counts.IsValid = isValid

	generated, err := countSuffixWalk(root, "_templ.go", "generated")
	if err != nil {
		return err
	}

	counts.GeneratedFiles = generated

	htmlGoldens, err := countGoldenWalk(root)
	if err != nil {
		return err
	}

	counts.HTMLGoldens = htmlGoldens

	visualGoldens, err := countVisualGoldens(root)
	if err != nil {
		return err
	}

	counts.VisualGoldens = visualGoldens

	return nil
}

func countPackageTemplFuncs(root, pkg string) (int, error) {
	files, err := filepath.Glob(filepath.Join(root, pkg, "*.templ"))
	if err != nil {
		return 0, fmt.Errorf("glob %s: %w", pkg, err)
	}

	total := 0

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return 0, fmt.Errorf("read %s: %w", file, err)
		}

		for line := range strings.SplitSeq(string(data), "\n") {
			if templFuncRe.MatchString(line) {
				total++
			}
		}
	}

	return total, nil
}

func countIcons(root string) (int, error) {
	data, err := os.ReadFile(filepath.Join(root, "icons", "icon_paths.go"))
	if err != nil {
		return 0, fmt.Errorf("read icons/icon_paths.go: %w", err)
	}

	return len(iconEntryRe.FindAll(data, -1)) + 1, nil // +1: Spinner
}

// countMatchesWalk mirrors utils' countIsValidMethods: non-test .go files,
// skipping .git, website, and the cmd/tc _sources mirror.
func countMatchesWalk(root string, pattern *regexp.Regexp, label string) (int, error) {
	total := 0

	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			if base := filepath.Base(path); base == ".git" || base == "website" || base == "_sources" {
				return filepath.SkipDir
			}

			return nil
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		total += len(pattern.FindAll(data, -1))

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("walk %s: %w", label, err)
	}

	return total, nil
}

// countSuffixWalk mirrors utils' countGeneratedFiles: skip .git and website.
func countSuffixWalk(root, suffix, label string) (int, error) {
	total := 0

	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			if base := filepath.Base(path); base == ".git" || base == "website" {
				return filepath.SkipDir
			}

			return nil
		}

		if strings.HasSuffix(path, suffix) {
			total++
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("walk %s: %w", label, err)
	}

	return total, nil
}

// countGoldenWalk mirrors utils' countHTMLGoldens: .golden files, visualtest
// and website excluded.
func countGoldenWalk(root string) (int, error) {
	total := 0

	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			if base := filepath.Base(path); base == ".git" || base == "website" || base == "visualtest" {
				return filepath.SkipDir
			}

			return nil
		}

		if strings.HasSuffix(path, ".golden") {
			total++
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("walk HTML goldens: %w", err)
	}

	return total, nil
}

// countVisualGoldens mirrors utils' countVisualGoldens: .png files under
// visualtest/testdata, dotted (failure-artifact) dirs skipped.
func countVisualGoldens(root string) (int, error) {
	total := 0
	visualDir := filepath.Join(root, "visualtest", "testdata")

	err := filepath.WalkDir(visualDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			if strings.HasPrefix(entry.Name(), ".") && entry.Name() != "." && path != visualDir {
				return filepath.SkipDir
			}

			return nil
		}

		if strings.HasSuffix(path, ".png") && !strings.HasSuffix(path, ".fail.png") {
			total++
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("walk visual goldens: %w", err)
	}

	return total, nil
}

// renderFunc rebuilds a site's matched text from the derived counts; prev
// carries the regex's captured groups for numbers countgen does not own
// (hand-typed totals) and exact whitespace (markdown table padding).
type renderFunc func(counts Counts, prev []string) string

// site is one count-bearing doc location: a single-match regex countgen
// rewrites with its render function.
type site struct {
	name    string
	file    string
	pattern *regexp.Regexp
	render  renderFunc
}

func featuresSites() []site {
	backtick := "`"

	totals := site{
		name:    "Totals line",
		file:    featuresDoc,
		pattern: regexp.MustCompile(`\*\*Totals:\*\* (\d+) templ components \(primitives, drift-guard verified\) \+ (\d+) recipe screens = (\d+), (\d+) icon names, (\d+) typed enums \((\d+) with ` + backtick + `IsValid\(\)` + backtick + `, (\d+) generated ` + backtick + `\*_templ\.go` + backtick + ` files`),
		render: func(c Counts, prev []string) string {
			return fmt.Sprintf(
				"**Totals:** %d templ components (primitives, drift-guard verified) + %d recipe screens = %d, %d icon names, %s typed enums (%d with %sIsValid()%s, %d generated %s*_templ.go%s files",
				c.Components, c.Recipes, c.ComponentsPlusRecipes, c.Icons,
				prev[5], c.IsValid, backtick, backtick, c.GeneratedFiles, backtick, backtick,
			)
		},
	}

	iconsRow := site{
		name:    "icons table row",
		file:    featuresDoc,
		pattern: regexp.MustCompile(`(?m)^\| (` + backtick + `icons` + backtick + `)(\s+)\| (\d+) \((\d+) icons\)`),
		render: func(c Counts, prev []string) string {
			return fmt.Sprintf("| %s%s| %d (%d icons)", prev[1], prev[2], c.PackageCounts["icons"], c.Icons)
		},
	}

	recipesRow := site{
		name:    "recipes table row",
		file:    featuresDoc,
		pattern: regexp.MustCompile(`(?m)^\| (` + backtick + `recipes` + backtick + `)(\s+)\| (\d+) screens`),
		render: func(c Counts, prev []string) string {
			return fmt.Sprintf("| %s%s| %d screens", prev[1], prev[2], c.Recipes)
		},
	}

	isValidBullet := site{
		name:    "Type Safety bullet (IsValid only; enums total is hand-typed)",
		file:    featuresDoc,
		pattern: regexp.MustCompile(`(\d+) typed string enums \((\d+) with ` + backtick + `IsValid\(\)` + backtick + ` methods \+ tests\)`),
		render: func(c Counts, prev []string) string {
			return fmt.Sprintf("%s typed string enums (%d with %sIsValid()%s methods + tests)",
				prev[1], c.IsValid, backtick, backtick)
		},
	}

	htmlBaselines := site{
		name:    "Test Coverage bullet (HTML baselines)",
		file:    featuresDoc,
		pattern: regexp.MustCompile(`in 5 packages, (\d+) baselines\)`),
		render: func(c Counts, prev []string) string {
			return fmt.Sprintf("in 5 packages, %d baselines)", c.HTMLGoldens)
		},
	}

	visualGoldens := site{
		name:    "Test Coverage bullet (visual goldens)",
		file:    featuresDoc,
		pattern: regexp.MustCompile(`; (\d+) goldens incl\. charts`),
		render: func(c Counts, prev []string) string {
			return fmt.Sprintf("; %d goldens incl. charts", c.VisualGoldens)
		},
	}

	sites := []site{totals, iconsRow, recipesRow, isValidBullet, htmlBaselines, visualGoldens}
	for _, pkg := range countedPackages {
		sites = append(sites, countTableRow(pkg))
	}

	return sites
}

// countTableRow renders one FEATURES.md package-table row. The captured gap
// preserves the markdown table's column padding byte-for-byte.
func countTableRow(pkg string) site {
	backtick := "`"

	return site{
		name:    "table row " + pkg,
		file:    featuresDoc,
		pattern: regexp.MustCompile(`(?m)^\| (` + regexp.QuoteMeta(backtick+pkg+backtick) + `)(\s+)\| (\d+)`),
		render: func(c Counts, prev []string) string {
			return fmt.Sprintf("| %s%s| %d", prev[1], prev[2], c.PackageCounts[pkg])
		},
	}
}

func skillSites() []site {
	backtick := "`"

	catalogueTotal := site{
		name:    "catalogue total line",
		file:    skillDoc,
		pattern: regexp.MustCompile(`(?m)^(\d+) components across (\d+) packages \((\w+) counted templ packages, plus`),
		render: func(c Counts, prev []string) string {
			return fmt.Sprintf("%d components across %d packages (%s counted templ packages, plus",
				c.Components, len(countedPackages)+2, numberWord(len(countedPackages)))
		},
	}

	componentsPlusIcons := site{
		name:    "components+icons line",
		file:    skillDoc,
		pattern: regexp.MustCompile(`(?m)^- (\d+) templ components \+ (\d+) icons\. Full signatures`),
		render: func(c Counts, prev []string) string {
			return fmt.Sprintf("- %d templ components + %d icons. Full signatures",
				c.Components, c.Icons)
		},
	}

	errorpageHeading := site{
		name:    "errorpage heading",
		file:    skillDoc,
		pattern: regexp.MustCompile(`(?m)^(#### ` + backtick + `errorpage` + backtick + ` — )(\d+)( components \+ )(\d+)( constructors \+ handler)`),
		render: func(c Counts, prev []string) string {
			return fmt.Sprintf("%s%d%s%d%s",
				prev[1], c.PackageCounts["errorpage"], prev[3], c.ErrorpageConstructors, prev[5])
		},
	}

	chartsHeading := site{
		name:    "charts/echarts heading",
		file:    skillDoc,
		pattern: regexp.MustCompile(`(?m)^(#### ` + backtick + `charts/echarts` + backtick + ` — )(\d+)( components \(opt-in module\))`),
		render: func(c Counts, prev []string) string {
			return fmt.Sprintf("%s%d%s", prev[1], c.PackageCounts["charts/echarts"], prev[3])
		},
	}

	sites := []site{catalogueTotal, componentsPlusIcons, errorpageHeading, chartsHeading}
	for _, pkg := range []string{"display", "forms", "feedback", "layout", "navigation", "htmx", "datastar"} {
		sites = append(sites, skillPkgHeading(pkg))
	}

	return sites
}

func skillPkgHeading(pkg string) site {
	backtick := "`"

	return site{
		name:    "heading " + pkg,
		file:    skillDoc,
		pattern: regexp.MustCompile(`(?m)^(#### ` + regexp.QuoteMeta(backtick+pkg+backtick) + ` — )(\d+)( components)`),
		render: func(c Counts, prev []string) string {
			return fmt.Sprintf("%s%d%s", prev[1], c.PackageCounts[pkg], prev[3])
		},
	}
}

func sites() []site {
	return append(featuresSites(), skillSites()...)
}

// numberWord renders small counts as words for the SKILL.md prose.
func numberWord(total int) string {
	words := []string{"zero", "one", "two", "three", "four", "five", "six",
		"seven", "eight", "nine", "ten", "eleven", "twelve", "thirteen",
		"fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen", "twenty"}

	if total < len(words) {
		return words[total]
	}

	return fmt.Sprintf("%d", total)
}

// siteMatch locates the single occurrence of a site in content.
func siteMatch(content []byte, s site) ([]int, error) {
	locations := s.pattern.FindAllIndex(content, -1)
	if len(locations) != 1 {
		return nil, fmt.Errorf(
			"%s: site %q matched %d times (want exactly 1) — structural drift, fix the doc or the site pattern by hand",
			s.file, s.name, len(locations),
		)
	}

	return locations[0], nil
}

// lineAt returns the 1-based line number of a byte offset.
func lineAt(content []byte, idx int) int {
	return strings.Count(string(content[:idx]), "\n") + 1
}

// docContent reads and caches one doc's bytes (cache holds read errors too).
type docContent struct {
	cache map[string][]byte
	errs  map[string]error
	root  string
}

func newDocContent(root string) *docContent {
	return &docContent{cache: map[string][]byte{}, errs: map[string]error{}, root: root}
}

func (d *docContent) get(name string) ([]byte, error) {
	if b, ok := d.cache[name]; ok {
		return b, nil
	}

	if err, ok := d.errs[name]; ok {
		return nil, err
	}

	b, err := os.ReadFile(filepath.Join(d.root, name))
	if err != nil {
		d.errs[name] = fmt.Errorf("read %s: %w", name, err)

		return nil, d.errs[name]
	}

	d.cache[name] = b

	return b, nil
}

// evaluate runs every site against the current docs. It yields, per site, the
// byte location of the match, its current text, and the expected text.
func evaluate(root string) (Counts, []evalResult, error) {
	counts, err := Derive(root)
	if err != nil {
		return counts, nil, err
	}

	docs := newDocContent(root)

	var results []evalResult

	for _, s := range sites() {
		content, err := docs.get(s.file)
		if err != nil {
			return counts, results, err
		}

		location, err := siteMatch(content, s)
		if err != nil {
			return counts, results, err
		}

		prev := captureGroups(content, location, s.pattern)
		results = append(results, evalResult{
			site:     s,
			content:  content,
			location: location,
			current:  string(content[location[0]:location[1]]),
			expected: s.render(counts, prev),
		})
	}

	return counts, results, nil
}

type evalResult struct {
	site     site
	content  []byte
	location []int
	current  string
	expected string
}

// Check verifies both docs against the tree. Every drift is a separate error
// naming the location and telling the user how to fix it mechanically.
func Check(root string) error {
	_, results, err := evaluate(root)
	if err != nil {
		return err
	}

	var problems []error

	for _, r := range results {
		if r.current == r.expected {
			continue
		}

		problems = append(problems, fmt.Errorf(
			"%s:%d: site %q drifts from the tree (expected %q, found %q) — run: go run ./tools/countgen/cmd/countgen -w",
			r.site.file, lineAt(r.content, r.location[0]), r.site.name, r.expected, r.current,
		))
	}

	return joinErrors(problems)
}

// Apply rewrites both docs to match the tree, returning the changed file names.
func Apply(root string) ([]string, error) {
	_, results, err := evaluate(root)
	if err != nil {
		return nil, err
	}

	updated := map[string][]byte{}
	dirty := map[string]bool{}

	for _, r := range results {
		if r.current == r.expected {
			continue
		}

		next := append([]byte{}, r.content[:r.location[0]]...)
		next = append(next, r.expected...)
		next = append(next, r.content[r.location[1]:]...)
		updated[r.site.file] = next
		dirty[r.site.file] = true
	}

	var written []string

	for _, name := range []string{featuresDoc, skillDoc} {
		if !dirty[name] {
			continue
		}

		if err := os.WriteFile(filepath.Join(root, name), updated[name], 0o644); err != nil {
			return written, fmt.Errorf("write %s: %w", name, err)
		}

		written = append(written, name)
	}

	return written, nil
}

// captureGroups extracts the site pattern's capture groups for passthrough.
func captureGroups(content []byte, location []int, pattern *regexp.Regexp) []string {
	match := pattern.FindSubmatch(content[location[0]:location[1]])
	if match == nil {
		return nil
	}

	out := make([]string, len(match))

	for i, m := range match {
		out[i] = string(m)
	}

	return out
}

func joinErrors(errs []error) error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		messages := make([]string, 0, len(errs))
		for _, err := range errs {
			messages = append(messages, err.Error())
		}

		return fmt.Errorf("%d count drifts:\n%s", len(errs), strings.Join(messages, "\n"))
	}
}
