# Count definitions — what each sold number means and where it is computed

Every public-facing count in this repo must be either live-computed by a guard
or pinned to a source that is. This page is the definitional home: when a
number is questioned, answer from here, not from memory.

## Components — **123**

- **Definition:** exported `templ` components declared as `templ Name(...)` in
  the canonical 9-package scan set: `display`, `feedback`, `forms`, `layout`,
  `navigation`, `charts/echarts`, `htmx`, `datastar`, `errorpage`. The 4
  `recipes` screens are compositions, counted separately in FEATURES.md
  ("123 + 4 recipe screens = 127") but the sold headline is 123.
- **Computed by:** `website/internal/build.CountStats` (site hero, live) over
  `statsDirs` (order-pinned by `TestStatsDirsAreCanonical`); the same set is
  the canonical package list in `utils.TestDocsCountDrift`.
- **Pinned in:** README By-the-Numbers row, website hero, `llms.txt`
  blockquote, demo hero (`componentCount` const, guarded against FEATURES.md).

## Icons — **102**

- **Definition:** entries of `icons/iconPathData` (path literals or
  `utils/svg` shared constants) + Spinner.
- **Computed by:** `icons.AllIconNames()` (live, in-process); the website and
  demo compute from it at render time.

## Typed enums — **64 (63 with IsValid())**

- **Definition:** exported closed-set string types; the IsValid count is
  functions matching `<Name>IsValid(` across library source.
- **Computed by:** `CountStats` (site, live); `utils.TestDocsCountDrift`
  guards the prose claims.

## Packages — **17**

- **Definition:** importable Go packages across the 7 published modules
  (module roots with Go files + first-level subpackages with Go files),
  excluding `internal/`, `cmd/`, `examples/`, and every repo-local module
  (own `go.mod`: `visualtest/`, `website/`, and the sub-module roots when
  scanned under the repo root — the same filter that prevents double-counting).
- **Computed by:** `utils.TestDocsCountDrift` (`countPublicPackages`) against
  the README row; the demo hero's `packageCount` chains to that row
  (`TestHeroCountsMatchFeatures`).

## Go modules — **7**

- **Definition:** the published modules in `publishedModules` (root, utils,
  icons, errorpage, charts/echarts, htmx, datastar). `visualtest/` and
  `website/` are repo-local and never counted or published.

## Tests — **~1,600 functions + ~1,650 subtests**

- **Definition:** `func Test/Fuzz/Benchmark` declarations and `t.Run(`
  occurrences in `*_test.go` files of the public packages (same scope as
  Packages), rounded to the nearest hundred — the `~` band is guard-enforced
  (multiple of 50 within 49 of live).
- **Computed by:** `utils.TestDocsCountDrift` (`countLibraryTests`).

## Goldens — **272 HTML / 200 pixel**

- **Definition:** committed baseline files under the golden testdata trees.
- **Computed by:** `utils.TestDocsCountDrift` (file counts per tree).

## Competitor facts (stars, claims) — dated, not constant

- **Definition:** external numbers about shadcn-templ/goshipit/etc. They are
  never guarded (they are not ours) — they are DATED and re-verified on the
  cadence in `docs/comparison.md` (quarterly + after known competitor events).
- **Lives in:** `docs/comparison.md` (canonical) and the website matrix
  (rendered from `data.go` + caption in `sections.templ`). Nowhere else may
  carry competitor claims; new surfaces must link, not restate.
