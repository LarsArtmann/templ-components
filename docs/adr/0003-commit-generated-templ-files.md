# Report: Committing Generated `*_templ.go` Files for Module Proxy Compatibility

**Date:** 2026-05-18
**Status:** Resolved — all 31 files staged and ready for commit

## Problem

The Go module proxy (proxy.golang.org) serves source archives directly from Git tags. It does **not**
execute build steps like `templ generate`. When `*_templ.go` files are absent from the repository,
any consumer running `go get github.com/larsartmann/templ-components` receives uncompilable source
code — every exported component function is undefined.

This is the single most critical publishing blocker for a templ library. Unlike applications (which
generate at build/deploy time), a library's generated code **is** its distributable artifact.

## Root Cause

Two gitignore entries suppressed the files:

1. **Local `.gitignore`** — contained `*_templ.go` (removed in this change)
2. **Global `~/.config/git/ignore`** — also contains `*_templ.go` (overridden with `!*_templ.go` negation rule)

The local entry was straightforward to remove. The global entry required a negation rule in
`.gitignore` (`!*_templ.go`) to force Git to track these files regardless of global config.

## Changes Made

### `.gitignore`

- Removed the `*_templ.go` ignore entry
- Added `!*_templ.go` negation rule to override the global gitignore's `*_templ.go` pattern
- This ensures the files are tracked regardless of any developer's global Git configuration

### `*_templ.go` (31 files, 8 packages)

All generated files committed as-is. No manual edits — these are pure `templ generate` output:

| Package         | Files                                                                                       |
| --------------- | ------------------------------------------------------------------------------------------- |
| `display/`      | accordion, avatar, badge, card, dropdown, empty_state, helpers, modal, table, tabs, tooltip |
| `feedback/`     | alert, loading, progress, toast                                                             |
| `forms/`        | input, label, select, textarea                                                              |
| `htmx/`         | error_handling, helpers, loading                                                            |
| `icons/`        | icon                                                                                        |
| `internal/svg/` | svg                                                                                         |
| `layout/`       | base, theme                                                                                 |
| `navigation/`   | breadcrumbs, mobile_menu, nav_link, nav, pagination                                         |

### `AGENTS.md`

- Added a prominent **CRITICAL** section at the top of the Architecture section explaining:
  - Why generated files must be committed for this library
  - The `.gitignore` negation rule and its purpose
  - The workflow requirement to commit generated files alongside `.templ` source edits
  - A permanent prohibition against re-adding `*_templ.go` to `.gitignore`

## Verification

- `go build ./...` — passes (all 9 packages compile)
- `go test ./...` — passes (all 701 tests pass)
- `git check-ignore` — confirms files are no longer ignored
- `git add -n` — confirms all 31 files would be staged

## Going Forward

After any `.templ` file edit, developers must:

1. Run `templ generate ./...`
2. Commit the updated `*_templ.go` files **in the same commit** as the `.templ` source changes
3. Never rely on CI or consumers to run `templ generate`
