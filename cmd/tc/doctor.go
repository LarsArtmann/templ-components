package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	for _, line := range strings.Split(string(content), "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) == 2 && (fields[0] == "go" || fields[0] == "toolchain") {
			version := strings.TrimPrefix(fields[1], "go")
			major, minor, ok := parseGoVersion(version)
			if ok && (major > 1 || (major == 1 && minor >= 27)) {
				return true
			}
		}
	}

	return false
}

// parseGoVersion extracts major.minor from "1.26.7"-style strings.
func parseGoVersion(version string) (int, int, bool) {
	parts := strings.SplitN(version, ".", 3)
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

	for _, line := range strings.Split(string(content), "\n") {
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

// templGeneratedWith is the generator version the library pins.
const templGeneratedWith = "v0.3.1020"

// compareSemver compares dotted numeric versions; -1, 0, or 1.
func compareSemver(a, b string) int {
	aParts, bParts := strings.Split(a, "."), strings.Split(b, ".")

	for i := range max(len(aParts), len(bParts)) {
		av, bv := 0, 0

		if i < len(aParts) {
			_, _ = fmt.Sscanf(aParts[i], "%d", &av)
		}

		if i < len(bParts) {
			_, _ = fmt.Sscanf(bParts[i], "%d", &bv)
		}

		switch {
		case av < bv:
			return -1
		case av > bv:
			return 1
		}
	}

	return 0
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

	out, err := exec.CommandContext(ctx, "git", "config", "--get", key).Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
