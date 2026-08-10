# Testing Guide

How templ-components verifies that every component renders correctly, stays
accessible, and does not regress visually.

## Three-tier strategy

| Tier                    | Mechanism                                                      | Location                                     | Catches                                                                     | Cost                |
| ----------------------- | -------------------------------------------------------------- | -------------------------------------------- | --------------------------------------------------------------------------- | ------------------- |
| 1. HTML golden          | Render component → normalize → diff against `.golden` file     | `utils/golden`, `<pkg>/testdata/*.golden`    | Structure, attributes, class changes                                        | Fast, deterministic |
| 2. Drift-guard scanners | Repo-wide invariant tests in `utils/`                          | `utils/*_compliance_test.go`, `integration/` | Cross-cutting regressions (dark-mode gaps, missing motion-reduce, RTL, CSP) | Fast                |
| 3. Visual regression    | Render in headless Chromium → pixel diff against committed PNG | `visualtest/` (separate module)              | Layout shifts, color regressions, RTL mirroring                             | ~4s full suite      |

### Tier 1 — HTML golden tests

The backbone. Every component has at least one golden snapshot — the exact
rendered HTML after two normalizations:

1. **CSS class tokens are sorted alphabetically** inside each `class="..."` so
   that `utils.Class()` / tailwind-merge reordering does not cause spurious diffs.
2. **Auto-generated `EnsureID` values** (`tc-<prefix>-<16hex>`) are replaced with
   `tc-<prefix>-NORMALIZED`, so components that auto-generate IDs (Accordion,
   Tabs, Dropdown, Tooltip, Carousel, ContextMenu, Combobox, Nav) are stable
   across runs without needing explicit IDs.

Golden files live in `<package>/testdata/*.golden` and use an LCS-based diff
with line numbers for readable failures.

```go
// Table-driven sweep — the recommended pattern (see golden_sweep_test.go files)
func TestGoldenSweepButton(t *testing.T) {
    golden.AssertSnapshots(t, []golden.Snapshot{
        {Name: "button_primary", HTML: utils.Render(t, Button(DefaultButtonProps()))},
        {Name: "button_disabled", HTML: utils.Render(t, Button(ButtonProps{Disabled: true}))},
    })
}
```

**Updating goldens** after an intentional visual change:

```bash
go test ./display/... -run TestGoldenSweep -update
```

Then review the `.golden` diff before committing.

### Tier 2 — Drift-guard scanners

These live in `utils/` (and `integration/`) and enforce cross-cutting invariants
that per-component tests cannot. They are the reason `go test ./...` stays
trustworthy across refactors.

| Test                                         | File                                       | Prevents                                                                             |
| -------------------------------------------- | ------------------------------------------ | ------------------------------------------------------------------------------------ |
| `TestDarkModeCompliance` / `…SemanticColors` | `utils/darkmode_compliance_test.go`        | Neutral/semantic colors without `dark:` variants                                     |
| `TestMotionReduceCompliance`                 | `utils/motion_compliance_test.go`          | `transition-*`/`animate-*` without `motion-reduce:`                                  |
| `TestRTLLogicalProperties`                   | `utils/rtl_compliance_test.go`             | Physical properties (`ml-`/`mr-`/`left-`) instead of logical                         |
| `TestNoOrderedTailwindSubstringsInTests`     | `utils/ordered_substring_test.go`          | Brittle ordered-class substring assertions that flake under `utils.Class` reordering |
| `TestGolangciDisabledLinters`                | `utils/lint_config_test.go`                | Incompatible linters re-entering `.golangci.yml`                                     |
| `TestTemplGeneratedInSync`                   | `utils/templ_sync_test.go`                 | `.templ` edit committed without regenerating `*_templ.go`                            |
| `TestContainerQueryCompliance`               | `utils/container_query_compliance_test.go` | Viewport breakpoints without `ContainerAware` opt-in                                 |
| `TestTailwindGoSourceScanning`               | `utils/tailwind_source_test.go`            | Tailwind classes in `.go` map literals missing from compiled CSS                     |
| `TestCSSFreshness`                           | `utils/css_freshness_test.go`              | Committed demo CSS older than newest source (fails in CI)                            |
| CSP nonce test                               | `integration/csp_nonce_test.go`            | Inline `<script>` without `nonce=`                                                   |

Run all of them: `go test ./utils/... ./integration/...`

### Tier 3 — Visual regression (`visualtest/`)

Pixel-level PNG comparison in headless Chromium (chromedp). Lives in a
**separate Go module** so chromedp never pollutes the consumer dependency graph.

```go
func TestButtonPrimary(t *testing.T) {
    visualtest.AssertScreenshot(t, "button/primary_dark",
        display.Button(display.DefaultButtonProps()),
        visualtest.Options{Dark: true})
}
```

- A 0.1% pixel mismatch threshold absorbs anti-aliasing noise across Chromium
  builds while still flagging any real visual change.
- A single shared Chromium process is reused across tests (~4s full suite).
- Supports `Dark`, `RTL`, `Viewport`, `State` (`Hover`/`Focus`/`Click`), and
  `FullViewport` (for top-layer overlays like modal/drawer).
- Skips automatically if no Chromium binary is found — `go test ./...` stays
  green wherever Chromium is absent.

**Running:**

```bash
nix run .#visual                   # run all visual tests (Nix provides Chromium)
nix run .#visual -- -update        # regenerate every golden PNG
nix run .#visual -- -run TestButton # run a single test
```

See [`docs/visual-testing.md`](visual-testing.md) for the full visual testing
reference.

## Running the full suite

```bash
# The canonical "done" check — generate + build + test + lint in one shot:
nix run .#verify

# Just tests:
nix run .#test      # go test ./... -count=1 -race

# Coverage:
nix run .#coverage  # go test -coverprofile + summary
```

## Adding tests for a new component

Every new component should have:

```
[ ] golden_test.go / golden_sweep_test.go — exact rendered HTML snapshot
[ ] a11y_test.go        — ARIA, roles, keyboard, motion-reduce, screen-reader text
[ ] edge_cases_test.go  — empty inputs, unknown enum values, ID collisions
[ ] example_test.go     — godoc ExampleXxx() compiles and renders
[ ] bdd_test.go         — behaviour spec (user-visible behaviour, not markup)
[ ] snapshot_test.go    — broader composition snapshot (optional)
[ ] coverage_*_test.go  — targeted coverage of private helpers and branches
```

After adding golden snapshots:

```bash
go test ./yourpackage/... -run TestGoldenSweep -update
```

Review the generated `.golden` files, then commit them.

## Fuzz and benchmark tests

- **Fuzz tests:** `forms.FuzzInputType`, `forms.FuzzFormMethod`,
  `display.FuzzButtonHTMLType` verify enum validation never panics on arbitrary
  input.
  ```bash
  go test -fuzz=. -run=Fuzz ./...
  ```
- **Benchmarks:** 7 packages (display, feedback, navigation, forms, layout,
  htmx, icons, utils).
  ```bash
  go test -bench=. -benchmem ./...
  ```
