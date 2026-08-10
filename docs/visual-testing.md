# Visual Regression Testing

This project has **two layers** of screenshot-style testing, each catching
different regressions:

| Layer             | Where                        | Compares   | Catches                                                                        |
| ----------------- | ---------------------------- | ---------- | ------------------------------------------------------------------------------ |
| HTML golden       | `utils/golden`               | HTML text  | Structure / class changes (classes are sorted)                                 |
| **Visual golden** | `visualtest/` (this package) | **Pixels** | Layout shifts, dark-mode color regressions, RTL mirroring, responsive collapse |

The HTML golden tests normalize CSS class order and diff strings, so they are
blind to anything visual. The visual tests render each component in a real
headless Chromium, capture an element screenshot, and diff the pixels against a
committed PNG. A 0.1% pixel threshold absorbs anti-aliasing noise across
Chromium builds while still flagging any real visual change.

## How it works

render component into an isolated HTML page (compiled Tailwind CSS inlined)
│
▼
serve page on an ephemeral httptest.Server
│
▼
chromedp: navigate → wait for #tc-root → (apply State) → settle → screenshot
│ (element OR full viewport)
▼
compare PNG pixels against testdata/<name>.png (0.1% threshold)
│
├── match → pass
└── differ → fail + write testdata/.fail/<name>.{actual,diff}.png

````

### Shared Chromium process

A **single** headless Chromium process is launched lazily on the first
`AssertScreenshot` call and reused for every test. Each test gets a fresh tab
(`chromedp.NewContext`, ~10ms) instead of a full browser launch (~1s), so the
full ~30-test suite finishes in ~4s. `TestMain` tears the process down after
all tests complete. If `CHROMEDP_CHROME_PATH` is unset or invalid, the
allocator skips (and every test skips with it) — `go test ./...` stays green
wherever Chromium is absent.

## Running

The visual tests need a Chromium binary. The Nix app wires it up automatically:

```bash
# Run all visual tests (headless Chromium provided by Nix)
nix run .#visual

# Regenerate every golden PNG after an intentional visual change
nix run .#visual -- -update

# Run a single test
nix run .#visual -- -run TestButtons
````

Without Nix, point `CHROMEDP_CHROME_PATH` at any Chromium/Chrome binary and run
the module directly (the `GOWORK=off` is required because the repo's parent
`go.work` would otherwise shadow this module's local replace directive):

```bash
export GOEXPERIMENT=jsonv2
export GOWORK=off
export CHROMEDP_CHROME_PATH="$(which chromium)"
cd visualtest
go test ./...              # compare
go test ./... -update      # regenerate goldens
```

If no browser is found, every test **skips** (not fails), so `go test ./...`
from the repo root stays green in environments without Chromium.

## Writing a new visual test

```go
func TestMyComponent(t *testing.T) {
    t.Parallel()

    props := display.DefaultCardProps()
    props.Title = "Revenue"


```

- `StateClick` clicks the first `[popovertarget]`/button/link inside `#tc-root`.
- `StateContext` dispatches a `contextmenu` event (for `ContextMenu`).
- Raise `MaxMismatch` to ~1% for JS-positioned overlays: the menu is placed
  from the trigger's `getBoundingClientRect()`, so the threshold must absorb
  Chromium-version micro-drift. A 10x serialized calibration confirmed 0.0000%
  run-to-run mismatch (fully deterministic in pinned Chromium), so 1% is pure
  headroom for version drift. A real regression blows past 1%.
- Pass `Nonce` on the component props so positioning scripts (and
  `ContextMenu`'s menu + handler, which are gated on `Nonce != ""`) render.

### Golden naming

Names map to files under `testdata/` (e.g. `"button/primary_dark"` →
`testdata/button/primary_dark.png`). Keep one subdirectory level per component
family so the directory tree mirrors the component packages.

## MaxMismatch calibration

The default `MaxMismatch` is 0.1% (catches any real visual change while
absorbing cross-build anti-aliasing noise). JS-positioned overlays (Dropdown,
Popover, ContextMenu) and native `<dialog>` overlays (Modal, Drawer) use a
raised 1% threshold because their menus are placed from the trigger's
`getBoundingClientRect()` and the threshold must absorb Chromium-version
micro-drift (a `nixpkgs-chromium` bump shifts rendered pixels by a fraction of
a percent).

### Calibration methodology (2026-08-04)

Each overlay test was run 10x serialized (`-count=10 -parallel 1`) under the
pinned Chromium to measure run-to-run variance:

| Golden                 | Helper      | 10x result |
| ---------------------- | ----------- | ---------- |
| dropdown/open_light    | overlayOpen | 0.0000%    |
| dropdown/open_dark     | overlayOpen | 0.0000% *  |
| popover/open_light     | overlayOpen | 0.0000%    |
| contextmenu/open_light | overlayOpen | 0.0000%    |
| modal/open_light       | dialogOpen  | 0.0000%    |
| modal/open_dark        | dialogOpen  | 0.0000%    |
| drawer/right_light     | dialogOpen  | 0.0000%    |
| drawer/left_dark       | dialogOpen  | 0.0000%    |

\* `dropdown/open_dark` measured a stable 0.7442% systematic diff before this
calibration — a stale golden, not anti-aliasing variance (the prior comment
misattributed it to AA noise). The golden was regenerated against the pinned
Chromium, after which it reads 0.0000% deterministically.

**Conclusion:** rendering is fully deterministic in the pinned headless Chromium
(zero run-to-run variance), so the 1% threshold is pure headroom for
Chromium-version drift, not anti-aliasing noise. A real regression (missing
menu, wrong colors, broken layout) blows far past 1%.

### Animation settle (race fix)

The initial serialized calibration showed 0% mismatch, but the **full parallel
suite** flaked ~20% of the time on Modal/Drawer captures (~90% false mismatch).
Root cause: `WaitVisible("dialog")` returns the instant `showModal()` makes the
`<dialog>` `display:block`, but the `@starting-style` slide-in transition
(200ms, defined in `templates/custom.css`) can still be mid-flight under
parallel load — capturing the drawer off-screen. The harness now calls
`waitAnimationSettled` after `WaitVisible`, which polls `getAnimations()` until
all CSS transitions finish before capturing. After the fix the full parallel
suite passes 8/8 with 0.0000% overlay mismatch.

## When a visual test fails

1. Open `testdata/.fail/<name>.diff.png` — red pixels show what changed.
2. Open `testdata/.fail/<name>.actual.png` next to the committed golden.
3. If the change is **intended**, regenerate: `nix run .#visual -- -update`,
   review the diff, and commit the new golden.
4. If the change is **unintended**, fix the component.

The `.fail/` directory is gitignored — it holds only the last failing run's
artifacts for inspection, never the source of truth.

## Why a separate Go module?

`chromedp` and its transitive dependencies are heavy and test-only. Adding
them to the library's `go.mod` would pollute every consumer's dependency graph
(they'd land in `go.sum` even though they're never compiled). The `visualtest/`
directory is its **own module** with a local `replace` directive pointing at the
parent, so `go get github.com/larsartmann/templ-components` never pulls
chromedp. This mirrors the existing single-module boundary while keeping the
test tooling cleanly separated.

## Keeping the compiled CSS fresh

The harness reads `examples/demo/static/app.css` (the demo's compiled Tailwind
output) at test time. If you add Tailwind classes to a component, recompile the
demo CSS so the visual tests render the new classes:

```bash
cd examples/demo && tailwindcss -i demo.css -o static/app.css
```

(or `nix run .#build`, which regenerates everything).
