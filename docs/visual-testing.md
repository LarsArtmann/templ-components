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
  chromedp: navigate → wait for #tc-root → (apply State) → settle → screenshot
        │                                                                  (element OR full viewport)
        ▼
  compare PNG pixels against testdata/<name>.png (0.1% threshold)
        │
        ├── match  → pass
        └── differ → fail + write testdata/.fail/<name>.{actual,diff}.png
```

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

| Option          | Default  | Purpose                                                                                          |
| --------------- | -------- | ------------------------------------------------------------------------------------------------ |
| `Dark`          | false    | Adds `class="dark"` to `<html>` (the library's dark strategy)                                   |
| `RTL`           | false    | Sets `dir="rtl"` to test logical-property mirroring                                              |
| `Viewport`      | 1280×800 | Emulated window size for responsive variants                                                     |
| `MaxMismatch`   | 0.001    | Max fraction of mismatched pixels (0–1) that still passes                                         |
| `Threshold`     | 0.1      | Pixelmatch perceptual color-distance threshold (0–1); higher tolerates more rendering noise       |
| `State`         | `Rest`   | Interaction to apply before capture: `StateHover`, `StateFocus`, `StateClick`, `StateContext`    |
| `FullViewport`  | false    | Capture the full viewport instead of the `#tc-root` element (required for top-layer overlays)    |
| `WaitSelector`  | `""`     | CSS selector to wait for (visible) after applying `State` — e.g. `[popover]` once a menu opens    |

### Testing overlays (Dropdown / Popover / ContextMenu / Modal / Drawer)

Components whose open state renders in the browser **top layer** (native
Popover API menus and `<dialog>`) paint *outside* `#tc-root`'s bounding box,
so a normal element screenshot crops them. Use `StateClick` (or
`StateContext` for right-click menus) with `FullViewport: true` and
`WaitSelector: "[popover]"`:

```go
visualtest.AssertScreenshot(t, "dropdown/open_light", display.Dropdown(props),
    visualtest.Options{
        State:        visualtest.StateClick,      // click the popovertarget trigger
        WaitSelector: "[popover]",                // wait for the menu to appear
        FullViewport: true,                       // capture the top-layer menu
        Viewport:     visualtest.Viewport{Width: 480, Height: 360},
        MaxMismatch:  0.02,                       // JS-positioned overlays have ~1px jitter
    })
```

- `StateClick` clicks the first `[popovertarget]`/button/link inside `#tc-root`.
- `StateContext` dispatches a `contextmenu` event (for `ContextMenu`).
- Raise `MaxMismatch` to ~2% for JS-positioned overlays: the menu is placed
  from the trigger's `getBoundingClientRect()`, so a 1px layout-timing shift
  shows up as edge anti-aliasing variance. A real regression blows past 2%.
- Pass `Nonce` on the component props so positioning scripts (and
  `ContextMenu`'s menu + handler, which are gated on `Nonce != ""`) render.

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
