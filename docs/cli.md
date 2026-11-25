# `tc` CLI — templ-components consumer tool

`tc` is the library's consumer CLI. It scaffolds and verifies the stylesheet
inputs your Tailwind build needs (`tc init`, `tc css`), copies component
source files into your project (`tc add`), and diagnoses integration traps
(`tc doctor`, `tc lint`).

## Install

Run it at the version your `go.mod` pins (no install, works in vendored projects too):

```bash
GOFLAGS=-mod=mod go run github.com/larsartmann/templ-components/cmd/tc <command>
```

Or install a binary:

```bash
go install github.com/larsartmann/templ-components/cmd/tc@latest
```

An installed binary carries the files of ITS version. `tc css sync` refuses
to write them into a project that pins another version and names the command
that matches.

## Commands

### `tc init`

Scaffolds the Tailwind entry `app.css` plus the four library-owned stylesheet
files (`templ-components-classes.txt`, `custom.css`,
`templ-components-theme.css`, `safelist.css`) in the current directory.
Skips files that already exist.

```bash
mkdir myapp && cd myapp
tc init
# tc: wrote app.css
# tc: wrote custom.css
# ...
```

`app.css` is yours: point its two `@source` lines at your own templates. The
other four are library-owned (see `tc css`).

### `tc ls`

Prints every component, grouped by package. Use it to discover what's
available before copying.

```bash
tc ls
# display
#   accordion
#   avatar
#   ...
```

### `tc add <component>`

Copies a component's `.templ` file (and its `_types.go` sibling, if present)
to `./components/` by default. Override the destination with `--out`.

```bash
tc add button
# tc: wrote components/button.templ

tc add dropdown --out ./src/components
# tc: wrote src/components/dropdown.templ
# tc: wrote src/components/dropdown_templ.go   (if present)
```

Only the source is copied: run `templ generate` in your project to produce the
`_templ.go` file, then import your local copy instead of the library package.
A component that relies on sibling helpers (class lookups, enums, shared
sub-templates) needs those too; `tc add <component> --list-deps` names them.

### `tc doctor`

Diagnoses the top consumer-integration traps in one command — every check
prints pass/fail with a fix hint, and the exit code is non-zero on failure
(so CI can gate on it):

- **Tailwind `@source` scans `.templ` files** — the #1 integration trap:
  the `@source` path resolves relative to the CSS file, not the CWD, so a
  wrong path silently compiles a stylesheet with zero component classes.
- **`GOEXPERIMENT=jsonv2`** — only needed on toolchains older than the Go
  1.27 floor (json/v2 is stable there; nothing to set on 1.27+).
- **templ version pin** — a generator older than the one the library
  generates with (v0.3.1100) emits incompatible `*_templ.go` runtime calls.
- **Committed `*_templ.go`** (library repo only) — a `.templ` source
  without its generated file ships uncompilable code through the proxy.
- **`core.hooksPath` adoption** (library repo only) — the tracked
  `.githooks/` guards are inactive until `scripts/setup-hooks.sh` runs.

```bash
tc doctor
# tc doctor — templ-components environment check
#   [pass] GOEXPERIMENT=jsonv2 (env)
#   ...
# result: all checks passed.
```

### `tc lint`

Scans your source files for the same compliance rules the library enforces on
itself with its own guard tests. Each scanner is the exact engine behind a
library test, so adopting `tc lint` in CI pins your markup to the library's
accessibility and RTL contract.

```bash
tc lint [--rtl] [--dark] [--motion] [--icons] [DIR ...]
```

Without `DIR`, the current directory is scanned recursively (icons: direct
children only). Select at least one scanner:

| Scanner    | Rule(s)                                 | Catches                                                                                            |
| ---------- | --------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `--rtl`    | `rtl-physical`, `rtl-inset-alias`       | Physical direction classes (`ml-`, `pr-`, `text-left`, …) and deprecated `start-*`/`end-*` aliases |
| `--dark`   | `darkmode-neutral`, `darkmode-semantic` | Color classes without a same-line `dark:` variant                                                  |
| `--motion` | `motion-reduce`                         | `transition-*`/`animate-*` without a `motion-reduce:` fallback                                     |
| `--icons`  | `icon-path-inline`                      | Inlined SVG `d="M..."` literals outside the canonical tables                                       |

Findings print as `file:line: rule: message` with the offending line. Exit
codes: `0` clean, `1` findings or scan failure, `2` usage error — so a CI job
or pre-commit hook can gate on it:

```bash
tc lint --rtl --motion ./components
# components/banner.templ:12: rtl-physical: RTL physical-property violation ...
# tc lint: 1 finding(s) (--rtl --motion) in 1 directory
```

#### Suppressing findings: `tc-lint:ignore`

Deliberate exceptions carry a per-line marker, in any comment syntax the file
language offers. It must sit on the SAME line as the finding:

```html
<div class="ml-2"><!-- tc-lint:ignore --></div>
<span class="ml-2"><!-- tc-lint:ignore rtl-physical --></span>
```

```go
var arrow = "translate-x-1" // tc-lint:ignore motion-reduce
```

A bare `tc-lint:ignore` suppresses every rule on the line; listing rule ids
suppresses exactly those. Prefer the narrow form — a bare marker hides future
rules too.

### `tc css sync`

Writes the stylesheet inputs of the library version `go.mod` pins into the
directory of your Tailwind entry (or `--dir DIR`), then prints the lines that
entry still needs:

```bash
GOFLAGS=-mod=mod go run github.com/larsartmann/templ-components/cmd/tc css sync
# tc: wrote css/templ-components-classes.txt
# tc: wrote css/custom.css
# tc: wrote css/templ-components-theme.css
# tc: wrote css/safelist.css
# tc css: css/app.css does not reach every library file yet; add these lines to it:
#   @source "./templ-components-classes.txt";
#   @import "./custom.css";
#   @import "./templ-components-theme.css";
#   @import "./safelist.css";
```

`go get` updates the Go code only. Re-run the sync after every library bump.
Each written file carries a first-line stamp naming the version it came from,
so the next sync refreshes it without asking. A differing copy WITHOUT a
stamp may hold your own rules and is kept: move those rules to a stylesheet of
your own, then re-run with `--force`.

### `tc css check`

Fails (exit 1) when the stylesheet inputs are not right, so CI can gate on it:

- `go.mod` pins a different library version than this `tc`.
- A library file is missing or differs from the pinned version.
- No Tailwind entry (a CSS file importing `tailwindcss`) scans the class corpus
  or imports the component CSS.

```bash
GOFLAGS=-mod=mod go run github.com/larsartmann/templ-components/cmd/tc css check
#   [pass] go.mod pin equals this tc (v1.27.2)
#   [FAIL] custom.css (css/custom.css) differs from v1.27.2 (copy is stamped v1.27.0)
#          fix: run `tc css sync`
# result: 1 check(s) FAILED; see the fix hints above
```

`--dir DIR` names where the library files live, `--css PATH` one Tailwind
entry to inspect instead of every discovered one. A stale copy raises nothing
on its own: the page just renders unstyled. That is what this check is for.

## Why use this?

The library is designed to cover the common cases with sensible defaults.
When you need full control — different DOM structure, different accessibility
patterns, integration with a design system the library doesn't support — copy
the component into your project and make it yours. The `tc` tool keeps the
copied source in sync with the library's version at install time.

## Why NOT use this?

If the library's `BaseProps.Class` and `BaseProps.Attrs` are enough to
customize the component, you don't need `tc add`. The library was designed so
that ~95% of customization happens via class override; `tc add` is for the
remaining ~5%.

## Updating

Re-run `tc add <component>` to overwrite your local copy with the latest
library version. **This destroys your local changes** — only do it if you
haven't customized the file, or if you're willing to re-apply your changes.

For a more controlled update, copy the new file to a temp location and diff:

```bash
tc add button --out /tmp/tc-new
diff components/button.templ /tmp/tc-new/button.templ
```
