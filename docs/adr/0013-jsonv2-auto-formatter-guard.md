# ADR 0013: encoding/json/v2 auto-formatter guard

## Date

2026-07-12

## Status

Accepted

## Context

This library uses `encoding/json/v2` + `encoding/json/jsontext` (Go 1.26+ with
`GOEXPERIMENT=jsonv2`). The `errorpage` package uses `json.MarshalEncode` +
`jsontext.NewEncoder` for JSON error responses.

The import was accidentally introduced **three times** by auto-formatters
running under `GOEXPERIMENT=jsonv2`. When `gofmt` or `golangci-lint` runs with
the experiment flag enabled, it rewrites `encoding/json` imports to
`encoding/json/v2` automatically — even in files that intentionally use v1
(e.g., breadcrumbs JSON-LD marshaling, test helpers).

This creates a cycle:

1. Developer runs formatter with `GOEXPERIMENT=jsonv2` set
2. Formatter rewrites `encoding/json` to `encoding/json/v2` in all files
3. Commit includes unintended v2 migration
4. CI without the experiment flag fails (v2 not available without the flag)
5. Developer reverts, but the next formatter run reintroduces it

## Decision

### 1. Pre-commit hook sets the experiment flag consistently

The pre-commit hook (`scripts/pre-commit.sh`) already sets
`GOEXPERIMENT=jsonv2` before running formatters. This ensures the formatter
produces consistent output regardless of the developer's shell environment.

### 2. Files using v1 intentionally are documented

The following files use `encoding/json` v1 by design (they don't need v2
features and should remain on v1 for maximum compatibility):

- `navigation/breadcrumbs.templ` — `json.Marshal` for JSON-LD structured data
- Test files that marshal/unmarshal simple structs

### 3. Only `errorpage` uses v2

The `errorpage` package is the only package that imports
`encoding/json/v2` + `encoding/json/jsontext`. It uses `json.MarshalEncode`
for streaming JSON error responses and `jsontext.NewEncoder` for fine-grained
control over the output format.

### 4. `.golangci.yml` enables the experiment flag

The linter config enables `goexperiment.jsonv2` build tag so CI is consistent
with the pre-commit hook.

## Consequences

- Developers must run `nix develop` (or set `GOEXPERIMENT=jsonv2` manually)
  before running formatters, to match CI.
- When Go 1.27 ships with json/v2 as stable, the experiment flag becomes
  unnecessary and this ADR can be retired.
- The breadcrumbs package stays on v1 until v2 is the default — it marshals
  a simple struct and gains nothing from v2.

## Prevention checklist

When adding a new `encoding/json` import:

1. Use v1 (`encoding/json`) unless you need v2 features (streaming encoder,
   custom marshalers with `jsontext.Value`).
2. If the formatter rewrites your import to v2, check whether your file
   actually needs it. If not, revert to v1.
3. The pre-commit hook will run the formatter with the correct flag —
   trust it over your shell's environment.
