package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// cmdDoctor diagnoses the top consumer-integration traps in one command:
// Tailwind source scanning, the jsonv2 build flag, templ version pinning,
// committed generated files, and (inside the library repo) hook adoption.
// Every check prints pass/fail + a fix hint; the exit code is non-zero when
// any check fails so CI can gate on it.
//
// It runs in TWO contexts and adapts:
//   - inside the templ-components repo (checks library-side invariants),
//   - in a consumer project (checks the consumer-side integration traps).
func cmdDoctor(_ *registry, _ []string) {
	fmt.Fprintln(os.Stdout, "tc doctor — templ-components environment check")
	fmt.Fprintln(os.Stdout)

	failures := 0
	failures += checkTailwindSources()
	failures += checkGOEXPERIMENT()
	failures += checkTemplVersionPinned()
	failures += checkGeneratedCommitted()
	failures += checkHooksAdoption()
	failures += checkNoncelessScripts()

	fmt.Fprintln(os.Stdout)

	if failures > 0 {
		fmt.Fprintf(os.Stdout, "result: %d check(s) FAILED — see the fix hints above.\n", failures)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, "result: all checks passed.")
}

// doctorCheck prints one pass/fail line and returns 1 on failure.
func doctorCheck(name string, ok bool, hint string) int {
	status := "pass"
	if !ok {
		status = "FAIL"
	}

	fmt.Fprintf(os.Stdout, "  [%s] %s\n", status, name)

	if !ok && hint != "" {
		fmt.Fprintf(os.Stdout, "         fix: %s\n", hint)
	}

	if !ok {
		return 1
	}

	return 0
}

// checkTailwindSources verifies the consumer's CSS entry point scans the
// directory that holds copied .templ components (the #1 integration trap:
// @source path resolved relative to the CSS file, missing every class).
func checkTailwindSources() int {
	cssCandidates := []string{"app.css", "static/app.css", "assets/app.css", "css/app.css"}

	for _, candidate := range cssCandidates {
		content, err := os.ReadFile(candidate)
		if err != nil {
			continue
		}

		text := string(content)
		hasSource := strings.Contains(text, "@source")
		scansTempl := strings.Contains(text, "**/*.templ")

		return doctorCheck(
			"Tailwind @source scans .templ files ("+candidate+")",
			hasSource && scansTempl,
			"add `@source \"./**/*.templ\";` (path relative to the CSS FILE, not the CWD — the demo uses @source \"../../**/*.templ\")",
		)
	}

	return doctorCheck(
		"Tailwind entry (app.css) found",
		false,
		"no app.css in ./, ./static, ./assets, ./css — run `tc init` and compile with tailwindcss",
	)
}

// checkGOEXPERIMENT verifies the jsonv2 build flag is set (env, go.mod
// toolchain directive, or Go >= 1.27 where it is stable). Without it the
// build fails with `build constraints exclude all Go files in .../encoding/json/v2`.
func checkGOEXPERIMENT() int {
	if os.Getenv("GOEXPERIMENT") == "jsonv2" {
		return doctorCheck("GOEXPERIMENT=jsonv2 (env)", true, "")
	}

	if goModDeclaresJSONv2() {
		return doctorCheck("GOEXPERIMENT=jsonv2 (go.mod toolchain directive)", true, "")
	}

	return doctorCheck(
		"GOEXPERIMENT=jsonv2 required (errorpage uses encoding/json/v2)",
		false,
		"`export GOEXPERIMENT=jsonv2` (or use Go 1.27+, or set `toolchain`/`go` accordingly); symptom without it: 'build constraints exclude all Go files in .../encoding/json/v2'",
	)
}

// goModDeclaresJSONv2 reports whether a go.mod in the CWD declares a Go
// version where jsonv2 is stable (1.27+).
func goModDeclaresJSONv2() bool {
	content, err := os.ReadFile("go.mod")
	if err != nil {
		return false
	}

	for line := range strings.SplitSeq(string(content), "\n") {
		fields := strings.Fields(strings.TrimSpace(line))

		if len(fields) == 2 && (fields[0] == "go" || fields[0] == "toolchain") {
			version := strings.TrimPrefix(fields[1], "go")
			major, minor, ok := parseGoVersion(version)

			if ok && (major > 1 || (major == 1 && minor >= jsonv2StableMinor)) {
				return true
			}
		}
	}

	return false
}

// parseGoVersion extracts major.minor from "1.26.7"-style strings.
func parseGoVersion(version string) (int, int, bool) {
	parts := strings.SplitN(version, ".", goVersionParts)
	if len(parts) < 2 {
		return 0, 0, false
	}

	var major, minor int

	if _, err := fmt.Sscanf(parts[0]+"."+parts[1], "%d.%d", &major, &minor); err != nil {
		return 0, 0, false
	}

	return major, minor, true
}

// checkTemplVersionPinned warns when a consumer's go.mod pulls a templ
// version older than the one this library generates with (v0.3.1020).
func checkTemplVersionPinned() int {
	content, err := os.ReadFile("go.mod")
	if err != nil {
		return doctorCheck("go.mod readable (templ pin check)", false, "run inside a Go module")
	}

	for line := range strings.SplitSeq(string(content), "\n") {
		fields := strings.Fields(strings.TrimSpace(line))

		if len(fields) < 3 || fields[0] != "require" {
			continue
		}

		if fields[1] != "github.com/a-h/templ" {
			continue
		}

		pinned := strings.TrimPrefix(fields[2], "v")
		want := strings.TrimPrefix(templGeneratedWith, "v")

		return doctorCheck(
			"templ "+pinned+" compatible with generated code ("+templGeneratedWith+")",
			compareSemver(pinned, want) >= 0,
			"`go get github.com/a-h/templ@"+templGeneratedWith+"` — older generators emit incompatible *_templ.go runtime calls",
		)
	}

	return doctorCheck(
		"templ dependency present",
		false,
		"add github.com/a-h/templ — the generated components need its runtime",
	)
}

const (
	// goVersionParts bounds version parsing to major.minor.patch.
	goVersionParts = 3
	// jsonv2StableMinor is the Go minor where encoding/json/v2 needs no flag.
	jsonv2StableMinor = 27
)

// templGeneratedWith is the generator version the library pins.
const templGeneratedWith = "v0.3.1020"

// compareSemver compares dotted numeric versions; -1, 0, or 1.
func compareSemver(a, b string) int {
	aParts, bParts := strings.Split(a, "."), strings.Split(b, ".")

	for i := range max(len(aParts), len(bParts)) {
		aPart, bPart := 0, 0

		if i < len(aParts) {
			_, _ = fmt.Sscanf(aParts[i], "%d", &aPart)
		}

		if i < len(bParts) {
			_, _ = fmt.Sscanf(bParts[i], "%d", &bPart)
		}

		switch {
		case aPart < bPart:
			return -1
		case aPart > bPart:
			return 1
		}
	}

	return 0
}

// nonceConsumers lists the library components that render an inline
// <script> gated on the consumer-supplied Nonce. A call site that passes an
// empty nonce renders markup that silently dies under a nonce-locked CSP —
// the 2026-10-12 fleet-audit bug class (cqrs-htmx CopyButton ×13, GmbH
// ThemeToggle, KeyHolderAI ToastContainer). docs/nonce-consumers.md is the
// human-readable twin of this table; TestNonceConsumerDocsSync pins them
// together.
var nonceConsumers = []struct{ Qualified, Note string }{ //nolint:gochecknoglobals // Package-level scan table, the docs twin pins it
	{"display.Accordion", "chevron rotation script"},
	{"display.Carousel", "prev/next/dot navigation script"},
	{"display.ContextMenu", "menu keyboard-nav + positioning script"},
	{"display.CopyButton", "clipboard singleton script"},
	{"display.Dropdown", "menu keyboard-nav + positioning script"},
	{"display.Drawer", "auto-open + backdrop-click script"},
	{"display.KanbanBoard", "drag-and-drop + optimistic move script"},
	{"display.Lightbox", "viewer navigation script"},
	{"display.Modal", "auto-open + backdrop-click script"},
	{"display.Table", "row-click script when any TableRow.Href is set"},
	{"display.Tabs", "client-side tab switching (when not Wire-driven)"},
	{"errorpage.ErrorPage", "go-back singleton script"},
	{"errorpage.NotFound404", "go-back singleton script"},
	{"echarts.EChart", "chart init script"},
	{"feedback.Alert", "dismiss script when Dismissible"},
	{"feedback.Toast", "dismiss + auto-dismiss timer script"},
	{"feedback.ToastContainer", "container for injected toasts"},
	{"forms.Combobox", "filter/selection singleton script"},
	{"forms.DirtyGuard", "unsaved-changes guard script"},
	{"forms.TagsInput", "tag add/remove script"},
	{"htmx.GlobalErrorHandling", "global error→toast pipeline script"},
	{"htmx.PolledRegion", "poll loop script when Eager"},
	{"layout.ThemeScript", "FOUC-prevention head script"},
	{"layout.ThemeToggle", "theme toggle script"},
}

// Note: MobileMenu's toggle script also consumes a nonce; consumers reach it
// through navigation.Nav/SimpleNav (pass the nonce to the nav component).
// It is documented in docs/nonce-consumers.md rather than the table above
// because the consuming constructor is positional-arg based.

// checkNoncelessScripts scans the consumer's .templ sources for call sites
// that pass an empty nonce to a script-emitting component — markup that
// silently dies under a nonce-locked CSP. Findings FAIL when the repo
// declares a Content-Security-Policy anywhere in its .go sources, and stay
// informational otherwise (nothing is broken yet, but the trap is armed).
func checkNoncelessScripts() int {
	workDir, _ := os.Getwd()

	if _, err := os.Stat(filepath.Join(workDir, "go.mod")); err != nil {
		return 0 // not a Go module — check does not apply
	}

	if len(walkTemplSources(workDir)) == 0 {
		return 0 // no templ sources — check does not apply
	}

	findings := scanNonceFindings(workDir)

	if len(findings) == 0 {
		return doctorCheck("no nonce-less script call sites in .templ sources", true, "")
	}

	cspDeclared := declaresCSP(workDir)

	hint := strings.Join(findings, "; ") +
		" — see docs/nonce-consumers.md for the full consuming-component table"
	if cspDeclared {
		return doctorCheck(
			fmt.Sprintf("no nonce-less script call sites in .templ sources (%d found, CSP declared)", len(findings)),
			false,
			hint,
		)
	}

	fmt.Fprintf(os.Stdout, "  info: %d nonce-less script call site(s), but no CSP declared yet — "+
		"the scripts run today and will silently die the day you adopt one\n", len(findings))

	return 0
}

// scanNonceFindings returns file:line findings for nonce-less call sites to
// script-emitting components under root. Kept separate from the check so the
// tests can assert on the raw findings.
func scanNonceFindings(root string) []string {
	propsRe := regexp.MustCompile(`Nonce:\s*""`)
	positionalRe := regexp.MustCompile(
		`(ThemeScript|ThemeToggle|ToastContainer|MobileMenu)\(([^()]|"[^"]*")*?,?\s*(""|"[a-z][a-z0-9-]*")`)

	findings := []string{}

	for _, file := range walkTemplSources(root) {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			code := stripTemplComment(line)
			if propsRe.MatchString(code) {
				findings = append(findings, fmt.Sprintf("%s:%d: Nonce: \"\" — thread the request nonce", file, i+1))

				continue
			}

			if positionalRe.MatchString(code) {
				findings = append(findings, fmt.Sprintf(
					"%s:%d: empty/fake literal nonce passed to a nonce-consuming constructor", file, i+1))
			}
		}
	}

	return findings
}

// walkTemplSources returns .templ files under root, skipping generated
// files, vendor/node_modules, and the library's own packages.
func walkTemplSources(root string) []string {
	files := []string{}

	skip := map[string]bool{
		"vendor": true, "node_modules": true, ".git": true,
		"display": true, "feedback": true, "forms": true, "layout": true,
		"navigation": true, "errorpage": true, "htmx": true, "datastar": true,
		"utils": true, "icons": true, "cmd": true, "examples": true,
		"internal": true, "charts": true, "recipes": true, "website": true,
		"integration": true,
	}

	_ = filepath.WalkDir(root, func(path string, dirEntry os.DirEntry, err error) error {
		if err != nil {
			return nil //nolint:nilerr // skip unreadable entries; a scan must not abort on one bad path
		}

		if dirEntry.IsDir() {
			if skip[dirEntry.Name()] {
				return filepath.SkipDir
			}

			return nil
		}

		if strings.HasSuffix(path, ".templ") && !strings.HasSuffix(path, "_templ.go") {
			files = append(files, path)
		}

		return nil
	})

	return files
}

// stripTemplComment removes an HTML comment tail so prose never trips the
// nonce regexes (same defense as the utils compliance scans).
func stripTemplComment(line string) string {
	if before, _, ok := strings.Cut(line, "<!--"); ok {
		return before
	}

	return line
}

// declaresCSP reports whether any .go file under root (excluding tests,
// vendor, and generated code) mentions a Content-Security-Policy.
func declaresCSP(root string) bool {
	cspRe := regexp.MustCompile(`Content-Security-Policy`)
	found := false

	_ = filepath.WalkDir(root, func(path string, dirEntry os.DirEntry, err error) error {
		if err != nil || found {
			return filepath.SkipAll
		}

		if dirEntry.IsDir() {
			switch dirEntry.Name() {
			case "vendor", "node_modules", ".git", "testdata":
				return filepath.SkipDir
			}

			return nil
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") ||
			strings.HasSuffix(path, "_templ.go") {
			return nil
		}

		data, readErr := os.ReadFile(path) //nolint:gosec // diagnostic scan over the consumer's own tree
		if readErr == nil && cspRe.Match(data) {
			found = true
		}

		return nil
	})

	return found
}

// checkGeneratedCommitted applies only inside the library repo: a .templ
// source without its committed *_templ.go means the proxy will serve
// uncompilable code.
func checkGeneratedCommitted() int {
	wd, _ := os.Getwd()

	// Heuristic: inside the library repo a component package dir exists.
	if _, err := os.Stat(filepath.Join(wd, "display")); err != nil {
		return 0 // consumer project — check does not apply
	}

	missing := []string{}

	for _, pkg := range []string{"display", "feedback", "forms", "layout", "navigation", "errorpage", "htmx", "datastar"} {
		entries, err := os.ReadDir(pkg)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			name := entry.Name()
			if !strings.HasSuffix(name, ".templ") || entry.IsDir() {
				continue
			}

			generated := strings.TrimSuffix(name, ".templ") + "_templ.go"
			if _, statErr := os.Stat(filepath.Join(pkg, generated)); statErr != nil {
				missing = append(missing, filepath.Join(pkg, generated))
			}
		}
	}

	return doctorCheck(
		"every .templ has its *_templ.go (library repo)",
		len(missing) == 0,
		"run `templ generate ./...` and COMMIT the *_templ.go files (proxy serves the tag as-is): "+strings.Join(
			missing,
			", ",
		),
	)
}

// checkHooksAdoption applies only inside the library repo: the tracked
// .githooks/ guards are active only when core.hooksPath points there.
func checkHooksAdoption() int {
	wd, _ := os.Getwd()

	if _, err := os.Stat(filepath.Join(wd, ".githooks")); err != nil {
		return 0 // consumer project — check does not apply
	}

	hooksPath := gitConfig("core.hooksPath")

	return doctorCheck(
		"core.hooksPath=.githooks (library repo)",
		hooksPath == ".githooks",
		"run `scripts/setup-hooks.sh` — otherwise the tracked pre-commit guards are inactive",
	)
}

// gitConfig reads a git config value via git itself (respects includes,
// scopes, and the `key = value` spacing variants).
func gitConfig(key string) string {
	ctx := context.Background()

	//nolint:gosec // key is a hardcoded literal from this file, never user input
	out, err := exec.CommandContext(ctx, "git", "config", "--get", key).Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
