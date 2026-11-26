//nolint:err113,gosec // CLI plumbing: messages are printed verbatim (never matched); every path is the consumer's own tree
package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/larsartmann/templ-components/utils"
)

// `tc css` keeps the stylesheet inputs templ-components ships in step with the
// version a consumer's go.mod pins, and verifies that the consumer's Tailwind
// build reaches them.
//
// The library delivers markup as Go source and its styling through the
// consumer's Tailwind build: a class corpus the scanner reads plus hand-written
// CSS that build imports. A stale or missing copy renders a page unstyled and
// fails nowhere, so every copy `tc css sync` writes carries a stamp naming the
// library version it came from, and `tc css check` fails on version skew, on a
// missing or drifted copy, and on an entry that never reaches them.

const (
	// libraryModulePath is the module that ships the stylesheet inputs.
	libraryModulePath = "github.com/larsartmann/templ-components"
	// stampMarker opens the stamp line of every library-owned copy.
	stampMarker = "tc:sync templ-components "
	// starterEntryName is the consumer-owned Tailwind entry `tc init` scaffolds.
	starterEntryName = "app.css"
	// corpusName is the class corpus a Tailwind entry scans in the synced layout.
	corpusName = "templ-components-classes.txt"
	// stylesheetName is the library's stylesheet: one file that imports the rest.
	stylesheetName = "templ-components.css"
	// customName is the hand-written component CSS.
	customName = "custom.css"
	// vendoredStylesheetPath is where `go mod vendor` ships the stylesheet, relative
	// to the library directory (layout/ embeds it so Go carries it).
	vendoredStylesheetPath = "layout/styles/" + stylesheetName
	// vendoredLibraryPath is the library directory inside a vendor tree, relative to
	// a Tailwind entry that lives one directory below the project root.
	vendoredLibraryPath = "../vendor/" + libraryModulePath
	// cssImportDepth bounds how far the wiring check follows local @import chains.
	cssImportDepth = 3
	// goModSearchDepth bounds the upward search for go.mod.
	goModSearchDepth = 32
	// syncCommand writes the library files of the version go.mod pins, for builds
	// that compile CSS without a vendor tree. The explicit -mod=mod matters:
	// cmd/tc is never vendored (a dependency's command packages are not part of
	// the main module), so a project with a vendor/ directory cannot run it in
	// the default vendor mode.
	syncCommand = "GOFLAGS=-mod=mod go run " + libraryModulePath + "/cmd/tc css sync"
)

// libraryPinModules are the modules whose go.mod version identifies the
// stylesheet inputs: the root module ships them, and the sibling pins ride
// the same release, so an icons-only consumer still resolves a version.
//
//nolint:gochecknoglobals // static table, the only configuration of tc css
var libraryPinModules = []string{
	libraryModulePath,
	libraryModulePath + "/utils",
	libraryModulePath + "/icons",
}

// libraryFile is one stylesheet input the library owns inside a consumer's CSS
// directory: `tc css sync` writes it and `tc css check` verifies it byte for byte.
type libraryFile struct {
	name string
	// required marks an input every component consumer needs: check fails when it is absent.
	required bool
	// plainText marks the class corpus, whose stamp is a `#` line instead of a comment block.
	plainText bool
}

// libraryFiles lists the inputs `tc css sync` writes, in output order.
//
//nolint:gochecknoglobals,exhaustruct_v5 // static table; the optional inputs leave required at its zero value
var libraryFiles = []libraryFile{
	{name: corpusName, required: true, plainText: true},
	{name: stylesheetName, required: true},
	{name: customName, required: true},
	{name: "templ-components-theme.css"},
	{name: "safelist.css"},
}

// cssOptions carries the flags of `tc css sync` and `tc css check`.
type cssOptions struct {
	// dir overrides where the library-owned files live (default: derived from the entry).
	dir string
	// force lets sync replace a differing copy that carries no stamp.
	force bool
	// css names one Tailwind entry to inspect instead of discovering every entry.
	css string
	// scaffold is the `tc init` contract: fill gaps only, never replace a file.
	scaffold bool
}

// cssWiring reports which library inputs the Tailwind entries reach.
type cssWiring struct {
	// corpus is set when an entry scans the class corpus.
	corpus bool
	// custom is set when an entry imports the library's stylesheet or component CSS.
	custom bool
	// libraryDirectory is the library tree an entry scans (a vendor tree or a
	// checkout); it replaces the corpus, because the tree holds the real sources.
	libraryDirectory string
	// vendoredVersions are the library versions vendor/modules.txt records for
	// the scanned tree (empty outside a vendor tree).
	vendoredVersions []string
	// missingVendor is set when an entry scans the library's vendor tree and that
	// tree does not exist: `go mod vendor` did not run where CSS is compiled.
	missingVendor bool
	// curatedVendorFiles counts @source lines that name single files of the
	// vendor tree: a hand-curated list ships the next component unstyled.
	curatedVendorFiles int
}

// syncAction is what cssSync does with one library file in the target directory.
type syncAction struct {
	verb    string
	target  string
	version string
	content []byte
}

// cssDirective is one @import or @source of a stylesheet that names a local file.
type cssDirective struct {
	isImport bool
	target   string
}

func cmdCSS(args []string) {
	mode, opts, err := parseCSSArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tc css: %v\n\n", err)
		cssUsage()

		os.Exit(2)
	}

	switch mode {
	case "sync":
		if err := cssSync(".", opts, os.Stdout); err != nil {
			failf("css sync: %v", err)
		}
	case "check":
		if failures := cssCheck(".", opts, os.Stdout); failures > 0 {
			os.Exit(1)
		}
	}
}

func cssUsage() {
	fmt.Fprintln(os.Stderr, `Usage:
  tc css sync [--dir DIR] [--force]   Write the library's files for a CSS build without a vendor tree.
  tc css check [--dir DIR] [--css PATH]
                                      Fail when Tailwind cannot reach the library.

  When go mod vendor runs before Tailwind, no sync is needed. The entry reads
  the library straight from the vendor tree:
    @source "../vendor/github.com/larsartmann/templ-components";
    @import "../vendor/github.com/larsartmann/templ-components/layout/styles/templ-components.css";

  sync is the fallback for a CSS build that has no Go toolchain and no vendor
  tree: it writes the class corpus, templ-components.css, custom.css, the theme
  and the safelist into the directory of your Tailwind entry (or DIR), stamped
  with the library version, then prints the lines that entry still needs.
  check verifies either layout and exits 1 on any finding, so CI can gate on it.

Run tc at the version go.mod pins (this form also works in vendored projects):
  `+syncCommand)
}

func parseCSSArgs(args []string) (string, cssOptions, error) {
	opts := defaultCSSOptions()

	if len(args) == 0 {
		return "", opts, errors.New("expected `sync` or `check`")
	}

	mode := args[0]
	if mode != "sync" && mode != "check" {
		return "", opts, fmt.Errorf("unknown subcommand %q", mode)
	}

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--force":
			opts.force = true
		case "--dir", "--css":
			if i+1 >= len(args) {
				return "", opts, fmt.Errorf("%s requires a value", args[i])
			}

			i++

			if args[i-1] == "--dir" {
				opts.dir = args[i]
			} else {
				opts.css = args[i]
			}
		default:
			return "", opts, fmt.Errorf("unknown flag %q", args[i])
		}
	}

	return mode, opts, nil
}

// defaultCSSOptions is the flag configuration both css commands start from.
func defaultCSSOptions() cssOptions {
	return cssOptions{dir: "", force: false, css: "", scaffold: false}
}

// cssSync makes the CSS directory hold the stylesheet inputs of the library
// version go.mod pins. With scaffold set it only fills gaps (`tc init`) and
// tolerates a project that does not require the library yet. The root is the
// directory the go.mod search and the entry discovery start from.
func cssSync(root string, opts cssOptions, out io.Writer) error {
	if err := syncPinGate(root, opts, out); err != nil {
		return err
	}

	directory, entries, err := resolveCSSDirectory(root, opts)
	if err != nil {
		return err
	}

	actions, err := planLibraryFiles(directory, opts)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(directory, 0o750); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	for _, action := range actions {
		if action.content != nil {
			if err := os.WriteFile(action.target, action.content, 0o600); err != nil {
				return fmt.Errorf("write %s: %w", action.target, err)
			}
		}

		note := ""
		if action.version != "" {
			note = " (was stamped " + action.version + ")"
		}

		fmt.Fprintf(out, "tc: %s %s%s\n", action.verb, action.target, note)
	}

	hintWiring(out, entries, directory)

	return nil
}

// syncPinGate decides whether a sync may write: go.mod must pin the version
// this tc embeds, or replace the library with a directory (the tc built from
// that tree carries that tree's files). Scaffolding tolerates a missing pin.
func syncPinGate(root string, opts cssOptions, out io.Writer) error {
	pin, err := libraryPin(root)

	switch {
	case err != nil && !opts.scaffold:
		return err
	case err != nil:
		fmt.Fprintf(out, "tc css: %v; scaffolding anyway, run `tc css check` after `go get`\n", err)
	case libraryReplacedLocally(root):
		fmt.Fprintf(out, "tc css: go.mod replaces %s with a local directory, so the version check is skipped\n",
			libraryModulePath)
	case pin != libraryVersion():
		return errPinMismatch(pin)
	}

	return nil
}

// planLibraryFiles decides per library file whether dir already holds the
// pinned bytes and what to do otherwise. A copy that carries a stamp belongs
// to `tc css sync` and is refreshed; an unstamped copy that differs may be the
// consumer's own (an older scaffold, hand edits) and is only replaced with
// force. A refusal aborts the whole plan, so a run is all or nothing.
func planLibraryFiles(dir string, opts cssOptions) ([]syncAction, error) {
	actions := make([]syncAction, 0, len(libraryFiles))

	for _, file := range libraryFiles {
		content, err := syncedContent(file)
		if err != nil {
			return nil, err
		}

		action := syncAction{
			verb: "wrote", target: filepath.Join(dir, file.name), version: "", content: content,
		}

		existing, readErr := os.ReadFile(action.target)

		switch {
		case errors.Is(readErr, os.ErrNotExist):
			// Nothing to compare against: write.
		case readErr != nil:
			return nil, fmt.Errorf("read %s: %w", action.target, readErr)
		default:
			stamped, body := splitStamp(existing)

			canonical, err := libraryContent(file.name)
			if err != nil {
				return nil, err
			}

			switch {
			case bytes.Equal(body, canonical):
				action.verb, action.content = "kept", nil
			case opts.scaffold:
				action.verb, action.content = "skipped", nil
			case opts.force || stamped != "":
				action.verb, action.version = "refreshed", stamped
			default:
				return nil, fmt.Errorf(
					"refusing to overwrite %s: it differs from %s and carries no tc css stamp, "+
						"so it may hold your own rules. Move those to a stylesheet of your own, then re-run with --force",
					action.target, libraryVersion(),
				)
			}
		}

		actions = append(actions, action)
	}

	return actions, nil
}

// cssCheck fails on version skew, on a stale vendor tree, and when no Tailwind
// entry reaches the library. Two layouts pass, and nothing else has to:
//
//   - vendored: an entry scans the library tree `go mod vendor` shipped and
//     imports the stylesheet it carries (layout/styles/). Synced files are not
//     read in this layout, so they are not checked.
//   - synced: an entry scans the synced class corpus and imports the synced
//     stylesheet, and every synced file matches the pinned version.
//
// It returns the failure count so the caller maps it to an exit code. The root
// is the directory the go.mod search and the entry discovery start from.
func cssCheck(root string, opts cssOptions, out io.Writer) int {
	fmt.Fprintln(out, "tc css check: templ-components stylesheet delivery")
	fmt.Fprintln(out)

	failures := checkPin(root, out)

	directory, entries, err := resolveCSSDirectory(root, opts)
	if err != nil {
		failures += cssReport(out, "Tailwind entry found", false, err.Error())

		return finishCSSCheck(out, failures)
	}

	wiring := entriesWiring(entries)

	switch {
	case wiring.libraryDirectory != "":
		failures += checkVendorVersion(out, root, wiring)
	case wiring.missingVendor || wiring.curatedVendorFiles > 0:
		break // The entry asked for the vendor tree; synced files are not what it reads.
	default:
		for _, file := range libraryFiles {
			failures += checkLibraryFile(out, file, filepath.Join(directory, file.name))
		}
	}

	failures += checkWiring(out, entries, directory, wiring)

	return finishCSSCheck(out, failures)
}

// checkPin verifies that go.mod requires exactly the library version this tc
// embeds, since the files would otherwise not match the module.
func checkPin(root string, out io.Writer) int {
	pin, err := libraryPin(root)

	switch {
	case err != nil:
		return cssReport(out, "go.mod requires templ-components", false, err.Error())
	case libraryReplacedLocally(root):
		return cssReport(
			out,
			"go.mod replaces templ-components with a local directory (version check skipped: that tree supplies the files)",
			true,
			"",
		)
	case pin != libraryVersion():
		return cssReport(out, "go.mod pin equals this tc ("+pin+" vs "+libraryVersion()+")", false,
			errPinMismatch(pin).Error())
	}

	return cssReport(out, "go.mod pin equals this tc ("+pin+")", true, "")
}

func finishCSSCheck(out io.Writer, failures int) int {
	fmt.Fprintln(out)

	if failures > 0 {
		fmt.Fprintf(out, "result: %d check(s) FAILED; see the fix hints above\n", failures)
	} else {
		fmt.Fprintln(out, "result: all checks passed")
	}

	return failures
}

func checkLibraryFile(out io.Writer, file libraryFile, target string) int {
	label := file.name + " (" + target + ")"

	existing, err := os.ReadFile(target)
	if errors.Is(err, os.ErrNotExist) {
		if !file.required {
			return 0
		}

		return cssReport(out, label+" present", false, "run `tc css sync`")
	}

	if err != nil {
		return cssReport(out, label+" readable", false, err.Error())
	}

	stamped, body := splitStamp(existing)

	canonical, err := libraryContent(file.name)
	if err != nil {
		return cssReport(out, label, false, err.Error())
	}

	switch {
	case bytes.Equal(body, canonical):
		return cssReport(out, label+" matches "+libraryVersion(), true, "")
	case stamped != "":
		return cssReport(out, label+" differs from "+libraryVersion()+" (copy is stamped "+stamped+")", false,
			"run `tc css sync`")
	}

	return cssReport(out, label+" differs from "+libraryVersion()+" and carries no tc css stamp", false,
		"if it holds your own rules, move them to a stylesheet of your own; then run `tc css sync --force`")
}

// checkVendorVersion verifies that a scanned vendor tree carries the library
// version go.mod pins: `go mod vendor` has to run after every bump, and a CSS
// build alone does not notice a stale tree.
func checkVendorVersion(out io.Writer, root string, wiring cssWiring) int {
	if len(wiring.vendoredVersions) == 0 || libraryReplacedLocally(root) {
		return 0
	}

	pin, err := libraryPin(root)
	if err != nil {
		return 0
	}

	for _, version := range wiring.vendoredVersions {
		if version != pin {
			return cssReport(
				out,
				"vendor tree matches the go.mod pin ("+strings.Join(wiring.vendoredVersions, ", ")+" vs "+pin+")",
				false,
				"run `go mod vendor`: the vendor tree was not refreshed after the bump",
			)
		}
	}

	return cssReport(out, "vendor tree matches the go.mod pin ("+pin+")", true, "")
}

// entriesWiring merges what the Tailwind entries reach. An entry that skips
// the library does not fail the check: a project may compile stylesheets that
// never render it.
func entriesWiring(entries []string) cssWiring {
	merged := emptyWiring()

	for _, entry := range entries {
		wiring, err := entryWiring(entry)
		if err != nil {
			continue
		}

		merged.corpus = merged.corpus || wiring.corpus
		merged.custom = merged.custom || wiring.custom
		merged.missingVendor = merged.missingVendor || wiring.missingVendor
		merged.curatedVendorFiles += wiring.curatedVendorFiles

		if merged.libraryDirectory == "" {
			merged.libraryDirectory = wiring.libraryDirectory
			merged.vendoredVersions = wiring.vendoredVersions
		}
	}

	return merged
}

func emptyWiring() cssWiring {
	return cssWiring{
		corpus: false, custom: false, libraryDirectory: "", vendoredVersions: nil, missingVendor: false,
		curatedVendorFiles: 0,
	}
}

// checkWiring verifies that a Tailwind entry reaches the library's classes and
// its stylesheet, and names one line per layout for whatever is missing.
func checkWiring(out io.Writer, entries []string, directory string, wiring cssWiring) int {
	failures := 0

	switch {
	case wiring.libraryDirectory != "":
		failures += cssReport(out, "a Tailwind entry scans the library sources ("+wiring.libraryDirectory+")", true, "")
	case wiring.corpus:
		failures += cssReport(out, "a Tailwind entry scans "+corpusName, true, "")
	case wiring.missingVendor:
		failures += cssReport(out, "a Tailwind entry scans the library's vendor tree", false,
			"the entry names the vendor tree but it does not exist here; run `go mod vendor` before Tailwind")
	case wiring.curatedVendorFiles > 0:
		failures += cssReport(out, "a Tailwind entry reaches the library's classes", false,
			strconv.Itoa(wiring.curatedVendorFiles)+" single library files are listed by hand; a curated list ships "+
				"the next component unstyled. Scan the directory instead: "+
				vendoredSourceLine(entryDirectoryOf(entries)))
	default:
		failures += cssReport(out, "a Tailwind entry reaches the library's classes", false,
			strings.Join(classesHints(entries, directory), "; or "))
	}

	if wiring.custom {
		failures += cssReport(out, "a Tailwind entry imports the library stylesheet", true, "")
	} else {
		failures += cssReport(out, "a Tailwind entry imports the library stylesheet", false,
			strings.Join(stylesheetHints(entries, directory), "; or "))
	}

	return failures
}

// classesHints names the line that lets an entry reach the library's classes,
// once per layout.
func classesHints(entries []string, directory string) []string {
	entryDirectory := entryDirectoryOf(entries)

	return []string{
		syncedSourceLine(entryDirectory, directory) + " (after `tc css sync`)",
		vendoredSourceLine(entryDirectory) + " (when `go mod vendor` runs before Tailwind)",
	}
}

// stylesheetHints names the line that lets an entry import the library
// stylesheet, once per layout.
func stylesheetHints(entries []string, directory string) []string {
	entryDirectory := entryDirectoryOf(entries)

	return []string{
		syncedImportLine(entryDirectory, directory) + " (after `tc css sync`)",
		vendoredImportLine(entryDirectory) + " (when `go mod vendor` runs before Tailwind)",
	}
}

// hintWiring prints the lines the Tailwind entry still needs after a sync.
func hintWiring(out io.Writer, entries []string, directory string) {
	wiring := emptyWiring()
	entry := ""

	if len(entries) > 0 {
		entry = entries[0]

		resolved, err := entryWiring(entry)
		if err != nil {
			return
		}

		wiring = resolved
	}

	entryDirectory := entryDirectoryOf(entries)

	lines := []string{}
	if !wiring.corpus {
		lines = append(lines, syncedSourceLine(entryDirectory, directory))
	}

	if !wiring.custom {
		lines = append(lines, syncedImportLine(entryDirectory, directory))
	}

	if len(lines) == 0 {
		fmt.Fprintln(out, "tc css: your Tailwind entry already reaches every library file")

		return
	}

	if entry == "" {
		fmt.Fprintln(out, "tc css: no Tailwind entry found; add these lines to yours (paths relative to that file):")
	} else {
		fmt.Fprintf(out, "tc css: %s does not reach every library file yet; add these lines to it:\n", entry)
	}

	for _, line := range lines {
		fmt.Fprintf(out, "  %s\n", line)
	}

	fmt.Fprintln(out, "then verify with: tc css check")
}

func entryDirectoryOf(entries []string) string {
	if len(entries) == 0 {
		return ""
	}

	return filepath.Dir(entries[0])
}

func syncedSourceLine(entryDirectory, directory string) string {
	return fmt.Sprintf("@source %q;", entryRelative(entryDirectory, directory)+corpusName)
}

func syncedImportLine(entryDirectory, directory string) string {
	return fmt.Sprintf("@import %q;", entryRelative(entryDirectory, directory)+stylesheetName)
}

func vendoredSourceLine(entryDirectory string) string {
	return fmt.Sprintf("@source %q;", vendoredLibraryRelative(entryDirectory))
}

func vendoredImportLine(entryDirectory string) string {
	return fmt.Sprintf("@import %q;", vendoredLibraryRelative(entryDirectory)+"/"+vendoredStylesheetPath)
}

// entryRelative is the synced library directory as a CSS path relative to the
// entry, with a trailing slash.
func entryRelative(entryDirectory, directory string) string {
	relative, err := filepath.Rel(entryDirectory, directory)
	if err != nil || relative == "." {
		return "./"
	}

	return cssRelative(relative) + "/"
}

// vendoredLibraryRelative is the library's vendor directory as a CSS path
// relative to the entry. When the project has no vendor tree yet, it is the
// conventional path of an entry one directory below the project root, which
// is where `go mod vendor` will create it.
func vendoredLibraryRelative(entryDirectory string) string {
	directory := entryDirectory

	for range goModSearchDepth {
		candidate := filepath.Join(directory, "vendor", libraryModulePath)
		if isLibraryDirectory(candidate) {
			relative, err := filepath.Rel(entryDirectory, candidate)
			if err == nil {
				return cssRelative(relative)
			}
		}

		parent := filepath.Dir(directory)
		if parent == directory {
			break
		}

		directory = parent
	}

	return vendoredLibraryPath
}

// cssRelative renders a relative path the way a CSS directive wants it:
// forward slashes, and a leading "./" for plain sibling names.
func cssRelative(relative string) string {
	relative = filepath.ToSlash(filepath.Clean(relative))
	if !strings.HasPrefix(relative, ".") {
		relative = "./" + relative
	}

	return relative
}

// cssReport prints one pass/fail line and returns 1 on failure.
func cssReport(out io.Writer, name string, ok bool, hint string) int {
	state := "pass"
	if !ok {
		state = "FAIL"
	}

	fmt.Fprintf(out, "  [%s] %s\n", state, name)

	if !ok && hint != "" {
		fmt.Fprintf(out, "         fix: %s\n", hint)
	}

	if !ok {
		return 1
	}

	return 0
}

// resolveCSSDirectory returns where the library-owned files live and the
// Tailwind entries to inspect. Without --dir the directory is derived from the
// first entry: the folder holding a library file the entry references, else
// the folder of the entry itself.
func resolveCSSDirectory(root string, opts cssOptions) (string, []string, error) {
	entries, err := tailwindEntries(root, opts)
	if err != nil {
		return "", nil, err
	}

	if opts.dir != "" {
		absolute, absErr := filepath.Abs(opts.dir)
		if absErr != nil {
			return "", entries, fmt.Errorf("resolve %s: %w", opts.dir, absErr)
		}

		return absolute, entries, nil
	}

	if len(entries) == 0 {
		return "", nil, errors.New(
			"no Tailwind entry found (a css file importing tailwindcss, searched from the current directory, " +
				"skipping node_modules/vendor/.git); run `tc init` first or pass --dir DIR",
		)
	}

	absolute, err := filepath.Abs(libraryDirectory(entries[0]))
	if err != nil {
		return "", entries, fmt.Errorf("resolve %s: %w", entries[0], err)
	}

	return absolute, entries, nil
}

// libraryDirectory derives the library file directory from what an entry
// references: a referenced file whose content is a library file, else a
// referenced path that names one (not synced yet), else the entry's own
// directory. Matching content first keeps a project's own custom.css from
// being mistaken for the library's.
func libraryDirectory(entry string) string {
	chain := stylesheetChain(entry)

	for _, stylesheet := range chain {
		for _, reference := range stylesheetReferences(stylesheet) {
			if fileMatches(reference, isLibraryFile(filepath.Base(reference))) {
				return filepath.Dir(reference)
			}
		}
	}

	for _, stylesheet := range chain {
		for _, reference := range stylesheetReferences(stylesheet) {
			if isLibraryFile(filepath.Base(reference)) != "" {
				return filepath.Dir(reference)
			}
		}
	}

	return filepath.Dir(entry)
}

// isLibraryFile maps a canonical file name to itself, or "" for anything else.
func isLibraryFile(name string) string {
	for _, file := range libraryFiles {
		if file.name == name {
			return name
		}
	}

	return ""
}

// stylesheetReferences lists the local paths a stylesheet's @import and
// @source directives name, whether or not the files exist.
func stylesheetReferences(stylesheet string) []string {
	directives := cssDirectives(textFile(stylesheet))
	references := make([]string, 0, len(directives))

	for _, directive := range directives {
		path := directive.target
		if !filepath.IsAbs(path) {
			path = filepath.Join(filepath.Dir(stylesheet), filepath.FromSlash(directive.target))
		}

		references = append(references, path)
	}

	return references
}

// tailwindEntries returns the stylesheets the consumer's build compiles: the
// --css override, or every css file under root that imports tailwindcss.
func tailwindEntries(root string, opts cssOptions) ([]string, error) {
	if opts.css != "" {
		if _, err := os.Stat(opts.css); err != nil {
			return nil, fmt.Errorf("read %s: %w", opts.css, err)
		}

		return []string{opts.css}, nil
	}

	entries := []string{}

	err := filepath.WalkDir(root, func(path string, dirent os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil //nolint:nilerr // unreadable entries are not this check's business
		}

		if dirent.IsDir() && cssDirPruned(root, path, dirent.Name()) {
			return filepath.SkipDir
		}

		if dirent.IsDir() || !strings.HasSuffix(path, ".css") {
			return nil
		}

		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil //nolint:nilerr // an unreadable stylesheet is not an entry either
		}

		text := stripCSSComments(string(content))
		if strings.Contains(text, `@import "tailwindcss`) || strings.Contains(text, `@import 'tailwindcss`) {
			entries = append(entries, path)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk %s: %w", root, err)
	}

	return entries, nil
}

// stylesheetChain returns the entry followed by every local stylesheet it
// imports, transitively, up to cssImportDepth levels.
func stylesheetChain(entry string) []string {
	chain := []string{entry}

	for range cssImportDepth {
		added := false

		for _, stylesheet := range slices.Clone(chain) {
			content, err := os.ReadFile(stylesheet)
			if err != nil {
				continue
			}

			for _, directive := range cssDirectives(string(content)) {
				if !directive.isImport {
					continue
				}

				target, ok := resolveCSSReference(stylesheet, directive.target, true)
				if !ok || slices.Contains(chain, target) {
					continue
				}

				chain = append(chain, target)
				added = true
			}
		}

		if !added {
			break
		}
	}

	return chain
}

// entryWiring reports what an entry reaches, looking through the stylesheets
// it imports: the synced class corpus or the library's own tree for classes,
// and the library stylesheet (or its component CSS) for the hand-written rules.
func entryWiring(entry string) (cssWiring, error) {
	wiring := emptyWiring()

	for _, stylesheet := range stylesheetChain(entry) {
		if err := wiring.read(stylesheet); err != nil {
			return wiring, err
		}
	}

	return wiring, nil
}

// read folds one stylesheet's @import and @source directives into the wiring.
func (wiring *cssWiring) read(stylesheet string) error {
	content, err := os.ReadFile(stylesheet)
	if err != nil {
		return fmt.Errorf("read %s: %w", stylesheet, err)
	}

	for _, directive := range cssDirectives(string(content)) {
		path := lexicalReference(stylesheet, directive.target)

		if directive.isImport {
			wiring.custom = wiring.custom || fileMatches(path, stylesheetName) || fileMatches(path, customName)

			continue
		}

		if fileMatches(path, corpusName) {
			wiring.corpus = true

			continue
		}

		if isVendoredLibraryPath(path) {
			info, statErr := os.Stat(path)
			if statErr != nil {
				wiring.missingVendor = true

				continue
			}

			if !info.IsDir() {
				wiring.curatedVendorFiles++

				continue
			}
		}

		if isLibraryDirectory(path) {
			wiring.corpus = true
			wiring.libraryDirectory = path
			wiring.vendoredVersions = vendoredVersions(path)
		}
	}

	return nil
}

// isVendoredLibraryPath reports whether a path lies in the library's vendor
// directory, whatever exists there. The surrounding slashes make the test
// work for a relative path too (entries are discovered relative to the project
// root) and stop a sibling module with a longer name from matching.
func isVendoredLibraryPath(path string) bool {
	normalized := "/" + filepath.ToSlash(path) + "/"

	return strings.Contains(normalized, "/vendor/"+libraryModulePath+"/")
}

// lexicalReference is the path a directive names, relative to the stylesheet
// that carries it. A glob narrows to the directory it starts from, which is
// what identifies a library tree.
func lexicalReference(stylesheet, target string) string {
	if cut := strings.IndexAny(target, "*?[{"); cut >= 0 {
		target = target[:cut]
	}

	if !filepath.IsAbs(target) {
		target = filepath.Join(filepath.Dir(stylesheet), filepath.FromSlash(target))
	}

	return filepath.Clean(target)
}

// isLibraryDirectory reports whether path is the library's own tree: a
// checkout (its go.mod names the library module) or a vendored copy (the
// enclosing vendor/modules.txt lists it). Directory names alone would do:
// consumers have `layout/` and `utils/` directories of their own.
func isLibraryDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	}

	content, err := os.ReadFile(filepath.Join(path, "go.mod"))
	if err == nil {
		for line := range strings.SplitSeq(string(content), "\n") {
			if strings.TrimSpace(line) == "module "+libraryModulePath {
				return true
			}
		}
	}

	return len(vendoredVersions(path)) > 0
}

// vendoredVersions are the library versions vendor/modules.txt records for a
// vendored library directory; empty when the directory is not in a vendor tree.
func vendoredVersions(path string) []string {
	directory := path

	for range goModSearchDepth {
		if filepath.Base(directory) == "vendor" {
			return moduleVersions(filepath.Join(directory, "modules.txt"))
		}

		parent := filepath.Dir(directory)
		if parent == directory {
			break
		}

		directory = parent
	}

	return nil
}

// moduleVersions lists the distinct versions vendor/modules.txt records for
// the library and its sibling modules.
func moduleVersions(modulesTxt string) []string {
	versions := []string{}

	for line := range strings.SplitSeq(textFile(modulesTxt), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 3 || fields[0] != "#" || !isLibraryModule(fields[1]) {
			continue
		}

		if !slices.Contains(versions, fields[2]) {
			versions = append(versions, fields[2])
		}
	}

	return versions
}

func isLibraryModule(module string) bool {
	return module == libraryModulePath || strings.HasPrefix(module, libraryModulePath+"/")
}

// textFile reads a file, or "" when it cannot be read.
func textFile(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return string(content)
}

// cssDirectives lists the @import and @source directives of a stylesheet that
// name a local file. Comments, exclusions (`@source not`), inline candidates,
// the tailwindcss module itself and remote urls are left out.
func cssDirectives(text string) []cssDirective {
	directives := []cssDirective{}

	for line := range strings.SplitSeq(stripCSSComments(text), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 || fields[1] == "not" || strings.HasPrefix(fields[1], "inline(") {
			continue
		}

		directive := cssDirective{isImport: fields[0] == "@import", target: ""}
		if fields[0] != "@import" && fields[0] != "@source" {
			continue
		}

		directive.target = directiveTarget(line)
		if namesLocalFile(directive) {
			directives = append(directives, directive)
		}
	}

	return directives
}

// directiveTarget extracts the path of one directive: the first quoted string,
// else the bare token after the keyword.
func directiveTarget(line string) string {
	if quote := strings.IndexAny(line, `"'`); quote >= 0 {
		if end := strings.IndexByte(line[quote+1:], line[quote]); end >= 0 {
			return line[quote+1 : quote+1+end]
		}
	}

	return strings.TrimSuffix(strings.Fields(line)[1], ";")
}

func namesLocalFile(directive cssDirective) bool {
	if directive.target == "" || strings.Contains(directive.target, "url(") ||
		strings.Contains(directive.target, "//") {
		return false
	}

	if !directive.isImport {
		return true
	}

	return directive.target != "tailwindcss" && !strings.HasPrefix(directive.target, "tailwindcss/")
}

// resolveCSSReference maps a directive target to the file it names, relative
// to the stylesheet that carries it (Tailwind resolves every @source and
// @import path against the stylesheet's own directory, never the CWD).
func resolveCSSReference(stylesheet, target string, isImport bool) (string, bool) {
	path := target
	if !filepath.IsAbs(path) {
		path = filepath.Join(filepath.Dir(stylesheet), filepath.FromSlash(target))
	}

	if _, err := os.Stat(path); err == nil {
		return path, true
	}

	if isImport {
		if _, err := os.Stat(path + ".css"); err == nil {
			return path + ".css", true
		}
	}

	return "", false
}

// fileMatches reports whether path holds exactly the canonical library file.
func fileMatches(path, name string) bool {
	if name == "" {
		return false
	}

	existing, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	canonical, err := libraryContent(name)
	if err != nil {
		return false
	}

	_, body := splitStamp(existing)

	return bytes.Equal(body, canonical)
}

// stripCSSComments removes /* */ comments, honoring quoted strings: a glob
// such as "../**/*.templ" contains both comment delimiters and must survive.
func stripCSSComments(text string) string {
	var builder strings.Builder

	builder.Grow(len(text))

	var quote byte

	for i := 0; i < len(text); i++ {
		current := text[i]

		switch {
		case quote != 0:
			if current == quote {
				quote = 0
			}

			builder.WriteByte(current)
		case current == '"' || current == '\'':
			quote = current

			builder.WriteByte(current)
		case strings.HasPrefix(text[i:], "/*"):
			end := strings.Index(text[i:], "*/")
			if end < 0 {
				return builder.String()
			}

			i += end + len("*/") - 1
		default:
			builder.WriteByte(current)
		}
	}

	return builder.String()
}

// libraryPin returns the templ-components version the nearest go.mod requires.
func libraryPin(start string) (string, error) {
	goModPath, found := findGoMod(start)
	if !found {
		return "", fmt.Errorf(
			"no go.mod found from %s upward; run inside the Go module that requires %s",
			start, libraryModulePath,
		)
	}

	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", goModPath, err)
	}

	for _, module := range libraryPinModules {
		if version, found := requireVersion(string(content), module); found {
			return version, nil
		}
	}

	return "", fmt.Errorf(
		"%s does not require %s; run `go get %s` first", goModPath, libraryModulePath, libraryModulePath,
	)
}

// libraryReplacedLocally reports whether the nearest go.mod replaces the
// library with a directory. Go then records no version for it (tidy writes
// the zero pseudo-version), and the tc built from the replaced tree carries
// the files of exactly that tree, so a version comparison would only refuse
// a correct setup.
func libraryReplacedLocally(start string) bool {
	goModPath, found := findGoMod(start)
	if !found {
		return false
	}

	content, err := os.ReadFile(goModPath)
	if err != nil {
		return false
	}

	for raw := range strings.SplitSeq(string(content), "\n") {
		fields := strings.Fields(stripLineComment(raw))

		singleLine := len(fields) == 4 && fields[0] == "replace" && fields[1] == libraryModulePath &&
			fields[2] == "=>"
		inBlock := len(fields) == 3 && fields[0] == libraryModulePath && fields[1] == "=>"

		switch {
		case singleLine && isLocalReplaceTarget(fields[3]), inBlock && isLocalReplaceTarget(fields[2]):
			return true
		}
	}

	return false
}

// isLocalReplaceTarget reports whether a replace target is a directory on
// disk rather than a module version.
func isLocalReplaceTarget(target string) bool {
	return filepath.IsAbs(target) || strings.HasPrefix(target, "./") || strings.HasPrefix(target, "../")
}

// findGoMod locates the nearest go.mod at or above start.
func findGoMod(start string) (string, bool) {
	directory := start

	for range goModSearchDepth {
		candidate := filepath.Join(directory, "go.mod")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}

		parent := filepath.Dir(directory)
		if parent == directory {
			break
		}

		directory = parent
	}

	return "", false
}

func libraryVersion() string {
	return "v" + utils.Version
}

// errPinMismatch explains how to reach a tc that matches the pinned module.
func errPinMismatch(pin string) error {
	return fmt.Errorf(
		"go.mod pins %s but this tc is %s, so its files would not match your module. Run: %s",
		pin, libraryVersion(), syncCommand,
	)
}

// requireVersion extracts the version of modulePath from go.mod content
// WITHOUT golang.org/x/mod (the dependency budget is closed, #35): both the
// single-line and the block require forms resolve, `// indirect` suffixes and
// trailing comments are stripped, and other directives (replace/exclude) are
// ignored because their lines never carry a require-shaped version field.
func requireVersion(content, modulePath string) (string, bool) {
	scanner := requireScanner{modulePath: modulePath, inRequireBlock: false}

	for raw := range strings.SplitSeq(content, "\n") {
		line := stripLineComment(strings.TrimSpace(raw))
		if line == "" {
			continue
		}

		if version, found := scanner.line(strings.Fields(line)); found {
			return version, true
		}
	}

	return "", false
}

// requireScanner tracks whether the go.mod lines it sees lie inside a require
// block, and recognizes the lines that pin its module.
type requireScanner struct {
	modulePath     string
	inRequireBlock bool
}

func (scanner *requireScanner) line(fields []string) (string, bool) {
	blockStart, blockEnd := requireBlockMarkers(fields)

	switch {
	case blockStart:
		scanner.inRequireBlock = true
	case blockEnd:
		scanner.inRequireBlock = false
	case scanner.inRequireBlock && scanner.isBlockPin(fields):
		return fields[1], true
	case scanner.isSingleLinePin(fields):
		return fields[2], true
	}

	return "", false
}

// isBlockPin matches `github.com/x/y v1.2.3 [// indirect]` inside a require block.
func (scanner *requireScanner) isBlockPin(fields []string) bool {
	return len(fields) >= 2 && fields[0] == scanner.modulePath && isVersionLiteral(fields[1])
}

// isSingleLinePin matches `require github.com/x/y v1.2.3`.
func (scanner *requireScanner) isSingleLinePin(fields []string) bool {
	return len(fields) >= 3 && fields[0] == "require" && fields[1] == scanner.modulePath &&
		isVersionLiteral(fields[2])
}

// libraryContent returns the canonical bytes of a library file as this tc
// embeds them.
func libraryContent(name string) ([]byte, error) {
	content, err := sourcesFS.ReadFile(filepath.Join("_sources", "starter", name))
	if err != nil {
		return nil, fmt.Errorf("read embedded %s: %w", name, err)
	}

	return content, nil
}

// syncedContent is what `tc css sync` writes: the stamp line above the
// canonical bytes.
func syncedContent(file libraryFile) ([]byte, error) {
	content, err := libraryContent(file.name)
	if err != nil {
		return nil, err
	}

	stamp := "/* " + stampText() + " */"
	if file.plainText {
		stamp = "# " + stampText()
	}

	return append([]byte(stamp+"\n"), content...), nil
}

// stampText names the library version and the command that regenerates a copy.
func stampText() string {
	return stampMarker + libraryVersion() +
		" (library-owned, do not edit; regenerate with: " + syncCommand + ")"
}

// splitStamp separates the stamp line of a library-owned copy from its body
// and returns the version the stamp names ("" when the file is unstamped).
func splitStamp(content []byte) (string, []byte) {
	first, rest, _ := bytes.Cut(content, []byte("\n"))

	stamp, found := bytes.CutPrefix(first, []byte("/* "+stampMarker))
	if !found {
		stamp, found = bytes.CutPrefix(first, []byte("# "+stampMarker))
	}

	if !found {
		return "", content
	}

	version, _, _ := bytes.Cut(stamp, []byte(" "))

	return string(version), rest
}
