# Visual Regression Testing

This project has **two layers** of screenshot-style testing, each catching
different regressions:

| Layer             | Where                        | Compares   | Catches                                                                        |
| ----------------- | ---------------------------- | ---------- | ------------------------------------------------------------------------------ |
| HTML golden       | `internal/golden`            | HTML text  | Structure / class changes (classes are sorted)                                 |
| **Visual golden** | `visualtest/` (this package) | **Pixels** | Layout shifts, dark-mode color regressions, RTL mirroring, responsive collapse |

The HTML golden tests normalize CSS class order and diff strings, so they are
blind to anything visual. The visual tests render each component in a real
headless Chromium, capture an element screenshot, and diff the pixels against a
committed PNG. A 0.1% pixel threshold absorbs anti-aliasing noise across
Chromium builds while still flagging any real visual change.

## How it works

```
AssertScreenshot(t, "button/primary_dark", display.Button(props), visualtest.Options{Dark: true})
        │
        ▼
  render component into an isolated HTML page (compiled Tailwind CSS inlined)
        │
        ▼
  serve page on an ephemeral httptest.Server
        │
        ▼
  chromedp: navigate → wait for #tc-root → settle → element screenshot
        │
        ▼
  compare PNG pixels against testdata/<name>.png (0.1% threshold)
        │
        ├── match  → pass
        └── differ → fail + write testdata/.fail/<name>.{actual,diff}.png
```

## Running

The visual tests need a Chromium binary. The Nix app wires it up automatically:

```bash
# Run all visual tests (headless Chromium provided by Nix)
nix run .#visual

# Regenerate every golden PNG after an intentional visual change
nix run .#visual -- -update

# Run a single test
nix run .#visual -- -run TestButtons
```

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

    visualtest.AssertScreenshot(t, "card/basic_light", display.Card(props))
    visualtest.AssertScreenshot(t, "card/basic_dark",  display.Card(props), visualtest.Options{Dark: true})
    visualtest.AssertScreenshot(t, "card/basic_rtl",   display.Card(props), visualtest.Options{RTL: true})
}
```

### Options

| Option              | Default  | Purpose                                                       |
| ------------------- | -------- | ------------------------------------------------------------- |
| `Dark`              | false    | Adds `class="dark"` to `<html>` (the library's dark strategy) |
| `RTL`               | false    | Sets `dir="rtl"` to test logical-property mirroring           |
| `Viewport`          | 1280×800 | Emulated window size for responsive variants                  |
| `MaxMismatch`       | 0.001    | Max fraction of mismatched pixels (0–1) that still passes     |
| `PerPixelTolerance` | 32       | Max per-channel diff (0–255) that counts as a matching pixel  |

### Golden naming

Names map to files under `testdata/` (e.g. `"button/primary_dark"` →
`testdata/button/primary_dark.png`). Keep one subdirectory level per component
family so the directory tree mirrors the component packages.

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
