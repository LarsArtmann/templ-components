# Tailwind Scanning for a Vendored templ-components

## Problem

Your CSS build needs two things from the library, and Tailwind cannot find
either by itself:

- every utility class its components can emit (Tailwind generates CSS only for
  class names its scanner sees), and
- its hand-written CSS (`custom.css`, the semantic theme, a safelist), which
  only your own Tailwind run can compile.

The library's sources live in the Go module tree, outside Tailwind's automatic
source detection. Copying them by hand works until the first bump: a stale
copy raises nothing, the page just renders unstyled.

## The pattern: read the library from the vendor tree

Go already has a carrier for a dependency's files: `go mod vendor`. It ships a
dependency's package directories and the files those packages embed, and the
library puts its stylesheet in exactly that place (`layout/styles/`, embedded
by `layout/embed.go`). So the vendor tree holds everything, at the version
`go.mod` pins.

Run `go mod vendor` before Tailwind, then write two native lines:

```css
/* app.css: every path is relative to THIS file */
@import "tailwindcss" source(none);

/* Your own templates */
@source "../internal/**/*.templ";
@source "../internal/**/*.go";

/* The library: its sources (every module you import, one line) ... */
@source "../vendor/github.com/larsartmann/templ-components";
/* ... and its stylesheet: theme, component CSS, safelist */
@import "../vendor/github.com/larsartmann/templ-components/layout/styles/templ-components.css";

@custom-variant dark (&:where(.dark, .dark *));
```

Nothing is copied and nothing is generated. `go get -u` followed by
`go mod vendor` moves the markup and the CSS together.

A project that commits its `vendor/` directory can leave out the `@source`
line: automatic detection then finds the tree. Keep it explicit anyway when
`vendor/` is gitignored, because automatic detection skips gitignored paths.

## Verify it: `tc css check`

Gate the setup in CI (exit 1 on any finding):

```bash
GOFLAGS=-mod=mod go run github.com/larsartmann/templ-components/cmd/tc css check
```

It fails when the vendor tree carries another library version than `go.mod`
pins (a bump without a re-vendor), when no Tailwind entry scans the library,
and when none imports its stylesheet. `GOFLAGS=-mod=mod` matters in a vendored
project: `cmd/tc` is never part of `vendor/`.

## When there is no vendor tree where CSS is compiled

A CSS-only sandbox has no Go toolchain, no module cache and no vendor tree.
There the library's files must be plain files in the repository:

```bash
GOFLAGS=-mod=mod go run github.com/larsartmann/templ-components/cmd/tc css sync
```

```css
@source "./templ-components-classes.txt"; /* the library's classes */
@import "./templ-components.css";          /* its hand-written CSS */
```

`tc css sync` writes five tracked files at the pinned version and stamps each
with that version, so the next sync refreshes them without asking. `tc css
check` verifies this layout too.

## What `.gitignore` really does to a scan

Measured on Tailwind 4.3.3, 4.3.4 and 4.3.5 (2026-11-25), which corrects an
earlier claim in ADR-0045:

| Scan form                                                    | gitignored `vendor/` scanned |
| ------------------------------------------------------------ | ---------------------------- |
| automatic source detection                                   | no                           |
| `@source "."` (the ignored directory lies inside the walk)   | no                           |
| `@source "../vendor/github.com/larsartmann/templ-components` | yes                          |
| `@import "../vendor/.../layout/styles/templ-components.css"` | yes                          |

An explicit `@source` into a gitignored vendor tree is scanned. What made
vendoring unreliable before ADR-0050 was not `.gitignore`: the vendor tree
never carried the stylesheet, because `go mod vendor` ships package directories
and embedded files only.

## Related

- [`docs/tailwind-v4-adoption-guide.md`](../tailwind-v4-adoption-guide.md): the
  base setup, including "Deterministic scanning" for `source(none)`.
- [`docs/cli.md`](../cli.md): `tc init`, `tc css sync`, `tc css check`.
- [`docs/adr/0050-css-delivery-sync-and-verification.md`](../adr/0050-css-delivery-sync-and-verification.md):
  why the library is vendor-complete.
- [`docs/recipes/theme-bridge.md`](theme-bridge.md): remapping library colors
  onto a custom palette.
