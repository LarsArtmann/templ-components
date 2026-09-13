package contract

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// conventionPackages are the library packages the AST convention linter
// sweeps. Sources = .templ files + hand-written .go files (generated
// *_templ.go are excluded: they restate declarations and would double-count).
var conventionPackages = []string{
	"charts/echarts", "datastar", "display", "errorpage",
	"feedback", "forms", "htmx", "icons", "layout", "navigation",
}

// ratchetPackages adds utils/wire to the sweep for the IsValid ratchet —
// the wire package carries 5 of the 58 IsValid enums (Transport, Method,
// Event, ContentType, PatchMode) but no component Props.
var ratchetPackages = append([]string{"utils/wire"}, conventionPackages...)

// propsEmbedExemptions lists Props types that deliberately do NOT embed
// utils.BaseProps. Every entry needs a reason; the list must shrink, not grow.
var propsEmbedExemptions = map[string]string{
	"layout.PageProps":               "page shell predates BaseProps; carries its own ID/Class/Attr surface (documented in layout package)",
	"layout.MinimalProps":            "zero-dependency static/print shell; deliberately attribute-free (the Minimal use case), same family as PageProps",
	"forms.FormFieldProps":           "shared sub-template wrapper (label + control chrome), not a component; composed INTO real components",
	"feedback.SkeletonCardGridProps": "skeleton grid primitive renders no root attributes of its own; revisit if it grows a consumer-facing surface",
}

// packageSources parses every .go file of a package except tests. Generated
// *_templ.go files are INCLUDED deliberately: they restate the .templ type
// declarations verbatim as valid Go, so the sweep covers what templ authors
// write without needing a templ parser.
func packageSources(t *testing.T, pkgDir string) map[string]*ast.File {
	t.Helper()

	fset := token.NewFileSet()

	entries, err := os.ReadDir(pkgDir)
	if err != nil {
		t.Fatalf("read %s: %v — the convention linter owns its fixture (fail loud, never skip)", pkgDir, err)
	}

	files := map[string]*ast.File{}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}

		file, parseErr := parser.ParseFile(fset, filepath.Join(pkgDir, name), nil, 0)
		if parseErr != nil {
			t.Fatalf("parse %s/%s: %v", pkgDir, name, parseErr)
		}

		files[name] = file
	}

	return files
}

// declaredProps returns "<pkg>.<Name>" for every `type XProps struct` in the
// package, paired with whether the struct embeds utils.BaseProps.
func declaredProps(files map[string]*ast.File, pkgPath string) map[string]bool {
	props := map[string]bool{}

	for _, file := range files {
		ast.Inspect(file, func(node ast.Node) bool {
			typeSpec, ok := node.(*ast.TypeSpec)
			if !ok || !strings.HasSuffix(typeSpec.Name.Name, "Props") {
				return true
			}

			structType, isStruct := typeSpec.Type.(*ast.StructType)

			// Unexported *Props structs are internal sub-template plumbing, not
			// the public component contract — out of scope.
			if isStruct && typeSpec.Name.IsExported() {
				props[pkgPath+"."+typeSpec.Name.Name] = structEmbedsBaseProps(structType)
			}

			return true
		})
	}

	return props
}

// structEmbedsBaseProps reports whether a struct embeds utils.BaseProps (as a
// bare ident or the qualified utils.BaseProps selector).
func structEmbedsBaseProps(structType *ast.StructType) bool {
	for _, field := range structType.Fields.List {
		if len(field.Names) != 0 {
			continue // named field, not an embed
		}

		switch fieldType := field.Type.(type) {
		case *ast.Ident:
			if fieldType.Name == "BaseProps" {
				return true
			}
		case *ast.SelectorExpr:
			if x, ok := fieldType.X.(*ast.Ident); ok && x.Name == "utils" && fieldType.Sel.Name == "BaseProps" {
				return true
			}
		}
	}

	return false
}

// TestPropsEmbedBaseProps is convention check 1: every *Props struct in the
// library embeds utils.BaseProps unless it carries a reasoned exemption.
// Unlike the manual componentTypes() inventory, this sweeps the AST — a new
// component cannot silently skip the convention.
func TestPropsEmbedBaseProps(t *testing.T) {
	t.Parallel()

	for _, pkg := range conventionPackages {
		pkgDir := filepath.Join("..", "..", pkg)
		sources := packageSources(t, pkgDir)

		for prop, embeds := range declaredProps(sources, filepath.Base(pkg)) {
			if embeds {
				continue
			}

			if reason, exempt := propsEmbedExemptions[prop]; exempt {
				t.Logf("exempt: %s (%s)", prop, reason)

				continue
			}

			t.Errorf(
				"%s does not embed utils.BaseProps — every component props struct embeds it (props.Class/Attrs/ID/AriaLabel propagation); add the embed or a reasoned exemption in conventions_test.go",
				prop,
			)
		}
	}
}

// propsInventoryExemptions lists declared Props types that are
// deliberately NOT in componentTypes() because they do not embed BaseProps
// (the interface inventory only carries embedders). Mirrors
// propsEmbedExemptions plus any type the old interface test cannot take.
var propsInventoryExemptions = map[string]bool{
	"layout.PageProps":               true,
	"forms.FormFieldProps":           true,
	"feedback.SkeletonCardGridProps": true,
	"layout.MinimalProps":            true,
}

// TestPropsInventoryCoversDeclarations cross-checks the manual componentTypes()
// inventory against the AST: a Props type missing from the inventory is the
// first step of a ghost system (the interface contract silently stops
// covering the new component).
func TestPropsInventoryCoversDeclarations(t *testing.T) {
	t.Parallel()

	inventory := map[string]bool{}
	for _, propsValue := range componentTypes() {
		inventory[reflect.TypeOf(propsValue).String()] = true
	}

	for _, pkg := range conventionPackages {
		sources := packageSources(t, filepath.Join("..", "..", pkg))

		for prop := range declaredProps(sources, filepath.Base(pkg)) {
			if inventory[prop] || propsInventoryExemptions[prop] {
				continue
			}

			t.Errorf(
				"%s is declared but missing from componentTypes() in component_props_test.go — add it so the ComponentProps interface contract covers it",
				prop,
			)
		}
	}
}

// enumWithIsValid returns "<pkg>.<Type>" for every closed-set enum that has
// an IsValid function (repo convention: `func XxxIsValid(v Xxx) bool`).
func enumWithIsValid(files map[string]*ast.File, pkgPath string) map[string]bool {
	enums := map[string]bool{}

	for _, file := range files {
		ast.Inspect(file, func(node ast.Node) bool {
			funcDecl, ok := node.(*ast.FuncDecl)
			if !ok || funcDecl.Recv != nil {
				return true
			}

			name := funcDecl.Name.Name
			if !strings.HasSuffix(name, "IsValid") || name == "IsValid" {
				return true
			}

			enumName := strings.TrimSuffix(name, "IsValid")
			if enumName != "" {
				enums[pkgPath+"."+enumName] = true
			}

			return true
		})
	}

	return enums
}

// TestEnumIsValidRatchet is convention check 2 (ratchet form): the count of
// IsValid-bearing enums must never DECREASE. Every closed-set enum ships an
// IsValid method + test in the same commit (AGENTS.md); removals are
// deliberate breaking changes that must update this floor consciously.
func TestEnumIsValidRatchet(t *testing.T) {
	t.Parallel()

	const floor = 58 // README "with IsValid()" count, guarded by TestDocsCountDrift

	count := 0

	for _, pkg := range ratchetPackages {
		sources := packageSources(t, filepath.Join("..", "..", pkg))
		count += len(enumWithIsValid(sources, filepath.Base(pkg)))
	}

	if count < floor {
		t.Errorf(
			"IsValid enum count = %d, floor = %d — an IsValid method was removed; removals are breaking changes (update README + this floor together, deliberately)",
			count,
			floor,
		)
	}
}

// TestIsValidMethodsAreTested is the "no IsValid without a test" half of
// convention check 2: a package that declares IsValid methods must reference
// them from its test files (the dead-code ghost-system rule).
func TestIsValidMethodsAreTested(t *testing.T) {
	t.Parallel()

	for _, pkg := range conventionPackages {
		pkgDir := filepath.Join("..", "..", pkg)
		sources := packageSources(t, pkgDir)

		if len(enumWithIsValid(sources, pkg)) == 0 {
			continue
		}

		testReferences := packageTestReferences(t, pkgDir)

		if !testReferences {
			t.Errorf(
				"%s declares IsValid methods but no test file references IsValid — the repo rule is IsValid + test in the same commit (AGENTS.md \"no IsValid without a test\")",
				pkg,
			)
		}
	}
}

// packageTestReferences reports whether any _test.go file mentions IsValid.
func packageTestReferences(t *testing.T, pkgDir string) bool {
	t.Helper()

	files, err := filepath.Glob(filepath.Join(pkgDir, "*_test.go"))
	if err != nil {
		t.Fatalf("glob %s tests: %v", pkgDir, err)
	}

	for _, file := range files {
		content, readErr := os.ReadFile(file)
		if readErr != nil {
			t.Fatalf("read %s: %v", file, readErr)
		}

		if strings.Contains(string(content), "IsValid") {
			return true
		}
	}

	return false
}

// TestLookupMapsUseTypedEnumKeys is convention check 3: package-level lookup
// maps (name matching lookup/StyleMap/StyleSet/Map conventions) keyed by
// bare `string` while their name embeds a same-package enum type name are
// the documented map[string]X anti-pattern (AGENTS.md: typed enum keys,
// never map[string]X).
func TestLookupMapsUseTypedEnumKeys(t *testing.T) {
	t.Parallel()

	for _, pkg := range conventionPackages {
		pkgDir := filepath.Join("..", "..", pkg)
		sources := packageSources(t, pkgDir)
		enumNames := packageTypeNames(sources)

		for _, file := range sources {
			for _, flag := range stringKeyedMapsReferencingEnums(file, enumNames) {
				t.Errorf(
					"%s: %s is keyed by bare string while its name references the enum %s — use map[%s]… (typed enum keys, utils.Lookup for access)",
					pkg,
					flag.mapName,
					flag.enumName,
					flag.enumName,
				)
			}
		}
	}
}

// typedMapFlag names a lookup map keyed by bare string that references a
// same-package enum type in its own name.
type typedMapFlag struct {
	mapName  string
	enumName string
}

// packageTypeNames collects every declared type name of a package.
func packageTypeNames(files map[string]*ast.File) map[string]bool {
	names := map[string]bool{}

	for _, file := range files {
		ast.Inspect(file, func(node ast.Node) bool {
			if typeSpec, ok := node.(*ast.TypeSpec); ok {
				names[typeSpec.Name.Name] = true
			}

			return true
		})
	}

	return names
}

// stringKeyedMapsReferencingEnums returns one flag per package-level map
// composite whose key type is bare `string` and whose name both ends in
// "map" and contains a declared enum type name.
func stringKeyedMapsReferencingEnums(file *ast.File, enumNames map[string]bool) []typedMapFlag {
	var flags []typedMapFlag

	ast.Inspect(file, func(node ast.Node) bool {
		genDecl, ok := node.(*ast.GenDecl)
		if !ok {
			return true
		}

		for _, spec := range genDecl.Specs {
			valueSpec, isValue := spec.(*ast.ValueSpec)
			if !isValue || valueSpec.Values == nil {
				continue
			}

			if !mapHasKeyString(valueSpec.Values[0]) {
				continue
			}

			flags = append(flags, namesReferencingEnums(valueSpec.Names, enumNames)...)
		}

		return true
	})

	return flags
}

// mapHasKeyString reports whether the expression is a map composite keyed
// by the bare `string` type.
func mapHasKeyString(expr ast.Expr) bool {
	mapLit, isMap := expr.(*ast.CompositeLit)
	if !isMap {
		return false
	}

	mapType, isMapType := mapLit.Type.(*ast.MapType)
	if !isMapType {
		return false
	}

	keyIdent, keyIsIdent := mapType.Key.(*ast.Ident)

	return keyIsIdent && keyIdent.Name == "string"
}

// namesReferencingEnums flags identifiers ending in "map" that contain a
// declared enum type name.
func namesReferencingEnums(names []*ast.Ident, enumNames map[string]bool) []typedMapFlag {
	var flags []typedMapFlag

	for _, name := range names {
		lower := strings.ToLower(name.Name)
		if !strings.HasSuffix(lower, "map") {
			continue
		}

		for enumName := range enumNames {
			if enumName != "string" && strings.Contains(lower, strings.ToLower(enumName)) {
				flags = append(flags, typedMapFlag{mapName: name.Name, enumName: enumName})
			}
		}
	}

	return flags
}
