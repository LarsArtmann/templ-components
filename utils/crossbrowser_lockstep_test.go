package utils

import (
	"os"
	"regexp"
	"testing"
)

// expectedPlaywrightVersion is the playwright version every crossbrowser
// consumer must agree on. Playwright hardcodes browser revision numbers per
// release, so the spike's npm pin (tools/crossbrowser-spike/package.json +
// package-lock.json) and the LOCKED nixpkgs playwright-driver the flake
// builds browsers from must move in LOCKSTEP. When nixpkgs bumps
// playwright-driver, verify the new version
// (nix eval --inputs-from . nixpkgs#playwright-driver.version), then update
// THIS constant, the package.json pin, and the lockfile in one commit — the
// same one-commit lockstep rule as the Tailwind lanes.
const expectedPlaywrightVersion = "1.63.0"

// TestCrossbrowserLockstep pins the crossbrowser spike's playwright-core
// version against the nix-side playwright-driver version. Drift here fails
// at runtime with "Executable doesn't exist" (browser revision mismatch) —
// cheap to catch at test time instead. The npm-install hazard is real: an
// npm install inside the spike once rewrote the pin silently.
func TestCrossbrowserLockstep(t *testing.T) {
	t.Parallel()

	pkg, err := os.ReadFile("../tools/crossbrowser-spike/package.json")
	if err != nil {
		t.Fatalf("read spike package.json: %v", err)
	}
	pkgVersion := regexp.MustCompile(`"playwright-core"\s*:\s*"([^"]+)"`).FindSubmatch(pkg)
	if pkgVersion == nil {
		t.Fatalf("no playwright-core dependency in spike package.json:\n%s", pkg)
	}
	if string(pkgVersion[1]) != expectedPlaywrightVersion {
		t.Errorf("spike package.json playwright-core = %q, want %q (nixpkgs playwright-driver lockstep; update both + the flake in one commit)",
			pkgVersion[1], expectedPlaywrightVersion)
	}

	lock, err := os.ReadFile("../tools/crossbrowser-spike/package-lock.json")
	if err != nil {
		t.Fatalf("read spike package-lock.json: %v", err)
	}
	lockVersion := regexp.MustCompile(`"node_modules/playwright-core":\s*\{[^}]*"version":\s*"([^"]+)"`).FindSubmatch(lock)
	if lockVersion == nil {
		t.Fatalf("no node_modules/playwright-core entry in spike package-lock.json")
	}
	if string(lockVersion[1]) != expectedPlaywrightVersion {
		t.Errorf("spike package-lock.json playwright-core = %q, want %q (an npm install rewrote the pin — restore it)",
			lockVersion[1], expectedPlaywrightVersion)
	}
}
