# crossbrowser-spike — 3-engine diagnostic screenshots (M06)

Playwright-driven full-page screenshots of every demo route in Chromium,
Firefox, and WebKit. **Diagnostic lane only** — produces artifacts and an
error report; it never gates CI and is deliberately NOT part of the Go test
surface (plan guardrail: never two pixel suites).

## Run

**Blessed path (M13):** `nix run .#crossbrowser` — builds the pinned browsers,
applies the pure fontconfig pin (same as the visual app), starts the demo binary,
and shoots all 3 engines. Findings land in `artifacts/` inside the spike dir.

Manual equivalent (what the app automates):

```bash
# 1. Browser binaries from the LOCKED nixpkgs (version-matched, see below)
nix build --inputs-from . nixpkgs#playwright-driver.browsers -o /tmp/pw-browsers

# 2. Node deps (playwright-core version below MUST equal nixpkgs playwright-driver)
cd tools/crossbrowser-spike && npm install --no-save

# 3. Route manifest from the demo registry (single source of truth)
./gen-routes.sh

# 4. Demo server (any free port)
go build -o /tmp/tc-demo ./examples/demo && (PORT=8093 /tmp/tc-demo &) 

# 5. Shoot all routes, all engines
for e in chromium firefox webkit; do
  PLAYWRIGHT_BROWSERS_PATH=/tmp/pw-browsers \
  PLAYWRIGHT_SKIP_VALIDATE_HOST_REQUIREMENTS=true \
  node screenshot.mjs --engine $e --base http://localhost:8093 --theme light
done

# 6. M07: base-select fallback assert over /stylable-select (Chromium + Firefox)
PLAYWRIGHT_BROWSERS_PATH=/tmp/pw-browsers \
PLAYWRIGHT_SKIP_VALIDATE_HOST_REQUIREMENTS=true \
node stylable-assert.mjs --base http://localhost:8093
```

Artifacts land in `artifacts/` (gitignored): `<engine>-<theme>/<route>.png` +
`report.json` with per-route console errors and pageerrors; the M07 fallback
report + per-select crops land in `artifacts/m07-fallback/`. The pre-fix M06
run is preserved at `artifacts-m06-initial/` (evidence for eval doc §5.1).

## Version lockstep (same-edit rule, like the templ pin)

`package.json` pins `playwright-core` to the nixpkgs `playwright-driver`
version (1.63.0 as of 2026-10-15). Playwright hardcodes browser revision
numbers per release (`chromium-1243`, `firefox-1543`, `webkit-2359`); a
mismatch fails with "Executable doesn't exist". Bump BOTH in one commit when
nixpkgs moves — `nix eval --inputs-from . nixpkgs#playwright-driver.version`.

## Findings

2026-10-15 spike results and the adoption verdict live in
`docs/testing/tooling-evaluation.md` (§5.1) and ADR 0048.

## Determinism (μ51)

Consecutive runs under the flake app's env are **pixel-identical**
(0/13,949,440 pixels differ on repeat captures) — the localStorage theme pin +
pure fontconfig pin make engine captures stable. Raw PNG BYTES may differ
between runs (encoder chunk noise): compare pixels with tolerance, never bytes.
The optional CI lane (`.github/workflows/crossbrowser.yml`) is
workflow_dispatch-only and stays INERT until ⫱G1.
