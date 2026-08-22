# Deterministic Tailwind Scanning for a Vendored templ-components

## Problem

You vendor `templ-components` (`go mod vendor`) and point Tailwind at the vendor
tree:

```css
@source "../vendor/github.com/larsartmann/templ-components";
```

It works on your machine — then components render unstyled in CI, in the Nix
sandbox, or on a fresh clone.

**Root cause:** `vendor/` is conventionally gitignored, and Tailwind v4 skips
gitignored paths — including directories listed explicitly via `@source`.
Whether the scan works depends on whether a `.gitignore` exists in the build
environment. A Nix sandbox has none (scan works); a fresh clone has one
(scan silently drops every library class). The output is
environment-dependent, which is the worst kind of build failure: it passes
everywhere you test and breaks where it ships.

## The pattern: scan a tracked class-inventory file

Concatenate the library's source files into one **committed** plain-text file
and scan that instead. The scan set becomes byte-identical in every
environment because the file is versioned like any other source.

### 1. Generate the inventory

```bash
#!/usr/bin/env bash
# scripts/gen-library-classes.sh — regenerate the vendored class inventory.
set -euo pipefail
cd "$(dirname "$0")/.."

dest=internal/server/views/library-classes.txt

# LC_ALL=C pins the sort order so the concatenation is byte-identical across
# environments (the Nix sandbox sorts under the C locale).
LC_ALL=C find vendor/github.com/larsartmann/templ-components \
	-type f \( -name '*.templ' -o -name '*.go' \) ! -name '*_templ.go' \
	-print0 | LC_ALL=C sort -z | xargs -0 cat >"$dest"

echo "wrote $dest ($(wc -c <"$dest") bytes)"
```

Notes:

- `! -name '*_templ.go'` — generated files mirror the `.templ` sources and leak
  utility-looking identifiers that inflate the output with dead CSS.
- `LC_ALL=C` + `sort -z` — deterministic byte order regardless of locale.

### 2. Scan it from your CSS

```css
/* Disable automatic detection: everything scanned is declared explicitly. */
@import "tailwindcss" source(none);

/* Your own views */
@source ".";

/* Never scan generated templ output or previous CSS builds */
@source not "**/*_templ.go";
@source not "app.min.css";
@source not "styles.css";

/* The vendored library's class inventory (committed, deterministic) */
@source "./library-classes.txt";

@custom-variant dark (&:where(.dark, .dark *));
```

### 3. Wire regeneration into your flow

Run the script:

- after `go mod vendor` (a `check-vendor-hash.sh` style guard works well), and
- after bumping the `templ-components` version.

Then recompile the CSS artifacts and commit all regenerated files together —
a CI check that rebuilds the CSS and diffs it catches a forgotten regeneration.

## Why a `.txt` file and not the vendor tree?

| Approach                     | Fresh clone | Nix sandbox | Docker | Hermetic |
| ---------------------------- | ----------- | ----------- | ------ | -------- |
| `@source "vendor/..."`       | ❌ skipped (gitignored) | ✅ | depends | no |
| `@source` on module cache    | ✅          | ❌ (no cache) | ❌     | no       |
| Tracked `library-classes.txt` | ✅         | ✅          | ✅     | **yes**  |

The inventory file is a build input like any other: versioned, hashable, and
identical everywhere. Proven in production by
[dnsblockd](https://github.com/larsartmann/dnsblockd)'s Nix build, which
byte-verifies the compiled CSS artifacts.

## Recipe: drift-guard test

If your project has Go tests, pin the inventory so it cannot silently rot:

```go
func TestLibraryClassesFresh(t *testing.T) {
	// Recompute the expected inventory in-process and compare to the
	// committed file — fails when someone bumps the vendored version
	// without regenerating library-classes.txt.
}
```

(Or simpler: a CI step that runs the script and asserts `git diff --exit-code`.)

## Related

- [`docs/tailwind-v4-adoption-guide.md`](../tailwind-v4-adoption-guide.md) — the
  base setup; see the "Deterministic scanning" section for `source(none)`.
- [`docs/recipes/theme-bridge.md`](theme-bridge.md) — remapping library colors
  onto a custom palette once the classes are being scanned correctly.
