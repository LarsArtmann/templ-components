# Tailwind Scanning for a Vendored templ-components

## Problem

You vendor `templ-components` (`go mod vendor`) and point Tailwind at the vendor
tree:

```css
@source "../vendor/github.com/larsartmann/templ-components";
```

The utility classes arrive, but something is always off: components render
unstyled in CI, on a fresh clone, or in the Nix sandbox, and every rule that
starts with `.tc-` is missing everywhere.

**Root cause:** the vendor tree is not a delivery channel for the stylesheet
inputs.

- It exists only after `go mod vendor` ran, so the scan depends on the
  environment: a fresh clone has no tree, a Nix sandbox has a different one.
- Go vendors package directories only. `templates/` is not a package, so
  `custom.css`, the theme and the safelist never arrive through `vendor/`,
  whatever the `@source` line says.

## The pattern: let `tc css` own the library files

The library ships every stylesheet input in the module and a command that
writes them at the version `go.mod` pins (ADR-0045, ADR-0050). No per-project
generator, no vendor path:

```bash
GOFLAGS=-mod=mod go run github.com/larsartmann/templ-components/cmd/tc css sync
```

`tc css sync` writes four tracked files next to your Tailwind entry and prints
the lines the entry still needs:

| File                           | What it is                                                                     |
| ------------------------------ | ------------------------------------------------------------------------------ |
| `templ-components-classes.txt` | Class corpus: every utility class the library can emit, for Tailwind's scanner |
| `custom.css`                   | Component CSS (dialog animations, stylable select, ...)                        |
| `templ-components-theme.css`   | Semantic color aliases (`tc-primary`, `tc-danger`, ...)                        |
| `safelist.css`                 | `GridColsAutoFit` classes the scanner cannot see                               |

```css
@import "tailwindcss" source(none);

/* Your own views (paths resolve relative to THIS file) */
@source "../internal/**/*.templ";
@source "../internal/**/*.go";

/* The library's class corpus */
@source "./templ-components-classes.txt";

/* The library's hand-written CSS */
@import "./custom.css";
@import "./templ-components-theme.css";
@import "./safelist.css";

@custom-variant dark (&:where(.dark, .dark *));
```

`GOFLAGS=-mod=mod` matters for a vendored project: `cmd/tc` is never part of
`vendor/`, so the default vendor mode cannot build it.

## Verify it: `tc css check`

Gate the setup in CI (exit 1 on any finding):

```bash
GOFLAGS=-mod=mod go run github.com/larsartmann/templ-components/cmd/tc css check
```

It fails when `go.mod` pins another library version than the files, when a
file is missing or drifted, or when no Tailwind entry reaches the corpus or
the component CSS. A stale copy raises nothing on its own: the page just
renders unstyled.

Run `tc css sync` after every library bump. The copies carry a first-line
stamp naming the version they came from, so the sync refreshes them without
asking; a differing copy without a stamp is treated as yours and kept.

## Why a tracked file and not the vendor tree or the module cache?

| Approach                         | Fresh clone      | Nix sandbox   | Hermetic | Carries `custom.css` |
| -------------------------------- | ---------------- | ------------- | -------- | -------------------- |
| `@source "vendor/..."`           | no (tree absent) | depends       | no       | never                |
| `@source` on module cache        | yes              | no (no cache) | no       | no                   |
| Tracked files from `tc css sync` | yes              | yes           | yes      | yes                  |

A tracked plain file is a build input like any other: versioned, hashable, and
identical everywhere.

## What `.gitignore` really does to a scan

Measured on Tailwind 4.3.3, 4.3.4 and 4.3.5 (2026-11-25), which corrects an
earlier claim in ADR-0045:

- Automatic source detection skips gitignored paths.
- A directory source skips the gitignored directories it _contains_
  (`@source "."` with a gitignored `vendor/` inside never sees it).
- An `@source` that _names_ the gitignored directory (or a glob into it) is
  scanned, and `@import` of a gitignored file resolves.

So an explicit vendor `@source` is not silently inert. It fails for the two
reasons above instead: the tree is absent where `go mod vendor` did not run,
and it never carries the hand-written CSS.

## Related

- [`docs/tailwind-v4-adoption-guide.md`](../tailwind-v4-adoption-guide.md): the
  base setup, including "Deterministic scanning" for `source(none)`.
- [`docs/cli.md`](../cli.md): `tc init`, `tc css sync`, `tc css check`.
- [`docs/recipes/theme-bridge.md`](theme-bridge.md): remapping library colors
  onto a custom palette once the classes are being scanned correctly.
