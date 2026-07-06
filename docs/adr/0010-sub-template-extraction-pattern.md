# ADR 0010: Sub-template extraction pattern

## Date

2026-07-06

## Status

Accepted

## Context

The templ-components library uses templ's sub-template capability to share markup
between components. As the library grew to 80+ components, duplicated markup
patterns emerged across packages (e.g., the flex-row icon+content layout in
`ErrorPage` and `ErrorDetail`, the `history.back()` script in `ErrorPage` and
`NotFound404`).

A `art-dupl --semantic` deduplication sprint identified 12 production clone groups
at threshold t≥7. Six were extracted into shared sub-templates. Six remain as
accepted clones (documented in [ADR 0009](0009-accepted-clones.md)).

This ADR formalizes **when** to extract a sub-template and **when** to leave
duplication in place.

## Decision

### Extract a sub-template when ALL three conditions are met:

1. **Two or more production callers** exist in different files.
2. **Five or more lines** of shared markup (single-line patterns aren't worth
   the indirection).
3. **A clear domain name** describes the shared concept (e.g., `errorHeader`,
   `skeletonContainer`, `definitionDetailContent`).

### Do NOT extract when:

- **Only one caller exists** (YAGNI — wait for the second caller).
- **The clone is in demo/example code** (showcase content, not production).
- **The shared concept has no clean domain name** (if you can't name it, it
  isn't a real abstraction).
- **The extraction would create a parameter explosion** (e.g., 10+ parameters
  to parameterize the variation).
- **The callers are in the same file** (templ already supports local
  templates — no cross-file extraction needed).

### Naming convention

- Private sub-templates use `camelCase`: `errorHeader`, `goBackScript`.
- Parameter count ≤ 7 is ideal. At 8+, consider a props struct.
- Document the trigger condition for promotion/demotion in a comment.

### Promotion trigger

When a private sub-template in one package is needed by a second package,
promote it to the shared layer (`utils` for Go helpers, `display/shared.templ`
for templ templates). The `DismissButton` promotion (feedback → utils) is
the established precedent.

## Consequences

- **Positive:** Reduces duplication, enforces consistent markup patterns,
  makes future changes atomic (edit one sub-template, all callers update).
- **Negative:** Adds indirection — reading a component requires following
  the sub-template call chain. Mitigated by clear names and godoc.
- **Neutral:** The `art-dupl` clone count is a measurable metric for when
  extraction is needed (threshold t≥7 for production code).
